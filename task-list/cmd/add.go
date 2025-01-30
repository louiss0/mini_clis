/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

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

		description := cmd.Flag("description").Value.String()

		tasks := lo.Map(args, func(item string, index int) task.Task {

			if description != "" {
				return task.NewTask(item, description)
			}

			return task.NewTask(item, "")

		})

		err := task.SaveTasks(tasks)

		if err != nil {

			log.Fatal(err)
		}

		fmt.Print("Added a task to the task list")

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
	addCmd.Flags().StringP(
		"description",
		"d",
		"",
		`The task description.
		What is the task about?
		What are the requirements?`,
	)
}
