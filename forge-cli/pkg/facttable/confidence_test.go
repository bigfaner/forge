package facttable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- ConfidenceRating model tests ---

func TestConfidenceRating_Fields(t *testing.T) {
	cr := ConfidenceRating{
		Level:              ConfidenceHigh,
		ConfirmedFactRatio: 0.85,
		TotalOutcomes:      10,
		ConfirmedOutcomes:  9,
		EvalSkipped:        false,
		EvalBypassed:       false,
	}
	assert.Equal(t, ConfidenceHigh, cr.Level)
	assert.InDelta(t, 0.85, cr.ConfirmedFactRatio, 0.001)
	assert.Equal(t, 10, cr.TotalOutcomes)
	assert.Equal(t, 9, cr.ConfirmedOutcomes)
	assert.False(t, cr.EvalSkipped)
	assert.False(t, cr.EvalBypassed)
}

// --- ComputeConfidenceRating tests ---

func TestComputeConfidenceRating_High(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindOutputFormat, ConfidenceConfirmed, "t3", "v"),
		makeFactEntry("4", SourceRuntime, "d", KindOutputFormat, ConfidenceConfirmed, "t4", "v"),
	}
	outcomes := []string{"a", "b", "c", "d", "e"}

	cr := table.ComputeConfidenceRating(outcomes, false, false)
	assert.Equal(t, ConfidenceHigh, cr.Level)
	assert.InDelta(t, 0.80, cr.ConfirmedFactRatio, 0.001)
	assert.Equal(t, 5, cr.TotalOutcomes)
	assert.Equal(t, 4, cr.ConfirmedOutcomes)
}

func TestComputeConfidenceRating_Medium(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
	}
	outcomes := []string{"a", "b", "c", "d"}

	cr := table.ComputeConfidenceRating(outcomes, false, false)
	assert.Equal(t, ConfidenceMedium, cr.Level)
	assert.InDelta(t, 0.50, cr.ConfirmedFactRatio, 0.001)
	assert.Equal(t, 4, cr.TotalOutcomes)
	assert.Equal(t, 2, cr.ConfirmedOutcomes)
}

func TestComputeConfidenceRating_Low(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
	}
	outcomes := []string{"a", "b", "c"}

	cr := table.ComputeConfidenceRating(outcomes, false, false)
	assert.Equal(t, ConfidenceLow, cr.Level)
	assert.InDelta(t, 0.333, cr.ConfirmedFactRatio, 0.01)
	assert.Equal(t, 3, cr.TotalOutcomes)
	assert.Equal(t, 1, cr.ConfirmedOutcomes)
}

func TestComputeConfidenceRating_EvalSkipped_ForcesLow(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindOutputFormat, ConfidenceConfirmed, "t3", "v"),
		makeFactEntry("4", SourceRuntime, "d", KindOutputFormat, ConfidenceConfirmed, "t4", "v"),
	}
	outcomes := []string{"a", "b", "c", "d", "e"}

	// Ratio is 0.80 which would normally be HIGH, but eval_skipped forces LOW
	cr := table.ComputeConfidenceRating(outcomes, true, false)
	assert.Equal(t, ConfidenceLow, cr.Level)
	assert.True(t, cr.EvalSkipped)
	assert.InDelta(t, 0.80, cr.ConfirmedFactRatio, 0.001)
}

func TestComputeConfidenceRating_EvalBypassed_ForcesLow(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindOutputFormat, ConfidenceConfirmed, "t3", "v"),
		makeFactEntry("4", SourceRuntime, "d", KindOutputFormat, ConfidenceConfirmed, "t4", "v"),
	}
	outcomes := []string{"a", "b", "c", "d", "e"}

	// Ratio is 0.80 which would normally be HIGH, but eval_bypassed forces LOW
	cr := table.ComputeConfidenceRating(outcomes, false, true)
	assert.Equal(t, ConfidenceLow, cr.Level)
	assert.True(t, cr.EvalBypassed)
}

func TestComputeConfidenceRating_EmptyOutcomes(t *testing.T) {
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
	}

	cr := table.ComputeConfidenceRating(nil, false, false)
	assert.Equal(t, ConfidenceLow, cr.Level)
	assert.Equal(t, 0.0, cr.ConfirmedFactRatio)
	assert.Equal(t, 0, cr.TotalOutcomes)
	assert.Equal(t, 0, cr.ConfirmedOutcomes)
}

func TestComputeConfidenceRating_ExactBoundary_High(t *testing.T) {
	// Exactly 0.80 should be HIGH
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
		makeFactEntry("3", SourceRuntime, "c", KindOutputFormat, ConfidenceConfirmed, "t3", "v"),
		makeFactEntry("4", SourceRuntime, "d", KindOutputFormat, ConfidenceConfirmed, "t4", "v"),
	}
	outcomes := []string{"a", "b", "c", "d", "e"}

	cr := table.ComputeConfidenceRating(outcomes, false, false)
	assert.Equal(t, ConfidenceHigh, cr.Level)
	assert.InDelta(t, 0.80, cr.ConfirmedFactRatio, 0.001)
}

func TestComputeConfidenceRating_ExactBoundary_Medium(t *testing.T) {
	// Exactly 0.40 should be MEDIUM
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
		makeFactEntry("2", SourceRuntime, "b", KindOutputFormat, ConfidenceConfirmed, "t2", "v"),
	}
	outcomes := []string{"a", "b", "c", "d", "e"}

	cr := table.ComputeConfidenceRating(outcomes, false, false)
	assert.Equal(t, ConfidenceMedium, cr.Level)
	assert.InDelta(t, 0.40, cr.ConfirmedFactRatio, 0.001)
}

func TestComputeConfidenceRating_JustBelowMedium(t *testing.T) {
	// 0.3999... should be LOW (< 0.40 boundary)
	table := FactTable{
		makeFactEntry("1", SourceRuntime, "a", KindOutputFormat, ConfidenceConfirmed, "t1", "v"),
	}
	outcomes := []string{"a", "b", "c"}

	cr := table.ComputeConfidenceRating(outcomes, false, false)
	assert.Equal(t, ConfidenceLow, cr.Level)
}

// --- NeedsReview tests ---

func TestNeedsReview_LowReturnsTrue(t *testing.T) {
	cr := ConfidenceRating{Level: ConfidenceLow}
	assert.True(t, cr.NeedsReview())
}

func TestNeedsReview_MediumReturnsFalse(t *testing.T) {
	cr := ConfidenceRating{Level: ConfidenceMedium}
	assert.False(t, cr.NeedsReview())
}

func TestNeedsReview_HighReturnsFalse(t *testing.T) {
	cr := ConfidenceRating{Level: ConfidenceHigh}
	assert.False(t, cr.NeedsReview())
}

// --- VerifyMark tests ---

func TestVerifyMark_LowIsReview(t *testing.T) {
	cr := ConfidenceRating{Level: ConfidenceLow}
	assert.Equal(t, "REVIEW", cr.VerifyMark())
}

func TestVerifyMark_HighIsVerify(t *testing.T) {
	cr := ConfidenceRating{Level: ConfidenceHigh}
	assert.Equal(t, "VERIFY", cr.VerifyMark())
}

func TestVerifyMark_MediumIsVerify(t *testing.T) {
	cr := ConfidenceRating{Level: ConfidenceMedium}
	assert.Equal(t, "VERIFY", cr.VerifyMark())
}

// --- Confidence level constant tests ---

func TestConfidenceLevelConstants(t *testing.T) {
	assert.Equal(t, ConfidenceLevel("HIGH"), ConfidenceHigh)
	assert.Equal(t, ConfidenceLevel("MEDIUM"), ConfidenceMedium)
	assert.Equal(t, ConfidenceLevel("LOW"), ConfidenceLow)
}
