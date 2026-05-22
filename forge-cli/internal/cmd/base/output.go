package base

import (
	"fmt"
	"os"
	"strings"
)

// OutputBlockSeparator is the separator for structured output blocks.
const OutputBlockSeparator = "---"

// PrintBlockStart prints the opening separator for a structured output block.
func PrintBlockStart() {
	fmt.Println(OutputBlockSeparator)
}

// PrintBlockEnd prints the closing separator for a structured output block.
func PrintBlockEnd() {
	fmt.Println(OutputBlockSeparator)
}

// PrintBlock prints a complete structured output block with key-value pairs.
func PrintBlock(key, value string) {
	PrintBlockStart()
	fmt.Printf("%s: %s\n", key, value)
	PrintBlockEnd()
}

// PrintFields prints multiple key-value pairs within a block.
func PrintFields(pairs ...string) {
	if len(pairs)%2 != 0 {
		panic("PrintFields requires even number of arguments (key-value pairs)")
	}
	PrintBlockStart()
	for i := 0; i < len(pairs); i += 2 {
		fmt.Printf("%s: %s\n", pairs[i], pairs[i+1])
	}
	PrintBlockEnd()
}

// PrintField prints a single key-value line (without separators).
func PrintField(key, value string) {
	fmt.Printf("%s: %s\n", key, value)
}

// PrintFieldIfNotEmpty prints a key-value line only if value is not empty.
func PrintFieldIfNotEmpty(key, value string) {
	if value != "" {
		fmt.Printf("%s: %s\n", key, value)
	}
}

// PrintFieldIfNotEmptySlice prints a key-value line for a slice only if not empty.
func PrintFieldIfNotEmptySlice(key string, values []string) {
	if len(values) > 0 {
		fmt.Printf("%s: %s\n", key, strings.Join(values, ", "))
	}
}

// PrintSection prints a section header (uppercase with colon).
func PrintSection(name string) {
	fmt.Printf("\n[%s]\n", name)
}

// PrintResult prints a result line: "RESULT: <status> [<details>]"
func PrintResult(status, details string) {
	if details != "" {
		fmt.Printf("RESULT: %s (%s)\n", status, details)
	} else {
		fmt.Printf("RESULT: %s\n", status)
	}
}

// PrintWarning prints a warning line with prefix.
func PrintWarning(msg string) {
	fmt.Printf("WARNING: %s\n", msg)
}

// PrintListItem prints an indented list item.
func PrintListItem(item string) {
	fmt.Printf("  %s\n", item)
}

// Debugf prints a debug line to stderr if verbose is true.
func Debugf(verbose bool, format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args)
	}
}

// Slug column sizing constants for dynamic table formatting.
const (
	SlugColMinWidth = 30
	SlugColMaxWidth = 60
)

// CalcSlugColWidth returns the dynamic column width for slug/name display.
// Width = clamp(max(30, maxSlugLen+2), 60).
func CalcSlugColWidth(slugLens []int) int {
	maxLen := 0
	for _, l := range slugLens {
		if l > maxLen {
			maxLen = l
		}
	}
	width := maxLen + 2
	if width < SlugColMinWidth {
		width = SlugColMinWidth
	}
	if width > SlugColMaxWidth {
		width = SlugColMaxWidth
	}
	return width
}

// TruncateSlug shortens a string to maxLen with ellipsis.
func TruncateSlug(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// PadRight pads a string to exactly n characters with trailing spaces.
func PadRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
