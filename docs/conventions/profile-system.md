---
title: "Testing System Conventions"
---

# Testing System Conventions

_Source: feature/simplify-testing-model_

## Configuration

### TECH-testing-system-001: Language-Aware E2E Architecture

**Requirement**: E2E commands MUST auto-detect the project's language from project root files (`go.mod` -> Go, `package.json` + `@playwright/test` -> JavaScript, etc.) to determine which test strategy to execute. The `languages` field in `.forge/config.yaml` overrides auto-detection when set. Supported languages are defined in the `languageCapabilities` map; unknown language values MUST be rejected.
**Source**: feature/simplify-testing-model TECH-001 (migrated from feature/forge-cli-v3 TECH-004)
