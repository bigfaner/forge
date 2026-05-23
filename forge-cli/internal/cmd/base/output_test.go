package base

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncateSlug_CJK(t *testing.T) {
	t.Run("bug: byte-level truncation splits CJK characters producing invalid UTF-8", func(t *testing.T) {
		// Chinese title where byte-level cut at position 47 lands mid-character
		title := "autogen.go Quick 模式：替换 gen-and-run 为 gen-journeys + gen-contracts"
		maxLen := 50

		result := TruncateSlug(title, maxLen)

		if !utf8.ValidString(result) {
			t.Errorf("TruncateSlug produced invalid UTF-8: %q (len=%d)", result, len(result))
		}
	})

	t.Run("does not truncate short CJK titles", func(t *testing.T) {
		title := "新增集成测试"
		result := TruncateSlug(title, 50)
		if result != title {
			t.Errorf("short CJK title should not be truncated, got %q", result)
		}
	})

	t.Run("truncated result is valid UTF-8 for various CJK strings", func(t *testing.T) {
		titles := []string{
			"这是一个非常长的中文标题用来测试截断功能是否正常工作应该要被截断的",
			"新增 test.gen-journeys 和 test.gen-contracts 的集成测试",
			"创建 embed 模板 test-gen-journeys.md 和 test-gen-contracts.md",
		}
		for _, title := range titles {
			result := TruncateSlug(title, 30)
			if !utf8.ValidString(result) {
				t.Errorf("invalid UTF-8 for %q: got %q", title, result)
			}
		}
	})

	t.Run("truncated result display width does not exceed maxLen", func(t *testing.T) {
		title := "这是一个非常长的中文标题用来测试截断功能是否正常工作应该要被截断的"
		maxLen := 20
		result := TruncateSlug(title, maxLen)
		dw := testDisplayWidth(result)
		if dw > maxLen {
			t.Errorf("display width %d exceeds maxLen %d for result %q", dw, maxLen, result)
		}
	})

	t.Run("ASCII truncation unchanged", func(t *testing.T) {
		result := TruncateSlug("short", 10)
		if result != "short" {
			t.Errorf("ASCII short string should pass through, got %q", result)
		}

		longASCII := "this-is-a-very-long-slug-name-that-should-be-cut"
		result = TruncateSlug(longASCII, 15)
		if !strings.HasSuffix(result, "...") {
			t.Errorf("long ASCII should be truncated with '...', got %q", result)
		}
	})
}

func TestPadRight_CJK(t *testing.T) {
	t.Run("bug: pads by byte length instead of display width for CJK", func(t *testing.T) {
		s := "新增"
		// "新增" has byte length 6 but display width 4
		targetWidth := 10

		padded := PadRight(s, targetWidth)
		dw := testDisplayWidth(padded)

		if dw != targetWidth {
			t.Errorf("PadRight(%q, %d): display width = %d, want %d (result = %q, len=%d)",
				s, targetWidth, dw, targetWidth, padded, len(padded))
		}
	})

	t.Run("ASCII padding unchanged", func(t *testing.T) {
		padded := PadRight("hello", 10)
		if len(padded) != 10 {
			t.Errorf("PadRight('hello', 10) len = %d, want 10", len(padded))
		}
	})

	t.Run("no padding when already at target width", func(t *testing.T) {
		padded := PadRight("exact", 5)
		if padded != "exact" {
			t.Errorf("PadRight('exact', 5) = %q, want 'exact'", padded)
		}
	})
}

// testDisplayWidth counts the display width of a string (CJK runes = 2 columns).
func testDisplayWidth(s string) int {
	w := 0
	for _, r := range s {
		if testIsWide(r) {
			w += 2
		} else {
			w += 1
		}
	}
	return w
}

// isWide returns true for runes that occupy 2 columns in terminal output.
// Simplified: covers CJK Unified Ideographs, Hiragana, Katakana, Hangul,
// fullwidth forms, and other common double-width ranges.
func testIsWide(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) ||
		r == 0x2329 || r == 0x232A ||
		(r >= 0x2E80 && r <= 0xA4CF && r != 0x303F) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE10 && r <= 0xFE19) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF01 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6) ||
		(r >= 0x1F300 && r <= 0x1F64F) ||
		(r >= 0x20000 && r <= 0x2FFFD) ||
		(r >= 0x30000 && r <= 0x3FFFD)
}
