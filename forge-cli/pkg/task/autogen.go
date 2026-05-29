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
	"forge-cli/pkg/types"
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

// ValidateAutogenTemplates validates the pipeline registry and autogen template integrity.
//
// Phase 1 (static) validation runs in init() and panics on failure — by the time this
// function is called from main(), the registry is already validated. This function
// performs additional template-level checks (parse/execute validation against embedded
// template files) that cannot be expressed in the registry's structural validation.
//
// Must be called from the CLI main() startup path, NOT from an init() function.
func ValidateAutogenTemplates() error {
	// Phase 1 (registry structural validation) already ran in init().
	// Perform template-specific validation here.
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
var uiSurfaceTypes = map[types.SurfaceType]bool{
	types.SurfaceTUI:    true,
	types.SurfaceWeb:    true,
	types.SurfaceMobile: true,
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

// isSkipTestIntent returns true when the intent should skip test pipeline tasks.
func isSkipTestIntent(intent string) bool {
	return intent == "refactor" || intent == "cleanup"
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
		Priority:      types.Priority(d.Priority),
		EstimatedTime: d.EstimatedTime,
		Dependencies:  d.Dependencies,
		Status:        types.StatusPending,
		File:          fileName,
		Record:        path.Join("records", fileName),
		Breaking:      d.Breaking,
		SurfaceKey:    d.SurfaceKey,
		SurfaceType:   d.SurfaceType,
		MainSession:   d.MainSession,
		Type:          d.Type,
	}
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

// ---------------------------------------------------------------------------
// Legacy bridge functions (deprecated — use GenerateTestTasks directly)
// ---------------------------------------------------------------------------

// GetBreakdownTestTasks generates breakdown-mode auto tasks via the pipeline registry.
// Kept for backward compatibility with existing tests.
func GetBreakdownTestTasks(surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig, intent string) []AutoGenTaskDef {
	return GenerateTestTasks("breakdown", surfaces, executionOrder, auto, intent, nil, nil)
}

// GetQuickTestTasks generates quick-mode auto tasks via the pipeline registry.
// Kept for backward compatibility with existing tests.
func GetQuickTestTasks(surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig, intent string) []AutoGenTaskDef {
	return GenerateTestTasks("quick", surfaces, executionOrder, auto, intent, nil, nil)
}

// findTaskIndex returns the index of the task with the given ID, or -1.
// Kept for backward compatibility with existing tests.
func findTaskIndex(tasks []AutoGenTaskDef, id string) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

