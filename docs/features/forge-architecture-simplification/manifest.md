---
feature: "forge-architecture-simplification"
status: design
---

# Feature: forge-architecture-simplification

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 7 个价值域（数据可靠性、行为正确性、错误清晰度、Eval 安全性、配置能力、CLI 一致性、代码健康）的重构 PRD |
| User Stories | prd/prd-user-stories.md | 6 个用户故事覆盖 3 个角色（CLI 开发者、Plugin 开发者、终端用户） |
| Tech Design | design/tech-design.md | 4 接口设计（StateMachine、AtomicWrite、PreserveFields、ConfigCRUD）+ 错误处理策略 + 测试策略 |

## Traceability

| PRD Section | Design Section | Interface / Model |
|-------------|----------------|-------------------|
| DR-1~4 Data Reliability | Interface 2: Atomic Write | SaveIndexLocked, SaveStateAtomic |
| BC-1~10 Behavioral Correctness | Interface 1: StateMachine | ValidateTransition, CanAutoUnblock |
| EC-1~8 Error Clarity | Error Handling | AIError factory functions, Exit() |
| ES-1~4 Eval Safety | Plugin Changes | eval/SKILL.md backup+rollback |
| CE-1~5 Configuration Empowerment | Interface 4: Config CRUD | SetAutoKey, GetAutoKeyValue |
| CC-1~4 CLI Consistency | W8 Implementation | RunE migration, config init merge |
| CH-1~4 Code Health | W1/W2 Implementation | Naming, dead code, constants |
