---
created: 2026-06-01
author: faner
status: Approved
---

# Proposal: Forge CLI Binary Distribution

## Problem

当前 `/init-forge` 命令从 Go 源码编译安装 forge CLI，存在两个根本性缺陷：

1. **路径错误**：使用 `${CLAUDE_PROJECT_ROOT}/forge-cli` 定位源码，但 Marketplace 安装后源码在 `~/.claude/plugins/marketplaces/forge/forge-cli/`，三个路径无一匹配，用户必定安装失败
2. **Go 依赖门槛**：要求用户安装 Go 1.25+ 才能使用 forge，对非 Go 项目用户不合理

## Proposed Solution

翻转分发模型：**CLI 是入口，Plugin 是 CLI 管理的依赖**。

### 分发模型对比

| | 当前模型 | 新模型 |
|---|---|---|
| 入口 | Plugin（用户从 Marketplace 安装） | CLI binary（用户 curl 下载） |
| CLI 安装 | `/init-forge` 从源码编译 | `install.sh` 下载预编译 binary |
| Plugin 安装 | Claude Code Marketplace | `forge upgrade` 自动管理 |
| Go 依赖 | 必须 | 不需要 |

### 前置依赖

- **curl**：系统自带（macOS/Linux），Windows 10+ 自带
- **Claude Code CLI**：`forge upgrade` 内部调用 `claude` 命令管理 Plugin，用户需先安装 Claude Code

### 用户流程

**首次安装：**

```bash
# 步骤 1: 安装 forge CLI
curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash

# 步骤 2: 安装 forge Plugin（CLI + Plugin 一步到位）
forge upgrade

# 步骤 3: 在项目中初始化（已有命令，无需修改）
cd my-project && forge init
```

**后续升级：**

```bash
forge upgrade    # 一条命令同时升级 CLI binary 和 Plugin
```

## Implementation

### 1. install.sh — curl 一行安装脚本

新建 `forge-cli/scripts/install.sh`（与现有 `install-local.sh` 并存，后者供开发者本地编译使用）。随 GitHub Release 作为 asset 分发。

> **tag 使用 `v` 前缀（如 `forge-cli/v5.17.0`），binary 文件名不使用 `v` 前缀（如 `forge-5.17.0-darwin-arm64`）。**

```
URL: https://github.com/bigfaner/forge/releases/latest/download/install.sh
```

职责：

1. 检测 OS 和架构（darwin/linux × amd64/arm64，windows × amd64/arm64）
2. 从 GitHub Release API 获取最新版本号
3. 拼接下载 URL：tag = `forge-cli/v{version}`，完整 URL = `https://github.com/bigfaner/forge/releases/download/forge-cli/v{version}/forge-{version}-{os}-{arch}`
4. 下载 binary 到 `~/.forge/bin/forge.new`
5. 原子替换：`mv ~/.forge/bin/forge.new ~/.forge/bin/forge`
6. 添加 `~/.forge/bin/` 到 PATH（写入 shell RC 文件）
7. 输出验证指令

Windows 版本 `install.ps1` 同理：检测平台、拼接 URL（tag 格式与 binary 文件名规则相同）、下载到 `%USERPROFILE%\.forge\bin\`、更新 User PATH。

### 2. forge upgrade — CLI 新子命令

新的 CLI 子命令，统一处理首次安装和升级。

**前置条件**：`claude` CLI 在 PATH 中可用。

```
forge upgrade
  ├── CLI binary 升级:
  │   ├── 获取当前版本: forge --version
  │   ├── 获取最新版本: GitHub Release API（解析 tag forge-cli/v{version} 提取版本号）
  │   ├── 版本相同 → 跳过
  │   └── 版本不同 → 下载最新 binary → 原子替换 ~/.forge/bin/forge
  │       └── Windows 特殊处理: 先重命名旧 binary 为 forge.old，再写入新 binary，
  │           完成后删除 forge.old（Windows 不允许替换正在运行的 exe）
  │
  └── Plugin 管理:
      ├── 检测 marketplace 是否已添加:
      │   └── 未添加 → claude plugin marketplace add https://github.com/bigfaner/forge.git --sparse .claude-plugin plugins
      ├── 检测 plugin 是否已安装:
      │   ├── 未安装 → claude plugin install forge
      │   └── 已安装 → claude plugin update forge
      └── 输出结果
```

### 3. GitHub Actions Workflow — 自动构建发布

放在 `.github/workflows/release-cli.yml`。

**触发条件：**

```yaml
on:
  push:
    tags: ['forge-cli/v*']   # 匹配 forge-cli/v5.17.0
```

**发版流程：**

```bash
# 开发者（由 /release-cli 命令自动执行）：
# 1. 修改 forge-cli/scripts/version.txt: 5.16.0 → 5.17.0
# 2. git commit -m "chore(forge-cli): bump version to 5.17.0"
# 3. git tag forge-cli/v5.17.0
# 4. git push origin HEAD forge-cli/v5.17.0
# → Actions 自动构建 + 发布
```

**构建矩阵（6 平台并行）：**

| GOOS | GOARCH | 目标平台 |
|------|--------|---------|
| darwin | arm64 | macOS Apple Silicon |
| darwin | amd64 | macOS Intel |
| linux | arm64 | Linux ARM64 |
| linux | amd64 | Linux AMD64 |
| windows | amd64 | Windows AMD64 |
| windows | arm64 | Windows ARM64 |

**编译命令：**

```bash
VERSION=$(cat forge-cli/scripts/version.txt | tr -d '[:space:]')
LDFLAGS="-s -w -X forge-cli/pkg/types.Version=${VERSION}"
CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
  go build -ldflags="$LDFLAGS" -o forge ./cmd/forge
```

- `-s -w`: 去掉调试符号，减小 binary 体积（约 30%）
- `-X forge-cli/pkg/types.Version`: 注入版本号，`forge --version` 读取
- `CGO_ENABLED=0`: 纯静态编译，无外部 C 依赖

**Release 产物：**

```
forge-{version}-darwin-arm64
forge-{version}-darwin-amd64
forge-{version}-linux-arm64
forge-{version}-linux-amd64
forge-{version}-windows-amd64.exe
forge-{version}-windows-arm64.exe
install.sh
checksums.txt              # SHA256 校验
```

**Workflow 结构：**

```yaml
name: Release forge-cli

on:
  push:
    tags: ['forge-cli/v*']

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - { goos: darwin,  goarch: arm64 }
          - { goos: darwin,  goarch: amd64 }
          - { goos: linux,   goarch: arm64 }
          - { goos: linux,   goarch: amd64 }
          - { goos: windows, goarch: amd64 }
          - { goos: windows, goarch: arm64 }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - name: Build
        run: |
          VERSION=$(cat scripts/version.txt | tr -d '[:space:]')
          EXT=$([ "${{ matrix.goos }}" = "windows" ] && echo ".exe" || echo "")
          BINARY="forge-${VERSION}-${{ matrix.goos }}-${{ matrix.goarch }}${EXT}"
          LDFLAGS="-s -w -X forge-cli/pkg/types.Version=${VERSION}"
          CGO_ENABLED=0 GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} \
            go build -ldflags="$LDFLAGS" -o "bin/${BINARY}" ./cmd/forge
      - uses: actions/upload-artifact@v4
        with:
          name: forge-${{ matrix.goos }}-${{ matrix.goarch }}
          path: bin/forge-*

  release:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/download-artifact@v4
        with:
          path: artifacts/
          merge-multiple: true
      - name: Generate checksums
        run: |
          cd artifacts
          shasum -a 256 forge-* > checksums.txt
      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          files: |
            artifacts/forge-*
            artifacts/checksums.txt
            forge-cli/scripts/install.sh
          generate_release_notes: true
```

**下载 URL 模式：**

```
# 特定版本（tag 含 v 前缀，文件名不含 v 前缀）
https://github.com/bigfaner/forge/releases/download/forge-cli/v5.17.0/forge-5.17.0-darwin-arm64

# latest 别名（自动指向最新 Release）
https://github.com/bigfaner/forge/releases/latest/download/install.sh
https://github.com/bigfaner/forge/releases/latest/download/forge-5.17.0-darwin-arm64
```

### 4. /release-cli — CLI 发布命令

本地 `.claude/commands/release-cli.md` 命令，将版本号修改 + tag + push 自动化。

**开发者执行：**

```
/release-cli
```

**命令流程：**

```
/release-cli
  ├── 1. 读取 forge-cli/scripts/version.txt 获取当前版本
  ├── 2. 询问新版本号（建议 semver bump）
  ├── 3. 更新 forge-cli/scripts/version.txt
  ├── 4. 提交: git commit -m "chore(forge-cli): bump version to {version}"
  ├── 5. 打 tag: git tag forge-cli/v{version}
  └── 6. 推送: git push origin HEAD forge-cli/v{version}
       └── GitHub Actions 自动触发构建 + 发布 Release
```

**与现有 `/upgrade-forge` 的关系：**

| 命令 | 位置 | 职责 |
|------|------|------|
| `/upgrade-forge` | `.claude/commands/upgrade-forge.md` | Plugin 版本 bump（plugin.json + marketplace.json） |
| `/release-cli` | `.claude/commands/release-cli.md` | CLI 版本 bump + tag + push → 触发 CI 发布 |

两者独立运作，版本号独立管理（Plugin: 3.0.0-rc.x，CLI: 5.x.x）。CLI 版本号起始于 5.x.x 是历史累积，与 Plugin 版本无对应关系。

### 5. 删除 forge-cli/CLAUDE.md 中的 Version Bump 规则

当前 `forge-cli/CLAUDE.md` 第 29-33 行要求每次代码改动都手动 bump `scripts/version.txt`：

```markdown
### Version Bump

Code changes must bump the version in `scripts/version.txt`. Follow semver:
- Patch: bug fixes, dead code removal (x.y.Z)
- Minor: new features, new commands (x.Y.z)
- Major: breaking CLI changes (X.y.z)
```

**删除此规则**。版本号由 `/release-cli` 统一管理——开发者在准备发布时执行 `/release-cli`，由该命令负责 bump 版本、打 tag、推送。日常代码改动不再需要手动修改 `version.txt`。

### 6. 删除 /init-forge

从 `plugins/forge/commands/init-forge.md` 移除。安装职责由 `install.sh` + `forge upgrade` 承担。现有 `forge-cli/scripts/install-local.sh` 保留，供开发者本地编译使用。

## Component Summary

| 组件 | 类型 | 职责 |
|------|------|------|
| `forge-cli/scripts/install.sh` | 安装脚本 | curl 一行安装 CLI binary（新建） |
| `forge-cli/scripts/install.ps1` | 安装脚本 | Windows PowerShell 版安装脚本（新建） |
| `forge upgrade` | CLI 新子命令 | CLI binary 升级 + Plugin 安装/升级 |
| `.github/workflows/release-cli.yml` | CI | tag 触发，6 平台构建 + Release 发布（新建） |
| `.claude/commands/release-cli.md` | 本地命令 | CLI 版本 bump + tag + push，触发 CI 发布（新建） |
| 删除 `forge-cli/CLAUDE.md` Version Bump 规则 | 清理 | 版本号由 /release-cli 统一管理 |
| 删除 `plugins/forge/commands/init-forge.md` | 清理 | 移除旧的编译安装命令 |
| `forge-cli/scripts/install-local.sh` | 保留 | 开发者本地编译安装（不受影响） |

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────┐
│  开发者                                                       │
│                                                               │
│  Plugin 发布: /upgrade-forge → commit                         │
│  CLI 发布:   /release-cli   → commit + tag + push             │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│  GitHub Actions                                               │
│  tag forge-cli/v* 触发 → 6 平台并行编译 → Release 发布         │
└───────────────────────────┬──────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│  用户                                                         │
│                                                               │
│  首次:                                                        │
│    curl .../install.sh | bash   → 安装 forge CLI               │
│    forge upgrade                → 安装 forge Plugin             │
│    cd project && forge init     → 初始化项目（已有命令）         │
│                                                               │
│  后续:                                                        │
│    forge upgrade                → 升级 CLI + Plugin             │
└──────────────────────────────────────────────────────────────┘
```
