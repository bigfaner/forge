package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	if !strings.Contains(err.Hint, "forge feature") {
		t.Errorf("Hint should mention forge feature: %s", err.Hint)
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
			t.Fatal("expected non-nil error")
		}
		if !strings.Contains(err.Cause, "statusEnum") {
			t.Errorf("Cause should mention statusEnum: %s", err.Cause)
		}
	})
}

// --- New error code constants ---

func TestNewErrorCodeConstants(t *testing.T) {
	assert.Equal(t, ErrorCode("INVALID_TRANSITION"), ErrInvalidTransition)
	assert.Equal(t, ErrorCode("INVALID_PATH"), ErrInvalidPath)
	assert.Equal(t, ErrorCode("EVAL_PARSE_FAILURE"), ErrEvalParseFailure)
	assert.Equal(t, ErrorCode("CONTRACT_UNVERIFIABLE"), ErrContractUnverifiable)
}

// --- New factory functions ---

func TestNewErrInvalidTransition(t *testing.T) {
	err := NewErrInvalidTransition("completed", "pending", "use reopen instead")
	assert.Equal(t, ErrInvalidTransition, err.Code)
	assert.Contains(t, err.Message, "completed")
	assert.Contains(t, err.Message, "pending")
	assert.NotEmpty(t, err.Cause)
	assert.Contains(t, err.Hint, "reopen")
	assert.NotEmpty(t, err.Action)
}

func TestNewErrInvalidPath(t *testing.T) {
	err := NewErrInvalidPath("../../../etc/passwd")
	assert.Equal(t, ErrInvalidPath, err.Code)
	assert.Contains(t, err.Message, "path")
	assert.Contains(t, err.Cause, "..")
	assert.NotEmpty(t, err.Hint)
	assert.NotEmpty(t, err.Action)
}

func TestNewErrEvalParseFailure(t *testing.T) {
	err := NewErrEvalParseFailure("score: not-a-number")
	assert.Equal(t, ErrEvalParseFailure, err.Code)
	assert.Contains(t, err.Message, "parse")
	assert.Contains(t, err.Cause, "score")
	assert.NotEmpty(t, err.Hint)
	assert.NotEmpty(t, err.Action)
}

func TestNewErrContractUnverifiable(t *testing.T) {
	err := NewErrContractUnverifiable("contracts/api.md")
	assert.Equal(t, ErrContractUnverifiable, err.Code)
	assert.Contains(t, err.Message, "unverifiable")
	assert.Contains(t, err.Cause, "contracts/api.md")
	assert.NotEmpty(t, err.Hint)
	assert.NotEmpty(t, err.Action)
}

// --- ExitCode method ---

func TestExitCode_BlockingErrors(t *testing.T) {
	blockingCodes := []ErrorCode{
		ErrInvalidTransition,
		ErrInvalidPath,
		ErrContractUnverifiable,
	}
	for _, code := range blockingCodes {
		err := &AIError{Code: code, Message: "test"}
		assert.Equal(t, 2, err.ExitCode(), "expected exit code 2 for %s", code)
	}
}

func TestExitCode_RetryableErrors(t *testing.T) {
	retryableCodes := []ErrorCode{
		ErrEvalParseFailure,
		ErrNoProject,
		ErrNoFeature,
		ErrNotFound,
		ErrConflict,
		ErrInvalidInput,
		ErrValidation,
	}
	for _, code := range retryableCodes {
		err := &AIError{Code: code, Message: "test"}
		assert.Equal(t, 1, err.ExitCode(), "expected exit code 1 for %s", code)
	}
}

func TestExitCode_DefaultRetryable(t *testing.T) {
	err := &AIError{Code: ErrorCode("UNKNOWN_CODE"), Message: "test"}
	assert.Equal(t, 1, err.ExitCode(), "unknown error codes should default to exit code 1")
}

// --- Existing error behavior unchanged ---

func TestExistingFactoryFunctionsUnchanged(t *testing.T) {
	tests := []struct {
		name     string
		err      *AIError
		wantCode ErrorCode
	}{
		{"ErrProjectNotFound", ErrProjectNotFound(), ErrNoProject},
		{"ErrFeatureNotSet", ErrFeatureNotSet(), ErrNoFeature},
		{"ErrTaskNotFound", ErrTaskNotFound("1.0"), ErrNotFound},
		{"ErrNoInput", ErrNoInput("detail"), ErrInvalidInput},
		{"ErrInvalidJSON", ErrInvalidJSON("file.json", "bad"), ErrValidation},
		{"ErrFileNotFound", ErrFileNotFound("x.txt"), ErrNotFound},
		{"ErrNoPendingTasks", ErrNoPendingTasks(), ErrNotFound},
		{"ErrDependenciesNotMet", ErrDependenciesNotMet("1.0", []string{"0.1"}), ErrConflict},
		{"ErrDataIntegrity", ErrDataIntegrity([]string{"x"}), ErrConflict},
		{"ErrInvalidStatus", ErrInvalidStatus("bad", []string{"ok"}), ErrValidation},
		{"ErrMissingFields", ErrMissingFields([]string{"x"}), ErrValidation},
		{"ErrFeatureNotFound", ErrFeatureNotFound("slug"), ErrNotFound},
		{"ErrNoTestEvidence", ErrNoTestEvidence(), ErrValidation},
		{"ErrUnmetAcceptanceCriteria", ErrUnmetAcceptanceCriteria([]string{"x"}), ErrValidation},
		{"ErrTaskIDConflict", ErrTaskIDConflict("1.0"), ErrConflict},
		{"ErrInvalidDependency", ErrInvalidDependency([]string{"9.9"}), ErrValidation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantCode, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message)
		})
	}
}
