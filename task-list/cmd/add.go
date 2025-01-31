/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/mini-clis/task-list/task"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Args:  cobra.ExactArgs(1),
	Short: "Adds a task to the task list",
	Long: `Adds a task name to the task list.
 A task can also have a description by using the --description flag.
 The time the task was created is automatcally stored.
 `,
	Run: func(cmd *cobra.Command, args []string) {

		description := cmd.Flag("description").Value.String()

		priority := cmd.Flag("priority").Value.String()

		title := args[0]

		tasks, error := task.ReadTasks()

		if error != nil {

			log.Fatal(error)

		}

		newTask := task.NewTask(title, description)

		if priority != "" {

			parsedPriority, error := task.ParsePriority(priority)

			if error != nil {

				log.Fatal(error)

			}
			newTask.Priority = parsedPriority
		}

		tasks = append(tasks, newTask)

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

	addCmd.Flags().StringP(
		"priority",
		"p",
		"",
		"Set a task to either a high, low or medium priority default is low",
	)
}
