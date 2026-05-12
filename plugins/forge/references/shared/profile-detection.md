# Test Profile Auto-Detection Rules

When `.forge/config.yaml` does not exist, use these rules to detect the appropriate test profile.

## Detection Order

Check signals top-to-bottom; first match wins.

| # | Signal | Profile |
|---|--------|---------|
| 1 | `package.json` exists AND (`playwright.config.ts` or `playwright.config.js` exists OR `@playwright/test` in devDependencies) | `web-playwright` |
| 2 | `go.mod` exists | `go-test` |
| 3 | `android/` or `ios/` directory exists at project root | `maestro` |
| 4 | `pom.xml` or `build.gradle` or `build.gradle.kts` exists | `java-junit` |
| 5 | `Cargo.toml` exists | `rust-test` |
| 6 | (`requirements.txt` or `pyproject.toml`) exists AND `pytest` listed as dependency | `pytest` |
| 7 | `package.json` exists AND frontend router detected (`src/App.tsx`, `src/router.tsx`, `next.config.*`, `nuxt.config.*`, `vite.config.*`) | `web-playwright` |

## No Match

If none of the above signals match, **ask the user explicitly**. Do NOT silently default to any profile.

Present the user with the list of available profiles and ask them to choose:

```
Could not auto-detect a test profile for this project. Please select one:

1. web-playwright — Web UI/API/CLI (Playwright + TypeScript)
2. go-test — Go CLI/TUI/backend (go test)
3. maestro — React Native / Flutter (Maestro CLI)
4. java-junit — Java CLI/backend (JUnit 5 + Maven)
5. rust-test — Rust CLI/backend (cargo test)
6. pytest — Python CLI/backend (pytest)
```

After the user selects, write `.forge/config.yaml` with their choice so future runs don't need to ask again.

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
