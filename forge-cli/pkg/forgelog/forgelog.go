// Package forgelog provides a unified diagnostic output gateway for forge-cli.
// It implements the Backend abstraction pattern where console and file are
// independent backends. Console output preserves the original format byte-for-byte;
// file output adds timestamp+level prefix with level filtering.
//
// Hard rules:
//   - forgelog functions do NOT append \n -- the formatted message is output exactly as-is
//   - FileBackend write errors are silently ignored (never propagate to caller)
//   - No external dependencies -- stdlib only
package forgelog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

// Log level constants ordered by severity.
const (
	// DEBUG is the lowest severity level for verbose diagnostic output.
	DEBUG LogLevel = iota
	// INFO is the default level for informational messages.
	INFO
	// WARN is for warning messages about potential issues.
	WARN
	// ERROR is for error messages about failures.
	ERROR
)

// String returns the string representation of a LogLevel.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// parseLogLevel parses a level string, defaulting to INFO for unrecognized values.
func parseLogLevel(s string) LogLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// Backend is a log output target.
type Backend interface {
	// Write writes a log message at the given level and timestamp.
	// The message is output exactly as-is (no \n appended).
	Write(level LogLevel, timestamp time.Time, msg string)
	// Close releases backend resources.
	Close() error
}

// ConsoleBackend writes to stderr with original format.
// Output: just the message as-is (preserves current behavior exactly).
// No level filtering -- always outputs all messages sent to it.
type ConsoleBackend struct{}

// Write outputs the message to stderr exactly as-is.
func (c *ConsoleBackend) Write(_ LogLevel, _ time.Time, msg string) {
	fmt.Fprint(os.Stderr, msg)
}

// Close is a no-op for ConsoleBackend.
func (c *ConsoleBackend) Close() error {
	return nil
}

// FileBackend writes to a log file with structured format.
// Output: 2006-01-02T15:04:05.000 [LEVEL] message
// Level filtering suppresses messages below the configured level.
type FileBackend struct {
	mu       sync.Mutex
	file     *os.File
	minLevel LogLevel
}

// NewFileBackend creates a FileBackend that writes to the given file path.
// Messages below minLevel are suppressed.
func NewFileBackend(path string, minLevel LogLevel) (*FileBackend, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}
	return &FileBackend{
		file:     f,
		minLevel: minLevel,
	}, nil
}

// Write writes a log message to the file with timestamp+level prefix.
// Messages below the configured level are silently dropped.
// Write errors are silently ignored (hard rule).
func (fb *FileBackend) Write(level LogLevel, timestamp time.Time, msg string) {
	if level < fb.minLevel {
		return
	}

	line := fmt.Sprintf("%s [%s] %s",
		timestamp.Format("2006-01-02T15:04:05.000"),
		level.String(),
		msg,
	)

	fb.mu.Lock()
	defer fb.mu.Unlock()
	if fb.file != nil {
		// Silently ignore write errors (hard rule)
		_, _ = fb.file.WriteString(line)
	}
}

// Close releases the file handle.
func (fb *FileBackend) Close() error {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	if fb.file != nil {
		err := fb.file.Close()
		fb.file = nil
		return err
	}
	return nil
}

// LogsConfig holds the logging configuration.
// When integrated with forgeconfig, this struct will be added to Config as:
//
//	Logs *LogsConfig `yaml:"logs,omitempty"`
type LogsConfig struct {
	Enabled       bool   `yaml:"enabled"`       // default: true; set false to disable file logging
	Level         string `yaml:"level"`         // default: "info"
	RetentionDays int    `yaml:"retentionDays"` // default: 7
}

// Global state managed by Init/Close.
var (
	globalMu sync.Mutex
	backends []Backend
)

// Init initializes the logging layer.
//   - Creates logsDir on demand via os.MkdirAll(logsDir, 0700).
//   - Falls back to console-only if directory creation fails.
//   - Checks FORGE_NO_LOG=1 -- if set, skips FileBackend.
//   - config may be nil (defaults: level=info, retentionDays=7).
func Init(config *LogsConfig, logsDir string) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	// Always register ConsoleBackend
	backends = []Backend{&ConsoleBackend{}}

	// Check env var disable
	if os.Getenv("FORGE_NO_LOG") == "1" {
		return nil
	}

	// Check config disable
	if config != nil && !config.Enabled {
		return nil
	}

	// Resolve level
	levelStr := "info"
	retentionDays := 7
	if config != nil {
		if config.Level != "" {
			levelStr = config.Level
		}
		if config.RetentionDays >= 1 {
			retentionDays = config.RetentionDays
		}
	}
	minLevel := parseLogLevel(levelStr)

	// Try to create log directory; fall back to console-only on failure
	if err := os.MkdirAll(logsDir, 0o700); err != nil {
		// Fallback: console-only mode
		return nil
	}

	// Create log file: <ISO-8601-datetime>-<pid>.log
	filename := fmt.Sprintf("%s-%d.log",
		time.Now().Format("2006-01-02T15-04-05"),
		os.Getpid(),
	)
	logPath := filepath.Join(logsDir, filename)

	fb, err := NewFileBackend(logPath, minLevel)
	if err != nil {
		// Fallback: console-only mode
		return nil
	}
	backends = append(backends, fb)

	// Auto-cleanup old log files
	cleanupOldLogs(logsDir, retentionDays)

	return nil
}

// cleanupOldLogs deletes log files older than retentionDays.
// Errors are silently ignored.
func cleanupOldLogs(logsDir string, retentionDays int) {
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(logsDir, entry.Name()))
		}
	}
}

// Close releases all backend resources.
func Close() {
	globalMu.Lock()
	defer globalMu.Unlock()

	for _, b := range backends {
		_ = b.Close()
	}
	backends = nil
}

// dispatch sends a formatted message to all backends.
func dispatch(level LogLevel, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ts := time.Now()

	globalMu.Lock()
	bs := backends
	globalMu.Unlock()

	for _, b := range bs {
		b.Write(level, ts, msg)
	}
}

// Debug logs a message at DEBUG level.
func Debug(format string, args ...any) {
	dispatch(DEBUG, format, args...)
}

// Info logs a message at INFO level.
func Info(format string, args ...any) {
	dispatch(INFO, format, args...)
}

// Warn logs a message at WARN level.
func Warn(format string, args ...any) {
	dispatch(WARN, format, args...)
}

// Error logs a message at ERROR level.
func Error(format string, args ...any) {
	dispatch(ERROR, format, args...)
}
