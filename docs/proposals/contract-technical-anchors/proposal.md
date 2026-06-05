---
created: "2026-06-05"
author: fanhuifeng
status: Draft
intent: "enhancement"
---

# Proposal: Contract Technical Anchors

## Problem

Forge 测试管道生成的 Contract 规格缺少技术锚点（API endpoint、CLI command、Web page、Mobile screen），导致 gen-test-scripts 只能依赖 LLM 推断技术细节（HTTP method、命令名、页面路由），当推断错误时测试脚本与实际代码不匹配，三层测试（E2E + 单元 + 集成）均无法捕获。

### Evidence

pm-work-tracker 项目中 Move sub-item 操作的 bug（gotcha-mock-repo-skips-whitelist.md）：
- Contract 未指定 HTTP method → gen-test-scripts 猜测 POST，实际路由注册为 PUT
- api-handbook 正确定义了 `PUT /teams/:teamId/sub-items/:subId/move`，但 Contract 不引用它
- 三层测试（E2E + service unit + API handbook）全部漏掉，直到生产环境返回 422

### Urgency

此类问题会在每个有 API 或 CLI surface 的项目中重复出现。当前 gen-test-scripts 的 Fact Table 代码侦察虽然能扫描路由注册，但不会与 Contract 做交叉比对，侦察结果和 Contract 之间的矛盾不会被发现。

## Proposed Solution

为 Contract 规格建立"设计文档 → Contract → 测试代码"的技术锚点信息链，消除信息断层：

1. **tech-design 自动生成全 surface 手册**：API → api-handbook（已有）、CLI/TUI → cli-handbook、Web → page-map、Mobile → screen-map
2. **Contract frontmatter 增加锚点字段**：`endpoint`（API）、`command`（CLI/TUI）、`page`（Web）、`screen`（Mobile）
3. **gen-contracts 从设计文档填充锚点**：读取对应 handbook，提取技术细节写入 Contract frontmatter

4. **gen-test-scripts 交叉验证并建议修复**：将 Fact Table（代码侦察）与 Contract frontmatter 比对，不匹配时输出分类结果（高置信度/低置信度/无法验证），以设计文档为准生成建议修复，由用户确认后写入 Contract；若设计文档与代码实现不一致，标记为代码 bug

5. **handbook 新鲜度检查**：gen-contracts 和交叉验证时比对 handbook 生成时间戳与 tech-design 最后修改时间，若 handbook 过期则提示用户重新生成

### Innovation Highlights

将测试管道从"LLM 推断技术细节"升级为"设计文档驱动的锚点验证"。核心洞察：Contract 是设计意图的规格说明，技术锚点应该来自设计阶段确定的接口定义，而非从代码中逆向提取或让 LLM 猜测。交叉验证以设计文档为 authority source，设计-实现不一致时定位为代码 bug 而非测试问题。


### Known Limitations

- **双源可靠性**：交叉验证依赖 Fact Table（静态代码侦察）和 handbook（设计文档生成）两个源的比对，两者各自有不可靠的场景（动态路由注册、设计文档滞后）。当两者都不可靠时，交叉验证结果为"低置信度"或"无法验证"，降级为提示用户人工确认，而非自动处理
- **代码侦察覆盖度**：静态分析无法覆盖所有路由注册模式（插件系统、动态加载、反射机制等），侦察结果不完整时会标记为"低置信度"或"无法验证"

## Requirements Analysis

### Key Scenarios

- **新功能先设计后实现**：tech-design 生成 handbook → gen-contracts 填充锚点 → 代码实现 → gen-test-scripts 交叉验证（代码存在时）
- **已有功能补全测试**：gen-contracts 从现有 handbook 填充锚点 → gen-test-scripts 交叉验证发现不匹配 → 建议修复，用户确认后写入 Contract
- **设计文档与代码不一致**：交叉验证发现 handbook 说 PUT 但代码是 POST → 建议修复 Contract 以 handbook 为准（用户确认），标记代码 bug
- **Handbook 不存在**：gen-contracts 阶段跳过锚点填充并提示用户"缺少 handbook，建议运行 tech-design 生成"；gen-test-scripts 发现缺少锚点时降级为 Fact Table 推断（向后兼容）
- **Handbook 过期**：gen-contracts 检测到 handbook 生成时间早于 tech-design 最后修改时间时，提示用户"handbook 可能过期，建议重新生成"

### Non-Functional Requirements

- 向后兼容：缺少 handbook 或锚点字段时，管道不中断，降级为现有行为并提示用户
- 性能影响：交叉验证在 gen-test-scripts Step 1（代码侦察）中执行，无显著额外开销（仅增加内存比对）


### Anchor Field Schema

每种 surface 类型的锚点字段定义如下：

| Surface | 必填字段 | 可选字段 | 说明 |
|---------|---------|---------|------|
| API | `endpoint` (string)、`method` (string) | `content_type`、`auth_required` (boolean) | endpoint 格式：`/path/:param`，method 为 HTTP verb |
| CLI | `command` (string) | `subcommand` (string)、`flags` (string[])、`aliases` (string[]) | 支持子命令嵌套，如 `forge surfaces list` |
| TUI | `command` (string) | `interactive_prompt` (string)、`keybindings` (string[]) | 终端交互界面的入口命令 |
| Web | `page` (string) | `route` (string)、`requires_auth` (boolean)、`layout` (string) | page 为页面名称，route 为 URL 路径 |
| Mobile | `screen` (string) | `navigation_path` (string[])、`deeplink` (string)、`platform` (string) | screen 为屏幕名称，navigation_path 为导航层级 |

### Constraints & Dependencies

- 依赖 tech-design skill 支持全 surface 手册生成（需扩展）
- 依赖 api-handbook 的现有格式稳定性
- CLI/TUI/Web/Mobile handbook 格式需要新定义

## Alternatives & Industry Benchmarking

### Industry Solutions

行业中常见的做法是 Contract Testing（Pact）或 OpenAPI spec 驱动的测试生成。Forge 的 Contract 是语义级别的规格说明，不直接等同于 OpenAPI，但"从权威源提取技术细节"的思路一致。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | 每个 API/CLI 项目都会重复出现 LLM 猜测问题 | Rejected: lesson 已证明三层测试漏掉 |
| 仅增强 Fact Table | 代码侦察 | 简单，不改 Contract 结构 | 不解决设计-实现一致性问题，无法在设计阶段暴露缺口 | Rejected: 治标不治本 |
| OpenAPI spec 驱动 | Pact/Swagger | 行业标准 | 要求项目维护 OpenAPI spec，与 Forge 的语义 Contract 模型不兼容 | Rejected: 架构不匹配 |
| **Contract Technical Anchors** | 设计文档 | 完整信息链、建议修复+人工确认、全 surface 覆盖 | 需要扩展 tech-design 和 gen-test-scripts | **Selected: 最小改动覆盖最大范围** |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动集中在三个现有 skill 内部：
- tech-design：增加 handbook 生成分支（已有 api-handbook 作为参考模式）
- gen-contracts：Step 2 后增加 handbook 读取和锚点填充
- gen-test-scripts：Step 1 代码侦察后增加比对逻辑

### Resource & Timeline

单项 enhancement，改动点明确，无外部依赖。

### Dependency Readiness

api-handbook 已稳定运行。cli-handbook / page-map / screen-map 是新增文档类型，无前置依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Fact Table 代码侦察足以保证测试准确性 | 5 Whys: 为什么侦察到了 PUT 路由但测试仍用 POST？因为侦察结果没有与 Contract 做交叉验证 | Confirmed: 侦察存在但验证链缺失 |
| Contract 只需要语义描述，不需要技术细节 | XY Detection: 用户真正需要的是"测试脚本使用正确的 HTTP method"，而非"Contract 描述更详细" | Refined: Contract 需要技术锚点作为设计-实现的桥梁 |
| Web/Mobile surface 不需要类似机制 | Assumption Flip: 如果 Web 测试不知道测哪个页面，也会猜测 | Confirmed: 全 surface 都需要锚点 |

## Scope

### In Scope

- tech-design 增加 CLI/TUI cli-handbook、Web page-map、Mobile screen-map 自动生成
- Contract 模板 frontmatter 增加 `endpoint`、`command`、`page`、`screen` 字段及 `last_anchor_sync` 时间戳
- gen-contracts 从 handbook 填充锚点字段，包含 handbook 新鲜度检查（比对 handbook 与 tech-design 时间戳）
- eval-contract 评分规则增加技术锚点完整性检查，包含 handbook 内部一致性检查（设计阶段检测 endpoint 冲突）
- gen-test-scripts Step 1 增加交叉验证：分类为高置信度/低置信度/无法验证，以设计文档为准生成建议修复（用户确认后写入），设计-代码不一致时标记代码 bug
- 交叉验证输出 surface 覆盖报告，明确列出已验证和未验证的 surface 类型

### Out of Scope

- Mock 层级指导（task 执行 agent 行为，不在管道范围）
- 代码不存在时的强制交叉验证（先设计场景跳过，仅做格式校验）
- gen-journeys 改动（Journey 是语义层，不需要技术锚点）
- api-handbook 格式变更（保持现有格式）
- Contract 手动编辑后的锚点漂移检测（由 `last_anchor_sync` 时间戳辅助用户自行判断，初始版本不做自动检测）


### Phased Implementation Roadmap

| Phase | Scope | Rationale |
|-------|-------|-----------|
| Phase 1 | API surface (`endpoint` + `method`) | 已有 api-handbook 成熟模式，验证最小闭环 |
| Phase 2 | CLI/TUI surface (`command` + `subcommand`) | 复用 Phase 1 的交叉验证框架，扩展 cli-handbook |
| Phase 3 | Web + Mobile surface (`page` + `screen`) | page-map / screen-map 格式定义复杂度最高，最后实现 |

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|

| handbook 格式设计不当导致填充不准确 | M | M | 复用 api-handbook 的成熟格式模式，新增类型尽量对齐 |
| 建议修复仍有误（设计文档本身有误） | L | H | 用户确认环节作为最终防线，修复前展示 diff 供审阅 |
| 全 surface 覆盖导致 scope 膨胀 | M | L | 分批实现（Phase 1: API → Phase 2: CLI → Phase 3: Web/Mobile） |
| 交叉验证因设计文档滞后而误判 | M | H | handbook 新鲜度检查，过期时提示重新生成而非静默使用 |
| 部分 surface 有 handbook 导致虚假安全感 | M | M | 交叉验证报告明确列出已验证和未验证的 surface |
| 代码侦察不完整导致低置信度结果 | M | L | 结果分类为高/低/无法验证，低置信度不自动处理 |

## Success Criteria

consistency_check_result:
  status: pass
  pairs_checked: 15
  conflicts_found: 0

- [ ] API surface 的 Contract 100% 包含 `endpoint` 字段（当 api-handbook 存在时）
- [ ] CLI/TUI surface 的 Contract 100% 包含 `command` 字段（当 cli-handbook 存在时）

- [ ] gen-test-scripts 交叉验证能捕获 lesson 场景（POST vs PUT 不匹配），建议修复为 api-handbook 定义的 PUT，用户确认后写入
- [ ] 缺少 handbook 时管道正常运行（向后兼容，零中断），并提示用户建议生成 handbook
- [ ] 交叉验证输出 surface 覆盖报告，区分已验证和未验证的 surface
- [ ] 设计文档与代码实现不一致时，生成明确的代码 bug 标记报告

## Next Steps

- Proceed to `/write-prd` to formalize requirements
