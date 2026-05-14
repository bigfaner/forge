//go:build darwin && amd64

package just

import _ "embed"

//go:embed binaries/just-darwin-amd64
var justBinary []byte

// Binary returns the embedded just binary for darwin/amd64.
func Binary() []byte {
	return justBinary
}
