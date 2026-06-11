---
created: "2026-05-21"
tags: [testing, local-dev-deployment, architecture]
---

# forge task index 不生成测试管道任务

## Problem

运行 `forge task index --feature <slug>` 后，只生成了业务任务和 stage-gates，没有生成测试管道任务（gen-test-cases、run-e2e-tests、graduate-tests、verify-regression、clean-code、consolidate-specs 等）。index 验证通过但任务数明显偏少。

## Root Cause

因果链（3 层）：

1. **测试管道任务由 `GetBreakdownTestTasks(interfaces, auto)` 生成，首行守卫 `if len(interfaces) == 0 { return nil }` 直接返回空**
2. **`interfaces` 由 `ReadInterfaces(projectRoot)` 从 `config.yaml` 的 `interfaces` 字段读取，未配置则为空 slice**
3. **项目未在 `.forge/config.yaml` 中声明 `interfaces` 字段**

### Historical Context (已解决)

此问题最初源于旧的语言检测系统。`DetectLanguages` 在项目根目录检查 `go.mod`/`package.json` 等文件是否存在。当 `go.mod` 位于子目录（如 `forge-cli/go.mod`）而非项目根目录时，检测失败，返回空 slice，导致 `interfaces` 为空，测试管道任务被静默跳过。

该问题通过 **Decouple Test Tasks from Languages** 重构彻底解决：移除了语言检测逻辑，改为用户在 `config.yaml` 中显式声明 `interfaces`。测试任务的粒度由接口类型（CLI、API、WebUI、TUI）决定，而非由实现语言决定。

## Solution

在 `.forge/config.yaml` 中添加 `interfaces` 字段：

```yaml
interfaces:
  - api
  - cli
```

配置后 `ReadInterfaces` 直接读取用户声明的接口类型列表，`GetBreakdownTestTasks` 按接口类型生成测试管道任务。不再依赖文件检测或语言推断。

如果未配置 `interfaces`，`BuildIndex` 会输出明确警告提示用户配置，测试管道任务被跳过（有警告，非静默）。

## Reusable Pattern

**当 `forge task index` 缺少预期的自动生成任务时，检查链路**：

1. `forge task index` → `BuildIndex` → `detectMode`（breakdown/quick/""）
2. `needsTestPipeline`（是否有 `coding.*` 类型的任务）
3. `GenerateTestTasks(mode, interfaces, auto)`
4. `GetBreakdownTestTasks` 守卫：`len(interfaces) == 0`
5. `ReadInterfaces` → 读取 `config.yaml` 的 `interfaces` 字段

**关键检查点**：`.forge/config.yaml` 中是否配置了 `interfaces`。如果没有，`BuildIndex` 会输出警告并跳过测试管道任务。

## Example

```bash
# 诊断：检查 interfaces 配置
cat .forge/config.yaml | grep -A 5 interfaces

# 如果为空，添加接口类型
echo 'interfaces:
  - api
  - cli' >> .forge/config.yaml
# 然后重新生成
rm docs/features/<slug>/tasks/index.json
forge task index --feature <slug>
```

## Related Files

- `forge-cli/pkg/forgeconfig/detect.go` — `ReadInterfaces`（已移除 `DetectLanguages`, `ReadLanguages`, `UnionLanguageInterfaces`）
- `forge-cli/pkg/task/build.go` — `BuildIndex`, `GenerateTestTasks`, `GetBreakdownTestTasks`
- `forge-cli/pkg/task/testgen.go` — 接口类型驱动的任务生成（已移除 `profileSuffix`, `suffixLetter`）
- `.forge/config.yaml` — `interfaces` 字段

## References

- 发现于 `forge-architecture-simplification` feature 的 breakdown-tasks 流程中
- 旧 index.json 有 32 tasks（含测试管道），重新生成后仅 24 tasks，缺少 11 个测试管道任务
- 通过 `decouple-test-tasks-from-languages` feature 彻底重构解决
