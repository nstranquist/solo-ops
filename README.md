# Solo Ops

Offline inquiry CRM (JSONL) + policy pack generator for one-person service businesses.

## Why this tech

Go + JSONL ledger (no cloud). Rejected Next/Postgres SaaS for offline operator reliability.

## Usage

```sh
go test ./...
go run ./cmd/solo-ops policy-pack --business "Painted by Demo"
go run ./cmd/solo-ops inquiry-list --ledger testdata/inquiries.jsonl
```
