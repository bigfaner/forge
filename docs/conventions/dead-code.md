---
title: "Dead Code Identification and Cleanup"
domains: [dead-code, deprecation, cleanup, maintenance]
---

# Dead Code Identification and Cleanup

本文档定义 forge-cli 的死代码识别标准、分类策略和清理流程。这是**规范性（normative）**文档，描述目标状态而非当前状态。所有新增代码必须避免引入死代码；已有代码的偏差记录在偏差分析表中，按计划逐步收敛。

## 1. 死代码分类

死代码并非单一概念。根据风险等级和处理方式，划分为三个类别：

### 1.1 纯粹死代码（Category A：可直接删除）

**定义**：无任何生产代码调用、无测试引用、无外部接口依赖的代码。删除后不影响编译或运行。

**识别标准**：
- 函数/方法：在 Go 编译范围内零调用（排除定义处和测试文件）
- 变量/常量：无读取引用
- 类型/结构体：无实例化或嵌入
- 导入包：未被使用的 import

**处理方式**：确认无引用后直接删除，常规 PR 即可。

### 1.2 Test-Bridge 别名（Category B：需评估后处理）

**定义**：在 `internal/cmd/` 中定义的包级变量别名，指向 `pkg/` 中的实际实现。最初是为了让测试可以通过 `internal/cmd` 包内路径 mock 这些函数。

**识别标准**：
- 声明形式为 `var funcName = pkg.SomeFunc`
- 注释中包含 "delegates to" 或 "alias for" 等关键字
- 测试文件通过包内名称访问（而非直接调用 pkg 函数）

**处理方式**：
1. 调查别名在生产代码中的实际调用点数量
2. 如有生产调用：先迁移调用点到直接使用 `pkg.` 路径，再删除别名
3. 如仅测试使用：重写测试为直接导入 `pkg/` 包并调用，再删除别名
4. 删除别名后运行全量测试确认无回归

**关键约束**：不可将仍有生产调用的别名视为"纯粹死代码"直接删除。

### 1.3 Deprecated 保留字段（Category C：需迁移计划）

**定义**：标记为 `Deprecated` 但仍需保留的字段或函数，通常用于数据迁移检测或向后兼容。

**识别标准**：
- 代码注释包含 `Deprecated:` 标记
- 字段仍有运行时读取逻辑（如迁移检测、告警生成）
- YAML/JSON 标签仍存在，反序列化仍在解析该字段

**处理方式**：
1. 记录字段的所有消费者（读取点）
2. 制定迁移时间线：何时完成所有旧数据的迁移
3. 迁移完成后，先移除消费者代码，再移除字段本身
4. 字段移除属于 breaking change，需 bump major 或 minor 版本

## 2. 目标状态

| 属性 | 目标 |
|------|------|
| 纯粹死代码 | 零容忍。`go vet ./...` 通过即保证无未使用导入；代码审查保证无未调用函数 |
| Test-Bridge 别名 | 零存在。测试直接导入 `pkg/` 包；生产代码直接使用 `pkg.` 路径 |
| Deprecated 字段 | 最多存在 1 个活跃的 deprecated 字段（迁移窗口内）；迁移完成即移除 |
| 重复定义 | 零容忍。相同签名的函数只定义一次，消费方统一导入 |

## 3. 当前偏差分析

### 模块级偏差摘要

| 编号 | 偏差项 | 类别 | 当前状态 | 目标状态 | 风险等级 |
|:----:|--------|------|----------|----------|:--------:|
| ~~DC-1~~ | `cmd.Debugf` 重复定义 | A（纯粹死代码/重复） | **已修复** — `internal/cmd/output.go` 中的重复定义已删除，仅保留 `base.Debugf` | N/A | N/A |
| ~~DC-2~~ | `getTaskPhase` test-bridge 别名 | B（Test-Bridge） | **已修复** — 别名已删除，5 处调用迁移为 `task.GetTaskPhase()` 直接调用，文件已更名为 `validate.go` | N/A | N/A |
| ~~DC-3~~ | `checkExistingTaskState` test-bridge 别名 | B（Test-Bridge） | **已修复** — 别名已删除，调用迁移为 `task.CheckExistingTaskState()` 直接调用 | N/A | N/A |
| ~~DC-4~~ | `compareVersionIDs` test-bridge 别名 | B（Test-Bridge） | **已修复** — 别名已删除，调用迁移为 `task.CompareVersionIDs()` 直接调用 | N/A | N/A |
| ~~DC-5~~ | `FrontmatterData.Scope` deprecated 字段 | C（Deprecated） | **已修复** — `FrontmatterData.Scope` 已从 `frontmatter.go` 移除，`CheckLegacyScope` 保留用于迁移检测 | N/A | N/A |

### 偏差详细说明

#### ~~DC-2~~ 已修复：`getTaskPhase` test-bridge 别名

`getTaskPhase` 别名已删除。5 处生产调用已全部迁移为 `task.GetTaskPhase()` 直接调用（位于 `internal/cmd/task/validate.go`）。文件已从 `validate_index.go` 更名为 `validate.go`。

#### DC-5 迁移路径

`Scope` 字段的完整消费者链：
- `CheckLegacyScope`（`pkg/task/migrate.go`）：扫描所有任务，检测仍使用 `Scope` 但无 `SurfaceKey` 的任务
- `migrateScopeToSurface`（`internal/cmd/task/migrate.go`）：执行 scope → surface 的数据迁移
- `FrontmatterData.Scope`（`pkg/task/frontmatter.go`）：YAML 反序列化目标字段
- `Task.Scope`（`pkg/task/types.go`）：内存数据模型字段

迁移完成条件：所有现有任务的 YAML frontmatter 中 `scope` 字段已被清空或替换为 `surface-key`/`surface-type`。

## 4. 清理流程

### 4.1 识别阶段

```
1. 运行静态分析：
   go vet ./...                    # 编译器级未使用检测
   deadcode -test ./...            # golang.org/x/tools/deadcode（含测试引用）

2. 人工审查补充：
   - 搜索 "Deprecated:" 注释
   - 搜索 "delegates to" / "alias for" / "Kept as alias" 注释
   - 对比 internal/cmd/ 与 pkg/ 中的同名函数

3. 分类（按第 1 节标准）：
   - Category A → 直接进入清理
   - Category B → 评估生产调用点
   - Category C → 制定迁移计划
```

### 4.2 清理执行

#### Category A（纯粹死代码）

```
1. 确认无引用：grep -rn "funcName\|TypeName\|varName" --include="*.go"
2. 删除代码
3. 运行 go build ./... && go test ./... 确认无回归
```

#### Category B（Test-Bridge 别名）

```
1. 统计生产调用点：grep -rn "aliasName" --include="*.go" | grep -v "_test.go" | grep -v "var aliasName"
2. 如有生产调用：
   a. 逐个迁移调用为 pkg.SomeFunc 直接调用
   b. 更新对应的测试代码
3. 如仅测试使用：
   a. 重写测试，import pkg 包并直接调用
4. 删除别名声明
5. 运行全量测试：go test -race -cover ./...
```

#### Category C（Deprecated 保留字段）

```
1. 列出所有消费者（参见偏差详细说明中的消费者链）
2. 确认迁移时间线（通常不超过 2 个 sprint）
3. 执行迁移：
   a. 运行迁移命令将旧数据转换为新格式
   b. 验证迁移结果
4. 移除消费者代码
5. 移除 deprecated 字段
6. 版本 bump（至少 minor）
```

### 4.3 防止新增死代码

| 检查点 | 方式 |
|--------|------|
| PR 审查 | 确认新增函数有调用方；确认不引入重复定义 |
| CI 门禁 | `go vet ./...` 作为构建步骤 |
| 代码注释 | test-bridge 别名必须注释 "delegates to" + 目标函数路径 |
| deprecated 标记 | 使用标准 `Deprecated:` 前缀，说明替代方案和计划移除时间 |

## 5. 参考

- [package-organization.md](./package-organization.md) -- 包组织规范，定义 `cmd → internal → pkg` 依赖方向
- [code-structure.md](./code-structure.md) -- 代码结构规范（嵌套、缩进、控制流）
- [pkg-dependency-graph.md](../../features/forge-cli-codebase-standards/pkg-dependency-graph.md) -- pkg/ 依赖图事实基线
