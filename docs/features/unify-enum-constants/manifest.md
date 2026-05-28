---
feature: "unify-enum-constants"
created: "2026-05-28"
status: design
---

# Feature: unify-enum-constants

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 统一枚举常量，250+ 处魔法值替换为 typed constants + 全量签名升级 |
| User Stories | prd/prd-user-stories.md | 4 个用户故事：类型安全、集中定义、可靠重构、验证整合 |
| Tech Design | design/tech-design.md | 3 个 typed constant 类型，6 phase 按 package 迁移，type alias 重导出 |

## Traceability

| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
| US1: Type-Safe Status | Interface 1: Status, Model 1: Transition | Phase 1, 3, 5 |
| US2: Centralized SurfaceType | Interface 2: SurfaceType, Model 3 | Phase 1, 4 |
| US3: Reliable Refactoring | Architecture (leaf package), Testing Strategy | Phase 1-6 |
| US4: Validation Consolidation | Interface 1: AllStatuses(), Interface 3: AllPriorities() | Phase 3, 5 |
