---
created: "2026-05-28"
topic: "MCP Browser Automation Tools"
mode: "multi-candidate-comparison"
candidates: [chrome-devtools-mcp, playwright-mcp, browserbase-mcp, puppeteer-mcp]
dimensions: [overview, architecture, capabilities, ecosystem, performance, developer-experience, learning-curve, security, migration-cost]
---

# MCP 浏览器自动化工具对比调研报告

## Overview

**Chrome DevTools MCP 和 Playwright MCP 是当前 AI agent 浏览器自动化的两大官方方案，互补而非竞争——Playwright 擅长"驱动"（自动化、测试），Chrome DevTools 擅长"调试"（性能分析、内存诊断）。Puppeteer MCP 已弃用，Browserbase MCP 面向云端场景。**

**Research mode:** 多候选对比

**Key question:** 各主流 MCP 浏览器自动化工具的技术细节、架构差异和适用场景是什么？

## Research Background & Objectives

MCP (Model Context Protocol) 已成为连接 AI agent 与浏览器的事实标准。Google (Chrome DevTools MCP) 和 Microsoft (Playwright MCP) 分别在 2025 年发布了官方 MCP 浏览器工具，社区也在快速演进。本报告旨在深入理解各工具的技术架构、能力边界和实际取舍，帮助技术选型。

### Research Scope

| Dimension | Value |
|---|---|
| Topic | MCP 浏览器自动化工具对比 |
| Mode | 多候选对比 |
| Dimensions covered | 概述定位、架构能力、生态社区、性能、开发者体验、学习曲线、安全性、迁移成本 |
| Candidates | chrome-devtools-mcp, playwright-mcp, browserbase-mcp, puppeteer-mcp |
| Project adaptation | 否（纯技术研究） |

---

## Candidate: Chrome DevTools MCP

### 概述与定位

Google Chrome DevTools 团队官方出品的 MCP 服务器，2025 年 9 月发布公开预览版。定位为"面向 coding agent 的 Chrome DevTools"——将浏览器调试的完整能力暴露给 AI agent。

- **创建者**: Google LLC (ChromeDevTools 组织)
- **许可证**: Apache-2.0
- **当前版本**: v1.1.1 (2026-05-27)
- **GitHub Stars**: 42,053
- **npm 周下载量**: ~230 万

### 架构与核心概念

**协议栈**: MCP (stdio) → Puppeteer → CDP (Chrome DevTools Protocol) → Chrome

**关键依赖**:
- `puppeteer` v25.0.4 — 浏览器自动化核心
- `@modelcontextprotocol/sdk` v1.29.0 — MCP 协议实现
- `lighthouse` v13.3.0 — 性能/可访问性审计
- `chrome-devtools-frontend` — DevTools 前端组件

**核心抽象**: `McpPage`（页面封装）、`McpContext`（上下文管理）、`ToolHandler`（工具调用处理）、`HeapSnapshotManager`（堆快照管理）、`TextSnapshot`（基于 a11y 树的文本快照）

**设计原则**:
1. **Agent-Agnostic** — 不锁定特定 LLM
2. **Token-Optimized** — 返回语义摘要，重型资源走文件路径
3. **Small, Deterministic Blocks** — 可组合的小工具而非魔法按钮
4. **Self-Healing Errors** — 错误消息包含修复建议
5. **Progressive Complexity** — 默认简单，高级参数可选

**连接模式**: 自动启动（默认）→ autoConnect (Chrome 144+) → browser-url → wsEndpoint

### 能力清单

**完整模式**: 45 个工具，10 个类别

| 类别 | 工具数 | 代表工具 |
|---|---|---|
| 输入自动化 | 10 | click, fill, type_text, hover, drag, press_key |
| 导航自动化 | 6 | navigate_page, new_page, close_page, list_pages |
| 模拟 | 2 | emulate (CPU/网络/地理/视口/暗色模式), resize_page |
| 性能分析 | 3 | performance_start_trace, performance_stop_trace, performance_analyze_insight |
| 网络 | 2 | list_network_requests, get_network_request |
| 调试 | 8 | evaluate_script, take_snapshot, take_screenshot, lighthouse_audit, screencast |
| 内存 | 5 | take_heapsnapshot, get_heapsnapshot_summary, get_heapsnapshot_retainers |
| 扩展 | 5 | install_extension, list_extensions, trigger_extension_action |
| 第三方 | 2 | execute_3p_developer_tool, list_3p_developer_tools |
| WebMCP | 2 | execute_webmcp_tool, list_webmcp_tools |

**Slim 模式**: 3 个工具 (navigate, evaluate, screenshot) — 适合基础浏览器任务，极大减少 token 消耗

**独特能力**: V8 堆快照与内存泄漏检测、Core Web Vitals (LCP/INP/CLS) 追踪、CrUX 字段数据集成、Lighthouse 审计、浏览器扩展管理、视频录制 (screencast)

### 生态与社区

- 月 npm 下载量 ~800 万，极高活跃度
- 8 个月内从 v0 发展到 v1.1.1
- 主要贡献者来自 Google Puppeteer 团队 (OrKoN — 315 commits)
- 提供 6 个官方 Skill (a11y-debugging, debug-optimize-lcp, memory-leak-debugging 等)
- 支持 20+ MCP 客户端 (Claude Code, VS Code Copilot, Cursor, Gemini CLI 等)

### 安全

- 浏览器内容完全暴露给 MCP client — 官方明确警告避免打开敏感页面
- 默认收集使用统计（可 `--no-usage-statistics` 退出）
- 性能工具会调用 Google CrUX API（可 `--no-performance-crux` 禁用）
- 远程调试端口需非默认 user-data-dir
- 无已知 CVE

### 开发者体验

- **零安装**: `npx chrome-devtools-mcp@latest` 直接运行
- **Claude Code 一行集成**: `claude mcp add chrome-devtools --scope user npx chrome-devtools-mcp@latest`
- **Slim 模式**降低认知负担和 token 消耗
- 工具命名清晰，返回值包含语义摘要
- 官方 troubleshooting 文档覆盖常见平台问题 (WSL, Windows 10, macOS)

### 限制

- 仅支持 Chrome/Chromium（不支持 Firefox/Safari）
- autoConnect 需 Chrome 144+
- 大量标签页时可能有性能问题
- 扩展管理当前仅支持 pipe 连接
- 内存分析、screencast 等为实验性功能
- WSL/沙盒环境有已知兼容性问题

**Sources:**
- [GitHub: ChromeDevTools/chrome-devtools-mcp](https://github.com/ChromeDevTools/chrome-devtools-mcp)
- [Chrome 官方博客公告](https://developer.chrome.com/blog/chrome-devtools-mcp)
- [设计原则文档](https://github.com/ChromeDevTools/chrome-devtools-mcp/blob/main/docs/design-principles.md)
- [工具参考文档](https://github.com/ChromeDevTools/chrome-devtools-mcp/blob/main/docs/tool-reference.md)

---

## Candidate: Playwright MCP

### 概述与定位

Microsoft 官方的 Playwright MCP 服务器，2025 年 3 月发布。定位为浏览器**驱动和自动化**工具——基于 Playwright 的 accessibility tree 进行确定性操作，无需视觉模型。

- **创建者**: Microsoft
- **许可证**: Apache-2.0
- **仓库**: [microsoft/playwright-mcp](https://github.com/microsoft/playwright-mcp)
- **官方文档**: [playwright.dev/docs/getting-started-mcp](https://playwright.dev/docs/getting-started-mcp)

### 架构与核心概念

**协议栈**: MCP (stdio/HTTP) → Playwright → 多浏览器引擎

**核心设计**: 基于 accessibility tree（无障碍树）的结构化操作，LLM 友好，确定性高

**浏览器支持**: Chromium, Firefox, WebKit, MS Edge — 唯一支持跨浏览器测试的 MCP 方案

**三种配置文件**: Persistent（保持登录）、Isolated（每次隔离）、Extension（连接已有浏览器标签页）

**工具数量**: ~21 个工具 (MCP 模式)，40+ 命令 (CLI 模式)

### 能力清单

| 类别 | 工具 |
|---|---|
| 导航 | 打开 URL、前进/后退、刷新 |
| 点击与输入 | 点击、输入文本、填写表单、选择下拉框 |
| 截图 | 整页或指定元素截图 |
| 键盘与鼠标 | 按键、悬停、拖放 |
| 对话框 | 接受或关闭浏览器对话框 |
| 标签页 | 创建、关闭、切换 |
| 网络 | 查看网络请求、模拟 API 响应 (URL 模式 mock) |
| 存储 | 保存/恢复 cookie 和 localStorage |
| 代码执行 | `browser_run_code_unsafe`（直接运行 Playwright 脚本） |

**独特能力**: 跨浏览器测试 (含 WebKit/Safari)、网络 mock/拦截、Playwright Test Agents (planner/generator/healer)、CLI 模式 (token 效率更高)

### 生态与社区

- Microsoft 官方维护，与 Playwright 主仓库协同开发
- 社区讨论热度最高（Reddit r/mcp 和 r/Anthropic）
- Microsoft 官方推荐 CLI 模式用于 coding agent，MCP 用于需要持久状态的场景
- 支持所有主流 MCP 客户端

### 安全

- **CVE-2025-9611**: CSRF 漏洞（已披露）
- `browser_run_code_unsafe` 等效 RCE，需谨慎使用
- 需要本地安装浏览器

### 开发者体验

- **一行安装**: `claude mcp add playwright npx @playwright/mcp@latest`
- 支持 headless、浏览器选择 (`--browser=firefox`)、HTTP 传输 (`--port 8931`)
- CLI 模式更适合 coding agent（Microsoft 官方建议）
- accessibility tree 信息量大，token 消耗是主要痛点

### 限制

- 工具定义消耗 ~14k tokens，用户报告 MCP 工具消耗高达 83.3k tokens
- 超过 10 个工具时 agent 性能急剧下降
- 性能分析能力有限（仅 trace，无 Lighthouse/内存分析）
- 无 Chrome DevTools 级深度调试能力

**Sources:**
- [GitHub: microsoft/playwright-mcp](https://github.com/microsoft/playwright-mcp)
- [Playwright 官方 MCP 文档](https://playwright.dev/docs/getting-started-mcp)
- [Driving vs Debugging the Browser — Steve Kinney](https://stevekinney.com/writing/driving-vs-debugging-the-browser)

---

## Candidate: Browserbase MCP

### 概述与定位

Browserbase 公司（云浏览器基础设施创业公司）出品的 MCP 服务器，基于其开源 AI 浏览器自动化框架 Stagehand 构建。定位为**云端浏览器自动化**——浏览器实例运行在 Browserbase 云端。

- **创建者**: Browserbase
- **许可证**: Apache-2.0
- **仓库**: [browserbase/mcp-server-browserbase](https://github.com/browserbase/mcp-server-browserbase)

### 架构与核心概念

**协议栈**: MCP (stdio/SHTTP) → Stagehand v3 → Playwright/Puppeteer → 云端浏览器

**核心特点**: 浏览器在云端运行，支持 AI 原生操作（自然语言指令驱动）

**双传输模式**: STDIO（本地）和 SHTTP（远程托管，推荐）

**工具数量**: 6 个 (start, end, navigate, act, observe, extract)

### 能力清单

| 工具 | 描述 | 特点 |
|---|---|---|
| `start` | 创建/复用 Browserbase 会话 | 云端浏览器实例 |
| `end` | 关闭会话 | — |
| `navigate` | 导航到 URL | — |
| `act` | 自然语言执行操作 | AI 驱动 |
| `observe` | 观察可操作元素 | AI 驱动 |
| `extract` | 从页面提取数据 | AI 驱动 |

**独特能力**: 云端浏览器（无需本地安装）、反检测 (advancedStealth)、代理轮换、AI 原生操作（自然语言驱动而非确定性 API）

### 生态与社区

- Stagehand v3 (2026-02) 性能提升 20-40%，支持 iframe 和 shadow DOM
- 托管模式 Browserbase 承担 LLM 费用（默认 Gemini）
- 有免费额度（约 60 分钟/月）

### 定价

| 方案 | 费用 | 内容 |
|---|---|---|
| Free | $0/月 | ~60 分钟浏览器时间/月，3 并发 |
| Developer | $20/月 | ~100 浏览器小时 |
| Startup | $99/月 | ~500 浏览器小时 |
| Scale/Enterprise | 定制 | 无限扩展 |

### 限制

- 依赖云服务（免费额度有限）
- 需管理 API Key 和 Project ID
- 工具粒度粗（6 个），不适合精细控制
- AI 操作准确度依赖模型质量

**Sources:**
- [GitHub: browserbase/mcp-server-browserbase](https://github.com/browserbase/mcp-server-browserbase)
- [Browserbase 定价](https://www.browserbase.com/pricing/)
- [Browserbase MCP 文档](https://docs.browserbase.com/integrations/mcp/introduction)

---

## Candidate: Puppeteer MCP

### 概述与定位

MCP 官方组织发布的参考 Puppeteer MCP 服务器。**已弃用且归档** — 2025 年 5 月被归档至只读仓库 [modelcontextprotocol/servers-archived](https://github.com/modelcontextprotocol/servers-archived)。

- **创建者**: @modelcontextprotocol (Anthropic 关联)
- **状态**: 已弃用，不再维护
- **许可证**: MIT

### 架构

MCP (stdio) → Puppeteer → Chromium。支持 Docker (headless) 和 npx (有头浏览器) 两种运行方式。

### 能力清单

7 个工具: navigate, screenshot, click, hover, fill, select, evaluate

2 种资源: `console://logs`（控制台日志）、`screenshot://`（截图图片）

### 结论

**不推荐使用**。已弃用且无安全更新。如果需要 Puppeteer 级能力，应使用 Chrome DevTools MCP（同样基于 Puppeteer/CDP 但活跃维护、功能远超）。

**Sources:**
- [npm: @modelcontextprotocol/server-puppeteer](https://www.npmjs.com/package/@modelcontextprotocol/server-puppeteer)
- [GitHub: modelcontextprotocol/servers-archived](https://github.com/modelcontextprotocol/servers-archived)

---

## Comparison Matrix

| 维度 | Chrome DevTools MCP | Playwright MCP | Browserbase MCP | Puppeteer MCP |
|---|---|---|---|---|
| **维护者** | Google (官方) | Microsoft (官方) | Browserbase (创业公司) | MCP 官方 (已弃用) |
| **状态** | 活跃 (v1.1.1) | 活跃 | 活跃 | **已弃用** |
| **浏览器** | Chrome/Chromium only | Chrome/Firefox/WebKit/Edge | 云端 Chromium | Chromium |
| **底层引擎** | Puppeteer + CDP | Playwright | Stagehand + Playwright | Puppeteer |
| **工具数量** | 45 (完整) / 3 (slim) | ~21 (MCP) / 40+ (CLI) | 6 | 7 |
| **核心定位** | 调试 + 检查 | 驱动 + 自动化 + 测试 | 云端 AI 自动化 | 已弃用 |
| **性能分析** | Lighthouse, Core Web Vitals, V8 堆快照, CrUX | 仅 trace | 无 | 无 |
| **网络 mock** | 仅检查 | mock + 拦截 | 无 | 无 |
| **跨浏览器** | 不支持 | Chrome/Firefox/WebKit/Edge | 不适用 | 不支持 |
| **内存分析** | V8 堆快照 + retainers | 无 | 无 | 无 |
| **AI 辅助** | 无（确定性工具） | 无（确定性工具） | Stagehand AI 操作 | 无 |
| **浏览器扩展** | 安装/卸载/管理 | 无 | 无 | 无 |
| **Token 效率** | Slim 模式极优；完整模式中等 | 工具定义 ~14k tokens；较重 | 6 工具，轻量 | 7 工具，轻量 |
| **连接模式** | 自动启动/autoConnect/browser-url/wsEndpoint | Persistent/Isolated/Extension | 云端 + STDIO | 本地 |
| **GitHub Stars** | 42,053 | (Playwright 主仓库 70k+) | 较少 | 已归档 |
| **npm 周下载** | ~230 万 | (Playwright 包级) | 较少 | 已归档 |
| **安全性** | 无已知 CVE，浏览器内容完全暴露 | CVE-2025-9611 (CSRF) | 依赖云端安全 | 无维护，不安全 |
| **学习曲线** | 低（npx 直跑） | 低（npx 直跑） | 中（需 API Key + 账号） | — |
| **迁移成本** | 低（标准 MCP 协议） | 低（标准 MCP 协议） | 中（云端依赖） | — |

---

## Risks & Caveats

| Risk | Severity | Mitigation |
|---|---|---|
| Token 消耗过高 | 高 | 使用 slim 模式/CLI；默认禁用 MCP，需要时启用 |
| 浏览器内容暴露给 AI | 中 | 避免在连接的浏览器中打开敏感页面；使用隔离 profile |
| MCP server 安全漏洞 | 中 | 仅使用官方维护的 server；关注 CVE 通报 |
| Chrome DevTools MCP 仅支持 Chrome | 中 | 跨浏览器测试场景使用 Playwright MCP 补充 |
| Puppeteer MCP 已弃用 | 高 | 不使用；迁移到 Chrome DevTools MCP |
| Browserbase 云端依赖 | 低 | 本地场景使用 Playwright/Chrome DevTools 替代 |
| 工具数量过多导致 agent 性能下降 | 中 | 控制同时启用的 MCP 工具数量在 10 个以内 |

### 社区痛点

Token 消耗是社区讨论最多的问题：
- 有用户报告 MCP 工具消耗了 83.3k tokens (41.6% 上下文窗口)
- 5 个 MCP server 加起来可达 45k tokens
- 工具数量超过 ~10 个时 agent 性能急剧下降
- 社区方案：Token Saver MCP、安装但默认禁用、使用 slim/CLI 模式

**Sources:**
- [MCP Server tools using 83.3k tokens — Reddit](https://www.reddit.com/r/ClaudeCode/comments/1mwxfit/mcp_server_tools_using_up_833k_tokens_416_of/)
- [3 Major pain points with MCP servers — Reddit](https://www.reddit.com/r/mcp/comments/1ncozfz/3_major_pain_points_with_mcp_servers_are_context/)

---

## Recommendation

**核心结论：Playwright 负责驱动，Chrome DevTools 负责调试。按动词选择，不按品牌选择。**

### 具体建议

1. **自动化/测试场景** — 使用 Playwright MCP (或更好的 Playwright CLI)
   - 跨浏览器测试、表单填写、E2E 测试生成
   - Microsoft 官方推荐 CLI 模式用于 coding agent

2. **调试/性能分析场景** — 使用 Chrome DevTools MCP
   - 性能 trace、Core Web Vitals 分析、内存泄漏检测
   - Lighthouse 审计、V8 堆快照

3. **云端/反检测场景** — 使用 Browserbase MCP
   - 需要云端浏览器、代理轮换、大规模并发
   - 不适合需要精细控制的场景

4. **Token 敏感场景** — 使用 slim 模式或 CLI
   - Chrome DevTools MCP 的 slim 模式仅 3 个工具
   - Playwright CLI 比 MCP 更 token 高效

5. **两者都安装，按需启用** — 这是社区公认的最佳实践
   - 默认禁用，需要调试时启用 Chrome DevTools
   - 需要自动化时启用 Playwright

### 不推荐

- **Puppeteer MCP** — 已弃用，无安全更新
- **BrowserTools MCP (AgentDesk)** — 项目已停止维护

### 前瞻趋势

| 趋势 | 影响 |
|---|---|
| WebMCP (Google, W3C) | 网站显式支持 agent 交互，下一个范式转变 |
| CLI 优于 MCP | Token 效率驱动，coding agent 趋向 CLI |
| Chrome DevTools CLI | Google 也在推 token 高效的 CLI 替代方案 |
| Stagehand v3 | AI 原生浏览器自动化，性能提升 20-40% |

**Confidence level:** 高 — 基于 Google/Microsoft 官方文档、npm 数据、GitHub 统计和社区实际体验交叉验证

---

## Sources

| Source | URL | Used for |
|---|---|---|
| Chrome DevTools MCP GitHub | https://github.com/ChromeDevTools/chrome-devtools-mcp | 架构、能力、社区数据 |
| Chrome 官方博客公告 | https://developer.chrome.com/blog/chrome-devtools-mcp | 定位与发布信息 |
| Chrome DevTools MCP 设计原则 | https://github.com/ChromeDevTools/chrome-devtools-mcp/blob/main/docs/design-principles.md | 架构设计 |
| Chrome DevTools MCP 工具参考 | https://github.com/ChromeDevTools/chrome-devtools-mcp/blob/main/docs/tool-reference.md | 能力清单 |
| Chrome DevTools MCP npm | https://www.npmjs.com/package/chrome-devtools-mcp | 下载量、版本 |
| Playwright MCP GitHub | https://github.com/microsoft/playwright-mcp | 架构、能力 |
| Playwright 官方 MCP 文档 | https://playwright.dev/docs/getting-started-mcp | 使用指南 |
| Driving vs Debugging — Steve Kinney | https://stevekinney.com/writing/driving-vs-debugging-the-browser | 核心对比框架 |
| Playwright vs Chrome DevTools — Mastalerz | https://mastalerz.it/comparing-playwright-mcp-vs-chrome-devtools-mcp-what-they-are-how-to-use-them-and-configuration-details/ | 详细对比 |
| Browserbase MCP GitHub | https://github.com/browserbase/mcp-server-browserbase | 架构、能力、定价 |
| Browserbase 定价 | https://www.browserbase.com/pricing/ | 定价信息 |
| Puppeteer MCP npm | https://www.npmjs.com/package/@modelcontextprotocol/server-puppeteer | 弃用状态 |
| MCP servers-archived | https://github.com/modelcontextprotocol/servers-archived | 弃用确认 |
| Reddit: MCP token 消耗讨论 | https://www.reddit.com/r/ClaudeCode/comments/1mwxfit/ | Token 痛点 |
| Reddit: MCP 三大痛点 | https://www.reddit.com/r/mcp/comments/1ncozfz/ | 社区反馈 |
| The Agentic Browser Landscape 2026 | https://nohacks.co/blog/agentic-browser-landscape-2026 | 行业趋势 |
| What is WebMCP | https://nohacks.co/blog/what-is-webMCP/ | WebMCP 前瞻 |
| CVE-2025-9611 Playwright | https://www.sentinelone.com/vulnerability-database/cve-2025-9611/ | 安全漏洞 |
| Chrome DevTools MCP Security | https://github.com/ChromeDevTools/chrome-devtools-mcp/blob/main/SECURITY.md | 安全模型 |
