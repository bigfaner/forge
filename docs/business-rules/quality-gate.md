---
title: "Quality Gate Rules"
---

# Quality Gate Rules

_Source: feature/forge-cli-v3_

## Pipeline

### BIZ-quality-gate-001: Quality Gate Multi-Phase Pipeline

**Rule**: `forge quality-gate` executes a multi-phase pipeline after all tasks complete: (1) compile -> fmt -> lint gate, (2) project-wide unit/integration tests, (3) e2e regression. Each failing step can auto-create a fix task (P0, breaking) and outputs a hook JSON block reason. Exit code 0 in all cases (hook JSON signals the actual decision). The submit command uses a separate inline gate: compile -> fmt -> lint -> test, where failure exits with code 1.
**Context**: Provides project health enforcement after all feature tasks complete. The quality-gate hook evolved from a single pipeline into a multi-phase gate with fix-task auto-creation and e2e regression support.
**Source**: feature/forge-cli-v3 BIZ-004
