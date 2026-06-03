---
created: "2026-06-03"
author: faner
status: Draft
intent: fix
---

# Proposal: Surface Scalar Form Fix — 统一文本模式贯通测试 Skill 链路

## Problem

当用户在 `.forge/config.yaml` 中使用 scalar 简写形式配置 surface（如 `surfaces: tui`），Go 解析器将其转换为哨兵 key `"."`。`forge surfaces --json` 将此哨兵直接输出为 `"key": "."`，导致依赖 JSON 输出解析 key 的 skill 拒绝该值，surface-aware 功能全部中断。

### Evidence

- `init-justfile` 以正则 `[a-zA-Z0-9_-]+` 校验从 JSON 解析的 key，`"."` 不匹配，跳过所有 surface recipe 生成
- Gotcha 文档记录了完整复现路径：`surfaces: tui` → `forge surfaces --json` 输出 `"key": "."` → init-justfile 报错 `Error: invalid surface-key "."`
- `test-guide` 直接读取 config.yaml，其解析逻辑（逗号分隔 type 字符串）与 Go CLI 的 YAML 解析不一致

### Urgency

Scalar 形式是最简配置方式，新用户最可能使用。当前行为让 surface-aware 功能对这类用户完全不可用，是"开箱即用"体验的硬阻断。

## Proposed Solution

**统一所有 skill 使用 `forge surfaces`（文本模式）替代 `forge surfaces --json`。**

文本模式的输出天然区分了 scalar 和 named 两种形式，CLI 内部已吸收 `"."` 哨兵：

| 配置形式 | 文本输出 | 含义 |
|---------|---------|------|
| Scalar 单 surface | `tui` | 无 key，单 surface → 无 prefix |
| Named 单 surface | `myapp=tui` | key=`myapp` → `myapp-` prefix |
| Multi-surface | `backend=api\nfrontend=web` | 各自 key prefix |

解析规则：每行按 `=` 分割，无 `=` 则行为 scalar（只有 type，无 key）。

### 核心变更

1. **数据源统一**：所有 skill 从 `forge surfaces --json` 切换到 `forge surfaces`（文本模式）
2. **test-guide 数据源修正**：从直接读取 config.yaml 改为调用 `forge surfaces`，与其他 skill 一致
3. **无 prefix 模式**：scalar 形式（文本输出无 `=`）下，recipe 名称直接使用动词（`test`、`build`、`dev`、`teardown`）
4. **Named key 行为不变**：有 `=` 的行按 key prefix 生成 recipe

## Requirements Analysis

### Key Scenarios

- **Scalar 单 surface**：`surfaces: tui` → `forge surfaces` 输出 `tui` → 无 `=` → 生成 `test`、`build` 等 recipe
- **Named 单 surface**：`surfaces: [{key: myapp, type: tui}]` → 输出 `myapp=tui` → recipe 名 `myapp-test`
- **Multi-surface**：输出 `backend=api\nfrontend=web` → 各自 `<key>-<verb>` recipe
- **Task frontmatter 传播**：breakdown-tasks 从文本模式解析 surface 信息，scalar 形式无 key → frontmatter 的 `surface-key` 留空或省略

### Non-Functional Requirements

- **兼容性**：已有 named key 配置的 recipe 命名不变
- **一致性**：所有 skill 使用同一 CLI 命令和同一解析规则

### Constraints & Dependencies

- 不修改 Go CLI 代码（`forge surfaces` 文本输出格式不变）
- 仅修改 skill 层（SKILL.md 文件）

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零成本 | Scalar 形式用户无法使用 surface 功能 | Rejected |
| `--json` + `"."` 防御 | 精确控制 | 每个 skill 需处理哨兵，维护成本高 | Rejected: 复杂度高 |
| 统一文本模式 | CLI 已吸收哨兵，解析简单，零防御代码 | 文本格式无结构化字段扩展性 | **Selected: 从根源消除问题** |

## Feasibility Assessment

5 个 skill 需要修改解析逻辑，预估 1 个 coding task。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "需要 `--json` 获取结构化数据" | 代码审计 | Overturned: 文本模式 `key=type` 格式已包含所有 skill 需要的 key 和 type 信息 |
| "scalar 形式需要 `"."` 哨兵防御" | 第一性原理 | Overturned: `"."` 是 JSON 的内部泄漏，文本模式在 CLI 层已处理，skill 层无需感知 |
| "test-guide 直接读 config.yaml 没问题" | 代码审计 | Overturned: 其解析逻辑（逗号分隔）与 Go CLI 不一致 |

## Scope

### In Scope

- **init-justfile**：`forge surfaces --json` → `forge surfaces`（文本）；更新 key 解析和 recipe 命名逻辑；scalar 无 prefix
- **run-tests**：`forge surfaces --json` → `forge surfaces`（文本）；更新 recipe-prefix 推导；scalar 无 prefix
- **test-guide**：从直接读 config.yaml → `forge surfaces`（文本）；统一数据源
- **breakdown-tasks**：`forge surfaces --json` → `forge surfaces`（文本）；更新 surface-key/type 解析
- **quick-tasks**：`forge surfaces --json` → `forge surfaces`（文本）；更新 surface-key/type 解析

### Out of Scope

- Go CLI 代码修改
- gen-journeys / gen-contracts / gen-test-scripts（已使用文本模式或只用 type，无需改动）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 文本解析逻辑在各 skill 中不一致 | L | M | 在各 skill 中统一写明相同的解析规则 |
| 未来 `forge surfaces` 输出格式变更 | L | H | 文本格式 `key=type` 已是 CLI 稳定接口 |

## Success Criteria

- [ ] `surfaces: tui` 配置下，`init-justfile` 生成 `test`、`build`、`dev`、`teardown` recipe（而非报错或跳过）
- [ ] `surfaces: tui` 配置下，`run-tests` 调用 `just test` 而非 `just tui-test`
- [ ] `surfaces: tui` 配置下，`test-guide` 通过 `forge surfaces` 获取 type，正确生成 `docs/conventions/testing/tui/core.md`
- [ ] Named key 配置（如 `[{key: app, type: tui}]`）输出 `app=tui`，recipe 名 `app-test`，行为不变
- [ ] Multi-surface 配置输出多行 `key=type`，各 recipe 用各自 key prefix，行为不变

## Next Steps

- Proceed to `/quick-tasks` to generate tasks (intent: fix, coding task)
