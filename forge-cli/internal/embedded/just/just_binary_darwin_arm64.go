//go:build darwin && arm64

package just

import _ "embed"

//go:embed binaries/just-darwin-arm64
var justBinary []byte

// Binary returns the embedded just binary for darwin/arm64.
func Binary() []byte {
	return justBinary
}
