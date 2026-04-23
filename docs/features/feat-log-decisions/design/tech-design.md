---
created: 2026-04-23
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Tech Design Skill 改进：重命名与决策归档

## Overview

本 feature 涉及三类变更，均为 markdown 文件操作，无运行时代码：

1. **目录重命名**：`plugins/zcode/skills/design-tech/` → `plugins/zcode/skills/tech-design/`，更新所有引用
2. **新数据存储**：`docs/decisions/` 目录，含 8 个类型文件 + `manifest.md`
3. **新 skill**：`plugins/zcode/skills/record-decision/`，共享逻辑提取至 `plugins/zcode/references/shared/decision-logging.md`

## Architecture

### Layer Placement

Plugin layer（`plugins/zcode/`）和 docs layer（`docs/decisions/`）。无业务逻辑层，所有"执行"由 Claude 读取 markdown 文件后完成。

### Component Diagram

```
plugins/zcode/
├── references/
│   └── shared/
│       └── decision-logging.md    ← 共享归档逻辑（两个 skill 共用）
├── skills/
│   ├── tech-design/               ← 从 design-tech 重命名
│   │   ├── SKILL.md               ← 更新：name 字段 + 决策归档步骤
│   │   ├── templates/
│   │   │   ├── tech-design.md
│   │   │   ├── api-handbook.md
│   │   │   ├── manifest-update-design.md
│   │   │   └── decision-entry.md  ← NEW: 单条决策行模板
│   │   └── examples/
│   │       ├── ask-question.md
│   │       └── exploration.md     ← 更新：DECISIONS.md 引用
│   └── record-decision/           ← NEW skill
│       └── SKILL.md               ← 引用 references/shared/decision-logging.md

docs/decisions/                    ← NEW 数据存储
├── manifest.md
├── architecture.md
├── interface.md
├── data-model.md
├── dependencies.md
├── error-handling.md
├── testing.md
├── security.md
└── local-dev-deployment.md

plugins/zcode/hooks/guide.md       ← 更新：DECISIONS.md 引用
zcode/CLAUDE.md                    ← 更新：skill 索引
```

### Component Interactions

```
/zcode:tech-design
    → reads SKILL.md
    → (用户批准后) reads plugins/zcode/references/shared/decision-logging.md
    → appends rows to docs/decisions/<type>.md
    → updates docs/decisions/manifest.md

/zcode:record-decision
    → reads SKILL.md
    → reads plugins/zcode/references/shared/decision-logging.md
    → appends rows to docs/decisions/<type>.md
    → updates docs/decisions/manifest.md
```

### Dependencies

- 无新外部依赖
- 内部依赖：`record-decision/SKILL.md` → `references/shared/decision-logging.md`（跨 skill 引用，路径为相对于 plugin 根目录的绝对路径）

## Interfaces

所有"接口"为 markdown 文件契约，无代码接口。

### Interface 1: Decision Entry（决策条目追加）

目标文件：`docs/decisions/<type>.md`

操作：在文件末尾追加一行表格行。

```
| Date       | Feature   | Decision              | Rationale             | Source                              |
|------------|-----------|-----------------------|-----------------------|-------------------------------------|
| YYYY-MM-DD | <slug>    | <一句话决策描述>       | <一句话决策理由>       | <feature>/tech-design.md §<Section> |
```

字段约束：
- `Date`：ISO 8601 格式（YYYY-MM-DD）
- `Feature`：feature slug，如 `feat-log-decisions`
- `Decision`：单句，不超过 80 字符
- `Rationale`：单句，不超过 80 字符
- `Source`：格式为 `<feature-slug>/<file>.md §<Section>` 或 `manual`

### Interface 2: Manifest Update（manifest 更新）

目标文件：`docs/decisions/manifest.md`

操作 A（Categories 表）：找到对应类型行，将 `Decisions` 计数 +1，`Last Updated` 更新为当前日期。

操作 B（Recent Decisions 表）：在表头下方插入新行（最新在前），保留最近 10 条。

### Interface 3: decision-logging.md（共享归档协议）

文件：`plugins/zcode/references/shared/decision-logging.md`

内容契约：
- 定义 8 个类型名称及对应文件路径映射
- 定义 tech-design 归档步骤（候选列表展示 → 用户选择 → 写入 → manifest 更新）
- 定义 record-decision 4 轮交互步骤
- 定义 `edit:<编号>` 子流程（重新编辑 Decision 或 Rationale 后归档）
- 定义无决策时的跳过逻辑

## Data Models

### Model 1: Decision Entry Row

```
DecisionEntry = {
    Date:      string    // ISO 8601, e.g. "2026-04-23"
    Feature:   string    // feature slug, e.g. "feat-log-decisions"
    Decision:  string    // one sentence, max 80 chars
    Rationale: string    // one sentence, max 80 chars
    Source:    string    // "<slug>/<file>.md §<Section>" or "manual"
}
```

### Model 2: Decisions Manifest

```
DecisionsManifest = {
    frontmatter: {
        updated: string  // ISO 8601 date of last write
    }
    categories: CategoryRow[]
    recentDecisions: RecentDecisionRow[]  // max 10, newest first
}

CategoryRow = {
    Category:    string  // e.g. "Architecture"
    File:        string  // e.g. "architecture.md"
    Decisions:   int     // count of rows in the file
    LastUpdated: string  // ISO 8601 or "-" if empty
}

RecentDecisionRow = {
    Date:     string
    Feature:  string
    Category: string
    Decision: string
    Source:   string
}
```

### Model 3: Type File Initial State

每个类型文件（如 `architecture.md`）初始内容：

```markdown
# <Category> Decisions

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
```

## Error Handling

### Error Types & Codes

| Scenario | Handling |
|----------|----------|
| `docs/decisions/` 目录不存在 | 首次归档时自动创建目录 + 8 个类型文件 + manifest.md |
| `manifest.md` 缺失 | 归档前从模板重建 |
| `/zcode:record-decision` 输入非法类型编号 | 重新提示："请输入 1-8 之间的数字" |
| `edit:<编号>` 引用不存在的候选编号 | 重新提示："编号 X 不在候选列表中，请重新输入" |
| 类型文件表头缺失（文件损坏） | 追加前先补全表头行 |

### Propagation Strategy

所有操作为本地文件写入，错误在操作点处理，不向上传播。用户通过 AskUserQuestion 交互感知错误并重试。

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| 重命名完整性 | 静态检查 | `grep -r "design-tech" plugins/ docs/` | 0 残留引用 | 100% |
| record-decision happy path | 手动 | 执行 skill，输入 4 字段 | 类型文件 +1 行；manifest 计数 +1 | - |
| tech-design 归档步骤 | 手动 | 在含关键决策的 test feature 上运行 | 候选列表展示；选中条目写入 | - |
| 无决策分支 | 手动 | 在无关键决策的 feature 上运行 | 归档步骤静默跳过 | - |
| validate-manifest CI | 脚本 | `bash scripts/validate-manifest.sh` | 计数一致时 exit 0 | - |

### Key Test Scenarios

1. **Happy path A**：`/zcode:tech-design` → 批准 → 选择部分决策归档 → 验证文件和 manifest
2. **Happy path B**：`/zcode:record-decision` → 4 轮输入 → 验证文件和 manifest
3. **edit 子流程**：输入 `edit:1` → 修改 Decision → 确认归档 → 验证修改后内容写入
4. **无决策跳过**：tech-design 无关键决策 → 归档步骤不出现
5. **首次初始化**：`docs/decisions/` 不存在 → 归档后目录和文件自动创建

### Overall Coverage Target

手动验证 5 个关键场景全部通过。

## Security Considerations

### Threat Model

本地文件操作，无网络传输，无用户认证场景。

### Mitigations

无特殊安全措施需要。

## PRD Coverage Map

| PRD Requirement / AC | Design Component | File / Interface |
|----------------------|------------------|-----------------|
| `/zcode:tech-design` 可用，旧命令失效 | 重命名 skill 目录 + 更新所有引用 | `plugins/zcode/skills/tech-design/SKILL.md` |
| tech-design 批准后展示候选决策列表 | decision-logging.md 归档步骤 | `plugins/zcode/references/shared/decision-logging.md` |
| 选择归档后写入类型文件 + 更新 manifest | Decision Entry + Manifest Update 接口 | `docs/decisions/<type>.md` + `manifest.md` |
| 无决策时跳过归档步骤 | decision-logging.md 条件分支 | `plugins/zcode/references/shared/decision-logging.md` |
| `/zcode:record-decision` 4 轮交互完成归档 | New skill | `plugins/zcode/skills/record-decision/SKILL.md` |
| manifest 计数 +1，Recent Decisions 更新 | Manifest Update 接口 | `docs/decisions/manifest.md` |
| hooks guide 不再引用 `docs/DECISIONS.md` | 引用更新 | `plugins/zcode/hooks/guide.md` |
| CLAUDE.md skill 索引包含 `tech-design` | 文档更新 | `zcode/CLAUDE.md` |
| `edit:<编号>` 子流程支持编辑后归档 | decision-logging.md edit 子流程 | `plugins/zcode/references/shared/decision-logging.md` |

## Open Questions

- [ ] `validate-manifest` CI 脚本放在 `scripts/` 还是 `plugins/zcode/scripts/`？（当前设计放 `scripts/`，待确认）

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| decision-logging.md 保留在 tech-design/references/ | 路径更短，无需新目录 | record-decision 需跨 skill 目录引用，路径语义不对称 | 用户选择提取到 shared/ |
| per-decision ADR 文件 | 单条可独立引用，字段完整 | 文件数量膨胀快，填写成本高 | 提案中已排除 |
| 单一 DECISIONS.md | 最简单 | 20+ 条后检索效率下降 | 提案中已排除 |

### References

- `docs/proposals/tech-design-with-decision-logging/proposal.md`
- `docs/features/feat-log-decisions/prd/prd-spec.md`
