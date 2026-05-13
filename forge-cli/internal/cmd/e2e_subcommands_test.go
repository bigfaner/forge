package cmd

import (
	"testing"
)

func TestE2EGroupStructure(t *testing.T) {
	subcommands := e2eCmd.Commands()
	subNames := make(map[string]bool)
	for _, cmd := range subcommands {
		subNames[cmd.Name()] = true
	}

	expected := []string{"validate-specs", "run", "setup", "verify", "compile", "discover"}
	for _, name := range expected {
		if !subNames[name] {
			t.Errorf("e2e group missing subcommand: %s (have: %v)", name, subNames)
		}
	}
}

func TestE2ERunCmd_HasFeatureFlag(t *testing.T) {
	f := e2eRunCmd.Flags().Lookup("feature")
	if f == nil {
		t.Fatal("e2e run command should have --feature flag")
	}
	if f.DefValue != "" {
		t.Errorf("expected default empty, got %q", f.DefValue)
	}
}

func TestE2ESetupCmd_HasForceFlag(t *testing.T) {
	f := e2eSetupCmd.Flags().Lookup("force")
	if f == nil {
		t.Fatal("e2e setup command should have --force flag")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %q", f.DefValue)
	}
}

func TestE2EVerifyCmd_FeatureFlagRequired(t *testing.T) {
	f := e2eVerifyCmd.Flags().Lookup("feature")
	if f == nil {
		t.Fatal("e2e verify command should have --feature flag")
	}
}

func TestE2ECompileCmd_NoFlags(t *testing.T) {
	if e2eCompileCmd.Flags().HasFlags() {
		t.Error("e2e compile command should have no flags")
	}
}

func TestE2EDiscoverCmd_NoFlags(t *testing.T) {
	if e2eDiscoverCmd.Flags().HasFlags() {
		t.Error("e2e discover command should have no flags")
	}
}

func TestE2ERunCmd_UseField(t *testing.T) {
	if e2eRunCmd.Use != "run" {
		t.Errorf("expected Use 'run', got %q", e2eRunCmd.Use)
	}
}

func TestE2ESetupCmd_UseField(t *testing.T) {
	if e2eSetupCmd.Use != "setup" {
		t.Errorf("expected Use 'setup', got %q", e2eSetupCmd.Use)
	}
}

func TestE2EVerifyCmd_UseField(t *testing.T) {
	if e2eVerifyCmd.Use != "verify" {
		t.Errorf("expected Use 'verify', got %q", e2eVerifyCmd.Use)
	}
}

func TestE2ECompileCmd_UseField(t *testing.T) {
	if e2eCompileCmd.Use != "compile" {
		t.Errorf("expected Use 'compile', got %q", e2eCompileCmd.Use)
	}
}

func TestE2EDiscoverCmd_UseField(t *testing.T) {
	if e2eDiscoverCmd.Use != "discover" {
		t.Errorf("expected Use 'discover', got %q", e2eDiscoverCmd.Use)
	}
}
