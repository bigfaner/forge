---
title: "Naming Conventions"
domains: [file-names, function-names, constant-names, package-names, go]
---

# Naming Conventions

本文档定义 forge-cli 的命名规范。这是**规范性（normative）**文档，描述目标状态而非当前状态。所有新增代码必须符合此规范；已有代码的偏差记录在模块级偏差摘要中，按计划逐步收敛。

## 1. 文件名规则

### 1.1 Go 源文件

| 规则 | 目标态 | 示例 |
|------|--------|------|
| 使用 `snake_case` | 全小写，单词间用下划线连接 | `quality_gate.go`, `init_surfaces.go`, `verify_task_done.go` |
| 测试文件以 `_test.go` 结尾 | 与被测文件同名加 `_test` 后缀 | `quality_gate_test.go`, `init_surfaces_test.go` |
| 一文件一职责 | 文件内容围绕单一概念或命令 | `cleanup.go` 只含清理逻辑，`proposal.go` 只含 proposal 命令 |

**cobra 命令文件的命名模式**：

| 位置 | 命名模式 | 说明 |
|------|----------|------|
| `internal/cmd/<command>.go` | 命令名直接作为文件名 | `cleanup.go`, `research.go`, `lesson.go` |
| `internal/cmd/<command>_<subcommand>.go` | 命令与子命令用下划线连接 | `init_surfaces.go`, `surfaces_detect.go` |
| `internal/cmd/<group>/cmd_<subcommand>.go` | 子包中子命令文件加 `cmd_` 前缀 | `worktree/cmd_start.go`, `worktree/cmd_list.go` |
| `internal/cmd/<group>/register.go` | cobra 命令注册入口 | 每个 cmd 子包必须有 |

**pkg/ 文件的命名模式**：

| 规则 | 示例 |
|------|------|
| 包的主文件以包名命名 | `facttable/facttable.go`, `serverprobe/serverprobe.go` |
| 按职责拆分文件 | `feature/paths.go`（路径计算）、`feature/constants.go`（常量）、`task/category.go`（分类逻辑） |
| 平台特定文件加后缀 | `index/lock_unix.go`, `index/lock_windows.go` |

### 1.2 验证方式

```bash
# 检查是否存在大写字母或连字符的文件名（不含平台后缀）
find forge-cli/ -name '*.go' | grep -E '[A-Z-]' | grep -v _test.go

# 检查 _test.go 文件是否有对应的源文件
find forge-cli/ -name '*_test.go' | while read f; do
  base="${f%_test.go}.go"
  [ ! -f "$base" ] && echo "Orphan test: $f"
done
```

## 2. 函数名规则

### 2.1 导出函数

| 规则 | 目标态 | 示例 |
|------|--------|------|
| 使用 `PascalCase` | 大驼峰，首字母大写 | `LoadIndex`, `ResolveSourceTask`, `GenerateNonce` |
| 名字自描述 | 函数名传达意图而非实现 | `RequireFeature`（非 `GetFeatureOrError`） |
| 动词开头 | 函数名以动词开头表示动作 | `Load`, `Save`, `Ensure`, `Validate`, `Read`, `Write` |

**常用动词约定**：

| 动词 | 含义 | 示例 |
|------|------|------|
| `Get` | 查询/获取值，可能返回零值 | `GetCurrentFeature`, `GetFeatureDir` |
| `Require` | 获取值，失败时返回错误 | `RequireFeature` |
| `Ensure` | 确保存在，不存在则创建 | `EnsureFeatureDir`, `EnsureForgeDir` |
| `Resolve` | 解析/转换一个值为另一个 | `ResolveSourceTask`, `resolveProjectRoot` |
| `Validate` | 校验合法性 | `ValidateSurfaceTypes`, `ValidateAutogenTemplates` |
| `Load` | 从持久化介质读取 | `LoadIndex`, `Load` |
| `Save` | 写入持久化介质 | `SaveIndex`, `Save` |
| `Read` | 读取但不解码为结构体 | `ReadSurfaces`, `ReadExecutionOrder` |
| `Write` | 写入原始数据 | `WriteForgeState`, `writeConfigFile` |
| `Detect` | 运行检测逻辑 | `DetectSurface` |
| `Format` | 格式化输出 | `formatDurationMs`, `formatValue` |
| `Parse` | 解析字符串为结构 | `parseYAMLTagName`, `parseTimestamp` |
| `Build` | 构建复杂对象 | `buildTaskMarkdown`, `buildDisplayLines` |
| `Generate` | 生成新内容 | `GenerateNonce`, `GenerateTestTaskMD` |
| `Check` | 布尔检查或轻量校验 | `checkDependenciesMet` |
| `Is` / `Has` | 布尔判断（返回 bool） | `isLeafType`, `hasActiveFixTasks` |

### 2.2 未导出函数

| 规则 | 目标态 | 示例 |
|------|--------|------|
| 使用 `camelCase` | 小驼峰，首字母小写 | `resolveProjectRoot`, `buildGitignoreAppend` |
| 同样遵循动词约定 | 与导出函数一致 | `detectModeFromPath`, `normalizeSurfaceKeyValue` |

### 2.3 cobra 命令运行函数

| 规则 | 目标态 | 示例 |
|------|--------|------|
| 顶层命令：`run<Command>` | 动词 + 命令名 | `runInit`, `runCleanup`, `runResearch` |
| 子命令：`run<Command><Subcommand>` | 动词 + 命令 + 子命令 | `runConfigSet`, `runConfigGet`, `runSurfacesList` |
| 子包命令：`run<Subcommand>` | 动词 + 子命令（子包名已在路径中） | `worktree/` 中用 `runWorktreeRemove`, `task/` 中用 `runAdd`, `runClaim` |

### 2.4 验证方式

```bash
# 检查导出函数是否以标准动词开头（人工 review 项）
grep -rn 'func [A-Z]' forge-cli/pkg/*/*.go | grep -v _test.go | grep -v 'func ('

# 检查 cobra 运行函数命名一致性
grep -rn 'func run[A-Z]' forge-cli/internal/cmd/
```

## 3. 常量名规则

### 3.1 导出常量

| 规则 | 目标态 | 示例 |
|------|--------|------|
| 使用 `PascalCase` 或 `UPPER_SNAKE_CASE` | 视语义选择 | 枚举值用 `PascalCase`，物理常量/路径用 `UPPER_SNAKE_CASE` |
| 枚举常量用 `PascalCase` | 类型名 + 值名 | `StatusPending`, `PriorityP0`, `SourceStatic`, `KindSignature` |
| 路径/文件名常量用 `PascalCase` | 描述性名称 | `FeaturesDir`, `StateFileName`, `RecordFileName` |
| 常量名自包含类型信息 | 名字本身说明含义 | `ConfidenceConfirmed`, `KindErrorCode` |

**枚举常量前缀约定**：

```
<TypeName><ValueName>
```

示例：
```go
const (
    StatusPending    Status = "pending"
    StatusCompleted  Status = "completed"
)

const (
    PriorityP0 Priority = "P0"
    PriorityP1 Priority = "P1"
)

const (
    SourceStatic  = "static"
    SourceRuntime = "runtime"
)
```

### 3.2 未导出常量

| 规则 | 目标态 | 示例 |
|------|--------|------|
| 使用 `camelCase` | 小驼峰 | `maxFixTasksPerStep`, `defaultLockTimeout`, `escapeHatchLimit` |
| 魔法数字必须提取为常量 | 不允许裸数字/字符串出现在逻辑中 | `const maxFixTasksPerStep = 3` |

### 3.3 验证方式

```bash
# 检查是否存在 iota 用法（当前项目不使用 iota，使用显式赋值）
grep -rn 'iota' forge-cli/

# 检查 pkg/types/ 中枚举常量是否遵循 <Type><Value> 模式
grep -rn 'const (' forge-cli/pkg/types/

# 检查魔法数字（人工 review 项，重点关注 > 0 和 != 0 的裸数字）
grep -rn '[^0-9a-zA-Z_][2-9][0-9]*[^0-9a-zA-Z_]' forge-cli/pkg/*/ forge-cli/internal/ | grep -v _test.go | grep -v 'const '
```

## 4. 包名规则

### 4.1 基本规则

| 规则 | 目标态 | 示例 |
|------|--------|------|
| 单个小写单词 | 不使用下划线、不使用驼峰 | `git`, `task`, `feature`, `just` |
| 简短、有意义 | 一个词即可表达包的职责 | `index`, `lesson`, `prompt` |
| 不使用 `common` / `util` / `helpers` | 避免无意义名称 | 用 `base` 代替 `common`（仅限 cmd 基础设施） |
| 避免与标准库冲突 | 如冲突则加限定词 | `index` 而非 `sync`, `git` 而非 `http` |

### 4.2 包名与目录名一致

| 规则 | 目标态 |
|------|--------|
| 目录名 = 包名 | `pkg/git/` 目录中 `package git` |
| 唯一例外：`internal/cmd/<group>/` 子包 | 目录名是命令组名，包名也是命令组名（`package task`, `package worktree`） |

### 4.3 pkg/ 各包命名说明

| 包名 | 含义 | 命名合理性 |
|------|------|-----------|
| `types` | 共享类型定义 | 标准 Go 做法，leaf 包 |
| `git` | Git 操作封装 | 简洁，无冲突 |
| `index` | 索引文件管理 | 简洁，描述性强 |
| `infocmd` | 命令元数据描述 | info + cmd 的复合词，见偏差 N1 |
| `just` | Just 任务运行器 | 对标工具名 |
| `project` | 项目根目录检测 | 标准 Go 命名 |
| `facttable` | 事实表数据结构 | 见偏差 N2 |
| `forgeconfig` | Forge 配置文件解析 | 见偏差 N3 |
| `feature` | Feature 生命周期管理 | 标准 Go 命名 |
| `task` | 任务管理 | 标准 Go 命名 |
| `prompt` | Prompt 模板管理 | 标准 Go 命名 |
| `proposal` | Proposal 文档管理 | 标准 Go 命名 |
| `serverprobe` | 服务器探测 | 见偏差 N2 |
| `testrunner` | 测试运行器 | 见偏差 N2 |

### 4.4 验证方式

```bash
# 检查目录名与 package 声明是否一致
for dir in forge-cli/pkg/*/; do
  pkg=$(basename "$dir")
  head -5 "${dir}"*.go 2>/dev/null | grep "package $pkg" > /dev/null || \
    echo "Mismatch: dir=$pkg, package=$(grep '^package' ${dir}*.go | head -1)"
done

# 检查包名是否包含下划线或大写
find forge-cli/pkg/ -name '*.go' -exec grep '^package ' {} \; | grep -E '[A-Z_]' && echo "VIOLATION"

# go vet 静态检查
cd forge-cli && go vet ./...
```

## 5. 模块级偏差摘要

以下是当前代码与目标状态之间的命名偏差，按领域分类。

### 偏差 N1：复合词包名

| 包 | 当前命名 | 问题 | 取舍 |
|---|---------|------|------|
| `pkg/infocmd` | 两个词拼接 `info` + `cmd` | 违反"单个小写单词"规则 | **保留**：`info` 过于宽泛，`cmd` 与标准库 `os/exec` 语境冲突。`infocmd` 在项目内有明确语义（命令元数据信息），且已有 4 个消费者，重命名成本高收益低 |
| `pkg/forgeconfig` | 两个词拼接 `forge` + `config` | 违反"单个小写单词"规则 | **保留**：`config` 过于通用（Go 生态中有大量同名包），加 `forge` 前缀消除歧义。该包是 infrastructure 层核心包，重命名影响面大 |

### 偏差 N2：描述性复合词包名

| 包 | 当前命名 | 问题 | 取舍 |
|---|---------|------|------|
| `pkg/facttable` | 两个词拼接 `fact` + `table` | 违反"单个小写单词"规则 | **保留**：`fact` 过于通用，`table` 与 `database/sql` 语境冲突。`facttable` 明确指向事实表这一特定领域概念 |
| `pkg/serverprobe` | 两个词拼接 `server` + `probe` | 违反"单个小写单词"规则 | **保留**：`probe` 或 `server` 单独使用都不足以表达"服务器健康检查"语义 |
| `pkg/testrunner` | 两个词拼接 `test` + `runner` | 违反"单个小写单词"规则 | **保留**：`testing` 与标准库冲突，`runner` 过于通用。`testrunner` 明确指向测试运行器 |

**总结**：4 个复合词包名均因单单词替代方案存在歧义或冲突而保留。新增包应优先选择单单词命名；当单单词无法消除歧义时，允许使用两个小写词拼接，但必须在 PR 中说明理由。

### 偏差 N3：`internal/cmd/` 子包命令运行函数前缀不一致

| 位置 | 当前模式 | 目标模式 | 说明 |
|------|---------|---------|------|
| `internal/cmd/task/` | `runAdd`, `runClaim`, `runList` | 无需变更 | 子包内函数名不需要重复包名前缀，上下文已明确 |
| `internal/cmd/worktree/` | `runWorktreeRemove`, `runWorktreeList` | 可简化为 `runRemove`, `runList` | 加了 `Worktree` 前缀 — 当前保留，避免大范围重命名 |
| `internal/cmd/fact/` | `runList`, `runGet`, `runSummary` | 无需变更 | 符合目标模式 |
| `internal/cmd/forensic/` | `runExtract`, `runSearch`, `runSubagents` | 无需变更 | 符合目标模式 |

**结论**：`internal/cmd/worktree/` 中的 `runWorktree*` 前缀与目标模式不一致，但当前保留以避免不必要的大范围重命名。新代码应在子包中使用 `run<Subcommand>` 模式（不加包名前缀）。

## 6. 开发者工作流

### 6.1 新增文件

1. 确定文件名：使用 `snake_case`，反映文件内容的概念
2. cobra 命令文件：`<command>.go` 或 `<command>_<subcommand>.go`
3. pkg/ 文件：主文件用包名，辅助文件按职责命名

### 6.2 新增函数

1. 导出函数：`PascalCase`，以标准动词开头
2. 未导出函数：`camelCase`，同样以动词开头
3. cobra 运行函数：`run<Command>` 或 `run<Subcommand>`

### 6.3 新增常量

1. 枚举常量：`<Type><Value>` 模式（如 `StatusPending`），放在 `const` 块中
2. 内部常量：`camelCase`（如 `maxFixTasksPerStep`）
3. 不使用 `iota`：显式赋值，便于搜索和代码审查

### 6.4 新增包

1. 优先选择单个小写单词
2. 当单单词存在歧义时，允许两个小写词拼接
3. 目录名与 package 声明保持一致

## 7. 参考

- [Go Code Review Comments - naming](https://github.com/golang/go/wiki/CodeReviewComments#naming)
- [Effective Go - names](https://go.dev/doc/effective_go#names)
- [package-organization.md](./package-organization.md) -- 包组织规范
- [enum-constants.md](./enum-constants.md) -- 枚举常量组织规范
- [code-structure.md](./code-structure.md) -- 代码结构规范
