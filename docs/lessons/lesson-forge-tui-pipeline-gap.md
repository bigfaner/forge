# Lesson: Forge Pipeline 对 TUI Feature 的防护缺失 — 系统性诊断与改进计划

## 背景

deep-drill-analytics feature 通过完整 Forge pipeline（PRD → design → tasks → e2e tests），毕业后产生 23 个 vibe coding 提交（11 fix + 5 style + 3 feat + 4 refactor/chore/docs）。事后提取了 5 个 convention 文件和 3 个 lesson 文件。

两轮 sub-agent 深度评估揭示了核心问题：

> **Conventions 写对了，但 Forge pipeline 不加载、不执行。下一个 TUI feature 将重复相同的 vibe coding 循环。**

## 系统性诊断

### 问题链条

```
Conventions 提取了 → 但 Forge skill 不读取 → pipeline 各阶段无 TUI 意识
                                                    ↓
/tech-design 不要求 mockup → /eval-design 不检查视觉规格 → /breakdown-tasks 不注入 TUI verify
                                                    ↓
Agent 凭猜测实现 → 运行时暴露 bug → vibe coding 循环
```

### Pipeline 逐 stage 演练（当前 vs 应有）

| Stage | 当前状态 | TUI 防护 | 影响 |
|-------|---------|----------|------|
| `/brainstorm` | 无 TUI 意识 | ASPIRATIONAL | 不问视觉复杂度 |
| `/write-prd` | 无 TUI flag | ASPIRATIONAL | PRD 不标注"需要视觉设计" |
| **`/tech-design`** | 列出 conventions 作为可选探索源，无硬 gate | **ASPIRATIONAL** | **无 mockup、无字符调色板、无边界场景** |
| **`/eval-design`** | rubric 6 维度全盲 TUI | **ASPIRATIONAL** | **视觉缺陷不被扣分** |
| **`/breakdown-tasks`** | verify criteria 纯功能 | **ASPIRATIONAL** | **无 golden test、无维度检查** |
| `/execute-task` | convention 读取弱（keyword 映射不含 TUI） | PARTIAL | 可能读到，但不保证 |
| Vibe coding 阶段 | 无 scope guard | ASPIRATIONAL | feat 提交不走 pipeline |

### 防护效果评估

| 原始 23 提交的根因 | Convention/Lesson 存在 | Skill 强制执行？ | 能阻止？ |
|-------------------|----------------------|----------------|---------|
| 无视觉规格 → 5 个 style 迭代 | lesson-tui-tech-design-mockup | No | No |
| 无边界场景 → 11 个 fix | lesson-tui-visual-verify | No | No |
| 无字符调色板 → 3 次字符迭代 | lesson-tui-tech-design-mockup | No | No |
| 无 golden 维度检查 → bug 不被捕获 | lesson-tui-visual-verify | No | No |
| scope creep → 3 个 feat 提交 | lesson-vibe-coding-scope-control | No | No |
| len() vs runewidth | tui-dynamic-content §3 | 弱（skill 可选读取） | Maybe |
| scrollbar 宽度 | tui-layout-ui §1, tui-dynamic-content §2 | 弱 | Maybe |
| 内容溢出 | tui-dynamic-content §5 | 弱 | Maybe |

**总计: 0/11 根因完全阻止，3/11 弱缓解。**

### 根因的根因

三个 lesson 的 "How to apply / 后续" 部分共提出 **9 条 Forge skill 修改建议，0 条已实施**：

1. `/tech-design` prompt 加入 TUI mockup 要求 → **未改**
2. `/eval-design` rubric 加入 Visual Specification 维度 → **未改**
3. `/breakdown-tasks` 注入 TUI verify 模板 → **未改**
4. `/execute-task` quality gate 增加 golden 维度检查 → **未改**
5. `/eval-design` 检查 boundary scenarios → **未改**
6. task-executor 增加 TUI keyword → convention 映射 → **未改**
7. `/git-commit` 增加 scope guard → **未改**
8. task-executor Step 1 增加 convention 读取摘要 → **未改**
9. `/write-prd` 增加 TUI complexity flag → **未改**

这些 lesson 本质上是**愿望清单**，不是操作文档。

## 改进计划

按优先级排列。P0 阻止 60-70% 的 vibe coding，P1 再阻止 20-30%。

### P0: 必须在下个 TUI feature 前完成

#### P0-1: `/tech-design` SKILL.md — 加入 TUI mockup 硬 gate

**文件**: `skills/tech-design/SKILL.md`

在 Step 3 (Design Template) 中，当探索上下文检测到 TUI 相关文件（`internal/model/*.go` 含 `View()` 函数，或 bubbletea/lipgloss import），要求设计文档必须包含：

- **ASCII layout mockup** — 每个面板一个，用 box-drawing 字符
- **精确尺寸数值** — 不允许"大约"/"适当"
- **5 个强制边界场景** — 窄终端(80×24)、宽终端(140+)、混合数字宽度、长路径(>50 chars)、无数据
- **字符调色板** — 每个视觉元素指定 Unicode 字符 + code point + 选择理由
- **颜色映射** — 从 `docs/conventions/tui-layout-ui.md` 色表选取

详细规格见 → `lesson-tui-tech-design-mockup.md`

#### P0-2: `/eval-design` rubric.md — 新增 Visual Specification 维度

**文件**: `skills/eval-design/templates/rubric.md`

新增第 7 维度 "Visual Specification (15 pts)"，当 design 引用 TUI 渲染时激活：

| 检查项 | 分值 |
|--------|------|
| Mockup present per panel | 0-4 pts |
| Dimensions are numerical, not vague | 0-4 pts |
| 5 boundary scenarios covered | 0-4 pts |
| Character palette complete | 0-3 pts |

同时将 Architecture Clarity 从 20 调整为 15，总分保持 100。

#### P0-3: `/breakdown-tasks` SKILL.md — TUI verify 模板注入

**文件**: `skills/breakdown-tasks/SKILL.md`

在 Step 4a (Generate Tasks)，当 task 的 affected files 包含 `internal/model/*.go` 且 scope 含 `frontend`，自动在 verify criteria 追加：

```markdown
### TUI Rendering
- [ ] Golden test exists for new/modified View()/Render() function
- [ ] Dimension check: lines == height, each line width <= terminal width
- [ ] Test data includes: CJK string, long path (>50 chars), multi-digit number (>9), empty field
- [ ] No hardcoded widths — all derived from m.width
- [ ] Colors from palette only (docs/conventions/tui-layout-ui.md)
```

### P1: 应在近期完成

#### P1-1: task-executor / execute-task — 补充 TUI keyword 映射

**文件**: `agents/task-executor.md`, `commands/execute-task.md`

在 convention 加载的 keyword 示例映射中补充：

```
"tui"/"panel"/"view"/"render" → tui-layout-ui.md, tui-dynamic-content.md
"lipgloss"/"bubbletea"        → lipgloss-panel-width.md
"scrollbar"/"truncation"      → tui-dynamic-content.md
"keyboard"/"navigation"       → tui-ux-interaction.md
"parser"/"jsonl"              → tui-data-contracts.md
```

#### P1-2: task-executor Step 1 — Convention 读取确认

在 convention 读取后，agent 输出一行摘要：
```
Loaded conventions: tui-layout-ui (§1,3,7), tui-dynamic-content (§2,3,5)
```
使 convention 读取可审计，增加 agent 实际应用的概率。

#### P1-3: `/quick-tasks` SKILL.md — TUI 快速模式 guard

Quick mode 跳过 design，对 TUI feature 危险性更高。当 proposal 涉及 TUI 渲染时：
- 在 task 中标注 "Quick mode may be insufficient for TUI features"
- 在 task verify criteria 中附加 TUI Rendering 模板

### P2: 锦上添花

#### P2-1: `/git-commit` SKILL.md — Scope guard

当 `git diff --stat HEAD` 显示单个 .go 文件 +50 行非 test 代码，提示:
"This appears to be new functionality. Is this within current task scope? Consider creating a task."

#### P2-2: `/write-prd` SKILL.md — TUI complexity flag

在 PRD Self-Check (Step 9.5) 增加:
"If feature involves terminal rendering or TUI, note that tech design requires visual specification."

## 实施路线

```
Phase 1 (P0, 下个 TUI feature 前):
  ├── P0-1: /tech-design SKILL.md 加 TUI mockup gate
  ├── P0-2: /eval-design rubric.md 加 Visual Spec 维度
  └── P0-3: /breakdown-tasks SKILL.md 加 TUI verify 注入

Phase 2 (P1, 1-2 周内):
  ├── P1-1: task-executor keyword 映射
  ├── P1-2: convention 读取确认
  └── P1-3: /quick-tasks TUI guard

Phase 3 (P2, 按需):
  ├── P2-1: /git-commit scope guard
  └── P2-2: /write-prd TUI flag
```

## 子 Lesson 索引

| Lesson | 内容 | 关联 Forge 改动 |
|--------|------|----------------|
| [lesson-tui-tech-design-mockup](lesson-tui-tech-design-mockup.md) | ASCII mockup + 字符调色板 + 边界场景的详细规格 | P0-1, P0-2 |
| [lesson-tui-visual-verify](lesson-tui-visual-verify.md) | Golden test + 维度检查 + 测试数据真实性的详细规格 | P0-3, P1-1 |
| [lesson-vibe-coding-scope-control](lesson-vibe-coding-scope-control.md) | Scope guard + 反思阈值 + parser 层保护 | P2-1, P1-2 |

## 验证方法

下个 TUI feature 开始前：
1. 确认 P0-1/2/3 已修改对应 skill 文件
2. 走一遍 mental walkthrough: `/tech-design` → mockup 生成了吗？→ `/eval-design` → 视觉维度打分了吗？→ `/breakdown-tasks` → TUI verify 注入了？
3. 如果任一步 answer 为 No，则 vibe coding 循环将重复

Feature 完成后，对比 vibe coding 提交数。目标：从 23 个降到 ≤8 个。

## Tags

`forge-improvement`, `tui`, `pipeline`, `meta-lesson`, `process`, `scope-control`
