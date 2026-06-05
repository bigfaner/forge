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
4. **gen-test-scripts 交叉验证并自动修复**：将 Fact Table（代码侦察）与 Contract frontmatter 比对，不匹配时以设计文档为准自动修复 Contract；若设计文档与代码实现不一致，标记为代码 bug

### Innovation Highlights

将测试管道从"LLM 推断技术细节"升级为"设计文档驱动的锚点验证"。核心洞察：Contract 是设计意图的规格说明，技术锚点应该来自设计阶段确定的接口定义，而非从代码中逆向提取或让 LLM 猜测。交叉验证以设计文档为 authority source，设计-实现不一致时定位为代码 bug 而非测试问题。

## Requirements Analysis

### Key Scenarios

- **新功能先设计后实现**：tech-design 生成 handbook → gen-contracts 填充锚点 → 代码实现 → gen-test-scripts 交叉验证（代码存在时）
- **已有功能补全测试**：gen-contracts 从现有 handbook 填充锚点 → gen-test-scripts 交叉验证发现不匹配 → 自动修复 Contract
- **设计文档与代码不一致**：交叉验证发现 handbook 说 PUT 但代码是 POST → 自动修复 Contract 以 handbook 为准，标记代码 bug
- **Handbook 不存在**：gen-contracts 阶段跳过锚点填充，gen-test-scripts 发现缺少锚点时降级为 Fact Table 推断（向后兼容）

### Non-Functional Requirements

- 向后兼容：缺少 handbook 或锚点字段时，管道不中断，降级为现有行为
- 性能影响：交叉验证在 gen-test-scripts Step 1（代码侦察）中执行，无额外网络或 IO 开销

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
| **Contract Technical Anchors** | 设计文档 | 完整信息链、自动修复、全 surface 覆盖 | 需要扩展 tech-design 和 gen-test-scripts | **Selected: 最小改动覆盖最大范围** |

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
- Contract 模板 frontmatter 增加 `endpoint`、`command`、`page`、`screen` 字段
- gen-contracts 从 handbook 填充锚点字段
- eval-contract 评分规则增加技术锚点完整性检查
- gen-test-scripts Step 1 增加交叉验证：以设计文档为准自动修复 Contract，设计-代码不一致时标记代码 bug

### Out of Scope

- Mock 层级指导（task 执行 agent 行为，不在管道范围）
- 代码不存在时的强制交叉验证（先设计场景跳过，仅做格式校验）
- gen-journeys 改动（Journey 是语义层，不需要技术锚点）
- api-handbook 格式变更（保持现有格式）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| handbook 格式设计不当导致填充不准确 | M | M | 复用 api-handbook 的成熟格式模式，新增类型尽量对齐 |
| 自动修复覆盖了正确的 Contract（设计文档本身有误） | L | H | 修复前保存原始值到 Contract 的注释中，可回溯 |
| 全 surface 覆盖导致 scope 膨胀 | M | L | 每种 surface 的锚点字段独立、互不影响，可分批实现 |

## Success Criteria

consistency_check_result:
  status: pass
  pairs_checked: 15
  conflicts_found: 0

- [ ] API surface 的 Contract 100% 包含 `endpoint` 字段（当 api-handbook 存在时）
- [ ] CLI/TUI surface 的 Contract 100% 包含 `command` 字段（当 cli-handbook 存在时）
- [ ] gen-test-scripts 交叉验证能捕获 lesson 场景（POST vs PUT 不匹配），自动修复为 api-handbook 定义的 PUT
- [ ] 缺少 handbook 时管道正常运行（向后兼容，零中断）
- [ ] 设计文档与代码实现不一致时，生成明确的代码 bug 标记报告

## Next Steps

- Proceed to `/write-prd` to formalize requirements
