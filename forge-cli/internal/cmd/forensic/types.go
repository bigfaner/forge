package forensic

import (
	"encoding/json"
)

// -- data types --

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
	Type      string `json:"type"`
	Thinking  string `json:"thinking,omitempty"`
	Text      string `json:"text,omitempty"`
	Name      string `json:"name,omitempty"`
	ID        string `json:"id,omitempty"`
	ToolUseID string `json:"tool_use_id,omitempty"`
	Input     any    `json:"input,omitempty"`
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
		ID         string     `json:"id"`
		Role       string     `json:"role"`
		Model      string     `json:"model"`
		StopReason string     `json:"stop_reason"`
		Usage      *usageInfo `json:"usage,omitempty"`
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

type snapshotData struct {
	Timestamp string `json:"timestamp,omitempty"`
}

type jsonlEntry struct {
	Type          string         `json:"type"`
	Message       jsonlMessage   `json:"message"`
	Content       string         `json:"content"`
	GitBranch     string         `json:"gitBranch"`
	Attachment    attachment     `json:"attachment"`
	ToolUseResult *toolUseResult `json:"toolUseResult"`
	Timestamp     string         `json:"timestamp,omitempty"`
	SessionID     string         `json:"sessionId,omitempty"`
	Snapshot      snapshotData   `json:"snapshot"`
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

// timingEntry records a single tool call's execution duration.
type timingEntry struct {
	Tool    string  `json:"tool"`
	Line    int     `json:"line"`
	Seconds float64 `json:"seconds"`
	Detail  string  `json:"detail,omitempty"`
}

// timingAgg aggregates timing statistics per tool type.
type timingAgg struct {
	Tool    string  `json:"tool"`
	Count   int     `json:"count"`
	Total   float64 `json:"total"`
	Average float64 `json:"average"`
	Max     float64 `json:"max"`
}

type thinkingTurn struct {
	Line       int     `json:"line"`
	Seconds    float64 `json:"seconds"`
	StopReason string  `json:"stopReason,omitempty"`
	Detail     string  `json:"detail,omitempty"`
}

type extractSummary struct {
	TotalThinking    int            `json:"totalThinking"`
	TotalToolCalls   int            `json:"totalToolCalls"`
	TotalToolResults int            `json:"totalToolResults"`
	TotalUserMsgs    int            `json:"totalUserMsgs"`
	ToolBreakdown    map[string]int `json:"toolBreakdown"`

	// Tool-specific aggregations
	FilesRead     []string       `json:"filesRead"`
	FilesWritten  []string       `json:"filesWritten"`
	GrepPatterns  []string       `json:"grepPatterns"`
	AgentsSpawned []toolAggEntry `json:"agentsSpawned"`
	Commands      []string       `json:"commands"`

	// Hook statistics
	HookBreakdown []toolAggEntry `json:"hookBreakdown"`
	HookFailures  int            `json:"hookFailures"`

	// Session lifecycle
	CompactCount  int            `json:"compactCount"`
	PlanModeCount int            `json:"planModeCount"`
	StopReasons   map[string]int `json:"stopReasons"`

	// Skill invocations (from attachment.invoked_skills)
	SkillInvocations []toolAggEntry `json:"skillInvocations"`

	// Sub-sessions
	SubagentCount int `json:"subagentCount"`

	// Time range
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Duration  string `json:"duration"`

	// Timing statistics (tool execution time from tool_use -> tool_result)
	TopSlowest      []timingEntry  `json:"topSlowest"`
	TimingByTool    []timingAgg    `json:"timingByTool"`
	TotalToolMs     int64          `json:"totalToolMs"`
	ThinkingTurns   []thinkingTurn `json:"thinkingTurns"`
	TotalThinkingMs int64          `json:"totalThinkingMs"`
}

type extractResult struct {
	File        string            `json:"file"`
	Lines       int               `json:"lines"`
	Model       string            `json:"model,omitempty"`
	GitBranch   string            `json:"gitBranch,omitempty"`
	Thinking    []thinkingEntry   `json:"thinking"`
	ToolCalls   []toolCallEntry   `json:"toolCalls"`
	ToolResults []toolResultEntry `json:"toolResults"`
	UserMsgs    []userMsgEntry    `json:"userMsgs"`
	SkillsUsed  []string          `json:"skillsUsed"`
	Hooks       []hookEventEntry  `json:"hooks"`
	FilesEdited []string          `json:"filesEdited"`
	Summary     extractSummary    `json:"summary"`
}

type subagentInfo struct {
	AgentID    string `json:"agentId"`
	AgentType  string `json:"agentType"`
	Transcript string `json:"transcript"`
}
