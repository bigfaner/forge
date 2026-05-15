---
status: draft
created: 2026-05-15
---

# 重构 eval-forge 为运行时可靠性审计

## Problem

当前 eval-forge 评估"结构一致性"（frontmatter、目录名、引用完整性），12 个维度 1000 分，历史最高分 965/1000。

手动深度审计（5 个并行 subagent）发现 **39 个运行时问题**（竞态条件、绕过漏洞、指令冲突、token 浪费），eval-forge 的 12 维度**完全没有覆盖**。根因：

1. **维度面向"文件对不对"，不面向"agent 跑起来会不会出问题"**
   - 290 分给 frontmatter/目录名/plugin metadata，这些对运行时可靠性几乎无影响
   - 0 分给工作流完整性、绕过抵抗、指令精确性、token 效率

2. **Dimension 7 (Task CLI Alignment, 240 分) 最接近运行时检查，但只检查 flag/status 是否匹配**，不检查：
   - agent 能否跳过质量门禁（`--force` 绕过一切）
   - eval 循环中 agent 能否伪造分数（主会话可自评自判）
   - 用户确认点是否可跳过（全部是 advisory text）
   - 条件分支是否有遗漏（if 有 then 无 else）

3. **Scorer 方法论是平面检查列表**，不构建工作流图、不做对抗性测试

## Proposed Solution

### 1. 重构 Rubric：6 维度替代 12 维度

| # | 维度 | 分值 | 核心 |
|---|------|------|------|
| 1 | 工作流完整性 | 250 | 全链路状态图可达性验证 |
| 2 | 工作流绕过抵抗 | 250 | 5 类绕过路径对抗测试 |
| 3 | 指令精确性 | 200 | 指令冲突第一顺位 |
| 4 | 跨文件无冗余 | 150 | 三类冗余不同标准 |
| 5 | 引用完整性 | 100 | 模板/agent/cross-skill/hook 引用 |
| 6 | 结构约定 | 50 | frontmatter/eval 模板/目录名 |

**分数分配逻辑**：运行时可靠性（D1+D2+D3）占 700 分（70%），信息效率（D4）占 150 分，基础完整性（D5+D6）占 150 分。旧 rubric 的结构检查从 290 分降到 50 分。

#### Dimension 1: 工作流完整性 (250 pts)

检查全链路状态图：每个 skill 的前置条件、产出物、后继步骤；quick mode 和 full mode 两条路径完整性。

**内嵌工作流规范（Ground Truth）：**

##### Full Mode Pipeline

```
brainstorm → write-prd → eval-prd → [has UI?] → ui-design → eval-ui → prototype → [human review] → tech-design → [db-schema?] → [schema review] → eval-design → breakdown-tasks
```

```
breakdown-tasks → forge task index (auto T-test tasks) → T-test-1 (gen-sitemap + gen-test-cases) → T-test-1b (eval-test-cases) → T-test-2 (gen-test-scripts) → T-test-3 (run-e2e-tests) → T-test-4 (graduate-tests) → T-test-4.5 (verify-regression) → T-test-5 (consolidate-specs)
```

##### Quick Mode Pipeline

```
/quick → /brainstorm → [human confirm] → /quick-tasks → /run-tasks
```

Quick test chain: T-quick-1 (gen-test-cases) → T-quick-2 (gen-test-scripts) → T-quick-3 (run-e2e-tests) → T-quick-4 (graduate-tests) → T-quick-5 (verify-regression)。Skips: gen-sitemap, eval-test-cases, consolidate-specs.

##### Manifest Status Machine

```
prd → design → tasks → in-progress → completed
```

Legal transitions only forward. Quick mode starts at tasks (no prd/design).

##### Per-Skill Precondition/Output Matrix

| Skill | Hard Prerequisites | Outputs | Conditionals |
|-------|-------------------|---------|-------------|
| brainstorm | None | `docs/proposals/<slug>/proposal.md` | → eval-proposal (optional) or write-prd |
| write-prd | Optional: proposal.md, sitemap.json | `prd/prd-spec.md`, `prd/prd-user-stories.md`, `prd/prd-ui-functions.md` (if UI), `manifest.md` (status: prd) | has UI → ui-design; no UI → tech-design |
| eval-prd | `prd/prd-spec.md`, `prd/prd-user-stories.md` | `prd/eval/iteration-{N}.md`, `prd/eval/report.md` | score gate pass/fail |
| ui-design | `prd/prd-ui-functions.md` (hard) | `ui/ui-design.md`, `ui/prototype/` | multi-platform → separate files |
| eval-ui | `ui/ui-design.md` | `ui/eval/iteration-{N}.md`, `ui/eval/report.md` | platform → rubric variant |
| tech-design | `prd/prd-spec.md` (hard) | `design/tech-design.md`, `design/er-diagram.md` + `design/schema.sql` (if db), `manifest.md` (status: design) | db-schema: yes → mandatory ER+schema |
| eval-design | `design/tech-design.md` | `design/eval/iteration-{N}.md`, `design/eval/report.md` | score gate |
| breakdown-tasks | `prd/prd-spec.md` + `design/tech-design.md` (both hard) | `tasks/*.md`, `tasks/index.json`, `manifest.md` (status: tasks) | HAS_UI/NO_UI/HAS_DB/HAS_PLACEMENT tags |
| gen-test-cases | `prd/prd-user-stories.md` + `prd/prd-spec.md` (both hard) | `testing/test-cases.md` | profile → interface types |
| eval-test-cases | `testing/test-cases.md` + PRD docs | `testing/eval/iteration-{N}.md` | Step Actionability < 200 blocks |
| gen-test-scripts | `testing/test-cases.md` (hard) | `tests/e2e/features/<slug>/` | profile → framework; Step Actionability gate |
| run-e2e-tests | justfile + staging area | `tests/e2e/features/<slug>/results/latest.md` | >30% failure → stop |
| graduate-tests | staging area + PASS results + no marker | `tests/e2e/<module>/`, `.graduated/<slug>` | profile → import rewriting |
| consolidate-specs | PRD + design (both hard) | `specs/`, updated `docs/business-rules/`, `docs/conventions/` | CROSS vs LOCAL |
| /quick | User idea | Orchestrates quick pipeline | >10 tasks → STOP |
| quick-tasks | proposal.md (hard) | `tasks/*.md`, `tasks/index.json`, `manifest.md` | max 10 tasks |
| /run-tasks | `tasks/index.json` | Execution loop | 3 consecutive failures → STOP |
| submit-task | Task executed + record.json | `records/*.md`, updated index.json | --force bypasses gate (not auto-downgrade) |

**评分标准：**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 1a. Full mode chain complete | 0-80 | 每个 skill 有前置条件定义和产出物。前置文件由前序 skill 产出。链路无断点。断点 = -20/个 |
| 1b. Quick mode chain complete | 0-40 | Quick 链路完整。缺步 = -20/个 |
| 1c. Conditional branching correct | 0-50 | 每个条件分支有 true-path 和 false-path。缺分支 = -10/个 |
| 1d. Manifest status transitions valid | 0-30 | 状态迁移合法。非法迁移 = -15/个 |
| 1e. Test lifecycle chain intact | 0-50 | 测试链完整。断链 = -15/个 |

#### Dimension 2: 工作流绕过抵抗 (250 pts)

假设自己是"想偷懒的 agent"，逐节点寻找绕过路径。

**5 种绕过类型：**

| Type | Points | Description |
|------|--------|-------------|
| Type 2: Skip quality gates | 0-70 | `--force` 绕过 compile/test/AC。无 justfile 则 gate 静默通过。noTest 跳过。 |
| Type 3: Fake eval results | 0-70 | 主会话可伪造 scorer 输出。分数解析无完整性检查。可跳过 scorer subagent。 |
| Type 1: Skip mandatory interaction | 0-45 | 用户确认点全部是 advisory text（HARD-RULE）。可跳过 brainstorm/write-prd/ui-design/tech-design/consolidate-specs 的确认。 |
| Type 4: Skip required steps | 0-35 | 条件需求依赖 agent 自报（db-schema、placement）。gen-test-scripts Step Actionability gate 仅在 eval report 存在时触发。 |
| Type 5: Lazy shortcuts | 0-30 | 禁止模式（no mock、no sleep）纯 advisory。record.json 指标自报。直接编辑 index.json 绕过所有验证。 |

**评分标准：**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 2a. Quality gate enforcement | 0-70 | 每个 gate 点是否有 CLI 强制或仅 advisory text。零执行且无文档化理由 = -15/个 |
| 2b. Eval integrity | 0-70 | 每个 eval skill 是否要求独立 subagent 评分。Decision gate 是否解析结构化输出。主会话能否伪造分数。弱点 = -25/个 |
| 2c. User interaction enforcement | 0-45 | 每个确认点是否有 enforcement mechanism。纯 advisory = -5/个 |
| 2d. Required step enforcement | 0-35 | 条件需求是否有后置检查。无检查 = -10/个 |
| 2e. Prohibition enforcement | 0-30 | 每个 HARD-RULE 禁止项是否有机械检查。纯 advisory = -5/个 |

**已知绕过向量（来自手动审计，scorer 应验证当前状态）：**

| Vector | Type | Severity | Description |
|--------|------|----------|-------------|
| BV-2.1 | T2 | HIGH | `forge task submit --force` 绕过 compile/test/AC 全部验证 |
| BV-2.4 | T2 | HIGH | Agent 可在 record.json 中伪造 testsFailed: 0（CLI 不验证数字来源） |
| BV-3.1 | T3 | HIGH | 主会课可跳过 scorer subagent 直接声明 SCORE: 950 |
| BV-3.2 | T3 | HIGH | 分数解析无完整性检查，主会话可篡改 scorer 返回值 |
| BV-1.1 | T1 | MED | brainstorm 用户审批可跳过（纯 HARD-RULE） |
| BV-1.2 | T1 | MED | write-prd 用户审批可跳过 |
| BV-1.3 | T1 | MED | /quick 用户确认可跳过 |
| BV-1.4 | T1 | MED | ui-design 原型审查可跳过 |
| BV-1.5 | T1 | MED | tech-design DB schema 审查可跳过 |
| BV-1.7 | T1 | MED | consolidate-specs spec 集成确认可跳过 |
| BV-4.2 | T4 | MED | gen-test-scripts Step Actionability gate 仅在 eval report 存在时触发 |
| BV-4.5 | T4 | MED | placement validation 依赖 sitemap 存在，缺失则跳过 |
| BV-5.1 | T5 | LOW | 禁止模式（no mock/no sleep/no hardcoded URL）纯 advisory |
| BV-5.2 | T5 | LOW | record.json 指标（coverage/testsPassed）自报，无交叉验证 |

#### Dimension 3: 指令精确性 (200 pts)

**优先级：指令冲突 > 步骤歧义 > 条件不完整 > 变量未定义**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 3a. 指令冲突（跨文件） | 0-80 | guide.md vs SKILL.md vs command files 对同一概念的不同描述。如 guide 说"lint blocks"但 skill 说"lint non-blocking"。冲突 = -25/个。**第一顺位检查。** |
| 3b. 步骤歧义 | 0-50 | SKILL.md 步骤是否只有一种理解。模糊动词（"check tests"、"verify quality"）无具体命令 = -10/个 |
| 3c. 条件不完整 | 0-40 | 每个 if-then 有 else 路径或显式"skip"指令。缺 else = -10/个 |
| 3d. 变量解析清晰度 | 0-30 | Agent 填充变量须有来源说明。CLI 填充变量须匹配 `prompt.go` 的 typeToTemplate。未定义的 agent 变量 = -10/个 |

**CLI 填充变量（来自 `prompt.go` Synthesize — 不标记为"未定义"）：**

| Variable | Source | Used by |
|----------|--------|---------|
| `{{TASK_ID}}` | task.ID | All 14 templates |
| `{{TASK_FILE}}` | feature.GetTaskFile() | All 14 templates |
| `{{SCOPE}}` | task.Scope | Most templates |
| `{{PHASE_SUMMARY}}` | PhaseDetect() | 10 templates |
| `{{FEATURE_SLUG}}` | SynthesizeOpts.FeatureSlug | 4 templates |
| `{{PROFILE}}` | task.Profile | 4 test pipeline templates |

#### Dimension 4: 跨文件无冗余 (150 pts)

**三类不同标准：**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 4a. 内容复制 | 0-60 | 相同/近似文本块出现在 3+ 文件。已知：(1) "Step 0: Resolve Profile" 9 个 SKILL.md ~8 行相同；(2) Eval Iron Laws + Steps 2-4 共 6 个 eval skill ~60 行相同；(3) Eval report 共享部分 5 个 report.md ~42 行相同。实例 = -10/个 |
| 4b. guide.md 与 SKILL.md 重叠 | 0-50 | guide.md 是 single source of truth。SKILL.md 复制 guide.md 内容（质量门禁序列、scope resolution）应改为引用。重复 = -10/个 |
| 4c. 不合理内联 | 0-40 | 内容有独立文件但也在 SKILL.md 内联。判断标准：agent 是否需要在单次上下文中看到完整内容（合理内联）vs 可通过 Read 工具按需获取（应引用）。不合理内联 = -10/个 |

**已知冗余实例：**

| Instance | Category | Files | Lines |
|----------|----------|-------|-------|
| "Step 0: Resolve Profile" | A | 9 SKILL.md files | ~90 total |
| Eval Iron Laws + Steps 2-4 | A | 6 eval SKILL.md | ~360 total |
| Eval report shared sections | A | 5 report.md | ~210 total |
| Quality gate sequence | B | guide.md, submit-task, fix-bug | ~12 |
| Scope resolution paraphrase | B | init-justfile SKILL.md | ~4 |

#### Dimension 5: 引用完整性 (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 5a. Agent references valid | 0-30 | `forge:<agent>` 或 `subagent_type` 指向 `plugins/forge/agents/` 中存在的文件。断链 = -15/个 |
| 5b. Template references valid | 0-25 | SKILL.md 中模板路径指向存在的文件。断链 = -15/个 |
| 5c. Cross-skill references valid | 0-25 | `invoke /<name>` 指向存在的 skill/command。断链 = -15/个 |
| 5d. Hook references valid | 0-20 | hooks.json 中路径/CLI 命令存在。断链 = -15/个 |

#### Dimension 6: 结构约定 (50 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 6a. Frontmatter completeness | 0-25 | SKILL.md 有 name + description。Command 有 name + description。Agent 有 name + description + model。缺失 = -5/个 |
| 6b. Eval template convention | 0-15 | eval-* 目录有 templates/rubric.md + report.md。缺失 = -10/个 |
| 6c. Name-directory alignment | 0-10 | skill name = 目录名，command name = 文件名。不匹配 = -5/个 |

### 2. Scorer 方法论：4-Phase 流程

替代当前的平面检查列表：

**Phase 1: 构建工作流图 (D1)**
1. 读取 rubric 内嵌的工作流规范
2. 扫描所有 skill/command/agent，提取实际的 prerequisites、outputs、gate points
3. 对比规范与实际，找出断点、死路、不可达状态

**Phase 2: 逐节点对抗测试 (D2)**
1. 列出每个节点的 gate/confirm 点
2. 对每个 gate，假设自己是懒 agent，提出"如何绕过？"
3. 检查 HARD-RULE 是否有可执行后果（非空话）
4. eval 循环：检查是否必须由独立 subagent 评分
5. quality gate：检查是否有 CLI 层面强制执行

**Phase 3: 逐文件精确性审查 (D3 + D4)**
1. 指令冲突（最高优先）：跨文件搜索同一概念的不同描述
2. 步骤歧义：每个 SKILL.md 步骤是否有唯一定义
3. 条件不完整：if-then 是否都有 else
4. 变量未定义：Agent 填充变量是否有来源说明
5. 内容冗余：按 A/B/C 三类标准检查

**Phase 4: 基础完整性 (D5 + D6)**
1. 引用完整性
2. Frontmatter、eval 模板、name 对齐

### 3. Reviser 两层修复

#### safe-fix（机械修复，不改语义）

- 修复 frontmatter 缺失字段
- 修正 name-directory 不匹配
- 修正 CLI flag/output field 名称
- 移除指向不存在文件的引用

#### guided-fix（有规则的修复）

**规则 1: 指令冲突 → guide.md 优先**
- guide.md 版本视为权威
- SKILL.md 改为引用："Follow the [Quality Gate Protocol](../hooks/guide.md)"
- guide.md 缺少该概念时，将最完整版本迁入 guide.md，其余引用

**规则 2: 内容复制 → 保留权威版本**
- guide.md 中已有 → SKILL.md 改为引用
- 无权威版本 → 保留最早/最完整的，其余改为引用
- eval loop protocol → 提取为 `references/shared/eval-loop-protocol.md`

**规则 3: 绕过漏洞 → 加最小 HARD-RULE**
- 添加具体后果描述："If you skip X, Y will fail because Z"
- 而非空话："You must not skip X"

### 4. SKILL.md 更新

Final Report 维度表从 12 维更新为 6 维。Parameters 保持 `--target 950 / --iterations 3` 不变。

## 预期效果

| Metric | Before | After (预期) |
|--------|--------|-------------|
| 维度数 | 12 | 6 |
| 首次评分 | ~965/1000 | ~500-650/1000 |
| 绕过向量检测 | 0 | ~15-20 |
| 指令冲突检测 | 0 | ~5-10 |
| 工作流断点检测 | 0 | ~3-5 |
| 冗余检测 | 弱（-10 per instance，仅 D3 0-10 分） | 强（D4 150 分） |

## Files to Modify

| File | Change |
|------|--------|
| `.claude/skills/eval-forge/templates/rubric.md` | 完全重写 |
| `.claude/skills/eval-forge/templates/scorer-prompt.md` | 完全重写 |
| `.claude/skills/eval-forge/templates/reviser-prompt.md` | 重写 |
| `.claude/skills/eval-forge/templates/report.md` | 更新 scorecard |
| `.claude/skills/eval-forge/SKILL.md` | 更新维度表 + Final Report |
