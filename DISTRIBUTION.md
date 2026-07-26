# Distribution readiness — solo-ops

**Version:** v0.1.0
**License:** MIT
**Module:** `github.com/nstranquist/solo-ops`

## Install

```sh
go install github.com/nstranquist/solo-ops/cmd/solo-ops@latest
```

Requires Go 1.26+.

Local build:

```sh
make install   # → ~/.local/bin/solo-ops
```

## Verify before publishing

```sh
make publish-ready
```

Runs `go vet`, the unit suite, `make smoke` over every README-documented
command, and the leak scan plus its self-test.

## Release scan

CI enforces the generic classes via `make leak-scan`. A private,
project-specific banned-fragment list is maintained outside this repository and
applied before publishing. **That list is deliberately not committed here** — a
scan that embeds the strings it forbids publishes them.

## Boundary

The local-first CLI is MIT and lives here. Any hosted or multi-tenant service
built on it stays proprietary and out of this repository.
