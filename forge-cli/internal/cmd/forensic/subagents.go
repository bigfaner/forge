package forensic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/internal/cmd/base"

	"github.com/spf13/cobra"
)

func runSubagents(_ *cobra.Command, args []string) error {
	sessionDir := args[0]
	subDir := filepath.Join(sessionDir, "subagents")

	entries, err := os.ReadDir(subDir)
	if err != nil {
		return base.NewAIError(base.ErrNotFound, "No subagents directory", err.Error(), "", "")
	}

	agents := []subagentInfo{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".meta.json") {
			continue
		}

		metaPath := filepath.Join(subDir, e.Name())
		data, err := os.ReadFile(metaPath)
		if err != nil {
			continue
		}

		var meta map[string]string
		_ = json.Unmarshal(data, &meta)

		base := strings.TrimSuffix(e.Name(), ".meta.json")
		transcriptPath := filepath.Join(subDir, base+".jsonl")

		agents = append(agents, subagentInfo{
			AgentID:    strings.TrimPrefix(base, "agent-"),
			AgentType:  meta["agentType"],
			Transcript: transcriptPath,
		})
	}

	out, _ := json.MarshalIndent(agents, "", "  ")
	fmt.Println(string(out))
	return nil
}
