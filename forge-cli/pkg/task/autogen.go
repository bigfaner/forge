package task

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"forge-cli/pkg/forgeconfig"
)

//go:embed data/*.md
var autogenTemplateFS embed.FS

// autogenTypeToFile maps task type constants to their embed template filenames.
// Filename convention: type name with '.' replaced by '-' (e.g., test.gen-scripts -> test-gen-scripts.md).
var autogenTypeToFile = map[string]string{
	TypeTestGenScripts:       "data/test-gen-scripts.md",
	TypeTestGenAndRun:        "data/test-gen-and-run.md",
	TypeTestRun:              "data/test-run.md",
	TypeTestVerifyRegression: "data/test-verify-regression.md",
	TypeEvalJourney:          "data/eval-journey.md",
	TypeEvalContract:         "data/eval-contract.md",
	TypeTestGenJourneys:      "data/test-gen-journeys.md",
	TypeTestGenContracts:     "data/test-gen-contracts.md",
	TypeValidationCode:       "data/validation-code.md",
	TypeValidationUx:         "data/validation-ux.md",
	TypeDocReview:            "data/doc-review.md",
	TypeDocConsolidate:       "data/doc-consolidate.md",
	TypeDocDrift:             "data/doc-drift.md",
	TypeCleanCode:            "data/code-quality-simplify.md",
}

// uiInterfaces is the set of interface types that have a visual UI
// and therefore require UX validation.
var uiInterfaces = map[string]bool{
	"tui":    true,
	"web":    true,
	"mobile": true,
}

// hasUIInterface returns true if any interface has a visual UI.
func hasUIInterface(interfaces []string) bool {
	for _, typ := range interfaces {
		if uiInterfaces[typ] {
			return true
		}
	}
	return false
}

// BodyContext carries planning-time data from BuildIndex() to template rendering.
// It is populated by BuildIndex() and consumed by renderBody() to substitute
// {{PLACEHOLDER}} tokens in embed template content.
type BodyContext struct {
	FeatureSlug        string   // feature slug from opts
	Mode               string   // "quick" or "breakdown"
	Scope              []string // in-scope items from proposal/PRD
	SuccessCriteria    []string // success criteria from proposal/PRD
	AcceptanceCriteria []string // PRD acceptance criteria (breakdown mode)
	ProjectType        string   // from .forge/config.yaml
	Interfaces         []string // test interfaces from config
}

// AutoGenTaskDef defines an auto-generated task definition.
type AutoGenTaskDef struct {
	ID              string
	Key             string // map key in index.json (e.g., "gen-scripts", "gen-scripts-api")
	Title           string
	Priority        string
	EstimatedTime   string
	Dependencies    []string
	Type            string
	Scope           string
	MainSession     bool
	Breaking        bool
	TestType        string // per-type interface (e.g., "api", "tui", "cli"); empty for non-per-type tasks
	FileName        string // .md filename (derived from key)
	StrategyKind    string // "generate", "run" or "" for generic
	StrategyContent []byte // resolved by caller from convention files
}

// GetBreakdownTestTasks returns test task definitions for breakdown mode.
// Interfaces are config-driven test types (e.g., "cli", "api"). Empty interfaces returns nil.
// auto controls which task categories are generated.
func GetBreakdownTestTasks(interfaces []string, auto forgeconfig.AutoConfig) []AutoGenTaskDef {
	if len(interfaces) == 0 {
		return nil
	}

	var tasks []AutoGenTaskDef

	// Shared tasks (gated by auto.E2eTest.Full)
	if auto.E2eTest.Full {
		// Per-type gen-journeys (first in pipeline)
		for _, typ := range interfaces {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "gen-journeys-" + typ, ID: "T-test-gen-journeys-" + typ,
				Title: fmt.Sprintf("Generate Test Journeys (%s)", typ), Priority: "P1", EstimatedTime: "20-30min",
				Type: TypeTestGenJourneys, Scope: "all", TestType: typ,
				StrategyKind: "interface",
			})
		}

		// Eval Journeys (after gen-journeys, before gen-contracts)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "eval-journey", ID: "T-eval-journey",
			Title: "Evaluate Journey Quality", Priority: "P1", EstimatedTime: "20-30min",
			Type: TypeEvalJourney, Scope: "all", MainSession: true,
		})

		// Gen Contracts (after eval-journey, before eval-contract)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "gen-contracts", ID: "T-test-gen-contracts",
			Title: "Generate Test Contracts", Priority: "P1", EstimatedTime: "30-45min",
			Type: TypeTestGenContracts, Scope: "all",
		})

		// Eval Contracts (after gen-contracts, before gen-test-scripts)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "eval-contract", ID: "T-eval-contract",
			Title: "Evaluate Contract Quality", Priority: "P1", EstimatedTime: "20-30min",
			Type: TypeEvalContract, Scope: "all", MainSession: true,
		})

		// Per-type gen-scripts (interface-only, no language loop)
		for _, typ := range interfaces {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "gen-test-scripts-" + typ, ID: "T-test-gen-scripts-" + typ,
				Title: fmt.Sprintf("Generate Test Scripts (%s)", typ), Priority: "P1", EstimatedTime: "1-2h",
				Type: TypeTestGenScripts, Scope: "all", TestType: typ,
				StrategyKind: "generate",
			})
		}

		// Single run (no language suffix)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "run-e2e-tests", ID: "T-test-run",
			Title: "Run e2e Tests", Priority: "P1", EstimatedTime: "30min-1h",
			Type: TypeTestRun, Scope: "all",
			StrategyKind: "run",
		})

		// Shared verify-regression
		tasks = append(tasks, AutoGenTaskDef{
			Key: "verify-regression", ID: "T-test-verify-regression",
			Title: "Verify Full E2E Regression", Priority: "P1", EstimatedTime: "15-30min",
			Type: TypeTestVerifyRegression, Scope: "all",
		})
	}

	// Validation tasks (gated by auto.Validation.Full)
	if auto.Validation.Full {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "validate-code", ID: "T-validate-code",
			Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationCode, Scope: "all", MainSession: false,
		})
		if hasUIInterface(interfaces) {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "validate-ux", ID: "T-validate-ux",
				Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
				Type: TypeValidationUx, Scope: "all", MainSession: true,
			})
		}
	}

	// Spec consolidation (gated by auto.ConsolidateSpecs.Full)
	if auto.ConsolidateSpecs.Full {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "consolidate-specs", ID: "T-specs-consolidate",
			Title: "Consolidate Specs", Priority: "P2", EstimatedTime: "20min",
			Type: TypeDocConsolidate, Scope: "all",
		})
	}

	// Clean code task (gated by auto.CleanCode.Full)
	if auto.CleanCode.Full {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "clean-code", ID: "T-clean-code",
			Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
			Type: TypeCleanCode, Scope: "all",
		})
	}

	// Set dependency chains
	resolveBreakdownDeps(tasks, interfaces, auto)

	return tasks
}

// GetQuickTestTasks returns test task definitions for quick mode.
// Interfaces are config-driven test types (e.g., "cli", "api"). Empty interfaces returns nil.
// auto controls which task categories are generated.
//
// Quick mode uses staged across types topology:
//
//	gen-journeys-per-type (parallel) -> gen-contracts -> gen-scripts-per-type (parallel) -> run -> verify-regression
//
// This replaces the old gen-and-run combined tasks with independent staged tasks,
// sharing the same task definitions as Breakdown mode (without eval quality gates).
func GetQuickTestTasks(interfaces []string, auto forgeconfig.AutoConfig) []AutoGenTaskDef {
	if len(interfaces) == 0 {
		return nil
	}

	var tasks []AutoGenTaskDef

	// Staged test pipeline (gated by auto.E2eTest.Quick)
	if auto.E2eTest.Quick {
		// Per-type gen-journeys (Stage 1: all parallel)
		for _, typ := range interfaces {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "gen-journeys-" + typ, ID: "T-test-gen-journeys-" + typ,
				Title: fmt.Sprintf("Generate Test Journeys (%s)", typ), Priority: "P1", EstimatedTime: "20-30min",
				Type: TypeTestGenJourneys, Scope: "all", TestType: typ,
				StrategyKind: "interface",
			})
		}

		// Gen Contracts (Stage 2: depends on all gen-journeys)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "gen-contracts", ID: "T-test-gen-contracts",
			Title: "Generate Test Contracts", Priority: "P1", EstimatedTime: "30-45min",
			Type: TypeTestGenContracts, Scope: "all",
		})

		// Per-type gen-scripts (Stage 3: all parallel, depend on gen-contracts)
		for _, typ := range interfaces {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "gen-test-scripts-" + typ, ID: "T-test-gen-scripts-" + typ,
				Title: fmt.Sprintf("Generate Test Scripts (%s)", typ), Priority: "P1", EstimatedTime: "1-2h",
				Type: TypeTestGenScripts, Scope: "all", TestType: typ,
				StrategyKind: "generate",
			})
		}

		// Single run (Stage 4: depends on all gen-scripts)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "run-e2e-tests", ID: "T-test-run",
			Title: "Run e2e Tests", Priority: "P1", EstimatedTime: "30min-1h",
			Type: TypeTestRun, Scope: "all",
			StrategyKind: "run",
		})

		// Shared verify-regression (Stage 5: depends on run)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "verify-regression", ID: "T-test-verify-regression",
			Title: "Verify Full E2E Regression", Priority: "P1", EstimatedTime: "15-30min",
			Type: TypeTestVerifyRegression, Scope: "all",
		})
	}

	// Validation tasks (gated by auto.Validation.Quick)
	if auto.Validation.Quick {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "validate-code", ID: "T-validate-code",
			Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationCode, Scope: "all", MainSession: false,
		})
		if hasUIInterface(interfaces) {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "validate-ux", ID: "T-validate-ux",
				Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
				Type: TypeValidationUx, Scope: "all", MainSession: true,
			})
		}
	}

	// Spec drift detection (gated by auto.ConsolidateSpecs.Quick)
	if auto.ConsolidateSpecs.Quick {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "quick-drift-detection", ID: "T-quick-doc-drift",
			Title: "Detect Spec Drift", Priority: "P2", EstimatedTime: "15min",
			Type: TypeDocDrift, Scope: "all",
		})
	}

	// Clean code task (gated by auto.CleanCode.Quick)
	if auto.CleanCode.Quick {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "quick-clean-code", ID: "T-clean-code",
			Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
			Type: TypeCleanCode, Scope: "all",
		})
	}

	resolveQuickDeps(tasks, interfaces, auto)

	return tasks
}

// renderBody substitutes {{PLACEHOLDER}} tokens in templateContent with BodyContext fields.
// Empty fields are handled per spec:
//   - {{MODE}} with empty Mode: omits the line containing {{MODE}}
//   - {{SCOPE}} with empty Scope: omits the section (## Scope ... next ## heading)
//   - {{INTERFACES}} with empty Interfaces: "See .forge/config.yaml"
//   - {{TEST_TYPE}} with empty TestType: omits the line containing {{TEST_TYPE}}
//   - {{ACCEPTANCE_CRITERIA}} with empty AcceptanceCriteria: "- [ ] All acceptance criteria met"
func renderBody(templateContent string, def AutoGenTaskDef, ctx BodyContext) string {
	s := templateContent

	// FEATURE_SLUG — always substituted (required field)
	s = strings.ReplaceAll(s, "{{FEATURE_SLUG}}", ctx.FeatureSlug)

	// MODE — omit line when empty
	modeLine := ctx.Mode
	if modeLine == "" {
		s = removeLineContaining(s, "{{MODE}}")
	} else {
		s = strings.ReplaceAll(s, "{{MODE}}", modeLine)
	}

	// SCOPE — omit section when empty
	if len(ctx.Scope) == 0 {
		s = removeSection(s, "Scope")
		// If no ## Scope heading was found, remove any residual placeholder
		s = strings.ReplaceAll(s, "{{SCOPE}}", "")
	} else {
		var scopeLines []string
		for _, item := range ctx.Scope {
			scopeLines = append(scopeLines, "- "+item)
		}
		s = strings.ReplaceAll(s, "{{SCOPE}}", strings.Join(scopeLines, "\n"))
	}

	// INTERFACES — default when empty
	if len(ctx.Interfaces) == 0 {
		s = strings.ReplaceAll(s, "{{INTERFACES}}", "See .forge/config.yaml")
	} else {
		var ifaceLines []string
		for _, iface := range ctx.Interfaces {
			ifaceLines = append(ifaceLines, "- "+iface)
		}
		s = strings.ReplaceAll(s, "{{INTERFACES}}", strings.Join(ifaceLines, "\n"))
	}

	// TEST_TYPE — omit line when empty
	testType := def.TestType
	if testType == "" {
		s = removeLineContaining(s, "{{TEST_TYPE}}")
	} else {
		s = strings.ReplaceAll(s, "{{TEST_TYPE}}", testType)
	}

	// ACCEPTANCE_CRITERIA — default when empty
	if len(ctx.AcceptanceCriteria) == 0 {
		s = strings.ReplaceAll(s, "{{ACCEPTANCE_CRITERIA}}", "- [ ] All acceptance criteria met")
	} else {
		var acLines []string
		for _, ac := range ctx.AcceptanceCriteria {
			acLines = append(acLines, "- [ ] "+ac)
		}
		s = strings.ReplaceAll(s, "{{ACCEPTANCE_CRITERIA}}", strings.Join(acLines, "\n"))
	}

	return s
}

// removeLineContaining removes the line that contains the target substring.
func removeLineContaining(s, target string) string {
	lines := strings.Split(s, "\n")
	var kept []string
	for _, line := range lines {
		if !strings.Contains(line, target) {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// removeSection removes a ## heading section by title,
// from the ## heading line up to (but not including) the next ## heading or end of string.
func removeSection(s, headingTitle string) string {
	lines := strings.Split(s, "\n")
	var result []string
	skip := false

	for _, line := range lines {
		if strings.HasPrefix(line, "## "+headingTitle) {
			skip = true
			continue
		}
		if skip && strings.HasPrefix(line, "## ") {
			skip = false
		}
		if !skip {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// GenerateTestTaskMD generates the .md file content for a test task.
func GenerateTestTaskMD(def AutoGenTaskDef, ctx BodyContext) ([]byte, error) {
	var buf strings.Builder

	// Frontmatter
	buf.WriteString("---\n")
	fmt.Fprintf(&buf, "id: %q\n", def.ID)
	fmt.Fprintf(&buf, "title: %q\n", def.Title)
	fmt.Fprintf(&buf, "priority: %q\n", def.Priority)
	fmt.Fprintf(&buf, "estimated_time: %q\n", def.EstimatedTime)
	fmt.Fprintf(&buf, "dependencies: %v\n", formatYAMLList(def.Dependencies))
	fmt.Fprintf(&buf, "type: %q\n", def.Type)
	fmt.Fprintf(&buf, "scope: %q\n", def.Scope)
	if def.MainSession {
		buf.WriteString("mainSession: true\n")
	}
	buf.WriteString("---\n\n")

	// Body — try embed template first, fallback to legacy behavior
	templateFile, hasTemplate := autogenTypeToFile[def.Type]
	if hasTemplate {
		data, err := autogenTemplateFS.ReadFile(templateFile)
		if err == nil {
			// Template loaded successfully — substitute placeholders and use as body
			rendered := renderBody(string(data), def, ctx)
			buf.WriteString(rendered)

			// Append TestType note if present
			if def.TestType != "" {
				fmt.Fprintf(&buf, "\nType: **%s**\n", def.TestType)
			}

			// Append StrategyContent after template content if present
			if len(def.StrategyContent) > 0 {
				buf.WriteString("\n\n")
				buf.Write(def.StrategyContent)
			}

			return []byte(buf.String()), nil
		}
		// Template file read failed — fall through to legacy behavior
	}

	// Legacy fallback body generation
	if def.StrategyKind != "" {
		if len(def.StrategyContent) > 0 {
			fmt.Fprintf(&buf, "# %s\n\n", def.Title)
			if def.TestType != "" {
				fmt.Fprintf(&buf, "Type: **%s**\n\n", def.TestType)
			}
			buf.Write(def.StrategyContent)
		} else {
			fmt.Fprintf(&buf, "# %s\n\nRead docs/conventions/testing/ for test generation strategy.", def.Title)
			if def.TestType != "" {
				fmt.Fprintf(&buf, " Type: %q.", def.TestType)
			}
			buf.WriteString("\n")
		}
	} else {
		fmt.Fprintf(&buf, "# %s\n\nExecute this test pipeline task.\n", def.Title)
	}

	return []byte(buf.String()), nil
}

// formatYAMLList formats a string slice as a YAML inline list.
func formatYAMLList(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

// resolveBreakdownDeps sets dependency chains for breakdown test tasks.
func resolveBreakdownDeps(tasks []AutoGenTaskDef, interfaces []string, auto forgeconfig.AutoConfig) {
	if !auto.E2eTest.Full && !auto.ConsolidateSpecs.Full && !auto.CleanCode.Full && !auto.Validation.Full {
		return // no tasks to wire
	}

	if auto.E2eTest.Full {
		// Pipeline: gen-journeys-per-type -> eval-journey -> gen-contracts -> eval-contract -> gen-scripts-per-type -> run -> verify-regression
		evalJourneyIdx := findTaskIndexOrPanic(tasks, "T-eval-journey")
		genContractsIdx := findTaskIndexOrPanic(tasks, "T-test-gen-contracts")
		evalContractIdx := findTaskIndexOrPanic(tasks, "T-eval-contract")
		runIdx := findTaskIndexOrPanic(tasks, "T-test-run")
		verifyIdx := findTaskIndexOrPanic(tasks, "T-test-verify-regression")

		// eval-journey depends on all gen-journeys tasks
		var genJourneysDeps []string
		for _, typ := range interfaces {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-journeys-"+typ)
			genJourneysDeps = append(genJourneysDeps, tasks[idx].ID)
		}
		tasks[evalJourneyIdx].Dependencies = genJourneysDeps

		// gen-contracts depends on eval-journey
		tasks[genContractsIdx].Dependencies = []string{tasks[evalJourneyIdx].ID}

		// eval-contract depends on gen-contracts
		tasks[evalContractIdx].Dependencies = []string{tasks[genContractsIdx].ID}

		// gen-scripts depend on eval-contract
		for _, typ := range interfaces {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-scripts-"+typ)
			tasks[idx].Dependencies = []string{tasks[evalContractIdx].ID}
		}

		// Run depends on all gen-scripts
		var genDeps []string
		for _, typ := range interfaces {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-scripts-"+typ)
			genDeps = append(genDeps, tasks[idx].ID)
		}
		tasks[runIdx].Dependencies = genDeps

		// Verify-regression depends on run
		tasks[verifyIdx].Dependencies = []string{tasks[runIdx].ID}
	}
	// T-validate-code depends on T-test-verify-regression (if e2e tasks exist)
	validateIdx := findTaskIndex(tasks, "T-validate-code")
	if validateIdx >= 0 && auto.E2eTest.Full {
		tasks[validateIdx].Dependencies = []string{"T-test-verify-regression"}
	}

	// T-specs-consolidate depends on T-test-verify-regression (if e2e tasks exist) or nothing
	if auto.ConsolidateSpecs.Full {
		specsIdx := findTaskIndex(tasks, "T-specs-consolidate")
		if specsIdx >= 0 && auto.E2eTest.Full {
			tasks[specsIdx].Dependencies = []string{"T-test-verify-regression"}
		}
	}

	// T-clean-code depends on last business task (resolved by caller via ResolveFirstTestDep)
	// The first test task depends on T-clean-code when both exist (resolved in BuildIndex)
}

// resolveQuickDeps sets dependency chains for quick test tasks using staged across types topology.
// Pipeline: gen-journeys-per-type (parallel) -> gen-contracts -> gen-scripts-per-type (parallel) -> run -> verify-regression
func resolveQuickDeps(tasks []AutoGenTaskDef, interfaces []string, auto forgeconfig.AutoConfig) {
	if !auto.E2eTest.Quick && !auto.ConsolidateSpecs.Quick && !auto.CleanCode.Quick && !auto.Validation.Quick {
		return // no tasks to wire
	}

	if auto.E2eTest.Quick {
		genContractsIdx := findTaskIndexOrPanic(tasks, "T-test-gen-contracts")
		runIdx := findTaskIndexOrPanic(tasks, "T-test-run")
		verifyIdx := findTaskIndexOrPanic(tasks, "T-test-verify-regression")

		// gen-contracts depends on all gen-journeys tasks (Stage 2)
		var genJourneysDeps []string
		for _, typ := range interfaces {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-journeys-"+typ)
			genJourneysDeps = append(genJourneysDeps, tasks[idx].ID)
		}
		tasks[genContractsIdx].Dependencies = genJourneysDeps

		// gen-scripts depend on gen-contracts (Stage 3)
		for _, typ := range interfaces {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-scripts-"+typ)
			tasks[idx].Dependencies = []string{tasks[genContractsIdx].ID}
		}

		// run depends on all gen-scripts (Stage 4)
		var genDeps []string
		for _, typ := range interfaces {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-scripts-"+typ)
			genDeps = append(genDeps, tasks[idx].ID)
		}
		tasks[runIdx].Dependencies = genDeps

		// verify-regression depends on run (Stage 5)
		tasks[verifyIdx].Dependencies = []string{tasks[runIdx].ID}
	}

	// T-validate-code depends on T-test-verify-regression (if e2e tasks exist) or nothing
	if auto.Validation.Quick {
		validateIdx := findTaskIndex(tasks, "T-validate-code")
		if validateIdx >= 0 && auto.E2eTest.Quick {
			tasks[validateIdx].Dependencies = []string{"T-test-verify-regression"}
		}
	}

	// T-quick-doc-drift depends on T-test-verify-regression (if e2e tasks exist) or nothing
	if auto.ConsolidateSpecs.Quick {
		idx := findTaskIndex(tasks, "T-quick-doc-drift")
		if idx >= 0 && auto.E2eTest.Quick {
			tasks[idx].Dependencies = []string{"T-test-verify-regression"}
		}
	}

	// T-clean-code depends on last business task (resolved by caller via ResolveFirstTestDep)
}

// findTaskIndex finds the index of the task with the given ID. Returns -1 if not found.
func findTaskIndex(tasks []AutoGenTaskDef, id string) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

// findTaskIndexByPrefix finds the index of the first task whose ID starts with the given prefix.
func findTaskIndexByPrefix(tasks []AutoGenTaskDef, prefix string) int {
	for i, t := range tasks {
		if strings.HasPrefix(t.ID, prefix) {
			return i
		}
	}
	return -1
}

// findTaskIndexOrPanic finds the index of the task with the given ID.
// Panics with a descriptive message (including all task IDs) if not found.
func findTaskIndexOrPanic(tasks []AutoGenTaskDef, id string) int {
	idx := findTaskIndex(tasks, id)
	if idx < 0 {
		var allIDs []string
		for _, t := range tasks {
			allIDs = append(allIDs, t.ID)
		}
		panic(fmt.Sprintf("findTaskIndex: task %q not found in tasks %v", id, allIDs))
	}
	return idx
}

// findTaskIndexByPrefixOrPanic finds the index of the first task whose ID starts with the given prefix.
// Panics with a descriptive message (including all task IDs) if not found.
func findTaskIndexByPrefixOrPanic(tasks []AutoGenTaskDef, prefix string) int {
	idx := findTaskIndexByPrefix(tasks, prefix)
	if idx < 0 {
		var allIDs []string
		for _, t := range tasks {
			allIDs = append(allIDs, t.ID)
		}
		panic(fmt.Sprintf("findTaskIndexByPrefix: task with prefix %q not found in tasks %v", prefix, allIDs))
	}
	return idx
}

// ResolveFirstTestDep resolves the first test task's dependency.
// For breakdown: depends on the highest-phase gate, or last summary if no gate.
// For quick: depends on the max business task ID.
// When T-clean-code exists, it is inserted between business tasks and test tasks.
// Returns the updated tasks with first-test-task deps set.
func ResolveFirstTestDep(tasks []AutoGenTaskDef, existingTasks map[string]Task, mode string) {
	if len(tasks) == 0 {
		return
	}

	switch mode {
	case "breakdown":
		dep := findHighestGateOrSummary(existingTasks)
		lastBiz := findMaxBusinessTaskID(existingTasks)
		// Prefer last business task when it's in a higher phase than the highest gate.
		// Handles the case where the final phase has only 1 task (no gate generated),
		// so the test chain must depend on that task, not an earlier-phase gate.
		if lastBiz != "" && (dep == "" || phaseFromID(lastBiz) > phaseFromID(dep)) {
			dep = lastBiz
		}
		if dep == "" {
			return
		}

		cleanIdx := findTaskIndex(tasks, "T-clean-code")
		firstTestIdx := findTaskIndexByPrefixOrPanic(tasks, "T-test-gen-journeys")

		if cleanIdx >= 0 {
			tasks[cleanIdx].Dependencies = []string{dep}
			tasks[firstTestIdx].Dependencies = []string{"T-clean-code"}
		} else {
			tasks[firstTestIdx].Dependencies = []string{dep}
		}

	case "quick":
		dep := findMaxBusinessTaskID(existingTasks)
		if dep == "" {
			return
		}

		cleanIdx := findTaskIndex(tasks, "T-clean-code")
		firstTestIdx := findTaskIndexByPrefixOrPanic(tasks, "T-test-gen-journeys")

		if cleanIdx >= 0 {
			tasks[cleanIdx].Dependencies = []string{dep}
			tasks[firstTestIdx].Dependencies = []string{"T-clean-code"}
		} else {
			tasks[firstTestIdx].Dependencies = []string{dep}
		}
	}
}

// findHighestGateOrSummary finds the highest-phase gate ID, falling back to last summary.
func findHighestGateOrSummary(tasks map[string]Task) string {
	var bestID string
	bestPhase := 0

	for _, t := range tasks {
		if strings.HasSuffix(t.ID, IDSuffixGate) {
			phase := phaseFromID(t.ID)
			if phase > bestPhase {
				bestPhase = phase
				bestID = t.ID
			}
		}
	}
	if bestID != "" {
		return bestID
	}

	// Fallback to last summary
	bestPhase = 0
	for _, t := range tasks {
		if strings.HasSuffix(t.ID, IDSuffixSummary) {
			phase := phaseFromID(t.ID)
			if phase > bestPhase {
				bestPhase = phase
				bestID = t.ID
			}
		}
	}
	return bestID
}

// findMaxBusinessTaskID finds the business task with the highest numeric ID.
func findMaxBusinessTaskID(tasks map[string]Task) string {
	maxN := 0
	var bestID string
	for _, t := range tasks {
		id := t.ID
		if strings.HasPrefix(id, IDPrefixTestPipeline) || strings.HasPrefix(id, "fix-") || strings.HasPrefix(id, "disc-") {
			continue
		}
		if strings.HasSuffix(id, IDSuffixGate) || strings.HasSuffix(id, IDSuffixSummary) {
			continue
		}
		n := numericID(id)
		if n > maxN {
			maxN = n
			bestID = id
		}
	}
	return bestID
}

// phaseFromID extracts the phase number from IDs like "2.gate" or "1.summary".
func phaseFromID(id string) int {
	dot := strings.LastIndex(id, ".")
	if dot < 0 {
		return 0
	}
	n := 0
	for _, c := range id[:dot] {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return n
}

// numericID extracts the leading numeric value from an ID like "3" or "2.1".
func numericID(id string) int {
	n := 0
	for _, c := range id {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

// TaskFromFile builds a Task struct from a AutoGenTaskDef.
func (d AutoGenTaskDef) TaskFromFile() Task {
	fileName := d.Key + ".md"
	return Task{
		ID:            d.ID,
		Title:         d.Title,
		Priority:      d.Priority,
		EstimatedTime: d.EstimatedTime,
		Dependencies:  d.Dependencies,
		Status:        "pending",
		File:          fileName,
		Record:        path.Join("records", fileName),
		Breaking:      d.Breaking,
		Scope:         d.Scope,
		MainSession:   d.MainSession,
		Type:          d.Type,
	}
}

// GetReviewDocTask returns a AutoGenTaskDef for the docs-only review task (T-review-doc).
// Dependencies are resolved separately by ResolveReviewDocDep.
func GetReviewDocTask() AutoGenTaskDef {
	return AutoGenTaskDef{
		Key:           "review-doc",
		ID:            "T-review-doc",
		Title:         "Review Documentation Quality",
		Priority:      "P1",
		EstimatedTime: "30min",
		Type:          TypeDocReview,
		Scope:         "all",
	}
}

// ResolveReviewDocDep sets the dependency of T-review-doc on the last business task.
// Uses lexicographic ordering to find the maximum task ID among business tasks.
func ResolveReviewDocDep(task *AutoGenTaskDef, existingTasks map[string]Task) {
	var bestID string
	for _, t := range existingTasks {
		if isAutoGenForDep(t.ID) {
			continue
		}
		if t.ID > bestID {
			bestID = t.ID
		}
	}
	if bestID != "" {
		task.Dependencies = []string{bestID}
	}
}

// isAutoGenForDep returns true for auto-generated task IDs that should be
// excluded from dependency resolution (they are not business tasks).
func isAutoGenForDep(id string) bool {
	if isTestTaskID(id) {
		return true
	}
	if id == "T-review-doc" || id == "T-validate-code" || id == "T-validate-ux" {
		return true
	}
	if strings.HasSuffix(id, IDSuffixGate) || strings.HasSuffix(id, IDSuffixSummary) {
		return true
	}
	return false
}
