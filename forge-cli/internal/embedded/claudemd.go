// Package embedded provides go:embed access to bundled template files.
package embedded

import _ "embed"

// CLAUDEmdTemplate holds the embedded CLAUDE.md template content.
//
//go:embed claudemd_template.md
var CLAUDEmdTemplate string
