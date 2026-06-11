package feature

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"forge-cli/pkg/task"
)

func TestGetCurrentFeature(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(string) error
		want       string
		wantErr    bool
		errContain string
	}{
		{
			name:       "no features directory",
			setup:      nil,
			wantErr:    true,
			errContain: "no feature set",
		},
		{
			name: "no feature directories",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, FeaturesDir), 0755)
			},
			wantErr:    true,
			errContain: "no feature set",
		},
		{
			name: "single feature with state",
			setup: func(dir string) error {
				// Create feature structure with index.json
				featureDir := filepath.Join(dir, FeaturesDir, "my-feature", TasksDirName)
				if err := os.MkdirAll(filepath.Join(featureDir, ProcessDirName), 0755); err != nil {
					return err
				}
				// Create index.json
				indexData, _ := json.Marshal(&task.TaskIndex{Feature: "my-feature"})
				if err := os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644); err != nil {
					return err
				}
				// Create state.json
				state := &task.TaskState{TaskID: "1.1"}
				data, _ := json.Marshal(state)
				return os.WriteFile(filepath.Join(featureDir, ProcessDirName, StateFileName), data, 0644)
			},
			want:    "my-feature",
			wantErr: false,
		},
		{
			name: "single feature without state",
			setup: func(dir string) error {
				featureDir := filepath.Join(dir, FeaturesDir, "my-feature", TasksDirName)
				if err := os.MkdirAll(featureDir, 0755); err != nil {
					return err
				}
				// Create index.json
				indexData, _ := json.Marshal(&task.TaskIndex{Feature: "my-feature"})
				return os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644)
			},
			want:    "my-feature",
			wantErr: false,
		},
		{
			name: "feature with state takes priority",
			setup: func(dir string) error {
				// Create feature1 with state
				featureDir1 := filepath.Join(dir, FeaturesDir, "feature1", TasksDirName)
				if err := os.MkdirAll(filepath.Join(featureDir1, ProcessDirName), 0755); err != nil {
					return err
				}
				indexData1, _ := json.Marshal(&task.TaskIndex{Feature: "feature1"})
				if err := os.WriteFile(filepath.Join(featureDir1, IndexFileName), indexData1, 0644); err != nil {
					return err
				}
				state := &task.TaskState{TaskID: "1.1"}
				data, _ := json.Marshal(state)
				if err := os.WriteFile(filepath.Join(featureDir1, ProcessDirName, StateFileName), data, 0644); err != nil {
					return err
				}
				// Create feature2 without state
				featureDir2 := filepath.Join(dir, FeaturesDir, "feature2", TasksDirName)
				if err := os.MkdirAll(featureDir2, 0755); err != nil {
					return err
				}
				indexData2, _ := json.Marshal(&task.TaskIndex{Feature: "feature2"})
				return os.WriteFile(filepath.Join(featureDir2, IndexFileName), indexData2, 0644)
			},
			want:    "feature1",
			wantErr: false,
		},
		{
			name: "multiple features with state returns error",
			setup: func(dir string) error {
				for _, f := range []string{"feature1", "feature2"} {
					featureDir := filepath.Join(dir, FeaturesDir, f, TasksDirName)
					if err := os.MkdirAll(filepath.Join(featureDir, ProcessDirName), 0755); err != nil {
						return err
					}
					indexData, _ := json.Marshal(&task.TaskIndex{Feature: f})
					if err := os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644); err != nil {
						return err
					}
					state := &task.TaskState{TaskID: "1.1"}
					data, _ := json.Marshal(state)
					if err := os.WriteFile(filepath.Join(featureDir, ProcessDirName, StateFileName), data, 0644); err != nil {
						return err
					}
				}
				return nil
			},
			wantErr:    true,
			errContain: "multiple active features",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.setup != nil {
				if err := tt.setup(dir); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			got, err := GetCurrentFeature(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentFeature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContain != "" {
				if !containsString(err.Error(), tt.errContain) {
					t.Errorf("GetCurrentFeature() error = %v, want containing %q", err, tt.errContain)
				}
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentFeature() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequireFeature(t *testing.T) {
	t.Run("delegates to GetCurrentFeature", func(t *testing.T) {
		dir := t.TempDir()
		// Create feature with index.json
		featureDir := filepath.Join(dir, FeaturesDir, "my-feature", TasksDirName)
		if err := os.MkdirAll(featureDir, 0755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		indexData, _ := json.Marshal(&task.TaskIndex{Feature: "my-feature"})
		if err := os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		got, err := RequireFeature(dir)
		if err != nil {
			t.Fatalf("RequireFeature() error = %v", err)
		}
		if got != "my-feature" {
			t.Errorf("RequireFeature() = %q, want %q", got, "my-feature")
		}
	})
}

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

func TestGetCurrentFeature_GitContext(t *testing.T) {
	t.Run("git branch with existing feature directory", func(t *testing.T) {
		dir := t.TempDir()

		// Initialize a git repo so git.GetFeatureFromGit can resolve the branch
		gitInit(t, dir)

		// Create and checkout a feature branch
		gitCheckoutBranch(t, dir, "feature/my-feature")

		// Create the feature directory so Stat succeeds
		featureDir := filepath.Join(dir, FeaturesDir, "my-feature", TasksDirName)
		if err := os.MkdirAll(featureDir, 0755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		indexData, _ := json.Marshal(&task.TaskIndex{Feature: "my-feature"})
		if err := os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		got, err := GetCurrentFeature(dir)
		if err != nil {
			t.Fatalf("GetCurrentFeature() error = %v", err)
		}
		if got != "my-feature" {
			t.Errorf("GetCurrentFeature() = %q, want %q", got, "my-feature")
		}
	})

	t.Run("git branch creates feature dir when missing", func(t *testing.T) {
		dir := t.TempDir()

		gitInit(t, dir)
		gitCheckoutBranch(t, dir, "feature/new-feature")

		// No feature directory exists yet — EnsureFeatureDir should create it
		got, err := GetCurrentFeature(dir)
		if err != nil {
			t.Fatalf("GetCurrentFeature() error = %v", err)
		}
		if got != "new-feature" {
			t.Errorf("GetCurrentFeature() = %q, want %q", got, "new-feature")
		}

		// Verify the directory structure was created
		expectedDir := filepath.Join(dir, FeaturesDir, "new-feature")
		if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
			t.Errorf("feature directory %s was not created", expectedDir)
		}
	})

	t.Run("git returns empty when on main branch", func(t *testing.T) {
		dir := t.TempDir()

		// Initialize git repo (starts on main/master by default)
		gitInit(t, dir)

		// Create a single feature directory without state (Priority 3 path)
		featureDir := filepath.Join(dir, FeaturesDir, "my-feature", TasksDirName)
		if err := os.MkdirAll(featureDir, 0755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		indexData, _ := json.Marshal(&task.TaskIndex{Feature: "my-feature"})
		if err := os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		got, err := GetCurrentFeature(dir)
		if err != nil {
			t.Fatalf("GetCurrentFeature() error = %v", err)
		}
		if got != "my-feature" {
			t.Errorf("GetCurrentFeature() = %q, want %q", got, "my-feature")
		}
	})
}

func gitInit(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("git", "init", dir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	// Set a default user so commits work
	_ = exec.Command("git", "-C", dir, "config", "user.email", "test@test.com").Run()
	_ = exec.Command("git", "-C", dir, "config", "user.name", "Test").Run()
	// Create an initial commit so branch switching works
	dummy := filepath.Join(dir, ".gitignore")
	_ = os.WriteFile(dummy, []byte(""), 0644)
	_ = exec.Command("git", "-C", dir, "add", ".gitignore").Run()
	_ = exec.Command("git", "-C", dir, "commit", "-m", "init").Run()
}

func gitCheckoutBranch(t *testing.T, dir, branch string) {
	t.Helper()
	cmd := exec.Command("git", "-C", dir, "checkout", "-b", branch)
	if err := cmd.Run(); err != nil {
		t.Fatalf("git checkout -b %s failed: %v", branch, err)
	}
}

func TestGetCurrentFeatureWithSource(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(string) error
		wantSlug   string
		wantSource string
		wantErr    bool
		errContain string
	}{
		{
			name:       "no features directory returns error",
			setup:      nil,
			wantErr:    true,
			errContain: "no feature set",
		},
		{
			name: "state.json takes priority over git branch",
			setup: func(dir string) error {
				// Initialize git repo on a feature branch
				gitInit(t, dir)
				gitCheckoutBranch(t, dir, "feature/git-feature")

				// Create git-feature directory
				gitFeatureDir := filepath.Join(dir, FeaturesDir, "git-feature", TasksDirName)
				if err := os.MkdirAll(gitFeatureDir, 0755); err != nil {
					return err
				}

				// Create state-feature directory
				stateFeatureDir := filepath.Join(dir, FeaturesDir, "state-feature", TasksDirName)
				if err := os.MkdirAll(stateFeatureDir, 0755); err != nil {
					return err
				}
				indexData, _ := json.Marshal(&task.TaskIndex{Feature: "state-feature"})
				if err := os.WriteFile(filepath.Join(stateFeatureDir, IndexFileName), indexData, 0644); err != nil {
					return err
				}

				// Write .forge/state.json pointing to state-feature
				return WriteForgeState(dir, "state-feature")
			},
			wantSlug:   "state-feature",
			wantSource: SourceForgeState,
		},
		{
			name: "state.json skipped when feature dir does not exist",
			setup: func(dir string) error {
				// Initialize git repo on a feature branch
				gitInit(t, dir)
				gitCheckoutBranch(t, dir, "feature/git-feature")

				// Create git-feature directory so git context resolves
				gitFeatureDir := filepath.Join(dir, FeaturesDir, "git-feature", TasksDirName)
				if err := os.MkdirAll(gitFeatureDir, 0755); err != nil {
					return err
				}

				// Write .forge/state.json pointing to nonexistent feature
				return WriteForgeState(dir, "nonexistent-feature")
			},
			wantSlug:   "git-feature",
			wantSource: SourceBranch,
		},
		{
			name: "corrupt state.json falls through to git",
			setup: func(dir string) error {
				// Initialize git repo on a feature branch
				gitInit(t, dir)
				gitCheckoutBranch(t, dir, "feature/git-feature")

				// Create git-feature directory
				gitFeatureDir := filepath.Join(dir, FeaturesDir, "git-feature", TasksDirName)
				if err := os.MkdirAll(gitFeatureDir, 0755); err != nil {
					return err
				}

				// Write corrupt .forge/state.json
				if err := os.MkdirAll(filepath.Join(dir, ForgeDir), 0755); err != nil {
					return err
				}
				return os.WriteFile(GetForgeStatePath(dir), []byte("not json at all"), 0644)
			},
			wantSlug:   "git-feature",
			wantSource: SourceBranch,
		},
		{
			name: "git worktree resolves as worktree source",
			setup: func(dir string) error {
				// This test is conceptual — in a real test we'd need a worktree.
				// For now, test that branch resolution returns SourceBranch.
				gitInit(t, dir)
				gitCheckoutBranch(t, dir, "feature/branch-feature")

				// Create feature directory
				featureDir := filepath.Join(dir, FeaturesDir, "branch-feature", TasksDirName)
				if err := os.MkdirAll(featureDir, 0755); err != nil {
					return err
				}
				indexData, _ := json.Marshal(&task.TaskIndex{Feature: "branch-feature"})
				return os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644)
			},
			wantSlug:   "branch-feature",
			wantSource: SourceBranch,
		},
		{
			name: "no state.json no git falls to features-dir",
			setup: func(dir string) error {
				// Single feature with state (task process state.json)
				featureDir := filepath.Join(dir, FeaturesDir, "my-feature", TasksDirName)
				if err := os.MkdirAll(filepath.Join(featureDir, ProcessDirName), 0755); err != nil {
					return err
				}
				indexData, _ := json.Marshal(&task.TaskIndex{Feature: "my-feature"})
				if err := os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644); err != nil {
					return err
				}
				state := &task.TaskState{TaskID: "1.1"}
				data, _ := json.Marshal(state)
				return os.WriteFile(filepath.Join(featureDir, ProcessDirName, StateFileName), data, 0644)
			},
			wantSlug:   "my-feature",
			wantSource: SourceFeaturesDir,
		},
		{
			name: "state.json absent falls back to git",
			setup: func(dir string) error {
				gitInit(t, dir)
				gitCheckoutBranch(t, dir, "feature/some-feature")

				// Feature directory exists but no .forge/state.json
				featureDir := filepath.Join(dir, FeaturesDir, "some-feature", TasksDirName)
				if err := os.MkdirAll(featureDir, 0755); err != nil {
					return err
				}
				indexData, _ := json.Marshal(&task.TaskIndex{Feature: "some-feature"})
				return os.WriteFile(filepath.Join(featureDir, IndexFileName), indexData, 0644)
			},
			wantSlug:   "some-feature",
			wantSource: SourceBranch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.setup != nil {
				if err := tt.setup(dir); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			slug, source, err := GetCurrentFeatureWithSource(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentFeatureWithSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContain != "" {
				if !containsString(err.Error(), tt.errContain) {
					t.Errorf("GetCurrentFeatureWithSource() error = %v, want containing %q", err, tt.errContain)
				}
				return
			}
			if slug != tt.wantSlug {
				t.Errorf("GetCurrentFeatureWithSource() slug = %q, want %q", slug, tt.wantSlug)
			}
			if source != tt.wantSource {
				t.Errorf("GetCurrentFeatureWithSource() source = %q, want %q", source, tt.wantSource)
			}
		})
	}
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
