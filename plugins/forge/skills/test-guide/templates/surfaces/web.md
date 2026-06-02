---
title: "Web E2E 测试策略"
domains: [testing, web]
---

<!-- Surface strategy template for Web. Filled by test-guide skill at runtime. -->

# Web E2E 测试策略

## 文件位置

- **目录**: `tests/e2e/` 或 `tests/<journey>/`（当 Journey 仅包含 Web 测试时）
- **文件命名**: `<feature>.<page>.spec.<ext>` 或 `<feature>.<journey>.test.<ext>`
- **Build tag**: `@web-e2e`（BDD tag）——这是"e2e"术语正确使用的两个 surface 之一
- **约束**: 仅 Web 和 Mobile surface 允许使用"e2e"术语

## 隔离模型

- **浏览器隔离**: 每个测试使用独立的浏览器上下文（browser context），不共享 cookie、localStorage 或 session
- **会话复用**: 认证场景在 Journey 级别登录一次，所有后续测试继承认证会话。非认证流程测试不得独立登录
- **网络拦截**: 依赖外部服务的测试必须拦截网络请求并返回固定响应，测试不得要求网络访问
- **视口管理**: 所有测试必须声明明确的默认视口尺寸。响应式测试使用显式视口切换作为前置步骤

## 断言重点

Web E2E 测试断言以下维度:

| 维度 | 断言模式 | 示例 |
|------|----------|------|
| DOM 元素可见性 | 元素在页面中可见 | 目标元素出现在 DOM 中且可见 |
| 用户操作响应 | 交互后 UI 状态变更 | 点击后表单提交成功提示出现 |
| 页面 URL 变更 | 导航后 URL 匹配预期 | 跳转后 URL 包含预期路径 |
| 元素属性值 | 元素属性匹配预期 | 输入框的 value 属性包含预期文本 |

**元素定位策略优先级**:
1. Test ID（data-testid 或等效稳定测试属性）
2. Accessible role + name（语义角色和可访问名称）
3. Label（表单标签、aria-label、placeholder）
4. Text content（可见文本作为最后手段）

**禁止**: CSS class 选择器——CSS 类是样式实现细节，样式重构会导致测试无功能回归地失败。

## 超时策略

- **元素等待超时**: 等待元素可见/可操作的显式超时
- **导航超时**: 页面导航完成的显式超时
- **测试函数级超时**: 测试运行器内置超时机制限制总执行时间
- **约束**: 每个阻塞操作必须指定超时，超时值来自 Convention 配置而非测试代码中的魔术数字

## 生命周期

1. **Setup**: 创建浏览器上下文，配置视口，建立认证状态（如需要）
2. **Navigate**: 导航到目标页面
3. **Interact**: 按测试步骤执行用户操作（点击、输入、选择）
4. **Assert**: 验证 DOM 元素可见性、内容、属性、URL 状态
5. **Teardown**: 关闭浏览器上下文，清理测试数据

## Contract/Journey 比例

Web surface 目标 **50/50 平衡比例**。

- **公式**: `Contract 测试函数数 / (Contract 测试函数数 + Journey 冒烟测试函数数) * 100%`
- **目标**: 约 50% Contract 测试，50% Journey 冒烟测试
- **实现**: 为每个 Outcome 生成 Contract 测试，同时生成覆盖 happy path 和至少 1 个失败路径的增强型 Journey 冒烟测试
- **最低要求**: 每个 Journey 必须至少有 1 个冒烟测试（happy path）

## 反模式

| 反模式 | 危害 | 替代方案 |
|--------|------|----------|
| CSS class 选择器 | 样式重构导致测试无功能回归地失败 | 使用 data-testid 或语义定位器 |
| 截图唯一断言 | 跨环境不确定，信号弱 | 结构断言（可见性、文本内容）为主，截图为辅 |
| 每测试独立登录 | 执行时间倍增，引入 flakiness | Journey 级别登录一次，复用会话 |
| 无 testid 时回退 CSS 选择器 | 产生脆弱测试 | 回退到语义定位器（role、label、text） |
| 调试输出（console.log） | 填充 CI 输出，掩盖真实失败 | 生成的测试代码仅包含断言 |

## 断言偏好表

| 断言库 | mock 机制 | fixture 模式 |
|--------|-----------|-------------|
| {{ASSERTION_LIBRARY}} | {{MOCK_MECHANISM}} | {{FIXTURE_PATTERN}} |
