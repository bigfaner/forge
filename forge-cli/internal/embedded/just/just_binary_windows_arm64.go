//go:build windows && arm64

package just

import _ "embed"

//go:embed binaries/just-windows-arm64.exe
var justBinary []byte

// Binary returns the embedded just binary for windows/arm64.
func Binary() []byte {
	return justBinary
}
