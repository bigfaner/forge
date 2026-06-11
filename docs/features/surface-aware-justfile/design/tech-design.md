---
created: 2026-05-25
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Surface-Aware Justfile

## Overview

本特性将 Forge 的测试编排模型从"语言感知 + 固定枚举"迁移到"surface 感知 + 用户自定义 key"。核心变更分三线并进：

1. **Go 数据模型迁移**（prompt.go / task 包）：`Scope` → `SurfaceKey`，`TestType` → `SurfaceType`，`resolveScope()` 删除。涉及 `task/types.go`、`autogen.go`、`frontmatter.go`、`build.go`、`prompt.go` 五个核心文件
2. **Skill 文档体系**（init-justfile / run-tests）：新增 10 个 surface 规则文件（每个 skill 5 个），两个 SKILL.md 大幅重写
3. **CLI 输出增强**（`forge surfaces` 命令已存在）：新增 `--json` flag 输出结构化 JSON，供 skill 通过 Bash 工具消费

设计遵循三阶段严格顺序：数据模型 → 上游适配 → 下游消费+清理。同阶段内可并行。

## Architecture

### Layer Placement

Forge CLI 层（Go 数据模型 + prompt 合成）+ Forge Plugin 层（Skill 文档 + 规则文件）

### Component Diagram

```
+---------------------------+
| forge surfaces CLI       |  (已存在: 命令+3模式; 新增: --json flag)
|  -- ReadSurfaces()       |  (已存在，不变)
|  -- MatchSurface(path)   |  (需修改: 返回 SurfaceMatch{Key,Type})
+-------------+-------------+
              | JSON output
              v
+---------------------------+     +---------------------------+
| breakdown-tasks / quick-  |     | init-justfile SKILL.md    |
| tasks SKILL.md            |     |  + rules/surfaces/<type>  |
|  填充 surface-key/type    |     |  生成差异化配方           |
+-------------+-------------+     +-------------+-------------+
              | task frontmatter                |
              | surface-key + surface-type      |
              v                                 |
+---------------------------+                 |
| Task Go struct            |                 |
|  SurfaceKey (string)      |                 |
|  SurfaceType (string)     |                 |
+-------------+-------------+                 |
              |                                 |
              v                                 v
+---------------------------+     +---------------------------+
| prompt.go renderTemplate()|     | run-tests SKILL.md        |
|  {{SURFACE_KEY}} direct   |     |  + rules/surfaces/<type>  |
|  {{TEST_TYPE_ARG}} direct |     |  调度器模式编排           |
+---------------------------+     +---------------------------+
```

### Dependencies

- 内部：`forgeconfig.ReadSurfaces()`（已存在）+ `MatchSurface()`（需修改签名，见 Interface 1a）
- 内部：`task` 包所有相关结构体
- 内部：`internal/cmd/surfaces.go`（需新增 `--json` flag 和 JSON 输出逻辑，见 Interface 1b）
- 外部：just >= 1.4.0（`[linux]`/`[windows]` recipe attribute）
- 无新外部依赖

### Phase Component Map

**阶段 1 — 数据模型与查询基础设施（无外部依赖）**：

| 文件 | 变更内容 | 对应 Interface/Model |
|------|---------|---------------------|
| `pkg/forgeconfig/match.go` | `MatchSurface()` 签名改为返回 `SurfaceMatch{Key, Type}` | Interface 1a |
| `internal/cmd/surfaces.go` | 新增 `--json` flag + 三模式 JSON 输出逻辑 | Interface 1b |
| `task/types.go` | Scope→SurfaceKey，新增 SurfaceType | Model 1, Model 5 |
| `task/autogen.go` | TestType→SurfaceType，传播链更新 | Model 2 |
| `task/frontmatter.go` | 新增 SurfaceKey/SurfaceType 字段 | Model 3 |
| `task/add.go` | 新增 SurfaceKey/SurfaceType 字段 | Model 4 |
| `prompt.go` | resolveScope() 删除，直接读 SurfaceKey | Model 1 |
| `task/migrate.go` | 新增 `CheckLegacyScope()` 共享函数 + `forge task migrate` 子命令 | Migration Notes |

**阶段 2 — 上游组件适配（依赖阶段 1）**：

| 文件 | 变更内容 | 对应 Interface/Model |
|------|---------|---------------------|
| `skills/breakdown-tasks/SKILL.md` | 任务生成流程新增：调用 `forge surfaces --json <path>` 获取 SurfaceKey+SurfaceType，写入 frontmatter | Interface 1b |
| `skills/breakdown-tasks/rules/scope-to-surface-key.md` | 新增规则文件：指导 agent 将 scope 枚举值映射为 surface-key（调用 CLI 而非硬编码映射） | Interface 1b |
| `skills/quick-tasks/SKILL.md` | 同 breakdown-tasks，任务生成时调用 `forge surfaces --json` 填充双字段 | Interface 1b |
| `task/add.go` | `AddTaskOpts` 新增 SurfaceKey/SurfaceType 字段，从源任务复制到新任务 | Model 4 |
| `quality-gate` 相关逻辑 | fix-task 流程中调用 `forge surfaces <path> --json` 推断 surface-key，替代原 scope 推断 | Interface 1a, 1b |

**阶段 3 — 下游消费与清理（依赖阶段 2）**：

| 文件 | 变更内容 | 对应 Interface/Model |
|------|---------|---------------------|
| `skills/init-justfile/SKILL.md` | 重写：增加 surface 检测流程，按 surface-type 加载对应规则文件生成差异化配方 | Interface 2 |
| `skills/init-justfile/rules/surfaces/web.md` | 新增：web surface 配方生成规则（完整内容见 Interface 2 示例） | Interface 2 |
| `skills/init-justfile/rules/surfaces/api.md` | 新增：api surface 配方生成规则 | Interface 2 |
| `skills/init-justfile/rules/surfaces/cli.md` | 新增：cli surface 配方生成规则（无 dev/probe 步骤） | Interface 2 |
| `skills/init-justfile/rules/surfaces/tui.md` | 新增：tui surface 配方生成规则（无 dev/probe 步骤） | Interface 2 |
| `skills/init-justfile/rules/surfaces/mobile.md` | 新增：mobile surface 配方生成规则 | Interface 2 |
| `skills/run-tests/SKILL.md` | 重写：调度器模式编排，按 surface-type 加载编排规则 | Interface 2 |
| `skills/run-tests/rules/surfaces/web.md` | 新增：web 编排序列（dev→probe→test→teardown） | Interface 2 |
| `skills/run-tests/rules/surfaces/api.md` | 新增：api 编排序列 | Interface 2 |
| `skills/run-tests/rules/surfaces/cli.md` | 新增：cli 编排序列（test→teardown） | Interface 2 |
| `skills/run-tests/rules/surfaces/tui.md` | 新增：tui 编排序列（test→teardown） | Interface 2 |
| `skills/run-tests/rules/surfaces/mobile.md` | 新增：mobile 编排序列 | Interface 2 |
| `skills/*/templates/*.md`（共 16 个 prompt 模板） | `{{SCOPE}}` → `{{SURFACE_KEY}}`，`{{TEST_TYPE}}` → `{{TEST_TYPE_ARG}}` 变量值域同步 | Model 1 |

**Prompt 模板完整清单**（`forge-cli/pkg/prompt/data/`，共 18 个使用 `{{SCOPE}}`，1 个同时使用 `{{TEST_TYPE_ARG}}`）：

| 模板文件 | `{{SCOPE}}` | `{{TEST_TYPE_ARG}}` |
|---------|:-----------:|:-------------------:|
| coding-enhancement.md | ✓ | |
| coding-cleanup.md | ✓ | |
| code-quality-simplify.md | ✓ | |
| coding-feature.md | ✓ | |
| coding-fix.md | ✓ | |
| coding-refactor.md | ✓ | |
| doc-summary.md | ✓ | |
| eval-contract.md | ✓ | |
| eval-journey.md | ✓ | |
| gate.md | ✓ | |
| fix-record-missed.md | ✓ | |
| test-gen-journeys.md | ✓ | |
| test-gen-contracts.md | ✓ | |
| test-gen-scripts.md | ✓ | ✓ |
| test-verify-regression.md | ✓ | |
| test-run.md | ✓ | |
| validation-code.md | ✓ | |
| validation-ux.md | ✓ | |

**Skill 模板清单**（使用 `{{SCOPE}}`，共 3 个）：

| 模板文件 | 说明 |
|---------|------|
| breakdown-tasks/templates/task.md | 任务 frontmatter 模板，`scope: "{{SCOPE}}"` → `surface-key: "{{SURFACE_KEY}}"` |
| quick-tasks/templates/task.md | 同上 |
| test-guide/templates/convention-template.md | `domains: [testing, {{SCOPE}}]` → `domains: [testing, {{SURFACE_KEY}}]` |

替换方式：`prompt.go` 的 `renderTemplate()` 中 `{{SCOPE}}` 替换逻辑直接改为 `{{SURFACE_KEY}}`，`{{TEST_TYPE_ARG}}` 保持不变（由 `extractTestTypeArg()` → 直接读 `task.SurfaceType` 生成）。模板文件内容中所有 `{{SCOPE}}` 字面量全局替换为 `{{SURFACE_KEY}}`。
| `prompt.go` 死代码清理 | 删除 `extractTestTypeArg()`、`genScriptBases()` 等已被 surface 机制替代的辅助函数 | Model 1 |

## Interfaces

### Interface 1: `forge surfaces [--json] [path] [--types]`

#### 1a. Go 层 API 变更：`forgeconfig/match.go`

当前 `MatchSurface()` 仅返回 surface-type（map value），不返回 surface-key（map key）。需修改签名以同时返回两者：

```
// 现有签名（仅返回 type）:
func MatchSurface(surfaces map[string]string, query string) (string, error)

// 新签名（返回 key + type）:
// 定义在 pkg/forgeconfig 包中；消费方：surfaces.go（CLI 输出）、autogen.go（任务生成传播）
type SurfaceMatch struct {
    Key   string  // surface-key（map key，如 "admin-panel"）
    Type  string  // surface-type（map value，如 "web"）
}
func MatchSurface(surfaces map[string]string, query string) (SurfaceMatch, error)
```

标量形式（单 key "."）时 `Key` 返回 `"."`，调用方按需处理。现有调用方仅 `surfaces.go` 的 `runSurfacesQuery`，需适配新签名从 `SurfaceMatch.Type` 取值。

#### 1b. CLI `--json` flag 规格（`internal/cmd/surfaces.go`）

在现有 cobra command 上注册 `--json` bool flag：

```
surfacesCmd.Flags().BoolVar(&jsonFlag, "json", false, "output in JSON format")
```

`--json` flag 与三种子调用模式的组合行为：

| 子调用模式 | 触发条件 | `--json` 输出格式 |
|-----------|---------|------------------|
| **列表模式** | 无 path 参数，无 `--types` | `{"surfaces": [{"key": "admin-panel", "type": "web"}, {"key": "payment-service", "type": "api"}]}` |
| **查询模式** | 有 path 参数 | `[{"key": "admin-panel", "type": "web"}]`；无匹配时 `[]`，exit 0 |
| **类型模式** | `--types` flag | `{"types": ["api", "web"]}` |

```
// 列表模式 --json 示例:
{"surfaces": [{"key": "admin-panel", "type": "web"}, {"key": "payment-service", "type": "api"}]}

// 查询模式 --json 示例（有匹配）:
[{"key": "admin-panel", "type": "web"}]

// 查询模式 --json 示例（无匹配）:
[]

// 类型模式 --json 示例:
{"types": ["api", "web"]}

// surfaces 配置缺失时（所有模式）:
stderr: {"error": "no surface configured; run `forge init` to configure surfaces"}
exit 1

// 无 --json 时保持现有文本输出行为不变
```

**`--json` 模式 stderr 格式覆盖声明**：当 `--json` flag 激活时，所有输出（包括 stdout 和 stderr）统一使用结构化 JSON 格式。这是对 TECH-error-handling-001 规定的 `<context>: <specific-detail>` 纯文本 stderr 格式的**显式例外**。理由：`--json` 的消费者是机器（skill 通过 Bash 工具解析），需要保证 stderr 同样可被 JSON 解析器无歧义消费，避免混合纯文本与 JSON 导致解析失败。`--json` 模式下所有错误路径必须通过 `json.NewEncoder(cmd.ErrOrStderr()).Encode()` 输出 `{"error": "..."}` 格式，不得使用 `fmt.Fprintf(os.Stderr, ...)`。

实现模式：在 `runSurfacesList`/`runSurfacesQuery`/`runSurfacesTypes` 中增加 `if jsonFlag` 分支，使用 `json.NewEncoder(cmd.OutOrStdout()).Encode()` 序列化输出，复用现有逻辑的 match/list 结果。

### Interface 2: Surface 规则文件格式

文件路径：`rules/surfaces/<type>.md`，消费者为 init-justfile（配方生成）和 run-tests（编排序列）。每种 surface 类型对应一个独立规则文件。以下为 `rules/surfaces/web.md` 的完整示例，其余 4 种 surface 类型（api / cli / tui / mobile）参照此模板编写。

```markdown
# Surface: web

## 编排序列

| 步骤 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|---------|---------|---------|---------|
| dev  | 服务启动成功，等待就绪 | 启动失败（依赖缺失/端口占用） | — | 进入 probe |
| probe | 健康检查通过 | 健康检查超时（服务未就绪） | — | 进入 test |
| test | 测试通过 | 测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | 清理完成 | 清理失败（残留进程） | — | 结束 |

注意事项：
- dev 失败时**不继续**后续步骤，直接 teardown 并退出
- probe 最多重试 3 次，间隔 5 秒；3 次均失败视为退出码 1
- test 退出码 2 允许重跑，skill 应提示用户 "测试环境异常，建议重试"

## 配方调用契约

| 配方名 | just 签名 | 退出码 0 语义 | 退出码 1 语义 |
|--------|----------|--------------|--------------|
| web-dev | `just web-dev` | 开发服务器就绪，监听端口 | 启动失败，stderr 含错误详情 |
| web-probe | `just web-probe` | HTTP 健康检查返回 2xx | 连接拒绝或超时 |
| web-test | `just web-test` | 所有测试用例通过 | 至少一个测试失败 |
| web-teardown | `just web-teardown` | 进程终止，端口释放 | 进程残留或清理异常 |
| web | `just web` | 聚合配方：dev→probe→test→teardown 完整流程 | 任一子步骤失败 |

实现约束：
- 每个配方必须支持 `[linux]` 和 `[windows]` 双平台变体
- `web` 聚合配方按编排序列顺序调用子配方，遇到非零退出码立即中断
- `web-teardown` 必须用 `just --dry-run` 验证语法

## journey 过滤策略

| journey 标签 | 匹配规则 | 说明 |
|-------------|---------|------|
| `@web` | 精确匹配 | web surface 的专用 journey |
| `@e2e` | 精确匹配 | 端到端测试，归入 web surface |
| `@smoke` | 精确匹配 | 冒烟测试，归入 web surface |
| 其他 | 忽略 | 非 web 相关 journey 不由本规则处理 |
```

其余 surface 类型差异点：
- **api**：dev 步骤启动 API 服务；probe 使用 HTTP `GET /health`；支持 `@api` journey
- **cli**：无 dev/probe 步骤，仅有 test + teardown；配方名前缀为 `cli-`；无聚合配方
- **tui**：与 cli 相同模式（无 dev/probe）；支持 `@tui` journey
- **mobile**：dev 步骤启动模拟器；probe 使用 appium 健康检查；支持 `@mobile` journey

## Data Models

### Model 1: Task struct (task/types.go)

```
Task {
    // ... 现有字段保持不变 ...
    SurfaceKey  string    // 新增：用户自定义 surface 标识（如 "admin-panel"），空值表示跨 surface
    SurfaceType string    // 新增：surface 类型枚举（web/api/cli/tui/mobile），空值表示未知
    // Scope 字段删除，不保留兼容层
}
```

JSON tag: `"surface-key,omitempty"` / `"surface-type,omitempty"`

### Model 2: AutoGenTaskDef struct (task/autogen.go)

```
AutoGenTaskDef {
    // ... 现有字段保持不变 ...
    SurfaceKey  string    // 新增（原 Scope 字段删除后新增）
    SurfaceType string    // 新增（原 TestType 字段删除后新增）
    // Scope 字段删除
    // TestType 字段删除
}
```

### Model 3: FrontmatterData struct (task/frontmatter.go)

```
FrontmatterData {
    // ... 现有字段保持不变 ...
    SurfaceKey  string    // 新增
    SurfaceType string    // 新增
    // Scope 字段删除
}
```

### Model 4: AddTaskOpts struct (task/add.go)

```
AddTaskOpts {
    // ... 现有字段保持不变 ...
    SurfaceKey  string    // 新增：从源任务继承
    SurfaceType string    // 新增：从源任务继承
}
```

### Model 5: TaskState struct (task/types.go)

```
TaskState {
    // ... 现有字段保持不变 ...
    SurfaceKey  string    // 新增
    SurfaceType string    // 新增
    // Scope 字段删除
}
```

## Error Handling

### Error Types & Codes

遵循 BIZ-error-reporting-001 统一 exit code 语义。

**Go 层**:

| 场景 | 行为 | Exit Code |
|------|------|-----------|
| `forge surfaces --json` surfaces 配置缺失 | stderr JSON 错误 + 恢复提示 | 1 (retryable) |
| `forge surfaces --json` 路径无匹配 | stdout `[]` + exit 0 | 0 |
| `build.go` 解析任务 frontmatter 缺 surface-key | 记录 warning 日志，赋空值 | 0 |
| `build.go` 解析旧 frontmatter 含 `scope` 但无 `surface-key` | **阻塞错误**：stderr 输出迁移提示，返回 exit 2 | 2 (blocking) |

### Migration Notes: Scope 字段删除

本设计**不保留 Scope 向后兼容层**，这是经过权衡的显式设计决策。理由：

1. 双字段维护负担高，查询逻辑复杂度随兼容期线性增长
2. Scope 固定枚举（frontend/backend）与用户自定义 surface-key 语义冲突
3. 迁移通过阻塞检查强制执行：`build.go` 检测到旧 `scope` 字段时返回 exit 2，用户必须重跑 `breakdown-tasks`/`quick-tasks` 后才能继续操作

**对 PRD Story 4 AC 的影响**：PRD Story 4 中 AC 条目"旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行"将更新为：旧任务文件需通过重新运行 `breakdown-tasks`/`quick-tasks` 获取 `surface-key`/`surface-type` 字段。这是 PRD 层面的需求变更，需在实现前同步更新 PRD。

**迁移验证步骤**（Phase 1 内执行）：迁移检查提取为 `task` 包的共享函数 `CheckLegacyScope(tasks []Task) error`，供所有任务读取路径调用。初始覆盖范围：

| 读取路径 | 入口函数 | 迁移检查调用点 |
|---------|---------|--------------|
| `build.go` 加载 `index.json` | `BuildIndex()` | 加载后遍历 tasks 列表 |
| `forge task list` | `ListTasks()` | 列表前检查 |
| `forge task show` | `ShowTask()` | 显示前检查单个 task |
| `forge task add` | `AddTask()` | 添加前检查源 task |

函数行为：扫描含 `scope` 字段但不含 `surface-key` 的任务，若存在此类任务，输出错误信息至 stderr（格式：`migration required: found N tasks with legacy 'scope' field but no 'surface-key' — run 'forge breakdown-tasks' or 'forge quick-tasks' to regenerate tasks`），返回 exit code 2（阻塞错误）。此检查确保升级后不存在数据不一致的任务，防止 run-tests 因缺失 surface-key 而加载错误的编排策略。

**CI 升级影响与缓解**：阻塞式迁移意味着 Forge CLI 升级后、任务重新生成前，所有依赖 task 读取的命令（包括 CI 中的 `forge run-tests`）将返回 exit 2。为缓解 CI 中断，提供 `forge task migrate` 子命令：扫描现有 `index.json`，将 `scope` 字段值通过 `forge surfaces --json <path>` 映射为 `surface-key` + `surface-type`，就地更新 `index.json` 和对应 frontmatter 文件。CI 升级流程：`升级 Forge CLI → forge task migrate → 正常操作`。`forge task migrate` 在 Phase 1 中与数据模型变更同步实现。

**Skill 层**（LLM agent 执行）:

PRD Error Handling Paths 表定义了 7 种场景。Skill 层通过 SKILL.md 中 HARD-GATE 规则约束 agent 行为。

### Propagation Strategy

- Go 层错误通过 exit code + stderr 传播给 skill（LLM agent 读取）
- Skill 层错误通过 just 命令 exit code 和 SKILL.md HARD-GATE 规则控制

## Cross-Layer Data Map

| Field | Go Struct | Task Frontmatter | Skill 消费 |
|-------|-----------|------------------|------------|
| surface-key | `Task.SurfaceKey` (string) | `surface-key: admin-panel` | skill 读取 task frontmatter |
| surface-type | `Task.SurfaceType` (string) | `surface-type: web` | skill 加载 `rules/surfaces/<type>.md` |

## Integration Specs

No existing-page integrations — not applicable.

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| Go CLI | Unit | `go test` | Task struct 序列化/deserialization、SurfaceKey/SurfaceType 字段读写 | 80% |
| Go CLI | Unit | `go test` | build.go frontmatter 解析（含/不含 surface-key/type） | 80% |
| Go CLI | Unit | `go test` | autogen.go 传播链（SurfaceKey/SurfaceType 从 AutoGenTaskDef → Task → frontmatter） | 80% |
| Go CLI | Unit | `go test` | renderTemplate() 的 {{SURFACE_KEY}}/{{TEST_TYPE_ARG}} 替换 | 80% |
| Go CLI | Unit | `go test` | forge surfaces --json 输出格式（有匹配/无匹配/无配置） | 80% |
| Skill | Manual | init-justfile | 5 种 surface 各生成一次 justfile，验证配方签名和内容 | N/A |
| Skill | Manual | run-tests | web/api 项目端到端编排流程（dev→probe→test→teardown） | N/A |
| Skill | Dry-run | `just --dry-run` | 生成的 justfile 语法验证 | N/A |

### Key Test Scenarios

1. 单 surface 项目（web）：生成 dev/probe/test/teardown 配方，`[linux]`/`[windows]` 双变体
2. 混合项目（web+api）：生成带前缀的独立配方 + 聚合配方
3. 无 surface 配置：输出与当前一致（零回归）
4. CLI/TUI surface：不生成 run 配方
5. `# user-customized` 保护：已标记的配方不被覆盖
6. `forge surfaces --json`：路径匹配返回正确 JSON，无匹配返回空数组

### Overall Coverage Target

Go 代码 80%

## Security Considerations

### Threat Model

无特殊安全风险。本特性不涉及用户输入验证、认证、授权。

### Mitigations

Surface-key 命名约束 `[a-zA-Z0-9_-]` 由 init-justfile 配方生成时强制执行，防止 just 配方名注入。

## PRD Coverage Map

| PRD AC | Design Component | Interface/Model |
|--------|-----------------|-----------------|
| Story1: web 项目生成 dev/probe/test/teardown | init-justfile SKILL.md + `rules/surfaces/web.md` | Surface 规则文件格式 Interface 2 |
| Story1: CLI/TUI 不生成 run | init-justfile SKILL.md | Surface 规则文件 (cli.md/tui.md 不含 run 契约) |
| Story1: 双平台变体 | init-justfile SKILL.md | 配方生成指导含 `[linux]`/`[windows]` |
| Story2: 调度器模式编排 | run-tests SKILL.md | Surface 规则文件 Interface 2 |
| Story2: probe 失败区分 exit code | run-tests SKILL.md | HARD-GATE 规则 + Exit Code 表 |
| Story2: surface 信息不可用时错误退出 | `forge surfaces --json` | Interface 1 + Error Handling 表 |
| Story3: surface-key 迁移为用户自定义 | `Task.SurfaceKey` + `AutoGenTaskDef.SurfaceKey` | Data Models |
| Story3: resolveScope 重写 | **删除 resolveScope**，直接读 SurfaceKey | `renderTemplate()` 直接注入 |
| Story4: Task 新增双字段 | `Task.SurfaceKey` + `Task.SurfaceType` | Data Models |
| Story4: forge task add 继承 | `AddTaskOpts` 新增 SurfaceKey/SurfaceType | `add.go` 从源任务复制 |
| Story4: fix-task 推断 | `addFixTask()` 调 `forge surfaces <path> --json` | Interface 1 |
| Story4: 旧任务 scope 兼容 | **不保留兼容层**，PRD AC 将更新（见 Migration Notes） | Migration Notes |

## Open Questions

无（PRD 评估 3 轮已解决所有疑问）。

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| 保留 Scope 兼容层 | 旧任务无需更新 | 双字段维护负担，查询逻辑复杂 | 用户选择干净迁移，不做向后兼容 |
| Go 代码管理进程生命周期 | 确定性最高 | 开发成本高，与 Forge 设计哲学冲突 | 长期方向，v3.0.0 范围外 |
| 从 config.yaml 读取 surface 信息（不通过 CLI） | 无 CLI 依赖 | LLM agent 无法直接读取 Go 函数，必须通过 CLI | Skill 层只能通过 Bash 工具执行 CLI |

### References

- PRD Spec: `docs/features/surface-aware-justfile/prd/prd-spec.md`
- Proposal: `docs/proposals/surface-aware-justfile/proposal.md`
- Forge Distribution Convention: `docs/conventions/forge-distribution.md`
- Error Reporting Rules: `docs/business-rules/error-reporting.md`
- Quality Gate Rules: `docs/business-rules/quality-gate.md`
