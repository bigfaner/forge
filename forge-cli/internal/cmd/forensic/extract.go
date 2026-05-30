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

func runExtract(_ *cobra.Command, args []string) error {
	jsonlPath := args[0]

	// Resolve output directory: --slug > --out > auto-derive from session ID
	if slug != "" {
		outDir = filepath.Join("docs", "forensics", slug, "evidence")
	} else if outDir == "" {
		base := strings.TrimSuffix(filepath.Base(jsonlPath), ".jsonl")
		if len(base) >= 8 {
			outDir = filepath.Join("docs", "forensics", base, "evidence")
		}
	}

	f, err := os.Open(jsonlPath)
	if err != nil {
		return base.NewAIError(base.ErrNotFound, "Cannot open transcript", err.Error(), "", "")
	}
	defer func() { _ = f.Close() }()

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
			StopReasons:   map[string]int{},
		},
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

	var firstTS, lastTS, prevTS string
	// pending tool_use calls: toolUseID -> {timestamp, tool, line, detail}
	type pendingCall struct {
		ts     string
		tool   string
		line   int
		detail string
	}
	pending := map[string]pendingCall{}

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
					// Record pending tool_use for timing
					if block.ID != "" && entry.Timestamp != "" {
						pending[block.ID] = pendingCall{
							ts:     entry.Timestamp,
							tool:   block.Name,
							line:   lineNum,
							detail: truncate(string(inputJSON), 120),
						}
					}
				}
			}

			// Track stop reasons from assistant messages
			if entry.Message.StopReason != "" {
				result.Summary.StopReasons[entry.Message.StopReason]++
			}

			// Track timestamps for duration
			ts := entry.Timestamp
			if ts == "" && entry.Snapshot.Timestamp != "" {
				ts = entry.Snapshot.Timestamp
			}
			if ts != "" {
				if firstTS == "" {
					firstTS = ts
				}
				lastTS = ts
				if prevTS != "" {
					if dur := computeDurationMs(prevTS, ts); dur > 0 {
						result.Summary.ThinkingTurns = append(result.Summary.ThinkingTurns, thinkingTurn{
							Line:       lineNum,
							Seconds:    float64(dur) / 1000.0,
							StopReason: entry.Message.StopReason,
							Detail:     truncate(firstThinking(entry.Message.Content), 80),
						})
						result.Summary.TotalThinkingMs += dur
					}
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
					// Match with pending tool_use for timing
					if pc, ok := pending[block.ToolUseID]; ok {
						dur := computeDurationMs(pc.ts, entry.Timestamp)
						if dur >= 0 {
							result.Summary.TotalToolMs += dur
							result.Summary.TopSlowest = append(result.Summary.TopSlowest, timingEntry{
								Tool:    pc.tool,
								Line:    pc.line,
								Seconds: float64(dur) / 1000.0,
								Detail:  pc.detail,
							})
						}
						delete(pending, block.ToolUseID)
					}
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
				addToAgg(&result.Summary.HookBreakdown, att.HookName)
				if att.ExitCode != 0 {
					result.Summary.HookFailures++
				}
			case "edited_text_file":
				if att.Filename != "" {
					result.FilesEdited = append(result.FilesEdited, att.Filename)
				}
			case "compact_file_reference":
				result.Summary.CompactCount++
			case "plan_mode", "plan_mode_exit", "plan_mode_reentry":
				result.Summary.PlanModeCount++
			}
			// Track skill invocations from all attachment types
			for _, skill := range att.Skills {
				addToAgg(&result.Summary.SkillInvocations, skill.Name)
			}
		}
		entryTS := entry.Timestamp
		if entryTS == "" && entry.Snapshot.Timestamp != "" {
			entryTS = entry.Snapshot.Timestamp
		}
		if entryTS != "" {
			prevTS = entryTS
		}
	}

	result.Lines = lineNum

	// Aggregate timing by tool type and sort TopSlowest
	toolTiming := map[string]*timingAgg{}
	for _, t := range result.Summary.TopSlowest {
		agg, ok := toolTiming[t.Tool]
		if !ok {
			agg = &timingAgg{Tool: t.Tool}
			toolTiming[t.Tool] = agg
		}
		agg.Count++
		agg.Total += t.Seconds
		if t.Seconds > agg.Max {
			agg.Max = t.Seconds
		}
	}
	for _, agg := range toolTiming {
		agg.Average = agg.Total / float64(agg.Count)
		result.Summary.TimingByTool = append(result.Summary.TimingByTool, *agg)
	}
	sort.Slice(result.Summary.TimingByTool, func(i, j int) bool {
		return result.Summary.TimingByTool[i].Total > result.Summary.TimingByTool[j].Total
	})
	sort.Slice(result.Summary.TopSlowest, func(i, j int) bool {
		return result.Summary.TopSlowest[i].Seconds > result.Summary.TopSlowest[j].Seconds
	})
	// Keep top 20 slowest
	if len(result.Summary.TopSlowest) > 20 {
		result.Summary.TopSlowest = result.Summary.TopSlowest[:20]
	}

	// Compute subagent count from Agent tool calls
	for _, a := range result.Summary.AgentsSpawned {
		result.Summary.SubagentCount += a.Count
	}

	// Compute start/end time and duration
	if firstTS != "" {
		if t, err := parseTimestamp(firstTS); err == nil {
			result.Summary.StartTime = t.Format("2006-01-02 15:04:05")
		} else {
			result.Summary.StartTime = firstTS
		}
	}
	if lastTS != "" {
		if t, err := parseTimestamp(lastTS); err == nil {
			result.Summary.EndTime = t.Format("2006-01-02 15:04:05")
		} else {
			result.Summary.EndTime = lastTS
		}
	}
	if firstTS != "" && lastTS != "" {
		t1, err1 := parseTimestamp(firstTS)
		t2, err2 := parseTimestamp(lastTS)
		if err1 == nil && err2 == nil {
			d := t2.Sub(t1)
			switch {
			case d < time.Minute:
				result.Summary.Duration = fmt.Sprintf("%.0fs", d.Seconds())
			case d < time.Hour:
				result.Summary.Duration = fmt.Sprintf("%.1fmin", d.Minutes())
			default:
				result.Summary.Duration = fmt.Sprintf("%.1fh", d.Hours())
			}
		} else {
			result.Summary.Duration = fmt.Sprintf("%s / %s", firstTS, lastTS)
		}
	}
	out, _ := json.MarshalIndent(result, "", "  ")

	if outDir != "" {
		_ = os.MkdirAll(outDir, 0o755)
		outPath := filepath.Join(outDir, "evidence.json")
		if err := os.WriteFile(outPath, out, 0o644); err != nil {
			return base.NewAIError(base.ErrNotFound, "Cannot write evidence file", err.Error(), "", "")
		}
		copyFile(jsonlPath, filepath.Join(outDir, filepath.Base(jsonlPath)))
		fmt.Println(outPath)
		printTimingSummary(&result.Summary)
	} else {
		fmt.Println(string(out))
	}
	return nil
}

func printTimingSummary(s *extractSummary) {
	fmt.Fprintln(os.Stderr, "\nTiming Summary:")
	fmt.Fprintf(os.Stderr, "  Session: %s -> %s (%s)\n", s.StartTime, s.EndTime, s.Duration)
	fmt.Fprintf(os.Stderr, "  Tool time: %s  Thinking time: %s\n", formatDurationMs(s.TotalToolMs), formatDurationMs(s.TotalThinkingMs))
	if len(s.TimingByTool) > 0 {
		fmt.Fprintln(os.Stderr, "  By tool:")
		for _, t := range s.TimingByTool {
			fmt.Fprintf(os.Stderr, "    %-12s %dx  total=%s  avg=%s  max=%s\n",
				t.Tool, t.Count, formatSec(t.Total), formatSec(t.Average), formatSec(t.Max))
		}
	}
	if len(s.TopSlowest) > 0 {
		fmt.Fprintln(os.Stderr, "  Top slowest actions:")
		limit := 5
		if len(s.TopSlowest) < limit {
			limit = len(s.TopSlowest)
		}
		for _, t := range s.TopSlowest[:limit] {
			fmt.Fprintf(os.Stderr, "    %6s  %-12s  %s\n", formatSec(t.Seconds), t.Tool, truncate(t.Detail, 60))
		}
	}
	if len(s.ThinkingTurns) > 0 {
		fmt.Fprintln(os.Stderr, "  Thinking turns:")
		limit := 5
		if len(s.ThinkingTurns) < limit {
			limit = len(s.ThinkingTurns)
		}
		for _, t := range s.ThinkingTurns[:limit] {
			fmt.Fprintf(os.Stderr, "    %6s  %s\n", formatSec(t.Seconds), truncate(t.Detail, 60))
		}
	}
	fmt.Fprintln(os.Stderr)
}

func formatDurationMs(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.1fs", float64(ms)/1000.0)
}

func formatSec(s float64) string {
	if s < 1 {
		return fmt.Sprintf("%.0fms", s*1000)
	}
	return fmt.Sprintf("%.1fs", s)
}

func firstThinking(blocks []contentBlock) string {
	for _, block := range blocks {
		if block.Type == "thinking" && block.Thinking != "" {
			return block.Thinking
		}
	}
	for _, block := range blocks {
		if block.Type == "tool_use" {
			return block.Name + "(...)"
		}
	}
	return ""
}
