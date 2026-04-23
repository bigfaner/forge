---
feature: "feat-log-decisions"
---

# User Stories: Tech Design Skill 改进：重命名与决策归档

## Story 1: 统一 skill 命名

**As a** 使用 zcode 工具链的开发者
**I want to** 通过 `/zcode:tech-design` 调用技术设计 skill
**So that** skill 名称与输出文件名方向一致，降低认知负担

**Acceptance Criteria:**
- Given 我在项目中执行 `/zcode:tech-design`
- When skill 被调用
- Then skill 正常执行，输出 `tech-design.md`，且 `/zcode:design-tech` 不再可用

---

## Story 2: 在 tech-design 流程中归档关键决策

**As a** 使用 zcode 工具链的开发者
**I want to** 在批准 tech-design 文档后，选择性地将关键决策归档到 `docs/decisions/`
**So that** 跨 feature 的技术决策可以集中追溯，不再需要逐一翻阅各 feature 的 tech-design.md

**Acceptance Criteria:**
- Given tech-design 文档已获我批准，且文档中存在被标记的关键决策
- When skill 展示编号候选决策列表
- Then 我可以输入编号（逗号分隔）、`all` 或 `none` 来选择归档范围
- Given 我选择了若干决策
- When 归档执行
- Then 对应 `docs/decisions/<type>.md` 新增表格行，`docs/decisions/manifest.md` 计数和 Recent Decisions 表同步更新
- Given tech-design 文档中没有关键决策
- When skill 完成审批后
- Then 跳过归档步骤，直接进入 manifest 更新

---

## Story 3: 独立记录技术决策

**As a** 使用 zcode 工具链的开发者
**I want to** 在任意阶段通过 `/zcode:record-decision` 主动记录一条技术决策
**So that** 不必等到 tech-design 阶段才能归档决策，brainstorm 或实现过程中的重要决策也能及时保存

**Acceptance Criteria:**
- Given 我执行 `/zcode:record-decision`
- When skill 通过 4 轮 AskUserQuestion 收集信息（类型、描述、理由、关联 feature）
- Then 对应 `docs/decisions/<type>.md` 新增 1 条表格行，包含 Date、Feature、Decision、Rationale、Source 五个字段
- Given 归档完成
- When 我查看 `docs/decisions/manifest.md`
- Then 对应类型的 Decisions 计数 +1，Recent Decisions 表包含该条记录

---

## Story 4: 查阅历史决策

**As a** 使用 zcode 工具链的开发者
**I want to** 打开 `docs/decisions/<type>.md` 查看该类型的所有历史决策
**So that** 在 30 秒内定位到相关历史决策，避免重复踩坑或做出矛盾决策

**Acceptance Criteria:**
- Given `docs/decisions/` 目录已初始化，且至少有 1 条归档记录
- When 我打开对应类型文件（如 `architecture.md`）
- Then 文件以表格形式展示所有该类型决策，包含 Date、Feature、Decision、Rationale、Source
- Given 我查看 `docs/decisions/manifest.md`
- When 打开文件
- Then 可以看到 8 个类型的决策数量汇总和最近归档的决策列表
