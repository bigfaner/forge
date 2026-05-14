// Package proposal provides discovery and parsing for project proposals.
package proposal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"

	"gopkg.in/yaml.v3"
)

// Metadata holds the parsed frontmatter fields from a proposal.md file.
type Metadata struct {
	Created string `yaml:"created"`
	Author  string `yaml:"author"`
	Status  string `yaml:"status"`
}

// Proposal represents a discovered proposal with its metadata and derived info.
type Proposal struct {
	Slug          string
	Created       string
	Status        string
	Author        string
	HasPRD        bool
	FeatureStatus string
	FilePath      string
}

// Discover walks docs/proposals/*/proposal.md and returns all proposals.
// Created date reads from frontmatter `created` field, falls back to file birth time.
func Discover(projectRoot string) ([]Proposal, error) {
	proposalsDir := filepath.Join(projectRoot, feature.ProposalBaseDir)
	entries, err := os.ReadDir(proposalsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read proposals directory: %w", err)
	}

	var proposals []Proposal
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		proposalFile := filepath.Join(proposalsDir, entry.Name(), feature.ProposalFileName)
		data, err := os.ReadFile(proposalFile)
		if err != nil {
			continue // skip directories without proposal.md
		}

		var meta Metadata
		if err := parseFrontmatter(data, &meta); err != nil {
			continue // skip malformed frontmatter
		}

		slug := entry.Name()
		created := meta.Created
		if created == "" {
			info, err := os.Stat(proposalFile)
			if err == nil {
				// Use ModTime as fallback (birth time is not portable)
				created = info.ModTime().Format("2006-01-02")
			}
		}

		// Check PRD existence
		prdPath := filepath.Join(projectRoot, feature.FeaturesDir, slug, feature.PRDDirName, feature.PRDSpecFile)
		hasPRD := fileExists(prdPath)

		// Check feature status
		featureStatus := ""
		manifestPath := filepath.Join(projectRoot, feature.FeaturesDir, slug, feature.ManifestFileName)
		if data, err := os.ReadFile(manifestPath); err == nil {
			var manifestMeta struct {
				Status string `yaml:"status"`
			}
			if err := parseFrontmatter(data, &manifestMeta); err == nil {
				featureStatus = manifestMeta.Status
			}
		}

		proposals = append(proposals, Proposal{
			Slug:          slug,
			Created:       created,
			Status:        meta.Status,
			Author:        meta.Author,
			HasPRD:        hasPRD,
			FeatureStatus: featureStatus,
			FilePath:      filepath.Join(feature.ProposalBaseDir, slug, feature.ProposalFileName),
		})
	}

	return proposals, nil
}

// FindBySlug returns a single proposal by slug, or an error if not found.
func FindBySlug(projectRoot, slug string) (*Proposal, error) {
	proposals, err := Discover(projectRoot)
	if err != nil {
		return nil, err
	}
	for _, p := range proposals {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("proposal not found: %s", slug)
}

// parseFrontmatter extracts YAML frontmatter from markdown content.
func parseFrontmatter(content []byte, target any) error {
	text := string(content)

	// Find opening ---
	if !strings.HasPrefix(text, "---") {
		return nil
	}
	text = text[3:]

	// Find closing ---
	closeIdx := strings.Index(text, "\n---")
	if closeIdx < 0 {
		return nil
	}

	yamlContent := text[:closeIdx]
	return yaml.Unmarshal([]byte(yamlContent), target)
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
