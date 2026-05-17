# Test Language Auto-Detection Rules

When `.forge/config.yaml` does not specify `languages`, `forge testing detect` scans the project root for file signals. Multiple signals may produce multiple languages (multi-language support).

## Detection Rules

All matching rules apply — languages are accumulated, not exclusive.

| # | Signal | Language |
|---|--------|----------|
| 1 | `package.json` exists AND (`playwright.config.*` exists OR `@playwright/test` in devDependencies/dependencies) | `javascript` |
| 2 | `go.mod` exists | `go` |
| 3 | `android/` or `ios/` directory exists at project root | `mobile` |
| 4 | `pom.xml` or `build.gradle` or `build.gradle.kts` exists | `java` |
| 5 | `Cargo.toml` exists | `rust` |
| 6 | (`requirements.txt` or `pyproject.toml`) contains `pytest` | `python` |

## Resolution Priority

1. **Config file** (`.forge/config.yaml` `languages` key) — always wins
2. **Auto-detection** — scans project root using rules above
3. **No match** — outputs empty result, skills ask user to configure `languages`

## Multi-Language

A project may match multiple rules. For example, a Go backend with a Python CLI tool:

```
go.mod          → go
requirements.txt with pytest → python
Result: [go, python]
```

When multiple languages are active, test tasks are expanded per-language with letter suffixes (a, b, c, ...).

## Interfaces Reference

Each language supports a fixed set of interfaces:

| Language | Interfaces |
|----------|-----------|
| javascript | web-ui, api, cli |
| go | tui, api, cli |
| mobile | mobile-ui, api |
| java | tui, api, cli |
| rust | tui, api, cli |
| python | api, cli |

Interface meanings:
- `web-ui` — Browser UI (DOM interaction)
- `tui` — Terminal UI (text rendering, keyboard interaction)
- `mobile-ui` — Mobile UI (touch, gestures)
- `api` — HTTP/network interface
- `cli` — Command-line interface
