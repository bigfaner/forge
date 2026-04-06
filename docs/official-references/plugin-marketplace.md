> ## Documentation Index
> Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> Use this file to discover all available pages before exploring further.

# 创建和分发 plugin marketplace

> 构建和托管 plugin marketplace，以在团队和社区中分发 Claude Code 扩展。

**plugin marketplace** 是一个目录，让你能够将 plugins 分发给他人。Marketplace 提供集中式发现、版本跟踪、自动更新以及对多种源类型（git 存储库、本地路径等）的支持。本指南展示了如何创建自己的 marketplace，与你的团队或社区共享 plugins。

想要从现有 marketplace 安装 plugins？请参阅[发现和安装预构建的 plugins](/zh-CN/discover-plugins)。

## 概述

创建和分发 marketplace 涉及：

1. **创建 plugins**：使用命令、agents、hooks、MCP servers 或 LSP servers 构建一个或多个 plugins。本指南假设你已经有要分发的 plugins；有关如何创建 plugins 的详细信息，请参阅[创建 plugins](/zh-CN/plugins)。
2. **创建 marketplace 文件**：定义一个 `marketplace.json`，列出你的 plugins 及其位置（请参阅[创建 marketplace 文件](#create-the-marketplace-file)）。
3. **托管 marketplace**：推送到 GitHub、GitLab 或其他 git 主机（请参阅[托管和分发 marketplaces](#host-and-distribute-marketplaces)）。
4. **与用户共享**：用户使用 `/plugin marketplace add` 添加你的 marketplace 并安装单个 plugins（请参阅[发现和安装 plugins](/zh-CN/discover-plugins)）。

一旦你的 marketplace 上线，你可以通过推送更改到你的存储库来更新它。用户使用 `/plugin marketplace update` 刷新他们的本地副本。

## 演练：创建本地 marketplace

此示例创建一个包含一个 plugin 的 marketplace：一个用于代码审查的 `/quality-review` skill。你将创建目录结构、添加 skill、创建 plugin manifest 和 marketplace 目录，然后安装并测试它。

<Steps>
  <Step title="创建目录结构">
    ```bash  theme={null}
    mkdir -p my-marketplace/.claude-plugin
    mkdir -p my-marketplace/plugins/quality-review-plugin/.claude-plugin
    mkdir -p my-marketplace/plugins/quality-review-plugin/skills/quality-review
    ```
  </Step>

  <Step title="创建 skill">
    创建一个 `SKILL.md` 文件，定义 `/quality-review` skill 的功能。

    ```markdown my-marketplace/plugins/quality-review-plugin/skills/quality-review/SKILL.md theme={null}
    ---
    description: Review code for bugs, security, and performance
    disable-model-invocation: true
    ---

    Review the code I've selected or the recent changes for:
    - Potential bugs or edge cases
    - Security concerns
    - Performance issues
    - Readability improvements

    Be concise and actionable.
    ```
  </Step>

  <Step title="创建 plugin manifest">
    创建一个 `plugin.json` 文件，描述该 plugin。manifest 位于 `.claude-plugin/` 目录中。

    ```json my-marketplace/plugins/quality-review-plugin/.claude-plugin/plugin.json theme={null}
    {
      "name": "quality-review-plugin",
      "description": "Adds a /quality-review skill for quick code reviews",
      "version": "1.0.0"
    }
    ```
  </Step>

  <Step title="创建 marketplace 文件">
    创建列出你的 plugin 的 marketplace 目录。

    ```json my-marketplace/.claude-plugin/marketplace.json theme={null}
    {
      "name": "my-plugins",
      "owner": {
        "name": "Your Name"
      },
      "plugins": [
        {
          "name": "quality-review-plugin",
          "source": "./plugins/quality-review-plugin",
          "description": "Adds a /quality-review skill for quick code reviews"
        }
      ]
    }
    ```
  </Step>

  <Step title="添加和安装">
    添加 marketplace 并安装 plugin。

    ```shell  theme={null}
    /plugin marketplace add ./my-marketplace
    /plugin install quality-review-plugin@my-plugins
    ```
  </Step>

  <Step title="尝试一下">
    在编辑器中选择一些代码并运行你的新命令。

    ```shell  theme={null}
    /quality-review
    ```
  </Step>
</Steps>

要了解更多关于 plugins 可以做什么的信息，包括 hooks、agents、MCP servers 和 LSP servers，请参阅 [Plugins](/zh-CN/plugins)。

<Note>
  **plugins 如何安装**：当用户安装 plugin 时，Claude Code 将 plugin 目录复制到缓存位置。这意味着 plugins 无法使用 `../shared-utils` 之类的路径引用其目录外的文件，因为这些文件不会被复制。

  如果你需要在 plugins 之间共享文件，请使用符号链接（在复制期间会被跟踪）。有关详细信息，请参阅 [Plugin 缓存和文件解析](/zh-CN/plugins-reference#plugin-caching-and-file-resolution)。
</Note>

## 创建 marketplace 文件

在你的存储库根目录中创建 `.claude-plugin/marketplace.json`。此文件定义你的 marketplace 的名称、所有者信息以及包含其源的 plugins 列表。

每个 plugin 条目至少需要一个 `name` 和 `source`（从哪里获取它）。有关所有可用字段，请参阅下面的[完整架构](#marketplace-schema)。

```json  theme={null}
{
  "name": "company-tools",
  "owner": {
    "name": "DevTools Team",
    "email": "devtools@example.com"
  },
  "plugins": [
    {
      "name": "code-formatter",
      "source": "./plugins/formatter",
      "description": "Automatic code formatting on save",
      "version": "2.1.0",
      "author": {
        "name": "DevTools Team"
      }
    },
    {
      "name": "deployment-tools",
      "source": {
        "source": "github",
        "repo": "company/deploy-plugin"
      },
      "description": "Deployment automation tools"
    }
  ]
}
```

## Marketplace 架构

### 必需字段

| 字段        | 类型     | 描述                                                                                                          | 示例             |
| :-------- | :----- | :---------------------------------------------------------------------------------------------------------- | :------------- |
| `name`    | string | Marketplace 标识符（kebab-case，无空格）。这是面向公众的：用户在安装 plugins 时会看到它（例如，`/plugin install my-tool@your-marketplace`）。 | `"acme-tools"` |
| `owner`   | object | Marketplace 维护者信息（[见下面的字段](#owner-fields)）                                                                  |                |
| `plugins` | array  | 可用 plugins 列表                                                                                               | 见下文            |

<Note>
  **保留名称**：以下 marketplace 名称为 Anthropic 官方使用保留，第三方 marketplaces 无法使用：`claude-code-marketplace`、`claude-code-plugins`、`claude-plugins-official`、`anthropic-marketplace`、`anthropic-plugins`、`agent-skills`、`knowledge-work-plugins`、`life-sciences`。冒充官方 marketplaces 的名称（如 `official-claude-plugins` 或 `anthropic-tools-v2`）也被阻止。
</Note>

### 所有者字段

| 字段      | 类型     | 必需 | 描述         |
| :------ | :----- | :- | :--------- |
| `name`  | string | 是  | 维护者或团队的名称  |
| `email` | string | 否  | 维护者的联系电子邮件 |

### 可选元数据

| 字段                     | 类型     | 描述                                                                                                       |
| :--------------------- | :----- | :------------------------------------------------------------------------------------------------------- |
| `metadata.description` | string | 简短的 marketplace 描述                                                                                       |
| `metadata.version`     | string | Marketplace 版本                                                                                           |
| `metadata.pluginRoot`  | string | 前置到相对 plugin 源路径的基目录（例如，`"./plugins"` 让你写 `"source": "formatter"` 而不是 `"source": "./plugins/formatter"`） |

## Plugin 条目

`plugins` 数组中的每个 plugin 条目描述一个 plugin 及其位置。你可以包含 [plugin manifest 架构](/zh-CN/plugins-reference#plugin-manifest-schema)中的任何字段（如 `description`、`version`、`author`、`commands`、`hooks` 等），加上这些 marketplace 特定的字段：`source`、`category`、`tags` 和 `strict`。

### 必需字段

| 字段       | 类型             | 描述                                                                                         |
| :------- | :------------- | :----------------------------------------------------------------------------------------- |
| `name`   | string         | Plugin 标识符（kebab-case，无空格）。这是面向公众的：用户在安装时会看到它（例如，`/plugin install my-plugin@marketplace`）。 |
| `source` | string\|object | 从哪里获取 plugin（见下面的 [Plugin 源](#plugin-sources)）                                             |

### 可选 plugin 字段

**标准元数据字段：**

| 字段            | 类型      | 描述                                                                   |
| :------------ | :------ | :------------------------------------------------------------------- |
| `description` | string  | 简短的 plugin 描述                                                        |
| `version`     | string  | Plugin 版本                                                            |
| `author`      | object  | Plugin 作者信息（`name` 必需，`email` 可选）                                    |
| `homepage`    | string  | Plugin 主页或文档 URL                                                     |
| `repository`  | string  | 源代码存储库 URL                                                           |
| `license`     | string  | SPDX 许可证标识符（例如，MIT、Apache-2.0）                                       |
| `keywords`    | array   | 用于 plugin 发现和分类的标签                                                   |
| `category`    | string  | Plugin 类别以供组织                                                        |
| `tags`        | array   | 用于可搜索性的标签                                                            |
| `strict`      | boolean | 控制 `plugin.json` 是否是组件定义的权威（默认：true）。见下面的 [Strict 模式](#strict-mode)。 |

**组件配置字段：**

| 字段           | 类型             | 描述                        |
| :----------- | :------------- | :------------------------ |
| `commands`   | string\|array  | 命令文件或目录的自定义路径             |
| `agents`     | string\|array  | agent 文件的自定义路径            |
| `hooks`      | string\|object | 自定义 hooks 配置或 hooks 文件的路径 |
| `mcpServers` | string\|object | MCP server 配置或 MCP 配置的路径  |
| `lspServers` | string\|object | LSP server 配置或 LSP 配置的路径  |

## Plugin 源

Plugin 源告诉 Claude Code 在你的 marketplace 中列出的每个单独 plugin 从哪里获取。这些在 `marketplace.json` 中每个 plugin 条目的 `source` 字段中设置。

一旦 plugin 被克隆或复制到本地机器，它就会被复制到本地版本化 plugin 缓存中，位置为 `~/.claude/plugins/cache`。

| 源            | 类型                           | 字段                               | 注释                                      |
| ------------ | ---------------------------- | -------------------------------- | --------------------------------------- |
| 相对路径         | `string`（例如 `"./my-plugin"`） | —                                | marketplace repo 中的本地目录。必须以 `./` 开头     |
| `github`     | object                       | `repo`、`ref?`、`sha?`             |                                         |
| `url`        | object                       | `url`、`ref?`、`sha?`              | Git URL 源                               |
| `git-subdir` | object                       | `url`、`path`、`ref?`、`sha?`       | git repo 中的子目录。稀疏克隆以最小化大型 monorepos 的带宽 |
| `npm`        | object                       | `package`、`version?`、`registry?` | 通过 `npm install` 安装                     |

<Note>
  **Marketplace 源与 plugin 源**：这些是控制不同事物的不同概念。

  * **Marketplace 源** — 从哪里获取 `marketplace.json` 目录本身。在用户运行 `/plugin marketplace add` 或在 `extraKnownMarketplaces` 设置中设置。支持 `ref`（分支/标签）但不支持 `sha`。
  * **Plugin 源** — 从哪里获取 marketplace 中列出的单个 plugin。在 `marketplace.json` 内每个 plugin 条目的 `source` 字段中设置。支持 `ref`（分支/标签）和 `sha`（精确提交）。

  例如，托管在 `acme-corp/plugin-catalog` 的 marketplace（marketplace 源）可以列出从 `acme-corp/code-formatter` 获取的 plugin（plugin 源）。marketplace 源和 plugin 源指向不同的存储库，并独立固定。
</Note>

### 相对路径

对于同一存储库中的 plugins，使用以 `./` 开头的路径：

```json  theme={null}
{
  "name": "my-plugin",
  "source": "./plugins/my-plugin"
}
```

路径相对于 marketplace 根目录解析，即包含 `.claude-plugin/` 的目录。在上面的示例中，`./plugins/my-plugin` 指向 `<repo>/plugins/my-plugin`，即使 `marketplace.json` 位于 `<repo>/.claude-plugin/marketplace.json`。不要使用 `../` 来爬出 `.claude-plugin/`。

<Note>
  相对路径仅在用户通过 Git（GitHub、GitLab 或 git URL）添加你的 marketplace 时有效。如果用户通过直接 URL 添加你的 marketplace 到 `marketplace.json` 文件，相对路径将无法正确解析。对于基于 URL 的分发，请改用 GitHub、npm 或 git URL 源。有关详细信息，请参阅[故障排除](#plugins-with-relative-paths-fail-in-url-based-marketplaces)。
</Note>

### GitHub 存储库

```json  theme={null}
{
  "name": "github-plugin",
  "source": {
    "source": "github",
    "repo": "owner/plugin-repo"
  }
}
```

你可以固定到特定的分支、标签或提交：

```json  theme={null}
{
  "name": "github-plugin",
  "source": {
    "source": "github",
    "repo": "owner/plugin-repo",
    "ref": "v2.0.0",
    "sha": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
  }
}
```

| 字段     | 类型     | 描述                               |
| :----- | :----- | :------------------------------- |
| `repo` | string | 必需。`owner/repo` 格式的 GitHub 存储库   |
| `ref`  | string | 可选。Git 分支或标签（默认为存储库默认分支）         |
| `sha`  | string | 可选。完整的 40 字符 git 提交 SHA 以固定到精确版本 |

### Git 存储库

```json  theme={null}
{
  "name": "git-plugin",
  "source": {
    "source": "url",
    "url": "https://gitlab.com/team/plugin.git"
  }
}
```

你可以固定到特定的分支、标签或提交：

```json  theme={null}
{
  "name": "git-plugin",
  "source": {
    "source": "url",
    "url": "https://gitlab.com/team/plugin.git",
    "ref": "main",
    "sha": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
  }
}
```

| 字段    | 类型     | 描述                                                                                                   |
| :---- | :----- | :--------------------------------------------------------------------------------------------------- |
| `url` | string | 必需。完整的 git 存储库 URL（`https://` 或 `git@`）。`.git` 后缀是可选的，所以 Azure DevOps 和 AWS CodeCommit URL 不带后缀也可以工作 |
| `ref` | string | 可选。Git 分支或标签（默认为存储库默认分支）                                                                             |
| `sha` | string | 可选。完整的 40 字符 git 提交 SHA 以固定到精确版本                                                                     |

### Git 子目录

使用 `git-subdir` 指向位于 git 存储库子目录中的 plugin。Claude Code 使用稀疏的部分克隆来仅获取子目录，最小化大型 monorepos 的带宽。

```json  theme={null}
{
  "name": "my-plugin",
  "source": {
    "source": "git-subdir",
    "url": "https://github.com/acme-corp/monorepo.git",
    "path": "tools/claude-plugin"
  }
}
```

你可以固定到特定的分支、标签或提交：

```json  theme={null}
{
  "name": "my-plugin",
  "source": {
    "source": "git-subdir",
    "url": "https://github.com/acme-corp/monorepo.git",
    "path": "tools/claude-plugin",
    "ref": "v2.0.0",
    "sha": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
  }
}
```

`url` 字段也接受 GitHub 简写（`owner/repo`）或 SSH URL（`git@github.com:owner/repo.git`）。

| 字段     | 类型     | 描述                                                    |
| :----- | :----- | :---------------------------------------------------- |
| `url`  | string | 必需。Git 存储库 URL、GitHub `owner/repo` 简写或 SSH URL        |
| `path` | string | 必需。repo 中包含 plugin 的子目录路径（例如，`"tools/claude-plugin"`） |
| `ref`  | string | 可选。Git 分支或标签（默认为存储库默认分支）                              |
| `sha`  | string | 可选。完整的 40 字符 git 提交 SHA 以固定到精确版本                      |

### npm 包

作为 npm 包分发的 Plugins 使用 `npm install` 安装。这适用于公共 npm registry 上的任何包或你的团队托管的私有 registry。

```json  theme={null}
{
  "name": "my-npm-plugin",
  "source": {
    "source": "npm",
    "package": "@acme/claude-plugin"
  }
}
```

要固定到特定版本，请添加 `version` 字段：

```json  theme={null}
{
  "name": "my-npm-plugin",
  "source": {
    "source": "npm",
    "package": "@acme/claude-plugin",
    "version": "2.1.0"
  }
}
```

要从私有或内部 registry 安装，请添加 `registry` 字段：

```json  theme={null}
{
  "name": "my-npm-plugin",
  "source": {
    "source": "npm",
    "package": "@acme/claude-plugin",
    "version": "^2.0.0",
    "registry": "https://npm.example.com"
  }
}
```

| 字段         | 类型     | 描述                                                        |
| :--------- | :----- | :-------------------------------------------------------- |
| `package`  | string | 必需。包名称或作用域包（例如，`@org/plugin`）                             |
| `version`  | string | 可选。版本或版本范围（例如，`2.1.0`、`^2.0.0`、`~1.5.0`）                  |
| `registry` | string | 可选。自定义 npm registry URL。默认为系统 npm registry（通常为 npmjs.org） |

### 高级 plugin 条目

此示例显示了使用许多可选字段的 plugin 条目，包括命令、agents、hooks 和 MCP servers 的自定义路径：

```json  theme={null}
{
  "name": "enterprise-tools",
  "source": {
    "source": "github",
    "repo": "company/enterprise-plugin"
  },
  "description": "Enterprise workflow automation tools",
  "version": "2.1.0",
  "author": {
    "name": "Enterprise Team",
    "email": "enterprise@example.com"
  },
  "homepage": "https://docs.example.com/plugins/enterprise-tools",
  "repository": "https://github.com/company/enterprise-plugin",
  "license": "MIT",
  "keywords": ["enterprise", "workflow", "automation"],
  "category": "productivity",
  "commands": [
    "./commands/core/",
    "./commands/enterprise/",
    "./commands/experimental/preview.md"
  ],
  "agents": ["./agents/security-reviewer.md", "./agents/compliance-checker.md"],
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/validate.sh"
          }
        ]
      }
    ]
  },
  "mcpServers": {
    "enterprise-db": {
      "command": "${CLAUDE_PLUGIN_ROOT}/servers/db-server",
      "args": ["--config", "${CLAUDE_PLUGIN_ROOT}/config.json"]
    }
  },
  "strict": false
}
```

需要注意的关键事项：

* **`commands` 和 `agents`**：你可以指定多个目录或单个文件。路径相对于 plugin 根目录。
* **`${CLAUDE_PLUGIN_ROOT}`**：在 hooks 和 MCP server 配置中使用此变量来引用 plugin 安装目录中的文件。这是必要的，因为 plugins 在安装时被复制到缓存位置。对于应该在 plugin 更新后保留的依赖项或状态，请改用 [`${CLAUDE_PLUGIN_DATA}`](/zh-CN/plugins-reference#persistent-data-directory)。
* **`strict: false`**：由于这设置为 false，plugin 不需要自己的 `plugin.json`。marketplace 条目定义了一切。见下面的 [Strict 模式](#strict-mode)。

### Strict 模式

`strict` 字段控制 `plugin.json` 是否是组件定义（命令、agents、hooks、skills、MCP servers、输出样式）的权威。

| 值          | 行为                                                                      |
| :--------- | :---------------------------------------------------------------------- |
| `true`（默认） | `plugin.json` 是权威。marketplace 条目可以用额外的组件补充它，两个源都被合并。                    |
| `false`    | marketplace 条目是完整的定义。如果 plugin 也有声明组件的 `plugin.json`，那就是冲突，plugin 无法加载。 |

**何时使用每种模式：**

* **`strict: true`**：plugin 有自己的 `plugin.json` 并管理自己的组件。marketplace 条目可以在顶部添加额外的命令或 hooks。这是默认值，适用于大多数 plugins。
* **`strict: false`**：marketplace 操作员想要完全控制。plugin repo 提供原始文件，marketplace 条目定义这些文件中的哪些被公开为命令、agents、hooks 等。当 marketplace 以不同于 plugin 作者意图的方式重组或策划 plugin 的组件时很有用。

## 托管和分发 marketplaces

### 在 GitHub 上托管（推荐）

GitHub 提供最简单的分发方法：

1. **创建存储库**：为你的 marketplace 设置一个新存储库
2. **添加 marketplace 文件**：使用你的 plugin 定义创建 `.claude-plugin/marketplace.json`
3. **与团队共享**：用户使用 `/plugin marketplace add owner/repo` 添加你的 marketplace

**优点**：内置版本控制、问题跟踪和团队协作功能。

### 在其他 git 服务上托管

任何 git 托管服务都可以工作，例如 GitLab、Bitbucket 和自托管服务器。用户使用完整的存储库 URL 添加：

```shell  theme={null}
/plugin marketplace add https://gitlab.com/company/plugins.git
```

### 私有存储库

Claude Code 支持从私有存储库安装 plugins。对于手动安装和更新，Claude Code 使用你现有的 git 凭证助手。如果 `git clone` 对你终端中的私有存储库有效，它在 Claude Code 中也有效。常见的凭证助手包括用于 GitHub 的 `gh auth login`、macOS Keychain 和 `git-credential-store`。

后台自动更新在启动时运行，不使用凭证助手，因为交互式提示会阻止 Claude Code 启动。要为私有 marketplaces 启用自动更新，请在你的环境中设置适当的身份验证令牌：

| 提供商       | 环境变量                        | 注释                    |
| :-------- | :-------------------------- | :-------------------- |
| GitHub    | `GITHUB_TOKEN` 或 `GH_TOKEN` | 个人访问令牌或 GitHub App 令牌 |
| GitLab    | `GITLAB_TOKEN` 或 `GL_TOKEN` | 个人访问令牌或项目令牌           |
| Bitbucket | `BITBUCKET_TOKEN`           | 应用密码或存储库访问令牌          |

在你的 shell 配置中设置令牌（例如，`.bashrc`、`.zshrc`）或在运行 Claude Code 时传递它：

```bash  theme={null}
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
```

<Note>
  对于 CI/CD 环境，将令牌配置为秘密环境变量。GitHub Actions 自动为同一组织中的存储库提供 `GITHUB_TOKEN`。
</Note>

### 在分发前本地测试

在共享前本地测试你的 marketplace：

```shell  theme={null}
/plugin marketplace add ./my-local-marketplace
/plugin install test-plugin@my-local-marketplace
```

有关完整的添加命令范围（GitHub、Git URL、本地路径、远程 URL），请参阅[添加 marketplaces](/zh-CN/discover-plugins#add-marketplaces)。

### 为你的团队要求 marketplaces

你可以配置你的存储库，以便当团队成员信任项目文件夹时，他们会自动被提示安装你的 marketplace。将你的 marketplace 添加到 `.claude/settings.json`：

```json  theme={null}
{
  "extraKnownMarketplaces": {
    "company-tools": {
      "source": {
        "source": "github",
        "repo": "your-org/claude-plugins"
      }
    }
  }
}
```

你也可以指定默认应启用哪些 plugins：

```json  theme={null}
{
  "enabledPlugins": {
    "code-formatter@company-tools": true,
    "deployment-tools@company-tools": true
  }
}
```

有关完整的配置选项，请参阅 [Plugin 设置](/zh-CN/settings#plugin-settings)。

<Note>
  如果你使用带有相对路径的本地 `directory` 或 `file` 源，路径将相对于你的存储库的主检出解析。当你从 git worktree 运行 Claude Code 时，路径仍然指向主检出，所以所有 worktrees 共享相同的 marketplace 位置。Marketplace 状态存储一次每个用户在 `~/.claude/plugins/known_marketplaces.json` 中，而不是每个项目。
</Note>

### 为容器预填充 plugins

对于容器镜像和 CI 环境，你可以在构建时预填充 plugins 目录，以便 Claude Code 启动时已经有 marketplaces 和 plugins 可用，无需在运行时克隆任何内容。设置 `CLAUDE_CODE_PLUGIN_SEED_DIR` 环境变量以指向此目录。

要分层多个种子目录，请在 Unix 上用 `:` 分隔路径，或在 Windows 上用 `;` 分隔。Claude Code 按顺序搜索每个目录，第一个包含给定 marketplace 或 plugin 缓存的种子获胜。

种子目录镜像 `~/.claude/plugins` 的结构：

```
$CLAUDE_CODE_PLUGIN_SEED_DIR/
  known_marketplaces.json
  marketplaces/<name>/...
  cache/<marketplace>/<plugin>/<version>/...
```

构建种子目录的最简单方法是在镜像构建期间运行 Claude Code 一次，安装你需要的 plugins，然后将生成的 `~/.claude/plugins` 目录复制到你的镜像中，并将 `CLAUDE_CODE_PLUGIN_SEED_DIR` 指向它。

在启动时，Claude Code 将种子的 `known_marketplaces.json` 中找到的 marketplaces 注册到主配置中，并使用在 `cache/` 下找到的 plugin 缓存，而无需重新克隆。这在交互模式和使用 `-p` 标志的非交互模式中都有效。

行为详情：

* **只读**：种子目录永远不会被写入。由于 git pull 会在只读文件系统上失败，种子 marketplaces 的自动更新被禁用。
* **种子条目优先**：在每次启动时，种子中声明的 marketplaces 会覆盖用户配置中的任何匹配条目。要选择退出种子 plugin，请使用 `/plugin disable` 而不是删除 marketplace。
* **路径解析**：Claude Code 通过在运行时探测 `$CLAUDE_CODE_PLUGIN_SEED_DIR/marketplaces/<name>/` 来定位 marketplace 内容，而不是信任存储在种子 JSON 内的路径。这意味着即使在与构建时不同的路径上挂载，种子也能正确工作。
* **与设置组合**：如果 `extraKnownMarketplaces` 或 `enabledPlugins` 声明的 marketplace 已经存在于种子中，Claude Code 使用种子副本而不是克隆。

### 托管 marketplace 限制

对于需要严格控制 plugin 源的组织，管理员可以使用托管设置中的 [`strictKnownMarketplaces`](/zh-CN/settings#strictknownmarketplaces) 设置限制用户允许添加哪些 plugin marketplaces。

当在托管设置中配置 `strictKnownMarketplaces` 时，限制行为取决于值：

| 值        | 行为                            |
| -------- | ----------------------------- |
| 未定义（默认）  | 无限制。用户可以添加任何 marketplace      |
| 空数组 `[]` | 完全锁定。用户无法添加任何新 marketplaces   |
| 源列表      | 用户只能添加与允许列表完全匹配的 marketplaces |

#### 常见配置

禁用所有 marketplace 添加：

```json  theme={null}
{
  "strictKnownMarketplaces": []
}
```

仅允许特定 marketplaces：

```json  theme={null}
{
  "strictKnownMarketplaces": [
    {
      "source": "github",
      "repo": "acme-corp/approved-plugins"
    },
    {
      "source": "github",
      "repo": "acme-corp/security-tools",
      "ref": "v2.0"
    },
    {
      "source": "url",
      "url": "https://plugins.example.com/marketplace.json"
    }
  ]
}
```

使用主机上的正则表达式模式匹配允许来自内部 git 服务器的所有 marketplaces。这是 [GitHub Enterprise Server](/zh-CN/github-enterprise-server#plugin-marketplaces-on-ghes) 或自托管 GitLab 实例的推荐方法：

```json  theme={null}
{
  "strictKnownMarketplaces": [
    {
      "source": "hostPattern",
      "hostPattern": "^github\\.example\\.com$"
    }
  ]
}
```

使用路径上的正则表达式模式匹配允许来自特定目录的基于文件系统的 marketplaces：

```json  theme={null}
{
  "strictKnownMarketplaces": [
    {
      "source": "pathPattern",
      "pathPattern": "^/opt/approved/"
    }
  ]
}
```

使用 `".*"` 作为 `pathPattern` 来允许任何文件系统路径，同时仍然使用 `hostPattern` 控制网络源。

<Note>
  `strictKnownMarketplaces` 限制用户可以添加的内容，但不会自行注册 marketplaces。要使允许的 marketplaces 自动可用而无需用户运行 `/plugin marketplace add`，请在同一 `managed-settings.json` 中将其与 [`extraKnownMarketplaces`](/zh-CN/settings#extraknownmarketplaces) 配对。见[同时使用两者](/zh-CN/settings#strictknownmarketplaces)。
</Note>

#### 限制如何工作

限制在 plugin 安装过程的早期进行验证，在任何网络请求或文件系统操作之前。这可以防止未授权的 marketplace 访问尝试。

允许列表对大多数源类型使用精确匹配。要允许 marketplace，所有指定的字段必须完全匹配：

* 对于 GitHub 源：`repo` 是必需的，如果在允许列表中指定，`ref` 或 `path` 也必须匹配
* 对于 URL 源：完整 URL 必须完全匹配
* 对于 `hostPattern` 源：marketplace 主机与正则表达式模式匹配
* 对于 `pathPattern` 源：marketplace 的文件系统路径与正则表达式模式匹配

因为 `strictKnownMarketplaces` 在[托管设置](/zh-CN/settings#settings-files)中设置，个别用户和项目配置无法覆盖这些限制。

有关完整的配置详细信息，包括所有支持的源类型和与 `extraKnownMarketplaces` 的比较，请参阅 [strictKnownMarketplaces 参考](/zh-CN/settings#strictknownmarketplaces)。

### 版本解析和发布渠道

Plugin 版本确定缓存路径和更新检测。你可以在 plugin manifest（`plugin.json`）或 marketplace 条目（`marketplace.json`）中指定版本。

<Warning>
  如果可能，避免在两个地方都设置版本。plugin manifest 总是无声地获胜，这可能导致 marketplace 版本被忽略。对于相对路径 plugins，在 marketplace 条目中设置版本。对于所有其他 plugin 源，在 plugin manifest 中设置它。
</Warning>

#### 设置发布渠道

要为你的 plugins 支持"稳定"和"最新"发布渠道，你可以设置两个指向同一 repo 的不同 refs 或 SHAs 的 marketplaces。然后，你可以通过[托管设置](/zh-CN/settings#settings-files)将两个 marketplaces 分配给不同的用户组。

<Warning>
  plugin 的 `plugin.json` 必须在每个固定的 ref 或提交处声明不同的 `version`。如果两个 refs 或提交具有相同的 manifest 版本，Claude Code 会将它们视为相同并跳过更新。
</Warning>

##### 示例

```json  theme={null}
{
  "name": "stable-tools",
  "plugins": [
    {
      "name": "code-formatter",
      "source": {
        "source": "github",
        "repo": "acme-corp/code-formatter",
        "ref": "stable"
      }
    }
  ]
}
```

```json  theme={null}
{
  "name": "latest-tools",
  "plugins": [
    {
      "name": "code-formatter",
      "source": {
        "source": "github",
        "repo": "acme-corp/code-formatter",
        "ref": "latest"
      }
    }
  ]
}
```

##### 将渠道分配给用户组

通过托管设置将每个 marketplace 分配给适当的用户组。例如，稳定组接收：

```json  theme={null}
{
  "extraKnownMarketplaces": {
    "stable-tools": {
      "source": {
        "source": "github",
        "repo": "acme-corp/stable-tools"
      }
    }
  }
}
```

早期访问组改为接收 `latest-tools`：

```json  theme={null}
{
  "extraKnownMarketplaces": {
    "latest-tools": {
      "source": {
        "source": "github",
        "repo": "acme-corp/latest-tools"
      }
    }
  }
}
```

## 验证和测试

在共享前测试你的 marketplace。

验证你的 marketplace JSON 语法：

```bash  theme={null}
claude plugin validate .
```

或从 Claude Code 内：

```shell  theme={null}
/plugin validate .
```

添加 marketplace 进行测试：

```shell  theme={null}
/plugin marketplace add ./path/to/marketplace
```

安装测试 plugin 以验证一切正常：

```shell  theme={null}
/plugin install test-plugin@marketplace-name
```

有关完整的 plugin 测试工作流，请参阅[本地测试你的 plugins](/zh-CN/plugins#test-your-plugins-locally)。有关技术故障排除，请参阅 [Plugins 参考](/zh-CN/plugins-reference)。

## 故障排除

### Marketplace 未加载

**症状**：无法添加 marketplace 或从中看到 plugins

**解决方案**：

* 验证 marketplace URL 是否可访问
* 检查 `.claude-plugin/marketplace.json` 是否存在于指定路径
* 使用 `claude plugin validate` 或 `/plugin validate` 确保 JSON 语法有效且 frontmatter 格式正确
* 对于私有存储库，确认你有访问权限

### Marketplace 验证错误

从你的 marketplace 目录运行 `claude plugin validate .` 或 `/plugin validate .` 来检查问题。验证器检查 `plugin.json`、skill/agent/command frontmatter 和 `hooks/hooks.json` 的语法和架构错误。常见错误：

| 错误                                                | 原因                                 | 解决方案                                                       |
| :------------------------------------------------ | :--------------------------------- | :--------------------------------------------------------- |
| `File not found: .claude-plugin/marketplace.json` | 缺少 manifest                        | 使用必需字段创建 `.claude-plugin/marketplace.json`                 |
| `Invalid JSON syntax: Unexpected token...`        | JSON 语法错误                          | 检查缺少的逗号、多余的逗号或未引用的字符串                                      |
| `Duplicate plugin name "x" found in marketplace`  | 两个 plugins 共享相同的名称                 | 给每个 plugin 一个唯一的 `name` 值                                  |
| `plugins[0].source: Path contains ".."`           | 源路径包含 `..`                         | 使用相对于 marketplace 根目录的路径，不包含 `..`。见[相对路径](#relative-paths) |
| `YAML frontmatter failed to parse: ...`           | skill、agent 或 command 文件中的 YAML 无效 | 修复 frontmatter 块中的 YAML 语法。在运行时，此文件加载时不带元数据。               |
| `Invalid JSON syntax: ...`（hooks.json）            | 格式错误的 `hooks/hooks.json`           | 修复 JSON 语法。格式错误的 `hooks/hooks.json` 会阻止整个 plugin 加载。       |

**警告**（非阻止）：

* `Marketplace has no plugins defined`：将至少一个 plugin 添加到 `plugins` 数组
* `No marketplace description provided`：添加 `metadata.description` 以帮助用户理解你的 marketplace
* `Plugin name "x" is not kebab-case`：plugin 名称包含大写字母、空格或特殊字符。重命名为仅包含小写字母、数字和连字符（例如，`my-plugin`）。Claude Code 接受其他形式，但 Claude.ai marketplace 同步会拒绝它们。

### Plugin 安装失败

**症状**：Marketplace 出现但 plugin 安装失败

**解决方案**：

* 验证 plugin 源 URL 是否可访问
* 检查 plugin 目录是否包含必需的文件
* 对于 GitHub 源，确保存储库是公开的或你有访问权限
* 通过手动克隆/下载来测试 plugin 源

### 私有存储库身份验证失败

**症状**：从私有存储库安装 plugins 时出现身份验证错误

**解决方案**：

对于手动安装和更新：

* 验证你已使用你的 git 提供商进行身份验证（例如，对于 GitHub 运行 `gh auth status`）
* 检查你的凭证助手是否配置正确：`git config --global credential.helper`
* 尝试手动克隆存储库以验证你的凭证有效

对于后台自动更新：

* 在你的环境中设置适当的令牌：`echo $GITHUB_TOKEN`
* 检查令牌是否具有所需的权限（对存储库的读取访问权限）
* 对于 GitHub，确保令牌对私有存储库具有 `repo` 范围
* 对于 GitLab，确保令牌至少具有 `read_repository` 范围
* 验证令牌未过期

### Marketplace 更新在离线环境中失败

**症状**：Marketplace `git pull` 失败，Claude Code 清除现有缓存，导致 plugins 变得不可用。

**原因**：默认情况下，当 `git pull` 失败时，Claude Code 会删除陈旧的克隆并尝试重新克隆。在离线或隔离的环境中，重新克隆以相同的方式失败，导致 marketplace 目录为空。

**解决方案**：设置 `CLAUDE_CODE_PLUGIN_KEEP_MARKETPLACE_ON_FAILURE=1` 以在拉取失败时保留现有缓存，而不是清除它：

```bash  theme={null}
export CLAUDE_CODE_PLUGIN_KEEP_MARKETPLACE_ON_FAILURE=1
```

设置此变量后，Claude Code 在 `git pull` 失败时保留陈旧的 marketplace 克隆，并继续使用最后已知的良好状态。对于存储库永远无法访问的完全离线部署，请改用 [`CLAUDE_CODE_PLUGIN_SEED_DIR`](#pre-populate-plugins-for-containers) 在构建时预填充 plugins 目录。

### Git 操作超时

**症状**：Plugin 安装或 marketplace 更新失败，出现超时错误，如"Git clone timed out after 120s"或"Git pull timed out after 120s"。

**原因**：Claude Code 对所有 git 操作使用 120 秒超时，包括克隆 plugin 存储库和拉取 marketplace 更新。大型存储库或缓慢的网络连接可能超过此限制。

**解决方案**：使用 `CLAUDE_CODE_PLUGIN_GIT_TIMEOUT_MS` 环境变量增加超时。该值以毫秒为单位：

```bash  theme={null}
export CLAUDE_CODE_PLUGIN_GIT_TIMEOUT_MS=300000  # 5 minutes
```

### 相对路径 Plugins 在基于 URL 的 marketplaces 中失败

**症状**：通过 URL（如 `https://example.com/marketplace.json`）添加了 marketplace，但具有相对路径源（如 `"./plugins/my-plugin"`）的 plugins 无法安装，出现"path not found"错误。

**原因**：基于 URL 的 marketplaces 仅下载 `marketplace.json` 文件本身。它们不从服务器下载 plugin 文件。marketplace 条目中的相对路径引用远程服务器上未下载的文件。

**解决方案**：

* **使用外部源**：将 plugin 条目更改为使用 GitHub、npm 或 git URL 源而不是相对路径：
  ```json  theme={null}
  { "name": "my-plugin", "source": { "source": "github", "repo": "owner/repo" } }
  ```
* **使用基于 Git 的 marketplace**：在 Git 存储库中托管你的 marketplace 并使用 git URL 添加它。基于 Git 的 marketplaces 克隆整个存储库，使相对路径有效。

### 安装后文件未找到

**症状**：Plugin 安装但对文件的引用失败，特别是 plugin 目录外的文件

**原因**：Plugins 被复制到缓存目录而不是就地使用。引用 plugin 目录外文件的路径（如 `../shared-utils`）不会工作，因为这些文件不会被复制。

**解决方案**：见 [Plugin 缓存和文件解析](/zh-CN/plugins-reference#plugin-caching-and-file-resolution) 了解解决方法，包括符号链接和目录重组。

有关其他调试工具和常见问题，请参阅[调试和开发工具](/zh-CN/plugins-reference#debugging-and-development-tools)。

## 另见

* [发现和安装预构建的 plugins](/zh-CN/discover-plugins) - 从现有 marketplaces 安装 plugins
* [Plugins](/zh-CN/plugins) - 创建你自己的 plugins
* [Plugins 参考](/zh-CN/plugins-reference) - 完整的技术规范和架构
* [Plugin 设置](/zh-CN/settings#plugin-settings) - Plugin 配置选项
* [strictKnownMarketplaces 参考](/zh-CN/settings#strictknownmarketplaces) - 托管 marketplace 限制
