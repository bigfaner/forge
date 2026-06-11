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

Build Fact Table with source citations. Each fact entry follows the canonical JSON schema (defined in `forge-cli/pkg/facttable/facttable.go`):

```json
{
  "fact_id": "CLI_COMMAND_FEATURE",
  "source": "static",
  "subject": "CLI command for creating features",
  "kind": "signature",
  "value": "forge feature",
  "confidence": "inferred",
  "updated_at": "<ISO8601 timestamp>"
}
```

```json
{
  "fact_id": "TUI_AWAIT_TIMEOUT",
  "source": "static",
  "subject": "TUI async await timeout duration",
  "kind": "output_format",
  "value": "3000ms",
  "confidence": "inferred",
  "updated_at": "<ISO8601 timestamp>"
}
```

All entries use `"source": "static"` to distinguish from runtime facts (added by Run-to-Learn with `"source": "runtime"`). Static entries default to `"confidence": "inferred"` (confirmed at runtime by R2L). The `kind` field must be one of: `signature`, `output_format`, `error_code`, `side_effect`, `precondition`.

During reconnaissance, a Markdown summary table may be used as an intermediate scratchpad for AI reasoning, but the final output written to `.forge/fact-table.json` MUST use this JSON format.

<HARD-RULE>
- Every Fact Table value must cite source file and line number. Unknown sources -> `UNKNOWN`. Do not fabricate.
- Fact Table values inform State dimension and Input dimension declarations. Use them to ground semantic descriptors in real code.
- When the project does not expose a state query interface, set `state-verification: partial` or `state-verification: deferred` in the Contract.
</HARD-RULE>
