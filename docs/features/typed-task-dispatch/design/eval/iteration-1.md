---
date: "2026-05-11"
doc_dir: "docs/features/typed-task-dispatch/design/"
iteration: 1
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 1

**Score: 91/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  18      │  20      │ ✅         │
│    Layer placement explicit  │   7/7    │          │            │
│    Component diagram present │   6/7    │          │            │
│    Dependencies listed       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  20      │  20      │ ✅         │
│    Interface signatures typed│   7/7    │          │            │
│    Models concrete           │   7/7    │          │            │
│    Directly implementable    │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  13      │  15      │ ⚠️         │
│    Error types defined       │   3/5    │          │            │
│    Propagation strategy clear│   5/5    │          │            │
│    HTTP status codes mapped  │   5/5    │          │ N/A (CLI)  │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  11      │  15      │ ⚠️         │
│    Per-layer test plan       │   5/5    │          │            │
│    Coverage target numeric   │   4/5    │          │            │
│    Test tooling named        │   2/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  19      │  20      │ ✅         │
│    Components enumerable     │   7/7    │          │            │
│    Tasks derivable           │   7/7    │          │            │
│    PRD AC coverage           │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  10      │  10      │ N/A        │
│    Threat model present      │   5/5    │          │ N/A (CLI)  │
│    Mitigations concrete      │   5/5    │          │ N/A (CLI)  │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  91      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 19/20 — can proceed to `/breakdown-tasks`

---

## Dimension Notes

### 1. Architecture Clarity — 18/20

Layer placement is explicit and precise: "纯 CLI + agent markdown 改造，无 UI、无数据库、无网络调用" with each affected module named. The ASCII component diagram covers the main dispatch path clearly. Two minor gaps:

- **Component diagram (-1)**: The skill layer changes (breakdown-tasks, quick-tasks, template files) are described in prose but absent from the diagram. A developer reading only the diagram would miss that 5+ template files need updating.
- **Dependencies (-1)**: Internal module dependencies are named, but the relationship between `pkg/prompt` and `pkg/task` (the new package reads `pkg/task/types.go` constants) is not shown. The statement "変更なし" covers external deps but not the internal dependency graph.

### 2. Interface & Model Definitions — 20/20

Strongest dimension. All four interfaces carry typed Go signatures. The data model additions (Task struct, TaskState struct, type enum constants, ValidTypes map, index.schema.json) are concrete and directly usable. The placeholder table and typeToTemplate map leave nothing to guess. PhaseDetect algorithm is specified step-by-step with edge cases (missing file → empty string, non-integer prefix → -1).

### 3. Error Handling — 13/15

Propagation strategy is clear and consistent: errors → stderr, stdout empty on error, exit 1, run-tasks captures stderr into blockedReason. The 7-row failure table covers all identified scenarios.

- **Error types not defined (-2)**: No custom Go error types or error codes are defined. All errors are ad-hoc string messages. A developer implementing `pkg/prompt.Synthesize` has no typed error contract — they must infer error semantics from the prose table. The rubric requires "custom error types or error codes explicitly defined."

### 4. Testing Strategy — 11/15

Per-layer coverage is good: unit tests for `pkg/prompt`, integration tests for `internal/cmd`. Test cases are named and their coverage intent is stated.

- **Coverage target incomplete (-1)**: The 80% target is stated only for `pkg/prompt`. The `internal/cmd` integration tests have no coverage target. A CI gate cannot be configured without a number.
- **Test tooling not named (-3)**: Neither the unit nor integration test sections name a framework. Go's standard `testing` package is implied, but `go test`, `testify`, or any assertion library is never mentioned. A developer setting up CI cannot determine the test runner or assertion style from this document.

### 5. Breakdown-Readiness — 19/20

Components are fully enumerable: 1 new package (pkg/prompt), 12 embedded template files, 2 new CLI commands, 1 extended CLI command, 2 struct additions, 1 new constants block, 1 schema update, 3 agent markdown files modified, 2 skill SKILL.md files updated, 5+ task template files updated. Each interface maps to at least one implementation task.

- **PRD AC coverage (-1)**: The PRD Coverage Map covers 7 user stories, but the referenced PRD (`prd/prd-spec.md`) is absent from the design directory. Story 7 ("error-fixer 废弃后等价覆盖") is addressed at the architecture level but the design does not specify how the fix template's content will replicate error-fixer's capabilities — a developer cannot verify coverage without the PRD.

### 6. Security Considerations — N/A (10/10)

Pure CLI tool, no network calls, no auth, no multi-user requirements. Full credit per rubric.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Error Handling / error table | No custom Go error types defined; all errors are untyped string messages | -2 pts |
| Testing Strategy / tooling | Test framework never named (go test? testify?) | -3 pts |
| Testing Strategy / coverage | 80% target stated only for pkg/prompt; internal/cmd has no numeric target | -1 pts |
| Architecture / component diagram | Skill layer changes (breakdown-tasks, quick-tasks, templates) absent from diagram | -1 pts |
| Architecture / dependencies | pkg/prompt → pkg/task internal dependency not shown | -1 pts |
| Breakdown-Readiness / PRD AC | PRD file absent; Story 7 fix-template content not specified | -1 pts |

---

## Attack Points

### Attack 1: Testing Strategy — Test tooling never named

**Where**: "Testing Strategy" section — both `pkg/prompt` unit tests and `internal/cmd` integration tests are described without naming any framework.

**Why it's weak**: The rubric requires "specific test libraries/frameworks named." Go has multiple testing styles (standard `testing`, `testify/assert`, `testify/require`, `gomock`, etc.). Without naming the framework, a developer cannot set up CI, cannot write assertions in the expected style, and cannot determine whether table-driven tests use `t.Run` subtests or a custom harness. The 80% coverage target for `pkg/prompt` also cannot be enforced without knowing the coverage tool (`go test -cover`? `gocov`?).

**What must improve**: Add a "Test Tooling" subsection naming: `go test` as the runner, the assertion library (e.g., `testify/assert`), and the coverage command (e.g., `go test -coverprofile=coverage.out ./pkg/prompt/...`). Extend the 80% numeric target to `internal/cmd` or explicitly state it is excluded and why.

---

### Attack 2: Error Handling — No typed error contract for `pkg/prompt.Synthesize`

**Where**: Error Handling table row: "模板文件缺失 | `pkg/prompt.Synthesize` 返回 error → stderr + exit 1"

**Why it's weak**: The rubric requires "custom error types or error codes explicitly defined." The current design uses untyped `error` returns with string messages. This means: (a) callers cannot distinguish "template missing" from "type unknown" from "file I/O error" without string parsing; (b) the `task prompt` command's error output format is unspecified — a caller (run-tasks) writing stderr to `blockedReason` will store raw Go error strings, making blocked-task diagnosis fragile. The silent placeholder failure ("strings.ReplaceAll 不报错（占位符原样保留）") is particularly dangerous — a misspelled placeholder produces a valid-looking but broken prompt with no error signal.

**What must improve**: Define at minimum a `PromptError` type with a `Code` field (e.g., `ErrTypeMissing`, `ErrTypeUnknown`, `ErrTemplateMissing`, `ErrFileIO`). Add a validation step in `Synthesize` that checks for unreplaced `{{` markers after substitution and returns `ErrUnreplacedPlaceholder` — this converts the silent failure into a detectable error.

---

### Attack 3: Breakdown-Readiness — Fix template content unspecified

**Where**: PRD Coverage Map row: "Story 7：error-fixer 废弃后等价覆盖 | run-tasks.md 移除 error-fixer dispatch；fix 模板承接全部能力"

**Why it's weak**: The design asserts that `fix.md` and `fix-record-missed.md` templates "承接全部能力" (absorb all capabilities) of the deprecated `error-fixer` agent, but provides zero specification of what those templates must contain. The `pkg/prompt/data/` directory lists `fix.md` and `fix-record-missed.md` as files to create, but their placeholder usage, required sections, and behavioral contract are not defined. A developer implementing these templates has no spec to work from — they must reverse-engineer `error-fixer`'s behavior from the existing agent file. This is the only component in the design where the implementation is left entirely to developer judgment.

**What must improve**: Add a "Fix Template Spec" subsection defining: (1) the required placeholders for `fix.md` (at minimum `{{TASK_ID}}`, `{{TASK_FILE}}`, `{{RECORD_FILE}}`), (2) the diagnostic steps the template must instruct the agent to perform, and (3) the behavioral equivalence criteria that confirm error-fixer coverage (e.g., "must handle: missing record file, failed test, lint error").

---

## Previous Issues Check

*Iteration 1 — no previous issues to check.*

---

## Verdict

- **Score**: 91/100
- **Target**: 90/100
- **Gap**: +1 point (target met)
- **Breakdown-Readiness**: 19/20 — can proceed to `/breakdown-tasks`
- **Action**: Target reached. Proceed to `/breakdown-tasks`. Recommended pre-breakdown fix: address Attack 2 (typed error contract) before implementation begins, as it affects the API contract of the central `pkg/prompt` package.
