---
feature: "unify-enum-constants"
---

# User Stories: unify-enum-constants

## Story 1: Type-Safe Status Constants

**As a** Forge CLI 开发者
**I want to** 使用 `types.Status` 类型常量而非字符串字面量来表示任务状态
**So that** 拼写错误或无效状态值在编译期被捕获，而非运行时静默失败

**Acceptance Criteria:**
- Given `pkg/types/status.go` 定义了 `type Status string` 和 7 个常量
- When 开发者在任意 .go 文件中赋值 `task.Status = "complteed"`（拼写错误）
- Then 编译器报错：`cannot use "complteed" (untyped string constant) as types.Status in assignment`
- Given `statemachine.go` 状态转移表使用 `types.StatusXxx` 常量
- When 开发者添加新的状态转移
- Then 编译器验证 `From`/`To` 字段类型为 `types.Status`，拼写错误被拦截

---

## Story 2: Centralized Surface Type Constants

**As a** Forge CLI 开发者
**I want to** 所有 Surface Type（web/api/cli/tui/mobile）定义在 `pkg/types/surface.go` 中的 typed constants
**So that** 添加新 Surface Type 时只需修改一处常量定义，所有引用点自动生效

**Acceptance Criteria:**
- Given `pkg/types/surface.go` 定义了 `type SurfaceType string` 和 5 个常量
- When 新增一种 Surface Type（如 `"desktop"`）
- Then 只需在 `surface.go` 中添加 `SurfaceDesktop SurfaceType = "desktop"` 并更新 `AllSurfaceTypes()`
- Given `detect_surface.go` 中的映射表使用 `types.SurfaceWeb` 等常量
- When 运行 `go build ./...`
- Then 编译通过，所有 Surface Type 引用类型正确

---

## Story 3: Reliable Enum Refactoring

**As a** Forge CLI 维护者
**I want to** 所有枚举值通过 typed constants 引用而非散落的字符串字面量
**So that** 通过 `go build` 即可验证枚举引用完整性，无需手动 grep 检查遗漏

**Acceptance Criteria:**
- Given 生产代码中 Status、SurfaceType、Priority 魔法值数量为 0
- When 运行 `go build ./...`
- Then 零编译错误（所有类型签名一致）
- When 运行 `go test ./...`
- Then 所有测试通过（行为零变更）
- Given `pkg/types/` 不导入任何 forge-cli 内部包
- When 检查依赖方向
- Then `pkg/types/` 是叶包，无循环依赖

---

## Story 4: Validation Map Consolidation

**As a** Forge CLI 开发者
**I want to** `validate_index.go` 中的 `validStatus`/`validPriority` map 使用 `types.AllStatuses()`/`types.AllPriorities()` 生成
**So that** 枚举验证逻辑与常量定义单一来源同步，不会遗漏新增的枚举值

**Acceptance Criteria:**
- Given `types.AllStatuses()` 返回所有 7 个 Status 常量
- When `validate_index.go` 构建验证 map
- Then 使用 `types.AllStatuses()` 而非硬编码的 `map[string]bool{"pending": true, ...}`
- Given 未来新增一个 Status 常量
- When 更新 `AllStatuses()` 实现
- Then `validate_index.go` 的验证逻辑自动包含新值，无需手动同步
