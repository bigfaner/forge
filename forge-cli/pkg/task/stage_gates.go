package task

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// PhaseInfo holds detected phase information for stage-gate generation.
type PhaseInfo struct {
	Number  int
	TaskIDs []string // business task IDs in this phase (sorted)
}

// Qualifies returns true if the phase has enough business tasks (>=2) for gate generation.
func (p PhaseInfo) Qualifies() bool {
	return len(p.TaskIDs) >= 2
}

// DetectPhases scans task IDs and groups them into phases.
// Only IDs matching the pattern <digit>.<digit> are considered.
// IDs prefixed with T-test or T-quick, or ending in .summary/.gate, are excluded.
// Returns phases sorted by phase number, with task IDs sorted within each phase.
func DetectPhases(taskIDs []string) []PhaseInfo {
	phaseMap := make(map[int][]string)

	for _, id := range taskIDs {
		// Skip test/quick task IDs
		if strings.HasPrefix(id, "T-test") || strings.HasPrefix(id, "T-quick") {
			continue
		}
		// Skip gate/summary IDs
		if strings.HasSuffix(id, IDSuffixSummary) || strings.HasSuffix(id, IDSuffixGate) {
			continue
		}

		// Parse phase from ID: must be exactly <digit>.<digit>
		parts := strings.Split(id, ".")
		if len(parts) != 2 {
			continue
		}
		phaseNum, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		if _, err := strconv.Atoi(parts[1]); err != nil {
			continue
		}
		if phaseNum <= 0 {
			continue
		}

		phaseMap[phaseNum] = append(phaseMap[phaseNum], id)
	}

	phases := make([]PhaseInfo, 0, len(phaseMap))
	for num, ids := range phaseMap {
		sort.Strings(ids)
		phases = append(phases, PhaseInfo{Number: num, TaskIDs: ids})
	}
	sort.Slice(phases, func(i, j int) bool {
		return phases[i].Number < phases[j].Number
	})

	return phases
}

// GenerateStageGates detects phases from taskIDs and generates .summary/.gate files
// for each qualifying phase. Returns the count of files generated.
// Skips files that already exist (idempotent).
func GenerateStageGates(taskIDs []string, tasksDir string, featureSlug string) (int, error) {
	phases := DetectPhases(taskIDs)
	generated := 0

	for _, phase := range phases {
		if !phase.Qualifies() {
			continue
		}

		// Generate summary if missing
		summaryKey := fmt.Sprintf("%d.summary", phase.Number)
		summaryPath := filepath.Join(tasksDir, summaryKey+".md")
		if _, err := os.Stat(summaryPath); os.IsNotExist(err) {
			content, err := GenerateSummaryMD(phase, featureSlug)
			if err != nil {
				return generated, fmt.Errorf("generate %s: %w", summaryKey, err)
			}
			if err := os.WriteFile(summaryPath, content, 0644); err != nil {
				return generated, fmt.Errorf("write %s: %w", summaryKey, err)
			}
			generated++
		}

		// Generate gate if missing
		gateKey := fmt.Sprintf("%d.gate", phase.Number)
		gatePath := filepath.Join(tasksDir, gateKey+".md")
		if _, err := os.Stat(gatePath); os.IsNotExist(err) {
			content, err := GenerateGateMD(phase, featureSlug)
			if err != nil {
				return generated, fmt.Errorf("generate %s: %w", gateKey, err)
			}
			if err := os.WriteFile(gatePath, content, 0644); err != nil {
				return generated, fmt.Errorf("write %s: %w", gateKey, err)
			}
			generated++
		}
	}

	return generated, nil
}

// GenerateSummaryMD generates the .md content for a phase summary task.
// Content is generated programmatically (Go string literals), not from embedded templates.
func GenerateSummaryMD(phase PhaseInfo, featureSlug string) ([]byte, error) {
	var buf strings.Builder

	// Ensure deterministic dependency order
	deps := make([]string, len(phase.TaskIDs))
	copy(deps, phase.TaskIDs)
	sort.Strings(deps)

	buf.WriteString("---\n")
	fmt.Fprintf(&buf, "id: %q\n", fmt.Sprintf("%d.summary", phase.Number))
	fmt.Fprintf(&buf, "title: %q\n", fmt.Sprintf("Phase %d Summary", phase.Number))
	buf.WriteString("priority: \"P0\"\n")
	buf.WriteString("estimated_time: \"15min\"\n")
	fmt.Fprintf(&buf, "dependencies: %v\n", formatYAMLList(deps))
	buf.WriteString("type: \"doc.summary\"\n")
	buf.WriteString("mainSession: false\n")
	buf.WriteString("---\n\n")

	fmt.Fprintf(&buf, "# %d.summary: Phase %d Summary\n\n", phase.Number, phase.Number)
	buf.WriteString("## Description\n\n")
	fmt.Fprintf(&buf, "Generate a structured summary of all completed tasks in phase %d. ", phase.Number)
	buf.WriteString("This summary is read by subsequent phase tasks to maintain cross-phase consistency.\n\n")

	buf.WriteString("## Instructions\n\n")
	buf.WriteString("### Step 1: Read all phase task records\n\n")
	fmt.Fprintf(&buf, "Read each record file from `docs/features/%s/tasks/records/` ", featureSlug)
	fmt.Fprintf(&buf, "whose filename starts with `%d.` and does NOT contain `.summary` ", phase.Number)
	fmt.Fprintf(&buf, "(e.g., `%d.1-*.md`, `%d.2-*.md`). ", phase.Number, phase.Number)
	buf.WriteString("Exclude the phase summary's own prior record if one exists.\n\n")

	buf.WriteString("### Step 2: Extract structured data into the summary field\n\n")
	buf.WriteString("The `summary` field in `record.json` MUST follow this exact template. ")
	buf.WriteString("Copy it verbatim and fill in the values.\n\n")

	fmt.Fprintf(&buf, "## Reference Files\n\n")
	fmt.Fprintf(&buf, "- All phase task records: `docs/features/%s/tasks/records/%d.*.md`\n", featureSlug, phase.Number)
	fmt.Fprintf(&buf, "- Design reference: `docs/features/%s/design/tech-design.md`\n", featureSlug)

	buf.WriteString("\n## Acceptance Criteria\n\n")
	buf.WriteString("- [ ] All phase task records have been read\n")
	buf.WriteString("- [ ] Summary follows the exact 5-section template\n")
	buf.WriteString("- [ ] Types & Interfaces Changed table is populated\n")
	buf.WriteString("- [ ] Record created via `/submit-task`\n")

	buf.WriteString("\n## Hard Rules\n\n")
	buf.WriteString("- MUST NOT write new feature code — this is documentation only\n")

	return []byte(buf.String()), nil
}

// GenerateGateMD generates the .md content for a phase gate task.
// Content is generated programmatically (Go string literals), not from embedded templates.
func GenerateGateMD(phase PhaseInfo, featureSlug string) ([]byte, error) {
	var buf strings.Builder

	summaryID := fmt.Sprintf("%d.summary", phase.Number)

	buf.WriteString("---\n")
	fmt.Fprintf(&buf, "id: %q\n", fmt.Sprintf("%d.gate", phase.Number))
	fmt.Fprintf(&buf, "title: %q\n", fmt.Sprintf("Phase %d Gate", phase.Number))
	buf.WriteString("priority: \"P0\"\n")
	buf.WriteString("estimated_time: \"1h\"\n")
	fmt.Fprintf(&buf, "dependencies: %v\n", formatYAMLList([]string{summaryID}))
	buf.WriteString("breaking: true\n")
	buf.WriteString("type: \"gate\"\n")
	buf.WriteString("mainSession: false\n")
	buf.WriteString("---\n\n")

	fmt.Fprintf(&buf, "# %d.gate: Phase %d Gate\n\n", phase.Number, phase.Number)

	buf.WriteString("## Description\n\n")
	fmt.Fprintf(&buf, "Exit verification gate for phase %d. ", phase.Number)
	buf.WriteString("Confirms that all outputs are complete, internally consistent, and match the design specification before the next phase begins.\n\n")

	buf.WriteString("## Verification Checklist\n\n")
	buf.WriteString("1. [ ] All interfaces from this phase compile without errors\n")
	buf.WriteString("2. [ ] Data models match `design/tech-design.md` (skip if single-layer feature — mark N/A)\n")
	buf.WriteString("3. [ ] No type mismatches between adjacent layers (skip if single-layer feature — mark N/A)\n")
	buf.WriteString("4. [ ] Project builds successfully\n")
	buf.WriteString("5. [ ] All existing tests pass\n")
	buf.WriteString("6. [ ] No deviations from design spec (or deviations are documented as decisions)\n")
	buf.WriteString("7. [ ] All Integration Specs from `tech-design.md` have corresponding code changes\n")
	buf.WriteString("8. [ ] All integration test cases pass (if gen-test-cases already ran)\n\n")

	fmt.Fprintf(&buf, "## Reference Files\n\n")
	fmt.Fprintf(&buf, "- Design: `docs/features/%s/design/tech-design.md`\n", featureSlug)
	fmt.Fprintf(&buf, "- This phase's task records: `docs/features/%s/tasks/records/%d.*.md`\n", featureSlug, phase.Number)
	fmt.Fprintf(&buf, "- This phase's summary: `docs/features/%s/tasks/records/%d-summary.md`\n\n", featureSlug, phase.Number)

	buf.WriteString("## Acceptance Criteria\n\n")
	buf.WriteString("- [ ] All applicable verification checklist items pass\n")
	buf.WriteString("- [ ] Any deviations from design are documented as decisions in the record\n")
	buf.WriteString("- [ ] Record created via `/submit-task` with test evidence\n\n")

	buf.WriteString("## Hard Rules\n\n")
	buf.WriteString("- MUST NOT write new feature code — this is verification only\n")

	return []byte(buf.String()), nil
}
