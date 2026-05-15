package task

import (
	"fmt"
	"path"
	"strings"
)

// TestTaskDef defines a test task to be generated.
type TestTaskDef struct {
	ID              string
	Key             string // map key in index.json (e.g., "gen-test-cases", "gen-test-scripts-go-test")
	Title           string
	Priority        string
	EstimatedTime   string
	Dependencies    []string
	Type            string
	Scope           string
	NoTest          bool
	MainSession     bool
	Breaking        bool
	ProfileName     string // empty for shared tasks
	TestType        string // per-type capability (e.g., "api", "tui", "cli"); empty for non-per-type tasks
	FileName        string // .md filename (derived from key)
	StrategyKind    string // "generate", "run", "graduate", or "" for generic
	StrategyContent []byte // resolved by caller from profile package
}

// GetBreakdownTestTasks returns test task definitions for breakdown mode.
// With 0 or 1 profile, uses no suffix. With 2+ profiles, uses letter suffixes.
// When detectedTypes is non-empty, creates per-type gen-scripts tasks instead of a single T-test-2 per profile.
func GetBreakdownTestTasks(profiles []string, detectedTypes []string) []TestTaskDef {
	suffix := profileSuffix(profiles)

	var tasks []TestTaskDef

	// Shared tasks
	tasks = append(tasks, TestTaskDef{
		Key: "gen-test-cases", ID: "T-test-1",
		Title: "Generate e2e Test Cases", Priority: "P1", EstimatedTime: "1-2h",
		Type: TypeTestPipelineGenCases, Scope: "all", NoTest: true,
		StrategyKind: "generate",
	})
	tasks = append(tasks, TestTaskDef{
		Key: "eval-test-cases", ID: "T-test-1b",
		Title: "Evaluate e2e Test Cases", Priority: "P1", EstimatedTime: "30min",
		Type: TypeTestPipelineEvalCases, Scope: "all", NoTest: true, MainSession: true,
	})

	if len(detectedTypes) > 0 {
		// Per-profile tasks with per-type gen-scripts
		for i, p := range profiles {
			s := suffixLetter(i, suffix)
			for _, typ := range detectedTypes {
				tasks = append(tasks, TestTaskDef{
					Key: "gen-test-scripts-" + p + "-" + typ, ID: "T-test-2" + s + "-" + typ,
					Title: fmt.Sprintf("Generate Test Scripts (%s, %s)", p, typ), Priority: "P1", EstimatedTime: "1-2h",
					Type: TypeTestPipelineGenScripts, Scope: "all", ProfileName: p, TestType: typ,
					StrategyKind: "generate",
				})
			}
			tasks = append(tasks, TestTaskDef{
				Key: "run-e2e-tests-" + p, ID: "T-test-3" + s,
				Title: fmt.Sprintf("Run e2e Tests (%s)", p), Priority: "P1", EstimatedTime: "30min-1h",
				Type: TypeTestPipelineRun, Scope: "all", ProfileName: p,
				StrategyKind: "run",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "graduate-tests-" + p, ID: "T-test-4" + s,
				Title: fmt.Sprintf("Graduate Test Scripts (%s)", p), Priority: "P1", EstimatedTime: "30min",
				Type: TypeTestPipelineGraduate, Scope: "all", ProfileName: p,
				StrategyKind: "graduate",
			})
		}
	} else {
		// Per-profile tasks: gen-scripts, run, graduate (legacy single gen-scripts)
		for i, p := range profiles {
			s := suffixLetter(i, suffix)
			tasks = append(tasks, TestTaskDef{
				Key: "gen-test-scripts-" + p, ID: "T-test-2" + s,
				Title: fmt.Sprintf("Generate Test Scripts (%s)", p), Priority: "P1", EstimatedTime: "1-2h",
				Type: TypeTestPipelineGenScripts, Scope: "all", ProfileName: p,
				StrategyKind: "generate",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "run-e2e-tests-" + p, ID: "T-test-3" + s,
				Title: fmt.Sprintf("Run e2e Tests (%s)", p), Priority: "P1", EstimatedTime: "30min-1h",
				Type: TypeTestPipelineRun, Scope: "all", ProfileName: p,
				StrategyKind: "run",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "graduate-tests-" + p, ID: "T-test-4" + s,
				Title: fmt.Sprintf("Graduate Test Scripts (%s)", p), Priority: "P1", EstimatedTime: "30min",
				Type: TypeTestPipelineGraduate, Scope: "all", ProfileName: p,
				StrategyKind: "graduate",
			})
		}
	}

	// More shared tasks
	tasks = append(tasks, TestTaskDef{
		Key: "verify-regression", ID: "T-test-4.5",
		Title: "Verify Full E2E Regression", Priority: "P1", EstimatedTime: "15-30min",
		Type: TypeTestPipelineVerifyRegression, Scope: "all",
	})
	tasks = append(tasks, TestTaskDef{
		Key: "consolidate-specs", ID: "T-test-5",
		Title: "Consolidate Specs", Priority: "P2", EstimatedTime: "20min",
		Type: TypeDocGenerationConsolidate, Scope: "all", NoTest: true,
	})

	// Set filenames and dependency chains
	resolveBreakdownDeps(tasks, profiles, suffix, detectedTypes)

	return tasks
}

// GetQuickTestTasks returns test task definitions for quick mode.
// When detectedTypes is non-empty, creates per-type gen-scripts tasks instead of a single T-quick-2 per profile.
func GetQuickTestTasks(profiles []string, detectedTypes []string) []TestTaskDef {
	suffix := profileSuffix(profiles)

	var tasks []TestTaskDef

	if len(detectedTypes) > 0 {
		// Per-profile with per-type gen-scripts
		for i, p := range profiles {
			s := suffixLetter(i, suffix)
			tasks = append(tasks, TestTaskDef{
				Key: "quick-test-cases-" + p, ID: "T-quick-1" + s,
				Title: fmt.Sprintf("Generate Quick Test Cases (%s)", p), Priority: "P1", EstimatedTime: "30min-1h",
				Type: TypeTestPipelineGenCases, Scope: "all", NoTest: true, ProfileName: p,
				StrategyKind: "generate",
			})
			for _, typ := range detectedTypes {
				tasks = append(tasks, TestTaskDef{
					Key: "quick-gen-scripts-" + p + "-" + typ, ID: "T-quick-2" + s + "-" + typ,
					Title: fmt.Sprintf("Generate Quick Test Scripts (%s, %s)", p, typ), Priority: "P1", EstimatedTime: "30min-1h",
					Type: TypeTestPipelineGenScripts, Scope: "all", ProfileName: p, TestType: typ,
					StrategyKind: "generate",
				})
			}
			tasks = append(tasks, TestTaskDef{
				Key: "quick-run-tests-" + p, ID: "T-quick-3" + s,
				Title: fmt.Sprintf("Run Quick E2E Tests (%s)", p), Priority: "P1", EstimatedTime: "15-30min",
				Type: TypeTestPipelineRun, Scope: "all", ProfileName: p,
				StrategyKind: "run",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "quick-graduate-" + p, ID: "T-quick-4" + s,
				Title: fmt.Sprintf("Graduate Quick Test Scripts (%s)", p), Priority: "P1", EstimatedTime: "15min",
				Type: TypeTestPipelineGraduate, Scope: "all", ProfileName: p,
				StrategyKind: "graduate",
			})
		}
	} else {
		// Per-profile: gen-cases, gen-scripts, run, graduate (legacy single gen-scripts)
		for i, p := range profiles {
			s := suffixLetter(i, suffix)
			tasks = append(tasks, TestTaskDef{
				Key: "quick-test-cases-" + p, ID: "T-quick-1" + s,
				Title: fmt.Sprintf("Generate Quick Test Cases (%s)", p), Priority: "P1", EstimatedTime: "30min-1h",
				Type: TypeTestPipelineGenCases, Scope: "all", NoTest: true, ProfileName: p,
				StrategyKind: "generate",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "quick-gen-scripts-" + p, ID: "T-quick-2" + s,
				Title: fmt.Sprintf("Generate Quick Test Scripts (%s)", p), Priority: "P1", EstimatedTime: "30min-1h",
				Type: TypeTestPipelineGenScripts, Scope: "all", ProfileName: p,
				StrategyKind: "generate",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "quick-run-tests-" + p, ID: "T-quick-3" + s,
				Title: fmt.Sprintf("Run Quick E2E Tests (%s)", p), Priority: "P1", EstimatedTime: "15-30min",
				Type: TypeTestPipelineRun, Scope: "all", ProfileName: p,
				StrategyKind: "run",
			})
			tasks = append(tasks, TestTaskDef{
				Key: "quick-graduate-" + p, ID: "T-quick-4" + s,
				Title: fmt.Sprintf("Graduate Quick Test Scripts (%s)", p), Priority: "P1", EstimatedTime: "15min",
				Type: TypeTestPipelineGraduate, Scope: "all", ProfileName: p,
				StrategyKind: "graduate",
			})
		}
	}

	// Shared
	tasks = append(tasks, TestTaskDef{
		Key: "quick-verify-regression", ID: "T-quick-5",
		Title: "Verify Quick E2E Regression", Priority: "P1", EstimatedTime: "15min",
		Type: TypeTestPipelineVerifyRegression, Scope: "all",
	})

	resolveQuickDeps(tasks, profiles, suffix, detectedTypes)

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
	if def.ProfileName != "" {
		fmt.Fprintf(&buf, "profile: %q\n", def.ProfileName)
	}
	if def.NoTest {
		buf.WriteString("noTest: true\n")
	}
	if def.MainSession {
		buf.WriteString("mainSession: true\n")
	}
	buf.WriteString("---\n\n")
	// Body
	if def.ProfileName != "" && def.StrategyKind != "" {
		if len(def.StrategyContent) > 0 {
			fmt.Fprintf(&buf, "# %s\n\n", def.Title)
			fmt.Fprintf(&buf, "Profile: **%s**\n\n", def.ProfileName)
			if def.TestType != "" {
				fmt.Fprintf(&buf, "Type: **%s**\n\n", def.TestType)
			}
			buf.Write(def.StrategyContent)
		} else {
			// Fallback: generic body
			fmt.Fprintf(&buf, "# %s\n\nCall the appropriate skill for profile %q", def.Title, def.ProfileName)
			if def.TestType != "" {
				fmt.Fprintf(&buf, " with type %q", def.TestType)
			}
			buf.WriteString(".\n")
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
func resolveBreakdownDeps(tasks []TestTaskDef, profiles []string, _ bool, detectedTypes []string) {
	// T-test-1 depends on last gate or last summary (placeholder, caller resolves)
	// T-test-1b depends on T-test-1
	// Per-profile: T-test-2<L>-<type> depends on T-test-1b
	//              T-test-3<L> depends on ALL T-test-2<L>-<type> for its profile
	//              T-test-4<L> depends on T-test-3<L>
	// T-test-4.5 depends on all T-test-4<L> (or T-test-4 if single)
	// T-test-5 depends on T-test-4.5

	// T-test-1b -> T-test-1
	if len(tasks) > 1 {
		tasks[1].Dependencies = []string{"T-test-1"}
	}

	profileStart := 2 // index 2 is first per-profile task
	if len(detectedTypes) > 0 {
		// Per-type mode: per-profile block is N gen-tasks + run + graduate
		nTypes := len(detectedTypes)
		blockSize := nTypes + 2 // gen-per-type + run + graduate
		for i := range profiles {
			blockStart := profileStart + i*blockSize

			// All per-type gen-scripts depend on T-test-1b
			for j := 0; j < nTypes; j++ {
				tasks[blockStart+j].Dependencies = []string{"T-test-1b"}
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

		// T-test-4.5 depends on all graduate tasks
		sharedStart := profileStart + len(profiles)*blockSize
		if len(tasks) > sharedStart {
			verifyReg := &tasks[sharedStart]
			var gradDeps []string
			for i := range profiles {
				gradDeps = append(gradDeps, tasks[profileStart+i*blockSize+nTypes+1].ID)
			}
			verifyReg.Dependencies = gradDeps
		}

		// T-test-5 depends on T-test-4.5
		if len(tasks) > sharedStart+1 {
			tasks[sharedStart+1].Dependencies = []string{"T-test-4.5"}
		}
	} else {
		// Legacy mode: per-profile block is 3 tasks (gen-scripts, run, graduate)
		for i := range profiles {
			genScripts := &tasks[profileStart+i*3]
			run := &tasks[profileStart+i*3+1]
			graduate := &tasks[profileStart+i*3+2]

			genScripts.Dependencies = []string{"T-test-1b"}
			run.Dependencies = []string{genScripts.ID}
			graduate.Dependencies = []string{run.ID}
		}

		// T-test-4.5 depends on all graduate tasks
		if len(tasks) > profileStart+len(profiles)*3 {
			verifyReg := &tasks[profileStart+len(profiles)*3]
			var gradDeps []string
			for i := range profiles {
				gradDeps = append(gradDeps, tasks[profileStart+i*3+2].ID)
			}
			if len(gradDeps) == 0 {
				gradDeps = []string{"T-test-4"}
			}
			verifyReg.Dependencies = gradDeps
		}

		// T-test-5 depends on T-test-4.5
		if len(tasks) > profileStart+len(profiles)*3+1 {
			tasks[profileStart+len(profiles)*3+1].Dependencies = []string{"T-test-4.5"}
		}
	}
}

// resolveQuickDeps sets dependency chains for quick test tasks.
func resolveQuickDeps(tasks []TestTaskDef, profiles []string, _ bool, detectedTypes []string) {
	// Per-profile: T-quick-1<L> depends on last business task (placeholder)
	// Per-type mode:   T-quick-2<L>-<type> depends on T-quick-1<L>
	//                  T-quick-3<L> depends on ALL T-quick-2<L>-<type> for its profile
	//                  T-quick-4<L> depends on T-quick-3<L>
	// Legacy mode:     T-quick-2<L> depends on T-quick-1<L>
	//                  T-quick-3<L> depends on T-quick-2<L>
	//                  T-quick-4<L> depends on T-quick-3<L>
	// T-quick-5 depends on all T-quick-4<L>

	if len(detectedTypes) > 0 {
		// Per-type mode: per-profile block is N gen-tasks + run + graduate
		nTypes := len(detectedTypes)
		blockSize := nTypes + 3 // gen-cases + gen-per-type + run + graduate
		for i := range profiles {
			blockStart := i * blockSize

			genCases := &tasks[blockStart]
			// genCases deps are placeholder (resolved by BuildIndex)

			// All per-type gen-scripts depend on gen-cases
			for j := 0; j < nTypes; j++ {
				tasks[blockStart+1+j].Dependencies = []string{genCases.ID}
			}

			// Run depends on all per-type gen-scripts for this profile
			run := &tasks[blockStart+1+nTypes]
			var genDeps []string
			for j := 0; j < nTypes; j++ {
				genDeps = append(genDeps, tasks[blockStart+1+j].ID)
			}
			run.Dependencies = genDeps

			// Graduate depends on run
			graduate := &tasks[blockStart+1+nTypes+1]
			graduate.Dependencies = []string{run.ID}
		}

		// T-quick-5 depends on all graduate tasks
		sharedStart := len(profiles) * blockSize
		if len(tasks) > sharedStart {
			var gradDeps []string
			for i := range profiles {
				gradDeps = append(gradDeps, tasks[i*blockSize+1+nTypes+1].ID)
			}
			tasks[sharedStart].Dependencies = gradDeps
		}
	} else {
		// Legacy mode: per-profile block is 4 tasks (gen-cases, gen-scripts, run, graduate)
		for i := range profiles {
			genCases := &tasks[i*4]
			genScripts := &tasks[i*4+1]
			run := &tasks[i*4+2]
			graduate := &tasks[i*4+3]

			// genCases deps are placeholder (resolved by BuildIndex)
			genScripts.Dependencies = []string{genCases.ID}
			run.Dependencies = []string{genScripts.ID}
			graduate.Dependencies = []string{run.ID}
		}

		// T-quick-5 depends on all graduate tasks
		if len(tasks) > len(profiles)*4 {
			verifyIdx := len(profiles) * 4
			var gradDeps []string
			for i := range profiles {
				gradDeps = append(gradDeps, tasks[i*4+3].ID)
			}
			tasks[verifyIdx].Dependencies = gradDeps
		}
	}
}

// profileSuffix returns true if profiles need letter suffixes (2+).
func profileSuffix(profiles []string) bool {
	return len(profiles) > 1
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
// Returns the updated tasks with first-test-task deps set.
func ResolveFirstTestDep(tasks []TestTaskDef, existingTasks map[string]Task, mode string) {
	if len(tasks) == 0 {
		return
	}

	switch mode {
	case "breakdown":
		dep := findHighestGateOrSummary(existingTasks)
		if dep != "" {
			tasks[0].Dependencies = []string{dep}
		}
	case "quick":
		dep := findMaxBusinessTaskID(existingTasks)
		if dep != "" {
			// Set deps on all first-per-profile tasks (T-quick-1<L>)
			for i := 0; i < len(tasks) && strings.HasPrefix(tasks[i].ID, "T-quick-1"); i++ {
				tasks[i].Dependencies = []string{dep}
			}
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
		NoTest:        d.NoTest,
		Type:          d.Type,
		Profile:       d.ProfileName,
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
		Type:          TypeDocEvaluation,
		Scope:         "all",
		NoTest:        true,
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
func isAutoGenForDep(id string) bool {
	if isTestTaskID(id) {
		return true
	}
	if id == "T-eval-doc" {
		return true
	}
	if strings.HasSuffix(id, ".gate") || strings.HasSuffix(id, ".summary") {
		return true
	}
	return false
}
