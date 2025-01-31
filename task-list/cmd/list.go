/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of all of your tasks",
	Long: `Get a list of all the tasks that you need to do today.
	You will see the
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		tasks, error := task.ReadTasks()

		if error != nil {

			log.Fatal(error)
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
	listCmd.Flags().String("filter-priority", "", "Filter tasks by priority")
	listCmd.Flags().String("filter-complete", "", "Filter tasks by completed")
	listCmd.Flags().String("filter-incomplete", "", "Filter tasks by incompleted")

	listCmd.MarkFlagsMutuallyExclusive("filter-complete", "filter-incomplete")

	listCmd.Flags().String(
		"sort-date",
		"",
		fmt.Sprintf("Sort by  date %s", strings.Join(allowedDateSortValues, ",")),
	)

	listCmd.Flags().String(
		"sort-priority",
		"",
		fmt.Sprintf("Sort by  priority %s", strings.Join(allowedPriortySortValues, ",")),
	)

	listCmd.MarkFlagsMutuallyExclusive("sort-priority", "filter-priority")

}
