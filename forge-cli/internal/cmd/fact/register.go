package fact

// Register adds all fact subcommands to Cmd (the parent fact command).
func Register() {
	Cmd.AddCommand(
		listCmd,
		getCmd,
		summaryCmd,
	)
}
