package cmd

import (
	"fmt"
	"os"

	"forge-cli/internal/cmd/base"
)

// Re-export output functions from base package for backward compatibility.
var (
	PrintBlockStart           = base.PrintBlockStart
	PrintBlockEnd             = base.PrintBlockEnd
	PrintBlock                = base.PrintBlock
	PrintFields               = base.PrintFields
	PrintField                = base.PrintField
	PrintFieldIfNotEmpty      = base.PrintFieldIfNotEmpty
	PrintFieldIfNotEmptySlice = base.PrintFieldIfNotEmptySlice
	PrintSection              = base.PrintSection
	PrintResult               = base.PrintResult
	PrintWarning              = base.PrintWarning
	PrintListItem             = base.PrintListItem
)

// Debugf prints a debug line to stderr if verbose is true.
// Inlined from base to preserve variadic call semantics across package boundaries.
func Debugf(verbose bool, format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}

// Re-export slug formatting utilities from base package for backward compatibility.
var (
	CalcSlugColWidth = base.CalcSlugColWidth
	TruncateSlug     = base.TruncateSlug
	PadRight         = base.PadRight
)
