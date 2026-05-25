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
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/task"
)

//go:embed data/*.md
var templateFS embed.FS

// templatePath derives the embed template filename from a task type constant
// using the naming convention: "data/" + typeName with '.' replaced by '-' + ".md".
func templatePath(typeName string) string {
	return "data/" + strings.ReplaceAll(typeName, ".", "-") + ".md"
}

// ValidatePromptTemplates checks that all task types used by Synthesize()
// have a corresponding template file in the prompt embed FS, and that no two types
// map to the same filename. Types without a prompt template are skipped (they may
// exist only in the autogen FS).
// Must be called from the CLI main() startup path, NOT from an init() function.
func ValidatePromptTemplates() error {
	seen := make(map[string]string) // filename -> type name (for collision detection)

	for typeName := range task.ValidTypes {
		filename := templatePath(typeName)
		data, err := templateFS.ReadFile(filename)
		if err != nil {
			// Type has no template in prompt FS — skip (may exist in autogen FS)
			continue
		}
		if len(data) == 0 {
			return fmt.Errorf("template convention error: type %q maps to %q but file is empty", typeName, filename)
		}

		if prev, collision := seen[filename]; collision {
			return fmt.Errorf("template convention error: types %q and %q both map to %q", prev, typeName, filename)
		}
		seen[filename] = typeName
	}

	return nil
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

	if !task.ValidTypes[t.Type] {
		return "", fmt.Errorf("unknown type: %s", t.Type)
	}

	return renderTemplate(templatePath(t.Type), opts, t)
}

// renderTemplate reads the embed template and substitutes placeholders.
//
// WARNING: Placeholder substitution uses strings.ReplaceAll with no escaping
// mechanism. Template content must not contain bare placeholder strings like
// {{TASK_ID}}, {{TASK_FILE}}, etc., or they will be silently replaced at
// runtime. This includes code examples, documentation snippets, or any text
// that coincidentally matches the {{...}} pattern. If literal {{...}} is ever
// needed in a template, an escaping mechanism must be implemented first.
//
// Available placeholders (all use {{NAME}} syntax):
//
//	{{TASK_ID}}         — task ID (e.g. "2.1", "T-test-gen-cases")
//	{{TASK_FILE}}       — absolute path to the task markdown file
//	{{SCOPE}}           — task scope (empty string when "all" or unspecified)
//	{{FEATURE_SLUG}}    — feature slug (e.g. "auth-refresh")
//	{{PHASE_SUMMARY}}   — "PHASE_SUMMARY: <path>" or empty (injected line)
//	{{TEST_TYPE_ARG}}   — " --type <capability>" or empty (for per-type gen-scripts)
func renderTemplate(templateFile string, opts SynthesizeOpts, t task.Task) (string, error) {
	data, err := templateFS.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", templateFile, err)
	}

	taskFile := filepath.Join(opts.ProjectRoot, feature.GetTaskFile(opts.FeatureSlug, t.File))

	scope := "" // TODO: resolve scope from SurfaceKey (task 1.2b)

	phaseSummaryPath := PhaseDetect(opts.ProjectRoot, opts.FeatureSlug, opts.TaskID)

	phaseSummaryLine := ""
	if phaseSummaryPath != "" {
		phaseSummaryLine = "PHASE_SUMMARY: " + phaseSummaryPath
	}

	result := string(data)
	result = strings.ReplaceAll(result, "{{TASK_ID}}", t.ID)
	result = strings.ReplaceAll(result, "{{TASK_FILE}}", taskFile)

	// Inject TASK_CATEGORY after TASK_FILE line for submit-task skill routing.
	category := task.CategoryForType(t.Type)
	result = strings.Replace(result, "TASK_FILE: "+taskFile, "TASK_FILE: "+taskFile+"\nTASK_CATEGORY: "+category, 1)
	result = strings.ReplaceAll(result, "{{SCOPE}}", scope)
	result = strings.ReplaceAll(result, "{{FEATURE_SLUG}}", opts.FeatureSlug)
	result = strings.ReplaceAll(result, "{{PHASE_SUMMARY}}", phaseSummaryLine)

	// Extract type suffix from task ID for per-type gen-scripts tasks.
	testTypeArg := extractTestTypeArg(t.ID)
	result = strings.ReplaceAll(result, "{{TEST_TYPE_ARG}}", testTypeArg)

	// Inject coverage target for testable (coding.*) task types.
	// Non-testable types get empty strings; cleanTemplateOutput removes the empty labels.
	coverageStrategy := ""
	coverageTarget := ""
	if task.IsTestableType(t.Type) {
		coverageStrategy, coverageTarget = resolveCoverage(opts.ProjectRoot, t)
	}
	result = strings.ReplaceAll(result, "{{COVERAGE_STRATEGY}}", coverageStrategy)
	result = strings.ReplaceAll(result, "{{COVERAGE_TARGET}}", coverageTarget)

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
		if !task.IsBusinessTask(t.ID) {
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
// "2.1" → 2, "1.gate" → 1, "T-test-gen-cases" → -1 (non-integer prefix).
func phaseOf(id string) int {
	parts := strings.SplitN(id, ".", 2)
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1
	}
	return n
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
	before, after, found := strings.Cut(line, ":")
	if !found {
		return false
	}
	before = strings.TrimSpace(before)
	after = strings.TrimSpace(after)
	if before == "" || strings.Contains(before, " ") {
		return false
	}
	return after == ""
}

// genScriptBases lists the task ID bases that support per-type gen-scripts.
//
// Each base corresponds to a specific task ID format in the index:
//   - "T-test-gen-scripts"     → tasks like "T-test-gen-scripts-api", "T-test-gen-scripts-ui"
//
// The type suffix (the part after the base) determines the --type argument passed
// to the gen-script command at runtime. Adding a new base here requires that the
// corresponding task ID format is also recognized by task.ExtractTypeSuffix.
var genScriptBases = []string{
	"T-test-gen-scripts",
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

// resolveCoverage determines the effective coverage strategy and target text
// for a given task. Priority: task frontmatter coverage > config per-type > built-in default.
// Returns (strategy, targetText) where strategy is "percentage" or "maintain",
// and targetText is the human-readable instruction for the agent.
//
// For coding.cleanup and coding.refactor types, percentage strategies are overridden
// to "maintain" because their templates prescribe "No new tests" — a percentage target
// would contradict that directive.
func resolveCoverage(projectRoot string, t task.Task) (string, string) {
	// coding.cleanup and coding.refactor always use maintain strategy
	// because their templates say "No new tests" — a percentage target contradicts that.
	if t.Type == task.TypeCodingCleanup || t.Type == task.TypeCodingRefactor {
		return "maintain", "Maintain existing coverage, no more than 2% decrease"
	}

	// Priority 1: task frontmatter coverage field overrides everything.
	if t.Coverage != nil {
		return "percentage", fmt.Sprintf("Achieve %d%% test coverage", *t.Coverage)
	}

	// Priority 2: config per-type, falling back to built-in defaults.
	coverageConfig, _ := forgeconfig.ReadCoverageConfig(projectRoot)
	strategy, ok := coverageConfig.ByType[t.Type]
	if !ok {
		// Unknown type: no coverage instruction
		return "", ""
	}

	switch strategy.Type {
	case "percentage":
		if strategy.Percentage != nil {
			return "percentage", fmt.Sprintf("Achieve %d%% test coverage", *strategy.Percentage)
		}
		return "percentage", ""
	case "maintain":
		return "maintain", "Maintain existing coverage, no more than 2% decrease"
	default:
		return "", ""
	}
}

// resolveScope determines the effective scope for template rendering.
// When the project is single-scope (e.g. project-type "backend") and the task
// scope doesn't match (e.g. scope "frontend"), the scope is cleared to prevent
// generating invalid commands like "just compile frontend" on a backend-only project.
// For multi-scope projects (fullstack, mixed) or when no project-type is configured,
// the task scope is preserved as-is.
func resolveScope(projectRoot, taskScope string) string {
	// Empty or "all" scope is always cleared.
	if taskScope == "" || taskScope == "all" {
		return ""
	}

	// Read project-type from config.
	cfg, err := forgeconfig.ReadConfig(projectRoot)
	if err != nil || cfg == nil {
		// No config file or read error: preserve scope as-is.
		return taskScope
	}

	projectType := cfg.ProjectType
	if projectType == "" {
		// No project-type configured: preserve scope as-is.
		return taskScope
	}

	// Single-scope project types: if the task scope doesn't match the project type,
	// fall back to empty (no scope) so commands like "just compile" are generated
	// without a scope suffix that would fail.
	switch projectType {
	case "backend":
		if taskScope != "backend" {
			return ""
		}
	case "frontend":
		if taskScope != "frontend" {
			return ""
		}
	}

	return taskScope
}
