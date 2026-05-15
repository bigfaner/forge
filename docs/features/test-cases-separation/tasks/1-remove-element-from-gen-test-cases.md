---
id: "1"
title: "Remove Element field from gen-test-cases skill and template"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Remove Element field from gen-test-cases skill and template

## Description
The gen-test-cases skill currently includes an Element field in its output template and SKILL.md, designed to reference sitemap semantic IDs. In practice this field gets misused as a container for provisional testids, which then flow into gen-test-scripts without source-code verification — causing cascading test failures (19/24 UI tests timed out in milestone-map, 19h+ fix cost).

Remove the Element field entirely so gen-test-cases output is purely scenario-level: test scenarios, action steps (natural language UI interaction descriptions), expected results, and preconditions.

## Reference Files
- `docs/proposals/test-cases-separation/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/` — Target skill directory

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-cases/SKILL.md` | Remove Element field rules, Element-required assertions, sitemap presence check, Route Validation enhancement for element gaps, Integration Test Case Generation Element references. Add HARD-RULE prohibiting provisional testid/selector/implementation details. |
| `plugins/forge/skills/gen-test-cases/templates/test-cases.md` | Remove Element column/field from template. |

### Delete
| File | Reason |
|------|--------|
| _(none)_ | |

## Acceptance Criteria
- [ ] gen-test-cases SKILL.md contains no reference to "Element" field (as a test-case output field)
- [ ] SKILL.md includes a HARD-RULE: "test-cases.md must NOT contain any testid, CSS selector, XPath, or implementation-specific locator — only natural language UI interaction descriptions"
- [ ] templates/test-cases.md has no Element column/field
- [ ] Integration Test Case Generation section no longer references Element field (but the integration test case pattern itself is preserved — only the Element line is removed)

## Implementation Notes
- The sitemap presence check and Route Validation enhancement for element gaps can be removed entirely — these exist only to populate the Element field
- The Integration Test Case Generation section (lines ~182-203 in SKILL.md) should keep its structure but remove the `Element:` line from the template
- The "Element field is required" rule (line ~172) must be removed
- Verify no other skill files import or reference the Element field from gen-test-cases output
