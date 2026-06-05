package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/forgelog"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// modeHighlight styles mode keywords for terminal display.
var modeHighlight = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorModeHighlight))

// hl returns a highlighted version of text using modeHighlight.
func hl(text string) string {
	return modeHighlight.Render(text)
}

// hlMode returns "Quick mode" or "Full mode" with the whole phrase highlighted.
func hlMode(mode string) string {
	return hl(mode + " mode")
}

// autoBehaviorPrompt defines one question in the auto-behavior config flow.
type autoBehaviorPrompt struct {
	title string                                       // question shown to user
	desc  string                                       // description shown below the question
	def   bool                                         // default value for the confirm prompt
	set   func(auto *forgeconfig.AutoConfig, val bool) // assigns the answer
}

// autoBehaviorPrompts is the ordered list of prompts for askAutoBehavior.
// Each prompt preserves the exact question text and defaults from the original
// per-block implementation to maintain behavioral equivalence.
func autoBehaviorPrompts(defaults forgeconfig.AutoConfig) []autoBehaviorPrompt {
	return []autoBehaviorPrompt{
		{
			title: fmt.Sprintf("%s: auto-run advanced tests?", hlMode("Quick")),
			desc:  fmt.Sprintf("Automatically run surface-level advanced tests during %s (lightweight verification after each task).", hl("quick mode")),
			def:   defaults.Test.Quick,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Test.Quick = v },
		},
		{
			title: fmt.Sprintf("%s: auto-run advanced tests?", hlMode("Full")),
			desc:  fmt.Sprintf("Automatically run surface-level advanced tests during %s (comprehensive coverage).", hl("full mode")),
			def:   defaults.Test.Full,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Test.Full = v },
		},
		{
			title: fmt.Sprintf("%s: auto-consolidate specs?", hlMode("Quick")),
			desc:  fmt.Sprintf("Automatically extract and consolidate specs from code after %s tasks.", hl("quick-mode")),
			def:   defaults.ConsolidateSpecs.Quick,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.ConsolidateSpecs.Quick = v },
		},
		{
			title: fmt.Sprintf("%s: auto-consolidate specs?", hlMode("Full")),
			desc:  fmt.Sprintf("Automatically extract and consolidate specs from code after %s tasks.", hl("full-mode")),
			def:   defaults.ConsolidateSpecs.Full,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.ConsolidateSpecs.Full = v },
		},
		{
			title: fmt.Sprintf("%s: auto code cleanup?", hlMode("Quick")),
			desc:  fmt.Sprintf("Automatically simplify and clean code during %s.", hl("quick mode")),
			def:   defaults.CleanCode.Quick,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.CleanCode.Quick = v },
		},
		{
			title: fmt.Sprintf("%s: auto code cleanup?", hlMode("Full")),
			desc:  fmt.Sprintf("Automatically simplify and clean code during %s.", hl("full mode")),
			def:   defaults.CleanCode.Full,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.CleanCode.Full = v },
		},
		{
			title: fmt.Sprintf("%s: auto validation?", hlMode("Quick")),
			desc:  fmt.Sprintf("Automatically run validation checks during %s (lightweight quality gates after each task).", hl("quick mode")),
			def:   defaults.Validation.Quick,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Validation.Quick = v },
		},
		{
			title: fmt.Sprintf("%s: auto validation?", hlMode("Full")),
			desc:  fmt.Sprintf("Automatically run validation checks during %s (comprehensive quality gates).", hl("full mode")),
			def:   defaults.Validation.Full,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Validation.Full = v },
		},
		{
			title: fmt.Sprintf("%s: auto-run tasks?", hlMode("Quick")),
			desc:  fmt.Sprintf("Automatically claim and execute tasks during %s.", hl("quick mode")),
			def:   defaults.RunTasks.Quick,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.RunTasks.Quick = v },
		},
		{
			title: fmt.Sprintf("%s: auto-run tasks?", hlMode("Full")),
			desc:  fmt.Sprintf("Automatically claim and execute tasks during %s.", hl("full mode")),
			def:   defaults.RunTasks.Full,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.RunTasks.Full = v },
		},
		{
			title: fmt.Sprintf("%s: auto knowledge save?", hlMode("Quick")),
			desc:  fmt.Sprintf("Automatically save knowledge after %s tasks.", hl("quick mode")),
			def:   defaults.KnowledgeSave.Quick,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.KnowledgeSave.Quick = v },
		},
		{
			title: fmt.Sprintf("%s: auto knowledge save?", hlMode("Full")),
			desc:  fmt.Sprintf("Automatically save knowledge after %s tasks.", hl("full mode")),
			def:   defaults.KnowledgeSave.Full,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.KnowledgeSave.Full = v },
		},
		{
			title: "Auto-evaluate proposals?",
			desc:  "Automatically run proposal evaluation after generation.",
			def:   defaults.Eval.Proposal,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Eval.Proposal = v },
		},
		{
			title: "Auto-evaluate PRD documents?",
			desc:  "Automatically run PRD evaluation after generation.",
			def:   defaults.Eval.Prd,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Eval.Prd = v },
		},
		{
			title: "Auto-evaluate UI designs?",
			desc:  "Automatically run UI design evaluation after generation.",
			def:   defaults.Eval.UiDesign,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Eval.UiDesign = v },
		},
		{
			title: "Auto-evaluate tech designs?",
			desc:  "Automatically run tech design evaluation after generation.",
			def:   defaults.Eval.TechDesign,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.Eval.TechDesign = v },
		},
		{
			title: "Auto git push after all tasks complete?",
			desc:  "Push to remote automatically when every task in a run finishes successfully.",
			def:   defaults.GitPush,
			set:   func(a *forgeconfig.AutoConfig, v bool) { a.GitPush = v },
		},
	}
}

// askAutoBehavior runs the auto-behavior config steps, one question per screen.
// Returns the config and whether the user cancelled.
func askAutoBehavior() (*forgeconfig.AutoConfig, bool) {
	defaults := forgeconfig.AutoConfigDefaults()
	auto := &forgeconfig.AutoConfig{}

	for _, p := range autoBehaviorPrompts(defaults) {
		val, ok := askConfirm(p.title, p.desc, p.def)
		if !ok {
			return nil, true
		}
		p.set(auto, val)
	}

	return auto, false
}

// askWorktreeConfig runs the optional worktree config steps.
// Returns nil if both source-branch and copy-files are empty (skippable).
// Returns (config, cancelled).
func askWorktreeConfig() (*forgeconfig.WorktreeConfig, bool) {
	var sourceBranch string
	err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Worktree source branch (leave empty to skip)").
			Description("Branch to use as the base for new worktrees (e.g. main, develop).").
			Value(&sourceBranch),
	)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, true
		}
		return nil, true
	}

	var copyFiles []string
	// Only ask about copy-files if user provided a source branch
	if sourceBranch != "" {
		err = huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Files to copy into worktrees").
				Description("Select files that should be copied from the source branch when creating a worktree.").
				Options(
					huh.NewOption(".env", ".env"),
					huh.NewOption(".env.local", ".env.local"),
					huh.NewOption(".env.development", ".env.development"),
				).
				Value(&copyFiles),
		)).Run()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil, true
			}
			return nil, true
		}
	}

	// Both empty means no worktree config block
	if sourceBranch == "" && len(copyFiles) == 0 {
		return nil, false
	}

	return &forgeconfig.WorktreeConfig{
		SourceBranch: sourceBranch,
		CopyFiles:    copyFiles,
	}, false
}

// askConfirm shows a single confirm prompt. Returns (value, ok).
// ok is false when the user pressed Ctrl+C.
func askConfirm(title, desc string, defaultVal bool) (bool, bool) {
	val := defaultVal
	err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(title).
			Description(desc).
			Affirmative("Yes").
			Negative("No").
			Value(&val),
	)).Run()
	if err != nil {
		return defaultVal, false
	}
	return val, true
}

func runConfigInitIfNeeded(projectRoot string) initAction {
	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)

	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeCharDevice == 0 {
		return initAction{status: "SKIPPED", target: ".forge/config.yaml", detail: "non-interactive terminal"}
	}

	// When config exists, ask if user wants to reconfigure
	if _, err := os.Stat(configFile); err == nil {
		reconfigure := false
		if err := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title("Config already exists. Reconfigure?").
				Description("Select Yes to overwrite .forge/config.yaml with new settings.").
				Affirmative("Yes").
				Negative("No").
				Value(&reconfigure),
		)).Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
			}
			return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
		}
		if !reconfigure {
			return initAction{status: "SKIPPED", target: ".forge/config.yaml", detail: "kept existing"}
		}
	}

	// Auto-behavior config
	auto, cancelled := askAutoBehavior()
	if cancelled {
		return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
	}

	// Worktree config (optional)
	worktree, cancelled := askWorktreeConfig()
	if cancelled {
		return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
	}

	cfg := forgeconfig.Config{
		Auto:     auto,
		Worktree: worktree,
	}

	if err := writeConfigFile(configFile, &cfg); err != nil {
		forgelog.Error("ERROR: failed to write config: %v\n", err)
		return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
	}

	return initAction{status: "CREATED", target: ".forge/config.yaml", detail: "interactive"}
}
