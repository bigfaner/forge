---
scale: 100
target: 70
iterations: 1
type: harness
context:
  conventions: []
  business-rules: []
---

# Harness Evaluation Rubric

**Total: 100 points**

## What This Evaluates

The project's "harness" — the scaffolding that makes agents productive: instruction files, documentation structure, architectural constraints, tooling, and feedback mechanisms. Based on [OpenAI's Harness Engineering](https://openai.com/index/harness-engineering/) practices.

## Input

The scorer receives a **harness snapshot** (a single markdown file) containing: entry point content, configuration, documentation structure, scripts, CI setup, skills/agents list, and test infrastructure.

## Dimensions

### 1. Progressive Disclosure (25 pts)

> "Give agents a map, not a 1,000-page instruction manual."

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Entry point concise | 0-8 | Is the main instruction file (CLAUDE.md / AGENTS.md) brief (~100 lines) and structured as a table of contents with pointers? Or is it a monolithic wall of rules? |
| Knowledge base structured | 0-9 | Is there a structured docs/ directory with indexed, discoverable content? Can agents find deeper documentation without reading everything? Is there an index or catalog file? |
| Doc validation mechanized | 0-8 | Are there automated checks for documentation freshness, cross-links, or structure? (CI jobs, lint scripts, doc-gardening agents, freshness detection scripts) |

### 2. Architectural Boundaries (25 pts)

> "Enforce invariants, not micromanage implementations."

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Dependency direction defined | 0-9 | Is there a stated architectural model with defined dependency directions? (layers, module boundaries, dependency flow). Documented as prose, diagrams, or code. |
| Boundaries mechanically enforced | 0-8 | Are architectural constraints enforced by tooling (linters, structural tests, CI checks) — not just written in docs? Can a developer accidentally violate a boundary without getting caught? |
| Error messages guide remediation | 0-8 | Do linters, tests, and CI errors include actionable fix instructions? Not just "error on line X" but "import from layer Y is forbidden, use Z instead". |

### 3. Golden Principles (25 pts)

> "Opinionated, mechanical rules that keep the codebase legible for agents."

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Shared tools over ad-hoc | 0-9 | Does the project centralize common patterns into reusable tools (shared utilities, skills, agents, templates)? Or does each task reinvent the wheel? |
| Boundary validation | 0-8 | Are data boundaries validated (schemas, types, interfaces) rather than probed raw? Is "YOLO-style" access discouraged by tooling or convention? |
| Tech debt GC process | 0-8 | Is there a recurring mechanism to detect and clean up drift? (quality scans, stale doc detection, duplicate code detection, code simplification hooks, background cleanup tasks) |

### 4. Plan Artifacts & Feedback (25 pts)

> "Plans as first-class artifacts; feedback loops for continuous correction."

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Plans are versioned artifacts | 0-9 | Are execution plans, task definitions, and decisions versioned in the repo? Can agents resume or reproduce work without external context (Slack threads, verbal instructions)? |
| Execution records structured | 0-8 | Are there structured records of what was done, why, and what was decided? (task records with schema validation, decision logs, change history) |
| Test results structured & observable | 0-8 | Are test results parseable and structured? Can errors be classified, coverage be queried, and task progress be tracked programmatically? |

## Deduction Rules

- **No instruction file (CLAUDE.md / AGENTS.md)**: 0 pts for dimension 1
- **No docs/ directory**: -10 pts from dimension 1
- **No mechanical enforcement of any kind**: 0 pts for "mechanically enforced" criterion
- **All knowledge lives only in external systems** (wiki, Slack, heads): -5 pts from dimension 4
- **Placeholder text ("TBD", "TODO", "fill in later")**: -2 pts per instance

## Scoring Guidance

- **90-100**: Exemplary harness. Agents can operate autonomously with high confidence.
- **70-89**: Solid foundation. Some gaps that would slow agents down on specific tasks.
- **50-69**: Functional but fragile. Agents can work but will struggle with edge cases and drift.
- **30-49**: Underspecified. Agents need frequent human intervention.
- **0-29**: Hostile environment. Agents cannot be productive without major harness investment.
