---
id: "1.gate"
title: "Phase 1 Exit Gate"
priority: "P0"
estimated_time: "30min"
dependencies: ["1.summary"]
status: pending
breaking: true
---

# 1.gate: Phase 1 Exit Gate

## Description

Exit verification gate for Phase 1. Confirms that the new `e2e-setup` and `e2e-verify` recipes are syntactically correct and complete in `init-justfile.md` before Phase 2 begins editing skill files that reference these targets.

## Verification Checklist

1. [ ] `grep -c 'e2e-setup\|e2e-verify' plugins/forge/commands/init-justfile.md` >= 4
2. [ ] `e2e-setup` recipe contains `set -euo pipefail`, package.json existence check, idempotent node_modules check, `playwright install chromium`, and `echo "OK: e2e dependencies ready"`
3. [ ] `e2e-verify` recipe contains `[arg("feature", long)]`, empty-string guard, directory-existence check (`tests/e2e/{{feature}}/`), grep scan for `// VERIFY:`, and correct exit 0/1 outputs
4. [ ] Step 4 Output Confirmation in init-justfile.md lists both new targets
5. [ ] No deviations from `design/tech-design.md` Interface 1 and Interface 2 specs

## Reference Files

- `plugins/forge/commands/init-justfile.md`
- `docs/features/justfile-e2e-integration/design/tech-design.md` — Interface 1 and Interface 2

## Acceptance Criteria

- [ ] All verification checklist items pass
- [ ] Record created via `/record-task` with `coverage: -1.0`

## Implementation Notes

Verification-only task. No new content should be written. If recipe syntax is wrong, fix inline and document as a decision.
