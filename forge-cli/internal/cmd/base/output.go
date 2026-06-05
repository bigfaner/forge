package base

import (
	"fmt"
	"strings"

	"forge-cli/pkg/forgelog"
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
		forgelog.Debug("[debug] "+format+"\n", args...)
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

// TruncateSlug shortens a string to maxLen display width with ellipsis.
// Truncates at rune boundaries to avoid splitting multi-byte characters.
func TruncateSlug(s string, maxLen int) string {
	if DisplayWidth(s) <= maxLen {
		return s
	}
	// Walk runes, accumulate display width, stop before exceeding maxLen-3
	target := maxLen - 3 // reserve space for "..."
	w := 0
	for i, r := range s {
		rw := runeWidth(r)
		if w+rw > target {
			return s[:i] + "..."
		}
		w += rw
	}
	return s
}

// PadRight pads a string to exactly n display columns with trailing spaces.
func PadRight(s string, n int) string {
	dw := DisplayWidth(s)
	if dw >= n {
		return s
	}
	return s + strings.Repeat(" ", n-dw)
}

// DisplayWidth returns the terminal display width of a string,
// counting East Asian wide/fullwidth characters as 2 columns.
func DisplayWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeWidth(r)
	}
	return w
}

// runeWidth returns 2 for wide runes (CJK, fullwidth, etc.), 1 otherwise.
func runeWidth(r rune) int {
	if isWide(r) {
		return 2
	}
	return 1
}

func isWide(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) ||
		r == 0x2329 || r == 0x232A ||
		(r >= 0x2E80 && r <= 0xA4CF && r != 0x303F) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE10 && r <= 0xFE19) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF01 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6) ||
		(r >= 0x20000 && r <= 0x2FFFD) ||
		(r >= 0x30000 && r <= 0x3FFFD)
}
