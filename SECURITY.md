# Security — solo-ops

## Reporting

Report suspected vulnerabilities privately via a GitHub security advisory on
this repository. Please do not open a public issue for an unpatched issue.

## Threat model

solo-ops is a local-first CLI. It reads and writes files on your machine, makes no
network requests, and has no telemetry, update check, or remote index.

**The inquiry ledger is customer data.** `~/.solo-ops` and any
`--ledger` file hold names, contact details, and event notes for real people.
The generated cockpit HTML embeds those rows verbatim — it is a local artifact,
not something to publish. Nothing is uploaded anywhere.

## No absolute paths or private hosts in releases

`make leak-scan` rejects absolute home paths and personal hosting subdomains,
and runs in CI. `make leak-scan-selftest` proves the scan can actually fail —
a gate never observed rejecting anything is not known to work.
