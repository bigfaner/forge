---
id: "6"
title: "移除 test.gen-and-run 废弃代码 + 更新 isTestTaskID + 更新测试文件"
priority: "P2"
estimated_time: "2h"
dependencies: ["1", "5"]
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 6: 移除 test.gen-and-run 废弃代码

## Description

按 Removal Checklist 逐一移除 `test.gen-and-run` / `T-quick-gen-and-run` 在所有文件中的引用。更新 `isTestTaskID` 覆盖 `T-review-doc`。更新 `validate_index.go` 提供迁移感知错误提示。更新所有引用废弃类型的测试文件。

## Reference Files
- `docs/proposals/pipeline-integration-stitch/proposal.md` — Removal Checklist
- `forge-cli/pkg/task/types.go` — TypeTestGenAndRun 常量
- `forge-cli/pkg/task/infer.go` — InferType gen-and-run 分支
- `forge-cli/pkg/task/build.go` — isTestTaskID 函数
- `forge-cli/internal/cmd/task/validate_index.go` — T-quick-gen-and-run- 前缀检查

## Acceptance Criteria
- [ ] `types.go` 移除 `TypeTestGenAndRun` 常量
- [ ] `infer.go` 移除 InferType 中 gen-and-run 分支
- [ ] `prompt.go` genScriptBases 移除 `T-quick-gen-and-run` 条目
- [ ] `autogen.go` 移除 gen-and-run 相关逻辑（如仍有）
- [ ] `validate_index.go` 对引用 `test.gen-and-run` 的旧 index.json 返回包含 "deprecated" 和 "regenerate" 的迁移指引错误信息
- [ ] `isTestTaskID` 覆盖 `T-review-doc`，函数文档注释更新
- [ ] 删除 `data/` 下 gen-and-run 相关 .md 模板文件
- [ ] `grep -r "gen-and-run" forge-cli/ plugins/` 返回零结果（不含 validate_index.go 中的迁移提示文本）
- [ ] 所有现有测试通过

## Hard Rules
- 移除必须按清单逐项执行，不可遗漏
- validate_index.go 的迁移提示必须包含 "deprecated" 和 "regenerate" 关键词

## Implementation Notes
- Removal Checklist 共 10 项，见 proposal.md Key Risks 章节
- 测试文件更新：`grep -rl "gen-and-run\|quick-gen-and-run\|T-quick-gen" forge-cli/ plugins/ --include="*_test.go"` 确定需修改的文件列表
- isTestTaskID 扩展后与 IsAutoGenTaskID 覆盖范围应一致
