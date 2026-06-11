---
feature: "forge-cli-v3"
status: tasks
---

# Feature: forge-cli-v3

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | Task CLI 扩展为 Forge CLI：5 组 + 5 顶层命令，AI-first 命名，e2e 从 justfile 迁移 |
| User Stories | prd/prd-user-stories.md | 6 个用户故事覆盖 AI agent、hook 系统、开发者三类角色 |
| Tech Design | design/tech-design.md | 四阶段重构：基础重命名→命令重组→e2e迁移→引用更新，新增 pkg/e2e 和 fix-task 上限 |
| Design Eval | design/eval/report.md | 966/1000 分，3 轮对抗迭代达标，Breakdown-Readiness 195/200 |

## Traceability

| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| Background + Scope (prd-spec §1-2) | Overview + Architecture (tech-design §1-2) | — | — | 1.1 |
| 命令结构规格 (prd-spec §3) | Interfaces §1-5 (tech-design §3) | — | — | 2.1, 2.2 |
| 命名变更规格 (prd-spec §3) | Renamed Commands (tech-design §3.2) | — | — | 2.2 |
| e2e 从 justfile 迁移 (prd-spec §3) | New E2E Subcommands (tech-design §3.5) + pkg/e2e (§3.6) | — | — | 3.1, 3.2 |
| Error Handling (prd-spec §4) | Error Handling (tech-design §4) | — | — | 3.1, 2.4, 2.5 |
| Performance Requirements (prd-spec §5) | Architecture (tech-design §2) | — | — | 1.1 |
| Hook 自动触发 (user-stories §4) | Quality-Gate Max Fix-Task Cap (tech-design §3.7) | — | — | 2.4, 4.1 |
| e2e 测试 profile (user-stories §5, §8) | New E2E Subcommands (tech-design §3.5) | — | — | 3.1, 3.2 |
| task list-types (user-stories §6) | list-types command (tech-design §3.3) | — | — | 2.3 |
| 并发写冲突 (user-stories §3) | Concurrent Write Conflict (tech-design §9) | — | — | 2.5 |
| 引用更新 (prd-spec §3 Related Changes) | Phase 4 Reference Update Map (tech-design Appendix) | — | — | 4.1, 4.2, 4.3 |
