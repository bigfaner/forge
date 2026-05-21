---
feature: "forge-architecture-simplification"
status: tasks
---

# Feature: forge-architecture-simplification

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 7 个价值域（数据可靠性、行为正确性、错误清晰度、Eval 安全性、配置能力、CLI 一致性、代码健康）的重构 PRD |
| User Stories | prd/prd-user-stories.md | 7 个用户故事覆盖 3 个角色（CLI 开发者、Plugin 开发者、终端用户） |
| Tech Design | design/tech-design.md | 基于代码库现状重写：修正 internal/cmd/ 架构、标注已有组件（SaveIndexAtomic/LockFile/AIError）、4 接口设计 + 错误处理策略 + 测试策略 |
| State Diagram | design/state-transition-diagram.md | 任务状态流转图（6 状态 + 4 角色 + reopen 机制） |
| Tasks | tasks/index.json | 33 任务（4 Phase 渐进式 + 辅助任务），总工时约 35h |

## Traceability

| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
| DR-1~4 Data Reliability | Interface 2 (WithLock + SaveStateAtomic) | 2.3, 2.6 |
| BC-1~10 Behavioral Correctness | Interface 1 (StateMachine) | 2.1, 2.4, 2.5, 2.6, 2.7, 2.8 |
| BC-11 Reopen Command | Interface 1 (RoleReopen) | 2.5 |
| EC-1~8 Error Clarity | Error Handling | 2.2, 2.4, 2.8, 2.10 |
| ES-1~4 Eval Safety | Plugin Changes | 2.9 |
| CE-1~5 Configuration | Interface 4 (Config CRUD) | 3.2, 3.3, 3.4 |
| CC-1~4 CLI Consistency | RunE Migration | 3.1, 3.2 |
| CH-1~4 Code Health | Constants + Naming | 1.1, 1.2, 1.3 |
