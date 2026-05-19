---
name: guide-common-knowledge
status: Draft
created: 2026-05-19
---

# Guide 公共知识提取

## Problem

Forge plugin 的 19 个 skill 和 18 个 command 文件中，存在大量重复的公共知识。这些知识没有统一的定义来源，而是各自完整地重复书写。这导致：

1. **维护负担**：修改一个公共协议需要在 9+ 个文件中同步更新
2. **一致性风险**：各文件描述的同一协议已出现措辞分歧（如 language detection 在不同 skill 中的 fallback 描述略有不同）
3. **token 浪费**：agent 读取 skill 时，每次都要加载大段重复内容

## Solution

将高频重复的公共知识提取到 `plugins/forge/hooks/guide.md` 中，各 skill/command 改为引用（"See guide.md → Language Detection Protocol"），仅保留 skill 特有的变体说明。

## Findings: 重复知识清单

### F1: Language Detection Protocol (9 个 skill)

**完全重复的流程**：

```
1. Detect language: forge test detect
2. On failure: ask user to add languages to .forge/config.yaml
3. Do NOT silently default to any language
```

**出现在**：
- `skills/breakdown-tasks/SKILL.md` (Step 0)
- `skills/gen-contracts/SKILL.md` (Step 0)
- `skills/gen-test-cases/SKILL.md` (Step 0)
- `skills/gen-test-scripts/SKILL.md` (Step 0) — 额外加载 strategy + framework
- `skills/run-e2e-tests/SKILL.md` (Step 0) — 额外加载 strategy
- `skills/init-justfile/SKILL.md` (Step 0)
- `skills/quick-tasks/SKILL.md` (Step 0)
- `skills/tech-design/SKILL.md` (Step 0)
- `skills/eval/SKILL.md` (Prerequisites table, 仅 test-cases 类型)

**变体**：gen-test-scripts 和 run-e2e-tests 在 detect 之后额外调用 `forge test get generate/run` 和 `forge test framework`，这些是 skill 特有的，应保留在各自 skill 中。

### F2: Docs-Only Fast Path (2 个 skill + 多个 template)

**完全重复的逻辑**：

```
When all tasks are type: "documentation" (non-compilable output):
- Skip Step 0 (language detection)
- Skip standard test tasks (breakdown) 或 Step 4 (quick-tasks)
- forge task index is always mandatory
```

**出现在**：
- `skills/breakdown-tasks/SKILL.md` (Docs-Only Fast Path section)
- `skills/quick-tasks/SKILL.md` (Docs-Only Fast Path section)

**相关 template**：breakdown-tasks 和 quick-tasks 下的多个 template 也引用 `type: "documentation"` 来决定是否生成测试任务。

### F3: HARD-GATE / HARD-RULE / EXTREMELY-IMPORTANT 标记协议 (30+ 文件)

**现状**：这三个标记在 skill/command 文件中广泛使用，但含义没有统一定义。各文件各自解释其语义：

- `HARD-GATE`：产出约束（不写代码、只生成文档等）— 每次都重新解释
- `HARD-RULE`：流程约束（不选技术、不静默默认等）— 每次都重新解释
- `EXTREMELY-IMPORTANT`：关键行为约束 — 每次都重新解释

**出现在**：30+ 文件（所有 skill/command 几乎都有至少一个）

### F4: Prerequisites 模式 (14 个 skill)

**完全重复的格式**：

```
| Artifact | Missing prompt |
|----------|----------------|
| path     | Run /xxx first |
```

**出现在**：breakdown-tasks, eval, forensic, gen-test-cases, quick-tasks, consolidate-specs, gen-contracts, tech-design, ui-design, gen-journeys, run-e2e-tests, improve-harness, init-justfile, write-prd

**没有统一说明**各 skill 独立描述格式，但语义相同。

### F5: Pipeline Position (4 个 skill)

**完全重复的 pipeline 概览图**：

```
gen-journeys → gen-contracts → gen-test-scripts → run-tests
```

**出现在**：
- `skills/gen-journeys/SKILL.md`
- `skills/gen-contracts/SKILL.md`
- `skills/gen-test-scripts/SKILL.md`
- `skills/run-e2e-tests/SKILL.md`

每个 skill 都画了完整的 pipeline 图，只是标注了 "YOU ARE HERE"。

## Deduplication Plan

### Phase 1: Extract to guide.md

在 guide.md 的 `## Execution Rules` 下新增以下子节：

#### 1. Language Detection Protocol

提取通用流程，保留各 skill 的变体说明。

#### 2. Docs-Only Fast Path

提取统一的检测逻辑和跳过规则。

#### 3. Skill Annotation Protocol

统一定义 `HARD-GATE`、`HARD-RULE`、`EXTREMELY-IMPORTANT` 三种标记的语义和用法规范。

#### 4. Test Pipeline Overview

提取 pipeline 全景图，各 skill 只标注自己的位置。

#### 5. Prerequisites Convention

提取格式约定，各 skill 的具体 artifact 列表保留在各自文件中。

### Phase 2: Update Skills/Commands

对每个重复知识点，将 skill 中的完整描述替换为引导引用：

**替换前**（以 gen-test-cases 为例）：
```markdown
## Step 0: Resolve Language and Interfaces

1. **Detect language**: Run `forge test detect` to auto-detect the project's test language(s) from file signals.
2. **On failure** (no language detected): ask the user to add `languages` to `.forge/config.yaml` (e.g., `languages: [go]`).
3. **Load interfaces**: Run `forge test interfaces` to get the project's active interface types.

<HARD-RULE>
Do NOT silently default to any language. If `forge test detect` returns no result and the user cannot configure `languages`, abort the skill.
</HARD-RULE>
```

**替换后**：
```markdown
## Step 0: Resolve Language and Interfaces

Follow the Language Detection Protocol (guide.md), then run `forge test interfaces` to load active interface types.
```

### 具体文件变更清单

| File | F1 Language | F2 Docs-Only | F3 Annotations | F5 Pipeline |
|------|------------|-------------|----------------|-------------|
| `skills/breakdown-tasks/SKILL.md` | ✅ replace | ✅ replace | - | - |
| `skills/gen-contracts/SKILL.md` | ✅ replace | - | - | ✅ replace |
| `skills/gen-test-cases/SKILL.md` | ✅ replace | - | - | - |
| `skills/gen-test-scripts/SKILL.md` | ✅ replace | - | - | ✅ replace |
| `skills/run-e2e-tests/SKILL.md` | ✅ replace | - | - | ✅ replace |
| `skills/init-justfile/SKILL.md` | ✅ replace | - | - | - |
| `skills/quick-tasks/SKILL.md` | ✅ replace | ✅ replace | - | - |
| `skills/tech-design/SKILL.md` | ✅ replace | - | - | - |
| `skills/eval/SKILL.md` | ✅ replace | - | - | - |
| `skills/gen-journeys/SKILL.md` | - | - | - | ✅ replace |
| 30+ files (HARD-GATE/RULE) | - | - | ✅ replace | - |

### Phase 3: F4 Prerequisites 暂不提取

Prerequisites 的 artifact 列表因 skill 而异，提取格式约定收益较低。保留现状，仅在 guide.md 中记录格式约定即可。

## Scope

### In Scope

- 向 guide.md 新增 5 个公共知识节（约 80-120 行）
- 更新 9 个 skill 的 Language Detection 引用
- 更新 2 个 skill 的 Docs-Only 引用
- 更新 4 个 skill 的 Pipeline Position 引用
- 统一定义 3 种 annotation 标记的语义（guide.md 中新增约 20 行）

### Out of Scope

- 不改变任何 skill 的实际行为逻辑
- 不重构 skill 文件结构
- 不修改 template 文件（它们由 skill 引用，不是独立文档）
- F4 Prerequisites 的具体内容不做去重（格式约定可记录但不强制迁移）

## Risks

| Risk | Mitigation |
|------|-----------|
| Agent 不读取 guide.md | guide.md 已通过 hook 注入 system prompt，agent 必然读取 |
| 去重后 skill 上下文不完整 | 保留 skill 特有变体，仅替换通用部分 |
| guide.md 膨胀 | 新增约 100 行，当前 111 行 → 约 210 行，仍在可接受范围 |
| Phase 2 修改量大 | 可分批进行，每个知识类型独立提交 |

## Success Criteria

- [ ] guide.md 包含 Language Detection Protocol 定义
- [ ] guide.md 包含 Docs-Only Fast Path 定义
- [ ] guide.md 包含 Skill Annotation Protocol 定义
- [ ] guide.md 包含 Test Pipeline Overview
- [ ] 9 个 skill 的 Language Detection 替换为引用
- [ ] 2 个 skill 的 Docs-Only Fast Path 替换为引用
- [ ] 4 个 skill 的 Pipeline Position 替换为引用
- [ ] 所有功能行为不变（guide.md 变更不影响运行时行为）
