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
	"text/template"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/task"
)

//go:embed data/*.md
var templateFS embed.FS

// promptTemplateData holds all fields exposed to prompt templates via text/template.
// All fields are value types (string); zero values are empty strings.
// In templates, {{if .Field}} evaluates to false for empty strings, requiring no nil checks.
type promptTemplateData struct {
	TaskID           string // task ID (e.g. "2.1", "T-test-gen-cases")
	TaskFile         string // absolute path to the task markdown file
	TaskCategory     string // task category (fix/cleanup/doc/test etc.), empty string omits the category paragraph
	FeatureSlug      string // feature slug (e.g. "auth-refresh")
	PhaseSummary     string // phase summary text, empty string omits the entire paragraph
	CoverageStrategy string // coverage strategy: "percentage"/"maintain"/empty
	CoverageTarget   string // coverage target value (e.g. "Achieve 80% test coverage")
	TestTypeArg      string // test type argument (e.g. " --type api"), empty for no type
	SurfaceKey       string // surface key, empty string omits surface label line
	SurfaceType      string // surface type
	Complexity       string // task complexity: low/medium/high (defaults to "medium")
}

// placeholderReplacer maps legacy {{PLACEHOLDER}} syntax to dot-notation {{.Placeholder}}
// for text/template compatibility. This bridge allows the engine to use text/template
// while template files still use the old {{X}} format (Task 2 migrates the files).
var placeholderReplacer = strings.NewReplacer(
	"{{TASK_ID}}", "{{.TaskID}}",
	"{{TASK_FILE}}", "{{.TaskFile}}",
	"{{SURFACE_KEY}}", "{{.SurfaceKey}}",
	"{{FEATURE_SLUG}}", "{{.FeatureSlug}}",
	"{{PHASE_SUMMARY}}", "{{.PhaseSummary}}",
	"{{TEST_TYPE_ARG}}", "{{.TestTypeArg}}",
	"{{COVERAGE_STRATEGY}}", "{{.CoverageStrategy}}",
	"{{COVERAGE_TARGET}}", "{{.CoverageTarget}}",
	"{{COMPLEXITY}}", "{{.Complexity}}",
)

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

	if !task.IsValidType(t.Type) {
		return "", fmt.Errorf("unknown type: %s", t.Type)
	}

	return renderTemplate(templatePath(t.Type), opts, t)
}

// renderTemplate reads the embed template and renders it using text/template.
//
// Legacy {{PLACEHOLDER}} syntax in template files is bridge-converted to {{.Placeholder}}
// dot-notation before parsing. Task 2 migrates the template files permanently.
//
// Available template fields (accessed as {{.FieldName}} after migration):
//
//	{{.TaskID}}           — task ID (e.g. "2.1", "T-test-gen-cases")
//	{{.TaskFile}}         — absolute path to the task markdown file
//	{{.TaskCategory}}     — task category for submit-task routing (fix/cleanup/doc/test etc.)
//	{{.SurfaceKey}}       — task surface key (empty string means cross-surface)
//	{{.FeatureSlug}}      — feature slug (e.g. "auth-refresh")
//	{{.PhaseSummary}}     — "PHASE_SUMMARY: <path>" or empty (injected line)
//	{{.TestTypeArg}}      — " --type <surfaceType>" or empty (for per-type gen-scripts)
//	{{.CoverageStrategy}} — coverage strategy text, empty for non-testable types
//	{{.CoverageTarget}}   — coverage target instruction text
//	{{.Complexity}}       — task complexity level ("low", "medium", or "high")
func renderTemplate(templateFile string, opts SynthesizeOpts, t task.Task) (string, error) {
	data, err := templateFS.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", templateFile, err)
	}

	taskFile := filepath.Join(opts.ProjectRoot, feature.GetTaskFile(opts.FeatureSlug, t.File))

	phaseSummaryPath := PhaseDetect(opts.ProjectRoot, opts.FeatureSlug, opts.TaskID)
	phaseSummaryLine := ""
	if phaseSummaryPath != "" {
		phaseSummaryLine = "PHASE_SUMMARY: " + phaseSummaryPath
	}

	// Build --type argument from task SurfaceType for per-type gen-scripts tasks.
	testTypeArg := ""
	if t.SurfaceType != "" {
		testTypeArg = " --type " + t.SurfaceType
	}

	// Inject coverage target for testable (coding.*) task types.
	coverageStrategy := ""
	coverageTarget := ""
	if task.IsTestableType(t.Type) {
		coverageStrategy, coverageTarget = resolveCoverage(opts.ProjectRoot, t)
	}

	// Inject complexity field (defaults to "medium" when empty for backward compatibility).
	complexity := t.Complexity
	if complexity == "" {
		complexity = "medium"
	}

	// Build template data struct. TaskCategory is set here (migrated from the old
	// strings.Replace injection that appended TASK_CATEGORY after TASK_FILE line).
	category := task.CategoryForType(t.Type)

	td := promptTemplateData{
		TaskID:           t.ID,
		TaskFile:         taskFile,
		TaskCategory:     category,
		FeatureSlug:      opts.FeatureSlug,
		PhaseSummary:     phaseSummaryLine,
		CoverageStrategy: coverageStrategy,
		CoverageTarget:   coverageTarget,
		TestTypeArg:      testTypeArg,
		SurfaceKey:       t.SurfaceKey,
		SurfaceType:      t.SurfaceType,
		Complexity:       complexity,
	}

	// Bridge-convert legacy {{PLACEHOLDER}} to dot-notation {{.Placeholder}}
	// for text/template compatibility. Task 2 permanently migrates the template files.
	converted := placeholderReplacer.Replace(string(data))

	tmpl, err := template.New(templateFile).Option("missingkey=error").Parse(converted)
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", templateFile, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, td); err != nil {
		return "", fmt.Errorf("execute template %s: %w", templateFile, err)
	}

	// Transitional: inject TASK_CATEGORY line after TASK_FILE line for submit-task routing.
	// Task 2 will add {{.TaskCategory}} to template files, making this post-processing unnecessary.
	result := buf.String()
	if td.TaskCategory != "" {
		result = strings.Replace(result, "TASK_FILE: "+td.TaskFile, "TASK_FILE: "+td.TaskFile+"\nTASK_CATEGORY: "+td.TaskCategory, 1)
	}

	result = cleanTemplateOutput(result, complexity)

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
// Delegates to task.InferType with nil surfaces. Kept for backward compatibility.
func InferType(id string) string {
	return task.InferType(id, nil)
}

// cleanTemplateOutput removes residual artifacts left when template variables
// are substituted with empty strings:
//
//  1. Lines that are only a label with an empty value (e.g. "SURFACE_KEY: " or "PROFILE: ")
//     are removed entirely.
//  2. Lines containing conditional sentences with empty backticks
//     (e.g. "If “ is non-empty, ...") are removed entirely.
//  3. Trailing whitespace on "just <cmd> " lines is stripped.
//  4. Collapsed consecutive blank lines are reduced to a single blank line.
//  5. Conditional paragraphs wrapped in <!-- IF NOT_LOW -->...<!-- END_IF --> markers
//     are removed when complexity is "low". This allows templates to include sections
//     (e.g., Step 1.5 spec-code scan) that are conditionally omitted for low-complexity tasks.
//     The marker format is: <!-- IF NOT_LOW --> before the paragraph and <!-- END_IF --> after it.
func cleanTemplateOutput(s string, complexity string) string {
	// Process conditional paragraph blocks first, before line-level cleanup.
	// This handles <!-- IF NOT_LOW -->...<!-- END_IF --> markers.
	if complexity == "low" {
		for {
			start := strings.Index(s, "<!-- IF NOT_LOW -->")
			if start == -1 {
				break
			}
			end := strings.Index(s, "<!-- END_IF -->")
			if end == -1 {
				break
			}
			// Remove the entire block from start marker to end of END_IF marker
			s = s[:start] + s[end+len("<!-- END_IF -->"):]
		}
	} else {
		// For non-low complexity, just strip the markers themselves (keep content)
		s = strings.ReplaceAll(s, "<!-- IF NOT_LOW -->", "")
		s = strings.ReplaceAll(s, "<!-- END_IF -->", "")
	}

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

// isLabelWithEmptyValue detects lines like "SURFACE_KEY:" or "PROFILE: " or "PHASE_SUMMARY:"
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
