package feature

import (
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/git"
)

// Feature resolution sources, from highest to lowest priority.
const (
	SourceForgeState  = "state.json"   // Explicit selection via .forge/state.json
	SourceWorktree    = "worktree"     // Git worktree name
	SourceBranch      = "branch"       // Git branch name
	SourceFeaturesDir = "features-dir" // Scanning docs/features/ directories
)

// GetCurrentFeature determines the current active feature.
// Priority:
// 1. .forge/state.json (explicit selection via "forge feature set")
// 2. Git context (worktree name -> branch name)
// 3. Feature with tasks/process/state.json
// 4. Single feature directory without state
// 5. Otherwise, return error
func GetCurrentFeature(projectRoot string) (string, error) {
	slug, _, err := GetCurrentFeatureWithSource(projectRoot)
	return slug, err
}

// GetCurrentFeatureWithSource returns the current feature slug and the
// resolution source. The source is one of: "state.json", "worktree",
// "branch", "features-dir".
func GetCurrentFeatureWithSource(projectRoot string) (string, string, error) {
	// Priority 1: Try .forge/state.json (explicit feature selection)
	state := ReadForgeState(projectRoot)
	if state != nil && state.Feature != "" {
		featureDir := filepath.Join(projectRoot, FeaturesDir, state.Feature)
		if _, err := os.Stat(featureDir); err == nil {
			return state.Feature, SourceForgeState, nil
		}
		// Feature directory doesn't exist — skip to next priority (don't auto-create)
	}

	// Priority 2: Try git context (worktree or branch)
	feature := git.GetFeatureFromGit(projectRoot)
	if feature != "" {
		// Determine source: worktree vs branch
		source := SourceBranch
		if git.GetWorktreeName(projectRoot) != "" {
			source = SourceWorktree
		}

		// Validate feature exists
		featureDir := filepath.Join(projectRoot, FeaturesDir, feature)
		if _, err := os.Stat(featureDir); err == nil {
			return feature, source, nil
		}
		// Feature doesn't exist but we inferred it from git
		// Create the feature directory structure for it
		if err := EnsureFeatureDir(projectRoot, feature); err == nil {
			return feature, source, nil
		}
	}

	// Priority 3-4: Fall back to feature directory scanning
	slug, err := getFeatureFromFeaturesDir(projectRoot)
	if err != nil {
		return "", "", err
	}
	return slug, SourceFeaturesDir, nil
}

// getFeatureFromFeaturesDir scans features directories for active features.
func getFeatureFromFeaturesDir(projectRoot string) (string, error) {
	featuresDir := filepath.Join(projectRoot, FeaturesDir)
	entries, err := os.ReadDir(featuresDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no feature set. Run: task feature <slug>")
		}
		return "", fmt.Errorf("failed to read features directory: %w", err)
	}

	var featuresWithState []string
	var featuresWithoutState []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Check for tasks/process/state.json
		statePath := filepath.Join(featuresDir, entry.Name(), TasksDirName, ProcessDirName, StateFileName)
		if _, err := os.Stat(statePath); err == nil {
			featuresWithState = append(featuresWithState, entry.Name())
		} else {
			// Check if it's a valid feature (has index.json in tasks/)
			indexPath := filepath.Join(featuresDir, entry.Name(), TasksDirName, IndexFileName)
			if _, err := os.Stat(indexPath); err == nil {
				featuresWithoutState = append(featuresWithoutState, entry.Name())
			}
		}
	}

	// Priority 2: Use feature with active state
	if len(featuresWithState) == 1 {
		return featuresWithState[0], nil
	}
	if len(featuresWithState) > 1 {
		return "", fmt.Errorf("multiple active features: %v. Complete current task first", featuresWithState)
	}

	// Priority 3: Use single feature without state
	if len(featuresWithoutState) == 1 {
		return featuresWithoutState[0], nil
	}

	if len(featuresWithoutState) > 1 {
		return "", fmt.Errorf("multiple features without active task: %v. Run: task feature <slug>", featuresWithoutState)
	}

	return "", fmt.Errorf("no feature set. Run: task feature <slug>")
}

// RequireFeature returns the current feature or error if not set.
func RequireFeature(projectRoot string) (string, error) {
	return GetCurrentFeature(projectRoot)
}

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
		if err := os.MkdirAll(fullPath, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
