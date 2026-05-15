//go:build linux && arm64

package just

import _ "embed"

//go:embed binaries/just-linux-arm64
var justBinary []byte

// Binary returns the embedded just binary for linux/arm64.
func Binary() []byte {
	return justBinary
}
