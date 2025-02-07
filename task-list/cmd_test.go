package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	. "github.com/mini-clis/task-list/cmd"
	"github.com/mini-clis/task-list/custom_errors"
	. "github.com/onsi/ginkgo/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type mockPersistedTask struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Complete    bool   `json:"complete"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

var executeCommand = func(cmd *cobra.Command, args ...string) (string, error) {

	var buffer bytes.Buffer

	cmd.SetOut(&buffer)

	cmd.SetArgs(nil)

	cmd.SetArgs(args)

	error := cmd.Execute()

	return buffer.String(), error

}

var createFlag = func(flagName string) string {

	return fmt.Sprintf("--%s", flagName)
}

var _ = Describe("Cmd", func() {

	assert := assert.New(GinkgoT())

	rootCmd := RootCmd()

	BeforeEach(func() {

		if len(rootCmd.Commands()) == 0 {
			rootCmd.AddCommand(
				CreateListCommand(),
			)

		}

	})

	AfterEach(func() {

		rootCmd.ResetCommands()

	})

	Context("List", func() {

		extractPersistedTasksFromOutput := func(commandOutput string, error error) ([]mockPersistedTask, error) {

			var tasks []mockPersistedTask

			if error != nil {

				return tasks, error
			}

			unmarshalError := json.Unmarshal([]byte(commandOutput), &tasks)

			if error != nil {

				return tasks, unmarshalError
			}

			return tasks, nil
		}

		It("works", func() {

			output, error := executeCommand(rootCmd, "list")

			assert.NotEmpty(output)

			assert.Nil(error)

			var tasks []mockPersistedTask

			json.Unmarshal([]byte(output), &tasks)

			assert.Greater(len(tasks), 0)

		})

		lo.ForEach([]string{
			FILTER_PRIORITY,
			SORT_DATE,
			SORT_PRIORITY,
		}, func(item string, index int) {

			It(
				fmt.Sprintf("creates an error when wrong value is passed to %s", item),
				func() {

					output, error := executeCommand(
						rootCmd,
						"list",
						createFlag(item),
						"foo",
					)

					assert.Empty(output)

					assert.ErrorIs(error, custom_errors.InvalidFlag)

				},
			)

		})

		Context("Organizing Tasks", Ordered, func() {

			It(
				"sorts tasks that by the ones that were inserted at the earliest times when sort-priority flag is passed 'highest'",
				func() {
					tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(SORT_PRIORITY), HIGHEST))

					assert.NoError(error)

					assert.Greater(len(tasks), 1)

					priorityMap := map[string]int{
						"high":   3,
						"medium": 2,
						"low":    1,
					}

					allTasksAreSortedByTheHigestOrder := lo.EveryBy(
						lo.Chunk(tasks, 2),
						func(item []mockPersistedTask) bool {
							first, second := item[0], item[1]

							if reflect.TypeOf(second).Kind() != reflect.Struct {
								return true
							}

							return priorityMap[first.Priority] >= priorityMap[second.Priority]

						})

					assert.True(allTasksAreSortedByTheHigestOrder)
				},
			)

			It(
				"sorts tasks by the ones that were inserted at the latest times when sort-priority flag is passed 'lowest'",
				func() {

					tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(SORT_PRIORITY), LOWEST))

					assert.NoError(error)

					assert.Greater(len(tasks), 1)

					priorityMap := map[string]int{
						"high":   3,
						"medium": 2,
						"low":    1,
					}

					allTasksAreSortedByTheLowestOrder := lo.EveryBy(
						lo.Chunk(tasks, 2),
						func(item []mockPersistedTask) bool {
							first, second := item[0], item[1]

							if reflect.TypeOf(second).Kind() != reflect.Struct {
								return true
							}

							return priorityMap[first.Priority] <= priorityMap[second.Priority]

						})

					assert.True(allTasksAreSortedByTheLowestOrder)

				},
			)

			It(
				"sorts tasks that by the ones that were inserted at the earliest times when sort-date flag is passed 'latest'",
				func() {

					tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(SORT_DATE), LATEST))

					assert.NoError(error)

					assert.Greater(len(tasks), 1)

					allTasksAreSortedByTheHigestOrder := lo.EveryBy(lo.Chunk(tasks, 2), func(item []mockPersistedTask) bool {
						first, second := item[0], item[1]

						if reflect.TypeOf(second).Kind() != reflect.Struct {
							return true
						}

						firstCreatedTime, firstCreatedError := time.Parse(time.UnixDate, first.CreatedAt)

						secondCreatedTime, secondCreatedError := time.Parse(time.UnixDate, second.CreatedAt)

						if firstCreatedError != nil || secondCreatedError != nil {

							return false
						}

						return firstCreatedTime.After(secondCreatedTime)

					})

					assert.True(allTasksAreSortedByTheHigestOrder)

				},
			)

			It(
				"sorts tasks by the ones that were inserted at the latest times when sort-date flag is passed 'earliest'",
				func() {

					tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(SORT_DATE), EARLIEST))

					assert.NoError(error)

					assert.Greater(len(tasks), 1)

					allTasksAreSortedByTheHigestOrder := lo.EveryBy(lo.Chunk(tasks, 2), func(item []mockPersistedTask) bool {
						first, second := item[0], item[1]

						if reflect.TypeOf(second).Kind() != reflect.Struct {
							return true
						}

						firstCreatedTime, firstCreatedError := time.Parse(time.UnixDate, first.CreatedAt)

						secondCreatedTime, secondCreatedError := time.Parse(time.UnixDate, second.CreatedAt)

						if firstCreatedError != nil || secondCreatedError != nil {

							return false
						}

						return firstCreatedTime.Before(secondCreatedTime)

					})

					assert.True(allTasksAreSortedByTheHigestOrder)

				},
			)

			PIt(
				"filters tasks by the highest priority when the --filter-priority is passed 'highest'",
				func() {

				},
			)

			PIt("filters tasks by the highest priority when the --filter-priority is passed 'lowest'",
				func() {

				},
			)

			PIt(
				"filters only tasks that are complete when the --filter-incomplete flag is passed",
				func() {

				})

			PIt(
				"filters only tasks that are incomplete whe the --filter-complete flag is passed",
				func() {

				})

		})

	})

})
