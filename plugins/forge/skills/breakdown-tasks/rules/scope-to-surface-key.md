---
name: scope-to-surface-key
description: Surface-key/type resolution via CLI instead of hardcoded path classification
---

# Scope-to-Surface-Key Migration

This rule replaces the old `rules/scope-assignment.md` hardcoded path classification logic with dynamic CLI-based surface resolution.

## Resolution Procedure

For each task, resolve `surface-key` and `surface-type` using `forge surfaces --json`:

### 1. Collect Affected File Paths

Gather all file paths from the task's affected files list (implementation files, test files, config files).

### 2. Query Surface for Each File

```bash
forge surfaces --json <file-path>
```

This returns one of:

| Output | Meaning | Action |
|--------|---------|--------|
| `[{"key": "admin-panel", "type": "web"}]` | Single match | Use this key+type |
| `[{"key": ".", "type": "web"}]` | Scalar form (single surface project) | Use key `"."`, type `"web"` |
| `[]` | No match for this file | File is not covered by any surface |
| stderr: `{"error": "no surface configured..."}` | No surfaces configured at all | See Error Handling below |

### 3. Merge Results

- **All files return the same surface-key** → use that key and type
- **Files return different surface-keys** → task spans multiple surfaces → set `surface-key: ""` and `surface-type: ""` (empty, meaning cross-surface)
- **Some files have no match (`[]`)** → ignore unmatched files; if ALL files are unmatched, set both to empty

### 4. Write Frontmatter

Populate task frontmatter:
```yaml
surface-key: "admin-panel"    # or "" for cross-surface / unknown
surface-type: "web"           # or "" for cross-surface / unknown
```

## Error Handling

### `forge surfaces --json` invocation fails (non-zero exit, no JSON output)

Output to the agent's execution context:
```
ERROR: forge surfaces --json failed with exit code <N>.
Recovery: verify that forge-cli is built and installed (run `forge version`).
If the error persists, set surface-key and surface-type to empty strings and continue.
```

### No surfaces configured (stderr JSON error, exit 1)

This is expected for projects without `.forge/config.yaml` surfaces field. Action:
- Set both `surface-key: ""` and `surface-type: ""` (empty values)
- Continue task generation without surface information
- Do NOT block or fail the task generation process

### Mixed or ambiguous results

When a task touches files from multiple surfaces, leave both fields empty. This signals "cross-surface" to downstream consumers (run-tests, init-justfile).

## Relationship to Old Scope Assignment

This rule fully replaces `rules/scope-assignment.md`. The old hardcoded path classification (`frontend`/`backend`/`all`) is deprecated. The new surface-key/type values come from user configuration in `.forge/config.yaml` via the CLI, not from path heuristics.

| Old Field | New Field | Source |
|-----------|-----------|--------|
| `scope: "frontend"` | `surface-key: "admin-panel"` | `forge surfaces --json <path>` |
| (no equivalent) | `surface-type: "web"` | `forge surfaces --json <path>` |
| `scope: "all"` | `surface-key: ""` + `surface-type: ""` | Cross-surface or unconfigured |
