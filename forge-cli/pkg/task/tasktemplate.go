package task

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	"forge-cli/pkg/types"
)

// TemplateData holds the data for rendering task creation templates.
type TemplateData struct {
	ID               string // task ID
	Title            string // task title
	Priority         string // task priority (P0, P1, P2)
	EstimatedTime    string // estimated time (e.g. "30min")
	Description      string // task description / root cause
	SourceTaskID     string // source task ID
	SurfaceKey       string // surface key, empty string -> omit field + omit Surface Inference section
	SurfaceType      string // surface type, empty string -> omit field
	SourceFiles      string // source files (from Vars)
	TestScript       string // test script (from Vars)
	TestResults      string // test results path (from Vars)
	ScopeDescription string // task-level scope description (from Vars)
}

// Defaults holds fixed field values from a template.
// These are applied to AddTaskOpts when --type matches a template,
// so callers don't need to pass them as flags.
type Defaults struct {
	Priority      string
	Breaking      bool
	EstimatedTime string
	Type          string // default task type (e.g. "coding.fix")
	IDPrefix      string // auto-generate ID as prefix-N (e.g. "fix" -> fix-1, fix-2)
}

// taskTemplateDefaults defines the fixed values for each template.
// Keys match the template filename (without .md extension), which is also the type value.
var taskTemplateDefaults = map[string]Defaults{
	"coding.fix": {
		Priority:      string(types.PriorityP0),
		Breaking:      true,
		EstimatedTime: "30min",
		Type:          "coding.fix",
		IDPrefix:      "fix",
	},
	"coding.cleanup": {
		Priority:      string(types.PriorityP0),
		Breaking:      false,
		EstimatedTime: "15min",
		Type:          "coding.cleanup",
		IDPrefix:      "fix",
	},
}

// GetTaskTemplate returns the template content for the given name (without .md extension).
// It reads from the autogen embed FS (templates/ directory).
func GetTaskTemplate(name string) (string, error) {
	data, err := autogenTemplateFS.ReadFile("templates/" + name + ".md")
	if err != nil {
		return "", fmt.Errorf("template %q not found", name)
	}
	return string(data), nil
}

// GetTaskTemplateDefaults returns the fixed field values for a template.
func GetTaskTemplateDefaults(name string) (Defaults, error) {
	defs, ok := taskTemplateDefaults[name]
	if !ok {
		return Defaults{}, fmt.Errorf("no defaults for template %q", name)
	}
	return defs, nil
}

// ExecuteTaskTemplate renders the named template with the given data using text/template.
// It uses missingkey=error to catch typos at render time.
// Metadata frontmatter is stripped before parsing.
func ExecuteTaskTemplate(name string, data TemplateData) (string, error) {
	raw, err := GetTaskTemplate(name)
	if err != nil {
		return "", err
	}

	// Strip metadata frontmatter before parsing (metadata is not part of rendered output)
	body := stripAutogenMetadata(raw)

	tmpl, err := template.New(name).Option("missingkey=error").Parse(body)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %q: %w", name, err)
	}
	return buf.String(), nil
}

// ListTaskTemplates returns all available task template names (without .md extension).
// This lists templates from the autogen embed FS that match known task template types.
func ListTaskTemplates() []string {
	var names []string
	for name := range taskTemplateDefaults {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
