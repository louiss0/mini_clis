/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mini-clis/task-list/error_log"
	"github.com/mini-clis/task-list/task"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const LATEST = "latest"

const EARLIEST = "earliest"

var allowedDateSortValues = []string{
	LATEST,
	EARLIEST,
}

const HIGHEST = "highest"

const LOWEST = "lowest"

var allowedPriortySortValues = []string{
	HIGHEST,
	LOWEST,
}

const FILTER_PRIORITY = "filter-priority"
const FILTER_COMPLETE = "filter-complete"
const FILTER_INCOMPLETE = "filter-incomplete"
const SORT_DATE = "sort-date"
const SORT_PRIORITY = "sort-priority"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of all of your tasks",
	Long: `Get a list of all the tasks that you need to do today.
	You will see the
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		filterPriority := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
			cmd.Flags().GetString(FILTER_PRIORITY),
		)

		filterComplete := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
			cmd.Flags().GetBool(FILTER_COMPLETE),
		)

		filterIncomplete := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
			cmd.Flags().GetBool(FILTER_INCOMPLETE),
		)

		sortDate := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
			cmd.Flags().GetString(SORT_DATE),
		)

		sortPriority := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
			cmd.Flags().GetString(SORT_PRIORITY),
		)

		tasks := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(task.ReadTasks())

		if sortPriority == HIGHEST {

			slices.SortFunc(tasks, func(a task.Task, b task.Task) int {
				return a.Priority.Order() - b.Priority.Order()
			})

		}

		if sortPriority == LOWEST {

			slices.SortFunc(tasks, func(a task.Task, b task.Task) int {

				return b.Priority.Order() - a.Priority.Order()
			})

		}

		if sortDate == LATEST {

			slices.SortFunc(tasks, func(a task.Task, b task.Task) int {

				aTime := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
					time.Parse(time.UnixDate, a.CreatedAt()),
				)

				bTime := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
					time.Parse(time.UnixDate, b.CreatedAt()),
				)

				if aTime.After(bTime) {
					return -1
				}
				if bTime.After(aTime) {
					return 1
				}
				return 0
			})

		}

		if sortDate == EARLIEST {

			slices.SortFunc(tasks, func(a task.Task, b task.Task) int {

				aTime := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
					time.Parse(time.UnixDate, a.CreatedAt()),
				)

				bTime := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
					time.Parse(time.UnixDate, b.CreatedAt()),
				)

				if aTime.After(bTime) {
					return 1
				}

				if bTime.After(aTime) {
					return -1
				}

				return 0
			})

		}

		if filterComplete {

			tasks = lo.Filter(tasks, func(item task.Task, index int) bool {

				return item.Complete
			})

		}

		if filterIncomplete {

			tasks = lo.Filter(tasks, func(item task.Task, index int) bool {

				return !item.Complete
			})

		}

		if filterPriority != "" {

			priority := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
				task.ParsePriority(filterPriority),
			)

			tasks = lo.Filter(tasks, func(item task.Task, index int) bool {

				return item.Priority == priority
			})

		}

		stringifiedTasks := error_log.ReturnValueIfErrorIsNotNilLogFatalIfError(
			task.MarshallTasks(tasks),
		)

		// fmt.Println("Here is the list of tasks you have to do")
		fmt.Fprintln(
			cmd.OutOrStdout(),
			stringifiedTasks,
		)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().String(FILTER_PRIORITY, "", "Filter tasks by priority")
	listCmd.Flags().Bool(FILTER_COMPLETE, false, "Filter tasks by completed")
	listCmd.Flags().Bool(FILTER_INCOMPLETE, false, "Filter tasks by incompleted")

	listCmd.MarkFlagsMutuallyExclusive(FILTER_COMPLETE, FILTER_INCOMPLETE)

	listCmd.Flags().String(
		SORT_DATE,
		"",
		fmt.Sprintf("Sort by  date %s", strings.Join(allowedDateSortValues, ",")),
	)

	listCmd.Flags().String(
		SORT_PRIORITY,
		"",
		fmt.Sprintf("Sort by  priority %s", strings.Join(allowedPriortySortValues, ",")),
	)

	listCmd.MarkFlagsMutuallyExclusive(SORT_PRIORITY, FILTER_PRIORITY)

}
