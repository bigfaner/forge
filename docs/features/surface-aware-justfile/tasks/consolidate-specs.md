---
id: "T-specs-consolidate"
title: "Consolidate Specs"
priority: "P2"
estimated_time: "20min"
dependencies: ["T-test-verify-regression"]
type: "doc.consolidate"
scope: "all"
---

Extract and consolidate business rules and tech specs from the surface-aware-justfile feature.

## Feature Context
- Scope: - 5 个 surface 规则文件：`skills/init-justfile/rules/surfaces/{web,api,cli,tui,mobile}.md`
- SKILL.md 新增 surface 检测步骤和 surface 感知配方生成流程
- CLI/TUI 只生成 `dev`，不生成 `run`
- 混合项目 dev 配方接受 surface-key 参数
- `# user-customized` 保护机制（差异摘要 + `--force-regenerate`）
- SKILL.md 改为调度器模式，检测 surface type 后加载对应执行策略规则
- 5 个执行策略规则文件：`skills/run-tests/rules/surfaces/{web,api,cli,tui,mobile}.md`
- 编排序列由规则文件定义，run-tests 按规则执行
- 移除 `test.execution` 节点文档（残留配置被 Go YAML 宽松解析模式静默忽略，不影响功能；无需迁移或告警）
- **前置依赖：`forge surfaces` CLI 命令** — 当前状态：需新建。此命令接受文件路径参数，返回 longest-prefix-match 的 surface-key 和 surface-type。breakdown-tasks、quick-tasks、run-tests、quality-gate fix-task 均依赖此 CLI 查询 surface 信息。需在实现本特性前或同时完成开发。
- prompt.go resolveScope() 完全重写为 surfaces map 集合查询
- Task Go struct：Scope→SurfaceKey，新增 SurfaceType；AutoGenTaskDef.TestType→SurfaceType
- 任务模板 frontmatter：scope→surface-key，新增 surface-type
- breakdown-tasks/quick-tasks 生成任务时填充 surface-key 和 surface-type
- forge task add CLI：从源任务继承 surface-key/surface-type
- quality-gate fix-task：从失败文件路径推断 surface-key/type
- init-justfile 混合项目配方 case 分支更新
- 16 个 prompt 模板 SURFACE_KEY 变量值域同步
- 死代码清理：extractTestTypeArg()、genScriptBases

## Discovery Strategy
1. Scan docs/features/surface-aware-justfile/ for all feature documents (PRD, design, task records)
2. Scan docs/proposals/surface-aware-justfile/ for proposal
3. Extract rules and specs from discovered documents
4. Compare against existing specs in docs/business-rules/ and docs/conventions/

Run in non-interactive mode: auto-integrate all CROSS items. Commit with [auto-specs] tag.
