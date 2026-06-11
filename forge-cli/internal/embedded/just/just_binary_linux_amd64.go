//go:build linux && amd64

package just

import _ "embed"

//go:embed binaries/just-linux-amd64
var justBinary []byte

// Binary returns the embedded just binary for linux/amd64.
func Binary() []byte {
	return justBinary
}
