/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/mini-clis/task-list/custom_errors"
	"github.com/mini-clis/task-list/flags"
	"github.com/mini-clis/task-list/task"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const UI = "ui"

// addCmd represents the add command
func CreateAddCmd() *cobra.Command {

	priorityFlag := flags.NewUnionFlag(task.AllowedProrities, PRIORITY)

	command := &cobra.Command{
		Use:   "add",
		Short: "Add a task to the list of tasks",
		Long: `This command allows you to add a task to the list.
    When you do you must supply a title for your task. you decide to store a task you can set other things using flags.
    The first argument will be the task title the second is the description.
    You can decide a priority by passing in the --priority flag.
    `,
		Args: func(cmd *cobra.Command, args []string) error {
			ui, error := cmd.Flags().GetBool(UI)

			if error != nil {
				return error
			}

			if len(args) == 0 && !ui {
				return custom_errors.CreateInvalidArgumentErrorWithMessage(
					"You must pass in the UI flag if you want to use the UI while passing no arguments",
				)
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			tasks, error := task.ReadTasks()

			if error != nil {
				return error
			}

			var title, description string
			priority := priorityFlag.String()

			ui, _ := cmd.Flags().GetBool(UI)

			if ui {

				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							CharLimit(80).
							Title("Title").
							Placeholder("What is your task?").
							Validate(huh.ValidateNotEmpty()).
							Value(&title),
						huh.NewText().
							Title("Description").
							Placeholder("Why is this task important?").
							Value(&description),
						huh.NewSelect[string]().
							Title("Priority").
							Description("How important is this task?").
							Value(&priority).
							Validate(huh.ValidateOneOf(task.AllowedProrities...)).
							Options(lo.Map(
								task.AllowedProrities,
								func(item string, index int) huh.Option[string] {
									return huh.NewOption(lo.Capitalize(item), item)
								})...),
					),
				)

				if error := form.Run(); error != nil {
					return error
				}

			} else {

				title, description = args[0], lo.TernaryF(
					len(args) == 2,
					func() string { return args[1] },
					func() string { return "" },
				)
			}

			newTask := task.NewTask(title, description)

			if priority != "" {

				parsedPriority, error := task.ParsePriority(priority)

				if error != nil {
					return error
				}

				newTask.Priority = parsedPriority
			}

			if error := task.SaveTasks(slices.Insert(tasks, 0, newTask)); error != nil {
				return error
			}

			plain, error := cmd.Flags().GetBool(PLAIN)

			if error != nil {
				return error
			}

			if plain {
				taskAsJSON, error := newTask.ToJSON()

				if error != nil {
					return error
				}
				fmt.Fprint(cmd.OutOrStdout(), taskAsJSON)

				return nil
			}

			taskAsJSON, error := newTask.ToPrettyJSON()

			if error != nil {
				return error
			}

			fmt.Fprintln(
				cmd.OutOrStdout(),
				taskAsJSON,
			)

			return nil
		},
	}

	command.Flags().VarP(&priorityFlag, PRIORITY, "p", "Decide the priority of a task")

	command.Flags().Bool(UI, false, "Render a ui for creating a tasks instead of passing arguments")

	command.RegisterFlagCompletionFunc(PRIORITY, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return task.AllowedProrities, cobra.ShellCompDirectiveDefault
	})

	command.MarkFlagsMutuallyExclusive(UI, PRIORITY)

	return command
}

func init() {
	rootCmd.AddCommand(CreateAddCmd())
}
