---
name: init-forge
description: 自动编译并安装 task-cli 工具
---

# /init-forge

自动编译安装 task-cli 工具。

## 流程

1. 检测操作系统（Windows/Linux/macOS）
2. 定位 task-cli 路径（zcode/task-cli/）
3. 调用对应安装脚本编译并安装
4. 提示用户重新打开终端

## 执行步骤

### Step 1: 定位 task-cli

```bash
# task-cli 在 zcode/task-cli/ 目录下
TASK_CLI_DIR="${CLAUDE_PROJECT_ROOT}/task-cli"
if [ ! -d "$TASK_CLI_DIR" ]; then
  echo "ERROR: task-cli not found at $TASK_CLI_DIR"
  exit 1
fi
echo "Found task-cli at: $TASK_CLI_DIR"
```

### Step 2: 检测操作系统并安装

**Windows (PowerShell):**
```powershell
cd $TASK_CLI_DIR
powershell -ExecutionPolicy Bypass -File scripts/install-local.ps1
```

**Linux/macOS:**
```bash
cd "$TASK_CLI_DIR" && bash scripts/install-local.sh
```

### Step 3: 验证安装

```bash
task --version
```

### Step 4: 提示用户

安装完成后，输出：

```
╔════════════════════════════════════════════════════════════════╗
║  ✅ task-cli 安装成功                                           ║
╠════════════════════════════════════════════════════════════════╣
║  请重新打开终端以刷新环境变量，然后运行:                           ║
║  task --version                                                ║
╚════════════════════════════════════════════════════════════════╝
```

## 错误处理

| 错误 | 解决方案 |
|------|----------|
| task-cli 未找到 | 确保 task-cli 目录存在于 zcode/ 下 |
| 编译失败 | 检查 Go 环境 |
| 权限错误 | 检查安装目录写入权限 |
