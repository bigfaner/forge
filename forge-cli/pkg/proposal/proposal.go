// Package proposal provides discovery and parsing for project proposals.
package proposal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/infocmd"
)

// Proposal represents a discovered proposal with its metadata and derived info.
type Proposal struct {
	Slug          string
	Created       string
	Status        string
	Author        string
	Intent        string
	HasPRD        bool
	FeatureStatus string
	FilePath      string
}

// metadata holds the parsed frontmatter fields from a proposal.md file.
type metadata struct {
	Created string `yaml:"created"`
	Author  string `yaml:"author"`
	Status  string `yaml:"status"`
	Intent  string `yaml:"intent"`
}

// scanConfig returns the infocmd.ScanConfig for proposal discovery.
func scanConfig(projectRoot string) infocmd.ScanConfig[Proposal] {
	return infocmd.ScanConfig[Proposal]{
		BaseDir:    feature.ProposalBaseDir,
		IsSubdir:   true,
		FileName:   feature.ProposalFileName,
		IDKey:      func(p Proposal) string { return p.Slug },
		CreatedKey: func(p Proposal) string { return p.Created },
		ParseEntry: func(slug, _ string, content []byte, modTime time.Time) (Proposal, error) {
			var meta metadata
			if err := infocmd.ParseFrontmatter(content, &meta); err != nil {
				return Proposal{}, err
			}

			created := meta.Created
			if created == "" {
				created = modTime.Format("2006-01-02")
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
				if err := infocmd.ParseFrontmatter(data, &manifestMeta); err == nil {
					featureStatus = manifestMeta.Status
				}
			}

			return Proposal{
				Slug:          slug,
				Created:       created,
				Status:        meta.Status,
				Author:        meta.Author,
				Intent:        meta.Intent,
				HasPRD:        hasPRD,
				FeatureStatus: featureStatus,
				FilePath:      filepath.Join(feature.ProposalBaseDir, slug, feature.ProposalFileName),
			}, nil
		},
	}
}

// Discover walks docs/proposals/*/proposal.md and returns all proposals.
// Created date reads from frontmatter `created` field, falls back to file birth time.
// Results are sorted by Created descending (newest first).
func Discover(projectRoot string) ([]Proposal, error) {
	return infocmd.Discover(projectRoot, scanConfig(projectRoot))
}

// FindBySlug returns a single proposal by slug, or an error if not found.
func FindBySlug(projectRoot, slug string) (*Proposal, error) {
	p, err := infocmd.FindByID(projectRoot, slug, scanConfig(projectRoot))
	if err != nil {
		return nil, fmt.Errorf("proposal not found: %s", slug)
	}
	return p, nil
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
