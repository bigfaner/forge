package main

import (
	"testing"
)

func TestMain_CallsRunFn(t *testing.T) {
	called := false
	orig := runFn
	runFn = func() { called = true }
	defer func() { runFn = orig }()

	main()
	if !called {
		t.Error("main() did not call runFn")
	}
}
