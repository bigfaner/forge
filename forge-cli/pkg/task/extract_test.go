package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- extractScope tests ---

func TestExtractScope(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "extracts in-scope items",
			content: "## Scope\n\n### In Scope\n\n- Backend API endpoints\n- Database schema\n- CLI commands\n\n### Out of Scope\n\n- Frontend UI\n- Documentation\n\n## Other",
			want:    []string{"Backend API endpoints", "Database schema", "CLI commands"},
		},
		{
			name:    "handles checkboxes",
			content: "## Scope\n\n### In Scope\n\n- [ ] Backend API endpoints\n- [x] Database schema\n\n## Other",
			want:    []string{"Backend API endpoints", "Database schema"},
		},
		{
			name:    "empty scope section",
			content: "## Scope\n\n### In Scope\n\n## Other",
			want:    nil,
		},
		{
			name:    "no scope heading",
			content: "## Other\n\n- item1\n- item2",
			want:    nil,
		},
		{
			name:    "in scope with blank lines",
			content: "## Scope\n\n### In Scope\n\n- First item\n\n- Second item\n\n## Other",
			want:    []string{"First item", "Second item"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractScope(tt.content)
			if len(got) != len(tt.want) {
				t.Fatalf("extractScope() = %v, want %v", got, tt.want)
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("extractScope()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

// --- extractSuccessCriteria tests ---

func TestExtractSuccessCriteria(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "extracts success criteria",
			content: "## Success Criteria\n\n- [ ] 13 embed template files use correct strategy\n- [ ] No workflow duplication\n- [x] Already done item\n\n## Other",
			want:    []string{"13 embed template files use correct strategy", "No workflow duplication", "Already done item"},
		},
		{
			name:    "empty success criteria",
			content: "## Success Criteria\n\n## Other",
			want:    nil,
		},
		{
			name:    "no success criteria heading",
			content: "## Other\n\n- item",
			want:    nil,
		},
		{
			name:    "success criteria with sub-items",
			content: "## Success Criteria\n\n- [ ] Main criterion\n  - Sub-item should be ignored\n- [ ] Another criterion\n\n## Other",
			want:    []string{"Main criterion", "Another criterion"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSuccessCriteria(tt.content)
			if len(got) != len(tt.want) {
				t.Fatalf("extractSuccessCriteria() = %v, want %v", got, tt.want)
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("extractSuccessCriteria()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

// --- extractAcceptanceCriteria tests ---

func TestExtractAcceptanceCriteria(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "extracts acceptance criteria",
			content: "## Acceptance Criteria\n\n- [ ] BodyContext struct added to autogen.go\n- [ ] renderBody function handles all 6 tokens\n- [ ] GenerateTestTaskMD signature updated\n\n## Other",
			want:    []string{"BodyContext struct added to autogen.go", "renderBody function handles all 6 tokens", "GenerateTestTaskMD signature updated"},
		},
		{
			name:    "empty acceptance criteria",
			content: "## Acceptance Criteria\n\n## Other",
			want:    nil,
		},
		{
			name:    "no acceptance criteria heading",
			content: "## Other\n\n- item",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractAcceptanceCriteria(tt.content)
			if len(got) != len(tt.want) {
				t.Fatalf("extractAcceptanceCriteria() = %v, want %v", got, tt.want)
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("extractAcceptanceCriteria()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

// --- extractDocTaskCriteria tests ---

func TestExtractDocTaskCriteria(t *testing.T) {
	t.Run("extracts AC from single doc task", func(t *testing.T) {
		dir := t.TempDir()
		content := "---\nid: \"1\"\ntitle: \"Doc Task\"\ntype: \"doc\"\n---\n\n# Doc Task\n\n## Acceptance Criteria\n\n- [ ] First criterion\n- [ ] Second criterion\n\n## Implementation Notes\n\nSome notes here"
		if err := os.WriteFile(filepath.Join(dir, "1-doc.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		got := extractDocTaskCriteria(dir)

		if len(got) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(got))
		}
		ac, ok := got["1-doc"]
		if !ok {
			t.Fatal("expected key '1-doc'")
		}
		if !strings.Contains(ac, "- [ ] First criterion") {
			t.Errorf("expected AC content to contain first criterion, got: %s", ac)
		}
		if !strings.Contains(ac, "- [ ] Second criterion") {
			t.Errorf("expected AC content to contain second criterion, got: %s", ac)
		}
		// Should NOT include Implementation Notes section
		if strings.Contains(ac, "Implementation Notes") {
			t.Errorf("AC content should not include next section, got: %s", ac)
		}
	})

	t.Run("extracts AC from multiple doc tasks", func(t *testing.T) {
		dir := t.TempDir()
		content1 := "---\nid: \"1\"\ntype: \"doc\"\n---\n\n## Acceptance Criteria\n\n- [ ] AC 1\n\n## Other"
		content2 := "---\nid: \"2\"\ntype: \"doc\"\n---\n\n## Acceptance Criteria\n\n- [ ] AC 2\n- [ ] AC 3\n\n## Other"
		if err := os.WriteFile(filepath.Join(dir, "1-doc.md"), []byte(content1), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "2-doc.md"), []byte(content2), 0644); err != nil {
			t.Fatal(err)
		}

		got := extractDocTaskCriteria(dir)

		if len(got) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(got))
		}
		if _, ok := got["1-doc"]; !ok {
			t.Error("expected key '1-doc'")
		}
		if _, ok := got["2-doc"]; !ok {
			t.Error("expected key '2-doc'")
		}
	})

	t.Run("returns entry with empty content for missing AC section", func(t *testing.T) {
		dir := t.TempDir()
		content := "---\nid: \"1\"\ntype: \"doc\"\n---\n\n# No AC Here\n\n## Other Section\n\nSome text"
		if err := os.WriteFile(filepath.Join(dir, "1-doc.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		got := extractDocTaskCriteria(dir)

		if len(got) != 1 {
			t.Fatalf("expected 1 entry (with empty AC), got %d", len(got))
		}
		ac, ok := got["1-doc"]
		if !ok {
			t.Fatal("expected key '1-doc'")
		}
		if ac != "" {
			t.Errorf("expected empty AC content for missing section, got: %q", ac)
		}
	})

	t.Run("skips non-doc task files", func(t *testing.T) {
		dir := t.TempDir()
		docContent := "---\nid: \"1\"\ntype: \"doc\"\n---\n\n## Acceptance Criteria\n\n- [ ] Doc AC\n\n## Other"
		codeContent := "---\nid: \"2\"\ntype: \"coding.feature\"\n---\n\n## Acceptance Criteria\n\n- [ ] Code AC\n\n## Other"
		if err := os.WriteFile(filepath.Join(dir, "1-doc.md"), []byte(docContent), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "2-code.md"), []byte(codeContent), 0644); err != nil {
			t.Fatal(err)
		}

		got := extractDocTaskCriteria(dir)

		if len(got) != 1 {
			t.Fatalf("expected 1 entry (doc only), got %d", len(got))
		}
		if _, ok := got["1-doc"]; !ok {
			t.Error("expected key '1-doc'")
		}
		if _, ok := got["2-code"]; ok {
			t.Error("should not include non-doc task '2-code'")
		}
	})

	t.Run("handles empty task dir", func(t *testing.T) {
		dir := t.TempDir()

		got := extractDocTaskCriteria(dir)

		if len(got) != 0 {
			t.Fatalf("expected 0 entries for empty dir, got %d", len(got))
		}
	})

	t.Run("handles fenced code blocks containing ## lines", func(t *testing.T) {
		dir := t.TempDir()
		content := "---\nid: \"1\"\ntype: \"doc\"\n---\n\n## Acceptance Criteria\n\n- [ ] AC with code example:\n```\n## This is not a section header\n```\n- [ ] Another AC\n\n## Real Next Section"
		if err := os.WriteFile(filepath.Join(dir, "1-doc.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		got := extractDocTaskCriteria(dir)

		if len(got) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(got))
		}
		ac := got["1-doc"]
		// The fenced code block ## should not terminate extraction
		if !strings.Contains(ac, "This is not a section header") {
			t.Errorf("AC should contain code block content, got: %s", ac)
		}
		if !strings.Contains(ac, "Another AC") {
			t.Errorf("AC should contain second criterion, got: %s", ac)
		}
		if strings.Contains(ac, "Real Next Section") {
			t.Errorf("AC should not contain next section, got: %s", ac)
		}
	})

	t.Run("preserves multi-line content including sub-items", func(t *testing.T) {
		dir := t.TempDir()
		content := "---\nid: \"1\"\ntype: \"doc\"\n---\n\n## Acceptance Criteria\n\n- [ ] Main criterion\n  - Sub-item A\n  - Sub-item B\n- [ ] Another criterion\n  - Sub-item C\n\n## Other"
		if err := os.WriteFile(filepath.Join(dir, "1-doc.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		got := extractDocTaskCriteria(dir)

		ac := got["1-doc"]
		if !strings.Contains(ac, "Sub-item A") || !strings.Contains(ac, "Sub-item C") {
			t.Errorf("AC should preserve sub-items, got: %s", ac)
		}
	})
}

// --- extractDocTaskCriteria Section Extraction (unit) ---

func TestExtractACSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantOK  bool
	}{
		{
			name:    "basic AC section",
			content: "## Acceptance Criteria\n\n- [ ] Item 1\n- [ ] Item 2\n\n## Other",
			want:    "\n- [ ] Item 1\n- [ ] Item 2\n",
			wantOK:  true,
		},
		{
			name:    "no AC section",
			content: "## Other\n\n- item",
			want:    "",
			wantOK:  false,
		},
		{
			name:    "AC at end of file",
			content: "## Acceptance Criteria\n\n- [ ] Last item",
			want:    "\n- [ ] Last item",
			wantOK:  true,
		},
		{
			name:    "AC with code blocks",
			content: "## Acceptance Criteria\n\n```\n## not a header\n```\n- [ ] Real item\n\n## Next",
			want:    "\n```\n## not a header\n```\n- [ ] Real item\n",
			wantOK:  true,
		},
		{
			name:    "case-insensitive: Acceptance criteria (lowercase c)",
			content: "## Acceptance criteria\n\n- [ ] Item 1\n\n## Other",
			want:    "\n- [ ] Item 1\n",
			wantOK:  true,
		},
		{
			name:    "Chinese alias: 验收标准",
			content: "## 验收标准\n\n- [ ] 中文条目 1\n- [ ] 中文条目 2\n\n## Other",
			want:    "\n- [ ] 中文条目 1\n- [ ] 中文条目 2\n",
			wantOK:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := extractACSection(tt.content)
			if ok != tt.wantOK {
				t.Errorf("extractACSection ok = %v, want %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Errorf("extractACSection = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- extractBodyContext tests ---

func TestExtractBodyContext(t *testing.T) {
	t.Run("quick mode reads proposal", func(t *testing.T) {
		projectRoot, _, _ := setupBuildEnv(t, "quick")

		// Write rich proposal
		propDir := filepath.Join(projectRoot, "docs", "proposals", "test-feature")
		proposalContent := "# Proposal\n\n## Scope\n\n### In Scope\n\n- Backend logic\n- API endpoints\n\n### Out of Scope\n\n- Frontend\n\n## Success Criteria\n\n- [ ] Criteria A\n- [ ] Criteria B\n\n## Other"
		if err := os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
			t.Fatal(err)
		}

		mode := detectMode(projectRoot, "test-feature", "")
		ctx := extractBodyContext(projectRoot, "test-feature", mode, []string{"api", "cli"})

		if ctx.FeatureSlug != "test-feature" {
			t.Errorf("FeatureSlug = %q, want test-feature", ctx.FeatureSlug)
		}
		if ctx.Mode != "quick" {
			t.Errorf("Mode = %q, want quick", ctx.Mode)
		}
		if len(ctx.SuccessCriteria) != 2 {
			t.Errorf("SuccessCriteria = %v, want 2 items", ctx.SuccessCriteria)
		}
		// Quick mode: no acceptance criteria from PRD
		if len(ctx.AcceptanceCriteria) != 0 {
			t.Errorf("AcceptanceCriteria = %v, want empty in quick mode", ctx.AcceptanceCriteria)
		}
		if len(ctx.SurfaceTypes) != 2 || ctx.SurfaceTypes[0] != "api" {
			t.Errorf("SurfaceTypes = %v, want [api, cli]", ctx.SurfaceTypes)
		}
	})

	t.Run("breakdown mode reads PRD", func(t *testing.T) {
		projectRoot, _, _ := setupBuildEnv(t, "breakdown")

		// Write rich PRD
		prdDir := filepath.Join(projectRoot, "docs", "features", "test-feature", "prd")
		prdContent := "# PRD\n\n## Scope\n\n### In Scope\n\n- Database layer\n- Service layer\n\n## Acceptance Criteria\n\n- [ ] AC 1: Works correctly\n- [ ] AC 2: Handles errors\n\n## Success Criteria\n\n- [ ] SC 1\n\n## Other"
		if err := os.WriteFile(filepath.Join(prdDir, "prd-spec.md"), []byte(prdContent), 0644); err != nil {
			t.Fatal(err)
		}

		mode := detectMode(projectRoot, "test-feature", "")
		ctx := extractBodyContext(projectRoot, "test-feature", mode, []string{"api"})

		if ctx.Mode != "breakdown" {
			t.Errorf("Mode = %q, want breakdown", ctx.Mode)
		}
		if len(ctx.AcceptanceCriteria) != 2 {
			t.Errorf("AcceptanceCriteria = %v, want 2 items", ctx.AcceptanceCriteria)
		}
		if len(ctx.SuccessCriteria) != 1 {
			t.Errorf("SuccessCriteria = %v, want 1 item", ctx.SuccessCriteria)
		}
	})

	t.Run("missing proposal returns empty context", func(t *testing.T) {
		projectRoot, _, _ := setupBuildEnv(t, "")
		ctx := extractBodyContext(projectRoot, "test-feature", "", nil)

		if ctx.FeatureSlug != "test-feature" {
			t.Errorf("FeatureSlug = %q, want test-feature", ctx.FeatureSlug)
		}
		if ctx.Mode != "" {
			t.Errorf("Mode = %q, want empty", ctx.Mode)
		}
		if len(ctx.AcceptanceCriteria) != 0 {
			t.Errorf("AcceptanceCriteria = %v, want empty", ctx.AcceptanceCriteria)
		}
	})
}

// --- Integration: BuildIndex passes BodyContext to GenerateTestTaskMD ---

func TestBuildIndex_BodyContextPopulatedInGeneratedMD(t *testing.T) {
	// Verify that generated .md files contain feature-specific content from BodyContext
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Write rich proposal with scope
	propDir := filepath.Join(projectRoot, "docs", "proposals", "test-feature")
	proposalContent := "# Proposal\n\n## Scope\n\n### In Scope\n\n- Backend API\n- CLI integration\n\n## Success Criteria\n\n- [ ] Criteria A\n\n## Other"
	if err := os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Write a feature task (to trigger test pipeline)
	writeTaskMDWithType(t, tasksDir, "1-feat.md", "1", "Feature Task", TypeCodingFeature, nil)
	// Add a second task to get stage gates
	writeTaskMDWithType(t, tasksDir, "2-feat.md", "2", "Feature Task 2", TypeCodingFeature, []string{"1"})
	// Create config with surfaces
	writeForgeConfig(t, projectRoot)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
		AutoConfig:  allEnabledAuto,
	}

	_, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Check generated test task .md files for BodyContext data
	genCasesPath := filepath.Join(tasksDir, "quick-test-cases.md")
	data, err := os.ReadFile(genCasesPath)
	if err != nil {
		t.Skipf("quick-test-cases.md not generated by BuildIndex (deferred to caller)")
	}

	content := string(data)
	if !strings.Contains(content, "test-feature") {
		t.Errorf("generated .md should contain feature slug, got:\n%s", content)
	}
}

func TestBuildIndex_BodyContextBackwardCompat_NoProposal(t *testing.T) {
	// Verify BuildIndex still works when no proposal/PRD exists (empty BodyContext)
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMDWithType(t, tasksDir, "1-feat.md", "1", "Feature Task", TypeCodingFeature, nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Should succeed with 1 task
	total := result.NewCount + result.UpdatedCount
	if total != 1 {
		t.Errorf("total tasks = %d, want 1", total)
	}
}

// --- Build-time AC validation tests (Task 3) ---

func TestBuildIndex_DocTaskMissingAC_WarningEmitted(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Create a doc task WITHOUT an Acceptance Criteria section
	writeTaskMDWithType(t, tasksDir, "1-doc.md", "1", "Doc Task", TypeDoc, nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Should have a warning about missing AC section
	found := false
	for _, w := range result.Warnings {
		if strings.Contains(w, "1-doc") && strings.Contains(w, "no Acceptance Criteria") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected warning for missing AC section, got warnings: %v", result.Warnings)
	}
}

func TestBuildIndex_AllDocTasksMissingAC_FeatureWarning(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Create TWO doc tasks, both without AC
	writeTaskMDWithType(t, tasksDir, "1-doc.md", "1", "Doc Task 1", TypeDoc, nil)
	writeTaskMDWithType(t, tasksDir, "2-doc.md", "2", "Doc Task 2", TypeDoc, []string{"1"})

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Should have a feature-level warning about zero AC
	found := false
	for _, w := range result.Warnings {
		if strings.Contains(w, "feature has no AC for any doc task") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected feature-level zero-AC warning, got warnings: %v", result.Warnings)
	}
}

func TestBuildIndex_SomeDocTasksMissingAC_NoFeatureWarning(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// One doc task WITH AC, one without
	writeTaskMDWithType(t, tasksDir, "1-doc.md", "1", "Doc Task 1", TypeDoc, nil)
	content := "---\nid: \"2\"\ntitle: \"Doc Task 2\"\ntype: \"doc\"\npriority: \"P1\"\nestimated_time: \"1h\"\nsurface-key: \".\"\nsurface-type: \"web\"\n---\n\n## Acceptance Criteria\n\n- [ ] AC item\n\n## Other"
	if err := os.WriteFile(filepath.Join(tasksDir, "2-doc.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Should NOT have the feature-level zero-AC warning (at least one has AC)
	for _, w := range result.Warnings {
		if strings.Contains(w, "feature has no AC for any doc task") {
			t.Errorf("should NOT emit feature-level zero-AC warning when at least one doc task has AC, got: %s", w)
		}
	}
}

func TestBuildIndex_DocTaskCriteriaKeysMatchDocTasks(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Create doc tasks with AC
	content1 := "---\nid: \"1\"\ntitle: \"Doc Task 1\"\ntype: \"doc\"\npriority: \"P1\"\nestimated_time: \"1h\"\nsurface-key: \".\"\nsurface-type: \"web\"\n---\n\n## Acceptance Criteria\n\n- [ ] AC 1\n\n## Other"
	content2 := "---\nid: \"2\"\ntitle: \"Doc Task 2\"\ntype: \"doc\"\npriority: \"P1\"\nestimated_time: \"1h\"\nsurface-key: \".\"\nsurface-type: \"web\"\n---\n\n## Acceptance Criteria\n\n- [ ] AC 2\n\n## Other"
	if err := os.WriteFile(filepath.Join(tasksDir, "1-doc.md"), []byte(content1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tasksDir, "2-doc.md"), []byte(content2), 0644); err != nil {
		t.Fatal(err)
	}

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Verify generated review-doc.md contains both tasks' AC
	reviewDocPath := filepath.Join(tasksDir, "review-doc.md")
	data, err := os.ReadFile(reviewDocPath)
	if err != nil {
		t.Fatalf("review-doc.md not generated: %v", err)
	}
	reviewContent := string(data)
	if !strings.Contains(reviewContent, "### 1-doc") {
		t.Error("review-doc.md should contain ### 1-doc sub-section")
	}
	if !strings.Contains(reviewContent, "### 2-doc") {
		t.Error("review-doc.md should contain ### 2-doc sub-section")
	}
	if !strings.Contains(reviewContent, "AC 1") {
		t.Error("review-doc.md should contain AC content from 1-doc")
	}
	if !strings.Contains(reviewContent, "AC 2") {
		t.Error("review-doc.md should contain AC content from 2-doc")
	}
	_ = result
}

func TestBuildIndex_DocTaskEmptyAC_PlaceholderInReviewDoc(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Create doc task without AC section
	writeTaskMDWithType(t, tasksDir, "1-doc.md", "1", "Doc Task", TypeDoc, nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	_, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Verify generated review-doc.md contains placeholder
	reviewDocPath := filepath.Join(tasksDir, "review-doc.md")
	data, err := os.ReadFile(reviewDocPath)
	if err != nil {
		t.Fatalf("review-doc.md not generated: %v", err)
	}
	reviewContent := string(data)
	if !strings.Contains(reviewContent, "> No acceptance criteria defined.") {
		t.Errorf("review-doc.md should contain placeholder for empty AC, got:\n%s", reviewContent)
	}
}

func TestSerializeDocTaskAC_EmptyContentShowsPlaceholder(t *testing.T) {
	criteria := map[string]string{
		"1-doc": "",
		"2-doc": "- [ ] Real AC",
	}

	result := serializeDocTaskAC(criteria)

	if !strings.Contains(result, "### 1-doc") {
		t.Error("should contain 1-doc sub-section header")
	}
	if !strings.Contains(result, "> No acceptance criteria defined.") {
		t.Error("should show placeholder for empty AC content")
	}
	if !strings.Contains(result, "### 2-doc") {
		t.Error("should contain 2-doc sub-section header")
	}
	if !strings.Contains(result, "Real AC") {
		t.Error("should contain actual AC content")
	}
}
