---
id: "3"
title: "Remove Element evaluation from eval-test-cases skill and rubric"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 3: Remove Element evaluation from eval-test-cases skill and rubric

## Description
eval-test-cases currently evaluates Element field quality as part of its Dimension 3 web-ui scoring ("Route & Element Accuracy", including "Elements are identifiable"). Since Element field is being removed from gen-test-cases output, the evaluation rubric must be updated to reflect the new output format.

Rename Dimension 3 web-ui from "Route & Element Accuracy" to "Route Accuracy" and remove all Element-related evaluation criteria.

## Reference Files
- `docs/proposals/test-cases-separation/proposal.md` — Source proposal
- `plugins/forge/skills/eval-test-cases/` — Target skill directory

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval-test-cases/SKILL.md` | Remove Element field references from evaluation instructions. Update Dimension 3 web-ui section. |
| `plugins/forge/skills/eval-test-cases/templates/rubric.md` | Rename Dimension 3 web-ui title from "Route & Element Accuracy" to "Route Accuracy". Remove "Elements are identifiable" evaluation item. Redistribute points within the dimension (total score 200pts for web-ui must be preserved). |

## Acceptance Criteria
- [ ] eval-test-cases SKILL.md no longer instructs agents to evaluate Element field quality
- [ ] rubric.md Dimension 3 web-ui title is "Route Accuracy" (not "Route & Element Accuracy")
- [ ] "Elements are identifiable" evaluation item is removed from rubric
- [ ] Total scoring for web-ui dimension remains 200pts (points redistributed to remaining items)

## Implementation Notes
- The rubric uses a 1000-point scale with dimension-level breakdowns. Ensure the web-ui dimension total stays at 200pts after removing Element items
- Check if the eval-test-cases SKILL.md has any sample output or reference that includes Element field — these should also be updated
- The non-web-ui profile evaluation (e.g., go-test CLI tests) should be unaffected as it doesn't have the Element dimension
