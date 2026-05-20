package e2e

import (
	"testing"
)

func TestRunOpts(t *testing.T) {
	opts := RunOpts{
		ProjectRoot: "/tmp/project",
		Feature:     "my-feature",
		Force:       true,
	}
	if opts.ProjectRoot != "/tmp/project" {
		t.Fatalf("unexpected ProjectRoot: %q", opts.ProjectRoot)
	}
	if opts.Feature != "my-feature" {
		t.Fatalf("unexpected Feature: %q", opts.Feature)
	}
	if !opts.Force {
		t.Fatal("expected Force to be true")
	}
}
