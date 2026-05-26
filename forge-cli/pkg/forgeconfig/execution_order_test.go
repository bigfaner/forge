package forgeconfig

import (
	"strings"
	"testing"
)

func TestNormalizeSurfaceKey(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// AC3: "ADMIN PANEL" normalizes to "admin-panel"
		{name: "spaces to hyphens and lowercase", input: "ADMIN PANEL", want: "admin-panel"},
		{name: "mixed case to lowercase", input: "MyBackend", want: "mybackend"},
		{name: "forward slash to hyphen", input: "frontend/api", want: "frontend-api"},
		{name: "multiple special chars normalized", input: "My Service / API", want: "my-service---api"},
		{name: "already valid key", input: "backend", want: "backend"},
		{name: "hyphens preserved", input: "admin-panel", want: "admin-panel"},
		{name: "numbers allowed after first char", input: "api-v2", want: "api-v2"},
		{name: "single char lowercase passes", input: "a", want: "a"},

		// AC3: "123bad" fails after normalization
		{name: "starts with digit fails", input: "123bad", want: "123bad", wantErr: true},
		{name: "starts with hyphen fails", input: "-bad", want: "-bad", wantErr: true},
		{name: "starts with digit after normalization", input: "123 API", want: "123-api", wantErr: true},
		{name: "empty string fails", input: "", want: "", wantErr: true},
		{name: "only spaces normalizes to empty fails", input: "   ", want: "---", wantErr: true},
		{name: "uppercase first char normalizes valid", input: "Frontend", want: "frontend"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeSurfaceKey(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NormalizeSurfaceKey(%q) expected error, got nil (result=%q)", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Errorf("NormalizeSurfaceKey(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeSurfaceKey(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateSurfaceKeys(t *testing.T) {
	t.Run("valid keys pass", func(t *testing.T) {
		surfaces := map[string]string{
			"backend":  "api",
			"frontend": "web",
		}
		err := ValidateSurfaceKeys(surfaces)
		if err != nil {
			t.Errorf("expected no error for valid keys, got: %v", err)
		}
	})

	t.Run("scalar dot key is valid", func(t *testing.T) {
		surfaces := map[string]string{".": "api"}
		err := ValidateSurfaceKeys(surfaces)
		if err != nil {
			t.Errorf("expected no error for scalar dot key, got: %v", err)
		}
	})

	t.Run("invalid key returns error", func(t *testing.T) {
		surfaces := map[string]string{
			"123bad": "web",
		}
		err := ValidateSurfaceKeys(surfaces)
		if err == nil {
			t.Fatal("expected error for invalid key")
		}
		if !strings.Contains(err.Error(), "invalid surface-key") {
			t.Errorf("error should mention 'invalid surface-key', got: %v", err)
		}
	})
}

func TestValidateExecutionOrder(t *testing.T) {
	t.Run("nil surfaces returns nil", func(t *testing.T) {
		err := ValidateExecutionOrder(nil, nil)
		if err != nil {
			t.Errorf("expected nil for nil surfaces, got: %v", err)
		}
	})

	t.Run("scalar surfaces with nil order returns nil", func(t *testing.T) {
		surfaces := map[string]string{".": "api"}
		err := ValidateExecutionOrder(surfaces, nil)
		if err != nil {
			t.Errorf("expected nil for scalar surfaces, got: %v", err)
		}
	})

	t.Run("single surface with nil order returns nil", func(t *testing.T) {
		surfaces := map[string]string{"backend": "api"}
		err := ValidateExecutionOrder(surfaces, nil)
		if err != nil {
			t.Errorf("expected nil for single surface, got: %v", err)
		}
	})

	// AC1: execution-order references non-existent surface-key
	t.Run("order references non-existent key errors", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
			"backend":  "api",
		}
		err := ValidateExecutionOrder(surfaces, []string{"backend", "nonexistent"})
		if err == nil {
			t.Fatal("expected error for non-existent key in execution-order")
		}
		if !strings.Contains(err.Error(), "nonexistent") {
			t.Errorf("error should mention the invalid key, got: %v", err)
		}
	})

	t.Run("valid order with all keys present passes", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
			"backend":  "api",
		}
		err := ValidateExecutionOrder(surfaces, []string{"backend", "frontend"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	// AC2: same type conflict detected
	t.Run("same type conflict without execution-order errors", func(t *testing.T) {
		surfaces := map[string]string{
			"auth-service":    "api",
			"payment-service": "api",
			"admin":           "web",
		}
		err := ValidateExecutionOrder(surfaces, nil)
		if err == nil {
			t.Fatal("expected error for same-type conflict without execution-order")
		}
		if !strings.Contains(err.Error(), "execution-order") {
			t.Errorf("error should mention 'execution-order', got: %v", err)
		}
	})

	// AC2: same type conflict resolved with execution-order
	t.Run("same type conflict resolved with execution-order passes", func(t *testing.T) {
		surfaces := map[string]string{
			"auth-service":    "api",
			"payment-service": "api",
			"admin":           "web",
		}
		err := ValidateExecutionOrder(surfaces, []string{"auth-service", "payment-service", "admin"})
		if err != nil {
			t.Errorf("expected no error when execution-order resolves conflict, got: %v", err)
		}
	})
}

func TestResolveExecutionOrder(t *testing.T) {
	// AC4: default priority api > web > cli > tui > mobile
	t.Run("default priority api web cli mobile", func(t *testing.T) {
		surfaces := map[string]string{
			"mobile": "mobile",
			"cli":    "cli",
			"web":    "web",
			"api":    "api",
		}
		order, err := ResolveExecutionOrder(surfaces, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"api", "web", "cli", "mobile"}
		if len(order) != len(expected) {
			t.Fatalf("expected %d, got %d: %v", len(expected), len(order), order)
		}
		for i, key := range expected {
			if order[i] != key {
				t.Errorf("position %d: got %q, want %q", i, order[i], key)
			}
		}
	})

	t.Run("explicit order overrides default", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
			"backend":  "api",
		}
		order, err := ResolveExecutionOrder(surfaces, []string{"frontend", "backend"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"frontend", "backend"}
		if len(order) != len(expected) {
			t.Fatalf("expected %d, got %d: %v", len(expected), len(order), order)
		}
		for i, key := range expected {
			if order[i] != key {
				t.Errorf("position %d: got %q, want %q", i, order[i], key)
			}
		}
	})

	t.Run("nil surfaces returns nil", func(t *testing.T) {
		order, err := ResolveExecutionOrder(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if order != nil {
			t.Errorf("expected nil, got %v", order)
		}
	})

	t.Run("scalar surfaces returns nil", func(t *testing.T) {
		order, err := ResolveExecutionOrder(map[string]string{".": "api"}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if order != nil {
			t.Errorf("expected nil for scalar surfaces, got %v", order)
		}
	})

	t.Run("single surface returns nil", func(t *testing.T) {
		order, err := ResolveExecutionOrder(map[string]string{"backend": "api"}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if order != nil {
			t.Errorf("expected nil for single surface, got %v", order)
		}
	})

	t.Run("uncovered combination uses YAML map order", func(t *testing.T) {
		// tui + cli has no default priority rule; should use insertion order
		surfaces := map[string]string{
			"terminal": "tui",
			"tool":     "cli",
		}
		order, err := ResolveExecutionOrder(surfaces, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order) != 2 {
			t.Fatalf("expected 2 keys, got %d: %v", len(order), order)
		}
		// Both keys should be present; exact order is YAML map order
		seen := map[string]bool{}
		for _, k := range order {
			seen[k] = true
		}
		if !seen["terminal"] || !seen["tool"] {
			t.Errorf("expected both keys present, got %v", order)
		}
	})

	t.Run("mixed types with default priority", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
			"backend":  "api",
		}
		order, err := ResolveExecutionOrder(surfaces, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order) != 2 {
			t.Fatalf("expected 2, got %d: %v", len(order), order)
		}
		// api should come before web by default
		if order[0] != "backend" {
			t.Errorf("expected backend (api) first, got %q", order[0])
		}
		if order[1] != "frontend" {
			t.Errorf("expected frontend (web) second, got %q", order[1])
		}
	})
}

func TestReadConfig_ExecutionOrder(t *testing.T) {
	t.Run("execution-order field parsed", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  frontend: web
  backend: api
execution-order:
  - backend
  - frontend
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ExecutionOrder == nil {
			t.Fatal("expected ExecutionOrder non-nil")
		}
		if len(cfg.ExecutionOrder) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(cfg.ExecutionOrder))
		}
		if cfg.ExecutionOrder[0] != "backend" {
			t.Errorf("expected first entry 'backend', got %q", cfg.ExecutionOrder[0])
		}
		if cfg.ExecutionOrder[1] != "frontend" {
			t.Errorf("expected second entry 'frontend', got %q", cfg.ExecutionOrder[1])
		}
	})

	t.Run("execution-order absent is nil", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  frontend: web
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ExecutionOrder != nil {
			t.Errorf("expected ExecutionOrder nil when absent, got %v", cfg.ExecutionOrder)
		}
	})
}

func TestReadConfig_SurfaceKeyNormalization(t *testing.T) {
	// AC3: "ADMIN PANEL" normalizes to "admin-panel" at config load time
	t.Run("ADMIN PANEL normalizes to admin-panel", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  ADMIN PANEL: web
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := cfg.Surfaces["admin-panel"]; !ok {
			t.Errorf("expected key 'admin-panel' after normalization, got surfaces: %v", cfg.Surfaces)
		}
	})

	// AC3: "123bad" fails at config load time
	t.Run("123bad fails at config load time", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  123bad: web
`)
		_, err := ReadConfig(dir)
		if err == nil {
			t.Fatal("expected error for invalid surface-key '123bad'")
		}
		if !strings.Contains(err.Error(), "invalid surface-key") {
			t.Errorf("error should mention 'invalid surface-key', got: %v", err)
		}
	})

	t.Run("normalized key with slash becomes hyphen", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  frontend/api: web
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := cfg.Surfaces["frontend-api"]; !ok {
			t.Errorf("expected key 'frontend-api' after normalization, got surfaces: %v", cfg.Surfaces)
		}
	})

	t.Run("scalar surfaces bypass key normalization", func(t *testing.T) {
		dir := setupConfig(t, "surfaces: api\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Surfaces["."] != "api" {
			t.Errorf("expected scalar surfaces '.', got %v", cfg.Surfaces)
		}
	})
}

func TestReadConfig_ExecutionOrderValidation(t *testing.T) {
	// AC1: execution-order references non-existent surface-key
	t.Run("execution-order references non-existent key errors at load time", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  frontend: web
  backend: api
execution-order:
  - backend
  - nonexistent
`)
		_, err := ReadConfig(dir)
		if err == nil {
			t.Fatal("expected error for non-existent key in execution-order")
		}
		if !strings.Contains(err.Error(), "nonexistent") {
			t.Errorf("error should mention the invalid key, got: %v", err)
		}
	})

	// AC2: same type conflict without execution-order
	t.Run("same type conflict without execution-order errors at load time", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  auth-service: api
  payment-service: api
  admin: web
`)
		_, err := ReadConfig(dir)
		if err == nil {
			t.Fatal("expected error for same-type conflict")
		}
		if !strings.Contains(err.Error(), "execution-order") {
			t.Errorf("error should mention 'execution-order', got: %v", err)
		}
	})

	// AC2: same type conflict resolved with execution-order
	t.Run("same type conflict with execution-order passes", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  auth-service: api
  payment-service: api
  admin: web
execution-order:
  - auth-service
  - payment-service
  - admin
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("expected no error with execution-order, got: %v", err)
		}
		if len(cfg.Surfaces) != 3 {
			t.Errorf("expected 3 surfaces, got %d", len(cfg.Surfaces))
		}
	})

	// AC4: default priority verification via ResolveExecutionOrder after successful load
	t.Run("default priority api web cli mobile via resolve", func(t *testing.T) {
		dir := setupConfig(t, `surfaces:
  mobile: mobile
  cli: cli
  web: web
  api: api
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		order, err := ResolveExecutionOrder(cfg.Surfaces, cfg.ExecutionOrder)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"api", "web", "cli", "mobile"}
		for i, key := range expected {
			if order[i] != key {
				t.Errorf("position %d: got %q, want %q", i, order[i], key)
			}
		}
	})
}
