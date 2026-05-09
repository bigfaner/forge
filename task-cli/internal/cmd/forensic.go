package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	forensicKeyword string
	forensicSession string
	forensicSkill   string
	forensicLast    int
	forensicOutDir  string
	forensicSlug    string
)

var forensicCmd = &cobra.Command{
	Use:   "forensic",
	Short: "Analyze Claude Code session transcripts for agent deviation forensics",
	Long: `Extract and analyze evidence from Claude Code session transcripts.

Subcommands:
  search    Find sessions in ~/.claude/history.jsonl
  extract   Extract thinking/tool chains from a session JSONL
  subagents List subagent transcripts for a session`,
}

var forensicSearchCmd = &cobra.Command{
	Use:   "search [project-path]",
	Short: "Search history.jsonl for matching sessions",
	Args:  cobra.MaximumNArgs(1),
	Run:   runForensicSearch,
}

var forensicExtractCmd = &cobra.Command{
	Use:   "extract <session-jsonl-path>",
	Short: "Extract compact evidence from a session transcript",
	Long: `Extract thinking blocks, tool calls, hook events, and file edits from a session JSONL.

Output modes:
  --slug <name>   Write to docs/forensics/<name>/evidence/ (auto-creates dirs)
  --out <dir>     Write to arbitrary directory
  (default)       Print JSON to stdout`,
	Args: cobra.ExactArgs(1),
	Run:  runForensicExtract,
}

var forensicSubagentsCmd = &cobra.Command{
	Use:   "subagents <session-dir-path>",
	Short: "List subagent transcripts for a session",
	Args:  cobra.ExactArgs(1),
	Run:   runForensicSubagents,
}

func init() {
	forensicSearchCmd.Flags().StringVar(&forensicKeyword, "keyword", "", "Filter sessions by keyword in user messages")
	forensicSearchCmd.Flags().StringVar(&forensicSession, "session", "", "Filter by session ID prefix")
	forensicSearchCmd.Flags().StringVar(&forensicSkill, "skill", "", "Filter by skill name invoked in session")
	forensicSearchCmd.Flags().IntVar(&forensicLast, "last", 20, "Limit number of results")

	forensicExtractCmd.Flags().StringVar(&forensicOutDir, "out", "", "Write evidence JSON to directory (default: stdout)")
	forensicExtractCmd.Flags().StringVar(&forensicSlug, "slug", "", "Write to docs/forensics/<slug>/evidence/ (default: session ID prefix)")

	forensicCmd.AddCommand(forensicSearchCmd)
	forensicCmd.AddCommand(forensicExtractCmd)
	forensicCmd.AddCommand(forensicSubagentsCmd)
	rootCmd.AddCommand(forensicCmd)
}

// ── data types ──────────────────────────────────────────────────────

type historyEntry struct {
	Display   string `json:"display"`
	Timestamp int64  `json:"timestamp"`
	Project   string `json:"project"`
	SessionID string `json:"sessionId"`
}

type sessionSummary struct {
	SessionID string `json:"sessionId"`
	Project   string `json:"project"`
	DateTime  string `json:"dateTime"`
	MsgCount  int    `json:"msgCount"`
	FirstMsg  string `json:"firstMsg"`
}

type contentBlock struct {
	Type     string `json:"type"`
	Thinking string `json:"thinking,omitempty"`
	Text     string `json:"text,omitempty"`
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	ToolUseID string `json:"tool_use_id,omitempty"`
	Input    any    `json:"input,omitempty"`
}

type usageInfo struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type jsonlMessage struct {
	ID         string          `json:"id"`
	Role       string          `json:"role"`
	Content    []contentBlock  `json:"content"`
	RawContent json.RawMessage `json:"-"`
	StopReason string          `json:"stop_reason"`
	Model      string          `json:"model"`
	Usage      *usageInfo      `json:"usage,omitempty"`
}

func (m *jsonlMessage) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID         string          `json:"id"`
		Role       string          `json:"role"`
		Model      string          `json:"model"`
		StopReason string          `json:"stop_reason"`
		Usage      *usageInfo      `json:"usage,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	m.ID = raw.ID
	m.Role = raw.Role
	m.Model = raw.Model
	m.StopReason = raw.StopReason
	m.Usage = raw.Usage

	var rawContent struct {
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &rawContent); err == nil && len(rawContent.Content) > 0 {
		m.RawContent = rawContent.Content
		if rawContent.Content[0] == '[' {
			var blocks []contentBlock
			if json.Unmarshal(rawContent.Content, &blocks) == nil {
				m.Content = blocks
			}
		}
	}
	return nil
}

func (m *jsonlMessage) textContent() string {
	for _, block := range m.Content {
		if block.Type == "text" && block.Text != "" {
			return block.Text
		}
	}
	if len(m.RawContent) > 0 && m.RawContent[0] == '"' {
		var s string
		if json.Unmarshal(m.RawContent, &s) == nil && s != "" {
			return s
		}
	}
	return ""
}

type jsonlEntry struct {
	Type         string       `json:"type"`
	Message      jsonlMessage `json:"message"`
	Content      string       `json:"content"`
	GitBranch    string       `json:"gitBranch"`
	Attachment   attachment   `json:"attachment"`
	ToolUseResult *toolUseResult `json:"toolUseResult"`
}

type attachment struct {
	Type string `json:"type"`
	// invoked_skills
	Skills []invokedSkill `json:"skills,omitempty"`
	// hook_success
	HookName   string `json:"hookName,omitempty"`
	HookEvent  string `json:"hookEvent,omitempty"`
	DurationMs int    `json:"durationMs,omitempty"`
	ExitCode   int    `json:"exitCode,omitempty"`
	Command    string `json:"command,omitempty"`
	// edited_text_file / file
	Filename    string `json:"filename,omitempty"`
	DisplayPath string `json:"displayPath,omitempty"`
	// plan_mode
	PlanFilePath string `json:"planFilePath,omitempty"`
	// skill_listing
	SkillCount int  `json:"skillCount,omitempty"`
	IsInitial  bool `json:"isInitial,omitempty"`
}

type invokedSkill struct {
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
}

type toolUseResult struct {
	Type     string `json:"type"`
	FilePath string `json:"filePath,omitempty"`
}

type thinkingEntry struct {
	Line       int    `json:"line"`
	Thinking   string `json:"thinking"`
	StopReason string `json:"stopReason,omitempty"`
	Model      string `json:"model,omitempty"`
	MsgID      string `json:"msgId,omitempty"`
}

type toolCallEntry struct {
	Line       int    `json:"line"`
	Tool       string `json:"tool"`
	Input      string `json:"input"`
	StopReason string `json:"stopReason,omitempty"`
	MsgID      string `json:"msgId,omitempty"`
}

type toolResultEntry struct {
	Line       int    `json:"line"`
	ToolUseID  string `json:"toolUseId"`
	ResultType string `json:"resultType,omitempty"`
	FilePath   string `json:"filePath,omitempty"`
}

type userMsgEntry struct {
	Line    int    `json:"line"`
	Content string `json:"content"`
	IsMeta  bool   `json:"isMeta"`
}

type hookEventEntry struct {
	Line       int    `json:"line"`
	HookName   string `json:"hookName"`
	HookEvent  string `json:"hookEvent"`
	DurationMs int    `json:"durationMs"`
	ExitCode   int    `json:"exitCode"`
	Command    string `json:"command"`
}

type toolAggEntry struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	type extractSummary struct {
		TotalThinking    int            `json:"totalThinking"`
		TotalToolCalls   int            `json:"totalToolCalls"`
		TotalToolResults int            `json:"totalToolResults"`
		TotalUserMsgs    int            `json:"totalUserMsgs"`
		ToolBreakdown    map[string]int `json:"toolBreakdown"`
		FilesRead        []string       `json:"filesRead"`
		FilesWritten     []string       `json:"filesWritten"`
		GrepPatterns     []string       `json:"grepPatterns"`
		AgentsSpawned    []toolAggEntry `json:"agentsSpawned"`
		Commands         []string       `json:"commands"`
	}

type extractResult struct {
	File         string            `json:"file"`
	Lines        int               `json:"lines"`
	Model        string            `json:"model,omitempty"`
	GitBranch    string            `json:"gitBranch,omitempty"`
	Thinking     []thinkingEntry   `json:"thinking"`
	ToolCalls    []toolCallEntry   `json:"toolCalls"`
	ToolResults  []toolResultEntry `json:"toolResults"`
	UserMsgs     []userMsgEntry    `json:"userMsgs"`
	SkillsUsed   []string          `json:"skillsUsed"`
	Hooks        []hookEventEntry  `json:"hooks"`
	FilesEdited  []string          `json:"filesEdited"`
	Summary      extractSummary    `json:"summary"`
}

type subagentInfo struct {
	AgentID    string `json:"agentId"`
	AgentType  string `json:"agentType"`
	Transcript string `json:"transcript"`
}

// ── search ──────────────────────────────────────────────────────────

func runForensicSearch(cmd *cobra.Command, args []string) {
	projectPath := ""
	if len(args) > 0 {
		projectPath = args[0]
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		Exit(NewAIError(ErrNotFound, "Cannot determine home directory", err.Error(), "", ""))
	}
	histPath := filepath.Join(homeDir, ".claude", "history.jsonl")

	f, err := os.Open(histPath)
	if err != nil {
		Exit(NewAIError(ErrNotFound, "Cannot open history.jsonl", err.Error(), "", ""))
	}
	defer f.Close()

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
		if forensicSession != "" && !strings.HasPrefix(entry.SessionID, forensicSession) {
			continue
		}
		if forensicKeyword != "" && !strings.Contains(strings.ToLower(entry.Display), strings.ToLower(forensicKeyword)) {
			continue
		}
		if forensicSkill != "" {
			lower := strings.ToLower(entry.Display)
			if !strings.Contains(lower, "/"+strings.ToLower(forensicSkill)) &&
				!strings.Contains(lower, "forge:"+strings.ToLower(forensicSkill)) {
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

	if forensicLast < len(sorted) {
		sorted = sorted[:forensicLast]
	}

	out, _ := json.MarshalIndent(sorted, "", "  ")
	fmt.Println(string(out))
}

// ── extract ─────────────────────────────────────────────────────────

func runForensicExtract(cmd *cobra.Command, args []string) {
	jsonlPath := args[0]

	// Resolve output directory: --slug > --out > auto-derive from session ID
	if forensicSlug != "" {
		forensicOutDir = filepath.Join("docs", "forensics", forensicSlug, "evidence")
	} else if forensicOutDir == "" {
		base := strings.TrimSuffix(filepath.Base(jsonlPath), ".jsonl")
		if len(base) >= 8 {
			forensicOutDir = filepath.Join("docs", "forensics", base, "evidence")
		}
	}

	f, err := os.Open(jsonlPath)
	if err != nil {
		Exit(NewAIError(ErrNotFound, "Cannot open transcript", err.Error(), "", ""))
	}
	defer f.Close()

	result := extractResult{
		File:        jsonlPath,
		Thinking:    []thinkingEntry{},
		ToolCalls:   []toolCallEntry{},
		ToolResults: []toolResultEntry{},
		UserMsgs:    []userMsgEntry{},
		SkillsUsed:  []string{},
		Hooks:       []hookEventEntry{},
		FilesEdited: []string{},
		Summary: extractSummary{
			ToolBreakdown: map[string]int{},
		},
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		var entry jsonlEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if result.GitBranch == "" && entry.GitBranch != "" {
			result.GitBranch = entry.GitBranch
		}

		switch entry.Type {
		case "assistant":
			for _, block := range entry.Message.Content {
				switch block.Type {
				case "thinking":
					result.Summary.TotalThinking++
					result.Thinking = append(result.Thinking, thinkingEntry{
						Line:       lineNum,
						Thinking:   truncate(block.Thinking, 500),
						StopReason: entry.Message.StopReason,
						Model:      entry.Message.Model,
						MsgID:      entry.Message.ID,
					})
					if result.Model == "" && entry.Message.Model != "" {
						result.Model = entry.Message.Model
					}
				case "tool_use":
					result.Summary.TotalToolCalls++
					result.Summary.ToolBreakdown[block.Name]++
					inputJSON, _ := json.Marshal(block.Input)
					result.ToolCalls = append(result.ToolCalls, toolCallEntry{
						Line:       lineNum,
						Tool:       block.Name,
						Input:      truncate(string(inputJSON), 300),
						StopReason: entry.Message.StopReason,
						MsgID:      entry.Message.ID,
					})
					aggregateToolInput(block.Name, block.Input, &result.Summary)
				}
			}

		case "user":
			result.Summary.TotalUserMsgs++
			content := extractUserContent(entry)
			if content != "" {
				result.UserMsgs = append(result.UserMsgs, userMsgEntry{
					Line:    lineNum,
					Content: truncate(content, 300),
				})
				detectSkills(content, &result.SkillsUsed)
			}

			// Extract tool_result metadata
			for _, block := range entry.Message.Content {
				if block.Type == "tool_result" {
					result.Summary.TotalToolResults++
					tre := toolResultEntry{
						Line:      lineNum,
						ToolUseID: block.ToolUseID,
					}
					if entry.ToolUseResult != nil {
						tre.ResultType = entry.ToolUseResult.Type
						tre.FilePath = entry.ToolUseResult.FilePath
					}
					result.ToolResults = append(result.ToolResults, tre)
				}
			}

		case "attachment":
			att := entry.Attachment
			switch att.Type {
			case "invoked_skills":
				for _, skill := range att.Skills {
					name := strings.ToLower(strings.TrimPrefix(skill.Name, "forge:"))
					found := false
					for _, s := range result.SkillsUsed {
						if s == name {
							found = true
							break
						}
					}
					if !found {
						result.SkillsUsed = append(result.SkillsUsed, name)
					}
				}
			case "hook_success":
				result.Hooks = append(result.Hooks, hookEventEntry{
					Line:       lineNum,
					HookName:   att.HookName,
					HookEvent:  att.HookEvent,
					DurationMs: att.DurationMs,
					ExitCode:   att.ExitCode,
					Command:    att.Command,
				})
			case "edited_text_file":
				if att.Filename != "" {
					result.FilesEdited = append(result.FilesEdited, att.Filename)
				}
			}
		}
	}

	result.Lines = lineNum

	out, _ := json.MarshalIndent(result, "", "  ")

	if forensicOutDir != "" {
		os.MkdirAll(forensicOutDir, 0755)
		outPath := filepath.Join(forensicOutDir, "evidence.json")
		if err := os.WriteFile(outPath, out, 0644); err != nil {
			Exit(NewAIError(ErrNotFound, "Cannot write evidence file", err.Error(), "", ""))
		}
		copyFile(jsonlPath, filepath.Join(forensicOutDir, filepath.Base(jsonlPath)))
		fmt.Println(outPath)
	} else {
		fmt.Println(string(out))
	}
}

// ── subagents ───────────────────────────────────────────────────────

func runForensicSubagents(cmd *cobra.Command, args []string) {
	sessionDir := args[0]
	subDir := filepath.Join(sessionDir, "subagents")

	entries, err := os.ReadDir(subDir)
	if err != nil {
		Exit(NewAIError(ErrNotFound, "No subagents directory", err.Error(), "", ""))
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
		json.Unmarshal(data, &meta)

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
}

// ── helpers ─────────────────────────────────────────────────────────

func extractUserContent(entry jsonlEntry) string {
	if entry.Message.Role != "user" {
		return ""
	}
	if text := entry.Message.textContent(); text != "" {
		return text
	}
	if entry.Content != "" {
		return entry.Content
	}
	return ""
}

func detectSkills(content string, skills *[]string) {
	lower := strings.ToLower(content)
	prefixes := []string{"/forge:", "<command-name>/", "<command-name>"}

	for _, prefix := range prefixes {
		searchFrom := 0
		for {
			idx := strings.Index(lower[searchFrom:], prefix)
			if idx == -1 {
				break
			}
			idx += searchFrom
			start := idx + len(prefix)
			end := start
			for end < len(lower) && (isAlphaNumeric(lower[end]) || lower[end] == '-' || lower[end] == '_') {
				end++
			}
			if end > start {
				name := lower[start:end]
				found := false
				for _, s := range *skills {
					if s == name {
						found = true
						break
					}
				}
				if !found {
					*skills = append(*skills, name)
				}
			}
			searchFrom = end
		}
	}
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

func truncate(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "..."
}

func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	// Best-effort copy; don't fail the extract if copy fails
	io.Copy(out, in)
}

// aggregateToolInput extracts tool-specific details into the summary.
func aggregateToolInput(tool string, input any, s *extractSummary) {
	m, ok := input.(map[string]any)
	if !ok {
		return
	}

	switch tool {
	case "Read":
		if fp, _ := m["file_path"].(string); fp != "" {
			s.FilesRead = appendUniq(s.FilesRead, fp)
		}
	case "Edit", "Write":
		if fp, _ := m["file_path"].(string); fp != "" {
			s.FilesWritten = appendUniq(s.FilesWritten, fp)
		}
	case "Grep":
		if p, _ := m["pattern"].(string); p != "" {
			s.GrepPatterns = appendUniq(s.GrepPatterns, p)
		}
	case "Bash":
		if cmd, _ := m["command"].(string); cmd != "" {
			s.Commands = appendUniq(s.Commands, truncate(cmd, 200))
		}
	case "Agent":
		name := "unknown"
		if t, _ := m["subagent_type"].(string); t != "" {
			name = t
		}
		found := false
		for i := range s.AgentsSpawned {
			if s.AgentsSpawned[i].Name == name {
				s.AgentsSpawned[i].Count++
				found = true
				break
			}
		}
		if !found {
			s.AgentsSpawned = append(s.AgentsSpawned, toolAggEntry{Name: name, Count: 1})
		}
	}
}

func appendUniq(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}
