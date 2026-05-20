package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	_ = w.Close()
	os.Stdout = old
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func captureStderr2(f func()) string {
	var buf bytes.Buffer
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	f()
	_ = w.Close()
	os.Stderr = old
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestPrintBlock(t *testing.T) {
	out := captureStdout(func() { PrintBlock("KEY", "value") })
	if !strings.Contains(out, "---") {
		t.Errorf("expected separator, got: %s", out)
	}
	if !strings.Contains(out, "KEY: value") {
		t.Errorf("expected key-value, got: %s", out)
	}
}

func TestPrintFields(t *testing.T) {
	out := captureStdout(func() {
		PrintFields("K1", "v1", "K2", "v2")
	})
	if !strings.Contains(out, "K1: v1") || !strings.Contains(out, "K2: v2") {
		t.Errorf("expected both pairs, got: %s", out)
	}
}

func TestPrintFields_PanicsOnOddArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for odd number of args")
		}
	}()
	PrintFields("K1", "v1", "K2") //nolint:staticcheck // intentionally odd to test panic
}

func TestPrintResult(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		details string
		want    string
	}{
		{"with details", "FAIL", "2 errors", "RESULT: FAIL (2 errors)"},
		{"without details", "PASS", "", "RESULT: PASS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureStdout(func() { PrintResult(tt.status, tt.details) })
			if !strings.Contains(out, tt.want) {
				t.Errorf("expected %q, got: %s", tt.want, out)
			}
		})
	}
}

func TestPrintWarning(t *testing.T) {
	out := captureStdout(func() { PrintWarning("careful") })
	if !strings.Contains(out, "WARNING: careful") {
		t.Errorf("expected warning prefix, got: %s", out)
	}
}

func TestDebugf(t *testing.T) {
	t.Run("verbose true", func(t *testing.T) {
		out := captureStderr2(func() { Debugf(true, "val=%d", 42) })
		if !strings.Contains(out, "[debug] val=42") {
			t.Errorf("expected debug output, got: %s", out)
		}
	})
	t.Run("verbose false", func(t *testing.T) {
		out := captureStderr2(func() { Debugf(false, "val=%d", 42) })
		if out != "" {
			t.Errorf("expected no output, got: %s", out)
		}
	})
}

func TestPrintField(t *testing.T) {
	out := captureStdout(func() { PrintField("KEY", "val") })
	if !strings.Contains(out, "KEY: val") {
		t.Errorf("expected KEY: val, got: %s", out)
	}
}

func TestPrintFieldIfNotEmpty(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		out := captureStdout(func() { PrintFieldIfNotEmpty("K", "v") })
		if !strings.Contains(out, "K: v") {
			t.Errorf("expected output, got: %s", out)
		}
	})
	t.Run("empty", func(t *testing.T) {
		out := captureStdout(func() { PrintFieldIfNotEmpty("K", "") })
		if strings.Contains(out, "K:") {
			t.Errorf("expected no output for empty value, got: %s", out)
		}
	})
}

func TestPrintFieldIfNotEmptySlice(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		out := captureStdout(func() { PrintFieldIfNotEmptySlice("K", []string{"a", "b"}) })
		if !strings.Contains(out, "K: a, b") {
			t.Errorf("expected output, got: %s", out)
		}
	})
	t.Run("empty", func(t *testing.T) {
		out := captureStdout(func() { PrintFieldIfNotEmptySlice("K", nil) })
		if strings.Contains(out, "K:") {
			t.Errorf("expected no output, got: %s", out)
		}
	})
}

func TestPrintListItem(t *testing.T) {
	out := captureStdout(func() { PrintListItem("item") })
	if !strings.HasPrefix(out, "  ") {
		t.Errorf("expected indented output, got: %s", out)
	}
}

func TestPrintSection(t *testing.T) {
	out := captureStdout(func() { PrintSection("ERRORS") })
	if !strings.Contains(out, "[ERRORS]") {
		t.Errorf("expected section header, got: %s", out)
	}
}

// Suppress unused import
var _ = fmt.Sprintf
