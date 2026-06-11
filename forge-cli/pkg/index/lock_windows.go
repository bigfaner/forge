package index

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	lockFileExclusiveLock                 = 0x00000002
	lockFileFailImmediately               = 0x00000001
	errorLockViolation      syscall.Errno = 33
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = kernel32.NewProc("LockFileEx")
	procUnlockFileEx = kernel32.NewProc("UnlockFileEx")
)

// tryLock attempts a non-blocking exclusive lock on Windows using LockFileEx.
func tryLock(f *os.File) (bool, error) {
	var overlapped syscall.Overlapped
	handle := syscall.Handle(f.Fd())

	// Try non-blocking exclusive lock
	ret, _, err := procLockFileEx.Call(
		uintptr(handle),
		lockFileExclusiveLock|lockFileFailImmediately,
		0,
		0xFFFFFFFF, // bytesLow: lock entire file
		0xFFFFFFFF, // bytesHigh
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if ret == 0 {
		// LockFileEx returns ERROR_LOCK_VIOLATION (33) if already locked
		if err == errorLockViolation {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// unlockPlatform releases the lock on Windows using UnlockFileEx.
func unlockPlatform(f *os.File) error {
	var overlapped syscall.Overlapped
	handle := syscall.Handle(f.Fd())

	ret, _, err := procUnlockFileEx.Call(
		uintptr(handle),
		0,
		0xFFFFFFFF, // bytesLow
		0xFFFFFFFF, // bytesHigh
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if ret == 0 {
		return err
	}
	return nil
}
