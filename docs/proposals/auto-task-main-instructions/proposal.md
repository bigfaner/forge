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

1. **Skill-delegating tasks** (gen-cases, gen-scripts, run, graduate, clean-code, consolidate, drift): The invoked skill handles everything. The task file body needs concise purpose + discovery pointers.
2. **Criteria-driven tasks** (validation-code, validation-ux): The prompt template expects "validation criteria listed in the task file." Without criteria, the agent has nothing to validate against.
3. **Target-driven tasks** (doc-eval): The prompt template expects "the list of documents to evaluate." Without discovery pointers, the agent doesn't know where to look.

## Solution

Use **embed template files with placeholders**, dynamically filled from a `BodyContext` struct. Templates contain feature-context structure (NOT workflow — that's already in prompt templates), and placeholders are substituted at generation time with data from PRD/proposal/config.

### Why Not Previous Approaches

| Approach | Why Rejected |
|----------|-------------|
| Embed templates with full workflow (v1) | Duplicates prompt templates. Agent reads same instructions twice. Wastes context. |
| `generateBody()` Go string building | Content buried in code. Not reviewable. Mixing content and logic. |
| Prompt template enhancement only | Prompt templates can't carry per-instance data (PRD criteria, feature slug). |

### Template + Placeholder Pattern

Templates live in `forge-cli/pkg/task/data/` as embed files. Each template contains:
- **Purpose statement** — what this task type does
- **Placeholders** — `{{FEATURE_SLUG}}`, `{{SCOPE}}`, `{{ACCEPTANCE_CRITERIA}}`, etc.
- **Discovery pointers** — where to find runtime information (for dynamic data not available at generation time)

NO workflow steps, NO Skill invocation instructions — those are in prompt templates.

### BodyContext

Carries feature-specific data from `BuildIndex()` to template rendering:

```go
type BodyContext struct {
    FeatureSlug        string
    Mode               string        // "quick" or "breakdown"
    Scope              []string      // In-scope items from proposal/PRD
    SuccessCriteria    []string      // Success criteria from proposal/PRD
    AcceptanceCriteria []string      // PRD acceptance criteria (breakdown mode only)
    ProjectType        string        // From .forge/config.yaml
    Interfaces         []string      // Test interfaces from config
}
```

### Template Files Per Task Type

**Skill-delegating tasks** — purpose + feature context:

`test-gen-cases.md`:
```markdown
Generate structured test cases for the {{FEATURE_SLUG}} feature ({{MODE}} mode).
Test interfaces: {{INTERFACES}}.
Read the PRD/proposal acceptance criteria and generate test cases with full traceability.
```

`test-gen-scripts.md`:
```markdown
Generate executable test scripts from approved test cases for the {{FEATURE_SLUG}} feature.
Test type: {{TEST_TYPE}}.
Use the framework specified by the active test profile.
```

`test-run.md`:
```markdown
Execute staged e2e test scripts for the {{FEATURE_SLUG}} feature.
If tests fail, identify failing tests and root cause, apply minimal fix, then re-run.
```

`test-graduate.md`:
```markdown
Promote feature test scripts to the project regression suite for the {{FEATURE_SLUG}} feature.
Read scripts from the staging directory and migrate to the permanent test suite.
```

`test-verify-regression.md`:
```markdown
Run full e2e regression suite to verify no regressions for the {{FEATURE_SLUG}} feature.
Use `just test-e2e` for regression verification.
```

`test-gen-and-run.md`:
```markdown
Generate and run test scripts for the {{FEATURE_SLUG}} feature.
Phase 1: Generate scripts. Phase 2: Execute and verify.
Test type: {{TEST_TYPE}}.
```

`test-eval-cases.md`:
```markdown
Evaluate generated test cases for executability for the {{FEATURE_SLUG}} feature.
Verify each test case has clear steps, expected results, and can drive script generation.
```

`code-quality-simplify.md`:
```markdown
Simplify and clean up code for the {{FEATURE_SLUG}} feature.
The skill resolves scope automatically (git diff > feature context).
```

`doc-consolidate.md`:
```markdown
Extract and consolidate business rules and tech specs for the {{FEATURE_SLUG}} feature.
Run in non-interactive mode: auto-integrate all CROSS items.
```

`doc-drift.md`:
```markdown
Detect spec drift for the {{FEATURE_SLUG}} feature.
Compare existing specs in docs/business-rules/ and docs/conventions/ against current code.
Auto-fix drifted specs and commit.
```

**Criteria-driven tasks** — purpose + PRD acceptance criteria:

`validation-code.md`:
```markdown
Validate code quality for the {{FEATURE_SLUG}} feature.

## Validation Criteria
{{ACCEPTANCE_CRITERIA}}

Check docs/conventions/ for project-specific quality standards.
Run the quality gate: compile, fmt, lint, test.
```

`validation-ux.md`:
```markdown
Validate user experience for the {{FEATURE_SLUG}} feature.

## Validation Criteria
{{ACCEPTANCE_CRITERIA}}

Verify that user-facing behavior matches the expected experience.
Check accessibility, usability, and consistency.
```

**Target-driven tasks** — purpose + discovery pointers:

`doc-eval.md`:
```markdown
Evaluate documentation quality for the {{FEATURE_SLUG}} feature.

Scan these directories for documents to evaluate:
{{DOCUMENT_DIRS}}

Score each document using the 8-dimension rubric (1000-point scale).
```

### Placeholder Substitution

`BodyContext` fields map to template placeholders:

| Placeholder | Source | When Empty |
|-------------|--------|-----------|
| `{{FEATURE_SLUG}}` | BodyContext.FeatureSlug | Never empty (required) |
| `{{MODE}}` | BodyContext.Mode | Omit mode line |
| `{{INTERFACES}}` | BodyContext.Interfaces | "See .forge/config.yaml" |
| `{{TEST_TYPE}}` | AutoGenTaskDef.TestType | Omit type line |
| `{{ACCEPTANCE_CRITERIA}}` | BodyContext.AcceptanceCriteria | "- [ ] All acceptance criteria met" |
| `{{DOCUMENT_DIRS}}` | Derived from FeatureSlug + Mode | Scan feature directory |
| `{{SCOPE}}` | BodyContext.Scope | Omit scope section |

### Information Availability: What Goes Where

| Information | Inject at generation time? | Why |
|-------------|:---:|------|
| Feature slug, mode | Yes | Static, known from BuildIndex |
| Scope, success criteria | Yes | Read from proposal/PRD |
| PRD acceptance criteria | Yes | Read from PRD, used by validation tasks |
| Test interfaces, project type | Yes | Read from .forge/config.yaml |
| Document directories | Yes | Derived from feature slug |
| Changed files (clean-code) | No | Skill discovers via git diff at runtime |
| Convention summaries | No | Agent reads + matches at runtime (prompt template handles this) |
| Workflow steps | No | Already in prompt templates |

## Implementation Steps

1. **Update** the 13 embed template files in `forge-cli/pkg/task/data/` — replace full workflow content with purpose + placeholders
2. **Add** `BodyContext` struct to `autogen.go`
3. **Add** `renderBody(templateContent string, def AutoGenTaskDef, ctx BodyContext) string` — reads embed template, substitutes placeholders
4. **Update** `GenerateTestTaskMD()` — use `renderBody()` instead of inline body generation
5. **Update** `BuildIndex()` — populate `BodyContext` from PRD/proposal/config, pass to `GenerateTestTaskMD()`
6. **Update** `GenerateTestTaskMD()` signature: `(def AutoGenTaskDef, ctx BodyContext) ([]byte, error)`
7. **Update** tests

### Changes to BuildIndex()

Before calling `GenerateTestTaskMD()`, `BuildIndex()`:
1. Reads proposal (`docs/proposals/<slug>/proposal.md`) or PRD (`docs/features/<slug>/prd/prd-spec.md`)
2. Extracts scope, success criteria, acceptance criteria (simple section parsing)
3. Reads `.forge/config.yaml` for project type and interfaces
4. Constructs `BodyContext` and passes it through

### Changes to GenerateTestTaskMD()

Current signature: `(def AutoGenTaskDef, _ string) ([]byte, error)`
New signature: `(def AutoGenTaskDef, ctx BodyContext) ([]byte, error)`

Body generation:
1. Look up template file via `autogenTypeToFile[def.Type]`
2. Read template from `autogenTemplateFS`
3. Call `renderBody()` with template, def, and ctx
4. Append TestType note if present
5. Append StrategyContent if present (legacy)

## Risks

| Risk | Mitigation |
|------|-----------|
| PRD/proposal not found at generation time | BodyContext fields default to empty; renderBody() omits sections with empty placeholders |
| doc-eval document list outdated after tasks run | Use discovery pointers (directory paths) instead of hardcoded file lists |
| Acceptance criteria extraction is fragile | Simple section header parsing; fall back to empty on failure |
| BodyContext adds complexity to BuildIndex() | Plain struct with sequential reads; clear lifecycle |

## Success Criteria

- [ ] Embed template files contain purpose + placeholders (NOT full workflow)
- [ ] No workflow duplication between task body and prompt templates
- [ ] Validation tasks have PRD acceptance criteria as validation checklist
- [ ] Doc-eval tasks have discovery pointers to feature directories
- [ ] Skill-delegating tasks have concise purpose + feature context
- [ ] Existing tests pass (backward compatible frontmatter)
- [ ] New tests verify body content per task type
