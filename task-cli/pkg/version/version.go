// Package version provides version information for the CLI.
package version

// Version is the CLI version, injected at build time via ldflags.
// Example: go build -ldflags "-X task-cli/pkg/version.Version=v1.0.0"
var Version = "dev"

// Name is the CLI name.
var Name = "task"

// GetVersion returns the CLI version.
func GetVersion() string {
	return Version
}

// GetName returns the CLI name.
func GetName() string {
	return Name
}
