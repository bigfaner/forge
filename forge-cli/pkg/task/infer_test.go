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
		{"T-quick-2a", TypeTestPipelineGenAndRun},
		{"T-quick-2b", TypeTestPipelineGenAndRun},
		{"T-quick-3a", TypeTestPipelineGraduate},
		{"T-quick-3b", TypeTestPipelineGraduate},
		{"T-quick-4a", TypeTestPipelineVerifyRegression},
		{"T-quick-4b", TypeTestPipelineVerifyRegression},
		{"T-quick-5", TypeDocGenerationDrift},
		{"T-quick-5a", TypeDocGenerationDrift},
		{"T-quick-5b", TypeDocGenerationDrift},

		// Fix tasks
		{"fix-1", TypeFix},
		{"fix-2", TypeFix},
		{"disc-1", TypeFix},

		// Doc evaluation task
		{"T-eval-doc", TypeDocEvaluation},

		// Type-suffixed test tasks (per-type split)
		{"T-test-2-api", TypeTestPipelineGenScripts},
		{"T-test-2-tui", TypeTestPipelineGenScripts},
		{"T-test-2-cli", TypeTestPipelineGenScripts},
		{"T-test-2-web-ui", TypeTestPipelineGenScripts},
		{"T-test-3-api", TypeTestPipelineRun},
		{"T-test-4-api", TypeTestPipelineGraduate},
		{"T-quick-2-api", TypeTestPipelineGenAndRun},
		{"T-quick-2-tui", TypeTestPipelineGenAndRun},
		{"T-quick-3-cli", TypeTestPipelineGraduate},
		{"T-quick-4-api", ""},

		// Profile-suffixed + type-suffixed
		{"T-test-2a-api", TypeTestPipelineGenScripts},
		{"T-test-2b-tui", TypeTestPipelineGenScripts},
		{"T-test-3a-cli", TypeTestPipelineRun},
		{"T-quick-2a-api", TypeTestPipelineGenAndRun},
		{"T-quick-2b-tui", TypeTestPipelineGenAndRun},
		{"T-quick-3a-cli", TypeTestPipelineGraduate},

		// Type suffix on T-test-1 should NOT match (exact + profileSuffixed only)
		{"T-test-1-api", ""},
		{"T-test-4.5-api", ""},
		{"T-test-5-api", ""},

		// Unknown IDs return empty string (no TypeFeature fallback)
		{"1.1", ""},
		{"2.3", ""},
		{"", ""},
		{"random-task", ""},

		// New business types are explicit — not inferred from patterns
		{"feature", ""},
		{"enhancement", ""},
		{"cleanup", ""},
		{"refactor", ""},

		// Renamed and new IDs
		{"T-specs-1", TypeDocGenerationConsolidate},
		{"T-quick-specs-1", TypeDocGenerationDrift},
		{"T-clean-code-1", TypeCleanCode},
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

func TestExtractTypeSuffix(t *testing.T) {
	tests := []struct {
		id   string
		base string
		want string
	}{
		{"T-test-2-api", "T-test-2", "api"},
		{"T-test-2a-tui", "T-test-2", "tui"},
		{"T-test-2-web-ui", "T-test-2", "web-ui"},
		{"T-test-2a-web-ui", "T-test-2", "web-ui"},
		{"T-quick-2-cli", "T-quick-2", "cli"},
		{"T-quick-2b-api", "T-quick-2", "api"},
		{"T-test-2", "T-test-2", ""},        // exact match
		{"T-test-2a", "T-test-2", ""},       // profile suffix only
		{"T-test-1-api", "T-test-1", "api"}, // syntactically valid, but InferType won't route it
		{"random", "T-test-2", ""},          // wrong base
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := ExtractTypeSuffix(tt.id, tt.base)
			if got != tt.want {
				t.Errorf("ExtractTypeSuffix(%q, %q) = %q, want %q", tt.id, tt.base, got, tt.want)
			}
		})
	}
}
