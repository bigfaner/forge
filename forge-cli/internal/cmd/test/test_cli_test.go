package test

import (
	"strings"
	"testing"
)

func TestTestCommand_NoE2EReferences(t *testing.T) {
	t.Run("parent command long description has no e2e references", func(t *testing.T) {
		long := Cmd.Long
		lower := strings.ToLower(long)
		if strings.Contains(lower, "e2e") {
			t.Errorf("Cmd.Long should not contain 'e2e', got: %q", long)
		}
		if strings.Contains(lower, "end-to-end") {
			t.Errorf("Cmd.Long should not contain 'end-to-end', got: %q", long)
		}
	})

	t.Run("run-journey command long description has no e2e references", func(t *testing.T) {
		long := testRunJourneyCmd.Long
		lower := strings.ToLower(long)
		if strings.Contains(lower, "e2e") {
			t.Errorf("run-journey Long should not contain 'e2e', got: %q", long)
		}
		if strings.Contains(lower, "end-to-end") {
			t.Errorf("run-journey Long should not contain 'end-to-end', got: %q", long)
		}
	})
}

func TestRunJourneyCommand_ReferencesJustTest(t *testing.T) {
	t.Run("run-journey description references just test", func(t *testing.T) {
		long := testRunJourneyCmd.Long
		if !strings.Contains(long, "just test") {
			t.Errorf("run-journey Long should reference 'just test', got: %q", long)
		}
	})
}
