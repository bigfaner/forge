package contract

import (
	"strings"
	"testing"
)

func TestRenderContract(t *testing.T) {
	t.Run("renders valid Contract to markdown", func(t *testing.T) {
		c := Contract{
			Journey: "task-lifecycle",
			Step:    2,
			Action:  "forge task claim",
			Outcomes: []Outcome{
				{
					Name:          "success",
					Preconditions: "feature exists with slug matching arg; at least one task available",
					Input:         "no positional args; no flags",
					Output:        "stdout contains claimed task identifier, exit code 0",
					State:         "tasks/task-id/status changed to in_progress; index.json updated",
					SideEffect:    "none",
				},
				{
					Name:          "no-tasks-available",
					Preconditions: "feature exists; no tasks available for claiming",
					Input:         "no positional args; no flags",
					Output:        "stderr contains no tasks available message, exit code 1",
					State:         "unchanged",
				},
			},
			Invariants: []string{
				"feature_slug consistent across all steps",
				"task_id stable once assigned",
			},
		}

		md := RenderContract(c)

		// Verify frontmatter
		if !strings.Contains(md, "journey: \"task-lifecycle\"") {
			t.Fatal("expected journey in frontmatter")
		}
		if !strings.Contains(md, "step: 2") {
			t.Fatal("expected step number in frontmatter")
		}

		// Verify Outcome blocks
		if !strings.Contains(md, `## Outcome "success"`) {
			t.Fatal("expected success Outcome heading")
		}
		if !strings.Contains(md, `## Outcome "no-tasks-available"`) {
			t.Fatal("expected no-tasks-available Outcome heading")
		}

		// Verify dimensions
		if !strings.Contains(md, "- Preconditions:") {
			t.Fatal("expected Preconditions dimension")
		}
		if !strings.Contains(md, "- Input:") {
			t.Fatal("expected Input dimension")
		}
		if !strings.Contains(md, "- Output:") {
			t.Fatal("expected Output dimension")
		}
		if !strings.Contains(md, "- State:") {
			t.Fatal("expected State dimension")
		}
		if !strings.Contains(md, "- Side-effect: none") {
			t.Fatal("expected Side-effect: none default")
		}

		// Verify Journey Invariants section
		if !strings.Contains(md, "## Journey Invariants") {
			t.Fatal("expected Journey Invariants heading")
		}
		if !strings.Contains(md, "feature_slug consistent across all steps") {
			t.Fatal("expected first invariant")
		}
		if !strings.Contains(md, "task_id stable once assigned") {
			t.Fatal("expected second invariant")
		}
	})

	t.Run("renders TUI await Outcome", func(t *testing.T) {
		c := Contract{
			Journey: "session-diagnostics",
			Step:    4,
			Action:  "open diagnosis panel",
			Outcomes: []Outcome{
				{
					Name:          "diagnosis-loaded",
					Preconditions: "session loaded, call tree visible, entry expanded",
					Input:         "key \"d\" await 3000ms",
					Output:        "view contains diagnosis summary panel",
					State:         "Model.diagnosis_panel field set to visible",
					IsAsyncTUI:    true,
					AwaitTimeout:  3000,
				},
				{
					Name:          "diagnosis-timeout",
					Preconditions: "async Cmd exceeds await duration of 3000ms",
					Input:         "key \"d\" await 3000ms",
					Output:        "error message containing timed-out Cmd name, fail-fast",
					State:         "unchanged from pre-Cmd state",
					IsAsyncTUI:    true,
					AwaitTimeout:  3000,
					TimedOutCmd:   "diagnosis-loader",
				},
			},
			Invariants: []string{
				"session_id consistent across all steps",
			},
		}

		md := RenderContract(c)

		if !strings.Contains(md, "await 3000ms") {
			t.Fatal("expected await timeout in Input dimension")
		}
		if !strings.Contains(md, "diagnosis-timeout") {
			t.Fatal("expected timeout Outcome name")
		}
	})

	t.Run("renders state-verification annotation", func(t *testing.T) {
		c := Contract{
			Journey: "task-lifecycle",
			Step:    1,
			Action:  "forge feature create",
			Outcomes: []Outcome{
				{
					Name:          "success",
					Preconditions: "no feature with this slug exists",
					Input:         "feature-slug as positional arg",
					Output:        "success confirmation",
					State:         "feature directory created",
				},
			},
			Invariants:       []string{"slug consistency"},
			StateVerifyLevel: "partial",
		}

		md := RenderContract(c)
		if !strings.Contains(md, "state-verification: partial") {
			t.Fatal("expected state-verification annotation")
		}
	})
}
