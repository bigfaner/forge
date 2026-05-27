---
name: step-0.5-validation
description: Surface detection and validation logic for determining project interface type and per-surface generation strategy
---

# Step 0.5: Surface Detection

Determine the project's interface surface type to drive per-surface generation strategy.

## 0.5.1 Read Surface Configuration

Read `.forge/config.yaml` from the project root and extract the `surface` field.

```bash
forge surfaces
```

| Result | Action |
|--------|--------|
| Surface type returned (e.g., `cli`, `tui`, `web`, `mobile`, `api`) | Use this as the active surface type for generation |
| No config file or field missing | Proceed to auto-detection (Step 0.5.2) |

## 0.5.2 Auto-Detection Fallback

If `.forge/config.yaml` does not contain a `surface` field, infer the surface type from code reconnaissance signals in Step 1. Use the Verification Method defined in each type file (`types/cli.md`, `types/api.md`, etc.) to probe for interface indicators.

Priority: config value > auto-detection. If auto-detection is ambiguous (multiple types detected), ask the user which surface type to prioritize.

## 0.5.3 Surface Strategy Application

The detected surface type determines the **test ratio strategy** for Step 3 generation:

| Surface Type | Contract : Journey Ratio | Key Generation Constraint |
|-------------|--------------------------|---------------------------|
| CLI | >= 80% Contract | Subprocess execution model, binary isolation, environment hermeticity |
| TUI | >= 80% Contract | Terminal I/O testing, non-interactive stdin pipe, ANSI sanitization |
| Web | Balanced 50/50 | Convention-defined browser framework, session reuse, network interception |
| API | Balanced 50/50 | HTTP client testing, status code coverage, content-type verification |
| Mobile | Best-effort | Maestro YAML skeleton + deep link tests, complex scenarios marked `manual-only` |

**Contract test ratio formula**: `Contract test functions / (Contract test functions + Journey smoke test functions) x 100%`

<HARD-RULE>
The surface type determines generation strategy -- test ratio, execution model, and assertion patterns. Type-specific Golden Rules (from `types/<type>.md`) take precedence over generic generation rules. Convention provides framework implementation details, Surface type provides strategy constraints. These two are orthogonal and merged at generation time.
</HARD-RULE>

## Surface-Driven Generation Strategy (Step 3.0)

Apply the surface type detected in Step 0.5 to constrain the generation plan. Each surface type has distinct ratio targets, execution models, and generation constraints.

### CLI Surface (Contract >= 80%)

- **Primary focus**: Contract test functions -- one test per Outcome per step
- **Execution model**: Subprocess execution (binary isolation, environment hermeticity per `types/cli.md`)
- **Journey smoke tests**: Generate exactly 1 smoke test per Journey (happy path only)
- **Ratio enforcement**: For N Contract steps with M total Outcomes, generate M Contract test functions + 1 Journey smoke test
- **Binary check**: Verify the binary can be built before generating tests. Auto-detect binary name and build command from Fact Table.

### TUI Surface (Contract >= 80%)

- **Primary focus**: Contract test functions -- one test per Outcome per step
- **Execution model**: Non-interactive stdin pipe with terminal output capture (per `types/tui.md`)
- **Journey smoke tests**: Generate exactly 1 smoke test per Journey (happy path only)
- **Ratio enforcement**: Same formula as CLI -- M Contract functions + 1 Journey smoke test

### Web Surface (Balanced 50/50)

- **Balanced approach**: Generate Contract tests for each Outcome AND enrich the Journey smoke test with multi-step verification
- **Execution model**: Convention-defined browser framework (per `types/ui.md`)
- **Journey smoke tests**: Generate 1 comprehensive smoke test that verifies the happy path AND at least 1 failure path
- **Ratio target**: Approximately equal Contract test functions and Journey smoke test functions

### API Surface (Balanced 50/50)

- **Balanced approach**: Generate Contract tests for each Outcome AND enrich the Journey smoke test
- **Execution model**: HTTP client testing (per `types/api.md`)
- **Journey smoke tests**: Generate 1 comprehensive smoke test that verifies the happy path AND at least 1 error path
- **Ratio target**: Approximately equal Contract test functions and Journey smoke test functions

### Mobile Surface (Best-Effort)

- **Maestro YAML skeleton**: Generate Maestro YAML flows instead of code-based test functions
- **Skeleton structure**: Each generated Maestro YAML file MUST contain:
  1. `appId` declaration (from Fact Table `MOBILE_APP_ID`)
  2. `onFlowStart: [launchApp]` lifecycle hook
  3. `onFlowEnd: [killApp]` lifecycle hook
  4. Navigation flow for the Contract step's happy path
  5. Deep link test variant when the Contract involves navigation to a specific screen
- **Deep link tests**: For each Journey step that navigates to a specific screen, generate an additional Maestro YAML that opens the app via URL scheme and asserts the target screen is visible
- **Complex scenario handling**: If a test case involves gestures not expressible in Maestro (pinch, rotate, multi-finger swipe), or requires physical device capabilities (sensors, camera), mark the test with `manual-only` annotation and skip generation. Add a comment explaining which capability requires manual testing.
- **Convention reference**: Use Maestro YAML syntax conventions. If no Mobile Convention file exists, use the Maestro reference from `types/mobile.md` as the authoritative syntax guide.

<HARD-RULE>
**Mobile best-effort**: Do not aim for comprehensive coverage. Generate skeleton flows + deep link tests for core Journeys only. Complex scenarios MUST be marked `manual-only` rather than generating incomplete or fragile tests. Mobile test generation MUST NOT fail the pipeline -- any generation issue should result in a skeleton with `manual-only` markers.
</HARD-RULE>
