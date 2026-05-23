package forensic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"forge-cli/internal/cmd/base"

	"github.com/spf13/cobra"
)

func runSearch(_ *cobra.Command, args []string) error {
	projectPath := ""
	if len(args) > 0 {
		projectPath = args[0]
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return base.NewAIError(base.ErrNotFound, "Cannot determine home directory", err.Error(), "", "")
	}
	histPath := filepath.Join(homeDir, ".claude", "history.jsonl")

	f, err := os.Open(histPath)
	if err != nil {
		return base.NewAIError(base.ErrNotFound, "Cannot open history.jsonl", err.Error(), "", "")
	}
	defer func() { _ = f.Close() }()

	searchWithProjectPath(projectPath, f)
	return nil
}

// searchWithProjectPath runs search logic against an already-opened file handle.
// Exported for test use.
func searchWithProjectPath(projectPath string, f *os.File) {
	sessions := map[string]*sessionSummary{}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

	for scanner.Scan() {
		var entry historyEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if projectPath != "" && !strings.Contains(entry.Project, projectPath) {
			continue
		}
		if entry.SessionID == "" {
			continue
		}
		if session != "" && !strings.HasPrefix(entry.SessionID, session) {
			continue
		}
		if keyword != "" && !strings.Contains(strings.ToLower(entry.Display), strings.ToLower(keyword)) {
			continue
		}
		if skill != "" {
			lower := strings.ToLower(entry.Display)
			if !strings.Contains(lower, "/"+strings.ToLower(skill)) &&
				!strings.Contains(lower, "forge:"+strings.ToLower(skill)) {
				continue
			}
		}

		ss, exists := sessions[entry.SessionID]
		if !exists {
			ss = &sessionSummary{
				SessionID: entry.SessionID,
				Project:   entry.Project,
			}
			sessions[entry.SessionID] = ss
		}
		ss.MsgCount++
		if entry.Timestamp > 0 {
			ts := time.UnixMilli(entry.Timestamp)
			if ss.DateTime == "" || ts.Format("2006-01-02 15:04") > ss.DateTime {
				ss.DateTime = ts.Format("2006-01-02 15:04")
			}
		}
		if ss.FirstMsg == "" {
			ss.FirstMsg = truncate(entry.Display, 80)
		}
	}

	sorted := make([]*sessionSummary, 0, len(sessions))
	for _, ss := range sessions {
		sorted = append(sorted, ss)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].DateTime > sorted[j].DateTime
	})

	if last < len(sorted) {
		sorted = sorted[:last]
	}

	out, _ := json.MarshalIndent(sorted, "", "  ")
	fmt.Println(string(out))
}
