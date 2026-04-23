---
feature: "feat-log-decisions"
---

# User Stories: feat-log-decisions

## Story 1: 使用重命名后的 tech-design skill

**As a** Skill 使用者
**I want to** 通过 `/zcode:tech-design` 调用技术设计 skill，且 skill 名称与输出文件 `tech-design.md` 方向一致
**So that** 降低记忆成本，减少调用错误

**Acceptance Criteria:**
- Given skill 目录已重命名为 `tech-design/`
- When 用户输入 `/zcode:tech-design`
- Then skill 正常启动，完成设计后在 feature 目录下生成 `tech-design.md`

---

## Story 2: 在 tech-design 流程中归档关键决策

**As a** Skill 使用者
**I want to** 在技术设计审批后，可选地将关键决策归档到 `docs/decisions/` 目录
**So that** 后续可以跨 feature 追溯技术决策，无需逐个翻阅各 feature 的 tech-design.md

**Acceptance Criteria:**
- Given tech-design 文档已获用户批准，且 AI 识别到关键决策
- When 展示编号候选决策列表后
- Then 用户可选择编号归档（或 `all`/`none`），选中决策写入对应 `docs/decisions/<type>.md`

- Given tech-design 文档已获用户批准，且 AI 未识别到关键决策
- When 进入归档步骤
- Then 跳过归档，直接进入 manifest 更新，无需用户额外确认

---

## Story 3: 主动记录技术决策

**As a** Skill 使用者
**I want to** 在任意阶段通过 `/zcode:record-decision` 主动记录一条技术决策
**So that** 非设计阶段（实现、brainstorm、PRD）产生的重要决策也能集中归档

**Acceptance Criteria:**
- Given 用户调用 `/zcode:record-decision`
- When 通过 4 轮交互输入类型、描述、理由、feature
- Then 对应 `docs/decisions/<type>.md` 新增 1 条表格行，`manifest.md` 计数 +1 且 Recent Decisions 表更新

---

## Story 4: 重命名 skill 目录并更新所有引用

**As a** Skill 开发者
**I want to** 将 `design-tech/` 重命名为 `tech-design/` 并同步更新 SKILL.md frontmatter、hooks guide、exploration 示例和 CLAUDE.md 中的引用
**So that** 所有入口和文档一致指向新名称，用户通过 `/zcode:tech-design` 即可正常调用

**Acceptance Criteria:**
- Given skill 目录 `plugins/zcode/skills/design-tech/` 已重命名为 `tech-design/`
- When 检查 SKILL.md frontmatter
- Then `name` 字段为 `tech-design`，旧名称 `design-tech` 不再出现在任何 skill 配置文件中

- Given 重命名已完成
- When 检查 `plugins/zcode/hooks/guide.md` 和 `plugins/zcode/skills/tech-design/examples/exploration.md`
- Then 所有 `DECISIONS.md` 引用已替换为 `docs/decisions/`，无残留

- Given 重命名已完成
- When 检查 `zcode/CLAUDE.md`
- Then skill 索引中包含 `tech-design` 和 `record-decision`，不包含 `design-tech`

---

## Story 5: 创建共享的 decision-logging reference 和 decision-entry template

**As a** Skill 开发者
**I want to** 创建 `references/decision-logging.md` 和 `templates/decision-entry.md`，由 tech-design 和 record-decision 两个 skill 共同引用
**So that** 决策提取和记录逻辑只维护一份，避免两个 skill 之间逻辑重复

**Acceptance Criteria:**
- Given `references/decision-logging.md` 已创建
- When tech-design 和 record-decision 各自执行决策归档/记录
- Then 两者引用同一份 reference 文件，提取逻辑无重复

- Given `templates/decision-entry.md` 已创建
- When 新增决策被写入类型文件
- Then 决策行格式与 template 一致（Date, Feature, Decision, Rationale, Source 五字段）

---

## Story 6: 查看决策索引

**As a** Skill 使用者
**I want to** 查看 `docs/decisions/manifest.md` 了解所有已归档决策的全局视图
**So that** 快速定位某类决策的历史记录，判断是否已有相关决策可复用

**Acceptance Criteria:**
- Given 至少有一条决策已归档
- When 打开 `docs/decisions/manifest.md`
- Then Categories 表显示每个类型的决策计数和最近更新时间，Recent Decisions 表显示最近 5 条决策

---
