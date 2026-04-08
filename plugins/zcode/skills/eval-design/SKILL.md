---
name: eval-design
description: Evaluate a design.md document against quality standards. Checks structure completeness, architecture clarity, interface/model concreteness, error handling, testing strategy, and breakdown-readiness. Outputs a scored report with actionable improvements.
---

# Eval Design

评估 design.md 是否满足规范，重点检查能否直接驱动 `/breakdown-tasks`。

## When to Use

**Trigger:**
- User asks to "evaluate design" or "check design quality"
- User provides `/eval-design` command
- Before handing off design.md to `/breakdown-tasks`

**Skip:**
- design.md doesn't exist yet (use `/design-tech` first)

## Workflow

```
1. 定位 design.md → 2. 检查结构 → 3. 检查内容质量 → 4. 生成报告
```

## Step 1: Locate design.md

Check in order:
1. Path provided by user
2. `docs/features/<current-feature>/design.md`
3. Ask user for path if not found

Also check if a PRD exists at `docs/features/<slug>/prd.md` — used for traceability checks.

## Step 2: Check Structure Completeness

Required sections — mark missing ones as F immediately:

| Section               | Required | Notes                                      |
| --------------------- | -------- | ------------------------------------------ |
| Overview              | ✓        | High-level approach + tech stack           |
| Architecture          | ✓        | Layer placement + component diagram        |
| Interfaces            | ✓        | At least one interface with method sigs    |
| Data Models           | ✓        | Concrete struct/type definitions           |
| Error Handling        | ✓        | Error types + propagation strategy         |
| Testing Strategy      | ✓        | Per-layer plan + coverage target           |
| Security Considerations | ○      | Required if PRD has auth/data requirements |
| Open Questions        | ○        | Optional but recommended                   |
| Alternatives Considered | ○      | Optional but recommended                   |

## Step 3: Check Content Quality

### Dimension 1: Architecture Clarity (架构清晰度)

| Check | Criteria |
|-------|----------|
| Layer placement | Explicitly states which layer(s) this feature belongs to |
| Component diagram | ASCII or text diagram showing components and data flow |
| Dependencies | Lists internal modules and external packages used |
| Consistency | Architecture matches project's existing patterns (check ARCHITECTURE.md if present) |

**Grading:**
- A: Layer placement explicit, diagram present, dependencies listed, consistent with project
- B: Diagram present, minor gaps in dependencies or layer description
- C: Prose description only, no diagram, or missing layer placement
- F: No architecture section

### Dimension 2: Interface & Model Definitions (接口与模型定义)

| Check | Criteria |
|-------|----------|
| Interface signatures | Methods have typed parameters and return values (not just names) |
| Model fields | Structs have field names, types, and constraints (not just descriptions) |
| Completeness | All major components have interfaces or models defined |
| Implementable | A developer can write code directly from these definitions without guessing |

**Grading:**
- A: All interfaces typed, all models concrete, directly implementable
- B: Most defined, 1-2 missing types or constraints
- C: Interfaces/models described in prose, not as code definitions
- F: No interface or model definitions

### Dimension 3: Error Handling (错误处理)

| Check | Criteria |
|-------|----------|
| Error types | Custom error types or error codes defined |
| Propagation | Clear strategy for how errors flow between layers |
| HTTP mapping | If API: HTTP status codes mapped to error types |
| Client behavior | What callers should do on each error |

**Grading:**
- A: Error types defined, propagation strategy clear, HTTP codes mapped
- B: Error types defined, propagation implicit
- C: Only mentions "handle errors" without specifics
- F: No error handling section

### Dimension 4: Testing Strategy (测试策略)

| Check | Criteria |
|-------|----------|
| Per-layer plan | Each layer (service, API, CLI, frontend, etc.) has a test approach |
| Test types | Specifies unit vs integration vs e2e per layer |
| Coverage target | Numeric coverage target stated |
| Test tooling | Testing libraries/frameworks named |

**Grading:**
- A: Per-layer plan, test types specified, coverage target, tooling named
- B: Per-layer plan, coverage target missing or no tooling
- C: Generic "write tests" without layer breakdown
- F: No testing strategy section

### Dimension 5: Breakdown-Readiness (可拆解性)

This is the most critical dimension — design.md is the direct input to `/breakdown-tasks`.

| Check | Criteria |
|-------|----------|
| Enumerable components | Components/modules can be listed and counted |
| Interface tasks derivable | Each interface → at least one implementation task |
| Model tasks derivable | Each data model → at least one schema/migration task |
| No ambiguous ownership | Each component has a clear boundary (not "shared logic") |
| PRD traceability | If PRD exists: all acceptance criteria are addressed in design |

**Grading:**
- A: All components enumerable, tasks clearly derivable, PRD fully covered
- B: Most components clear, 1-2 ambiguous areas, PRD mostly covered
- C: Components described but not enumerable, or significant PRD gaps
- F: Design is too high-level to derive tasks from

### Dimension 6: Security Considerations (安全考量)

*Only required if PRD has auth, data privacy, or multi-user requirements.*

| Check | Criteria |
|-------|----------|
| Threat model | Identifies what could go wrong (injection, unauthorized access, etc.) |
| Mitigations | Concrete countermeasures for each threat |
| Scope-appropriate | Depth matches the feature's actual risk surface |

**Grading:**
- A: Threats identified, mitigations concrete
- B: Threats identified, mitigations vague
- C: Section exists but only says "will add auth later"
- N/A: Feature has no security surface (mark as N/A, not F)

## Step 4: Generate Report

### Grading Rules

**Overall:**

| Grade | Condition |
|-------|-----------|
| A | All required dimensions A/B, at least 3 A's, Breakdown-Readiness ≥ B |
| B | No F on required dimensions, Breakdown-Readiness ≥ B |
| C | 1 F on non-critical dimension, or Breakdown-Readiness = C |
| D | Breakdown-Readiness = F, or 2 F's on required dimensions |
| F | 3+ F's, or Interfaces + Models both F |

> Breakdown-Readiness is weighted higher because it's the direct gate to `/breakdown-tasks`.

Save report to `docs/features/<feature-slug>/design-eval.md` using `templates/report.md`.

## Related

- `/design-tech` — Create or revise the design.md
- `/eval-prd` — Evaluate PRD before design starts
- `/breakdown-tasks` — Next step after design passes evaluation
