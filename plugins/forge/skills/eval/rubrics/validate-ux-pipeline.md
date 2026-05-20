# validate-ux Sub-Pipeline

Reference file for eval's validate-ux pre-processing. Loaded by `SKILL.md`.

## Project Type Detection

Resolve project type from `forge test interfaces`:

| Interface | Project Type | Execution Method | Operation Unit | Capture |
|------------|-------------|-----------------|----------------|---------|
| `cli` | CLI | Bash command | Shell command | stdout/stderr/exit code |
| `web-ui` | Web | agent-browser | URL + element selector + action | Screenshot + accessibility tree |
| `tui` | TUI | Bash stdin pipe | Key sequence (non-interactive only) | Terminal output |

Detection priority: project interfaces -> `forge test detect` -> ask user.

TUI constraint: first version covers non-interactive scenarios only (initial render, help output, invalid input response).

## PRD-to-Operation Translation

All project types use a hybrid translation strategy:

1. **Direct extraction**: scan PRD for code blocks, commands, URLs, key-binding descriptions
2. **Inference**: for missing concrete operations, the agent infers from auxiliary information

| Type | Auxiliary Information | Inference Method |
|------|----------------------|-----------------|
| CLI | Recursive `forge --help` subcommand discovery | Match PRD description -> subcommand -> argument format |
| Web | `sitemap.json` (accessibility tree + element IDs) | Match PRD description -> route -> DOM selector |
| TUI | Run program to capture initial screen + help output | Match PRD description -> menu option -> key-binding |

## ux-snapshot.md Format

```markdown
# UX Snapshot: <feature-name>

## Project Info
- Type: cli | web | tui
- Binary/URL: <path or url>
- PRD Reference: <path to PRD>
- Generated: <timestamp>

## Flow: <flow-name-from-PRD>

### Step 1: <action-description>
**Command/Navigate**: <what was executed>
**Input**: <what was sent>
**Output**:
`<raw stdout/stderr, screenshot path, or terminal capture>`
**Exit Code**: <cli only>

**Effect Verification**:
- Data: <expected data change> -> <actual result> pass/fail
- Side Effect: <unexpected changes checked via git diff --stat> -> pass/fail
- Output Consistency: <output claim vs reality> -> pass/fail
- Cascade: <downstream behavior triggered?> -> pass/fail

**Idempotency Check**:
- Re-run: <result of running same command again>

**State Integrity**:
- <consistency check between related state>

### Step 2: ...

## Standalone Checks

### Help Text
**Command**: `<binary> --help`
**Output**:
`<full help output>`

### Error Handling
**Command**: `<binary> invalid-command`
**Output**:
`<error output>`

### Version Info
**Command**: `<binary> --version`
**Output**:
`<version output>`
```

## Operation Impact Verification (7 Types)

| Impact Type | Verification Method | Example |
|-------------|-------------------|---------|
| Data Effect | Compare file/db/state before and after operation | `submit` updates index.json status |
| Side Effect | `git diff --stat` checks for unexpected file changes | `delete task` does not affect adjacent tasks |
| Idempotency | Re-execute the same operation | `submit` returns "already submitted" on second run |
| Output-Reality Consistency | Verify output claims match actual state | Output "created: X.md" -> file exists on disk |
| State Integrity | Check system-wide consistency after multi-step operations | Record file count matches index.json count |
| Cascade Effect | Check if downstream behavior is triggered | `submit` triggers quality-gate |
| Rollback Feasibility | Check state recoverability after operation failure | Failed operation leaves no residual dirty state |
