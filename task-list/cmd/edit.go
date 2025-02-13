/*
Copyright © 2024 Your Name
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

// CreateEditCmd represents the creation of the edit command
var CreateEditCmd = func() *cobra.Command {

	editCommand := &cobra.Command{
		Use:          "edit",
		Short:        "Edits a task",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		Long: `A task can be edited by using it's id.
			When editing a task you can pass in a flag to tell this command which property you want to change.
			The only ones that are supported are title, description, complete, and priority.
			If there are no flags passed through then you will see a form allowing you to edit all four of the following props.
		`,
		RunE: func(cmd *cobra.Command, args []string) error {

			id := args[0]

			tasks, err := task.ReadTasks()

			if err != nil {
				return err
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

			complete, completeErr := cmd.Flags().GetString(COMPLETE)
			title, titleErr := cmd.Flags().GetString(TITLE)
			description, descriptionErr := cmd.Flags().GetString(DESCRIPTION)
			priority, priorityErr := cmd.Flags().GetString(PRIORITY)

			flagErrors := errors.Join(
				titleErr,
				completeErr,
				descriptionErr,
				priorityErr,
			)

			if flagErrors != nil {
				return custom_errors.CreateInvalidFlagErrorWithMessage(
					flagErrors.Error(),
				)
			}

			if title != "" && foundTask.Title != title {
				foundTask.Title = title
				foundTask.UpdatedAt = time.Now()
			}

			if description != "" && foundTask.Description != description {
				foundTask.Description = description
				foundTask.UpdatedAt = time.Now()
			}

			if priority != "" {
				parsedPriority, err := task.ParsePriority(priority)

				if err != nil {
					return err
				}

				if parsedPriority.Value() != foundTask.Priority.Value() {
					foundTask.Priority = parsedPriority
					foundTask.UpdatedAt = time.Now()
				}
			}

			if complete != "" {
				if complete != "true" && complete != "false" {
					return custom_errors.CreateInvalidFlagErrorWithMessage(
						"complete flag must be either 'true' or 'false'",
					)
				}

				parsedBool, err := strconv.ParseBool(complete)

				if err != nil {
					return err
				}

				if strconv.FormatBool(foundTask.Complete) != complete {
					foundTask.Complete = parsedBool
					foundTask.UpdatedAt = time.Now()
				}
			}

			if err := task.SaveTasks(lo.Map(tasks, func(item task.Task, index int) task.Task {
				return lo.If(item.Id() == foundTask.Id(), foundTask).Else(item)
			})); err != nil {
				return err
			}

			taskAsJSON, err := foundTask.ToJSON()

			if err != nil {
				return err
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
	editCommand.Flags().String(DESCRIPTION, "", "Set the description of the task")
	editCommand.Flags().String(PRIORITY, "", "Set the priority of the task")
	editCommand.RegisterFlagCompletionFunc(
		PRIORITY,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return task.AllowedProrities, cobra.ShellCompDirectiveDefault
		},
	)
	editCommand.Flags().String(COMPLETE, "", "Mark task complete or not")

	return editCommand
}

func init() {
	rootCmd.AddCommand(CreateEditCmd())
}
