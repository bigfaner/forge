package cmd

import (
	"strings"
	"testing"
)

func TestAIError_Error(t *testing.T) {
	err := &AIError{Message: "something went wrong"}
	if got := err.Error(); got != "something went wrong" {
		t.Errorf("Error() = %q, want %q", got, "something went wrong")
	}
}

func TestNewAIError(t *testing.T) {
	err := NewAIError(ErrConflict, "msg", "cause", "hint", "action")
	if err.Code != ErrConflict {
		t.Errorf("Code = %q, want %q", err.Code, ErrConflict)
	}
	if err.Message != "msg" {
		t.Errorf("Message = %q, want %q", err.Message, "msg")
	}
}

func TestErrProjectNotFound(t *testing.T) {
	err := ErrProjectNotFound()
	if err.Code != ErrNoProject {
		t.Errorf("Code = %q, want %q", err.Code, ErrNoProject)
	}
	if !strings.Contains(err.Message, "Project root") {
		t.Errorf("Message should mention project root: %s", err.Message)
	}
}

func TestErrFeatureNotSet(t *testing.T) {
	err := ErrFeatureNotSet()
	if err.Code != ErrNoFeature {
		t.Errorf("Code = %q, want %q", err.Code, ErrNoFeature)
	}
	if !strings.Contains(err.Hint, "task feature") {
		t.Errorf("Hint should mention task feature: %s", err.Hint)
	}
}

func TestErrTaskNotFound(t *testing.T) {
	err := ErrTaskNotFound("1.2.3")
	if err.Code != ErrNotFound {
		t.Errorf("Code = %q, want %q", err.Code, ErrNotFound)
	}
	if !strings.Contains(err.Message, "1.2.3") {
		t.Errorf("Message should contain task ID: %s", err.Message)
	}
}

func TestErrNoInput(t *testing.T) {
	err := ErrNoInput("need a file path")
	if err.Code != ErrInvalidInput {
		t.Errorf("Code = %q, want %q", err.Code, ErrInvalidInput)
	}
	if err.Cause != "need a file path" {
		t.Errorf("Cause = %q, want %q", err.Cause, "need a file path")
	}
}

func TestErrDependenciesNotMet(t *testing.T) {
	err := ErrDependenciesNotMet("2.1", []string{"1.1", "1.2"})
	if err.Code != ErrConflict {
		t.Errorf("Code = %q, want %q", err.Code, ErrConflict)
	}
	if !strings.Contains(err.Cause, "1.1") || !strings.Contains(err.Cause, "1.2") {
		t.Errorf("Cause should list unmet deps: %s", err.Cause)
	}
}

func TestErrDataIntegrity(t *testing.T) {
	err := ErrDataIntegrity([]string{"state mismatch", "index corrupt"})
	if err.Code != ErrConflict {
		t.Errorf("Code = %q, want %q", err.Code, ErrConflict)
	}
	if !strings.Contains(err.Cause, "state mismatch") {
		t.Errorf("Cause should include issues: %s", err.Cause)
	}
}

func TestErrFeatureNotFound(t *testing.T) {
	err := ErrFeatureNotFound("my-feature")
	if err.Code != ErrNotFound {
		t.Errorf("Code = %q, want %q", err.Code, ErrNotFound)
	}
	if !strings.Contains(err.Message, "my-feature") {
		t.Errorf("Message should contain slug: %s", err.Message)
	}
}

func TestErrMissingFields(t *testing.T) {
	err := ErrMissingFields([]string{"summary", "taskId"})
	if err.Code != ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, ErrValidation)
	}
	if !strings.Contains(err.Message, "summary") {
		t.Errorf("Message should list missing fields: %s", err.Message)
	}
}

func TestErrNoTestEvidence(t *testing.T) {
	err := ErrNoTestEvidence()
	if err.Code != ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, ErrValidation)
	}
	if !strings.Contains(err.Cause, "testsPassed") {
		t.Errorf("Cause should mention test metrics: %s", err.Cause)
	}
}

func TestErrUnmetAcceptanceCriteria(t *testing.T) {
	err := ErrUnmetAcceptanceCriteria([]string{"login works", "data persists"})
	if err.Code != ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, ErrValidation)
	}
	if !strings.Contains(err.Cause, "login works") {
		t.Errorf("Cause should list unmet criteria: %s", err.Cause)
	}
}

func TestErrInvalidStatus(t *testing.T) {
	t.Run("with valid statuses", func(t *testing.T) {
		err := ErrInvalidStatus("bad", []string{"pending", "completed"})
		if err.Code != ErrValidation {
			t.Errorf("Code = %q, want %q", err.Code, ErrValidation)
		}
		if !strings.Contains(err.Cause, "pending") {
			t.Errorf("Cause should list valid statuses: %s", err.Cause)
		}
		if !strings.Contains(err.Action, "pending") {
			t.Errorf("Action should suggest first valid status: %s", err.Action)
		}
	})

	t.Run("empty valid statuses", func(t *testing.T) {
		err := ErrInvalidStatus("bad", []string{})
		if err == nil {
			t.Error("expected non-nil error")
		}
		if !strings.Contains(err.Cause, "statusEnum") {
			t.Errorf("Cause should mention statusEnum: %s", err.Cause)
		}
	})
}
