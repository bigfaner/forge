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
)

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
	backend, err := NewFileBackend(logFile, INFO)
	if err != nil {
		t.Fatal(err)
	}

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
	backend, err := NewFileBackend(logFile, WARN)
	if err != nil {
		t.Fatal(err)
	}

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

func TestFileBackendWriteErrorSilenced(t *testing.T) {
	// Hard rule: FileBackend write errors are silently ignored
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	backend, err := NewFileBackend(logFile, INFO)
	if err != nil {
		t.Fatal(err)
	}

	// Close the underlying file to cause write errors
	_ = backend.Close()

	// This should not panic or return error
	ts := time.Now()
	backend.Write(INFO, ts, "should not panic\n")
}

func TestFileBackendConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	backend, err := NewFileBackend(logFile, INFO)
	if err != nil {
		t.Fatal(err)
	}

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
	// File naming: ISO-8601 datetime + PID
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
	pid := os.Getpid()
	expectedSuffix := fmt.Sprintf("-%d.log", pid)
	if !strings.HasSuffix(name, expectedSuffix) {
		t.Errorf("log file name %q should end with %q", name, expectedSuffix)
	}
	// Should match pattern: 2006-01-02T15-04-05-<pid>.log
	if !strings.Contains(name, "T") {
		t.Errorf("log file name %q should contain ISO date separator T", name)
	}
}

func TestConcurrentInitProducesSeparateFiles(t *testing.T) {
	// AC-4: Two concurrent Init() calls produce separate log files with distinct PIDs
	// Since we can't fork in tests, simulate by calling Init twice sequentially
	// which creates two log files (different timestamps)
	dir := t.TempDir()
	logsDir := filepath.Join(dir, ".forge", "logs")

	// First init
	err := Init(nil, logsDir)
	if err != nil {
		t.Fatal(err)
	}
	Info("first\n")
	Close()

	// Small delay to ensure different timestamp
	time.Sleep(1100 * time.Millisecond)

	// Second init
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
	if len(entries) < 2 {
		t.Errorf("expected at least 2 log files, got %d", len(entries))
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
