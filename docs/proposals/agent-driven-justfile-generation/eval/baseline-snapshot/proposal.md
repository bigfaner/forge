---
created: 2026-06-08
author: "faner"
status: Draft
intent: refactor
---

# Proposal: Agent-Driven Justfile Generation

## Problem

`init-justfile` skill 的 recipe 生成依赖 6 个语言模板文件（go/node/python/rust/mixed/generic.just），模板硬编码了每种语言的命令（`go vet`、`npm run build`、`pytest` 等）。当项目使用非标准结构（自定义构建工具、monorepo 子目录、非主流语言）时，模板默认值反而是障碍，agent 需要大量覆盖模板内容。

### Evidence

- 6 个模板共 ~800 行，mixed.just 单独 230 行
- 新增语言（Java、Kotlin、Ruby）需从零编写模板
- agent 的三层生成流程（加载模板 → Convention 覆盖 → LLM 微调）中，第一步（模板）和第三步（LLM 微调）经常冲突——模板给出的默认命令与实际项目不匹配，LLM 需要全部替换

### Urgency

模板是灵活性的主要瓶颈。非标准项目生成后必须手动编辑 justfile，削弱了自动化脚手架的价值。移除模板后，agent 直接针对实际项目结构生成正确命令，消除这层摩擦。

## Proposed Solution

移除所有语言模板，以 surface 为唯一驱动源：
1. **枚举 recipe**：agent 根据 surfaceKey + surfaceType，从 surface rule 文件读取 recipe 清单和契约
2. **填充内容**：agent 根据检测到的语言/框架 + Convention 知识 + 自身知识生成每个 recipe 的具体命令
3. **复杂 bash 模式**：server lifecycle（PID 追踪、幂等启动、健康检查）提取到独立的 `rules/server-lifecycle.md`

### Innovation Highlights

传统脚手架工具（cookiecutter、yeoman、hygen）使用模板驱动生成。本方案用 **LLM 驱动 + 结构约束** 替代模板，实现零模板维护的脚手架生成。Surface rule 文件定义"要生成什么"（契约），agent 决定"怎么生成"（内容）。这与 Forge 的 surface-first 设计理念一致——surface type 决定编排序列，语言/框架只影响具体命令。

## Requirements Analysis

### Key Scenarios

1. **单 surface 标量项目**（`surfaces: cli`）：生成无前缀 recipe（compile, unit-test, test, teardown）
2. **多 surface 命名项目**（`frontend=web + backend=api`）：生成 surface 前缀 recipe（frontend-compile, backend-test 等）
3. **混合语言项目**（Go backend + Node frontend）：每个 surface 独立检测语言，生成对应命令
4. **非标准项目结构**（自定义构建工具、非主流框架）：agent 直接适配，无模板默认值干扰
5. **空 surfaces**：提示用户运行 `forge init` 配置 surfaces

### Non-Functional Requirements

- **一致性**：相同项目多次运行生成的 justfile 结构一致（recipe 名称、分组、边界标记不变；具体命令可能因 LLM 变化而有细微差异）
- **向后兼容**：生成的 justfile 结构（boundary markers、recipe 命名、分组、user-customized 标记）与当前模板输出一致

### Constraints & Dependencies

- 依赖 `forge surfaces` 命令输出（已稳定）
- 依赖 surface rule 文件定义 recipe 契约（已有 5 个）
- 依赖 Convention 文件提供框架特定知识（已有机制）
- 不依赖任何语言模板文件

## Alternatives & Industry Benchmarking

### Industry Solutions

Make/Just 脚手架工具普遍使用模板（template → fill → output）。更现代的工具（如 Earthly、Taskfile）通过 DSL 抽象部分命令，但仍需用户编写。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | 模板维护负担持续增长，非标准项目体验差 | Rejected: 灵活性瓶颈未解决 |
| Parameterized templates | hygen/cookiecutter 模式 | 可配置性提升 | 仍需维护模板，复杂度转移为参数爆炸 | Rejected: 治标不治本 |
| Agent-driven + structural rules | 本方案 | 零模板维护，无限语言覆盖，天然适配非标准结构 | LLM 生成有轻微不确定性 | **Selected: 彻底解决灵活性问题** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Agent（Claude/GPT）已具备主流语言构建命令的知识。Surface rule 文件已定义完整的 recipe 契约。Convention 机制已提供框架特定知识的注入通道。唯一的技术挑战是 server lifecycle bash 代码的可靠性，通过提取为独立 rule 文件解决。

### Resource & Timeline

改动集中在 `init-justfile` skill 目录内：
- 删除 6 个模板 + 1 个 rule
- 新增 1 个 rule（server-lifecycle.md）
- 重写 SKILL.md 流程
- 简化 5 个 surface rule 文件

预计单个 skill 改动，可在一次 session 内完成。

### Dependency Readiness

所有依赖已就绪：`forge surfaces` 命令、surface rules、Convention 机制、just >= 1.50.0 要求。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 语言模板是生成 recipe 的必要起点 | XY Detection: 模板是手段（X），正确命令是目标（Y）。Agent 知识已足够覆盖 Y | Overturned: agent 可直接从 surface rule + 项目检测生成命令，无需模板中间层 |
| server lifecycle bash 必须由模板提供 | Occam's Razor: lifecycle 代码是语言无关的结构性模式，不应分散在 6 个模板中 | Confirmed: 提取为独立 rule 是更简洁的方案 |
| 需要项目类型分类（frontend/backend/mixed）| Assumption Flip: 如果不分类会怎样？agent 根据 surface 配置独立处理每个 surface，无需知道整体"类型" | Overturned: project type 分类可移除 |

## Scope

### In Scope

- 删除 6 个语言模板（go/node/python/rust/mixed/generic.just）
- 删除 `rules/project-detection.md`
- 新增 `rules/server-lifecycle.md`（PID 追踪、幂等启动、健康检查通用 bash 模式）
- 重写 SKILL.md：移除 `--type` 参数、移除项目类型检测步骤、surfaces 为前提（空时提示 `forge init`）、Step 3 改为 agent 驱动生成
- 简化 5 个 surface rule 文件：保留编排序列/recipe 契约/journey 策略，替换 TODO stub 模板为 "Recipe Generation Requirements" section

### Out of Scope

- forge CLI surfaces detect 命令改动
- 新增 surface type
- Convention 加载机制改动
- 其他 skill 更新（它们消费生成的 just recipe，不直接调用 init-justfile；recipe 契约不变）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| LLM 生成命令不一致（不同 run 产生不同命令） | M | L | Surface rule 的 recipe 契约 + Convention 约束了结构和语义；verification step (dry-run + actual) 捕获错误 |
| Server lifecycle bash 复杂度高，agent 可能遗漏边界情况 | M | M | 提取为独立 rule 文件，提供完整参考模式；现有验证步骤覆盖 |
| 冷启动（无 Convention）时生成质量下降 | L | M | agent 具备主流语言的默认知识；verification step 捕获错误并自修正 |
| 罕见语言/框架生成失败 | L | L | agent 回退到 error stub（与当前 generic.just 行为一致） |

## Success Criteria

- [ ] 所有 5 个 surface rule 文件的 TODO stub 模板已替换为 "Recipe Generation Requirements" section
- [ ] 6 个语言模板文件已删除，`rules/project-detection.md` 已删除
- [ ] `rules/server-lifecycle.md` 已创建，包含 PID 追踪、幂等启动、健康检查的完整 bash 模式
- [ ] `forge surfaces` 为空时输出提示用户运行 `forge init` 配置 surfaces
- [ ] 对 Go/Node/Python/Rust 项目，agent 生成的 recipe 通过 verification step（dry-run + actual execution）
- [ ] 混合语言多 surface 项目（如 Go backend=api + Node frontend=web），每个 surface 独立生成正确命令
- [ ] 生成的 justfile 结构（boundary markers、recipe 命名、[linux]/[windows] 双平台、user-customized 标记）与当前输出一致

consistency_check_result:
  status: pass
  pairs_checked: 21
  conflicts_found: 0

## Next Steps

- Proceed to `/tech-design` to define the detailed implementation plan
