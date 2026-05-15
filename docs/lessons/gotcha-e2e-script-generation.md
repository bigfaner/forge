# E2E 测试脚本生成：SKILL.md 流程结构性盲区导致测试与代码脱节

## Problem

生成的 e2e 测试脚本存在大量基础错误，导致需要多轮修复才能运行：
- 相对路径计算错误（多了一层 `../`）
- API 路径前缀假设错误（`/api/v1/` vs 实际的 `/v1/`）
- 端口号与实际运行服务不符
- 路由名称错误（`/items` vs `/main-items`，`/members/invite` vs `/members`）
- 测试用例引用了不存在的路由
- 测试依赖硬编码的预置数据，但数据库中没有这些数据
- UI testid 猜测错误（`map-card` vs 实际的 `map-card-${bizKey}`），24 个测试全部超时

## Root Cause

**不是 agent 不遵守指令，而是 SKILL.md 流程存在结构性盲区。**

### 第一层：原始版本没有 Code Reconnaissance 步骤

SKILL.md 最初的流程是从 "读 test cases" 直接到 "生成 spec files"，中间没有读源码的步骤。agent 完全按照指令执行——但指令本身缺少关键环节。

**Fix**：2026-04-29 加入 Step 1.5 Code Reconnaissance，要求读 router/config/handler/auth/CLI 源码并建 Fact Table。

### 第二层：Step 1.5 只覆盖后端，遗漏了 Frontend UI

Step 1.5 的 Required reads 表有 5 类源文件（Router、Config、API handlers、Auth、CLI），但没有 "Frontend UI components" 类别。UI 选择器验证依赖两条路径：
1. **sitemap.json**（Step 2）——但 sitemap 不含 `data-testid`，且新页面可能不在 sitemap 中
2. **Step 1.5 Fact Table** ——但没有 frontend 源码读取类别

当两条路径都失效时（新页面 + Step 1.5 不读前端），agent 标注 "provisional" 后继续生成——因为 SKILL.md 没有要求它在无法验证 testid 时停止。

**Fix**：在 Step 1.5 Required reads 表加入 "Frontend UI components" 行，要求 grep `data-testid` 并记录到 Fact Table。

### 第三层：sitemap 缺失路由时的 fallback 是空话

SKILL.md 写了 "use Fact Table DOM structure from Step 1.5"，但 Step 1.5 没有 frontend 读取——这个 fallback 引用了一个不存在的数据来源。

**Fix**：sitemap 缺失路由时 emit WARNING 并建议 re-run `/gen-sitemap`，同时用 Step 1.5 的 Fact Table 继续推断（现在 Fact Table 有 frontend 数据了）。

## Key Takeaway

1. **设计流程时必须覆盖完整的 fallback 链**——每一条 fallback 路径都必须有实际的数据来源，不能是空话
2. **"provisional" 不是安全的退出策略**——如果关键值无法验证，应该警告并建议修复前置条件，而不是标注后继续
3. **Frontend 和 Backend 应该对称地出现在侦察步骤中**——不能只读后端代码、不读前端代码就生成 UI 测试

## Timeline

| 日期 | 事件 |
|------|------|
| 2026-04-28 | 原始问题发生（pm-work-tracker），API/路径/端口/路由/数据全部猜错 |
| 2026-04-29 | Step 1.5 Code Reconnaissance 加入 SKILL.md，覆盖后端源码 |
| 2026-05-13 | UI testid 问题发生（pm-work-tracker milestones 页面），Step 1.5 不覆盖 frontend |
| 2026-05-14 | 修复：前端 testid 被改为静态值以适配测试（方向有误，应反过来） |
| 2026-05-15 | Step 1.5 加入 Frontend UI components 行，sitemap-missing fallback 补全 |

**反模式**：根据 PRD/设计文档推断实现细节 → 测试与实际代码脱节。
**更深的反模式**：SKILL.md 设计 fallback 路径但不提供 fallback 所需的数据来源 → agent 按"指令"执行却产出错误结果。
