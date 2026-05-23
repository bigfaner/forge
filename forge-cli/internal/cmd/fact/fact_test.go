package fact

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/facttable"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// helper to create a temp project root with .forge/ dir
func newFactTestRoot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	forgeDir := filepath.Join(dir, ".forge")
	err := os.MkdirAll(forgeDir, 0o755)
	assert.NoError(t, err)
	return dir
}

func writeFactTable(t *testing.T, root string, table facttable.FactTable) {
	t.Helper()
	data, err := json.MarshalIndent(table, "", "  ")
	assert.NoError(t, err)
	err = os.WriteFile(facttable.FactFilePath(root), data, 0o644)
	assert.NoError(t, err)
}

func makeEntry(id, source, subject, kind, confidence, updatedAt string, value interface{}) *facttable.FactEntry {
	valBytes, _ := json.Marshal(value)
	return &facttable.FactEntry{
		FactID:     id,
		Source:     source,
		Subject:    subject,
		Kind:       kind,
		Value:      valBytes,
		Confidence: confidence,
		UpdatedAt:  updatedAt,
	}
}

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	// Reset global flag variables before each test execution
	listSource = ""
	listConfidence = ""
	// Reset parsed flags so cobra re-parses from args
	_ = cmd.Flags().Set("source", "")
	_ = cmd.Flags().Set("confidence", "")
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

// --- List tests ---

func TestList_Empty(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	out, err := executeCommand(listCmd)
	assert.NoError(t, err)
	assert.Contains(t, out, "no facts found")
}

func TestList_All(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	writeFactTable(t, root, facttable.FactTable{
		makeEntry("s-k-001", facttable.SourceStatic, "cli.forge", facttable.KindSignature, facttable.ConfidenceConfirmed, "2026-01-01T00:00:00Z", "val"),
		makeEntry("r-k-001", facttable.SourceRuntime, "api.GET /tasks", facttable.KindOutputFormat, facttable.ConfidenceInferred, "2026-01-02T00:00:00Z", "val"),
	})

	out, err := executeCommand(listCmd)
	assert.NoError(t, err)
	assert.Contains(t, out, "2 facts found")
	assert.Contains(t, out, "s-k-001")
	assert.Contains(t, out, "r-k-001")
	assert.Contains(t, out, "static")
	assert.Contains(t, out, "runtime")
}

func TestList_FilterBySource(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	writeFactTable(t, root, facttable.FactTable{
		makeEntry("s-1", facttable.SourceStatic, "a", facttable.KindSignature, facttable.ConfidenceConfirmed, "t1", "v"),
		makeEntry("r-1", facttable.SourceRuntime, "b", facttable.KindSignature, facttable.ConfidenceInferred, "t2", "v"),
	})

	out, err := executeCommand(listCmd, "--source", "runtime")
	assert.NoError(t, err)
	assert.Contains(t, out, "1 facts found")
	assert.Contains(t, out, "r-1")
	assert.NotContains(t, out, "s-1")
}

func TestList_FilterByConfidence(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	writeFactTable(t, root, facttable.FactTable{
		makeEntry("s-1", facttable.SourceStatic, "a", facttable.KindSignature, facttable.ConfidenceConfirmed, "t1", "v"),
		makeEntry("s-2", facttable.SourceStatic, "b", facttable.KindSignature, facttable.ConfidenceInferred, "t2", "v"),
	})

	out, err := executeCommand(listCmd, "--confidence", "confirmed")
	assert.NoError(t, err)
	assert.Contains(t, out, "1 facts found")
	assert.Contains(t, out, "s-1")
	assert.NotContains(t, out, "s-2")
}

func TestList_NoMatch(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	writeFactTable(t, root, facttable.FactTable{
		makeEntry("s-1", facttable.SourceStatic, "a", facttable.KindSignature, facttable.ConfidenceConfirmed, "t1", "v"),
	})

	out, err := executeCommand(listCmd, "--source", "runtime")
	assert.NoError(t, err)
	assert.Contains(t, out, "no facts found")
}

// --- Get tests ---

func TestGet_Found(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	writeFactTable(t, root, facttable.FactTable{
		makeEntry("cli-sig-001", facttable.SourceStatic, "cli.forge", facttable.KindSignature, facttable.ConfidenceConfirmed, "2026-01-01T00:00:00Z", map[string]interface{}{"params": []string{"arg1", "arg2"}}),
	})

	out, err := executeCommand(getCmd, "cli-sig-001")
	assert.NoError(t, err)
	assert.Contains(t, out, "FACT_ID:     cli-sig-001")
	assert.Contains(t, out, "SOURCE:      static")
	assert.Contains(t, out, "SUBJECT:     cli.forge")
	assert.Contains(t, out, "KIND:        signature")
	assert.Contains(t, out, "CONFIDENCE:  confirmed")
	assert.Contains(t, out, `"params"`)
}

func TestGet_NotFound(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	writeFactTable(t, root, facttable.FactTable{
		makeEntry("cli-sig-001", facttable.SourceStatic, "cli.forge", facttable.KindSignature, facttable.ConfidenceConfirmed, "t", "v"),
	})

	_, err := executeCommand(getCmd, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Fact not found")
}

func TestGet_CorruptFile(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	err := os.WriteFile(facttable.FactFilePath(root), []byte("{bad}"), 0o644)
	assert.NoError(t, err)

	_, err = executeCommand(getCmd, "any-id")
	assert.Error(t, err)
}

func TestGet_NoArgs(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	_, err := executeCommand(getCmd)
	assert.Error(t, err)
}

// --- Summary tests ---

func TestSummary_Empty(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	out, err := executeCommand(summaryCmd)
	assert.NoError(t, err)
	assert.Contains(t, out, "TOTAL: 0 facts")
}

func TestSummary_WithData(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	writeFactTable(t, root, facttable.FactTable{
		makeEntry("1", facttable.SourceStatic, "a", facttable.KindSignature, facttable.ConfidenceConfirmed, "t1", "v"),
		makeEntry("2", facttable.SourceStatic, "b", facttable.KindSignature, facttable.ConfidenceInferred, "t2", "v"),
		makeEntry("3", facttable.SourceRuntime, "c", facttable.KindOutputFormat, facttable.ConfidenceConfirmed, "t3", "v"),
	})

	out, err := executeCommand(summaryCmd)
	assert.NoError(t, err)
	assert.Contains(t, out, "TOTAL: 3 facts")
	assert.Contains(t, out, "[BY SOURCE]")
	assert.Contains(t, out, "[BY CONFIDENCE]")
	assert.Contains(t, out, "[BY KIND]")
	assert.True(t, strings.Contains(out, "STATIC") || strings.Contains(out, "static"), "should show static source count")
}

func TestSummary_CorruptFile(t *testing.T) {
	root := newFactTestRoot(t)
	t.Setenv("CLAUDE_PROJECT_DIR", root)

	err := os.WriteFile(facttable.FactFilePath(root), []byte("not json"), 0o644)
	assert.NoError(t, err)

	_, err = executeCommand(summaryCmd)
	assert.Error(t, err)
}
