# Solo Ops

Offline inquiry CRM (JSONL) + policy pack + local HTML cockpit for one-person service businesses.

## Usage

```sh
go test ./...
go run ./cmd/solo-ops policy-pack --business "Painted by Demo"
go run ./cmd/solo-ops inquiry-list --ledger testdata/inquiries.jsonl
go run ./cmd/solo-ops inquiry-set --ledger testdata/inquiries.jsonl --id inq-demo --status booked
go run ./cmd/solo-ops cockpit-html --ledger testdata/inquiries.jsonl --out /tmp/solo-ops.html
# optional live serve:
# go run ./cmd/solo-ops serve --ledger testdata/inquiries.jsonl --addr 127.0.0.1:8791
```
