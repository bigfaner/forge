---
feature: "tui-ui-design"
sources:
  - docs/proposals/tui-ui-design/proposal.md
  - docs/features/tui-ui-design/tasks/1-tui-platform-themes.md
  - docs/features/tui-ui-design/tasks/2-prd-tui-navigation.md
  - docs/features/tui-ui-design/tasks/3-ui-design-tui-core.md
  - docs/features/tui-ui-design/tasks/4-tui-html-prototype.md
  - docs/features/tui-ui-design/tasks/5-eval-ui-multi-platform.md
  - docs/features/tui-ui-design/tasks/6-multi-platform-manifest.md
generated: "2026-05-15"
---

# Test Cases: tui-ui-design

> **Note**: Quick-mode feature -- no PRD files. Acceptance criteria extracted from task definitions and proposal success criteria.

> **WARNING**: sitemap.json not found -- Element set to `sitemap-missing`. Run `/gen-sitemap` for precise element references. (Skipped: profile has no web-ui capability.)

## Summary

| Type | Count |
|------|-------|
| TUI  | 0     |
| **Integration** | **0** |
| API  | 0     |
| CLI  | 31    |
| **Total** | **31** |

---

## CLI Test Cases

## TC-001: TUI platform file defines navigation structure
- **Source**: Task 1 AC-1
- **Type**: CLI
- **Target**: cli/platforms-tui
- **Test ID**: cli/platforms-tui/tui-platform-file-defines-navigation-structure
- **Pre-conditions**: Task 1 implementation complete; `platforms/tui.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/platforms/tui.md`
  2. Assert file contains "Keymap" table with columns `[Key | Action | Context/Mode]`
  3. Assert file contains "Panel Layout" table with columns `[Panel | View | Position | Size Hint]`
  4. Assert file contains "Modes" table with columns `[Mode | Description | Default Keybindings]`
  5. Assert file contains "Navigation Rules" section
- **Expected**: tui.md contains all four TUI navigation structures with correct table column headers
- **Priority**: P0

## TC-002: Modern Dark TUI theme specifies correct properties
- **Source**: Task 1 AC-2
- **Type**: CLI
- **Target**: cli/styles-modern-dark-tui
- **Test ID**: cli/styles-modern-dark-tui/modern-dark-tui-theme-specifies-correct-properties
- **Pre-conditions**: Task 1 implementation complete; `styles/modern-dark-tui.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/styles/modern-dark-tui.md`
  2. Assert file specifies color space as "256-color"
  3. Assert file specifies character set including box-drawing + block elements with examples (`▄▪─│┃`)
  4. Assert file specifies dark background palette with high contrast and green/red/blue semantic colors
  5. Assert file specifies density as "compact"
  6. Assert file includes applicable scenarios section
- **Expected**: modern-dark-tui.md contains all required theme properties per proposal D3
- **Priority**: P0

## TC-003: Minimal ASCII TUI theme specifies correct properties
- **Source**: Task 1 AC-3
- **Type**: CLI
- **Target**: cli/styles-minimal-ascii-tui
- **Test ID**: cli/styles-minimal-ascii-tui/minimal-ascii-tui-theme-specifies-correct-properties
- **Pre-conditions**: Task 1 implementation complete; `styles/minimal-ascii-tui.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/styles/minimal-ascii-tui.md`
  2. Assert file specifies color space as "16-color"
  3. Assert file specifies character set as pure ASCII (`#=-\|*+.`)
  4. Assert file specifies default terminal background palette with weight/spacing distinction
  5. Assert file specifies density as "loose"
  6. Assert file includes applicable scenarios section
- **Expected**: minimal-ascii-tui.md contains all required theme properties per proposal D3
- **Priority**: P0

## TC-004: TUI theme files follow existing style file format
- **Source**: Task 1 AC-4
- **Type**: CLI
- **Target**: cli/styles-tui-format
- **Test ID**: cli/styles-tui-format/tui-theme-files-follow-existing-style-file-format
- **Pre-conditions**: Task 1 implementation complete; existing style files (e.g., `apple.md`) available as reference
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `plugins/forge/skills/ui-design/templates/styles/apple.md` to identify format structure
  2. Read `plugins/forge/skills/ui-design/templates/styles/modern-dark-tui.md`
  3. Read `plugins/forge/skills/ui-design/templates/styles/minimal-ascii-tui.md`
  4. Assert both TUI theme files have frontmatter matching the existing style file format
  5. Assert both TUI theme files have section headings matching the established pattern
- **Expected**: TUI theme files use the same structural format as existing web/mobile style files
- **Priority**: P1

## TC-005: PRD UI functions template includes TUI Navigation Architecture
- **Source**: Task 2 AC-1
- **Type**: CLI
- **Target**: cli/prd-ui-functions-tui
- **Test ID**: cli/prd-ui-functions-tui/prd-ui-functions-template-includes-tui-navigation-architecture
- **Pre-conditions**: Task 2 implementation complete; `prd-ui-functions.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/write-prd/templates/prd-ui-functions.md`
  2. Assert file contains "Platform: tui" indicator
  3. Assert file contains Keymap table with columns `[Key | Action | Context/Mode]`
  4. Assert file contains Panel Layout table with columns `[Panel | View | Position | Size Hint]`
  5. Assert file contains Modes table with columns `[Mode | Description | Default Keybindings]`
  6. Assert file contains Navigation Rules section for TUI
- **Expected**: prd-ui-functions.md includes TUI Navigation Architecture section per proposal D4
- **Priority**: P0

## TC-006: TUI navigation section is conditionally rendered
- **Source**: Task 2 AC-2
- **Type**: CLI
- **Target**: cli/prd-ui-functions-conditional
- **Test ID**: cli/prd-ui-functions-conditional/tui-navigation-section-is-conditionally-rendered
- **Pre-conditions**: Task 2 implementation complete; `prd-ui-functions.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/write-prd/templates/prd-ui-functions.md`
  2. Assert TUI navigation section is guarded by a platform=tui condition (conditional marker or template variable)
  3. Assert web navigation section is guarded by platform=web condition
  4. Assert mobile navigation section is guarded by platform=mobile condition
  5. Verify TUI section does not appear in the web/mobile rendering path
- **Expected**: TUI navigation section only renders when platform=tui, web/mobile templates unaffected
- **Priority**: P0

## TC-007: write-prd SKILL.md references TUI navigation template
- **Source**: Task 2 AC-3
- **Type**: CLI
- **Target**: cli/write-prd-skill-tui
- **Test ID**: cli/write-prd-skill-tui/write-prd-skill-references-tui-navigation-template
- **Pre-conditions**: Task 2 implementation complete; `write-prd/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/write-prd/SKILL.md`
  2. Assert SKILL.md contains reference to TUI navigation rendering
  3. Assert SKILL.md contains platform=tui detection logic
  4. Assert SKILL.md references the TUI navigation section in prd-ui-functions.md template
- **Expected**: write-prd SKILL.md handles platform=tui and triggers TUI navigation section
- **Priority**: P0

## TC-008: Existing web/mobile PRD generation behavior unchanged
- **Source**: Task 2 AC-4
- **Type**: CLI
- **Target**: cli/write-prd-regression
- **Test ID**: cli/write-prd-regression/existing-web-mobile-prd-generation-behavior-unchanged
- **Pre-conditions**: Task 2 implementation complete
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `plugins/forge/skills/write-prd/templates/prd-ui-functions.md`
  2. Assert web Navigation Architecture section content is unchanged from pre-task state
  3. Assert mobile Navigation Architecture section content is unchanged from pre-task state
  4. Verify no TUI-specific content leaks into web/mobile sections
- **Expected**: Web and mobile PRD template content identical to pre-task state
- **Priority**: P0

## TC-009: ui-design SKILL.md detects platform=tui
- **Source**: Task 3 AC-1
- **Type**: CLI
- **Target**: cli/ui-design-skill-tui-detection
- **Test ID**: cli/ui-design-skill-tui-detection/ui-design-skill-detects-platform-tui
- **Pre-conditions**: Task 3 implementation complete; `ui-design/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/SKILL.md`
  2. Assert SKILL.md contains platform detection logic that identifies platform=tui from PRD
  3. Assert SKILL.md has a distinct TUI branch/flow separate from web and mobile
- **Expected**: SKILL.md detects platform=tui and enters TUI-specific processing branch
- **Priority**: P0

## TC-010: TUI branch presents theme selection
- **Source**: Task 3 AC-2
- **Type**: CLI
- **Target**: cli/ui-design-theme-selection
- **Test ID**: cli/ui-design-theme-selection/tui-branch-presents-theme-selection
- **Pre-conditions**: Task 3 implementation complete; `ui-design/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/SKILL.md`
  2. Assert TUI branch offers three theme options: Modern Dark, Minimal ASCII, DESIGN.md custom
  3. Assert theme selection is prompted during TUI flow
- **Expected**: SKILL.md TUI branch presents all three theme choices per proposal D3
- **Priority**: P0

## TC-011: TUI branch uses TUI platform and theme files
- **Source**: Task 3 AC-3
- **Type**: CLI
- **Target**: cli/ui-design-tui-resources
- **Test ID**: cli/ui-design-tui-resources/tui-branch-uses-tui-platform-and-theme-files
- **Pre-conditions**: Task 3 implementation complete; `ui-design/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/SKILL.md`
  2. Assert TUI branch references `platforms/tui.md` for navigation rules
  3. Assert TUI branch references the selected theme file for visual style
  4. Assert TUI branch reads theme properties from styles directory
- **Expected**: TUI branch loads and uses platforms/tui.md and the selected theme file
- **Priority**: P0

## TC-012: ui-design template includes TUI component template with 5 structural requirements
- **Source**: Task 3 AC-4
- **Type**: CLI
- **Target**: cli/ui-design-tui-template
- **Test ID**: cli/ui-design-tui-template/ui-design-template-includes-tui-component-template-with-5-structural-requirements
- **Pre-conditions**: Task 3 implementation complete; `ui-design.md` template modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/ui-design.md`
  2. Assert template contains TUI component section with "Panel Placement" subsection
  3. Assert template contains "ASCII Layout Mockup" subsection
  4. Assert template contains "Dimensions" subsection (concrete numbers)
  5. Assert template contains "Character Palette" subsection (Unicode + reason)
  6. Assert template contains "Color Mapping" subsection (fg/bg color codes)
  7. Assert template contains "Edge Cases" subsection (5 mandatory scenarios)
- **Expected**: ui-design.md template has TUI component template with all 5 structural requirements from lesson
- **Priority**: P0

## TC-013: Multi-platform features produce separate ui-design files
- **Source**: Task 3 AC-5
- **Type**: CLI
- **Target**: cli/ui-design-multi-platform
- **Test ID**: cli/ui-design-multi-platform/multi-platform-features-produce-separate-ui-design-files
- **Pre-conditions**: Task 3 implementation complete; `ui-design/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/SKILL.md`
  2. Assert SKILL.md multi-platform logic produces separate output files per platform
  3. Assert output filenames follow pattern `ui-design-web.md` + `ui-design-tui.md` per proposal D7
- **Expected**: Multi-platform feature produces `ui-design-web.md` and `ui-design-tui.md` as separate files
- **Priority**: P0

## TC-014: Single TUI feature produces ui-design-tui.md
- **Source**: Task 3 AC-6
- **Type**: CLI
- **Target**: cli/ui-design-single-tui
- **Test ID**: cli/ui-design-single-tui/single-tui-feature-produces-ui-design-tui-md
- **Pre-conditions**: Task 3 implementation complete; `ui-design/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/SKILL.md`
  2. Assert single TUI platform feature produces `ui-design-tui.md` (not `ui-design-web.md`)
  3. Assert output goes to correct feature directory
- **Expected**: Single TUI feature produces `ui-design-tui.md`
- **Priority**: P0

## TC-015: Existing web/mobile ui-design behavior unchanged
- **Source**: Task 3 AC-7
- **Type**: CLI
- **Target**: cli/ui-design-regression
- **Test ID**: cli/ui-design-regression/existing-web-mobile-ui-design-behavior-unchanged
- **Pre-conditions**: Task 3 implementation complete
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `plugins/forge/skills/ui-design/SKILL.md`
  2. Assert web platform branch logic is unchanged from pre-task state
  3. Assert mobile platform branch logic is unchanged from pre-task state
  4. Read `plugins/forge/skills/ui-design/templates/ui-design.md`
  5. Assert web/mobile component template sections are unchanged
- **Expected**: Web and mobile ui-design behavior identical to pre-task state
- **Priority**: P0

## TC-016: Prototype template includes TUI-specific generation rules
- **Source**: Task 4 AC-1
- **Type**: CLI
- **Target**: cli/prototype-tui-rules
- **Test ID**: cli/prototype-tui-rules/prototype-template-includes-tui-specific-generation-rules
- **Pre-conditions**: Task 4 implementation complete; `prototype.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/prototype.md`
  2. Assert file contains TUI-specific prototype generation rules section
  3. Assert TUI rules are distinct from web/mobile prototype rules
- **Expected**: prototype.md includes TUI prototype generation rules
- **Priority**: P0

## TC-017: TUI prototype is single index.html with terminal window div
- **Source**: Task 4 AC-2
- **Type**: CLI
- **Target**: cli/prototype-tui-structure
- **Test ID**: cli/prototype-tui-structure/tui-prototype-is-single-index-html-with-terminal-window-div
- **Pre-conditions**: Task 4 implementation complete; `prototype.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/prototype.md`
  2. Assert TUI prototype rules specify a single `index.html` output file
  3. Assert rules specify a "terminal window" div container for all panel rendering
  4. Assert rules specify all panels rendered inside the terminal window div
- **Expected**: TUI prototype rules describe single index.html with terminal-window div containing all panels
- **Priority**: P0

## TC-018: TUI prototype includes simulated key buttons
- **Source**: Task 4 AC-3
- **Type**: CLI
- **Target**: cli/prototype-tui-keys
- **Test ID**: cli/prototype-tui-keys/tui-prototype-includes-simulated-key-buttons
- **Pre-conditions**: Task 4 implementation complete; `prototype.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/prototype.md`
  2. Assert TUI prototype rules specify simulated key buttons at the bottom
  3. Assert rules include buttons: `[Tab]`, `[1]`, `[q]`, `[:command]`
  4. Assert buttons trigger panel switching
- **Expected**: TUI prototype has simulated key buttons for panel navigation per proposal D6
- **Priority**: P1

## TC-019: TUI prototype uses monospace font and dark background
- **Source**: Task 4 AC-4
- **Type**: CLI
- **Target**: cli/prototype-tui-styling
- **Test ID**: cli/prototype-tui-styling/tui-prototype-uses-monospace-font-and-dark-background
- **Pre-conditions**: Task 4 implementation complete; `prototype.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/prototype.md`
  2. Assert TUI prototype rules specify monospace font (e.g., Menlo/Consolas)
  3. Assert rules specify dark background color (e.g., #1e1e1e)
  4. Assert rules specify fixed-width character rendering
- **Expected**: TUI prototype CSS simulates terminal appearance with monospace font and dark background
- **Priority**: P1

## TC-020: TUI prototype panel layout matches ASCII mockup
- **Source**: Task 4 AC-5
- **Type**: CLI
- **Target**: cli/prototype-tui-layout-match
- **Test ID**: cli/prototype-tui-layout-match/tui-prototype-panel-layout-matches-ascii-mockup
- **Pre-conditions**: Task 4 implementation complete; `prototype.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/prototype.md`
  2. Assert TUI prototype rules state that HTML panel layout must match ASCII mockup from ui-design.md
  3. Assert rules reference ui-design.md dimensions as the source of truth
- **Expected**: Prototype template explicitly requires HTML layout to match ASCII mockup layout
- **Priority**: P0

## TC-021: TUI prototypes output to correct directories
- **Source**: Task 4 AC-6
- **Type**: CLI
- **Target**: cli/prototype-tui-output-dir
- **Test ID**: cli/prototype-tui-output-dir/tui-prototypes-output-to-correct-directories
- **Pre-conditions**: Task 4 implementation complete; `prototype.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/prototype.md`
  2. Assert TUI prototype output path for multi-platform is `prototype/tui/`
  3. Assert TUI prototype output path for single TUI feature is `prototype/`
- **Expected**: TUI prototype outputs to `prototype/tui/` (multi-platform) or `prototype/` (single TUI)
- **Priority**: P1

## TC-022: rubric-web.md contains existing web rubric content
- **Source**: Task 5 AC-1
- **Type**: CLI
- **Target**: cli/eval-ui-rubric-web
- **Test ID**: cli/eval-ui-rubric-web/rubric-web-md-contains-existing-web-rubric-content
- **Pre-conditions**: Task 5 implementation complete; `rubric-web.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/eval-ui/templates/rubric-web.md`
  2. Assert file contains rubric content (4 dimensions, scoring criteria)
  3. Assert total score equals 1000 points
- **Expected**: rubric-web.md contains valid web evaluation rubric content
- **Priority**: P0

## TC-023: rubric-tui.md has 4 correct dimensions
- **Source**: Task 5 AC-2
- **Type**: CLI
- **Target**: cli/eval-ui-rubric-tui
- **Test ID**: cli/eval-ui-rubric-tui/rubric-tui-md-has-4-correct-dimensions
- **Pre-conditions**: Task 5 implementation complete; `rubric-tui.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/eval-ui/templates/rubric-tui.md`
  2. Assert file contains "Requirement Coverage" dimension with 250 points
  3. Assert file contains "Terminal Experience" dimension with 250 points
  4. Assert file contains "Visual Specification" dimension with 250 points
  5. Assert file contains "Implementability" dimension with 250 points
  6. Assert total score equals 1000 points
- **Expected**: rubric-tui.md has exactly 4 dimensions at 250 points each per proposal D9
- **Priority**: P0

## TC-024: rubric-tui.md deduction rules are correct
- **Source**: Task 5 AC-3
- **Type**: CLI
- **Target**: cli/eval-ui-rubric-tui-deductions
- **Test ID**: cli/eval-ui-rubric-tui-deductions/rubric-tui-md-deduction-rules-are-correct
- **Pre-conditions**: Task 5 implementation complete; `rubric-tui.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/eval-ui/templates/rubric-tui.md`
  2. Assert deduction rule: missing ASCII mockup sets panel Visual Specification to 0
  3. Assert deduction rule: characters marked as pending/unspecified incur -30 per instance
  4. Assert deduction rule: missing mandatory edge case incurs -50 per case
  5. Assert deduction rule: vague dimensions incur -20 per instance
- **Expected**: rubric-tui.md contains all 4 deduction rules per proposal D9
- **Priority**: P0

## TC-025: rubric-mobile.md has 4 correct dimensions
- **Source**: Task 5 AC-4
- **Type**: CLI
- **Target**: cli/eval-ui-rubric-mobile
- **Test ID**: cli/eval-ui-rubric-mobile/rubric-mobile-md-has-4-correct-dimensions
- **Pre-conditions**: Task 5 implementation complete; `rubric-mobile.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/eval-ui/templates/rubric-mobile.md`
  2. Assert file contains "Requirement Coverage" dimension with 250 points
  3. Assert file contains "Touch Experience" dimension with 250 points
  4. Assert file contains "Adaptive Layout" dimension with 250 points
  5. Assert file contains "Implementability" dimension with 250 points
  6. Assert total score equals 1000 points
- **Expected**: rubric-mobile.md has exactly 4 dimensions at 250 points each per proposal D10
- **Priority**: P0

## TC-026: rubric-mobile.md deduction rules are correct
- **Source**: Task 5 AC-5
- **Type**: CLI
- **Target**: cli/eval-ui-rubric-mobile-deductions
- **Test ID**: cli/eval-ui-rubric-mobile-deductions/rubric-mobile-md-deduction-rules-are-correct
- **Pre-conditions**: Task 5 implementation complete; `rubric-mobile.md` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/eval-ui/templates/rubric-mobile.md`
  2. Assert deduction rule: touch targets without size annotation incur -30 per instance
  3. Assert deduction rule: missing landscape/portrait adaptation incurs -50
  4. Assert deduction rule: missing safe area handling incurs -40
- **Expected**: rubric-mobile.md contains all 3 deduction rules per proposal D10
- **Priority**: P0

## TC-027: eval-ui SKILL.md detects platform and selects rubric
- **Source**: Task 5 AC-6
- **Type**: CLI
- **Target**: cli/eval-ui-skill-platform
- **Test ID**: cli/eval-ui-skill-platform/eval-ui-skill-detects-platform-and-selects-rubric
- **Pre-conditions**: Task 5 implementation complete; `eval-ui/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/eval-ui/SKILL.md`
  2. Assert SKILL.md detects platform from ui-design document
  3. Assert platform=web selects `rubric-web.md`
  4. Assert platform=mobile selects `rubric-mobile.md`
  5. Assert platform=tui selects `rubric-tui.md`
- **Expected**: eval-ui SKILL.md selects correct rubric file based on detected platform per proposal D8
- **Priority**: P0

## TC-028: Multi-platform features evaluate with respective rubrics
- **Source**: Task 5 AC-7
- **Type**: CLI
- **Target**: cli/eval-ui-multi-platform
- **Test ID**: cli/eval-ui-multi-platform/multi-platform-features-evaluate-with-respective-rubrics
- **Pre-conditions**: Task 5 implementation complete; `eval-ui/SKILL.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/eval-ui/SKILL.md`
  2. Assert multi-platform logic evaluates each platform's ui-design file independently
  3. Assert each platform uses its respective rubric (web rubric for web, tui rubric for tui)
- **Expected**: Multi-platform features run separate evaluations per platform with matched rubrics
- **Priority**: P1

## TC-029: Single-platform web manifest unchanged
- **Source**: Task 6 AC-1
- **Type**: CLI
- **Target**: cli/manifest-single-web
- **Test ID**: cli/manifest-single-web/single-platform-web-manifest-unchanged
- **Pre-conditions**: Task 6 implementation complete; `manifest-update-ui.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/manifest-update-ui.md`
  2. Assert template lists `ui-design.md` for single web platform
  3. Assert template lists `prototype/` directory for single web platform
- **Expected**: Single-platform web manifest lists ui-design.md and prototype/ unchanged
- **Priority**: P0

## TC-030: Multi-platform manifest lists platform-specific files
- **Source**: Task 6 AC-2
- **Type**: CLI
- **Target**: cli/manifest-multi-platform
- **Test ID**: cli/manifest-multi-platform/multi-platform-manifest-lists-platform-specific-files
- **Pre-conditions**: Task 6 implementation complete; `manifest-update-ui.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/manifest-update-ui.md`
  2. Assert multi-platform section lists `ui-design-web.md`
  3. Assert multi-platform section lists `ui-design-tui.md`
  4. Assert multi-platform section lists `prototype/web/`
  5. Assert multi-platform section lists `prototype/tui/`
- **Expected**: Multi-platform manifest lists all platform-specific files per proposal D7
- **Priority**: P0

## TC-031: Single TUI manifest lists correct files
- **Source**: Task 6 AC-3
- **Type**: CLI
- **Target**: cli/manifest-single-tui
- **Test ID**: cli/manifest-single-tui/single-tui-manifest-lists-correct-files
- **Pre-conditions**: Task 6 implementation complete; `manifest-update-ui.md` modified
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read file `plugins/forge/skills/ui-design/templates/manifest-update-ui.md`
  2. Assert single TUI section lists `ui-design-tui.md`
  3. Assert single TUI section lists `prototype/`
- **Expected**: Single TUI manifest lists ui-design-tui.md and prototype/
- **Priority**: P0

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Task 1 AC-1 | CLI | cli/platforms-tui | P0 |
| TC-002 | Task 1 AC-2 | CLI | cli/styles-modern-dark-tui | P0 |
| TC-003 | Task 1 AC-3 | CLI | cli/styles-minimal-ascii-tui | P0 |
| TC-004 | Task 1 AC-4 | CLI | cli/styles-tui-format | P1 |
| TC-005 | Task 2 AC-1 | CLI | cli/prd-ui-functions-tui | P0 |
| TC-006 | Task 2 AC-2 | CLI | cli/prd-ui-functions-conditional | P0 |
| TC-007 | Task 2 AC-3 | CLI | cli/write-prd-skill-tui | P0 |
| TC-008 | Task 2 AC-4 | CLI | cli/write-prd-regression | P0 |
| TC-009 | Task 3 AC-1 | CLI | cli/ui-design-skill-tui-detection | P0 |
| TC-010 | Task 3 AC-2 | CLI | cli/ui-design-theme-selection | P0 |
| TC-011 | Task 3 AC-3 | CLI | cli/ui-design-tui-resources | P0 |
| TC-012 | Task 3 AC-4 | CLI | cli/ui-design-tui-template | P0 |
| TC-013 | Task 3 AC-5 | CLI | cli/ui-design-multi-platform | P0 |
| TC-014 | Task 3 AC-6 | CLI | cli/ui-design-single-tui | P0 |
| TC-015 | Task 3 AC-7 | CLI | cli/ui-design-regression | P0 |
| TC-016 | Task 4 AC-1 | CLI | cli/prototype-tui-rules | P0 |
| TC-017 | Task 4 AC-2 | CLI | cli/prototype-tui-structure | P0 |
| TC-018 | Task 4 AC-3 | CLI | cli/prototype-tui-keys | P1 |
| TC-019 | Task 4 AC-4 | CLI | cli/prototype-tui-styling | P1 |
| TC-020 | Task 4 AC-5 | CLI | cli/prototype-tui-layout-match | P0 |
| TC-021 | Task 4 AC-6 | CLI | cli/prototype-tui-output-dir | P1 |
| TC-022 | Task 5 AC-1 | CLI | cli/eval-ui-rubric-web | P0 |
| TC-023 | Task 5 AC-2 | CLI | cli/eval-ui-rubric-tui | P0 |
| TC-024 | Task 5 AC-3 | CLI | cli/eval-ui-rubric-tui-deductions | P0 |
| TC-025 | Task 5 AC-4 | CLI | cli/eval-ui-rubric-mobile | P0 |
| TC-026 | Task 5 AC-5 | CLI | cli/eval-ui-rubric-mobile-deductions | P0 |
| TC-027 | Task 5 AC-6 | CLI | cli/eval-ui-skill-platform | P0 |
| TC-028 | Task 5 AC-7 | CLI | cli/eval-ui-multi-platform | P1 |
| TC-029 | Task 6 AC-1 | CLI | cli/manifest-single-web | P0 |
| TC-030 | Task 6 AC-2 | CLI | cli/manifest-multi-platform | P0 |
| TC-031 | Task 6 AC-3 | CLI | cli/manifest-single-tui | P0 |
