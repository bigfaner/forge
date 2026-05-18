package contract

import (
	"regexp"
	"strings"
)

// regexPatterns detects common regex syntax elements in semantic descriptors.
var regexPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\\d`),
	regexp.MustCompile(`\\w`),
	regexp.MustCompile(`\\s`),
	regexp.MustCompile(`\\b`),
	regexp.MustCompile(`\.\*`),
	regexp.MustCompile(`\[[\w-]+\]`),
	regexp.MustCompile(`\(\?:`),
	regexp.MustCompile(`^\^`),
	regexp.MustCompile(`\$$`),
}

// Validate checks a Contract against all validation rules and returns
// any violations found. An empty slice means the Contract is valid.
func Validate(c Contract) []ValidationError {
	var errs []ValidationError

	// Validate Outcomes
	outcomeNames := make(map[string]bool)
	for i := range c.Outcomes {
		o := &c.Outcomes[i]

		// Check Outcome name uniqueness
		if outcomeNames[o.Name] {
			errs = append(errs, ValidationError{
				Step:    c.Step,
				Outcome: o.Name,
				Rule:    "Outcome name must be unique within a Step: " + o.Name,
			})
		}
		outcomeNames[o.Name] = true

		// Check mandatory dimensions
		if strings.TrimSpace(o.Preconditions) == "" {
			errs = append(errs, ValidationError{
				Step:      c.Step,
				Outcome:   o.Name,
				Dimension: DimensionPreconditions,
				Rule:      "Preconditions must not be empty (mandatory dimension)",
			})
		}
		if strings.TrimSpace(o.Input) == "" {
			errs = append(errs, ValidationError{
				Step:      c.Step,
				Outcome:   o.Name,
				Dimension: DimensionInput,
				Rule:      "Input must not be empty (mandatory dimension)",
			})
		}
		if strings.TrimSpace(o.Output) == "" {
			errs = append(errs, ValidationError{
				Step:      c.Step,
				Outcome:   o.Name,
				Dimension: DimensionOutput,
				Rule:      "Output must not be empty (mandatory dimension)",
			})
		}
		if strings.TrimSpace(o.State) == "" {
			errs = append(errs, ValidationError{
				Step:      c.Step,
				Outcome:   o.Name,
				Dimension: DimensionState,
				Rule:      "State must not be empty (mandatory dimension)",
			})
		}

		// Check semantic descriptor purity (no regex) for all dimensions
		dimensions := map[string]string{
			DimensionPreconditions: o.Preconditions,
			DimensionInput:         o.Input,
			DimensionOutput:        o.Output,
			DimensionState:         o.State,
			DimensionSideEffect:    o.SideEffect,
		}
		for dim, value := range dimensions {
			if ContainsRegex(value) {
				errs = append(errs, ValidationError{
					Step:      c.Step,
					Outcome:   o.Name,
					Dimension: dim,
					Rule:      "Semantic descriptor must not contain regex syntax in " + dim,
				})
			}
		}
	}

	// Check Preconditions mutual exclusivity
	if !ArePreconditionsMutuallyExclusive(c.Outcomes) {
		errs = append(errs, ValidationError{
			Step: c.Step,
			Rule: "Outcome Preconditions must be mutually exclusive",
		})
	}

	// Check Outcome count checkpoint
	if len(c.Outcomes) > 5 {
		errs = append(errs, ValidationError{
			Step: c.Step,
			Rule: "Outcome count exceeds 5: review and consider merging semantically similar Outcomes",
		})
	}

	// Check Journey Invariants presence
	if len(c.Invariants) == 0 {
		errs = append(errs, ValidationError{
			Step: c.Step,
			Rule: "Journey Invariants must have at least 1 entry",
		})
	}

	// Check TUI async await semantics
	errs = append(errs, validateTUIAwait(c)...)

	return errs
}

// validateTUIAwait checks that async TUI Outcomes have corresponding timeout Outcomes.
func validateTUIAwait(c Contract) []ValidationError {
	var errs []ValidationError

	// Find async TUI Outcomes
	var asyncOutcomes []int
	var timeoutOutcomes []int
	for i, o := range c.Outcomes {
		if !o.IsAsyncTUI {
			continue
		}
		if strings.Contains(strings.ToLower(o.Output), "timeout") ||
			strings.Contains(strings.ToLower(o.Name), "timeout") {
			timeoutOutcomes = append(timeoutOutcomes, i)
		} else {
			asyncOutcomes = append(asyncOutcomes, i)
		}
	}

	// If there are async Outcomes, at least one timeout Outcome must exist
	if len(asyncOutcomes) > 0 && len(timeoutOutcomes) == 0 {
		for _, idx := range asyncOutcomes {
			errs = append(errs, ValidationError{
				Step:    c.Step,
				Outcome: c.Outcomes[idx].Name,
				Rule:    "TUI async Outcome must have a corresponding timeout Outcome that reports the timed-out Cmd name",
			})
		}
	}

	// Timeout Outcomes must report the timed-out Cmd name
	for _, idx := range timeoutOutcomes {
		if c.Outcomes[idx].TimedOutCmd == "" {
			errs = append(errs, ValidationError{
				Step:    c.Step,
				Outcome: c.Outcomes[idx].Name,
				Rule:    "TUI timeout Outcome must specify the timed-out Cmd name in TimedOutCmd field",
			})
		}
	}

	return errs
}

// ResolveAwaitTimeout determines the effective await timeout for a TUI Outcome.
// Priority: Outcome-specific > config > default (3000ms).
func ResolveAwaitTimeout(o Outcome, configTimeout int) int {
	if o.AwaitTimeout > 0 {
		return o.AwaitTimeout
	}
	if configTimeout > 0 {
		return configTimeout
	}
	return DefaultTUIAwaitTimeout
}

// ContainsRegex checks whether a string contains common regex syntax patterns.
// Returns true if any regex metacharacters are detected, false for pure natural language.
func ContainsRegex(s string) bool {
	for _, p := range regexPatterns {
		if p.MatchString(s) {
			return true
		}
	}
	return false
}

// ArePreconditionsMutuallyExclusive checks whether all Outcomes have
// distinguishable Preconditions. Two Preconditions are considered overlapping
// if one is a substring of the other (subset relationship) or they are identical.
func ArePreconditionsMutuallyExclusive(outcomes []Outcome) bool {
	if len(outcomes) <= 1 {
		return true
	}

	for i := 0; i < len(outcomes); i++ {
		for j := i + 1; j < len(outcomes); j++ {
			a := strings.TrimSpace(outcomes[i].Preconditions)
			b := strings.TrimSpace(outcomes[j].Preconditions)
			if a == "" || b == "" {
				continue
			}
			// Identical Preconditions
			if a == b {
				return false
			}
			// Subset check: if one contains the other as a substring,
			// they are overlapping
			if strings.Contains(a, b) || strings.Contains(b, a) {
				return false
			}
		}
	}
	return true
}

// StateVerificationLevel represents the degree to which State can be verified.
type StateVerificationLevel string

const (
	// StateVerificationFull means all state fields can be independently verified.
	StateVerificationFull StateVerificationLevel = "full"
	// StateVerificationPartial means state fields can be inferred from Output only.
	StateVerificationPartial StateVerificationLevel = "partial"
	// StateVerificationDeferred means some state fields cannot be inferred from Output.
	StateVerificationDeferred StateVerificationLevel = "deferred"
)

// DetermineStateVerificationLevel determines the appropriate State verification
// level based on the State description and whether a state query interface exists.
func DetermineStateVerificationLevel(stateDesc string, hasStateQueryInterface bool) StateVerificationLevel {
	if hasStateQueryInterface {
		return StateVerificationFull
	}
	stateLower := strings.ToLower(stateDesc)
	if strings.Contains(stateLower, "inferred") || strings.Contains(stateLower, "from output") {
		return StateVerificationPartial
	}
	return StateVerificationDeferred
}

// BatchNeeded returns true if the number of Contracts exceeds the batching
// threshold (15).
func BatchNeeded(contracts []Contract) bool {
	return len(contracts) >= 15
}

// BatchSplit splits Contracts into batches: first batch contains only
// success Outcomes, subsequent batches contain remaining Outcomes.
func BatchSplit(contracts []Contract) [][]Contract {
	if len(contracts) == 0 {
		return nil
	}

	// Batch 1: happy path (success Outcomes only)
	var happyPath []Contract
	var edgeCases []Contract

	for _, c := range contracts {
		var successOutcomes []Outcome
		var otherOutcomes []Outcome
		for _, o := range c.Outcomes {
			if o.Name == "success" {
				successOutcomes = append(successOutcomes, o)
			} else {
				otherOutcomes = append(otherOutcomes, o)
			}
		}
		if len(successOutcomes) > 0 {
			happyPath = append(happyPath, Contract{
				Journey:    c.Journey,
				Step:       c.Step,
				Action:     c.Action,
				Outcomes:   successOutcomes,
				Invariants: c.Invariants,
			})
		}
		if len(otherOutcomes) > 0 {
			edgeCases = append(edgeCases, Contract{
				Journey:    c.Journey,
				Step:       c.Step,
				Action:     c.Action,
				Outcomes:   otherOutcomes,
				Invariants: c.Invariants,
			})
		}
	}

	var batches [][]Contract
	if len(happyPath) > 0 {
		batches = append(batches, happyPath)
	}
	if len(edgeCases) > 0 {
		batches = append(batches, edgeCases)
	}
	return batches
}
