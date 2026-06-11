---
created: "2026-06-09"
tags: [testing, architecture]
---

# 纯 Web 旅程生成零条 E2E 测试脚本

## Problem

经过 forge full 流水线（gen-journeys → gen-contracts → gen-test-scripts）后，一个定义了 7 个旅程的 feature，最终只有 2 个旅程生成了测试脚本。缺失的 5 个旅程包含 36 个正向路径 + 50 个边界场景，完全没有被自动化测试覆盖。

以 milestone-map 为例：生成了 7 个旅程，但只有 `milestone-lifecycle` 和 `milestone-map-lifecycle` 获得了 API 和 E2E 测试脚本。其余 5 个旅程（`item-milestone-binding`、`item-list-milestone-integration`、`milestone-map-visualization`、`milestone-item-management`、`read-only-milestone-access`）虽然拥有完整的旅程定义和评估报告，却没有任何测试脚本。

## Root Cause

流水线存在结构性断层：**gen-contracts 只为拥有 API surface 的旅程生成契约，gen-test-scripts 只从契约生成脚本**。

逐层追溯因果链：

**Level 1 — 7 个旅程中只有 2 个获得了测试脚本。**
gen-test-scripts 任务只为 `milestone-lifecycle` 和 `milestone-map-lifecycle` 生成了脚本。

**Level 2 — 只有 2 个旅程拥有可喂给 gen-test-scripts 的契约文件。**
gen-contracts 任务创建了 14 个契约步骤文件，但全部在这同样的 2 个旅程下。其余 5 个旅程的 `contracts/` 目录是空的。

**Level 3 — gen-contracts 跳过了纯 Web 旅程。**
被覆盖的 2 个旅程配置为 `surface_types: ["web", "api"]`（双 surface）。未覆盖的 5 个旅程配置为 `surface_types: ["web"]`（纯 Web）。gen-contracts 按步骤为 API 端点生成契约——它没有机制能为纯视觉/交互场景生成契约。

**Level 4 — 流水线架构将"契约"等同于"API 契约"，而非"可测试结果规格"。**
六维度契约格式（前置条件、输入、输出、状态、副作用、不变量）以 API 为中心。纯 Web 旅程描述的是用户交互（点击、导航、筛选、拖拽），无法自然映射到该格式。

**Level 5 — 纯 Web 旅程没有从 journey 到 test-script 的替代路径。**
流水线是：`journey → contract → test-script`。不存在 `journey → test-script` 的直达路径。一旦某个旅程未能生成契约，它就被静默地从整个下游流水线中丢弃。

**Level 6 — 质量门禁未检测旅程覆盖缺口。**
`forge quality-gate` 执行 compile/fmt/lint → 单元测试 → 回归测试。它验证的是已生成脚本能编译且通过，但从不检查"我们是否为所有已定义旅程都生成了脚本？"。`gen-contracts` 的验收标准写的是"每个 Journey 至少生成 1 个 Contract 文件"，但这一条没有被强制执行——任务在只覆盖了 7 个旅程中 2 个的情况下就成功完成了。

## Solution

### 即时应对方案
gen-test-scripts 完成后，手动审计旅程覆盖率：
```bash
# 列出拥有 journey.md 的旅程
ls docs/features/<slug>/testing/*/journey.md
# 列出拥有测试脚本的旅程
ls tests/*/journey-name*/
# 对比——任何有 journey.md 但没有对应 tests/ 目录的旅程即为缺口
```

### Contract 适配性分析（按 Surface 执行模型）

分界线是执行模型，不是 surface 类型：

| Surface | 交互模型 | Contract 适配 | 原因 |
|---------|----------|--------------|------|
| API | 请求-响应 | 原生适配 | 输入输出都是结构化数据 |
| CLI | 子进程调用 | 原生适配 | 同 API，exit code + stdout |
| TUI | 有状态管道 | 可适配，摩擦可控 | 纯文本输出，按 step 可拆分 |
| Web | 用户交互序列 | **不适用** | Input 是 DOM 动作，Output 是视觉状态 |
| Mobile | 用户交互序列 | **不适用** | 同 Web，Maestro 是交互脚本 |

- **协议级（API、CLI）**：Contract 六维度是自然抽象，输入是结构化数据，输出是结构化响应，保留现有路径。
- **准协议级（TUI）**：forge 模型下 TUI 测试是 `subprocess + stdin pipe`，不涉及真正的终端渲染。单步交互适配良好，多步有状态会话有摩擦但可控——终端输出是纯文本，无需处理 DOM/CSS。
- **交互级（Web、Mobile）**：Input 是用户动作（点击、滑动），Output 是视觉状态（元素可见、文本匹配），Side-effect 是页面导航。将六维度翻译过去等于在写测试脚本本身——中间层不增加信息量，是仪式而非实质。journey.md 已经包含了生成 E2E 脚本所需的全部信息（用户动作步骤、期望页面状态、边界场景、业务不变量）。

### 流水线修复（需要修改 forge skill）

gen-test-scripts 应按 surface 执行模型路由生成路径：

```
协议级 surface (api, cli)   → contract → test script（现有路径，保留）
交互级 surface (web, mobile) → journey.md → test script（直达路径，缺失）
准协议级 (tui)               → 可走 contract，也可直达
```

1. **增加 journey.md 直达路径**：当旅程的 surface_types 只包含 web/mobile 时，gen-test-scripts 直接从 `journey.md` 的 happy-path + edge cases 生成 Playwright/Maestro 脚本，不要求 contract 文件存在。
2. **增加覆盖率完整性门禁**：gen-test-scripts 完成后，将 `testing/*/journey.md` 与 `tests/**/` 目录对比。缺少测试脚本的旅程应导致失败或告警。

### 任务设计修复
保持现有 2 个任务（per surface），但在任务定义中强制加入覆盖率自检步骤：任务完成后枚举所有匹配该 surface 的 journey，检查每个 journey 是否有对应的测试脚本，遗漏则 FAIL。一个 journey 一个任务会产生 10+ 个任务，管理成本远超收益——问题不是粒度太粗，而是缺少完整性校验。

## Reusable Pattern

**从多阶段流水线生成测试时，始终包含一个"覆盖率审计"步骤：将输入（已定义旅程数）与输出（已生成测试脚本数）进行对比，发现缺口即失败。**

模式：`count(输入) == count(输出)` 是最低完整性检查。缺少这一步时，部分生成看起来是成功的——因为流水线验证的是已生成内容的质量，而非本应生成内容的覆盖率。

## Example

milestone-map 的情况：
- 输入：7 个旅程定义在 `testing/` 中
- 输出：仅 2 个旅程有测试脚本
- 覆盖缺口：5 个旅程（71%）零测试脚本——对质量门禁不可见

## Related Files

- docs/features/milestone-map/testing/*/journey.md（7 个旅程定义）
- docs/features/milestone-map/testing/*/contracts/（仅 2 个有内容）
- tests/backend/milestone-*/（仅 2 个旅程目录）
- tests/frontend/milestone-*/（仅 2 个旅程目录）
- docs/features/milestone-map/tasks/gen-contracts.md
- docs/features/milestone-map/tasks/gen-test-scripts-backend.md
- docs/features/milestone-map/tasks/gen-test-scripts-frontend.md
