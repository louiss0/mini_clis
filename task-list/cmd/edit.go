/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mini-clis/task-list/custom_errors"
	"github.com/mini-clis/task-list/task"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const (
	TITLE       = "title"
	DESCRIPTION = "description"
	PRIORITY    = "priority"
	COMPLETE    = "complete"
)

// CreateEditCmd represents the creation of the  edit command
var CreateEditCmd = func() *cobra.Command {

	editCommand := &cobra.Command{
		Use:          "edit",
		Short:        "Edit's a task",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		Long: `A task can be edited by using it's id.
			When editing a task you can pass in a flag to tell this command which property you want to change.
			The only ones that are supported are title, description, complete, and priority.
			If there are no flags passed through then you will see a form allowing you to edit all four of the following props.
		`,
		RunE: func(cmd *cobra.Command, args []string) error {

			id := args[0]

			tasks, error := task.ReadTasks()

			if error != nil {
				return error
			}

			foundTask, ok := lo.Find(tasks, func(task task.Task) bool {

				return task.Id() == id

			})

			if !ok {

				return fmt.Errorf(
					"%w Task with this id wasn't found %s",
					custom_errors.InvalidArgument,
					id,
				)

			}

			complete, completeError := cmd.Flags().GetString(COMPLETE)

			title, titleError := cmd.Flags().GetString(TITLE)

			description, descriptionError := cmd.Flags().GetString(DESCRIPTION)

			priority, priorityError := cmd.Flags().GetString(PRIORITY)

			flagErrors := errors.Join(
				titleError,
				completeError,
				descriptionError,
				priorityError,
			)

			if flagErrors != nil {

				return custom_errors.CreateInvalidFlagErrorWithMessage(
					flagErrors.Error(),
				)
			}

			if title != "" {

				foundTask.Title = title
				foundTask.UpdatedAt = time.Now()
			}

			if description != "" {
				foundTask.Description = description
				foundTask.UpdatedAt = time.Now()

			}

			if priority != "" {
				parsedPriority, _ := task.ParsePriority(priority)
				foundTask.Priority = parsedPriority
				foundTask.UpdatedAt = time.Now()

			}

			if complete != "" {

				parsedBool, _ := strconv.ParseBool(complete)

				foundTask.Complete = parsedBool
				foundTask.UpdatedAt = time.Now()
			}

			taskAsJSON, error := foundTask.ToJSON()

			if error != nil {
				return error
			}

			fmt.Printf("Here is the task %s\n", id)
			fmt.Fprintln(
				cmd.OutOrStdout(),
				taskAsJSON,
			)

			return nil
		},
	}

	editCommand.Flags().String(TITLE, "", "Set the title of the task")
	editCommand.Flags().String(DESCRIPTION, "", "Set the desctiption of the task")
	editCommand.Flags().String(PRIORITY, "", "Set the priority of the task")

	editCommand.Flags().String(COMPLETE, "", "Mark task complete or not")

	return editCommand

}

func init() {
	rootCmd.AddCommand(CreateEditCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
