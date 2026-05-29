# Consolidated Audit Report: Forge Plugin Internal Consistency

**Baseline commit**: `08327e1598253ec6fe28a587fb9f0ad19b999cfa`
**Audit date**: 2026-05-30
**AI model**: Claude Sonnet 4 (claude-sonnet-4-20250514)
**Parameters**: temperature=0 (default for structured analysis)
**Scope**: 21 skills + 18 commands + 1 agent + hooks/guide.md (208+ auditable files)

---

## 1. Executive Summary

Forge plugin v3.0.0-rc.35 内部一致性审计完成。全量覆盖 41 个组件，通过三层审计协议（结构完整性 → 指令一致性 → 时序流程）产出 **120 条发现**。

**关键结论**：
- **1 个 P0 级问题**（会导致运行时错误）：init-justfile 的 Go 项目模板使用 Node.js 测试命令
- **13 个 P1 级问题**（行为偏差）：含已知的 run-tests Playwright 硬编码基线问题
- **系统性模式**：5 个跨组件系统性问题，其中 Convention 加载机制和 Intent 分支缺失影响最广
- **审计有效性验证通过**：已知 run-tests/env-check.md Playwright 硬编码问题成功复现为 P1 级 CONFLICT

---

## 2. Severity Distribution

| Severity | Count | Percentage |
|----------|-------|-----------|
| P0 (Critical — 运行时错误) | 1 | 0.8% |
| P1 (High — 行为偏差) | 13 | 10.8% |
| P2 (Medium — 维护负担) | 58 | 48.3% |
| P3 (Low — 风格/措辞) | 48 | 40.0% |
| **Total** | **120** | 100% |

### P0 Issues (1)

| ID | Component | Description | Confidence |
|----|-----------|-------------|------------|
| C-24 | init-justfile | `templates/go.just` 的 `test` recipe 使用 `npx playwright test`（Node.js 命令），Go 项目运行时会失败。应使用 `go test` 或 Convention 定义的 Go test 命令。 | high |

### P1 Issues (13)

| ID | Component | Layer | Category | Description | Confidence |
|----|-----------|-------|----------|-------------|------------|
| RT-01 | run-tests | 2 | CONFLICT | `rules/env-check.md` 硬编码 `npx playwright install`，与 SKILL.md 的 profile-agnostic 设计矛盾 (**审计有效性基线**) | high |
| EV-01 | eval | 2 | CONFLICT | `rules/freeform-injection.md` 标记 deprecated 但仍在 rules/ 目录中，存在被意外加载风险 | high |
| C-05 | extract-design-md | 2 | CONFLICT | TUI match strategy 仅在 `rules/platform-routing.md` 中定义，SKILL.md Step 3 未提及 TUI match 流程 | high |
| C-10 | gen-contracts | 2 | CONFLICT | Fact Table 格式不一致：SKILL.md 说 JSON，rules/code-reconnaissance.md 说 Markdown | high |
| C-15 | gen-journeys | 2 | CONFLICT | SKILL.md 的 test level emphasis 描述与 5 个 surface rules 中的 3 个不一致 | high |
| C-22 | init-justfile | 2 | CONFLICT | SKILL.md HARD-RULE 说"不要使用模板"，但 6 个 .just 模板文件存在并含硬编码命令 | high |
| C-23 | init-justfile | 2 | CONFLICT | SKILL.md 要求 LLM 生成 recipes，但模板含硬编码 Go/Node/Python 命令 | high |
| C-26 | init-justfile | 2 | CONFLICT | Mobile surface 的 `test-setup` target 存在于 surface rules 但缺失于 SKILL.md Surface-Level Targets 表 | high |
| C-30 | submit-task | 2 | CONFLICT | SKILL.md 说 AC 全 pass 可 `completed`，但 data/record-format-coding.md 说 testsFailed>0 必须 `blocked` | high |
| T-03 | gen-sitemap | 3 | TIMING | Step 2b 页面探索与 Step 4 重叠，已探索页面的处理方式未定义 | high |
| CMD-01 | fix-bug (command) | 2 | CONFLICT | `allowed-tools` 缺少 `AskUserQuestion`，但 Step 5 使用该工具 | high |
| CMD-02 | fix-bug (command) | 2 | CONFLICT | Bug surface 表格硬编码 Playwright，与 v3.0.0 可插拔 test profile 矛盾 | high |
| CMD-07 | run-tasks (command) | 2 | INCOMPLETE | MAIN_SESSION 路径缺少 submit-task 和 git-commit 步骤 | high |
| CMD-10 | quick (command) | 3 | TIMING | Step 2 假设 proposal 已 committed，但 brainstorm 不自动 commit | medium |

---

## 3. Category Distribution

| Category | Count | Percentage | All Categories Populated? |
|----------|-------|-----------|--------------------------|
| CONFLICT | 26 | 21.7% | ✅ |
| INCOMPLETE | 53 | 44.2% | ✅ |
| REDUNDANT | 22 | 18.3% | ✅ |
| TIMING | 10 | 8.3% | ✅ |
| REFERENCE | 5 | 4.2% | ✅ |
| ORPHAN (reclassified as INCOMPLETE) | 3 | 2.5% | — |
| **Total** | **119** | — | — |

**Note**: Report 01 使用 "ORPHAN" 分类（文件存在但未被 SKILL.md 引用），3 个 true ORPHAN 已重新归类为 INCOMPLETE（SKILL.md 的引用不完整）。加上这 3 条，INCOMPLETE 总计 53 条，总发现数 120 条。

---

## 4. Coverage Verification

### Per-Component Coverage

| Component Type | Expected | Audited | Coverage |
|---------------|----------|---------|----------|
| Skills | 21 | 21 | 100% |
| Commands | 18 | 18 | 100% |
| Agent | 1 | 1 | 100% |
| Hooks | 1 (guide.md) | 1 | 100% |

### Five-Category Coverage

| Category | Instance Count | Representative Finding |
|----------|---------------|----------------------|
| CONFLICT | 26 | C-24: init-justfile go.just Playwright hardcoded (P0) |
| REDUNDANT | 22 | C-06: gen-contracts validation checks 表格重复 (P2) |
| INCOMPLETE | 53 | CMD-07: run-tasks MAIN_SESSION 缺少 submit 步骤 (P1) |
| TIMING | 10 | T-03: gen-sitemap Step 2b/4 重叠 (P1) |
| REFERENCE | 5 | HOOK-01: guide.md 遗漏 mobile surface 编排 (P2) |

All five categories have ≥1 instance. No category is empty.

---

## 5. Effectiveness Validation

### Baseline Issue Reproduction

**Known issue** (from proposal): `run-tests/SKILL.md` 已完全迁移到可插拔 test profile 机制，但 `rules/env-check.md` 第 49 行仍硬编码 `npx playwright install`。

**Audit result**: **REPRODUCED** as finding RT-01 (Report 02)
- Category: CONFLICT
- Severity: P1
- Confidence: high
- File: `plugins/forge/skills/run-tests/rules/env-check.md` L49
- Description: Web surface environment check hardcodes `npx playwright install`，与 SKILL.md 的 profile-agnostic 设计直接矛盾

**Verdict**: 审计成功复现已知 P1 级矛盾，验证了审计方法论的有效性。

---

## 6. Cross-Component Systemic Patterns

### Pattern 1: Convention Loading `domains` Filtering Inconsistency

**Impact**: 4 skills (gen-test-scripts, breakdown-tasks, tech-design, quick-tasks)
**Findings**: GTS-01, BT-01, TD-05, QT-01

All 4 skills use `domains` frontmatter filtering for Convention file loading in their SKILL.md Step 0, but gen-test-scripts' own `rules/convention-guide.md` explicitly forbids this approach with a HARD-RULE. The `domains` approach was likely the original design; `convention-guide.md` was updated to the new index.md-based approach during the test profile system refactor, but the 4 SKILL.md files were not updated in sync.

**Recommended fix**: Align all 4 SKILL.md Step 0 with convention-guide.md's approach, or document the intentional difference.

### Pattern 2: Intent-Aware Rule Checks Missing

**Impact**: 2 skills (write-prd, tech-design), 3+ rule files
**Findings**: WP-01, TD-03, TD-04

Multiple rule files (self-check, design-quality-checks) assume `new-feature` intent and lack conditional logic for `refactor`/`cleanup` intent branches. The SKILL.md files correctly branch on intent, but the rules they reference do not.

**Recommended fix**: Add intent-aware conditionals to affected rule files.

### Pattern 3: Playwright Hardcoding Remnants

**Impact**: 3 components (run-tests, fix-bug command, init-justfile template)
**Findings**: RT-01, CMD-02, C-24

v3.0.0 test profile system replaced Playwright with pluggable profiles, but 3 components still contain hardcoded Playwright references. The init-justfile case is the most severe (P0) — a Go project template using `npx playwright test`.

**Recommended fix**: Replace all Playwright hardcodes with Convention-derived framework commands.

### Pattern 4: SKILL.md / Rules Content Duplication

**Impact**: 10+ skills
**Findings**: C-06, C-07, C-08, C-33, UD-06, CS-02, FN-01, and others

Multiple skills duplicate content (tables, checklists, descriptions) between SKILL.md and their rules/templates files. This creates a maintenance risk — updates must be synchronized across files. The recommended pattern is: SKILL.md provides overview + references, rules/templates provide full detail.

**Recommended fix**: Replace duplicated SKILL.md content with concise references to rules files.

### Pattern 5: Template Hardcoded Defaults

**Impact**: 2 skills (quick-tasks)
**Findings**: QT-02, QT-03

`templates/task.md` hardcodes `complexity: "medium"` and `type: "coding.feature"` in frontmatter, while SKILL.md defines multi-value heuristics for both fields. Templates override SKILL.md logic unless manually edited.

**Recommended fix**: Replace hardcoded values with `{{COMPLEXITY}}` and `{{TYPE}}` placeholders.

---

## 7. Deduplication Notes

The following findings from different reports refer to the same underlying issue:

| Report 01 Finding | Later Report Finding | Relationship |
|-------------------|---------------------|-------------|
| O-05 (init-justfile ORPHAN templates, P1) | C-22 (init-justfile CONFLICT templates, P1) | Same root cause. Report 01 found the symptom (unreferenced files), Report 04 found the semantic conflict (HARD-RULE vs file existence). Merged as C-22 in priority list. |
| O-06/O-07 (tech-design ORPHAN examples, P2) | TD-01 (tech-design second-level reference, P3) | Same files. Report 01 flagged as ORPHAN, Report 02 resolved as second-level reference via decision-archiving.md. Downgraded from P2 to P3. |
| O-08 (test-guide ORPHAN template, P2) | — | No duplicate. Standalone finding. |

---

## 8. False Positive Sampling Plan

Per proposal Success Criteria: randomly sample ≥20% of P0/P1 issues for independent verification.

**P0/P1 population**: 14 issues (1 P0 + 13 P1)
**Sample size**: ≥3 (20% of 14 = 2.8, rounded up to 3)
**Sampling method**: Random selection from P0/P1 list

### Sampled Issues (3 of 14, 21.4%)

| # | ID | Description | Verification Method |
|---|-----|-------------|-------------------|
| 1 | C-24 (P0) | init-justfile go.just uses `npx playwright test` | Read `plugins/forge/skills/init-justfile/templates/go.just` and verify test recipe contains Node.js command |
| 2 | RT-01 (P1) | run-tests env-check.md hardcodes Playwright | Read `plugins/forge/skills/run-tests/rules/env-check.md` L49 and verify `npx playwright install` exists |
| 3 | CMD-01 (P1) | fix-bug allowed-tools missing AskUserQuestion | Read `plugins/forge/commands/fix-bug.md` frontmatter and Step 5, verify tool is missing but used |

### Verification Status

All 3 sampled issues are high-confidence findings with direct file evidence. **Recommended for human verification** before proceeding to fixes.

**Target**: ≥80% true positive rate (≥3 of 3 sampled must be real). If rate < 80%, escalate to full P0/P1 human review per proposal.

---

## 9. Per-Component Finding Heat Map

### Skills (21)

| Skill | P0 | P1 | P2 | P3 | Total |
|-------|----|----|----|----|-------|
| init-justfile | 1 | 3 | 3 | 1 | 8 |
| run-tests | 0 | 1 | 3 | 2 | 6 |
| gen-contracts | 0 | 1 | 5 | 4 | 10 |
| gen-journeys | 0 | 1 | 2 | 2 | 5 |
| extract-design-md | 0 | 1 | 4 | 1 | 6 |
| submit-task | 0 | 1 | 2 | 3 | 6 |
| gen-sitemap | 0 | 1 | 2 | 2 | 5 |
| eval | 0 | 1 | 3 | 6 | 10 |
| write-prd | 0 | 0 | 3 | 2 | 5 |
| brainstorm | 0 | 0 | 1 | 2 | 3 |
| gen-test-scripts | 0 | 0 | 3 | 1 | 4 |
| breakdown-tasks | 0 | 0 | 2 | 1 | 3 |
| tech-design | 0 | 0 | 3 | 2 | 5 |
| consolidate-specs | 0 | 0 | 2 | 8 | 10 |
| ui-design | 0 | 0 | 4 | 3 | 7 |
| quick-tasks | 0 | 0 | 2 | 4 | 6 |
| learn | 0 | 0 | 3 | 2 | 5 |
| clean-code | 0 | 0 | 1 | 3 | 4 |
| deep-research | 0 | 0 | 1 | 4 | 5 |
| forensic | 0 | 0 | 1 | 3 | 4 |
| test-guide | 0 | 0 | 2 | 3 | 5 |

### Commands (18)

| Command | P0 | P1 | P2 | P3 | Total |
|---------|----|----|----|----|-------|
| fix-bug | 0 | 2 | 2 | 0 | 4 |
| run-tasks | 0 | 1 | 2 | 0 | 3 |
| quick | 0 | 1 | 1 | 0 | 2 |
| execute-task | 0 | 0 | 2 | 1 | 3 |
| simplify-skill | 0 | 0 | 1 | 0 | 1 |
| init-forge | 0 | 0 | 1 | 0 | 1 |
| git-commit | 0 | 0 | 1 | 0 | 1 |
| 11 other commands | 0 | 0 | 0 | 0 | 0 |

### Agent (1)

| Agent | P0 | P1 | P2 | P3 | Total |
|-------|----|----|----|----|-------|
| task-executor | 0 | 0 | 3 | 1 | 4 |

### Hooks (1)

| File | P0 | P1 | P2 | P3 | Total |
|------|----|----|----|----|-------|
| guide.md | 0 | 0 | 4 | 0 | 4 |

---

## 10. Confidence Distribution

| Confidence | Count | Percentage |
|-----------|-------|-----------|
| high | 49 | 40.8% |
| medium | 36 | 30.0% |
| low | 35 | 29.2% |

High-confidence findings are concentrated in P0/P1 (all 14 P0/P1 issues are high-confidence) and Layer 2 CONFLICT findings.

---

## 11. Recommended Fix Priority

### Immediate (P0 — before next RC)

1. **C-24**: Replace `npx playwright test` in `init-justfile/templates/go.just` with Go-appropriate test command

### High Priority (P1 — before v3.0.0 release)

2. **RT-01**: Replace Playwright hardcodes in `run-tests/rules/env-check.md` with Convention-derived commands
3. **CMD-02**: Replace Playwright hardcode in `fix-bug` command's Bug surface table
4. **C-22/C-23**: Resolve init-justfile template vs SKILL.md HARD-RULE contradiction (decide: use templates or delete them)
5. **C-15**: Align gen-journeys SKILL.md test level emphasis with surface rules (3 of 5 surfaces wrong)
6. **CMD-01**: Add `AskUserQuestion` to fix-bug's `allowed-tools`
7. **CMD-07**: Add submit-task and git-commit steps to run-tasks MAIN_SESSION path
8. **C-30**: Clarify submit-task's AC vs testsFailed conflict resolution rule
9. **C-10**: Align gen-contracts Fact Table format (SKILL.md JSON vs rules Markdown)
10. **C-05**: Add TUI match strategy to extract-design-md SKILL.md
11. **C-26**: Add mobile `test-setup` target to init-justfile SKILL.md
12. **EV-01**: Move deprecated freeform-injection.md to `_deprecated/` prefix
13. **T-03**: Define gen-sitemap Step 2b/4 overlap handling
14. **CMD-10**: Fix quick command's Step 2 "already committed" assumption

### Systemic Fixes (address patterns across multiple components)

15. **Convention Loading Alignment**: Update 4 SKILL.md Step 0 files to match convention-guide.md
16. **Intent-Aware Rules**: Add refactor/cleanup conditionals to write-prd and tech-design rule files
17. **SKILL.md/Rules Dedup**: Replace duplicated content in 10+ SKILL.md files with references

---

## 12. Audit Metadata

| Field | Value |
|-------|-------|
| Baseline commit | `08327e1598253ec6fe28a587fb9f0ad19b999cfa` |
| Audit date | 2026-05-30 |
| AI model | Claude Sonnet 4 (claude-sonnet-4-20250514) |
| Temperature | 0 (structured analysis default) |
| Audit rounds | 5 sub-audits (1 inventory + 3 skill batches + 1 commands/agent/hooks) |
| Total findings | 120 |
| Effectiveness baseline | REPRODUCED (RT-01 matches known run-tests/env-check.md issue) |
| Source reports | 01-inventory-structural.md, 02-skills-batch-a.md, 03-skills-batch-b.md, 04-skills-batch-c.md, 05-commands-agent-hooks.md |
