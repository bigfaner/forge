---
date: 2026-05-16
type: investigation
trigger: eval-forge D2 audit flagged "noTest frontmatter injection" as potential bypass vector
status: closed
---

# noTest Field Lifecycle Investigation

## Trigger

eval-forge D2 bypass resistance audit identified `noTest: true` frontmatter injection as a potential HIGH-severity bypass vector: a lazy agent could add `noTest: true` to any task's frontmatter to skip quality gates. Investigation was needed to determine whether this field is deprecated or still active.

## Finding

**`noTest` is NOT deprecated.** Two distinct concepts were conflated:

| Concept | Location | Status |
|---------|----------|--------|
| `--no-test` CLI flag (`BuildIndexOpts.NoTest`) | `forge-cli/internal/cmd/` | **Removed** in commit `2d82321`, replaced by `isDocsOnlyFeature()` auto-detection |
| `Task.NoTest` per-task field | `forge-cli/pkg/task/types.go:91` | **Active** — fully functional, extensively used |

## Commit History

| Commit | Action | Scope |
|--------|--------|-------|
| `16f9eba` | Added `NoTest` to `Task`/`TaskState` structs, coverage auto-set, quality gate skip | Full feature |
| `31786d1` | Removed noTest references from command/SKILL.md documentation | Docs only |
| `ebc413a`/`15a629a` | Removed noTest from task template frontmatter | Templates |
| `2d82321` | Removed `--no-test` CLI flag, replaced by `isDocsOnlyFeature()` auto-detection | CLI flag only |
| `b8c7748` | Attempted full removal of `Task.NoTest` from structs | Breaking change |
| `23a6dd1`+ | Re-added `Task.NoTest` — auto-generated test tasks need it | Restored |

The removal in `b8c7748` was effectively reversed because auto-generated tasks (T-test-1, T-test-1b, T-eval-doc, T-quick-6) require the field to correctly skip quality gates.

## Current Active Usage

### Definition

```
types.go:91    → NoTest bool `json:"noTest,omitempty"`
frontmatter.go:20 → NoTest bool `yaml:"noTest"` (parsed from task .md YAML)
```

### Consumers

| File | Line(s) | Behavior |
|------|---------|----------|
| `submit.go` | 137 | Quality gate skip: `!t.NoTest` required to run gate |
| `submit.go` | 128-131 | Coverage auto-set to `-1.0` when `noTest=true` |
| `submit.go` | 406-408 | `formatTestsExecuted()` returns `"No (noTest task)"` |
| `claim.go` | 128 | Copies `NoTest` into `TaskState` |
| `build.go` | 115, 219 | Parses `fm.NoTest` from frontmatter into `Task` struct |

### Producers (auto-generated tasks only)

| File | Line | Task | Why |
|------|------|------|-----|
| `testgen.go` | 43 | T-test-1 (gen-test-cases) | Generates test artifacts, not runnable code |
| `testgen.go` | 49 | T-test-1b (eval-test-cases) | Evaluates test cases |
| `testgen.go` | 111, 615 | T-eval-doc | Docs-only consolidation |
| `testgen.go` | 134, 165 | Per-profile gen-test-cases | Test artifact generation |
| `testgen.go` | 198 | T-quick-6 (drift detection) | Spec analysis task |

### Documentation references

- `guide.md:115` — "Documentation tasks (noTest: true in task frontmatter) skip the quality gate entirely"
- `guide.md:135` — "forge quality-gate automatically skips docs-only features (all tasks have noTest: true)"
- `submit-task/SKILL.md:129` — Documents quality gate condition includes `noTest` check
- `eval-forge/rubric.md:169` — Lists "noTest skips" as known bypass context

## Conclusion

`noTest` is a legitimate field with clear producers (testgen.go auto-generated tasks) and consumers (submit.go quality gate). The D2 audit's "frontmatter injection" concern is technically valid — an agent CAN write `noTest: true` into a hand-written task's frontmatter — but this matches the documented design: task templates use frontmatter to declare properties, and the CLI trusts the frontmatter.

The "deprecation" that was recalled was the `--no-test` CLI flag, not the per-task field.

## Unresolved

Whether `noTest: true` on a task with `type: implementation` or `type: fix` should be rejected by the CLI (currently accepted without warning). This is a design choice, not a bug — the field serves docs-only and test-generation tasks that legitimately skip quality gates.
