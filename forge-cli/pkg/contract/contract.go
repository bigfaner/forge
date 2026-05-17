// Package contract provides types and validation for Journey-Driven Contract specifications.
package contract

// DefaultTUIAwaitTimeout is the default timeout for TUI async Cmd await in milliseconds.
const DefaultTUIAwaitTimeout = 3000

// Outcome represents a single Outcome within a Step Contract.
// Each Outcome declares its own set of dimension values.
type Outcome struct {
	Name          string // Descriptive label (e.g., "success", "not-in-progress")
	Preconditions string // Mandatory: state that must hold before execution
	Input         string // Mandatory: what goes into the system
	Output        string // Mandatory: what the system produces
	State         string // Mandatory: how system state changes
	SideEffect    string // Optional: external effects (empty = "none")
	Invariants    string // Optional: step-level invariants (empty = no constraint)
	IsAsyncTUI    bool   // True if this Outcome involves TUI async Cmd await
	AwaitTimeout  int    // Await timeout in ms (0 = use default from config)
	TimedOutCmd   string // Name of Cmd that timed out (for timeout Outcomes)
}

// Contract represents a Contract specification for a single Journey Step.
type Contract struct {
	Journey          string    // Journey name (kebab-case)
	Step             int       // 1-based step ordinal
	Action           string    // Human-readable step action description
	Outcomes         []Outcome // One or more Outcomes for this Step
	Invariants       []string  // Journey-level Invariants (at least 1 required)
	StateVerifyLevel string    // State verification level: "full", "partial", or "deferred"
}

// Dimension name constants for validation error reporting.
const (
	DimensionPreconditions = "Preconditions"
	DimensionInput         = "Input"
	DimensionOutput        = "Output"
	DimensionState         = "State"
	DimensionSideEffect    = "Side-effect"
	DimensionInvariants    = "Invariants"
)

// ValidationError represents a single validation failure in a Contract.
type ValidationError struct {
	Step      int    // Step number (1-based)
	Outcome   string // Outcome name (empty for Journey-level errors)
	Dimension string // Dimension name (empty for non-dimension errors)
	Rule      string // Human-readable rule that was violated
}

func (e ValidationError) Error() string {
	if e.Outcome != "" {
		return e.Rule
	}
	return e.Rule
}
