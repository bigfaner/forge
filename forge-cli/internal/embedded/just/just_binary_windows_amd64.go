//go:build windows && amd64

package just

import _ "embed"

//go:embed binaries/just-windows-amd64.exe
var justBinary []byte

// Binary returns the embedded just binary for windows/amd64.
func Binary() []byte {
	return justBinary
}
