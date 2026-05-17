---
created: 2026-05-16
updated: 2026-05-17
author: faner
status: Draft
---

# Proposal: Project Knowledge Discovery via Frontmatter Domain Tags

## Problem

Forge consumers (prompt templates, commands, agents) benefit from user-provided project knowledge in `docs/conventions/` and `docs/business-rules/`, but cannot discover what files exist or determine which are relevant to the current task. The filenames are unknown to forge — they are created by users (via `/consolidate-specs`) and vary across projects.

### Current State

- `fix-bug.md:50-51` and `error-fixer.md:64-65` hardcode identical keyword→filename mapping tables, assuming specific files exist
- 5 prompt templates (`fix.md`, `cleanup.md`, `enhancement.md`, `feature.md`, `refactor.md`) say "Read relevant project knowledge files" with no guidance on how to discover or evaluate relevance
- `gen-test-scripts`, `gen-test-cases`, and `eval` skills solved this differently — via consumer-declared frontmatter (`conventions: [filename]`) where the skill explicitly lists its dependencies. This works because those skills know their domain (testing). Prompt templates handle arbitrary tasks and cannot predict which conventions they need.

### Urgency

Each new consumer that needs project knowledge invents its own approach — hardcoded tables, vague instructions, or nothing. Without a shared discovery mechanism, the pattern diverges further with each new skill.

**Concrete failure case**: During `/fix-bug` execution for an error-handling bug, the agent in `fix.md` received the instruction "Read relevant project knowledge files from `docs/business-rules/` and `docs/conventions/` based on the affected files and error context." With no guidance on *which* files or *how* to determine relevance, the agent loaded all 6 convention files — consuming ~2000 tokens of context on unrelated topics (profile-system, data-model) — while missing `business-rules/auth.md` which was directly relevant to the auth error. The result was a fix that violated the project's error-response conventions (documented in `conventions/error-handling.md`) because the agent treated the generic instruction as optional and skipped reading the most relevant file.

## Proposed Solution

Make convention/business-rule files **self-describing** via a `domains` field in YAML frontmatter. Embed a lightweight discovery instruction in each consumer: check what files exist, read their frontmatter, load what's relevant, skip the rest.

`/consolidate-specs` auto-generates `domains` when creating or updating files — users never manage it manually.

**Inspiration**: This approach adapts Hugo/Jekyll's frontmatter tag taxonomy (files self-describe their topic via metadata) and VS Code's context-based extension activation (load only what matches the current context). The key insight from these systems is that flat keyword metadata, combined with a lightweight consumer-side filter, eliminates the need for central indexing while remaining discoverable.

### Frontmatter Schema

```yaml
---
title: "Error Handling Conventions"  # existing field, unchanged
domains: [error, status, response, stderr]  # NEW: keywords this file covers
---
```

### Discovery Instruction

A brief instruction embedded in consumers, replacing hardcoded tables and vague "read relevant" text:

```
Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge relevant to this task.
Read each file's YAML frontmatter `domains` field to determine relevance.
Load files whose domains overlap with the task context.
If no files match, skip — no matching convention files for this task.
```

This is not a rigid protocol — it's guidance for the agent to discover user-provided knowledge. The agent uses its own judgment to match domains against task context. The tradeoff is explicit: we accept non-deterministic matching (agent may occasionally miss or over-include) in exchange for zero infrastructure and semantic flexibility. The failure mode is bounded — worst case, the agent loads an extra file or skips a relevant one, which is no worse than the current state of vague "read relevant" instructions. When deterministic matching is needed (e.g., skills with known dependencies like gen-test-scripts), the existing consumer-declared frontmatter pattern should be used instead.

### User-Facing Impact

After implementation, a human running `/fix-bug "TypeError in task claim"` will notice the agent's responses reference project-specific conventions (e.g., "Following the error handling convention in `error-handling.md`...") rather than producing generic code that contradicts established patterns. The user no longer needs to manually mention convention files in their prompt — the agent discovers them automatically. The observable change is fewer "I don't know the project conventions" failures and less back-and-forth to correct convention-violating code.

## Requirements Analysis

### Key Scenarios

- **Task execution (prompt template)**: Agent receives T-impl-3 via `enhancement.md`. Checks `docs/conventions/`, finds `error-handling.md` (domains: [error, status]) — matches task context, loads it. Skips `profile-system.md` (domains: [profile, config]) — no match.
- **Bug fix (command)**: User runs `/fix-bug "TypeError in task claim"`. Agent discovers `task-lifecycle.md` (domains: [task, claim, status]) and `error-handling.md` (domains: [error]) — both relevant, loads both.
- **Error fixer (agent)**: Subagent launched via `error-fixer.md`. Same discovery instruction, no duplicated mapping table.
- **User adds new convention**: User runs `/consolidate-specs`, which creates `docs/conventions/api.md` with `domains: [api, endpoint, route, rest]`. Next task automatically discovers it — no template/command updates needed.
- **No conventions exist**: New project with empty `docs/conventions/` and `docs/business-rules/`. Agent finds nothing to load, proceeds normally.
- **File without domains**: A convention file missing the `domains` field. Agent falls back to filename and title as relevance signals.
- **Ambiguous domain match**: Task involves error handling. Both `error-handling.md` (domains: [error, status, response]) and `error-reporting.md` (domains: [error, status, log]) match equally. Agent loads both — this is correct behavior since both files cover different aspects of the error domain (handling vs. reporting). Redundant overlap is acceptable; the cost is a few extra tokens for complementary knowledge.
- **Domain drift**: A convention file's content changes substantially after `domains` were written (e.g., `error-handling.md` gains a new section on retry logic). The stale `domains` field may not include "retry" — agent misses the new content for retry-related tasks. Handled by consolidate-specs drift detection (Step 9-11): when drift is detected and fixed, `domains` are re-derived from the updated content. No separate lifecycle step needed.

### Non-Functional Requirements

- **Token efficiency**: Agent reads only frontmatter (first ~5 lines) before deciding to load the full file
- **Graceful degradation**: Missing files, empty directories, or no matching domains never block task execution
- **Zero maintenance**: Adding new convention files requires no changes to any forge consumer

### Constraints & Dependencies

- Convention/business-rule files need `domains` frontmatter (managed by consolidate-specs, not users)
- Prompt templates are plain text — the discovery instruction must be embeddable as-is
- The mechanism is guidance for the agent, not a hard protocol — reliability depends on agent comprehension, not code-level guarantees

## Alternatives & Industry Benchmarking

### Industry References

The discovery problem — finding relevant context among a set of documents — is well-studied across domains:

- **RAG (Retrieval-Augmented Generation) systems**: Query the full document corpus via vector similarity, then inject top-K results into the prompt. Industry standard for LLM context retrieval (LangChain, LlamaIndex). Strength: semantic matching handles phrasing variations. Weakness: requires embedding model + vector store infrastructure. Overkill for <20 local Markdown files.
- **Static site generators (Hugo, Jekyll)**: Use YAML frontmatter tags/categories to organize and discover content. Hugo's `.Pages.GetPage` and Jekyll's `site.tags[tag]` use frontmatter metadata to filter pages at build time. This is the closest analogy — self-describing files with lightweight metadata, no indexing infrastructure.
- **Package registries (npm keywords, Cargo categories)**: Packages declare `keywords` in `package.json` or `categories` in `Cargo.toml`. Consumers search via keyword match. The pattern proves that flat keyword metadata scales to millions of packages with simple matching logic.
- **VS Code extension contributions**: Extensions declare `contributes` fields in `package.json` (languages, commands, debuggers). VS Code discovers and activates extensions based on workspace context — similar to our "load files whose domains overlap with task context."

### Comparison

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | Zero cost | Each consumer invents its own approach; knowledge remains undiscoverable | Rejected |
| RAG-style semantic retrieval | Handles phrasing variations; industry-proven at scale | Requires embedding model + vector store; heavy infrastructure for <20 local files | Rejected: disproportionate infrastructure cost |
| Hugo/Jekyll frontmatter tag indexing | Industry-validated pattern; self-describing files; zero infrastructure | Requires build step or code-level indexer | Partial fit: adopt the metadata pattern, skip the build step (agent reads frontmatter directly) |
| Consumer-declared dependencies (gen-test-* pattern, analogous to npm peerDependencies) | Explicit, reliable | Consumer must know filenames in advance; doesn't work for arbitrary-task prompt templates | Viable for domain-specific skills (already used), not for general consumers |
| **Frontmatter domain tags + agent discovery** | Self-describing (like Hugo tags); auto-managed; works across projects; agent uses semantic judgment (like VS Code context-based activation); no infrastructure required | Depends on agent correctly executing discovery; no code-level guarantee | **Selected: combines Hugo's metadata pattern with VS Code's context-driven activation, without requiring either's infrastructure** |

## Feasibility Assessment

### Technical Feasibility

All convention/business-rule files already have YAML frontmatter (`title` field). Adding `domains` is a mechanical extension. The discovery instruction is plain text embeddable in any consumer.

### Resource & Timeline

6 existing files need `domains` frontmatter. 5 prompt templates + 2 command/agent files need the discovery instruction. 1 skill (consolidate-specs) needs frontmatter management logic. All changes are small and well-bounded.

### Dependency Readiness

No external dependencies. All files are local and under version control.

## Scope

### In Scope

1. Define frontmatter `domains` schema for convention and business-rule files
2. Add `domains` frontmatter to all 6 existing files (3 conventions + 3 business-rules)
3. Update 5 prompt templates (`fix.md`, `cleanup.md`, `enhancement.md`, `feature.md`, `refactor.md`) — replace vague "read relevant" with discovery instruction
4. Update 2 command/agent files (`fix-bug.md`, `error-fixer.md`) — remove hardcoded mapping tables, replace with discovery instruction
5. Update `plugins/forge/skills/consolidate-specs/SKILL.md` — generate/maintain `domains` frontmatter when creating or updating files
6. Update `plugins/forge/hooks/guide.md` — update the project knowledge note to reference `domains` frontmatter

### Out of Scope

- Changes to `gen-test-scripts`, `gen-test-cases`, or `eval` skills — they use consumer-declared frontmatter and don't need discovery
- Changes to `gen-test-cases-per-type-dispatch` proposal or its implementation
- Creating new convention or business-rule files
- Changes to `run-tasks`, `breakdown-tasks`, or other orchestration skills
- CLI-level code changes (Go template variables, etc.)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent skips discovery step | M | M | If task output appears to ignore known project conventions (e.g., agent uses a pattern contradicted by a convention file), add a reminder line to the discovery instruction emphasizing "check conventions before writing code." Monitor across 5+ tasks and adjust wording if skip rate remains high. Graceful degradation: missing knowledge doesn't block execution |
| Agent loads irrelevant files | L | L | Acceptable: over-inclusion wastes tokens but doesn't break functionality. If token waste becomes measurable (>20% of convention tokens unused across 10+ tasks), narrow domain keywords or split overly broad files |
| `domains` tags too broad/narrow | M | M | consolidate-specs derives domains from spec ID and source keywords. Guidelines: 3-7 specific keywords. Fallback: agent uses filename/title if `domains` missing |
| consolidate-specs writes inaccurate domains | L | M | Domains derived from spec content. Human reviews via consolidate-specs confirmation flow |
| Domain keyword overlap between files | M | L | Two files claim overlapping domains (e.g., `error-handling.md` and `error-reporting.md` both have [error, status]). Agent loads both, consuming extra tokens for partially redundant content. Mitigation: consolidate-specs detects domain overlap >50% between files and warns the user during confirmation; the user can split or merge before committing |

## Success Criteria

- [ ] All 6 existing convention/business-rule files have valid `domains` frontmatter with 3-7 keywords each, where each keyword appears in the file's own content (source code identifier, file path, or spec term) at least once
- [ ] All 5 prompt templates contain the discovery instruction and no vague "read relevant" text
- [ ] Both command/agent files contain the discovery instruction and no hardcoded keyword→filename mapping tables
- [ ] `consolidate-specs` SKILL.md generates `domains` frontmatter when creating new files and updates it when drifting existing files
- [ ] `plugins/forge/hooks/guide.md` references `domains` frontmatter in the project knowledge note (scope item 6)
- [ ] For a task about "error handling", the agent loads `error-handling.md` and `error-reporting.md` but not `profile-system.md` or `quality-gate.md`
- [ ] An empty `docs/conventions/` directory doesn't block task execution
- [ ] No keyword→filename mapping tables remain in any command, agent, or prompt template

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
