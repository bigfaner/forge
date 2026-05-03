package feature

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"task-cli/pkg/task"
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

func TestSetFeature(t *testing.T) {
	t.Run("creates feature directory structure", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetFeature(dir, "test-feature"); err != nil {
			t.Fatalf("SetFeature() error = %v", err)
		}

		expectedDirs := []string{
			filepath.Join(dir, GetFeatureDir("test-feature")),
			filepath.Join(dir, GetFeaturePRDDir("test-feature")),
			filepath.Join(dir, GetFeatureDesignDir("test-feature")),
			filepath.Join(dir, GetFeatureUIDesignDir("test-feature")),
			filepath.Join(dir, GetFeatureTasksDir("test-feature")),
			filepath.Join(dir, GetFeatureRecordsDir("test-feature")),
			filepath.Join(dir, FeaturesDir, "test-feature", TasksDirName, ProcessDirName),
		}

		for _, expectedDir := range expectedDirs {
			if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
				t.Errorf("directory %s was not created", expectedDir)
			}
		}
	})

	t.Run("idempotent", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetFeature(dir, "test-feature"); err != nil {
			t.Fatalf("first SetFeature() error = %v", err)
		}
		if err := SetFeature(dir, "test-feature"); err != nil {
			t.Fatalf("second SetFeature() error = %v", err)
		}
	})
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
	exec.Command("git", "-C", dir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Test").Run()
	// Create an initial commit so branch switching works
	dummy := filepath.Join(dir, ".gitignore")
	os.WriteFile(dummy, []byte(""), 0644)
	exec.Command("git", "-C", dir, "add", ".gitignore").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "init").Run()
}

func gitCheckoutBranch(t *testing.T, dir, branch string) {
	t.Helper()
	cmd := exec.Command("git", "-C", dir, "checkout", "-b", branch)
	if err := cmd.Run(); err != nil {
		t.Fatalf("git checkout -b %s failed: %v", branch, err)
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
