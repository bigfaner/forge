---
title: "Error Handling Conventions"
---

# Error Handling Conventions

_Source: feature/forge-cli-v3_

## Output Channels

### TECH-error-handling-001: Stderr-Only Error Output Pattern

**Requirement**: All error diagnostics MUST be written to stderr, never to stdout. stdout is reserved for command output data. Error messages MUST follow the format: `<context>: <specific-detail>` (e.g., "task not found: T-impl-1", "unknown profile: bad-value").
**Source**: feature/forge-cli-v3 TECH-006
