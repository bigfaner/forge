---
title: "Profile System Conventions"
domains: [profile, config, e2e, framework, registry]
---

# Profile System Conventions

_Source: feature/forge-cli-v3_

## Configuration

### TECH-profile-system-001: Profile-Aware E2E Architecture

**Requirement**: E2E commands MUST read the `profile` field from `.forge/config.yaml` to determine which test suite to execute. Profile detection scans project structure for framework-specific config files (e.g., `playwright.config.ts` -> web-playwright). Supported profiles are defined in a registry; unknown profile values MUST be rejected.
**Source**: feature/forge-cli-v3 TECH-004
