# Eval-Proposal Report — Iteration 1

**Score: 789/1000** (target: 900)

## DIMENSIONS

### 1. Problem Definition (88/110)
- Problem stated clearly: 34/40 — Core problem is unambiguous: "Forge consumers...cannot discover what files exist or determine which are relevant." Concrete file/line references ground the problem. Deduction: no quantitative data on failure frequency or user-visible impact.
- Evidence provided: 32/40 — Strong code-level evidence with specific file references (`fix-bug.md:50-51`, `error-fixer.md:64-65`, 5 named prompt templates). Deduction: no user feedback, no failure reports, no data on how often agents load wrong/no knowledge in practice.
- Urgency justified: 22/30 — "Each new consumer that needs project knowledge invents its own approach...the pattern diverges further with each new skill." Deduction: urgency is argued qualitatively but not quantified. How many new consumers are planned? What concrete failures has the divergence already caused?

### 2. Solution Clarity (105/120)
- Approach is concrete: 35/40 — Frontmatter schema with YAML example, discovery instruction text, consolidate-specs auto-generation. A reader can explain back exactly what will be built. Deduction: unclear whether the discovery instruction is embedded verbatim or adapted per consumer.
- User-facing behavior described: 38/45 — Key Scenarios section describes agent behavior: reads frontmatter, loads relevant files, skips irrelevant ones. Edge cases covered (empty dirs, missing domains). Deduction: the end-user perspective is thin — what does the human running `/fix-bug` actually notice differently? Better responses? Faster execution? This is implied but never stated.
- Technical direction clear: 32/35 — YAML frontmatter `domains` field, plain text instruction, one skill update. Clear enough to implement.

### 3. Industry Benchmarking (60/120)
- Industry solutions referenced: 12/40 — No external tools, open-source projects, or published patterns are cited. The entire comparison is against internally invented approaches. Missing references: VS Code extension discovery, RAG retrieval patterns, package metadata indexing (npm keywords, Cargo categories), static site generator tag systems, IDE workspace indexing.
- At least 3 meaningful alternatives: 18/30 — Five alternatives listed including "do nothing," but none reference industry-validated solutions. "Load all files" is close to a straw man — trivially rejected without serious consideration. "CLI-level injection" is another straw man framed as "too rigid." Only the consumer-declared pattern is genuinely battle-tested (within the codebase).
- Honest trade-off comparison: 20/25 — The selected approach's weakness is honestly stated: "Depends on agent correctly executing discovery; no code-level guarantee." Internal comparisons are fair.
- Chosen approach justified against benchmarks: 10/25 — Justification is against internal alternatives only. No external benchmark exists to justify against. The rationale "best fit for plugin context" is asserted but not demonstrated against industry standards.

### 4. Requirements Completeness (92/110)
- Scenario coverage: 34/40 — Six scenarios including happy path, no-conventions, file-without-domains. Deduction: missing scenarios for (1) ambiguous domain matches where multiple files score equally, (2) large directories with 20+ convention files, (3) consolidate-specs run on files with manually edited domains, (4) conflicting domain tags across files.
- Non-functional requirements: 32/40 — Token efficiency, graceful degradation, zero maintenance. Deduction: no performance consideration for scanning directories with many files, no discussion of frontmatter trust/reliability, no compatibility analysis with existing consolidate-specs behavior.
- Constraints & dependencies: 26/30 — Convention files need domains frontmatter, prompt templates are plain text, mechanism is guidance not protocol. Good. Deduction: "agent comprehension" as a dependency is vague — what level of model capability is assumed?

### 5. Solution Creativity (56/100)
- Novelty over industry baseline: 22/40 — Frontmatter metadata tagging is a well-established pattern (Jekyll tags, npm keywords, Cargo categories, Atom feed categories). Applying it to agent prompt templates is a reasonable adaptation but not strongly differentiated from the industry baseline. The proposal does not articulate what makes this application novel beyond context.
- Cross-domain inspiration: 12/35 — No cross-domain references or inspiration are cited. The idea of metadata tagging is generic. Missing opportunity to reference how search engines, IDEs, or knowledge management systems solve analogous discovery problems.
- Simplicity of insight: 22/25 — The insight "make files self-describing" is genuinely elegant. The solution has a "why didn't I think of that" quality. Clean minimal design.

### 6. Feasibility (94/100)
- Technical feasibility: 38/40 — All components exist in the current system. Adding a YAML field to existing frontmatter, updating text templates, updating one skill. No showstoppers.
- Resource & timeline feasibility: 28/30 — 6 files + 5 templates + 2 commands + 1 skill = 14 small, well-bounded changes. Clearly achievable.
- Dependency readiness: 28/30 — No external dependencies. All files are local and version controlled.

### 7. Scope Definition (76/80)
- In-scope items are concrete: 28/30 — Six numbered items, each naming specific files and the exact change needed. Excellent specificity.
- Out-of-scope explicitly listed: 24/25 — Five items explicitly deferred with reasons. Very good.
- Scope is bounded: 24/25 — 14 files total, all changes described as "small and well-bounded." Tight, executable scope.

### 8. Risk Assessment (72/90)
- Risks identified: 24/30 — Four risks identified. Deduction: missing risks for (1) consolidate-specs corrupting existing frontmatter, (2) domain keyword collision where two files claim the same domain, (3) user confusion about the domains field and manual edits, (4) frontmatter parsing failures in non-standard file encodings.
- Likelihood + impact rated: 26/30 — Ratings are varied and honest (M/L mix). No "all high" inflation pattern.
- Mitigations are actionable: 22/30 — Some mitigations are concrete ("agent uses filename/title if domains missing", "Human reviews via consolidate-specs confirmation flow"). Others are design descriptions rather than mitigations: "Instruction is brief and action-oriented" describes the instruction, not a response to the risk. "domains frontmatter provides concrete matching criteria" is a feature description, not a mitigation strategy.

### 9. Success Criteria (62/80)
- Criteria are measurable and testable: 42/55 — Most criteria are objectively verifiable: "all 6 files have valid domains," "no keyword→filename mapping tables remain." Deduction: "valid domains frontmatter with 3-7 specific keywords" — "specific" is subjective. What counts as specific enough? The behavioral criterion about agent loading correct files is observable but not automatable.
- Coverage is complete: 20/25 — Seven criteria cover scope items 1-5. Deduction: no criterion for scope item 6 (guide.md update). No criterion for graceful degradation when files lack domains (the scenario is described in requirements but not in success criteria).

### 10. Logical Consistency (84/90)
- Solution addresses the stated problem: 34/35 — Problem: consumers can't discover files. Solution: files self-describe via domains, consumers use discovery instruction. Direct, tight alignment.
- Scope ↔ Solution ↔ Success Criteria aligned: 26/30 — Six in-scope items map to solution components. Most success criteria correspond to scope items. Deduction: scope item 6 (guide.md update) has no corresponding success criterion — an alignment gap.
- Requirements ↔ Solution coherent: 24/25 — Requirements map cleanly to solution components. No orphan requirements or solution features without requirements.

## ATTACKS

1. [Industry Benchmarking]: Zero external references — The entire Alternatives section compares internal approaches only. No product, open-source project, or published pattern is cited. Quote: the table lists "Do nothing", "CLI-level injection", "Load all files", "Consumer-declared dependencies", and "Frontmatter domain tags + agent discovery" — all self-invented. Must add at least 2-3 industry references (e.g., RAG retrieval, package metadata indexing, IDE context discovery) and evaluate the chosen approach against them.

2. [Industry Benchmarking]: Straw-man alternatives — "Load all files unconditionally" is trivially dismissed with "Token waste at scale; loads irrelevant knowledge; no matching." This is not a serious alternative in any system with <20 files. Similarly, "CLI-level injection" is rejected as "too rigid for a plugin" without exploring partial adoption. Must replace straw men with genuine industry alternatives.

3. [Solution Creativity]: No cross-domain inspiration cited — The proposal does not reference how any other domain solves analogous discovery problems. Quote: the entire section "Proposed Solution" and "Alternatives & Industry Benchmarking" contain zero external references. Must explicitly identify at least one external domain/product that inspired the approach (e.g., "Similar to how Jekyll uses frontmatter tags for content discovery...").

4. [Success Criteria]: Missing criterion for scope item 6 — Scope item 6 is "Update `plugins/forge/hooks/guide.md` — update the project knowledge note to reference `domains` frontmatter." No success criterion verifies this change. Must add a criterion: "guide.md contains a reference to domains frontmatter in the project knowledge section."

5. [Success Criteria]: "Specific keywords" is untestable — Quote: "valid `domains` frontmatter with 3-7 specific keywords each." The word "specific" is subjective. Must define what makes a keyword specific (e.g., "each keyword appears in the file's content at least once" or "each keyword is a domain term, not a generic word").

6. [Risk Assessment]: Mitigations describe features, not responses — Quote for "Agent skips discovery step": "Instruction is brief and action-oriented. Graceful degradation: missing knowledge doesn't block execution." The first part describes the instruction design, not how to handle the risk. Must replace with actionable mitigation: "Add a post-task check that logs whether convention files were consulted, and surface a warning if not."

7. [Risk Assessment]: Missing collision risk — No risk identified for two convention files claiming overlapping domain keywords. If `error-handling.md` has domains: [error, status] and `error-reporting.md` has domains: [error, status, log], both match an error-related task, potentially loading redundant content. Must add a risk for domain keyword overlap/redundancy with mitigation.

8. [Requirements Completeness]: Missing scenario for ambiguous matches — Six scenarios are described but none address what happens when domain matching is ambiguous (multiple files with similar domains). Quote: the scenarios describe clear matches (error-handling matches, profile-system doesn't) but never the borderline case. Must add a scenario where multiple files score equally and the agent must decide.

9. [Solution Clarity]: User-facing impact not described — The proposal describes what the agent does but not what the human user notices. Quote: scenarios are all from the agent's perspective ("Agent receives T-impl-3... Agent discovers..."). Must add explicit user-facing outcome: "User running `/fix-bug` receives responses informed by relevant project conventions without needing to manually reference them."

10. [Problem Definition]: No failure evidence — The Current State section identifies the divergence pattern but provides no evidence of actual failures caused by it. Quote: "5 prompt templates...say 'Read relevant project knowledge files' with no guidance." Has this led to wrong agent behavior? Missing conventions? User complaints? Must include at least one concrete failure case or user report.
