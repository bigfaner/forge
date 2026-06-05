package cmd

import (
	"bytes"
	"io"
	"os"

	"forge-cli/pkg/forgelog"
)

// captureOutput captures stdout and stderr during a function execution.
// Also closes forgelog file handles to prevent Windows "file in use" errors
// during test temp directory cleanup.
func captureOutput(f func() error) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, err := os.Pipe()
	if err != nil {
		return "", err
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = wOut
	os.Stderr = wErr

	outCh := make(chan string)
	errCh := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rOut)
		outCh <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rErr)
		errCh <- buf.String()
	}()

	runErr := f()

	// Close forgelog file handles before restoring stderr,
	// so temp directory cleanup can delete log files on Windows.
	forgelog.Close()

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	stdout := <-outCh
	stderr := <-errCh

	return stdout + stderr, runErr
}
