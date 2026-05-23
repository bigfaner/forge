package forensic

import (
	"io"
	"os"
	"strings"
	"time"
)

func extractUserContent(entry jsonlEntry) string {
	if entry.Message.Role != "user" {
		return ""
	}
	if text := entry.Message.textContent(); text != "" {
		return text
	}
	if entry.Content != "" {
		return entry.Content
	}
	return ""
}

func detectSkills(content string, skills *[]string) {
	lower := strings.ToLower(content)
	prefixes := []string{"/forge:", "<command-name>/", "<command-name>"}

	for _, prefix := range prefixes {
		searchFrom := 0
		for {
			idx := strings.Index(lower[searchFrom:], prefix)
			if idx == -1 {
				break
			}
			idx += searchFrom
			start := idx + len(prefix)
			end := start
			for end < len(lower) && (isAlphaNumeric(lower[end]) || lower[end] == '-' || lower[end] == '_') {
				end++
			}
			if end > start {
				name := lower[start:end]
				found := false
				for _, s := range *skills {
					if s == name {
						found = true
						break
					}
				}
				if !found {
					*skills = append(*skills, name)
				}
			}
			searchFrom = end
		}
	}
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

func truncate(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "..."
}

func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() { _ = in.Close() }()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() { _ = out.Close() }()

	// Best-effort copy; don't fail the extract if copy fails
	_, _ = io.Copy(out, in)
}

// aggregateToolInput extracts tool-specific details into the summary.
func aggregateToolInput(tool string, input any, s *extractSummary) {
	m, ok := input.(map[string]any)
	if !ok {
		return
	}

	switch tool {
	case "Read":
		if fp, _ := m["file_path"].(string); fp != "" {
			s.FilesRead = appendUniq(s.FilesRead, fp)
		}
	case "Edit", "Write":
		if fp, _ := m["file_path"].(string); fp != "" {
			s.FilesWritten = appendUniq(s.FilesWritten, fp)
		}
	case "Grep":
		if p, _ := m["pattern"].(string); p != "" {
			s.GrepPatterns = appendUniq(s.GrepPatterns, p)
		}
	case "Bash":
		if cmd, _ := m["command"].(string); cmd != "" {
			s.Commands = appendUniq(s.Commands, truncate(cmd, 200))
		}
	case "Agent":
		name := "unknown"
		if t, _ := m["subagent_type"].(string); t != "" {
			name = t
		}
		found := false
		for i := range s.AgentsSpawned {
			if s.AgentsSpawned[i].Name == name {
				s.AgentsSpawned[i].Count++
				found = true
				break
			}
		}
		if !found {
			s.AgentsSpawned = append(s.AgentsSpawned, toolAggEntry{Name: name, Count: 1})
		}
	}
}

func appendUniq(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}

func computeDurationMs(startTS, endTS string) int64 {
	t1, err1 := parseTimestamp(startTS)
	if err1 != nil {
		return -1
	}
	t2, err2 := parseTimestamp(endTS)
	if err2 != nil {
		return -1
	}
	return t2.Sub(t1).Milliseconds()
}

func parseTimestamp(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
	}
	return t, err
}

func addToAgg(entries *[]toolAggEntry, name string) {
	for i := range *entries {
		if (*entries)[i].Name == name {
			(*entries)[i].Count++
			return
		}
	}
	*entries = append(*entries, toolAggEntry{Name: name, Count: 1})
}
