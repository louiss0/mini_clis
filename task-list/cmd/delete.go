/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mini-clis/task-list/custom_errors"
	"github.com/mini-clis/task-list/task"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
func CreateDeleteCommand() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a task based on an id",
		Long: `You can delete a task based on an Id.
			The flags in this command allow you to pass in
			a title or delete tasks with specific properties.
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			tasks, error := task.ReadTasks()

			if error != nil {
				return error
			}

			firstArgument := args[0]

			filteredTasks := lo.Filter(tasks, func(item task.Task, index int) bool {
				return item.Id() != firstArgument
			})

			if len(tasks) == len(filteredTasks) {

				return custom_errors.CreateInvalidArgumentErrorWithMessage(
					fmt.Sprintf("A task with this id %s doesn't exist", firstArgument),
				)

			}

			task.SaveTasks(filteredTasks)

			fmt.Fprintln(
				cmd.OutOrStdout(),
				fmt.Sprintf(
					"A task with this ID was deleted %s",
					firstArgument,
				),
			)

			return nil
		},
	}
	return deleteCmd
}

func init() {
	rootCmd.AddCommand(CreateDeleteCommand())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
