package forgeconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WriteConfigWithSources writes a Config to .forge/config.yaml and appends YAML
// source annotations as line comments on the surfaces node. Uses yaml.Node
// round-trip API to preserve comments — never string concatenation or regex.
//
// Comment format: # source: inference:cmd-dir or # source: dependency:cobra
// The comment is purely informational and does not affect config parsing.
func WriteConfigWithSources(projectRoot string, cfg *Config, sources SourcesMap) error {
	path := configPath(projectRoot)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create .forge dir: %w", err)
	}

	// Marshal config to yaml.Node tree.
	// root.Encode(cfg) produces a MappingNode directly (not DocumentNode).
	var root yaml.Node
	if err := root.Encode(cfg); err != nil {
		return fmt.Errorf("encode config to node: %w", err)
	}

	if len(sources) > 0 {
		annotateSurfacesNode(&root, sources)
	}

	out, err := yaml.Marshal(&root)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// annotateSurfacesNode finds the surfaces node in the YAML tree and appends
// source annotations as line comments. Handles both scalar and map forms.
func annotateSurfacesNode(mappingNode *yaml.Node, sources SourcesMap) {
	if mappingNode.Kind != yaml.MappingNode {
		return
	}

	// Find the "surfaces" key in the mapping
	for i := 0; i < len(mappingNode.Content); i += 2 {
		keyNode := mappingNode.Content[i]
		if keyNode.Value != "surfaces" {
			continue
		}

		valNode := mappingNode.Content[i+1]

		switch valNode.Kind {
		case yaml.ScalarNode:
			// Scalar form: surfaces: cli
			// Use the "." key from sources
			if source, ok := sources["."]; ok {
				valNode.LineComment = formatSourceComment(source)
			}

		case yaml.MappingNode:
			// Map form: surfaces: {frontend: web, backend: api}
			// Annotate each value node with its source
			for j := 0; j < len(valNode.Content); j += 2 {
				entryKey := valNode.Content[j].Value
				entryVal := valNode.Content[j+1]
				if source, ok := sources[entryKey]; ok {
					entryVal.LineComment = formatSourceComment(source)
				}
			}
		}

		return
	}
}

// formatSourceComment formats a source annotation into a YAML line comment.
// Returns "# source: <annotation>" or empty if annotation is empty.
func formatSourceComment(annotation string) string {
	if annotation == "" {
		return ""
	}
	return "source: " + annotation
}

// ReadSurfaceComment reads the config file using yaml.Node and extracts the
// source comment from the surfaces node. Returns the comment text without the
// "# " prefix, or empty string if no comment is present.
func ReadSurfaceComment(projectRoot string) (string, error) {
	path := configPath(projectRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read config for comment: %w", err)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return "", fmt.Errorf("parse config for comment: %w", err)
	}

	if len(root.Content) == 0 {
		return "", nil
	}

	docMapping := root.Content[0]
	if docMapping.Kind != yaml.MappingNode {
		return "", nil
	}

	// Find the "surfaces" key
	for i := 0; i < len(docMapping.Content); i += 2 {
		keyNode := docMapping.Content[i]
		if keyNode.Value != "surfaces" {
			continue
		}

		valNode := docMapping.Content[i+1]

		// Collect comments from all surface entries
		switch valNode.Kind {
		case yaml.ScalarNode:
			return extractSourceFromComment(valNode.LineComment), nil
		case yaml.MappingNode:
			// For map form, return the first non-empty comment found
			for j := 0; j < len(valNode.Content); j += 2 {
				entryVal := valNode.Content[j+1]
				if c := extractSourceFromComment(entryVal.LineComment); c != "" {
					return c, nil
				}
			}
		}

		return "", nil
	}

	return "", nil
}

// extractSourceFromComment extracts the "source: ..." portion from a YAML line comment.
// Input: "source: inference:cmd-dir" (without "# " prefix)
// Returns: "source: inference:cmd-dir" or empty string
func extractSourceFromComment(comment string) string {
	comment = strings.TrimSpace(comment)
	// Strip leading "# " if present (yaml.v3 may or may not include it)
	comment = strings.TrimPrefix(comment, "# ")
	comment = strings.TrimPrefix(comment, "#")

	if strings.HasPrefix(comment, "source: ") {
		return comment
	}
	return ""
}
