---
created: 2026-05-15
author: "fanhuifeng"
status: Draft
---

# Proposal: Go Test Flat Staging -- 消除 go-test 的曲折毕业路径

## Problem

go-test profile 的测试脚本从 staging 到毕业需要经历不必要的曲折。核心矛盾：

**Go 的 package 是按目录隔离的编译单元，而 forge 的两阶段 staging/graduation 流程是为 TypeScript 的路径 import 设计的。**

### 具体表现

#### 1. 生成阶段：staging 目录导致 helpers 不可见

`gen-test-scripts` 将测试文件写到 `tests/e2e/features/<slug>/` 子目录。但 Go 在该子目录编译时，只看该目录下的文件，无法引用父目录 `tests/e2e/helpers.go` 中的函数。

```
tests/e2e/
  helpers.go           # package e2e, func runCLI()
  features/
    my-feature/
      foo_cli_test.go  # package e2e, 调用 runCLI() → 编译失败: undefined
```

**当前变通方案**：每个 feature 子目录复制一份 helpers（`tsptRunCLI`、`init()` 等），如 `tui-ui-design/` 和刚生成的 `test-scripts-per-type/` 都是如此。

#### 2. 毕业阶段：import rewrite 策略声明与实际不符

go-test 的 `graduate.md` 声明 "Import Rewrite: None required"，理由是 Go 使用 module path。但实际毕业后：

| 已毕业 feature | 实际做法 | 与策略的矛盾 |
|----------------|----------|-------------|
| `justfile-canonical-e2e/` | 独立 `go.mod`，独立 `package justfile_canonical_e2e`，独立 `helpers_test.go` | 非共享 helpers，需要复制 |
| `tui-ui-design/` | 复制到子目录，自包含 helpers | 非共享 helpers |
| `cli_lean_output_cli_test.go` | 直接放在 `tests/e2e/` 根目录，共享 helpers | 没有经过 features/ staging |

没有任何一个已毕业的 go-test feature 实际经历了 "features/ staging → 搬移 + import rewrite → 毕业" 这个流程。每个都在生成阶段就绕过了 staging。

#### 3. 每个 feature 重复造轮子

| Helper 函数 | tests/e2e/helpers.go | tui-ui-design/ | justfile-canonical-e2e/ | test-scripts-per-type/ |
|-------------|---------------------|----------------|------------------------|----------------------|
| `runCLI` / `tsptRunCLI` / `runForge` | 有 | 无（直接读文件） | `runForge`（不同实现） | `tsptRunCLI`（复制版） |
| `runCLIRaw` / `tsptRunCLIRaw` | 有 | 无 | 无 | `tsptRunCLIRaw`（复制版） |
| `withRetry` | 有 | 无 | `withRetry`（复制版） | 无 |
| `init()` chdir | 无 | 有 | 无 | 有（复制版） |
| `parseBlock` / `hasField` | 有 | 无 | 无 | 无 |

### Evidence

- 已毕业的 6 个 go-test feature 中，0 个走完 staging → graduation 的完整流程
- `gen-test-scripts` SKILL.md 的 HARD-GATE 硬编码 `tests/e2e/features/<feature>/` 为唯一输出路径
- `graduate-tests` 的 "Import Rewrite: None required" 暗示 graduation 是 no-op -- 既然 nothing to rewrite，为什么还需要 staging？

### Urgency

当前项目已有 6+ 个 go-test feature。每个新 feature 都需要复制 helpers、加 `init()` chdir hack。随着 feature 增多，helper 副本间的 diverge 风险增大。

## Proposed Solution

**核心思路**：go-test profile 生成测试文件直接放到 `tests/e2e/` 根目录，跳过 `features/` staging。由于文件已在最终位置，取消 graduation 任务（T-test-4 / T-quick-4），由 run-tests 通过后自动写 marker。

### 为什么这对 Go 是正确的

1. **Go package = directory**。同一个 `package e2e` 下的所有文件共享符号表。放在根目录天然共享 helpers。
2. **Go 没有 import path rewriting 问题**。TypeScript 需要 staging 是因为 import `'../../helpers.js'` 是相对路径，搬文件要改路径。Go 用 module path，不存在这个问题。
3. **Go test discovery 是递归的**。`go test ./...` 自动发现所有子目录的测试。staging 目录不阻碍也不帮助测试发现。
4. **功能模块分类在 Go 中无意义**。TypeScript 按业务域分目录（`tests/e2e/module/`）有意义，因为 Playwright 可以按目录过滤。Go 按目录分反而破坏了共享 package 的优势。

### 方案细节

#### 1. Profile manifest 增加 staging-mode 字段

在 `manifest.yaml` 中声明 staging 策略：

```yaml
# go-test/manifest.yaml
name: go-test
display: "Go CLI/TUI/Backend (go test)"
language: go
file-extension: _test.go
test-directory: tests/e2e/

staging-mode: flat   # 新字段。flat = 直接写到 test-directory 根目录，不生成 graduation 任务

capabilities: [tui, api, cli]

templates:
  test-file: templates/test-file.go
  helpers: templates/helpers.go
  additional: []
```

对比 web-playwright 保持现有行为：

```yaml
# web-playwright/manifest.yaml
staging-mode: nested   # 默认值，写到 features/<slug>/，生成 graduation 任务
```

**staging-mode 值**：

| 值 | 生成路径 | graduation 任务 | 适用场景 |
|----|---------|----------------|---------|
| `nested` (默认) | `tests/e2e/features/<slug>/` | 生成 T-test-4 / T-quick-4 | 需要 import rewrite 的语言 (TypeScript) |
| `flat` | `tests/e2e/` 根目录 | **不生成** | 使用 module path 或不需要 rewrite 的语言 (Go, Python, Rust) |

#### 2. gen-test-scripts 读取 staging-mode

`gen-test-scripts` SKILL.md 的 HARD-GATE 从硬编码 `features/<feature>/` 改为根据 profile 的 `staging-mode` 决定输出路径：

```
staging-mode=nested → tests/e2e/features/<slug>/     (现有行为)
staging-mode=flat   → tests/e2e/<slug>_cli_test.go   (直接放根目录)
```

flat 模式下的文件命名：<feature>_<type>_test.go，和现有惯例一致（如 `cli_lean_output_cli_test.go`）。

#### 3. run-e2e-tests 适配 flat 模式

`run-e2e-tests` 当前 prerequisite 硬编码 `tests/e2e/features/<slug>/`。增加 flat 模式检测：

```
staging-mode=flat → 检查 tests/e2e/<slug>_*_test.go 是否存在
```

flat 模式下，run-e2e-tests 通过后自动写 graduation marker（无需单独的 graduation 任务）。

#### 4. 取消 flat 模式的 graduation 任务

**核心决策**：flat 模式下不生成 T-test-4 / T-quick-4 任务。

**理由**：flat 模式的测试文件已经在最终位置（`tests/e2e/`），run-e2e-tests 通过即意味着脚本可用。Graduation 的两个原始目的 -- 文件搬移 + import rewrite -- 在 flat 模式下都不存在。生成一个不做任何事的任务只会增加 pipeline 开销和 review 负担。

**实现**：修改 `testgen.go` 中的 `GetBreakdownTestTasks` 和 `GetQuickTestTasks`，读取 profile 的 `staging-mode`，当为 `flat` 时跳过 graduate 任务生成，并调整依赖链。

改动前（go-test, breakdown mode）:
```
T-test-1 → T-test-1b → T-test-2-cli → T-test-3 → T-test-4 → T-test-4.5 → T-test-5
```

改动后（go-test with flat staging）:
```
T-test-1 → T-test-1b → T-test-2-cli → T-test-3 → T-test-4.5 → T-test-5
```

T-test-4 被移除。T-test-4.5 (verify-regression) 的依赖从 T-test-4 改为 T-test-3。

改动前（go-test, quick mode）:
```
T-quick-1 → T-quick-2-cli → T-quick-3 → T-quick-4 → T-quick-5
```

改动后（go-test with flat staging）:
```
T-quick-1 → T-quick-2-cli → T-quick-3 → T-quick-5
```

T-quick-4 被移除。T-quick-5 (verify-regression) 的依赖从 T-quick-4 改为 T-quick-3。

#### 5. run-e2e-tests 集成 graduation marker 写入

flat 模式下，run-e2e-tests 通过后自动写 graduation marker：

```yaml
# tests/e2e/.graduated/<slug> -- 由 run-e2e-tests 在 flat 模式下自动生成
schema_version: 1
status: completed
timestamp: <UTC ISO timestamp>
source: tests/e2e/<slug>_*_test.go
mode: flat
testCount: <N>
```

这保证了 `.graduated/` marker 的一致性，下游逻辑（如 docsync）不需要区分 flat 和 nested 模式的 marker。

#### 6. graduate-tests SKILL.md 增加 flat 模式跳过逻辑

当 agent 直接调用 `/graduate-tests` 时，检测 profile 的 staging-mode：
- `nested` → 执行现有 graduation 流程
- `flat` → 输出 "Flat staging mode: graduation not needed. Tests are already at final location. Marker written by run-e2e-tests."

### 其他受影响的 profile

| Profile | staging-mode 建议 | graduation 任务 | 理由 |
|---------|-------------------|----------------|------|
| go-test | `flat` | 不生成 | Module path，共享 package，无需搬移 |
| pytest | `flat` | 不生成 | Python import 用 module path，不依赖相对路径 |
| rust-test | `flat` | 不生成 | Rust use path 是 module path |
| java-junit | `nested` (保持) | 生成 | Java package 与目录强绑定，毕业时可能需要调整 package 声明 |
| web-playwright | `nested` (保持) | 生成 | TypeScript 相对 import 需要 rewrite |
| maestro | `nested` (保持) | 生成 | YAML 配置文件，可能引用相对路径 |

## Requirements Analysis

### Key Scenarios

1. **新 feature 生成**：`gen-test-scripts` 生成 `tests/e2e/my_feature_cli_test.go`，直接使用根目录的 `helpers.go`，无需复制 helpers
2. **运行测试**：`run-e2e-tests` 在 `tests/e2e/` 运行 `go test ./... -tags=e2e`，自动发现新文件
3. **测试通过后**：run-e2e-tests 自动写 graduation marker，无需单独的 graduation 任务
4. **pipeline 缩短**：go-test 从 6 步变为 5 步，quick mode 从 5 步变为 4 步

### Non-Functional Requirements

- **向后兼容**：已毕业的 feature 不受影响。已有的 `features/` 目录下的文件保持不变
- **减少代码重复**：所有 go-test feature 共享 `tests/e2e/helpers.go`，消除 helper 副本
- **pipeline 加速**：减少一个无意义的中间任务，缩短整体执行时间

### Constraints & Dependencies

- `manifest.yaml` schema 需要支持 `staging-mode` 字段（向后兼容：缺失时默认 `nested`）
- `testgen.go` 需要读取 profile 的 staging-mode 来决定是否生成 graduate 任务
- `profile` package 需要暴露 staging-mode 解析能力
- `run-e2e-tests` SKILL.md 需要增加 flat 模式下自动写 marker 的逻辑

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Flat staging + no graduation task (recommended)** | 零 helper 复制，零 import rewrite，pipeline 最短 | 所有 test 文件在同一目录，需要代码改动 | **Selected: 最彻底的简化** |
| Flat staging + marker-only graduation task | 零 helper 复制，pipeline 保持一致 | 多一个空转任务，增加 pipeline 开销 | Rejected: 任务不做任何事就不该存在 |
| Go workspace (go.work) | 保留子目录，共享 helpers 通过 module 引用 | 需改造 go.mod 为 library module，改动大，overkill | Rejected: 为 staging 问题引入 workspace 是杀鸡用牛刀 |
| 保持 nested + 改善 helper 复制 | 最小改动 | 本质问题未解决，helper 副本持续 diverge | Rejected: 治标不治本 |
| 保持 nested + symlink helpers | 子目录通过 symlink 引用根目录 helpers | 平台兼容性 (Windows symlink 需要 admin 权限) | Rejected: Windows 兼容性差 |

### 单目录文件数量质疑

"所有 test 文件在同一目录会不会太乱？" -- 不会。当前 `tests/e2e/` 根目录只有 2 个 `_test.go` 文件。按每个 feature 1 个文件计算，20 个 feature = 20 个文件。Go 标准库和大型项目中，单个目录包含 20-50 个文件是常见且可维护的。文件名 `<feature>_<type>_test.go` 自带分类信息。

## Feasibility Assessment

### Technical Feasibility

改动分为文档改动和代码改动两部分：

**文档改动**：

| 文件 | 改动 |
|------|------|
| `manifest.yaml` | 加 `staging-mode: flat` 字段 |
| `gen-test-scripts/SKILL.md` | HARD-GATE 读取 staging-mode，flat 模式输出到根目录 |
| `run-e2e-tests/SKILL.md` | prerequisite 适配 flat 模式 + 通过后自动写 marker |
| `graduate-tests/SKILL.md` | 增加 flat 模式跳过逻辑 |

**代码改动**：

| 文件 | 改动 |
|------|------|
| `profile/embed.go` | 解析 `manifest.yaml` 的 `staging-mode` 字段，暴露 `GetStagingMode(name) string` |
| `task/testgen.go` | `GetBreakdownTestTasks` / `GetQuickTestTasks` 读取 staging-mode，flat 模式跳过 graduate 任务，调整依赖链 |
| `task/testgen_test.go` | 新增 flat 模式测试用例 |

### Resource & Timeline

2-3 个任务：文档改动 + 代码改动 + 测试。估计 2-3 小时。

### Dependency Readiness

`manifest.yaml` 的 `staging-mode` 字段已可被 `profile` package 解析（YAML unmarshal 自动忽略未知字段，需要显式添加解析逻辑）。

## Scope

### In Scope

- `manifest.yaml` 加 `staging-mode` 字段（go-test、pytest、rust-test）
- `gen-test-scripts/SKILL.md` 适配 flat staging
- `run-e2e-tests/SKILL.md` 适配 flat prerequisite + 自动写 marker
- `graduate-tests/SKILL.md` 增加 flat 模式跳过逻辑
- `profile/embed.go` 增加 staging-mode 解析
- `task/testgen.go` flat 模式跳过 graduate 任务 + 调整依赖链

### Out of Scope

- 已毕业 feature 的迁移（保持现状）
- java-junit、maestro、web-playwright 的 staging-mode 变更
- `features/` 目录下已有文件的清理
- `infer.go` 中的 `TypeTestPipelineGraduate` 类型保留（向后兼容已存在的任务文件）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 根目录文件过多 | L | L | 文件名自带 feature 前缀，实际 20+ feature 才会有 20 个文件 |
| 与 nested 模式 feature 混合 | L | M | flat 和 nested 可以共存，互不影响 |
| 已有 helper 不够用需要扩展 | M | L | helpers.go 可通过 gen-test-scripts 的 merge 机制扩展 |
| testgen.go 依赖 profile package | M | M | testgen 需要 profile staging-mode 数据；传入方式可以是函数参数或全局配置 |
| 向后兼容性 | L | M | `staging-mode` 缺失时默认 `nested`，现有行为不变 |

## Success Criteria

- [ ] go-test profile 生成测试文件直接到 `tests/e2e/` 根目录，无需复制 helpers
- [ ] 生成的测试文件编译通过，可直接使用 `helpers.go` 中的 `runCLI`、`parseBlock` 等函数
- [ ] flat 模式下不生成 T-test-4 / T-quick-4 任务
- [ ] run-e2e-tests 通过后自动写 graduation marker
- [ ] 已毕业的 feature 不受影响
- [ ] T-test-4.5 / T-quick-5 的依赖链正确指向 T-test-3 / T-quick-3
- [ ] pytest、rust-test profile 可选择 opt-in flat staging

## Next Steps

- 通过 `/quick` 流程生成实施任务
- 或者直接实施（预计 2-3 个任务）
