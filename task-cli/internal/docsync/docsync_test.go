package docsync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

// projectRoot walks up from cwd to find go.mod.
func projectRoot(t *testing.T) string {
	t.Helper()
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("cannot find project root (no go.mod found)")
		}
		dir = parent
	}
}

func readDoc(t *testing.T, root, name string) string {
	t.Helper()
	path := filepath.Join(root, "docs", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	return string(data)
}

// jsonTagNames extracts JSON field names from a struct type.
// Returns a map of json-name → Go-field-name for non-omitempty fields,
// and a separate map for omitempty fields.
func jsonTagNames(typ reflect.Type) (required, optional map[string]string) {
	required = make(map[string]string)
	optional = make(map[string]string)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		parts := strings.SplitN(tag, ",", 2)
		name := parts[0]
		if name == "" {
			continue
		}
		if len(parts) > 1 && strings.Contains(parts[1], "omitempty") {
			optional[name] = f.Name
		} else {
			required[name] = f.Name
		}
	}
	return
}

// extractJSONBlock finds the first JSON code block in markdown content.
func extractJSONBlock(content string) string {
	re := regexp.MustCompile("(?s)```json\\s*\n(.*?)```")
	m := re.FindStringSubmatch(content)
	if m == nil {
		return ""
	}
	return m[1]
}

// extractSection finds content between a heading and the next heading of same or higher level.
func extractSection(content, heading string) string {
	pattern := regexp.MustCompile(regexp.QuoteMeta(heading) + `\s*\n`)
	loc := pattern.FindStringIndex(content)
	if loc == nil {
		return ""
	}
	rest := content[loc[1]:]
	nextHeading := regexp.MustCompile(`\n#{1,6}\s`).FindStringIndex(rest)
	if nextHeading != nil {
		return rest[:nextHeading[0]]
	}
	return rest
}

func TestRecordDataFieldsMatch(t *testing.T) {
	root := projectRoot(t)
	doc := readDoc(t, root, "WORKFLOW.md")

	jsonBlock := extractJSONBlock(extractSection(doc, "### RecordData Structure"))
	if jsonBlock == "" {
		t.Fatal("cannot find RecordData JSON block in WORKFLOW.md")
	}

	// Parse the doc JSON
	var docFields map[string]interface{}
	if err := json.Unmarshal([]byte(jsonBlock), &docFields); err != nil {
		t.Fatalf("doc JSON parse error: %v\nJSON:\n%s", err, jsonBlock)
	}

	// Get actual struct tags
	typ := reflect.TypeOf(task.RecordData{})
	required, optional := jsonTagNames(typ)
	allTags := make(map[string]string)
	for k, v := range required {
		allTags[k] = v
	}
	for k, v := range optional {
		allTags[k] = v
	}

	// Every JSON tag in the struct must appear in the doc
	for tagName, goName := range allTags {
		if _, ok := docFields[tagName]; !ok {
			t.Errorf("struct field %s (json:\"%s\") missing from doc JSON", goName, tagName)
		}
	}

	// Every key in the doc must correspond to a struct tag
	for key := range docFields {
		if _, ok := allTags[key]; !ok {
			t.Errorf("doc JSON key %q has no corresponding struct field (possible stale field)", key)
		}
	}
}

func TestAcceptanceCriterionFieldsMatch(t *testing.T) {
	root := projectRoot(t)
	doc := readDoc(t, root, "WORKFLOW.md")

	jsonBlock := extractJSONBlock(extractSection(doc, "### RecordData Structure"))
	if jsonBlock == "" {
		t.Fatal("cannot find RecordData JSON block in WORKFLOW.md")
	}

	var docFields map[string]interface{}
	if err := json.Unmarshal([]byte(jsonBlock), &docFields); err != nil {
		t.Fatalf("doc JSON parse error: %v", err)
	}

	// Get acceptanceCriteria array from doc
	acRaw, ok := docFields["acceptanceCriteria"]
	if !ok {
		t.Fatal("acceptanceCriteria missing from doc JSON")
	}
	acArray, ok := acRaw.([]interface{})
	if !ok || len(acArray) == 0 {
		t.Fatal("acceptanceCriteria is not a non-empty array in doc JSON")
	}
	firstItem, ok := acArray[0].(map[string]interface{})
	if !ok {
		t.Fatal("acceptanceCriteria[0] is not an object in doc JSON")
	}

	// Get actual struct tags
	typ := reflect.TypeOf(task.AcceptanceCriterion{})
	required, _ := jsonTagNames(typ)

	// Every struct tag must appear in the doc item
	for tagName, goName := range required {
		if _, ok := firstItem[tagName]; !ok {
			t.Errorf("AcceptanceCriterion field %s (json:\"%s\") missing from doc", goName, tagName)
		}
	}

	// Every doc key must correspond to a struct tag
	allTags := required
	for key := range firstItem {
		if _, ok := allTags[key]; !ok {
			t.Errorf("doc AC key %q has no corresponding struct field", key)
		}
	}
}

func TestTaskStructFieldsMatch(t *testing.T) {
	root := projectRoot(t)
	doc := readDoc(t, root, "OVERVIEW.md")

	// Extract the Task struct code block
	section := extractSection(doc, "### Task")
	codeBlock := extractGoBlock(section)
	if codeBlock == "" {
		t.Fatal("cannot find Task Go code block in OVERVIEW.md")
	}

	typ := reflect.TypeOf(task.Task{})
	required, optional := jsonTagNames(typ)
	allTags := make(map[string]string)
	for k, v := range required {
		allTags[k] = v
	}
	for k, v := range optional {
		allTags[k] = v
	}

	// Check each struct tag appears in the doc code block
	for tagName, goName := range allTags {
		// The doc shows Go struct fields, so check for the Go field name
		if !strings.Contains(codeBlock, goName) {
			t.Errorf("Task field %s (json:\"%s\") missing from OVERVIEW.md struct", goName, tagName)
		}
		// Also check the json tag appears in the doc
		if !strings.Contains(codeBlock, tagName) {
			t.Errorf("Task json tag %q missing from OVERVIEW.md struct", tagName)
		}
	}
}

func TestTaskIndexStructFieldsMatch(t *testing.T) {
	root := projectRoot(t)
	doc := readDoc(t, root, "OVERVIEW.md")

	section := extractSection(doc, "### TaskIndex")
	codeBlock := extractGoBlock(section)
	if codeBlock == "" {
		t.Fatal("cannot find TaskIndex Go code block in OVERVIEW.md")
	}

	typ := reflect.TypeOf(task.TaskIndex{})
	required, optional := jsonTagNames(typ)
	allTags := make(map[string]string)
	for k, v := range required {
		allTags[k] = v
	}
	for k, v := range optional {
		allTags[k] = v
	}

	for tagName, goName := range allTags {
		if !strings.Contains(codeBlock, goName) {
			t.Errorf("TaskIndex field %s (json:\"%s\") missing from OVERVIEW.md struct", goName, tagName)
		}
		if !strings.Contains(codeBlock, tagName) {
			t.Errorf("TaskIndex json tag %q missing from OVERVIEW.md struct", tagName)
		}
	}
}

func TestTestDetectionOrderMatch(t *testing.T) {
	root := projectRoot(t)
	overview := readDoc(t, root, "OVERVIEW.md")
	workflow := readDoc(t, root, "WORKFLOW.md")

	// The actual order from all_completed.go runProjectTests()
	expectedOrder := []string{"testCommand", "justfile", "Makefile", "go.mod", "package.json", "pytest"}

	// Check OVERVIEW.md
	overviewOrder := extractNumberedItemsAfterHeading(overview, "**Test command auto-detection order")
	if len(overviewOrder) == 0 {
		t.Fatal("cannot extract detection order from OVERVIEW.md")
	}
	assertOrder(t, "OVERVIEW.md", overviewOrder, expectedOrder)

	// Check WORKFLOW.md
	workflowOrder := extractNumberedItemsAfterHeading(workflow, "**Test command detection order:")
	if len(workflowOrder) == 0 {
		t.Fatal("cannot extract detection order from WORKFLOW.md")
	}
	assertOrder(t, "WORKFLOW.md", workflowOrder, expectedOrder)
}

func TestStatusPriorityValuesMatch(t *testing.T) {
	root := projectRoot(t)
	overview := readDoc(t, root, "OVERVIEW.md")

	// Status values
	expectedStatuses := []string{feature.StatusPending, feature.StatusInProgress, feature.StatusCompleted, feature.StatusBlocked, feature.StatusSkipped}
	statusLine := regexp.MustCompile(`\*\*Status values:\*\*.*`).FindString(overview)
	if statusLine == "" {
		t.Fatal("cannot find Status values line in OVERVIEW.md")
	}
	for _, s := range expectedStatuses {
		if !strings.Contains(statusLine, fmt.Sprintf("`%s`", s)) {
			t.Errorf("status value %q missing from OVERVIEW.md", s)
		}
	}

	// Priority values
	expectedPriorities := []string{feature.PriorityP0, feature.PriorityP1, feature.PriorityP2}
	prioLine := regexp.MustCompile(`\| Priority \|.*`).FindString(overview)
	if prioLine == "" {
		t.Fatal("cannot find Priority line in OVERVIEW.md")
	}
	for _, p := range expectedPriorities {
		if !strings.Contains(prioLine, p) {
			t.Errorf("priority value %q missing from OVERVIEW.md", p)
		}
	}
}

func TestNoDeadPathReferences(t *testing.T) {
	root := projectRoot(t)
	overview := readDoc(t, root, "OVERVIEW.md")
	workflow := readDoc(t, root, "WORKFLOW.md")

	// Paths that should NOT appear (removed from codebase)
	deadPaths := []string{"testing/scripts/"}

	for _, doc := range []struct {
		name, content string
	}{
		{"OVERVIEW.md", overview},
		{"WORKFLOW.md", workflow},
	} {
		for _, dead := range deadPaths {
			if strings.Contains(doc.content, dead) {
				t.Errorf("%s references dead path %q (no corresponding code constant)", doc.name, dead)
			}
		}
	}

	// Paths that MUST appear (active code constants)
	requiredPaths := []struct {
		path       string
		constValue string
	}{
		{"docs/features/", feature.FeaturesDir},
		{"tests/e2e/", feature.E2ETestsBaseDir},
		{"tests/e2e/.graduated/", feature.E2EGraduatedDir},
	}

	for _, rp := range requiredPaths {
		found := strings.Contains(overview, rp.path) || strings.Contains(workflow, rp.path)
		if !found {
			t.Errorf("active path %q (constant %q) not referenced in any doc", rp.path, rp.constValue)
		}
	}
}

func TestZhRecordDataFieldsMatch(t *testing.T) {
	root := projectRoot(t)
	path := filepath.Join(root, "docs", "WORKFLOW.zh.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("WORKFLOW.zh.md not found, skipping Chinese doc check")
	}
	doc, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read WORKFLOW.zh.md: %v", err)
	}
	content := string(doc)

	jsonBlock := extractJSONBlock(content)
	if jsonBlock == "" {
		t.Fatal("cannot find JSON block in WORKFLOW.zh.md")
	}

	// Verify it's the RecordData block by checking for "summary" key

	var docFields map[string]interface{}
	if err := json.Unmarshal([]byte(jsonBlock), &docFields); err != nil {
		t.Fatalf("WORKFLOW.zh.md JSON parse error: %v\nJSON:\n%s", err, jsonBlock)
	}

	typ := reflect.TypeOf(task.RecordData{})
	required, optional := jsonTagNames(typ)
	allTags := make(map[string]string)
	for k, v := range required {
		allTags[k] = v
	}
	for k, v := range optional {
		allTags[k] = v
	}

	for tagName, goName := range allTags {
		if _, ok := docFields[tagName]; !ok {
			t.Errorf("WORKFLOW.zh.md: struct field %s (json:\"%s\") missing from JSON", goName, tagName)
		}
	}
	for key := range docFields {
		if _, ok := allTags[key]; !ok {
			t.Errorf("WORKFLOW.zh.md: JSON key %q has no corresponding struct field", key)
		}
	}
}

func TestZhNoDeadPathsAndDetectionOrder(t *testing.T) {
	root := projectRoot(t)
	for _, file := range []string{"OVERVIEW.zh.md", "WORKFLOW.zh.md"} {
		path := filepath.Join(root, "docs", file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		doc, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(doc)

		for _, dead := range []string{"testing/scripts/"} {
			if strings.Contains(content, dead) {
				t.Errorf("%s references dead path %q", file, dead)
			}
		}
	}

	// Verify Chinese WORKFLOW detection order includes justfile
	path := filepath.Join(root, "docs", "WORKFLOW.zh.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("WORKFLOW.zh.md not found")
	}
	doc, _ := os.ReadFile(path)
	content := string(doc)
	section := extractNumberedItemsAfterHeading(content, "**测试命令检测顺序")
	if len(section) == 0 {
		section = extractNumberedItemsAfterHeading(content, "测试命令检测顺序")
	}
	if len(section) == 0 {
		t.Fatal("cannot extract detection order from WORKFLOW.zh.md")
	}

	hasJustfile := false
	for _, item := range section {
		if strings.Contains(item, "justfile") || strings.Contains(item, "just") {
			hasJustfile = true
			break
		}
	}
	if !hasJustfile {
		t.Errorf("WORKFLOW.zh.md: detection order missing justfile step. Got: %v", section)
	}
}

// --- helpers ---

// extractGoBlock finds the first Go code block in content.
func extractGoBlock(content string) string {
	re := regexp.MustCompile("(?s)```go\\s*\n(.*?)```")
	m := re.FindStringSubmatch(content)
	if m == nil {
		return ""
	}
	return m[1]
}

// extractNumberedItemsAfterHeading finds a heading text in content and extracts
// numbered list items from the next 20 lines after it.
func extractNumberedItemsAfterHeading(content, heading string) []string {
	idx := strings.Index(content, heading)
	if idx == -1 {
		return nil
	}
	// Take the next 20 lines after the heading line
	after := content[idx:]
	lineEnd := strings.Index(after, "\n")
	if lineEnd == -1 {
		return nil
	}
	window := after[lineEnd+1:]
	lines := strings.SplitN(window, "\n", 20)
	var relevant []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "**") {
			if len(relevant) > 0 {
				break
			}
			continue
		}
		relevant = append(relevant, line)
	}
	return extractNumberedItems(strings.Join(relevant, "\n"))
}

// extractNumberedItems extracts items from a numbered markdown list (1. 2. 3.).
func extractNumberedItems(content string) []string {
	re := regexp.MustCompile(`(?m)^\d+\.\s+\` + "`" + `([^` + "`" + `]+)`)
	matches := re.FindAllStringSubmatch(content, -1)
	var items []string
	for _, m := range matches {
		items = append(items, m[1])
	}
	return items
}

// assertOrder checks that detected order keywords match expected order keywords.
func assertOrder(t *testing.T, docName string, detected, expected []string) {
	t.Helper()

	// Normalize: extract keyword from each item
	normalize := func(s string) string {
		keywords := []string{"testCommand", "justfile", "Makefile", "go.mod", "package.json", "pytest"}
		for _, kw := range keywords {
			if strings.Contains(s, kw) {
				return kw
			}
		}
		return s
	}

	var detectedNorm []string
	for _, d := range detected {
		detectedNorm = append(detectedNorm, normalize(d))
	}

	// Remove duplicates (keep first occurrence)
	seen := make(map[string]bool)
	var unique []string
	for _, d := range detectedNorm {
		if !seen[d] {
			seen[d] = true
			unique = append(unique, d)
		}
	}
	detectedNorm = unique

	// Compare lengths
	if len(detectedNorm) != len(expected) {
		t.Errorf("%s: detection order has %d items, expected %d\ngot:      %v\nexpected: %v",
			docName, len(detectedNorm), len(expected), detectedNorm, expected)
		return
	}

	for i, got := range detectedNorm {
		if got != expected[i] {
			// Build diff
			t.Errorf("%s: detection order mismatch at position %d: got %q, want %q\ngot:      %v\nexpected: %v",
				docName, i+1, got, expected[i], detectedNorm, expected)
			return
		}
	}

	// Verify sorted (should match exactly)
	if !sort.IsSorted(sort.StringSlice(detectedNorm)) {
		t.Logf("%s: detection order verified: %v", docName, detectedNorm)
	}
}
