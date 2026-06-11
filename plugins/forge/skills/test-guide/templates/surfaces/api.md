---
title: "API 测试策略"
domains: [testing, api]
---

<!-- Surface strategy template for API. Filled by test-guide skill at runtime. -->

# API 测试策略

## 文件位置

- **目录**: `tests/<surfaceKey>/<journey>/`（多 surface 项目）或 `tests/<journey>/`（单 surface 项目）。Journey 名称由 gen-journeys 生成，surfaceKey 为 `forge surfaces` 输出的 key
- **文件命名**: `<feature>_<endpoint>_test.<ext>`（Go）、`test_<feature>_<endpoint>.<ext>`（Python）、`<feature>.<endpoint>.test.<ext>`（Node.js）
- **Build tag**: `//go:build api_functional`（Go）、`@api-functional`（BDD tag）
- **约束**: 不得使用 `e2e` 作为 build tag 或测试分类名

## 隔离模型

- **协议边界隔离**: 通过 HTTP 客户端发送请求、观测响应，验证 API 行为。测试工具在协议边界上观测被测系统的响应
- **服务隔离**: 测试启动独立的 HTTP 服务器（或使用测试服务器），不依赖共享服务实例
- **数据隔离**: 每个测试创建自己的测试数据并在完成后清理，或使用不会持久化的临时数据
- **认证隔离**: 认证场景按 Auth Plan 分类——login-test 使用独立认证，auth-required-test 使用缓存共享认证

## 断言重点

API 测试必须对以下每个维度包含具体断言:

| 维度 | 断言模式 | 示例 |
|------|----------|------|
| Status code | 精确断言 HTTP 状态码 | 200、401、404 等 |
| Response body fields | 断言具体字段值存在且匹配 | 响应中特定字段等于预期值 |
| Response body schema | 断言响应结构匹配预期形状 | 必填字段存在且类型正确 |
| Response headers | 断言 Header 值（如 Content-Type） | Content-Type 匹配预期值 |

**最低要求**: 每个 API 测试除状态码外必须断言至少一个具体的响应体字段。仅检查状态码是必要的但不是充分的。

**幂等性检查**: 对于 PUT 和 DELETE 端点，测试必须验证重复相同请求产生相同结果。

## 超时策略

- **连接超时**: HTTP 客户端必须设置连接超时
- **读写超时**: HTTP 客户端必须设置读/写超时
- **测试函数级超时**: 测试运行器内置超时机制限制总执行时间
- **约束**: 任何 HTTP 请求不得无限期等待连接或响应

## 生命周期

1. **Setup**: 启动测试服务器（或确认服务器已运行），获取基础 URL 和端口，配置认证
2. **Build request**: 组合 base URL 与端点路径，设置 HTTP 方法、Header 和请求体
3. **Send**: 使用 Convention 指定的 HTTP 客户端发送请求
4. **Assert**: 对状态码、响应体、响应头进行具体断言
5. **Teardown**: 清理测试创建的数据，关闭连接

## Contract/Journey 比例

API surface 目标 **50/50 平衡比例**。

- **公式**: `Contract 测试函数数 / (Contract 测试函数数 + Journey 冒烟测试函数数) * 100%`
- **目标**: 约 50% Contract 测试，50% Journey 冒烟测试
- **实现**: 为每个 Outcome 生成 Contract 测试，同时生成覆盖 happy path 和至少 1 个错误路径的增强型 Journey 冒烟测试
- **最低要求**: 每个 Journey 必须至少有 1 个冒烟测试（happy path）

## 反模式

| 反模式 | 危害 | 替代方案 |
|--------|------|----------|
| 硬编码 URL | 端口或主机变更即失败 | 使用 config/环境变量构造 URL |
| 缺少错误 Contract 测试 | 错误格式静默变更不被发现 | 每个端点至少一个错误场景测试 |
| 空洞的"返回成功"断言 | 200 + 空响应体也通过测试 | 必须断言具体响应体字段 |
| Sleep 等待 | 时序不稳定 | 使用事件驱动等待或轮询+超时 |
| 硬编码配置 | 环境变更即失败 | 所有配置来自环境变量或 Fact Table |

## 断言偏好表

| 断言库 | mock 机制 | fixture 模式 |
|--------|-----------|-------------|
| {{ASSERTION_LIBRARY}} | {{MOCK_MECHANISM}} | {{FIXTURE_PATTERN}} |
