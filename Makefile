.PHONY: test build install vet smoke leak-scan leak-scan-selftest gate publish-ready clean

BIN := bin/solo-ops

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -o $(BIN) ./cmd/solo-ops

install: build
	mkdir -p $(HOME)/.local/bin
	cp -f $(BIN) $(HOME)/.local/bin/solo-ops

# smoke runs every command documented in README.md against the committed
# fixtures. Unit tests exercise the internal packages and never the CLI, so a
# command that is documented but not run here is undefended against drift.
# Mutating verbs run against a copy in a temp dir so the fixture stays pristine.
smoke: build
	@echo "--- policy-pack"
	@./$(BIN) policy-pack --business "Painted by Demo" >/dev/null
	@echo "--- inquiry-list"
	@./$(BIN) inquiry-list --ledger testdata/inquiries.jsonl >/dev/null
	@echo "--- inquiry-add (against a copy)"
	@t=$$(mktemp -d); \
	  ./$(BIN) inquiry-add --ledger "$$t/new.jsonl" --name "Smoke" --event "smoke test" --guests 1 >/dev/null; \
	  ./$(BIN) inquiry-list --ledger "$$t/new.jsonl" | grep -q Smoke || { rm -rf "$$t"; echo "inquiry-add did not persist"; exit 1; }; \
	  rm -rf "$$t"
	@echo "--- inquiry-set (against a copy)"
	@t=$$(mktemp -d); cp -f testdata/inquiries.jsonl "$$t/led.jsonl"; \
	  ./$(BIN) inquiry-set --ledger "$$t/led.jsonl" --id inq-demo --status booked >/dev/null; \
	  ./$(BIN) inquiry-list --ledger "$$t/led.jsonl" | grep -q booked || { rm -rf "$$t"; echo "status did not persist"; exit 1; }; \
	  rm -rf "$$t"
	@echo "--- cockpit-html"
	@t=$$(mktemp -d); ./$(BIN) cockpit-html --ledger testdata/inquiries.jsonl --out "$$t/c.html" >/dev/null; \
	  test -s "$$t/c.html" || { rm -rf "$$t"; echo "cockpit html not written"; exit 1; }; \
	  rm -rf "$$t"
	@echo "--- serve (loopback only, then stop)"
	@./$(BIN) serve --ledger testdata/inquiries.jsonl --addr 127.0.0.1:8791 >/dev/null 2>&1 & \
	  pid=$$!; \
	  ok=0; \
	  for i in 1 2 3 4 5 6 7 8 9 10; do \
	    if curl -fsS --max-time 2 http://127.0.0.1:8791/ >/dev/null 2>&1; then ok=1; break; fi; \
	    sleep 0.3; \
	  done; \
	  kill $$pid 2>/dev/null; wait $$pid 2>/dev/null; \
	  [ $$ok = 1 ] || { echo "serve did not answer on loopback"; exit 1; }
	@echo "--- fixture unchanged by smoke"
	@git diff --quiet -- testdata/inquiries.jsonl
	@echo "smoke OK"

# leak-scan checks generic leak classes: absolute home paths and personal
# hosting subdomains. Patterns require a concrete leading label so they match a
# real hostname and not a prose mention of the wildcard form — a scan whose
# pattern allows an empty label matches its own documentation. The file list
# includes untracked-but-not-ignored files, because scanning only tracked files
# is blind to exactly the files a release adds. Do not write an example hostname
# anywhere in this file.
LEAK_FILES = git ls-files --cached --others --exclude-standard -z

leak-scan:
	@status=0; \
	set -- '/Users/[a-z]' '/home/[a-z]' '[a-z0-9][a-z0-9-]*\.workers\.dev'; \
	for pat in "$$@"; do \
	  found=$$($(LEAK_FILES) | xargs -0 grep -rilE -- "$$pat" 2>/dev/null | sort -u); \
	  hits=$$(printf '%s' "$$found" | grep -c . || true); \
	  printf '%-34s hits=%s\n' "$$pat" "$$hits"; \
	  if [ "$$hits" != "0" ]; then printf '%s\n' "$$found" | sed 's/^/    /'; status=1; fi; \
	done; \
	[ $$status = 0 ] && echo "leak-scan OK" || { echo "leak-scan FAILED"; exit 1; }

# A gate never observed rejecting anything is not known to work.
leak-scan-selftest:
	@w=workers; d=dev; probe=./.leak-probe.md; \
	printf 'contact https://example-private.%s.%s/admin\n' "$$w" "$$d" > "$$probe"; \
	if $(MAKE) --no-print-directory leak-scan >/dev/null 2>&1; then \
	  rm -f "$$probe"; echo "leak-scan-selftest FAILED: scan accepted a planted leak"; exit 1; \
	fi; \
	rm -f "$$probe"; \
	$(MAKE) --no-print-directory leak-scan >/dev/null 2>&1 || { \
	  echo "leak-scan-selftest FAILED: scan still red after probe removal"; exit 1; }; \
	echo "leak-scan-selftest OK — planted leak rejected, clean tree accepted"

gate: vet test smoke leak-scan leak-scan-selftest
	@echo "solo-ops gate OK"

publish-ready: gate
	@echo "solo-ops publish-ready OK"

clean:
	rm -rf bin
