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
		{"T-test-gen-cases", TypeTestGenCases},
		{"T-test-gen-scripts", TypeTestGenScripts},
		{"T-test-run", TypeTestRun},
		{"T-test-graduate", TypeTestGraduate},
		{"T-test-verify-regression", TypeTestVerifyRegression},
		{"T-test-eval-cases", TypeTestEvalCases},

		// Breakdown test tasks (profile-suffixed)
		{"T-test-gen-casesa", TypeTestGenCases},
		{"T-test-gen-casesc", TypeTestGenCases},
		{"T-test-gen-scriptsa", TypeTestGenScripts},
		{"T-test-gen-scripts", TypeTestGenScripts},
		{"T-test-runa", TypeTestRun},
		{"T-test-runb", TypeTestRun},
		{"T-test-graduatea", TypeTestGraduate},
		{"T-test-graduateb", TypeTestGraduate},

		// Quick test tasks (profile-suffixed)
		{"T-quick-gen-casesa", TypeTestGenCases},
		{"T-quick-gen-casesb", TypeTestGenCases},
		{"T-quick-gen-and-runa", TypeTestGenAndRun},
		{"T-quick-gen-and-runb", TypeTestGenAndRun},
		{"T-quick-graduatea", TypeTestGraduate},
		{"T-quick-graduateb", TypeTestGraduate},
		{"T-quick-verify-regressiona", TypeTestVerifyRegression},
		{"T-quick-verify-regressionb", TypeTestVerifyRegression},
		{"T-quick-doc-drift", TypeDocDrift},
		{"T-quick-doc-drifta", TypeDocDrift},
		{"T-quick-doc-driftb", TypeDocDrift},

		// Fix tasks
		{"fix-1", TypeCodingFix},
		{"fix-2", TypeCodingFix},
		{"disc-1", TypeCodingFix},

		// Doc evaluation task
		{"T-eval-doc", TypeDocEval},

		// Validation tasks
		{"T-validate-code", TypeValidationCode},
		{"T-validate-ux", TypeValidationUx},

		// Type-suffixed test tasks (per-type split)
		{"T-test-gen-scripts-api", TypeTestGenScripts},
		{"T-test-gen-scripts-tui", TypeTestGenScripts},
		{"T-test-gen-scripts-cli", TypeTestGenScripts},
		{"T-test-gen-scripts-web-ui", TypeTestGenScripts},
		{"T-test-run-api", TypeTestRun},
		{"T-test-graduate-api", TypeTestGraduate},
		{"T-quick-gen-and-run-api", TypeTestGenAndRun},
		{"T-quick-gen-and-run-tui", TypeTestGenAndRun},
		{"T-quick-graduate-cli", TypeTestGraduate},
		{"T-quick-verify-regression-api", ""},

		// Profile-suffixed + type-suffixed
		{"T-test-gen-scriptsa-api", TypeTestGenScripts},
		{"T-test-gen-scriptsb-tui", TypeTestGenScripts},
		{"T-test-runa-cli", TypeTestRun},
		{"T-quick-gen-and-runa-api", TypeTestGenAndRun},
		{"T-quick-gen-and-runb-tui", TypeTestGenAndRun},
		{"T-quick-graduatea-cli", TypeTestGraduate},

		// Type suffix on T-test-gen-cases should NOT match (exact + profileSuffixed only)
		{"T-test-gen-cases-api", ""},
		{"T-test-verify-regression-api", ""},
		{"T-specs-consolidate-api", ""},

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
		{"T-specs-consolidate", TypeDocConsolidate},
		{"T-quick-doc-drift", TypeDocDrift},
		{"T-clean-code", TypeCleanCode},
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
		{"T-test-gen-scripts-api", "T-test-gen-scripts", "api"},
		{"T-test-gen-scriptsa-tui", "T-test-gen-scripts", "tui"},
		{"T-test-gen-scripts-web-ui", "T-test-gen-scripts", "web-ui"},
		{"T-test-gen-scriptsa-web-ui", "T-test-gen-scripts", "web-ui"},
		{"T-quick-gen-and-run-cli", "T-quick-gen-and-run", "cli"},
		{"T-quick-gen-and-runb-api", "T-quick-gen-and-run", "api"},
		{"T-test-gen-scripts", "T-test-gen-scripts", ""},    // exact match
		{"T-test-gen-scriptsa", "T-test-gen-scripts", ""},   // profile suffix only
		{"T-test-gen-cases-api", "T-test-gen-cases", "api"}, // syntactically valid, but InferType won't route it
		{"random", "T-test-gen-scripts", ""},                // wrong base
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
