# Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values for enriching Contracts with real context.

**This step is REQUIRED** -- gen-contracts needs code context to produce accurate State dimensions, Input schemas, and Side-effect declarations.

**Generic reconnaissance reads**:

| Source | What to extract |
|--------|-----------------|
| CLI entry points | Command names, flag names, flag types, output patterns |
| API handlers | Request/response schemas, status codes, middleware |
| TUI model files | Model struct fields, Cmd definitions, Msg types, View rendering |
| Config files | Port numbers, base paths, timeout values, auth mechanisms |
| State storage | File paths, JSON schemas, database tables |
| Hook definitions | Hook names, trigger conditions, parameter schemas |

**TUI-specific reconnaissance** (when `tui` interface detected):

| Source | What to extract |
|--------|-----------------|
| Cmd definitions | Cmd function names, async behavior (do they return Msg?) |
| Batch usage | `tea.Batch()` calls and their Cmd arguments |
| Timeout configurations | Any timeout constants, default wait durations |
| Model transitions | Init -> Idle -> Processing -> Result states |

Build Fact Table with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| CLI_COMMAND_FEATURE | forge feature | cmd/feature.go:15 |
| TUI_AWAIT_TIMEOUT | 3000ms | internal/tui/config.go:8 |
```

<HARD-RULE>
- Every Fact Table value must cite source file and line number. Unknown sources -> `UNKNOWN`. Do not fabricate.
- Fact Table values inform State dimension and Input dimension declarations. Use them to ground semantic descriptors in real code.
- When the project does not expose a state query interface, set `state-verification: partial` or `state-verification: deferred` in the Contract.
</HARD-RULE>
