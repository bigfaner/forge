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

// GetFeaturePRDFile returns the path to prd/prd-spec.md.
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

// GetFeatureTestingScriptsDir returns the path to docs/features/{slug}/testing/scripts/.
func GetFeatureTestingScriptsDir(featureSlug string) string {
	return filepath.Join(FeaturesDir, featureSlug, TestingScriptsDirName)
}

// GetFeatureTestingResultsDir returns the path to docs/features/{slug}/testing/results/.
func GetFeatureTestingResultsDir(featureSlug string) string {
	return filepath.Join(FeaturesDir, featureSlug, TestingResultsDirName)
}

// GetFeatureTestCasesFile returns the path to docs/features/{slug}/testing/test-cases.md.
func GetFeatureTestCasesFile(featureSlug string) string {
	return filepath.Join(FeaturesDir, featureSlug, TestCasesFileName)
}

// GetE2EGraduatedMarker returns the path to tests/e2e/.graduated/<slug>.
func GetE2EGraduatedMarker(projectRoot, featureSlug string) string {
	return filepath.Join(projectRoot, E2EGraduatedDir, featureSlug)
}

// GetE2ETargetDir returns the path to tests/e2e/<target> (e.g. tests/e2e/ui/login).
func GetE2ETargetDir(projectRoot, target string) string {
	return filepath.Join(projectRoot, E2ETestsBaseDir, target)
}
func GetProposalDir(slug string) string {
	return filepath.Join(ProposalBaseDir, slug)
}

// GetProposalFile returns the path to a proposal.md file.
func GetProposalFile(slug string) string {
	return filepath.Join(ProposalBaseDir, slug, ProposalFileName)
}
