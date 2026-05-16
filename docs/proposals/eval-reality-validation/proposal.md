---
created: 2026-05-16
author: "fanhuifeng"
status: Draft
---

# Proposal: eval-* Problem Detection Enhancement

## Problem

eval-* 系统存在结构性盲区：它只评审文档的内在质量（格式、逻辑、结构），无法验证文档内容与项目现实的吻合度。这导致一个问题在所有 eval 步骤都通过后，仍然在实现后或用户使用时才被发现。

### Three Detection Levels

| Level | Description | eval-* 能力 |
|-------|-------------|-------------|
| A. 格式/结构 | 文档缺失章节、ID 不一致、引用断裂 | 强：模式匹配即可解决 |
| B. 内容质量 | 需求模糊、设计有逻辑漏洞、测试覆盖不足 | 中：依赖 rubric 维度覆盖和 LLM 推理 |
| C. 正确性 | 设计方案在实现时必然失败、业务规则与代码矛盾 | 弱：缺乏项目现实参照 |

当前 eval 系统在 A/B 上有效，在 C 上有结构性盲区。

### Evidence

- eval-prd、eval-design、eval-test-cases 全部通过的功能，完成后仍经常发现功能不完善（交互不友好、功能 Bug）
- 根因不是某个 eval 步骤做得不好，而是 **整条流水线都只看文档**——PRD 遗漏了真实场景，但所有 eval 只检查"这份不完整的 PRD 写得好不好"
- traceability 体系放大了这个问题：test-cases 100% 回溯到 PRD，但 PRD 可能只覆盖了 70% 的真实场景——30% 的场景从未进入流水线

### Urgency

每个 feature 都经过 eval-* 流程，但 eval 通过不等于功能正确。随着 feature 数量增加，返工成本累积。当前是修复 eval 盲区的最佳时机——spec 文件少、eval 架构刚稳定。

## Proposed Solution

三道防线：前置拦截（eval-* 增强）+ Layer 1 静态追踪 + Layer 2 动态验证。

### 防线 1：eval-* 上下文注入 + rubric 维度增强

**上下文注入**：给 scorer agent 注入项目级 conventions 和 business-rules，使其能检测文档内容与项目现实的矛盾。

当前 scorer 只看被评估文档 + rubric，是纯文档内评审。注入上下文后，scorer 能发现：
- PRD 违反已有 business-rules
- 设计使用了项目中不存在的依赖
- 测试用例遗漏了 conventions 要求的覆盖

**Rubric 维度增强**：在现有 rubric 中增加场景完整性相关维度。

### 防线 2：Layer 1 — 静态代码追踪（validate-code）

在 E2E 测试之前，拿着 PRD 的每个用户场景，验证实现代码中是否存在完整路径。

映射链：PRD 场景 → task records（粗粒度）→ git diff（细粒度）→ 定向代码阅读 → 验证路径完整性

### 防线 3：Layer 2 — 动态运行验证（validate-ux）

在 quality-gate 中，对运行中的系统执行两项检查：
1. UX rubric 规则检查（覆盖 CLI/Web/TUI 三类项目）
2. PRD 用户流程实际走通 + 操作影响验证

详见 [validate-ux 详细设计](#validate-ux-详细设计)。

### Architecture

所有改动统一在现有 `forge:eval` skill 内，通过 `--type` 参数区分：

| --type | 输入 | scorer 行为 | iterations |
|--------|------|-------------|------------|
| `prd` / `design` / `ui` / ... | 文档 + 注入上下文 | 文档打分 + 上下文矛盾检测 | 3 |
| `validate-code` | PRD + git diff + 代码 | 场景追踪（走通/走不通） | 1 |
| `validate-ux` | PRD + ux-snapshot.md | UX 规则检查 + 流程走通 + 影响验证 | 1 |

validate-* 使用 iterations=1，跳过 revise 循环（产出的是问题报告，不是修订后的文档）。

### Context Injection Design

Rubric frontmatter 声明所需上下文类别：

```yaml
context:
  conventions: [api, naming, ux]  # 按类别筛选
  business-rules: auto            # 按 PRD 内容自动匹配领域
```

- `docs/conventions/` 按 rubric 声明的类别筛选，不全量加载
- `docs/business-rules/` 按 PRD 涉及的领域匹配
- 不参考同类 feature PRD（成本高、ROI 低）

## validate-ux 详细设计

### 执行模型：两阶段

Phase 1（Pre-processing，main session）：编译运行系统，采集 ux-snapshot.md
Phase 2（Score，doc-scorer）：评估 ux-snapshot.md + rubric

### 项目类型适配

| 类型 | 运行方式 | 操作单元 | 采集内容 | TUI 特殊处理 |
|------|----------|----------|----------|-------------|
| CLI | Bash 执行命令 | Shell 命令 | stdout/stderr/exit code | — |
| Web | agent-browser | URL + element selector + action | 截图 + accessibility tree + 网络日志 | — |
| TUI | Bash 执行 | stdin pipe（按键序列） | 终端输出 | 第一版仅覆盖非交互式场景 |

### PRD 到操作序列的翻译（策略 3：混合翻译）

三种类型共用同一翻译策略：

1. **直接提取**：扫描 PRD 中的代码块、命令、URL、按键描述
2. **推断补全**：缺失具体操作的，agent 基于辅助信息推断

各类型的推断辅助信息：

| 类型 | 辅助信息来源 | 推断方式 |
|------|-------------|---------|
| CLI | `forge --help` 递归获取子命令 help | 匹配 PRD 描述 → 子命令 → 参数格式 |
| Web | sitemap.json（accessibility tree + element IDs） | 匹配 PRD 描述 → 路由 → DOM selector |
| TUI | 运行程序 capture 初始屏幕 + help 输出 | 匹配 PRD 描述 → 菜单选项 → key-binding |

### ux-snapshot.md 格式

```markdown
# UX Snapshot: <feature-name>

## Project Info
- Type: cli | web | tui
- Binary/URL: <path or url>
- PRD Reference: <path to PRD>
- Generated: <timestamp>

## Flow: <flow-name-from-PRD>

### Step 1: <action-description>
**Command/Navigate**: <what was executed>
**Input**: <what was sent>
**Output**:
```
<raw stdout/stderr, screenshot path, or terminal capture>
```
**Exit Code**: <cli only>

**Effect Verification**:
- Data: <expected data change> → <actual result> ✓/✗
- Side Effect: <unexpected changes checked via git diff --stat> → ✓/✗
- Output Consistency: <output claim vs reality> → ✓/✗
- Cascade: <downstream behavior triggered?> → ✓/✗

**Idempotency Check**:
- Re-run: <result of running same command again>

**State Integrity**:
- <consistency check between related state>

### Step 2: ...

## Standalone Checks

### Help Text
**Command**: `<binary> --help`
**Output**:
```
<full help output>
```

### Error Handling
**Command**: `<binary> invalid-command`
**Output**:
```
<error output>
```

### Version Info
**Command**: `<binary> --version`
**Output**:
```
<version output>
```
```

### 操作影响验证（7 类）

| 影响类型 | 验证方式 | 示例 |
|----------|---------|------|
| Data Effect | 操作前后对比文件/数据库/状态 | `submit` 后 index.json 状态更新 ✓ |
| Side Effect | `git diff --stat` 检查非预期文件变更 | `delete task` 未影响相邻 task ✓ |
| Idempotency | 重复执行同一操作 | `submit` 跑第二次返回 "already submitted" ✓ |
| Output-Reality Consistency | 验证输出信息与实际状态 | 输出 "created: X.md" → 文件确实存在 ✓ |
| State Integrity | 多步操作后检查系统整体一致性 | record 文件数 = index.json 计数 ✓ |
| Cascade Effect | 检查下游行为是否触发 | `submit` 后 quality-gate 被触发 ✓ |
| Rollback Feasibility | 操作失败后检查状态可恢复性 | 失败后无残留脏状态 ✓ |

### Rubric 维度（1000 分制）

```yaml
scale: 1000
target: 700
iterations: 1
type: validate-ux
context:
  conventions: [ux, cli, api]  # 根据 project-type 动态调整
  business-rules: auto
```

| # | 维度 | 分值 | 核心评判标准 |
|---|------|------|-------------|
| **板块 A：UX 规则检查** | | **400 分** | |
| 1 | Error Actionability | 120 | 错误信息是否包含可操作的修复建议 |
| 2 | Help Completeness | 120 | help 是否覆盖所有可用操作、参数有默认值标注 |
| 3 | Output Clarity | 90 | 输出是否可读、结构化、重点突出 |
| 4 | Platform UX Rules | 70 | CLI: exit code/pipe/progress; Web: loading/form/nav; TUI: key-binding/layout/status |
| **板块 B：PRD 流程走通 + 影响验证** | | **600 分** | |
| 5 | Flow Completeness | 120 | PRD 每个流程是否从第一步到最后一步都能走通 |
| 6 | Output-Reality Consistency | 120 | 输出信息与实际状态是否一致 |
| 7 | Data & Side Effect | 120 | 预期数据变化是否发生；非预期变化是否发生 |
| 8 | Idempotency & State Integrity | 100 | 重复执行是否安全；多步操作后状态是否自洽 |
| 9 | Cascade Effect | 60 | 下游行为是否被正确触发 |
| 10 | Friction Detection | 80 | 流程中是否有需要查文档才能继续的步骤 |

### TUI 限制

第一版 TUI 仅覆盖非交互式场景（初始渲染、help 屏幕、无效输入响应）。交互式流程验证需要 raw terminal 模拟（如 `expect` 脚本），业界尚无成熟 headless TUI 方案，作为后续增强。

### Pre-processing 执行流程

```
1. 读 PRD → 提取用户流程列表
2. 检测项目类型（CLI/Web/TUI）
3. 扫描 PRD 中的直接操作引用（代码块、命令、URL、按键描述）
4. 加载项目辅助信息：
   CLI:  递归获取所有子命令 help
   Web:  sitemap.json + agent-browser 基础扫描
   TUI:  运行程序 capture 初始屏幕 + help
5. 对缺失的具体操作，agent 基于 PRD 描述 + 辅助信息推断
6. 生成类型特定的操作序列
7. 逐步骤执行操作序列，采集 output
8. 每步执行 effect verification：
   a. Data Effect: 对比操作前后的文件/状态
   b. Side Effect: git diff --stat
   c. Idempotency: 重复执行关键操作
   d. State Integrity: 多步操作后一致性检查
   e. Cascade: 检查下游触发
9. 写入 ux-snapshot.md
10. Spawn doc-scorer 评估 ux-snapshot.md
```

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **扩展 eval skill（本方案）** | 复用 scorer→gate→revise 循环，架构统一 | eval 复杂度增加 | **Chosen** |
| 新建独立 validate skill | 职责分离 | skill 膨胀，validate 逻辑与 eval 高度重叠 | Rejected |
| 只做后置验证，不改 eval-* | 改动范围小 | 不在前端拦截，返工成本高 | Rejected |
| 只增强 eval-* rubric | 改动最小 | 无法覆盖实现后才发现的问题 | Rejected |
| 全量加载 conventions | 实现简单 | 噪音淹没 scorer 注意力 | Rejected |
| validate-ux 用 general-purpose agent 直接执行评估 | 执行上下文完整 | 不可控（错误命令、状态污染、无限循环） | Rejected |
| TUI 交互式验证（expect/pipe） | 覆盖完整 | raw terminal mode 不兼容，flaky | Deferred: v2 |

## Feasibility Assessment

### Technical Feasibility

完全可行。eval skill 已有多类型参数化（prd/design/ui/harness/consistency/test-cases），增加 validate-code 和 validate-ux 是已有模式的自然扩展。上下文注入只在 pre-processing 层改动，不影响 scorer/gate 核心逻辑。

validate-ux 的两阶段模型保持 doc-scorer 的文档评估范式不变——scorer 评估的是 ux-snapshot.md 这个"文档"，不是运行中的系统。复杂度集中在 pre-processing 阶段。

### Resource & Timeline

5 个批次，预计 8-12 sessions（validate-ux 复杂度高于初始估算）。

## Scope

### In Scope

**Batch 1 — 基础设施**
- rubric frontmatter 增加 `context` 字段规范
- eval skill pre-processing 支持读取 `context` 声明，按类别筛选 conventions/business-rules 并注入 scorer prompt

**Batch 2 — eval-prd 增强**
- eval-prd rubric 增加 "Scenario Completeness" 和 "Edge Case Coverage" 维度
- eval-prd rubric frontmatter 声明所需上下文类别
- eval-design、eval-test-cases 等后续按需补充

**Batch 3 — validate-code**
- 新增 `skills/eval/rubrics/validate-code.md` rubric
- eval skill 增加 validate-code pre-processing（组装 PRD + git diff + 代码文件列表）
- doc-scorer prompt 适配场景追踪模式
- 在 quick-tasks 和 breakdown-tasks 中增加对应 task 模板

**Batch 4 — validate-ux**
- 新增 `skills/eval/rubrics/validate-ux.md` rubric（1000 分制，10 维度）
- eval skill 增加 validate-ux pre-processing（两阶段：采集 ux-snapshot.md → doc-scorer 评估）
- 支持 CLI/Web/TUI 三种项目类型的项目类型检测和操作翻译
- 7 类操作影响验证（Data/Side/Idempotency/Consistency/Integrity/Cascade/Rollback）
- 在 quality-gate 中集成

**Batch 5 — 剩余 rubric 增强**
- eval-design、eval-test-cases、eval-ui 等 rubric 维度补充
- 对应 rubric frontmatter context 声明

### Out of Scope

- 新建独立 skill（validate 不脱离 eval）
- 修改 doc-scorer.md 或 doc-reviser.md agent 定义
- TUI 交互式流程验证（raw terminal mode，v2 增强）
- Rollback Feasibility 自动验证（需要 pre/post snapshot，复杂度高，v2 增强）
- 修改现有 eval 命令的 CLI 接口

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 上下文注入增加 scorer token 消耗 | High | Low | 按类别筛选控制量；conventions 通常不超过 5K tokens |
| 上下文噪声稀释 scorer 注意力 | Medium | Medium | 严格按 rubric 声明的类别筛选，不加载无关内容 |
| validate-code 场景追踪误报（代码实现了但 agent 找不到路径） | Medium | Medium | 用 git diff 缩小搜索范围；允许 agent 逐步扩大搜索 |
| validate-ux pre-processing 执行命令破坏项目状态 | Medium | High | 在 git worktree 中执行；操作前 stash/commit；所有写入操作在临时目录 |
| PRD 到操作序列翻译不准确 | Medium | Medium | 混合策略优先使用 PRD 直接引用；推断结果需在 snapshot 中展示供 scorer 审查 |
| Web 项目 agent-browser 采集成本高 | Medium | Low | 截图 + a11y tree 组合足够 LLM 判断；不做视频录制 |
| TUI 非交互式覆盖不足 | High | Low | 大部分 TUI UX 问题在初始屏幕可见；交互式留 v2 |
| 5 个批次跨度大，中途需求变化 | Medium | Medium | 每批独立可交付，不依赖后续批次 |

## Success Criteria

- [ ] rubric frontmatter 支持 `context` 字段，eval pre-processing 按声明注入筛选后的 conventions 和 business-rules
- [ ] eval-prd rubric 包含 "Scenario Completeness" 维度，scorer 能引用注入的 business-rules 发现矛盾
- [ ] `forge:eval --type validate-code` 可执行：输入 PRD + feature branch git diff，输出每个 PRD 场景的代码追踪报告（走通/走不通/部分实现）
- [ ] `forge:eval --type validate-ux` 可执行：编译安装后运行系统，输出 UX rubric 评分 + PRD 流程走通结果 + 7 类影响验证报告
- [ ] validate-ux 覆盖 CLI/Web/TUI 三种项目类型，自动检测并适配
- [ ] ux-snapshot.md 包含完整的 flow 步骤、standalone checks、effect verification
- [ ] 所有改动在现有 `forge:eval` skill 内，无新建 skill
- [ ] validate-* 使用 iterations=1，不触发 revise 循环
- [ ] 不修改 doc-scorer.md 或 doc-reviser.md

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
