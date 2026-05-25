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
| forge surfaces CLI       |  (已存在，新增 --json)
|  -- ReadSurfaces()       |
|  -- MatchSurface(path)   |
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

- 内部：`forgeconfig.ReadSurfaces()` + `MatchSurface()`（已存在）
- 内部：`task` 包所有相关结构体
- 外部：just >= 1.4.0（`[linux]`/`[windows]` recipe attribute）
- 无新外部依赖

## Interfaces

### Interface 1: `forge surfaces <path> --json`

```
// CLI: 查询文件路径所属 surface，输出 JSON
// 输入: 文件路径（相对于项目根目录）
// 输出: JSON array
[
  {"surface-key": "admin-panel", "surface-type": "web"}
]
// 无匹配时: stdout []，exit 0
// surfaces 配置缺失时: stderr 错误消息 + 恢复提示，exit 1
```

### Interface 2: Surface 规则文件格式

```
// 文件: rules/surfaces/<type>.md
// 消费者: init-justfile (配方生成) + run-tests (编排序列)

## 编排序列
// 步骤名 | 退出码 | 语义 | 后续动作

## 配方调用契约
// 配方名 | 参数签名 | 退出码语义

## journey 过滤策略
// journey 标签映射表
```

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
| `forge surfaces --json` surfaces 配置缺失 | stderr 错误 + 恢复提示 | 1 (retryable) |
| `forge surfaces --json` 路径无匹配 | stdout `[]` + exit 0 | 0 |
| `build.go` 解析任务 frontmatter 缺 surface-key | 静默赋空值 | 0 |

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
