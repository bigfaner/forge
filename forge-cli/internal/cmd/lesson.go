package cmd

import (
	"fmt"
	"os"
	"strings"

	"forge-cli/pkg/lesson"
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
	Run:  runLesson,
}

func runLesson(_ *cobra.Command, args []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	if len(args) == 0 {
		runLessonList(projectRoot)
	} else {
		runLessonDetail(projectRoot, args[0])
	}
}

func runLessonList(projectRoot string) {
	lessons, err := lesson.Discover(projectRoot)
	if err != nil {
		Exit(newErrLessonDiscovery(err))
	}

	if len(lessons) == 0 {
		fmt.Fprintln(os.Stderr, "no lessons found")
		return
	}

	// Calculate dynamic name column width.
	slugWidth := calcSlugColWidth(mapLessonsToNameLens(lessons))

	PrintBlockStart()
	PrintField("LESSONS", fmt.Sprintf("%d found", len(lessons)))
	fmt.Println()

	// Table header
	fmt.Printf("  %-s %-12s %-15s %-12s\n", padRight("NAME", slugWidth), "CREATED", "CATEGORY", "TAGS")
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
			padRight(truncateSlug(l.Name, slugWidth), slugWidth),
			l.Created,
			l.Category,
			tags)
	}

	fmt.Println()
	PrintBlockEnd()
}

func runLessonDetail(projectRoot, name string) {
	l, err := lesson.FindByName(projectRoot, name)
	if err != nil {
		Exit(newErrLessonNotFound(name))
	}

	PrintBlockStart()
	PrintField("NAME", l.Name)
	PrintFieldIfNotEmpty("TITLE", l.Title)
	PrintField("CREATED", l.Created)
	PrintField("CATEGORY", l.Category)
	if len(l.Tags) > 0 {
		PrintField("TAGS", strings.Join(l.Tags, ", "))
	}
	PrintField("FILE", l.FilePath)
	PrintBlockEnd()
}

func newErrLessonDiscovery(err error) *AIError {
	return NewAIError(
		ErrNotFound,
		"Failed to discover lessons",
		err.Error(),
		"Ensure docs/lessons/ directory exists",
		"ls docs/lessons/",
	)
}

func newErrLessonNotFound(name string) *AIError {
	return NewAIError(
		ErrNotFound,
		fmt.Sprintf("Lesson not found: %s", name),
		"No lesson .md file found with this name",
		"Check the name is correct (without .md extension)",
		"ls docs/lessons/",
	)
}

// mapLessonsToNameLens extracts name lengths from lesson list.
func mapLessonsToNameLens(lessons []lesson.Lesson) []int {
	lens := make([]int, len(lessons))
	for i, l := range lessons {
		lens[i] = len(l.Name)
	}
	return lens
}
