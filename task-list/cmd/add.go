/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slices"

	"github.com/mini-clis/task-list/task"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
func CreateAddCmd() *cobra.Command {
    command := &cobra.Command{
        Use:   "add",
        Short: "Add a task to the list of tasks",
        Long: `This command allows you to add a task to the list.
    When you do you must supply a title for your task. you decide to store a task you can set other things using flags.
    The first argument will be the task title the second is the description.
    You can decide a priority by passing in the --priority flag.
    `,
        Args: cobra.RangeArgs(1, 2),
        RunE: func(cmd *cobra.Command, args []string) error {
            tasks, error := task.ReadTasks()

            if error != nil {
                return error
            }

            title, description := args[0], lo.TernaryF(
                len(args) == 2,
                func() string { return args[1] },
                func() string { return "" },
            )

            newTask := task.NewTask(title, description)

            priorityFlag, error := cmd.Flags().GetString(PRIORITY)

            if error != nil {
                return error
            }

            if priorityFlag != "" {
                priority, error := task.ParsePriority(priorityFlag)

                if error != nil {
                    return error
                }

                newTask.Priority = priority
            }

            if error := task.SaveTasks(slices.Insert(tasks, 0, newTask)); error != nil {
                return error
            }

            fmt.Println("This is the task you added")

            taskAsJSON, error := newTask.ToJSON()

            if error != nil {
                return error
            }

            fmt.Fprint(cmd.OutOrStdout(), taskAsJSON)

            return nil
        },
    }

    command.Flags().StringP(PRIORITY, "p", "", "Decide the priority of a task")

    command.RegisterFlagCompletionFunc(PRIORITY, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        return task.AllowedProrities, cobra.ShellCompDirectiveDefault
    })

    return command
}

func init() {
    rootCmd.AddCommand(CreateAddCmd())
}
