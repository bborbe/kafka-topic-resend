# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards compatible manner, and
* PATCH version when you make backwards compatible bug fixes.

## v0.4.0

- feat: `make build` resolves `BUILD_TARGET` from a bounded sun probe - sun when reachable, local docker as the fallback. An explicit `BUILD_TARGET` always wins, and the decision prints at parse time so it is visible under `make -n`

## v0.3.2

- chore: restate the do-not-merge contract in the gate workflow

## v0.3.1

- fix: keep the rootless BuildKit build container startable. `allowPrivilegeEscalation: false` and `capabilities.drop: ["ALL"]` on that container stop `rootlesskit` from installing its uid/gid map (`newuidmap: operation not permitted`), so every cluster build died at startup; both are removed and the reason is recorded in `hack/build-job.yaml`. Also hardened the clone path: the token is scoped to its origin host, deleted on the error path as well as the success path, and refused over a non-https URL; `REF` can no longer be read as a git option; and `__CONTEXT__` is listed in the template's token table

## v0.3.0

- feat: build the image on a Kubernetes cluster with `make build BUILD_TARGET=<cluster> REF=<sha> VERSION=<tag>`, keeping the default local arm unchanged. `hack/build-cluster.sh`, `hack/cluster-alias.sh` and `hack/build-job.yaml` are vendored from `bborbe/agent-task-executor`; `check-version-tag` and `check-go-mod` stay local-arm prerequisites, since they describe the local HEAD rather than the commit the cluster builds
- chore: bump `golang.org/x/net` to v0.62.0 to clear `make vulncheck` — GO-2026-6860 and GO-2026-6866 are newly surfaced against v0.61.0 and present on master too, so this is a pre-existing advisory rather than a change in this PR

## v0.2.7

- chore: update Go to 1.27.3 and github.com/bborbe/kafka to v1.28.4, github.com/bborbe/sentry to v1.10.5, github.com/bborbe/service to v1.11.0, github.com/bborbe/time to v1.27.17

## v0.2.6

- chore: update Go to 1.27.2 and github.com/IBM/sarama to v1.61.1, github.com/bborbe/boltkv to v1.16.2, github.com/bborbe/errors to v1.6.3, github.com/bborbe/kafka to v1.27.1, github.com/bborbe/kv to v1.21.18, github.com/bborbe/log to v1.7.4, github.com/bborbe/metrics to v0.6.7, github.com/bborbe/run to v1.11.2, github.com/bborbe/sentry to v1.10.4, github.com/bborbe/service to v1.10.16, github.com/bborbe/time to v1.27.16, github.com/onsi/ginkgo/v2 to v2.33.1, github.com/onsi/gomega to v1.44.0

## v0.2.5

- fix: bump `osv-scanner` to v2.6.0 and `golang.org/x/net` to v0.60.0 so the Linux vulnerability gates stop failing. v2.3.1 pins `golang.org/x/tools` v0.38.0, whose SSA builder aborts with `unexpected expr: *ast.KeyValueExpr` on the promoted-field composite-literal key Go 1.27 permits in the Linux stdlib, so a repo on the old pin passes locally on darwin and fails only in Linux CI. `x/net` v0.58.0 carries `GO-2026-6603/6610/6611/6612/6617`, which fail both `vulncheck` and `trivy`.

## v0.2.4

- chore: update Go to 1.27.1 and github.com/bborbe/boltkv to v1.15.5, github.com/bborbe/errors to v1.6.1, github.com/bborbe/kafka to v1.25.16, github.com/bborbe/kv to v1.21.15, github.com/bborbe/metrics to v0.6.3, github.com/bborbe/run to v1.10.3, github.com/bborbe/sentry to v1.10.2, github.com/bborbe/service to v1.10.13, github.com/bborbe/time to v1.27.14

## v0.2.3

- fix: `make build` refuses to stamp a version onto a tree that is not that version's tag (`check-version-tag`, escape hatch `ALLOW_UNTAGGED_BUILD=1`). `VERSION` defaults to the newest tag repo-wide, so an operator-run build from an untagged or older tree silently republishes under the newest tag. The guard compares `git describe --exact-match HEAD` against `$(VERSION)` and exits non-zero on mismatch.

## v0.2.2

- ci: add ci.yml running `make precommit` -- the `test` required status check had no workflow producing it, permanently blocking every PR
- chore: update golang.org/x/crypto to v0.56.0 (fixes GO-2026-6354, GO-2026-6355)

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
