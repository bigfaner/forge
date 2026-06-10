---
name: forge-skill-audit
created: 2026-06-10
intent: "refactor"
status: proposed
---

# Forge Plugin Skill 系统一致性修复

## Problem

Forge Plugin v3.0.0-rc.53（21 个 skill、16 个命令、5 个 hook）经历大规模重构后，通过八维度深度审计（内部自洽性、Pipeline 契合度、提示词模板、任务模板、Surface 系统、Eval 系统、Hook 系统、Config 系统）发现了 23 处不一致问题，其中 4 个为 HIGH 级别静默错误——不会导致崩溃但会在用户执行 pipeline 时产出错误结果：

1. **错误的目标分数**（H-1）: `rubric-reference.md` 记录 journey scale=1000/target=850，实际 rubric 为 scale=1150/target=975。`eval-journey` 命令的 `argument-hint` 显示 `--target 850`，用户手动传入后将用 850 分去评判一个 1150 满分的 rubric——通过率被系统性高估
2. **死路径误导**（H-2）: `tech-design/SKILL.md` 引用从未被任何 skill 创建的 `docs/features/<slug>/proposal.md`，LLM 浪费上下文查找不存在的文件
3. **模板硬编码**（H-3）: `breakdown-tasks/templates/task.md` 硬编码 `complexity: "medium"` 和 `type: "coding.feature"`，LLM 在疲劳上下文下可能直接采用模板值而非覆盖
4. **类型分类文档错误**（H-4）: `record-format-coding.md` 列出 `doc.fix`，但 CLI 将其归类为 `doc` category，导致记录格式文档与实际运行时行为矛盾

**紧迫性**: v3.0.0-rc.53 是 release candidate，计划在 v3.0.0 正式发布前完成修复。H-1 意味着当用户参照 `argument-hint` 手动传入 `--target 850` 时，`eval-journey` 会以 850/1150=73.9% 作为通过标准，而非正确的 975/1150=84.8%，通过门槛被降低了 11 个百分点（注意：不传 `--target` 时，eval skill 会从 rubric frontmatter 正确读取 975，仅手动传参用户受影响）。

## Proposed Solution

按优先级执行文本级修复（无代码变更），每项修复后运行回归验证：

1. **数据同步**（H-1）: 更新 `rubric-reference.md` 表格 + `eval-journey`/`eval-contract` 的 `argument-hint` 和 `description` 字段（补充缺失维度）+ `eval/SKILL.md` description 以准确反映 scale 支持。在 `rubric-reference.md` 头部添加维护注释，标明此文件为 rubric scale/target 的二级缓存
2. **路径修正**（H-2）: 搜索所有 skill 中对 `docs/features/` 的引用，确认无并存路径后移除死路径
3. **模板修正**（H-3）: 将 `breakdown-tasks/templates/task.md` 的硬编码值替换为 `{{COMPLEXITY}}`/`{{TYPE}}` 占位符，添加与 `quick-tasks` 一致的注释块
4. **文档修正**（H-4）: 从 `record-format-coding.md` 移除 `doc.fix`
5. **M 级修复**（M-1~M-7, M-9）: 统一 config 命名、统一 eval 实现方式、补充路径、清理孤儿文件等

## Alternatives & Industry Benchmarking

| 方案 | 参考 | 优势 | 劣势 |
|------|------|------|------|
| **手动修正 + 回归验证（推荐）** | — | 最小变更、无代码风险、可立即执行 | 未来 rubric 变更仍需手动同步多处 |
| Schema 驱动的自动验证 | [promptfoo](https://promptfoo.dev) 对 prompt template 做断言验证（类似本提案的 `{{PLACEHOLDER}}` 完整性检查）；[conftest](https://conftest.dev) 对 YAML/JSON 配置文件做策略验证（可映射到 rubric-reference.md 与 rubric frontmatter 的一致性校验）；[Pact](https://pact.io) 的 provider-consumer 契约测试模式（可映射到 skill 间的输入输出契约验证） | 将 rubric-reference 验证、模板占位符检查、跨 skill 契约验证转化为 CI 检查，防止复发 | 需要编写 schema 定义和验证规则，投入较大；当前 forge 项目无 CI pipeline 集成点 |
| 仅修复 HIGH 级别 | — | 最小投入 | MEDIUM 问题（命名不一致、孤儿文件）积累后维护成本上升 |
| 延迟处理 | — | 无当前投入 | RC 阶段的静默错误随正式发布扩散，修复成本增加 |

**选择理由**: 手动修正是当前阶段最务实的选择——4 个 HIGH 问题均为文本/配置修正，无代码变更，无向后兼容风险。长期（v3.1+）应引入 conftest 风格的 schema 验证，将 rubric scale/target 的多真相源问题转化为 CI 阻断检查；中期可在回归验证中固化 grep 命令为 justfile recipe。

## Scope

### In Scope

- 修复 23 处审计发现中的所有 HIGH（4 项）和 MEDIUM（8 项，含 L-4 升级为 M-9）问题
- 范围覆盖 skill markdown 文件和 command markdown 文件；当发现验证需要引用 CLI Go 代码逻辑时（如 H-4 的 `CategoryForType` 函数、H-1 的 Config Resolution 逻辑），仅作为证据引用，不修改 Go 代码
- M-1 的 config key 重命名仅修改 skill markdown 中的引用；Go config reader 的 alias 兼容标记为后续任务

### Out of Scope

- Go 代码逻辑变更
- 已用错误 target 通过的 eval 历史结果修正（标注为 out-of-scope：无法追溯外部用户的 eval 结果，仅修正源数据防止未来复发）
- 用户项目级别的文件假设
- 性能优化
- 新功能添加

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 修复引入新不一致 | 中 | 中 | 每项修复后运行回归验证 grep 命令 |
| eval 生态多真相源同步 | 高 | 高 | 在 rubric-reference.md 头部添加维护注释；长期引入 CI 检查 |
| INLINE 跨 skill 引用过时 | 中 | 中 | 为 4 处 INLINE 引用添加源文件版本号标记（如 `@ v3.0.0-rc.53`），便于 grep 检测过时引用 |
| M-1 config key 重命名破坏现有用户配置 | 低 | 中 | 重命名仅修改 skill markdown 文件（in scope）；alias 兼容需要修改 Go config reader，标记为 out-of-scope 后续任务 |
| M-1 部分迁移导致新 key 静默失效 | 中 | 高 | 如果 Go config reader 不支持 kebab-case，用户按新 key 配置后 config 读取失败回退到默认值；必须先验证 Go config reader 再决定是否执行 M-1 |
| 审计遗漏 | 低 | 中 | 回归验证覆盖全量交叉检查；审计方法为逐文件完整读取，遗漏概率低 |

**回滚计划**: 所有修复均为 markdown 文件修改，通过 `git revert` 即可完整回滚，无数据迁移风险。

---

## Audit Results

### Summary Statistics

| 维度 | 审计范围 | 发现数 | HIGH | MEDIUM | MINOR |
|------|---------|--------|------|--------|-------|
| A: 内部自洽性 | 21 skills | 5 | 0 | 2 | 3 |
| B: Pipeline 契合度 | 11 pipeline skills | 6 | 1 | 3 | 2 |
| C: 提示词模板 | 56 模板文件 | 4 | 0 | 2 | 2 |
| D: 任务模板 | 6 模板 + 6 record format | 4 | 2 | 0 | 2 |
| E: Surface 系统 | 5 surface types x 5 skills | 0 | 0 | 0 | 0 |
| F: Eval 系统 | 11 rubrics + 8 experts + 7 commands | 2 | 1 | 0 | 1 |
| G: Hook 系统 | 5 hooks | 1 | 0 | 0 | 1 |
| H: Config 系统 | 20+ config keys | 1 | 0 | 1 | 0 |
| **Total** | | **23** | **4** | **8** | **11** |

> L-4 升级为 M-9，原 M-8 降级为 L-10（设计合理，非缺陷）。行合计 = 列总计 = 23。

---

### HIGH Severity (4 项)

#### H-1: rubric-reference.md 数据过时

- **文件**: `skills/eval/rules/rubric-reference.md`
- **问题**: journey 和 contract 的 scale/target 值与实际 rubric frontmatter 不一致
  - journey: rubric-reference 记录 `scale=1000, target=850`，实际 rubric 为 `scale=1150, target=975`
  - contract: rubric-reference 记录 `scale=1000, target=850`，实际 rubric 为 `scale=1100, target=935`
- **影响**: 依赖 rubric-reference.md 的代码/LLM 会使用错误的目标分数，导致文档"通过"了不该通过的评估
- **同时影响**: `commands/eval-journey.md` 和 `commands/eval-contract.md` 的 `argument-hint` 显示 `--target 850`，与实际 target 975/935 不符；两个命令文件的 `description` 字段仍声称 "1000-point rubric"，与实际 scale 1150/1100 不符，会误导 LLM 执行上下文；`eval-journey.md` description 仅列出 6 个评估维度，遗漏第 7 维度 "Workflow Coverage"（150 分），LLM scorer 可能跳过该维度；`eval-contract.md` description 声称 "six-dimension" 但实际 rubric 有 8 个维度，遗漏 Anchor Integrity 和 Fixture Specification；`skills/eval/SKILL.md` 声称 "Supports 100-point and 1000-point scales" 但实际不存在 100 分制 rubric，且 journey/contract 的 scale 超出 1000（1150/1100）
- **系统性风险**: eval 生态存在多个真相源（rubric frontmatter、rubric-reference.md、命令 argument-hint、命令 description、config 键默认值），缺乏单一真相源机制，未来变更需同步至少 5 处；其中 config 键默认值同步点不在本修复范围内，标记为已知缺口
- **修复**: 更新 rubric-reference.md 表格、两个命令文件的 argument-hint 和 description 字段（补充 eval-journey 的 Workflow Coverage 维度、eval-contract 的完整 8 维度）；更新 eval/SKILL.md description 以准确反映支持的 scale；在 rubric-reference.md 头部添加维护注释，标明此文件为 rubric scale/target 的二级缓存，任何 rubric frontmatter 变更必须同步更新此文件

#### H-2: tech-design intent 读取路径误导

- **文件**: `skills/tech-design/SKILL.md` 第 47 行
- **问题**: 文档列出 `docs/features/<slug>/proposal.md` 作为第一个 intent 读取路径，但该路径从未被任何 skill 创建。bash 命令（第 65 行）只检查 `docs/proposals/<slug>/proposal.md`
- **影响**: LLM 可能优先查找不存在的路径，浪费时间或报错
- **修复**: 移除死路径 `docs/features/<slug>/proposal.md`，只保留正确的 `docs/proposals/<slug>/proposal.md`。注意：`docs/features/` 路径被 200+ 目录和多 skill 广泛使用（tech-design 的 Prerequisites 检查 `docs/features/<slug>/prd/prd-spec.md`），需先搜索确认 proposal 文件是否可能存在于该路径（如 quick pipeline 的特殊行为），再决定修复方案

#### H-3: breakdown-tasks task.md 模板硬编码 complexity 和 type

- **文件**: `skills/breakdown-tasks/templates/task.md` 第 6 行和第 11 行
- **问题**: `complexity: "medium"` 和 `type: "coding.feature"` 是硬编码值，而非占位符。SKILL.md 有完整的 complexity 判定规则和 type 分配规则，但模板不使用占位符
- **对比**: `quick-tasks/templates/task.md` 正确使用了 `{{COMPLEXITY}}` 和 `{{TYPE}}` 占位符
- **影响**: LLM 看到模板中的 `"medium"` 后，需要意识到这是硬编码并覆盖——这依赖于 LLM 的推理能力而非系统保证，在疲劳上下文下可能直接采用模板值
- **修复**: 改为 `complexity: "{{COMPLEXITY}}"` 和 `type: "{{TYPE}}"`（与 quick-tasks 一致），并在 frontmatter 开头添加与 `quick-tasks/templates/task.md` 相同格式的注释块：
  ```
  # Template placeholders:
  #   COMPLEXITY — low | medium | high (default: medium)
  #   TYPE — coding.feature | coding.enhancement | coding.cleanup | coding.refactor | coding.fix | doc | doc.consolidate | doc.drift (default: coding.feature)
  ```

#### H-4: record-format-coding.md 错误列出 doc.fix

- **文件**: `skills/submit-task/data/record-format-coding.md` 第 3 行
- **问题**: 列出 `doc.fix` 作为 coding category 的任务类型，但 CLI 的 `CategoryForType` 函数将 `doc.fix` 归类为 `doc` category（`strings.HasPrefix("doc")` 匹配）
- **影响**: 记录格式文档声称覆盖 `doc.fix`，但实际执行时 `doc.fix` 任务会使用 `record-format-doc.md`，造成混淆
- **修复**: 从 record-format-coding.md 移除 `doc.fix`，在 record-format-doc.md 中补充 `doc.fix` 覆盖（当前 `doc.fix` 不在任何 record-format 文件中）
- **脆弱性分析**: `code-quality.simplify` 同样存在类似问题——它是一个非 `coding.` 前缀的任务类型，但通过硬编码特殊规则映射到 coding category。如果未来有人修改 `CategoryForType` 匹配逻辑而不知道这个特殊常量，分类会静默断裂。建议将 `code-quality.simplify` 重命名为 `coding.simplify` 以消除特殊映射需求

---

### MEDIUM Severity (9 项)

#### M-1: auto.eval 配置键命名风格不一致

- **文件**: brainstorm 第 153 行, write-prd 第 336 行, ui-design 第 140 行, tech-design 第 332 行
- **问题**: `auto.eval.proposal`/`auto.eval.prd` 使用全小写，`auto.eval.uiDesign`/`auto.eval.techDesign` 使用 camelCase
- **修复**: 统一为 kebab-case（`auto.eval.ui-design`, `auto.eval.tech-design`）。注意：alias 兼容需要修改 Go config reader，超出本 proposal 的 scope（仅修改 markdown 文件），因此标记为后续任务

#### M-2: ui-design auto.eval 实现方式不一致

- **文件**: `skills/ui-design/SKILL.md` 第 140 行
- **问题**: 其他三个 eval-capable skill 使用 bash script 模板实现三路分支（确定性执行），ui-design 使用自然语言描述（依赖 LLM 解释，引入变异性）
- **修复**: 统一为 bash script 模板

#### M-3: breakdown-tasks intent 读取路径不明确

- **文件**: `skills/breakdown-tasks/SKILL.md` 第 161 行
- **问题**: 只说 "If `proposal.md` has `intent`"，未给出完整路径
- **修复**: 添加完整路径 `docs/proposals/<slug>/proposal.md`

#### M-4: test-guide 存在 2 个孤儿规则文件

- **文件**: `skills/test-guide/rules/draft-generation.md`（248 行）和 `skills/test-guide/rules/pattern-extraction.md`（50 行）
- **问题**: SKILL.md 和其他 rules 文件均未引用。内容完整，暗示曾是流程一部分后被重构移除但文件残留
- **修复**: 移至 `skills/test-guide/rules/_deprecated/` 目录（与 eval 系统惯例一致）

#### M-5: run-tests test-isolation.md 跨 skill 依赖未声明

- **文件**: `skills/run-tests/rules/test-isolation.md`
- **问题**: run-tests SKILL.md 未引用，但 gen-test-scripts 通过 `<!-- INLINE -->` 引用。删除会破坏 gen-test-scripts
- **修复**: 在 test-isolation.md 头部添加 `<!-- OWNER: run-tests | CONSUMERS: gen-test-scripts (INLINE) -->` 注释

#### M-6: proposal 模板 {{AUTHOR}} 占位符无赋值逻辑

- **文件**: `skills/brainstorm/templates/proposal.md`
- **问题**: `{{AUTHOR}}` 占位符在 SKILL.md 中没有显式赋值指导，LLM 可能使用不一致的值
- **修复**: 在 brainstorm SKILL.md Step 5 中添加赋值指导：`Set {{AUTHOR}} to git config user.name output, or ask user if not available.`

#### M-7: manifest slug 占位符命名不一致

- **文件**: `skills/write-prd/templates/manifest.md` 使用 `{{FEATURE_SLUG}}`，`skills/quick-tasks/templates/manifest-quick.md` 使用 `{{SLUG}}`
- **修复**: 统一为 `{{SLUG}}`

#### L-10: breakdown-tasks task-doc.md 缺少 {{SLUG}} 占位符

- **文件**: `skills/breakdown-tasks/templates/task-doc.md`
- **问题**: Reference Files 使用通用 `{{REFERENCE_FILES}}` 占位符
- **评估**: breakdown-tasks 有更多输入文档（PRD、design、ui-design），通用占位符比硬编码 proposal 路径更灵活，当前设计合理

#### M-9: INLINE 跨 skill 引用同步风险（4 处，由 L-4 升级）

- gen-journeys 内联 gen-contracts 的 journey-contract-model.md
- gen-contracts 内联 gen-journeys 的 Surface Detection 逻辑（与上一项形成双向依赖）
- gen-test-scripts 内联 run-tests 的 test-isolation.md
- init-justfile 内联 test-guide 的 test-type-model.md
- **风险升级理由**: 如果源文件更新但内联副本未同步，LLM 会基于过时的模型定义生成 journeys，而 gen-contracts 会基于最新模型定义生成 contracts，产生不会报错的语义间隙
- **修复**: 在每个 INLINE 引用处添加源文件版本号标记（如 `<!-- INLINE from skills/gen-contracts/rules/journey-contract-model.md @ v3.0.0-rc.53 -->`），便于 grep 检测过时引用

---

### MINOR Severity (11 项)

#### L-1: ui-design 多平台输出 vs tech-design 单路径检查

- ui-design 可输出 `ui-design-web.md`/`ui-design-tui.md`/`ui-design-mobile.md`，tech-design 只检查 `ui/ui-design.md`
- 影响低：tech-design 的 UI 读取是建议性而非阻断性

#### L-2: run-tests 缺少 fix-task 创建逻辑

- execute-task/run-tasks/submit-task 都有 fix-task 创建模板，run-tests 没有
- 可能是有意设计：测试运行器只报告不做修复

#### L-3: eval 命令配置键与 auto.eval 键命名差异

- `eval.design.target` vs `auto.eval.techDesign`，`eval.ui.target` vs `auto.eval.uiDesign`
- 用途不同（目标分数 vs 自动运行开关），差异可接受但可能造成用户困惑

#### L-5: eval _deprecated 文件残留

- `skills/eval/rules/_deprecated/freeform-injection.md` 已正确标记废弃，无活跃引用
- 可在未来清理中移除

#### L-6: clean-code summary.md 无 YAML frontmatter

- 作为输出格式模板可接受，非文档模板

#### L-7: breakdown-tasks task.md 有 User Stories 段，quick-tasks 没有

- 设计意图：quick-tasks 无 PRD 输入

#### L-8: quick-tasks 与 breakdown-tasks manifest 结构不同

- 设计意图：通过 `mode: quick` 标记区分

#### L-9: run-tests 引用项目级文件

- `docs/business-rules/error-reporting.md` 是假设用户项目已生成的文件，不影响执行逻辑

#### L-11: hooks/guide.md 缺少 forge config get/set 命令说明

- **文件**: `hooks/guide.md`
- **问题**: `forge config get/set` 是 skill 中最高频使用的配置读取机制（36+ 处引用，用于 eval auto-run 检查、feature 上下文、target/iterations 解析等），但 guide.md 中未列出该命令
- **影响**: 新 skill 开发者查看 guide.md 时无法了解 config 读取的标准模式
- **修复**: 在 guide.md 的 CLI 命令参考部分补充 `forge config get <key>` 和 `forge config set <key> <value>` 说明

#### L-12: tech-design/examples/ 存在 2 个死文件

- **文件**: `skills/tech-design/examples/ask-question.md` 和 `skills/tech-design/examples/exploration.md`
- **问题**: SKILL.md 和 rules/ 中均未引用这两个 example 文件
- **修复**: 移除或在 SKILL.md 中补充引用

---

### Verified Healthy Areas

| 维度 | 状态 | 说明 |
|------|------|------|
| E: Surface 系统 | 健康 | 5 种 surface type 在所有相关 skill 中全覆盖，Surface Key vs Type 概念区分正确，测试目录映射一致 |
| G: Hook 系统 | 健康 | 5 个 hook 生命周期完整，SessionStart/SubagentStart 注入内容一致，Stop hook 链路正确 |
| H: Config 系统 | 存在命名风格不一致 | 默认值行为一致，无同 key 不同名问题；`auto.eval.*` 键存在 kebab-case 与 camelCase 混用（见 M-1） |
| Intent 传播 | 健康 | 6 种 intent 值在所有读取 intent 的 skill 中全部正确处理 |
| Fix-type 派生 | 健康 | 6 个位置的 fix-task 创建逻辑完全一致 |
| Manifest 状态机 | 健康 | 各 skill 正确设置状态，无覆盖冲突 |
| 命令 vs Skill 对齐 | 健康 | eval 命令正确引用 eval skill，quick 命令正确编排三个 skill |

---

## Success Criteria

- [ ] H-1: rubric-reference.md 数据与实际 rubric frontmatter 完全一致
- [ ] H-1: eval-journey/eval-contract 的 argument-hint 和 description 字段均反映正确的 target/scale 值，description 列出完整维度
- [ ] H-1: eval/SKILL.md description 准确反映实际支持的 scale 范围
- [ ] H-2: tech-design SKILL.md 仅引用 `docs/proposals/<slug>/proposal.md`，无死路径
- [ ] H-3: breakdown-tasks 模板使用 `{{COMPLEXITY}}`/`{{TYPE}}` 占位符，包含与 quick-tasks 一致的注释块
- [ ] H-4a: record-format-coding.md 不再列出 doc.fix
- [ ] H-4b: record-format-doc.md 包含 doc.fix 覆盖
- [ ] M-1: auto.eval 配置键统一为 kebab-case（承诺执行：先验证 Go config reader 是否支持 kebab-case 查询；如不支持，M-1 不可单独执行，必须与 Go alias 兼容绑定为原子操作并推迟至 Go 代码变更窗口，但不跳过——仅在 markdown 侧标记 TODO 并创建跟踪 issue）
- [ ] M-2: ui-design auto.eval 使用 bash script 模板
- [ ] M-3: breakdown-tasks intent 读取路径明确为完整路径
- [ ] M-4: test-guide 孤儿文件移至 _deprecated/ 或在 SKILL.md 中补充引用
- [ ] M-5: test-isolation.md 头部添加 OWNER/CONSUMERS 注释
- [ ] M-6: brainstorm SKILL.md 包含 {{AUTHOR}} 赋值指导
- [ ] M-7: manifest slug 占位符统一为 {{SLUG}}
- [ ] M-9: 4 处 INLINE 引用均标注源文件版本号
- [ ] Config 系统审计结论与实际发现一致
- [ ] 回归验证通过：`grep -r "1000-point" plugins/forge/commands/` 无残留；`grep -r 'complexity: "medium"' plugins/forge/skills/breakdown-tasks/templates/` 无残留

## Proposed Fix Order

1. **H-1**: 更新 rubric-reference.md + eval-journey/eval-contract argument-hint + description 字段（含维度补全）+ eval/SKILL.md description（数据修正，无风险）
2. **H-3**: breakdown-tasks task.md 改用占位符并添加注释块（模板修改，低风险）
3. **H-4**: record-format-coding.md 移除 doc.fix（文档修正，无风险）
4. **H-2**: tech-design SKILL.md 移除死路径（修复前先搜索确认完整路径拓扑）
5. **M-9**: 为 4 处 INLINE 引用添加版本号标记
6. **M-1~M-7, M-9**: 按 M 编号逐一处理

## Regression Verification

每项修复完成后，执行以下回归验证：

1. **数据一致性**: `grep "1000-point" plugins/forge/commands/eval-journey.md plugins/forge/commands/eval-contract.md` 确认无残留的过时 scale 描述；`grep "target.*850" plugins/forge/commands/eval-journey.md plugins/forge/commands/eval-contract.md` 确认无残留的过时 target 值（注意：其他 eval 命令的 scale 确实是 1000，不应修改）
2. **模板占位符**: `grep -r 'complexity: "medium"' plugins/forge/skills/breakdown-tasks/templates/` 确认硬编码已替换
3. **跨 skill 引用**: `grep -r "INLINE" plugins/forge/skills/` 确认所有内联引用都有版本号标记
4. **record-format 覆盖**: `grep -r "doc.fix" plugins/forge/skills/submit-task/data/` 确认仅出现在 record-format-doc.md 中
5. **全量交叉验证**: 重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致
6. **端到端验证**: 实际运行一次 `eval-journey` 命令，确认 target 值从 rubric frontmatter 正确读取
