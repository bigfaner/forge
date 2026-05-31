package cmd

import (
	"fmt"
	"os"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/infocmd"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var lessonCmd = &cobra.Command{
	Use:   "lesson [name]",
	Short: "List or show lesson details",
	Long: `List all lessons or show details for a specific lesson.

Without arguments: lists all lessons in table format.
With a name argument: shows metadata and file path for that lesson.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLesson,
}

func runLesson(_ *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	if len(args) == 0 {
		return runLessonList(projectRoot)
	}
	return runLessonDetail(projectRoot, args[0])
}

func runLessonList(projectRoot string) error {
	lessons, err := infocmd.DiscoverLessons(projectRoot)
	if err != nil {
		return newErrLessonDiscovery(err)
	}

	if len(lessons) == 0 {
		fmt.Fprintln(os.Stderr, "no lessons found")
		return nil
	}

	// Calculate dynamic name column width.
	slugWidth := base.CalcSlugColWidth(mapLessonsToNameLens(lessons))

	base.PrintBlockStart()
	base.PrintField("LESSONS", fmt.Sprintf("%d found", len(lessons)))
	fmt.Println()

	// Table header
	fmt.Printf("  %-s %-12s %-15s %-12s\n", base.PadRight("NAME", slugWidth), "CREATED", "CATEGORY", "TAGS")
	fmt.Printf("  %-s %-12s %-15s %-12s\n",
		strings.Repeat("-", slugWidth),
		strings.Repeat("-", 10),
		strings.Repeat("-", 10),
		strings.Repeat("-", 10))

	for _, l := range lessons {
		tags := "-"
		if len(l.Tags) > 0 {
			tags = strings.Join(l.Tags, ", ")
		}
		fmt.Printf("  %-s %-12s %-15s %-12s\n",
			base.PadRight(base.TruncateSlug(l.Name, slugWidth), slugWidth),
			l.Created,
			l.Category,
			tags)
	}

	fmt.Println()
	base.PrintBlockEnd()
	return nil
}

func runLessonDetail(projectRoot, name string) error {
	l, err := infocmd.FindLessonByName(projectRoot, name)
	if err != nil {
		return newErrLessonNotFound(name)
	}

	base.PrintBlockStart()
	base.PrintField("NAME", l.Name)
	base.PrintFieldIfNotEmpty("TITLE", l.Title)
	base.PrintField("CREATED", l.Created)
	base.PrintField("CATEGORY", l.Category)
	if len(l.Tags) > 0 {
		base.PrintField("TAGS", strings.Join(l.Tags, ", "))
	}
	base.PrintField("FILE", l.FilePath)
	base.PrintBlockEnd()
	return nil
}

func newErrLessonDiscovery(err error) *base.AIError {
	return base.NewAIError(
		base.ErrNotFound,
		"Failed to discover lessons",
		err.Error(),
		"Ensure docs/lessons/ directory exists",
		"ls docs/lessons/",
	)
}

func newErrLessonNotFound(name string) *base.AIError {
	return base.NewAIError(
		base.ErrNotFound,
		fmt.Sprintf("Lesson not found: %s", name),
		"No lesson .md file found with this name",
		"Check the name is correct (without .md extension)",
		"ls docs/lessons/",
	)
}

// mapLessonsToNameLens extracts name lengths from lesson list.
func mapLessonsToNameLens(lessons []infocmd.Lesson) []int {
	lens := make([]int, len(lessons))
	for i, l := range lessons {
		lens[i] = len(l.Name)
	}
	return lens
}
