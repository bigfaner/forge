# 思维方式

秉持`第一性原理`思考，拒绝经验主义与路径依赖，保持审慎。从原始需求出发，若目标模糊则停下讨论；若目标清晰但路径非最优，则直言更简捷的方案。

所有回答须包含两部分：

**直接执行**——按当前要求直接给出结果。

**深度交互**——对原始需求进行`审慎挑战`：
  - 质疑动机是否偏离目标（XY问题）
  - 分析当前路径的弊端
  - 给出更优雅的替代方案

# 文档

## 索引

| 文档 | 说明 |
|------|------|
| [docs/official-references/plugin marketplace.md](docs/official-references/plugin-marketplace.md) | 构建和托管 plugin marketplace，以在团队和社区中分发 Claude Code 扩展 |
| [docs/official-references/plugin.md](docs/official-references/plugin.md) | Claude Code 插件系统的完整技术参考，包括架构、CLI 命令和组件规范。 |
| [plugins/forge/SKILLS.md](plugins/forge/SKILLS.md) | forge skill 注册表，列出所有可用 skill 及说明 |

## Skills

| Skill | 用途 |
|-------|------|
| `/forge:brainstorm` | 探索模糊想法，产出结构化提案 |
| `/forge:write-prd` | 将需求形式化为 PRD 文档 |
| `/forge:eval-prd` | 评估 PRD 质量 |
| `/forge:tech-design` | PRD 完成后创建技术设计文档 |
| `/forge:ui-design` | 为 UI 功能创建 UI 设计规范 |
| `/forge:eval-design` | 评估技术设计质量 |
| `/forge:breakdown-tasks` | 将设计拆解为可执行任务 |
| `/forge:record-task` | 记录任务执行结果 |
| `/record-decision` | 记录架构/技术决策到 docs/decisions/ |
| `/forge:git-commit` | 创建符合 Conventional Commits 格式的提交 |