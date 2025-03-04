/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mini-clis/shared/custom_errors"
	"github.com/mini-clis/task-list/task"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const COMPLETION = "completion"
const INCOMPLETE = "incomplete"

var allowedCompletionValues = []string{
	COMPLETE,
	INCOMPLETE,
}

// deleteCmd represents the delete command
func CreateDeleteCommand() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a task based on an id",
		Long: `You can delete a task based on an Id.
The flags in this command allow you to pass in
a title or delete tasks with specific properties.
`,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks, error := task.ReadTasks()

			if error != nil {
				return error
			}

			firstArgument := args[0]

			title, titleError := cmd.Flags().GetBool(TITLE)
			completion, completionError := cmd.Flags().GetBool(COMPLETION)
			priority, priorityError := cmd.Flags().GetBool(PRIORITY)

			flagErrors := errors.Join(titleError, completionError, priorityError)

			if flagErrors != nil {
				return custom_errors.CreateInvalidFlagErrorWithMessage(
					custom_errors.FlagName(""),
					flagErrors.Error(),
				)
			}

			if completion && !lo.Contains(allowedCompletionValues, firstArgument) {
				return custom_errors.CreateInvalidArgumentErrorWithMessage(
					fmt.Sprintf(
						"When you use the %s flag you must pass in the %s",
						COMPLETION,
						strings.Join(allowedCompletionValues, ","),
					),
				)
			}

			if priority && !lo.Contains(task.AllowedProrities, firstArgument) {
				return custom_errors.CreateInvalidArgumentErrorWithMessage(
					fmt.Sprintf(
						"When you use the %s flag you must pass in the %s",
						PRIORITY,
						strings.Join(task.AllowedProrities, ","),
					),
				)
			}

			filteredTasks := lo.If(
				priority,
				lo.Filter(tasks, func(item task.Task, index int) bool {
					parsedPriority, _ := task.ParsePriority(firstArgument)
					return item.Priority != parsedPriority
				})).
				ElseIf(
					completion && firstArgument == COMPLETE,
					lo.Filter(tasks, func(item task.Task, index int) bool {
						return item.Complete == true
					})).
				ElseIf(
					completion && firstArgument == INCOMPLETE,
					lo.Filter(tasks, func(item task.Task, index int) bool {
						return item.Complete == false
					})).
				ElseIf(
					title,
					lo.Filter(tasks, func(item task.Task, index int) bool {
						return item.Title != firstArgument
					})).
				Else(
					lo.Filter(tasks, func(item task.Task, index int) bool {
						return item.Id() != firstArgument
					}),
				)

			if len(tasks) == len(filteredTasks) {
				return custom_errors.CreateInvalidArgumentErrorWithMessage(
					fmt.Sprintf("A task with this id %s doesn't exist", firstArgument),
				)
			}

			if error := task.SaveTasks(filteredTasks); error != nil {
				return error
			}

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

	allowedFlagNames := []string{PRIORITY, COMPLETION, TITLE}

	deleteCmdFlags := deleteCmd.Flags()

	lo.ForEach(allowedFlagNames, func(item string, index int) {
		deleteCmdFlags.Bool(
			item,
			false,
			fmt.Sprintf("Delete tasks based on %s", item),
		)
	})

	deleteCmd.MarkFlagsMutuallyExclusive(allowedFlagNames...)

	return deleteCmd
}

func init() {
	rootCmd.AddCommand(CreateDeleteCommand())
}
