---
name: env-check
description: Per-surface environment readiness checks before test execution. Defines detection items and repair suggestions for each surface type.
---

# Environment Readiness Detection

Before executing tests, verify that the execution environment meets the requirements for the detected surface type. Each surface type has different readiness concerns.

## How It Works

1. Read the current surface type from `forge config get surface`
2. Look up the corresponding surface rule file in `skills/gen-journeys/rules/surface-<type>.md` -- specifically the "Environment Readiness Checks" table
3. Execute each check item listed for that surface type
4. Report results: all-pass -> proceed to test execution; any-fail -> output diagnostics and abort

<HARD-RULE>
Environment detection failure does NOT auto-fix. Only output diagnostic information and repair suggestions. The user must fix the environment themselves, then re-run.
</HARD-RULE>

## Per-Surface Detection Items

The following sections define the detection items for each built-in surface type. When adding a new surface type, include an "Environment Readiness Checks" table in its surface rule file; this skill reads that table automatically.

### CLI

| # | Check | How to Verify | Blocking | Repair Suggestion |
|---|-------|--------------|----------|-------------------|
| 1 | Binary compiles | Run `just compile` or the project's build command. Exit code 0 = pass | Yes | Fix compilation errors, then retry |
| 2 | Binary is executable | Check the compiled binary file exists and has execute permission (`test -x <path>`) | Yes | Run the build command to produce the binary |
| 3 | Required external tools available | Run `which <tool>` for each external tool referenced in test scripts. Exit code 0 = pass | Yes | Install the missing tool |
| 4 | Config files accessible | Verify config file paths referenced in tests exist or can be created | No | Create missing config files or adjust test paths |

### TUI

| # | Check | How to Verify | Blocking | Repair Suggestion |
|---|-------|--------------|----------|-------------------|
| 1 | Binary compiles | Run `just compile` or the project's build command. Exit code 0 = pass | Yes | Fix compilation errors, then retry |
| 2 | Binary is executable | Check the compiled binary file exists and has execute permission | Yes | Run the build command to produce the binary |
| 3 | Stdin pipe works | Spawn a test subprocess with stdin piped, verify it accepts input without error | Yes | Verify no terminal-exclusive mode is required for basic input |
| 4 | No GUI dependency | Verify the TUI does not require X11/Wayland/display server (check env vars or process dependencies) | No | Set `TERM=dumb` or use headless terminal if GUI deps are present |

### WebUI

| # | Check | How to Verify | Blocking | Repair Suggestion |
|---|-------|--------------|----------|-------------------|
| 1 | Dev server starts | Start the dev server (`npm run dev` or equivalent). Verify it responds on the expected port within 30s | Yes | Check port availability, fix startup errors in server logs |
| 2 | Dev server responds | HTTP GET to dev server root returns 200 | Yes | Check server health endpoint, verify application routes are configured |
| 3 | Browser automation framework installed | Run framework install command (e.g., `npx playwright install --dry-run`). Exit code 0 = pass | Yes | Run `npx playwright install` (or equivalent for your framework) |
| 4 | Test database seeded | Verify test data fixtures are loaded (check a known test record exists) | No | Run seed command or setup script |

### API

| # | Check | How to Verify | Blocking | Repair Suggestion |
|---|-------|--------------|----------|-------------------|
| 1 | Server starts | Start the application server. Verify it binds to the expected port without error | Yes | Check port availability, fix startup errors in server logs |
| 2 | Server responds | HTTP GET to health endpoint returns 200 | Yes | Check server configuration, verify health endpoint is registered |
| 3 | Database connected | Query a known table or run a health-check query. Returns result = pass | Yes | Verify DB connection string, check DB server is running, run migrations |
| 4 | Authentication configured | Verify test API keys or tokens are available (env vars or config files exist) | Yes | Set required auth environment variables or create test credentials |

### Mobile

| # | Check | How to Verify | Blocking | Repair Suggestion |
|---|-------|--------------|----------|-------------------|
| 1 | Maestro CLI installed | Run `which maestro`. Exit code 0 = installed | No | Install Maestro CLI: `curl -Ls "https://get.maestro.mobile.dev" \| bash` |
| 2 | Emulator/simulator available | Run `maestro devices`. Lists at least one device = pass | No | Start Android emulator or iOS simulator |
| 3 | App binary available | Check APK/IPA exists at expected path | No | Build the app binary first |

<HARD-RULE>
Mobile checks are **best-effort** -- all items are non-blocking. Missing Maestro CLI does not prevent test generation or block the pipeline. It only prevents test execution.
</HARD-RULE>

## Output Format

When environment is NOT ready, output a diagnostic report:

```
Environment Readiness: NOT READY (N/M checks passed)

Missing:
  - [CLI-1] Binary does not compile
    Suggestion: Fix compilation errors, then retry
  - [CLI-3] External tool "jq" not found
    Suggestion: Install the missing tool

Fix the issues above, then re-run /run-tests.
```

When environment IS ready:

```
Environment Readiness: READY (M/M checks passed)
```

## Integration with SKILL.md

This rule file is consumed during the run-tests workflow. The SKILL.md should:

1. Read surface type before test execution (between Step 3 Setup and Step 4 Pre-check, or as a new early step)
2. Execute the per-surface checks defined above
3. On failure: output diagnostics and abort (exit code 1, retryable)
4. On success: proceed to test execution

**Extensibility**: New surface types only need to add an "Environment Readiness Checks" table to their surface rule file in `skills/gen-journeys/rules/surface-<type>.md`. This env-check.md defines the detection framework; the surface rule files provide the per-type items.
