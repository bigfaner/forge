//go:build !windows

package index

import (
	"os"
	"syscall"
)

// tryLock attempts a non-blocking exclusive flock on Unix.
func tryLock(f *os.File) (bool, error) {
	err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		if err == syscall.EWOULDBLOCK {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// unlockPlatform releases the flock on Unix.
func unlockPlatform(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
