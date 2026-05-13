package task

import "testing"

func TestInferType(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		// Summary and gate
		{"1.summary", TypeDocGenerationSummary},
		{"2.summary", TypeDocGenerationSummary},
		{"1.gate", TypeGate},
		{"2.gate", TypeGate},

		// Breakdown test tasks (exact match)
		{"T-test-1", TypeTestPipelineGenCases},
		{"T-test-2", TypeTestPipelineGenScripts},
		{"T-test-3", TypeTestPipelineRun},
		{"T-test-4", TypeTestPipelineGraduate},
		{"T-test-4.5", TypeTestPipelineVerifyRegression},
		{"T-test-5", TypeDocGenerationConsolidate},

		// T-test-1b is a special case: exact match for eval-cases (NOT a profile suffix)
		// It's checked before profileSuffixedID, so it returns TypeTestPipelineEvalCases.
		{"T-test-1b", TypeTestPipelineEvalCases},

		// Breakdown test tasks (profile-suffixed)
		{"T-test-1a", TypeTestPipelineGenCases},
		{"T-test-1c", TypeTestPipelineGenCases},
		{"T-test-2a", TypeTestPipelineGenScripts},
		{"T-test-2b", TypeTestPipelineGenScripts},
		{"T-test-3a", TypeTestPipelineRun},
		{"T-test-3b", TypeTestPipelineRun},
		{"T-test-4a", TypeTestPipelineGraduate},
		{"T-test-4b", TypeTestPipelineGraduate},

		// Quick test tasks (profile-suffixed)
		{"T-quick-1a", TypeTestPipelineGenCases},
		{"T-quick-1b", TypeTestPipelineGenCases},
		{"T-quick-2a", TypeTestPipelineGenScripts},
		{"T-quick-2b", TypeTestPipelineGenScripts},
		{"T-quick-3a", TypeTestPipelineRun},
		{"T-quick-3b", TypeTestPipelineRun},
		{"T-quick-4a", TypeTestPipelineGraduate},
		{"T-quick-4b", TypeTestPipelineGraduate},
		{"T-quick-5a", TypeTestPipelineVerifyRegression},
		{"T-quick-5b", TypeTestPipelineVerifyRegression},

		// Fix tasks
		{"fix-1", TypeFix},
		{"fix-2", TypeFix},
		{"disc-1", TypeFix},

		// Default: implementation
		{"1.1", TypeImplementation},
		{"2.3", TypeImplementation},
		{"", TypeImplementation},
		{"random-task", TypeImplementation},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := InferType(tt.id)
			if got != tt.want {
				t.Errorf("InferType(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}
