# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards compatible manner, and
* PATCH version when you make backwards compatible bug fixes.

## Unreleased

- chore: update github.com/bborbe/boltkv to v1.15.0, github.com/bborbe/errors to v1.6.0, github.com/bborbe/kafka to v1.25.10, github.com/bborbe/kv to v1.21.12, github.com/bborbe/metrics to v0.6.0, github.com/bborbe/run to v1.10.1, github.com/bborbe/sentry to v1.10.0, github.com/bborbe/service to v1.10.10, github.com/bborbe/time to v1.27.11, github.com/onsi/gomega to v1.43.0

## v0.2.1

- chore: update github.com/IBM/sarama to v1.60.2, github.com/bborbe/boltkv to v1.14.10, github.com/bborbe/errors to v1.5.21, github.com/bborbe/log to v1.6.25, github.com/bborbe/metrics to v0.5.15

## v0.2.0

- feat: opt into `autoMerge.trivial` for mechanically-trivial update PRs

## v0.1.6

- chore: update go module dependencies

## v0.1.5

- chore: update Go to 1.27.0 dependencies (github.com/bborbe/math to v1.4.4, github.com/bborbe/sentry to v1.9.27)

## v0.1.4

- chore: update Go to 1.27.0 and github.com/bborbe/boltkv to v1.14.9, github.com/bborbe/errors to v1.5.20, github.com/bborbe/kafka to v1.25.9, github.com/bborbe/kv to v1.21.11, github.com/bborbe/log to v1.6.24, github.com/bborbe/metrics to v0.5.14, github.com/bborbe/run to v1.9.37, github.com/bborbe/sentry to v1.9.26, github.com/bborbe/service to v1.10.9, github.com/bborbe/time to v1.27.10

## v0.1.3

- chore: Bump errcheck to v1.20.0 and golangci-lint to v2.13.1 for Go 1.27 support
## v0.1.2

- update Go to 1.26.6 and update dependencies, fix GO-2026-6179, GO-2026-6180, CVE-2026-56864, CVE-2026-56865

## v0.1.1

- docs: add a License section to the README

## v0.1.0

- Initial release — extracted from `bborbe/trading` (`strimzi/topic-resend`) as a standalone public repo
- Consumes a full Kafka topic to a local BoltDB store, then re-sends all messages; broker-agnostic
- Decoupled from `trading/lib`: build-info via public `github.com/bborbe/metrics`, sync-producer via public `github.com/bborbe/kafka`
- Publish-only build → `docker.io/bborbe/kafka-topic-resend:vX.Y.Z`
