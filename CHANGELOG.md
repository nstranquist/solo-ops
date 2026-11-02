# Changelog — solo-ops

Offline inquiry CRM plus policy pack and a local read-only HTML cockpit.

## [0.2.0] — 2026-11-02

### Changed

- The read-only cockpit is now a summary-and-card layout (status counts,
  pills, dates) instead of a bare HTML table. Public screenshots use the
  fictional `testdata/showcase.jsonl` ledger.

## [0.1.0] — 2026-07-25

First public release.

### Added

- Release boundary: `CHANGELOG.md`, `SECURITY.md`, `CONTRIBUTING.md`,
  `DISTRIBUTION.md`, `Makefile`, and CI.
- `make smoke` — runs every command documented in the README against the
  committed fixtures, so documented commands cannot drift from the CLI.
- `make leak-scan` + `make leak-scan-selftest` — generic leak classes, with a
  self-test that plants a synthetic leak and asserts the scan rejects it.

### Notes

The CLI predates this release; it was built on 2026-07-21 and lived as a local
repository with no remote until now. Publishing it is the change — the code is
not new. Versioned `0.1.0` because no version was ever released.
