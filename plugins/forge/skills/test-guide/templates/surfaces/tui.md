---
title: "TUI 测试策略"
domains: [testing, tui]
---

<!-- Surface strategy template for TUI. Filled by test-guide skill at runtime. -->

# TUI 测试策略

## 文件位置

- **目录**: `tests/tui/` 或 `tests/<journey>/`（当 Journey 仅包含 TUI 测试时）
- **文件命名**: `<feature>_<screen>_test.<ext>`（Go）、`test_<feature>_<screen>.<ext>`（Python）、`<feature>.<screen>.test.<ext>`（Node.js）
- **Build tag**: `//go:build tui_functional`（Go）、`@tui-functional`（BDD tag）
- **约束**: 不得使用 `e2e` 作为 build tag 或测试分类名

## 隔离模型

- **非交互执行模型**: TUI 测试脚本必须使用非交互执行（stdin pipe，非真实终端），不使用交互测试模式
- **终端尺寸契约**: 测试必须控制终端尺寸环境以消除渲染差异，设置 `TERM=dumb` 或固定 `LINES` 和 `COLUMNS` 环境变量
- **进程边界隔离**: 编译独立二进制，通过子进程 + stdin pipe 模拟用户输入
- **工作目录隔离**: 每个测试使用临时目录作为工作目录

## 断言重点

TUI 测试断言以下维度:

| 维度 | 断言模式 | 示例 |
|------|----------|------|
| 精确文本 | 输出包含精确字符串 | 捕获输出包含预期屏幕文本 |
| 正则匹配 | 输出匹配模式 | 捕获输出匹配预期模式 |
| 快照 | 输出匹配 golden file | 与参考文件对比 |
| 缺失 | 输出不包含文本 | 捕获输出不包含错误文本 |

**ANSI 净化**: 断言前必须从捕获输出中剥离 ANSI 转义序列，或使用专用终端输出解析。不得对包含控制序列的原始终端输出进行断言。

**稳定状态检测**: 测试必须定义"屏幕渲染完成"的可观测信号（stdout 稳定、子进程退出、特定标记字符串出现），而非依赖基于时间的假设。

## 超时策略

- **进程级超时**: 生成的子进程必须在可配置秒数内退出
- **测试函数级超时**: 测试运行器内置超时机制限制总执行时间
- **稳定状态等待**: 轮询检查输出内容直到超时或匹配预期，不使用固定延时

## 生命周期

1. **Setup**: 编译二进制，设置终端环境（TERM、LINES、COLUMNS），准备 stdin pipe 中的按键序列
2. **Execute**: 启动二进制并传入 stdin 按键序列
3. **Capture**: 捕获 stdout、stderr 和 exit code
4. **Assert**: 对捕获输出进行 ANSI 净化后断言内容，同时断言 exit code
5. **Teardown**: 终止残留子进程，清理临时资源

**按键编码**:

| 按键 | stdin 编码 | 说明 |
|------|-----------|------|
| Enter | `\n` 或 `\r` | 确认/提交 |
| Escape | `\x1b` | 取消/返回 |
| Tab | `\t` | 下一个字段 |
| 方向键 | `\x1b[A/B/C/D` | 上/下/右/左 |
| 普通字符 | 字面字符 | 字母、数字、符号 |
| Ctrl+C | `\x03` | 中断 |

## Contract/Journey 比例

TUI surface 目标 **Contract 测试比例 >= 80%**。

- **公式**: `Contract 测试函数数 / (Contract 测试函数数 + Journey 冒烟测试函数数) * 100%`
- **最低要求**: 每个 Journey 生成 M 个 Contract 测试函数和恰好 1 个 Journey 冒烟测试（happy path）
- **小型 Journey 调整**: 若 Journey 总 Outcomes < 5，1 个冒烟测试仍只计为 1 个函数，比例自然保持较高
- **禁止**: 不得跳过冒烟测试来膨胀比例——每个 Journey 必须至少有 1 个冒烟测试

## 反模式

| 反模式 | 危害 | 替代方案 |
|--------|------|----------|
| Sleep 等待屏幕过渡 | 时序不稳定 | 轮询+超时检查输出内容 |
| 源码检查替代运行时 | 测试实现结构而非行为 | 始终执行 TUI 二进制，捕获输出并断言 |
| 无 stdin pipe 的交互提示 | 测试挂起等待输入，CI 超时 | 始终将完整按键序列 pipe 到 stdin |
| 对原始 ANSI 输出断言 | 终端模拟器差异导致 flaky | 断言前剥离 ANSI 序列 |
| 忽略终端尺寸 | 不同环境渲染不一致 | 固定 TERM/LINES/COLUMNS |

**需要真实终端的测试**: 鼠标交互、窗口缩放、颜色渲染等无法通过 stdin pipe 测试的场景，必须显式标记为 manual-only。

## 断言偏好表

| 断言库 | mock 机制 | fixture 模式 |
|--------|-----------|-------------|
| {{ASSERTION_LIBRARY}} | {{MOCK_MECHANISM}} | {{FIXTURE_PATTERN}} |
