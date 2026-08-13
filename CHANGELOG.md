# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards compatible manner, and
* PATCH version when you make backwards compatible bug fixes.

## Unreleased

- docs: add a License section to the README

## v0.1.0

- Initial release — extracted from `bborbe/trading` (`strimzi/topic-resend`) as a standalone public repo
- Consumes a full Kafka topic to a local BoltDB store, then re-sends all messages; broker-agnostic
- Decoupled from `trading/lib`: build-info via public `github.com/bborbe/metrics`, sync-producer via public `github.com/bborbe/kafka`
- Publish-only build → `docker.io/bborbe/kafka-topic-resend:vX.Y.Z`
