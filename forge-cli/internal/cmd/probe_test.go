package cmd

import (
	"bytes"
	"testing"
)

func TestProbeCmd_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "probe" {
			found = true
			break
		}
	}
	if !found {
		t.Error("probe command should be registered as top-level command")
	}
}

func TestProbeCmd_Use(t *testing.T) {
	if probeCmd.Use != "probe [path]" {
		t.Errorf("probeCmd.Use = %q, want %q", probeCmd.Use, "probe [path]")
	}
}

func TestProbeCmd_Short(t *testing.T) {
	if probeCmd.Short == "" {
		t.Error("probeCmd.Short should not be empty")
	}
}

func TestProbeCmd_MaxArgs(t *testing.T) {
	// Probe accepts at most 1 argument (optional path)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"probe", "too", "many"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for too many arguments")
	}
}
