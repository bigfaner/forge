---
name: gen-test-cases
description: Generate structured test cases from PRD acceptance criteria. Classifies by type (UI/API/CLI) with full traceability to PRD sections.
---

# Gen Test Cases

从 PRD 验收标准生成结构化测试用例。

**核心原则**：PRD 是唯一输入源。每个测试用例必须可追溯到 PRD 中的具体验收标准，不发明 PRD 中不存在的验收标准。

<HARD-GATE>
本 skill 只生成测试用例文档（testing/test-cases.md），不生成可执行的测试脚本。
测试脚本的生成由 `/gen-test-scripts` skill 负责。
</HARD-GATE>

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

| 产物 | 缺失时提示 |
|------|-----------|
| `prd/prd-user-stories.md` | 先执行 `/write-prd` |
| `prd/prd-spec.md` | 先执行 `/write-prd` |
| `docs/sitemap/sitemap.json`（可选，仅 UI 测试） | 先执行 `/gen-sitemap` 可获得更精确的元素引用 |

**注意**：本 skill 既可以手动调用，也可以作为 `/breakdown-tasks` 追加的标准任务 T-test-1 被 agent 调用。

```bash
task feature
ls docs/features/<slug>/prd/prd-user-stories.md
ls docs/features/<slug>/prd/prd-spec.md
```

## When to Use

**Trigger:**
- User asks to "generate test cases" or "create test cases"
- User provides `/gen-test-cases` command
- After PRD is finalized, before or after implementation

**Skip:**
- No PRD exists yet (use `/write-prd` first)

## Workflow

```
1. Read PRD sources → 2. Extract AC → 3. Classify & generate → 4. Write test-cases.md
```

### Step 1: Read PRD Sources

Read all available PRD documents：

1. `prd/prd-user-stories.md` — primary source for acceptance criteria (Given/When/Then format)
2. `prd/prd-spec.md` — functional specs, scope, quality checks at the end
3. `prd/prd-ui-functions.md` — UI-specific criteria (if exists)

Also read `ui/ui-design.md` if it exists — provides component-level verification points for UI tests.

### Step 2: Extract Acceptance Criteria

From each source, extract every verifiable criterion：

**From user stories** (`prd-user-stories.md`):
- Each `Given/When/Then` block is one acceptance criterion
- Each story may have multiple AC blocks
- Preserve the story reference (e.g., "Story 1 / AC-1")

**From PRD spec** (`prd/prd-spec.md`):
- Quality check items at the end (checkboxes)
- Functional requirements in Section 5 (列表页, 按钮操作, 表单)
- Performance/security requirements if testable

**From UI functions** (`prd/prd-ui-functions.md`):
- Each UI function's behavior description
- Interaction requirements
- State requirements (loading, empty, error)

<EXTREMELY-IMPORTANT>
只提取 PRD 中**明确存在**的验收标准。禁止：
- 发明 PRD 中没有提到的测试场景
- 将模糊描述臆断为具体验收标准
- 遗漏任何明确的 Given/When/Then 条件
</EXTREMELY-IMPORTANT>

### Step 3: Classify & Generate Test Cases

For each extracted criterion, classify by type and generate a test case.

<HARD-RULE>
每个测试用例必须包含 `Target` 和 `Test ID` 字段：
- **Target**: `<type>/<page-or-resource>` (例如 `ui/login`、`api/auth`、`cli/deploy`)
- **Test ID**: `<target>/<title-slug>`，其中 title-slug = 标题小写 + 空格转连字符 + 去除标点
</HARD-RULE>

**Type classification rules:**

| Type | Indicators |
|------|-----------|
| **UI** | Page rendering, navigation, visual state, interactions, responsive behavior, component visibility, form input, modals, tabs, dropdowns |
| **API** | Endpoints, request/response, status codes, data contracts, HTTP methods, authentication headers |
| **CLI** | Commands, flags, output format, exit codes, arguments, stdin/stdout |

**Priority assignment:**
- **P0**: Criteria tied to core user stories or critical path
- **P1**: Criteria tied to secondary features or edge cases in core flow
- **P2**: Nice-to-have verifications, performance checks, edge cases

For each criterion, generate：

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y} or {UI Function Name}
- **Type**: UI | API | CLI
- **Target**: <type>/<page-or-resource>          ← e.g. ui/login, api/auth, cli/deploy
- **Test ID**: <target>/<title-slug>            ← e.g. ui/login/login-with-valid-credentials
- **Pre-conditions**: {What must be true before testing}
- **Route**: {Page route for UI tests}            ← e.g. /login, /settings
- **Element**: {Optional: sitemap element IDs}    ← e.g. E-001, L-003 (only if sitemap exists)
- **Steps**:
  1. {Step 1}
  2. {Step 2}
  ...
- **Expected**: {What the correct result looks like}
- **Priority**: P0 | P1 | P2
```

**Element 字段规则**：
- 仅当 `docs/sitemap/sitemap.json` 存在时生成
- 引用 sitemap 中的元素 ID（E-NNN 为页面元素，L-NNN 为布局元素）
- 列出测试步骤中直接操作的元素 ID，多个用逗号分隔
- 无 sitemap 时省略此字段，gen-test-scripts 会使用页面全部元素

<HARD-RULE>
**Numbering**: Start from TC-001, sequential. Group by type (UI first, then API, then CLI).

**Traceability**: 每个测试用例的 `Source` 字段必须指向 PRD 中的具体位置（Story 编号、Spec 章节号、UI Function 名称）。文件末尾必须包含完整的追溯表（TC ID → Source → Type → Target → Priority）。

**Target 推导规则**：
- UI 测试：`ui/<page-name>`（从 URL 或组件名推导，如 login 页 → `ui/login`）
- API 测试：`api/<resource>`（从端点推导，如 `/api/auth` → `api/auth`）
- CLI 测试：`cli/<command>`（从命令名推导，如 `task claim` → `cli/claim`）

**Test ID 生成规则**：`<target>/<title-slug>`，其中 title-slug = 标题小写 + 空格转连字符 + 去除标点符号。
</HARD-RULE>

### Step 4: Write Output

Read the template at `plugins/zcode/skills/gen-test-cases/templates/test-cases.md`.

Fill in:
- Frontmatter with feature slug, source references, generation date
- All test cases
- Traceability table at the end

Write to: `docs/features/<slug>/testing/test-cases.md`

Create the `testing/` directory if it doesn't exist.

## Overwrite Policy

If `testing/test-cases.md` already exists:
- **Overwrite without asking** — this skill regenerates from current PRD state
- The old file is replaced; PRD is the source of truth
- If user wants to preserve, they should commit the previous version first

## Related Skills

| Skill | Usage |
|-------|-------|
| `/write-prd` | Create PRD with acceptance criteria |
| `/gen-test-scripts` | Generate executable scripts from test cases |
| `/run-e2e-tests` | Execute test scripts and report results |
