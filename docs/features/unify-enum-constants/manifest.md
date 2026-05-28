---
feature: "unify-enum-constants"
created: "2026-05-28"
status: tasks
---

# Feature: unify-enum-constants

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 统一枚举常量，250+ 处魔法值替换为 typed constants + 全量签名升级 |
| User Stories | prd/prd-user-stories.md | 4 个用户故事：类型安全、集中定义、可靠重构、验证整合 |
| Tech Design | design/tech-design.md | 3 个 typed constant 类型，5 phase 按 package 迁移，type alias 重导出 |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1.1 | Create pkg/types/ package with typed constants | pending | tasks/1.1-define-types.md |
| 2.1 | Migrate pkg/task/ core types and statemachine | pending | tasks/2.1-migrate-task-core.md |
| 2.2 | Migrate pkg/task/ remaining files | pending | tasks/2.2-migrate-task-remaining.md |
| 3.1 | Migrate pkg/forgeconfig/ SurfaceType constants | pending | tasks/3.1-migrate-forgeconfig.md |
| 4.1 | Migrate internal/cmd/task/ Status and Priority constants | pending | tasks/4.1-migrate-cmd-task.md |
| 4.2 | Migrate internal/cmd/ other files (all enums) | pending | tasks/4.2-migrate-cmd-other.md |
| 5.1 | Re-export typed constants from pkg/feature/ | pending | tasks/5.1-reexport-feature.md |

## Traceability

| PRD Section | Design Section | UI Component | Tasks |
|-------------|----------------|--------------|-------|
| US1: Type-Safe Status | Interface 1: Status, Model 1: TransitionRule | — | 1.1, 2.1, 4.1 |
| US2: Centralized SurfaceType | Interface 2: SurfaceType, Model 3 | — | 1.1, 3.1, 4.2 |
| US3: Reliable Refactoring | Architecture (leaf package), Testing Strategy | — | 1.1, 2.1-2.2, 5.1 |
| US4: Validation Consolidation | Interface 1: AllStatuses(), Interface 3: AllPriorities() | — | 2.2, 4.1 |
