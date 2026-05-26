---
id: "T-test-gen-contracts"
title: "Generate Test Contracts"
priority: "P1"
estimated_time: "30-45min"
dependencies: ["T-eval-journey"]
type: "test.gen-contracts"
scope: "all"
---

Generate test Contract specifications for the surface-aware-justfile feature.
Mode: breakdown

## Scope

- 5 个 surface 规则文件：`skills/init-justfile/rules/surfaces/{web,api,cli,tui,mobile}.md`
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

Invoke the `/gen-contracts` skill to generate Contract specifications from Journey documents and code reconnaissance.

### Eval Gate by Mode

- **Breakdown mode**: An eval-journey report must exist for all Journeys before proceeding. Check for `testing/<journey>/.eval-report.md` files. If any Journey lacks an eval report or scored below target, abort this task.
- **Quick mode**: The eval-journey gate is skipped. Proceed directly to Contract generation.

## SKIP_EVAL_GATE Directive

When this task runs in Quick mode as an automated pipeline task, SKIP_EVAL_GATE=true is in effect:

- **If SKIP_EVAL_GATE=true**: Skip the eval-journey prerequisite check. Do not require `testing/<journey>/.eval-report.md` files. Proceed directly to code reconnaissance and Contract generation.
- **If SKIP_EVAL_GATE is not set** (Breakdown mode): Require eval reports for all Journeys. Abort if any Journey scored below the eval target threshold.

## Process

Follow the `/gen-contracts` skill process flow:

1. **Resolve Language & Interfaces**: Detect project language and interface types from config
2. **Read Journeys**: Enumerate Journey directories under `docs/features/surface-aware-justfile/testing/` and read each `journey.md`
3. **Code Reconnaissance**: Build the Fact Table by reading source code per the reconnaissance rules
4. **Generate Contracts**: For each Journey, generate one Contract file per Step with six-dimension declarations (Preconditions, Input, Output, State, Side-effect, Invariants). Apply risk-driven Outcome density.
5. **Validate Contracts**: Schema validation for structural completeness. Retry once on failure.
6. **Write Output**: Write Contract files to `docs/features/surface-aware-justfile/testing/<journey>/contracts/` and Fact Table to `.forge/fact-table.json`

## Acceptance Criteria

- [ ] At least 1 Contract file generated per Journey
- [ ] Each Contract has six-dimension declarations with semantic descriptors (no regex)
- [ ] Risk-driven Outcome density targets met per Journey risk level
- [ ] Fact Table written to `.forge/fact-table.json`
- [ ] All Contracts passed schema validation
