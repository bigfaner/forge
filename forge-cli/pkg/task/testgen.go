package task

import (
	"fmt"
	"path"
	"strings"

	"forge-cli/pkg/forgeconfig"
)

// TestTaskDef defines a test task to be generated.
type TestTaskDef struct {
	ID              string
	Key             string // map key in index.json (e.g., "gen-test-cases", "gen-test-scripts-go")
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
	StrategyKind    string // "generate", "run", "graduate", or "" for generic
	StrategyContent []byte // resolved by caller from convention files
}

// GetBreakdownTestTasks returns test task definitions for breakdown mode.
// With 0 or 1 language, uses no suffix. With 2+ languages, uses letter suffixes.
// Interfaces are config-driven test types (e.g., "cli", "api"). Empty interfaces returns nil.
// auto controls which task categories are generated.
func GetBreakdownTestTasks(languages []string, interfaces []string, auto forgeconfig.AutoConfig) []TestTaskDef {
	if len(interfaces) == 0 {
		return nil
	}

	suffix := profileSuffix(languages)

	var tasks []TestTaskDef

	// Shared tasks (gated by auto.E2eTest.Full)
	if auto.E2eTest.Full {
		tasks = append(tasks, TestTaskDef{
			Key: "gen-test-cases", ID: "T-test-gen-cases",
			Title: "Generate e2e Test Cases", Priority: "P1", EstimatedTime: "1-2h",
			Type: TypeTestGenCases, Scope: "all",
			StrategyKind: "generate",
		})
		tasks = append(tasks, TestTaskDef{
			Key: "eval-test-cases", ID: "T-test-eval-cases",
			Title: "Evaluate e2e Test Cases", Priority: "P1", EstimatedTime: "30min",
			Type: TypeTestEvalCases, Scope: "all", MainSession: true,
		})

		// Per-language tasks with per-type gen-scripts
		for i, lang := range languages {
			s := suffixLetter(i, suffix)
			for _, typ := range interfaces {
				tasks = append(tasks, TestTaskDef{
					Key: "gen-test-scripts-" + lang + "-" + typ, ID: "T-test-gen-scripts" + s + "-" + typ,
					Title: fmt.Sprintf("Generate Test Scripts (%s, %s)", lang, typ), Priority: "P1", EstimatedTime: "1-2h",
					Type: TypeTestGenScripts, Scope: "all", TestType: typ,
					StrategyKind: "generate",
				})
			}
			tasks = append(tasks, TestTaskDef{
				Key: "run-e2e-tests-" + lang, ID: "T-test-run" + s,
				Title: fmt.Sprintf("Run e2e Tests (%s)", lang), Priority: "P1", EstimatedTime: "30min-1h",
				Type: TypeTestRun, Scope: "all",
				StrategyKind: "run",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "graduate-tests-" + lang, ID: "T-test-graduate" + s,
				Title: fmt.Sprintf("Graduate Test Scripts (%s)", lang), Priority: "P1", EstimatedTime: "30min",
				Type: TypeTestGraduate, Scope: "all",
				StrategyKind: "graduate",
			})
		}

		// More shared tasks
		tasks = append(tasks, TestTaskDef{
			Key: "verify-regression", ID: "T-test-verify-regression",
			Title: "Verify Full E2E Regression", Priority: "P1", EstimatedTime: "15-30min",
			Type: TypeTestVerifyRegression, Scope: "all",
		})
	}

	// Validation tasks (gated by auto.Validation.Full)
	if auto.Validation.Full {
		tasks = append(tasks, TestTaskDef{
			Key: "validate-code", ID: "T-validate-code",
			Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationCode, Scope: "all", MainSession: false,
		})
		tasks = append(tasks, TestTaskDef{
			Key: "validate-ux", ID: "T-validate-ux",
			Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationUx, Scope: "all", MainSession: true,
		})
	}

	// Spec consolidation (gated by auto.ConsolidateSpecs.Full)
	if auto.ConsolidateSpecs.Full {
		tasks = append(tasks, TestTaskDef{
			Key: "consolidate-specs", ID: "T-specs-consolidate",
			Title: "Consolidate Specs", Priority: "P2", EstimatedTime: "20min",
			Type: TypeDocConsolidate, Scope: "all",
		})
	}

	// Clean code task (gated by auto.CleanCode.Full)
	if auto.CleanCode.Full {
		tasks = append(tasks, TestTaskDef{
			Key: "clean-code", ID: "T-clean-code",
			Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
			Type: TypeCleanCode, Scope: "all",
		})
	}

	// Set filenames and dependency chains
	resolveBreakdownDeps(tasks, languages, suffix, interfaces, auto)

	return tasks
}

// GetQuickTestTasks returns test task definitions for quick mode.
// Interfaces are config-driven test types (e.g., "cli", "api"). Empty interfaces returns nil.
// auto controls which task categories are generated.
func GetQuickTestTasks(languages []string, interfaces []string, auto forgeconfig.AutoConfig) []TestTaskDef {
	if len(interfaces) == 0 {
		return nil
	}

	suffix := profileSuffix(languages)

	var tasks []TestTaskDef

	// Per-profile with per-type gen-and-run (gated by auto.E2eTest.Quick)
	if auto.E2eTest.Quick {
		for i, lang := range languages {
			s := suffixLetter(i, suffix)
			tasks = append(tasks, TestTaskDef{
				Key: "quick-test-cases-" + lang, ID: "T-quick-gen-cases" + s,
				Title: fmt.Sprintf("Generate Quick Test Cases (%s)", lang), Priority: "P1", EstimatedTime: "30min-1h",
				Type: TypeTestGenCases, Scope: "all",
				StrategyKind: "generate",
			})
			for _, typ := range interfaces {
				tasks = append(tasks, TestTaskDef{
					Key: "quick-gen-and-run-" + lang + "-" + typ, ID: "T-quick-gen-and-run" + s + "-" + typ,
					Title: fmt.Sprintf("Generate and Run Quick Test Scripts (%s, %s)", lang, typ), Priority: "P1", EstimatedTime: "1-2h",
					Type: TypeTestGenAndRun, Scope: "all", TestType: typ,
					StrategyKind: "generate",
				})
			}
			tasks = append(tasks, TestTaskDef{
				Key: "quick-graduate-" + lang, ID: "T-quick-graduate" + s,
				Title: fmt.Sprintf("Graduate Quick Test Scripts (%s)", lang), Priority: "P1", EstimatedTime: "15min",
				Type: TypeTestGraduate, Scope: "all",
				StrategyKind: "graduate",
			})
		}

		// Shared
		tasks = append(tasks, TestTaskDef{
			Key: "quick-verify-regression", ID: "T-quick-verify-regression",
			Title: "Verify Quick E2E Regression", Priority: "P1", EstimatedTime: "15min",
			Type: TypeTestVerifyRegression, Scope: "all",
		})
	}

	// Validation tasks (gated by auto.Validation.Quick)
	if auto.Validation.Quick {
		tasks = append(tasks, TestTaskDef{
			Key: "validate-code", ID: "T-validate-code",
			Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationCode, Scope: "all", MainSession: false,
		})
		tasks = append(tasks, TestTaskDef{
			Key: "validate-ux", ID: "T-validate-ux",
			Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationUx, Scope: "all", MainSession: true,
		})
	}

	// Spec drift detection (gated by auto.ConsolidateSpecs.Quick)
	if auto.ConsolidateSpecs.Quick {
		tasks = append(tasks, TestTaskDef{
			Key: "quick-drift-detection", ID: "T-quick-doc-drift",
			Title: "Detect Spec Drift", Priority: "P2", EstimatedTime: "15min",
			Type: TypeDocDrift, Scope: "all",
		})
	}

	// Clean code task (gated by auto.CleanCode.Quick)
	if auto.CleanCode.Quick {
		tasks = append(tasks, TestTaskDef{
			Key: "quick-clean-code", ID: "T-clean-code",
			Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
			Type: TypeCleanCode, Scope: "all",
		})
	}

	resolveQuickDeps(tasks, languages, suffix, interfaces, auto)

	return tasks
}

// GenerateTestTaskMD generates the .md file content for a test task.
func GenerateTestTaskMD(def TestTaskDef, _ string) ([]byte, error) {
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

	// Body — Convention-driven: tasks reference docs/conventions/ for test strategy
	if def.StrategyKind != "" {
		if len(def.StrategyContent) > 0 {
			fmt.Fprintf(&buf, "# %s\n\n", def.Title)
			if def.TestType != "" {
				fmt.Fprintf(&buf, "Type: **%s**\n\n", def.TestType)
			}
			buf.Write(def.StrategyContent)
		} else {
			fmt.Fprintf(&buf, "# %s\n\nRead docs/conventions/testing-*.md for test generation strategy.", def.Title)
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
func resolveBreakdownDeps(tasks []TestTaskDef, languages []string, _ bool, interfaces []string, auto forgeconfig.AutoConfig) {
	if !auto.E2eTest.Full && !auto.ConsolidateSpecs.Full && !auto.CleanCode.Full && !auto.Validation.Full {
		return // no tasks to wire
	}

	e2eStart := 0
	e2eCount := 0

	if auto.E2eTest.Full {
		// T-test-eval-cases -> T-test-gen-cases
		if len(tasks) > 1 {
			tasks[1].Dependencies = []string{"T-test-gen-cases"}
		}

		profileStart := 2 // index 2 is first per-language task
		nTypes := len(interfaces)
		blockSize := nTypes + 2 // gen-per-type + run + graduate
		for i := range languages {
			blockStart := profileStart + i*blockSize

			// All per-type gen-scripts depend on T-test-eval-cases
			for j := 0; j < nTypes; j++ {
				tasks[blockStart+j].Dependencies = []string{"T-test-eval-cases"}
			}

			// Run depends on all per-type gen-scripts for this profile
			run := &tasks[blockStart+nTypes]
			var genDeps []string
			for j := 0; j < nTypes; j++ {
				genDeps = append(genDeps, tasks[blockStart+j].ID)
			}
			run.Dependencies = genDeps

			// Graduate depends on run
			graduate := &tasks[blockStart+nTypes+1]
			graduate.Dependencies = []string{run.ID}
		}

		// T-test-verify-regression depends on all graduate tasks
		sharedStart := profileStart + len(languages)*blockSize
		if len(tasks) > sharedStart {
			verifyReg := &tasks[sharedStart]
			var gradDeps []string
			for i := range languages {
				gradDeps = append(gradDeps, tasks[profileStart+i*blockSize+nTypes+1].ID)
			}
			verifyReg.Dependencies = gradDeps
		}
		e2eCount = sharedStart + 1 // number of e2e tasks
	}

	// T-validate-code depends on T-test-verify-regression (if e2e tasks exist)
	validateIdx := findTaskIndex(tasks, "T-validate-code")
	if validateIdx >= 0 && auto.E2eTest.Full && e2eCount > 0 {
		tasks[validateIdx].Dependencies = []string{"T-test-verify-regression"}
	}

	// T-specs-consolidate depends on T-test-verify-regression (if e2e tasks exist) or nothing
	if auto.ConsolidateSpecs.Full {
		specsIdx := findTaskIndex(tasks, "T-specs-consolidate")
		if specsIdx >= 0 && auto.E2eTest.Full && e2eCount > 0 {
			tasks[specsIdx].Dependencies = []string{"T-test-verify-regression"}
		}
	}

	// T-clean-code depends on last business task (resolved by caller via ResolveFirstTestDep)
	// The first test task depends on T-clean-code when both exist (resolved in BuildIndex)
	_ = e2eStart
}

// resolveQuickDeps sets dependency chains for quick test tasks.
func resolveQuickDeps(tasks []TestTaskDef, languages []string, _ bool, interfaces []string, auto forgeconfig.AutoConfig) {
	if !auto.E2eTest.Quick && !auto.ConsolidateSpecs.Quick && !auto.CleanCode.Quick && !auto.Validation.Quick {
		return // no tasks to wire
	}

	if auto.E2eTest.Quick {
		nTypes := len(interfaces)
		blockSize := 1 + nTypes + 1 // gen-cases + gen-per-type + graduate
		for i := range languages {
			blockStart := i * blockSize

			genCases := &tasks[blockStart]
			// genCases deps are placeholder (resolved by BuildIndex)

			// All per-type gen-and-run depend on gen-cases
			for j := 0; j < nTypes; j++ {
				tasks[blockStart+1+j].Dependencies = []string{genCases.ID}
			}

			// Graduate depends on all per-type gen-and-run for this profile
			graduate := &tasks[blockStart+1+nTypes]
			var genDeps []string
			for j := 0; j < nTypes; j++ {
				genDeps = append(genDeps, tasks[blockStart+1+j].ID)
			}
			graduate.Dependencies = genDeps
		}

		// T-quick-verify-regression depends on all graduate tasks
		sharedStart := len(languages) * blockSize
		if len(tasks) > sharedStart {
			var gradDeps []string
			for i := range languages {
				gradDeps = append(gradDeps, tasks[i*blockSize+1+nTypes].ID)
			}
			tasks[sharedStart].Dependencies = gradDeps
		}
	}

	// T-validate-code depends on T-quick-verify-regression (if e2e tasks exist) or nothing
	if auto.Validation.Quick {
		validateIdx := findTaskIndex(tasks, "T-validate-code")
		if validateIdx >= 0 && auto.E2eTest.Quick {
			tasks[validateIdx].Dependencies = []string{"T-quick-verify-regression"}
		}
	}

	// T-quick-doc-drift depends on T-quick-verify-regression (if e2e tasks exist) or nothing
	if auto.ConsolidateSpecs.Quick {
		idx := findTaskIndex(tasks, "T-quick-doc-drift")
		if idx >= 0 && auto.E2eTest.Quick {
			tasks[idx].Dependencies = []string{"T-quick-verify-regression"}
		}
	}

	// T-clean-code depends on last business task (resolved by caller via ResolveFirstTestDep)
}

// findTaskIndex finds the index of the task with the given ID. Returns -1 if not found.
func findTaskIndex(tasks []TestTaskDef, id string) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

// findTaskIndexByPrefix finds the index of the first task whose ID starts with the given prefix.
func findTaskIndexByPrefix(tasks []TestTaskDef, prefix string) int {
	for i, t := range tasks {
		if strings.HasPrefix(t.ID, prefix) {
			return i
		}
	}
	return -1
}

// profileSuffix returns true if languages need letter suffixes (2+).
func profileSuffix(languages []string) bool {
	return len(languages) > 1
}

// suffixLetter returns the letter suffix for the i-th profile.
func suffixLetter(i int, useSuffix bool) string {
	if !useSuffix {
		return ""
	}
	return string(rune('a' + i))
}

// ResolveFirstTestDep resolves the first test task's dependency.
// For breakdown: depends on the highest-phase gate, or last summary if no gate.
// For quick: depends on the max business task ID.
// When T-clean-code exists, it is inserted between business tasks and test tasks.
// Returns the updated tasks with first-test-task deps set.
func ResolveFirstTestDep(tasks []TestTaskDef, existingTasks map[string]Task, mode string) {
	if len(tasks) == 0 {
		return
	}

	switch mode {
	case "breakdown":
		dep := findHighestGateOrSummary(existingTasks)
		if dep == "" {
			dep = findMaxBusinessTaskID(existingTasks)
		}
		if dep == "" {
			return
		}

		cleanIdx := findTaskIndex(tasks, "T-clean-code")
		firstTestIdx := findTaskIndex(tasks, "T-test-gen-cases")

		if cleanIdx >= 0 {
			tasks[cleanIdx].Dependencies = []string{dep}
			if firstTestIdx >= 0 {
				tasks[firstTestIdx].Dependencies = []string{"T-clean-code"}
			}
		} else if firstTestIdx >= 0 {
			tasks[firstTestIdx].Dependencies = []string{dep}
		}

	case "quick":
		dep := findMaxBusinessTaskID(existingTasks)
		if dep == "" {
			return
		}

		cleanIdx := findTaskIndex(tasks, "T-clean-code")
		firstTestIdx := findTaskIndexByPrefix(tasks, "T-quick-gen-cases")

		if cleanIdx >= 0 {
			tasks[cleanIdx].Dependencies = []string{dep}
			if firstTestIdx >= 0 {
				tasks[firstTestIdx].Dependencies = []string{"T-clean-code"}
			}
		} else if firstTestIdx >= 0 {
			tasks[firstTestIdx].Dependencies = []string{dep}
		}
	}
}

// findHighestGateOrSummary finds the highest-phase gate ID, falling back to last summary.
func findHighestGateOrSummary(tasks map[string]Task) string {
	var bestID string
	bestPhase := 0

	for _, t := range tasks {
		if strings.HasSuffix(t.ID, ".gate") {
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
		if strings.HasSuffix(t.ID, ".summary") {
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
		// Skip test tasks, gates, summaries, fix tasks
		id := t.ID
		if strings.HasPrefix(id, "T-") || strings.HasPrefix(id, "fix-") || strings.HasPrefix(id, "disc-") {
			continue
		}
		if strings.HasSuffix(id, ".gate") || strings.HasSuffix(id, ".summary") {
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

// TaskFromFile builds a Task struct from a TestTaskDef.
func (d TestTaskDef) TaskFromFile() Task {
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

// GetDocEvalTask returns a TestTaskDef for the docs-only evaluation task (T-eval-doc).
// Dependencies are resolved separately by ResolveDocEvalDep.
func GetDocEvalTask() TestTaskDef {
	return TestTaskDef{
		Key:           "eval-doc",
		ID:            "T-eval-doc",
		Title:         "Evaluate Documentation Quality",
		Priority:      "P1",
		EstimatedTime: "30min",
		Type:          TypeDocEval,
		Scope:         "all",
	}
}

// ResolveDocEvalDep sets the dependency of T-eval-doc on the last business task.
// Uses lexicographic ordering to find the maximum task ID among business tasks.
func ResolveDocEvalDep(task *TestTaskDef, existingTasks map[string]Task) {
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
// Includes validation tasks which are auto-generated pipeline tasks.
func isAutoGenForDep(id string) bool {
	if isTestTaskID(id) {
		return true
	}
	if id == "T-eval-doc" || id == "T-validate-code" || id == "T-validate-ux" {
		return true
	}
	if strings.HasSuffix(id, ".gate") || strings.HasSuffix(id, ".summary") {
		return true
	}
	return false
}
