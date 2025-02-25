/*
Copyright © 2024 Your Name
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/mini-clis/shared/custom_errors"
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

			title := titleFlag.String()
			description := descriptionFlag.String()
			priority := priorityFlag.String()
			complete := completeFlag.String()

			everyFlagValueIsEmpty := lo.EveryBy(
				[]string{title, description, priority, complete},
				func(flagValue string) bool {
					return flagValue == ""
				},
			)

			if everyFlagValueIsEmpty {

				title = foundTask.Title
				description = foundTask.Description
				priority = foundTask.Priority.Value()
				internalComplete := foundTask.Complete

				form := huh.NewForm(
					huh.NewGroup(
						huh.NewText().
							Title("Title").
							CharLimit(80).
							Placeholder("What are you doing?").
							Lines(1).
							Value(&title).
							Validate(huh.ValidateNotEmpty()),
						huh.NewText().
							Title("Description").
							Value(&description).
							CharLimit(80).
							Lines(2).
							Placeholder("Why are you doing this task?").
							Validate(huh.ValidateNotEmpty()),
						huh.NewSelect[string]().
							Title("Priority").
							Value(&priority).
							Validate(huh.ValidateNotEmpty()).
							Description("How important is this task").
							Options(
								huh.NewOption("Low", "low"),
								huh.NewOption("Medium", "medium"),
								huh.NewOption("High", "high"),
							),
						huh.NewConfirm().
							Key("complete").
							Value(&internalComplete).
							Description("Is this task complete?").
							Title("Complete").
							Affirmative("Yes").
							Negative("No"),
					),
				)

				if err := form.Run(); err != nil {
					return err
				}

				complete = fmt.Sprintf("%t", internalComplete)

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

				if strconv.FormatBool(foundTask.Complete) != complete {

					parsedComplete, err := strconv.ParseBool(complete)

					if err != nil {
						return err
					}

					foundTask.Complete = parsedComplete

					foundTask.UpdatedAt = time.Now()
				}
			}

			updatedTasks := lo.Map(tasks, func(item task.Task, index int) task.Task {
				return lo.If(item.Id() == foundTask.Id(), foundTask).Else(item)
			})

			if err := task.SaveTasks(updatedTasks); err != nil {
				return err
			}

			plain, error := cmd.Flags().GetBool(PLAIN)

			if error != nil {
				return error
			}

			if plain {
				taskAsJSON, err := foundTask.ToJSON()

				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), taskAsJSON)

				return nil
			}

			taskAsJSON, err := foundTask.ToPrettyJSON()

			if err != nil {
				return err
			}

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
