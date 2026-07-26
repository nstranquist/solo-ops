# Contributing to solo-ops

## The loop

```sh
make gate           # vet + test + smoke + leak-scan + selftest
make publish-ready  # run before tagging
```

## Rules that matter here

1. **`make smoke` must exercise every command the README documents.** The unit
   tests cover the internal packages and never execute the CLI, so a documented
   command that smoke does not run is undefended — that is precisely how a
   sibling project shipped a README naming four verbs that did not exist.
2. **Keep behaviour deterministic.** No network, no clock dependence, and no LLM
   in any path that decides pass/fail.
3. **Fixtures stay synthetic.** Nothing under `testdata/` may contain real
   personal data — see `SECURITY.md`.
4. **No absolute local paths in committed files**, including README examples.
   `make leak-scan` enforces this.
