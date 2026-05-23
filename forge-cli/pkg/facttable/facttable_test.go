package facttable

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to create a temp project root with .forge/ dir
func newTestProjectRoot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	forgeDir := filepath.Join(dir, ".forge")
	err := os.MkdirAll(forgeDir, 0o755)
	assert.NoError(t, err)
	return dir
}

func makeFactEntry(id, source, subject, kind, confidence, updatedAt string, value interface{}) *FactEntry {
	valBytes, _ := json.Marshal(value)
	return &FactEntry{
		FactID:     id,
		Source:     source,
		Subject:    subject,
		Kind:       kind,
		Value:      valBytes,
		Confidence: confidence,
		UpdatedAt:  updatedAt,
	}
}

// --- Load tests ---

func TestLoad_EmptyTable(t *testing.T) {
	root := newTestProjectRoot(t)

	table, err := Load(root)
	assert.NoError(t, err)
	assert.Empty(t, table)
}

func TestLoad_FileNotExist(t *testing.T) {
	root := t.TempDir()
	// No .forge dir at all

	table, err := Load(root)
	assert.NoError(t, err)
	assert.Empty(t, table)
}

func TestLoad_ValidFile(t *testing.T) {
	root := newTestProjectRoot(t)
	entries := FactTable{
		makeFactEntry("sub-kind-1", SourceStatic, "sub", KindSignature, ConfidenceConfirmed, "2026-01-01T00:00:00Z", map[string]string{"key": "val"}),
	}
	data, _ := json.MarshalIndent(entries, "", "  ")
	err := os.WriteFile(FactFilePath(root), data, 0o644)
	assert.NoError(t, err)

	table, err := Load(root)
	assert.NoError(t, err)
	assert.Len(t, table, 1)
	assert.Equal(t, "sub-kind-1", table[0].FactID)
}

func TestLoad_CorruptJSON(t *testing.T) {
	root := newTestProjectRoot(t)
	err := os.WriteFile(FactFilePath(root), []byte("{bad json}"), 0o644)
	assert.NoError(t, err)

	table, err := Load(root)
	assert.Nil(t, table)
	assert.Error(t, err)

	corrupt, ok := err.(*CorruptError)
	assert.True(t, ok, "expected CorruptError")
	assert.Contains(t, corrupt.Hint(), "Fix JSON syntax")
}

func TestLoad_EmptyFile(t *testing.T) {
	root := newTestProjectRoot(t)
	err := os.WriteFile(FactFilePath(root), []byte(""), 0o644)
	assert.NoError(t, err)

	table, err := Load(root)
	assert.NoError(t, err)
	assert.Empty(t, table)
}

// --- Save tests ---

func TestSave_CreatesDir(t *testing.T) {
	root := t.TempDir()
	// No .forge dir

	table := FactTable{
		makeFactEntry("a-b-1", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "2026-01-01T00:00:00Z", "val"),
	}

	err := Save(root, table)
	assert.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(FactFilePath(root))
	assert.NoError(t, err)
}

func TestSave_RoundTrip(t *testing.T) {
	root := newTestProjectRoot(t)

	original := FactTable{
		makeFactEntry("s-k-1", SourceStatic, "s", KindSignature, ConfidenceConfirmed, "2026-01-01T00:00:00Z", map[string]int{"n": 1}),
		makeFactEntry("s-k-2", SourceRuntime, "s", KindOutputFormat, ConfidenceInferred, "2026-01-02T00:00:00Z", map[string]string{"format": "json"}),
	}

	err := Save(root, original)
	assert.NoError(t, err)

	loaded, err := Load(root)
	assert.NoError(t, err)
	assert.Len(t, loaded, 2)
	assert.Equal(t, original[0].FactID, loaded[0].FactID)
	assert.Equal(t, original[1].FactID, loaded[1].FactID)
}

// --- Filter tests ---

func TestFilter_NoFilters(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceInferred, "t2", "v"),
	}

	result := table.Filter("", "")
	assert.Len(t, result, 2)
}

func TestFilter_BySource(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceInferred, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindErrorCode, ConfidenceAssumed, "t3", "v"),
	}

	result := table.Filter(SourceRuntime, "")
	assert.Len(t, result, 2)
	assert.Equal(t, SourceRuntime, result[0].Source)
	assert.Equal(t, SourceRuntime, result[1].Source)
}

func TestFilter_ByConfidence(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceInferred, "t2", "v"),
	}

	result := table.Filter("", ConfidenceConfirmed)
	assert.Len(t, result, 1)
	assert.Equal(t, ConfidenceConfirmed, result[0].Confidence)
}

func TestFilter_ByBoth(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceInferred, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindErrorCode, ConfidenceConfirmed, "t3", "v"),
	}

	result := table.Filter(SourceRuntime, ConfidenceConfirmed)
	assert.Len(t, result, 1)
	assert.Equal(t, "3", result[0].FactID)
}

func TestFilter_NoMatch(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
	}

	result := table.Filter(SourceRuntime, "")
	assert.Empty(t, result)
}

// --- GetByID tests ---

func TestGetByID_Found(t *testing.T) {
	table := FactTable{
		makeFactEntry("abc-123", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
	}

	entry := table.GetByID("abc-123")
	assert.NotNil(t, entry)
	assert.Equal(t, "abc-123", entry.FactID)
}

func TestGetByID_NotFound(t *testing.T) {
	table := FactTable{
		makeFactEntry("abc-123", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
	}

	entry := table.GetByID("nonexistent")
	assert.Nil(t, entry)
}

// --- Summary tests ---

func TestSummary_Empty(t *testing.T) {
	table := FactTable{}

	stats := table.Summary()
	assert.Equal(t, 0, stats.Total)
	assert.Empty(t, stats.BySource)
	assert.Empty(t, stats.ByConfidence)
	assert.Empty(t, stats.ByKind)
}

func TestSummary_Counts(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceStatic, "b", KindSignature, ConfidenceInferred, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindOutputFormat, ConfidenceConfirmed, "t3", "v"),
	}

	stats := table.Summary()
	assert.Equal(t, 3, stats.Total)
	assert.Equal(t, 2, stats.BySource[SourceStatic])
	assert.Equal(t, 1, stats.BySource[SourceRuntime])
	assert.Equal(t, 2, stats.ByConfidence[ConfidenceConfirmed])
	assert.Equal(t, 1, stats.ByConfidence[ConfidenceInferred])
	assert.Equal(t, 2, stats.ByKind[KindSignature])
	assert.Equal(t, 1, stats.ByKind[KindOutputFormat])
}

// --- EffectiveEntries tests ---

func TestEffectiveEntries_RuntimeConfirmedReplacesStatic(t *testing.T) {
	table := FactTable{
		makeFactEntry("s1", SourceStatic, "cli.forge", KindSignature, ConfidenceInferred, "t1", "static-val"),
		makeFactEntry("r1", SourceRuntime, "cli.forge", KindSignature, ConfidenceConfirmed, "t2", "runtime-val"),
	}

	effective := table.EffectiveEntries()
	assert.Len(t, effective, 1)
	assert.Equal(t, "r1", effective[0].FactID)
	assert.Equal(t, SourceRuntime, effective[0].Source)
}

func TestEffectiveEntries_RuntimeNonConfirmedKeepsFallback(t *testing.T) {
	table := FactTable{
		makeFactEntry("s1", SourceStatic, "cli.forge", KindSignature, ConfidenceInferred, "t1", "static-val"),
		makeFactEntry("r1", SourceRuntime, "cli.forge", KindSignature, ConfidenceInferred, "t2", "runtime-val"),
	}

	effective := table.EffectiveEntries()
	assert.Len(t, effective, 2)

	// Runtime entry is present
	runtimeFound := false
	staticFound := false
	for _, e := range effective {
		if e.Source == SourceRuntime {
			runtimeFound = true
		}
		if e.Source == SourceStatic {
			staticFound = true
		}
	}
	assert.True(t, runtimeFound, "runtime entry should be in effective entries")
	assert.True(t, staticFound, "static entry should be kept as fallback")
}

func TestEffectiveEntries_NoRuntimeKeepsStatic(t *testing.T) {
	table := FactTable{
		makeFactEntry("s1", SourceStatic, "cli.forge", KindSignature, ConfidenceInferred, "t1", "static-val"),
	}

	effective := table.EffectiveEntries()
	assert.Len(t, effective, 1)
	assert.Equal(t, SourceStatic, effective[0].Source)
}

func TestEffectiveEntries_DifferentKindsCoexist(t *testing.T) {
	table := FactTable{
		makeFactEntry("s1", SourceStatic, "cli.forge", KindSignature, ConfidenceConfirmed, "t1", "v1"),
		makeFactEntry("s2", SourceStatic, "cli.forge", KindOutputFormat, ConfidenceConfirmed, "t2", "v2"),
		makeFactEntry("r1", SourceRuntime, "cli.forge", KindSignature, ConfidenceConfirmed, "t3", "v3"),
	}

	effective := table.EffectiveEntries()
	assert.Len(t, effective, 2)

	// Kind=signature: runtime replaces static
	// Kind=output_format: no runtime, static kept
	kinds := map[string]string{}
	for _, e := range effective {
		kinds[e.Kind] = e.Source
	}
	assert.Equal(t, SourceRuntime, kinds[KindSignature])
	assert.Equal(t, SourceStatic, kinds[KindOutputFormat])
}

func TestEffectiveEntries_ManualSourcePreserved(t *testing.T) {
	table := FactTable{
		makeFactEntry("m1", SourceManual, "cli.forge", KindSignature, ConfidenceConfirmed, "t1", "manual-val"),
	}

	effective := table.EffectiveEntries()
	assert.Len(t, effective, 1)
	assert.Equal(t, SourceManual, effective[0].Source)
}

// --- Validate tests ---

func TestValidate_ValidEntry(t *testing.T) {
	entry := makeFactEntry("id", SourceStatic, "subj", KindSignature, ConfidenceConfirmed, "2026-01-01T00:00:00Z", "v")
	assert.NoError(t, entry.Validate())
}

func TestValidate_MissingFactID(t *testing.T) {
	entry := makeFactEntry("", SourceStatic, "subj", KindSignature, ConfidenceConfirmed, "t", "v")
	assert.Error(t, entry.Validate())
}

func TestValidate_InvalidSource(t *testing.T) {
	entry := makeFactEntry("id", "invalid", "subj", KindSignature, ConfidenceConfirmed, "t", "v")
	assert.Error(t, entry.Validate())
}

func TestValidate_MissingSubject(t *testing.T) {
	entry := makeFactEntry("id", SourceStatic, "", KindSignature, ConfidenceConfirmed, "t", "v")
	assert.Error(t, entry.Validate())
}

func TestValidate_InvalidKind(t *testing.T) {
	entry := makeFactEntry("id", SourceStatic, "subj", "invalid", ConfidenceConfirmed, "t", "v")
	assert.Error(t, entry.Validate())
}

func TestValidate_InvalidConfidence(t *testing.T) {
	entry := makeFactEntry("id", SourceStatic, "subj", KindSignature, "invalid", "t", "v")
	assert.Error(t, entry.Validate())
}

func TestValidate_MissingValue(t *testing.T) {
	entry := &FactEntry{
		FactID:     "id",
		Source:     SourceStatic,
		Subject:    "subj",
		Kind:       KindSignature,
		Confidence: ConfidenceConfirmed,
		UpdatedAt:  "t",
	}
	assert.Error(t, entry.Validate())
}

func TestValidate_MissingUpdatedAt(t *testing.T) {
	entry := makeFactEntry("id", SourceStatic, "subj", KindSignature, ConfidenceConfirmed, "", "v")
	assert.Error(t, entry.Validate())
}

// --- SortedEntries tests ---

func TestSortedEntries(t *testing.T) {
	table := FactTable{
		makeFactEntry("c", SourceStatic, "a", KindSignature, ConfidenceConfirmed, "t", "v"),
		makeFactEntry("a", SourceStatic, "b", KindSignature, ConfidenceConfirmed, "t", "v"),
		makeFactEntry("b", SourceStatic, "c", KindSignature, ConfidenceConfirmed, "t", "v"),
	}

	sorted := table.SortedEntries()
	assert.Equal(t, "a", sorted[0].FactID)
	assert.Equal(t, "b", sorted[1].FactID)
	assert.Equal(t, "c", sorted[2].FactID)

	// Original is unchanged
	assert.Equal(t, "c", table[0].FactID)
}

// --- GenerateNonce tests ---

func TestGenerateNonce(t *testing.T) {
	nonce := GenerateNonce()
	assert.NotEmpty(t, nonce)
}

// --- RuntimeCoverage tests (R2L support) ---

func TestConfirmedRuntimeSubjects_Empty(t *testing.T) {
	table := FactTable{}
	subjects := table.ConfirmedRuntimeSubjects()
	assert.Empty(t, subjects)
}

func TestConfirmedRuntimeSubjects_OnlyStatic(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "cli.forge", KindSignature, ConfidenceConfirmed, "t1", "v"),
	}
	subjects := table.ConfirmedRuntimeSubjects()
	assert.Empty(t, subjects)
}

func TestConfirmedRuntimeSubjects_RuntimeConfirmed(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "cli.forge task claim", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "cli.forge task submit", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "cli.forge task status", KindOutputFormat, ConfidenceInferred, "t3", "v"),
	}
	subjects := table.ConfirmedRuntimeSubjects()
	assert.Len(t, subjects, 2)
	assert.True(t, subjects["cli.forge task claim"])
	assert.True(t, subjects["cli.forge task submit"])
	assert.False(t, subjects["cli.forge task status"])
}

func TestRuntimeCoverageRatio_NoOutcomes(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
	}
	ratio := table.RuntimeCoverageRatio(nil)
	assert.Equal(t, 0.0, ratio)
}

func TestRuntimeCoverageRatio_FullCoverage(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "cli.forge task claim", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "cli.forge task submit", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
	}
	outcomes := []string{"cli.forge task claim", "cli.forge task submit"}
	ratio := table.RuntimeCoverageRatio(outcomes)
	assert.Equal(t, 1.0, ratio)
}

func TestRuntimeCoverageRatio_PartialCoverage(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "cli.forge task claim", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
	}
	outcomes := []string{"cli.forge task claim", "cli.forge task submit", "cli.forge task status"}
	ratio := table.RuntimeCoverageRatio(outcomes)

	// 1 out of 3 = 0.333...
	assert.InDelta(t, 0.333, ratio, 0.01)
}

func TestRuntimeCoverageRatio_ZeroCoverage(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceStatic, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
	}
	outcomes := []string{"cli.forge task claim", "cli.forge task submit"}
	ratio := table.RuntimeCoverageRatio(outcomes)
	assert.Equal(t, 0.0, ratio)
}

func TestRuntimeCoverageRatio_OnlyConfirmedCounts(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "cli.forge task claim", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "cli.forge task submit", KindOutputFormat, ConfidenceInferred, "t2", "v"),
	}
	outcomes := []string{"cli.forge task claim", "cli.forge task submit"}
	ratio := table.RuntimeCoverageRatio(outcomes)

	// Only "task claim" is confirmed runtime, "task submit" is inferred
	// 1 out of 2 = 0.5
	assert.Equal(t, 0.5, ratio)
}

func TestRuntimeCoverageRatio_EightyPercent(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindOutputFormat, ConfidenceConfirmed, "t3", "v"),
		makeFactEntry("4", SourceRuntime, "d", KindOutputFormat, ConfidenceConfirmed, "t4", "v"),
	}
	outcomes := []string{"a", "b", "c", "d", "e"}
	ratio := table.RuntimeCoverageRatio(outcomes)
	assert.Equal(t, 0.8, ratio)
}

// --- CorruptError tests ---

func TestCorruptError_Error(t *testing.T) {
	err := &CorruptError{Path: "/path/to/file", Details: "bad json"}
	assert.Contains(t, err.Error(), "corrupted fact table")
	assert.Contains(t, err.Error(), "bad json")
}

func TestCorruptError_Hint(t *testing.T) {
	err := &CorruptError{Path: "/path/to/file", Details: "bad json"}
	assert.Contains(t, err.Hint(), "Fix JSON syntax")
}

// --- Valid constants tests ---

func TestValidConstants(t *testing.T) {
	assert.Contains(t, ValidSources, SourceStatic)
	assert.Contains(t, ValidSources, SourceRuntime)
	assert.Contains(t, ValidSources, SourceManual)

	assert.Contains(t, ValidConfidences, ConfidenceConfirmed)
	assert.Contains(t, ValidConfidences, ConfidenceInferred)
	assert.Contains(t, ValidConfidences, ConfidenceAssumed)

	assert.Contains(t, ValidKinds, KindSignature)
	assert.Contains(t, ValidKinds, KindOutputFormat)
	assert.Contains(t, ValidKinds, KindErrorCode)
	assert.Contains(t, ValidKinds, KindSideEffect)
	assert.Contains(t, ValidKinds, KindPrecondition)
	assert.Contains(t, ValidKinds, KindCompilationError)
	assert.Contains(t, ValidKinds, KindRuntimeCrash)
}
