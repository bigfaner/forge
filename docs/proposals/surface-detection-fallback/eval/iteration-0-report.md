---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Iteration 0: Pre-Revision Report

## ATTACK_POINTS

### Factual Corrections

- **[high]** Sources map key 在 success criteria 中不一致 | quote: `"Sources` map correctly populated: `{\"forge-cli/cli\": \"inference:cmd-dir\"}` for inferred paths, `{\"forge-cli/cli\": \"dependency:cobra\"}` for detected paths" vs 另一 criterion 用 `"cli"` 作为 key | improvement: 统一 Sources map key 为路径格式（与 Surfaces map key 一致），修正所有 success criteria 中的 key 示例

- **[medium]** Node.js `index.html` 推断规则说 "at root" 但未约束推断函数仅在根目录生效 | quote: "`index.html` at root -> `web`" | improvement: 在 Key Scenarios 中明确 `inferNodeSurface` 的 `index.html` 检测仅在项目根目录（与 `package.json` 同级）触发，不在子目录扫描中生效

### Structural/Architectural Suggestions

- **[high]** `forge surfaces detect` 是写命令但嵌套在只读的 `forge surfaces` 下，与 list/query/types 的只读语义矛盾 | quote: "`forge surfaces detect` command: new subcommand that runs detection + inference and displays results, independent of `forge init`" | improvement: 将 `forge surfaces detect` 改为默认只读（显示结果），添加 `--apply` flag 启用配置写入，`--dry-run` 变为默认行为的冗余别名

- **[high]** Re-run 只有 Confirm/Re-detect 两选项，缺少 Edit，与现有 `askMapConfirmation` 的四选项设计不一致 | quote: "TUI shows `"Surfaces already configured: cli. Re-detect?"` with Confirm (keep existing) / Re-detect options" | improvement: 增加第三个选项 Edit（进入 manualSurfaceEntry 流程），形成 Confirm / Re-detect / Edit 三选项

- **[high]** Source 不持久化在 `forge surfaces detect` 的 re-detection 场景中造成信息丢失 | quote: "Source information is display-only and discarded after TUI confirmation completes" | improvement: 在 config 中以 YAML 注释形式保留来源元数据（如 `surfaces: api  # source: inference:cmd-dir`），不修改 SurfacesMap schema

## SKIPPED_FINDINGS (Subjective Preference)

- 性能预算缺乏测量基准（已有 NFR 声明 <50ms，分拆方案留给实现阶段）
- Sources map value 格式规范（留给实现阶段定义常量）
- 非交互模式 exit code 不一致（两命令不同入口点，exit code 差异合理）
- 代码复用边界文档化（task breakdown 阶段处理）

## Classification Audit

- Factual corrections: 2
- Structural/architectural suggestions: 3
- Subjective preferences: 4 (skipped)
