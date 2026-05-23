package infocmd

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testItem is a minimal struct for testing generic functions.
type testItem struct {
	Slug     string
	Created  string
	FilePath string
}

// idBySlug extracts ID from Slug field.
func idBySlug(t testItem) string { return t.Slug }

// --- ParseFrontmatter tests ---

func TestParseFrontmatter_ValidYAML(t *testing.T) {
	content := []byte("---\ncreated: 2026-01-15\ntopic: test\n---\nBody text")
	var result struct {
		Created string `yaml:"created"`
		Topic   string `yaml:"topic"`
	}
	err := ParseFrontmatter(content, &result)
	assert.NoError(t, err)
	assert.Equal(t, "2026-01-15", result.Created)
	assert.Equal(t, "test", result.Topic)
}

func TestParseFrontmatter_NoOpeningDelimiter(t *testing.T) {
	content := []byte("no frontmatter here\n---\nmore text")
	var result struct {
		Created string `yaml:"created"`
	}
	err := ParseFrontmatter(content, &result)
	assert.NoError(t, err)
	assert.Equal(t, "", result.Created)
}

func TestParseFrontmatter_NoClosingDelimiter(t *testing.T) {
	content := []byte("---\ncreated: 2026-01-15\nno closing delimiter")
	var result struct {
		Created string `yaml:"created"`
	}
	err := ParseFrontmatter(content, &result)
	assert.NoError(t, err)
	assert.Equal(t, "", result.Created)
}

func TestParseFrontmatter_EmptyContent(t *testing.T) {
	content := []byte("")
	var result struct {
		Created string `yaml:"created"`
	}
	err := ParseFrontmatter(content, &result)
	assert.NoError(t, err)
}

func TestParseFrontmatter_InvalidYAML(t *testing.T) {
	content := []byte("---\n: invalid yaml [[[\n---\nbody")
	var result struct {
		Created string `yaml:"created"`
	}
	err := ParseFrontmatter(content, &result)
	assert.Error(t, err)
}

// --- Discover tests ---

func TestDiscover_FlatMode(t *testing.T) {
	dir := t.TempDir()
	researchDir := filepath.Join(dir, "docs", "research")
	assert.NoError(t, os.MkdirAll(researchDir, 0755))

	// Create two .md files
	assert.NoError(t, os.WriteFile(filepath.Join(researchDir, "report-a.md"), []byte("---\ncreated: 2026-01-10\ntopic: alpha\n---\nbody"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(researchDir, "report-b.md"), []byte("---\ncreated: 2026-02-20\ntopic: beta\n---\nbody"), 0644))
	// Non-.md file should be skipped
	assert.NoError(t, os.WriteFile(filepath.Join(researchDir, "ignore.txt"), []byte("text"), 0644))

	type researchMeta struct {
		Created string `yaml:"created"`
		Topic   string `yaml:"topic"`
	}

	type researchItem struct {
		Slug    string
		Created string
		Topic   string
	}

	cfg := ScanConfig[researchItem]{
		BaseDir:    "docs/research",
		IsSubdir:   false,
		IDKey:      func(r researchItem) string { return r.Slug },
		CreatedKey: func(r researchItem) string { return r.Created },
		ParseEntry: func(name, _ string, content []byte, _ time.Time) (researchItem, error) {
			var meta researchMeta
			if err := ParseFrontmatter(content, &meta); err != nil {
				return researchItem{}, err
			}
			return researchItem{
				Slug:    name,
				Created: meta.Created,
				Topic:   meta.Topic,
			}, nil
		},
	}

	items, err := Discover(dir, cfg)
	assert.NoError(t, err)
	assert.Len(t, items, 2)

	// Sorted by created descending: report-b (2026-02-20) first
	assert.Equal(t, "report-b", items[0].Slug)
	assert.Equal(t, "report-a", items[1].Slug)
}

func TestDiscover_SubdirMode(t *testing.T) {
	dir := t.TempDir()
	proposalsDir := filepath.Join(dir, "docs", "proposals")

	// Create two proposal subdirectories
	for _, slug := range []string{"my-proposal", "other-proposal"} {
		pDir := filepath.Join(proposalsDir, slug)
		assert.NoError(t, os.MkdirAll(pDir, 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(pDir, "proposal.md"), []byte("---\ncreated: 2026-03-01\nstatus: draft\n---\nbody"), 0644))
	}
	// Non-directory entries should be skipped
	assert.NoError(t, os.WriteFile(filepath.Join(proposalsDir, "README.md"), []byte("readme"), 0644))

	type proposalItem struct {
		Slug    string
		Created string
		Status  string
	}

	type proposalMeta struct {
		Created string `yaml:"created"`
		Status  string `yaml:"status"`
	}

	cfg := ScanConfig[proposalItem]{
		BaseDir:    "docs/proposals",
		IsSubdir:   true,
		FileName:   "proposal.md",
		IDKey:      func(p proposalItem) string { return p.Slug },
		CreatedKey: func(p proposalItem) string { return p.Created },
		ParseEntry: func(name, _ string, content []byte, _ time.Time) (proposalItem, error) {
			var meta proposalMeta
			if err := ParseFrontmatter(content, &meta); err != nil {
				return proposalItem{}, err
			}
			return proposalItem{
				Slug:    name,
				Created: meta.Created,
				Status:  meta.Status,
			}, nil
		},
	}

	items, err := Discover(dir, cfg)
	assert.NoError(t, err)
	assert.Len(t, items, 2)

	slugs := []string{items[0].Slug, items[1].Slug}
	sort.Strings(slugs)
	assert.Contains(t, slugs, "my-proposal")
	assert.Contains(t, slugs, "other-proposal")
}

func TestDiscover_DirNotExist(t *testing.T) {
	dir := t.TempDir()

	cfg := ScanConfig[testItem]{
		BaseDir:    "docs/nonexistent",
		IsSubdir:   false,
		IDKey:      idBySlug,
		CreatedKey: func(t testItem) string { return t.Created },
		ParseEntry: func(name, _ string, _ []byte, _ time.Time) (testItem, error) {
			return testItem{Slug: name}, nil
		},
	}

	items, err := Discover(dir, cfg)
	assert.NoError(t, err)
	assert.Empty(t, items)
}

func TestDiscover_SkipInvalidFrontmatter(t *testing.T) {
	dir := t.TempDir()
	testDir := filepath.Join(dir, "docs", "test")
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	// Valid file
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "valid.md"), []byte("---\ncreated: 2026-01-01\n---\nbody"), 0644))
	// Invalid frontmatter file (should be skipped)
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "bad.md"), []byte("---\n: broken [[[\n---\nbody"), 0644))

	type simpleItem struct {
		Slug    string
		Created string
	}

	type simpleMeta struct {
		Created string `yaml:"created"`
	}

	cfg := ScanConfig[simpleItem]{
		BaseDir:    "docs/test",
		IsSubdir:   false,
		IDKey:      func(s simpleItem) string { return s.Slug },
		CreatedKey: func(s simpleItem) string { return s.Created },
		ParseEntry: func(name, _ string, content []byte, _ time.Time) (simpleItem, error) {
			var meta simpleMeta
			if err := ParseFrontmatter(content, &meta); err != nil {
				return simpleItem{}, err
			}
			return simpleItem{Slug: name, Created: meta.Created}, nil
		},
	}

	items, err := Discover(dir, cfg)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "valid", items[0].Slug)
}

func TestDiscover_SortByCreatedDescWithMtimeFallback(t *testing.T) {
	dir := t.TempDir()
	testDir := filepath.Join(dir, "docs", "items")
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	// File with created field
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "newest.md"), []byte("---\ncreated: 2026-05-01\n---\nbody"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "older.md"), []byte("---\ncreated: 2026-01-01\n---\nbody"), 0644))
	// File without created field (will use mtime)
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "no-date.md"), []byte("---\n---\nbody"), 0644))

	type sortItem struct {
		Slug    string
		Created string
	}

	type sortMeta struct {
		Created string `yaml:"created"`
	}

	cfg := ScanConfig[sortItem]{
		BaseDir:    "docs/items",
		IsSubdir:   false,
		IDKey:      func(s sortItem) string { return s.Slug },
		CreatedKey: func(s sortItem) string { return s.Created },
		ParseEntry: func(name, _ string, content []byte, _ time.Time) (sortItem, error) {
			var meta sortMeta
			if err := ParseFrontmatter(content, &meta); err != nil {
				return sortItem{}, err
			}
			return sortItem{Slug: name, Created: meta.Created}, nil
		},
	}

	items, err := Discover(dir, cfg)
	assert.NoError(t, err)
	assert.Len(t, items, 3)

	// Items with created field sort first (descending), items without go last
	assert.Equal(t, "newest", items[0].Slug)  // 2026-05-01
	assert.Equal(t, "older", items[1].Slug)   // 2026-01-01
	assert.Equal(t, "no-date", items[2].Slug) // no created, falls to mtime
}

// --- FindByID tests ---

func TestFindByID_Found(t *testing.T) {
	dir := t.TempDir()
	testDir := filepath.Join(dir, "docs", "items")
	assert.NoError(t, os.MkdirAll(testDir, 0755))
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "target.md"), []byte("---\ncreated: 2026-01-01\n---\nbody"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "other.md"), []byte("---\ncreated: 2026-02-01\n---\nbody"), 0644))

	type item struct {
		Slug    string
		Created string
	}

	type itemMeta struct {
		Created string `yaml:"created"`
	}

	cfg := ScanConfig[item]{
		BaseDir:    "docs/items",
		IsSubdir:   false,
		IDKey:      func(i item) string { return i.Slug },
		CreatedKey: func(i item) string { return i.Created },
		ParseEntry: func(name, _ string, content []byte, _ time.Time) (item, error) {
			var meta itemMeta
			if err := ParseFrontmatter(content, &meta); err != nil {
				return item{}, err
			}
			return item{Slug: name, Created: meta.Created}, nil
		},
	}

	result, err := FindByID(dir, "target", cfg)
	assert.NoError(t, err)
	assert.Equal(t, "target", result.Slug)
}

func TestFindByID_NotFound(t *testing.T) {
	dir := t.TempDir()
	testDir := filepath.Join(dir, "docs", "items")
	assert.NoError(t, os.MkdirAll(testDir, 0755))
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "exists.md"), []byte("---\n---\nbody"), 0644))

	type item struct {
		Slug string
	}

	cfg := ScanConfig[item]{
		BaseDir:    "docs/items",
		IsSubdir:   false,
		IDKey:      func(i item) string { return i.Slug },
		CreatedKey: func(_ item) string { return "" },
		ParseEntry: func(name, _ string, _ []byte, _ time.Time) (item, error) {
			return item{Slug: name}, nil
		},
	}

	_, err := FindByID(dir, "nonexistent", cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestFindByID_WithName(t *testing.T) {
	dir := t.TempDir()
	testDir := filepath.Join(dir, "docs", "items")
	assert.NoError(t, os.MkdirAll(testDir, 0755))
	assert.NoError(t, os.WriteFile(filepath.Join(testDir, "my-lesson.md"), []byte("---\n---\nbody"), 0644))

	type item struct {
		Name string
	}

	cfg := ScanConfig[item]{
		BaseDir:    "docs/items",
		IsSubdir:   false,
		IDKey:      func(i item) string { return i.Name },
		CreatedKey: func(_ item) string { return "" },
		ParseEntry: func(name, _ string, _ []byte, _ time.Time) (item, error) {
			return item{Name: name}, nil
		},
	}

	result, err := FindByID(dir, "my-lesson", cfg)
	assert.NoError(t, err)
	assert.Equal(t, "my-lesson", result.Name)
}
