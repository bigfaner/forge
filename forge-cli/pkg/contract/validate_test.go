package contract

import (
	"strings"
	"testing"
)

// --- Test helpers ---

func validOutcome(name string) Outcome {
	return Outcome{
		Name:          name,
		Preconditions: "feature exists with slug matching arg",
		Input:         "feature-slug as positional arg",
		Output:        "success confirmation containing feature-slug",
		State:         "feature directory created with manifest.md",
		SideEffect:    "",
		Invariants:    "",
	}
}

func validContract() Contract {
	return Contract{
		Journey: "task-lifecycle",
		Step:    1,
		Action:  "forge feature my-feature",
		Outcomes: []Outcome{
			validOutcome("success"),
		},
		Invariants: []string{
			"feature_slug consistent across all steps",
			"task_id stable once assigned",
		},
	}
}

// --- Test: Mandatory dimensions must be non-empty ---

func TestValidate_MandatoryDimensionsNonEmpty(t *testing.T) {
	t.Run("all mandatory dimensions present passes", func(t *testing.T) {
		errs := Validate(validContract())
		assertNoErrors(t, errs)
	})

	t.Run("missing Preconditions fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].Preconditions = ""
		errs := Validate(c)
		assertHasDimensionError(t, errs, DimensionPreconditions)
	})

	t.Run("missing Input fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].Input = ""
		errs := Validate(c)
		assertHasDimensionError(t, errs, DimensionInput)
	})

	t.Run("missing Output fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].Output = ""
		errs := Validate(c)
		assertHasDimensionError(t, errs, DimensionOutput)
	})

	t.Run("missing State fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].State = ""
		errs := Validate(c)
		assertHasDimensionError(t, errs, DimensionState)
	})

	t.Run("whitespace-only Preconditions fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].Preconditions = "   "
		errs := Validate(c)
		assertHasDimensionError(t, errs, DimensionPreconditions)
	})
}

// --- Test: Semantic descriptor purity (no regex) ---

func TestValidate_SemanticDescriptorPurity(t *testing.T) {
	regexPatterns := []struct {
		name    string
		content string
	}{
		{"backslash-d", "output matches \\d+ tasks"},
		{"dot-star", "output matches .* successfully"},
		{"character-class", "output matches [a-z]+ task"},
		{"non-capturing group", "output matches (?:task_)(\\d+)"},
		{"backslash-s", "output matches \\s+ separator"},
		{"backslash-w", "output matches \\w+ identifier"},
		{"backslash-b", "output matches \\bword\\b"},
		{"caret anchor", "^output starts here"},
		{"dollar anchor", "output ends here$"},
	}

	for _, tt := range regexPatterns {
		t.Run("Preconditions with "+tt.name+" fails", func(t *testing.T) {
			c := validContract()
			c.Outcomes[0].Preconditions = tt.content
			errs := Validate(c)
			assertHasSemanticError(t, errs, DimensionPreconditions)
		})
		t.Run("Output with "+tt.name+" fails", func(t *testing.T) {
			c := validContract()
			c.Outcomes[0].Output = tt.content
			errs := Validate(c)
			assertHasSemanticError(t, errs, DimensionOutput)
		})
	}

	t.Run("natural language passes", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].Output = "success confirmation containing feature-slug"
		errs := Validate(c)
		assertNoErrors(t, errs)
	})
}

// --- Test: Outcome Preconditions mutual exclusivity ---

func TestValidate_OutcomePreconditionsMutualExclusivity(t *testing.T) {
	t.Run("distinct Preconditions passes", func(t *testing.T) {
		c := validContract()
		c.Outcomes = append(c.Outcomes, Outcome{
			Name:          "no-feature",
			Preconditions: "no feature exists with given slug",
			Input:         "feature-slug as positional arg",
			Output:        "error message about missing feature",
			State:         "unchanged",
		})
		errs := Validate(c)
		assertNoErrors(t, errs)
	})

	t.Run("identical Preconditions fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes = append(c.Outcomes, Outcome{
			Name:          "duplicate",
			Preconditions: "feature exists with slug matching arg", // same as "success"
			Input:         "feature-slug as positional arg",
			Output:        "some other output",
			State:         "unchanged",
		})
		errs := Validate(c)
		assertHasMutualExclusivityError(t, errs)
	})

	t.Run("overlapping Preconditions fails", func(t *testing.T) {
		c := validContract()
		// "feature exists with slug matching arg" and "feature exists with slug matching arg and tasks available"
		// are overlapping because the first is a subset of the second
		c.Outcomes = append(c.Outcomes, Outcome{
			Name:          "overlapping",
			Preconditions: "feature exists with slug matching arg and tasks available",
			Input:         "feature-slug as positional arg",
			Output:        "success with tasks",
			State:         "feature directory with tasks",
		})
		errs := Validate(c)
		assertHasMutualExclusivityError(t, errs)
	})
}

// --- Test: Journey Invariants presence ---

func TestValidate_JourneyInvariantsPresence(t *testing.T) {
	t.Run("at least one invariant passes", func(t *testing.T) {
		c := validContract()
		c.Invariants = []string{"feature_slug consistent across all steps"}
		errs := Validate(c)
		assertNoErrors(t, errs)
	})

	t.Run("no invariants fails", func(t *testing.T) {
		c := validContract()
		c.Invariants = nil
		errs := Validate(c)
		assertHasJourneyInvariantError(t, errs)
	})

	t.Run("empty invariants list fails", func(t *testing.T) {
		c := validContract()
		c.Invariants = []string{}
		errs := Validate(c)
		assertHasJourneyInvariantError(t, errs)
	})
}

// --- Test: Outcome name uniqueness ---

func TestValidate_OutcomeNameUniqueness(t *testing.T) {
	t.Run("unique names passes", func(t *testing.T) {
		c := validContract()
		c.Outcomes = append(c.Outcomes, Outcome{
			Name:          "not-found",
			Preconditions: "no feature exists",
			Input:         "feature-slug",
			Output:        "error",
			State:         "unchanged",
		})
		errs := Validate(c)
		assertNoErrors(t, errs)
	})

	t.Run("duplicate names fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes = append(c.Outcomes, Outcome{
			Name:          "success", // duplicate
			Preconditions: "different preconditions here",
			Input:         "different input",
			Output:        "different output",
			State:         "different state",
		})
		errs := Validate(c)
		assertHasOutcomeNameError(t, errs)
	})
}

// --- Test: Outcome count checkpoint (> 5 Outcomes) ---

func TestValidate_OutcomeCountCheckpoint(t *testing.T) {
	t.Run("5 or fewer Outcomes passes", func(t *testing.T) {
		c := validContract()
		// Use truly exclusive Preconditions (different status values)
		extraPreconditions := []string{
			"no feature exists",
			"feature exists but empty",
			"feature exists with invalid slug",
			"feature exists with duplicate slug",
		}
		for i, pre := range extraPreconditions {
			c.Outcomes = append(c.Outcomes, Outcome{
				Name:          "outcome-" + string(rune('a'+i)),
				Preconditions: pre,
				Input:         "input",
				Output:        "output",
				State:         "state",
			})
		}
		errs := Validate(c)
		assertNoErrors(t, errs)
	})

	t.Run("more than 5 Outcomes triggers warning", func(t *testing.T) {
		c := validContract()
		// Use truly exclusive Preconditions
		extraPreconditions := []string{
			"no feature exists",
			"feature exists but empty",
			"feature exists with invalid slug",
			"feature exists with duplicate slug",
			"feature exists but permission denied",
			"feature exists with corrupted index",
		}
		for i, pre := range extraPreconditions {
			c.Outcomes = append(c.Outcomes, Outcome{
				Name:          "outcome-" + string(rune('a'+i)),
				Preconditions: pre,
				Input:         "input",
				Output:        "output",
				State:         "state",
			})
		}
		// 1 original + 6 added = 7 outcomes
		errs := Validate(c)
		assertHasOutcomeCountWarning(t, errs)
	})
}

// --- Test: Side-effect default ---

func TestValidate_SideEffectDefault(t *testing.T) {
	t.Run("empty Side-effect is acceptable", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].SideEffect = ""
		errs := Validate(c)
		assertNoErrors(t, errs)
	})

	t.Run("explicit none Side-effect is acceptable", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].SideEffect = "none"
		errs := Validate(c)
		assertNoErrors(t, errs)
	})

	t.Run("populated Side-effect is acceptable", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].SideEffect = "hook post-submit triggered with task_id"
		errs := Validate(c)
		assertNoErrors(t, errs)
	})
}

// --- Test: Batch splitting trigger ---

func TestBatchNeeded(t *testing.T) {
	t.Run("fewer than 15 Contracts does not need batching", func(t *testing.T) {
		contracts := make([]Contract, 14)
		for i := range contracts {
			contracts[i] = validContract()
			contracts[i].Step = i + 1
		}
		if BatchNeeded(contracts) {
			t.Fatal("expected no batching for 14 contracts")
		}
	})

	t.Run("15 or more Contracts needs batching", func(t *testing.T) {
		contracts := make([]Contract, 15)
		for i := range contracts {
			contracts[i] = validContract()
			contracts[i].Step = i + 1
		}
		if !BatchNeeded(contracts) {
			t.Fatal("expected batching for 15 contracts")
		}
	})

	t.Run("empty Contracts does not need batching", func(t *testing.T) {
		if BatchNeeded(nil) {
			t.Fatal("expected no batching for empty contracts")
		}
	})
}

// --- Test: Batch splitting strategy ---

func TestBatchSplit(t *testing.T) {
	t.Run("happy path batch contains success Outcomes", func(t *testing.T) {
		contracts := make([]Contract, 4)
		for i := range contracts {
			contracts[i] = validContract()
			contracts[i].Step = i + 1
			contracts[i].Outcomes = []Outcome{
				validOutcome("success"),
				{
					Name:          "error",
					Preconditions: "unique error precondition " + strings.Repeat("z", i),
					Input:         "input",
					Output:        "error output",
					State:         "unchanged",
				},
			}
		}
		batches := BatchSplit(contracts)
		if len(batches) < 1 {
			t.Fatal("expected at least one batch")
		}
		// First batch should contain only "success" Outcomes
		for _, c := range batches[0] {
			for _, o := range c.Outcomes {
				if o.Name != "success" {
					t.Fatalf("happy path batch should only contain success Outcomes, got %q", o.Name)
				}
			}
		}
	})
}

// --- Test: ContainsRegex ---

func TestContainsRegex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantReg bool
	}{
		{"plain text", "success confirmation", false},
		{"backslash-d", "matches \\d+", true},
		{"backslash-w", "contains \\w+", true},
		{"backslash-s", "has \\s+", true},
		{"backslash-b", "\\bword\\b", true},
		{"dot-star", ".* pattern", true},
		{"char class", "[a-z]+", true},
		{"non-capturing group", "(?:pattern)", true},
		{"caret anchor", "^start", true},
		{"dollar anchor", "end$", true},
		{"escaped backslash", "path\\\\to\\\\file", false},
		{"natural language", "success confirmation containing feature-slug", false},
		{"simple quotes", `"hello world"`, false},
		{"angle brackets", "<task_id>", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsRegex(tt.input)
			if got != tt.wantReg {
				t.Fatalf("ContainsRegex(%q) = %v, want %v", tt.input, got, tt.wantReg)
			}
		})
	}
}

// --- Test: Preconditions mutual exclusivity check ---

func TestArePreconditionsMutuallyExclusive(t *testing.T) {
	t.Run("single Outcome always exclusive", func(t *testing.T) {
		outcomes := []Outcome{validOutcome("success")}
		if !ArePreconditionsMutuallyExclusive(outcomes) {
			t.Fatal("single Outcome should always be mutually exclusive")
		}
	})

	t.Run("distinct Preconditions are exclusive", func(t *testing.T) {
		outcomes := []Outcome{
			{Name: "success", Preconditions: "task status is in_progress"},
			{Name: "not-claimed", Preconditions: "no task claimed"},
		}
		if !ArePreconditionsMutuallyExclusive(outcomes) {
			t.Fatal("distinct Preconditions should be mutually exclusive")
		}
	})

	t.Run("identical Preconditions are not exclusive", func(t *testing.T) {
		outcomes := []Outcome{
			{Name: "a", Preconditions: "task status is in_progress"},
			{Name: "b", Preconditions: "task status is in_progress"},
		}
		if ArePreconditionsMutuallyExclusive(outcomes) {
			t.Fatal("identical Preconditions should not be mutually exclusive")
		}
	})

	t.Run("subset Preconditions are not exclusive", func(t *testing.T) {
		outcomes := []Outcome{
			{Name: "a", Preconditions: "feature exists"},
			{Name: "b", Preconditions: "feature exists with slug matching arg"},
		}
		if ArePreconditionsMutuallyExclusive(outcomes) {
			t.Fatal("subset Preconditions should not be mutually exclusive")
		}
	})

	t.Run("truly exclusive Preconditions pass", func(t *testing.T) {
		outcomes := []Outcome{
			{Name: "success", Preconditions: "task status is in_progress"},
			{Name: "not-in-progress", Preconditions: "no task claimed"},
			{Name: "already-submitted", Preconditions: "task status is completed"},
		}
		if !ArePreconditionsMutuallyExclusive(outcomes) {
			t.Fatal("truly exclusive Preconditions should pass")
		}
	})
}

// --- Test: State verification level ---

func TestValidate_TUIAwaitSemantics(t *testing.T) {
	t.Run("async TUI Outcome without timeout Outcome fails", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].IsAsyncTUI = true
		c.Outcomes[0].AwaitTimeout = 3000
		c.Outcomes[0].Input = "key \"d\" await 3000ms"
		c.Outcomes[0].Output = "view contains diagnosis summary panel"
		errs := Validate(c)
		assertHasTUIAwaitError(t, errs)
	})

	t.Run("async TUI Outcome with timeout Outcome passes", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].IsAsyncTUI = true
		c.Outcomes[0].AwaitTimeout = 3000
		c.Outcomes[0].Input = "key \"d\" await 3000ms"
		c.Outcomes[0].Output = "view contains diagnosis summary panel"

		// Add the required timeout Outcome
		c.Outcomes = append(c.Outcomes, Outcome{
			Name:          "diagnosis-timeout",
			Preconditions: "async Cmd exceeds await duration of 3000ms",
			Input:         "key \"d\" await 3000ms",
			Output:        "error message containing timed-out Cmd name, fail-fast",
			State:         "unchanged from pre-Cmd state",
			IsAsyncTUI:    true,
			AwaitTimeout:  3000,
			TimedOutCmd:   "diagnosis-loader",
		})
		errs := Validate(c)
		assertNoErrors(t, errs)
	})

	t.Run("timeout Outcome must report timed-out Cmd name", func(t *testing.T) {
		c := validContract()
		c.Outcomes[0].IsAsyncTUI = true
		c.Outcomes[0].AwaitTimeout = 3000
		c.Outcomes[0].Input = "key \"d\" await 3000ms"
		c.Outcomes[0].Output = "view contains diagnosis panel"

		// Timeout Outcome without TimedOutCmd
		c.Outcomes = append(c.Outcomes, Outcome{
			Name:          "diagnosis-timeout",
			Preconditions: "async Cmd exceeds await duration of 3000ms",
			Input:         "key \"d\" await 3000ms",
			Output:        "error message about timeout",
			State:         "unchanged",
			IsAsyncTUI:    true,
			AwaitTimeout:  3000,
			TimedOutCmd:   "", // missing
		})
		errs := Validate(c)
		assertHasTimedOutCmdError(t, errs)
	})

	t.Run("default timeout when AwaitTimeout is 0", func(t *testing.T) {
		o := Outcome{
			IsAsyncTUI:   true,
			AwaitTimeout: 0,
		}
		timeout := ResolveAwaitTimeout(o, 0)
		if timeout != DefaultTUIAwaitTimeout {
			t.Fatalf("expected default %d, got %d", DefaultTUIAwaitTimeout, timeout)
		}
	})

	t.Run("config timeout overrides default", func(t *testing.T) {
		o := Outcome{
			IsAsyncTUI:   true,
			AwaitTimeout: 0,
		}
		timeout := ResolveAwaitTimeout(o, 5000)
		if timeout != 5000 {
			t.Fatalf("expected config 5000, got %d", timeout)
		}
	})

	t.Run("Outcome-specific timeout overrides config", func(t *testing.T) {
		o := Outcome{
			IsAsyncTUI:   true,
			AwaitTimeout: 10000,
		}
		timeout := ResolveAwaitTimeout(o, 5000)
		if timeout != 10000 {
			t.Fatalf("expected Outcome-specific 10000, got %d", timeout)
		}
	})
}

func TestStateVerificationLevel(t *testing.T) {
	t.Run("full state returns full", func(t *testing.T) {
		got := DetermineStateVerificationLevel("feature directory created with manifest.md", true)
		if got != StateVerificationFull {
			t.Fatalf("expected full, got %q", got)
		}
	})

	t.Run("inferred state without query returns partial", func(t *testing.T) {
		got := DetermineStateVerificationLevel("status can be inferred from output", false)
		if got != StateVerificationPartial {
			t.Fatalf("expected partial, got %q", got)
		}
	})

	t.Run("deferred state returns deferred", func(t *testing.T) {
		got := DetermineStateVerificationLevel("internal task index ordering", false)
		if got != StateVerificationDeferred {
			t.Fatalf("expected deferred, got %q", got)
		}
	})
}

// --- Assertions ---

func assertNoErrors(t *testing.T, errs []ValidationError) {
	t.Helper()
	if len(errs) > 0 {
		t.Fatalf("expected no validation errors, got %d: %v", len(errs), errs)
	}
}

func assertHasDimensionError(t *testing.T, errs []ValidationError, dimension string) {
	t.Helper()
	for _, e := range errs {
		if e.Dimension == dimension {
			return
		}
	}
	t.Fatalf("expected validation error for dimension %q, got %v", dimension, errs)
}

func assertHasSemanticError(t *testing.T, errs []ValidationError, dimension string) {
	t.Helper()
	for _, e := range errs {
		if e.Dimension == dimension && strings.Contains(e.Rule, "regex") {
			return
		}
	}
	t.Fatalf("expected semantic descriptor error for dimension %q, got %v", dimension, errs)
}

func assertHasMutualExclusivityError(t *testing.T, errs []ValidationError) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e.Rule, "mutually exclusive") {
			return
		}
	}
	t.Fatalf("expected mutual exclusivity error, got %v", errs)
}

func assertHasJourneyInvariantError(t *testing.T, errs []ValidationError) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e.Rule, "Journey Invariants") {
			return
		}
	}
	t.Fatalf("expected Journey Invariants error, got %v", errs)
}

func assertHasOutcomeNameError(t *testing.T, errs []ValidationError) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e.Rule, "Outcome name") {
			return
		}
	}
	t.Fatalf("expected Outcome name uniqueness error, got %v", errs)
}

func assertHasOutcomeCountWarning(t *testing.T, errs []ValidationError) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e.Rule, "Outcome count") {
			return
		}
	}
	t.Fatalf("expected Outcome count warning, got %v", errs)
}

func assertHasTUIAwaitError(t *testing.T, errs []ValidationError) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e.Rule, "TUI async") {
			return
		}
	}
	t.Fatalf("expected TUI async await error, got %v", errs)
}

func assertHasTimedOutCmdError(t *testing.T, errs []ValidationError) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e.Rule, "timed-out Cmd") {
			return
		}
	}
	t.Fatalf("expected timed-out Cmd name error, got %v", errs)
}
