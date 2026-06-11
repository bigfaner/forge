---
id: "2"
title: "Update documentation and project config for interface-only model"
priority: "P1"
estimated_time: "15min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Update documentation and project config for interface-only model

## Description

Update the gotcha lesson doc to reflect the new interface-only model, and configure `interfaces` in this project's `.forge/config.yaml` so that test pipeline tasks are generated correctly.

## Reference Files
- `docs/proposals/decouple-test-tasks-from-languages/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/lessons/gotcha-test-pipeline-no-languages.md` | Rewrite to reflect new model: interfaces config-driven, no language detection |
| `.forge/config.yaml` | Add `interfaces: [api, cli]` |

## Acceptance Criteria

- [ ] Gotcha doc updated: root cause reflects old language detection system, solution points to new `interfaces` config field
- [ ] `.forge/config.yaml` has `interfaces: [api, cli]`
- [ ] `forge task index --feature <any-feature>` now generates test pipeline tasks for this project

## Hard Rules

- Do not remove the gotcha doc — it remains as a historical lesson. Update it to reflect the resolution.

## Implementation Notes

For `.forge/config.yaml`, add:

```yaml
interfaces:
  - api
  - cli
```

This replaces the previous workaround of adding `languages: [go]`.
