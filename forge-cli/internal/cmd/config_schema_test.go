package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// schemaPath resolves the forge-config schema relative to the test binary.
func schemaPath(t *testing.T) string {
	t.Helper()
	// Test runs from forge-cli/internal/cmd/; schema is in plugins/forge/references/shared/
	p := filepath.Join("..", "..", "..", "plugins", "forge", "references", "shared", "forge-config.schema.json")
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("resolve schema path: %v", err)
	}
	return abs
}

func TestConfigSchemaAutoBlock(t *testing.T) {
	data, err := os.ReadFile(schemaPath(t))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("parse schema JSON: %v", err)
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("schema missing 'properties' object")
	}

	autoProp, ok := props["auto"]
	if !ok {
		t.Fatal("schema properties missing 'auto' key")
	}
	autoObj, ok := autoProp.(map[string]any)
	if !ok {
		t.Fatal("'auto' property is not an object")
	}

	if autoObj["type"] != "object" {
		t.Errorf("auto.type = %v, want 'object'", autoObj["type"])
	}

	if autoObj["additionalProperties"] != false {
		t.Errorf("auto.additionalProperties = %v, want false", autoObj["additionalProperties"])
	}

	autoProps, ok := autoObj["properties"].(map[string]any)
	if !ok {
		t.Fatal("auto missing 'properties' object")
	}

	// Verify all expected sub-objects exist with quick/full booleans
	modeFields := []string{"e2eTest", "consolidateSpecs", "cleanCode"}
	for _, field := range modeFields {
		fieldObj, ok := autoProps[field].(map[string]any)
		if !ok {
			t.Fatalf("auto.properties missing '%s' or not an object", field)
		}
		if fieldObj["type"] != "object" {
			t.Errorf("auto.%s.type = %v, want 'object'", field, fieldObj["type"])
		}
		if fieldObj["additionalProperties"] != false {
			t.Errorf("auto.%s.additionalProperties = %v, want false", field, fieldObj["additionalProperties"])
		}

		fieldProps, ok := fieldObj["properties"].(map[string]any)
		if !ok {
			t.Fatalf("auto.%s missing 'properties' object", field)
		}

		for _, mode := range []string{"quick", "full"} {
			modeProp, ok := fieldProps[mode].(map[string]any)
			if !ok {
				t.Fatalf("auto.%s.properties missing '%s' or not an object", field, mode)
			}
			if modeProp["type"] != "boolean" {
				t.Errorf("auto.%s.%s.type = %v, want 'boolean'", field, mode, modeProp["type"])
			}
		}
	}

	// Verify gitPush is a boolean
	gitPush, ok := autoProps["gitPush"].(map[string]any)
	if !ok {
		t.Fatal("auto.properties missing 'gitPush' or not an object")
	}
	if gitPush["type"] != "boolean" {
		t.Errorf("auto.gitPush.type = %v, want 'boolean'", gitPush["type"])
	}
}

func TestConfigSchemaAutoDefaults(t *testing.T) {
	data, err := os.ReadFile(schemaPath(t))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("parse schema JSON: %v", err)
	}

	autoRaw, ok := schema["properties"].(map[string]any)["auto"]
	if !ok {
		t.Fatal("schema properties missing 'auto' key — cannot verify defaults")
	}
	autoObj, ok := autoRaw.(map[string]any)
	if !ok {
		t.Fatal("'auto' property is not an object — cannot verify defaults")
	}
	autoProps, ok := autoObj["properties"].(map[string]any)
	if !ok {
		t.Fatal("auto missing 'properties' object — cannot verify defaults")
	}

	// Verify defaults per Hard Rules:
	// e2eTest.{quick,true; full,true}
	// consolidateSpecs.{quick,true; full,true}
	// cleanCode.{quick,false; full,false}
	// gitPush: false
	expectedDefaults := map[string]map[string]bool{
		"e2eTest":          {"quick": true, "full": true},
		"consolidateSpecs": {"quick": true, "full": true},
		"cleanCode":        {"quick": false, "full": false},
	}

	for field, modes := range expectedDefaults {
		fieldProps := autoProps[field].(map[string]any)["properties"].(map[string]any)
		for mode, expected := range modes {
			modeProp := fieldProps[mode].(map[string]any)
			defaultVal, ok := modeProp["default"]
			if !ok {
				t.Errorf("auto.%s.%s missing 'default'", field, mode)
				continue
			}
			if defaultVal != expected {
				t.Errorf("auto.%s.%s.default = %v, want %v", field, mode, defaultVal, expected)
			}
		}
	}

	gitPushDefault := autoProps["gitPush"].(map[string]any)["default"]
	if gitPushDefault != false {
		t.Errorf("auto.gitPush.default = %v, want false", gitPushDefault)
	}
}

func TestConfigSchemaBackwardCompatible(t *testing.T) {
	data, err := os.ReadFile(schemaPath(t))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("parse schema JSON: %v", err)
	}

	// 'auto' must not be in 'required' — existing configs without auto block must pass
	required, ok := schema["required"]
	if !ok {
		t.Fatal("schema missing 'required' field")
	}
	reqArr, ok := required.([]any)
	if !ok {
		t.Fatalf("schema 'required' is %T, want array", required)
	}
	for _, r := range reqArr {
		if r == "auto" {
			t.Error("'auto' must not be required — existing configs without auto block must continue to work")
		}
	}

	// Root additionalProperties: false must be preserved
	if schema["additionalProperties"] != false {
		t.Error("root schema additionalProperties must be false")
	}
}

func TestConfigExampleDocumentsAllAutoFields(t *testing.T) {
	// Verify the example YAML contains all 7 auto fields with comments
	examplePath := filepath.Join("..", "..", "..", "plugins", "forge", "references", "shared", "forge-config.example.yaml")
	abs, err := filepath.Abs(examplePath)
	if err != nil {
		t.Fatalf("resolve example path: %v", err)
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	content := string(data)

	// All 7 fields must appear in the example
	requiredFields := []string{
		"e2eTest:",
		"consolidateSpecs:",
		"cleanCode:",
		"gitPush:",
		"quick:",
		"full:",
	}
	for _, field := range requiredFields {
		if !strings.Contains(content, field) {
			t.Errorf("example YAML missing field %q", field)
		}
	}

	// Verify specific default values appear in the example
	expectedValues := map[string]bool{
		"quick: true":    false, // must appear at least once (e2eTest, consolidateSpecs)
		"quick: false":   false, // must appear at least once (cleanCode)
		"full: true":     false, // must appear at least once
		"full: false":    false, // must appear at least once
		"gitPush: false": true,  // must appear exactly
	}
	for val, required := range expectedValues {
		count := strings.Count(content, val)
		if required && count < 1 {
			t.Errorf("example YAML must contain %q", val)
		}
		if !required && count < 1 {
			t.Errorf("example YAML should contain %q at least once", val)
		}
	}
}
