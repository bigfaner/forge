---
created: "2026-06-08"
tags: [architecture, testing, gotcha]
---

# Pipeline-modification features don't fit the surface-based test pipeline

## Problem

When a feature modifies the forge pipeline itself (e.g., `GenerateTestTasks` dependency wiring, auto-gen templates, prompt templates), `forge task index` still generates T-test-gen-journeys → T-test-gen-contracts → T-test-gen-scripts → T-test-run for the configured surface. These tasks are designed to test user-facing runtime behavior, but pipeline-modification features produce no runnable surface to test against.

## Root Cause

1. `forge task index` generates test pipeline tasks based on `.forge/config.yaml` surfaces, regardless of feature content.
2. The surface-based test pipeline assumes the feature adds/modifies user-facing behavior that can be exercised via CLI/API/Web.
3. Pipeline-modification features (Go code in `pipeline.go`, `.md` templates) don't expose any surface behavior — they change how tasks are generated, not what the application does at runtime.
4. Running CLI functional tests on a pipeline-modification feature would test the forge binary itself, not the feature's changes.

## Solution

For features that modify the forge pipeline/build system:
- **Unit tests** for pipeline logic (Go `*_test.go`) — already the primary coverage mechanism.
- **Template content verification** — check that `{{.AcceptanceCriteria}}` renders correctly in generated .md files.
- **Integration verification** — run `forge task index` with multi-surface config and inspect `index.json` for correct dependency chains.

Consider skipping or marking auto-generated test pipeline tasks (T-test-gen-journeys through T-test-run) as not applicable for pure pipeline-modification features.

## Reusable Pattern

When classifying a feature, check whether the output is "forge pipeline changes" vs "application behavior changes". If the former, the standard surface-based test pipeline is a mismatch — use unit tests + structural verification instead.
