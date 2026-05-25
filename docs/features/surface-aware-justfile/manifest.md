---
feature: "surface-aware-justfile"
created: "2026-05-25"
status: tasks
---

# Feature: surface-aware-justfile

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | init-justfile surface 感知配方生成 + run-tests 编排简化 + surface-key 统一迁移 + Task 数据模型扩展 |
| User Stories | prd/prd-user-stories.md | 4 个用户故事覆盖配方生成、测试编排、surface-key 迁移、数据模型扩展 |
| Tech Design | design/tech-design.md | 三线并行迁移：Go 数据模型 + Skill 规则体系 + CLI --json 增强，干净迁移无兼容层 |

## Spec Consolidation

| Spec File | Rules | Integrated To |
|-----------|-------|---------------|
| specs/biz-specs.md | 9 rules (6 CROSS, 3 LOCAL) | docs/business-rules/surface-orchestration.md |
| specs/tech-specs.md | 8 specs (5 CROSS, 3 LOCAL) | docs/conventions/surface-cli.md, docs/conventions/surface-rules.md |

Consolidated on 2026-05-26.

## Traceability

| PRD Section | Design Section | UI Component | Tasks |
|-------------|----------------|--------------|-------|
| Story1: surface 感知配方生成 | Interface 2: Surface 规则文件格式 | — | 3.1, 3.2 |
| Story1: CLI/TUI 不生成 run | Interface 2: cli/tui 差异点 | — | 3.1, 3.2 |
| Story1: 双平台变体 | Interface 2: 实现约束 | — | 3.1 |
| Story2: 调度器模式编排 | Interface 2: 编排序列 | — | 3.3, 3.4 |
| Story2: probe 失败 exit code | Error Handling: Exit Code 语义 | — | 3.3 |
| Story2: surface 信息不可用 | Interface 1b: --json 错误格式 | — | 1.1, 3.3 |
| Story3: surface-key 迁移 | Data Models: Task.SurfaceKey | — | 1.2a, 1.2b, 2.1 |
| Story3: resolveScope 删除 | Phase Component Map: prompt.go | — | 1.4 |
| Story4: Task 新增双字段 | Data Models: Task/AutoGen/Frontmatter | — | 1.2a, 1.2b |
| Story4: forge task add 继承 | Model 4: AddTaskOpts | — | 1.2b |
| Story4: fix-task 推断 | Phase 2: quality-gate | — | 2.3 |
| Story4: 旧任务 scope 迁移 | Migration Notes + forge task migrate | — | 1.3 |
