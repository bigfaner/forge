---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: Forge Info Commands — init/proposal/feature/lesson/config 命令

## Problem

Forge CLI 目前的命令以「写操作」为主（task claim/submit、profile set、feature <slug>）。Agent 和用户在运行时需要查询项目状态，但缺乏对应的「读操作」命令：

1. **无法快速了解项目全貌** — 有哪些 proposal、哪些 feature、状态如何？只能手动 `ls` + 读文件
2. **配置分散在多个来源** — `project-type` 硬编码在 justfile，`capabilities` 埋在 profile manifest 中，agent 需要多步查询才能拼出完整配置
3. **lesson 知识库不可发现** — `docs/lessons/` 有 41 个文件，agent 无法枚举或按标签检索
4. **skill 依赖 justfile 获取配置** — scope resolution 调用 `just project-type`，绕过了统一的 config.yaml 体系
5. **项目初始化无统一入口** — 新项目接入 forge 需要：手动创建 `.forge/` 目录、手动配置 `.gitignore`、手动编写 `CLAUDE.md`、手动运行 `forge config init`。没有一步到位的命令

### Evidence

- `just project-type` 在 justfile 中硬编码返回 `"backend"`，配置无法被 CLI 管理
- 10 个 skill 遵循相同的 "Step 0: Resolve Profile" 模式，其中 capabilities 需要两步获取（`task profile` → `task profile get --manifest`）
- guide.md、fix-bug.md、execute-task.md、run-tasks.md 等 6 处 scope resolution 逻辑都依赖 `just project-type`
- 无任何命令可列出 proposals 或 features 的状态概览

### Urgency

随着 feature 和 proposal 数量增长，项目状态查询的需求越来越频繁。尽早建立查询命令体系，避免 agent 和用户依赖 `ls` + 文件读取的脆弱模式。

## Proposed Solution

新增 5 组命令，覆盖项目初始化、状态查询和配置管理：

### 1. `forge init` — 项目初始化

一站式初始化 forge 项目环境。依次执行：

| 步骤 | 操作 | 冲突处理 |
|------|------|---------|
| 创建 `.forge/` 目录 | `mkdir -p .forge` | 已存在则跳过 |
| 写入 `CLAUDE.md` | 从 CLI 内嵌模板复制到项目根目录 | 已存在则跳过 |
| 更新 `.gitignore` | 追加 forge 运行时忽略规则（去重检查） | 追加不覆盖 |
| 更新 `justfile` | 追加 `claude` / `claude-c` recipe（去重检查） | 追加不覆盖 |
| 交互式配置 | 等同 `forge config init`（project-type + test-profiles + capabilities） | `.forge/config.yaml` 已存在则跳过 |

**追加到 `.gitignore` 的条目**：

```
# Forge runtime
docs/features/*/tasks/process/
.forge/state.json
tests/results/.last-run.json
tests/e2e/results/.last-run.json
tests/e2e/results/*/error-context.md
```

**追加到 `justfile` 的 recipe**：

```just
claude:
    claude --dangerously-skip-permissions

claude-c:
    claude --dangerously-skip-permissions -c
```

**CLAUDE.md 模板**：内嵌在 CLI 二进制中（使用 `go:embed`），内容为通用行为准则（Think Before Coding / Simplicity First / Surgical Changes / Goal-Driven Execution）。用户可在此基础上自定义。

**执行结果报告**：

```
>>>
CREATED   .forge/
CREATED   CLAUDE.md (from template)
APPENDED  .gitignore (5 entries)
APPENDED  justfile (2 recipes: claude, claude-c)
CREATED   .forge/config.yaml (interactive)
<<<
```

### 2. `forge proposal` — 提案查询

| 命令 | 用途 | 用户 |
|------|------|------|
| `forge proposal` | 表格列出所有 proposals | 人类 / Agent |
| `forge proposal <slug>` | 展示单个 proposal 详情 | 人类 / Agent |

**列表表格列**：Slug | Created | Status | PRD关联 | Feature状态

- **Created**：优先读取 `proposal.md` frontmatter 的 `created` 字段，缺失时回退到文件系统 birth time
- **PRD关联**：检查 `docs/features/{同slug}/prd/prd-spec.md` 是否存在
- **Feature状态**：检查 `docs/features/{同slug}/manifest.md` 的 `status` 字段

**详情视图**：元数据（created, author, status）+ 内容摘要（Problem + Proposed Solution 摘要）+ 关联状态（PRD/Feature/Task 进度）+ 文件路径

### 3. `forge feature` — Feature 查询（扩展现有命令）

保留现有行为不变，新增子命令：

| 命令 | 用途 | 用户 |
|------|------|------|
| `forge feature` | 显示当前 feature（现有） | 人类 |
| `forge feature <slug>` | 设置当前 feature（现有） | 人类 |
| `forge feature list` | 表格列出所有 features | 人类 / Agent |
| `forge feature status <slug>` | 展示 feature 详细状态 | 人类 / Agent |

**列表表格列**：Slug | 状态 | 任务进度 | PRD(分值) | Design(分值) | UI(分值) | Tests(分值)

- **任务进度**：从 `tasks/index.json` 统计 completed/total
- **分值**：从各产物 frontmatter 的 `score` 字段读取（需 eval skill 评估后回写）
- 分值缺失时显示 `—`

**status 详情**：manifest 摘要 + 任务统计（pending/in_progress/completed/blocked/skipped/rejected 计数）+ 产物清单+分值 + 当前进行中任务信息

### 4. `forge lesson` — Lesson 查询

| 命令 | 用途 | 用户 |
|------|------|------|
| `forge lesson` | 表格列出所有 lessons | 人类 / Agent |
| `forge lesson <name>` | 展示 lesson 元数据和路径 | 人类 / Agent |

**列表表格列**：名称 | 创建时间 | 标签 | 分类

- **分类**：从文件名前缀推断（gotcha-/arch-/pattern-/tool-/lesson-/hook-）
- **创建时间**：从 frontmatter `created` 字段读取
- **标签**：从 frontmatter `tags` 字段读取

**详情视图**：元数据（created, tags）+ 文件路径（不打印完整内容）

### 5. `forge config` — 配置管理

| 命令 | 用途 | 用户 |
|------|------|------|
| `forge config init` | 交互式初始化 `.forge/config.yaml` | 人类 |
| `forge config get <key>` | 获取配置值，纯文本输出 | Agent |

**`forge config init` 交互流程**：

1. **project-type** — 选择：frontend / backend / mixed
2. **test-profiles** — 多选，内置列表（go-test / web-playwright / pytest / java-junit / rust-test / maestro）
3. **capabilities** — 从所选 profiles 的 capabilities 取并集，让用户勾选实际需要的子集

最终 `.forge/config.yaml` 示例：

```yaml
project-type: backend
test-profiles:
  - go-test
capabilities:
  - tui
  - api
  - cli
```

**`forge config get <key>`**：

- 支持任意 key（project-type / capabilities / test-profiles / 其他自定义字段）
- 标量直接输出纯文本，数组每行一个值
- key 不存在时退出码 1，无输出
- 已存在 config 时 `init` 提示是否重新配置

### 6. 配置迁移 — `just project-type` → `forge config get project-type`

**迁移范围**：

| 文件 | 当前方式 | 迁移后 |
|------|---------|--------|
| `plugins/forge/hooks/guide.md` | `just project-type` | `forge config get project-type` |
| `plugins/forge/commands/fix-bug.md` | `just project-type` | `forge config get project-type` |
| `plugins/forge/commands/execute-task.md` | `just project-type` | `forge config get project-type` |
| `plugins/forge/commands/run-tasks.md` | `just project-type` | `forge config get project-type` |
| `plugins/forge/skills/*/SKILL.md`（6 处） | `just project-type` | `forge config get project-type` |
| `justfile` | `project-type:` recipe | 删除该 recipe |
| `plugins/forge/skills/init-justfile/templates/*.just` | 各模板中的 `project-type:` recipe | 删除 |

**Go 代码改造**：

- `just.ResolveScope()` 从调用 `just project-type` 子进程改为直接读取 `.forge/config.yaml` 的 `project-type` 字段
- `ForgeConfig` struct 新增 `ProjectType` 和 `Capabilities` 字段

## Scope

### In Scope

- `forge init` 一站式初始化（.forge/ 目录 + CLAUDE.md + .gitignore + justfile recipe + config.yaml）
- CLAUDE.md 模板内嵌到 CLI 二进制
- `forge proposal` 列表 + 详情命令
- `forge feature list` + `forge feature status <slug>` 子命令
- `forge lesson` 列表 + 详情命令
- `forge config init` 交互式初始化
- `forge config get <key>` 配置查询
- `.forge/config.yaml` schema 扩展（project-type、capabilities）
- `just.ResolveScope()` 改造为直接读 config.yaml
- guide.md 和所有 skill 中的 scope resolution 迁移
- justfile 和模板中删除 `project-type` recipe

### Out of Scope

- eval skill 回写产物 frontmatter `score` 字段（独立工作）
- `forge config set <key> <value>` 命令（后续迭代）
- lesson 搜索功能（后续迭代）
- `forge proposal` 的创建/编辑功能（由 `/brainstorm` skill 负责）
- CLAUDE.md 内容定制（用户手动编辑模板生成的基础版本）

## Key Risks

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| 交互式依赖增加二进制体积 | 低 | 低 | 使用 lightweight 库（如 bubbletea）或纯 stdin 读取 |
| 旧 justfile 模板兼容性 | 中 | 中 | `forge config init` 作为初始化步骤，旧项目需手动运行一次 |
| feature list 分值列全为空（eval 未回写） | 高 | 低 | 分值缺失时显示 `—`，不影响其他信息 |
| forge init 覆盖用户自定义 CLAUDE.md | 低 | 高 | 已存在时跳过，不覆盖 |
| justfile 格式多样性导致追加失败 | 中 | 中 | 仅追加简单 recipe，不做格式化 |

## Success Criteria

- [ ] `forge init` 创建 `.forge/` 目录
- [ ] `forge init` 从内嵌模板生成 `CLAUDE.md`
- [ ] `forge init` 追加 forge 运行时条目到 `.gitignore`（去重）
- [ ] `forge init` 追加 `claude` / `claude-c` recipe 到 `justfile`（去重）
- [ ] `forge init` 在 `.forge/config.yaml` 不存在时运行交互式配置
- [ ] `forge init` 对已存在文件跳过并报告结果
- [ ] `forge proposal` 列出所有 proposals，含正确的创建时间、PRD关联、Feature状态
- [ ] `forge proposal <slug>` 展示完整详情
- [ ] `forge feature list` 列出所有 features，含状态和任务进度
- [ ] `forge feature status <slug>` 展示完整状态详情
- [ ] `forge lesson` 列出所有 lessons，含分类和标签
- [ ] `forge lesson <name>` 展示元数据和路径
- [ ] `forge config init` 交互式收集 project-type + test-profiles + capabilities
- [ ] `forge config get project-type` 返回正确的值
- [ ] `forge config get capabilities` 返回正确的值（每行一个）
- [ ] `ResolveScope()` 直接读 config.yaml，不再调用 just 子进程
- [ ] 所有 skill 中的 scope resolution 使用 `forge config get project-type`
- [ ] justfile 中无 `project-type` recipe
- [ ] 测试覆盖率 ≥ 80%

## Next Steps

1. 运行 `/eval-proposal` 评估本提案
2. 通过后进入 `/write-prd` → `/tech-design` → `/breakdown-tasks`
