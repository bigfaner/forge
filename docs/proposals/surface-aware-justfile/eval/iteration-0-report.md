---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report (Iteration 0)

## ATTACK_POINTS

- **[high]** 规则文件双重职责的职责边界不清晰 | quote: "每个 surface 类型的规则文件定义两部分内容：编排序列（该 surface 类型的测试执行流程）和 just 配方调用契约（序列中每个步骤调用的 just 配方名、参数、退出码语义）" | improvement: 明确规则文件物理归属和职责分离机制

- **[high]** 规则文件物理归属未澄清——共享文件还是独立副本 | quote: "将 init-justfile（配方生产者）和 run-tests（调度器）的 surface 感知通过规则文件统一设计" | improvement: 澄清规则文件是物理共享还是逻辑同构但物理独立

- **[high]** npm wrapper 进程退出导致 PID 指向已死进程而非实际 dev server | quote: "每次 probe 重试失败后，检查 .forge/dev-server.pid 中记录的进程是否仍然存活。如果进程已死，probe 立即退出并报告'dev server 崩溃'，不等待剩余重试" | improvement: 在 teardown 中增加端口反查机制定位实际进程 PID

- **[medium]** 调度器模式同构性声明掩盖了 gen-test-cases 与 run-tests 的本质差异 | quote: "run-tests 不自行推断编排模式，也不依赖跨 skill 文件传递编排参数。它遵循与 gen-test-cases 相同的调度模式：1. 检测 surface type 2. 加载执行策略 3. 按策略执行" | improvement: 重新审视同构性声明，避免在共享抽象层时过度泛化

- **[medium]** scope 迁移原子提交涉及 7+ 组件，review 认知负荷极高 | quote: "阶段 1-4 必须在同一提交中完成。如果分批提交，中间状态会导致 scope-assignment 输出 admin-panel 但 init-justfile 仍期望 frontend 的不一致" | improvement: 将同一提交约束弱化为同一 PR 允许逻辑分组多提交

- **[medium]** 兼容层在 surfaces map 中找不到对应类型时行为未定义 | quote: "若 scope 值为旧枚举（frontend/backend），在 surfaces map 中查找对应的 surface key 并映射（frontend → 找到 type=web 的 key，backend → 找到 type=api 的 key）" | improvement: 定义类型完全不存在时的兼容层回退行为

- **[medium]** macOS ps 命令截断导致 teardown 命令行匹配失败 | quote: "Linux/macOS：读取 /proc/<pid>/cmdline（Linux）或 ps -p <pid> -o command=（macOS），确认命令行包含预期的 dev server 关键词" | improvement: 改用 ps -o command= -w 禁用截断，关键词匹配同时覆盖用户命令和实际进程路径

- **[medium]** HARD-GATE 规则作为文本指令而非运行时约束，可被 LLM 忽略 | quote: "HARD-GATE 规则：probe 失败后禁止重试 probe 或重试 dev——唯一允许的下一步是执行 teardown 后中止" | improvement: 为 HARD-GATE 规则增加运行时层面的强制保障机制

- **[medium]** 后台启动退出码不可靠，probe 依赖 PID 存活检查存在竞态 | quote: "just dev 后台启动 dev server 后配方本身以 exit 0 返回（因为后台进程已分离），但 dev server 可能在启动后的几秒内崩溃" | improvement: 增强 PID 追踪机制覆盖 npm fork 子进程场景

- **[medium]** 多 scope 混合项目 probe 等待时间线性叠加，用户体验差 | quote: "顺序启动策略：just dev 配方体按顺序启动各 scope 的 dev server（先启动后端，再启动前端），每个启动后将 PID 写入 .forge/dev-server.<scope>.pid" | improvement: 在 probe 重试中区分连接被拒绝和连接超时，提供更快失败反馈

- **[medium]** 新增 surface 类型声称无需更新 Go 代码但未验证 config 校验逻辑 | quote: "新增 surface 类型只需两步：1. 在 init-justfile 的 rules/surfaces/ 下新增 <type>.md 2. 在 run-tests 的 rules/surfaces/ 下新增 <type>.md。无需修改 config schema、无需新增中间文件、无需更新 Go 代码" | improvement: 验证 Go 代码是否校验 surface 类型值域，定义 run-tests 加载不到规则文件时的行为

- **[low]** Windows Get-CimInstance 失败时缺少 fallback 行为定义 | quote: "通过 PowerShell Get-CimInstance Win32_Process -Filter \"ProcessId=<pid>\" | Select-Object CommandLine 校验命令行" | improvement: 定义 Get-CimInstance 失败时的 fallback 行为

- **[low]** 端口冲突场景下 probe 超时等待过长缺乏早期反馈 | quote: "probe 超时作为端口冲突的最终兜底检测机制" | improvement: 提供 probe 进度输出和早期退出条件

- **[low]** probe 顺序定义为先 api 后 web 但等待时间无优化 | quote: "probe 顺序：run-tests 先 probe 后端（api），再 probe 前端（web）。后端就绪是前端可用的前提条件" | improvement: 在后端 probe 完成后前端 probe 期间提供进度反馈

## BORDERLINE_FINDINGS

（无）

## SKIPPED_FINDINGS

（无）

## Classification Audit

- Factual correction: 0
- Structural suggestion: 14 (规则文件归属、进程管理可靠性、迁移策略、调度器同构性等)
- Subjective preference: 0
