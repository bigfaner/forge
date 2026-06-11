// Package infocmd provides shared utilities for info-commands: generic directory
// scanning, frontmatter parsing, sorting, and lookup. It also provides concrete
// discovery functions for lessons and research reports built on top of the
// generic scanning framework.
//
// # Sub-domains
//
//   - Generic scanning framework: ScanConfig[T], Discover[T], FindByID[T],
//     ExtractFrontmatter, ParseFrontmatter -- reusable infrastructure for any
//     info-command that needs to discover markdown files with YAML frontmatter
//   - Lesson discovery: DiscoverLessons, FindLessonByName, Lesson type --
//     discovers lessons from docs/lessons/*.md with category inference
//   - Research discovery: DiscoverReports, FindReportBySlug, Report type --
//     discovers research reports from docs/research/*.md with topic/mode metadata
//
// # Responsibility Boundaries
//
// The generic scanning framework must remain domain-agnostic -- it must not
// contain any business logic specific to lessons, research, or proposals.
// Domain-specific parsing logic belongs in dedicated functions/types within
// this package (DiscoverLessons, DiscoverReports) or in consumer packages
// (pkg/proposal, pkg/task).
//
// This package is a leaf package: it imports only the Go standard library and
// gopkg.in/yaml.v3. It must not import any other forge-cli pkg/ package.
package infocmd
