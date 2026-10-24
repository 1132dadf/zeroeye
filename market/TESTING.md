# Market module testing guide

Run commands from the repository root unless noted.

## Focused unit tests

```bash
cd market
go test ./...
```

## Race detector

```bash
cd market
go test -race ./...
```

Use the race detector for packages that share order books, trade history, WebSocket clients, or mutable configuration. If a race is found, capture the package and reproducer rather than suppressing the test.

## Static checks

```bash
cd market
go vet ./...
```

## Full repository validation

```bash
cd ..
python3 build.py
```

The root build is the repository's documented diagnostic workflow. Include the real generated diagnostic artifact in the PR as required by the root contribution guide. Never create a placeholder log or report a check as passing if it did not run.

## Before opening a PR

- [ ] Add a regression test for each behavior change.
- [ ] Include boundary and malformed-input cases.
- [ ] Run focused package tests.
- [ ] Run `go test -race ./...` when concurrency is involved.
- [ ] Run `go vet ./...`.
- [ ] Run `python3 build.py` and include the real diagnostic artifact.
- [ ] Review the root PR template and state what remains unverified.
