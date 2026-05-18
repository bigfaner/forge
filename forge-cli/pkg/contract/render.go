package contract

import (
	"fmt"
	"strings"
)

// RenderContract renders a Contract to its canonical Markdown representation.
// The output format matches the Contract specification file format defined in
// the model-and-directory-spec.md document.
func RenderContract(c Contract) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	fmt.Fprintf(&sb, "journey: %q\n", c.Journey)
	fmt.Fprintf(&sb, "step: %d\n", c.Step)
	fmt.Fprintf(&sb, "step-action: %q\n", c.Action)
	sb.WriteString("---\n\n")

	// Title
	fmt.Fprintf(&sb, "# Contract: %s / Step %d: %s\n\n", c.Journey, c.Step, c.Action)

	// Outcome blocks
	for _, o := range c.Outcomes {
		fmt.Fprintf(&sb, "## Outcome %q\n", o.Name)
		fmt.Fprintf(&sb, "- Preconditions: %q\n", o.Preconditions)
		fmt.Fprintf(&sb, "- Input: %s\n", o.Input)
		fmt.Fprintf(&sb, "- Output: %s\n", o.Output)
		fmt.Fprintf(&sb, "- State: %s\n", o.State)

		if o.SideEffect != "" {
			fmt.Fprintf(&sb, "- Side-effect: %s\n", o.SideEffect)
		} else {
			sb.WriteString("- Side-effect: none\n")
		}

		if o.Invariants != "" {
			fmt.Fprintf(&sb, "- Invariants: %s\n", o.Invariants)
		}

		sb.WriteString("\n")
	}

	// State verification annotation
	if c.StateVerifyLevel == "partial" || c.StateVerifyLevel == "deferred" {
		fmt.Fprintf(&sb, "<!-- state-verification: %s -->\n", c.StateVerifyLevel)
		sb.WriteString("\n")
	}

	// Journey Invariants
	sb.WriteString("## Journey Invariants\n")
	for _, inv := range c.Invariants {
		fmt.Fprintf(&sb, "- %s\n", inv)
	}
	sb.WriteString("\n")

	return sb.String()
}
