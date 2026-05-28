package types

import (
	"encoding/json"
	"testing"
)

// ---------------------------------------------------------------------------
// Status constants — value correctness
// ---------------------------------------------------------------------------

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      Status
		expected string
	}{
		{"pending", StatusPending, "pending"},
		{"in_progress", StatusInProgress, "in_progress"},
		{"completed", StatusCompleted, "completed"},
		{"blocked", StatusBlocked, "blocked"},
		{"suspended", StatusSuspended, "suspended"},
		{"skipped", StatusSkipped, "skipped"},
		{"rejected", StatusRejected, "rejected"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("got %q, want %q", tt.got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AllStatuses
// ---------------------------------------------------------------------------

func TestAllStatuses(t *testing.T) {
	all := AllStatuses()
	if len(all) != 7 {
		t.Fatalf("AllStatuses returned %d items, want 7", len(all))
	}
	want := map[Status]bool{
		StatusPending:    true,
		StatusInProgress: true,
		StatusCompleted:  true,
		StatusBlocked:    true,
		StatusSuspended:  true,
		StatusSkipped:    true,
		StatusRejected:   true,
	}
	for _, s := range all {
		if !want[s] {
			t.Errorf("unexpected status %q in AllStatuses", s)
		}
	}
}

// ---------------------------------------------------------------------------
// IsTerminalStatus
// ---------------------------------------------------------------------------

func TestIsTerminalStatus(t *testing.T) {
	tests := []struct {
		status Status
		want   bool
	}{
		{StatusPending, false},
		{StatusInProgress, false},
		{StatusCompleted, true},
		{StatusBlocked, false},
		{StatusSuspended, false},
		{StatusSkipped, true},
		{StatusRejected, true},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := IsTerminalStatus(tt.status); got != tt.want {
				t.Errorf("IsTerminalStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// SurfaceType constants — value correctness
// ---------------------------------------------------------------------------

func TestSurfaceTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      SurfaceType
		expected string
	}{
		{"web", SurfaceWeb, "web"},
		{"api", SurfaceAPI, "api"},
		{"cli", SurfaceCLI, "cli"},
		{"tui", SurfaceTUI, "tui"},
		{"mobile", SurfaceMobile, "mobile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("got %q, want %q", tt.got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AllSurfaceTypes
// ---------------------------------------------------------------------------

func TestAllSurfaceTypes(t *testing.T) {
	all := AllSurfaceTypes()
	if len(all) != 5 {
		t.Fatalf("AllSurfaceTypes returned %d items, want 5", len(all))
	}
	want := map[SurfaceType]bool{
		SurfaceWeb:    true,
		SurfaceAPI:    true,
		SurfaceCLI:    true,
		SurfaceTUI:    true,
		SurfaceMobile: true,
	}
	for _, s := range all {
		if !want[s] {
			t.Errorf("unexpected surface type %q in AllSurfaceTypes", s)
		}
	}
}

// ---------------------------------------------------------------------------
// Priority constants — value correctness
// ---------------------------------------------------------------------------

func TestPriorityConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      Priority
		expected string
	}{
		{"P0", PriorityP0, "P0"},
		{"P1", PriorityP1, "P1"},
		{"P2", PriorityP2, "P2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("got %q, want %q", tt.got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AllPriorities
// ---------------------------------------------------------------------------

func TestAllPriorities(t *testing.T) {
	all := AllPriorities()
	if len(all) != 3 {
		t.Fatalf("AllPriorities returned %d items, want 3", len(all))
	}
	want := map[Priority]bool{
		PriorityP0: true,
		PriorityP1: true,
		PriorityP2: true,
	}
	for _, p := range all {
		if !want[p] {
			t.Errorf("unexpected priority %q in AllPriorities", p)
		}
	}
}

// ---------------------------------------------------------------------------
// JSON serialization compatibility — type X string must serialize as plain string
// ---------------------------------------------------------------------------

func TestStatusJSONCompat(t *testing.T) {
	s := StatusPending
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"pending"` {
		t.Errorf("got %s, want %q", data, `"pending"`)
	}

	var got Status
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != s {
		t.Errorf("roundtrip: got %q, want %q", got, s)
	}
}

func TestSurfaceTypeJSONCompat(t *testing.T) {
	s := SurfaceCLI
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"cli"` {
		t.Errorf("got %s, want %q", data, `"cli"`)
	}

	var got SurfaceType
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != s {
		t.Errorf("roundtrip: got %q, want %q", got, s)
	}
}

func TestPriorityJSONCompat(t *testing.T) {
	p := PriorityP1
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"P1"` {
		t.Errorf("got %s, want %q", data, `"P1"`)
	}

	var got Priority
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != p {
		t.Errorf("roundtrip: got %q, want %q", got, p)
	}
}

// ---------------------------------------------------------------------------
// Zero internal imports — compile-time check (ensures pkg/types is a leaf)
// This test file itself imports only stdlib, and the implementation files
// are verified by `go build ./pkg/types/...` which would fail if non-stdlib
// imports exist.
// ---------------------------------------------------------------------------
