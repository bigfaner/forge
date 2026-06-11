---
created: 2026-05-12
author: faner
status: Draft
---

# Proposal: 测试 Profile 多策略系统

## Problem

Forge 的 E2E 测试管线硬编码为 TypeScript + Playwright，无法适配非 Web 项目。三种典型失配场景：

| 项目类型 | 现状 | 问题 |
|---|---|---|
| Go TUI (agent-forensic) | 生成 TS/Playwright 壳，内部用 `execSync` 调 `go test` | Playwright 仅当测试运行器，模板中浏览器相关代码全部废弃 |
| Mobile App (train-recorder) | 放弃 Forge 生成，手写 Maestro YAML | Forge 管线完全不可用 |
| Java/Rust/Python CLI | 不支持 | 无路径 |

**根因**：Forge 假设所有项目都是 Web 应用，测试框架、模板、目录结构、执行命令全部耦合为 Playwright/TypeScript。

## Solution

引入 **Test Profile** 概念——定义一套可插拔的测试策略系统。每个 profile 封装一种技术栈的测试方案，Forge 框架作为调度器按 profile 分发。

### Core Decisions

#### D1. 架构模式：全模板化

当前 TS/Playwright 降级为 `web-playwright` profile（与其他 profile 平等）。不引入 flag 开关或外部插件机制。

#### D2. 多 Profile 支持

一个项目可声明多个 profile（如 monorepo 前后端分离）。大多数项目只需一个 profile。

#### D3. Profile 列表（v3.0.0 全部交付）

| Profile | 适用场景 | 测试工具 | 语言 | Capabilities |
|---|---|---|---|---|
| `web-playwright` | Web 前端 + API + CLI | @playwright/test | TypeScript | web-ui, api, cli |
| `go-test` | Go CLI/TUI/后端 | go test | Go | tui, api, cli |
| `maestro` | React Native / Flutter | Maestro CLI | YAML | mobile-ui, api |
| `java-junit` | Java CLI/后端 | JUnit 5 + Maven | Java | tui, api, cli |
| `rust-test` | Rust CLI/后端 | cargo test | Rust | tui, api, cli |
| `pytest` | Python CLI/后端 | pytest | Python | tui, api, cli |

#### D4. 配置位置：`.forge/config.yaml`

```yaml
test-profiles:
  - web-playwright
  - go-test
```

- 声明性配置，提交到 git
- 运行时配置（baseUrl、credentials）仍留在 `tests/e2e/config.yaml`

#### D5. 配置时机

- **已设置**：tech-design 沿用，/quick 沿用
- **未设置**：
  1. tech-design 阶段：自动探测项目结构 → 探测不到 → 主动询问用户
  2. /quick 启动时：同样的探测 → 询问流程
- **无静默默认**：不 fallback 到 web-playwright

自动探测规则：

| 信号 | 推断 Profile |
|---|---|
| `package.json` + `@playwright/test` | web-playwright |
| `go.mod` | go-test |
| `android/` 或 `ios/` 目录 | maestro |
| `pom.xml` 或 `build.gradle` | java-junit |
| `Cargo.toml` | rust-test |
| `requirements.txt` + `pytest` | pytest |

#### D6. Capabilities 封闭枚举

| Capability | 含义 |
|---|---|
| `web-ui` | 浏览器 UI（DOM 交互） |
| `tui` | 终端 UI（文本渲染、键盘交互） |
| `mobile-ui` | 移动端 UI（触摸、手势） |
| `api` | HTTP/网络接口 |
| `cli` | 命令行界面 |

- 当前封闭，后续再考虑开放
- 新增 capability 需要改 forge 核心代码

#### D7. Profile-Agnostic Test Cases

`gen-test-cases` 生成的 `test-cases.md` 只描述**行为**（测什么），不涉及**实现**（怎么测）。

- gen-test-cases 读 `.forge/config.yaml` 获取 capabilities，按集合运算标注 test case 为 `automated` 或 `manual-only`
- 不匹配规则：PRD 有 web-ui 功能但 profile 只有 `tui` → 标记 `manual-only`
- 无 UI 功能 → 跳过
- gen-test-scripts 根据标注和 profile 策略生成具体实现

#### D8. Auth 分类：框架分类 + Profile 实现

4 种 auth 类型（`login-test`、`auth-required-test`、`public-test`、`custom-auth-test`）由 gen-test-cases 从 PRD 推导，与 profile 无关。具体实现由各 profile 的 `generate.md` 决定。

#### D9. 任务管线：按 Profile 展开（方案 B）

breakdown-tasks 读 `.forge/config.yaml` 的 profiles，为每个 profile 生成独立的 T-test-2~4 任务：

```
T-test-1: gen-test-cases              (单任务，profile-agnostic)
T-test-1b: eval-test-cases            (单任务，profile-aware 评分)
T-test-2a: gen-test-scripts (profile-a)
T-test-2b: gen-test-scripts (profile-b)
T-test-3a: run-e2e-tests (profile-a)
T-test-3b: run-e2e-tests (profile-b)
T-test-4a: graduate-tests (profile-a)
T-test-4b: graduate-tests (profile-b)
T-test-4.5: verify-regression         (单任务，全量)
T-test-5: consolidate-specs           (单任务，只从文档提取)
```

quick-tasks 同理，T-quick 任务按 profile 展开。

#### D10. Skill 架构：调度器 + Profile 策略

三个核心测试 skill 统一采用调度器模式：

| Skill | 职责 | Profile 策略文件 |
|---|---|---|
| gen-test-scripts | 读 test-cases + config → 遍历 profiles → 委托 | `profile/generate.md` |
| run-e2e-tests | 调 justfile → 读结果 → 遍历 profiles → 委托 | `profile/run.md` |
| graduate-tests | 读 staging → 遍历 profiles → 委托 | `profile/graduate.md` |

skill.md 做通用部分（读配置、遍历 profile、调用 justfile），profile 策略文件做专有部分。

#### D11. Profile 目录结构

```
plugins/forge/profiles/<name>/
  manifest.yaml          # 元数据 + capabilities + 命令声明
  generate.md            # gen-test-scripts 的 profile 策略
  run.md                 # run-e2e-tests 的 profile 策略
  graduate.md            # graduate-tests 的 profile 策略
  templates/             # 代码模板（spec 文件、helpers、config 等）
```

每个 profile 自包含，新增 profile 只需加一个目录。

#### D12. Manifest Schema

```yaml
name: go-test
display: "Go Test"
language: go
file-extension: .go
test-directory: tests/e2e/

capabilities: [tui, api, cli]

templates:
  test-file: templates/test-file.go
  helpers: templates/helpers.go
  config-file: null
  additional: []

run:
  command: "go test ./tests/e2e/... -v -tags=e2e"
  compile: "go build ./..."
  result-format: go-json

graduate:
  target-directory: tests/e2e/
  merge-strategy: package
  import-rewrite: null
```

#### D13. 统一目录约定

所有 profile 遵循 forge 的目录结构：

```
tests/e2e/
  features/<slug>/        # staging（所有 profile）
  <module>/               # graduated regression（所有 profile）
```

Profile 只声明文件扩展名，不自定义目录。justfile recipe 知道如何从该目录运行对应框架。

#### D14. 执行路径：经过 Justfile

run-e2e-tests 通过 `just test-e2e --feature <slug>` 执行，不直接调用测试命令。

- `init-justfile` skill 读 `.forge/config.yaml`，为每个 active profile 生成对应的 justfile recipe
- manifest 中声明的 `run.command` 供 init-justfile 参考生成 recipe

#### D15. eval-test-cases：按 Capabilities 动态评分

评分维度根据项目 capabilities 动态选择：

| 维度 | 分值 | 适用 Capability |
|---|---|---|
| PRD Traceability | 30 | 通用 |
| Step Actionability | 30 | 通用 |
| Interface Accuracy | 20 | 按 capability 替换内容（见下表） |
| Completeness | 20 | 通用 |

Interface Accuracy 按 capability 替换评分内容：

| Capability | 评分重点 |
|---|---|
| web-ui | Route & Element Accuracy（sitemap + locator） |
| tui | Output Assertion Accuracy（golden file / snapshot 对比点） |
| mobile-ui | Interaction Accuracy（触摸/手势/导航流程） |
| api | Contract Accuracy（请求/响应结构匹配度） |
| cli | Command Coverage（flag / subcommand / error 覆盖度） |

#### D16. Validate：当前砍掉

- 移除 `validate-specs.mjs` 及相关校验逻辑
- 记入方案备忘：未来可用 tree-sitter 构建通用 AST 校验引擎（声明式规则 + tree query），当前不引入依赖

#### D17. 非 Test Skill 的改动

| Skill | 改动 |
|---|---|
| tech-design | 新增 profile 选择步骤（探测 → 确认 → 写入 .forge/config.yaml） |
| init-justfile | 读 .forge/config.yaml，为每个 profile 生成 test-e2e recipe |
| breakdown-tasks | 读 profiles 动态展开 T-test 任务 |
| quick-tasks | 读 profiles 动态展开 T-quick 任务 |
| /quick | 启动时检查 .forge/config.yaml 是否存在 |
| gen-sitemap | 非 web-ui 项目跳过 |
| consolidate-specs | 只从文档提取（PRD、design、test-cases），不解析测试代码 |
| execute-task | 无改动（Quality Gate 走 justfile，天然 profile-agnostic） |
| run-tasks | 无改动（只调度任务） |

### `.forge/` 目录定位

```
.forge/
  config.yaml           # 声明性配置（test-profiles、项目类型等）
```

- 声明性、不随环境变化 → `.forge/`
- 运行时、随环境变化 → `tests/e2e/config.yaml`

### Impact Assessment

#### 改动的 Skills（9 个）

1. gen-test-scripts — 重构为调度器
2. run-e2e-tests — 重构为调度器
3. graduate-tests — 重构为调度器
4. gen-test-cases — 增加 capabilities 读取和标注逻辑
5. eval-test-cases — 动态评分维度
6. tech-design — 增加 profile 选择步骤
7. breakdown-tasks — 动态展开测试任务
8. quick-tasks — 动态展开测试任务
9. init-justfile — profile-aware recipe 生成

#### 新增的 Profiles（6 个）

1. web-playwright（从现有代码迁移）
2. go-test（新增）
3. maestro（新增）
4. java-junit（新增）
5. rust-test（新增）
6. pytest（新增）

#### 新增的配置

- `.forge/config.yaml` — 项目级声明性配置

#### 不变的部分

- Quality Gate（`just compile → just fmt → just lint → just test`）— 走 justfile，profile-agnostic
- all-completed hook — 走 justfile
- execute-task / run-tasks — 不涉及
- task-cli — 不涉及
- brainstorm / write-prd — 不涉及

### Risks

| 风险 | 缓解 |
|---|---|
| 6 个 profile 全交付周期长 | web-playwright 已有代码迁移，其余按相似度优先（go-test/rust-test 相似，java-junit/pytest 相似） |
| Profile 策略文件质量不均 | 每个 profile 用真实项目验证（go-test → agent-forensic，maestro → train-recorder） |
| 现有项目升级断裂 | tech-design 自动探测 + 主动询问，不静默默认 |
| gen-test-cases 的 capabilities 映射不准 | 先从 PRD 关键词粗匹配，后续迭代优化 |
