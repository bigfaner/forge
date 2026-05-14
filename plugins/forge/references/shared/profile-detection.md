# Test Profile Auto-Detection Rules

When `.forge/config.yaml` does not specify `test-profiles`, `forge profile detect` scans the project root for file signals. Multiple signals may produce multiple profiles (multi-profile support).

## Detection Rules

All matching rules apply — profiles are accumulated, not exclusive.

| # | Signal | Profile |
|---|--------|---------|
| 1 | `package.json` exists AND (`playwright.config.*` exists OR `@playwright/test` in devDependencies/dependencies) | `web-playwright` |
| 2 | `go.mod` exists | `go-test` |
| 3 | `android/` or `ios/` directory exists at project root | `maestro` |
| 4 | `pom.xml` or `build.gradle` or `build.gradle.kts` exists | `java-junit` |
| 5 | `Cargo.toml` exists | `rust-test` |
| 6 | (`requirements.txt` or `pyproject.toml`) contains `pytest` | `pytest` |
| 7 | `package.json` exists AND no Playwright detected (fallback) | `web-playwright` |

## Resolution Priority

1. **Config file** (`.forge/config.yaml` `test-profiles` key) — always wins
2. **Auto-detection** — scans project root using rules above
3. **No match** — outputs `PROFILE: (none)`, skills ask user to choose

## Multi-Profile

A project may match multiple rules. For example, a Go backend with a Python CLI tool:

```
go.mod          → go-test
requirements.txt with pytest → pytest
Result: [go-test, pytest]
```

When multiple profiles are active, test tasks are expanded per-profile with letter suffixes (a, b, c, ...).

## Capabilities Reference

Each profile declares a fixed set of capabilities:

| Profile | Capabilities |
|---------|-------------|
| web-playwright | web-ui, api, cli |
| go-test | tui, api, cli |
| maestro | mobile-ui, api |
| java-junit | tui, api, cli |
| rust-test | tui, api, cli |
| pytest | api, cli |

Capability meanings:
- `web-ui` — Browser UI (DOM interaction)
- `tui` — Terminal UI (text rendering, keyboard interaction)
- `mobile-ui` — Mobile UI (touch, gestures)
- `api` — HTTP/network interface
- `cli` — Command-line interface
