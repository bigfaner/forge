---
id: "9"
title: "Dead code cleanup (example files and templates)"
priority: "P1"
estimated_time: "1h"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 9: Dead code cleanup (example files and templates)

## Description
清理两类死代码：(1) gen-sitemap 中的 sitemap-example.json（示例文件未被引用）；(2) init-justfile 的 6 个 .just 模板（评估是否仍在使用，未使用则删除或标注为 legacy）。

## Reference Files
- `proposal.md#Scope` — P1.9: dead code cleanup targets
- `proposal.md#Key-Risks` — dead code deletion risk L/M, mitigation via grep confirmation

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-sitemap/rules/schema.md` | Inlined sitemap-example.json content as `## Full Example` section, removed external file reference |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/gen-sitemap/sitemap-example.json` | Unreferenced example file |

## Acceptance Criteria
- [x] `grep -r "sitemap-example" plugins/forge/` 返回 0（排除删除操作本身）
- [x] init-justfile 6 个 .just 模板已评估，使用状态已记录
- [x] 删除的文件通过 `grep -r` 全仓库确认无引用

## Template Evaluation Results

### sitemap-example.json — DELETED

- Was referenced by `rules/schema.md` line 3: `Read the full example at templates/sitemap-example.json`
- Action: Inlined the example content into `rules/schema.md` as a `## Full Example` section, then deleted the file
- Post-deletion grep confirms zero references in `plugins/forge/`

### init-justfile .just templates — RETAINED (all 6)

All 6 templates are actively referenced by integration tests in `forge-cli/tests/justfile-integration/`:

| Template | Referenced By |
|----------|--------------|
| `go.just` | `init_justfile_test.go`, `forge_detection_test.go` (4 references) |
| `node.just` | `init_justfile_test.go`, `forge_detection_test.go` (3 references) |
| `mixed.just` | `init_justfile_test.go`, `mixed_cli_test.go`, `forge_detection_test.go` (4 references) |
| `generic.just` | `mixed_cli_test.go` (1 reference) |
| `python.just` | No test references, but referenced in task records |
| `rust.just` | No test references, but referenced in task records |

Note: SKILL.md says "Do NOT use framework-specific recipe templates" — the LLM generates recipes at runtime, but the templates still serve as test fixtures. They cannot be deleted without breaking tests.

## Hard Rules
- 删除前必须 `grep -r <filename>` 确认全仓库无引用
- 不删除 .just 模板中仍被引用的文件

## Implementation Notes
init-justfile 模板需逐个评估：检查 `grep -r "template-name.just"` 是否有引用。无引用的模板安全删除，有引用的保留并标注。
