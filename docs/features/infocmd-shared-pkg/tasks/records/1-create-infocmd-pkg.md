---
status: "completed"
started: "2026-05-23 11:29"
completed: "2026-05-23 11:40"
time_spent: "~11m"
---

# Task Record: 1 创建 pkg/infocmd/ 共享工具包

## Summary
创建 pkg/infocmd/ 共享工具包，提供泛型 Discover[T]、FindByID[T]、ParseFrontmatter 和统一排序逻辑。支持 flat file（research/lesson）和 subdirectory+fixed-file（proposal）两种目录扫描模式，通过 ScanConfig[T] 配置结构体封装差异。

## Changes

### Files Created
- forge-cli/pkg/infocmd/infocmd.go
- forge-cli/pkg/infocmd/infocmd_test.go

### Files Modified
无

### Key Decisions
- 使用 ScanConfig[T] 结构体封装目录扫描差异（IsSubdir + FileName 字段控制模式），而非两个独立函数
- 通过 CreatedKey func(T) string 回调提取排序键，避免反射或接口约束，符合 Hard Rules 中 '泛型约束使用 any' 的要求
- 通过 IDKey func(T) string 回调统一 Slug/Name 标识符提取，适配 lesson 用 Name 而非 Slug 的差异
- ParseFrontmatter 作为导出函数直接提取自三个包中逐字节相同的实现，零行为变更

## Test Results
- **Tests Executed**: Yes
- **Passed**: 13
- **Failed**: 0
- **Coverage**: 84.4%

## Acceptance Criteria
- [x] parseFrontmatter 仅存在于 pkg/infocmd/ 一处
- [x] Discover[T] 泛型函数支持两种目录模式：flat files 和 subdirectory+fixed-file
- [x] FindByID[T] 泛型函数支持配置化的标识符提取（Slug 或 Name）
- [x] 排序逻辑统一：按 Created 降序，mtime 回退
- [x] 包通过 go vet 和编译
- [x] 现有三个命令的行为不受影响（此 task 仅创建新包，不修改现有代码）

## Notes
Hard Rules 全部遵守：未修改现有三个包的代码，泛型约束使用 any，新包未导入现有三个包。
