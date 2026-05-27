// Package template handles task template loading and processing.
package template

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"text/template"
)

//go:embed data/*.md
var templateFS embed.FS

// TaskTemplateData holds the data for rendering task creation templates.
type TaskTemplateData struct {
	ID               string // task ID
	Title            string // task title
	Priority         string // task priority (P0, P1, P2)
	EstimatedTime    string // estimated time (e.g. "30min")
	Description      string // task description / root cause
	SourceTaskID     string // source task ID
	SurfaceKey       string // surface key, empty string → omit field + omit Surface Inference section
	SurfaceType      string // surface type, empty string → omit field
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
	IDPrefix      string // auto-generate ID as prefix-N (e.g. "fix" → fix-1, fix-2)
}

// templateDefaults defines the fixed values for each template.
// Keys match the template filename (without .md extension), which is also the type value.
var templateDefaults = map[string]Defaults{
	"coding.fix": {
		Priority:      "P0",
		Breaking:      true,
		EstimatedTime: "30min",
		Type:          "coding.fix",
		IDPrefix:      "fix",
	},
	"coding.cleanup": {
		Priority:      "P0",
		Breaking:      false,
		EstimatedTime: "15min",
		Type:          "coding.cleanup",
		IDPrefix:      "fix",
	},
}

// Get returns the template content for the given name (without .md extension).
func Get(name string) (string, error) {
	data, err := templateFS.ReadFile("data/" + name + ".md")
	if err != nil {
		return "", fmt.Errorf("template %q not found", name)
	}
	return string(data), nil
}

// GetDefaults returns the fixed field values for a template.
func GetDefaults(name string) (Defaults, error) {
	defs, ok := templateDefaults[name]
	if !ok {
		return Defaults{}, fmt.Errorf("no defaults for template %q", name)
	}
	return defs, nil
}

// Execute renders the named template with the given data using text/template.
// It uses missingkey=error to catch typos at render time.
// Metadata frontmatter is stripped before parsing.
func Execute(name string, data TaskTemplateData) (string, error) {
	raw, err := Get(name)
	if err != nil {
		return "", err
	}

	// Strip metadata frontmatter before parsing (metadata is not part of rendered output)
	body := stripTaskTemplateMetadata(raw)

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

// List returns all available template names (without .md extension).
func List() []string {
	entries, err := fs.ReadDir(templateFS, "data")
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		name, ok := strings.CutSuffix(e.Name(), ".md")
		if ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// stripTaskTemplateMetadata removes metadata frontmatter from template content.
// If no metadata frontmatter is found (no leading ---), returns the original content.
func stripTaskTemplateMetadata(content string) string {
	trimmed := strings.TrimLeft(content, " \t\n")
	if !strings.HasPrefix(trimmed, "---") {
		return content
	}

	afterOpen := trimmed[3:]
	if len(afterOpen) > 0 && afterOpen[0] == '\n' {
		afterOpen = afterOpen[1:]
	} else if len(afterOpen) > 1 && afterOpen[0] == '\r' && afterOpen[1] == '\n' {
		afterOpen = afterOpen[2:]
	}

	closeIdx := strings.Index(afterOpen, "\n---")
	if closeIdx < 0 {
		return content
	}

	remaining := afterOpen[closeIdx+4:]
	if len(remaining) > 0 && remaining[0] == '\n' {
		remaining = remaining[1:]
	} else if len(remaining) > 1 && remaining[0] == '\r' && remaining[1] == '\n' {
		remaining = remaining[2:]
	}
	return remaining
}
