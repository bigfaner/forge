# Contract: skill-ops / Step 1: Plugin Validation

## Outcome "skill-files-no-raw-commands"
- Preconditions: "plugin files exist at expected paths"
- Input: "read content of skill/agent/command files"
- Output: "zero occurrences of raw toolchain commands (go test, npm run build, etc.)"
- State: "no state changes"
- Side-effect: none

## Forbidden commands
- `go test ./...`, `go build ./...`, `go vet ./...`
- `npm run build`, `npm test`, `npm test -- --coverage`
- `npx serve`, `cargo build`, `pytest --cov=`
- `go test -cover ./...`, `go test -race -cover ./...`
- `npm run build && npm test`, `cd tests/e2e && npm install`

## Journey Invariants
- file paths relative to project root
- all skill/agent/command files validated against same forbidden list
