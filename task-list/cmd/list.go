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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of all of your tasks",
	Long: `Get a list of all the tasks that you need to do today.
	You will see the
	`,
	Run: func(cmd *cobra.Command, args []string) {

		tasks, error := task.ReadTasks()

		if error != nil {

			log.Fatal(error)
		}

		fmt.Println("Here is the list of tasks you have to do")
		fmt.Printf("%#v", tasks)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
