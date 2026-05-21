---
created: "2026-05-21"
tags: [testing, local-dev-deployment, architecture]
---

# forge task index 不生成测试管道任务

## Problem

运行 `forge task index --feature <slug>` 后，只生成了业务任务和 stage-gates，没有生成测试管道任务（gen-test-cases、run-e2e-tests、graduate-tests、verify-regression、clean-code、consolidate-specs 等）。index 验证通过但任务数明显偏少。

## Root Cause

因果链（3 层）：

1. **测试管道任务由 `GetBreakdownTestTasks(languages, interfaces, auto)` 生成，但首行守卫 `if len(interfaces) == 0 { return nil }` 直接返回空**
2. **`interfaces` 由 `ReadInterfaces(projectRoot)` 计算，它先查 `config.yaml` 的 `interfaces` 字段，为空则 fallback 到 `defaultInterfaces` → `UnionLanguageInterfaces(ReadLanguages(projectRoot))`**
3. **`ReadLanguages` 先查 `config.yaml` 的 `languages` 字段，为空则 fallback 到 `DetectLanguages`，后者在项目根目录检查 `go.mod`/`package.json` 等文件是否存在。当 `go.mod` 位于子目录（如 `forge-cli/go.mod`）而非项目根目录时，检测失败，返回空 slice**

根本原因：**项目结构与 `DetectLanguages` 的假设不匹配**——monorepo 或子目录模块的 `go.mod` 不在项目根目录，且 config.yaml 未显式声明 `languages`。

## Solution

在 `.forge/config.yaml` 中显式添加 `languages` 字段：

```yaml
languages:
  - go
```

添加后 `ReadLanguages` 直接返回 `["go"]`，不再依赖文件检测。`UnionLanguageInterfaces(["go"])` 返回 `["api", "cli"]`，`GetBreakdownTestTasks` 正常生成测试管道任务。

## Reusable Pattern

**当 `forge task index` 缺少预期的自动生成任务时，检查链路**：

1. `forge task index` → `BuildIndex` → `detectMode`（breakdown/quick/""）
2. `needsTestPipeline`（是否有 `coding.*` 类型的任务）
3. `GenerateTestTasks(mode, languages, capabilities, auto)`
4. `GetBreakdownTestTasks` 守卫：`len(interfaces) == 0`
5. `ReadInterfaces` → `ReadLanguages` → `DetectLanguages` 文件检测

**关键检查点**：`.forge/config.yaml` 中是否显式配置了 `languages`。如果没有，且 `go.mod`/`package.json` 不在项目根目录，测试管道任务会静默跳过（无错误、无警告）。

## Example

```bash
# 诊断：检查 languages 检测结果
# 方法：查看 config.yaml 是否有 languages 字段
cat .forge/config.yaml | grep -A 5 languages

# 如果为空且 go.mod 在子目录，显式添加
echo 'languages:
  - go' >> .forge/config.yaml
# 然后重新生成
rm docs/features/<slug>/tasks/index.json
forge task index --feature <slug>
```

## Related Files

- `forge-cli/pkg/forgeconfig/detect.go` — `DetectLanguages`, `ReadInterfaces`, `UnionLanguageInterfaces`
- `forge-cli/pkg/task/build.go` — `BuildIndex`, `GenerateTestTasks`, `GetBreakdownTestTasks`
- `.forge/config.yaml` — `languages` 字段

## References

- 发现于 `forge-architecture-simplification` feature 的 breakdown-tasks 流程中
- 旧 index.json 有 32 tasks（含测试管道），重新生成后仅 24 tasks，缺少 11 个测试管道任务
