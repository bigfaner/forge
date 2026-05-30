---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Eval Report (Iteration 0)

## ATTACK_POINTS

- **[high]** Phase 2 将包重组、死代码删除、魔法值提取、重复消除四类变更混合在同一阶段，blast radius 不可隔离 | quote: "全面重新设计 internal/cmd/ 和 pkg/ 两层包结构，同时彻底清除所有已识别的死代码和魔法值" | improvement: 将 Phase 2 拆分为 Phase 2a（死代码删除）→ Phase 2b（魔法值提取）→ Phase 2c（包结构重组），按 blast radius 从小到大排序

- **[high]** test-bridge 别名函数被错误归类为"死代码"，删除将破坏测试编译 | quote: "删除所有死代码：deprecated Scope 字段、别名函数、兼容层、构建产物（.out 文件）" | improvement: 将 test-bridge 清理从"死代码"中分离为独立 Scope 项，区分纯粹重导出（可直接删除）和内部函数导出（需评估迁移策略）

- **[high]** Phase 1 规范产出目标模糊——"提炼现有模式"产出描述性文档，Phase 2 需要规范性目标态 | quote: "分析现有代码库模式，扩展 docs/conventions/ 下的规范文件" | improvement: 明确 Phase 1 产出必须包含目标态定义（而非仅现状描述），并要求偏差分析和合规迁移路径

- **[high]** 未检查 monorepo 内是否存在跨 Go 模块的 import 依赖 | quote: "无外部依赖。所有涉及的包都是 forge-cli 内部包。" | improvement: 添加跨模块依赖审计作为 Phase 2 前置条件

- **[medium]** SC-5 "零个顶层命令文件"定义模糊——root.go、output.go、surfaces.go 是基础设施而非命令文件 | quote: "SC-5: internal/cmd/ 下零个顶层命令文件（所有命令均已子包化）" | improvement: 明确列举豁免文件（root.go、output.go 等），或重新定义 SC-5 仅计数命令实现文件

- **[medium]** 缺少当前包到目标包的具体映射表 | quote: "按领域合并小包，明确每个包的职责边界" | improvement: 在 Scope 中添加当前包→目标包映射表

- **[medium]** Evidence 部分事实性错误——"tests/results/raw-output.txt 在 quality_gate.go 中出现 7 次"实际为 2 次（其余在测试文件） | quote: "tests/results/raw-output.txt 在 quality_gate.go 中出现 7 次" | improvement: 修正为准确计数，区分生产代码和测试文件

- **[medium]** 引用过时工具 gorename | quote: "工具链（gorename、IDE refactor）支持良好" | improvement: 替换为 gopls 或 IDE 内置重构

- **[medium]** quality_gate.go（1067 行）等巨型文件只移动不拆分，重组只是藏问题 | quote: "重组 internal/cmd/ 包结构" | improvement: 添加文件行数上限成功标准，或在 Scope 中明确巨型文件拆分计划

## BORDERLINE_FINDINGS

- 无回滚计划——"不保留兼容层"但无回退机制（介于结构性建议和补充风险之间）

## SKIPPED_FINDINGS

- Subjective: 建议增加 golangci-lint 作为硬性门控成功标准（属于实施细节，Scorer 周期会评估）
- Subjective: 建议增加文件行数成功标准的具体阈值（属于 SC 写法偏好）

## rubric
(all dimensions): N/A
