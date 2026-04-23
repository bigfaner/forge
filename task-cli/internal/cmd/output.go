// Package cmd provides the CLI commands for the task management tool.
package cmd

import (
	"fmt"
	"os"
	"strings"
)

// Output separator for structured output blocks.
const outputBlockSeparator = "---"

// PrintBlockStart prints the opening separator for a structured output block.
func PrintBlockStart() {
	fmt.Println(outputBlockSeparator)
}

// PrintBlockEnd prints the closing separator for a structured output block.
func PrintBlockEnd() {
	fmt.Println(outputBlockSeparator)
}

// PrintBlock prints a complete structured output block with key-value pairs.
// Example:
//
//	PrintBlock("FEATURE", "my-feature")
//
// Output:
//
//	---
//	FEATURE: my-feature
//	---
func PrintBlock(key, value string) {
	PrintBlockStart()
	fmt.Printf("%s: %s\n", key, value)
	PrintBlockEnd()
}

// PrintFields prints multiple key-value pairs within a block.
// Example:
//
//	PrintFields("KEY", "task1", "ID", "1.2.3", "STATUS", "pending")
//
// Output:
//
//	---
//	KEY: task1
//	ID: 1.2.3
//	STATUS: pending
//	---
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
// Use this inside custom blocks.
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

// PrintKeyValue prints a simple key-value line (alias for PrintField).
func PrintKeyValue(key, value string) {
	fmt.Printf("%s: %s\n", key, value)
}

// PrintSection prints a section header (uppercase with colon).
// Example: PrintSection("ERRORS") -> "ERRORS:"
func PrintSection(name string) {
	fmt.Printf("\n[%s]\n", name)
}

// PrintResult prints a result line: "RESULT: <status> [<details>]"
// Example: PrintResult("PASS", "") -> "RESULT: PASS"
// Example: PrintResult("FAIL", "2 errors") -> "RESULT: FAIL (2 errors)"
func PrintResult(status, details string) {
	if details != "" {
		fmt.Printf("RESULT: %s (%s)\n", status, details)
	} else {
		fmt.Printf("RESULT: %s\n", status)
	}
}

// PrintError prints an error line with prefix.
func PrintError(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
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
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}
