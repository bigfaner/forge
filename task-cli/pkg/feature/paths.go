package feature

import (
	"path/filepath"
)

// GetFeatureDir returns the base directory for a feature.
func GetFeatureDir(feature string) string {
	return filepath.Join(FeaturesDir, feature)
}

// GetFeatureIndexFile returns the path to the feature's index.json.
func GetFeatureIndexFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName, IndexFileName)
}

// GetFeaturePRDFile returns the path to the feature's prd.md.
func GetFeaturePRDFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, PRDFileName)
}

// GetFeatureDesignFile returns the path to the feature's design.md.
func GetFeatureDesignFile(feature string) string {
	return filepath.Join(FeaturesDir, feature, DesignFileName)
}

// GetFeatureTasksDir returns the path to the feature's tasks directory.
func GetFeatureTasksDir(feature string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName)
}

// GetFeatureRecordsDir returns the path to the feature's records directory.
func GetFeatureRecordsDir(feature string) string {
	return filepath.Join(FeaturesDir, feature, RecordsDirName)
}

// GetTaskFile returns the path to a specific task file.
func GetTaskFile(feature, filename string) string {
	return filepath.Join(FeaturesDir, feature, TasksDirName, filename)
}

// GetRecordFile returns the path to a specific record file (final record).
func GetRecordFile(feature, filename string) string {
	return filepath.Join(FeaturesDir, feature, RecordsDirName, filename)
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

