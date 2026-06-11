---
title: "Mobile E2E 测试策略"
domains: [testing, mobile]
---

<!-- Surface strategy template for Mobile. Filled by test-guide skill at runtime. -->

# Mobile E2E 测试策略

## 文件位置

- **目录**: `tests/<surfaceKey>/<journey>/`（多 surface 项目）或 `tests/<journey>/`（单 surface 项目）。Journey 名称由 gen-journeys 生成，surfaceKey 为 `forge surfaces` 输出的 key
- **文件命名**: `step<N>_<action>.yaml`（Maestro 格式）、`step<N>_<action>_deeplink.yaml`（深度链接变体）
- **Build tag**: `@mobile-e2e`（BDD tag）——这是"e2e"术语正确使用的两个 surface 之一
- **约束**: 仅 Web 和 Mobile surface 允许使用"e2e"术语

## 隔离模型

- **应用状态重置**: 每个测试必须在执行前清理应用状态（kill + clear data 或 uninstall + reinstall），不依赖前一个测试遗留的状态
- **权限处理**: 系统权限对话框（相机、位置、通知）必须作为前置步骤处理，不得在测试执行中临时忽略
- **生命周期声明**: 每个测试必须显式声明应用生命周期——启动时启动应用，结束时终止应用，不得假设应用已运行
- **元素定位优先级**:
  1. Accessibility ID（accessibility label、test ID、accessibility identifier）
  2. Resource ID（平台资源标识符）
  3. Text content（可见文本作为最后手段）

## 断言重点

Mobile E2E 测试断言以下维度:

| 维度 | 断言模式 | 示例 |
|------|----------|------|
| UI 元素可见性 | 目标元素在屏幕上可见 | 特定按钮或文本出现在当前屏幕 |
| 用户操作响应 | 交互后屏幕状态变更 | 点击后跳转到预期屏幕 |
| 屏幕 ID 变更 | 导航后屏幕标识匹配预期 | 切换后显示目标屏幕内容 |

**屏幕过渡断言**: 每个导航操作后必须断言目标屏幕可见再继续。移动端导航涉及动画、网络加载和状态转换——未确认目标屏幕就继续操作会导致后续交互命中错误屏幕元素。

## 超时策略

- **屏幕等待超时**: 等待目标屏幕元素出现的显式超时
- **测试函数级超时**: 测试运行器内置超时机制限制总执行时间
- **约束**: 超时值来自 Convention 配置而非测试代码中的魔术数字

## 生命周期

1. **Setup**: 启动应用（launchApp），清理应用状态，处理权限对话框
2. **Navigate**: 通过 UI 操作或深度链接导航到目标屏幕
3. **Assert**: 验证目标屏幕可见，确认预期内容
4. **Interact**: 执行用户操作（tap、swipe、input）
5. **Teardown**: 终止应用（killApp），清理应用数据

**Maestro YAML 骨架结构**:
- `appId` 声明（来自 Fact Table MOBILE_APP_ID）
- `onFlowStart: [launchApp]` 生命周期钩子
- `onFlowEnd: [killApp]` 生命周期钩子
- 命令序列
- `assertVisible` 断言

## Contract/Journey 比例

Mobile surface 遵循 **best-effort** 策略——不以 Contract 测试比例衡量。

- **输出格式**: Maestro YAML 文件
- **骨架要求**: 每个 Maestro YAML 必须包含 appId、onFlowStart/onFlowEnd 钩子、命令序列和 assertVisible 断言
- **深度链接测试**: 每个 Journey 步骤额外生成一个通过 URL scheme 打开应用的 Maestro YAML
- **Manual-only 标记**: 无法可靠自动化的场景（多指手势、物理传感器、生物识别、系统级交互）必须标记为 manual-only
- **硬性规则**: Mobile 测试生成不得导致管线失败——任何生成问题产生带 manual-only 标记的骨架

## 反模式

| 反模式 | 危害 | 替代方案 |
|--------|------|----------|
| 像素坐标交互 | 跨屏幕尺寸/分辨率失败 | 使用基于元素的定位（accessibility ID、resource ID、text） |
| 测试间状态泄漏 | 隔离通过但序列失败 | 每个测试前清理应用状态 |
| 未处理的权限对话框 | 阻塞测试执行，导致超时 | 预授权权限或作为显式前置步骤处理 |
| 无屏幕断言的导航 | 后续操作命中错误屏幕元素 | 每次导航后断言目标屏幕可见 |
| 物理设备能力假设 | CI 环境模拟器无传感器 | 限制为模拟器/仿真器支持的操作，标记需要物理设备的测试 |

## 断言偏好表

| 断言库 | mock 机制 | fixture 模式 |
|--------|-----------|-------------|
| {{ASSERTION_LIBRARY}} | {{MOCK_MECHANISM}} | {{FIXTURE_PATTERN}} |
