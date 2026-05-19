---
name: guide-simplify
status: Draft
created: 2026-05-19
---

# Simplify and Restructure guide.md

## Problem

`plugins/forge/hooks/guide.md` (241 lines) is loaded on every session and subagent start. It contains 6 sections of content that is duplicated in skill-specific files (Skill Workflow diagrams, Testing Lifecycle, Knowledge Accumulation, Eval Exceptions, Auxiliary Skills table, Quick Mode details). This wastes ~60% of the token budget on information agents can find on-demand.

## Solution

Remove duplicated content and restructure the remaining ~100-120 lines into 3 thematic sections: Directory Conventions → Execution Rules → Automation Config. Agents continue to find removed details in their respective skill files.

## Alternatives

| Approach | Trade-off |
|----------|-----------|
| **Do nothing** | 241 lines loaded per session; wasted tokens; no risk |
| **Section-by-section trim** | Safe but leaves structural debt; sections don't flow logically |
| **Thematic restructure (chosen)** | Cleaner flow; slightly larger diff but one-time cost |

## Scope

**In Scope:**
- Remove 6 sections: Skill Workflow mermaid diagrams, Quick Mode details, Testing Lifecycle, Evaluation Parameter Exceptions, Knowledge Accumulation details, Auxiliary Skills table
- 4 items are guide-only (not duplicated in skills): pipeline diagrams, Quick/Full comparison, 3-layer testing model, auto-extract trigger table. User confirmed: delete directly — these are reference docs, not hard rules agents need globally.
- Restructure remaining content into 3 themes: Directory Conventions → Execution Rules → Automation Config
- Preserve all factual rules accurately
- Target: ~100-120 lines (from 241)

**Out of Scope:**
- Changes to skill-specific SKILL.md files
- Migration of removed content to docs/reference/
- Changes to hooks.json or hook behavior
- Any functional changes to forge behavior

## Risks

| Risk | Mitigation |
|------|------------|
| Agents lose quick reference to pipeline flow | Each skill's SKILL.md already documents its own prerequisites; /quick command has full pipeline |
| Auto-extract triggers not known globally | /learn and /consolidate-specs skill files contain all details |
| Eval parameter exceptions not visible | /eval-ui and /eval-test-cases skill files contain their own parameters |

## Success Criteria

- [ ] guide.md is ~100-120 lines (from 241)
- [ ] All removed content exists in corresponding skill files (verify by grep)
- [ ] No functional behavior change — quality gate, scope resolution, auto-config remain accurate
- [ ] Three clear thematic sections with logical flow
