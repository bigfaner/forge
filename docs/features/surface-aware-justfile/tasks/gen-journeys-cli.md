---
id: "T-test-gen-journeys-cli"
title: "Generate Test Journeys (cli)"
priority: "P1"
estimated_time: "20-30min"
dependencies: []
type: "test.gen-journeys"
scope: "all"
---

Generate test Journey documents for the surface-aware-justfile feature.
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

Invoke the `/gen-journeys` skill to extract Journey narratives from specification documents.

### Input Source by Mode

- **Breakdown mode**: Read PRD user stories from `docs/features/surface-aware-justfile/prd/prd-user-stories.md` and functional specs from `docs/features/surface-aware-justfile/prd/prd-spec.md`. These are the primary input sources.
- **Quick mode**: Read the proposal from `docs/proposals/surface-aware-justfile/proposal.md`. Extract Key Scenarios as Journey candidates. If the proposal lacks `scope` or `success criteria` sections, abort the task with a diagnostic message — Journey generation requires these minimum inputs.

## Process

Follow the `/gen-journeys` skill process flow:

1. **Surface Detection**: Detect the project surface type and persist to `.forge/config.yaml`
2. **Read Sources**: Read PRD user stories (Breakdown) or proposal.md (Quick)
3. **Identify Workflows**: Map each user story or key scenario to a Journey candidate
4. **Classify Risk**: Assign High/Medium/Low risk to each Journey based on workflow characteristics
5. **Generate Files**: Output one `journey.md` per Journey to `docs/features/surface-aware-justfile/testing/<journey-name>/journey.md`
6. **Validate Output**: Check each Journey for required fields (name, risk level, happy path steps, edge cases, invariants)

## AUTO_COMMIT Directive

When this task runs as an automated pipeline task (not invoked manually by the user), AUTO_COMMIT=true is in effect:

- **If AUTO_COMMIT=true**: Skip the user review-and-approval step. After validation passes, directly commit all generated Journey files:
  ```bash
  git add docs/features/surface-aware-justfile/testing/
  git commit -m "docs: generate journeys for surface-aware-justfile"
  ```
- **If AUTO_COMMIT is not set** (manual invocation): Present all Journey files to the user for review. Wait for explicit approval before committing.

## Acceptance Criteria

- [ ] At least 1 Journey file generated under `docs/features/surface-aware-justfile/testing/`
- [ ] Each Journey has: name, risk level, happy path steps, edge cases, invariants
- [ ] High-risk Journeys have edge case count >= happy path step count
- [ ] All Journey files committed (AUTO_COMMIT=true) or awaiting user review (manual mode)

Type: **cli**
