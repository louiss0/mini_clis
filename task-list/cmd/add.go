/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
	Each flag will be dedicated to helping create a task.
	`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {

			tasks, error := task.ReadTasks()

			if error != nil {
				return error
			}

			title, description := args[0], lo.TernaryF(len(args) == 2,
				func() string { return args[1] },
				func() string { return "" },
			)

			newTask := task.NewTask(title, description)

			priorityFlag, priorityFlagError := cmd.Flags().GetString(PRIORITY)

			if priorityFlagError != nil {
				return priorityFlagError
			}

			if priorityFlag != "" {

				_, error := task.ParsePriority(priorityFlag)

				if error != nil {
					return error
				}

			}

			task.SaveTasks(append([]task.Task{newTask}, tasks...))

			fmt.Println("This is the task you added")

			taskAsJson, unmarshallError := newTask.ToJSON()

			if unmarshallError != nil {
				return unmarshallError
			}

			fmt.Fprint(cmd.OutOrStdout(), taskAsJson)

			return nil
		},
	}

	command.Flags().StringP(PRIORITY, "p", "", "Decide the priority of a task")

	return command
}

func init() {
	rootCmd.AddCommand(CreateAddCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
