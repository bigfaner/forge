// Package template handles task template loading and processing.
package template

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed data/*.md
var templateFS embed.FS

// Defaults holds fixed field values from a template.
// These are applied to AddTaskOpts when --template is used,
// so callers don't need to pass them as flags.
type Defaults struct {
	Priority      string
	Breaking      bool
	EstimatedTime string
	IDPrefix      string // auto-generate ID as prefix-N (e.g. "fix" → fix-1, fix-2)
}

// templateDefaults defines the fixed values for each template.
var templateDefaults = map[string]Defaults{
	"fix-task": {
		Priority:      "P0",
		Breaking:      true,
		EstimatedTime: "30min",
		IDPrefix:      "fix",
	},
	"cleanup-task": {
		Priority:      "P0",
		Breaking:      true,
		EstimatedTime: "15min",
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
