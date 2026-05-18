package journey

import (
	"strings"
	"testing"

	"forge-cli/pkg/contract"
	"forge-cli/pkg/descriptor"
	"forge-cli/pkg/profile"
)

// --- Fixtures ---

func sampleContracts() []contract.Contract {
	return []contract.Contract{
		{
			Journey: "task-lifecycle",
			Step:    1,
			Action:  "forge feature my-feature",
			Outcomes: []contract.Outcome{
				{
					Name:          "success",
					Preconditions: "no feature with this slug exists",
					Input:         "feature-slug as positional arg",
					Output:        "success confirmation containing feature-slug",
					State:         "feature directory created with manifest.md",
				},
			},
			Invariants: []string{
				"feature_slug consistent across all steps",
				"task_id stable once assigned",
			},
		},
		{
			Journey: "task-lifecycle",
			Step:    2,
			Action:  "forge task claim",
			Outcomes: []contract.Outcome{
				{
					Name:          "success",
					Preconditions: "feature exists with at least one task available",
					Input:         "no args",
					Output:        "claimed task with identifier",
					State:         "task status changed to in_progress",
				},
				{
					Name:          "no-tasks",
					Preconditions: "no tasks available for claiming",
					Input:         "no args",
					Output:        "no tasks available error message",
					State:         "unchanged",
				},
			},
			Invariants: []string{
				"feature_slug consistent across all steps",
				"task_id stable once assigned",
			},
		},
	}
}

func sampleFacts() []descriptor.FactEntry {
	return []descriptor.FactEntry{
		{
			Key:    "CLI_FEATURE_CREATE_OUTPUT",
			Value:  "Feature my-feature created successfully",
			Source: "internal/cmd/feature.go:45",
		},
		{
			Key:    "CLI_TASK_CLAIM_OUTPUT",
			Value:  "claimed task <task_id>",
			Source: "internal/cmd/claim.go:42",
		},
		{
			Key:    "CLI_TASK_CLAIM_NO_TASKS",
			Value:  "no tasks available for claiming",
			Source: "internal/cmd/claim.go:50",
		},
	}
}

func sampleOpts() TestGenerationOpts {
	return TestGenerationOpts{
		Journey:   "task-lifecycle",
		Contracts: sampleContracts(),
		Framework: profile.FrameworkInfo{
			Name:                "go-testing",
			TestFunctionPattern: "func Test*",
			FilePattern:         "*_test.go",
			LanguageHint:        "go",
		},
		Facts: sampleFacts(),
	}
}

// --- Test: FeatureTag ---

func TestFeatureTag(t *testing.T) {
	tests := []struct {
		framework string
		expected  string
	}{
		{"go-testing", "//go:build feature"},
		{"pytest", "@pytest.mark.feature"},
		{"mocha", `describe("@feature", ...)`},
		{"junit5", `@Tag("feature")`},
		{"rust-test", `#[cfg(feature = "feature")]`},
		{"maestro", "# feature"},
		{"unknown", "// @feature"},
	}

	for _, tt := range tests {
		t.Run(tt.framework, func(t *testing.T) {
			fw := profile.FrameworkInfo{Name: tt.framework}
			got := FeatureTag(fw)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// --- Test: RegressionTag ---

func TestRegressionTag(t *testing.T) {
	tests := []struct {
		framework string
		expected  string
	}{
		{"go-testing", "//go:build regression"},
		{"pytest", "@pytest.mark.regression"},
		{"mocha", `describe("@regression", ...)`},
	}

	for _, tt := range tests {
		t.Run(tt.framework, func(t *testing.T) {
			fw := profile.FrameworkInfo{Name: tt.framework}
			got := RegressionTag(fw)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// --- Test: TestOutputDir ---

func TestTestOutputDir(t *testing.T) {
	got := TestOutputDir("/project", "task-lifecycle")
	if !strings.Contains(got, "tests") || !strings.Contains(got, "task-lifecycle") {
		t.Fatalf("expected path containing tests/task-lifecycle, got %q", got)
	}
}

// --- Test: SmokeTestName ---

func TestSmokeTestName(t *testing.T) {
	tests := []struct {
		journey  string
		expected string
	}{
		{"task-lifecycle", "TestJourneyTaskLifecycleSmoke"},
		{"session-diagnostics", "TestJourneySessionDiagnosticsSmoke"},
		{"simple", "TestJourneySimpleSmoke"},
	}

	for _, tt := range tests {
		t.Run(tt.journey, func(t *testing.T) {
			got := SmokeTestName(tt.journey)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// --- Test: GenerateGoTest ---

func TestGenerateGoTest(t *testing.T) {
	opts := sampleOpts()
	fw := profile.FrameworkInfo{Name: "go-testing"}

	t.Run("generates step tests and smoke test", func(t *testing.T) {
		tests, err := GenerateGoTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tests) < 3 {
			t.Fatalf("expected at least 3 generated tests (2 steps + 1 smoke), got %d", len(tests))
		}
	})

	t.Run("step tests contain @feature build tag", func(t *testing.T) {
		tests, err := GenerateGoTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, test := range tests {
			if !strings.Contains(test.Content, "//go:build feature") {
				t.Fatalf("test %q missing @feature build tag", test.Filename)
			}
		}
	})

	t.Run("step tests contain Contract Outcome assertions", func(t *testing.T) {
		tests, err := GenerateGoTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Find the step 2 test (has 2 Outcomes)
		var step2Test *GeneratedTest
		for i := range tests {
			if strings.Contains(tests[i].Filename, "step2") && !tests[i].IsSmokeTest {
				step2Test = &tests[i]
				break
			}
		}
		if step2Test == nil {
			t.Fatal("step 2 test not found")
		}

		// Should contain assertions for both "success" and "no-tasks" Outcomes
		if !strings.Contains(step2Test.Content, "success") {
			t.Fatal("step 2 test should reference 'success' Outcome")
		}
		if !strings.Contains(step2Test.Content, "no-tasks") {
			t.Fatal("step 2 test should reference 'no-tasks' Outcome")
		}
	})

	t.Run("smoke test contains happy path steps", func(t *testing.T) {
		tests, err := GenerateGoTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var smokeTest *GeneratedTest
		for i := range tests {
			if tests[i].IsSmokeTest {
				smokeTest = &tests[i]
				break
			}
		}
		if smokeTest == nil {
			t.Fatal("smoke test not found")
		}

		// Should contain step references
		if !strings.Contains(smokeTest.Content, "Step 1") {
			t.Fatal("smoke test should reference Step 1")
		}
		if !strings.Contains(smokeTest.Content, "Step 2") {
			t.Fatal("smoke test should reference Step 2")
		}
		// Should contain Journey Invariants
		if !strings.Contains(smokeTest.Content, "feature_slug consistent") {
			t.Fatal("smoke test should reference Journey Invariants")
		}
	})

	t.Run("smoke test name follows convention", func(t *testing.T) {
		tests, err := GenerateGoTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var smokeTest *GeneratedTest
		for i := range tests {
			if tests[i].IsSmokeTest {
				smokeTest = &tests[i]
				break
			}
		}
		if !strings.Contains(smokeTest.Content, "TestJourneyTaskLifecycleSmoke") {
			t.Fatalf("smoke test should contain function TestJourneyTaskLifecycleSmoke, got:\n%s", smokeTest.Content)
		}
	})

	t.Run("no hardcoded secrets in generated code", func(t *testing.T) {
		tests, err := GenerateGoTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, test := range tests {
			// Check for common secret patterns
			if strings.Contains(test.Content, "password") && !strings.Contains(test.Content, "<from-env>") {
				t.Fatalf("test %q contains hardcoded password", test.Filename)
			}
		}
	})

	t.Run("wrong framework returns error", func(t *testing.T) {
		wrongOpts := sampleOpts()
		wrongOpts.Framework = profile.FrameworkInfo{Name: "pytest"}
		_, err := GenerateGoTest(wrongOpts)
		if err == nil {
			t.Fatal("expected error for wrong framework")
		}
	})

	_ = fw
}

// --- Test: GeneratePythonTest ---

func TestGeneratePythonTest(t *testing.T) {
	opts := sampleOpts()
	opts.Framework = profile.FrameworkInfo{
		Name:                "pytest",
		TestFunctionPattern: "def test_*",
		FilePattern:         "test_*.py",
		LanguageHint:        "python",
	}

	t.Run("generates tests with pytest.mark.feature", func(t *testing.T) {
		tests, err := GeneratePythonTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, test := range tests {
			if !strings.Contains(test.Content, "pytest.mark.feature") {
				t.Fatalf("test %q missing pytest.mark.feature", test.Filename)
			}
		}
	})

	t.Run("generates smoke test", func(t *testing.T) {
		tests, err := GeneratePythonTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var smokeTest *GeneratedTest
		for i := range tests {
			if tests[i].IsSmokeTest {
				smokeTest = &tests[i]
				break
			}
		}
		if smokeTest == nil {
			t.Fatal("smoke test not found")
		}
		if !strings.Contains(smokeTest.Filename, "smoke") {
			t.Fatalf("smoke test filename should contain 'smoke', got %q", smokeTest.Filename)
		}
	})
}

// --- Test: GenerateJSTest ---

func TestGenerateJSTest(t *testing.T) {
	opts := sampleOpts()
	opts.Framework = profile.FrameworkInfo{
		Name:                "mocha",
		TestFunctionPattern: "describe/it",
		FilePattern:         "*.spec.ts",
		LanguageHint:        "javascript",
	}

	t.Run("generates tests with @feature describe", func(t *testing.T) {
		tests, err := GenerateJSTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, test := range tests {
			if !strings.Contains(test.Content, `describe("@feature"`) {
				t.Fatalf("test %q missing @feature describe block", test.Filename)
			}
		}
	})

	t.Run("generates smoke test with .spec.ts extension", func(t *testing.T) {
		tests, err := GenerateJSTest(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var smokeTest *GeneratedTest
		for i := range tests {
			if tests[i].IsSmokeTest {
				smokeTest = &tests[i]
				break
			}
		}
		if smokeTest == nil {
			t.Fatal("smoke test not found")
		}
		if !strings.HasSuffix(smokeTest.Filename, ".spec.ts") {
			t.Fatalf("JS test should have .spec.ts extension, got %q", smokeTest.Filename)
		}
	})
}

// --- Test: GenerateDispatched ---

func TestGenerateDispatched(t *testing.T) {
	t.Run("dispatches to Go generator", func(t *testing.T) {
		opts := sampleOpts()
		tests, err := GenerateDispatched(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tests) == 0 {
			t.Fatal("expected generated tests")
		}
	})

	t.Run("dispatches to Python generator", func(t *testing.T) {
		opts := sampleOpts()
		opts.Framework = profile.FrameworkInfo{Name: "pytest"}
		tests, err := GenerateDispatched(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tests) == 0 {
			t.Fatal("expected generated tests")
		}
	})

	t.Run("dispatches to JS generator", func(t *testing.T) {
		opts := sampleOpts()
		opts.Framework = profile.FrameworkInfo{Name: "mocha"}
		tests, err := GenerateDispatched(opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tests) == 0 {
			t.Fatal("expected generated tests")
		}
	})

	t.Run("unsupported framework returns error", func(t *testing.T) {
		opts := sampleOpts()
		opts.Framework = profile.FrameworkInfo{Name: "unknown-fw"}
		_, err := GenerateDispatched(opts)
		if err == nil {
			t.Fatal("expected error for unsupported framework")
		}
	})
}

// --- Test: sanitizeName ---

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"forge feature my-feature", "forge_feature_my_feature"},
		{"forge task claim", "forge_task_claim"},
		{"camelCase", "camelcase"},
		{"with/slash", "with_slash"},
		{"with.dot", "with_dot"},
		{"with:colon", "with_colon"},
		{"extra  spaces", "extra_spaces"},
		{"trailing_", "trailing"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeName(tt.input)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
