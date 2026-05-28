package task

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"forge-cli/pkg/forgeconfig"
)

//go:embed templates/*.md
var autogenTemplateFS embed.FS

// autogenTemplatePath derives the embed template filename from a task type constant
// using the naming convention: "templates/" + typeName with '.' replaced by '-' + ".md".
// For surface-specific types (e.g. "test.gen-scripts.cli"), strips the last segment
// to find the base type template (e.g. "test.gen-scripts" -> "templates/test-gen-scripts.md").
func autogenTemplatePath(typeName string) string {
	// Try exact match first
	path := "templates/" + strings.ReplaceAll(typeName, ".", "-") + ".md"
	if _, err := autogenTemplateFS.ReadFile(path); err == nil {
		return path
	}
	// For surface-specific types, strip last segment and try base type
	if idx := strings.LastIndex(typeName, "."); idx >= 0 {
		base := typeName[:idx]
		return "templates/" + strings.ReplaceAll(base, ".", "-") + ".md"
	}
	return path
}

// ValidateAutogenTemplates checks that all task types used by GenerateTestTaskMD()
// have a corresponding template file in the autogen embed FS, that no two types
// map to the same filename, and that each template parses as valid text/template
// and executes without error against a zero-value autogenTemplateData struct.
//
// When a template contains metadata frontmatter, it is stripped before parsing
// and the variables field is cross-validated against autogenTemplateData using reflection.
//
// Types without an autogen template are skipped (they may exist only in the prompt FS).
// Must be called from the CLI main() startup path, NOT from an init() function.
func ValidateAutogenTemplates() error {
	seen := make(map[string]string) // filename -> type name (for collision detection)
	structType := reflect.TypeOf(autogenTemplateData{})

	for typeName := range ValidTypes {
		filename := autogenTemplatePath(typeName)
		data, err := autogenTemplateFS.ReadFile(filename)
		if err != nil {
			// Type has no template in autogen FS — skip (may exist in prompt FS)
			continue
		}
		if len(data) == 0 {
			return fmt.Errorf("autogen template convention error: type %q maps to %q but file is empty", typeName, filename)
		}

		if prev, collision := seen[filename]; collision {
			return fmt.Errorf("autogen template convention error: types %q and %q both map to %q", prev, typeName, filename)
		}
		seen[filename] = typeName

		// Strip metadata frontmatter before validation
		content := string(data)
		body, meta := parseAutogenMetadata(content)

		// Cross-validate metadata variables against struct fields
		if meta != nil {
			if err := validateAutogenVariables(meta, structType); err != nil {
				return fmt.Errorf("autogen template validation error: %s: %w", filename, err)
			}
		}

		// Validate template syntax and execution with missingkey=error
		tmpl, err := template.New(filename).Option("missingkey=error").Parse(body)
		if err != nil {
			return fmt.Errorf("autogen template parse error for %q: %w", filename, err)
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, autogenTemplateData{}); err != nil {
			return fmt.Errorf("autogen template execute error for %q with zero-value data: %w", filename, err)
		}
	}

	return nil
}

// uiSurfaceTypes is the set of surface types that have a visual UI
// and therefore require UX validation.
var uiSurfaceTypes = map[string]bool{
	"tui":    true,
	"web":    true,
	"mobile": true,
}

// hasUISurface returns true if any surface type has a visual UI.
func hasUISurface(types []string) bool {
	for _, typ := range types {
		if uiSurfaceTypes[typ] {
			return true
		}
	}
	return false
}

// BodyContext carries planning-time data from BuildIndex() to template rendering.
// It is populated by BuildIndex() and consumed by renderBody() to substitute
// template fields in embed template content.
type BodyContext struct {
	FeatureSlug        string            // feature slug from opts
	Mode               string            // "quick" or "breakdown"
	SuccessCriteria    []string          // success criteria from proposal/PRD
	AcceptanceCriteria []string          // PRD acceptance criteria (breakdown mode)
	ProjectType        string            // from .forge/config.yaml
	SurfaceTypes       []string          // deduplicated surface types from config
	DocTaskCriteria    map[string]string // doc task name -> raw AC markdown (key=filename without .md)
}

// autogenTemplateData is the data model for text/template rendering of autogen body templates.
// All fields are pre-formatted strings — the caller is responsible for serializing
// slice/map fields before passing them to the template engine.
type autogenTemplateData struct {
	TaskID             string // task ID
	TaskType           string // task type identifier (e.g., "test.gen-contracts")
	FeatureSlug        string // feature identifier for template references
	Mode               string // generation mode; empty string omits Mode line
	SurfaceKey         string // surface key for inline replacement and conditional sections
	SurfaceType        string // surface type; empty string omits TestType line
	SurfaceTypes       string // pre-formatted multi-surface type string (newline-separated "- type" items); empty defaults to "See .forge/config.yaml"
	AcceptanceCriteria string // pre-formatted acceptance criteria text; empty defaults to "- [ ] All acceptance criteria met"
	DocTaskCriteria    string // pre-formatted doc task criteria text; empty omits the section
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
	MainSession     bool
	Breaking        bool
	SurfaceKey      string // user-defined surface identifier
	SurfaceType     string // surface type (e.g., "api", "tui", "cli"); empty for non-per-type tasks
	FileName        string // .md filename (derived from key)
	StrategyKind    string // "generate", "run" or "" for generic
	StrategyContent []byte // resolved by caller from convention files
}

// isSingleSurface returns true when the surfaces map represents a single surface
// (scalar form with "." key, or map with exactly one entry).
func isSingleSurface(surfaces map[string]string) bool {
	if len(surfaces) == 0 {
		return false
	}
	if len(surfaces) == 1 {
		if _, ok := surfaces["."]; ok {
			return true
		}
		// Single map entry is also single-surface
		return true
	}
	return false
}

// runTestTitle generates the title for a run-test task based on surface type.
func runTestTitle(surfaceType string) string {
	return fmt.Sprintf("Run %ss", TestTypeTitle(surfaceType))
}

// GetBreakdownTestTasks returns test task definitions for breakdown mode.
// surfaces is the surfaces map from config (e.g., {".": "api"} or {"backend": "api", "frontend": "web"}).
// executionOrder is the resolved execution order of surface keys (may be nil for single-surface).
// auto controls which task categories are generated.
func GetBreakdownTestTasks(surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig) []AutoGenTaskDef {
	if len(surfaces) == 0 {
		return nil
	}

	surfaceTypes := forgeconfig.SurfaceTypes(surfaces)
	if len(surfaceTypes) == 0 {
		return nil
	}

	singleSurface := isSingleSurface(surfaces)

	var tasks []AutoGenTaskDef

	// Shared tasks (gated by auto.Test.Full)
	if auto.Test.Full {
		// Single gen-journeys task covering all configured surfaces
		tasks = append(tasks, AutoGenTaskDef{
			Key: "gen-journeys", ID: "T-test-gen-journeys",
			Title: "Generate Test Journeys", Priority: "P1", EstimatedTime: "20-30min",
			Type:         TypeTestGenJourneys,
			StrategyKind: "interface",
		})

		// Eval Journeys (after gen-journeys, before gen-contracts)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "eval-journey", ID: "T-eval-journey",
			Title: "Evaluate Journey Quality", Priority: "P1", EstimatedTime: "20-30min",
			Type: TypeEvalJourney, MainSession: true,
		})

		// Gen Contracts (after eval-journey, before eval-contract)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "gen-contracts", ID: "T-test-gen-contracts",
			Title: "Generate Test Contracts", Priority: "P1", EstimatedTime: "30-45min",
			Type: TypeTestGenContracts,
		})

		// Eval Contracts (after gen-contracts, before gen-test-scripts)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "eval-contract", ID: "T-eval-contract",
			Title: "Evaluate Contract Quality", Priority: "P1", EstimatedTime: "20-30min",
			Type: TypeEvalContract, MainSession: true,
		})

		// Per-type gen-scripts (interface-only, no language loop)
		for _, typ := range surfaceTypes {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "gen-test-scripts-" + typ, ID: "T-test-gen-scripts-" + typ,
				Title: fmt.Sprintf("Generate %s Scripts", TestTypeTitle(typ)), Priority: "P1", EstimatedTime: "1-2h",
				Type: GenSurfaceTestType(TypeTestGenScripts, typ), SurfaceType: typ,
				StrategyKind: "generate",
			})
		}

		// Per-surface-key run-test tasks (serial chain)
		if singleSurface {
			// Single surface: degenerate to no suffix T-test-run
			// Use surface type for type name (single surface: key is ".")
			singleType := surfaceTypes[0]
			tasks = append(tasks, AutoGenTaskDef{
				Key: "run-test", ID: "T-test-run",
				Title: runTestTitle(singleType), Priority: "P1", EstimatedTime: "30min-1h",
				Type:         GenSurfaceTestType(TypeTestRun, singleType),
				StrategyKind: "run",
				SurfaceType:  singleType,
			})
		} else {
			// Multi-surface: generate T-test-run-{surface-key} per surface key in execution order
			for _, key := range executionOrder {
				surfaceType := surfaces[key]
				tasks = append(tasks, AutoGenTaskDef{
					Key: "run-test-" + key, ID: "T-test-run-" + key,
					Title: runTestTitle(surfaceType), Priority: "P1", EstimatedTime: "30min-1h",
					Type:         GenSurfaceTestType(TypeTestRun, key),
					SurfaceKey:   key,
					SurfaceType:  surfaceType,
					StrategyKind: "run",
				})
			}
		}
	}

	// Validation tasks (gated by auto.Validation.Full)
	if auto.Validation.Full {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "validate-code", ID: "T-validate-code",
			Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationCode, MainSession: false,
		})
		if hasUISurface(surfaceTypes) {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "validate-ux", ID: "T-validate-ux",
				Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
				Type: TypeValidationUx, MainSession: true,
			})
		}
	}

	// Spec consolidation (gated by auto.ConsolidateSpecs.Full)
	if auto.ConsolidateSpecs.Full {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "consolidate-specs", ID: "T-specs-consolidate",
			Title: "Consolidate Specs", Priority: "P2", EstimatedTime: "20min",
			Type: TypeDocConsolidate,
		})
	}

	// Clean code task (gated by auto.CleanCode.Full)
	if auto.CleanCode.Full {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "clean-code", ID: "T-clean-code",
			Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
			Type: TypeCleanCode,
		})
	}

	// Set dependency chains
	resolveBreakdownDeps(tasks, surfaceTypes, surfaces, executionOrder, auto)

	return tasks
}

// GetQuickTestTasks returns test task definitions for quick mode.
// surfaces is the surfaces map from config (e.g., {".": "api"} or {"backend": "api", "frontend": "web"}).
// executionOrder is the resolved execution order of surface keys (may be nil for single-surface).
// auto controls which task categories are generated.
//
// Quick mode uses staged across types topology:
//
//	gen-journeys (single) -> run-test-{key1} -> run-test-{key2} -> ...
//
// This replaces the old gen-and-run combined tasks with independent staged tasks,
// sharing the same task definitions as Breakdown mode (without eval quality gates).
func GetQuickTestTasks(surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig) []AutoGenTaskDef {
	if len(surfaces) == 0 {
		return nil
	}

	surfaceTypes := forgeconfig.SurfaceTypes(surfaces)
	if len(surfaceTypes) == 0 {
		return nil
	}

	singleSurface := isSingleSurface(surfaces)

	var tasks []AutoGenTaskDef

	// Staged test pipeline (gated by auto.Test.Quick)
	// Quick mode: gen-journeys -> run-tests(serial)
	// (no gen-contracts or gen-scripts in Quick mode)
	if auto.Test.Quick {
		// Single gen-journeys task covering all configured surfaces (Stage 1)
		tasks = append(tasks, AutoGenTaskDef{
			Key: "gen-journeys", ID: "T-test-gen-journeys",
			Title: "Generate Test Journeys", Priority: "P1", EstimatedTime: "20-30min",
			Type:         TypeTestGenJourneys,
			StrategyKind: "interface",
		})

		// Per-surface-key run-test tasks (serial chain, Stage 2)
		if singleSurface {
			// Single surface: degenerate to no suffix T-test-run
			singleType := surfaceTypes[0]
			tasks = append(tasks, AutoGenTaskDef{
				Key: "run-test", ID: "T-test-run",
				Title: runTestTitle(singleType), Priority: "P1", EstimatedTime: "30min-1h",
				Type:         GenSurfaceTestType(TypeTestRun, singleType),
				StrategyKind: "run",
				SurfaceType:  singleType,
			})
		} else {
			// Multi-surface: generate T-test-run-{surface-key} per surface key in execution order
			for _, key := range executionOrder {
				surfaceType := surfaces[key]
				tasks = append(tasks, AutoGenTaskDef{
					Key: "run-test-" + key, ID: "T-test-run-" + key,
					Title: runTestTitle(surfaceType), Priority: "P1", EstimatedTime: "30min-1h",
					Type:         GenSurfaceTestType(TypeTestRun, key),
					SurfaceKey:   key,
					SurfaceType:  surfaceType,
					StrategyKind: "run",
				})
			}
		}
	}
	// Validation tasks (gated by auto.Validation.Quick)
	if auto.Validation.Quick {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "validate-code", ID: "T-validate-code",
			Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
			Type: TypeValidationCode, MainSession: false,
		})
		if hasUISurface(surfaceTypes) {
			tasks = append(tasks, AutoGenTaskDef{
				Key: "validate-ux", ID: "T-validate-ux",
				Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
				Type: TypeValidationUx, MainSession: true,
			})
		}
	}

	// Spec drift detection (gated by auto.ConsolidateSpecs.Quick)
	if auto.ConsolidateSpecs.Quick {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "quick-drift-detection", ID: "T-quick-doc-drift",
			Title: "Detect Spec Drift", Priority: "P2", EstimatedTime: "15min",
			Type: TypeDocDrift,
		})
	}

	// Clean code task (gated by auto.CleanCode.Quick)
	if auto.CleanCode.Quick {
		tasks = append(tasks, AutoGenTaskDef{
			Key: "quick-clean-code", ID: "T-clean-code",
			Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
			Type: TypeCleanCode,
		})
	}

	resolveQuickDeps(tasks, surfaceTypes, surfaces, executionOrder, auto)

	return tasks
}

// renderBody renders the template content using text/template with autogenTemplateData.
// The template data is pre-formatted by buildAutogenTemplateData before reaching this function.
// Metadata frontmatter is stripped before parsing.
func renderBody(templateContent string, data autogenTemplateData) (string, error) {
	// Strip metadata frontmatter before parsing (metadata is not part of rendered output)
	body := stripAutogenMetadata(templateContent)

	tmpl, err := template.New("autogen").Option("missingkey=error").Parse(body)
	if err != nil {
		return "", fmt.Errorf("autogen template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("autogen template execute error: %w", err)
	}

	return buf.String(), nil
}

// formatSurfaceTypes formats a slice of surface types as newline-separated "- type" items.
// Returns empty string when the slice is empty.
func formatSurfaceTypes(types []string) string {
	if len(types) == 0 {
		return ""
	}
	var lines []string
	for _, t := range types {
		lines = append(lines, "- "+t)
	}
	return strings.Join(lines, "\n")
}

// formatAcceptanceCriteria formats a slice of acceptance criteria as newline-separated
// unchecked checklist items. Returns empty string when the slice is empty.
func formatAcceptanceCriteria(criteria []string) string {
	if len(criteria) == 0 {
		return ""
	}
	var lines []string
	for _, ac := range criteria {
		lines = append(lines, "- [ ] "+ac)
	}
	return strings.Join(lines, "\n")
}

// buildAutogenTemplateData constructs the template data from a BodyContext and AutoGenTaskDef.
// Pre-formats all fields so the template engine only does simple field substitution.
func buildAutogenTemplateData(def AutoGenTaskDef, ctx BodyContext) autogenTemplateData {
	surfaceTypesStr := formatSurfaceTypes(ctx.SurfaceTypes)
	if surfaceTypesStr == "" {
		surfaceTypesStr = "See .forge/config.yaml"
	}

	acStr := formatAcceptanceCriteria(ctx.AcceptanceCriteria)
	if acStr == "" {
		acStr = "- [ ] All acceptance criteria met"
	}

	docTaskACStr := ""
	if len(ctx.DocTaskCriteria) > 0 {
		docTaskACStr = serializeDocTaskAC(ctx.DocTaskCriteria)
	}

	return autogenTemplateData{
		TaskID:             def.ID,
		TaskType:           def.Type,
		FeatureSlug:        ctx.FeatureSlug,
		Mode:               ctx.Mode,
		SurfaceKey:         def.SurfaceKey,
		SurfaceType:        def.SurfaceType,
		SurfaceTypes:       surfaceTypesStr,
		AcceptanceCriteria: acStr,
		DocTaskCriteria:    docTaskACStr,
	}
}

// serializeDocTaskAC serializes a DocTaskCriteria map into markdown sub-sections.
// Keys are sorted alphabetically for deterministic output.
// Format per entry:
//
//	### task-name
//	<raw AC content>
//
// When AC content is empty, displays "> No acceptance criteria defined." as placeholder.
func serializeDocTaskAC(criteria map[string]string) string {
	keys := make([]string, 0, len(criteria))
	for k := range criteria {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sections []string
	for _, key := range keys {
		content := criteria[key]
		if strings.TrimSpace(content) == "" {
			content = "> No acceptance criteria defined."
		}
		sections = append(sections, "### "+key+"\n"+content)
	}
	return strings.Join(sections, "\n\n")
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
	fmt.Fprintf(&buf, "surface-key: %q\n", def.SurfaceKey)
	fmt.Fprintf(&buf, "surface-type: %q\n", def.SurfaceType)
	if def.MainSession {
		buf.WriteString("mainSession: true\n")
	}
	buf.WriteString("---\n\n")

	// Body — try embed template first, fallback to legacy behavior
	templateFile := autogenTemplatePath(def.Type)
	tmplData, err := autogenTemplateFS.ReadFile(templateFile)
	if err == nil {
		// Template loaded successfully — render with text/template engine
		tplData := buildAutogenTemplateData(def, ctx)
		rendered, renderErr := renderBody(string(tmplData), tplData)
		if renderErr != nil {
			return nil, fmt.Errorf("rendering template %s: %w", templateFile, renderErr)
		}
		buf.WriteString(rendered)

		// Append TestType note if present
		if def.SurfaceType != "" {
			fmt.Fprintf(&buf, "\nType: **%s**\n", def.SurfaceType)
		}

		// Append StrategyContent after template content if present
		if len(def.StrategyContent) > 0 {
			buf.WriteString("\n\n")
			buf.Write(def.StrategyContent)
		}

		return []byte(buf.String()), nil
	}
	// Template file read failed — fall through to legacy behavior

	// Legacy fallback body generation
	if def.StrategyKind != "" {
		if len(def.StrategyContent) > 0 {
			fmt.Fprintf(&buf, "# %s\n\n", def.Title)
			if def.SurfaceType != "" {
				fmt.Fprintf(&buf, "Type: **%s**\n\n", def.SurfaceType)
			}
			buf.Write(def.StrategyContent)
		} else {
			fmt.Fprintf(&buf, "# %s\n\nRead docs/conventions/testing/ for test generation strategy.", def.Title)
			if def.SurfaceType != "" {
				fmt.Fprintf(&buf, " Type: %q.", def.SurfaceType)
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
func resolveBreakdownDeps(tasks []AutoGenTaskDef, surfaceTypes []string, surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig) {
	if !auto.Test.Full && !auto.ConsolidateSpecs.Full && !auto.CleanCode.Full && !auto.Validation.Full {
		return // no tasks to wire
	}

	var lastRunID string
	if auto.Test.Full {
		// Pipeline: gen-journeys -> eval-journey -> gen-contracts -> eval-contract -> gen-scripts-per-type -> run-test(s)
		evalJourneyIdx := findTaskIndexOrPanic(tasks, "T-eval-journey")
		genContractsIdx := findTaskIndexOrPanic(tasks, "T-test-gen-contracts")
		evalContractIdx := findTaskIndexOrPanic(tasks, "T-eval-contract")

		// eval-journey depends on single gen-journeys task
		genJourneysIdx := findTaskIndexOrPanic(tasks, "T-test-gen-journeys")
		tasks[evalJourneyIdx].Dependencies = []string{tasks[genJourneysIdx].ID}

		// gen-contracts depends on eval-journey
		tasks[genContractsIdx].Dependencies = []string{tasks[evalJourneyIdx].ID}

		// eval-contract depends on gen-contracts
		tasks[evalContractIdx].Dependencies = []string{tasks[genContractsIdx].ID}

		// gen-scripts depend on eval-contract
		for _, typ := range surfaceTypes {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-scripts-"+typ)
			tasks[idx].Dependencies = []string{tasks[evalContractIdx].ID}
		}

		// Wire run-test task(s)
		lastRunID = wireRunTestChain(tasks, surfaceTypes, surfaces, executionOrder)
	}
	// T-validate-code depends on last run-test (if e2e tasks exist)
	validateIdx := findTaskIndex(tasks, "T-validate-code")
	if validateIdx >= 0 && auto.Test.Full && lastRunID != "" {
		tasks[validateIdx].Dependencies = []string{lastRunID}
	}
	// T-validate-ux depends on last run-test (if e2e tasks exist)
	uxIdx := findTaskIndex(tasks, "T-validate-ux")
	if uxIdx >= 0 && auto.Test.Full && lastRunID != "" {
		tasks[uxIdx].Dependencies = []string{lastRunID}
	}

	// T-specs-consolidate depends on last run-test (if e2e tasks exist) or nothing
	if auto.ConsolidateSpecs.Full {
		specsIdx := findTaskIndex(tasks, "T-specs-consolidate")
		if specsIdx >= 0 && auto.Test.Full && lastRunID != "" {
			tasks[specsIdx].Dependencies = []string{lastRunID}
		}
	}

	// T-clean-code depends on last business task (resolved by caller via ResolveFirstTestDep)
	// The first test task depends on T-clean-code when both exist (resolved in BuildIndex)
}

// resolveQuickDeps sets dependency chains for quick test tasks using staged across types topology.
// For multi-surface projects: T-test-gen-journeys is the direct upstream of all T-test-run-{key} tasks.
// T-test-run-{key} tasks form a serial chain ordered by executionOrder.
// Downstream tasks (drift, validation) depend on the last run-test in the chain.
func resolveQuickDeps(tasks []AutoGenTaskDef, _ []string, surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig) {
	if !auto.Test.Quick && !auto.ConsolidateSpecs.Quick && !auto.CleanCode.Quick && !auto.Validation.Quick {
		return // no tasks to wire
	}

	var lastRunID string
	if auto.Test.Quick {
		// Wire run-test task(s): first run-test depends on gen-journeys
		// (no gen-contracts/gen-scripts in Quick mode)
		// Serial chain: T-test-run-{key1} -> T-test-run-{key2} -> ...
		lastRunID = wireQuickRunTestChain(tasks, surfaces, executionOrder)
	}

	// T-validate-code depends on last run-test (if e2e tasks exist) or nothing
	if auto.Validation.Quick {
		validateIdx := findTaskIndex(tasks, "T-validate-code")
		if validateIdx >= 0 && auto.Test.Quick && lastRunID != "" {
			tasks[validateIdx].Dependencies = []string{lastRunID}
		}
		// T-validate-ux depends on last run-test (if e2e tasks exist) or nothing
		uxIdx := findTaskIndex(tasks, "T-validate-ux")
		if uxIdx >= 0 && auto.Test.Quick && lastRunID != "" {
			tasks[uxIdx].Dependencies = []string{lastRunID}
		}
	}

	// T-quick-doc-drift depends on last run-test (if e2e tasks exist) or nothing
	if auto.ConsolidateSpecs.Quick {
		idx := findTaskIndex(tasks, "T-quick-doc-drift")
		if idx >= 0 && auto.Test.Quick && lastRunID != "" {
			tasks[idx].Dependencies = []string{lastRunID}
		}
	}

	// T-clean-code depends on last business task (resolved by caller via ResolveFirstTestDep)
}

// wireRunTestChain wires the run-test task(s) dependency chain for Breakdown mode.
// For single-surface projects: T-test-run depends on all gen-scripts.
// For multi-surface projects: first T-test-run-{key} depends on all gen-scripts,
// subsequent T-test-run-{key} tasks form a serial chain.
// Returns the ID of the last run-test task in the chain.
func wireRunTestChain(tasks []AutoGenTaskDef, surfaceTypes []string, surfaces map[string]string, executionOrder []string) string {
	singleSurface := isSingleSurface(surfaces)

	if singleSurface {
		// Single surface: T-test-run depends on all gen-scripts
		runIdx := findTaskIndexOrPanic(tasks, "T-test-run")
		var genDeps []string
		for _, typ := range surfaceTypes {
			idx := findTaskIndexOrPanic(tasks, "T-test-gen-scripts-"+typ)
			genDeps = append(genDeps, tasks[idx].ID)
		}
		tasks[runIdx].Dependencies = genDeps
		return tasks[runIdx].ID
	}

	// Multi-surface: serial chain
	// First run-test depends on all gen-scripts, subsequent ones depend on previous
	var prevRunID string
	for i, key := range executionOrder {
		runID := "T-test-run-" + key
		runIdx := findTaskIndexOrPanic(tasks, runID)

		if i == 0 {
			// First run-test depends on all gen-scripts
			var genDeps []string
			for _, typ := range surfaceTypes {
				idx := findTaskIndexOrPanic(tasks, "T-test-gen-scripts-"+typ)
				genDeps = append(genDeps, tasks[idx].ID)
			}
			tasks[runIdx].Dependencies = genDeps
		} else {
			// Subsequent run-test depends on previous run-test (serial chain)
			tasks[runIdx].Dependencies = []string{prevRunID}
		}
		prevRunID = runID
	}

	return prevRunID
}

// wireQuickRunTestChain wires the run-test task(s) dependency chain for Quick mode.
// Quick mode skips gen-contracts and gen-scripts: first run-test depends on gen-journeys.
// For single-surface projects: T-test-run depends on gen-journeys.
// For multi-surface projects: first T-test-run-{key} depends on gen-journeys,
// subsequent T-test-run-{key} tasks form a serial chain.
// Returns the ID of the last run-test task in the chain.
func wireQuickRunTestChain(tasks []AutoGenTaskDef, surfaces map[string]string, executionOrder []string) string {
	genJourneysID := "T-test-gen-journeys"
	singleSurface := isSingleSurface(surfaces)

	if singleSurface {
		// Single surface: T-test-run depends on gen-journeys
		runIdx := findTaskIndexOrPanic(tasks, "T-test-run")
		tasks[runIdx].Dependencies = []string{genJourneysID}
		return tasks[runIdx].ID
	}

	// Multi-surface: serial chain
	// First run-test depends on gen-journeys, subsequent ones depend on previous
	var prevRunID string
	for i, key := range executionOrder {
		runID := "T-test-run-" + key
		runIdx := findTaskIndexOrPanic(tasks, runID)

		if i == 0 {
			// First run-test depends on gen-journeys
			tasks[runIdx].Dependencies = []string{genJourneysID}
		} else {
			// Subsequent run-test depends on previous run-test (serial chain)
			tasks[runIdx].Dependencies = []string{prevRunID}
		}
		prevRunID = runID
	}

	return prevRunID
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

	// Only resolve when surface test tasks exist (they have T-test-gen-journeys prefix)
	if findTaskIndexByPrefix(tasks, "T-test-gen-journeys") < 0 {
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
		SurfaceKey:    d.SurfaceKey,
		SurfaceType:   d.SurfaceType,
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
	}
}

// ResolveReviewDocDep sets the dependency of T-review-doc on all doc-category business tasks.
// Root cause: previously depended only on the single highest-ID business task,
// allowing review-doc to run before doc tasks that were added later.
func ResolveReviewDocDep(task *AutoGenTaskDef, existingTasks map[string]Task) {
	var deps []string
	for _, t := range existingTasks {
		if isAutoGenForDep(t.ID) {
			continue
		}
		if CategoryForType(t.Type) == CategoryDoc {
			deps = append(deps, t.ID)
		}
	}
	sort.Strings(deps)
	if len(deps) > 0 {
		task.Dependencies = deps
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

// --- Metadata frontmatter support ---

// autogenMetadata holds parsed metadata frontmatter from autogen templates.
type autogenMetadata struct {
	Type      string
	Category  string
	Variables []string
}

// parseAutogenMetadata extracts metadata from between the first pair of ---
// markers in a template file. Returns the body and parsed metadata.
// If no frontmatter is found, returns the original content with nil metadata.
func parseAutogenMetadata(content string) (body string, meta *autogenMetadata) {
	trimmed := strings.TrimLeft(content, " \t\n")
	if !strings.HasPrefix(trimmed, "---") {
		return content, nil
	}

	afterOpen := trimmed[3:]
	if len(afterOpen) > 0 && afterOpen[0] == '\n' {
		afterOpen = afterOpen[1:]
	} else if len(afterOpen) > 1 && afterOpen[0] == '\r' && afterOpen[1] == '\n' {
		afterOpen = afterOpen[2:]
	}

	closeIdx := strings.Index(afterOpen, "\n---")
	if closeIdx < 0 {
		return content, nil
	}

	frontmatter := afterOpen[:closeIdx]
	remaining := afterOpen[closeIdx+4:]
	if len(remaining) > 0 && remaining[0] == '\n' {
		remaining = remaining[1:]
	} else if len(remaining) > 1 && remaining[0] == '\r' && remaining[1] == '\n' {
		remaining = remaining[2:]
	}

	meta = &autogenMetadata{}
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		switch {
		case strings.HasPrefix(line, "type:"):
			meta.Type = strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "type:")), "\"")
		case strings.HasPrefix(line, "category:"):
			meta.Category = strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "category:")), "\"")
		case strings.HasPrefix(line, "- ") && meta.Variables != nil:
			varName := strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "- ")), "\"")
			meta.Variables = append(meta.Variables, varName)
		case strings.HasPrefix(line, "variables:"):
			meta.Variables = []string{}
		}
	}

	return remaining, meta
}

// stripAutogenMetadata removes metadata frontmatter from template content.
func stripAutogenMetadata(content string) string {
	body, _ := parseAutogenMetadata(content)
	return body
}

// validateAutogenVariables checks that each variable declared in metadata
// exists as an exported field on the autogenTemplateData struct.
func validateAutogenVariables(meta *autogenMetadata, structType reflect.Type) error {
	if meta == nil || len(meta.Variables) == 0 {
		return nil
	}

	var mismatches []string
	for _, varName := range meta.Variables {
		if _, ok := structType.FieldByName(varName); !ok {
			mismatches = append(mismatches, varName)
		}
	}

	if len(mismatches) > 0 {
		return fmt.Errorf("metadata variables not found in %s struct: %s", structType.Name(), strings.Join(mismatches, ", "))
	}
	return nil
}
