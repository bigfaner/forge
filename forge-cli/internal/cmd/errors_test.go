package cmd

import (
	"strings"
	"testing"

	"forge-cli/internal/cmd/base"

	"github.com/stretchr/testify/assert"
)

func TestAIError_Error(t *testing.T) {
	err := &base.AIError{Message: "something went wrong"}
	if got := err.Error(); got != "something went wrong" {
		t.Errorf("Error() = %q, want %q", got, "something went wrong")
	}
}

func TestNewAIError(t *testing.T) {
	err := base.NewAIError(base.ErrConflict, "msg", "cause", "hint", "action")
	if err.Code != base.ErrConflict {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrConflict)
	}
	if err.Message != "msg" {
		t.Errorf("Message = %q, want %q", err.Message, "msg")
	}
}

func TestErrProjectNotFound(t *testing.T) {
	err := base.ErrProjectNotFound()
	if err.Code != base.ErrNoProject {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrNoProject)
	}
	if !strings.Contains(err.Message, "Project root") {
		t.Errorf("Message should mention project root: %s", err.Message)
	}
}

func TestErrFeatureNotSet(t *testing.T) {
	err := base.ErrFeatureNotSet()
	if err.Code != base.ErrNoFeature {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrNoFeature)
	}
	if !strings.Contains(err.Hint, "forge feature") {
		t.Errorf("Hint should mention forge feature: %s", err.Hint)
	}
}

func TestErrTaskNotFound(t *testing.T) {
	err := base.ErrTaskNotFound("1.2.3")
	if err.Code != base.ErrNotFound {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrNotFound)
	}
	if !strings.Contains(err.Message, "1.2.3") {
		t.Errorf("Message should contain task ID: %s", err.Message)
	}
}

func TestErrNoInput(t *testing.T) {
	err := base.ErrNoInput("need a file path")
	if err.Code != base.ErrInvalidInput {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrInvalidInput)
	}
	if err.Cause != "need a file path" {
		t.Errorf("Cause = %q, want %q", err.Cause, "need a file path")
	}
}

func TestErrDependenciesNotMet(t *testing.T) {
	err := base.ErrDependenciesNotMet("2.1", []string{"1.1", "1.2"})
	if err.Code != base.ErrConflict {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrConflict)
	}
	if !strings.Contains(err.Cause, "1.1") || !strings.Contains(err.Cause, "1.2") {
		t.Errorf("Cause should list unmet deps: %s", err.Cause)
	}
}

func TestErrDataIntegrity(t *testing.T) {
	err := base.ErrDataIntegrity([]string{"state mismatch", "index corrupt"})
	if err.Code != base.ErrConflict {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrConflict)
	}
	if !strings.Contains(err.Cause, "state mismatch") {
		t.Errorf("Cause should include issues: %s", err.Cause)
	}
}

func TestErrFeatureNotFound(t *testing.T) {
	err := base.ErrFeatureNotFound("my-feature")
	if err.Code != base.ErrNotFound {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrNotFound)
	}
	if !strings.Contains(err.Message, "my-feature") {
		t.Errorf("Message should contain slug: %s", err.Message)
	}
}

func TestErrMissingFields(t *testing.T) {
	err := base.ErrMissingFields([]string{"summary", "taskId"})
	if err.Code != base.ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrValidation)
	}
	if !strings.Contains(err.Message, "summary") {
		t.Errorf("Message should list missing fields: %s", err.Message)
	}
}

func TestErrNoTestEvidence(t *testing.T) {
	err := base.ErrNoTestEvidence()
	if err.Code != base.ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrValidation)
	}
	if !strings.Contains(err.Cause, "testsPassed") {
		t.Errorf("Cause should mention test metrics: %s", err.Cause)
	}
}

func TestErrUnmetAcceptanceCriteria(t *testing.T) {
	err := base.ErrUnmetAcceptanceCriteria([]string{"login works", "data persists"})
	if err.Code != base.ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, base.ErrValidation)
	}
	if !strings.Contains(err.Cause, "login works") {
		t.Errorf("Cause should list unmet criteria: %s", err.Cause)
	}
}

func TestErrInvalidStatus(t *testing.T) {
	t.Run("with valid statuses", func(t *testing.T) {
		err := base.ErrInvalidStatus("bad", []string{"pending", "completed"})
		if err.Code != base.ErrValidation {
			t.Errorf("Code = %q, want %q", err.Code, base.ErrValidation)
		}
		if !strings.Contains(err.Cause, "pending") {
			t.Errorf("Cause should list valid statuses: %s", err.Cause)
		}
		if !strings.Contains(err.Action, "pending") {
			t.Errorf("Action should suggest first valid status: %s", err.Action)
		}
	})

	t.Run("empty valid statuses", func(t *testing.T) {
		err := base.ErrInvalidStatus("bad", []string{})
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
	assert.Equal(t, base.ErrorCode("INVALID_TRANSITION"), base.ErrInvalidTransition)
	assert.Equal(t, base.ErrorCode("INVALID_PATH"), base.ErrInvalidPath)
	assert.Equal(t, base.ErrorCode("EVAL_PARSE_FAILURE"), base.ErrEvalParseFailure)
	assert.Equal(t, base.ErrorCode("CONTRACT_UNVERIFIABLE"), base.ErrContractUnverifiable)
}

// --- New factory functions ---

func TestNewErrInvalidTransition(t *testing.T) {
	err := base.NewErrInvalidTransition("completed", "pending", "use reopen instead")
	assert.Equal(t, base.ErrInvalidTransition, err.Code)
	assert.Contains(t, err.Message, "completed")
	assert.Contains(t, err.Message, "pending")
	assert.NotEmpty(t, err.Cause)
	assert.Contains(t, err.Hint, "reopen")
	assert.NotEmpty(t, err.Action)
}

func TestNewErrInvalidPath(t *testing.T) {
	err := base.NewErrInvalidPath("../../../etc/passwd")
	assert.Equal(t, base.ErrInvalidPath, err.Code)
	assert.Contains(t, err.Message, "path")
	assert.Contains(t, err.Cause, "..")
	assert.NotEmpty(t, err.Hint)
	assert.NotEmpty(t, err.Action)
}

func TestNewErrEvalParseFailure(t *testing.T) {
	err := base.NewErrEvalParseFailure("score: not-a-number")
	assert.Equal(t, base.ErrEvalParseFailure, err.Code)
	assert.Contains(t, err.Message, "parse")
	assert.Contains(t, err.Cause, "score")
	assert.NotEmpty(t, err.Hint)
	assert.NotEmpty(t, err.Action)
}

func TestNewErrContractUnverifiable(t *testing.T) {
	err := base.NewErrContractUnverifiable("contracts/api.md")
	assert.Equal(t, base.ErrContractUnverifiable, err.Code)
	assert.Contains(t, err.Message, "unverifiable")
	assert.Contains(t, err.Cause, "contracts/api.md")
	assert.NotEmpty(t, err.Hint)
	assert.NotEmpty(t, err.Action)
}

// --- ExitCode method ---

func TestExitCode_BlockingErrors(t *testing.T) {
	blockingCodes := []base.ErrorCode{
		base.ErrInvalidTransition,
		base.ErrInvalidPath,
		base.ErrContractUnverifiable,
	}
	for _, code := range blockingCodes {
		err := &base.AIError{Code: code, Message: "test"}
		assert.Equal(t, 2, err.ExitCode(), "expected exit code 2 for %s", code)
	}
}

func TestExitCode_RetryableErrors(t *testing.T) {
	retryableCodes := []base.ErrorCode{
		base.ErrEvalParseFailure,
		base.ErrNoProject,
		base.ErrNoFeature,
		base.ErrNotFound,
		base.ErrConflict,
		base.ErrInvalidInput,
		base.ErrValidation,
	}
	for _, code := range retryableCodes {
		err := &base.AIError{Code: code, Message: "test"}
		assert.Equal(t, 1, err.ExitCode(), "expected exit code 1 for %s", code)
	}
}

func TestExitCode_DefaultRetryable(t *testing.T) {
	err := &base.AIError{Code: base.ErrorCode("UNKNOWN_CODE"), Message: "test"}
	assert.Equal(t, 1, err.ExitCode(), "unknown error codes should default to exit code 1")
}

// --- Existing error behavior unchanged ---

func TestExistingFactoryFunctionsUnchanged(t *testing.T) {
	tests := []struct {
		name     string
		err      *base.AIError
		wantCode base.ErrorCode
	}{
		{"ErrProjectNotFound", base.ErrProjectNotFound(), base.ErrNoProject},
		{"ErrFeatureNotSet", base.ErrFeatureNotSet(), base.ErrNoFeature},
		{"ErrTaskNotFound", base.ErrTaskNotFound("1.0"), base.ErrNotFound},
		{"ErrNoInput", base.ErrNoInput("detail"), base.ErrInvalidInput},
		{"ErrInvalidJSON", base.ErrInvalidJSON("file.json", "bad"), base.ErrValidation},
		{"ErrFileNotFound", base.ErrFileNotFound("x.txt"), base.ErrNotFound},
		{"ErrNoPendingTasks", base.ErrNoPendingTasks(), base.ErrNotFound},
		{"ErrDependenciesNotMet", base.ErrDependenciesNotMet("1.0", []string{"0.1"}), base.ErrConflict},
		{"ErrDataIntegrity", base.ErrDataIntegrity([]string{"x"}), base.ErrConflict},
		{"ErrInvalidStatus", base.ErrInvalidStatus("bad", []string{"ok"}), base.ErrValidation},
		{"ErrMissingFields", base.ErrMissingFields([]string{"x"}), base.ErrValidation},
		{"ErrFeatureNotFound", base.ErrFeatureNotFound("slug"), base.ErrNotFound},
		{"ErrNoTestEvidence", base.ErrNoTestEvidence(), base.ErrValidation},
		{"ErrUnmetAcceptanceCriteria", base.ErrUnmetAcceptanceCriteria([]string{"x"}), base.ErrValidation},
		{"ErrTaskIDConflict", base.ErrTaskIDConflict("1.0"), base.ErrConflict},
		{"ErrInvalidDependency", base.ErrInvalidDependency([]string{"9.9"}), base.ErrValidation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantCode, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message)
		})
	}
}
