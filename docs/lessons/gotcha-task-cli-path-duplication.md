# task-cli Path Duplication Bug

## Problem

`task claim` 输出的 RECORD 路径包含重复的 `tasks/` 目录：

- 期望：`docs/features/agent-task-center/tasks/records/1.1-init-monorepo-scaffold.md`
- 实际：`docs/features/agent-task-center/tasks/tasks/records/1.1-init-monorepo-scaffold.md`

同样的问题影响 FILE 字段，以及 `query`、`record` 命令中的路径构造。

## Root Cause

**数据层和代码层各自添加了一次 `tasks/` 前缀，导致路径拼接时重复。**

1. `index.json` 中 `record` 字段值已包含 `tasks/` 前缀：`"tasks/records/1.1-xxx.md"`
2. `GetTaskFile()` 函数又会添加 `TasksDirName`（即 `"tasks"`）：
   ```go
   func GetTaskFile(feature, filename string) string {
       return filepath.Join(FeaturesDir, feature, TasksDirName, filename)
   }
   ```
3. 两者拼接后产生 `docs/features/<slug>/tasks/tasks/records/<file>.md`

**涉及位置（3处）：**
- `claim.go:286` — printTaskDetails 输出 RECORD 路径
- `query.go:58` — query 命令输出 RECORD 路径
- `record.go:143` — record 命令写入文件路径

## Solution

两种修复方向：

**方案 A（修数据）**：`index.json` 中 `record` 和 `file` 字段去掉 `tasks/` 前缀，仅保留 `records/xxx.md` 和 `xxx.md`。需要同步修改生成 index.json 的逻辑。

**方案 B（修代码）**：对 RECORD 等已含完整相对路径的字段，改用 `GetFeatureDir()` 直接拼接，不再经过 `GetTaskFile()`：
```go
// 修改前
filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, t.Record))
// 修改后
filepath.Join(projectRoot, feature.GetFeatureDir(featureSlug), t.Record)
```

## Key Takeaway

**当路径由数据+代码共同构造时，必须明确边界——谁负责哪一段路径。**

- 如果数据存储的是相对于 feature 目录的完整路径（如 `tasks/records/xxx.md`），代码不应再添加 `tasks/` 前缀
- 如果数据只存文件名（如 `records/xxx.md`），代码才需要添加目录前缀
- 路径构造函数（如 `GetTaskFile`）的调用方必须知道传入参数的"锚点"在哪里
