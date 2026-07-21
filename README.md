# kafka-topic-resend

Consumes a complete Kafka topic into a local BoltDB store, then re-sends every message back to the same topic. Useful for triggering re-processing by downstream consumers. Broker-agnostic.

## Run

```
make run
```

## Flags

| Flag | Env | Description |
|---|---|---|
| `-kafka-brokers` | `KAFKA_BROKERS` | Kafka brokers (comma-separated) |
| `-topic` | `TOPIC` | topic to resend |
| `-batch-size` | `BATCH_SIZE` | batch consume size (default 1) |
| `-no-sync` | `NO_SYNC` | disable BoltDB fsync (faster, less durable) |
| `-sentry-dsn` | `SENTRY_DSN` | optional Sentry DSN |

## Build

`make buca` builds and publishes `docker.io/bborbe/kafka-topic-resend:vX.Y.Z` (git-tag semver).
