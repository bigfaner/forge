---
created: 2026-05-13
author: faner
status: Draft
---

# Proposal: config.yaml 增加 project-type 字段

## Problem

项目的类型信息（frontend/backend/mixed）散落在 justfile 的 `project-type` recipe 中。每次 quality gate scope resolution 都要 `exec.Command("just", "project-type")` spawn 一个进程来读取一个静态字符串。

更深层的问题：**枚举值与 profile capabilities 词汇表不统一**。

| 来源 | 枚举值 |
|------|--------|
| `just project-type` | `frontend`, `backend`, `mixed` |
| Profile capabilities | `tui`, `web-ui`, `mobile-ui`, `api`, `cli` |

一个 Go TUI 项目说自己是 `backend`，但 profile capabilities 说它支持 `tui, api, cli`。两套词汇描述同一个东西，无法交叉过滤。

**根因**：项目元数据（project-type）和 profile 元数据（capabilities）独立演化，没有统一词汇表。

## Solution

### D1. 统一枚举值

`project-type` 采用与 profile capabilities 相同的词汇：

| 值 | 含义 | 旧值映射 |
|---|---|---|
| `tui` | 终端 UI（文本渲染、键盘交互） | `backend` 的子集 |
| `web-ui` | 浏览器 UI（DOM 交互） | `frontend` |
| `mobile-ui` | 移动端 UI（触摸、手势） | 无 |
| `api` | HTTP/网络接口 | `backend` 的子集 |
| `cli` | 命令行界面 | `backend` 的子集 |
| `mixed` | 多种类型混合 | `mixed` |

### D2. 配置位置：`.forge/config.yaml`

```yaml
test-profiles:
  - go-test
project-type: tui
```

- 与 `test-profiles` 统一在 `.forge/config.yaml`
- 声明性配置，提交到 git
- 单值（不是数组）——一个项目的主类型

### D3. 自动探测

`task profile detect` 同时探测 project-type：

| 信号 | 推断 project-type |
|---|---|
| `package.json` + 浏览器框架 (react/vue/svelte) | `web-ui` |
| `go.mod` + 无 web 框架 | `tui` |
| `go.mod` + gin/echo/fiber | `api` |
| `go.mod` + cobra/urfave | `cli` |
| `android/` 或 `ios/` | `mobile-ui` |
| 混合信号（frontend/ + backend/ 目录） | `mixed` |

### D4. ResolveScope 改读 config.yaml

```go
// 当前：spawn just process
output, success := RunCapture(projectRoot, "just", "project-type")

// 改为：读 config.yaml
cfg, _ := profile.ReadForgeConfig(projectRoot)
projectType := cfg.ProjectType
```

消除 scope resolution 对 `just` 的依赖。

### D5. test case 按 project-type 过滤

`gen-test-cases` 读 config.yaml 的 `project-type`，与 profile 的 `capabilities` 交叉：

- project-type = `tui`，profile supports `[tui, api, cli]` → 生成 tui/api/cli 相关 test case
- project-type = `tui`，profile supports `[web-ui, api, cli]` → 只生成 api/cli 相关 test case（web-ui 不匹配）
- project-type = `mixed` → 不过滤，生成全部

### D6. init-justfile 移除 `project-type` recipe

不再生成：

```makefile
# project-type: return project type identifier
project-type:
    @echo "backend"
```

`ResolveScope` 直接读 config.yaml，不调用 just。

### D7. 向后兼容

`ResolveScope` 迁移期支持旧值：

| 旧值 | 映射到 |
|---|---|
| `frontend` | `web-ui` |
| `backend` | 取决于探测（优先 `tui`，其次 `api`/`cli`） |

一个 major version 后移除旧值支持。

## Impact Assessment

### 改动的代码

| 文件 | 改动 |
|---|---|
| `pkg/profile/config.go` | `ForgeConfig` 增加 `ProjectType string` |
| `pkg/just/just.go` | `ResolveScope` 改读 config.yaml |
| `internal/cmd/profile.go` | `detect` / `set` 命令支持 project-type |
| `init-justfile` skill | 移除 `project-type` recipe 生成 |

### 新增的配置

`.forge/config.yaml` 新增 `project-type` 字段。

### 不变的部分

- Quality Gate 流程不变（只是 scope resolution 来源变了）
- Task scope 字段不变（`frontend`/`backend`/`all`）
- Profile manifest 的 capabilities 不变

## Risks

| 风险 | 缓解 |
|---|---|
| 旧项目 justfile 有 `project-type` recipe 但无 config.yaml | 迁移期 fallback 到 `just project-type` |
| 单值 `project-type` 不够描述 monorepo | `mixed` 覆盖；未来可扩展为数组 |
| `backend` 映射模糊（tui? api? cli?） | `task profile detect` 自动探测，用户确认 |
