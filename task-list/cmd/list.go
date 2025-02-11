/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mini-clis/task-list/custom_errors"
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

var allowedPrioritySortValues = []string{
	HIGHEST,
	LOWEST,
}

const FILTER_PRIORITY = "filter-priority"
const FILTER_COMPLETE = "filter-complete"
const FILTER_INCOMPLETE = "filter-incomplete"
const SORT_DATE = "sort-date"
const SORT_PRIORITY = "sort-priority"

// listCmd represents the list command
func CreateListCommand() *cobra.Command {

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of all of your tasks",
		Long: `Get a list of all the tasks that you need to do today.
			You will see the
			`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			filterPriority, filterPriorityErr := cmd.Flags().GetString(FILTER_PRIORITY)

			filterComplete, filterCompleteError := cmd.Flags().GetBool(FILTER_COMPLETE)

			filterIncomplete, filterIncompleteErr := cmd.Flags().GetBool(FILTER_INCOMPLETE)

			sortDate, sortDateErr := cmd.Flags().GetString(SORT_DATE)

			sortPriority, sortPriorityErr := cmd.Flags().GetString(SORT_PRIORITY)

			flagError := errors.Join(
				filterPriorityErr,
				filterCompleteError,
				filterIncompleteErr,
				sortDateErr,
				sortPriorityErr,
			)

			if flagError != nil {
				return flagError
			}

			tasks, tasksErr := task.ReadTasks()

			if tasksErr != nil {
				return tasksErr
			}

			sortPriorityIsEmptyAndItHasIncorrectValues :=
				sortPriority != "" && !lo.Contains(allowedPrioritySortValues, sortPriority)

			if sortPriorityIsEmptyAndItHasIncorrectValues {

				return custom_errors.CreateInvalidFlagErrorWithMessage(
					fmt.Sprintf(
						"The only allowed value for sort-priority are %s",
						strings.Join(allowedPrioritySortValues, ","),
					),
				)

			}

			if sortPriority == HIGHEST {
				slices.SortFunc(tasks, func(a task.Task, b task.Task) int {
					return b.Priority.Order() - a.Priority.Order()
				})
			}

			if sortPriority == LOWEST {
				slices.SortFunc(tasks, func(a task.Task, b task.Task) int {
					return a.Priority.Order() - b.Priority.Order()
				})
			}

			sortDateIsNotEmptyAndItHasInCorrrectValues :=
				sortDate != "" && !lo.Contains(allowedDateSortValues, sortDate)

			if sortDateIsNotEmptyAndItHasInCorrrectValues {

				return custom_errors.CreateInvalidFlagErrorWithMessage(
					fmt.Sprintf(
						"The only allowed value for sort-priority are %s",
						strings.Join(allowedDateSortValues, ","),
					),
				)

			}

			if sortDate == LATEST {
				slices.SortFunc(tasks, func(a task.Task, b task.Task) int {
					aTime, aTimeErr := time.Parse(time.DateOnly, a.CreatedAt())
					if aTimeErr != nil {
						return 0
					}

					bTime, bTimeErr := time.Parse(time.DateOnly, b.CreatedAt())
					if bTimeErr != nil {
						return 0
					}

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
					aTime, aTimeErr := time.Parse(time.DateOnly, a.CreatedAt())
					if aTimeErr != nil {
						return 0
					}

					bTime, bTimeErr := time.Parse(time.DateOnly, b.CreatedAt())
					if bTimeErr != nil {
						return 0
					}

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
				priority, priorityErr := task.ParsePriority(filterPriority)

				if priorityErr != nil {
					return priorityErr
				}

				tasks = lo.Filter(tasks, func(item task.Task, index int) bool {
					return item.Priority == priority
				})
			}

			stringifiedTasks, stringifiedTasksErr := task.MarshallTasks(tasks)

			if stringifiedTasksErr != nil {
				return stringifiedTasksErr
			}

			fmt.Println("Here is the list of tasks you have to do")
			fmt.Fprintln(
				cmd.OutOrStdout(),
				stringifiedTasks,
			)

			return nil

		},
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	listCmd.Flags().String(FILTER_PRIORITY, "", "Filter tasks by priority")
	listCmd.RegisterFlagCompletionFunc(
		FILTER_PRIORITY,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			return task.AllowedProrities, cobra.ShellCompDirectiveDefault
		},
	)
	listCmd.Flags().Bool(FILTER_COMPLETE, false, "Filter tasks by completed")
	listCmd.Flags().Bool(FILTER_INCOMPLETE, false, "Filter tasks by incompleted")

	listCmd.MarkFlagsMutuallyExclusive(FILTER_COMPLETE, FILTER_INCOMPLETE)

	listCmd.Flags().String(
		SORT_DATE,
		"",
		fmt.Sprintf("Sort by  date %s", strings.Join(allowedDateSortValues, ",")),
	)

	listCmd.RegisterFlagCompletionFunc(
		SORT_DATE,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			return allowedDateSortValues, cobra.ShellCompDirectiveDefault
		},
	)
	listCmd.Flags().String(
		SORT_PRIORITY,
		"",
		fmt.Sprintf("Sort by  priority %s", strings.Join(allowedPrioritySortValues, ",")),
	)
	listCmd.RegisterFlagCompletionFunc(
		SORT_PRIORITY,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			return allowedPrioritySortValues, cobra.ShellCompDirectiveDefault
		},
	)
	listCmd.MarkFlagsMutuallyExclusive(SORT_PRIORITY, FILTER_PRIORITY)

	return listCmd

}

func init() {

	rootCmd.AddCommand(CreateListCommand())

}
