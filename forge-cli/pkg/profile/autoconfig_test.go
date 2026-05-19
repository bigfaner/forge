package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadAutoConfig_MissingConfig(t *testing.T) {
	projectRoot := t.TempDir()

	auto, err := ReadAutoConfig(projectRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Missing config should return defaults: e2eTest quick=false/full=true, consolidateSpecs quick=true/full=true, cleanCode false
	if auto.E2eTest.Quick || !auto.E2eTest.Full {
		t.Errorf("E2eTest defaults = %+v, want {Quick:false Full:true}", auto.E2eTest)
	}
	if !auto.ConsolidateSpecs.Quick || !auto.ConsolidateSpecs.Full {
		t.Errorf("ConsolidateSpecs defaults = %+v, want {Quick:true Full:true}", auto.ConsolidateSpecs)
	}
	if auto.CleanCode.Quick || auto.CleanCode.Full {
		t.Errorf("CleanCode defaults = %+v, want {Quick:false Full:false}", auto.CleanCode)
	}
	if auto.Validation.Quick || auto.Validation.Full {
		t.Errorf("Validation defaults = %+v, want {Quick:false Full:false}", auto.Validation)
	}
	if auto.GitPush {
		t.Errorf("GitPush default = %v, want false", auto.GitPush)
	}
}

func TestReadAutoConfig_WithAutoBlock(t *testing.T) {
	projectRoot := t.TempDir()
	forgeDir := filepath.Join(projectRoot, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := `languages:
  - go

auto:
  e2eTest:
    quick: false
    full: true
  consolidateSpecs:
    quick: false
    full: false
  cleanCode:
    quick: true
    full: true
  gitPush: true
`
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	auto, err := ReadAutoConfig(projectRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if auto.E2eTest.Quick {
		t.Error("E2eTest.Quick should be false")
	}
	if !auto.E2eTest.Full {
		t.Error("E2eTest.Full should be true")
	}
	if auto.ConsolidateSpecs.Quick {
		t.Error("ConsolidateSpecs.Quick should be false")
	}
	if auto.ConsolidateSpecs.Full {
		t.Error("ConsolidateSpecs.Full should be false")
	}
	if !auto.CleanCode.Quick {
		t.Error("CleanCode.Quick should be true")
	}
	if !auto.CleanCode.Full {
		t.Error("CleanCode.Full should be true")
	}
	if auto.Validation.Quick || auto.Validation.Full {
		t.Errorf("Validation should default to false/false, got %+v", auto.Validation)
	}
	if !auto.GitPush {
		t.Error("GitPush should be true")
	}
}

func TestReadAutoConfig_PartialAutoBlock(t *testing.T) {
	projectRoot := t.TempDir()
	forgeDir := filepath.Join(projectRoot, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Only set e2eTest, others should get defaults
	configContent := `languages:
  - go

auto:
  e2eTest:
    quick: false
`
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	auto, err := ReadAutoConfig(projectRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if auto.E2eTest.Quick {
		t.Error("E2eTest.Quick should be false (explicitly set)")
	}
	if !auto.E2eTest.Full {
		t.Error("E2eTest.Full should be true (default)")
	}
	// consolidateSpecs not in YAML → defaults: quick=true, full=true
	if !auto.ConsolidateSpecs.Quick || !auto.ConsolidateSpecs.Full {
		t.Errorf("ConsolidateSpecs should default to true/true, got %+v", auto.ConsolidateSpecs)
	}
	if auto.CleanCode.Quick || auto.CleanCode.Full {
		t.Errorf("CleanCode should default to false/false, got %+v", auto.CleanCode)
	}
	if auto.Validation.Quick || auto.Validation.Full {
		t.Errorf("Validation should default to false/false, got %+v", auto.Validation)
	}
}

func TestReadAutoConfig_NoAutoBlock(t *testing.T) {
	projectRoot := t.TempDir()
	forgeDir := filepath.Join(projectRoot, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := `languages:
  - go

`
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	auto, err := ReadAutoConfig(projectRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All defaults should apply: e2eTest quick=false/full=true, consolidateSpecs quick=true/full=true
	if auto.E2eTest.Quick || !auto.E2eTest.Full {
		t.Errorf("E2eTest defaults = %+v, want {Quick:false Full:true}", auto.E2eTest)
	}
	if !auto.ConsolidateSpecs.Quick || !auto.ConsolidateSpecs.Full {
		t.Errorf("ConsolidateSpecs defaults = %+v, want {Quick:true Full:true}", auto.ConsolidateSpecs)
	}
	if auto.CleanCode.Quick || auto.CleanCode.Full {
		t.Errorf("CleanCode defaults = %+v, want {Quick:false Full:false}", auto.CleanCode)
	}
	if auto.Validation.Quick || auto.Validation.Full {
		t.Errorf("Validation defaults = %+v, want {Quick:false Full:false}", auto.Validation)
	}
}
