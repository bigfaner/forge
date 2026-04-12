# Directory Structure Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate zcode from flat file structure (`prd.md`, `design.md`) to nested subdirectory structure with source-prefixed filenames and manifest-based traceability.

**Architecture:** Go task-cli constants/paths layer drives all file location logic. Skill markdown files (SKILL.md, templates) consume these paths. A new `manifest.md` at feature root provides single-entry-point context loading and traceability mapping for AI agents.

**Tech Stack:** Go 1.21 (task-cli), Markdown (skills/templates), JSON (index.json schema)

---

## File Structure Map

### Go Source (task-cli)

| File | Action | Responsibility |
|------|--------|---------------|
| `task-cli/pkg/feature/constants.go` | Modify | Rename old constants, add new constants for subdirectories and manifest |
| `task-cli/pkg/feature/paths.go` | Modify | Rename old functions, add new path functions for subdirs, manifest, proposals |
| `task-cli/pkg/feature/feature.go` | Modify | Update `EnsureFeatureDir` to create `prd/`, `design/`, `ui/`, `tasks/` subdirs |
| `task-cli/pkg/feature/paths_test.go` | Modify | Update expected paths from flat to nested |
| `task-cli/pkg/feature/feature_test.go` | Modify | Update expected dirs created by `EnsureFeatureDir` |
| `task-cli/pkg/task/types.go` | Modify | Update `TaskIndex` field comments (PRD/Design path values) |
| `task-cli/pkg/task/types_test.go` | Modify | Update `"prd.md"` → `"prd/prd-spec.md"`, `"design.md"` → `"design/tech-design.md"` |
| `task-cli/internal/cmd/check_test.go` | Modify | Update PRD/Design field values |
| `task-cli/internal/cmd/claim_test.go` | Modify | Update PRD/Design field values |
| `task-cli/internal/cmd/feature_test.go` | Modify | Update PRD/Design field values |
| `task-cli/internal/cmd/runners_test.go` | Modify | Update PRD/Design field values |
| `task-cli/internal/cmd/validate_test.go` | Modify | Update PRD/Design field values |
| `task-cli/docs/OVERVIEW.md` | Modify | Update directory structure diagram |

### Skill Files (plugins/zcode)

| File | Action | Responsibility |
|------|--------|---------------|
| `plugins/zcode/skills/write-prd/SKILL.md` | Modify | New output paths (`prd/prd-spec.md`, `prd/prd-user-stories.md`, `prd/prd-ui-functions.md`), manifest creation |
| `plugins/zcode/skills/write-prd/templates/prd.md` | Rename → `prd-spec.md` | Rename template file |
| `plugins/zcode/skills/write-prd/templates/user-stories.md` | Rename → `prd-user-stories.md` | Rename template file |
| `plugins/zcode/skills/write-prd/templates/prd-ui-functions.md` | Create | NEW: UI functions requirements template |
| `plugins/zcode/skills/write-prd/templates/manifest.md` | Create | NEW: Manifest template with PRD section |
| `plugins/zcode/skills/design-tech/SKILL.md` | Modify | New input/output paths (`prd/prd-spec.md` → `design/tech-design.md`, `design/api-handbook.md`) |
| `plugins/zcode/skills/design-tech/templates/design.md` | Rename → `tech-design.md` | Rename template file |
| `plugins/zcode/skills/design-tech/templates/api-handbook.md` | Create | NEW: API handbook template |
| `plugins/zcode/skills/design-tech/templates/manifest-update-design.md` | Create | NEW: Manifest update snippet |
| `plugins/zcode/skills/eval-prd/SKILL.md` | Modify | Locate via manifest, new paths |
| `plugins/zcode/skills/eval-design/SKILL.md` | Modify | Locate via manifest, new paths |
| `plugins/zcode/skills/breakdown-tasks/SKILL.md` | Modify | Read manifest → all docs, new paths |
| `plugins/zcode/skills/breakdown-tasks/templates/index.json` | Modify | Update `prd`/`design` field values |
| `plugins/zcode/skills/breakdown-tasks/templates/manifest-update-tasks.md` | Create | NEW: Manifest update snippet |
| `plugins/zcode/skills/brainstorm/SKILL.md` | Create | NEW: Brainstorm skill |
| `plugins/zcode/skills/brainstorm/templates/proposal.md` | Create | NEW: Proposal template |
| `plugins/zcode/skills/ui-design/SKILL.md` | Create | NEW: UI design skill |
| `plugins/zcode/skills/ui-design/templates/ui-design.md` | Create | NEW: UI design template |
| `plugins/zcode/skills/ui-design/templates/manifest-update-ui.md` | Create | NEW: Manifest update snippet |

### Config & Docs

| File | Action | Responsibility |
|------|--------|---------------|
| `plugins/zcode/hooks/guide.md` | Modify | Update Document Index with new structure + manifest |
| `plugins/zcode/.claude-plugin/plugin.json` | Modify | Bump to v2.0.0, add keywords |
| `docs/zcode-redesign-plan.md` | Modify | Align parent plan with spec's directory structure |

---

### Task 1: Update Go Constants

**Files:**
- Modify: `task-cli/pkg/feature/constants.go`
- Test: `task-cli/pkg/feature/paths_test.go`

- [ ] **Step 1: Write the failing test for new constants and updated path functions**

Add new test cases to `paths_test.go` for the new path functions. Update existing test expectations.

```go
// In paths_test.go — update TestGetFeaturePRDFile
func TestGetFeaturePRDFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd", "prd-spec.md")
	if got := GetFeaturePRDFile(feature); got != want {
		t.Errorf("GetFeaturePRDFile(%q) = %q, want %q", feature, got, want)
	}
}

// In paths_test.go — update TestGetFeatureDesignFile
func TestGetFeatureDesignFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "design", "tech-design.md")
	if got := GetFeatureDesignFile(feature); got != want {
		t.Errorf("GetFeatureDesignFile(%q) = %q, want %q", feature, got, want)
	}
}

// NEW tests
func TestGetFeaturePRDDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd")
	if got := GetFeaturePRDDir(feature); got != want {
		t.Errorf("GetFeaturePRDDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureDesignDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "design")
	if got := GetFeatureDesignDir(feature); got != want {
		t.Errorf("GetFeatureDesignDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUserStoriesFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd", "prd-user-stories.md")
	if got := GetFeatureUserStoriesFile(feature); got != want {
		t.Errorf("GetFeatureUserStoriesFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUIFunctionsFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd", "prd-ui-functions.md")
	if got := GetFeatureUIFunctionsFile(feature); got != want {
		t.Errorf("GetFeatureUIFunctionsFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureAPIHandbookFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "design", "api-handbook.md")
	if got := GetFeatureAPIHandbookFile(feature); got != want {
		t.Errorf("GetFeatureAPIHandbookFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUIDesignDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "ui")
	if got := GetFeatureUIDesignDir(feature); got != want {
		t.Errorf("GetFeatureUIDesignDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUIDesignFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "ui", "ui-design.md")
	if got := GetFeatureUIDesignFile(feature); got != want {
		t.Errorf("GetFeatureUIDesignFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureManifest(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "manifest.md")
	if got := GetFeatureManifest(feature); got != want {
		t.Errorf("GetFeatureManifest(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetProposalDir(t *testing.T) {
	slug := "my-idea"
	want := filepath.Join("docs", "proposals", "my-idea")
	if got := GetProposalDir(slug); got != want {
		t.Errorf("GetProposalDir(%q) = %q, want %q", slug, got, want)
	}
}

func TestGetProposalFile(t *testing.T) {
	slug := "my-idea"
	want := filepath.Join("docs", "proposals", "my-idea", "proposal.md")
	if got := GetProposalFile(slug); got != want {
		t.Errorf("GetProposalFile(%q) = %q, want %q", slug, got, want)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd task-cli && go test ./pkg/feature/ -run "TestGetFeaturePRDFile|TestGetFeatureDesignFile|TestGetFeaturePRDDir|TestGetFeatureDesignDir|TestGetFeatureUserStoriesFile|TestGetFeatureUIFunctionsFile|TestGetFeatureAPIHandbookFile|TestGetFeatureUIDesignDir|TestGetFeatureUIDesignFile|TestGetFeatureManifest|TestGetProposalDir|TestGetProposalFile" -v`
Expected: FAIL (undefined functions, wrong paths)

- [ ] **Step 3: Update constants.go**

```go
// File names and directory names within a feature directory
const (
	// IndexFileName is the name of the task index file
	IndexFileName = "index.json"

	// Subdirectory names within a feature directory
	PRDDirName    = "prd"
	DesignDirName = "design"
	UIDirName     = "ui"

	// PRD file names (source-prefixed to prevent naming collisions)
	PRDSpecFile        = "prd-spec.md"
	PRDUserStoriesFile = "prd-user-stories.md"
	PRDUIFunctionsFile = "prd-ui-functions.md"

	// Design file names
	TechDesignFile  = "tech-design.md"
	APIHandbookFile = "api-handbook.md"

	// UI design file names
	UIDesignFile = "ui-design.md"

	// Manifest file
	ManifestFileName = "manifest.md"

	// Tasks subdirectory names
	TasksDirName   = "tasks"
	RecordsDirName = "records"

	// Template file
	TemplateFileName = "template.md"

	// Proposals directory
	ProposalBaseDir  = "docs/proposals"
	ProposalFileName = "proposal.md"
)
```

Remove the old `PRDFileName` and `DesignFileName` constants entirely.

- [ ] **Step 4: Update paths.go**

Rename existing functions and add new ones:

```go
package feature

import (
	"path/filepath"
)

// GetFeatureDir returns the base directory for a feature.
func GetFeatureDir(feature string) string {
	return filepath.Join(FeaturesDir, feature)
}

// GetFeatureManifest returns the path to the feature's manifest.md.
func GetFeatureManifest(feature string) string {
	return filepath.Join(FeaturesDir, feature, ManifestFileName)
}

// GetFeaturePRDDir returns the path to the feature's prd/ subdirectory.
func GetFeaturePRDDir(feature string) string {
	return filepath.Join(FeaturesDir, feature, PRDDirName)
}

// GetFeaturePRDFile returns the path to the feature's prd/prd-spec.md.
func GetFeaturePRDFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, PRDDirName, PRDSpecFile)
}

// GetFeatureUserStoriesFile returns the path to prd/prd-user-stories.md.
func GetFeatureUserStoriesFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, PRDDirName, PRDUserStoriesFile)
}

// GetFeatureUIFunctionsFile returns the path to prd/prd-ui-functions.md.
func GetFeatureUIFunctionsFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, PRDDirName, PRDUIFunctionsFile)
}

// GetFeatureDesignDir returns the path to the feature's design/ subdirectory.
func GetFeatureDesignDir(feature string) string {
	return filepath.Join(FeaturesDir, feature, DesignDirName)
}

// GetFeatureDesignFile returns the path to design/tech-design.md.
func GetFeatureDesignFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, DesignDirName, TechDesignFile)
}

// GetFeatureAPIHandbookFile returns the path to design/api-handbook.md.
func GetFeatureAPIHandbookFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, DesignDirName, APIHandbookFile)
}

// GetFeatureUIDesignDir returns the path to the feature's ui/ subdirectory.
func GetFeatureUIDesignDir(feature string) string {
	return filepath.Join(FeaturesDir, feature, UIDirName)
}

// GetFeatureUIDesignFile returns the path to ui/ui-design.md.
func GetFeatureUIDesignFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, UIDirName, UIDesignFile)
}

// GetFeatureIndexFile returns the path to the feature's index.json.
func GetFeatureIndexFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName, IndexFileName)
}

// GetFeatureTasksDir returns the path to the feature's tasks directory.
func GetFeatureTasksDir(feature string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName)
}

// GetFeatureRecordsDir returns the path to the feature's records directory (under tasks).
func GetFeatureRecordsDir(feature string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName, RecordsDirName)
}

// GetTaskFile returns the path to a specific task file.
func GetTaskFile(feature, filename string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName, filename)
}

// GetRecordFile returns the path to a specific record file (under tasks/records).
func GetRecordFile(feature, filename string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName, RecordsDirName, filename)
}

// GetTaskStatePath returns the absolute path to state.json for a feature.
func GetTaskStatePath(projectRoot, feature string) string {
	return filepath.Join(projectRoot, FeaturesDir, feature, TasksDirName, ProcessDirName, StateFileName)
}

// GetProcessRecordPath returns the absolute path to the in-progress record.json.
func GetProcessRecordPath(projectRoot, feature string) string {
	return filepath.Join(projectRoot, FeaturesDir, feature, TasksDirName, ProcessDirName, RecordFileName)
}

// GetProcessDir returns the absolute path to the process directory for a feature.
func GetProcessDir(projectRoot, feature string) string {
	return filepath.Join(projectRoot, FeaturesDir, feature, TasksDirName, ProcessDirName)
}

// GetProposalDir returns the path to a proposal directory.
func GetProposalDir(slug string) string {
	return filepath.Join(ProposalBaseDir, slug)
}

// GetProposalFile returns the path to a proposal.md file.
func GetProposalFile(slug string) string {
	return filepath.Join(ProposalBaseDir, slug, ProposalFileName)
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd task-cli && go test ./pkg/feature/ -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add task-cli/pkg/feature/constants.go task-cli/pkg/feature/paths.go task-cli/pkg/feature/paths_test.go
git commit -m "refactor(task-cli): rename constants and add nested subdirectory paths

Rename PRDFileName→PRDSpecFile, DesignFileName→TechDesignFile.
Add path functions for prd/, design/, ui/, manifest.md, proposals/.
Update tests to expect nested paths."
```

---

### Task 2: Update EnsureFeatureDir

**Files:**
- Modify: `task-cli/pkg/feature/feature.go`
- Test: `task-cli/pkg/feature/feature_test.go`

- [ ] **Step 1: Write the failing test**

Update `TestEnsureFeatureDir` in `feature_test.go`:

```go
func TestEnsureFeatureDir(t *testing.T) {
	t.Run("creates all directories", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		if err := EnsureFeatureDir(dir, featureSlug); err != nil {
			t.Fatalf("EnsureFeatureDir() error = %v", err)
		}

		expectedDirs := []string{
			filepath.Join(dir, GetFeatureDir(featureSlug)),
			filepath.Join(dir, GetFeaturePRDDir(featureSlug)),
			filepath.Join(dir, GetFeatureDesignDir(featureSlug)),
			filepath.Join(dir, GetFeatureUIDesignDir(featureSlug)),
			filepath.Join(dir, GetFeatureTasksDir(featureSlug)),
			filepath.Join(dir, GetFeatureRecordsDir(featureSlug)),
			filepath.Join(dir, FeaturesDir, featureSlug, TasksDirName, ProcessDirName),
		}

		for _, expectedDir := range expectedDirs {
			if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
				t.Errorf("directory %s was not created", expectedDir)
			}
		}
	})
}
```

Also update `TestSetFeature` similarly to check for `prd/`, `design/`, `ui/` dirs.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd task-cli && go test ./pkg/feature/ -run TestEnsureFeatureDir -v`
Expected: FAIL (missing `prd/`, `design/`, `ui/` directories)

- [ ] **Step 3: Update EnsureFeatureDir in feature.go**

```go
// EnsureFeatureDir ensures the feature directory structure exists.
func EnsureFeatureDir(projectRoot, featureSlug string) error {
	dirs := []string{
		GetFeatureDir(featureSlug),
		GetFeaturePRDDir(featureSlug),
		GetFeatureDesignDir(featureSlug),
		GetFeatureUIDesignDir(featureSlug),
		GetFeatureTasksDir(featureSlug),
		GetFeatureRecordsDir(featureSlug),
		filepath.Join(FeaturesDir, featureSlug, TasksDirName, ProcessDirName),
	}
	for _, dir := range dirs {
		fullPath := filepath.Join(projectRoot, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd task-cli && go test ./pkg/feature/ -run "TestEnsureFeatureDir|TestSetFeature" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add task-cli/pkg/feature/feature.go task-cli/pkg/feature/feature_test.go
git commit -m "refactor(task-cli): create prd/, design/, ui/ subdirs in EnsureFeatureDir"
```

---

### Task 3: Update Downstream Go Tests

**Files:**
- Modify: `task-cli/pkg/task/types_test.go`
- Modify: `task-cli/internal/cmd/check_test.go`
- Modify: `task-cli/internal/cmd/claim_test.go`
- Modify: `task-cli/internal/cmd/feature_test.go`
- Modify: `task-cli/internal/cmd/runners_test.go`
- Modify: `task-cli/internal/cmd/validate_test.go`

- [ ] **Step 1: Replace all `"prd.md"` → `"prd/prd-spec.md"` and `"design.md"` → `"design/tech-design.md"` in test files**

Use replace_all to update each file. In each file, replace:
- `PRD:     "prd.md"` → `PRD:     "prd/prd-spec.md"`
- `PRD:          "prd.md"` → `PRD:          "prd/prd-spec.md"`
- `Design:  "design.md"` → `Design:  "design/tech-design.md"`
- `Design:       "design.md"` → `Design:       "design/tech-design.md"`

Files to update (all under `task-cli/`):
1. `pkg/task/types_test.go` (lines 93-94)
2. `internal/cmd/check_test.go` (lines 30-31)
3. `internal/cmd/claim_test.go` (lines 492-493, 560-561)
4. `internal/cmd/feature_test.go` (lines 190-191, 258-259, 321-322)
5. `internal/cmd/runners_test.go` (lines 326-327)
6. `internal/cmd/validate_test.go` (lines 291-292)

- [ ] **Step 2: Run full test suite**

Run: `cd task-cli && go test -race -cover ./...`
Expected: PASS (all tests)

- [ ] **Step 3: Commit**

```bash
git add task-cli/pkg/task/types_test.go task-cli/internal/cmd/check_test.go task-cli/internal/cmd/claim_test.go task-cli/internal/cmd/feature_test.go task-cli/internal/cmd/runners_test.go task-cli/internal/cmd/validate_test.go
git commit -m "test(task-cli): update downstream test expectations for nested paths"
```

---

### Task 4: Update index.json Template and TaskIndex Types

**Files:**
- Modify: `plugins/zcode/skills/breakdown-tasks/templates/index.json`
- Modify: `task-cli/pkg/task/types.go` (comments only)

- [ ] **Step 1: Update index.json template**

```json
{
  "feature": "<feature-slug>",
  "prd": "prd/prd-spec.md",
  "design": "design/tech-design.md",
  "created": "YYYY-MM-DD",
  "status": "planning",
  "tasks": {
    "1.1-interface": {
      "id": "1.1",
      "title": "Interface Definition",
      "priority": "P0",
      "estimatedTime": "2-3h",
      "dependencies": [],
      "status": "pending",
      "file": "tasks/1.1-interface.md",
      "record": "records/1.1-interface.md"
    },
    "3.1-implementation": {
      "id": "3.1",
      "title": "Core Implementation",
      "priority": "P0",
      "estimatedTime": "4h",
      "dependencies": ["1.1", "2.1"],
      "status": "pending",
      "file": "tasks/3.1-implementation.md",
      "record": "records/3.1-implementation.md"
    }
  }
}
```

- [ ] **Step 2: Update types.go comments**

Update `TaskIndex` field comments in `task-cli/pkg/task/types.go`:

```go
// TaskIndex represents the index.json structure for a feature.
type TaskIndex struct {
	Feature      string          `json:"feature,omitempty"`
	PRD          string          `json:"prd,omitempty"`     // e.g. "prd/prd-spec.md"
	Design       string          `json:"design,omitempty"`  // e.g. "design/tech-design.md"
	Created      string          `json:"created,omitempty"`
	Status       string          `json:"status,omitempty"`
	Tasks        map[string]Task `json:"tasks"`
	StatusEnum   []string        `json:"statusEnum,omitempty"`
	PriorityEnum []string        `json:"priorityEnum,omitempty"`
}
```

- [ ] **Step 3: Run full Go test suite**

Run: `cd task-cli && go test -race -cover ./...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add plugins/zcode/skills/breakdown-tasks/templates/index.json task-cli/pkg/task/types.go
git commit -m "refactor: update index.json template and types.go for nested paths"
```

---

### Task 5: Create Manifest Template

**Files:**
- Create: `plugins/zcode/skills/write-prd/templates/manifest.md`

- [ ] **Step 1: Create manifest template**

```markdown
# Feature: {{FEATURE_SLUG}}

## Status

<!-- Auto-updated by skills. Do not edit manually. -->
prd -> design -> tasks -> in-progress -> done

Current: prd

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | {{PRD_SUMMARY}} |
| User Stories | prd/prd-user-stories.md | {{USER_STORIES_SUMMARY}} |
{{UI_FUNCTIONS_ROW}}
## Traceability

<!-- Auto-generated by skills. Entries link PRD sections to design sections and tasks. -->
| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
```

Where `{{UI_FUNCTIONS_ROW}}` is conditionally:
```
| UI Functions | prd/prd-ui-functions.md | {{UI_FUNCTIONS_SUMMARY}} |
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/write-prd/templates/manifest.md
git commit -m "feat(templates): add manifest.md template for feature traceability"
```

---

### Task 6: Rename and Create PRD Templates

**Files:**
- Rename: `plugins/zcode/skills/write-prd/templates/prd.md` → `prd-spec.md`
- Rename: `plugins/zcode/skills/write-prd/templates/user-stories.md` → `prd-user-stories.md`
- Create: `plugins/zcode/skills/write-prd/templates/prd-ui-functions.md`

- [ ] **Step 1: Rename prd.md → prd-spec.md (content update needed)**

The template content stays largely the same. Update the first line heading and metadata comment:

```markdown
# {{FEATURE_NAME}} — PRD Spec

> PRD Spec: defines WHAT the feature is and why it exists.

## 需求背景
<!-- ... rest of existing template content unchanged ... -->
```

Run: `cd plugins/zcode/skills/write-prd/templates && git mv prd.md prd-spec.md`

- [ ] **Step 2: Rename user-stories.md → prd-user-stories.md**

Content unchanged (just filename).

Run: `cd plugins/zcode/skills/write-prd/templates && git mv user-stories.md prd-user-stories.md`

- [ ] **Step 3: Create prd-ui-functions.md template**

```markdown
# {{FEATURE_NAME}} — UI Functions

> Requirements layer: defines WHAT the UI must do. Not HOW it looks (that's ui-design.md).

## UI Scope

<!-- List all UI surfaces this feature requires -->

## UI Function 1: {{Function Name}}

### Description
<!-- What this UI element does -->

### User Interaction Flow
<!-- Step-by-step interaction: user clicks X → system responds with Y -->

### Data Requirements
<!-- What data this UI element needs to display or collect -->

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |

### States
<!-- States this UI element can be in (loading, empty, error, populated, etc.) -->

| State | Display | Trigger |
|-------|---------|---------|
| <!-- --> | <!-- --> | <!-- --> |

### Validation Rules
<!-- Input validation, conditional display, etc. -->

---

## UI Function 2: {{Function Name}}

<!-- Repeat pattern above for each UI surface -->
```

- [ ] **Step 4: Commit**

```bash
git add plugins/zcode/skills/write-prd/templates/
git commit -m "refactor(templates): rename prd templates to source-prefixed names, add prd-ui-functions"
```

---

### Task 7: Rename and Create Design Templates

**Files:**
- Rename: `plugins/zcode/skills/design-tech/templates/design.md` → `tech-design.md`
- Create: `plugins/zcode/skills/design-tech/templates/api-handbook.md`
- Create: `plugins/zcode/skills/design-tech/templates/manifest-update-design.md`

- [ ] **Step 1: Rename design.md → tech-design.md**

Update template heading and PRD reference:

```markdown
# Technical Design: <Feature Name>

## Metadata
- Created: YYYY-MM-DD
- PRD: prd/prd-spec.md
- Status: Draft | Review | Approved
<!-- ... rest unchanged ... -->
```

Run: `cd plugins/zcode/skills/design-tech/templates && git mv design.md tech-design.md`

- [ ] **Step 2: Create api-handbook.md template**

```markdown
# API Handbook: <Feature Name>

## Metadata
- Created: YYYY-MM-DD
- Related: design/tech-design.md

## API Overview

<!-- High-level API design summary -->

## Endpoints

### {{Endpoint Name}}

**Method**: `{{GET|POST|PUT|DELETE}}`
**Path**: `{{/api/path}}`
**Auth**: <!-- Required role or "none" -->

#### Request

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |

#### Response (200)

| Field | Type | Description |
|-------|------|-------------|
| <!-- --> | <!-- --> | <!-- --> |

#### Error Responses

| Status | Code | Description |
|--------|------|-------------|
| 400 | <!-- --> | <!-- --> |
| 404 | <!-- --> | <!-- --> |

---

<!-- Repeat for each endpoint -->

## Data Contracts

<!-- Shared types used across endpoints -->

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| <!-- --> | <!-- --> | <!-- --> |
```

- [ ] **Step 3: Create manifest-update-design.md template**

```markdown
<!-- Snippet to append/update in manifest.md after /design-tech completes -->

## Documents (updated)

Add rows:
| Tech Design | design/tech-design.md | {{TECH_DESIGN_SUMMARY}} |
| API Handbook | design/api-handbook.md | {{API_HANDBOOK_SUMMARY}} |

## Traceability (updated)

Add entries linking PRD sections to design sections:
| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
| "{{PRD_SECTION}}" (prd-spec §N) | "{{DESIGN_SECTION}}" (tech-design §N) | <!-- task IDs added by /breakdown-tasks --> |

## Status

Advance to "design" if /ui-design already completed or if UI is not applicable.
```

- [ ] **Step 4: Commit**

```bash
git add plugins/zcode/skills/design-tech/templates/
git commit -m "refactor(templates): rename design.md→tech-design.md, add api-handbook and manifest-update templates"
```

---

### Task 8: Create Brainstorm Skill

**Files:**
- Create: `plugins/zcode/skills/brainstorm/SKILL.md`
- Create: `plugins/zcode/skills/brainstorm/templates/proposal.md`

- [ ] **Step 1: Create proposal.md template**

```markdown
# Proposal: {{PROPOSAL_TITLE}}

## Metadata
- Created: YYYY-MM-DD
- Author: <!-- who proposed this -->
- Status: Draft | Reviewed | Approved | Rejected

## Problem

<!-- What problem are we solving? Why now? -->

## Proposed Solution

<!-- High-level approach -->

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |

## Scope

### In Scope
- <!-- -->

### Out of Scope
- <!-- -->

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |

## Success Criteria

- [ ] <!-- measurable outcome -->

## Next Steps

- Proceed to `/write-prd` to formalize requirements
```

- [ ] **Step 2: Create brainstorm SKILL.md**

```markdown
---
name: brainstorm
description: Use when a user has a vague idea or feature request and needs to explore it before formalizing into a PRD. Outputs a structured proposal document.
---

# Brainstorm

## Overview

从模糊想法到结构化提案，通过协作对话探索问题空间。

**核心原则**：在投入 PRD 之前，先确认问题值得解决、方案值得投入。

<HARD-GATE>
Do NOT write any code or take implementation action. This skill produces a proposal document only.
</HARD-GATE>

## When to Use

**Trigger conditions:**
- User describes an idea without clear specifications
- User says "I'm thinking about..." or "What if we..."
- Starting exploration before committing to a feature

**Skip when:**
- Requirements are already clear (use `/write-prd` directly)
- Bug fix or small tweak

## Process Flow

```
Understand idea → Explore context → Challenge assumptions → Define scope → Write proposal → Commit
```

## Step 1: Understand the Idea

Listen actively. Ask clarifying questions one at a time via `AskUserQuestion`:
- What problem does this solve?
- Who is affected?
- What does success look like?

## Step 2: Explore Context

| Source | What to Look For |
|--------|-----------------|
| Existing features | Is this already solved elsewhere? |
| Recent commits | Related recent changes |
| Project docs | Architecture constraints, existing decisions |

## Step 3: Challenge Assumptions

Play devil's advocate:
- Is this the right problem to solve?
- Are there simpler alternatives?
- What if we did nothing?

## Step 4: Define Scope

Propose in-scope and out-of-scope boundaries. Get user agreement.

## Step 5: Write Proposal

Save to `docs/proposals/<slug>/proposal.md` using `templates/proposal.md`.

## Step 6: Commit

```bash
git add docs/proposals/<slug>/
git commit -m "docs: add proposal for <feature-slug>"
```

## Integration

Works well with:
- `/write-prd` — Takes proposal as optional input to formalize into PRD
```

- [ ] **Step 3: Commit**

```bash
git add plugins/zcode/skills/brainstorm/
git commit -m "feat(skills): add brainstorm skill with proposal template"
```

---

### Task 9: Create UI Design Skill

**Files:**
- Create: `plugins/zcode/skills/ui-design/SKILL.md`
- Create: `plugins/zcode/skills/ui-design/templates/ui-design.md`
- Create: `plugins/zcode/skills/ui-design/templates/manifest-update-ui.md`

- [ ] **Step 1: Create ui-design.md template**

```markdown
# UI Design: {{FEATURE_NAME}}

## Metadata
- Created: YYYY-MM-DD
- Source: prd/prd-ui-functions.md
- Status: Draft | Review | Approved

## Design System

<!-- Reference to existing design system or component library -->

## Component: {{Component Name}}

### Layout Structure
<!-- Component hierarchy, grid/flex layout description -->

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Default | <!-- --> | <!-- --> |
| Loading | <!-- --> | <!-- --> |
| Empty | <!-- --> | <!-- --> |
| Error | <!-- --> | <!-- --> |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| <!-- --> | <!-- --> | <!-- --> |

### Data Binding

| UI Element | Data Field | Source |
|------------|-----------|--------|
| <!-- --> | <!-- --> | <!-- --> |

---

<!-- Repeat for each component -->
```

- [ ] **Step 2: Create manifest-update-ui.md template**

```markdown
<!-- Snippet to update manifest.md after /ui-design completes -->

## Documents (updated)

Add row:
| UI Design | ui/ui-design.md | {{UI_DESIGN_SUMMARY}} |

## Traceability (updated)

Add entries linking PRD UI functions to UI design sections:
| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
| "UI Functions > {{Function Name}}" | "UI Design > {{Component Name}}" | <!-- task IDs added by /breakdown-tasks --> |

## Status

Advance to "design" if /design-tech already completed.
```

- [ ] **Step 3: Create ui-design SKILL.md**

```markdown
---
name: ui-design
description: Use after PRD ui-functions are defined to create UI design specifications. Parallel to /design-tech, reads prd/prd-ui-functions.md.
---

# UI Design

## Overview

从 PRD 的 UI 功能需求产出 UI 设计规格文档。

**核心原则**：定义 HOW 界面呈现和交互，与 PRD 的 WHAT（需求层）分离。

<HARD-GATE>
Do NOT write any implementation code. This skill produces a design specification document only.
</HARD-GATE>

## Position in Workflow

```
/write-prd → /design-tech ─→ /breakdown-tasks
     ↓            ↓
     ↓       /ui-design ──→ /breakdown-tasks
     ↓
prd/prd-ui-functions.md → ui/ui-design.md
```

Parallel to `/design-tech`. Both must complete before `/breakdown-tasks`.

## When to Use

**Trigger conditions:**
- PRD with `prd/prd-ui-functions.md` exists
- User asks to design the UI

**Skip when:**
- No UI functions defined (backend/API/CLI features)
- UI design already exists

## Process Flow

```
1. Read manifest → 2. Read UI functions → 3. Explore patterns → 4. Draft design → 5. Review → 6. Update manifest
```

## Step 1: Read Manifest

Read `manifest.md` to locate `prd/prd-ui-functions.md`.

## Step 2: Read UI Functions

Read `prd/prd-ui-functions.md` to understand UI requirements.

## Step 3: Explore Existing Patterns

- Check for existing design system or component library
- Review existing UI components in the project

## Step 4: Draft UI Design

For each UI function, define:
- Layout structure (component hierarchy)
- States (loading, empty, error, populated)
- Interactions (triggers, actions, feedback)
- Data binding (UI element → data field)

## Step 5: Write UI Design

Save to `ui/ui-design.md` using `templates/ui-design.md`.

## Step 6: Update Manifest

Update `manifest.md`:
- Add UI Design row to Documents table
- Add traceability links from UI Functions to UI Design sections
- Advance status to `design` if `/design-tech` already completed

## Integration

Works well with:
- `/write-prd` — Produces `prd/prd-ui-functions.md` input
- `/design-tech` — Parallel skill; both must complete before breakdown
- `/eval-design` — Evaluates UI design alongside tech design
```

- [ ] **Step 4: Commit**

```bash
git add plugins/zcode/skills/ui-design/
git commit -m "feat(skills): add ui-design skill with templates"
```

---

### Task 10: Create Breakdown-Tasks Manifest Update Template

**Files:**
- Create: `plugins/zcode/skills/breakdown-tasks/templates/manifest-update-tasks.md`

- [ ] **Step 1: Create manifest-update-tasks.md template**

```markdown
<!-- Snippet to update manifest.md after /breakdown-tasks completes -->

## Traceability (updated)

Fill Tasks column with task IDs linked to design sections:
| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
| "{{PRD_SECTION}}" (prd-spec §N) | "{{DESIGN_SECTION}}" (tech-design §N) | {{TASK_IDS}} |

## Status

Advance to "tasks".
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/breakdown-tasks/templates/manifest-update-tasks.md
git commit -m "feat(templates): add manifest-update-tasks template for breakdown skill"
```

---

### Task 11: Update write-prd SKILL.md

**Files:**
- Modify: `plugins/zcode/skills/write-prd/SKILL.md`

- [ ] **Step 1: Update SKILL.md**

Key changes to `plugins/zcode/skills/write-prd/SKILL.md`:

1. **Checklist** — Step 1 add: check `docs/proposals/<slug>/proposal.md` as optional input
2. **Step 6** — output path: `prd.md` → `prd/prd-spec.md`
3. **Step 7** — output path: `user-stories.md` → `prd/prd-user-stories.md`
4. **New Step 8** — output `prd/prd-ui-functions.md` (optional, only for features with UI)
5. **New Step 9** — create `manifest.md` at feature root
6. **Output Documents table** — 3 files + manifest
7. **Directory structure** — updated tree

Replace the relevant sections in SKILL.md:

```markdown
## Checklist

1. **Explore project context** — check files, docs, recent commits
2. **Check for existing proposal** — read `docs/proposals/<slug>/proposal.md` if it exists
3. **Assess scope** — determine if request needs decomposition
4. **Ask clarifying questions** — one at a time via AskUserQuestion tool
5. **Propose 2-3 approaches** — with trade-offs and your recommendation
6. **Present PRD sections** — get approval after each section
7. **Write PRD Spec** — save to `docs/features/<feature-slug>/prd/prd-spec.md`
8. **Write User Stories** — save to `docs/features/<feature-slug>/prd/prd-user-stories.md`
9. **Write UI Functions** (if applicable) — save to `docs/features/<feature-slug>/prd/prd-ui-functions.md`
10. **Create Manifest** — save to `docs/features/<feature-slug>/manifest.md`
11. **Commit** — commit all documents

## Output Documents

PRD 完成后输出以下文件：

| 文件 | 模板 | 说明 |
|------|------|------|
| `prd/prd-spec.md` | `templates/prd-spec.md` | 产品需求文档，包含背景、目标、Scope、流程、功能描述等 |
| `prd/prd-user-stories.md` | `templates/prd-user-stories.md` | 用户故事，从 PRD 背景中识别的用户角色推导而出 |
| `prd/prd-ui-functions.md` | `templates/prd-ui-functions.md` | UI 功能要点（需求层，仅适用于有 UI 表面的功能） |
| `manifest.md` | `templates/manifest.md` | Feature 索引和可追溯性映射 |

## Step 6: Write PRD Spec

使用 `templates/prd-spec.md` 模板填写。

**目录结构：**

```
docs/features/<feature-slug>/
├── manifest.md                # Feature index & traceability
├── prd/
│   ├── prd-spec.md            # PRD Spec
│   ├── prd-user-stories.md    # 用户故事
│   └── prd-ui-functions.md    # UI 功能要点（可选）
├── design/                    # (created by /design-tech)
├── ui/                        # (created by /ui-design)
└── tasks/                     # (created by /breakdown-tasks)
    └── records/
```

## Step 9: Create Manifest

Create `manifest.md` at the feature root using `templates/manifest.md`:
- Fill in PRD entries and summaries
- Set status to `prd`
- Include UI Functions row only if `prd/prd-ui-functions.md` was created
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/write-prd/SKILL.md
git commit -m "refactor(write-prd): update paths to nested prd/ subdirectory, add manifest creation"
```

---

### Task 12: Update design-tech SKILL.md

**Files:**
- Modify: `plugins/zcode/skills/design-tech/SKILL.md`

- [ ] **Step 1: Update SKILL.md**

Key changes to `plugins/zcode/skills/design-tech/SKILL.md`:

1. **Position in Workflow** — add manifest reference
2. **When to Use** — update trigger path to `prd/prd-spec.md`
3. **Step 1** — read via manifest → `prd/prd-spec.md`
4. **Step 7** — output split: `design/tech-design.md` + `design/api-handbook.md`
5. **New Step 8** — update `manifest.md` with design entries and traceability
6. **Integration** — updated references

Replace the relevant sections:

```markdown
## Position in Workflow

```
/write-prd → /design-tech → /eval-design → /breakdown-tasks
     ↓              ↓              ↓               ↓
  prd/*.{3}    design/*.{2}  eval report     tasks/*.md
  manifest.md  manifest.md                  manifest.md
```

## When to Use

**Trigger conditions:**
- Manifest exists at `docs/features/<slug>/manifest.md` with status `prd`
- PRD Spec exists at `prd/prd-spec.md`

**Skip when:**
- No manifest or PRD exists (use `/write-prd` first)
- Design already exists for the feature

## Step 1: Read Manifest → PRD

1. Read `manifest.md` to locate documents
2. Read `prd/prd-spec.md`:
   - Understand requirements
   - Note non-functional requirements
   - Identify acceptance criteria

## Step 7: Write Design Documents

Save to:
- `docs/features/<slug>/design/tech-design.md` — using `templates/tech-design.md`
- `docs/features/<slug>/design/api-handbook.md` — using `templates/api-handbook.md` (if feature has API surface)

## Step 8: Update Manifest

Update `manifest.md`:
- Add Tech Design and API Handbook rows to Documents table
- Add traceability links from PRD sections to design sections
- Advance status to `design` if `/ui-design` already completed or if UI is not applicable
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/design-tech/SKILL.md
git commit -m "refactor(design-tech): update paths to nested design/, read via manifest"
```

---

### Task 13: Update eval-prd SKILL.md

**Files:**
- Modify: `plugins/zcode/skills/eval-prd/SKILL.md`

- [ ] **Step 1: Update SKILL.md**

Key changes:
1. **Step 1** — locate via `manifest.md` first, then `prd/prd-spec.md`
2. **Agent prompt** — update all path references

Replace Step 1 and agent prompt paths:

```markdown
## Step 1: Locate Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate PRD documents
3. Fall back to `docs/features/<current-feature>/prd/prd-spec.md` + `prd/prd-user-stories.md`
4. Ask user for path if not found
```

In the agent prompt template, update paths:
- `{{PRD_PATH}}` default: `prd/prd-spec.md`
- `{{USER_STORIES_PATH}}` default: `prd/prd-user-stories.md`
- Add: `{{UI_FUNCTIONS_PATH}}` optional: `prd/prd-ui-functions.md`

Add a UI Functions dimension check (optional):
```markdown
## Dimension 6: UI Functions (optional)

Only checked if `prd/prd-ui-functions.md` exists.

Checks: each UI function has description, interaction flow, data requirements, states, validation.

- A: All functions fully specified with all sub-sections
- B: Most specified, 1-2 missing sub-sections
- C: Functions listed but incomplete
- N/A: File doesn't exist (not an F)
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-prd/SKILL.md
git commit -m "refactor(eval-prd): locate documents via manifest, add ui-functions check"
```

---

### Task 14: Update eval-design SKILL.md

**Files:**
- Modify: `plugins/zcode/skills/eval-design/SKILL.md`

- [ ] **Step 1: Update SKILL.md**

Key changes:
1. **Step 1** — locate via `manifest.md` first
2. **Agent prompt** — update all path references
3. Add checks for `design/api-handbook.md` and `ui/ui-design.md`

Replace Step 1:
```markdown
## Step 1: Locate Design Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate design documents
3. Fall back to `design/tech-design.md`, `design/api-handbook.md`, `ui/ui-design.md`
4. Ask user for path if not found
```

In agent prompt, update paths:
- `{{DESIGN_PATH}}` default: `design/tech-design.md`
- Add `{{API_HANDBOOK_PATH}}` default: `design/api-handbook.md`
- Add `{{UI_DESIGN_PATH}}` default: `ui/ui-design.md`
- `{{PRD_PATH}}` default: `prd/prd-spec.md`

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-design/SKILL.md
git commit -m "refactor(eval-design): locate documents via manifest, add api-handbook and ui-design checks"
```

---

### Task 15: Update breakdown-tasks SKILL.md

**Files:**
- Modify: `plugins/zcode/skills/breakdown-tasks/SKILL.md`

- [ ] **Step 1: Update SKILL.md**

Key changes:
1. **Position in Workflow** — add manifest
2. **Directory Structure** — updated tree
3. **Step 1** — read `manifest.md` → all docs
4. **New Step 6** — update manifest with task traceability

Replace relevant sections:

```markdown
## Position in Workflow

```
/write-prd → /eval-prd → /design-tech → /eval-design → /breakdown-tasks
     ↓             ↓            ↓              ↓               ↓
  prd/*.{3}   eval report  design/*.{2}  eval report    tasks/*.md
  manifest.md              manifest.md                  manifest.md
```

## Directory Structure

```
docs/features/<feature-slug>/
├── manifest.md                    # Feature index & traceability
├── prd/
│   ├── prd-spec.md
│   ├── prd-user-stories.md
│   └── prd-ui-functions.md
├── design/
│   ├── tech-design.md             # Technical design (input)
│   └── api-handbook.md
├── ui/
│   └── ui-design.md               # (if applicable)
├── tasks/
│   ├── index.json                 # Task index
│   ├── 1.1-<title>.md            # Task detail files
│   ├── process/
│   └── records/
```

## Step 1: Read Manifest → All Documents

1. Read `manifest.md` to locate all documents
2. Read `prd/prd-spec.md` — understand WHAT
3. Read `design/tech-design.md` — understand HOW
4. Read `design/api-handbook.md` — understand interfaces (if exists)
5. Read `ui/ui-design.md` — understand UI components (if exists)
6. Read `prd/prd-user-stories.md` — understand user scenarios (if exists)

## Step 6: Validate

```bash
task validate -file docs/features/<slug>/tasks/index.json
```

## Step 7: Update Manifest

Update `manifest.md`:
- Fill Tasks column in Traceability table with task IDs linked to design sections
- Advance status to `tasks`
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/breakdown-tasks/SKILL.md
git commit -m "refactor(breakdown-tasks): read all docs via manifest, add traceability update"
```

---

### Task 16: Update guide.md

**Files:**
- Modify: `plugins/zcode/hooks/guide.md`

- [ ] **Step 1: Update Document Index section**

Replace the Document Index in guide.md with the new structure:

```markdown
## Document Index

```
project-root/
└── docs/
    ├── proposals/<slug>/           # /brainstorm 产出
    │   └── proposal.md
    ├── features/<slug>/            # Feature 工作区
    │   ├── manifest.md             # Feature 索引 & 可追溯性映射
    │   ├── prd/
    │   │   ├── prd-spec.md         # PRD Spec (需求文档)
    │   │   ├── prd-user-stories.md # 用户故事
    │   │   └── prd-ui-functions.md # UI 功能要点（可选）
    │   ├── design/
    │   │   ├── tech-design.md      # 技术设计
    │   │   └── api-handbook.md     # API 文档
    │   ├── ui/
    │   │   └── ui-design.md        # UI 设计规格（可选）
    │   └── tasks/
    │       ├── index.json          # 任务定义（核心）
    │       ├── process/            # 运行时状态（不提交）
    │       │   ├── state.json
    │       │   └── record.json
    │       ├── 1.1-<title>.md     # 任务详情
    │       └── records/            # 执行记录
    ├── README.md
    ├── ARCHITECTURE.md
    ├── DECISIONS.md
    └── lessons/
```

### Manifest

`manifest.md` 是 Feature 的单一入口，AI agent 读取此文件即可了解完整上下文：
- **Documents** 表：列出所有文档路径和自动生成的摘要
- **Traceability** 表：PRD → Design → Tasks 的追溯映射
- **Status**：prd → design → tasks → in-progress → done
```

Keep the Task-CLI section, `task record` workflow, record.json format, and forbidden operations sections unchanged.

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/hooks/guide.md
git commit -m "docs(guide): update Document Index to nested directory structure with manifest"
```

---

### Task 17: Update plugin.json

**Files:**
- Modify: `plugins/zcode/.claude-plugin/plugin.json`

- [ ] **Step 1: Update version and keywords**

```json
{
	"name": "zcode",
	"version": "2.0.0",
	"description": "Task management and workflow helper tools for Claude Code",
	"keywords": ["task", "workflow", "productivity", "prd", "git", "brainstorm", "ui-design", "manifest"]
}
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/.claude-plugin/plugin.json
git commit -m "chore: bump plugin version to 2.0.0, add new keywords"
```

---

### Task 18: Update task-cli OVERVIEW.md

**Files:**
- Modify: `task-cli/docs/OVERVIEW.md`

- [ ] **Step 1: Update directory structure diagram**

Replace the directory structure section (lines 76-90) with:

```markdown
## 目录结构约定

```
project-root/
├── docs/
│   ├── proposals/<slug>/           # /brainstorm 产出
│   │   └── proposal.md
│   └── features/<slug>/            # Feature 工作区
│       ├── manifest.md             # Feature 索引 & 可追溯性映射
│       ├── prd/
│       │   ├── prd-spec.md         # PRD Spec
│       │   ├── prd-user-stories.md # 用户故事
│       │   └── prd-ui-functions.md # UI 功能要点（可选）
│       ├── design/
│       │   ├── tech-design.md      # 技术设计
│       │   └── api-handbook.md     # API 文档
│       ├── ui/
│       │   └── ui-design.md        # UI 设计规格（可选）
│       └── tasks/
│           ├── index.json          # 任务定义
│           ├── process/            # 运行时状态
│           │   ├── state.json
│           │   └── record.json
│           ├── 1.1-<title>.md     # 任务详情
│           └── records/            # 执行记录
```
```

- [ ] **Step 2: Commit**

```bash
git add task-cli/docs/OVERVIEW.md
git commit -m "docs(task-cli): update directory structure in OVERVIEW.md"
```

---

### Task 19: Update Parent Redesign Plan

**Files:**
- Modify: `docs/zcode-redesign-plan.md`

- [ ] **Step 1: Update the parent plan to align with the approved spec**

Apply the changes specified in the spec's "Impact on Redesign Plan" table:

| Location | Current | Updated |
|----------|---------|---------|
| Target structure (line ~16) | `tech/` with `overview.md` | `design/` with `tech-design.md` |
| `EnsureFeatureDir` (line ~63) | `prd/`, `design/`, `design/ui/` | `prd/`, `design/`, `ui/`, `tasks/` |
| `ProposalBaseDir` (line ~53) | `"docs/proposal"` | `"docs/proposals"` |
| Phase 2 brainstorm output (line ~119) | `docs/proposal/<slug>/proposal.md` | `docs/proposals/<slug>/proposal.md` |
| Phase 3 ui-design output (line ~133) | `design/ui/` | `ui/` |
| Phase 4.2 design-tech prose (line ~166) | "`design/ui/` 由 `/ui-design` skill 填充" | "`ui/` 由 `/ui-design` skill 填充" |
| Phase 6 e2e verification (line ~223) | "design/ui/ content" | "ui/ content" |

Also add a note at the top: "Directory structure superseded by `docs/superpowers/specs/2026-04-09-directory-structure-redesign.md`."

- [ ] **Step 2: Commit**

```bash
git add docs/zcode-redesign-plan.md
git commit -m "docs: align parent redesign plan with approved directory structure spec"
```

---

### Task 20: Final Verification

**Files:** None (verification only)

- [ ] **Step 1: Run full Go test suite**

Run: `cd task-cli && go build ./... && go vet ./... && go test -race -cover ./...`
Expected: PASS (all tests, 0 failures)

- [ ] **Step 2: Run lint**

Run: `cd task-cli && golangci-lint run ./...`
Expected: PASS (0 issues)

- [ ] **Step 3: Verify template file renames are consistent**

Run: `ls plugins/zcode/skills/write-prd/templates/ && ls plugins/zcode/skills/design-tech/templates/`
Expected:
- `prd-spec.md`, `prd-user-stories.md`, `prd-ui-functions.md`, `manifest.md`
- `tech-design.md`, `api-handbook.md`, `manifest-update-design.md`

- [ ] **Step 4: Verify new skill directories exist**

Run: `ls plugins/zcode/skills/brainstorm/ && ls plugins/zcode/skills/ui-design/`
Expected: Each contains `SKILL.md` and `templates/`

- [ ] **Step 5: Verify no stale references remain**

Run: `grep -r "prd\.md\|design\.md\|design/ui/" plugins/zcode/skills/ --include="*.md" | grep -v "templates/"`
Expected: No matches (all references updated)
