package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ── truncate ────────────────────────────────────────────────────────

func TestTruncate_Short(t *testing.T) {
	got := truncate("hello", 10)
	if got != "hello" {
		t.Errorf("truncate(%q, 10) = %q, want %q", "hello", got, "hello")
	}
}

func TestTruncate_Exact(t *testing.T) {
	got := truncate("hello", 5)
	if got != "hello" {
		t.Errorf("truncate(%q, 5) = %q, want %q", "hello", got, "hello")
	}
}

func TestTruncate_Long(t *testing.T) {
	got := truncate("hello world", 5)
	if got != "hello..." {
		t.Errorf("truncate(%q, 5) = %q, want %q", "hello world", got, "hello...")
	}
}

func TestTruncate_Zero(t *testing.T) {
	got := truncate("hello", 0)
	if got != "..." {
		t.Errorf("truncate(%q, 0) = %q, want %q", "hello", got, "...")
	}
}

// ── isAlphaNumeric ──────────────────────────────────────────────────

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		char byte
		want bool
	}{
		{'a', true},
		{'z', true},
		{'0', true},
		{'9', true},
		{'A', false},
		{'Z', false},
		{'-', false},
		{'_', false},
		{' ', false},
	}
	for _, tt := range tests {
		got := isAlphaNumeric(tt.char)
		if got != tt.want {
			t.Errorf("isAlphaNumeric(%q) = %v, want %v", tt.char, got, tt.want)
		}
	}
}

// ── extractUserContent ───────────────────────────────────────────────

func TestExtractUserContent_NonUserRole(t *testing.T) {
	entry := jsonlEntry{
		Message: jsonlMessage{Role: "assistant"},
	}
	got := extractUserContent(entry)
	if got != "" {
		t.Errorf("expected empty for non-user role, got %q", got)
	}
}

func TestExtractUserContent_ContentBlocks(t *testing.T) {
	entry := jsonlEntry{
		Message: jsonlMessage{
			Role: "user",
			Content: []contentBlock{
				{Type: "text", Text: "hello from blocks"},
			},
		},
	}
	got := extractUserContent(entry)
	if got != "hello from blocks" {
		t.Errorf("got %q, want %q", got, "hello from blocks")
	}
}

func TestExtractUserContent_RawString(t *testing.T) {
	// Simulate user message with content as raw string (not array of blocks)
	rawContent, _ := json.Marshal("raw string content")
	msg := jsonlMessage{Role: "user", RawContent: rawContent}
	entry := jsonlEntry{Message: msg}

	got := extractUserContent(entry)
	if got != "raw string content" {
		t.Errorf("got %q, want %q", got, "raw string content")
	}
}

func TestExtractUserContent_EntryContentFallback(t *testing.T) {
	entry := jsonlEntry{
		Message: jsonlMessage{Role: "user"},
		Content: "entry-level content",
	}
	got := extractUserContent(entry)
	if got != "entry-level content" {
		t.Errorf("got %q, want %q", got, "entry-level content")
	}
}

func TestExtractUserContent_EmptyAll(t *testing.T) {
	entry := jsonlEntry{
		Message: jsonlMessage{Role: "user"},
	}
	got := extractUserContent(entry)
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestExtractUserContent_BlockPriorityOverRaw(t *testing.T) {
	rawContent, _ := json.Marshal("raw fallback")
	entry := jsonlEntry{
		Message: jsonlMessage{
			Role: "user",
			Content: []contentBlock{
				{Type: "text", Text: "block text"},
			},
			RawContent: rawContent,
		},
	}
	got := extractUserContent(entry)
	if got != "block text" {
		t.Errorf("blocks should take priority, got %q", got)
	}
}

// ── detectSkills ────────────────────────────────────────────────────

func TestDetectSkills_ForgePrefix(t *testing.T) {
	skills := []string{}
	detectSkills("running /forge:eval-prd now", &skills)
	if len(skills) != 1 || skills[0] != "eval-prd" {
		t.Errorf("got %v, want [eval-prd]", skills)
	}
}

func TestDetectSkills_CommandNameTag(t *testing.T) {
	skills := []string{}
	detectSkills(`<command-name>/run-tasks</command-name>`, &skills)
	if len(skills) != 1 || skills[0] != "run-tasks" {
		t.Errorf("got %v, want [run-tasks]", skills)
	}
}

func TestDetectSkills_CommandNameNoSlash(t *testing.T) {
	skills := []string{}
	detectSkills(`<command-name>brainstorm</command-name>`, &skills)
	if len(skills) != 1 || skills[0] != "brainstorm" {
		t.Errorf("got %v, want [brainstorm]", skills)
	}
}

func TestDetectSkills_MultipleSamePrefix(t *testing.T) {
	skills := []string{}
	detectSkills("first /forge:eval-prd then /forge:eval-design", &skills)
	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %v", skills)
	}
	if skills[0] != "eval-prd" || skills[1] != "eval-design" {
		t.Errorf("got %v, want [eval-prd, eval-design]", skills)
	}
}

func TestDetectSkills_Dedup(t *testing.T) {
	skills := []string{}
	detectSkills("/forge:eval-prd and again /forge:eval-prd", &skills)
	if len(skills) != 1 {
		t.Errorf("expected dedup, got %v", skills)
	}
}

func TestDetectSkills_NoMatch(t *testing.T) {
	skills := []string{}
	detectSkills("just a regular message", &skills)
	if len(skills) != 0 {
		t.Errorf("expected empty, got %v", skills)
	}
}

func TestDetectSkills_CaseInsensitive(t *testing.T) {
	skills := []string{}
	detectSkills("/forge:Write-PRD here", &skills)
	if len(skills) != 1 || skills[0] != "write-prd" {
		t.Errorf("got %v, want [write-prd]", skills)
	}
}

func TestDetectSkills_HyphenAndUnderscore(t *testing.T) {
	skills := []string{}
	detectSkills("/forge:gen_test-scripts", &skills)
	if len(skills) != 1 || skills[0] != "gen_test-scripts" {
		t.Errorf("got %v, want [gen_test-scripts]", skills)
	}
}

func TestDetectSkills_SkillNameBoundary(t *testing.T) {
	skills := []string{}
	detectSkills("/forge:eval-prd some text after", &skills)
	if len(skills) != 1 || skills[0] != "eval-prd" {
		t.Errorf("got %v, want [eval-prd]", skills)
	}
}

// ── runForensicExtract ───────────────────────────────────────────────

func TestForensicExtract_ThinkingAndToolCalls(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	entries := []map[string]any{
		{
			"type": "assistant",
			"message": map[string]any{
				"id":          "msg1",
				"role":        "assistant",
				"model":       "glm-5",
				"stop_reason": "tool_use",
				"content": []map[string]any{
					{"type": "thinking", "thinking": "I should read the file first"},
					{"type": "tool_use", "name": "Read", "id": "call1", "input": map[string]any{"file_path": "/some/file.go"}},
				},
			},
		},
		{
			"type": "user",
			"message": map[string]any{
				"role": "user",
				"content": []map[string]any{
					{"type": "tool_result", "tool_use_id": "call1", "content": "file contents here"},
				},
			},
			"toolUseResult": map[string]any{
				"type":     "file",
				"filePath": "/some/file.go",
			},
		},
		{
			"type": "assistant",
			"message": map[string]any{
				"id":          "msg2",
				"role":        "assistant",
				"model":       "glm-5",
				"stop_reason": "end_turn",
				"content": []map[string]any{
					{"type": "tool_use", "name": "Edit", "id": "call2", "input": map[string]any{"file_path": "/some/file.go", "old_string": "old", "new_string": "new"}},
				},
			},
		},
		{
			"type":      "user",
			"gitBranch": "main",
			"message": map[string]any{
				"role":    "user",
				"content": "Running /forge:eval-prd now",
			},
		},
	}

	f, _ := os.Create(jsonlPath)
	for _, entry := range entries {
		line, _ := json.Marshal(entry)
		_, _ = f.Write(line)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Close()

	outDir := filepath.Join(dir, "evidence")
	forensicOutDir = outDir

	out := captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	if !strings.Contains(out, outDir) {
		t.Fatalf("expected output path, got: %s", out)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	if err != nil {
		t.Fatal(err)
	}

	var result extractResult
	_ = json.Unmarshal(data, &result)

	if result.Lines != 4 {
		t.Errorf("Lines = %d, want 4", result.Lines)
	}
	if result.Model != "glm-5" {
		t.Errorf("Model = %q, want glm-5", result.Model)
	}
	if result.GitBranch != "main" {
		t.Errorf("GitBranch = %q, want main", result.GitBranch)
	}
	if result.Summary.TotalThinking != 1 {
		t.Errorf("TotalThinking = %d, want 1", result.Summary.TotalThinking)
	}
	if result.Summary.TotalToolCalls != 2 {
		t.Errorf("TotalToolCalls = %d, want 2", result.Summary.TotalToolCalls)
	}
	if result.Summary.ToolBreakdown["Read"] != 1 {
		t.Errorf("ToolBreakdown[Read] = %d, want 1", result.Summary.ToolBreakdown["Read"])
	}
	if result.Summary.ToolBreakdown["Edit"] != 1 {
		t.Errorf("ToolBreakdown[Edit] = %d, want 1", result.Summary.ToolBreakdown["Edit"])
	}
	if len(result.Thinking) != 1 || !strings.Contains(result.Thinking[0].Thinking, "read the file") {
		t.Errorf("Thinking = %v, unexpected", result.Thinking)
	}
	if len(result.ToolCalls) != 2 {
		t.Errorf("ToolCalls len = %d, want 2", len(result.ToolCalls))
	}
	if len(result.SkillsUsed) != 1 || result.SkillsUsed[0] != "eval-prd" {
		t.Errorf("SkillsUsed = %v, want [eval-prd]", result.SkillsUsed)
	}

	// New fields: stop_reason on thinking/tool entries
	if result.Thinking[0].StopReason != "tool_use" {
		t.Errorf("Thinking[0].StopReason = %q, want tool_use", result.Thinking[0].StopReason)
	}
	if result.Thinking[0].Model != "glm-5" {
		t.Errorf("Thinking[0].Model = %q, want glm-5", result.Thinking[0].Model)
	}
	if result.Thinking[0].MsgID != "msg1" {
		t.Errorf("Thinking[0].MsgID = %q, want msg1", result.Thinking[0].MsgID)
	}
	if result.ToolCalls[0].StopReason != "tool_use" {
		t.Errorf("ToolCalls[0].StopReason = %q, want tool_use", result.ToolCalls[0].StopReason)
	}
	if result.ToolCalls[1].StopReason != "end_turn" {
		t.Errorf("ToolCalls[1].StopReason = %q, want end_turn", result.ToolCalls[1].StopReason)
	}

	// New fields: tool results
	if result.Summary.TotalToolResults != 1 {
		t.Errorf("TotalToolResults = %d, want 1", result.Summary.TotalToolResults)
	}
	if len(result.ToolResults) != 1 {
		t.Fatalf("ToolResults len = %d, want 1", len(result.ToolResults))
	}
	if result.ToolResults[0].ToolUseID != "call1" {
		t.Errorf("ToolResults[0].ToolUseID = %q, want call1", result.ToolResults[0].ToolUseID)
	}
	if result.ToolResults[0].ResultType != "file" {
		t.Errorf("ToolResults[0].ResultType = %q, want file", result.ToolResults[0].ResultType)
	}
	if result.ToolResults[0].FilePath != "/some/file.go" {
		t.Errorf("ToolResults[0].FilePath = %q, want /some/file.go", result.ToolResults[0].FilePath)
	}
}

func TestForensicExtract_InvalidJSONLines(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "bad.jsonl")

	f, _ := os.Create(jsonlPath)
	_, _ = f.WriteString("not json at all\n")
	_, _ = f.WriteString(`{"type":"user","message":{"role":"user","content":"hello"}}` + "\n")
	_ = f.Close()

	outDir := filepath.Join(t.TempDir(), "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})
	data, _ := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	var result extractResult
	_ = json.Unmarshal(data, &result)

	if result.Lines != 2 {
		t.Errorf("Lines = %d, want 2 (bad line still counted)", result.Lines)
	}
	if result.Summary.TotalUserMsgs != 1 {
		t.Errorf("TotalUserMsgs = %d, want 1", result.Summary.TotalUserMsgs)
	}
}

func TestForensicExtract_FileNotFound(t *testing.T) {
	if os.Getenv("TEST_FORENSIC_EXTRACT_MISSING") == "1" {
		runForensicExtract(nil, []string{"/nonexistent/file.jsonl"})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestForensicExtract_FileNotFound")
	cmd.Env = append(os.Environ(), "TEST_FORENSIC_EXTRACT_MISSING=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for missing file")
	}
	if !strings.Contains(string(output), "NOT_FOUND") {
		t.Errorf("expected NOT_FOUND error, got: %s", string(output))
	}
}

func TestForensicExtract_AttachmentInvokedSkills(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	entries := []map[string]any{
		{
			"type": "attachment",
			"attachment": map[string]any{
				"type": "invoked_skills",
				"skills": []map[string]any{
					{"name": "forge:eval-prd", "path": "plugins/forge/skills/eval-prd/SKILL.md"},
					{"name": "forge:run-tasks", "path": "plugins/forge/skills/run-tasks/SKILL.md"},
				},
			},
		},
		{
			"type": "user",
			"message": map[string]any{
				"role":    "user",
				"content": "some text without skill names",
			},
		},
	}

	f, _ := os.Create(jsonlPath)
	for _, entry := range entries {
		line, _ := json.Marshal(entry)
		_, _ = f.Write(line)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Close()

	outDir := filepath.Join(dir, "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	data, _ := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	var result extractResult
	_ = json.Unmarshal(data, &result)

	if len(result.SkillsUsed) != 2 {
		t.Fatalf("SkillsUsed = %v, want 2 skills", result.SkillsUsed)
	}
	if result.SkillsUsed[0] != "eval-prd" {
		t.Errorf("SkillsUsed[0] = %q, want eval-prd", result.SkillsUsed[0])
	}
	if result.SkillsUsed[1] != "run-tasks" {
		t.Errorf("SkillsUsed[1] = %q, want run-tasks", result.SkillsUsed[1])
	}
}

func TestForensicExtract_AttachmentInvokedSkillsDedup(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	entries := []map[string]any{
		{
			"type": "attachment",
			"attachment": map[string]any{
				"type": "invoked_skills",
				"skills": []map[string]any{
					{"name": "forge:eval-prd"},
				},
			},
		},
		{
			"type": "user",
			"message": map[string]any{
				"role":    "user",
				"content": "/forge:eval-prd again",
			},
		},
	}

	f, _ := os.Create(jsonlPath)
	for _, entry := range entries {
		line, _ := json.Marshal(entry)
		_, _ = f.Write(line)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Close()

	outDir := filepath.Join(dir, "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	data, _ := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	var result extractResult
	_ = json.Unmarshal(data, &result)

	if len(result.SkillsUsed) != 1 || result.SkillsUsed[0] != "eval-prd" {
		t.Errorf("SkillsUsed should dedup, got %v", result.SkillsUsed)
	}
}

func TestForensicExtract_HookEvents(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	entries := []map[string]any{
		{
			"type": "attachment",
			"attachment": map[string]any{
				"type":       "hook_success",
				"hookName":   "pre-commit",
				"hookEvent":  "PreToolUse",
				"durationMs": 150,
				"exitCode":   0,
				"command":    "golangci-lint run",
			},
		},
		{
			"type": "attachment",
			"attachment": map[string]any{
				"type":       "hook_success",
				"hookName":   "post-edit",
				"hookEvent":  "PostToolUse",
				"durationMs": 50,
				"exitCode":   1,
				"command":    "go vet ./...",
			},
		},
	}

	f, _ := os.Create(jsonlPath)
	for _, entry := range entries {
		line, _ := json.Marshal(entry)
		_, _ = f.Write(line)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Close()

	outDir := filepath.Join(dir, "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	data, _ := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	var result extractResult
	_ = json.Unmarshal(data, &result)

	if len(result.Hooks) != 2 {
		t.Fatalf("Hooks len = %d, want 2", len(result.Hooks))
	}
	if result.Hooks[0].HookName != "pre-commit" {
		t.Errorf("Hooks[0].HookName = %q, want pre-commit", result.Hooks[0].HookName)
	}
	if result.Hooks[0].HookEvent != "PreToolUse" {
		t.Errorf("Hooks[0].HookEvent = %q, want PreToolUse", result.Hooks[0].HookEvent)
	}
	if result.Hooks[0].DurationMs != 150 {
		t.Errorf("Hooks[0].DurationMs = %d, want 150", result.Hooks[0].DurationMs)
	}
	if result.Hooks[0].ExitCode != 0 {
		t.Errorf("Hooks[0].ExitCode = %d, want 0", result.Hooks[0].ExitCode)
	}
	if result.Hooks[1].ExitCode != 1 {
		t.Errorf("Hooks[1].ExitCode = %d, want 1", result.Hooks[1].ExitCode)
	}
	if result.Hooks[1].Command != "go vet ./..." {
		t.Errorf("Hooks[1].Command = %q, want go vet ./...", result.Hooks[1].Command)
	}
}

func TestForensicExtract_EditedFiles(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	entries := []map[string]any{
		{
			"type": "attachment",
			"attachment": map[string]any{
				"type":     "edited_text_file",
				"filename": "internal/cmd/forensic.go",
			},
		},
		{
			"type": "attachment",
			"attachment": map[string]any{
				"type":     "edited_text_file",
				"filename": "internal/cmd/forensic_test.go",
			},
		},
		{
			"type": "attachment",
			"attachment": map[string]any{
				"type":     "edited_text_file",
				"filename": "",
			},
		},
	}

	f, _ := os.Create(jsonlPath)
	for _, entry := range entries {
		line, _ := json.Marshal(entry)
		_, _ = f.Write(line)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Close()

	outDir := filepath.Join(dir, "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	data, _ := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	var result extractResult
	_ = json.Unmarshal(data, &result)

	if len(result.FilesEdited) != 2 {
		t.Fatalf("FilesEdited = %v, want 2 files", result.FilesEdited)
	}
	if result.FilesEdited[0] != "internal/cmd/forensic.go" {
		t.Errorf("FilesEdited[0] = %q, want internal/cmd/forensic.go", result.FilesEdited[0])
	}
	if result.FilesEdited[1] != "internal/cmd/forensic_test.go" {
		t.Errorf("FilesEdited[1] = %q, want internal/cmd/forensic_test.go", result.FilesEdited[1])
	}
}

func TestForensicExtract_ToolResultWithoutMetadata(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	entries := []map[string]any{
		{
			"type": "user",
			"message": map[string]any{
				"role": "user",
				"content": []map[string]any{
					{"type": "tool_result", "tool_use_id": "call99", "content": "done"},
				},
			},
			// no toolUseResult field
		},
	}

	f, _ := os.Create(jsonlPath)
	for _, entry := range entries {
		line, _ := json.Marshal(entry)
		_, _ = f.Write(line)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Close()

	outDir := filepath.Join(dir, "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	data, _ := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	var result extractResult
	_ = json.Unmarshal(data, &result)

	if result.Summary.TotalToolResults != 1 {
		t.Errorf("TotalToolResults = %d, want 1", result.Summary.TotalToolResults)
	}
	if len(result.ToolResults) != 1 {
		t.Fatalf("ToolResults len = %d, want 1", len(result.ToolResults))
	}
	if result.ToolResults[0].ToolUseID != "call99" {
		t.Errorf("ToolUseID = %q, want call99", result.ToolResults[0].ToolUseID)
	}
	if result.ToolResults[0].ResultType != "" {
		t.Errorf("ResultType should be empty without metadata, got %q", result.ToolResults[0].ResultType)
	}
}

func TestAppendUniq(t *testing.T) {
	s := []string{}
	s = appendUniq(s, "a")
	s = appendUniq(s, "b")
	s = appendUniq(s, "a")
	if len(s) != 2 {
		t.Errorf("expected 2, got %d: %v", len(s), s)
	}
}

func TestAggregateToolInput(t *testing.T) {
	s := extractSummary{ToolBreakdown: map[string]int{}}

	aggregateToolInput("Read", map[string]any{"file_path": "/a.go"}, &s)
	aggregateToolInput("Read", map[string]any{"file_path": "/b.go"}, &s)
	aggregateToolInput("Read", map[string]any{"file_path": "/a.go"}, &s)
	if len(s.FilesRead) != 2 {
		t.Errorf("FilesRead = %v, want 2 unique", s.FilesRead)
	}

	aggregateToolInput("Edit", map[string]any{"file_path": "/c.go"}, &s)
	aggregateToolInput("Write", map[string]any{"file_path": "/d.go"}, &s)
	if len(s.FilesWritten) != 2 {
		t.Errorf("FilesWritten = %v, want 2", s.FilesWritten)
	}

	aggregateToolInput("Grep", map[string]any{"pattern": "TODO", "path": "."}, &s)
	aggregateToolInput("Grep", map[string]any{"pattern": "FIXME", "path": "."}, &s)
	aggregateToolInput("Grep", map[string]any{"pattern": "TODO", "path": "."}, &s)
	if len(s.GrepPatterns) != 2 {
		t.Errorf("GrepPatterns = %v, want 2 unique", s.GrepPatterns)
	}

	aggregateToolInput("Bash", map[string]any{"command": "go test ./..."}, &s)
	aggregateToolInput("Bash", map[string]any{"command": "go build ./..."}, &s)
	if len(s.Commands) != 2 {
		t.Errorf("Commands = %v, want 2", s.Commands)
	}

	aggregateToolInput("Agent", map[string]any{"subagent_type": "Explore"}, &s)
	aggregateToolInput("Agent", map[string]any{"subagent_type": "Explore"}, &s)
	aggregateToolInput("Agent", map[string]any{"subagent_type": "general-purpose"}, &s)
	if len(s.AgentsSpawned) != 2 {
		t.Fatalf("AgentsSpawned = %v, want 2", s.AgentsSpawned)
	}
	if s.AgentsSpawned[0].Count != 2 {
		t.Errorf("Explore count = %d, want 2", s.AgentsSpawned[0].Count)
	}

	// Non-map input should be safe
	aggregateToolInput("Read", "not a map", &s)
}

func TestGolden_AggregationPopulated(t *testing.T) {
	result := extractTestdata(t, "fix-bug.jsonl")

	if len(result.Summary.FilesRead) == 0 {
		t.Error("FilesRead should be populated")
	}
	if len(result.Summary.Commands) == 0 {
		t.Error("Commands should be populated")
	}
	for _, f := range result.Summary.FilesRead {
		if f == "" {
			t.Error("FilesRead should not contain empty strings")
		}
	}
	for _, c := range result.Summary.Commands {
		if c == "" {
			t.Error("Commands should not contain empty strings")
		}
	}

	t.Logf("filesRead=%d filesWritten=%d agents=%v grepPatterns=%d commands=%d",
		len(result.Summary.FilesRead), len(result.Summary.FilesWritten),
		result.Summary.AgentsSpawned, len(result.Summary.GrepPatterns),
		len(result.Summary.Commands))
}

func TestForensicExtract_CopiesSourceJSONL(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	content := []byte(`{"type":"user","message":{"role":"user","content":"hello"}}` + "\n")
	_ = os.WriteFile(jsonlPath, content, 0644)

	outDir := filepath.Join(dir, "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	// Source JSONL should be copied alongside evidence.json
	copiedPath := filepath.Join(outDir, "session.jsonl")
	data, err := os.ReadFile(copiedPath)
	if err != nil {
		t.Fatalf("source JSONL not copied to output dir: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("copied content mismatch:\ngot  %q\nwant %q", string(data), string(content))
	}
}

func TestForensicExtract_NoCopyWithoutOutDir(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")

	_ = os.WriteFile(jsonlPath, []byte(`{"type":"user","message":{"role":"user","content":"hello"}}`+"\n"), 0644)

	forensicOutDir = ""
	// Stdout mode — no file operations, no copy
	out := captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})
	if !strings.Contains(out, "hello") {
		t.Errorf("stdout mode should output JSON, got: %s", out)
	}
}

// ── runForensicSubagents ────────────────────────────────────────────

func TestForensicSubagents_WithMeta(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "subagents")
	_ = os.MkdirAll(subDir, 0755)

	// Create meta file
	meta := map[string]string{"agentType": "Explore"}
	metaJSON, _ := json.Marshal(meta)
	_ = os.WriteFile(filepath.Join(subDir, "agent-abc123.meta.json"), metaJSON, 0644)

	// Create transcript file (empty)
	_ = os.WriteFile(filepath.Join(subDir, "agent-abc123.jsonl"), []byte(""), 0644)

	out := captureStdout(func() {
		runForensicSubagents(nil, []string{dir})
	})

	var agents []subagentInfo
	_ = json.Unmarshal([]byte(strings.TrimSpace(out)), &agents)

	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}
	if agents[0].AgentID != "abc123" {
		t.Errorf("AgentID = %q, want abc123", agents[0].AgentID)
	}
	if agents[0].AgentType != "Explore" {
		t.Errorf("AgentType = %q, want Explore", agents[0].AgentType)
	}
}

func TestForensicSubagents_SkipsDirs(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "subagents")
	_ = os.MkdirAll(filepath.Join(subDir, "somedir"), 0755)
	// Only a directory, no .meta.json files

	out := captureStdout(func() {
		runForensicSubagents(nil, []string{dir})
	})

	var agents []subagentInfo
	_ = json.Unmarshal([]byte(strings.TrimSpace(out)), &agents)

	if len(agents) != 0 {
		t.Errorf("expected 0 agents, got %d", len(agents))
	}
}

func TestForensicSubagents_NoDir(t *testing.T) {
	if os.Getenv("TEST_FORENSIC_SUBAGENTS_NODIR") == "1" {
		runForensicSubagents(nil, []string{"/nonexistent"})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestForensicSubagents_NoDir")
	cmd.Env = append(os.Environ(), "TEST_FORENSIC_SUBAGENTS_NODIR=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for missing dir")
	}
	if !strings.Contains(string(output), "NOT_FOUND") {
		t.Errorf("expected NOT_FOUND, got: %s", string(output))
	}
}

// ── runForensicSearch ────────────────────────────────────────────────

func TestForensicSearch_WithKeyword(t *testing.T) {
	// Create a temp history.jsonl
	dir := t.TempDir()
	histPath := filepath.Join(dir, "history.jsonl")

	entries := []historyEntry{
		{Display: "/forge:eval-prd started", Timestamp: 1700000000000, Project: "/forge", SessionID: "sess-1"},
		{Display: "some other command", Timestamp: 1700000001000, Project: "/forge", SessionID: "sess-2"},
		{Display: "eval-prd done", Timestamp: 1700000002000, Project: "/other", SessionID: "sess-3"},
	}
	f, _ := os.Create(histPath)
	for _, e := range entries {
		line, _ := json.Marshal(e)
		_, _ = f.Write(line)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Close()

	// Override home to use our temp dir
	t.Setenv("HOME", dir)

	forensicKeyword = "eval-prd"
	forensicSession = ""
	forensicSkill = ""
	forensicLast = 10

	out := captureStdout(func() {
		// Manually build the search with our temp history
		searchWithHistPath("forge", histPath)
	})

	var results []sessionSummary
	_ = json.Unmarshal([]byte(strings.TrimSpace(out)), &results)

	// Should find sess-1 (forge + eval-prd) but not sess-2 (forge, no keyword) or sess-3 (other project)
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d: %v", len(results), results)
	}
	if len(results) > 0 && results[0].SessionID != "sess-1" {
		t.Errorf("SessionID = %q, want sess-1", results[0].SessionID)
	}
}

// searchWithHistPath is a test helper that runs search logic with a custom history path.
func searchWithHistPath(projectPath, histPath string) {
	f, err := os.Open(histPath)
	if err != nil {
		Exit(NewAIError(ErrNotFound, "Cannot open history.jsonl", err.Error(), "", ""))
	}
	defer func() { _ = f.Close() }()

	sessions := map[string]*sessionSummary{}
	scanner := newBufScanner(f)

	for scanner.scan() {
		var entry historyEntry
		if err := json.Unmarshal(scanner.bytes(), &entry); err != nil {
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

		ss, exists := sessions[entry.SessionID]
		if !exists {
			ss = &sessionSummary{SessionID: entry.SessionID, Project: entry.Project}
			sessions[entry.SessionID] = ss
		}
		ss.MsgCount++
		if entry.Timestamp > 0 {
			ss.DateTime = "2023-11-14 22:13" // fixed for test determinism
		}
		if ss.FirstMsg == "" {
			ss.FirstMsg = truncate(entry.Display, 80)
		}
	}

	sorted := make([]*sessionSummary, 0, len(sessions))
	for _, ss := range sessions {
		sorted = append(sorted, ss)
	}
	if forensicLast < len(sorted) {
		sorted = sorted[:forensicLast]
	}
	out, _ := json.MarshalIndent(sorted, "", "  ")
	fmt.Println(string(out))
}

// bufScanner wraps bufio.Scanner for testability.
type bufScanner struct {
	*bufio.Scanner
}

func newBufScanner(f *os.File) *bufScanner {
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)
	return &bufScanner{s}
}

func (b *bufScanner) scan() bool    { return b.Scan() }
func (b *bufScanner) bytes() []byte { return b.Bytes() }

// ── jsonlMessage UnmarshalJSON ──────────────────────────────────────

func TestJsonlMessage_Unmarshal_StringContent(t *testing.T) {
	data := `{"id":"m1","role":"user","content":"hello world","model":"glm-5"}`
	var msg jsonlMessage
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		t.Fatal(err)
	}
	if msg.textContent() != "hello world" {
		t.Errorf("textContent() = %q, want %q", msg.textContent(), "hello world")
	}
}

func TestJsonlMessage_Unmarshal_ArrayContent(t *testing.T) {
	data := `{"id":"m1","role":"assistant","content":[{"type":"thinking","thinking":"hmm"},{"type":"text","text":"hello"}],"model":"glm-5"}`
	var msg jsonlMessage
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		t.Fatal(err)
	}
	if msg.textContent() != "hello" {
		t.Errorf("textContent() = %q, want %q", msg.textContent(), "hello")
	}
}

func TestJsonlMessage_Unmarshal_EmptyContent(t *testing.T) {
	data := `{"id":"m1","role":"user","content":"","model":"glm-5"}`
	var msg jsonlMessage
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		t.Fatal(err)
	}
	if msg.textContent() != "" {
		t.Errorf("textContent() = %q, want empty", msg.textContent())
	}
}

// ── Testdata-based integration tests ────────────────────────────────
// These tests use sampled session data from testdata/forensic/.
// The JSONL files are sampled from real Claude Code sessions, preserving
// all entry types but reduced to ~50-70 lines each.

var testdataDir = filepath.Join("testdata", "forensic")

// extractTestdata is a helper to extract evidence from a testdata JSONL.
func extractTestdata(t *testing.T, filename string) extractResult {
	t.Helper()
	jsonlPath := filepath.Join(testdataDir, filename)

	outDir := filepath.Join(t.TempDir(), "evidence")
	forensicOutDir = outDir
	forensicSlug = ""
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})
	data, err := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	if err != nil {
		t.Fatalf("evidence file not created: %v", err)
	}
	var result extractResult
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	return result
}

func TestGolden_SearchHistory(t *testing.T) {
	histPath := filepath.Join(testdataDir, "history.jsonl")

	forensicKeyword = ""
	forensicSession = ""
	forensicSkill = ""
	forensicLast = 10

	out := captureStdout(func() {
		searchWithHistPath("coding-harness/forge", histPath)
	})

	var results []sessionSummary
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected at least 1 session for forge project")
	}

	for _, r := range results {
		if r.SessionID == "" {
			t.Error("SessionID should not be empty")
		}
		if r.DateTime == "" {
			t.Error("DateTime should not be empty")
		}
		if r.Project == "" {
			t.Error("Project should not be empty")
		}
		t.Logf("session %s: %s msgs=%d first=%q", truncate(r.SessionID, 12), r.DateTime, r.MsgCount, truncate(r.FirstMsg, 40))
	}
}

func TestGolden_SearchByKeyword(t *testing.T) {
	histPath := filepath.Join(testdataDir, "history.jsonl")

	forensicKeyword = "eval-prd"
	forensicSession = ""
	forensicSkill = ""
	forensicLast = 10

	out := captureStdout(func() {
		searchWithHistPath("coding-harness/forge", histPath)
	})

	var results []sessionSummary
	_ = json.Unmarshal([]byte(strings.TrimSpace(out)), &results)

	if len(results) == 0 {
		t.Fatal("expected sessions matching 'eval-prd'")
	}
	for _, r := range results {
		if !strings.Contains(strings.ToLower(r.FirstMsg), "eval-prd") {
			t.Errorf("result firstMsg should contain 'eval-prd': %q", r.FirstMsg)
		}
	}
}

func TestGolden_ExtractFixBugSession(t *testing.T) {
	jsonlPath := filepath.Join(testdataDir, "fix-bug.jsonl")

	outDir := filepath.Join(t.TempDir(), "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	data, err := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	if err != nil {
		t.Fatalf("evidence file not created: %v", err)
	}
	var result extractResult
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result.Lines == 0 {
		t.Error("Lines should be > 0")
	}
	if result.Model == "" {
		t.Error("Model should not be empty")
	}
	if result.GitBranch == "" {
		t.Error("GitBranch should not be empty")
	}
	if result.Summary.TotalThinking == 0 {
		t.Error("should have thinking blocks")
	}
	if result.Summary.TotalToolCalls == 0 {
		t.Error("should have tool calls")
	}
	if result.Summary.TotalUserMsgs == 0 {
		t.Error("should have user messages")
	}

	// Consistency checks
	if len(result.Thinking) != result.Summary.TotalThinking {
		t.Errorf("thinking count mismatch: len=%d summary=%d",
			len(result.Thinking), result.Summary.TotalThinking)
	}
	if len(result.ToolCalls) != result.Summary.TotalToolCalls {
		t.Errorf("tool call count mismatch: len=%d summary=%d",
			len(result.ToolCalls), result.Summary.TotalToolCalls)
	}
	if len(result.ToolResults) != result.Summary.TotalToolResults {
		t.Errorf("tool results count mismatch: len=%d summary=%d",
			len(result.ToolResults), result.Summary.TotalToolResults)
	}

	// Tool breakdown sum must equal total
	total := 0
	for _, count := range result.Summary.ToolBreakdown {
		total += count
	}
	if total != result.Summary.TotalToolCalls {
		t.Errorf("tool breakdown total=%d != summary=%d", total, result.Summary.TotalToolCalls)
	}

	// Common tools
	for _, tool := range []string{"Read", "Bash"} {
		if result.Summary.ToolBreakdown[tool] == 0 {
			t.Errorf("session should use %s", tool)
		}
	}

	// Thinking entries: content, line numbers
	for i, th := range result.Thinking {
		if th.Thinking == "" {
			t.Errorf("thinking[%d] empty", i)
		}
		if th.Line == 0 {
			t.Errorf("thinking[%d] zero line", i)
		}
	}
	// Tool calls: valid tool names and line numbers
	for i, tc := range result.ToolCalls {
		if tc.Tool == "" {
			t.Errorf("toolCalls[%d] empty tool", i)
		}
		if tc.Line == 0 {
			t.Errorf("toolCalls[%d] zero line", i)
		}
	}

	// Hooks and files edited
	for i, h := range result.Hooks {
		if h.HookName == "" {
			t.Errorf("hooks[%d] empty HookName", i)
		}
	}
	for i, f := range result.FilesEdited {
		if f == "" {
			t.Errorf("filesEdited[%d] empty filename", i)
		}
	}

	// Source JSONL copied to output dir
	copiedPath := filepath.Join(outDir, "fix-bug.jsonl")
	if _, err := os.Stat(copiedPath); err != nil {
		t.Error("source JSONL should be copied to output directory")
	}

	t.Logf("fix-bug: lines=%d thinking=%d toolCalls=%d userMsgs=%d hooks=%d filesEdited=%d tools=%v",
		result.Lines, result.Summary.TotalThinking, result.Summary.TotalToolCalls,
		result.Summary.TotalUserMsgs, len(result.Hooks), len(result.FilesEdited),
		result.Summary.ToolBreakdown)
}

func TestGolden_ExtractSubagentEvalSession(t *testing.T) {
	result := extractTestdata(t, "subagent-eval.jsonl")

	if result.Lines == 0 {
		t.Error("Lines should be > 0")
	}
	if result.Model == "" {
		t.Error("Model should not be empty")
	}
	if result.Summary.TotalToolCalls == 0 {
		t.Error("should have tool calls")
	}
	if result.Summary.TotalToolResults == 0 {
		t.Error("should have tool results")
	}
	if result.Summary.ToolBreakdown["Agent"] == 0 {
		t.Error("subagent eval session should use Agent tool")
	}
	if len(result.Hooks) == 0 {
		t.Error("expected at least one hook event")
	}
	if len(result.FilesEdited) == 0 {
		t.Error("expected edited files from attachment data")
	}

	t.Logf("subagent-eval: lines=%d thinking=%d toolCalls=%d agents=%d hooks=%d filesEdited=%d",
		result.Lines, result.Summary.TotalThinking, result.Summary.TotalToolCalls,
		result.Summary.ToolBreakdown["Agent"], len(result.Hooks), len(result.FilesEdited))
}

func TestGolden_ExtractQuickModeSession(t *testing.T) {
	result := extractTestdata(t, "quick-mode.jsonl")

	if result.Lines == 0 {
		t.Error("Lines should be > 0")
	}
	if len(result.Hooks) == 0 {
		t.Error("quick-mode session should have hooks")
	}
	for i, h := range result.Hooks {
		if h.HookEvent == "" {
			t.Errorf("hooks[%d] empty HookEvent", i)
		}
	}
	if result.Summary.TotalToolResults == 0 {
		t.Error("should have tool results")
	}

	reasons := map[string]int{}
	for _, tc := range result.ToolCalls {
		reasons[tc.StopReason]++
	}
	if len(reasons) == 0 {
		t.Error("expected stop_reason values on tool calls")
	}

	t.Logf("quick-mode: lines=%d hooks=%d toolResults=%d stopReasons=%v",
		result.Lines, len(result.Hooks),
		result.Summary.TotalToolResults, reasons)
}

func TestGolden_ExtractLessonFixSession(t *testing.T) {
	result := extractTestdata(t, "lesson-fix.jsonl")

	if result.Lines == 0 {
		t.Error("Lines should be > 0")
	}
	if result.Summary.ToolBreakdown["Edit"] == 0 {
		t.Error("lesson-fix session should have Edit calls")
	}

	// Thinking chain integrity
	for i := 1; i < len(result.Thinking); i++ {
		if result.Thinking[i].Line <= result.Thinking[i-1].Line {
			t.Errorf("thinking[%d].Line=%d <= thinking[%d].Line=%d",
				i, result.Thinking[i].Line, i-1, result.Thinking[i-1].Line)
		}
	}

	for i, tc := range result.ToolCalls {
		if tc.Tool == "" {
			t.Errorf("toolCalls[%d] empty tool name", i)
		}
	}

	t.Logf("lesson-fix: lines=%d thinking=%d edits=%d",
		result.Lines, result.Summary.TotalThinking, result.Summary.ToolBreakdown["Edit"])
}

func TestGolden_Subagents(t *testing.T) {
	sessionDir := filepath.Join(testdataDir, "subagents-session")

	out := captureStdout(func() {
		runForensicSubagents(nil, []string{sessionDir})
	})

	var agents []subagentInfo
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &agents); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}

	types := map[string]string{}
	for _, a := range agents {
		types[a.AgentID] = a.AgentType
	}
	if types["a1b2c3d4"] != "Explore" {
		t.Errorf("agent a1b2c3d4 type = %q, want Explore", types["a1b2c3d4"])
	}
	if types["e5f6a7b8"] != "forge:task-executor" {
		t.Errorf("agent e5f6a7b8 type = %q, want forge:task-executor", types["e5f6a7b8"])
	}
}

func TestGolden_ExtractSubagentTranscript(t *testing.T) {
	transcriptPath := filepath.Join(testdataDir, "subagents-session", "subagents", "agent-a1b2c3d4.jsonl")

	outDir := filepath.Join(t.TempDir(), "evidence")
	forensicOutDir = outDir
	captureStdout(func() {
		runForensicExtract(nil, []string{transcriptPath})
	})

	data, err := os.ReadFile(filepath.Join(outDir, "evidence.json"))
	if err != nil {
		t.Fatalf("evidence file not created: %v", err)
	}
	var result extractResult
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result.Lines != 2 {
		t.Errorf("Lines = %d, want 2", result.Lines)
	}
	if result.Summary.TotalThinking != 1 {
		t.Errorf("TotalThinking = %d, want 1", result.Summary.TotalThinking)
	}
	if result.Summary.TotalToolCalls != 1 {
		t.Errorf("TotalToolCalls = %d, want 1", result.Summary.TotalToolCalls)
	}
	if result.Summary.ToolBreakdown["Grep"] != 1 {
		t.Errorf("expected Grep tool, got %v", result.Summary.ToolBreakdown)
	}

	// Source JSONL copied
	copied := filepath.Join(outDir, "agent-a1b2c3d4.jsonl")
	if _, err := os.Stat(copied); err != nil {
		t.Error("source transcript should be copied to output directory")
	}
}

func TestGolden_ThinkingTruncation(t *testing.T) {
	result := extractTestdata(t, "fix-bug.jsonl")

	for i, th := range result.Thinking {
		if len(th.Thinking) > 600 {
			t.Errorf("thinking[%d] not truncated: len=%d", i, len(th.Thinking))
		}
	}
}

func TestGolden_ToolInputTruncation(t *testing.T) {
	result := extractTestdata(t, "fix-bug.jsonl")

	for i, tc := range result.ToolCalls {
		if len(tc.Input) > 600 {
			t.Errorf("toolCalls[%d] input not truncated: len=%d", i, len(tc.Input))
		}
	}
}

func TestGolden_UserMessagesNotEmpty(t *testing.T) {
	result := extractTestdata(t, "fix-bug.jsonl")

	for i, msg := range result.UserMsgs {
		if strings.TrimSpace(msg.Content) == "" {
			t.Errorf("empty userMsg at index %d line %d", i, msg.Line)
		}
	}
}

func TestGolden_CrossSessionConsistency(t *testing.T) {
	fixtures := []string{"fix-bug.jsonl", "subagent-eval.jsonl", "quick-mode.jsonl", "lesson-fix.jsonl"}

	for _, name := range fixtures {
		result := extractTestdata(t, name)

		total := 0
		for _, count := range result.Summary.ToolBreakdown {
			total += count
		}
		if total != result.Summary.TotalToolCalls {
			t.Errorf("%s: tool breakdown total=%d != summary=%d",
				name, total, result.Summary.TotalToolCalls)
		}
		if len(result.Thinking) != result.Summary.TotalThinking {
			t.Errorf("%s: thinking len=%d != summary=%d",
				name, len(result.Thinking), result.Summary.TotalThinking)
		}
		if len(result.ToolCalls) != result.Summary.TotalToolCalls {
			t.Errorf("%s: toolCalls len=%d != summary=%d",
				name, len(result.ToolCalls), result.Summary.TotalToolCalls)
		}
		if len(result.ToolResults) != result.Summary.TotalToolResults {
			t.Errorf("%s: toolResults len=%d != summary=%d",
				name, len(result.ToolResults), result.Summary.TotalToolResults)
		}
	}
}

func TestGolden_SlugFlag(t *testing.T) {
	// Create a JSONL file with a session-ID-like name
	dir := t.TempDir()
	sessionID := "abc12345-6789-def0-abcd-ef1234567890"
	jsonlPath := filepath.Join(dir, sessionID+".jsonl")
	_ = os.WriteFile(jsonlPath, []byte(`{"type":"user","message":{"role":"user","content":"hello"}}`+"\n"), 0644)

	forensicOutDir = ""
	forensicSlug = "my-investigation"

	out := captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	expectedDir := filepath.Join("docs", "forensics", "my-investigation", "evidence")
	if !strings.Contains(out, expectedDir) {
		t.Errorf("expected output path %q, got %q", expectedDir, out)
	}

	data, err := os.ReadFile(filepath.Join(expectedDir, "evidence.json"))
	if err != nil {
		t.Fatalf("evidence not written to slug dir: %v", err)
	}
	var result extractResult
	_ = json.Unmarshal(data, &result)
	if result.Lines != 1 {
		t.Errorf("Lines = %d, want 1", result.Lines)
	}

	// Cleanup
	_ = os.RemoveAll("docs/forensics/my-investigation")
}

func TestGolden_AutoDeriveSlug(t *testing.T) {
	// Without --slug or --out, auto-derive from session ID filename
	dir := t.TempDir()
	sessionID := "abc12345-6789-def0-abcd-ef1234567890"
	jsonlPath := filepath.Join(dir, sessionID+".jsonl")
	_ = os.WriteFile(jsonlPath, []byte(`{"type":"user","message":{"role":"user","content":"test"}}`+"\n"), 0644)

	forensicOutDir = ""
	forensicSlug = ""

	out := captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	expectedDir := filepath.Join("docs", "forensics", sessionID, "evidence")
	if !strings.Contains(out, expectedDir) {
		t.Errorf("expected auto-derived path %q, got %q", expectedDir, out)
	}

	data, err := os.ReadFile(filepath.Join(expectedDir, "evidence.json"))
	if err != nil {
		t.Fatalf("evidence not written: %v", err)
	}
	var result extractResult
	_ = json.Unmarshal(data, &result)
	if result.Lines != 1 {
		t.Errorf("Lines = %d, want 1", result.Lines)
	}

	// Cleanup
	_ = os.RemoveAll(filepath.Join("docs", "forensics", sessionID))
}

func TestGolden_ExplicitOutWithoutSlug(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "session.jsonl")
	_ = os.WriteFile(jsonlPath, []byte(`{"type":"user","message":{"role":"user","content":"x"}}`+"\n"), 0644)

	customOut := filepath.Join(dir, "custom-output")
	forensicOutDir = customOut
	forensicSlug = ""

	captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	data, err := os.ReadFile(filepath.Join(customOut, "evidence.json"))
	if err != nil {
		t.Fatalf("evidence not written to --out dir: %v", err)
	}
	var result extractResult
	_ = json.Unmarshal(data, &result)
	if result.Lines != 1 {
		t.Errorf("Lines = %d, want 1", result.Lines)
	}
}

func TestAddToAgg(t *testing.T) {
	entries := []toolAggEntry{}

	addToAgg(&entries, "hook-a")
	addToAgg(&entries, "hook-b")
	addToAgg(&entries, "hook-a")

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(entries), entries)
	}
	found := false
	for _, e := range entries {
		if e.Name == "hook-a" && e.Count == 2 {
			found = true
		}
	}
	if !found {
		t.Errorf("hook-a count should be 2, got %v", entries)
	}
}

func TestGolden_ExtendedMetrics(t *testing.T) {
	// fix-bug.jsonl: has hooks (hook_success), compact_file_reference, plan_mode, stop_reasons, timestamps
	result := extractTestdata(t, "fix-bug.jsonl")

	// Hook breakdown
	if len(result.Summary.HookBreakdown) == 0 {
		t.Error("HookBreakdown should be populated for fix-bug session")
	}
	t.Logf("hookBreakdown=%v hookFailures=%d", result.Summary.HookBreakdown, result.Summary.HookFailures)

	// Stop reasons
	if len(result.Summary.StopReasons) == 0 {
		t.Error("StopReasons should be populated")
	}
	for reason, count := range result.Summary.StopReasons {
		if reason == "" {
			t.Error("StopReason key should not be empty")
		}
		if count <= 0 {
			t.Errorf("StopReason[%s] count should be > 0, got %d", reason, count)
		}
		t.Logf("  stopReason: %s = %d", reason, count)
	}

	// Compact count
	if result.Summary.CompactCount <= 0 {
		t.Errorf("CompactCount should be > 0 for fix-bug session, got %d", result.Summary.CompactCount)
	}
	t.Logf("compactCount=%d", result.Summary.CompactCount)

	// Plan mode count
	if result.Summary.PlanModeCount <= 0 {
		t.Errorf("PlanModeCount should be > 0 for fix-bug session, got %d", result.Summary.PlanModeCount)
	}
	t.Logf("planModeCount=%d", result.Summary.PlanModeCount)

	// Duration
	if result.Summary.Duration == "" {
		t.Error("Duration should be computed from timestamps")
	}
	t.Logf("duration=%s", result.Summary.Duration)
}

func TestGolden_QuickModeMetrics(t *testing.T) {
	result := extractTestdata(t, "quick-mode.jsonl")

	// quick-mode has hooks, plan_mode, stop_reasons
	if len(result.Summary.HookBreakdown) == 0 {
		t.Error("HookBreakdown should be populated")
	}
	if result.Summary.PlanModeCount <= 0 {
		t.Errorf("PlanModeCount should be > 0, got %d", result.Summary.PlanModeCount)
	}
	if len(result.Summary.StopReasons) == 0 {
		t.Error("StopReasons should be populated")
	}

	// quick-mode has Agent calls => subagent count
	if result.Summary.SubagentCount <= 0 {
		t.Errorf("SubagentCount should be > 0, got %d", result.Summary.SubagentCount)
	}
	t.Logf("subagentCount=%d planMode=%d hooks=%v stopReasons=%v",
		result.Summary.SubagentCount, result.Summary.PlanModeCount,
		result.Summary.HookBreakdown, result.Summary.StopReasons)
}

func TestGolden_SubagentCountConsistency(t *testing.T) {
	// SubagentCount should equal sum of AgentsSpawned counts
	fixtures := []string{"subagent-eval.jsonl", "quick-mode.jsonl"}
	for _, name := range fixtures {
		result := extractTestdata(t, name)
		expected := 0
		for _, a := range result.Summary.AgentsSpawned {
			expected += a.Count
		}
		if result.Summary.SubagentCount != expected {
			t.Errorf("%s: SubagentCount=%d but sum of AgentsSpawned=%d",
				name, result.Summary.SubagentCount, expected)
		}
		t.Logf("%s: subagentCount=%d agents=%v", name, result.Summary.SubagentCount, result.Summary.AgentsSpawned)
	}
}

func TestGolden_SlugOverridesOut(t *testing.T) {
	dir := t.TempDir()
	sessionID := "abc12345-6789-def0-abcd-ef1234567890"
	jsonlPath := filepath.Join(dir, sessionID+".jsonl")
	_ = os.WriteFile(jsonlPath, []byte(`{"type":"user","message":{"role":"user","content":"x"}}`+"\n"), 0644)

	customOut := filepath.Join(dir, "custom-output")
	forensicOutDir = customOut
	forensicSlug = "slug-wins"

	out := captureStdout(func() {
		runForensicExtract(nil, []string{jsonlPath})
	})

	// --slug should override --out
	expectedDir := filepath.Join("docs", "forensics", "slug-wins", "evidence")
	if !strings.Contains(out, expectedDir) {
		t.Errorf("expected slug path %q, got %q", expectedDir, out)
	}

	if _, err := os.Stat(filepath.Join(expectedDir, "evidence.json")); err != nil {
		t.Fatalf("evidence not written to slug dir: %v", err)
	}

	// --out dir should NOT have evidence
	if _, err := os.Stat(filepath.Join(customOut, "evidence.json")); err == nil {
		t.Error("--out dir should not be used when --slug is set")
	}

	_ = os.RemoveAll("docs/forensics/slug-wins")
}
