# Smith: Claude Code Plugin Development Meta-Plugin

## Motivation

forge 插件提供了完整的软件开发工作流（brainstorm → PRD → design → tasks → testing）。但在开发 Claude Code 插件本身时，缺少专门的工具链来脚手架、验证、评估和改进插件。

smith 是一个"元插件"——用于开发 Claude Code 插件的插件。与 forge 形成兄弟关系：forge 面向通用软件开发，smith 面向插件开发。

命名：**smith**（铁匠），与 forge（锻炉）呼应。

## Skill Pipeline

```
/create-plugin → /gen-skill|/gen-command|/gen-agent → /validate-plugin → /eval-skill|/eval-plugin → /improve-skill
       ↓                    ↓                           ↓                    ↓                        ↓
  plugin.json + dirs   SKILL.md + templates      validation-report    rubric score report      improved skill
```

## Skills (8)

### 生成类 (4)

| Skill | 说明 |
|-------|------|
| create-plugin | 交互式创建完整插件目录结构，生成 plugin.json + SKILLS.md + hooks/guide.md |
| gen-skill | 在已有插件中生成新 skill：SKILL.md + templates/ + examples/ |
| gen-command | 在已有插件中生成新 slash command |
| gen-agent | 在已有插件中生成新 agent 定义 |

### 验证类 (1)

| Skill | 说明 |
|-------|------|
| validate-plugin | 验证插件结构：plugin.json schema、路径存在性、frontmatter 语法、hook 事件名 |

### 评估类 (2)

| Skill | 说明 |
|-------|------|
| eval-skill | 评估单个 skill 质量（100 分制），含 adversarial score-revise 循环 |
| eval-plugin | 评估完整插件质量（100 分制），单次评估 + 维度报告 |

### 改进类 (1)

| Skill | 说明 |
|-------|------|
| improve-skill | 根据 eval-skill 报告逐条改进 skill |

## Agents (2)

| Agent | 说明 | 复用自 |
|-------|------|--------|
| plugin-scorer | 插件/skill 评估打分（对抗性），返回 SCORE/DIMENSIONS/ATTACKS | forge/doc-scorer |
| skill-reviser | 根据评估报告改进 skill（建设性），最多 3 轮自审 | forge/doc-reviser |

## Commands (2)

| Command | 说明 |
|---------|------|
| /init-smith | 初始化 smith 环境，验证依赖 |
| /scaffold-plugin | 快速一脚手架（create-plugin 的精简版，最少交互） |

## Hooks

| Event | Matcher | Action |
|-------|---------|--------|
| SessionStart | `startup\|clear\|compact` | 注入 guide.md 上下文 |
| SubagentStart | (always) | 注入 guide.md 上下文 |
| PostToolUse | `Edit\|Write` + `*plugin.json` | 自动校验 plugin.json 语法 |

## References (4)

| 文件 | 内容 |
|------|------|
| plugin-structure-guide.md | 插件标准目录结构（提取自官方文档） |
| skill-writing-guide.md | SKILL.md 编写规范（从 forge skills 提炼） |
| agent-frontmatter-reference.md | Agent frontmatter 字段速查 |
| hook-events-reference.md | 28 个 hook 事件名速查 |

## 目录结构

```
plugins/smith/
├── .claude-plugin/
│   └── plugin.json
├── SKILLS.md
├── agents/
│   ├── plugin-scorer.md
│   └── skill-reviser.md
├── commands/
│   ├── init-smith.md
│   └── scaffold-plugin.md
├── hooks/
│   ├── hooks.json
│   ├── guide.md
│   ├── run-hook.cmd              # 复用 forge
│   ├── session-start              # 复用 forge
│   └── debug                      # 复用 forge
├── scripts/
│   └── validate-plugin-json.sh
├── references/
│   └── shared/
│       ├── plugin-structure-guide.md
│       ├── skill-writing-guide.md
│       ├── agent-frontmatter-reference.md
│       └── hook-events-reference.md
└── skills/
    ├── create-plugin/
    │   ├── SKILL.md
    │   └── templates/
    │       └── plugin-scaffold.md
    ├── gen-skill/
    │   ├── SKILL.md
    │   ├── templates/
    │   │   └── skill-template.md
    │   └── examples/
    │       └── example-skill-walkthrough.md
    ├── gen-command/
    │   ├── SKILL.md
    │   └── templates/
    │       └── command-template.md
    ├── gen-agent/
    │   ├── SKILL.md
    │   └── templates/
    │       └── agent-template.md
    ├── validate-plugin/
    │   ├── SKILL.md
    │   └── templates/
    │       └── validation-report.md
    ├── eval-skill/
    │   ├── SKILL.md
    │   └── templates/
    │       ├── rubric.md
    │       └── report.md
    ├── eval-plugin/
    │   ├── SKILL.md
    │   └── templates/
    │       ├── rubric.md
    │       └── report.md
    └── improve-skill/
        ├── SKILL.md
        └── templates/
            └── improvements.md
```

## 实施顺序

### Phase 1: 骨架 & 基础设施 (7 files)
1. `plugin.json` — manifest
2. `SKILLS.md` — skill 注册表
3. 复制 `run-hook.cmd`, `session-start`, `debug` from forge
4. `hooks/guide.md` — smith pipeline 说明
5. `hooks/hooks.json`
6. `scripts/validate-plugin-json.sh`

### Phase 2: 参考资料 (4 files)
7. `plugin-structure-guide.md`
8. `skill-writing-guide.md`
9. `agent-frontmatter-reference.md`
10. `hook-events-reference.md`

### Phase 3: 生成 Skills (4 skills, ~10 files)
11. `gen-skill/` (SKILL.md + template + example)
12. `gen-command/` (SKILL.md + template)
13. `gen-agent/` (SKILL.md + template)
14. `create-plugin/` (SKILL.md + template)

### Phase 4: 验证 & 评估 (4 skills + 2 agents, ~14 files)
15. `validate-plugin/`
16. `agents/plugin-scorer.md`
17. `agents/skill-reviser.md`
18. `eval-skill/`
19. `eval-plugin/`
20. `improve-skill/`

### Phase 5: Commands (2 files)
21. `init-smith.md`
22. `scaffold-plugin.md`

**总计: ~37 个文件**

## 关键设计决策

| 决策 | 理由 |
|------|------|
| eval-skill 使用 adversarial loop | 单个 skill 可原地迭代改进，适合 score-gate-revise 模式 |
| eval-plugin 使用单次评估 | 整个插件的改进跨越多文件，报告引导用户用 improve-skill 或 gen-* 处理 |
| scorer 和 reviser 分离 | 保持 scorer 的对抗性，不让它"心软" |
| hooks 基础设施复用 forge | run-hook.cmd 的跨平台 polyglot 模式已验证，无需重写 |

## 验证方式

1. `claude --plugin-dir plugins/smith` 加载插件
2. `/validate-plugin` 对 forge 本身运行
3. `/create-plugin` 创建测试插件
4. `/eval-skill` 对 forge 的某个 skill 运行评估循环
5. `/scaffold-plugin test-plugin` 快速脚手架
