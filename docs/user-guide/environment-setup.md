# 环境配置指南

> 最后更新：2026-05-30 | 对应版本：v3.0.0

本指南帮助新用户从零开始配置 Forge 开发环境。完成本指南后，你将能够使用 Forge 的全部功能。

---

## 前置条件

在安装 Forge 之前，请确保系统满足以下要求。

### 操作系统

Forge 支持以下操作系统和架构：

| 操作系统 | 支持的架构 |
|---------|-----------|
| macOS | x86_64 (amd64)、Apple Silicon (arm64) |
| Linux | x86_64 (amd64)、ARM64 |
| Windows | x86_64 (amd64)、ARM64 |

### curl

Forge CLI 通过 curl 下载预编译二进制文件。macOS 和 Linux 自带 curl，Windows 10+ 也已内置。

### Claude Code CLI

Forge 是 Claude Code 的插件，必须先安装 Claude Code。

**安装 Claude Code：**

```bash
npm install -g @anthropic-ai/claude-code
```

**验证安装：**

```bash
claude --version
```

### just（任务运行器）

Forge 依赖 `just` 作为构建任务运行器。`forge init` 会在初始化时自动引导安装，你也可以提前手动安装。

**手动安装 just：**

```bash
# macOS
brew install just

# Linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to ~/.local/bin

# Windows
winget install --id Casey.Just
```

---

## 安装方式

### 方式一：Binary 安装（推荐）

适用于大多数用户，下载预编译二进制文件，无需安装 Go。

**步骤：**

1. 安装 forge CLI（macOS / Linux）：

```bash
curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash
```

Windows 用户（PowerShell）：

```powershell
irm https://github.com/bigfaner/forge/releases/latest/download/install.ps1 | iex
```

2. 刷新终端环境变量：

```bash
# zsh 用户
source ~/.zshrc

# bash 用户
source ~/.bashrc

# Windows — 重新打开终端
```

3. 安装 forge Plugin（CLI + Plugin 一步到位）：

```bash
forge upgrade
```

4. 验证安装：

```bash
forge --version
```

> **后续升级：** 运行 `forge upgrade` 即可同时升级 CLI binary 和 Plugin。

### 方式二：本地构建安装（开发者）

适用于 Forge 的贡献者和开发者，需要从源码构建并迭代开发。

**前置条件：** [Go 1.25+](https://golang.org/dl/)

**步骤：**

1. 克隆仓库并构建安装 CLI：

```bash
git clone git@github.com:bigfaner/forge.git
cd forge

# Linux / macOS
cd forge-cli && bash scripts/install-local.sh

# Windows (PowerShell)
cd forge-cli
powershell -ExecutionPolicy Bypass -File scripts/install-local.ps1
```

2. 刷新终端环境变量：

```bash
# zsh 用户
source ~/.zshrc

# bash 用户
source ~/.bashrc

# Windows — 重新打开终端
```

3. 安装 Plugin：

```bash
forge upgrade
```

4. 验证安装：

```bash
forge --version
```

> **开发提示：** 修改 CLI 源码后，需要重新执行步骤 1 中的 `bash scripts/install-local.sh` 重新构建安装。Forge CLI 安装位置为 `~/.forge/bin/forge`。

---

## 安装后验证

完成安装后，运行以下命令确认环境配置正确。

### 1. 检查 Forge CLI 版本

```bash
forge --version
```

应输出当前安装的版本号（如 `5.16.0`）。

### 2. 检查 Forge 帮助信息

```bash
forge --help
```

应显示所有可用命令列表，包括 `task`、`worktree`、`surfaces` 等。

### 3. 初始化项目

在项目目录中运行：

```bash
forge init
```

此命令将：
- 检测操作系统和项目类型
- 引导安装 `just`（如尚未安装）
- 创建 `.forge/config.yaml` 配置文件
- 通过 TUI 交互确认 surface 类型

### 4. 验证 just 安装

```bash
just --version
```

应输出 just 版本号。Forge 的构建、测试、代码质量检查等操作均依赖 just。

---

## 常见问题

### Go 版本不兼容（仅开发者构建）

**症状：** 使用 `install-local.sh` 构建时报错 `go: ...: module ... requires go >= 1.25`。

**解决方案：** Binary 安装用户不受影响。开发者需要升级 Go 至 1.25+：

```bash
# macOS
brew upgrade go

# Linux — 从 https://golang.org/dl/ 下载最新版本
# Windows — 从 https://golang.org/dl/ 下载安装包
```

### Claude Code 未安装或未找到

**症状：** 运行 `/plugin` 命令时提示 `command not found` 或无响应。

**解决方案：**

1. 确认 Claude Code 已安装：

```bash
claude --version
```

2. 如果未安装：

```bash
npm install -g @anthropic-ai/claude-code
```

3. 确保 Node.js 版本 >= 18：

```bash
node --version
```

### 权限问题

**症状：** 安装时报错 `Permission denied` 或 `forge --version` 提示命令未找到。

**解决方案：**

1. 检查 `~/.forge/bin` 目录权限：

```bash
ls -la ~/.forge/bin/
```

2. 如果目录不存在或权限不足：

```bash
mkdir -p ~/.forge/bin
chmod 755 ~/.forge/bin
```

3. 确保 `~/.forge/bin` 在 PATH 中：

```bash
echo $PATH | grep -o '[^:]*forge[^:]*'
```

4. 如果 PATH 中没有该路径，手动添加：

```bash
# zsh 用户
echo 'export PATH="${PATH}:${HOME}/.forge/bin"' >> ~/.zshrc
source ~/.zshrc

# bash 用户
echo 'export PATH="${PATH}:${HOME}/.forge/bin"' >> ~/.bashrc
source ~/.bashrc
```

5. Windows 用户检查系统 PATH 是否包含 `%USERPROFILE%\.forge\bin`。
