# Plan: `task index --feature <slug>`

## Context

AI agents currently generate `index.json` manually (via `breakdown-tasks`/`quick-tasks` skills), producing structural errors that require `task validate` to catch. The `.md` frontmatter already contains most structural fields, creating redundant maintenance. This plan moves index.json generation from AI to CLI code, making `.md` files the single source of truth and `index.json` a pure derived artifact.

### MR #59 基础设施

v3.0.0 分支已合并 test-profile-system（MR #59），提供：
- `task profile` 命令（resolve/set/detect）
- `pkg/profile` 包（config.go, detect.go）
- `gopkg.in/yaml.v3` 已是依赖
- 6 个 profile 在 `plugins/forge/profiles/`

本方案在此基础上：
1. 将 profile 数据从 `plugins/forge/profiles/` 迁入 `task-cli/pkg/profile/profiles/`（唯一源）
2. 新增 `task profile get <name> --<flag>` 供 skill 读取策略文件
3. 新增 `task index --feature <slug>` 从 .md frontmatter + 嵌入 profiles 生成 index.json

### 职责划分

| 职责 | 之前 | 本方案 |
|------|------|--------|
| 写业务 task .md | Skill | Skill（不变） |
| 写测试 task .md | Skill（手写 + 模板替换） | **CLI（从 embed 生成）** |
| 写 index.json | Skill（手拼 JSON） | **CLI（`task index`）** |
| 读 profile 策略文件 | Skill 直接读文件路径 | Skill 通过 `task profile get` |

## Architecture Overview

```
pkg/profile/
  config.go        — 已有：read/write .forge/config.yaml
  detect.go        — 已有：项目结构探测
  profiles/        — 新增：嵌入 profile 数据（从 plugins/forge/profiles/ 迁入）
  embed.go         — 新增：embed.FS + 加载 manifest/策略文件/模板
internal/cmd/
  profile.go       — 已有：新增 get 子命令
  index.go         — 新增：task index Cobra 命令
pkg/task/
  infer.go         — 新增：InferType（从 pkg/prompt 迁入，+pattern matching）
  frontmatter.go   — 新增：YAML frontmatter 解析
  testgen.go       — 新增：从嵌入 profiles 生成测试任务
  build.go         — 新增：BuildIndex 合并逻辑
```

## Tasks (execution order)

### T1: 迁移 profile 数据 + embed + `task profile get`

**Move:** `plugins/forge/profiles/*` → `task-cli/pkg/profile/profiles/`
**Create:** `pkg/profile/embed.go`, `pkg/profile/embed_test.go`
**Modify:** `internal/cmd/profile.go`（新增 get 子命令）

**Profile 数据（唯一源）：**

```
task-cli/pkg/profile/profiles/
  web-playwright/
    manifest.yaml
    generate.md
    run.md
    graduate.md
    justfile-recipes
    templates/
  go-test/
    ...
  maestro/
    ...
  java-junit/
    ...
  rust-test/
    ...
  pytest/
    ...
```

```go
//go:embed profiles/*/manifest.yaml profiles/*/generate.md profiles/*/run.md profiles/*/graduate.md profiles/*/justfile-recipes profiles/*/templates/*
var profileFS embed.FS
```

**`task profile get` 子命令：**

```
task profile get <name> --manifest       输出 manifest.yaml 内容
task profile get <name> --generate       输出 generate.md 内容
task profile get <name> --run            输出 run.md 内容
task profile get <name> --graduate       输出 graduate.md 内容
task profile get <name> --template <f>   输出 templates/<f> 内容
task profile get <name> --justfile       输出 justfile-recipes 内容
```

**Key functions in `embed.go`：**

```go
func GetManifest(name string) ([]byte, error)
func GetStrategy(name, kind string) ([]byte, error)   // kind: "generate", "run", "graduate"
func GetTemplate(name, filename string) ([]byte, error)
func GetJustfileRecipes(name string) ([]byte, error)
func ListProfileTemplates(name string) ([]string, error)
```

**Skill 消费方式变化：**

之前：`Read plugins/forge/profiles/go-test/generate.md`
之后：`task profile get go-test --generate`

**Tests:** embed 加载、manifest 解析、策略文件读取、未知 profile error、`task profile get` CLI 测试

### T2: Relocate InferType + pattern-based matching

**Create:** `pkg/task/infer.go`, `pkg/task/infer_test.go`
**Modify:** `pkg/prompt/prompt.go`（body 委托 `task.InferType`）

- Move `InferType()` from `pkg/prompt/prompt.go:191-216` to `pkg/task/infer.go`
- Add T-quick-* cases (5 个)
- Add profile-suffixed pattern matching：
  - `T-test-2[a-z]` / `T-quick-2[a-z]` → `TypeTestPipelineGenScripts`
  - `T-test-3[a-z]` / `T-quick-3[a-z]` → `TypeTestPipelineRun`
  - `T-test-4[a-z]` / `T-quick-4[a-z]` → `TypeTestPipelineGraduate`
- Keep `pkg/prompt.InferType` as thin wrapper
- InferType 降级为向后兼容兜底（BuildIndex 优先读 frontmatter `type`）

### T3: YAML frontmatter parser

**Create:** `pkg/task/frontmatter.go`, `pkg/task/frontmatter_test.go`

`gopkg.in/yaml.v3` 已是依赖（MR #59 引入），无需新增。

```go
type FrontmatterData struct {
    ID            string   `yaml:"id"`
    Title         string   `yaml:"title"`
    Priority      string   `yaml:"priority"`
    EstimatedTime string   `yaml:"estimated_time"`
    Dependencies  []string `yaml:"dependencies"`
    Scope         string   `yaml:"scope"`
    Breaking      bool     `yaml:"breaking"`
    MainSession   bool     `yaml:"mainSession"`
    NoTest        bool     `yaml:"noTest"`
    Type          string   `yaml:"type"`
}

func ParseFrontmatter(content []byte) (FrontmatterData, []byte, error)
```

### T4: Test task generation from profiles

**Create:** `pkg/task/testgen.go`, `pkg/task/testgen_test.go`

**测试任务展开规则（D9）：**

Breakdown 模式：
```
共享任务（单）：
  T-test-1    gen-test-cases
  T-test-1b   eval-test-cases
  T-test-4.5  verify-regression
  T-test-5    consolidate-specs

按 profile 展开：
  T-test-2<L>  gen-test-scripts (profile-N)   → L = a,b,c...
  T-test-3<L>  run-tests (profile-N)
  T-test-4<L>  graduate-tests (profile-N)
```

Quick 模式：
```
按 profile 展开：
  T-quick-1<L>  gen-test-cases (profile-N)    → L = a,b,c...
  T-quick-2<L>  gen-test-scripts (profile-N)
  T-quick-3<L>  run-tests (profile-N)
  T-quick-4<L>  graduate-tests (profile-N)

共享任务：
  T-quick-5    verify-regression
```

注意：quick 模式无 eval-test-cases 和 consolidate-specs（与 breakdown 的区别）。

**单 profile 时无后缀**（`T-test-2` 而非 `T-test-2a`），保持向后兼容。

**Key functions:**

```go
func GetBreakdownTestTasks(profiles []string) []TestTaskDef
func GetQuickTestTasks(profiles []string) []TestTaskDef
func GenerateTestTaskMD(def TestTaskDef, slug string) ([]byte, error)
```

- `GenerateTestTaskMD`: frontmatter（机械推导）+ body（从 embed 读取对应策略文件或通用模板）
- 共享任务 body 使用嵌入的通用模板
- Per-profile 任务的 body 使用 `profile.GetStrategy()` 读取

**Fallback：**
- `--no-test` flag：跳过测试任务生成
- `--test-profiles` flag：手动传入 profile 列表（覆盖 config.yaml）

### T5: BuildIndex merge logic (core)

**Create:** `pkg/task/build.go`, `pkg/task/build_test.go`

```go
type BuildIndexOpts struct {
    FeatureSlug  string
    ProjectRoot  string
    TasksDir     string
    IndexPath    string
    NoTest       bool
    TestProfiles []string   // flag > config.yaml > none
}

type BuildIndexResult struct {
    NewCount, UpdatedCount, PreservedCount int
    Warnings []string
}

func BuildIndex(opts BuildIndexOpts) (*BuildIndexResult, error)
```

**BuildIndex flow:**
1. Load existing index.json (or `NewTaskIndex`). Preserve `created` date.
2. Set feature metadata（检测 prd/design/proposal 存在性）
3. Detect mode: `prd/prd-spec.md` → breakdown; `proposal.md` → quick
4. Resolve profiles: flag > `.forge/config.yaml` > none
5. Scan `tasksDir` for `*.md` (exclude `_templates/`, `records/`, `process/`)
6. For each `.md`: parse frontmatter → build Task → merge（保留 status/sourceTaskID/blockedReason）
7. Detect orphans → WARNING
8. Unless `--no-test`: generate test tasks per mode + profiles
9. Save index
10. Return result

**Type 推导优先级：** frontmatter `type` > `InferType(id)` 兜底

### T6: `task index` CLI command

**Create:** `internal/cmd/index.go`, `internal/cmd/index_test.go`
**Modify:** `internal/cmd/root.go`（注册 indexCmd）

```go
var indexCmd = &cobra.Command{
    Use:   "index --feature <slug> [--no-test] [--test-profiles p1,p2]",
    Short: "Build or rebuild index.json from task markdown files",
    Run:   runIndex,
}
```

Flags: `--feature` (required), `--no-test`, `--test-profiles`

`runIndex` flow: resolve paths → resolve profiles → `BuildIndex` → print summary → validate → print hint

### T7: Refactor `task add` to call BuildIndex

**Modify:** `internal/cmd/add.go`

`AddTask` 保持签名不变（backward compat）。`executeAdd` 在 `CreateTaskMarkdown` 后调用 `BuildIndex`。

### T8: Version bump + doc sync

- `scripts/version.txt`: bump minor
- `docs/OVERVIEW.md`: add `task index`, `task profile get`
- `docs/WORKFLOW.md`: update workflow

### T9: 删除 `plugins/forge/profiles/` + 更新 skill 引用

**Delete:** `plugins/forge/profiles/`（数据已迁入 task-cli）
**Modify:** 所有引用 `plugins/forge/profiles/<name>/` 的 skill prompt → 改为 `task profile get <name> --<flag>`

涉及 skills：
- `breakdown-tasks` — 移除 Step 5（手拼 index.json），改为调用 `task index`
- `quick-tasks` — 同上
- `gen-test-scripts` — 读策略文件改为 `task profile get <name> --generate`
- `run-tests` — 读策略文件改为 `task profile get <name> --run`
- `graduate-tests` — 读策略文件改为 `task profile get <name> --graduate`
- `init-justfile` — 读 justfile-recipes 改为 `task profile get <name> --justfile`
- `tech-design` — profile 选择步骤调整

## Reused Functions

| Function | Location | Used In |
|----------|----------|---------|
| `profile.ReadTestProfiles` | `pkg/profile/config.go` | `build.go` |
| `profile.DetectProfiles` | `pkg/profile/detect.go` | `build.go` |
| `profile.GetStrategy/GetManifest` | `pkg/profile/embed.go` (new) | `testgen.go`, CLI `profile get` |
| `InferType` | `pkg/task/infer.go` (new) | `build.go` |
| `LoadIndex`/`SaveIndex` | `pkg/task/index.go` | `build.go` |
| `ParseFrontmatter` | `pkg/task/frontmatter.go` (new) | `build.go` |
| `PrintBlockStart/End` | `internal/cmd/output.go` | `cmd/index.go` |

## Verification

1. `go build ./...` — compiles
2. `go test -race -cover ./...` — all pass, 80%+ on new code
3. `golangci-lint run ./...` — clean
4. `bash ../claude-code-go/scripts/lint-arch.sh` — dependency direction OK
5. Manual: `task profile get web-playwright --generate` — outputs strategy content
6. Manual: create business .md + `.forge/config.yaml`, run `task index --feature <slug>`, verify test tasks + index.json
7. Manual: re-run `task index`, verify idempotent
8. `make check-docs` — docs in sync
