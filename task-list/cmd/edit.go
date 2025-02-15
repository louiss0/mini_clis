/*
Copyright © 2024 Your Name
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mini-clis/task-list/custom_errors"
	"github.com/mini-clis/task-list/flags"
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

	titleFlag := flags.NewEmptyStringFlag(TITLE)
	descriptionFlag := flags.NewEmptyStringFlag(DESCRIPTION)
	priorityFlag := flags.NewUnionFlag(task.AllowedProrities, PRIORITY)
	completeFlag := flags.NewBoolFlag(COMPLETE)

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

			if titleFlag.String() != "" && foundTask.Title != titleFlag.String() {
				foundTask.Title = titleFlag.String()
				foundTask.UpdatedAt = time.Now()
			}

			if descriptionFlag.String() != "" && foundTask.Description != descriptionFlag.String() {
				foundTask.Description = descriptionFlag.String()
				foundTask.UpdatedAt = time.Now()
			}

			if priorityFlag.String() != "" {
				parsedPriority, err := task.ParsePriority(priorityFlag.String())

				if err != nil {
					return err
				}

				if parsedPriority.Value() != foundTask.Priority.Value() {
					foundTask.Priority = parsedPriority
					foundTask.UpdatedAt = time.Now()
				}
			}

			if completeFlag.String() != "" {

				if strconv.FormatBool(foundTask.Complete) != completeFlag.String() {

					foundTask.Complete = completeFlag.Value()

					foundTask.UpdatedAt = time.Now()
				}
			}

			updatedTasks := lo.Map(tasks, func(item task.Task, index int) task.Task {
				return lo.If(item.Id() == foundTask.Id(), foundTask).Else(item)
			})

			if err := task.SaveTasks(updatedTasks); err != nil {
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

	editCommand.Flags().Var(&titleFlag, TITLE, "Set the title of the task")
	editCommand.Flags().Var(&descriptionFlag, DESCRIPTION, "Set the description of the task")
	editCommand.Flags().Var(&priorityFlag, PRIORITY, "Set the priority of the task")
	editCommand.RegisterFlagCompletionFunc(
		PRIORITY,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return task.AllowedProrities, cobra.ShellCompDirectiveDefault
		},
	)

	editCommand.Flags().Var(&completeFlag, COMPLETE, "Mark task complete or not")

	return editCommand
}

func init() {
	rootCmd.AddCommand(CreateEditCmd())
}
