package task

import "testing"

func TestInferType(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		// Summary and gate
		{"1.summary", TypeDocSummary},
		{"2.summary", TypeDocSummary},
		{"1.gate", TypeGate},
		{"2.gate", TypeGate},

		// Breakdown test tasks (exact match)
		{"T-test-1", TypeTestGenCases},
		{"T-test-2", TypeTestGenScripts},
		{"T-test-3", TypeTestRun},
		{"T-test-4", TypeTestGraduate},
		{"T-test-4.5", TypeTestVerifyRegression},
		{"T-test-5", TypeDocConsolidate},

		// T-test-1b is a special case: exact match for eval-cases (NOT a profile suffix)
		// It's checked before profileSuffixedID, so it returns TypeTestEvalCases.
		{"T-test-1b", TypeTestEvalCases},

		// Breakdown test tasks (profile-suffixed)
		{"T-test-1a", TypeTestGenCases},
		{"T-test-1c", TypeTestGenCases},
		{"T-test-2a", TypeTestGenScripts},
		{"T-test-2b", TypeTestGenScripts},
		{"T-test-3a", TypeTestRun},
		{"T-test-3b", TypeTestRun},
		{"T-test-4a", TypeTestGraduate},
		{"T-test-4b", TypeTestGraduate},

		// Quick test tasks (profile-suffixed)
		{"T-quick-1a", TypeTestGenCases},
		{"T-quick-1b", TypeTestGenCases},
		{"T-quick-2a", TypeTestGenAndRun},
		{"T-quick-2b", TypeTestGenAndRun},
		{"T-quick-3a", TypeTestGraduate},
		{"T-quick-3b", TypeTestGraduate},
		{"T-quick-4a", TypeTestVerifyRegression},
		{"T-quick-4b", TypeTestVerifyRegression},
		{"T-quick-5", TypeDocDrift},
		{"T-quick-5a", TypeDocDrift},
		{"T-quick-5b", TypeDocDrift},

		// Fix tasks
		{"fix-1", TypeCodingFix},
		{"fix-2", TypeCodingFix},
		{"disc-1", TypeCodingFix},

		// Doc evaluation task
		{"T-eval-doc", TypeDocEval},

		// Type-suffixed test tasks (per-type split)
		{"T-test-2-api", TypeTestGenScripts},
		{"T-test-2-tui", TypeTestGenScripts},
		{"T-test-2-cli", TypeTestGenScripts},
		{"T-test-2-web-ui", TypeTestGenScripts},
		{"T-test-3-api", TypeTestRun},
		{"T-test-4-api", TypeTestGraduate},
		{"T-quick-2-api", TypeTestGenAndRun},
		{"T-quick-2-tui", TypeTestGenAndRun},
		{"T-quick-3-cli", TypeTestGraduate},
		{"T-quick-4-api", ""},

		// Profile-suffixed + type-suffixed
		{"T-test-2a-api", TypeTestGenScripts},
		{"T-test-2b-tui", TypeTestGenScripts},
		{"T-test-3a-cli", TypeTestRun},
		{"T-quick-2a-api", TypeTestGenAndRun},
		{"T-quick-2b-tui", TypeTestGenAndRun},
		{"T-quick-3a-cli", TypeTestGraduate},

		// Type suffix on T-test-1 should NOT match (exact + profileSuffixed only)
		{"T-test-1-api", ""},
		{"T-test-4.5-api", ""},
		{"T-test-5-api", ""},

		// Unknown IDs return empty string (no TypeCodingFeature fallback)
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
		{"T-specs-1", TypeDocConsolidate},
		{"T-quick-specs-1", TypeDocDrift},
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
