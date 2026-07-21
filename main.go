// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/IBM/sarama"
	libboltkv "github.com/bborbe/boltkv"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/bborbe/metrics"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
	bolt "go.etcd.io/bbolt"
)

const serviceName = "kafka-topic-resend"

var bucketName = libkv.NewBucketName("data")

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	SentryDSN       string             `required:"true"  arg:"sentry-dsn"        env:"SENTRY_DSN"        usage:"SentryDSN"                             display:"length"`
	SentryProxy     string             `required:"false" arg:"sentry-proxy"      env:"SENTRY_PROXY"      usage:"Sentry Proxy"`
	KafkaBrokers    libkafka.Brokers   `required:"true"  arg:"kafka-brokers"     env:"KAFKA_BROKERS"     usage:"Comma separated list of Kafka brokers"`
	BatchSize       libkafka.BatchSize `required:"true"  arg:"batch-size"        env:"BATCH_SIZE"        usage:"batch consume size"                                     default:"1"`
	Topic           string             `required:"true"  arg:"topic"             env:"TOPIC"             usage:"topic to resend"`
	NoSync          bool               `required:"true"  arg:"no-sync"           env:"NO_SYNC"           usage:"no sync"                                                default:"false"`
	BuildGitVersion string             `required:"false" arg:"build-git-version" env:"BUILD_GIT_VERSION" usage:"Build Git version"                                      default:"dev"`
	BuildGitCommit  string             `required:"false" arg:"build-git-commit"  env:"BUILD_GIT_COMMIT"  usage:"Build Git commit hash"                                  default:"none"`
	BuildDate       *libtime.DateTime  `required:"false" arg:"build-date"        env:"BUILD_DATE"        usage:"Build timestamp (RFC3339)"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	metrics.NewBuildInfoMetrics().SetBuildInfo(a.BuildGitVersion, a.BuildGitCommit, a.BuildDate)

	saramaClientProvider, err := libkafka.NewSaramaClientProviderByType(
		ctx,
		libkafka.SaramaClientProviderTypeReused,
		a.KafkaBrokers,
	)
	if err != nil {
		return errors.Wrapf(ctx, err, "create sarama client provider failed")
	}
	defer saramaClientProvider.Close()

	rawSyncProducer, err := libkafka.NewSyncProducerWithName(
		ctx,
		a.KafkaBrokers,
		serviceName,
	)
	if err != nil {
		return errors.Wrapf(ctx, err, "create sync producer failed")
	}
	syncProducer := libkafka.NewSyncProducerMetrics(rawSyncProducer)
	defer syncProducer.Close()

	db, err := libboltkv.OpenTemp(ctx, func(opts *bolt.Options) {
		opts.NoSync = a.NoSync
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "open db failed")
	}
	defer func() {
		_ = db.Close()
		_ = db.Remove()
	}()

	topicConsumed := run.NewTrigger()
	topicSend := run.NewTrigger()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-topicSend.Done():
			cancel()
		}
	}()

	return service.Run(
		ctx,
		a.consumeTopic(saramaClientProvider, db, topicConsumed),
		run.Triggered(a.send(db, syncProducer, topicSend), topicConsumed.Done()),
	)
}

func (a *application) consumeTopic(
	saramaClientProvider libkafka.SaramaClientProvider,
	db libkv.DB,
	trigger run.Fire,
) run.Func {
	return func(ctx context.Context) error {
		logSampler := log.DefaultSamplerFactory.Sampler()
		return libkafka.NewOffsetConsumerHighwaterMarksBatchWithProvider(
			saramaClientProvider,
			libkafka.Topic(a.Topic),
			libkafka.NewStoreOffsetManager(
				libkafka.NewOffsetStore(db),
				libkafka.OffsetOldest,
				libkafka.OffsetNewest,
			),
			libkafka.NewMessageHandlerBatchTxUpdate(
				db,
				libkafka.NewMessageHandlerBatchTx(
					libkafka.NewMessageHandlerTxSkipErrors(
						libkafka.NewMessageHandlerTxMetrics(
							libkafka.MessageHandlerTxFunc(
								func(ctx context.Context, tx libkv.Tx, msg *sarama.ConsumerMessage) error {
									bucket, err := tx.CreateBucketIfNotExists(ctx, bucketName)
									if err != nil {
										return errors.Wrapf(ctx, err, "get bucket failed")
									}
									if err := bucket.Put(ctx, msg.Key, msg.Value); err != nil {
										return errors.Wrapf(ctx, err, "put failed")
									}
									if logSampler.IsSample() {
										glog.V(2).
											Infof("consumer message %d completed (sample)", msg.Offset)
									}
									return nil
								},
							),
							libkafka.NewMetrics(),
						),
						log.DefaultSamplerFactory,
					),
				),
			),
			a.BatchSize,
			trigger,
			log.DefaultSamplerFactory,
		).Consume(ctx)
	}
}

func (a *application) send(
	db libkv.DB,
	syncProducer libkafka.SyncProducer,
	trigger run.Fire,
) run.Func {
	return func(ctx context.Context) error {
		logSampler := log.DefaultSamplerFactory.Sampler()

		return db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
			bucket, err := tx.Bucket(ctx, bucketName)
			if err != nil {
				return errors.Wrapf(ctx, err, "get bucket failed")
			}
			it := bucket.Iterator()
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				err := item.Value(func(v []byte) error {
					partition, offset, err := syncProducer.SendMessage(
						ctx,
						&sarama.ProducerMessage{
							Topic: a.Topic,
							Key:   sarama.ByteEncoder(item.Key()),
							Value: sarama.ByteEncoder(v),
						},
					)
					if err != nil {
						return errors.Wrapf(ctx, err, "send failed")
					}
					if logSampler.IsSample() {
						glog.V(2).
							Infof("send update message successful to %s with partition %d offset %d (sample)", a.Topic, partition, offset)
					}
					return nil
				})
				if err != nil {
					return errors.Wrapf(ctx, err, "handle value failed")
				}
			}
			trigger.Fire()
			glog.V(2).Infof("send completed")
			return nil
		})
	}
}
