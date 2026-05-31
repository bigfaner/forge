package types

// Version is the CLI version, injected at build time via ldflags.
// Example: go build -ldflags "-X forge-cli/pkg/types.Version=v1.0.0"
var Version = "dev"

// Name is the CLI name.
var Name = "forge"

// GetVersion returns the CLI version.
func GetVersion() string {
	return Version
}

// GetName returns the CLI name.
func GetName() string {
	return Name
}
