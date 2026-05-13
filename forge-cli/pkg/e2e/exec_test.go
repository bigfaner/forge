package e2e

import (
	"fmt"
	"strings"
	"testing"
)

// stubExec is a hand-rolled mock matching project convention (no testify/gomock).
type stubExec struct {
	responses map[string]execResponse
}

type execResponse struct {
	output []byte
	err    error
}

func (s *stubExec) Run(name string, args ...string) ([]byte, error) {
	key := name + " " + strings.Join(args, " ")
	if r, ok := s.responses[key]; ok {
		return r.output, r.err
	}
	return nil, fmt.Errorf("stubExec: unexpected command: %s", key)
}

func TestStubExec(t *testing.T) {
	t.Run("returns configured response", func(t *testing.T) {
		s := &stubExec{responses: map[string]execResponse{
			"echo hello": {output: []byte("hello\n"), err: nil},
		}}

		out, err := s.Run("echo", "hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(out) != "hello\n" {
			t.Fatalf("expected 'hello\\n', got %q", string(out))
		}
	})

	t.Run("returns error for unexpected command", func(t *testing.T) {
		s := &stubExec{responses: map[string]execResponse{}}

		_, err := s.Run("unknown", "cmd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "stubExec: unexpected command") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns configured error", func(t *testing.T) {
		s := &stubExec{responses: map[string]execResponse{
			"fail cmd": {output: nil, err: fmt.Errorf("command failed")},
		}}

		_, err := s.Run("fail", "cmd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "command failed" {
			t.Fatalf("expected 'command failed', got %q", err.Error())
		}
	})
}

func TestRealExecImplementsExecRunner(_ *testing.T) {
	// Compile-time interface check
	var _ ExecRunner = RealExec{}
}

func TestRunnerDefault(t *testing.T) {
	// Verify runner is set to RealExec by default
	_, ok := runner.(RealExec)
	if !ok {
		t.Fatal("expected runner to be RealExec by default")
	}
}
