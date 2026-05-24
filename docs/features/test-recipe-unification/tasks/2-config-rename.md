---
id: "2"
title: "Rename config fields and add migration hint"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 2: Rename config fields and add migration hint

## Description

将 `forgeconfig.Config` 中的 `E2eTest` 字段重命名为 `Test`，YAML tag 从 `e2eTest` 改为 `test`。在 `parseAutoRaw()` 中检测旧键名 `e2eTest`，输出迁移提示并映射到新字段。同步更新 `autogen.go` 中的字段引用和任务 key（`run-e2e-tests` → `run-test`），以及 testdata 配置文件。

## Reference Files
- `proposal.md#Proposed-Solution` — defines config rename `auto.e2eTest` → `auto.test` and task key `run-e2e-tests` → `run-test`
- `proposal.md#Constraints-&-Dependencies` — parseAutoRaw migration: detect old key, output hint, map to new field; direct YAML tag rename
- `proposal.md#Non-Functional-Requirements` — migration prompt format, v3.1.0 removal timeline for old key detection logic
- `proposal.md#Impact-Analysis` — Config & Testdata section lists all config files requiring update

## Acceptance Criteria
- `forgeconfig.Config` struct field `E2eTest` renamed to `Test`, YAML tag `e2eTest` → `test`
- `parseAutoRaw()` detects old key `e2eTest` and outputs migration hint to stderr: `"config key 'auto.e2eTest' is renamed to 'auto.test' in v3.0.0; please update your config.yaml"`
- When old key detected, its value is mapped to `Test` field (not silently ignored)
- When only new key `test` is used, no migration hint output
- `autogen.go`: `auto.E2eTest` → `auto.Test`, `Key: "run-e2e-tests"` → `"run-test"`
- `.forge/config.yaml`: `e2eTest` → `test`
- `internal/cmd/testdata/forge-config.example.yaml`: `e2eTest` → `test`
- `internal/cmd/testdata/forge-config.schema.json`: `e2eTest` → `test`

## Hard Rules
- 不引入永久兼容层。`parseAutoRaw()` 中的旧键名检测逻辑是为 v3.0.x 过渡期设计的，v3.1.0 起遇到旧键名将直接报错

## Implementation Notes
- Config YAML 直接使用 `test` 键名，不保留 `e2eTest` 兼容
- testdata 文件是 Go 测试的 fixture，需与 config struct 同步更新以避免测试失败
