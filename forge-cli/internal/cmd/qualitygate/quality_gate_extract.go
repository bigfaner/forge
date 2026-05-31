package qualitygate

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// sourceFileRe matches source file paths followed by :line or :line:col patterns.
var sourceFileRe = regexp.MustCompile(`([\w][\w./-]*\.\w{1,10})(?::\d+){1,2}`)

// sourceExts is a whitelist of source code extensions for file extraction.
var sourceExts = map[string]bool{
	".go": true, ".ts": true, ".js": true, ".tsx": true, ".jsx": true,
	".py": true, ".rs": true, ".java": true, ".rb": true,
	".c": true, ".cpp": true, ".h": true, ".hpp": true,
	".css": true, ".scss": true, ".html": true, ".sql": true,
	".vue": true, ".svelte": true,
}

// ExtractSourceFiles parses error output and returns comma-separated file paths.
func ExtractSourceFiles(output string) string {
	seen := make(map[string]bool)
	var files []string
	for _, match := range sourceFileRe.FindAllStringSubmatch(output, -1) {
		path := strings.TrimPrefix(match[1], "./")
		if path == "" || seen[path] {
			continue
		}
		if !sourceExts[filepath.Ext(path)] {
			continue
		}
		seen[path] = true
		files = append(files, path)
	}

	if len(files) > maxSourceFiles {
		files = files[:maxSourceFiles]
	}
	if len(files) == 0 {
		return "See error output for affected files"
	}
	return strings.Join(files, ", ")
}

// isTestFile checks if a filename matches Go test file naming convention (*_test.go).
func isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

// extractFileLineMap extracts a mapping from test file paths to their related output lines.
// Only files with direct --- FAIL: entries (primary files) get independent map entries.
// Stack-trace-only references are excluded from independent entries but their lines
// appear in the primary file's context via line matching.
func extractFileLineMap(output string) map[string][]string {
	result := make(map[string][]string)
	if output == "" {
		return result
	}

	allLines := strings.Split(output, "\n")

	// Step 1: Collect primary files (first test file reference in each --- FAIL: block).
	primaryFiles := make(map[string]bool)
	inFailBlock := false
	foundPrimary := false
	for _, line := range allLines {
		if strings.HasPrefix(line, "--- FAIL:") {
			inFailBlock = true
			foundPrimary = false
			continue
		}
		if inFailBlock {
			if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
				if !foundPrimary {
					for _, match := range sourceFileRe.FindAllStringSubmatch(line, -1) {
						path := strings.TrimPrefix(match[1], "./")
						if isTestFile(path) {
							primaryFiles[path] = true
							foundPrimary = true
							break
						}
					}
				}
			} else {
				inFailBlock = false
				foundPrimary = false
			}
		}
	}

	if len(primaryFiles) == 0 {
		return result
	}

	// Step 2: For each line, check if it contains any primary file reference.
	// Add context window (+-2 lines) for matched lines.
	lineSet := make(map[string]map[int]bool)
	for file := range primaryFiles {
		lineSet[file] = make(map[int]bool)
	}

	for i, line := range allLines {
		for _, match := range sourceFileRe.FindAllStringSubmatch(line, -1) {
			path := strings.TrimPrefix(match[1], "./")
			if primaryFiles[path] {
				for j := i - 2; j <= i+2; j++ {
					if j >= 0 && j < len(allLines) {
						lineSet[path][j] = true
					}
				}
			}
		}
	}

	// Step 3: Convert to sorted lines per file.
	for file, indices := range lineSet {
		if len(indices) == 0 {
			continue
		}
		sorted := make([]int, 0, len(indices))
		for idx := range indices {
			sorted = append(sorted, idx)
		}
		sort.Ints(sorted)
		lines := make([]string, 0, len(sorted))
		for _, idx := range sorted {
			lines = append(lines, allLines[idx])
		}
		result[file] = lines
	}

	return result
}
