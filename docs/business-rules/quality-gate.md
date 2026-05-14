---
title: "Quality Gate Rules"
---

# Quality Gate Rules

_Source: feature/forge-cli-v3_

## Pipeline

### BIZ-quality-gate-001: Quality Gate Sequential Pipeline

**Rule**: `forge quality-gate` executes a sequential pipeline: compile -> fmt -> lint -> test. The first failing step terminates the pipeline with exit code 1. All steps passing yields exit code 0.
**Context**: Provides a single-command CI check that enforces code quality in order of dependency (code must compile before it can be linted).
**Source**: feature/forge-cli-v3 BIZ-004
