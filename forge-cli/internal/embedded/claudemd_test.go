package embedded

import (
	"strings"
	"testing"
)

func TestCLAUDEmdTemplate(t *testing.T) {
	t.Run("template is non-empty", func(t *testing.T) {
		if CLAUDEmdTemplate == "" {
			t.Fatal("CLAUDEmdTemplate must not be empty")
		}
	})

	t.Run("contains expected sections", func(t *testing.T) {
		expected := []string{
			"Think Before Coding",
			"Simplicity First",
			"Surgical Changes",
			"Goal-Driven Execution",
		}
		for _, section := range expected {
			if !strings.Contains(CLAUDEmdTemplate, section) {
				t.Errorf("template should contain %q", section)
			}
		}
	})

	t.Run("template is not just whitespace", func(t *testing.T) {
		trimmed := strings.TrimSpace(CLAUDEmdTemplate)
		if len(trimmed) < 100 {
			t.Errorf("template seems too short (%d chars), expected substantive content", len(trimmed))
		}
	})
}
