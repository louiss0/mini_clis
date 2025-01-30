/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mini-clis/task-list/task"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to the task list",
	Long: `Adds a task name to the task list.
 A task can also have a description by using the --description flag.
 The time the task was created is automatcally stored.
 `,
	Run: func(cmd *cobra.Command, args []string) {
		tasks := lo.Map(args, func(item string, index int) string {

			return task.NewTask(item, "").ToJSON()

		})

		fmt.Printf("%#v\n", tasks)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
