# Maestro Generate Strategy

Profile-specific test generation rules for the `gen-test-scripts` skill.

## Test Runner & Format

| Item | Value |
|------|-------|
| Test runner | Maestro CLI |
| Test format | YAML flows |
| File extension | `.yaml` |

## Flow Structure

```yaml
appId: ${APP_ID}
name: TC-NNN Description
env:
  KEY: ${VALUE}
onFlowStart:
  - launchApp
onFlowEnd:
  - killApp
commands:
  - tapOn: ...
  - assertVisible: ...
```

## Element Selectors

| Priority | Method | Example |
|----------|--------|---------|
| 1 | Text match | `tapOn: { text: "Login" }` |
| 2 | Resource ID | `tapOn: { id: "com.app:id/button" }` |
| 3 | Index | `tapOn: { index: 0 }` |
| 4 | Relative | `tapOn: { below: { text: "Label" } }` |

## Commands Reference

| Command | Purpose |
|---------|---------|
| `tapOn` | Tap element |
| `assertVisible` | Assert element visible |
| `assertNotVisible` | Assert element not visible |
| `inputText` | Type text into focused field |
| `swipe` | Swipe gesture (directions: UP, DOWN, LEFT, RIGHT) |
| `scroll` | Scroll to element |
| `back` | Press back button |
| `launchApp` | Launch app |
| `killApp` | Kill app |
| `takeScreenshot` | Capture screenshot |

## Auth

Use `env` section for credentials, `inputText` for login fields:

```yaml
env:
  USERNAME: ${TEST_USERNAME}
  PASSWORD: ${TEST_PASSWORD}
commands:
  - tapOn: { text: "Email" }
  - inputText: ${USERNAME}
  - tapOn: { text: "Password" }
  - inputText: ${PASSWORD}
  - tapOn: { text: "Sign in" }
```

## Anti-Patterns

- No `extendedWait` with arbitrary durations — use `assertVisible` which waits automatically
- No hardcoded device dimensions — Maestro is device-agnostic

## Traceability

Flow `name` field includes TC ID:

```yaml
name: TC-NNN Description
# Traceability: TC-NNN → {PRD Source}
```

## Screenshots

```yaml
- takeScreenshot: TC-NNN
```
