---
created: 2026-04-28
author: faner
status: Draft
---

# Proposal: E2E 测试脚本生成质量提升

## Problem

e2e 测试生成 pipeline 的两个阶段都存在"推断而非验证"的问题：

1. **gen-test-cases**：从 PRD 提取 test case 时，Route/Target 字段凭 PRD 描述推断，未与实际路由交叉验证
2. **gen-test-scripts**：基于 test-cases 生成脚本时，路径/端口/行为全靠猜测，不读源码
3. **生成结构冗余**：每个 feature 独立生成 helpers.ts / package.json / node_modules，维护成本高

质量传递链：`PRD → gen-test-cases (Route 可能错) → gen-test-scripts (继承错误 + 新增猜测) → 脚本不可运行`

### Evidence

实际项目中（pm-work-tracker，22 个 feature、13 个含 testing scripts）观察到以下问题：

#### A. 脚本质量错误（推断 vs 实际值）

| 错误类型 | 错误值 | 正确值 | 来源 | 影响范围 |
|---------|--------|--------|------|---------|
| API 端口错误 | `localhost:8083` | `localhost:8080` | backend/config.yaml:L3 | bizkey-unification 等 feature |
| 默认密码错误 | `'password'` | `'admin123'` | backend/config.yaml:L23 | api-permission-test-coverage 等 feature |
| 前端端口不一致 | `localhost:5173` | `localhost:5174` | tests/e2e/config.yaml:L2 | api-permission-test-coverage |
| 硬编码开发者路径 | `/Users/fanhuifeng/Projects/Go/...` | 应从 config 读取 | api-permission-test-coverage/api.spec.ts:L235 | 任何人克隆后无法运行 |
| 路由名称错误 | `/items` | `/main-items` | router.go:L104 | 多个 feature |
| API 路径前缀错误 | `/api/v1/` | `/v1/` | router.go:L73 (`r.Group("/v1")`) | 多个 feature |

#### B. 结构性冗余

- **13 个 per-feature `helpers.ts`**（~95% 内容相同，但存在 drift）：
  - 端口不一致：有的写 `8080`，有的写 `8083`
  - 密码不一致：有的写 `'password'`，有的写 `'admin123'`
  - `loginViaUI` 实现分化：一份用硬编码中文（`账号`/`密码`/`登录`），一份用 regex
  - `runCli` 签名分化：大部分是 `(cmd, cwd?)`，一份增加了 `timeout?` 参数
- **12 个独立的 `package.json`**（各自管理依赖）
- `tests/e2e/` 已有共享基础设施（helpers.ts、package.json、config.yaml），但只有 3 个 API 测试和 2 个 CLI 测试使用

#### C. 毕业率极低

- 13 个 feature test suite 中只有 **1 个** 通过毕业（`db-dialect-compat`）
- `tests/e2e/api/` 下存在字节级重复的 spec 文件（`item-pool/api.spec.ts` 和 `regression/api.spec.ts` 完全相同）
- **0 个 UI 测试**在回归套件中，尽管 per-feature helpers 全部包含 Playwright 代码

### Urgency

每次生成低质量脚本 → run-e2e-tests 失败 → 多轮 error-fixer 修复 → 单 feature 平均消耗 3-5 轮 agent 交互（每轮 ~10k tokens）。结构性冗余（独立 node_modules）随 feature 数量增长线性膨胀。

## Proposed Solution

### 1. gen-test-cases：增加 Route Validation 步骤

gen-test-cases 的核心原则"PRD 是唯一输入源"不变——test case 内容仍从 PRD 提取。但在 Step 3（Classify & Generate）之后，增加一个**轻量校验**步骤：

```
Step 3.5: Route Validation (after generating test cases, before writing output)

1. Locate route definition files in the project
   - Search patterns: router.go, routes.ts, app.ts, *routes*.go
   - If not found → skip validation, add WARNING to test-cases.md
2. Cross-reference each test case's Route field against actual route definitions
3. If a route doesn't exist in source code:
   - Flag with ⚠️ in test-cases.md: "Route /items not found in router — verify path"
   - Suggest closest matching route if available
4. For API test cases: validate that Target's resource path matches route group prefix
5. Record validation results in a summary table at the end of test-cases.md
```

这不会改变 gen-test-cases 的"PRD-first"原则——test case 仍然从 PRD 提取，只是多了一层交叉验证。校验结果作为警告写入 test-cases.md，供 gen-test-scripts 参考。

### 2. gen-test-scripts：Code Reconnaissance 步骤

在 Step 1（读 test cases）和 Step 2（Resolve Sitemap）之间，插入强制的 **Code Reconnaissance** 步骤：

```
Step 1.5: Code Reconnaissance (MANDATORY for API/CLI tests, RECOMMENDED for UI tests)

1. Check test-cases.md for ⚠️ route warnings → use corrected routes from validation
2. Locate and read route definitions → extract actual paths and HTTP methods
   - Search patterns: router.go, routes.ts, routes/*.go, app.ts
   - Extract: path prefixes, route groups, middleware bindings
3. Read config files → confirm port, base_path
   - Search patterns: config.yaml, .env, config.json
4. For API tests: read middleware/handler code → confirm which validations exist
5. Verify relative paths with: node -e "console.log(require('path').resolve(__dirname, '..'))"
6. Record findings in a fact table used by Step 4 (Generate Spec Files)
```

Agent 必须用事实表中的值填充脚本，禁止推断。

### 3. VERIFY 标记（模板）

在模板关键值位置插入 `// VERIFY:` 注释，agent 生成时必须替换为 Code Reconnaissance 确认的实际值：

- API 模板：`// VERIFY: confirm auth endpoint path from router → replace /v1/auth/login`
- API 模板：`// VERIFY: confirm API base path prefix from config → replace /v1/`
- UI 模板：`// VERIFY: confirm login redirect URL from router → replace **/dashboard`
- helpers.ts：`// VERIFY: confirm default port from config → replace :3456 / :8080`

生成后自检：脚本中不得残留 `// VERIFY:` 标记。

### 4. 生成目标迁移

```
# Before (现状)
docs/features/<slug>/testing/scripts/   ← 每个feature独立生成
  helpers.ts, package.json, tsconfig.json, *.spec.ts, node_modules/

# After (改进)
tests/e2e/                              ← 共享基础设施
  helpers.ts, package.json, tsconfig.json, config.yaml, node_modules/
  <feature>/                            ← gen-test-scripts 生成到此处
    *.spec.ts
  <target>/                             ← 毕业后整合到此处
    *.spec.ts
```

- gen-test-scripts 输出到 `tests/e2e/<feature>/`，共享 `tests/e2e/` 的 helpers/config/deps
- 不再生成 per-feature 的 helpers.ts / package.json / tsconfig.json
- 毕业时从 `<feature>/` 整合到 `<target>/`（按 ui/api/cli 分类）

### 5. Before/After 示例

#### test-cases.md — Route Validation 输出

Before（当前，无校验）:
```markdown
## TC-001: Create item
- Type: API
- Route: POST /api/v1/items
- Target: /items
```

After（增加 Route Validation）:
```markdown
## TC-001: Create item
- Type: API
- Route: POST /api/v1/items
- Target: /items
- Status: ⚠️ Route `/api/v1/items` not found in router.go
  - Nearest match: `POST /v1/main-items` (router.go:L42)
  - Action: verify correct path before script generation

## TC-002: Invite member
- Type: API
- Route: POST /api/v1/members/invite
- Target: /members
- Status: ✅ Matched (router.go:L78)

---

### Route Validation Summary
| Test Case | Route Checked | Result | Matched Route |
|-----------|--------------|--------|---------------|
| TC-001 | POST /api/v1/items | ⚠️ Not found | POST /v1/main-items (suggested) |
| TC-002 | POST /api/v1/members/invite | ✅ Matched | POST /v1/members/invite |
```

#### gen-test-scripts — Code Reconnaissance 事实表 + 脚本输出

Before（当前，推断值）:
```typescript
// helpers.ts (generated by gen-test-scripts)
const API_BASE = 'http://localhost:8080/api/v1';  // guessed
```

```typescript
// item-create.spec.ts
const res = await fetch(`${API_BASE}/items`, {  // guessed from PRD
  method: 'POST',
  body: JSON.stringify({ name: 'test-item' })
});
```

After（Code Reconnaissance 事实表驱动）:

Fact table (generated during Step 1.5):
```markdown
### Code Reconnaissance Fact Table
| Key | Value | Source |
|-----|-------|--------|
| API_HOST | localhost | .env:L3 `SERVER_HOST=localhost` |
| API_PORT | 3456 | config.yaml:L12 `port: 3456` |
| API_PREFIX | /v1 | router.go:L15 `r.Group("/v1")` |
| ROUTE_CREATE_ITEM | POST /v1/main-items | router.go:L42 `r.Post("/main-items", ...)` |
| AUTH_METHOD | Bearer token in Authorization header | middleware/auth.go:L8 |
| VALIDATION_BIZKEY | Required, int, min=1 | handlers/items.go:L23 `binding:"required,min=1"` |
```

Generated script (values from fact table):
```typescript
// helpers.ts (generated by gen-test-scripts)
const API_BASE = 'http://localhost:3456/v1';  // Source: config.yaml:L12, router.go:L15
```

```typescript
// item-create.spec.ts
const res = await fetch(`${API_BASE}/main-items`, {  // Source: router.go:L42
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },  // Source: middleware/auth.go:L8
  body: JSON.stringify({ bizKey: 1, name: 'test-item' })  // bizKey min=1 from handlers/items.go:L23
});
```

## Quality Gate Summary

```
gen-test-cases
  └─ Step 3.5: Route Validation ← 校验 Route/Target vs 实际路由
       ↓
gen-test-scripts
  └─ Step 1.5: Code Reconnaissance ← 读源码建事实表
  └─ Step 4: VERIFY markers ← 关键值必须有源码依据
       ↓
run-e2e-tests
  └─ 执行验证通过的脚本
```

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零改动 | 每次生成需多轮修复；冗余加剧 | Rejected: 成本持续累积 |
| 只改 gen-test-scripts | 解决脚本生成问题 | test-cases.md 的错误仍然传递 | Rejected: 治标不治本 |
| 只改 gen-test-cases | 错误更早发现 | gen-test-scripts 仍会推断端口/行为细节 | Rejected: 只解决一半 |
| 两阶段质量门 + VERIFY + 目录迁移 | 完整解决质量链；消除冗余 | 涉及修改 gen-test-cases、gen-test-scripts、run-e2e-tests、graduate-tests 共 4 个 skill 的 8 个文件，需协调测试 | **Selected** |
| 扩展 config.yaml 增加 routes 段 | 结构化 | 新增维护负担；routes 变更需要手动同步 | Rejected: 维护成本转嫁给用户 |

## Scope

### In Scope

- gen-test-cases SKILL.md：增加 Step 3.5 Route Validation
- gen-test-cases templates/test-cases.md：增加 validation warnings 格式
- gen-test-scripts SKILL.md：增加 Step 1.5 Code Reconnaissance + 事实表要求
- gen-test-scripts 模板（playwright-ui.spec.ts, api.spec.ts）：增加 VERIFY 标记
- helpers.ts 模板：增加 VERIFY 标记
- gen-test-scripts SKILL.md 输出路径：从 `docs/features/<slug>/testing/scripts/` 改为 `tests/e2e/<feature>/`
- gen-test-scripts SKILL.md 中删除 per-feature helpers/package.json/tsconfig.json 生成指令
- run-e2e-tests SKILL.md：适配新路径（从 `tests/e2e/<feature>/` 执行）
- graduate-tests SKILL.md：适配新路径（`<feature>/` → `<target>/` 整合）

### Out of Scope

- breakdown-tasks 模板（task T-test-1/T-test-2 的描述可能需微调，但不在本 proposal 范围）
- 现有 feature 的迁移（手动处理，不自动化）
- run-e2e-tests 的执行逻辑重构（只适配路径）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent 找不到路由文件（项目结构特殊） | Medium | **High** — test cases/scripts proceed with unvalidated routes, downstream scripts inherit incorrect paths and fail at runtime | 明确 search patterns；找不到时在 test-cases.md 写 WARNING，agent 回退到现有行为（无校验） |
| VERIFY 标记被遗漏未替换 | Low | **Medium** — 脚本包含占位值，运行时立即报错但浪费一轮 run-e2e-tests 执行 | SKILL.md 增加自检规则：生成后 `grep '// VERIFY:' *.spec.ts` 残留标记 = incomplete，阻止提交 |
| 目录迁移影响 run-e2e-tests / graduate-tests | Medium | **High** — 已有 feature 的脚本路径不匹配，run-e2e-tests 找不到脚本导致全部 FAIL | 两个 skill SKILL.md 适配新路径；已有 feature 手动迁移（Out of Scope 内已标注） |
| gen-test-cases 的 Route Validation 增加执行时间 | Low | **Low** — 多读 1-2 个文件，agent 执行时间增加 <10s | 验证步骤是 OPTIONAL（找不到路由文件则跳过，不阻塞流程） |
| 共享 helpers.ts 的 feature 间冲突 | Low | **Medium** — 不同 feature 需要不同 config 时，helpers.ts 行为不一致 | helpers.ts 从 `tests/e2e/config.yaml` 读取配置，config 是项目级单一来源 |

## Success Criteria

- [ ] gen-test-scripts 模板文件 playwright-ui.spec.ts、api.spec.ts、helpers.ts 各包含至少 1 个 `// VERIFY:` 标记（验证方法：对每个模板文件执行 `grep -c '// VERIFY:' <template-path>` 结果 >= 1）
- [ ] gen-test-cases 输出的 test-cases.md 末尾包含 Route Validation Summary 表格（格式见 Before/After 示例），每个 Route 字段标注 ✅ 或 ⚠️
- [ ] gen-test-scripts 输出的每个 spec 文件附带 Code Reconnaissance Fact Table（格式见 Before/After 示例），且脚本中每个 URL/port/path 的值均可在 Fact Table 中找到对应行（验证方法：对脚本中 URL 执行 `grep -f <fact-table-keys>` 全部命中）
- [ ] 生成的脚本中 `grep -r '// VERIFY:' tests/e2e/` 返回 0 行（自检通过，无残留标记）
- [ ] `tests/e2e/` 目录存在共享的 helpers.ts + package.json + config.yaml，且不存在 `tests/e2e/<feature>/helpers.ts`（验证方法：`find tests/e2e -name helpers.ts | wc -l` = 1）
- [ ] gen-test-scripts 输出路径为 `tests/e2e/<feature>/*.spec.ts`（不是 `docs/features/<slug>/testing/scripts/`）
- [ ] graduate-tests 执行后，`tests/e2e/<target>/` 目录存在且包含从 `<feature>/` 迁移的 spec 文件，且 import 路径指向共享 helpers（验证方法：迁移后 spec 文件中 `from './helpers'` 不存在，应为 `from '../helpers'` 或 `from '../config'`）
- [ ] run-e2e-tests 从 `tests/e2e/` 执行，能发现并运行 `tests/e2e/<target>/*.spec.ts`（验证方法：run-e2e-tests 输出的 test list 包含 target 目录下的 spec 文件）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
