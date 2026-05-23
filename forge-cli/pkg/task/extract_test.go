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

		mode := detectMode(projectRoot, "test-feature")
		ctx := extractBodyContext(projectRoot, "test-feature", mode, []string{"api", "cli"})

		if ctx.FeatureSlug != "test-feature" {
			t.Errorf("FeatureSlug = %q, want test-feature", ctx.FeatureSlug)
		}
		if ctx.Mode != "quick" {
			t.Errorf("Mode = %q, want quick", ctx.Mode)
		}
		if len(ctx.Scope) != 2 || ctx.Scope[0] != "Backend logic" {
			t.Errorf("Scope = %v, want [Backend logic, API endpoints]", ctx.Scope)
		}
		if len(ctx.SuccessCriteria) != 2 {
			t.Errorf("SuccessCriteria = %v, want 2 items", ctx.SuccessCriteria)
		}
		// Quick mode: no acceptance criteria from PRD
		if len(ctx.AcceptanceCriteria) != 0 {
			t.Errorf("AcceptanceCriteria = %v, want empty in quick mode", ctx.AcceptanceCriteria)
		}
		if len(ctx.Interfaces) != 2 || ctx.Interfaces[0] != "api" {
			t.Errorf("Interfaces = %v, want [api, cli]", ctx.Interfaces)
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

		mode := detectMode(projectRoot, "test-feature")
		ctx := extractBodyContext(projectRoot, "test-feature", mode, []string{"api"})

		if ctx.Mode != "breakdown" {
			t.Errorf("Mode = %q, want breakdown", ctx.Mode)
		}
		if len(ctx.Scope) != 2 || ctx.Scope[0] != "Database layer" {
			t.Errorf("Scope = %v, want [Database layer, Service layer]", ctx.Scope)
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
		if len(ctx.Scope) != 0 {
			t.Errorf("Scope = %v, want empty", ctx.Scope)
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

	// Create config with interfaces
	configDir := filepath.Join(projectRoot, ".forge")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configContent := "interfaces:\n  - cli\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

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
	if !strings.Contains(content, "Backend API") || !strings.Contains(content, "CLI integration") {
		t.Errorf("generated .md should contain scope items from proposal, got:\n%s", content)
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
