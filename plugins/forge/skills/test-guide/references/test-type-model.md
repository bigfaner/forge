---
title: "Surface 测试类型模型"
domains: [testing, surface, test-type]
---

<!-- Migrated from docs/reference/test-type-model.md. Complete reference for Surface -> Test Type mapping. -->

# Surface 测试类型模型

本文档定义 Surface -> Test Type 的映射模型，作为 Forge 测试管线的概念权威参考。所有 skill 文件、task type 命名和 justfile recipe 中的测试类型术语以本文档为准。

## 映射表

| Surface | Test Type（EN） | 测试类型（CN） | 验证维度 | 执行模型 |
|---------|-----------------|---------------|---------|---------|
| `cli` | CLI Functional Test | CLI 功能测试 | 进程退出码 + stdout 文本 + stderr 文本 | 子进程执行 |
| `tui` | Terminal Functional Test | 终端功能测试 | 终端输出文本 + stdin 交互响应序列 | 子进程 + stdin pipe |
| `api` | API Functional Test | API 功能测试 | HTTP 状态码 + 响应体 JSON + 响应 Header | HTTP 客户端 |
| `web` | Web E2E Test | Web 端到端测试 | DOM 元素可见性 + 用户操作响应 + 页面 URL 变更 + 元素属性值 | 浏览器自动化 |
| `mobile` | Mobile E2E Test | 移动端端到端测试 | UI 元素可见性 + 用户操作响应 + 屏幕 ID 变更 | Maestro YAML / 手动验证 |

## 分类标准

本模型采用两级分类：

1. **一级分类键：Surface**（cli / tui / api / web / mobile）— 决定测试的执行模型（子进程 / HTTP / 浏览器 / 设备自动化）
2. **二级属性：测试范围**（功能测试 / 端到端测试）— 由验证机制决定

二级属性的判定基于**验证机制**，而非技术栈覆盖深度：

- **功能测试**：通过协议级调用（子进程调用、HTTP 请求）验证输入-输出行为。测试工具在协议边界上观测被测系统的响应——CLI 测试观测子进程的退出码和 stdout/stderr，API 测试观测 HTTP 响应的状态码和 body。
- **端到端测试**：通过设备级自动化（浏览器驱动、移动设备自动化）模拟真实用户操作序列。Web 测试通过 Playwright 模拟浏览器操作，Mobile 测试通过 Maestro 模拟移动设备操作。

关键区分：CLI 测试可以遍历完整技术栈（如读写数据库后输出结果），API 测试也可以触发从 HTTP 请求到持久层的完整调用链。它们在技术栈覆盖上可能是"端到端"的，但验证机制是在协议边界上的单次调用观测，而非通过设备级自动化模拟用户操作流程。本模型中"功能测试"/"端到端测试"标签反映的是验证机制，不是技术栈覆盖深度。

## 语义定义

- **CLI 功能测试**：编译独立二进制，通过子进程调用，验证命令行参数解析、输出格式、退出码、错误处理。不测试内部函数，通过进程边界隔离。
- **终端功能测试**：编译独立二进制，通过 stdin pipe 模拟用户输入，验证终端渲染输出（ANSI 序列处理、布局、异步 Cmd 响应）。与 CLI 功能测试的区别在于需要模拟交互输入流。
- **API 功能测试**：启动 HTTP 服务器（或使用测试服务器），发送请求，验证响应符合 Contract 定义的六个维度。此处的 "Contract" 指 Forge 在 gen-contracts 阶段生成的 API 行为规约。
- **Web 端到端测试**：启动 dev server，通过浏览器自动化（Playwright）模拟用户操作，验证 UI 渲染、交互逻辑、跨页面导航和跨组件状态流转。覆盖从用户输入到持久层再回到 UI 的完整用户旅程。
- **移动端端到端测试**：通过 Maestro YAML 定义操作序列，驱动移动端 UI，验证渲染、交互和屏幕导航。与 Web 端到端测试逻辑同构，均通过设备级自动化覆盖完整用户旅程。Best-effort 模式，部分场景标记为 manual-only。

## "e2e" 术语使用约束

"e2e"（端到端）一词**仅用于** Web 和 Mobile surface 的端到端测试上下文。以下用法被禁止：

- 将 CLI / TUI / API surface 的测试称为 "e2e 测试" 或 "端到端测试"
- 在 justfile recipe、task type、测试报告或文档中将所有 Forge 生成的测试统称为 "e2e 测试"

新 surface 加入时，按分类标准的判定规则归入"功能测试"或"端到端测试"：验证机制为协议级调用则归入功能测试，为设备级自动化则归入端到端测试。若出现混合验证机制的 surface，可在本文档中新增第三分类。
