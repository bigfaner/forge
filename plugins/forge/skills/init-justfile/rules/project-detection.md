# Project Type Detection

## Detection Signal Mapping

```bash
ls package.json go.mod Cargo.toml pyproject.toml 2>/dev/null
```

**Detection signal mapping:**

| Marker File      | Signal   |
| ---------------- | -------- |
| `package.json`   | frontend |
| `go.mod`         | backend  |
| `Cargo.toml`     | backend  |
| `pyproject.toml` | backend  |

## Classification Algorithm

1. Check for each marker file's existence in the project root.
2. Count frontend vs backend signals and classify:
   - Exactly one frontend signal AND exactly one backend signal -> **`mixed`**
   - Exactly one frontend signal, no backend signals -> **`frontend`**
   - Exactly one backend signal, no frontend signals -> **`backend`**
   - Neither set has signals -> **Error**: "Error: no known project markers detected (expected one of: package.json, go.mod, Cargo.toml, pyproject.toml)" -- abort, do NOT generate a justfile.
   - Multiple backend signals (e.g. `go.mod` + `Cargo.toml`) -> **Error**: "Error: multiple backend markers detected -- not supported" -- abort.

## Mixed Project Root Path Detection

For `mixed` projects, detect root paths:

```bash
find . -name package.json -not -path '*/node_modules/*' -maxdepth 3 | head -1 | xargs dirname
find . \( -name go.mod -o -name Cargo.toml -o -name pyproject.toml \) -maxdepth 3 | head -1 | xargs dirname
```

Record these as `FRONTEND_DIR` and `BACKEND_DIR` (e.g. `./frontend`, `./backend`). Use `.` if the marker is in the project root.

## Backend Entry Point Detection

For **backend** and **mixed** projects, detect the entry point:

| Language | Marker           | Entry point detection (`BACKEND_ENTRY`)                                           |
| -------- | ---------------- | --------------------------------------------------------------------------------- |
| Go       | `go.mod`         | `ls cmd/*/main.go` -> `cmd/<name>/main.go`; else `ls main.go` -> `.`              |
| Rust     | `Cargo.toml`     | `grep '\[\[bin\]\]' Cargo.toml` -> `--bin <name>`; else empty                      |
| Python   | `pyproject.toml` | `ls src/__init__.py` -> `-m src`; `ls main.py` -> `main.py`; `ls app.py` -> `app.py` |

Record `BACKEND_ENTRY` from the detected entry point.

## Frontend Run Script Detection

For **frontend** and **mixed** projects, detect available npm scripts:

```bash
node -e "const s=JSON.parse(require('fs').readFileSync('package.json')).scripts||{}; console.log(s.start ? 'start' : s.preview ? 'preview' : 'dev')"
```

Record as `FRONTEND_RUN_SCRIPT` (e.g. `start`, `preview`, or `dev`).
