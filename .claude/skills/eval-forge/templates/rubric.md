# Forge Plugin Runtime Reliability Rubric

**Total: 1000 points**
**Report template:** `.claude/skills/eval-forge/templates/report.md`

## What This Rubric Measures

Runtime reliability of the forge plugin — not just structural consistency, but whether components work correctly together at runtime. Measures workflow completeness, bypass resistance, instruction precision, redundancy, reference integrity, and structural conventions.

## Scoring Methodology: 4-Phase Process

This rubric is designed for a 4-phase evaluation methodology:

**Phase 1: Construct Workflow Graph (D1)**
1. Read the ground-truth workflow specs embedded in Dimension 1 below
2. Scan all skill/command/agent files, extract actual prerequisites, outputs, gate points
3. Compare spec vs actual, find breakpoints, dead ends, unreachable states

**Phase 2: Per-Node Adversarial Testing (D2)**
1. List every gate/confirm point per node
2. For each gate, assume you are a lazy agent — ask "how can I bypass this?"
3. Check whether HARD-RULE has enforceable consequences (not just words)
4. Eval loops: check whether scoring must be done by independent subagent
5. Quality gates: check whether CLI-level enforcement exists

**Phase 3: Per-File Precision Review (D3 + D4)**
1. Instruction conflicts (highest priority): cross-file search for contradictory descriptions of the same concept
2. Step ambiguity: does each SKILL.md step have a single unambiguous interpretation
3. Incomplete conditionals: does every if-then have an else
4. Undefined variables: do agent-filled variables have source annotations
5. Content redundancy: check by A/B/C category standards

**Phase 4: Baseline Integrity (D5 + D6)**
1. Reference integrity
2. Frontmatter, eval templates, name alignment

## Scoring Dimensions

| Dimension | Points |
|-----------|--------|
| 1. Workflow Completeness | 250 |
| 2. Bypass Resistance | 250 |
| 3. Instruction Precision | 200 |
| 4. Cross-file Dedup | 150 |
| 5. Reference Integrity | 100 |
| 6. Structural Convention | 50 |
| **Total** | **1000** |

Score allocation logic: Runtime reliability (D1+D2+D3) = 700 pts (70%), Information efficiency (D4) = 150 pts, Baseline integrity (D5+D6) = 150 pts.

## Deduction Tiers

| Severity | Penalty | Examples |
|----------|---------|---------|
| Critical | -20 | Workflow breakpoint, eval integrity bypass, instruction conflict |
| High | -15 | Dangling reference, invalid manifest transition, missing eval template |
| Medium | -10 | Step ambiguity, conditional without else, content duplication, missing prerequisite check |
| Low | -5 | Missing frontmatter field, name-directory mismatch, advisory-only prohibition |

**Floor rule:** Each sub-criterion has a minimum score of 0. Total deductions for a sub-criterion cannot result in a negative score.

**Cross-dimension deduction:** The same issue may be deducted in multiple dimensions if it violates independent criteria (e.g., an advisory-only HARD-RULE is both a D2 bypass vector and a D3 incomplete conditional). Bypass resistance (D2) and instruction precision (D3) are independent dimensions.

## Dimensions

### 1. Workflow Completeness (250 pts)

Verify the full-chain state graph: every skill's prerequisites, outputs, and successor steps; both quick mode and full mode paths must be complete.

**Ground-Truth Workflow Specs:**

#### Full Mode Pipeline

```
brainstorm → write-prd → eval-prd → [has UI?] → ui-design → eval-ui → prototype → [human review] → tech-design → [db-schema?] → [schema review] → eval-design → breakdown-tasks
```

```
breakdown-tasks → forge task index (auto T-test tasks) → T-test-1 (gen-sitemap + gen-test-cases) → T-test-1b (eval-test-cases) → T-test-2 (gen-test-scripts) → T-test-3 (run-e2e-tests) → T-test-4 (graduate-tests) → T-test-4.5 (verify-regression) → T-test-5 (consolidate-specs)
```

#### Quick Mode Pipeline

```
/quick → /brainstorm → [human confirm] → /quick-tasks → /run-tasks
```

Quick test chain: T-quick-1 (gen-test-cases) → T-quick-2 (gen-test-scripts) → T-quick-3 (run-e2e-tests) → T-quick-4 (graduate-tests) → T-quick-5 (verify-regression). Skips: gen-sitemap, eval-test-cases, consolidate-specs.

#### Manifest Status Machine

```
prd → design → tasks → in-progress → completed
```

Legal transitions only forward. Quick mode starts at tasks (no prd/design).

#### Per-Skill Precondition/Output Matrix

| Skill | Hard Prerequisites | Outputs | Conditionals |
|-------|-------------------|---------|-------------|
| brainstorm | None | `docs/proposals/<slug>/proposal.md` | → eval-proposal (optional) or write-prd |
| write-prd | Optional: proposal.md, sitemap.json | `prd/prd-spec.md`, `prd/prd-user-stories.md`, `prd/prd-ui-functions.md` (if UI), `manifest.md` (status: prd) | has UI → ui-design; no UI → tech-design |
| eval-prd | `prd/prd-spec.md`, `prd/prd-user-stories.md` | `prd/eval/iteration-{N}.md`, `prd/eval/report.md` | score gate pass/fail; delegates to `skills/eval` with `rubrics/prd.md` |
| ui-design | `prd/prd-ui-functions.md` (hard) | `ui/ui-design.md`, `ui/prototype/` | multi-platform → separate files |
| eval-ui | `ui/ui-design.md` | `ui/eval/iteration-{N}.md`, `ui/eval/report.md` | platform → rubric variant; delegates to `skills/eval` with `rubrics/ui-{platform}.md` |
| tech-design | `prd/prd-spec.md` (hard) | `design/tech-design.md`, `design/er-diagram.md` + `design/schema.sql` (if db), `manifest.md` (status: design) | db-schema: yes → mandatory ER+schema |
| eval-design | `design/tech-design.md` | `design/eval/iteration-{N}.md`, `design/eval/report.md` | score gate; delegates to `skills/eval` with `rubrics/design.md` |
| breakdown-tasks | `prd/prd-spec.md` + `design/tech-design.md` (both hard) | `tasks/*.md`, `tasks/index.json`, `manifest.md` (status: tasks) | HAS_UI/NO_UI/HAS_DB/HAS_PLACEMENT tags |
| gen-test-cases | `prd/prd-user-stories.md` + `prd/prd-spec.md` (both hard) | `testing/test-cases.md` | profile → interface types |
| eval-test-cases | `testing/test-cases.md` + PRD docs | `testing/eval/iteration-{N}.md` | Step Actionability < 200 blocks; delegates to `skills/eval` with `rubrics/test-cases.md` |
| gen-test-scripts | `testing/test-cases.md` (hard) | `tests/e2e/features/<slug>/` | profile → framework; Step Actionability gate |
| run-e2e-tests | justfile + staging area | `tests/e2e/features/<slug>/results/latest.md` | >30% failure → stop |
| graduate-tests | staging area + PASS results + no marker | `tests/e2e/<module>/`, `.graduated/<slug>` | profile → import rewriting |
| consolidate-specs | PRD + design (both hard) | `specs/`, updated `docs/business-rules/`, `docs/conventions/` | CROSS vs LOCAL |
| /quick | User idea | Orchestrates quick pipeline | >10 tasks → STOP |
| quick-tasks | proposal.md (hard) | `tasks/*.md`, `tasks/index.json`, `manifest.md` | max 10 tasks |
| /run-tasks | `tasks/index.json` | Execution loop | 3 consecutive failures → STOP |
| submit-task | Task executed + record.json | `records/*.md`, updated index.json | --force bypasses gate (not auto-downgrade) |

**Scoring Criteria:**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 1a. Full mode chain complete | 0-80 | Every skill has prerequisite definitions and outputs. Prerequisite files are produced by predecessor skills. Chain has no breakpoints. Breakpoint = -20 each. |
| 1b. Quick mode chain complete | 0-40 | Quick chain is complete. Missing step = -20 each. |
| 1c. Conditional branching correct | 0-50 | Every conditional branch has a true-path and false-path. Missing branch = -10 each. |
| 1d. Manifest status transitions valid | 0-30 | Status transitions are legal. Illegal transition = -15 each. |
| 1e. Test lifecycle chain intact | 0-50 | Test chain is complete. Broken link = -15 each. |

### 2. Bypass Resistance (250 pts)

Assume you are a "lazy agent trying to cut corners" — test each node for bypass paths.

**5 Bypass Types:**

| Type | Points | Description |
|------|--------|-------------|
| Type 2: Skip quality gates | 0-70 | `--force` bypasses compile/test/AC. No justfile means gate silently passes. noTest skips. |
| Type 3: Fake eval results | 0-70 | Main session can fake scorer output. Score parsing lacks integrity checks. Scorer subagent can be skipped. |
| Type 1: Skip mandatory interaction | 0-45 | User confirmation points are all advisory text (HARD-RULE). Can skip brainstorm/write-prd/ui-design/tech-design/consolidate-specs confirmations. |
| Type 4: Skip required steps | 0-35 | Conditional requirements depend on agent self-reporting (db-schema, placement). gen-test-scripts Step Actionability gate only triggers when eval report exists. |
| Type 5: Lazy shortcuts | 0-30 | Prohibition patterns (no mock, no sleep) are purely advisory. record.json metrics are self-reported. Direct index.json editing bypasses all validation. |

**Scoring Criteria:**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 2a. Quality gate enforcement | 0-70 | Is each gate point enforced by CLI or merely advisory text. Zero enforcement with no documented rationale = -15 each. |
| 2b. Eval integrity | 0-70 | Does the generic eval skill (`skills/eval/SKILL.md`) require an independent subagent for scoring. Does the decision gate parse structured output. Can the main session fake the score. Weakness = -25 each. |
| 2c. User interaction enforcement | 0-45 | Does each confirmation point have an enforcement mechanism. Purely advisory = -5 each. |
| 2d. Required step enforcement | 0-35 | Do conditional requirements have downstream verification. No verification = -10 each. |
| 2e. Prohibition enforcement | 0-30 | Does each HARD-RULE prohibition have a mechanical check. Purely advisory = -5 each. |

**Known Bypass Vectors (from manual audit — scorer should verify current state):**

| Vector | Type | Severity | Class | Description |
|--------|------|----------|-------|-------------|
| BV-2.1 | T2 | HIGH | ARCHITECTURAL | `forge task submit --force` bypasses compile/test/AC all validation |
| BV-2.4 | T2 | HIGH | ARCHITECTURAL | Agent can fake testsFailed: 0 in record.json (CLI does not verify number source) |
| BV-3.1 | T3 | HIGH | ARCHITECTURAL | Main session can skip scorer subagent and directly declare SCORE: 950 |
| BV-3.2 | T3 | HIGH | ARCHITECTURAL | Score parsing has no integrity check, main session can tamper with scorer return values |
| BV-1.1 | T1 | MED | TEXT-FIXABLE | brainstorm user approval can be skipped (pure HARD-RULE) |
| BV-1.2 | T1 | MED | TEXT-FIXABLE | write-prd user approval can be skipped |
| BV-1.3 | T1 | MED | TEXT-FIXABLE | /quick user confirmation can be skipped |
| BV-1.4 | T1 | MED | TEXT-FIXABLE | ui-design prototype review can be skipped |
| BV-1.5 | T1 | MED | TEXT-FIXABLE | tech-design DB schema review can be skipped |
| BV-1.7 | T1 | MED | TEXT-FIXABLE | consolidate-specs spec integration confirmation can be skipped |
| BV-4.2 | T4 | MED | TEXT-FIXABLE | gen-test-scripts Step Actionability gate only triggers when eval report exists |
| BV-4.5 | T4 | MED | TEXT-FIXABLE | placement validation depends on sitemap existence, skipped if absent |
| BV-5.1 | T5 | LOW | ARCHITECTURAL | Prohibition patterns (no mock/no sleep/no hardcoded URL) are purely advisory |
| BV-5.2 | T5 | LOW | ARCHITECTURAL | record.json metrics (coverage/testsPassed) are self-reported, no cross-verification |

**Bypass classification for fix strategy:**

- **ARCHITECTURAL**: Cannot be fixed by adding text to SKILL.md/command files. These require code-level changes (CLI enforcement, cryptographic verification, etc.). Score them and report them, but do NOT generate reviser fix tasks. Adding HARD-RULE text for architectural bypasses is counterproductive — it inflates context without changing agent behavior.
- **TEXT-FIXABLE**: Can be mitigated by adding conditional branches, fallback paths, or actionable instructions. These are valid targets for the reviser.

### 3. Instruction Precision (200 pts)

**Priority order: Instruction conflicts > Step ambiguity > Incomplete conditionals > Undefined variables**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 3a. Instruction conflicts (cross-file) | 0-80 | guide.md vs SKILL.md vs command files describe the same concept differently. E.g., guide says "lint blocks" but skill says "lint non-blocking". Conflict = -25 each. **Check this first.** |
| 3b. Step ambiguity | 0-50 | SKILL.md steps must have a single unambiguous interpretation. Vague verbs ("check tests", "verify quality") without specific commands = -10 each. |
| 3c. Incomplete conditionals | 0-40 | Every if-then must have an else path or an explicit "skip" instruction. Missing else = -10 each. |
| 3d. Variable resolution clarity | 0-30 | Agent-filled variables must have source annotations. CLI-filled variables must match `prompt.go` typeToTemplate. Undefined agent variable = -10 each. |

**CLI-Filled Variables (from `prompt.go` Synthesize — do NOT mark as "undefined"):**

| Variable | Source | Used by |
|----------|--------|---------|
| `{{TASK_ID}}` | task.ID | All 14 templates |
| `{{TASK_FILE}}` | feature.GetTaskFile() | All 14 templates |
| `{{SCOPE}}` | task.Scope | Most templates |
| `{{PHASE_SUMMARY}}` | PhaseDetect() | 10 templates |
| `{{FEATURE_SLUG}}` | SynthesizeOpts.FeatureSlug | 4 templates |
| `{{PROFILE}}` | task.Profile | 4 test pipeline templates |

### 4. Cross-file Dedup (150 pts)

**Three categories with different standards:**

<PLUGIN-PORTABILITY>
Forge is a Claude Code plugin deployed to users' projects. Plugin SKILL.md and command files run in arbitrary working directories. Cross-file relative paths (`../../`) won't resolve at runtime. Therefore:
- Duplication that serves plugin portability is **acceptable and should NOT be deducted**.
- Only deduct when dedup is feasible without breaking portability (e.g., content within the same skill directory, or in non-plugin project-level files like `.claude/skills/`).
- guide.md cannot be referenced via relative path from plugin SKILL.md files. Inlining guide.md content in plugin files is a necessary trade-off.
</PLUGIN-PORTABILITY>

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 4a. Content copy | 0-60 | Identical/near-identical text blocks appearing in 3+ files. **Plugin portability exception:** Known instances (profile resolution in 9 files, eval protocol in 6 files, report sections in 5 files) serve plugin portability — flag as INFO but do NOT deduct. Only deduct when dedup is feasible within the same skill directory or in non-plugin files. Actionable instance = -10 each. |
| 4b. guide.md vs SKILL.md overlap | 0-50 | guide.md is the single source of truth. **But plugin files cannot reference guide.md via relative paths** — they must inline the content. Only deduct for non-plugin files or when referencing is feasible. Actionable duplication = -10 each. |
| 4c. Unreasonable inline | 0-40 | Content that has its own file but is also inlined in SKILL.md. Judgment: does the agent need to see the full content in a single context (reasonable inline) vs can it be obtained via Read tool on demand (should reference). **Plugin portability:** inlining content from outside the skill directory is reasonable when the alternative (relative path reference) would break at runtime. Unreasonable inline = -10 each. |

**Known Redundancy Instances:**

| Instance | Category | Portability-Required? | Deduct? |
|----------|----------|----------------------|---------|
| "Step 0: Resolve Profile" | A | YES (9 plugin SKILL.md files) | NO |
| Eval Iron Laws + Steps 2-4 | A | NO (consolidated to 1 `skills/eval/SKILL.md`) | NO |
| Eval report shared sections | A | YES (5 plugin report.md files) | NO |
| Quality gate sequence | B | YES (plugin skills need it inline) | NO |
| Scope resolution paraphrase | B | Partial | Only if in non-plugin file |

### 5. Reference Integrity (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 5a. Agent references valid | 0-30 | `forge:<agent>` or `subagent_type` points to an existing file in `plugins/forge/agents/`. Dangling = -15 each. |
| 5b. Template references valid | 0-25 | Template paths in SKILL.md point to existing files. Dangling = -15 each. |
| 5c. Cross-skill references valid | 0-25 | `invoke /<name>` points to an existing skill/command. Dangling = -15 each. |
| 5d. Hook references valid | 0-20 | Paths and CLI commands in hooks.json exist. Dangling = -15 each. |

### 6. Structural Convention (50 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 6a. Frontmatter completeness | 0-25 | SKILL.md has name + description. Command has name + description. Agent has name + description + model. Missing = -5 each. |
| 6b. Eval template convention | 0-15 | `skills/eval/SKILL.md` exists, `skills/eval/rubrics/` contains rubric files for each eval type, and each `commands/eval-*.md` delegates to `Skill("eval", ...)`. Missing rubric = -10 each. |
| 6c. Name-directory alignment | 0-10 | Skill name = directory name, command name = filename. Mismatch = -5 each. |
