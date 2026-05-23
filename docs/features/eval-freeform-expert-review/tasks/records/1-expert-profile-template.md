---
status: "completed"
started: "2026-05-23 20:02"
completed: "2026-05-23 20:04"
time_spent: "~2m"
---

# Task Record: 1 Expert Profile Template & Inference Prompt

## Summary
Created expert profile template (expert-template.md) and expert inference prompt (expert-inference.md) for freeform expert review Phase 0

## Changes

### Files Created
- plugins/forge/skills/eval/experts/freeform/expert-template.md
- plugins/forge/skills/eval/experts/freeform/expert-inference.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
2 files created: expert-template.md (32 lines), expert-inference.md (111 lines)

## Referenced Documents
- docs/proposals/eval-freeform-expert-review/proposal.md
- plugins/forge/skills/eval/experts/scorer/cto.md
- plugins/forge/skills/eval/experts/scorer/architect.md
- plugins/forge/skills/eval/experts/scorer/qa.md
- plugins/forge/skills/eval/rules/scorer-composition.md
- plugins/forge/skills/eval/SKILL.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] Expert template contains 6 required YAML fields: domain, background, review_style, generated_for, created_at, review_history
- [x] Template format is YAML front matter + Markdown body, compatible with existing scorer/*.md prompt format
- [x] Inference prompt extracts domain, tech stack, complexity, key decisions from proposal
- [x] Inference prompt includes AskUserQuestion 3-option confirmation (Accept/Modify/Regenerate)
- [x] Inference prompt limits to 3 modification rounds
- [x] Inference prompt includes degradation logic: 3 consecutive rejections prompts manual input or skip
- [x] Expert profile includes verifiable domain keywords and background descriptions (anti-hallucination)
- [x] Inference prompt generates cross-reference (expert keywords vs proposal terms) and 3-5 self-check questions

## Notes
Expert template adds 'deprecated' field per proposal success criteria (expert deprecation mechanism). Inference prompt includes anti-hallucination safeguards (Step 4) with keyword grounding, background verifiability, and coverage ratio reporting.
