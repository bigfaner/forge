//go:build !linux && !darwin && !windows

// Package just provides access to the embedded just binary for the current platform.
package just

// Binary returns the embedded just binary for the current OS and architecture.
// The correct platform-specific implementation is selected via Go build tags.
// If no platform binary is embedded, it returns nil.
func Binary() []byte {
	// No supported platform matched; return nil.
	return nil
}
