> ## Documentation Index
> Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> Use this file to discover all available pages before exploring further.

# 使用 worktrees 运行并行会话

> 在单独的 git worktrees 中隔离并行 Claude Code 会话，以便更改不会相互冲突。涵盖 `--worktree` 标志、子代理隔离、`.worktreeinclude`、清理和非 git VCS hooks。

[git worktree](https://git-scm.com/docs/git-worktree) 是一个单独的工作目录，具有自己的文件和分支，但与主检出共享相同的存储库历史和远程。在自己的 worktree 中运行每个 Claude Code 会话意味着一个会话中的编辑永远不会触及另一个会话中的文件，因此您可以让 Claude 在一个终端中构建功能，同时在第二个终端中修复错误。

本页涵盖 CLI 中的 worktree 隔离。下面的所有内容都假设使用 git 存储库。对于其他版本控制系统，请参阅[非 git 版本控制](#non-git-version-control)。[桌面应用](/zh-CN/desktop#work-in-parallel-with-sessions)会为每个新会话自动创建一个 worktree。

Worktrees 是运行 Claude 并行的几种方式之一。它们隔离文件编辑，而[子代理](/zh-CN/sub-agents)和[代理团队](/zh-CN/agent-teams)协调工作本身。请参阅[并行运行代理](/zh-CN/agents)来比较这些方法，或跳到[使用 worktrees 隔离子代理](#isolate-subagents-with-worktrees)以同时使用 worktrees 和子代理。

## 在 worktree 中启动 Claude

传递 `--worktree` 或 `-w` 来创建隔离的 worktree 并在其中启动 Claude。默认情况下，worktree 在您的存储库根目录下的 `.claude/worktrees/<value>/` 下创建，在名为 `worktree-<value>` 的新分支上：

```bash theme={null}
claude --worktree feature-auth
```

要将 worktrees 放在其他地方，请配置 [`WorktreeCreate` hook](#non-git-version-control)。在另一个终端中使用不同的名称再次运行该命令以启动第二个隔离会话：

```bash theme={null}
claude --worktree bugfix-123
```

如果您省略名称，Claude 会生成一个名称，例如 `bright-running-fox`：

```bash theme={null}
claude --worktree
```

您也可以在会话期间要求 Claude "在 worktree 中工作"，它将使用 [`EnterWorktree`](/zh-CN/tools-reference) 工具创建一个。

在首次在目录中使用 `--worktree` 之前，请通过在该目录中运行一次 `claude` 来接受工作区信任对话框。如果尚未接受信任，`--worktree` 将以错误退出并提示您首先在目录中运行 `claude`，包括与 `-p` 结合使用时。

<Tip>
  将 `.claude/worktrees/` 添加到您的 `.gitignore`，以便 worktree 内容不会在您的主检出中显示为未跟踪的文件。
</Tip>

### 选择基础分支

Worktrees 从您的存储库的默认分支 `origin/HEAD` 分支，因此它们从与远程匹配的干净树开始。如果未配置远程或获取失败，worktree 会回退到您当前的本地 `HEAD`。要始终从本地 `HEAD` 分支，请在[设置](/zh-CN/settings#worktree-settings)中将 `worktree.baseRef` 设置为 `"head"`。将 `baseRef` 设置为 `"head"` 会使新 worktrees 携带您未推送的提交和功能分支状态，这在隔离需要在进行中的工作上操作的子代理时很有用。该设置仅接受 `"fresh"` 或 `"head"`，不接受任意 git refs：

```json theme={null}
{
  "worktree": {
    "baseRef": "head"
  }
}
```

要从特定的拉取请求分支，请传递以 `#` 为前缀的 PR 编号或完整的 GitHub 拉取请求 URL。Claude Code 从 `origin` 获取 `pull/<number>/head` 并在 `.claude/worktrees/pr-<number>` 创建 worktree：

```bash theme={null}
claude --worktree "#1234"
```

要完全控制 worktrees 的创建方式，请配置 [`WorktreeCreate` hook](/zh-CN/hooks#worktreecreate)，它完全替代默认的 `git worktree` 逻辑。

## 将 gitignored 文件复制到 worktrees

Worktree 是一个新的检出，因此来自您主存储库的未跟踪文件（如 `.env` 或 `.env.local`）不存在。要在 Claude 创建 worktree 时自动复制它们，请将 `.worktreeinclude` 文件添加到您的项目根目录。

该文件使用 `.gitignore` 语法。只有匹配模式且也被 gitignored 的文件才会被复制，因此跟踪的文件永远不会被重复。

此 `.worktreeinclude` 将两个 env 文件和一个 secrets 配置复制到每个新 worktree：

```text .worktreeinclude theme={null}
.env
.env.local
config/secrets.json
```

这适用于使用 `--worktree` 创建的 worktrees、[子代理 worktrees](#isolate-subagents-with-worktrees) 和[桌面应用](/zh-CN/desktop#work-in-parallel-with-sessions)中的并行会话。

## 使用 worktrees 隔离子代理

子代理可以在自己的 worktrees 中运行，以便并行编辑不会冲突。要求 Claude "为您的代理使用 worktrees"，或通过向 frontmatter 添加 `isolation: worktree` 在[自定义子代理](/zh-CN/sub-agents#supported-frontmatter-fields)上永久设置它。每个子代理都会获得一个临时 worktree，当子代理完成且没有更改时会自动删除。

## 清理 worktrees

当您退出 worktree 会话时，清理取决于您是否进行了更改：

* **无更改**：worktree 及其分支会自动删除
* **存在更改或提交**：Claude 提示您保留或删除 worktree。保留会保留目录和分支，以便您稍后可以返回。删除会删除 worktree 目录及其分支，丢弃所有未提交的更改和提交
* **非交互式运行**：使用 `--worktree` 和 `-p` 创建的 worktrees 不会自动清理，因为没有退出提示。使用 `git worktree remove` 删除它们

由崩溃或中断的运行孤立的子代理 worktrees 在启动时会被删除，一旦它们的年龄超过您的 [`cleanupPeriodDays`](/zh-CN/settings#available-settings) 设置，前提是它们没有未提交的更改、没有未跟踪的文件和没有未推送的提交。使用 `--worktree` 创建的 Worktrees 永远不会被此扫描删除。

## 手动管理 worktrees

要完全控制 worktree 位置和分支配置，请直接使用 Git 创建 worktrees。当您需要检出特定的现有分支或将 worktree 放在存储库外时，这很有用。

在新分支上创建 worktree：

```bash theme={null}
git worktree add ../project-feature-a -b feature-a
```

从现有分支创建 worktree：

```bash theme={null}
git worktree add ../project-bugfix bugfix-123
```

在 worktree 中启动 Claude：

```bash theme={null}
cd ../project-feature-a && claude
```

列出您的 worktrees：

```bash theme={null}
git worktree list
```

完成后删除一个：

```bash theme={null}
git worktree remove ../project-feature-a
```

有关完整的命令参考，请参阅 [Git worktree 文档](https://git-scm.com/docs/git-worktree)。记住在每个新 worktree 中初始化您的开发环境：安装依赖项、设置虚拟环境或运行您的项目设置所需的任何内容。

## 非 git 版本控制

Worktree 隔离默认使用 git。对于 SVN、Perforce、Mercurial 或其他系统，请配置 [`WorktreeCreate` 和 `WorktreeRemove` hooks](/zh-CN/hooks#worktreecreate) 以提供自定义创建和清理逻辑。因为 hook 替代了默认的 git 行为，当您使用 `--worktree` 时，[`.worktreeinclude`](#copy-gitignored-files-into-worktrees) 不会被处理。改为在您的 hook 脚本内复制任何本地配置文件。

此 `WorktreeCreate` hook 从 stdin 读取 worktree 名称，检出一个新的 SVN 工作副本，并打印目录路径，以便 Claude Code 可以将其用作会话的工作目录：

```json theme={null}
{
  "hooks": {
    "WorktreeCreate": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash -c 'NAME=$(jq -r .name); DIR=\"$HOME/.claude/worktrees/$NAME\"; svn checkout https://svn.example.com/repo/trunk \"$DIR\" >&2 && echo \"$DIR\"'"
          }
        ]
      }
    ]
  }
}
```

将其与 `WorktreeRemove` hook 配对以在会话结束时进行清理。有关输入架构和删除示例，请参阅 [hooks 参考](/zh-CN/hooks#worktreecreate)。

## 另请参阅

Worktrees 处理文件隔离。下面的相关页面涵盖将工作委派到这些隔离的检出中以及在您创建的会话之间切换：

* [子代理](/zh-CN/sub-agents)：在会话内将工作委派给隔离的代理
* [代理团队](/zh-CN/agent-teams)：自动协调多个 Claude 会话
* [管理会话](/zh-CN/sessions)：命名、恢复和在对话之间切换
* [桌面并行会话](/zh-CN/desktop#work-in-parallel-with-sessions)：桌面应用中由 worktree 支持的会话
