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
	Long: `Search Claude's history.jsonl for sessions matching the given criteria.

Reads ~/.claude/history.jsonl and filters by project path, session ID prefix,
keyword, or skill name. Results are sorted by most recent first and output as
a JSON array of session summaries, each containing:
  sessionId, project, dateTime, msgCount, firstMsg

Flags:
  --keyword   Filter sessions by keyword in user messages
  --session   Filter by session ID prefix
  --skill     Filter by skill name invoked in session
  --last      Limit number of results (default: 20)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSearch,
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
	Long: `List subagent transcripts and metadata for a given session directory.

Reads <session-dir>/subagents/*.meta.json files and outputs a JSON array of
subagent entries, each containing:
  agentId    — subagent identifier (derived from filename)
  agentType  — type of agent (from metadata)
  transcript — path to the subagent's JSONL transcript file`,
	Args: cobra.ExactArgs(1),
	RunE: runSubagents,
}

func init() {
	searchCmd.Flags().StringVar(&keyword, "keyword", "", "Filter sessions by keyword in user messages")
	searchCmd.Flags().StringVar(&session, "session", "", "Filter by session ID prefix")
	searchCmd.Flags().StringVar(&skill, "skill", "", "Filter by skill name invoked in session")
	searchCmd.Flags().IntVar(&last, "last", 20, "Limit number of results")

	extractCmd.Flags().StringVar(&outDir, "out", "", "Write evidence JSON to directory (default: stdout)")
	extractCmd.Flags().StringVar(&slug, "slug", "", "Write to docs/forensics/<slug>/evidence/ (default: session ID prefix)")
}
