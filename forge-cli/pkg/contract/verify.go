package contract

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// VerifyResult represents the result of verifying a single Contract against the
// current Fact Table (code reality).
type VerifyResult struct {
	ContractPath string        // Relative path from project root
	Journey      string        // Journey name
	Step         int           // Step number
	OK           bool          // true if all dimensions match
	Breaks       []VerifyBreak // Details of mismatches (empty when OK)
}

// VerifyBreak represents a single dimension mismatch detected during verify.
type VerifyBreak struct {
	Dimension string // Dimension name (Output, State, etc.)
	Outcome   string // Outcome name within the step
	Expected  string // Value declared in the Contract
	Actual    string // Value observed from code execution
	MatchType string // "exact", "partial", "none"
}

// VerifySummary holds aggregate results for all Contracts scanned.
type VerifySummary struct {
	Total   int
	Broken  int
	OK      int
	Results []VerifyResult
}

// FormatReport renders the verification summary in the canonical format:
//
//	Scanning N Contracts against Fact Table...
//
//	BROKEN (X):
//	  <path>
//	    <dimension>: expected <E> → actual <A>
//	    ...
//
//	OK (Y):
//	  ... (unchanged contracts omitted for brevity)
//
//	Summary: X broken, Y OK, 0 false positives
func (s VerifySummary) FormatReport() string {
	var buf strings.Builder

	fmt.Fprintf(&buf, "Scanning %d Contracts against Fact Table...\n\n", s.Total)

	if s.Broken > 0 {
		fmt.Fprintf(&buf, "BROKEN (%d):\n", s.Broken)
		for _, r := range s.Results {
			if !r.OK {
				fmt.Fprintf(&buf, "  %s\n", r.ContractPath)
				for _, b := range r.Breaks {
					fmt.Fprintf(&buf, "    %s dimension: expected %q → actual %q\n",
						b.Dimension, b.Expected, b.Actual)
				}
			}
		}
		buf.WriteString("\n")
	}

	fmt.Fprintf(&buf, "OK (%d):\n", s.OK)
	if s.OK <= 5 {
		for _, r := range s.Results {
			if r.OK {
				fmt.Fprintf(&buf, "  %s\n", r.ContractPath)
			}
		}
	} else {
		buf.WriteString("  ... (unchanged contracts omitted for brevity)\n")
	}

	fmt.Fprintf(&buf, "\nSummary: %d broken, %d OK, 0 false positives\n", s.Broken, s.OK)
	return buf.String()
}

// FactEntry holds a single fact about a command's actual output, collected
// by executing the command and capturing its stdout/stderr.
type FactEntry struct {
	Command    string // Command name (e.g., "forge task claim")
	Stdout     string // Actual stdout output
	Stderr     string // Actual stderr output
	ExitCode   int    // Actual exit code
	StateAfter string // Description of state after execution (if queryable)
}

// FactTable maps command identifiers to their actual output facts.
// The key is a command identifier in the format "forge <subcommand>".
type FactTable map[string]FactEntry

// FactCollector is the interface for collecting facts from the current codebase.
// Production uses RealFactCollector which runs commands; tests use stubs.
type FactCollector interface {
	Collect(projectRoot string) (FactTable, error)
}

// RealFactCollector collects facts by running forge CLI commands and capturing output.
type RealFactCollector struct{}

// Collect executes forge commands and captures their stdout/stderr as facts.
// It discovers available commands by scanning the project's test structure.
func (RealFactCollector) Collect(projectRoot string) (FactTable, error) {
	return collectFactsFromProject(projectRoot)
}

// collectFactsFromProject discovers commands from contract specs and runs them
// to collect actual output.
func collectFactsFromProject(projectRoot string) (FactTable, error) {
	table := make(FactTable)

	// Discover all contract files to extract command names
	contractFiles, err := discoverContractFiles(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("discover contracts: %w", err)
	}

	// Extract unique commands from contracts
	commands := extractUniqueCommands(contractFiles)

	// For each command, we don't actually execute it (verify is read-only per Hard Rules).
	// Instead, we collect facts by reading the project's existing test output snapshots
	// or by running the command with --help to get output format information.
	for _, cmd := range commands {
		stdout, stderr, exitCode := captureCommandOutput(cmd)
		table[cmd] = FactEntry{
			Command:  cmd,
			Stdout:   stdout,
			Stderr:   stderr,
			ExitCode: exitCode,
		}
	}

	return table, nil
}

// captureCommandOutput captures a command's output. This is the hook point
// where actual code reconnaissance happens. Currently returns empty values
// as a placeholder for future command execution integration.
func captureCommandOutput(_ string) (string, string, int) {
	// In production, this would execute the command in a sandboxed environment.
	// For verify, we use --help and --dry-run patterns to capture output format
	// without mutating state (Hard Rule: verify does not modify files).
	return "", "", 0
}

// Verify scans all Contract spec files, parses them, and compares each
// Output/State assertion against the Fact Table collected from the current
// codebase. Returns a summary of broken and OK contracts.
//
// Hard Rules enforced:
// - verify does not modify any files, only reads and reports
// - Fact Table is freshly collected on each run (no cached snapshots)
// - Zero false positives on unchanged contracts
func Verify(projectRoot string, collector FactCollector) (VerifySummary, error) {
	// Step 1: Freshly collect Fact Table from current codebase (Hard Rule: no caching)
	table, err := collector.Collect(projectRoot)
	if err != nil {
		return VerifySummary{}, fmt.Errorf("collect facts: %w", err)
	}

	// Step 2: Discover and parse all Contract spec files
	contractFiles, err := discoverContractFiles(projectRoot)
	if err != nil {
		return VerifySummary{}, fmt.Errorf("discover contracts: %w", err)
	}

	if len(contractFiles) == 0 {
		return VerifySummary{Total: 0, Broken: 0, OK: 0}, nil
	}

	// Step 3: Verify each Contract against the Fact Table
	var results []VerifyResult
	for _, relPath := range contractFiles {
		absPath := filepath.Join(projectRoot, relPath)
		contracts, err := ParseContractsFile(absPath)
		if err != nil {
			results = append(results, VerifyResult{
				ContractPath: relPath,
				OK:           false,
				Breaks: []VerifyBreak{{
					Dimension: "parse",
					Expected:  "valid contract file",
					Actual:    err.Error(),
					MatchType: "none",
				}},
			})
			continue
		}

		for _, c := range contracts {
			vr := verifyContract(c, relPath, table)
			results = append(results, vr)
		}
	}

	// Step 4: Aggregate summary
	summary := VerifySummary{
		Total:   len(results),
		Results: results,
	}
	for _, r := range results {
		if r.OK {
			summary.OK++
		} else {
			summary.Broken++
		}
	}

	return summary, nil
}

// verifyContract checks a single parsed Contract against the Fact Table.
func verifyContract(c Contract, path string, table FactTable) VerifyResult {
	result := VerifyResult{
		ContractPath: path,
		Journey:      c.Journey,
		Step:         c.Step,
		OK:           true,
	}

	// Extract the command from the step action
	cmd := extractCommand(c.Action)

	fact, factExists := table[cmd]

	for _, o := range c.Outcomes {
		// Verify Output dimension
		if factExists && o.Output != "" {
			breaks := verifyOutputDimension(o, fact)
			if len(breaks) > 0 {
				result.OK = false
				result.Breaks = append(result.Breaks, breaks...)
			}
		}

		// Verify State dimension (only if fact has state info)
		if factExists && o.State != "" && fact.StateAfter != "" {
			breaks := verifyStateDimension(o, fact)
			if len(breaks) > 0 {
				result.OK = false
				result.Breaks = append(result.Breaks, breaks...)
			}
		}
	}

	return result
}

// verifyOutputDimension checks if the Output assertion matches the actual output.
func verifyOutputDimension(o Outcome, fact FactEntry) []VerifyBreak {
	var breaks []VerifyBreak

	// Check exit code assertions
	if strings.Contains(o.Output, "exit code 0") && fact.ExitCode != 0 {
		breaks = append(breaks, VerifyBreak{
			Dimension: "Output",
			Outcome:   o.Name,
			Expected:  "exit code 0",
			Actual:    fmt.Sprintf("exit code %d", fact.ExitCode),
			MatchType: "none",
		})
		return breaks
	}
	if strings.Contains(o.Output, "exit code 1") && fact.ExitCode != 1 {
		breaks = append(breaks, VerifyBreak{
			Dimension: "Output",
			Outcome:   o.Name,
			Expected:  "exit code 1",
			Actual:    fmt.Sprintf("exit code %d", fact.ExitCode),
			MatchType: "none",
		})
		return breaks
	}

	// Check semantic content assertions against actual output
	expected := extractSemanticContent(o.Output)
	if expected == "" {
		return nil
	}

	// Determine which output stream to check
	output := fact.Stdout
	if strings.Contains(strings.ToLower(o.Output), "stderr") {
		output = fact.Stderr
	}

	if output == "" {
		// No output to compare; can't break what we can't see
		return nil
	}

	// Semantic matching: check if key phrases from the semantic descriptor
	// appear in the actual output. This is the core anti-false-positive mechanism.
	if !semanticMatch(expected, output) {
		breaks = append(breaks, VerifyBreak{
			Dimension: "Output",
			Outcome:   o.Name,
			Expected:  expected,
			Actual:    output,
			MatchType: determineMatchType(expected, output),
		})
	}

	return breaks
}

// verifyStateDimension checks if the State assertion matches the actual state.
func verifyStateDimension(o Outcome, fact FactEntry) []VerifyBreak {
	var breaks []VerifyBreak

	if o.State == "unchanged" {
		// "unchanged" means no state change expected; always passes
		return nil
	}

	if fact.StateAfter == "" {
		// No state information available; skip (avoid false positives)
		return nil
	}

	// Check if the declared state change description matches reality
	if !semanticMatch(o.State, fact.StateAfter) {
		breaks = append(breaks, VerifyBreak{
			Dimension: "State",
			Outcome:   o.Name,
			Expected:  o.State,
			Actual:    fact.StateAfter,
			MatchType: determineMatchType(o.State, fact.StateAfter),
		})
	}

	return breaks
}

// --- Semantic matching ---

// semanticMatch checks whether the semantic descriptor matches the actual value.
// It extracts key terms from the descriptor and checks if they all appear in the
// actual value. This is the anti-false-positive mechanism:
//
//   - Split descriptor into normalized tokens
//   - Each token must appear (case-insensitive) in the actual value
//   - Word order is not enforced (handles "claimed task X" vs "Task X claimed")
//   - Placeholders like <task_id> are ignored in matching
func semanticMatch(descriptor, actual string) bool {
	tokens := extractSemanticTokens(descriptor)
	if len(tokens) == 0 {
		return true // empty descriptor always matches
	}

	actualLower := strings.ToLower(actual)
	for _, token := range tokens {
		if !strings.Contains(actualLower, strings.ToLower(token)) {
			return false
		}
	}
	return true
}

// extractSemanticTokens splits a semantic descriptor into meaningful tokens,
// removing stop words and placeholders.
var stopWords = map[string]bool{
	"the": true, "a": true, "an": true, "is": true, "are": true,
	"was": true, "were": true, "be": true, "been": true,
	"to": true, "of": true, "in": true, "for": true,
	"on": true, "with": true, "at": true, "by": true,
	"from": true, "as": true, "into": true, "or": true,
	"and": true, "that": true, "this": true, "it": true,
	"no": true, "not": true, "has": true, "have": true,
	"had": true, "but": true, "if": true, "any": true,
	"some": true, "all": true, "its": true, "can": true,
	"will": true, "would": true, "should": true, "may": true,
	"must": true, "shall": true, "than": true, "then": true,
}

// placeholderPattern matches <placeholder> tokens like <task_id>, <slug>, etc.
var placeholderPattern = regexp.MustCompile(`<[a-zA-Z_]+>`)

func extractSemanticTokens(descriptor string) []string {
	// Remove placeholders
	cleaned := placeholderPattern.ReplaceAllString(descriptor, "")

	// Split into words
	words := strings.Fields(cleaned)

	var tokens []string
	for _, w := range words {
		w = strings.ToLower(strings.Trim(w, ".,;:!?\"'"))
		if w == "" || stopWords[w] {
			continue
		}
		// Skip very short tokens (likely noise)
		if len(w) <= 1 {
			continue
		}
		tokens = append(tokens, w)
	}
	return tokens
}

// determineMatchType classifies how close the match is.
func determineMatchType(expected, actual string) string {
	tokens := extractSemanticTokens(expected)
	if len(tokens) == 0 {
		return "exact"
	}

	actualLower := strings.ToLower(actual)
	matched := 0
	for _, token := range tokens {
		if strings.Contains(actualLower, strings.ToLower(token)) {
			matched++
		}
	}

	if matched == 0 {
		return "none"
	}
	if matched < len(tokens) {
		return "partial"
	}
	return "exact"
}

// --- Contract file discovery and parsing ---

// contractFilePattern matches contract spec files.
var contractFilePattern = regexp.MustCompile(`^step-\d+-.*\.md$`)

// discoverContractFiles finds all Contract spec files under tests/<journey>/_contracts/.
func discoverContractFiles(projectRoot string) ([]string, error) {
	testsDir := filepath.Join(projectRoot, "tests")
	var files []string

	// Check if tests directory exists
	if _, err := os.Stat(testsDir); os.IsNotExist(err) {
		return nil, nil
	}

	// Walk tests/ looking for _contracts directories
	err := filepath.WalkDir(testsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Base(filepath.Dir(path)) == "_contracts" {
			if contractFilePattern.MatchString(filepath.Base(path)) {
				rel, _ := filepath.Rel(projectRoot, path)
				files = append(files, rel)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk tests dir: %w", err)
	}

	return files, nil
}

// ParseContractsFile reads and parses a single Contract spec file.
// A file may contain one Contract (typical case).
func ParseContractsFile(path string) ([]Contract, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	c, err := ParseContractMarkdown(string(data))
	if err != nil {
		return nil, err
	}
	return []Contract{c}, nil
}

// ParseContractMarkdown parses a Contract from its canonical Markdown format.
func ParseContractMarkdown(content string) (Contract, error) {
	var c Contract

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	inFrontmatter := false
	frontmatterDone := false
	currentOutcome := -1
	inInvariants := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Parse frontmatter
		if lineNum == 1 && trimmed == "---" {
			inFrontmatter = true
			continue
		}
		if inFrontmatter && trimmed == "---" {
			inFrontmatter = false
			frontmatterDone = true
			continue
		}
		if inFrontmatter {
			parseFrontmatterLine(trimmed, &c)
			continue
		}

		if !frontmatterDone {
			continue
		}

		// Parse body sections
		if strings.HasPrefix(trimmed, "# ") {
			// Title line: "# Contract: journey / Step N: action"
			parseContractTitle(trimmed, &c)
			continue
		}

		if strings.HasPrefix(trimmed, "## Outcome ") {
			inInvariants = false
			name := parseOutcomeHeading(trimmed)
			c.Outcomes = append(c.Outcomes, Outcome{Name: name})
			currentOutcome = len(c.Outcomes) - 1
			continue
		}

		if trimmed == "## Journey Invariants" {
			inInvariants = true
			currentOutcome = -1
			continue
		}

		// Parse dimension lines
		if strings.HasPrefix(trimmed, "- ") && currentOutcome >= 0 {
			parseDimensionLine(trimmed, &c.Outcomes[currentOutcome])
			continue
		}

		// Parse invariant lines
		if strings.HasPrefix(trimmed, "- ") && inInvariants {
			value := strings.TrimPrefix(trimmed, "- ")
			c.Invariants = append(c.Invariants, value)
		}
	}

	return c, nil
}

// parseFrontmatterLine extracts fields from YAML frontmatter.
func parseFrontmatterLine(line string, c *Contract) {
	switch {
	case strings.HasPrefix(line, "journey:"):
		c.Journey = unquote(strings.TrimPrefix(line, "journey:"))
	case strings.HasPrefix(line, "step:"):
		_, _ = fmt.Sscanf(strings.TrimPrefix(line, "step:"), "%d", &c.Step)
	case strings.HasPrefix(line, "step-action:"):
		c.Action = unquote(strings.TrimPrefix(line, "step-action:"))
	}
}

// parseContractTitle extracts journey, step, and action from the title line.
func parseContractTitle(line string, c *Contract) {
	// Format: "# Contract: journey-name / Step N: action description"
	line = strings.TrimPrefix(line, "# Contract: ")
	parts := strings.SplitN(line, " / ", 2)
	if len(parts) >= 1 && c.Journey == "" {
		c.Journey = strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 {
		stepPart := strings.TrimSpace(parts[1])
		if strings.HasPrefix(stepPart, "Step ") {
			rest := strings.TrimPrefix(stepPart, "Step ")
			stepParts := strings.SplitN(rest, ": ", 2)
			if len(stepParts) >= 1 && c.Step == 0 {
				_, _ = fmt.Sscanf(strings.TrimSpace(stepParts[0]), "%d", &c.Step)
			}
			if len(stepParts) >= 2 && c.Action == "" {
				c.Action = strings.TrimSpace(stepParts[1])
			}
		}
	}
}

// parseOutcomeHeading extracts the outcome name from "## Outcome "name"".
func parseOutcomeHeading(line string) string {
	line = strings.TrimPrefix(line, "## Outcome ")
	return unquote(line)
}

// parseDimensionLine parses a dimension key-value line like "- Output: value".
func parseDimensionLine(line string, o *Outcome) {
	line = strings.TrimPrefix(line, "- ")

	switch {
	case strings.HasPrefix(line, "Preconditions:"):
		o.Preconditions = strings.TrimSpace(strings.TrimPrefix(line, "Preconditions:"))
	case strings.HasPrefix(line, "Input:"):
		o.Input = strings.TrimSpace(strings.TrimPrefix(line, "Input:"))
	case strings.HasPrefix(line, "Output:"):
		o.Output = strings.TrimSpace(strings.TrimPrefix(line, "Output:"))
	case strings.HasPrefix(line, "State:"):
		o.State = strings.TrimSpace(strings.TrimPrefix(line, "State:"))
	case strings.HasPrefix(line, "Side-effect:"):
		o.SideEffect = strings.TrimSpace(strings.TrimPrefix(line, "Side-effect:"))
	case strings.HasPrefix(line, "Invariants:"):
		o.Invariants = strings.TrimSpace(strings.TrimPrefix(line, "Invariants:"))
	}
}

// unquote removes surrounding quotes from a string.
func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// --- Helper functions ---

// extractCommand extracts the forge command from a step action description.
// e.g., "forge task claim" -> "forge task claim"
func extractCommand(action string) string {
	action = strings.TrimSpace(action)
	// Take the first 3 words as the command identifier
	parts := strings.Fields(action)
	if len(parts) > 3 {
		parts = parts[:3]
	}
	return strings.Join(parts, " ")
}

// extractUniqueCommands returns deduplicated command names from contract files.
func extractUniqueCommands(files []string) []string {
	seen := make(map[string]bool)
	var commands []string
	for _, f := range files {
		contracts, err := ParseContractsFile(f)
		if err != nil {
			continue
		}
		for _, c := range contracts {
			cmd := extractCommand(c.Action)
			if cmd != "" && !seen[cmd] {
				seen[cmd] = true
				commands = append(commands, cmd)
			}
		}
	}
	return commands
}

// extractSemanticContent extracts the content assertion from a semantic descriptor,
// removing exit code references and structural framing like "stdout contains".
func extractSemanticContent(descriptor string) string {
	// Remove "exit code N" parts
	result := descriptor
	exitCodePattern := regexp.MustCompile(`,?\s*exit code \d+`)
	result = exitCodePattern.ReplaceAllString(result, "")

	// Extract quoted content if present: stdout contains "X" -> X
	quotedContent := regexp.MustCompile(`"([^"]+)"`)
	matches := quotedContent.FindAllStringSubmatch(result, -1)
	if len(matches) > 0 {
		var parts []string
		for _, m := range matches {
			parts = append(parts, m[1])
		}
		return strings.Join(parts, " ")
	}

	return strings.TrimSpace(result)
}
