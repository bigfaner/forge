// Package main is the entry point for the forge CLI.
package main

// runFn is the function called by main. Overridden in tests.
var runFn = Run

func main() {
	runFn()
}
