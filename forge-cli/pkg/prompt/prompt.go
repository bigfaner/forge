// Package prompt provides task prompt synthesis for the typed-task-dispatch system.
// It selects the correct agent prompt template based on task type and substitutes
// placeholders with runtime values.
package prompt

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

//go:embed data/*.md
var templateFS embed.FS

// typeToTemplate maps task type constants to their embed template filenames.
var typeToTemplate = map[string]string{
	task.TypeImplementation:               "data/implementation.md",
	task.TypeDocumentation:                "data/documentation.md",
	task.TypeDocEvaluation:                "data/doc-evaluation.md",
	task.TypeDocGenerationSummary:         "data/doc-generation-summary.md",
	task.TypeDocGenerationConsolidate:     "data/doc-generation-consolidate.md",
	task.TypeDocGenerationDrift:           "data/doc-generation-drift.md",
	task.TypeTestPipelineGenCases:         "data/test-pipeline-gen-cases.md",
	task.TypeTestPipelineEvalCases:        "data/test-pipeline-eval-cases.md",
	task.TypeTestPipelineGenScripts:       "data/test-pipeline-gen-scripts.md",
	task.TypeTestPipelineRun:              "data/test-pipeline-run.md",
	task.TypeTestPipelineGenAndRun:        "data/test-pipeline-gen-and-run.md",
	task.TypeTestPipelineGraduate:         "data/test-pipeline-graduate.md",
	task.TypeTestPipelineVerifyRegression: "data/test-pipeline-verify-regression.md",
	task.TypeFix:                          "data/fix.md",
	task.TypeGate:                         "data/gate.md",
}

// SynthesizeOpts holds inputs for prompt synthesis.
type SynthesizeOpts struct {
	ProjectRoot     string // absolute path to project root
	FeatureSlug     string // e.g. "auth-refresh"
	TaskID          string // e.g. "2.1"
	FixRecordMissed bool   // true → use fix-record-missed template
}

// Synthesize returns the synthesized agent prompt for the given task.
// On success: returns non-empty string, nil error.
// On failure: returns empty string, non-nil error.
func Synthesize(opts SynthesizeOpts) (string, error) {
	// Load the task index to look up the task.
	indexPath := filepath.Join(opts.ProjectRoot, feature.GetFeatureIndexFile(opts.FeatureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return "", fmt.Errorf("load index: %w", err)
	}

	t, found := index.ByID(opts.TaskID)
	if !found {
		return "", fmt.Errorf("task %q not found in index", opts.TaskID)
	}

	// fix-record-missed overrides normal type routing.
	if opts.FixRecordMissed {
		return renderTemplate("data/fix-record-missed.md", opts, t)
	}

	if t.Type == "" {
		return "", fmt.Errorf("type field missing for task %s", opts.TaskID)
	}

	templateFile, ok := typeToTemplate[t.Type]
	if !ok {
		return "", fmt.Errorf("unknown type: %s", t.Type)
	}

	return renderTemplate(templateFile, opts, t)
}

// renderTemplate reads the embed template and substitutes placeholders.
func renderTemplate(templateFile string, opts SynthesizeOpts, t task.Task) (string, error) {
	data, err := templateFS.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", templateFile, err)
	}

	taskFile := filepath.Join(opts.ProjectRoot, feature.GetTaskFile(opts.FeatureSlug, t.File))

	scope := t.Scope
	if scope == "" || scope == "all" {
		scope = ""
	}

	phaseSummaryPath := PhaseDetect(opts.ProjectRoot, opts.FeatureSlug, opts.TaskID)

	phaseSummaryLine := ""
	if phaseSummaryPath != "" {
		phaseSummaryLine = "PHASE_SUMMARY: " + phaseSummaryPath
	}

	result := string(data)
	result = strings.ReplaceAll(result, "{{TASK_ID}}", t.ID)
	result = strings.ReplaceAll(result, "{{TASK_FILE}}", taskFile)
	result = strings.ReplaceAll(result, "{{SCOPE}}", scope)
	result = strings.ReplaceAll(result, "{{FEATURE_SLUG}}", opts.FeatureSlug)
	result = strings.ReplaceAll(result, "{{PHASE_SUMMARY}}", phaseSummaryLine)
	result = strings.ReplaceAll(result, "{{PROFILE}}", t.Profile)

	// Extract type suffix from task ID for per-type gen-scripts tasks.
	testTypeArg := extractTestTypeArg(t.ID)
	result = strings.ReplaceAll(result, "{{TEST_TYPE_ARG}}", testTypeArg)

	result = cleanTemplateOutput(result)

	return result, nil
}

// PhaseDetect determines whether a phase summary should be injected for the given task.
// It returns the path to the previous phase's summary file if:
//   - currentPhase > maxCompletedPhase AND currentPhase > 1
//   - the summary file exists on disk
//
// Returns empty string (not an error) if no summary should be injected.
func PhaseDetect(projectRoot, featureSlug, taskID string) string {
	currentPhase := phaseOf(taskID)
	if currentPhase <= 1 {
		return ""
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return ""
	}

	maxCompleted := -1
	for _, t := range index.TasksMap() {
		if t.Status != feature.StatusCompleted {
			continue
		}
		if !isBusinessTask(t.ID) {
			continue
		}
		p := phaseOf(t.ID)
		if p > maxCompleted {
			maxCompleted = p
		}
	}

	if currentPhase <= maxCompleted {
		return ""
	}

	// currentPhase > maxCompleted AND currentPhase > 1: inject previous phase summary.
	summaryFile := fmt.Sprintf("%d-summary.md", currentPhase-1)
	summaryPath := filepath.Join(
		projectRoot,
		feature.GetFeatureRecordsDir(featureSlug),
		summaryFile,
	)
	if _, err := os.Stat(summaryPath); err != nil {
		return ""
	}

	// Return a project-relative path for portability.
	rel, err := filepath.Rel(projectRoot, summaryPath)
	if err != nil {
		return summaryPath
	}
	return rel
}

// phaseOf extracts the integer phase number from a task ID.
// "2.1" → 2, "1.gate" → 1, "T-test-1" → -1 (non-integer prefix).
func phaseOf(id string) int {
	parts := strings.SplitN(id, ".", 2)
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1
	}
	return n
}

// isBusinessTask returns true if the task ID is a regular business task
// (not a gate, summary, or T- prefixed test pipeline task).
func isBusinessTask(id string) bool {
	if strings.HasPrefix(id, "T-") {
		return false
	}
	if strings.HasSuffix(id, ".gate") {
		return false
	}
	if strings.HasSuffix(id, ".summary") {
		return false
	}
	return true
}

// InferType infers the task type from the task ID.
// Delegates to task.InferType. Kept for backward compatibility.
func InferType(id string) string {
	return task.InferType(id)
}

// cleanTemplateOutput removes residual artifacts left when template variables
// are substituted with empty strings:
//
//  1. Lines that are only a label with an empty value (e.g. "SCOPE: " or "PROFILE: ")
//     are removed entirely.
//  2. Lines containing conditional sentences with empty backticks
//     (e.g. "If “ is non-empty, ...") are removed entirely.
//  3. Trailing whitespace on "just <cmd> " lines is stripped.
//  4. Collapsed consecutive blank lines are reduced to a single blank line.
func cleanTemplateOutput(s string) string {
	lines := strings.Split(s, "\n")
	var out []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Remove conditional sentences referencing empty backticks.
		if strings.Contains(trimmed, "If `` is non-empty") {
			continue
		}

		// Remove label-only lines with empty values: "KEY:" or "KEY: " (no value after colon).
		if isLabelWithEmptyValue(trimmed) {
			continue
		}

		// Strip trailing whitespace on "just" command lines.
		if strings.HasPrefix(trimmed, "just ") && strings.HasSuffix(line, " ") {
			line = strings.TrimRight(line, " \t")
		}

		out = append(out, line)
	}

	// Collapse consecutive blank lines (3+ newlines → 2 newlines).
	result := strings.Join(out, "\n")
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	return result
}

// isLabelWithEmptyValue detects lines like "SCOPE:" or "PROFILE: " or "PHASE_SUMMARY:"
// where the label is followed by a colon and optional whitespace but no actual value.
func isLabelWithEmptyValue(line string) bool {
	if line == "" {
		return false
	}
	colonIdx := strings.Index(line, ":")
	if colonIdx < 0 {
		return false
	}
	// Everything after the colon must be empty or only whitespace.
	after := strings.TrimSpace(line[colonIdx+1:])
	// The part before the colon must look like a word (no spaces).
	before := strings.TrimSpace(line[:colonIdx])
	if before == "" || strings.Contains(before, " ") {
		return false
	}
	return after == ""
}

// genScriptBases lists the task ID bases that support per-type gen-scripts or gen-and-run.
var genScriptBases = []string{
	"T-test-2",
	"T-quick-2",
}

// extractTestTypeArg extracts the --type argument from a type-suffixed task ID.
// Returns ` --type <capability>` if a type suffix is found, or empty string otherwise.
func extractTestTypeArg(id string) string {
	for _, base := range genScriptBases {
		suffix := task.ExtractTypeSuffix(id, base)
		if suffix != "" {
			return " --type " + suffix
		}
	}
	return ""
}
