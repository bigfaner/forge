---
title: Proposal Baseline Drift During Execution
domains: [proposal, execution, data-consistency]
---

# Gotcha: Proposal Baseline Drift During Execution

## Root Cause

When a feature's tasks modify the same files that the proposal uses as baseline data, the proposal's numbers become stale during execution. In skill-slimming, the proposal claimed 6394 lines across 22 SKILL.md files, but by the time eval-doc ran, tasks had already reduced the total to 4421 lines.

## Symptom

- Eval or validation tasks report data that contradicts the proposal
- Proposal success criteria reference incorrect baseline numbers
- Task descriptions show original line counts that no longer match reality

## Prevention

- Proposal baseline data is a snapshot at creation time — it should not be treated as live data
- If tasks modify files tracked in the proposal, the eval-doc task should verify actual numbers rather than trusting the proposal
- For features where tasks progressively modify the same artifact, consider updating the proposal after each tier/phase completes
