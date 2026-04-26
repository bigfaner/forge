package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(tempDir string) (workDir string)
		wantMarker string
		wantType   RootType
	}{
		{
			name: "finds go.mod in current directory",
			setup: func(tempDir string) string {
				goModPath := filepath.Join(tempDir, "go.mod")
				os.WriteFile(goModPath, []byte("module test\n"), 0644)
				return tempDir
			},
			wantMarker: "go.mod",
			wantType:   RootTypeProject,
		},
		{
			name: "finds go.mod in parent directory",
			setup: func(tempDir string) string {
				subDir := filepath.Join(tempDir, "subdir")
				os.MkdirAll(subDir, 0755)
				goModPath := filepath.Join(tempDir, "go.mod")
				os.WriteFile(goModPath, []byte("module test\n"), 0644)
				return subDir
			},
			wantMarker: "go.mod",
			wantType:   RootTypeProject,
		},
		{
			name: "finds package.json for Node.js project",
			setup: func(tempDir string) string {
				pkgPath := filepath.Join(tempDir, "package.json")
				os.WriteFile(pkgPath, []byte(`{"name": "test"}`), 0644)
				return tempDir
			},
			wantMarker: "package.json",
			wantType:   RootTypeProject,
		},
		{
			name: "finds Cargo.toml for Rust project",
			setup: func(tempDir string) string {
				cargoPath := filepath.Join(tempDir, "Cargo.toml")
				os.WriteFile(cargoPath, []byte(`[package]\nname = "test"`), 0644)
				return tempDir
			},
			wantMarker: "Cargo.toml",
			wantType:   RootTypeProject,
		},
		{
			name: "finds pyproject.toml for Python project",
			setup: func(tempDir string) string {
				pyPath := filepath.Join(tempDir, "pyproject.toml")
				os.WriteFile(pyPath, []byte(`[project]\nname = "test"`), 0644)
				return tempDir
			},
			wantMarker: "pyproject.toml",
			wantType:   RootTypeProject,
		},
		{
			name: "finds pom.xml for Java/Maven project",
			setup: func(tempDir string) string {
				pomPath := filepath.Join(tempDir, "pom.xml")
				os.WriteFile(pomPath, []byte(`<project></project>`), 0644)
				return tempDir
			},
			wantMarker: "pom.xml",
			wantType:   RootTypeProject,
		},
		{
			name: "finds build.gradle for Java/Gradle project",
			setup: func(tempDir string) string {
				gradlePath := filepath.Join(tempDir, "build.gradle")
				os.WriteFile(gradlePath, []byte(`plugins { id("java") }`), 0644)
				return tempDir
			},
			wantMarker: "build.gradle",
			wantType:   RootTypeProject,
		},
		{
			name: "prefers workspace marker over project marker",
			setup: func(tempDir string) string {
				// Create workspace marker at root
				goWorkPath := filepath.Join(tempDir, "go.work")
				os.WriteFile(goWorkPath, []byte(`go 1.21`), 0644)
				// Create project marker in subdirectory
				subDir := filepath.Join(tempDir, "service", "auth")
				os.MkdirAll(subDir, 0755)
				goModPath := filepath.Join(subDir, "go.mod")
				os.WriteFile(goModPath, []byte("module test\n"), 0644)
				return subDir
			},
			wantMarker: "go.work",
			wantType:   RootTypeWorkspace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			workDir := tt.setup(tempDir)

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			if err := os.Chdir(workDir); err != nil {
				t.Fatalf("Chdir() error = %v", err)
			}

			root, err := FindProjectRoot()
			if err != nil {
				t.Fatalf("FindProjectRoot() error = %v", err)
			}
			if root == "" {
				t.Error("FindProjectRoot() returned empty path")
			}

			// Verify the marker type
			info, err := FindRootInfo()
			if err != nil {
				t.Fatalf("FindRootInfo() error = %v", err)
			}
			if info.Marker != tt.wantMarker {
				t.Errorf("Marker = %q, want %q", info.Marker, tt.wantMarker)
			}
			if info.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", info.Type, tt.wantType)
			}
		})
	}
}

func TestFindRootInfo(t *testing.T) {
	t.Run("returns detailed info for Go project", func(t *testing.T) {
		tempDir, _ := filepath.EvalSymlinks(t.TempDir())
		goModPath := filepath.Join(tempDir, "go.mod")
		os.WriteFile(goModPath, []byte("module test\n"), 0644)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Path != tempDir {
			t.Errorf("Path = %q, want %q", info.Path, tempDir)
		}
		if info.Type != RootTypeProject {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeProject)
		}
		if info.Marker != "go.mod" {
			t.Errorf("Marker = %q, want %q", info.Marker, "go.mod")
		}
		if len(info.Languages) == 0 || info.Languages[0] != "go" {
			t.Errorf("Languages = %v, want [go]", info.Languages)
		}
	})

	t.Run("returns workspace type for go.work", func(t *testing.T) {
		tempDir := t.TempDir()
		goWorkPath := filepath.Join(tempDir, "go.work")
		os.WriteFile(goWorkPath, []byte("go 1.21\n"), 0644)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Type != RootTypeWorkspace {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeWorkspace)
		}
		if info.Marker != "go.work" {
			t.Errorf("Marker = %q, want %q", info.Marker, "go.work")
		}
	})
}

func TestFindRootInfoFrom(t *testing.T) {
	t.Run("finds root from subdirectory", func(t *testing.T) {
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "a", "b", "c")
		os.MkdirAll(subDir, 0755)

		goModPath := filepath.Join(tempDir, "go.mod")
		os.WriteFile(goModPath, []byte("module test\n"), 0644)

		info, err := FindRootInfoFrom(subDir)
		if err != nil {
			t.Fatalf("FindRootInfoFrom() error = %v", err)
		}
		if info.Path != tempDir {
			t.Errorf("Path = %q, want %q", info.Path, tempDir)
		}
	})
}

func TestGetProjectRootFromEnv(t *testing.T) {
	t.Run("returns CLAUDE_PROJECT_DIR when set", func(t *testing.T) {
		testPath := filepath.Join("custom", "path")
		os.Setenv("CLAUDE_PROJECT_DIR", testPath)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		got := GetProjectRootFromEnv()
		want := filepath.Clean(testPath)
		if got != want {
			t.Errorf("GetProjectRootFromEnv() = %q, want %q", got, want)
		}
	})

	t.Run("returns PROJECT_ROOT when CLAUDE_PROJECT_DIR not set", func(t *testing.T) {
		testPath := filepath.Join("fallback", "path")
		os.Unsetenv("CLAUDE_PROJECT_DIR")
		os.Setenv("PROJECT_ROOT", testPath)
		defer os.Unsetenv("PROJECT_ROOT")

		got := GetProjectRootFromEnv()
		want := filepath.Clean(testPath)
		if got != want {
			t.Errorf("GetProjectRootFromEnv() = %q, want %q", got, want)
		}
	})

	t.Run("CLAUDE_PROJECT_DIR takes priority", func(t *testing.T) {
		priorityPath := filepath.Join("priority")
		fallbackPath := filepath.Join("fallback")
		os.Setenv("CLAUDE_PROJECT_DIR", priorityPath)
		os.Setenv("PROJECT_ROOT", fallbackPath)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")
		defer os.Unsetenv("PROJECT_ROOT")

		got := GetProjectRootFromEnv()
		want := filepath.Clean(priorityPath)
		if got != want {
			t.Errorf("GetProjectRootFromEnv() = %q, want %q", got, want)
		}
	})

	t.Run("returns empty when neither set", func(t *testing.T) {
		os.Unsetenv("CLAUDE_PROJECT_DIR")
		os.Unsetenv("PROJECT_ROOT")

		got := GetProjectRootFromEnv()
		if got != "" {
			t.Errorf("GetProjectRootFromEnv() = %q, want empty", got)
		}
	})
}

func TestEnvVarOverride(t *testing.T) {
	t.Run("environment variable overrides detection", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create go.mod in temp dir
		goModPath := filepath.Join(tempDir, "go.mod")
		os.WriteFile(goModPath, []byte("module test\n"), 0644)

		// Set env var to different path
		overridePath := filepath.Join("override", "path")
		os.Setenv("CLAUDE_PROJECT_DIR", overridePath)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		root, err := FindProjectRoot()
		if err != nil {
			t.Fatalf("FindProjectRoot() error = %v", err)
		}
		want := filepath.Clean(overridePath)
		if root != want {
			t.Errorf("FindProjectRoot() = %q, want %q (env override)", root, want)
		}
	})
}

func TestGitWorktree(t *testing.T) {
	t.Run("matchesMarker accepts .git as file (worktree)", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create .git as a file (worktree scenario)
		gitFile := filepath.Join(tempDir, ".git")
		os.WriteFile(gitFile, []byte("gitdir: /main/repo/.git/worktrees/feature\n"), 0644)

		// Test that matchesMarker accepts .git as a file
		gitMarker := Marker{Name: ".git", Type: RootTypeVCS, IsDirectory: false}
		if !matchesMarker(tempDir, gitMarker) {
			t.Error("matchesMarker() should accept .git as a file")
		}
	})

	t.Run("matchesMarker accepts .git as directory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create .git as a directory (normal repo)
		gitDir := filepath.Join(tempDir, ".git")
		os.MkdirAll(gitDir, 0755)

		// Test that matchesMarker accepts .git as a directory
		gitMarker := Marker{Name: ".git", Type: RootTypeVCS, IsDirectory: false}
		if !matchesMarker(tempDir, gitMarker) {
			t.Error("matchesMarker() should accept .git as a directory")
		}
	})
}

func TestFindVCSRoot(t *testing.T) {
	t.Run("finds .git directory", func(t *testing.T) {
		tempDir, _ := filepath.EvalSymlinks(t.TempDir())
		gitDir := filepath.Join(tempDir, ".git")
		os.MkdirAll(gitDir, 0755)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		root, err := FindVCSRoot()
		if err != nil {
			t.Fatalf("FindVCSRoot() error = %v", err)
		}
		if root != tempDir {
			t.Errorf("FindVCSRoot() = %q, want %q", root, tempDir)
		}
	})

	t.Run("finds .git in parent", func(t *testing.T) {
		tempDir, _ := filepath.EvalSymlinks(t.TempDir())
		subDir := filepath.Join(tempDir, "subdir")
		os.MkdirAll(subDir, 0755)
		gitDir := filepath.Join(tempDir, ".git")
		os.MkdirAll(gitDir, 0755)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(subDir)

		root, err := FindVCSRoot()
		if err != nil {
			t.Fatalf("FindVCSRoot() error = %v", err)
		}
		if root != tempDir {
			t.Errorf("FindVCSRoot() = %q, want %q", root, tempDir)
		}
	})
}

func TestMonorepoDetection(t *testing.T) {
	t.Run(".forge dir overrides go.mod in subdirectory", func(t *testing.T) {
		tempDir, _ := filepath.EvalSymlinks(t.TempDir())

		// Create .forge/ at project root (simulates task claim bootstrap)
		forgeDir := filepath.Join(tempDir, ".forge")
		os.MkdirAll(forgeDir, 0755)

		// Create go.mod in backend/ subdirectory
		backendDir := filepath.Join(tempDir, "backend")
		os.MkdirAll(backendDir, 0755)
		goModPath := filepath.Join(backendDir, "go.mod")
		os.WriteFile(goModPath, []byte("module backend\n"), 0644)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(backendDir)

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Path != tempDir {
			t.Errorf("Path = %q, want %q (project root, not backend/)", info.Path, tempDir)
		}
		if info.Type != RootTypeWorkspace {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeWorkspace)
		}
		if info.Marker != ".forge" {
			t.Errorf("Marker = %q, want %q", info.Marker, ".forge")
		}
	})

	t.Run(".forge dir overrides package.json in subdirectory", func(t *testing.T) {
		tempDir, _ := filepath.EvalSymlinks(t.TempDir())

		forgeDir := filepath.Join(tempDir, ".forge")
		os.MkdirAll(forgeDir, 0755)

		frontendDir := filepath.Join(tempDir, "frontend")
		os.MkdirAll(frontendDir, 0755)
		pkgJson := filepath.Join(frontendDir, "package.json")
		os.WriteFile(pkgJson, []byte(`{"name": "frontend"}`), 0644)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(frontendDir)

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Path != tempDir {
			t.Errorf("Path = %q, want %q (project root, not frontend/)", info.Path, tempDir)
		}
		if info.Type != RootTypeWorkspace {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeWorkspace)
		}
		if info.Marker != ".forge" {
			t.Errorf("Marker = %q, want %q", info.Marker, ".forge")
		}
	})

	t.Run(".forge and go.mod both at root — .forge wins", func(t *testing.T) {
		tempDir, _ := filepath.EvalSymlinks(t.TempDir())

		forgeDir := filepath.Join(tempDir, ".forge")
		os.MkdirAll(forgeDir, 0755)
		goModPath := filepath.Join(tempDir, "go.mod")
		os.WriteFile(goModPath, []byte("module test\n"), 0644)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Path != tempDir {
			t.Errorf("Path = %q, want %q", info.Path, tempDir)
		}
		if info.Type != RootTypeWorkspace {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeWorkspace)
		}
	})

	t.Run("pnpm monorepo returns workspace root", func(t *testing.T) {
		tempDir, _ := filepath.EvalSymlinks(t.TempDir())

		// Create pnpm workspace at root
		pnpmWorkspace := filepath.Join(tempDir, "pnpm-workspace.yaml")
		os.WriteFile(pnpmWorkspace, []byte("packages:\n  - 'apps/*'\n"), 0644)

		// Create package.json in subproject
		subDir := filepath.Join(tempDir, "apps", "web")
		os.MkdirAll(subDir, 0755)
		pkgJson := filepath.Join(subDir, "package.json")
		os.WriteFile(pkgJson, []byte(`{"name": "web"}`), 0644)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(subDir)

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Type != RootTypeWorkspace {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeWorkspace)
		}
		if info.Path != tempDir {
			t.Errorf("Path = %q, want %q (workspace root)", info.Path, tempDir)
		}
	})

	t.Run("Gradle multi-project returns settings.gradle location", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create settings.gradle at root
		settingsGradle := filepath.Join(tempDir, "settings.gradle")
		os.WriteFile(settingsGradle, []byte("include 'service-auth'"), 0644)

		// Create build.gradle in subproject
		subDir := filepath.Join(tempDir, "service-auth")
		os.MkdirAll(subDir, 0755)
		buildGradle := filepath.Join(subDir, "build.gradle")
		os.WriteFile(buildGradle, []byte("plugins { id('java') }"), 0644)

		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(subDir)

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Type != RootTypeWorkspace {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeWorkspace)
		}
		if info.Marker != "settings.gradle" {
			t.Errorf("Marker = %q, want %q", info.Marker, "settings.gradle")
		}
	})
}

func TestRootTypeString(t *testing.T) {
	tests := []struct {
		typ  RootType
		want string
	}{
		{RootTypeUnknown, "unknown"},
		{RootTypeVCS, "vcs"},
		{RootTypeWorkspace, "workspace"},
		{RootTypeProject, "project"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.typ.String(); got != tt.want {
				t.Errorf("RootType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
