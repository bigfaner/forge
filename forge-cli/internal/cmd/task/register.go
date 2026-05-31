// Package task contains all forge task subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
package task

// Register adds all task subcommands to Cmd (the parent task command).
func Register() {
	Cmd.AddCommand(
		claimCmd,
		submitCmd,
		statusCmd,
		queryCmd,
		checkDepsCmd,
		validateCmd,
		addCmd,
		indexCmd,
		migrateCmd,
		listTypesCmd,
		listCmd,
		reopenCmd,
		transitionCmd,
	)
}
