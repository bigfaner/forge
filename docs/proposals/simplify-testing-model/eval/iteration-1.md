---
iteration: 1
date: 2026-05-17
score: 640
target: 900
---

# Proposal Evaluation Report — Iteration 1

**Document**: `docs/proposals/simplify-testing-model/proposal.md`
**Score**: 640/1000
**Target**: 900
**Outcome**: NOT MET — 260 points short

## DIMENSIONS

### 1. Problem Definition: 70/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 30/40 | The core problem (3 overlapping concepts where 2 suffice) is understandable, but the description relies on assertion ("概念边界模糊", "用户配置复杂") without quoting actual user confusion. A new reader unfamiliar with Forge v2 must infer what "profile" means from a single code snippet. |
| Evidence provided | 10/40 | The "Evidence" table lists 4 rows, but every row is the author's own assessment — no user feedback, no support tickets, no quantitative data (e.g., "X% of new users misconfigure capabilities"), no external examples. The table is structured as evidence but functions as a restatement of the problem. |
| Urgency justified | 30/30 | Clear and concrete: "v3.0.0 重构窗口期，breaking change 可接受。越晚改，迁移成本越高。" Specific version, specific reason, specific cost of delay. |

### 2. Solution Clarity: 100/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | The before/after YAML snippets, the CLI command listing, and the directory tree make the approach highly concrete. A reader could explain it back. Minor deduction: the phrase "Profile 降级为 forge 内部实现细节" is slightly ambiguous — does it mean profiles are hidden but still exist internally, or completely removed? |
| User-facing behavior described | 45/45 | Excellent. D1 shows the exact config fields, D4 shows every CLI command, D5 shows multi-language handling, and D6 shows per-skill migration. Observable behavior is fully specified. |
| Technical direction clear | 17/35 | "自动检测逻辑已存在（detect.go），只需修改输出格式" is a high-level direction but lacks specifics on how detection signals map to language keys at the code level. The D3 directory structure is helpful, but there is no description of the Go package API, no data flow diagram, and no explanation of how `embed.go` or `config.go` change. The technical direction is a sketch, not a plan. |

### 3. Industry Benchmarking: 55/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 25/40 | Three references are named: GitHub Actions auto-detect, Jest/pytest, and Terraform/ESLint. However, none are described in any depth — no links, no explanation of how each solves the analogous problem, no analysis of what Forge can borrow specifically. The "参考" column has no URLs, no citations, no detail beyond a product name. |
| At least 3 meaningful alternatives | 10/30 | Four alternatives are listed, but 3 of the 4 are straw men: "Do nothing" (rejected for v3.0 window), "渐进废弃" (rejected for not fitting v3.0), "重命名不改结构" (rejected as "治标不治本"). Each is presented with a single-sentence dismissal that makes rejection appear preordained. Per the deduction rules, each straw man incurs -20. Only the selected alternative ("去除 Profile + 语言自动检测") is meaningfully explored. Three straw men: -60 from this sub-score. |
| Honest trade-off comparison | 10/25 | The comparison table is superficial. The "Cons" column for the selected approach says only "所有 skill 需迁移" — no analysis of migration risk, no estimation of effort, no discussion of what breaks. The other approaches have 3-5 word con descriptions. Trade-offs are listed but not analyzed. |
| Chosen approach justified against benchmarks | 10/25 | The justification for choosing auto-detection over explicit declaration (a la Terraform/ESLint) is a single sentence in Innovation Highlights: "行业实践（CI 系统、IDE）都倾向于从项目结构自动推导". This is an assertion, not a justification. Why is auto-detect better for Forge specifically? What about cases where auto-detect fails? No analysis. |

### 4. Requirements Completeness: 75/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 30/40 | Five key scenarios are listed (single-language, frontend, multi-language, single-feature, mobile). However, edge cases are missing: What happens when no language is detected? What if a project has conflicting signals (e.g., both `go.mod` and `pyproject.toml`)? What about monorepos with language-specific subdirectories? The mobile scenario is mentioned but deferred. Error scenarios (detection failure, ambiguous detection) are not covered. |
| Non-functional requirements | 25/40 | No explicit NFR section exists. Performance (detection speed), security (what if detection reads sensitive files?), compatibility (existing CI pipelines using `forge profile`), and accessibility are not discussed. The "Constraints & Dependencies" section addresses backward compatibility but not other NFRs. |
| Constraints & dependencies | 20/30 | Three constraints are listed (v3.0 breaking change, embedded strategy files, migrate 6 existing strategies). However, no mention of dependency on specific Go version, file system access patterns, CLI output format stability, or downstream consumers beyond the listed skills. |

### 5. Solution Creativity: 50/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 20/40 | The proposal acknowledges this is inspired by CI auto-detect patterns. The "innovation" is applying this to a test configuration system — which is reasonable but not novel. The document says "用户选择 profile 本质上是多余的" — if it's obviously redundant, removing it is cleanup, not creativity. No differentiation from the CI auto-detect baseline is articulated. |
| Cross-domain inspiration | 20/35 | CI auto-detect and IDE language detection are cited as inspiration. These are directly adjacent domains (developer tools), not cross-domain. No inspiration from more distant fields (e.g., how package managers handle dependency resolution, how type systems infer types). |
| Simplicity of insight | 10/25 | The core insight ("language = 1:1 with test framework, so auto-detect") is clean. However, the proposal does not address the cases where this 1:1 assumption breaks (multi-framework languages like JavaScript with Jest vs Playwright vs Cypress), making the insight feel incomplete rather than elegant. |

### 6. Feasibility: 65/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 30/40 | Claims "自动检测逻辑已存在（detect.go）" and "策略目录迁移是文件重命名 + 路径更新". These are plausible but unverified assertions. No code references (line numbers, function signatures), no prototype or proof-of-concept, no assessment of the detect.go changes needed. |
| Resource & timeline feasibility | 20/30 | "1 个 Go package 重构 + ~10 个 skill 文件更新 + config schema 更新" provides a rough scope but no timeline estimate, no staffing plan, no Sprint/week breakdown. "~10" is vague — is it 8 or 14? |
| Dependency readiness | 15/30 | No assessment of external dependency readiness. Are the 6 existing strategy files in a format that supports straightforward migration? Are there any API changes needed in the CLI framework? Is the config schema library (presumably a YAML parser) flexible enough? |

### 7. Scope Definition: 65/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | Seven concrete deliverables: config schema changes, global rename, internal refactoring, CLI commands, detection logic, skill updates, schema/example updates. Each is identifiable. Minor issue: "所有消费 profile/capability 的 skill 更新" is slightly vague (which skills specifically? D6 later enumerates them, but the in-scope section itself does not cross-reference). |
| Out-of-scope explicitly listed | 22/25 | Four items explicitly deferred: new strategy packages, per-feature language narrowing, mobile deep design, strategy content modification. Good. Minor gap: documentation updates (user-facing docs, migration guide) are not explicitly in-scope or out-of-scope. |
| Scope is bounded | 15/25 | The scope is described in terms of deliverables, not in terms of time or effort. There is no "this fits in X weeks" or "this is a single sprint" framing. Without a timeline anchor, the scope is technically bounded by what is listed but unbounded in terms of when "done" occurs. |

### 8. Risk Assessment: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 22/30 | Four risks listed. Three are meaningful (multi-language misdetection, skill migration gaps, mobile detection mismatch). The fourth (user upgrade confusion) is operational, not technical. Missing risks: breaking existing CI pipelines, performance impact of file-system scanning, strategy file format incompatibilities, edge cases in language detection (monorepos, Dockerfiles, mixed-language repos). |
| Likelihood + impact rated | 18/30 | Ratings use M/H/L but no scale is defined. The ratings feel honest (not everything is "low likelihood, high impact"), but there is no justification for why multi-language misdetection is "M" likelihood — is that based on data, experience, or gut feeling? |
| Mitigations are actionable | 15/30 | Two mitigations are specific: "用户可通过 languages 覆盖" and "全局 grep profile 和 capability 确保零遗漏". The others are vague: "检测规则可配置优先级" (how? where configured?), "strategies/mobile/ 作为特殊 case，检测逻辑单独处理" (separate from what?), "v3.0 文档提供迁移指南" (what format? where linked?). |

### 9. Success Criteria: 50/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 35/55 | Six checklist items. Five are verifiable ("不再有 test-profiles 和 capabilities 字段", "forge testing detect 能正确检测所有 6 种语言", "代码库中无 profile 或 capability 残留"). However: "用户配置 project-type + interfaces 即可运行完整测试管线" is somewhat vague — what constitutes "完整测试管线"? What inputs are used for verification? The criterion about CLI output "与 v2 对应策略内容一致" is testable but does not specify the verification method. |
| Coverage is complete | 15/25 | Success criteria cover config schema, detection, CLI, and skill migration. Missing: no criterion for multi-language detection accuracy, no criterion for the `languages` override field working correctly, no criterion for the `interfaces` defaulting behavior, no criterion for documentation/migration guide completeness. |

### 10. Logical Consistency: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 25/35 | The solution (remove profile, add auto-detect) directly targets the stated problem (3 overlapping concepts). However, the problem also mentions "命名不直觉" about "capability" — the rename to "interfaces" addresses this, but the document does not explicitly trace this problem-solution link. More critically: the problem states "Profile 声明 capabilities，config.yaml 又能覆盖 capabilities。两层配置互相覆盖" — but the proposed solution still has override semantics (`languages` overrides auto-detect; `interfaces` overrides defaults), introducing a new form of the same problem. |
| Scope / Solution / Success Criteria aligned | 15/30 | The scope says "Per-feature 语言窄化" is out of scope, but the solution (D5) describes per-feature narrowing as a design decision with a partial spec. The scope says "移动端策略深度设计" is out of scope, but scenario 5 describes mobile detection and the directory structure includes a `mobile/` language directory. These create ambiguity about what is actually in vs. out of scope. Success criteria do not mention mobile at all, creating a gap with the scenarios and directory structure. |
| Requirements / Solution coherent | 15/25 | Scenario 4 (per-feature language narrowing) is listed as a key scenario but is explicitly out of scope for delivery. Scenario 5 (mobile) has a partial solution but deferred scope. These orphan requirements are acknowledged but not resolved — they exist in the requirements section without a corresponding committed solution. |

## ATTACKS

1. **Problem Definition**: Evidence is entirely self-authored assertions — "用户需要同时理解 'profile 选择' 和 'capability 裁剪' 两个独立决策" — no user feedback data, no support tickets, no metrics. Must add quantitative or qualitative external evidence.

2. **Industry Benchmarking**: Three of four alternatives are straw men with one-line dismissals — "重命名不改结构 | 改动小 | 根本问题未解决 | Rejected: 治标不治本" — each alternative must have genuine pros/cons analysis, not just a pretext for rejection.

3. **Industry Benchmarking**: No URLs, no citations, no depth for referenced industry solutions — "GitHub Actions auto-detect" is a product name, not an analysis. Must add specific references and explain what Forge borrows from each.

4. **Industry Benchmarking**: The chosen approach is justified with a single assertion — "行业实践（CI 系统、IDE）都倾向于从项目结构自动推导" — must articulate why auto-detect is specifically better for Forge's context, including failure modes.

5. **Requirements Completeness**: Edge cases missing — no scenario for "no language detected", "conflicting detection signals" (e.g., `go.mod` + `package.json`), or "monorepo with subdirectory languages". Must add error/failure scenarios to the key scenarios list.

6. **Requirements Completeness**: No non-functional requirements section — performance (file scanning overhead), security (what files are read during detection), compatibility (CI pipeline disruption) are all absent. Must add explicit NFR subsection.

7. **Solution Clarity**: Technical direction is a sketch, not a plan — "自动检测逻辑已存在（detect.go），只需修改输出格式" — must describe the Go package API changes, data flow, and config reading changes at a level sufficient for a developer to start implementation.

8. **Solution Creativity**: The 1:1 language-to-framework assumption is stated without acknowledging counter-examples — JavaScript has Jest, Playwright, Cypress, Vitest, Mocha. The "innovation" of auto-detect must address this gap or acknowledge the limitation.

9. **Feasibility**: No timeline, no staffing, no Sprint breakdown — "1 个 Go package 重构 + ~10 个 skill 文件更新" is a scope description, not a feasibility assessment. Must estimate weeks/effort and validate team capacity.

10. **Feasibility**: Dependency readiness not assessed — are the 6 existing strategy files compatible with the new directory structure? What changes to the embedding mechanism? What CLI framework changes? Must audit actual code before claiming feasibility.

11. **Risk Assessment**: Missing significant risks — CI pipeline breakage for existing users, file-system scanning performance on large repos, strategy file format incompatibilities. Must expand the risk table.

12. **Risk Assessment**: Vague mitigations — "检测规则可配置优先级" (how?), "strategies/mobile/ 作为特殊 case，检测逻辑单独处理" (what does this mean in code?). Each mitigation must be specific enough to become a task.

13. **Logical Consistency**: The problem complains about "两层配置互相覆盖" but the solution introduces new override semantics (`languages` overrides auto-detect, `interfaces` overrides defaults). Must address whether the new override chain is cleaner or just different.

14. **Logical Consistency**: Per-feature narrowing and mobile detection appear in scenarios/directory structure but are marked out-of-scope — must either commit them to scope or remove from scenarios and design decisions.

15. **Success Criteria**: Missing criteria for `languages` override, multi-language detection accuracy, `interfaces` defaulting, and documentation/migration guide. Must add criteria that cover all in-scope deliverables.

16. **Scope Definition**: No timeline anchor — the scope lists deliverables but says nothing about when "done" is. Must add a time-box or Sprint framing to make scope genuinely bounded.
