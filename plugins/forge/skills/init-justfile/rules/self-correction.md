# Self-Correction Rules

When a recipe fails in Phase 2 (actual execution), analyze the error and apply corrections:

| Error Pattern                         | Recipe                        | Fix                                                                                                                                                                                        |
| ------------------------------------- | ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `npm error Missing script: "start"`   | `run` (node/mixed)            | Replace `npm run start` -> `npm run preview` in justfile, retry                                                                                                                             |
| `npm error Missing script: "preview"` | `run` (node/mixed)            | Replace -> `npm run dev` in justfile, retry                                                                                                                                                 |
| `npm error Missing script: "dev"`     | `dev` (node/mixed)            | Replace -> `npm run start` in justfile, retry                                                                                                                                               |
| `can't load package: no Go files`     | `run`/`dev`/`compile` (go)    | Scan for `cmd/*/main.go`, update entry point in justfile, retry                                                                                                                            |
| `CGO_ENABLED=1` available             | `test` (go)                   | Add `-race` flag to `go test` recipe for race detection, retry                                                                                                                             |
| `command not found: golangci-lint`    | `lint`/`check` (go)           | In `lint`: replace `golangci-lint run ./...` -> `go vet ./...`. In `check`: replace `golangci-lint run ./... && go vet ./...` -> `go vet ./...`. Retry both.                                 |
| `command not found: uvicorn`          | `dev` (python)                | Replace -> `python -m src --reload` or skip with comment, retry                                                                                                                             |
| `command not found: ruff`             | `lint`/`fmt`/`check` (python) | In `lint`: replace `ruff check .` -> `python -m flake8`. In `check`: replace `ruff check .` -> `python -m flake8` (keep `&& python -m py_compile src/`). In `fmt`: skip with comment. Retry. |

For each fix:

1. Edit the justfile to apply the correction.
2. Re-run the failed command (actual execution, same method as Phase 2).
3. If it still fails after 2 attempts, leave the recipe as-is and report the failure in the output.
