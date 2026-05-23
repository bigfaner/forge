package forensic

import "github.com/spf13/cobra"

var (
	keyword string
	session string
	skill   string
	last    int
	outDir  string
	slug    string
)

var searchCmd = &cobra.Command{
	Use:   "search [project-path]",
	Short: "Search history.jsonl for matching sessions",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSearch,
}

var extractCmd = &cobra.Command{
	Use:   "extract <session-jsonl-path>",
	Short: "Extract compact evidence from a session transcript",
	Long: `Extract thinking blocks, tool calls, hook events, and file edits from a session JSONL.

Output modes:
  --slug <name>   Write to docs/forensics/<name>/evidence/ (auto-creates dirs)
  --out <dir>     Write to arbitrary directory
  (default)       Print JSON to stdout`,
	Args: cobra.ExactArgs(1),
	RunE: runExtract,
}

var subagentsCmd = &cobra.Command{
	Use:   "subagents <session-dir-path>",
	Short: "List subagent transcripts for a session",
	Args:  cobra.ExactArgs(1),
	RunE:  runSubagents,
}

func init() {
	searchCmd.Flags().StringVar(&keyword, "keyword", "", "Filter sessions by keyword in user messages")
	searchCmd.Flags().StringVar(&session, "session", "", "Filter by session ID prefix")
	searchCmd.Flags().StringVar(&skill, "skill", "", "Filter by skill name invoked in session")
	searchCmd.Flags().IntVar(&last, "last", 20, "Limit number of results")

	extractCmd.Flags().StringVar(&outDir, "out", "", "Write evidence JSON to directory (default: stdout)")
	extractCmd.Flags().StringVar(&slug, "slug", "", "Write to docs/forensics/<slug>/evidence/ (default: session ID prefix)")
}
