package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcSlugColWidth(t *testing.T) {
	tests := []struct {
		name     string
		slugLens []int
		want     int
	}{
		{
			name:     "empty list defaults to minimum",
			slugLens: nil,
			want:     30,
		},
		{
			name:     "short slugs use minimum 30",
			slugLens: []int{5, 10, 15},
			want:     30,
		},
		{
			name:     "max slug 42 gives 44 (42+2), clamped to 44",
			slugLens: []int{10, 20, 42},
			want:     44,
		},
		{
			name:     "max slug exactly 28 gives 30 (minimum)",
			slugLens: []int{28},
			want:     30,
		},
		{
			name:     "max slug 60 gives max 60",
			slugLens: []int{60},
			want:     60,
		},
		{
			name:     "max slug 65 clamped to max 60",
			slugLens: []int{65},
			want:     60,
		},
		{
			name:     "max slug 58 gives 60 (58+2=60)",
			slugLens: []int{58},
			want:     60,
		},
		{
			name:     "max slug 59 clamped to 60",
			slugLens: []int{59},
			want:     60,
		},
		{
			name:     "max slug 50 gives 52",
			slugLens: []int{50},
			want:     52,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcSlugColWidth(tt.slugLens)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProposalList_DynamicWidth(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	// Create a proposal with a 42-char slug
	longSlug := "profile-aware-shared-infra-precise-staging" // 42 chars
	for _, slug := range []string{"short", longSlug} {
		proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
		require.NoError(t, os.MkdirAll(proposalDir, 0755))
		content := "---\ncreated: 2026-01-01\nstatus: Draft\n---\n"
		require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"proposal"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)

	// The long slug should appear untruncated in output
	assert.Contains(t, output, longSlug)

	// Verify header, separator and data lines all contain the same slug column width
	// by checking the separator line has dashes matching the expected dynamic width (44)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Find the separator line: starts with spaces then dashes
		if afterPrefix, ok := strings.CutPrefix(line, "  "); ok {
			if strings.HasPrefix(afterPrefix, "---") && !strings.Contains(afterPrefix, "CREATED") {
				// First dash group should be slugWidth=44 chars (42+2=44)
				firstDashes := strings.SplitN(afterPrefix, " ", 2)[0]
				assert.Equal(t, 44, len(firstDashes), "slug separator should be 44 dashes (42+2)")
				break
			}
		}
	}
}

func TestLessonList_DynamicWidth(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	// Create lessons with different name lengths
	longName := "profile-aware-shared-infra-precise-staging" // 42 chars
	content := "---\ndate: 2026-01-01\ntags: [testing]\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, longName+".md"), []byte(content), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "short.md"), []byte(content), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"lesson"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)

	// The long name should appear untruncated
	assert.Contains(t, output, longName)
	assert.Contains(t, output, "short")
}

func TestFeatureList_DynamicWidth(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	// Create features with different slug lengths
	longSlug := "profile-aware-shared-infra-precise-staging" // 42 chars
	for _, slug := range []string{"short", longSlug} {
		featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
		require.NoError(t, os.MkdirAll(featureDir, 0755))
		manifest := "---\nstatus: active\ncreated: 2026-01-01\n---\n"
		require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifest), 0644))
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "list"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)

	// The long slug should appear untruncated
	assert.Contains(t, output, longSlug)
	assert.Contains(t, output, "short")
}

func TestTruncateSlug_WithDynamicWidth(t *testing.T) {
	// When slug exceeds dynamic width, it should still truncate
	dynamicWidth := 44 // e.g., maxSlugLen=42, 42+2=44
	longSlug := "this-is-an-extremely-long-slug-that-exceeds-sixty-chars-and-should-be-truncated"
	result := truncateSlug(longSlug, dynamicWidth)
	assert.Equal(t, dynamicWidth, len(result))
	assert.True(t, strings.HasSuffix(result, "..."))
}
