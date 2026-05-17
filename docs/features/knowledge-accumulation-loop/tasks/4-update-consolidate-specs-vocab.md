---
id: "4"
title: "Update /consolidate-specs — add auto-vocabulary generation"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "enhancement"
mainSession: false
---

# 4: Update /consolidate-specs — add auto-vocabulary generation

## Description

Add a vocabulary generation step to `/consolidate-specs` that automatically归纳 domain and type vocabularies from existing knowledge files during its drift-detection pass. The generated vocabulary is stored as a derived artifact that `/learn` and auto-extract triggers can read at runtime for classification suggestions.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal (Part 3: Auto-Generated Vocabulary)
- `plugins/forge/skills/consolidate-specs/SKILL.md` — Existing consolidate-specs skill

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Add vocabulary generation step after drift detection |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] New step added to consolidate-specs workflow (after Step 11 drift fix, before Step 12 record task)
- [ ] Vocabulary generation scans all 4 knowledge directories:
  - `docs/decisions/*.md` — extract type names and domains from decision rows
  - `docs/lessons/*.md` — extract tags from frontmatter
  - `docs/conventions/*.md` — extract domains from frontmatter
  - `docs/business-rules/*.md` — extract domains from frontmatter
- [ ] Generates a vocabulary index with:
  - `types`: unique knowledge types found (decision, lesson, convention, business-rule)
  - `domains`: unique domain keywords aggregated from all files
  - `counts`: how many entries per type/domain for context
- [ ] Vocabulary output is a prompt-readable format (markdown or yaml)
- [ ] Step is idempotent — regenerates on every run, replacing previous vocabulary
- [ ] `/learn` skill and auto-extract triggers can reference this vocabulary for classification suggestions
- [ ] Vocabulary generation works even when knowledge directories are sparse or empty (produces base vocabulary)

## Hard Rules
- Do not change existing consolidate-specs steps — only add the new step
- Vocabulary must include the base 8-category vocabulary (architecture, interface, data-model, dependencies, error-handling, testing, security, local-dev-deployment) even when no knowledge files exist
- The generated vocabulary must be clearly marked as auto-generated (not user-editable)

## Implementation Notes
- The vocabulary generation step should run after drift detection (Step 10-11) so it reflects the latest state
- Consider storing the vocabulary in a predictable location that both `/learn` and triggers can find, e.g., referenced directly in the consolidate-specs SKILL.md as a section, or written to a file in `.forge/` or `docs/`
- Since this is prompt-level, the "vocabulary" is really a set of instructions for the agent to follow when generating suggestions — not a code data structure
