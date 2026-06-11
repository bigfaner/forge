---
status: "completed"
started: "2026-05-17 13:25"
completed: "2026-05-17 13:29"
time_spent: "~4m"
---

# Task Record: 5 Create types/mobile.md instruction file

## Summary
Created plugins/forge/skills/gen-test-scripts/types/mobile.md with all required sections: Reconnaissance Strategy (mobile-specific search patterns for app manifests, Maestro flow configs, accessibility labels), Fact Table Required Keys (MOBILE_APP_ID, MOBILE_FRAMEWORK, MOBILE_SCREEN_*), Verification Method (detect mobile framework and app manifests), Generation Patterns (touch/gesture simulation, screen transition assertions, app lifecycle events, element location via accessibility labels, Maestro YAML flow skeleton), and Mobile Antipattern Guards (device-dependent tests, hardcoded wait durations, physical device requirements). Frontmatter declares type: mobile and conventions: [testing-mobile.md]. Grounded in the maestro profile's YAML template format.

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/types/mobile.md

### Files Modified
无

### Key Decisions
- Grounded generation patterns in the maestro profile's YAML flow template structure (appId, onFlowStart, onFlowEnd, commands) rather than inventing abstract patterns
- Made MOBILE_APP_ID the minimum required Fact Table key since Maestro cannot launch the app without it
- Kept the file minimal but structurally complete per Hard Rules -- mobile testing is not actively used

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] plugins/forge/skills/gen-test-scripts/types/mobile.md exists
- [x] Frontmatter declares type: mobile and conventions: [testing-mobile.md]
- [x] Contains Reconnaissance Strategy section with mobile-specific search patterns
- [x] Contains Fact Table Required Keys section listing minimum keys for mobile type
- [x] Contains Verification Method section describing how to confirm project is a mobile app
- [x] Contains Generation Patterns section describing mobile test case to script translation
- [x] Contains Mobile Antipattern Guards section
- [x] At least 3 section headings are unique to this file

## Notes
Documentation-only task, no test metrics applicable. File follows same structure as cli.md/api.md/tui.md type files.
