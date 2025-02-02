/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/mini-clis/task-list/task"
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

		filterPriority := cmd.Flag(FILTER_PRIORITY).Value.String()
		filterComplete := cmd.Flag(FILTER_COMPLETE).Value.String()
		filterIncomplete := cmd.Flag(FILTER_INCOMPLETE).Value.String()
		sortDate := cmd.Flag(SORT_DATE).Value.String()
		sortPriority := cmd.Flag(SORT_PRIORITY).Value.String()

		tasks, error := task.ReadTasks()

		if error != nil {

			log.Fatal(error)
		}

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
				aTime, err := time.Parse(time.UnixDate, a.CreatedAt())
				if err != nil {
					log.Fatal(err)
				}
				bTime, err := time.Parse(time.UnixDate, b.CreatedAt())
				if err != nil {
					log.Fatal(err)
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
				aTime, err := time.Parse(time.UnixDate, a.CreatedAt())
				if err != nil {
					log.Fatal(err)
				}
				bTime, err := time.Parse(time.UnixDate, b.CreatedAt())
				if err != nil {
					log.Fatal(err)
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

		stringifiedTasks, error := task.MarshallTasks(tasks)

		if error != nil {

			log.Fatal(error)
		}

		fmt.Println("Here is the list of tasks you have to do")
		fmt.Printf("%s", stringifiedTasks)

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
	listCmd.Flags().String(FILTER_COMPLETE, "", "Filter tasks by completed")
	listCmd.Flags().String(FILTER_INCOMPLETE, "", "Filter tasks by incompleted")

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
