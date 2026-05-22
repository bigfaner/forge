---
name: auto-task-main-instructions
status: Approved
created: 2026-05-22
updated: 2026-05-22
---

# Proposal: Auto-Generated Task Body Content

## Problem

CLI auto-generated tasks (T-test-*, T-validate-*, T-clean-code, T-specs-*, T-eval-doc) produced by `GenerateTestTaskMD()` have useless body content. The prompt templates (`forge prompt get-by-task-id`) already provide complete workflow guidance, but their Step 1 instructs agents to "Read the task file to understand X." The task file then says either "Execute this test pipeline task." or "Read docs/conventions/testing-*.md" — providing zero value.

Three categories of task types have different needs:

1. **Skill-delegating tasks** (gen-cases, gen-scripts, run, graduate, clean-code, consolidate, drift): The invoked skill handles everything. The task file body adds no value beyond a one-line purpose statement.
2. **Criteria-driven tasks** (validation-code, validation-ux): The prompt template expects "validation criteria listed in the task file." Without criteria, the agent has nothing to validate against.
3. **Target-driven tasks** (doc-eval): The prompt template expects "the list of documents to evaluate." Without a document list, the agent doesn't know what to evaluate.

## Analysis: Why Embed Templates Duplicated Prompt Templates

A previous implementation created 13 embed template files (`forge-cli/pkg/task/data/*.md`) containing the same workflow as prompt templates. This was pure duplication — the agent received identical instructions twice, wasting context window.

The root cause: trying to inject **static workflow content** that prompt templates already provide, instead of injecting **dynamic context** that only the task file can carry.

### Information Availability Matrix

| Information | Known at generation time (BuildIndex) | Known at execution time | Already handled by |
|-------------|:---:|:---:|---|
| Feature slug, scope, success criteria | Yes | — | — |
| PRD acceptance criteria | Yes (read from PRD/proposal) | — | — |
| Test interfaces, project type | Yes (read from config) | — | — |
| Document list for doc-eval | Partial (may miss new docs) | Yes (scan directory) | — |
| Changed files for clean-code | No | Yes (git diff) | Skill handles at runtime |
| Convention summaries | No (domain matching needed) | Yes (read + match domains) | Prompt template Step 1 |
| Workflow steps (TDD, fix, verify) | — | — | Prompt template already provides |

### Three-Tier Content Strategy

| Tier | When to inject | What | Example |
|------|---------------|------|---------|
| **Static context** | Generation time (BuildIndex) | Feature slug, purpose, scope, success criteria, PRD AC | "Feature: auth-refresh. Scope: add token rotation. Success: tokens refresh without re-login." |
| **Discovery pointers** | Generation time | Where to find runtime information | "Scan `docs/features/<slug>/` for all documents." |
| **Nothing** | — | Information already handled by prompt template or skill | Clean-code: skill does git diff. Conventions: prompt template tells agent to read them. |

## Solution

Replace the 13 embed template files with a single `generateBody()` function in `GenerateTestTaskMD()` that produces different body content per task type using a `BodyContext` struct.

### BodyContext

Passed from `BuildIndex()` to `GenerateTestTaskMD()`, carrying information available at generation time:

```go
type BodyContext struct {
    FeatureSlug    string
    Mode           string        // "quick" or "breakdown"
    ProjectRoot    string
    Scope          []string      // In-scope items from proposal/PRD
    SuccessCriteria []string     // Success criteria from proposal/PRD
    AcceptanceCriteria []string  // PRD acceptance criteria (breakdown mode only)
    ProjectType    string        // From .forge/config.yaml
    Interfaces     []string      // Test interfaces from config
}
```

### Body Generation Per Task Type

**Skill-delegating tasks** (gen-cases, gen-scripts, run, graduate, clean-code, consolidate, drift):
- One-line purpose statement
- Feature slug and mode
- Discovery pointers for skill inputs

Example (test.gen-cases):
```markdown
Generate structured test cases for the {{SLUG}} feature ({{MODE}} mode).
Test interfaces: api, cli.
Read the PRD/proposal acceptance criteria and generate test cases with full traceability.
```

**Criteria-driven tasks** (validation-code, validation-ux):
- Purpose statement
- **PRD Acceptance Criteria as validation checklist** (read from PRD/proposal at generation time)
- Discovery pointer for conventions

Example (validation.code):
```markdown
Validate code quality for the {{SLUG}} feature.

## Validation Criteria
- [ ] Token rotation works without re-login
- [ ] Expired tokens are refreshed transparently
- [ ] Concurrent requests handle token refresh correctly
```

**Target-driven tasks** (doc-eval):
- Purpose statement
- Discovery pointers to feature directories (not hardcoded file list — agent scans at runtime to discover new documents)

Example (doc.eval):
```markdown
Evaluate documentation quality for the {{SLUG}} feature.
Scan these directories for documents to evaluate:
- docs/features/{{SLUG}}/prd/
- docs/features/{{SLUG}}/design/
- docs/features/{{SLUG}}/testing/
- docs/proposals/{{SLUG}}/
```

### Implementation Steps

1. **Delete** the 13 embed template files in `forge-cli/pkg/task/data/`
2. **Add** `BodyContext` struct to `autogen.go`
3. **Add** `generateBody(def AutoGenTaskDef, ctx BodyContext) string` function
4. **Update** `GenerateTestTaskMD()` signature: `(def AutoGenTaskDef, ctx BodyContext) ([]byte, error)`
5. **Update** `BuildIndex()` to populate `BodyContext` from PRD/proposal/config and pass it through
6. **Remove** `autogenTemplateFS` and `autogenTypeToFile` map (no longer needed)

### Changes to BuildIndex()

Before calling `GenerateTestTaskMD()`, `BuildIndex()` reads:
- Proposal/PRD → extracts scope, success criteria, acceptance criteria
- `.forge/config.yaml` → extracts project type, interfaces
- Constructs `BodyContext` and passes to `GenerateTestTaskMD()`

## Alternatives Considered

1. **Embed template files with full workflow** — Duplicates prompt templates. Agent reads same content twice. Wastes context. **Rejected.**
2. **Static per-type instructions (Go map)** — Same duplication problem. Content is workflow, not feature context. **Rejected.**
3. **Enhance prompt templates only** — Prompt templates can't carry feature-specific data (PRD criteria, document paths). The task file is the only per-instance artifact. **Rejected.**
4. **Pre-populate all runtime info** — doc-eval would miss new documents, clean-code can't know changed files, convention matching is complex. Most runtime info should stay dynamic. **Rejected.**

## Scope

### In Scope

- Delete `forge-cli/pkg/task/data/*.md` (13 files) and `autogenTemplateFS`/`autogenTypeToFile`
- Add `BodyContext` struct
- Add `generateBody()` with type-specific body generation
- Update `GenerateTestTaskMD()` to use `generateBody()`
- Update `BuildIndex()` to populate `BodyContext`
- Update tests

### Out of Scope

- Prompt templates (`forge-cli/pkg/prompt/data/`) — unchanged
- fix-task.md / cleanup-task.md templates — separate scope
- buildTaskMarkdown() — separate scope
- Skill-level templates — managed by skills

## Risks

| Risk | Mitigation |
|------|-----------|
| PRD/proposal not found at generation time | BodyContext fields default to empty; generateBody() omits sections with empty data |
| doc-eval document list outdated after tasks run | Use discovery pointers (scan directories) instead of hardcoded lists |
| Acceptance criteria extraction is fragile | Use simple section parsing ("## Acceptance Criteria" / "## Success Criteria" headers); fall back to empty on parse failure |
| BodyContext adds complexity to BuildIndex() | BodyContext is a plain struct with clear lifecycle; population is sequential reads |

## Success Criteria

- [ ] No embed template files remain in `forge-cli/pkg/task/data/`
- [ ] Skill-delegating tasks have concise purpose + discovery pointers in body
- [ ] Validation tasks have PRD acceptance criteria as validation checklist
- [ ] Doc-eval tasks have discovery pointers to feature directories
- [ ] No workflow duplication between task body and prompt templates
- [ ] Existing tests pass (backward compatible frontmatter)
- [ ] New tests verify body content per task type
