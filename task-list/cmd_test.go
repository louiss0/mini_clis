package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	. "github.com/mini-clis/task-list/cmd"
	"github.com/mini-clis/task-list/custom_errors"
	"github.com/mini-clis/task-list/task"
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
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
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

var getMockPersistedTasks = func() ([]mockPersistedTask, error) {

	var tasks []mockPersistedTask

	output, error := os.ReadFile(task.TASK_LIST_STORAGE_PATH)

	if error != nil {
		return tasks, error
	}

	unmarshalError := json.Unmarshal(output, &tasks)

	if unmarshalError != nil {
		return tasks, unmarshalError
	}

	return tasks, nil
}

var getMockPersistedTaskBasedOnOutput = func(output string, error error) (mockPersistedTask, error) {

	var task mockPersistedTask

	if error != nil {

		return task, error
	}

	unmarshalError := json.Unmarshal([]byte(output), &task)

	if unmarshalError != nil {
		return task, unmarshalError
	}

	return task, nil
}

var getRandomPersistedTask = func(tasks []mockPersistedTask) (mockPersistedTask, error) {

	if len(tasks) == 0 {
		return mockPersistedTask{}, fmt.Errorf("There are no tasks in the storage")
	}

	randomNumberBasedOnTaskLength := rand.New(
		rand.NewSource(time.Now().UnixNano())).
		Intn(len(tasks))

	return tasks[randomNumberBasedOnTaskLength], nil

}

var seedTasks = func(assert *assert.Assertions) {
	file, err := os.OpenFile(
		task.TASK_LIST_STORAGE_PATH,
		os.O_TRUNC|os.O_WRONLY,
		os.ModePerm,
	)
	if err != nil {
		log.Fatal(err)
	}

	fakeTasks := lo.Map(
		lo.Range(gofakeit.IntRange(3, 15)),
		func(item int, index int) mockPersistedTask {

			time.Sleep(time.Millisecond * 2)

			return mockPersistedTask{
				Id:    gofakeit.UUID(),
				Title: gofakeit.Sentence(gofakeit.IntRange(2, 12)),
				Description: gofakeit.Paragraph(
					gofakeit.IntRange(1, 5),
					gofakeit.IntRange(5, 10),
					gofakeit.IntRange(2, 15),
					"\n",
				),
				CreatedAt: time.Now().UnixMilli(),
				UpdatedAt: time.Now().UnixMilli(),
				Complete:  gofakeit.Bool(),
				Priority: gofakeit.RandomString([]string{
					task.HIGH.Value(),
					task.MEDIUM.Value(),
					task.LOW.Value(),
				}),
			}

		})

	byte, error := json.Marshal(fakeTasks)

	assert.NoError(error)

	file.Truncate(0)

	file.Write(byte)

	defer file.Close()
}

var _ = Describe("Cmd", func() {

	assert := assert.New(GinkgoT())

	rootCmd := RootCmd()

	seedTasks(assert)

	BeforeEach(func() {

		if len(rootCmd.Commands()) == 0 {
			rootCmd.AddCommand(
				CreateListCommand(),
				CreateEditCmd(),
				CreateAddCmd(),
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

		Context("Organizing Tasks", func() {

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

							if len(item) < 2 {
								return true
							}

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

							if len(item) < 2 {
								return true
							}

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

						if len(item) < 2 {
							return true
						}

						first, second := item[0], item[1]

						if reflect.TypeOf(second).Kind() != reflect.Struct {
							return true
						}

						return first.CreatedAt > second.CreatedAt

					})

					assert.True(allTasksAreSortedByTheHigestOrder)

				},
			)

			It(
				"sorts tasks by the ones that were inserted at the latest times when sort-date flag is passed 'earliest'",
				func() {

					tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(SORT_DATE), EARLIEST))

					assert.NoError(error)

					allTasksAreSortedByTheHigestOrder := lo.EveryBy(lo.Chunk(tasks, 2), func(item []mockPersistedTask) bool {

						if len(item) < 2 {
							return true
						}

						first, second := item[0], item[1]

						if reflect.TypeOf(second).Kind() != reflect.Struct {
							return true
						}

						return first.CreatedAt < second.CreatedAt

					})

					assert.True(allTasksAreSortedByTheHigestOrder)

				},
			)

			lo.ForEach([]string{
				task.HIGH.Value(),
				task.LOW.Value(),
				task.MEDIUM.Value(),
			},
				func(priority string, index int) {

					It(
						fmt.Sprintf("filters tasks by the highest priority when the --filter-priority is passed %s", priority),
						func() {

							tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(FILTER_PRIORITY), priority))

							assert.NoError(error)

							allTasksHaveTheSamePriority := lo.EveryBy(tasks, func(task mockPersistedTask) bool {

								return task.Priority == priority

							})

							assert.True(allTasksHaveTheSamePriority)

						},
					)

				})

			It(
				"filters only tasks that are complete when the --filter-complete flag is passed",
				func() {
					tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(FILTER_COMPLETE)))

					assert.NoError(error)

					allCompletedTasks := lo.EveryBy(
						tasks,
						func(item mockPersistedTask) bool {

							return item.Complete == true
						})

					assert.True(allCompletedTasks)

				})

			It(
				"filters only tasks that are incomplete whe the --filter-incomplete flag is passed",
				func() {

					tasks, error := extractPersistedTasksFromOutput(executeCommand(rootCmd, "list", createFlag(FILTER_INCOMPLETE)))

					assert.NoError(error)

					allIncompleteTasks := lo.EveryBy(
						tasks,
						func(item mockPersistedTask) bool {

							return item.Complete == false
						})

					assert.True(allIncompleteTasks)
				})

		})

	})

	Context("Editing tasks", Ordered, func() {

		var mockTasks []mockPersistedTask

		BeforeAll(func() {

			tasks, error := getMockPersistedTasks()

			assert.NoError(error)

			assert.NotEmpty(tasks)

			mockTasks = tasks

		})

		var mockTask mockPersistedTask

		BeforeEach(func() {

			storageTask, storageError := getRandomPersistedTask(mockTasks)

			mockTask = storageTask
			assert.NoError(storageError)
			assert.NotEmpty(storageTask)

		})

		type EditCase struct {
			FlagName string
			Argument string
		}

		fakeEditCase := func(flagName string) EditCase {

			return lo.Switch[string, EditCase](flagName).
				Case(TITLE,
					EditCase{
						TITLE,
						gofakeit.Sentence(gofakeit.Number(1, 16)),
					},
				).
				Case(
					DESCRIPTION,
					EditCase{
						DESCRIPTION,
						gofakeit.Paragraph(1, gofakeit.Number(3, 11), gofakeit.Number(1, 16), " "),
					},
				).
				Case(
					PRIORITY,
					EditCase{
						PRIORITY,
						gofakeit.RandomString(task.AllowedProrities),
					},
				).
				Case(
					COMPLETE,
					EditCase{
						COMPLETE,
						gofakeit.RandomString([]string{
							"true",
							"false",
						}),
					}).
				Default(EditCase{})

		}

		lo.ForEach([]EditCase{
			fakeEditCase(TITLE),
			fakeEditCase(DESCRIPTION),
			fakeEditCase(PRIORITY),
			fakeEditCase(COMPLETE),
		}, func(editCase EditCase, index int) {
			It(
				fmt.Sprintf(
					"edits a task's %s field when %s is passed through",
					editCase.FlagName,
					createFlag(editCase.FlagName)),
				func() {

					taskFromOutput, outputError := getMockPersistedTaskBasedOnOutput(
						executeCommand(
							rootCmd,
							"edit",
							mockTask.Id,
							createFlag(editCase.FlagName),
							editCase.Argument,
						),
					)

					assert.NoError(outputError)
					assert.NotEmpty(taskFromOutput)

					capitalisedFlagName := lo.Capitalize(editCase.FlagName)

					taskFieldValueBasedOnFlagName := reflect.ValueOf(mockTask).
						FieldByName(capitalisedFlagName).Interface()

					taskFromOutputFieldValueBasedOnFlagName := reflect.ValueOf(taskFromOutput).
						FieldByName(capitalisedFlagName).Interface()

					assert.Truef(lo.Ternary(
						taskFieldValueBasedOnFlagName != taskFromOutputFieldValueBasedOnFlagName,
						mockTask.UpdatedAt != taskFromOutput.UpdatedAt,
						mockTask.UpdatedAt == taskFromOutput.UpdatedAt,
					),
						strings.Join(
							[]string{
								"The value of the updatedAt field from a task in storage is only supposed to change when the the value a field changes",
								"%s before %s vs %s after %s",
								"Updated At Field Before %s vs Updated At Field After %s",
							},
							"\n",
						),
						capitalisedFlagName,
						taskFieldValueBasedOnFlagName,
						capitalisedFlagName,
						taskFromOutputFieldValueBasedOnFlagName,
						mockTask.UpdatedAt,
						taskFromOutput.UpdatedAt,
					)

				})

		})

		It("errors when --priority is passed the wrong value", func() {

			taskFromOutput, outputError := getMockPersistedTaskBasedOnOutput(
				executeCommand(
					rootCmd,
					"edit",
					mockTask.Id,
					createFlag(PRIORITY),
					"beem boom boom boom bop bam!",
				),
			)

			assert.Error(outputError)
			assert.Empty(taskFromOutput)

		})

		It("errors when --complete is passed the wrong value", func() {
			taskFromOutput, outputError := getMockPersistedTaskBasedOnOutput(
				executeCommand(
					rootCmd,
					"edit",
					mockTask.Id,
					createFlag(COMPLETE),
					"beem boom boom boom bop bam!",
				),
			)

			assert.Error(outputError)
			assert.Empty(taskFromOutput)

		})

		Context("Adding tasks", Ordered, func() {

			generateFakeTitle := func() string {

				return gofakeit.Sentence(gofakeit.IntRange(
					1,
					5,
				))
			}

			generateFakeDescription := func() string {

				return gofakeit.Paragraph(
					gofakeit.IntRange(1, 5),
					gofakeit.IntRange(5, 10),
					gofakeit.IntRange(7, 12),
					"\n",
				)
			}

			AfterEach(func() {

				time.Sleep(time.Second)

			})

			It("works", func() {

				previousTasksFromStorage, error := getMockPersistedTasks()
				assert.NoError(error)
				assert.NotEmpty(previousTasksFromStorage)

				task, error := getMockPersistedTaskBasedOnOutput(
					executeCommand(
						rootCmd,
						"add",
						generateFakeTitle(),
					),
				)

				assert.NoError(error)
				assert.NotEmpty(task)

				newMockPersistedTasks, error := getMockPersistedTasks()

				assert.NoError(error)
				assert.NotEmpty(newMockPersistedTasks)

				newMockPersistedTasksLength := len(newMockPersistedTasks)
				previousTasksFromStorageLength := len(previousTasksFromStorage)

				assert.Truef(
					newMockPersistedTasksLength > previousTasksFromStorageLength,
					strings.Join([]string{
						"The length of the previous tasks aren't longer than the length of the new ones",
						"Amount of previous tasks: %d",
						"Amount of new tasks %d",
					}, "\n"),
					previousTasksFromStorageLength,
					newMockPersistedTasksLength,
				)

			})

			It("adds a description to the task when a second argument is passed", func() {

				task, error := getMockPersistedTaskBasedOnOutput(
					executeCommand(
						rootCmd,
						"add",
						generateFakeTitle(),
						generateFakeDescription(),
					),
				)

				assert.NoError(error)
				assert.NotEmpty(task)

				assert.NotEmpty(task.Description)

			})

		})

	})

})
