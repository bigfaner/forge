package forgeconfig

import (
	"fmt"
	"strings"
	"testing"
)

// TestEvalSettings_Struct tests EvalTypeSettings and EvalSettings struct definitions.
func TestEvalSettings_Struct(t *testing.T) {
	t.Run("EvalTypeSettings has Target and Iterations *int fields", func(t *testing.T) {
		ts := EvalTypeSettings{}
		if ts.Target != nil {
			t.Error("expected Target nil by default")
		}
		if ts.Iterations != nil {
			t.Error("expected Iterations nil by default")
		}
		val := 900
		ts.Target = &val
		if *ts.Target != 900 {
			t.Errorf("expected Target 900, got %d", *ts.Target)
		}
	})

	t.Run("EvalSettings has 7 eval type fields", func(t *testing.T) {
		es := EvalSettings{}
		// All fields should be zero-value EvalTypeSettings with nil pointers
		if es.Proposal.Target != nil || es.Proposal.Iterations != nil {
			t.Error("expected Proposal fields nil by default")
		}
		if es.Prd.Target != nil || es.Prd.Iterations != nil {
			t.Error("expected Prd fields nil by default")
		}
		if es.Design.Target != nil || es.Design.Iterations != nil {
			t.Error("expected Design fields nil by default")
		}
		if es.Ui.Target != nil || es.Ui.Iterations != nil {
			t.Error("expected Ui fields nil by default")
		}
		if es.Journey.Target != nil || es.Journey.Iterations != nil {
			t.Error("expected Journey fields nil by default")
		}
		if es.Contract.Target != nil || es.Contract.Iterations != nil {
			t.Error("expected Contract fields nil by default")
		}
		if es.Consistency.Target != nil || es.Consistency.Iterations != nil {
			t.Error("expected Consistency fields nil by default")
		}
	})

	t.Run("Config has Eval *EvalSettings field", func(t *testing.T) {
		cfg := Config{}
		if cfg.Eval != nil {
			t.Error("expected Eval nil by default")
		}
	})
}

// TestGetConfigValue_EvalSettings tests eval settings get via reflection routing.
func TestGetConfigValue_EvalSettings(t *testing.T) {
	t.Run("eval.proposal.target returns correct integer when configured", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  proposal:\n    target: 900\n")
		val, err := GetConfigValue(dir, "eval.proposal.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "900" {
			t.Errorf("expected '900', got %q", val)
		}
	})

	t.Run("eval.proposal.target returns errKeyNotFound when not configured", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  proposal: {}\n")
		_, err := GetConfigValue(dir, "eval.proposal.target")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for nil *int, got %v", err)
		}
	})

	t.Run("eval.proposal.target returns errKeyNotFound with no eval block", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "eval.proposal.target")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("eval.journey.iterations returns correct integer when configured", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  journey:\n    iterations: 5\n")
		val, err := GetConfigValue(dir, "eval.journey.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "5" {
			t.Errorf("expected '5', got %q", val)
		}
	})

	t.Run("eval returns summary of configured types only", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  proposal:\n    target: 900\n  journey:\n    target: 850\n    iterations: 3\n")
		val, err := GetConfigValue(dir, "eval")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(val, "proposal:") {
			t.Errorf("expected 'proposal:' in summary, got %q", val)
		}
		if !strings.Contains(val, "journey:") {
			t.Errorf("expected 'journey:' in summary, got %q", val)
		}
	})

	t.Run("eval.design.target returns value", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  design:\n    target: 900\n")
		val, err := GetConfigValue(dir, "eval.design.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "900" {
			t.Errorf("expected '900', got %q", val)
		}
	})

	t.Run("eval.ui.target returns value", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  ui:\n    target: 950\n")
		val, err := GetConfigValue(dir, "eval.ui.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "950" {
			t.Errorf("expected '950', got %q", val)
		}
	})

	t.Run("eval.contract.target returns value", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  contract:\n    target: 850\n")
		val, err := GetConfigValue(dir, "eval.contract.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "850" {
			t.Errorf("expected '850', got %q", val)
		}
	})

	t.Run("eval.consistency.target returns value", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  consistency:\n    target: 900\n")
		val, err := GetConfigValue(dir, "eval.consistency.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "900" {
			t.Errorf("expected '900', got %q", val)
		}
	})

	t.Run("eval.prd.iterations returns value", func(t *testing.T) {
		dir := setupConfig(t, "eval:\n  prd:\n    iterations: 3\n")
		val, err := GetConfigValue(dir, "eval.prd.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "3" {
			t.Errorf("expected '3', got %q", val)
		}
	})
}

// TestSetConfigValue_EvalSettings tests eval settings set via reflection routing.
func TestSetConfigValue_EvalSettings(t *testing.T) {
	t.Run("set eval.proposal.target 850 writes and reads correctly", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "eval.proposal.target", "850"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.proposal.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "850" {
			t.Errorf("expected '850', got %q", val)
		}
	})

	t.Run("set eval.journey.iterations 5 writes and reads correctly", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "eval.journey.iterations", "5"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.journey.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "5" {
			t.Errorf("expected '5', got %q", val)
		}
	})

	t.Run("set eval.proposal.target persists across file read", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "eval.proposal.target", "900"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Eval == nil {
			t.Fatal("expected Eval non-nil after set")
		}
		if cfg.Eval.Proposal.Target == nil || *cfg.Eval.Proposal.Target != 900 {
			t.Errorf("expected Proposal.Target = 900, got %v", cfg.Eval.Proposal.Target)
		}
	})

	t.Run("set multiple eval fields coexist", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "eval.proposal.target", "900"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := SetConfigValue(dir, "eval.proposal.iterations", "3"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := SetConfigValue(dir, "eval.journey.target", "850"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify all values
		val, err := GetConfigValue(dir, "eval.proposal.target")
		if err != nil || val != "900" {
			t.Errorf("proposal.target: expected '900', got %q err=%v", val, err)
		}
		val, err = GetConfigValue(dir, "eval.proposal.iterations")
		if err != nil || val != "3" {
			t.Errorf("proposal.iterations: expected '3', got %q err=%v", val, err)
		}
		val, err = GetConfigValue(dir, "eval.journey.target")
		if err != nil || val != "850" {
			t.Errorf("journey.target: expected '850', got %q err=%v", val, err)
		}
	})

	t.Run("set non-integer value for eval target returns error", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "eval.proposal.target", "notanumber")
		if err == nil {
			t.Fatal("expected error for non-integer value")
		}
		if !strings.Contains(err.Error(), "expected integer") {
			t.Errorf("expected 'expected integer' in error, got %v", err)
		}
	})

	t.Run("set eval block rejected as non-leaf", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "eval", "900")
		if err == nil {
			t.Fatal("expected error for non-leaf set")
		}
		if !strings.Contains(err.Error(), "cannot set non-leaf key") {
			t.Errorf("expected 'cannot set non-leaf key' in error, got %v", err)
		}
	})

	t.Run("set eval.proposal rejected as non-leaf", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "eval.proposal", "900")
		if err == nil {
			t.Fatal("expected error for non-leaf set")
		}
		if !strings.Contains(err.Error(), "cannot set non-leaf key") {
			t.Errorf("expected 'cannot set non-leaf key' in error, got %v", err)
		}
	})

	t.Run("set eval.consistency.target overwrites existing value", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "eval.consistency.target", "900"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := SetConfigValue(dir, "eval.consistency.target", "850"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.consistency.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "850" {
			t.Errorf("expected '850' after overwrite, got %q", val)
		}
	})

	t.Run("set eval.prd.target and eval.prd.iterations independently", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "eval.prd.target", "800"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := SetConfigValue(dir, "eval.prd.iterations", "7"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		target, err := GetConfigValue(dir, "eval.prd.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != "800" {
			t.Errorf("expected '800', got %q", target)
		}
		iter, err := GetConfigValue(dir, "eval.prd.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if iter != "7" {
			t.Errorf("expected '7', got %q", iter)
		}
	})

	t.Run("set writes valid YAML readable by ReadConfig", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "eval.proposal.target", "850"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := SetConfigValue(dir, "eval.proposal.iterations", "5"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Eval == nil {
			t.Fatal("expected Eval non-nil")
		}
		if cfg.Eval.Proposal.Target == nil || *cfg.Eval.Proposal.Target != 850 {
			t.Errorf("Proposal.Target = %v, want 850", cfg.Eval.Proposal.Target)
		}
		if cfg.Eval.Proposal.Iterations == nil || *cfg.Eval.Proposal.Iterations != 5 {
			t.Errorf("Proposal.Iterations = %v, want 5", cfg.Eval.Proposal.Iterations)
		}
	})
}

// TestGetConfigValue_EvalSettings_PartialConfig tests nil *int pointer fallback behavior.
func TestGetConfigValue_EvalSettings_PartialConfig(t *testing.T) {
	t.Run("only target set, iterations nil returns errKeyNotFound", func(t *testing.T) {
		dir := t.TempDir()
		target := 850
		cfg := &Config{
			Eval: &EvalSettings{
				Proposal: EvalTypeSettings{Target: &target},
			},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.proposal.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "850" {
			t.Errorf("expected '850', got %q", val)
		}
		_, err = GetConfigValue(dir, "eval.proposal.iterations")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for nil *int iterations, got %v", err)
		}
	})

	t.Run("only iterations set, target nil returns errKeyNotFound", func(t *testing.T) {
		dir := t.TempDir()
		iter := 5
		cfg := &Config{
			Eval: &EvalSettings{
				Journey: EvalTypeSettings{Iterations: &iter},
			},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.journey.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "5" {
			t.Errorf("expected '5', got %q", val)
		}
		_, err = GetConfigValue(dir, "eval.journey.target")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for nil *int target, got %v", err)
		}
	})

	t.Run("full config both target and iterations set", func(t *testing.T) {
		dir := t.TempDir()
		target := 900
		iter := 3
		cfg := &Config{
			Eval: &EvalSettings{
				Contract: EvalTypeSettings{Target: &target, Iterations: &iter},
			},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		tVal, err := GetConfigValue(dir, "eval.contract.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tVal != "900" {
			t.Errorf("expected '900', got %q", tVal)
		}
		iVal, err := GetConfigValue(dir, "eval.contract.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if iVal != "3" {
			t.Errorf("expected '3', got %q", iVal)
		}
	})

	t.Run("eval not configured at all returns errKeyNotFound for all paths", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		for _, key := range []string{
			"eval.proposal.target",
			"eval.proposal.iterations",
			"eval.ui.target",
			"eval.consistency.iterations",
		} {
			_, err := GetConfigValue(dir, key)
			if err != errKeyNotFound {
				t.Errorf("expected errKeyNotFound for %q, got %v", key, err)
			}
		}
	})
}

// TestGetConfigValue_EvalSettings_AllTypes verifies all 7 eval types via GetConfigValue with defaults.
func TestGetConfigValue_EvalSettings_AllTypes(t *testing.T) {
	evalTypes := []struct {
		name       string
		target     int
		iterations int
	}{
		{"proposal", 900, 3},
		{"prd", 900, 3},
		{"design", 900, 3},
		{"ui", 950, 3},
		{"journey", 850, 3},
		{"contract", 850, 3},
		{"consistency", 900, 3},
	}

	dir := t.TempDir()
	defaults := EvalSettingsDefaults()
	cfg := &Config{Eval: &defaults}
	if err := writeConfig(dir, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, tt := range evalTypes {
		t.Run(fmt.Sprintf("eval.%s.target=%d", tt.name, tt.target), func(t *testing.T) {
			val, err := GetConfigValue(dir, fmt.Sprintf("eval.%s.target", tt.name))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if val != fmt.Sprintf("%d", tt.target) {
				t.Errorf("expected '%d', got %q", tt.target, val)
			}
		})
		t.Run(fmt.Sprintf("eval.%s.iterations=%d", tt.name, tt.iterations), func(t *testing.T) {
			val, err := GetConfigValue(dir, fmt.Sprintf("eval.%s.iterations", tt.name))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if val != fmt.Sprintf("%d", tt.iterations) {
				t.Errorf("expected '%d', got %q", tt.iterations, val)
			}
		})
	}
}
