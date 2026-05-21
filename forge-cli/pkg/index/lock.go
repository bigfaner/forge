// Package index provides file locking for concurrent write safety on index.json.
package index

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrLockConflict indicates that the lock could not be acquired within the timeout.
var ErrLockConflict = errors.New("concurrent write conflict, retry")

// defaultLockTimeout is the maximum time to wait for lock acquisition.
const defaultLockTimeout = 5 * time.Second

// LockFile acquires an exclusive advisory lock on <feature-dir>/tasks/index.json.lock.
// Creates the lock file if it does not exist. The lock file persists for reuse.
// Blocks for up to 5 seconds waiting for the lock; returns ErrLockConflict on timeout.
// The returned *os.File holds the lock - callers must call UnlockFile when done.
func LockFile(indexPath string) (*os.File, error) {
	lockPath := indexPath + ".lock"

	// Ensure parent directory exists
	dir := filepath.Dir(lockPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create lock directory: %w", err)
	}

	// Open or create the lock file
	f, err := os.OpenFile(lockPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create lock file: %w", err)
	}

	// Try to acquire lock with timeout
	deadline := time.Now().Add(defaultLockTimeout)
	for {
		acquired, tryErr := tryLock(f)
		if tryErr != nil {
			_ = f.Close()
			return nil, fmt.Errorf("failed to acquire lock: %w", tryErr)
		}
		if acquired {
			return f, nil
		}

		if time.Now().After(deadline) {
			_ = f.Close()
			return nil, ErrLockConflict
		}

		// Brief backoff before retry
		time.Sleep(50 * time.Millisecond)
	}
}

// UnlockFile releases the advisory lock and closes the file descriptor.
// The lock file itself is NOT deleted - it persists for reuse.
func UnlockFile(f *os.File) error {
	if f == nil {
		return nil
	}
	err := unlockPlatform(f)
	if closeErr := f.Close(); closeErr != nil && err == nil {
		return closeErr
	}
	return err
}

// WithLock acquires the advisory lock for indexPath, calls fn, then releases the lock.
// The lock is always released (even if fn panics) via defer.
// Returns ErrLockConflict if the lock cannot be acquired within the 5-second timeout.
func WithLock(indexPath string, fn func() error) error {
	lock, err := LockFile(indexPath)
	if err != nil {
		return err
	}
	defer func() { _ = UnlockFile(lock) }()

	return fn()
}
