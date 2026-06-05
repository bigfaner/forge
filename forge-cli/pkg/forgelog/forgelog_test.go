package forgelog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"forge-cli/pkg/forgeconfig"
)

// ptrBool is a test helper that returns a pointer to the given bool value.
func ptrBool(v bool) *bool { return &v }

// --- LogLevel tests ---

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel(%d).String() = %q, want %q", tt.level, got, tt.expected)
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", DEBUG},
		{"info", INFO},
		{"warn", WARN},
		{"error", ERROR},
		{"INFO", INFO},
		{"Warn", WARN},
		{"bogus", INFO}, // default
		{"", INFO},      // default
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseLogLevel(tt.input)
			if got != tt.expected {
				t.Errorf("parseLogLevel(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

// --- ConsoleBackend tests ---

func TestConsoleBackendOutput(t *testing.T) {
	// AC-1: ConsoleBackend outputs message exactly as-is to stderr
	backend := &ConsoleBackend{}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStderr := os.Stderr
	os.Stderr = w

	backend.Write(WARN, time.Now(), "WARNING: task x not found\n")

	_ = w.Close()
	os.Stderr = oldStderr

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	got := string(buf[:n])
	expected := "WARNING: task x not found\n"
	if got != expected {
		t.Errorf("ConsoleBackend.Write() = %q, want %q", got, expected)
	}
}

func TestConsoleBackendNoNewline(t *testing.T) {
	// Hard rule: forgelog does NOT append \n
	backend := &ConsoleBackend{}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStderr := os.Stderr
	os.Stderr = w

	backend.Write(INFO, time.Now(), "no newline here")

	_ = w.Close()
	os.Stderr = oldStderr

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	got := string(buf[:n])
	if strings.HasSuffix(got, "\n") {
		t.Errorf("ConsoleBackend appended \\n, but hard rule says it must not")
	}
	if got != "no newline here" {
		t.Errorf("ConsoleBackend.Write() = %q, want %q", got, "no newline here")
	}
}

func TestConsoleBackendCloseIsNoop(t *testing.T) {
	backend := &ConsoleBackend{}
	if err := backend.Close(); err != nil {
		t.Errorf("ConsoleBackend.Close() returned error: %v", err)
	}
}

// --- FileBackend tests ---

func TestFileBackendWriteFormat(t *testing.T) {
	// AC-2: FileBackend writes format: 2006-01-02T15:04:05.000 [LEVEL] message
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	backend := NewFileBackend(logFile, INFO)

	ts := time.Date(2026, 6, 4, 17, 30, 0, 123000000, time.Local)
	backend.Write(WARN, ts, "WARNING: task x not found\n")

	_ = backend.Close()

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	got := string(data)
	// Expected: 2026-06-04T17:30:00.123 [WARN] WARNING: task x not found\n
	if !strings.Contains(got, "[WARN]") {
		t.Errorf("FileBackend output missing [WARN]: %q", got)
	}
	if !strings.Contains(got, "WARNING: task x not found\n") {
		t.Errorf("FileBackend output missing message: %q", got)
	}
	// Verify timestamp format
	expectedPrefix := "2026-06-04T17:30:00.123 [WARN] WARNING: task x not found\n"
	if got != expectedPrefix {
		t.Errorf("FileBackend.Write() = %q, want %q", got, expectedPrefix)
	}
}

func TestFileBackendLevelFiltering(t *testing.T) {
	// AC-2: Level filtering suppresses below configured level in file only
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	backend := NewFileBackend(logFile, WARN)

	ts := time.Now()
	backend.Write(DEBUG, ts, "debug msg\n")
	backend.Write(INFO, ts, "info msg\n")
	backend.Write(WARN, ts, "warn msg\n")
	backend.Write(ERROR, ts, "error msg\n")

	_ = backend.Close()

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	got := string(data)
	if strings.Contains(got, "debug msg") {
		t.Error("FileBackend should suppress DEBUG when level is WARN")
	}
	if strings.Contains(got, "info msg") {
		t.Error("FileBackend should suppress INFO when level is WARN")
	}
	if !strings.Contains(got, "warn msg") {
		t.Error("FileBackend should include WARN when level is WARN")
	}
	if !strings.Contains(got, "error msg") {
		t.Error("FileBackend should include ERROR when level is WARN")
	}
}

func TestFileBackendWriteErrorSilenced(_ *testing.T) {
	// Hard rule: FileBackend write errors are silently ignored
	// Use an invalid path that cannot be opened
	backend := NewFileBackend("/nonexistent/dir/test.log", INFO)

	// This should not panic — ensureOpen fails silently
	ts := time.Now()
	backend.Write(INFO, ts, "should not panic\n")
}

func TestFileBackendConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	backend := NewFileBackend(logFile, INFO)

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			backend.Write(INFO, time.Now(), fmt.Sprintf("line %d\n", n))
		}(i)
	}
	wg.Wait()
	_ = backend.Close()

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Count(string(data), "\n")
	if lines != 100 {
		t.Errorf("Expected 100 lines, got %d", lines)
	}
}

// --- Init tests ---

func TestInitCreatesLogDir(t *testing.T) {
	// AC-3: Init() creates .forge/logs/ via os.MkdirAll(dir, 0700)
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	// Directory should not exist yet
	if _, err := os.Stat(logsDir); !os.IsNotExist(err) {
		t.Fatalf("logsDir should not exist before Init()")
	}

	err := Init(nil, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	info, err := os.Stat(logsDir)
	if err != nil {
		t.Fatalf("logsDir should exist after Init(): %v", err)
	}
	if !info.IsDir() {
		t.Error("logsDir should be a directory")
	}
	if runtime.GOOS != "windows" {
		// AC-5: directory mode 0700 (Unix only)
		if info.Mode().Perm() != 0o700 {
			t.Errorf("logsDir mode = %o, want 0700", info.Mode().Perm())
		}
	}
}

func TestInitFallbackOnDirCreationFailure(t *testing.T) {
	// AC-3: Falls back to console-only if creation fails
	// Use a path that cannot be created (e.g., under a non-existent read-only parent)
	dir := t.TempDir()
	// Create a file where a directory should be - blocks MkdirAll
	blockPath := filepath.Join(dir, "blocked")
	if err := os.WriteFile(blockPath, []byte("block"), 0o644); err != nil {
		t.Fatal(err)
	}
	logsDir := filepath.Join(blockPath, "sub", "logs")

	// Init should not fail - it falls back to console-only
	err := Init(nil, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() should not fail even when dir creation fails: %v", err)
	}

	// ConsoleBackend should still work
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStderr := os.Stderr
	os.Stderr = w

	Info("test message\n")

	_ = w.Close()
	os.Stderr = oldStderr

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	got := string(buf[:n])
	if !strings.Contains(got, "test message\n") {
		t.Errorf("Console fallback should output message, got %q", got)
	}
}

func TestInitWithEnvDisable(t *testing.T) {
	// FORGE_NO_LOG=1 skips FileBackend
	t.Setenv("FORGE_NO_LOG", "1")
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	err := Init(nil, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// No .forge/logs directory should be created
	if _, err := os.Stat(logsDir); !os.IsNotExist(err) {
		t.Error("logsDir should not be created when FORGE_NO_LOG=1")
	}
}

func TestInitFilePermissions(t *testing.T) {
	// AC-5: Log file created with mode 0600
	if runtime.GOOS == "windows" {
		t.Skip("file permissions are advisory on Windows")
	}

	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	err := Init(nil, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Write something to create the file
	Warn("test\n")
	Close()

	// Find the log file
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no log files created")
	}

	info, err := entries[0].Info()
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("log file mode = %o, want 0600", info.Mode().Perm())
	}
}

func TestInitLogFileNaming(t *testing.T) {
	// File naming: date-based (2006-01-02.log), one file per day
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	err := Init(nil, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	Info("test\n")
	Close()

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log file, got %d", len(entries))
	}

	name := entries[0].Name()
	expected := time.Now().Format("2006-01-02") + ".log"
	if name != expected {
		t.Errorf("log file name = %q, want %q", name, expected)
	}
}

func TestConcurrentInitAppendsSameDay(t *testing.T) {
	// Same-day inits append to the same date-based file
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	// First init
	err := Init(nil, logsDir)
	if err != nil {
		t.Fatal(err)
	}
	Info("first\n")
	Close()

	// Second init on same day appends to same file
	err = Init(nil, logsDir)
	if err != nil {
		t.Fatal(err)
	}
	Info("second\n")
	Close()

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatal(err)
	}
	// Same day → 1 file, not 2
	if len(entries) != 1 {
		t.Errorf("expected 1 log file (same day), got %d", len(entries))
	}

	// Verify both messages are in the file
	content, err := os.ReadFile(filepath.Join(logsDir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "first\n") || !strings.Contains(string(content), "second\n") {
		t.Errorf("expected both messages in log file, got:\n%s", string(content))
	}
}

// --- Printf-style API tests ---

func TestPrintfStyleAPI(t *testing.T) {
	// Verify Debug/Info/Warn/Error format correctly
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStderr := os.Stderr
	os.Stderr = w

	// Initialize console-only (no file backend)
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")
	t.Setenv("FORGE_NO_LOG", "1")
	err = Init(nil, logsDir)
	if err != nil {
		t.Fatal(err)
	}

	Warn("WARNING: task %s not found\n", "x")

	_ = w.Close()
	os.Stderr = oldStderr
	Close()

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	got := string(buf[:n])
	expected := "WARNING: task x not found\n"
	if got != expected {
		t.Errorf("Warn() = %q, want %q", got, expected)
	}
}

func TestAPIByteIdenticalToFmtFprintf(t *testing.T) {
	// AC-1: byte-identical to fmt.Fprintf(os.Stderr, ...)
	// Compare: forgelog.Warn("WARNING: task %s not found\n", "x")
	// vs: fmt.Fprintf(os.Stderr, "WARNING: task %s not found\n", "x")

	msg := "WARNING: task %s not found\n"
	arg := "x"

	// Capture forgelog output
	r1, w1, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w1

	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")
	t.Setenv("FORGE_NO_LOG", "1")
	_ = Init(nil, logsDir)
	Warn(msg, arg)
	_ = w1.Close()
	os.Stderr = oldStderr
	Close()

	buf1 := make([]byte, 1024)
	n1, _ := r1.Read(buf1)
	forgelogOutput := string(buf1[:n1])

	// Capture fmt.Fprintf output
	r2, w2, _ := os.Pipe()
	os.Stderr = w2
	fmt.Fprintf(os.Stderr, msg, arg)
	_ = w2.Close()
	os.Stderr = oldStderr

	buf2 := make([]byte, 1024)
	n2, _ := r2.Read(buf2)
	fmtOutput := string(buf2[:n2])

	if forgelogOutput != fmtOutput {
		t.Errorf("forgelog output %q != fmt.Fprintf output %q", forgelogOutput, fmtOutput)
	}
}

func TestFileBackendLazyCreation(t *testing.T) {
	// NewFileBackend does NOT create the file; it's created on first Write
	dir := t.TempDir()
	logFile := filepath.Join(dir, "lazy.log")

	backend := NewFileBackend(logFile, INFO)
	defer func() { _ = backend.Close() }()

	// File should NOT exist yet
	if _, err := os.Stat(logFile); !os.IsNotExist(err) {
		t.Fatal("lazy log file should not exist after NewFileBackend")
	}

	// Write creates the file
	backend.Write(INFO, time.Now(), "hello\n")

	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("log file should exist after Write: %v", err)
	}
}

func TestInitDoesNotCreateEmptyFile(t *testing.T) {
	// Init should not create an empty log file if no messages are dispatched
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	err := Init(nil, logsDir)
	if err != nil {
		t.Fatal(err)
	}

	// Close without writing any messages
	Close()

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatal(err)
	}
	// logsDir may exist (MkdirAll), but no .log files should be present
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".log") {
			t.Errorf("no log file should exist without writes, found %s", e.Name())
		}
	}
}

func TestCloseWithoutInit(_ *testing.T) {
	// Close should be safe to call even without Init
	// Reset global state first
	Close()
}

func TestMultipleClose(t *testing.T) {
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")
	_ = Init(nil, logsDir)
	Close()
	Close() // Should not panic
}

// --- Auto-cleanup tests (AC-1) ---

func TestCleanupOldLogs(t *testing.T) {
	// Create old and new log files; verify only old ones are deleted
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create an old log file (10 days ago)
	oldFile := filepath.Join(logsDir, "old.log")
	if err := os.WriteFile(oldFile, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Set modification time to 10 days ago
	oldTime := time.Now().AddDate(0, 0, -10)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	// Create a recent log file (1 day ago)
	recentFile := filepath.Join(logsDir, "recent.log")
	if err := os.WriteFile(recentFile, []byte("recent"), 0o600); err != nil {
		t.Fatal(err)
	}
	recentTime := time.Now().AddDate(0, 0, -1)
	if err := os.Chtimes(recentFile, recentTime, recentTime); err != nil {
		t.Fatal(err)
	}

	// Init with 7-day retention should delete old but keep recent
	err := Init(nil, logsDir)
	defer Close()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("old log file should have been deleted")
	}
	if _, err := os.Stat(recentFile); os.IsNotExist(err) {
		t.Error("recent log file should NOT have been deleted")
	}
}

func TestCleanupNeverDeletesActiveLog(t *testing.T) {
	// SC-4: active log file is never deleted by its own cleanup pass
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a file that would be "old" if deleted
	// But the active log (created by Init) has current timestamp, so it's safe
	err := Init(nil, logsDir)
	defer Close()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Write to create the log file
	Info("test\n")

	// Verify log file exists
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one log file")
	}

	// The newly created log file should still exist
	found := false
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".log") {
			found = true
			break
		}
	}
	if !found {
		t.Error("active log file was deleted by cleanup")
	}
}

// --- Config disable tests (AC-3) ---

func TestInitWithConfigDisabled(t *testing.T) {
	// logs.enabled: false in config skips FileBackend
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	err := Init(&forgeconfig.LogsConfig{Enabled: ptrBool(false)}, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// No .forge/logs directory should be created
	if _, err := os.Stat(logsDir); !os.IsNotExist(err) {
		t.Error("logsDir should not be created when logs.enabled=false")
	}
}

func TestInitWithEnvDisableOverridesConfig(t *testing.T) {
	// FORGE_NO_LOG=1 takes precedence over config
	t.Setenv("FORGE_NO_LOG", "1")
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	// Config says enabled=true, but env says disable
	err := Init(&forgeconfig.LogsConfig{Enabled: ptrBool(true)}, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// No .forge/logs directory should be created (env takes precedence)
	if _, err := os.Stat(logsDir); !os.IsNotExist(err) {
		t.Error("logsDir should not be created when FORGE_NO_LOG=1 overrides config")
	}
}

func TestInitWithConfigLevelAndRetention(t *testing.T) {
	// Config with custom level and retention
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	cfg := &forgeconfig.LogsConfig{
		Enabled:       ptrBool(true),
		Level:         "warn",
		RetentionDays: 3,
	}
	err := Init(cfg, logsDir)
	defer Close()

	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Write info and warn messages
	Info("info msg\n")
	Warn("warn msg\n")
	Close()

	// Read log file and check filtering
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one log file")
	}

	data, err := os.ReadFile(filepath.Join(logsDir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if strings.Contains(content, "info msg") {
		t.Error("INFO message should be suppressed when level is warn")
	}
	if !strings.Contains(content, "warn msg") {
		t.Error("WARN message should be present when level is warn")
	}
}

func TestInitWithInvalidRetentionDays(t *testing.T) {
	// AC-2: retentionDays < 1 falls back to 7
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	// Create an old file (5 days old, within default 7-day retention but outside 0-day)
	oldFile := filepath.Join(logsDir, "old.log")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(oldFile, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().AddDate(0, 0, -5)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	// Config with retentionDays=0 should fall back to 7, so the 5-day-old file survives
	cfg := &forgeconfig.LogsConfig{
		Enabled:       ptrBool(true),
		RetentionDays: 0,
	}
	err := Init(cfg, logsDir)
	defer Close()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 5-day-old file should still exist (within 7-day default retention)
	if _, err := os.Stat(oldFile); os.IsNotExist(err) {
		t.Error("5-day-old file should survive when retentionDays=0 falls back to 7")
	}
}
