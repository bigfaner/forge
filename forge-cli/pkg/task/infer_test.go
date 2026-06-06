package task

import "testing"

func TestInferType(t *testing.T) {
	// Surfaces map used for T-test-run-{key} prefix matching tests.
	multiSurfaces := map[string]string{
		"backend":  "api",
		"frontend": "web",
	}
	singleSurface := map[string]string{
		".": "api",
	}

	tests := []struct {
		id       string
		surfaces map[string]string
		want     string
	}{
		// Summary and gate
		{"1.summary", nil, TypeDocSummary},
		{"2.summary", nil, TypeDocSummary},
		{"1.gate", nil, TypeGate},
		{"2.gate", nil, TypeGate},

		// Breakdown test tasks (exact match, nil surfaces)
		{"T-test-gen-scripts", nil, TypeTestGenScripts},
		{"T-test-run", nil, TypeTestRun},

		// Quick test tasks (exact match)
		{"T-quick-doc-drift", nil, TypeDocDrift},

		// Fix tasks
		{"fix-1", nil, TypeCodingFix},
		{"fix-2", nil, TypeCodingFix},
		{"disc-1", nil, TypeCodingFix},

		// Doc fix tasks (doc-fix- prefix must match before fix- prefix)
		{"doc-fix-1", nil, TypeDocFix},
		{"doc-fix-2", nil, TypeDocFix},

		// Doc review task
		{"T-review-doc", nil, TypeDocReview},

		// Validation tasks
		{"T-validate-code", nil, TypeValidationCode},
		{"T-validate-ux", nil, TypeValidationUx},

		// Type-suffixed test tasks (now per-surface-key, require surfaces map for suffix matching)
		{"T-test-gen-scripts-api", map[string]string{"api": "api"}, TypeTestGenScripts},
		{"T-test-gen-scripts-tui", map[string]string{"tui": "tui"}, TypeTestGenScripts},
		{"T-test-gen-scripts-cli", map[string]string{"cli": "cli"}, TypeTestGenScripts},
		{"T-test-gen-scripts-web-ui", map[string]string{"web-ui": "web"}, TypeTestGenScripts},
		// Without surfaces map, per-surface-key suffixes don't match
		{"T-test-gen-scripts-api", nil, ""},
		{"T-test-gen-scripts-tui", nil, ""},
		// T-test-gen-journeys suffixed IDs no longer match (gen-journeys is not expanded per-surface-type in registry)
		{"T-test-gen-journeys-api", nil, ""},
		{"T-test-gen-journeys-tui", nil, ""},
		{"T-test-gen-journeys-cli", nil, ""},

		// Gen-contracts exact match
		{"T-test-gen-contracts", nil, TypeTestGenContracts},
		{"T-test-gen-journeys", nil, TypeTestGenJourneys},

		// Type suffix on tasks that don't support it should NOT match
		{"T-specs-consolidate-api", nil, ""},

		// T-test-run-api without surfaces map: no prefix matching, falls through
		{"T-test-run-api", nil, ""},

		// Unknown IDs return empty string
		{"1.1", nil, ""},
		{"2.3", nil, ""},
		{"", nil, ""},
		{"random-task", nil, ""},

		// New business types are explicit — not inferred from patterns
		{"feature", nil, ""},
		{"enhancement", nil, ""},
		{"cleanup", nil, ""},
		{"refactor", nil, ""},

		// Other IDs
		{"T-specs-consolidate", nil, TypeDocConsolidate},
		{"T-clean-code", nil, TypeCleanCode},

		// Old profile-suffixed IDs no longer match
		{"T-test-gen-scriptsa", nil, ""},
		{"T-test-runa", nil, ""},
		{"T-quick-doc-drifta", nil, ""},

		// T-test-gen-journeys-web-ui no longer matches (gen-journeys is not expanded per-surface-type)
		{"T-test-gen-journeys-web-ui", nil, ""},

		// Surface-key prefix matching: known key -> TypeTestRun
		{"T-test-run-backend", multiSurfaces, TypeTestRun},
		{"T-test-run-frontend", multiSurfaces, TypeTestRun},

		// Surface-key prefix matching: unknown key -> fallback (no match)
		{"T-test-run-unknown", multiSurfaces, ""},
		{"T-test-run-api", multiSurfaces, ""}, // "api" is a type, not a key in multiSurfaces

		// Single surface project: T-test-run (no suffix) still works
		{"T-test-run", singleSurface, TypeTestRun},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := InferType(tt.id, tt.surfaces)
			if got != tt.want {
				t.Errorf("InferType(%q, %v) = %q, want %q", tt.id, tt.surfaces, got, tt.want)
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
