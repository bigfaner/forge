---
feature: "forge-cli-readability-round2"
created: "2026-06-06"
status: completed
mode: quick
---

# Feature (Quick): forge-cli-readability-round2

<!-- Status flow: tasks -> in-progress -> completed -->

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/forge-cli-readability-round2/proposal.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | 删除死代码：requireSurfaceInference、extractScope、extractBulletItems | pending | tasks/1-delete-dead-code.md |
| 2 | 拆分 BuildIndex 390 行上帝函数 | pending | tasks/2-split-build-index.md |
| 3 | 拆分 config.go 为三文件 | pending | tasks/3-split-config-go.md |
| 4 | 拆分 pipeline.go 提取校验逻辑 | pending | tasks/4-split-pipeline-go.md |
| 5 | 拆分 detect_surface.go 提取信号表 | pending | tasks/5-split-detect-surface.md |
| 6 | 拆分 runExtract 304 行函数 | pending | tasks/6-split-run-extract.md |
| 7 | 拆分 runList 217 行函数 | pending | tasks/7-split-run-list.md |
| 8 | 拆分 doSubmit 131 行函数 | pending | tasks/8-split-do-submit.md |
| 9 | 平坦化 validate.go 嵌套过深的 validator 方法 | pending | tasks/9-flatten-validate-nesting.md |
| 10 | 重构 quality_gate.go 的 os.Exit 反模式 | pending | tasks/10-refactor-os-exit.md |
