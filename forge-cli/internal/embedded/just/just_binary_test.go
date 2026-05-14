// Package just_test tests the embedded just binary infrastructure.
package just

import (
	"testing"
)

func TestBinary(t *testing.T) {
	t.Run("returns non-nil byte slice", func(t *testing.T) {
		bin := Binary()
		if bin == nil {
			t.Fatal("Binary() must not return nil")
		}
	})

	t.Run("returns non-empty byte slice", func(t *testing.T) {
		bin := Binary()
		if len(bin) == 0 {
			t.Fatal("Binary() must not return an empty slice")
		}
	})

	t.Run("binary starts with expected ELF/Mach-O/PE magic", func(t *testing.T) {
		bin := Binary()
		if len(bin) < 4 {
			t.Fatalf("binary too short (%d bytes), expected at least 4 bytes for magic number", len(bin))
		}
		// Check for known binary magic numbers:
		// ELF:   0x7f 0x45 0x4c 0x46
		// Mach-O: 0xcf 0xfa 0xed 0xfe (little-endian) or 0xce 0xfa 0xed 0xfe
		// PE:     0x4d 0x5a ("MZ")
		isELF := bin[0] == 0x7f && bin[1] == 0x45 && bin[2] == 0x4c && bin[3] == 0x46
		isMachO := (bin[0] == 0xcf && bin[1] == 0xfa) || (bin[0] == 0xce && bin[1] == 0xfa) || (bin[0] == 0xca && bin[1] == 0xfe)
		isPE := bin[0] == 0x4d && bin[1] == 0x5a
		if !isELF && !isMachO && !isPE {
			t.Errorf("binary does not start with a recognized executable magic number (first 4 bytes: %x)", bin[:4])
		}
	})

	t.Run("binary size is reasonable", func(t *testing.T) {
		bin := Binary()
		size := len(bin)
		const minSize = 1_000_000  // 1 MB - minimum reasonable just binary
		const maxSize = 20_000_000 // 20 MB - upper bound for just binary
		if size < minSize {
			t.Errorf("binary seems too small (%d bytes), expected at least %d bytes", size, minSize)
		}
		if size > maxSize {
			t.Errorf("binary seems too large (%d bytes), expected at most %d bytes", size, maxSize)
		}
	})
}
