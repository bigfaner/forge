---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Iteration 0: Pre-Revision Report

## ATTACK_POINTS

- **[high]** 分类方案未覆盖约 15 个无已知前缀的 stderr 调用点（qualitygate progress、forensic output、upgrade.go lowercase error:），迁移将留下缺口 | quote: "The categorization scheme does not account for a significant number of stderr calls that lack any recognizable prefix" | improvement: produce exhaustive categorization table covering all 72 call sites with explicit rules for prefixless messages

- **[high]** AUTO-RESTORE-SKIP 被归为两个不同级别（ERROR 和 WARN），但三个调用点共享相同前缀，无法通过前缀机械区分 | quote: "The proposal lists AUTO-RESTORE-SKIP: under two different levels" | improvement: unify all AUTO-RESTORE-SKIP to WARN level or introduce sub-prefixes

- **[high]** forgelog.Init() 在 forge init 之前调用时 .forge/logs/ 目录不存在，形成启动悖论 | quote: "This creates a bootstrapping paradox" | improvement: Init() must create .forge/logs/ on demand via os.MkdirAll, falling back to stderr-only mode on failure

- **[medium]** 提案中 stderr 调用点计数 64 与实际代码库 72 不符，迁移计划工作量被低估 | quote: "64 fmt.Fprintf(os.Stderr, ...) call sites across the CLI" | improvement: re-audit and update the count to match current codebase

- **[medium]** forensic 命令的 stderr 是其主输出，若被迁移将产生重复输出并破坏其输出契约 | quote: "The proposal does not address the forensic command's stderr output pattern" | improvement: explicitly exclude forensic and intentional-stderr-output commands from migration scope

- **[medium]** upgrade.go 中 5 个小写 error: 前缀调用不匹配分类表中大写 ERROR: 模式 | quote: "The upgrade.go file uses lowercase error: prefix (5 call sites)" | improvement: normalize error: to ERROR: in upgrade.go or specify case-insensitive matching

- **[medium]** 秒级精度文件名在并发场景下会产生文件名冲突 | quote: "If two forge commands are invoked within the same second — they will collide on the same filename" | improvement: add PID or random suffix to filename for collision resistance

- **[medium]** 提案未定义日志文件写入失败时 stderr 输出的行为保证 | quote: "the proposal must guarantee that stderr output continues uninterrupted. But it does not specify the failure mode" | improvement: specify stderr-first-then-file write ordering and graceful degradation

- **[medium]** 小于 1KB 的性能估计未考虑 run-tasks 自动循环和多步骤测试场景 | quote: "This estimate of <1KB per invocation appears to be based on a typical interactive command run" | improvement: adopt buffered writes (bufio.Writer) with defer flush, update performance estimate

- **[medium]** 自动清理在创建新日志文件之前执行时可能删除当前活跃日志文件 | quote: "the command would delete its own active log file" | improvement: cleanup runs after new log file opened, errors are best-effort

- **[medium]** 提案未讨论为现有 forgeconfig.Config 添加 Logs 字段的配置迁移路径 | quote: "the proposal does not discuss the migration path for existing config files" | improvement: specify config struct extension and backward compatibility

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

- **[low]** 建议定义每行日志格式为带时间戳的标准格式 — subjective preference on format, not a factual defect
- **[low]** 建议将所有 AUTO-RESTORE-SKIP 统一归为 WARN 级别 — implementation suggestion, addressed by high-severity finding above
- **[low]** 建议日志文件名添加 PID 或随机后缀 — implementation suggestion, addressed by medium-severity finding above
- **[low]** 建议明确 forgelog.Init() 的失败语义 — implementation suggestion, addressed by high-severity finding above
- **[low]** 建议采用带显式 flush 的缓冲写入 — implementation suggestion, addressed by medium-severity finding above
- **[low]** 建议明确将 forensic 等命令排除在迁移范围之外 — implementation suggestion, addressed by medium-severity finding above
- **[low]** 建议统一 upgrade.go 中小写 error: 前缀 — implementation suggestion, addressed by medium-severity finding above
- **[low]** 建议明确清理操作在新日志文件打开之后执行 — implementation suggestion, addressed by medium-severity finding above
- **[low]** 建议重新审计所有调用点 — implementation suggestion, addressed by medium-severity finding above

## Rubric

(All dimensions: N/A — pre-revision driven by freeform findings only)
