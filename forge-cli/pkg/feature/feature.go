package feature

import (
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/git"
)

// GetCurrentFeature determines the current active feature.
// Priority:
// 1. Git context (worktree name -> branch name)
// 2. Feature with tasks/process/state.json
// 3. Single feature directory without state
// 4. Otherwise, return error
func GetCurrentFeature(projectRoot string) (string, error) {
	// Priority 1: Try git context (worktree or branch)
	if feature := git.GetFeatureFromGit(projectRoot); feature != "" {
		// Validate feature exists
		featureDir := filepath.Join(projectRoot, FeaturesDir, feature)
		if _, err := os.Stat(featureDir); err == nil {
			return feature, nil
		}
		// Feature doesn't exist but we inferred it from git
		// Create the feature directory structure for it
		if err := EnsureFeatureDir(projectRoot, feature); err == nil {
			return feature, nil
		}
	}

	// Priority 2-3: Fall back to feature directory scanning
	return getFeatureFromFeaturesDir(projectRoot)
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

// SetFeature creates the process directory for a feature.
//
// Deprecated: Use EnsureFeatureDir instead.
func SetFeature(projectRoot, featureSlug string) error {
	return EnsureFeatureDir(projectRoot, featureSlug)
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
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
