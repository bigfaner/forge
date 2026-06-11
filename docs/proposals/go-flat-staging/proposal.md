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

当前项目已有 6 个已毕业 go-test feature，每个在生成阶段就绕过了 staging（见 Evidence）。每新增一个 go-test feature 需要约 30 分钟手动工作（测量依据：最近 3 个 go-test feature 的 helper 同步工作——`tui-ui-design/`、`justfile-canonical-e2e/`、`test-scripts-per-type/`——从 `gen-test-scripts` 生成到编译通过，helper 复制 + package 调整 + chdir hack 分别耗时约 25、32、28 分钟）。按项目规划（见 `docs/roadmap/` 中的下季度 feature list）还有 8-10 个 go-test feature，累计约 4-5 小时纯重复劳动。

已观察到的 diverge 实例：`withRetry` 函数在 `tests/e2e/helpers.go` 和 `justfile-canonical-e2e/` 中存在两个不同实现（参数签名不同）；`tsptRunCLI` 在 `test-scripts-per-type/` 的复制版与 root 版本已差异数个 commit。如果不现在解决，每增加一个 feature 就多一个 helper 副本需要手动同步，bug 修复需要同步到 N 个位置而非一处。

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

**实现**：修改 `testgen.go` 中的 `GetBreakdownTestTasks` 和 `GetQuickTestTasks`，通过显式函数参数接收 staging-mode（`stagingMode string`），当为 `flat` 时跳过 graduate 任务生成，并调整依赖链。调用方通过 `profile.GetStagingMode(profileName)` 从 manifest 读取值后传入。

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

### Developer Walkthrough: Before vs After

#### Before (当前行为)

```
Developer: /gen-test-scripts --feature justfile-canonical-e2e
Output:    Writing test to tests/e2e/features/justfile-canonical-e2e/cli_test.go
           ERROR: compilation failed -- undefined: runCLI

Developer: (手动复制 helpers.go 到 features/justfile-canonical-e2e/)
           (手动添加 init() chdir hack)

Developer: /run-e2e-tests --feature justfile-canonical-e2e
Output:    Running go test ./tests/e2e/features/justfile-canonical-e2e/...
           PASS

Developer: /graduate-tests --feature justfile-canonical-e2e
Output:    Graduating from features/justfile-canonical-e2e/ to tests/e2e/
           Import rewrite: None required
           (实际是手动搬移 + 重写 package)

Pipeline: T-test-1 → T-test-1b → T-test-2-cli → T-test-3 → T-test-4 → T-test-4.5 → T-test-5
```

#### After (flat staging)

```
Developer: /gen-test-scripts --feature justfile-canonical-e2e
Output:    [flat staging mode detected for go-test profile]
           Writing test to tests/e2e/justfile_canonical_e2e_cli_test.go
           Compiling... OK (helpers.go shared via package e2e)

Developer: /run-e2e-tests --feature justfile-canonical-e2e
Output:    [flat mode] Checking tests/e2e/justfile_canonical_e2e_*_test.go
           Running go test ./... -run TestJustfileCanonicalE2e -tags=e2e
           PASS
           Writing graduation marker to tests/e2e/.graduated/justfile-canonical-e2e

Developer: /graduate-tests --feature justfile-canonical-e2e
Output:    Flat staging mode: graduation not needed.
           Tests are already at final location (tests/e2e/).
           Marker written by run-e2e-tests.

Pipeline: T-test-1 → T-test-1b → T-test-2-cli → T-test-3 → T-test-4.5 → T-test-5
```

**After - Error Scenario (helpers.go merge conflict)**:

```
Developer: /gen-test-scripts --feature another-feature
Output:    [flat staging mode detected for go-test profile]
           Writing test to tests/e2e/another_feature_cli_test.go
           Merging helper parseBlock into tests/e2e/helpers.go...
           ERROR: merge conflict detected in tests/e2e/helpers.go
           Conflict: helper parseBlock already modified by concurrent feature
           Diff:
             - existing:  func parseBlock(output string) map[string]string
             + incoming:  func parseBlock(output string, mode string) map[string]string
           Action required: resolve conflict manually, then re-run gen-test-scripts.

Developer: (手动 resolve helpers.go 中的 conflict)

Developer: /gen-test-scripts --feature another-feature
Output:    [flat staging mode detected for go-test profile]
           Writing test to tests/e2e/another_feature_cli_test.go
           Compiling... OK (helpers.go shared via package e2e)
```

**关键体验差异**：
- 消除手动复制 helpers 的步骤
- 消除 init() chdir hack
- 消除 graduation 手动搬移步骤
- `/graduate-tests` 明确告知 "不需要" 而非执行空操作
- pipeline 从 6 步减少到 5 步（breakdown mode）或 5 步减到 4 步（quick mode）

### 其他受影响的 profile

| Profile | staging-mode 建议 | graduation 任务 | 理由 |
|---------|-------------------|----------------|------|
| go-test | `flat` | 不生成 | Module path，共享 package，无需搬移 |
| pytest | `flat` (声明性占位，本次不实现运行时行为) | 不生成 | Python import 用 module path，不依赖相对路径。`gen-test-scripts` 不为 pytest/rust-test 处理 `staging-mode`（Out of Scope） |
| rust-test | `flat` (声明性占位，本次不实现运行时行为) | 不生成 | Rust use path 是 module path。`gen-test-scripts` 不为 pytest/rust-test 处理 `staging-mode`（Out of Scope） |
| java-junit | `nested` (保持) | 生成 | Java package 与目录强绑定，毕业时可能需要调整 package 声明 |
| web-playwright | `nested` (保持) | 生成 | TypeScript 相对 import 需要 rewrite |
| maestro | `nested` (保持) | 生成 | YAML 配置文件，可能引用相对路径 |

## Requirements Analysis

### Key Scenarios

#### Happy Path

1. **新 feature 生成**：`gen-test-scripts` 生成 `tests/e2e/my_feature_cli_test.go`，直接使用根目录的 `helpers.go`，无需复制 helpers
2. **运行测试**：`run-e2e-tests` 在 `tests/e2e/` 运行 `go test ./... -tags=e2e`，自动发现新文件
3. **测试通过后**：run-e2e-tests 自动写 graduation marker，无需单独的 graduation 任务
4. **pipeline 缩短**：go-test 从 6 步变为 5 步，quick mode 从 5 步变为 4 步

#### Error & Edge Cases

5. **测试失败后重新生成**：开发者在 flat 模式下运行 `gen-test-scripts` 对同一 feature 重新生成。系统覆盖 `tests/e2e/<slug>_cli_test.go`（同路径同文件名），不会产生重复文件。marker 未写入（因为测试未通过），所以不干扰已有记录。

6. **并发 helpers.go merge conflict**：两个 feature 同时开发，各自的 `gen-test-scripts` 尝试向 `helpers.go` 添加新 helper 函数。由于 `helpers.go` 是模板生成，`gen-test-scripts` 的 merge 机制应检测 conflict 并报错，要求开发者手动 resolve。**处理方式**：merge conflict 时 `gen-test-scripts` 中止并输出 conflict diff，由开发者手动合并后重新运行。

7. **feature 需要 root helpers.go 中不存在的 helper**：新 feature 需要 `runCLIRaw` 但 root `helpers.go` 中只有 `runCLI`。`gen-test-scripts` 检测到模板引用的 helper 不在 root `helpers.go` 中，将新 helper append 到 `helpers.go`（通过 merge 机制）。如果 helper 已被其他 feature 添加过，跳过追加。

8. **gen-test-scripts 对同一 feature 被调用两次**：第二次调用覆盖第一次生成的文件（同路径）。如果第一次调用已通过测试并写入 marker，marker 保留不变（marker 只在 run-e2e-tests 通过后更新）。不会产生两份测试文件或两个 marker。

9. **已存在的 features/ 目录残留**：切换到 flat 模式后，旧的 `features/<slug>/` 目录中的测试文件仍在。Go 的 `go test ./...` 会递归编译子目录，导致 old nested 和 new flat 文件同时被编译。由于两者属于不同目录（不同 `package` 声明或同一 `package e2e`），**如果函数名相同**（如 `TestMyFeatureCLI`），Go 不会报编译错误（不同 package 的同名函数不冲突），但 `go test` 会执行两遍测试——old nested 中的测试可能引用已 diverge 的 helper 副本，产生不一致结果。**建议**：在 scope 之外保留这些旧文件，但 marker 系统只跟踪 flat 模式生成的文件。旧文件需手动清理（Out of Scope 已声明）。

### Non-Functional Requirements

- **向后兼容**：已毕业的 feature 不受影响。已有的 `features/` 目录下的文件保持不变
- **减少代码重复**：所有 go-test feature 共享 `tests/e2e/helpers.go`，消除 helper 副本
- **pipeline 加速**：减少一个无意义的中间任务，缩短整体执行时间
- **测试隔离**：共享 `package e2e` 的各 feature 测试文件不得通过全局变量或 `init()` 污染彼此状态。要求：每个测试函数使用 `t.Parallel()` 独立运行，禁止在测试中修改 package-level 变量。**迁移策略**：现有 features（`tui-ui-design/`、`test-scripts-per-type/`）中的 `init()` chdir hack 在 flat 模式下不再需要——flat 模式的测试文件直接在 `tests/e2e/` 根目录，工作目录已正确。新 flat 模式生成的文件不包含 `init()`。旧的 nested 文件（`features/<slug>/`）保留在 Out of Scope，不迁移到 flat 目录，因此其 `init()` 不会污染 flat 模式测试。验证方式：在 `tests/e2e/` 下放置两个不同 feature 的 flat 模式测试文件，`go test ./... -count=1` 通过且无竞态检测报告
- **性能不退化**：`go test ./... -count=1` 在 `tests/e2e/` 根目录的执行时间不超过 flat 模式切换前的基线值（以当前 6 个已毕业 feature 为基准测量，20 个 feature 时 `go test` discovery + compile 时间应 < 5 秒）
- **可维护性**：flat 模式不适用于所有语言。当 profile 的测试文件依赖相对路径 import（TypeScript、YAML 配置文件引用）时，必须使用 `nested` 模式。`staging-mode` 字段的可选值由文档约束，不得随意扩展

### Constraints & Dependencies

- `manifest.yaml` schema 需要支持 `staging-mode` 字段（向后兼容：缺失时默认 `nested`）
- `testgen.go` 需要读取 profile 的 staging-mode 来决定是否生成 graduate 任务
- `profile` package 需要暴露 staging-mode 解析能力
- `run-e2e-tests` SKILL.md 需要增加 flat 模式下自动写 marker 的逻辑

## Alternatives & Industry Benchmarking

### Industry Patterns for Go Test Organization

Go 社区和大型项目在测试组织上有成熟的做法，均指向同一个结论：同一 package 的测试文件应放在同一目录。

**1. Go 标准库 (stdlib) 的 flat layout**

Go 标准库的 `src/strings/` 目录包含 `strings.go`（实现）和 `strings_test.go`、`reader_test.go`、`replace_test.go` 等多个测试文件，全部在同一目录，共享同一 `package strings`。没有按测试类型或功能划分子目录。这是 Go 官方推荐的 layout（参考 `golang.org/doc/go1.html` 中的 test organization 描述）。

Source: https://cs.opensource.google/go/go/+/refs/tags/go1.24.3:src/strings/strings_test.go

**2. Kubernetes 的 test helper 模式**

Kubernetes (`kubernetes/kubernetes`) 在 `test/e2e/` 目录下放置所有 e2e 测试文件（超过 200 个 `_test.go` 文件），共享 `test/e2e/framework/` 中的 helper 函数。测试文件通过 `framework` package 引用 helpers，而非复制。这种 "framework package + flat test files" 模式在大型 Go 项目中是标准做法。

Source: https://github.com/kubernetes/kubernetes/tree/master/test/e2e

**3. Bazel 的 language-specific test strategy**

Bazel 通过 `rules_go` 对 Go 测试提供 `go_test` rule，默认在同一 package 下编译测试，无需 staging 目录。Bazel 的 approach 是：每个语言有自己的编译和测试约定，build system 应适配语言，而非让语言适配 build system。这与本 proposal 的核心思路一致。

Source: https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/rules.md#go_test

**4. Hugo 的 single-directory e2e test layout**

Hugo (`gohugoio/hugo`) 的 integration test 全部放在 `hugolib/` 目录下，一个 `package hugolib` 包含数百个 `_test.go` 文件。Helper 函数（如 `newSiteBuilder`）定义在同一个 package 中，所有测试文件直接调用，无复制。

Source: https://github.com/gohugoio/hugo/tree/master/hugolib

**Relevance note**: 以上项目均无 staging/graduation 概念。引用它们是为了验证 Go 社区公认的测试文件组织方式（同 package 同目录），从而论证 flat layout 的合理性。staging/graduation 是 forge 的特有流程，其简化需要从 Go 的编译模型推导，而非从外部项目照搬。

### Alternatives Comparison

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Do nothing (status quo)** | 零改动成本 | 每个 feature 复制 helpers，已证明 diverge（6 个 feature，3 种不同 helper 实现） | Rejected: 问题已量化，随 feature 增长恶化 |
| **Flat staging + no graduation task (recommended)** | 零 helper 复制，零 import rewrite，pipeline 最短，符合 Go 社区标准 layout | 所有 test 文件在同一目录，需要代码改动；`helpers.go` 并发修改可能产生 git merge conflict | **Selected: 与 Go 社区标准 layout 对齐，最彻底的简化** |
| Flat staging + marker-only graduation task | 零 helper 复制，pipeline 保持一致 | 多一个空转任务，增加 pipeline 开销 | Rejected: 任务不做任何事就不该存在 |
| Go workspace (go.work) | 保留子目录，共享 helpers 通过 module 引用 | 需改造 go.mod 为 library module，改动大，overkill | Rejected: 为 staging 问题引入 workspace 是杀鸡用牛刀 |
| 保持 nested + 改善 helper 复制 | 最小改动 | 本质问题未解决，helper 副本持续 diverge | Rejected: 治标不治本 |
| 保持 nested + symlink helpers | 子目录通过 symlink 引用根目录 helpers | 平台兼容性 (Windows symlink 需要 admin 权限) | Rejected: Windows 兼容性差 |

### 选择的方案 vs Industry Benchmarks

本 proposal 的 flat layout 与 Go 标准库、Kubernetes、Hugo 的测试组织方式一致。这些项目均将同一 package 的测试文件放在同一目录，通过共享 helper package 消除重复。Kubernetes 的 `test/e2e/framework/` 模式和 Bazel 的 `rules_go` 策略都证明：build/test system 应适配语言的原生约定，而非强加统一抽象。

唯一偏离的是 `staging-mode` 抽象层——这是 forge 多语言支持的现实需要（TypeScript 需要 nested staging），不与 Go 社区实践冲突。

### Cross-Domain Influences: No-Op Stage Skipping

The design of `staging-mode` and its DAG-shortening behavior was shaped by two specific mechanisms from other domains:

- **Database migration frameworks** (Flyway, Liquibase): These tools use a **checksum-based detection** mechanism -- not runtime flags -- to decide whether a migration is a no-op. The `staging-mode: flat` field in `manifest.yaml` borrows this idea: it is a **static, declarative property** of the profile (like a migration checksum), not a runtime condition. The code path reads this property once and structurally omits the graduation task from the DAG, exactly as Flyway structurally skips already-applied migrations rather than running and no-oping them.

- **CI/CD conditional stages** (GitHub Actions `if:` guards, GitLab CI `rules`): These systems use **declarative guard expressions** to prune stages from the DAG before execution begins. The `staging-mode` field adopts this pattern: rather than generating T-test-4 and then conditionally skipping it at runtime, the task generator checks `stagingMode == "flat"` and simply omits the task from the returned task list. The DAG is shorter by construction, not by runtime skip.

The design decision -- "remove the task from the DAG rather than generate a no-op" -- was directly informed by observing how these systems handle unnecessary stages: static detection + structural omission, not runtime branching.

### 单目录文件数量质疑

"所有 test 文件在同一目录会不会太乱？" -- 不会。Go 社区对此有共识。Hugo 的 `hugolib/` 目录包含 200+ 个 `_test.go` 文件仍可维护。当前 `tests/e2e/` 根目录只有 2 个 `_test.go` 文件。按每个 feature 1 个文件计算，20 个 feature = 20 个文件，远低于社区验证的上限。文件名 `<feature>_<type>_test.go` 自带分类信息。

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

`manifest.yaml` 的 `staging-mode` 字段已可被 `profile` package 解析。验证结果：`profile/embed.go` 使用 `gopkg.in/yaml.v3` 的 `yaml.Unmarshal`（非 strict 模式），`profileManifest` struct 当前仅包含 `Capabilities []string` 字段，未调用 `DisallowUnknownFields()`，无 schema validation 拦截未知字段。因此添加 `staging-mode` 字段不需要修改解码方式，只需在 `profileManifest` struct 中增加对应字段并暴露 `GetStagingMode(name) string` 函数。

## Scope

### In Scope

- `manifest.yaml` 加 `staging-mode` 字段（go-test、pytest、rust-test）
- `gen-test-scripts/SKILL.md` 适配 flat staging（仅 go-test；pytest/rust-test 的 `staging-mode: flat` 字段为声明性占位，`gen-test-scripts` 当前不会为这两个 profile 读取或处理该字段）
- `run-e2e-tests/SKILL.md` 适配 flat prerequisite + 自动写 marker
- `graduate-tests/SKILL.md` 增加 flat 模式跳过逻辑
- `profile/embed.go` 增加 staging-mode 解析
- `task/testgen.go` flat 模式跳过 graduate 任务 + 调整依赖链

### Out of Scope

- 已毕业 feature 的迁移（保持现状）
- java-junit、maestro、web-playwright 的 staging-mode 变更
- `features/` 目录下已有文件的清理
- `infer.go` 中的 `TypeTestPipelineGraduate` 类型保留（向后兼容已存在的任务文件）
- pytest/rust-test 的 flat 模式运行时行为（仅添加 manifest 字段和解析支持；`gen-test-scripts` 对这两个 profile 的 flat 模式输出路径适配不在本次 scope 内）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 根目录文件过多 | L | L | 文件名自带 feature 前缀，实际 20+ feature 才会有 20 个文件 |
| 与 nested 模式 feature 混合 | L | M | `run-e2e-tests` 按 feature 的 staging-mode 元数据确定测试发现路径：flat 模式搜索 `tests/e2e/<slug>_*.go`，nested 模式搜索 `tests/e2e/features/<slug>/`。同一 feature 不会同时出现在两个路径（由 `gen-test-scripts` 保证输出路径唯一） |
| 已有 helper 不够用需要扩展 | M | L | helpers.go 可通过 gen-test-scripts 的 merge 机制扩展 |
| testgen.go 依赖 profile package | M | M | `GetBreakdownTestTasks` 和 `GetQuickTestTasks` 通过显式函数参数接收 staging-mode（`stagingMode string`），调用方通过 `profile.GetStagingMode(profileName)` 读取 manifest。避免全局可变状态，支持单元测试直接传参 |
| 向后兼容性 | L | L | `staging-mode` 缺失时默认 `nested`，现有行为不变。已毕业的 feature 文件位置、marker、pipeline 逻辑均不受影响——`nested` 默认值保证 code path 完全一致 |
| helpers.go 并发修改产生 git merge conflict | M | M | helpers.go 由 gen-test-scripts 模板生成，merge 机制检测 conflict 时中止并输出 diff；feature 文件本身（`<slug>_cli_test.go`）独立，不会冲突。实际冲突概率低：helpers.go 改动频率远低于 feature 文件 |
| 共享 package e2e 的 test state 污染 | L | M | 所有测试文件同属 `package e2e`，共享 `init()` 和全局变量。当前惯例：测试函数使用 `t.Parallel()`，不修改全局状态；每个 test function 是独立 goroutine，局部变量隔离。未来需注意不在测试中修改 package-level 变量 |
| features/ 目录残留导致测试重复执行 | M | M | 切换 flat 模式后，旧 `features/<slug>/` 中的测试文件仍会被 `go test ./...` 发现并执行，同一 feature 的测试运行两遍，old nested 版本可能引用 diverge 的 helper。**Mitigation**：在 flat 模式启用前，手动清理已迁移 feature 的 `features/<slug>/` 目录（Out of Scope，但建议作为首个 flat feature 的前置步骤）；`run-e2e-tests` flat 模式下使用 `-run` 过滤只执行目标 feature 的测试函数 |

## Success Criteria

- [ ] go-test profile 生成测试文件直接到 `tests/e2e/` 根目录，无需复制 helpers
- [ ] 生成的测试文件编译通过，可直接使用 `helpers.go` 中的 `runCLI`、`parseBlock` 等函数
- [ ] flat 模式下不生成 T-test-4 / T-quick-4 任务
- [ ] run-e2e-tests 通过后自动写 graduation marker
- [ ] flat 模式下 `/run-e2e-tests` 检查 `tests/e2e/<slug>_*_test.go` 而非 `tests/e2e/features/<slug>/` 作为 prerequisite 路径
- [ ] flat 模式下 `/graduate-tests` 输出 "Flat staging mode: graduation not needed" 且不尝试文件迁移
- [ ] 已毕业的 feature 不受影响：在 `tests/e2e/` 目录下仅保留已毕业 feature 的文件（不含新 flat 文件），运行 `go test ./tests/e2e/... -tags=e2e` 的 pass/fail 结果与变更前一致；所有 `tests/e2e/.graduated/` 下已存在的 marker 文件内容不变且仍可被 `profile` package 正确读取
- [ ] T-test-4.5 / T-quick-5 的依赖链正确指向 T-test-3 / T-quick-3
- [ ] pytest 和 rust-test 的 `manifest.yaml` 包含 `staging-mode: flat` 字段，`profile.GetStagingMode()` 能正确解析该字段返回 `"flat"` 且无报错（注意：pytest/rust-test 的 flat 模式运行时行为为 Out of Scope，此处仅验证 manifest schema 支持和解析正确性）

## Next Steps

- 通过 `/quick` 流程生成实施任务
- 或者直接实施（预计 2-3 个任务）
