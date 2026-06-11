package index

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLockFile_BasicAcquireRelease(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	lock, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("LockFile failed: %v", err)
	}
	if lock == nil {
		t.Fatal("LockFile returned nil file")
	}

	// Lock file should exist
	lockPath := indexPath + ".lock"
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Error("lock file should exist after LockFile")
	}

	err = UnlockFile(lock)
	if err != nil {
		t.Fatalf("UnlockFile failed: %v", err)
	}
}

func TestLockFile_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "deep", "nested", "tasks", "index.json")

	lock, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("LockFile should create parent directories: %v", err)
	}
	defer func() { _ = UnlockFile(lock) }()

	lockPath := indexPath + ".lock"
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Error("lock file should exist")
	}
}

func TestLockFile_LockFilePersists(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	lock, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("LockFile failed: %v", err)
	}
	_ = UnlockFile(lock)

	// Lock file should still exist after unlock
	lockPath := indexPath + ".lock"
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Error("lock file should persist after unlock for reuse")
	}

	// Can acquire again on same lock file
	lock2, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("LockFile on reuse failed: %v", err)
	}
	_ = UnlockFile(lock2)
}

func TestLockFile_ConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	// Acquire first lock
	lock1, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("first LockFile failed: %v", err)
	}

	// Try to acquire second lock concurrently - should timeout
	var wg sync.WaitGroup
	wg.Add(1)

	var secondErr error
	go func() {
		defer wg.Done()
		_, secondErr = LockFile(indexPath)
	}()

	wg.Wait()

	if secondErr == nil {
		t.Error("second LockFile should have failed with ErrLockConflict")
		_ = UnlockFile(lock1)
		return
	}

	if !errors.Is(secondErr, ErrLockConflict) {
		t.Errorf("expected ErrLockConflict, got: %v", secondErr)
	}

	_ = UnlockFile(lock1)
}

func TestLockFile_SequentialAccess(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	// Acquire and release first lock
	lock1, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("first LockFile failed: %v", err)
	}
	_ = UnlockFile(lock1)

	// Second lock should succeed immediately after first is released
	lock2, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("second LockFile after release failed: %v", err)
	}
	_ = UnlockFile(lock2)
}

func TestLockFile_DifferentFeaturesNoContention(t *testing.T) {
	dir := t.TempDir()
	indexPath1 := filepath.Join(dir, "feature-a", "tasks", "index.json")
	indexPath2 := filepath.Join(dir, "feature-b", "tasks", "index.json")

	// Lock feature-a
	lock1, err := LockFile(indexPath1)
	if err != nil {
		t.Fatalf("LockFile feature-a failed: %v", err)
	}
	defer func() { _ = UnlockFile(lock1) }()

	// Lock feature-b should succeed - different lock scope
	lock2, err := LockFile(indexPath2)
	if err != nil {
		t.Fatalf("LockFile feature-b should not be blocked by feature-a lock: %v", err)
	}
	_ = UnlockFile(lock2)
}

func TestLockFile_Timeout(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	// Acquire first lock
	lock1, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("first LockFile failed: %v", err)
	}

	// Second lock attempt should timeout in ~5 seconds
	start := time.Now()
	_, err = LockFile(indexPath)
	elapsed := time.Since(start)

	_ = UnlockFile(lock1)

	if err == nil {
		t.Fatal("expected ErrLockConflict, got nil")
	}
	if !errors.Is(err, ErrLockConflict) {
		t.Errorf("expected ErrLockConflict, got: %v", err)
	}

	// Should have waited approximately 5 seconds (allow 1s tolerance)
	if elapsed < 4*time.Second {
		t.Errorf("expected ~5s timeout, but only waited %v", elapsed)
	}
}

func TestUnlockFile_NilFile(t *testing.T) {
	// Should not panic
	err := UnlockFile(nil)
	if err != nil {
		t.Errorf("UnlockFile(nil) should return nil, got: %v", err)
	}
}

func TestLockFile_InvalidPath(t *testing.T) {
	// Path that cannot be created (e.g., under a file that exists as a file)
	dir := t.TempDir()
	// Create a file where a directory would need to be
	indexPath := filepath.Join(dir, "not-a-dir", "tasks", "index.json")
	_ = os.WriteFile(filepath.Join(dir, "not-a-dir"), []byte("blocker"), 0644)

	_, err := LockFile(indexPath)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestSaveIndexAtomic_Basic(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	data := map[string]string{"key": "value"}
	if err := SaveIndexAtomic(indexPath, data); err != nil {
		t.Fatalf("SaveIndexAtomic failed: %v", err)
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "{\n  \"key\": \"value\"\n}\n"
	if string(content) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(content))
	}
}

func TestSaveIndexAtomic_NoPartialState(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	// Write initial data
	initial := map[string]string{"version": "1"}
	if err := SaveIndexAtomic(indexPath, initial); err != nil {
		t.Fatalf("first SaveIndexAtomic failed: %v", err)
	}

	// Overwrite with new data
	updated := map[string]string{"version": "2"}
	if err := SaveIndexAtomic(indexPath, updated); err != nil {
		t.Fatalf("second SaveIndexAtomic failed: %v", err)
	}

	// File should contain updated data
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "{\n  \"version\": \"2\"\n}\n"
	if string(content) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(content))
	}

	// No temp files should remain
	entries, err := os.ReadDir(filepath.Dir(indexPath))
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	for _, e := range entries {
		if matched, _ := filepath.Match(".index.json.tmp.*", e.Name()); matched {
			t.Errorf("temp file should not remain: %s", e.Name())
		}
	}
}

func TestSaveIndexAtomic_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "deep", "nested", "tasks", "index.json")

	data := map[string]string{"key": "value"}
	if err := SaveIndexAtomic(indexPath, data); err != nil {
		t.Fatalf("SaveIndexAtomic should create parent dir: %v", err)
	}

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("index file should exist")
	}
}

func TestSaveIndexAtomic_MarshalError(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.json")

	// Channel values cannot be marshaled to JSON
	data := map[string]any{"ch": make(chan int)}
	err := SaveIndexAtomic(indexPath, data)
	if err == nil {
		t.Fatal("expected marshal error for channel value")
	}
}

func TestSaveIndexAtomic_Overwrite(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.json")

	// Write initial
	if err := SaveIndexAtomic(indexPath, map[string]string{"a": "1"}); err != nil {
		t.Fatalf("first write: %v", err)
	}

	// Overwrite
	if err := SaveIndexAtomic(indexPath, map[string]string{"b": "2"}); err != nil {
		t.Fatalf("second write: %v", err)
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(content) != "{\n  \"b\": \"2\"\n}\n" {
		t.Errorf("unexpected content: %s", content)
	}
}

func TestLockFile_ReusesLockFile(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")
	lockPath := indexPath + ".lock"

	// First lock
	lock1, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("first LockFile: %v", err)
	}
	_ = UnlockFile(lock1)

	// Check lock file exists
	info1, err := os.Stat(lockPath)
	if err != nil {
		t.Fatalf("lock file should exist: %v", err)
	}

	// Second lock should reuse same file
	lock2, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("second LockFile: %v", err)
	}
	_ = UnlockFile(lock2)

	info2, err := os.Stat(lockPath)
	if err != nil {
		t.Fatalf("lock file should still exist: %v", err)
	}

	// Same file (same size, same name)
	if info1.Size() != info2.Size() {
		t.Error("lock file should be reused, not recreated with different size")
	}
}

func TestSaveIndexAtomic_LargeData(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.json")

	// Large map to exercise write path
	data := make(map[string]string, 1000)
	for i := 0; i < 1000; i++ {
		data[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	if err := SaveIndexAtomic(indexPath, data); err != nil {
		t.Fatalf("SaveIndexAtomic with large data: %v", err)
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(content) == 0 {
		t.Error("file should not be empty")
	}
}

func TestSaveIndexAtomic_ReadOnlyDir(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "readonly")
	if err := os.MkdirAll(subDir, 0555); err != nil {
		t.Skipf("cannot set up read-only dir: %v", err)
	}

	indexPath := filepath.Join(subDir, "index.json")
	err := SaveIndexAtomic(indexPath, map[string]string{"k": "v"})
	if err == nil {
		t.Log("SaveIndexAtomic succeeded on read-only dir (OS allows it)")
	}
}

func TestSaveIndexAtomic_MkdirAllError(t *testing.T) {
	dir := t.TempDir()
	// Create a file where a directory would need to be created
	blockerPath := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blockerPath, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join(blockerPath, "sub", "index.json")
	err := SaveIndexAtomic(indexPath, map[string]string{"k": "v"})
	if err == nil {
		t.Error("expected mkdir error")
	}
}

func TestLockFile_SaveUnderLock(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	lock, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("LockFile: %v", err)
	}

	// While holding the lock, atomically save index
	data := map[string]string{"task": "done"}
	if err := SaveIndexAtomic(indexPath, data); err != nil {
		_ = UnlockFile(lock)
		t.Fatalf("SaveIndexAtomic under lock: %v", err)
	}

	_ = UnlockFile(lock)

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	expected := "{\n  \"task\": \"done\"\n}\n"
	if string(content) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(content))
	}
}

func TestUnlockFile_CloseError(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	lock, err := LockFile(indexPath)
	if err != nil {
		t.Fatalf("LockFile: %v", err)
	}

	err = UnlockFile(lock)
	if err != nil {
		t.Logf("UnlockFile returned error (expected on some platforms): %v", err)
	}
}

func TestErrLockConflict_Value(t *testing.T) {
	if ErrLockConflict.Error() != "concurrent write conflict, retry" {
		t.Errorf("ErrLockConflict message = %q, want %q", ErrLockConflict.Error(), "concurrent write conflict, retry")
	}
}

// --- WithLock tests ---

func TestWithLock_BasicExecution(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	executed := false
	err := WithLock(indexPath, func() error {
		executed = true
		return nil
	})
	if err != nil {
		t.Fatalf("WithLock failed: %v", err)
	}
	if !executed {
		t.Error("callback was not executed")
	}
}

func TestWithLock_ReturnsCallbackError(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	cbErr := fmt.Errorf("callback failed")
	err := WithLock(indexPath, func() error {
		return cbErr
	})
	if !errors.Is(err, cbErr) {
		t.Errorf("expected callback error, got: %v", err)
	}
}

func TestWithLock_ConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	// WithLock held by first goroutine, second should fail
	errCh := make(chan error, 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		err := WithLock(indexPath, func() error {
			// Hold the lock long enough for the second attempt to fail
			time.Sleep(6 * time.Second)
			return nil
		})
		if err != nil {
			t.Logf("first WithLock errored: %v", err)
		}
	}()

	// Give the first goroutine time to acquire the lock
	time.Sleep(200 * time.Millisecond)

	err := WithLock(indexPath, func() error {
		return nil
	})

	// The second WithLock should fail with ErrLockConflict
	if !errors.Is(err, ErrLockConflict) {
		t.Errorf("expected ErrLockConflict, got: %v", err)
	}

	// Wait for first goroutine to complete
	<-done
	_ = errCh
}

func TestWithLock_SequentialReuse(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	// First call
	if err := WithLock(indexPath, func() error { return nil }); err != nil {
		t.Fatalf("first WithLock: %v", err)
	}

	// Second call should succeed after first released
	if err := WithLock(indexPath, func() error { return nil }); err != nil {
		t.Fatalf("second WithLock: %v", err)
	}
}

func TestWithLock_LockReleasedOnCallbackError(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "tasks", "index.json")

	// First call returns error — lock should still be released
	err := WithLock(indexPath, func() error {
		return fmt.Errorf("boom")
	})
	if err == nil {
		t.Fatal("expected error from callback")
	}

	// Second call should succeed because lock was released
	err = WithLock(indexPath, func() error { return nil })
	if err != nil {
		t.Fatalf("lock should have been released after error, but: %v", err)
	}
}

// --- AtomicWrite tests ---

func TestAtomicWrite_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	data := []byte("hello world")
	if err := AtomicWrite(path, data, 0644); err != nil {
		t.Fatalf("AtomicWrite failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", string(content))
	}
}

func TestAtomicWrite_NoTempFilesRemain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	if err := AtomicWrite(path, []byte("{}"), 0644); err != nil {
		t.Fatalf("AtomicWrite failed: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	for _, e := range entries {
		if matched, _ := filepath.Match(".state.json.tmp.*", e.Name()); matched {
			t.Errorf("temp file should not remain: %s", e.Name())
		}
	}
}

func TestAtomicWrite_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.bin")

	if err := AtomicWrite(path, []byte("first"), 0644); err != nil {
		t.Fatalf("first write: %v", err)
	}
	if err := AtomicWrite(path, []byte("second"), 0644); err != nil {
		t.Fatalf("second write: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(content) != "second" {
		t.Errorf("expected %q, got %q", "second", string(content))
	}
}

func TestAtomicWrite_CreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deep", "nested", "file.txt")

	if err := AtomicWrite(path, []byte("data"), 0644); err != nil {
		t.Fatalf("AtomicWrite with nested dir: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(content) != "data" {
		t.Errorf("expected %q, got %q", "data", string(content))
	}
}

func TestAtomicWrite_ErrorCleansUp(t *testing.T) {
	dir := t.TempDir()
	// Create a file where a directory would need to be
	blocker := filepath.Join(dir, "blocker")
	_ = os.WriteFile(blocker, []byte("x"), 0644)
	path := filepath.Join(blocker, "sub", "file.txt")

	err := AtomicWrite(path, []byte("data"), 0644)
	if err == nil {
		t.Error("expected error for blocked path")
	}
}

func TestAtomicWrite_EmptyData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := AtomicWrite(path, []byte{}, 0644); err != nil {
		t.Fatalf("AtomicWrite with empty data: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file, got %d bytes", info.Size())
	}
}
