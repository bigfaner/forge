---
name: eval-design
description: Evaluate a design.md document against quality standards. Checks structure completeness, architecture clarity, interface/model concreteness, error handling, testing strategy, and breakdown-readiness. Outputs a scored report with actionable improvements.
---

# Eval Design

评估 tech-design.md 是否满足规范，重点检查能否直接驱动 `/breakdown-tasks`。

## When to Use

**Trigger:**
- User asks to "evaluate design" or "check design quality"
- User provides `/eval-design` command
- Before handing off tech-design.md to `/breakdown-tasks`

**Skip:**
- design.md doesn't exist yet (use `/design-tech` first)

## Workflow

```
1. 定位 design.md → 2. 启动评估 Agent → 3. 汇报结果
```

## Step 1: Locate Design Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate design documents
3. Fall back to `design/tech-design.md`, `design/api-handbook.md`, `ui/ui-design.md`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. Also check if a PRD exists at `prd/prd-spec.md` — used for traceability checks.

## Step 2: Launch Evaluation Agent

Use the **Agent tool** to spawn a subagent. Pass the full prompt below, substituting `{{DESIGN_PATH}}`, `{{PRD_PATH}}`, and `{{FEATURE_SLUG}}`:

---

**Agent prompt template:**

```
You are a technical design quality evaluator. Your job: read the design doc, apply the rubric, write the report, return a summary.

## Inputs
- Design path: {{DESIGN_PATH}} (default: design/tech-design.md)
- API Handbook path: {{API_HANDBOOK_PATH}} (default: design/api-handbook.md)
- UI Design path: {{UI_DESIGN_PATH}} (default: ui/ui-design.md)
- PRD path: {{PRD_PATH}} (default: prd/prd-spec.md, read if it exists)
- Feature slug: {{FEATURE_SLUG}}
- Report output: docs/features/{{FEATURE_SLUG}}/design-eval.md
- Report template: plugins/zcode/skills/eval-design/templates/report.md

## Steps
1. Read {{DESIGN_PATH}}
2. Read {{API_HANDBOOK_PATH}} (if exists)
3. Read {{UI_DESIGN_PATH}} (if exists)
4. If {{PRD_PATH}} exists, read it (needed for traceability checks)
5. Read the report template
6. Apply the rubric below to every dimension
7. Fill in the template and write to docs/features/{{FEATURE_SLUG}}/design-eval.md
8. Return: overall grade, top 2-3 issues, Breakdown-Readiness grade, and whether it can proceed to /breakdown-tasks

## Structure Check

Required sections — mark missing as F:

| Section                 | Required | Notes                                      |
|-------------------------|----------|--------------------------------------------|
| Overview                | ✓        | High-level approach + tech stack           |
| Architecture            | ✓        | Layer placement + component diagram        |
| Interfaces              | ✓        | At least one interface with method sigs    |
| Data Models             | ✓        | Concrete struct/type definitions           |
| Error Handling          | ✓        | Error types + propagation strategy         |
| Testing Strategy        | ✓        | Per-layer plan + coverage target           |
| Security Considerations | ○        | Required if PRD has auth/data requirements |
| Open Questions          | ○        | Optional                                   |
| Alternatives Considered | ○        | Optional                                   |

## Dimension 1: Architecture Clarity

Checks: layer placement (explicitly states which layer), component diagram (ASCII or text), dependencies (internal modules + external packages), consistency with project patterns.

- A: Layer placement explicit, diagram present, dependencies listed, consistent with project
- B: Diagram present, minor gaps in dependencies or layer description
- C: Prose description only, no diagram, or missing layer placement
- F: No architecture section

## Dimension 2: Interface & Model Definitions

Checks: interface signatures (typed params + return values), model fields (names, types, constraints), completeness (all major components defined), implementable (developer can code directly without guessing).

- A: All interfaces typed, all models concrete, directly implementable
- B: Most defined, 1-2 missing types or constraints
- C: Interfaces/models described in prose, not as code definitions
- F: No interface or model definitions

## Dimension 3: Error Handling

Checks: error types (custom types or codes defined), propagation (clear strategy between layers), HTTP mapping (if API: status codes mapped), client behavior (what callers do on each error).

- A: Error types defined, propagation strategy clear, HTTP codes mapped
- B: Error types defined, propagation implicit
- C: Only mentions "handle errors" without specifics
- F: No error handling section

## Dimension 4: Testing Strategy

Checks: per-layer plan (each layer has a test approach), test types (unit/integration/e2e specified per layer), coverage target (numeric), test tooling (libraries named).

- A: Per-layer plan, test types specified, coverage target, tooling named
- B: Per-layer plan, coverage target missing or no tooling
- C: Generic "write tests" without layer breakdown
- F: No testing strategy section

## Dimension 5: Breakdown-Readiness ★ (critical gate)

Checks: enumerable components (can be listed and counted), interface tasks derivable (each interface → at least one impl task), model tasks derivable (each model → at least one schema/migration task), no ambiguous ownership, PRD traceability (if PRD exists: all AC addressed in design).

- A: All components enumerable, tasks clearly derivable, PRD fully covered
- B: Most components clear, 1-2 ambiguous areas, PRD mostly covered
- C: Components described but not enumerable, or significant PRD gaps
- F: Design too high-level to derive tasks from

## Dimension 6: Security Considerations

Only required if PRD has auth, data privacy, or multi-user requirements.

Checks: threat model (identifies what could go wrong), mitigations (concrete countermeasures), scope-appropriate (depth matches actual risk surface).

- A: Threats identified, mitigations concrete
- B: Threats identified, mitigations vague
- C: Section exists but only says "will add auth later"
- N/A: Feature has no security surface (mark N/A, not F)

## Overall Grade

| Grade | Condition                                                        |
|-------|------------------------------------------------------------------|
| A     | All required dimensions A/B, at least 3 A's, Breakdown-Readiness ≥ B |
| B     | No F on required dimensions, Breakdown-Readiness ≥ B            |
| C     | 1 F on non-critical dimension, or Breakdown-Readiness = C        |
| D     | Breakdown-Readiness = F, or 2 F's on required dimensions         |
| F     | 3+ F's, or Interfaces + Models both F                            |

Breakdown-Readiness is weighted higher — it is the direct gate to /breakdown-tasks.
```

---

## Step 3: Report to User

After the agent completes, relay its summary: overall grade, Breakdown-Readiness grade, top issues, and next step recommendation.

## Related

- `/design-tech` — Create or revise the design.md
- `/eval-prd` — Evaluate PRD before design starts
- `/breakdown-tasks` — Next step after design passes evaluation
