---
name: eval-harness
description: Harness health evaluation centered on three core principles - Design Environment, Clarify Intent, Build Feedback Loops. Based on OpenAI's harness engineering practices.
---

# Eval Harness

评估项目 harness 健康状态，基于 OpenAI Harness Engineering 三大核心原则。

## When to Use

**Trigger:**
- User asks to "evaluate harness" or "check harness health"
- User provides `/eval-harness` command
- Periodic health check (recommended: weekly)

**Skip:**
- User wants to implement improvements directly (use `/improve-harness`)

## Core Framework

```
┌─────────────────────────────────────────────────────────────┐
│                    HARNESS ENGINEERING                       │
├─────────────────────────────────────────────────────────────┤
│   设计环境 (Environment) → 明确意图 (Intent) → 反馈回路 (Feedback) │
│         │                      │                    │        │
│         ▼                      ▼                    ▼        │
│   结构、边界、工具、抽象    计划、规格、原则、约束   验证、修正、可观测  │
└─────────────────────────────────────────────────────────────┘
```

## Workflow

```
1. 检查环境 ──▶ 2. 检查意图 ──▶ 3. 检查反馈 ──▶ 4. 生成报告
```

### Step 1: Design Environment (设计环境)

| 维度 | 检查项 | 命令 |
|------|--------|------|
| 渐进式披露 | CLAUDE.md < 100行, 有 paths frontmatter | `wc -l CLAUDE.md` |
| 架构边界 | 依赖方向定义, 边界检查脚本 | `ls scripts/lint*.sh 2>/dev/null` |
| 工具抽象 | skills, agents 存在 | `ls .claude/skills/ .claude/agents/ 2>/dev/null` |
| Agent可读性 | 架构在代码中可见 | 项目特定命令 |

### Step 2: Clarify Intent (明确意图)

| 维度 | 检查项 | 命令 |
|------|--------|------|
| 黄金原则 | 原则有强制机制 | 检查 CLAUDE.md 或项目规则 |
| 计划工件 | 任务有 schema 验证 | `ls docs/features/*/tasks/index.json 2>/dev/null` |
| 不变量 | Lint, format 自动化 | 检查 Makefile 或 CI 配置 |

### Step 3: Build Feedback Loops (构建反馈回路)

| 维度 | 检查项 | 命令 |
|------|--------|------|
| 即时反馈 | hooks 配置 | `cat .claude/settings.json 2>/dev/null` |
| 自动修复 | error-fixer agent | `ls .claude/agents/error-fixer.md 2>/dev/null` |
| 债务GC | 代码简化 skill | `grep -r simplify .claude/ 2>/dev/null` |
| 可观测性 | 测试结果结构化 | 项目特定命令 |

### Step 4: Generate Report

输出到 `docs/harness-reports/YYYY-MM-DD.md`，使用 `templates/report.md` 模板。

## Language Detection

自动检测项目语言并使用对应的检查命令：

| 语言 | 检测标记 | 格式检查 | Lint | 测试 |
|------|----------|----------|------|------|
| Go | `go.mod` | `gofmt` | `go vet`, `golangci-lint` | `go test` |
| Node.js | `package.json` | `prettier` | `eslint` | `npm test` |
| Python | `pyproject.toml` | `black` | `ruff`, `pylint` | `pytest` |
| Rust | `Cargo.toml` | `rustfmt` | `clippy` | `cargo test` |
| Java | `pom.xml` | `google-java-format` | `checkstyle` | `mvn test` |

## Grading Rules

### Sub-dimension

| Grade | Condition |
|-------|-----------|
| A | All checks pass with automation |
| B | Most pass, minor manual intervention |
| C | Defined but enforcement incomplete |
| F | Critical checks missing |

### Overall

| Grade | Condition |
|-------|-----------|
| A | All 3 dimensions A/B, at least 2 A's |
| B | No F, max 2 C's |
| C | 1 F or 3+ C's |
| D | 2 F's |
| F | 3+ F's or core capability missing |

## Output

After evaluation:

1. Create report: `docs/harness-reports/YYYY-MM-DD.md`
2. Update link in `docs/HARNESS-EVALUATION.md`
3. List priority improvements (P0/P1/P2)

## Related

- `/improve-harness` - Implement improvements from report
- `docs/HARNESS-EVALUATION.md` - Current evaluation summary
- `docs/harness-reports/` - Historical reports
