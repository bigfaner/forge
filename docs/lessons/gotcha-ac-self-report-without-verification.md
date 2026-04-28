# Gotcha: Agent Self-Reports AC as Met Without Running the Artifact

## Problem

A `gen-test-scripts` task was marked `completed` with all acceptance criteria `met: true`.
The generated scripts contained wrong API path prefixes, wrong directory depth, and wrong port.
The bugs were only caught in the next session when the scripts were actually inspected.

## Root Cause

`task record` validates **record structure**, not **AC truth**:

- It checks: `status=completed` → `testsPassed+testsFailed > 0` OR `coverage=-1`
- It does NOT check: whether the agent actually ran the generated artifact

The sub-agent that executed T-test-2 inspected the files it created (they existed, had traceability comments) and marked AC as met — without running `node --test` or even a `tsc --noEmit` to verify the scripts compiled and hit real endpoints.

Causal chain:
```
scripts generated with wrong paths
  → agent checks file existence only (AC text says "created", not "runs")
    → all AC marked met: true in record.json
      → task record accepts the record (structure is valid)
        → task status = completed, bug ships
```

## Solution

For tasks that generate executable artifacts (scripts, binaries, configs), AC must include a **run/compile check**, not just a file-existence check.

**Weak AC (what we had):**
```
- [x] testing/scripts/package.json created
- [x] At least one spec file generated (api.spec.ts)
```

**Strong AC (what to write):**
```
- [ ] testing/scripts/package.json created
- [ ] `cd testing/scripts && npm install && npx tsc --noEmit` exits 0
- [ ] At least one spec file generated (api.spec.ts)
```

If the environment doesn't allow running the scripts (no live server), use `coverage: -1` and add a note explaining why execution was skipped — don't silently mark AC as met.

## Key Takeaway

`task record` trusts the agent's self-report. The only enforcement is structural (field presence, test count). **AC quality is the real gate** — if AC only checks file existence, the CLI cannot catch a broken artifact. Write AC that require execution evidence for any task that produces runnable output.
