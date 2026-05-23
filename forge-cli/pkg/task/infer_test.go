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
		{"T-test-gen-scripts", TypeTestGenScripts},
		{"T-test-run", TypeTestRun},
		{"T-test-verify-regression", TypeTestVerifyRegression},

		// Quick test tasks (exact match)
		{"T-quick-gen-and-run", TypeTestGenAndRun},
		{"T-quick-verify-regression", TypeTestVerifyRegression},
		{"T-quick-doc-drift", TypeDocDrift},

		// Fix tasks
		{"fix-1", TypeCodingFix},
		{"fix-2", TypeCodingFix},
		{"disc-1", TypeCodingFix},

		// Doc review task
		{"T-review-doc", TypeDocReview},

		// Validation tasks
		{"T-validate-code", TypeValidationCode},
		{"T-validate-ux", TypeValidationUx},

		// Type-suffixed test tasks (per-type split, no profile letter)
		{"T-test-gen-scripts-api", TypeTestGenScripts},
		{"T-test-gen-scripts-tui", TypeTestGenScripts},
		{"T-test-gen-scripts-cli", TypeTestGenScripts},
		{"T-test-gen-scripts-web-ui", TypeTestGenScripts},
		{"T-test-gen-journeys-api", TypeTestGenJourneys},
		{"T-test-gen-journeys-tui", TypeTestGenJourneys},
		{"T-test-gen-journeys-cli", TypeTestGenJourneys},
		{"T-quick-gen-and-run-api", TypeTestGenAndRun},
		{"T-quick-gen-and-run-tui", TypeTestGenAndRun},

		// Gen-contracts exact match
		{"T-test-gen-contracts", TypeTestGenContracts},
		{"T-test-gen-journeys", TypeTestGenJourneys},

		// Type suffix on tasks that don't support it should NOT match
		{"T-test-verify-regression-api", ""},
		{"T-specs-consolidate-api", ""},
		{"T-test-run-api", ""},
		{"T-quick-verify-regression-api", ""},

		// Unknown IDs return empty string
		{"1.1", ""},
		{"2.3", ""},
		{"", ""},
		{"random-task", ""},

		// New business types are explicit — not inferred from patterns
		{"feature", ""},
		{"enhancement", ""},
		{"cleanup", ""},
		{"refactor", ""},

		// Other IDs
		{"T-specs-consolidate", TypeDocConsolidate},
		{"T-clean-code", TypeCleanCode},

		// Old profile-suffixed IDs no longer match
		{"T-test-gen-scriptsa", ""},
		{"T-test-runa", ""},
		{"T-quick-verify-regressiona", ""},
		{"T-quick-doc-drifta", ""},

		// Hard Rule: gen-journeys type suffix with hyphenated types
		{"T-test-gen-journeys-web-ui", TypeTestGenJourneys},
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
		{"T-test-gen-scripts-tui", "T-test-gen-scripts", "tui"},
		{"T-test-gen-scripts-web-ui", "T-test-gen-scripts", "web-ui"},
		{"T-quick-gen-and-run-cli", "T-quick-gen-and-run", "cli"},
		{"T-quick-gen-and-run-api", "T-quick-gen-and-run", "api"},
		{"T-test-gen-scripts", "T-test-gen-scripts", ""}, // exact match, no type suffix
		{"T-test-gen-journeys-api", "T-test-gen-journeys", "api"},
		{"T-test-gen-journeys-tui", "T-test-gen-journeys", "tui"},
		{"random", "T-test-gen-scripts", ""}, // wrong base

		// Old profile-suffixed IDs no longer extract correctly (profile letter is part of suffix)
		{"T-test-gen-scriptsa-tui", "T-test-gen-scripts", ""}, // 'a' before '-' is not a valid type start
		{"T-test-gen-scriptsa-api", "T-test-gen-scripts", ""}, // 'a' before '-' is not valid
		{"T-test-gen-scriptsa", "T-test-gen-scripts", ""},     // profile suffix only, no type
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
