# Forge Plugin Runtime Reliability Rubric

**Total: 1000 points**
**Report template:** `.claude/skills/eval-forge/templates/report.md`

## What This Rubric Measures

Runtime reliability of the forge plugin — not just structural consistency, but whether components work correctly together at runtime. Measures workflow completeness, bypass resistance, instruction precision, redundancy, reference integrity, and structural conventions.

## Scoring Methodology: 4-Phase Process

This rubric is designed for a 4-phase evaluation methodology:

**Phase 1: Construct Workflow Graph (D1)**
1. Read guide.md pipeline diagrams as the authoritative workflow description
2. Scan all SKILL.md/command files, extract actual prerequisites, outputs, gate points
3. Build directed graph, check connectivity, find breakpoints and unauditable steps
4. Build intra-skill step graphs: extract step sequences, conditional skip targets, and detection points; verify temporal ordering consistency (see criterion 1f)

**Phase 2: Per-Node Adversarial Testing (D2)**
1. List every gate/confirm point per node
2. For each gate, assume you are a lazy agent — ask "how can I bypass this?"
3. Check whether HARD-RULE has enforceable consequences (not just words)
4. Eval loops: check whether scoring must be done by independent subagent
5. Quality gates: check whether CLI-level enforcement exists

**Phase 3: Per-File Precision Review (D3 + D4)**
1. Instruction conflicts (highest priority): behavioral cross-reference with CLI source code, instruction chain tracing, then cross-file same-concept comparison
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
| 1. Workflow Completeness | 280 |
| 2. Bypass Resistance | 280 |
| 3. Instruction Precision | 280 |
| 4. Cross-file Dedup | 30 |
| 5. Reference Integrity | 80 |
| 6. Structural Convention | 50 |
| **Total** | **1000** |

Score allocation logic: Runtime reliability (D1+D2+D3) = 840 pts (84%), Baseline integrity (D4+D5+D6) = 160 pts (16%).

## Deduction Tiers

| Severity | Penalty | Examples |
|----------|---------|---------|
| Critical | -20 | Workflow breakpoint, eval integrity bypass, instruction conflict, temporal ordering contradiction |
| High | -15 | Dangling reference, invalid manifest transition, missing eval template |
| Medium | -10 | Step ambiguity, conditional without else, content duplication, missing prerequisite check |
| Low | -5 | Missing frontmatter field, name-directory mismatch, advisory-only prohibition, context inflation |

**Floor rule:** Each sub-criterion has a minimum score of 0. Total deductions for a sub-criterion cannot result in a negative score. When deductions exceed the sub-criterion max, report the uncapped total in the Deductions table as a severity indicator (e.g., "5 dangling references × -10 = -50, capped at 0/20").

**Cross-dimension deduction:** The same issue may be deducted in multiple dimensions if it violates independent criteria (e.g., an advisory-only HARD-RULE is both a D2 bypass vector and a D3 incomplete conditional). Bypass resistance (D2) and instruction precision (D3) are independent dimensions.

## Dimensions

### 1. Workflow Completeness (280 pts)

Build the workflow graph dynamically from actual plugin files and verify architectural soundness.

**Discovery procedure:**
1. Read guide.md — extract pipeline mermaid diagrams and quality gate protocol as the authoritative workflow description
2. Read all SKILL.md files — extract declared prerequisites (hard vs optional), outputs, gate points, conditional branches, successor skills
3. Read all command files — extract skill invocations, tool references
4. Read `forge-cli/pkg/prompt/prompt.go` + `forge-cli/pkg/prompt/data/*.md` — identify CLI-embedded task templates (steps with no SKILL.md)
5. Read `forge-cli/internal/cmd/testgen.go` — discover actual test task types generated for full mode and quick mode
6. Build directed graph: skill → prerequisite files → producer skill → successor skill

**Scoring Criteria:**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 1a. Graph connectivity | 0-80 | Every skill's hard prerequisites are produced by some predecessor skill/command. No disconnected subgraphs. Build the graph from SKILL.md declared prerequisites/outputs and verify every edge resolves. Breakpoint (unmet prerequisite) = -20 each. **CLI-embedded templates** (`forge-cli/pkg/prompt/data/*.md`): read each template and classify: (1) *thin dispatcher* — delegates to a SKILL.md, no real logic = no deduction; (2) *real logic* — multi-step workflow, rubric, validation, or behavioral instructions = -10 each (should be independently auditable via SKILL.md). For all templates: check instruction quality, error handling, consistency with corresponding SKILL.md, and stale references. Apply standard deductions for bugs/conflicts found. |
| 1b. Quick mode completeness | 0-40 | Quick pipeline from guide.md mermaid diagram is complete. Test tasks discovered from testgen.go quick-mode generation have corresponding SKILL.md or are flagged as CLI-embedded. Missing step = -20 each. |
| 1c. Conditional branching | 0-30 | Every conditional branch in SKILL.md has a true-path AND false-path. Missing branch = -10 each. |
| 1d. Status consistency | 0-30 | Status transitions are consistent across all files (guide.md, SKILL.md, commands). Cross-reference with `status.go` `isTransitionAllowed()` for actual guards. Inconsistency = -15 each. |
| 1e. Test chain connectivity | 0-70 | Test chain discovered from testgen.go task types + guide.md test lifecycle diagram is complete: every task type's prerequisites are satisfied by a predecessor. For test-related CLI-embedded templates in `prompt/data/test-pipeline-*.md`: read each, verify instruction quality and consistency with corresponding SKILL.md. Broken link = -15 each. |
| 1f. Intra-skill temporal ordering | 0-30 | Every conditional skip's detection point precedes the step it modifies. Detection-after-skip-target = -15 each. |

**1f Detailed Check Procedure:**

For each SKILL.md that defines numbered steps with conditional skip/fast-path logic:

1. **Extract step sequence**: Identify all numbered steps (Step 0, Step 1, Step 2, ...) and their execution order.
2. **Identify conditional skip targets**: Find instructions like "skip Step N if condition C" or "Steps N-M are unnecessary when condition C".
3. **Locate detection points**: Find where condition C is evaluated — look for "Detection:" paragraphs, "After Step N" markers, "During Step N" markers, or mermaid decision nodes.
4. **Verify temporal consistency**: For each (detection_point, skip_target) pair, check that `detection_step <= skip_target`. A detection point at "After Step 4a" cannot skip "Step 0" because Step 0 executes before the detection.
5. **Check mermaid diagrams**: If a mermaid flow diagram exists, trace the decision node's incoming edge — does it connect from a step that precedes the skipped step?

**Violation pattern:** Detection at "After Step N" with skip target "Step M where M < N" = contradiction (Step M already executed before detection).

### 2. Bypass Resistance (280 pts)

Phase 2 uses a two-stage approach: independent discovery first (systematically find all enforcement gaps), then regression verification (confirm previously-fixed TEXT-FIXABLE vectors stay fixed).

**5 Bypass Types:**

| Type | Points | Examples |
|------|--------|----------|
| Type 2: Skip quality gates | 0-80 | `--force` overrides, noTest bypasses, missing justfile silently passes, self-reported metrics |
| Type 3: Fake eval results | 0-80 | Skip scorer subagent, tamper with scores, general-purpose fallback loses adversarial constraints |
| Type 1: Skip mandatory interaction | 0-45 | User confirmations are advisory-only HARD-RULE text, no mechanical enforcement |
| Type 4: Skip required steps | 0-35 | Conditional requirements depend on self-reporting, gates only trigger when prerequisites exist |
| Type 5: Lazy shortcuts | 0-40 | Prohibitions are purely advisory, metrics are self-reported, direct file editing bypasses CLI |

**Scoring Criteria:**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 2a. Quality gate enforcement | 0-80 | Is each gate point enforced by CLI or merely advisory text. Zero enforcement with no documented rationale = -15 each. |
| 2b. Eval integrity | 0-80 | Does each eval skill require independent subagent scoring. Does the decision gate parse structured output. Can the main session fake the score. Weakness = -25 each. |
| 2c. User interaction enforcement | 0-45 | Does each confirmation point have a mechanical enforcement mechanism. Purely advisory = -10 each. |
| 2d. Required step enforcement | 0-35 | Do conditional requirements have downstream verification. No verification = -10 each. |
| 2e. Prohibition enforcement | 0-40 | Does each HARD-RULE prohibition have a mechanical check. Purely advisory = -5 each. |

**Bypass classification for fix strategy:**

- **ARCHITECTURAL**: Cannot be fixed by adding text to SKILL.md/command files. These require code-level changes (CLI enforcement, cryptographic verification, etc.). Score them and report them, but do NOT generate reviser fix tasks. Adding HARD-RULE text for architectural bypasses is counterproductive — it inflates context without changing agent behavior.
- **TEXT-FIXABLE**: Can be mitigated by adding conditional branches, fallback paths, or actionable instructions. These are valid targets for the reviser.

**Regression Guard (Phase 2b):**

The following TEXT-FIXABLE bypass vectors were previously identified. Verify each one's current state. Do not deduct for vectors that have been fixed.

| Vector | Type | Severity | Description |
|--------|------|----------|-------------|
| BV-1.1 | T1 | MED | brainstorm user approval can be skipped (pure HARD-RULE) |
| BV-1.2 | T1 | MED | write-prd user approval can be skipped |
| BV-1.3 | T1 | MED | /quick user confirmation can be skipped |
| BV-1.4 | T1 | MED | ui-design prototype review can be skipped |
| BV-1.5 | T1 | MED | tech-design DB schema review can be skipped |
| BV-1.7 | T1 | MED | consolidate-specs spec integration confirmation can be skipped |
| BV-4.2 | T4 | MED | gen-test-scripts Step Actionability gate only triggers when eval report exists |
| BV-4.5 | T4 | MED | placement validation depends on sitemap existence, skipped if absent |

### 3. Instruction Precision (280 pts)

**Priority order: Instruction conflicts > Step ambiguity > Incomplete conditionals > Undefined variables**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 3a. Instruction conflicts | 0-100 | Bidirectional verification: **(1) Plugin→CLI**: For every behavioral claim in SKILL.md/commands/guide.md (command refs, data model, status transitions, quality gates, enforcement claims, template variables, hook behavior), trace to CLI source and verify match. **(2) CLI→Plugin**: For every CLI behavior affecting agent workflow (auto-downgrade, transition guards, silent pass conditions, --force scope, docs-only detection, all-completed logic), verify it's documented in plugin files. **(3) Chain**: Step N's output mismatches Step N+1's expected input. **(4) Cross-file**: same concept described differently across files. Contract mismatch = -25 each. |
| 3b. Step ambiguity | 0-60 | SKILL.md steps must have a single unambiguous interpretation. Vague verbs ("check tests", "verify quality") without specific commands = -10 each. Also flag steps where chain tracing reveals the agent must make unstated choices. |
| 3c. Incomplete conditionals | 0-50 | Every if-then must have an else path or an explicit "skip" instruction. Missing else = -10 each. **Implicit-else exception:** if the false-path is the natural default (normal execution continues, zero-value default, or no-op), no explicit else is required. Only flag if-then patterns where the false-path requires distinct handling but none is documented. |
| 3d. Variable resolution clarity | 0-40 | Agent-filled variables must have source annotations. CLI-filled variables are discovered by reading `forge-cli/pkg/prompt/prompt.go` Synthesize function — do NOT mark prompt.go variables as undefined. Undefined agent variable = -10 each. |
| 3e. Narrative inflation | 0-30 | Text that inflates context without changing agent behavior: consequence/rationale paragraphs (why vs what), stale code/function references (pointing to moved or non-existent files), redundant re-explanation of what a table/step/code-block already states. Instance = -5 each. **Exempt**: content inside `<HARD-RULE>`, `<HARD-GATE>`, `<EXTREMELY-IMPORTANT>` tags. |

### 4. Cross-file Dedup (30 pts)

**Three categories with different standards:**

<PLUGIN-PORTABILITY>
Forge is a Claude Code plugin deployed to users' projects. Plugin SKILL.md and command files run in arbitrary working directories. Cross-file relative paths (`../../`) won't resolve at runtime. Therefore:
- Duplication that serves plugin portability is **acceptable and should NOT be deducted**.
- Only deduct when dedup is feasible without breaking portability (e.g., content within the same skill directory, or in non-plugin project-level files like `.claude/skills/`).
- guide.md cannot be referenced via relative path from plugin SKILL.md files. Inlining guide.md content in plugin files is a necessary trade-off.
</PLUGIN-PORTABILITY>

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 4a. Content copy | 0-10 | Identical/near-identical text blocks appearing in 3+ files. **Plugin portability exception:** Known instances serve plugin portability — flag as INFO but do NOT deduct. Only deduct when dedup is feasible within the same skill directory or in non-plugin files. Actionable instance = -10 each. |
| 4b. guide.md vs SKILL.md overlap | 0-10 | guide.md is the single source of truth. **But plugin files cannot reference guide.md via relative paths** — they must inline the content. Only deduct for non-plugin files or when referencing is feasible. Actionable duplication = -10 each. |
| 4c. Unreasonable inline | 0-10 | Content that has its own file but is also inlined in SKILL.md. Judgment: does the agent need to see the full content in a single context (reasonable inline) vs can it be obtained via Read tool on demand (should reference). **Plugin portability:** inlining content from outside the skill directory is reasonable when the alternative would break at runtime. Unreasonable inline = -10 each. |

**Known Redundancy Instances:**

| Instance | Category | Portability-Required? | Deduct? |
|----------|----------|----------------------|---------|
| "Step 0: Resolve Profile" | A | YES (9 plugin SKILL.md files) | NO |
| Eval Iron Laws + Steps 2-4 | A | NO (consolidated to 1 `skills/eval/SKILL.md`) | NO |
| Eval report shared sections | A | YES (5 plugin report.md files) | NO |
| Quality gate sequence | B | YES (plugin skills need it inline) | NO |
| Scope resolution paraphrase | B | Partial | Only if in non-plugin file |

### 5. Reference Integrity (80 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 5a. Agent references valid | 0-20 | `forge:<agent>` or `subagent_type` points to an existing file in `plugins/forge/agents/`. Dangling = -10 each. |
| 5b. Template references valid | 0-20 | Template paths in SKILL.md point to existing files. **Also read template content** — verify format/fields match SKILL.md output descriptions. Format mismatch = -10 each. Dangling path = -10 each. |
| 5c. Cross-skill references valid | 0-15 | `invoke /<name>` points to an existing skill/command. Dangling = -10 each. |
| 5d. Hook references valid | 0-15 | Paths and CLI commands in hooks.json exist. Dangling = -10 each. |
| 5e. Shared reference paths valid | 0-10 | `plugins/forge/references/*` paths in SKILL.md or commands must exist. Dangling = -5 each. |

### 6. Structural Convention (50 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 6a. Frontmatter completeness | 0-25 | SKILL.md has name + description. Command has name + description. Agent has name + description + model. Missing = -5 each. |
| 6b. Eval template convention | 0-15 | `skills/eval/SKILL.md` exists, `skills/eval/rubrics/` contains rubric files for each eval type, and each `commands/eval-*.md` delegates to `Skill("eval", ...)`. **Also read rubric content** — verify scoring dimensions match eval SKILL.md instructions. Missing or contradictory dimension = -10 each. |
| 6c. Name-directory alignment | 0-10 | Skill name = directory name, command name = filename. Mismatch = -5 each. |
