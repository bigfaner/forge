---
feature: "feat-log-decisions"
status: tasks
---

# Feature: feat-log-decisions

<!-- Status flow: prd → design → tasks → in-progress → done -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | Tech design skill 重命名（design-tech → tech-design）+ 建立 docs/decisions/ 集中决策归档目录，新增 /zcode:record-decision 命令 |
| User Stories | prd/prd-user-stories.md | 4 个用户故事：统一 skill 命名、tech-design 流程中归档决策、独立记录决策、查阅历史决策 |
| Tech Design | design/tech-design.md | 三类变更：skill 目录重命名、docs/decisions/ 数据存储、record-decision 新 skill；共享逻辑提取至 plugins/zcode/references/shared/decision-logging.md |

## Traceability

| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
| 5.1 决策归档目录结构 | Data Models §Model 2/3 | [2.1](tasks/2.1-create-decisions-directory.md), [2.2](tasks/2.2-create-decision-entry-template.md) |
| 5.2 tech-design 归档步骤 | Interfaces §Interface 3 | [3.1](tasks/3.1-create-decision-logging-shared.md), [3.2](tasks/3.2-update-tech-design-skill.md) |
| 5.3 /zcode:record-decision | Architecture §record-decision skill | [4.1](tasks/4.1-create-record-decision-skill.md), [4.2](tasks/4.2-register-record-decision-skill.md) |
| 5.4 关联性需求改动 | PRD Coverage Map | [1.1](tasks/1.1-rename-design-tech-skill.md), [4.3](tasks/4.3-update-reference-links.md), [5.1](tasks/5.1-create-validate-manifest-script.md) |
