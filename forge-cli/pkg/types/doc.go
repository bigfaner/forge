// Package types defines typed constants, enumerations, and version information
// used throughout forge-cli. It is a leaf package with zero internal dependencies
// -- only the Go standard library is imported.
//
// # Sub-domains
//
//   - Status constants: task lifecycle states (pending, in_progress, completed, etc.)
//   - SurfaceType constants: interface surface types (web, api, cli, tui, mobile)
//   - Priority constants: task urgency levels (P0, P1, P2)
//   - Version information: CLI version and name, injected at build time via ldflags
//
// # Responsibility Boundaries
//
// This package owns only pure data definitions -- no business logic, no I/O, no
// external dependencies. Other pkg/ packages may import types but types must never
// import any other forge-cli pkg/ package.
package types
