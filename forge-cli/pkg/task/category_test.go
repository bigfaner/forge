package task

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryForType_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		typ      string
		expected string
	}{
		// Coding category (coding.* prefix)
		{"coding.feature", TypeCodingFeature, CategoryCoding},
		{"coding.enhancement", TypeCodingEnhancement, CategoryCoding},
		{"coding.cleanup", TypeCodingCleanup, CategoryCoding},
		{"coding.refactor", TypeCodingRefactor, CategoryCoding},
		{"coding.fix", TypeCodingFix, CategoryCoding},
		// Doc category (doc* prefix)
		{"doc", TypeDoc, CategoryDoc},
		{"doc.fix", TypeDocFix, CategoryDoc},
		{"doc.review", TypeDocReview, CategoryDoc},
		{"doc.summary", TypeDocSummary, CategoryDoc},
		{"doc.consolidate", TypeDocConsolidate, CategoryDoc},
		{"doc.drift", TypeDocDrift, CategoryDoc},
		// Test category (test.* prefix)
		{"test.gen-scripts", TypeTestGenScripts, CategoryTest},
		{"test.run", TypeTestRun, CategoryTest},
		// Validation category (validation.* prefix)
		{"validation.code", TypeValidationCode, CategoryValidation},
		{"validation.ux", TypeValidationUx, CategoryValidation},
		// Eval category (eval.* prefix)
		{"eval.journey", TypeEvalJourney, CategoryEval},
		{"eval.contract", TypeEvalContract, CategoryEval},
		// Gate category (exact match)
		{"gate", TypeGate, CategoryGate},
		// code-quality.simplify maps to coding (explicit match)
		{"code-quality.simplify maps to coding", TypeCleanCode, CategoryCoding},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CategoryForType(tt.typ))
		})
	}
}

func TestCategoryForType_DefaultAndUnknown(t *testing.T) {
	assert.Equal(t, CategoryCoding, CategoryForType(""), "empty string defaults to coding")
	assert.Equal(t, CategoryCoding, CategoryForType("unknown.type"), "unknown type defaults to coding")
	assert.Equal(t, CategoryCoding, CategoryForType("totally-invalid"), "invalid type defaults to coding")
}

func TestCategoryForType_UnknownLogsWarning(t *testing.T) {
	// Capture stderr output (forgelog dispatches to stderr when no backends registered)
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	result := CategoryForType("unknown.type")
	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stderr = oldStderr

	assert.Equal(t, CategoryCoding, result)
	assert.Contains(t, buf.String(), `CategoryForType: unknown type "unknown.type", defaulting to coding`)
}

func TestCategoryForType_KnownTypeNoWarning(t *testing.T) {
	// Capture stderr output (forgelog dispatches to stderr when no backends registered)
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	_ = CategoryForType(TypeEvalJourney)
	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stderr = oldStderr

	assert.NotContains(t, buf.String(), "CategoryForType: unknown type")
}

func TestCategoryConstants(t *testing.T) {
	assert.Equal(t, "coding", CategoryCoding)
	assert.Equal(t, "doc", CategoryDoc)
	assert.Equal(t, "test", CategoryTest)
	assert.Equal(t, "validation", CategoryValidation)
	assert.Equal(t, "gate", CategoryGate)
	assert.Equal(t, "eval", CategoryEval)
}

func TestRecordData_NewFieldsOmitEmpty(t *testing.T) {
	rd := RecordData{
		TaskID: "1",
		Status: "completed",
	}
	data, err := json.Marshal(rd)
	assert.NoError(t, err)

	// When no new fields are set, they should not appear in JSON (omitempty)
	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &m))

	// Old fields present
	assert.Equal(t, "1", m["taskId"])
	assert.Equal(t, "completed", m["status"])

	// New fields absent (omitempty)
	assert.NotContains(t, m, "referencedDocs")
	assert.NotContains(t, m, "reviewStatus")
	assert.NotContains(t, m, "docMetrics")
	assert.NotContains(t, m, "casesGenerated")
	assert.NotContains(t, m, "casesEvaluated")
	assert.NotContains(t, m, "scriptsCreated")
	assert.NotContains(t, m, "testResults")
	assert.NotContains(t, m, "validationPassed")
	assert.NotContains(t, m, "issuesFound")
	assert.NotContains(t, m, "gatePassed")
	assert.NotContains(t, m, "gateChecks")
	assert.NotContains(t, m, "score")
	assert.NotContains(t, m, "findings")
	assert.NotContains(t, m, "severity")
	assert.NotContains(t, m, "passed")
}

func TestRecordData_BackwardCompatibility(t *testing.T) {
	// Old JSON (without new fields) should deserialize cleanly
	oldJSON := `{
		"taskId": "42",
		"status": "completed",
		"summary": "did stuff",
		"filesCreated": ["a.go"],
		"filesModified": ["b.go"],
		"keyDecisions": ["decided X"],
		"testsPassed": 5,
		"testsFailed": 0,
		"coverage": 85.5,
		"notes": "all good"
	}`

	var rd RecordData
	assert.NoError(t, json.Unmarshal([]byte(oldJSON), &rd))

	assert.Equal(t, "42", rd.TaskID)
	assert.Equal(t, "completed", rd.Status)
	assert.Equal(t, "did stuff", rd.Summary)
	assert.Equal(t, []string{"a.go"}, rd.FilesCreated)
	assert.Equal(t, []string{"b.go"}, rd.FilesModified)
	assert.Equal(t, []string{"decided X"}, rd.KeyDecisions)
	assert.Equal(t, 5, rd.TestsPassed)
	assert.Equal(t, 0, rd.TestsFailed)
	assert.Equal(t, 85.5, rd.Coverage)
	assert.Equal(t, "all good", rd.Notes)

	// New fields should be zero values
	assert.Nil(t, rd.ReferencedDocs)
	assert.Equal(t, "", rd.ReviewStatus)
	assert.Nil(t, rd.ScriptsCreated)
	assert.Equal(t, 0, rd.CasesGenerated)
	assert.Equal(t, 0, rd.CasesEvaluated)
	assert.Equal(t, "", rd.TestResults)
	assert.False(t, rd.ValidationPassed)
	assert.Nil(t, rd.IssuesFound)
	assert.False(t, rd.GatePassed)
	assert.Nil(t, rd.GateChecks)
}

func TestRecordData_NewFieldsRoundTrip(t *testing.T) {
	rd := RecordData{
		TaskID:           "10",
		Status:           "completed",
		ReferencedDocs:   []string{"doc1.md", "doc2.md"},
		ReviewStatus:     "approved",
		DocMetrics:       "50% coverage",
		CasesGenerated:   12,
		CasesEvaluated:   10,
		ScriptsCreated:   []string{"test1.sh", "test2.sh"},
		TestResults:      "10 passed, 2 failed",
		ValidationPassed: true,
		IssuesFound:      []string{"issue1", "issue2"},
		GatePassed:       true,
		GateChecks:       []string{"lint", "compile"},
		Score:            850,
		Findings:         []string{"finding1", "finding2"},
		Severity:         "major",
		Passed:           true,
	}

	data, err := json.Marshal(rd)
	assert.NoError(t, err)

	var rd2 RecordData
	assert.NoError(t, json.Unmarshal(data, &rd2))

	assert.Equal(t, rd.TaskID, rd2.TaskID)
	assert.Equal(t, rd.ReferencedDocs, rd2.ReferencedDocs)
	assert.Equal(t, rd.ReviewStatus, rd2.ReviewStatus)
	assert.Equal(t, rd.DocMetrics, rd2.DocMetrics)
	assert.Equal(t, rd.CasesGenerated, rd2.CasesGenerated)
	assert.Equal(t, rd.CasesEvaluated, rd2.CasesEvaluated)
	assert.Equal(t, rd.ScriptsCreated, rd2.ScriptsCreated)
	assert.Equal(t, rd.TestResults, rd2.TestResults)
	assert.Equal(t, rd.ValidationPassed, rd2.ValidationPassed)
	assert.Equal(t, rd.IssuesFound, rd2.IssuesFound)
	assert.Equal(t, rd.GatePassed, rd2.GatePassed)
	assert.Equal(t, rd.GateChecks, rd2.GateChecks)
	assert.Equal(t, rd.Score, rd2.Score)
	assert.Equal(t, rd.Findings, rd2.Findings)
	assert.Equal(t, rd.Severity, rd2.Severity)
	assert.Equal(t, rd.Passed, rd2.Passed)
}
