---
id: "4"
title: "Make gen-test-cases Element field required"
priority: "P1"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
noTest: false
mainSession: false
---

# 4: Make gen-test-cases Element field required

## Description

Change the Element field in gen-test-cases from optional to required. Currently, agents skip the Element field when sitemap data is incomplete, causing gen-test-scripts to use fuzzy locator matching that produces incorrect selectors.

This is a breaking change to the test-cases.md schema — downstream consumers (gen-test-scripts, eval-test-cases) must handle the "sitemap-missing" sentinel value.

## Reference Files
- `docs/proposals/forge-testing-optimization/proposal.md` — Source proposal (Phase 2, Section 2.4)
- `plugins/forge/skills/gen-test-cases/SKILL.md` — Target: mark Element as required
- `plugins/forge/skills/gen-test-cases/templates/test-cases.md` — Template: update Element field annotation

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-cases/SKILL.md` | Add HARD-RULE: Element field is required; add fallback for missing sitemap |
| `plugins/forge/skills/gen-test-cases/templates/test-cases.md` | Mark Element as `required` in the template schema |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] gen-test-cases SKILL.md has a HARD-RULE stating Element field is required for every test case
- [ ] SKILL.md documents the "sitemap-missing" sentinel: if sitemap.json doesn't exist for the target page, Element is set to `sitemap-missing` with a WARNING note in test-cases.md
- [ ] SKILL.md documents the Route Validation behavior: if sitemap exists but lacks element data for the target route, report the gap and suggest running `/gen-sitemap`
- [ ] test-cases.md template marks Element field as required (not optional)
- [ ] gen-test-scripts SKILL.md or templates have a note about handling `sitemap-missing` Element values: use Fact Table DOM structure to infer locators

## Implementation Notes

1. **Breaking change**: This modifies the test-cases.md contract. Existing features with optional Element fields will still work (gen-test-scripts already handles missing Element). The change ensures *new* generation always includes Element data
2. **sitemap-missing sentinel**: When sitemap.json doesn't exist at all, don't block — set Element to `sitemap-missing` and let gen-test-scripts handle it. This avoids a hard dependency on `/gen-sitemap` for projects that haven't adopted it yet
3. **Partial sitemap gap**: When sitemap exists but lacks data for a specific route, that's a gap worth flagging. The Route Validation step (already exists in gen-test-cases) should report this
4. **Template change**: In `test-cases.md` template, change `Element: (optional)` to `Element: (required)` and add a note about the sentinel value
