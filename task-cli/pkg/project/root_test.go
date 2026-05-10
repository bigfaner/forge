package project

import (
	"os"
	"os/exec"
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
				_ = os.WriteFile(goModPath, []byte("module test\n"), 0644)
				return tempDir
			},
			wantMarker: "go.mod",
			wantType:   RootTypeProject,
		},
		{
			name: "finds go.mod in parent directory",
			setup: func(tempDir string) string {
				subDir := filepath.Join(tempDir, "subdir")
				_ = os.MkdirAll(subDir, 0755)
				goModPath := filepath.Join(tempDir, "go.mod")
				_ = os.WriteFile(goModPath, []byte("module test\n"), 0644)
				return subDir
			},
			wantMarker: "go.mod",
			wantType:   RootTypeProject,
		},
		{
			name: "finds package.json for Node.js project",
			setup: func(tempDir string) string {
				pkgPath := filepath.Join(tempDir, "package.json")
				_ = os.WriteFile(pkgPath, []byte(`{"name": "test"}`), 0644)
				return tempDir
			},
			wantMarker: "package.json",
			wantType:   RootTypeProject,
		},
		{
			name: "finds Cargo.toml for Rust project",
			setup: func(tempDir string) string {
				cargoPath := filepath.Join(tempDir, "Cargo.toml")
				_ = os.WriteFile(cargoPath, []byte(`[package]\nname = "test"`), 0644)
				return tempDir
			},
			wantMarker: "Cargo.toml",
			wantType:   RootTypeProject,
		},
		{
			name: "finds pyproject.toml for Python project",
			setup: func(tempDir string) string {
				pyPath := filepath.Join(tempDir, "pyproject.toml")
				_ = os.WriteFile(pyPath, []byte(`[project]\nname = "test"`), 0644)
				return tempDir
			},
			wantMarker: "pyproject.toml",
			wantType:   RootTypeProject,
		},
		{
			name: "finds pom.xml for Java/Maven project",
			setup: func(tempDir string) string {
				pomPath := filepath.Join(tempDir, "pom.xml")
				_ = os.WriteFile(pomPath, []byte(`<project></project>`), 0644)
				return tempDir
			},
			wantMarker: "pom.xml",
			wantType:   RootTypeProject,
		},
		{
			name: "finds build.gradle for Java/Gradle project",
			setup: func(tempDir string) string {
				gradlePath := filepath.Join(tempDir, "build.gradle")
				_ = os.WriteFile(gradlePath, []byte(`plugins { id("java") }`), 0644)
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
				_ = os.WriteFile(goWorkPath, []byte(`go 1.21`), 0644)
				// Create project marker in subdirectory
				subDir := filepath.Join(tempDir, "service", "auth")
				_ = os.MkdirAll(subDir, 0755)
				goModPath := filepath.Join(subDir, "go.mod")
				_ = os.WriteFile(goModPath, []byte("module test\n"), 0644)
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
			defer func() { _ = os.Chdir(originalDir) }()
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
		_ = os.WriteFile(goModPath, []byte("module test\n"), 0644)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tempDir)

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
		_ = os.WriteFile(goWorkPath, []byte("go 1.21\n"), 0644)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tempDir)

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
		_ = os.MkdirAll(subDir, 0755)

		goModPath := filepath.Join(tempDir, "go.mod")
		_ = os.WriteFile(goModPath, []byte("module test\n"), 0644)

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
		_ = os.Setenv("CLAUDE_PROJECT_DIR", testPath)
		defer func() { _ = os.Unsetenv("CLAUDE_PROJECT_DIR") }()

		got := GetProjectRootFromEnv()
		want := filepath.Clean(testPath)
		if got != want {
			t.Errorf("GetProjectRootFromEnv() = %q, want %q", got, want)
		}
	})

	t.Run("returns PROJECT_ROOT when CLAUDE_PROJECT_DIR not set", func(t *testing.T) {
		testPath := filepath.Join("fallback", "path")
		_ = os.Unsetenv("CLAUDE_PROJECT_DIR")
		_ = os.Setenv("PROJECT_ROOT", testPath)
		defer func() { _ = os.Unsetenv("PROJECT_ROOT") }()

		got := GetProjectRootFromEnv()
		want := filepath.Clean(testPath)
		if got != want {
			t.Errorf("GetProjectRootFromEnv() = %q, want %q", got, want)
		}
	})

	t.Run("CLAUDE_PROJECT_DIR takes priority", func(t *testing.T) {
		priorityPath := "priority"
		fallbackPath := "fallback"
		_ = os.Setenv("CLAUDE_PROJECT_DIR", priorityPath)
		_ = os.Setenv("PROJECT_ROOT", fallbackPath)
		defer func() { _ = os.Unsetenv("CLAUDE_PROJECT_DIR") }()
		defer func() { _ = os.Unsetenv("PROJECT_ROOT") }()

		got := GetProjectRootFromEnv()
		want := filepath.Clean(priorityPath)
		if got != want {
			t.Errorf("GetProjectRootFromEnv() = %q, want %q", got, want)
		}
	})

	t.Run("returns empty when neither set", func(t *testing.T) {
		_ = os.Unsetenv("CLAUDE_PROJECT_DIR")
		_ = os.Unsetenv("PROJECT_ROOT")

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
		_ = os.WriteFile(goModPath, []byte("module test\n"), 0644)

		// Set env var to different path
		overridePath := filepath.Join("override", "path")
		_ = os.Setenv("CLAUDE_PROJECT_DIR", overridePath)
		defer func() { _ = os.Unsetenv("CLAUDE_PROJECT_DIR") }()

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tempDir)

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
		_ = os.WriteFile(gitFile, []byte("gitdir: /main/repo/.git/worktrees/feature\n"), 0644)

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
		_ = os.MkdirAll(gitDir, 0755)

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
		_ = os.MkdirAll(gitDir, 0755)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tempDir)

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
		_ = os.MkdirAll(subDir, 0755)
		gitDir := filepath.Join(tempDir, ".git")
		_ = os.MkdirAll(gitDir, 0755)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(subDir)

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
		_ = os.MkdirAll(forgeDir, 0755)

		// Create go.mod in backend/ subdirectory
		backendDir := filepath.Join(tempDir, "backend")
		_ = os.MkdirAll(backendDir, 0755)
		goModPath := filepath.Join(backendDir, "go.mod")
		_ = os.WriteFile(goModPath, []byte("module backend\n"), 0644)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(backendDir)

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
		_ = os.MkdirAll(forgeDir, 0755)

		frontendDir := filepath.Join(tempDir, "frontend")
		_ = os.MkdirAll(frontendDir, 0755)
		pkgJSON := filepath.Join(frontendDir, "package.json")
		_ = os.WriteFile(pkgJSON, []byte(`{"name": "frontend"}`), 0644)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(frontendDir)

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
		_ = os.MkdirAll(forgeDir, 0755)
		goModPath := filepath.Join(tempDir, "go.mod")
		_ = os.WriteFile(goModPath, []byte("module test\n"), 0644)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tempDir)

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
		_ = os.WriteFile(pnpmWorkspace, []byte("packages:\n  - 'apps/*'\n"), 0644)

		// Create package.json in subproject
		subDir := filepath.Join(tempDir, "apps", "web")
		_ = os.MkdirAll(subDir, 0755)
		pkgJSON := filepath.Join(subDir, "package.json")
		_ = os.WriteFile(pkgJSON, []byte(`{"name": "web"}`), 0644)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(subDir)

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
		_ = os.WriteFile(settingsGradle, []byte("include 'service-auth'"), 0644)

		// Create build.gradle in subproject
		subDir := filepath.Join(tempDir, "service-auth")
		_ = os.MkdirAll(subDir, 0755)
		buildGradle := filepath.Join(subDir, "build.gradle")
		_ = os.WriteFile(buildGradle, []byte("plugins { id('java') }"), 0644)

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(subDir)

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

// --- Additional tests to push coverage to 85%+ ---

func TestFindProjectRootFrom(t *testing.T) {
	t.Run("returns path from FindRootInfoFrom", func(t *testing.T) {
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "deep", "nested", "dir")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}

		goModPath := filepath.Join(tempDir, "go.mod")
		if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		root, err := FindProjectRootFrom(subDir)
		if err != nil {
			t.Fatalf("FindProjectRootFrom() error = %v", err)
		}
		if root != tempDir {
			t.Errorf("FindProjectRootFrom() = %q, want %q", root, tempDir)
		}
	})

	t.Run("returns env override path", func(t *testing.T) {
		overridePath := filepath.Join("env", "override")
		_ = os.Setenv("CLAUDE_PROJECT_DIR", overridePath)
		defer func() { _ = os.Unsetenv("CLAUDE_PROJECT_DIR") }()

		root, err := FindProjectRootFrom("/some/dir")
		if err != nil {
			t.Fatalf("FindProjectRootFrom() error = %v", err)
		}
		want := filepath.Clean(overridePath)
		if root != want {
			t.Errorf("FindProjectRootFrom() = %q, want %q", root, want)
		}
	})

	t.Run("returns project root from deep nesting", func(t *testing.T) {
		tempDir := t.TempDir()
		deepDir := filepath.Join(tempDir, "a", "b", "c", "d")
		if err := os.MkdirAll(deepDir, 0755); err != nil {
			t.Fatal(err)
		}
		goModPath := filepath.Join(tempDir, "go.mod")
		if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		root, err := FindProjectRootFrom(deepDir)
		if err != nil {
			t.Fatalf("FindProjectRootFrom() error = %v", err)
		}
		if root != tempDir {
			t.Errorf("FindProjectRootFrom() = %q, want %q", root, tempDir)
		}
	})
}

func TestFindRootInfoFrom_EnvOverride(t *testing.T) {
	t.Run("returns RootTypeUnknown when env var set", func(t *testing.T) {
		envDir := filepath.Join("custom", "env", "root")
		_ = os.Setenv("CLAUDE_PROJECT_DIR", envDir)
		defer func() { _ = os.Unsetenv("CLAUDE_PROJECT_DIR") }()

		info, err := FindRootInfoFrom("/any/path")
		if err != nil {
			t.Fatalf("FindRootInfoFrom() error = %v", err)
		}
		if info.Type != RootTypeUnknown {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeUnknown)
		}
		if info.Marker != "ENV" {
			t.Errorf("Marker = %q, want %q", info.Marker, "ENV")
		}
		if info.Path != filepath.Clean(envDir) {
			t.Errorf("Path = %q, want %q", info.Path, filepath.Clean(envDir))
		}
	})
}

func TestFindProjectRoot_EnvOverride(t *testing.T) {
	t.Run("CLAUDE_PROJECT_DIR returns env path directly", func(t *testing.T) {
		envDir := filepath.Join("from", "env")
		_ = os.Setenv("CLAUDE_PROJECT_DIR", envDir)
		defer func() { _ = os.Unsetenv("CLAUDE_PROJECT_DIR") }()

		root, err := FindProjectRoot()
		if err != nil {
			t.Fatalf("FindProjectRoot() error = %v", err)
		}
		want := filepath.Clean(envDir)
		if root != want {
			t.Errorf("FindProjectRoot() = %q, want %q", root, want)
		}
	})
}

func TestFindRootInfo_EnvOverride(t *testing.T) {
	t.Run("CLAUDE_PROJECT_DIR returns RootTypeUnknown info", func(t *testing.T) {
		envDir := filepath.Join("env", "based")
		_ = os.Setenv("CLAUDE_PROJECT_DIR", envDir)
		defer func() { _ = os.Unsetenv("CLAUDE_PROJECT_DIR") }()

		info, err := FindRootInfo()
		if err != nil {
			t.Fatalf("FindRootInfo() error = %v", err)
		}
		if info.Type != RootTypeUnknown {
			t.Errorf("Type = %v, want %v", info.Type, RootTypeUnknown)
		}
		if info.Marker != "ENV" {
			t.Errorf("Marker = %q, want %q", info.Marker, "ENV")
		}
	})
}

func TestFindVCSRootWithGitInit(t *testing.T) {
	t.Run("finds VCS root in real git repo", func(t *testing.T) {
		tempDir := t.TempDir()

		// Run git init to create a real .git directory
		cmd := exec.Command("git", "init")
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("git init failed: %v", err)
		}

		originalDir, _ := os.Getwd()
		defer func() { _ = os.Chdir(originalDir) }()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatalf("Chdir() error = %v", err)
		}

		root, err := FindVCSRoot()
		if err != nil {
			t.Fatalf("FindVCSRoot() error = %v", err)
		}
		// Resolve both to handle symlink differences on some platforms
		rootResolved, _ := filepath.EvalSymlinks(root)
		tempResolved, _ := filepath.EvalSymlinks(tempDir)
		if rootResolved != tempResolved {
			t.Errorf("FindVCSRoot() = %q, want %q", rootResolved, tempResolved)
		}
	})

	t.Run("finds VCS root from subdirectory of real git repo", func(t *testing.T) {
		tempDir := t.TempDir()
		cmd := exec.Command("git", "init")
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("git init failed: %v", err)
		}

		subDir := filepath.Join(tempDir, "src", "pkg")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}

		root, err := FindVCSRootFrom(subDir)
		if err != nil {
			t.Fatalf("FindVCSRootFrom() error = %v", err)
		}
		rootResolved, _ := filepath.EvalSymlinks(root)
		tempResolved, _ := filepath.EvalSymlinks(tempDir)
		if rootResolved != tempResolved {
			t.Errorf("FindVCSRootFrom() = %q, want %q", rootResolved, tempResolved)
		}
	})
}

func TestFindVCSRootFrom_Errors(t *testing.T) {
	t.Run("error when no VCS markers found", func(t *testing.T) {
		tempDir := t.TempDir()
		// No .git or .hg anywhere up the tree
		_, err := FindVCSRootFrom(tempDir)
		if err == nil {
			t.Error("FindVCSRootFrom() expected error, got nil")
		}
	})
}

func TestFindVCSRootFrom_Hg(t *testing.T) {
	t.Run("finds .hg directory", func(t *testing.T) {
		tempDir := t.TempDir()
		hgDir := filepath.Join(tempDir, ".hg")
		if err := os.MkdirAll(hgDir, 0755); err != nil {
			t.Fatal(err)
		}

		root, err := FindVCSRootFrom(tempDir)
		if err != nil {
			t.Fatalf("FindVCSRootFrom() error = %v", err)
		}
		if root != tempDir {
			t.Errorf("FindVCSRootFrom() = %q, want %q", root, tempDir)
		}
	})
}

func TestMatchesMarker_DirectoryRequired(t *testing.T) {
	t.Run("directory marker rejects file", func(t *testing.T) {
		tempDir := t.TempDir()
		// .hg must be a directory, create it as a file instead
		hgFile := filepath.Join(tempDir, ".hg")
		if err := os.WriteFile(hgFile, []byte("not a directory"), 0644); err != nil {
			t.Fatal(err)
		}

		hgMarker := Marker{Name: ".hg", Type: RootTypeVCS, IsDirectory: true}
		if matchesMarker(tempDir, hgMarker) {
			t.Error("matchesMarker() should reject .hg when it is a file, not a directory")
		}
	})

	t.Run("directory marker accepts directory", func(t *testing.T) {
		tempDir := t.TempDir()
		forgeDir := filepath.Join(tempDir, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}

		forgeMarker := Marker{Name: ".forge", Type: RootTypeWorkspace, IsDirectory: true}
		if !matchesMarker(tempDir, forgeMarker) {
			t.Error("matchesMarker() should accept .forge when it is a directory")
		}
	})

	t.Run("directory marker rejects when absent", func(t *testing.T) {
		tempDir := t.TempDir()
		forgeMarker := Marker{Name: ".forge", Type: RootTypeWorkspace, IsDirectory: true}
		if matchesMarker(tempDir, forgeMarker) {
			t.Error("matchesMarker() should reject absent marker")
		}
	})
}

func TestMatchesMarker_GlobPattern(t *testing.T) {
	t.Run("glob matches build.gradle.kts", func(t *testing.T) {
		tempDir := t.TempDir()
		buildFile := filepath.Join(tempDir, "build.gradle.kts")
		if err := os.WriteFile(buildFile, []byte("plugins {}"), 0644); err != nil {
			t.Fatal(err)
		}

		globMarker := Marker{Name: "build.gradle*", Type: RootTypeProject, IsFileGlob: true}
		if !matchesMarker(tempDir, globMarker) {
			t.Error("matchesMarker() should match build.gradle.kts via glob")
		}
	})

	t.Run("glob matches build.gradle", func(t *testing.T) {
		tempDir := t.TempDir()
		buildFile := filepath.Join(tempDir, "build.gradle")
		if err := os.WriteFile(buildFile, []byte("plugins {}"), 0644); err != nil {
			t.Fatal(err)
		}

		globMarker := Marker{Name: "build.gradle*", Type: RootTypeProject, IsFileGlob: true}
		if !matchesMarker(tempDir, globMarker) {
			t.Error("matchesMarker() should match build.gradle via glob")
		}
	})

	t.Run("glob returns false when no match", func(t *testing.T) {
		tempDir := t.TempDir()
		globMarker := Marker{Name: "build.gradle*", Type: RootTypeProject, IsFileGlob: true}
		if matchesMarker(tempDir, globMarker) {
			t.Error("matchesMarker() should return false when glob has no matches")
		}
	})
}

func TestFindRootInfoFrom_VCSDetected(t *testing.T) {
	t.Run("VCS marker is detected alongside project marker", func(t *testing.T) {
		tempDir := t.TempDir()
		// Create both .git and go.mod in the same directory
		gitDir := filepath.Join(tempDir, ".git")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}
		goModPath := filepath.Join(tempDir, "go.mod")
		if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		info, err := FindRootInfoFrom(tempDir)
		if err != nil {
			t.Fatalf("FindRootInfoFrom() error = %v", err)
		}
		// Project marker (go.mod) takes priority over VCS marker (.git)
		if info.Type != RootTypeProject {
			t.Errorf("Type = %v, want %v (project should take priority over VCS)", info.Type, RootTypeProject)
		}
		if info.Marker != "go.mod" {
			t.Errorf("Marker = %q, want %q", info.Marker, "go.mod")
		}
	})
}

func TestFindRootInfoFrom_NoMarkersInTree(t *testing.T) {
	t.Run("finds marker from ancestor directory", func(t *testing.T) {
		// Create a marker in a parent, then query from a nested child
		parentDir := t.TempDir()
		_ = os.WriteFile(filepath.Join(parentDir, "go.mod"), []byte("module test\n"), 0644)
		childDir := filepath.Join(parentDir, "sub", "deep")
		_ = os.MkdirAll(childDir, 0755)

		info, err := FindRootInfoFrom(childDir)
		if err != nil {
			t.Fatalf("FindRootInfoFrom() error = %v", err)
		}
		if info.Path != parentDir {
			t.Errorf("expected %s, got %s", parentDir, info.Path)
		}
		if info.Type != RootTypeProject {
			t.Errorf("expected RootTypeProject, got %v", info.Type)
		}
	})
}
