package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	. "github.com/mini-clis/task-list/cmd"
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

// Helper Functions
var executeCommand = func(cmd *cobra.Command, args ...string) (string, error) {
	var buffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.SetOut(&buffer)
	cmd.SetErr(&errBuffer)
	cmd.SetArgs(nil)
	cmd.SetArgs(args)

	if errBuffer.String() != "" {
		return "", fmt.Errorf("command failed: %s", errBuffer.String())
	}

	err := cmd.Execute()

	return buffer.String(), err
}

var createFlag = func(flagName string) string {
	return fmt.Sprintf("--%s", flagName)
}

var getMockPersistedTasks = func() ([]mockPersistedTask, error) {
	var tasks []mockPersistedTask

	data, err := os.ReadFile(task.TASK_LIST_STORAGE_PATH)
	if err != nil {
		return tasks, err
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return tasks, err
	}

	return tasks, nil
}

var getMockPersistedTaskBasedOnOutput = func(output string, err error) (mockPersistedTask, error) {
	var task mockPersistedTask

	if err != nil {
		return task, err
	}

	if err := json.Unmarshal([]byte(output), &task); err != nil {
		return task, err
	}

	return task, nil
}

var getRandomPersistedTask = func(tasks []mockPersistedTask) (mockPersistedTask, error) {
	if len(tasks) == 0 {
		return mockPersistedTask{}, fmt.Errorf("no tasks found in storage")
	}

	randomIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(tasks))
	return tasks[randomIndex], nil
}

var seedTasks = func() {
	file, err := os.OpenFile(
		task.TASK_LIST_STORAGE_PATH,
		os.O_TRUNC|os.O_WRONLY,
		os.ModePerm,
	)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fakeTasks := lo.Map(
		lo.Range(gofakeit.IntRange(25, 50)),
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

	data, err := json.Marshal(fakeTasks)
	if err != nil {
		log.Fatal(err)
	}

	if err := file.Truncate(0); err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

var _ = Describe("Cmd", func() {
	assert := assert.New(GinkgoT())
	rootCmd := RootCmd()

	seedTasks()

	BeforeEach(func() {
		if len(rootCmd.Commands()) == 0 {
			rootCmd.AddCommand(
				CreateListCommand(),
				CreateEditCmd(),
				CreateAddCmd(),
				CreateDeleteCommand(),
			)
		}
	})

	AfterEach(func() {
		rootCmd.ResetCommands()
	})

	Context("List", func() {
		extractPersistedTasksFromOutput := func(commandOutput string, err error) ([]mockPersistedTask, error) {
			var tasks []mockPersistedTask

			if err != nil {
				return tasks, err
			}

			if err := json.Unmarshal([]byte(commandOutput), &tasks); err != nil {
				return tasks, err
			}

			return tasks, nil
		}

		It("successfully lists all tasks", func() {
			output, err := executeCommand(rootCmd, "list")
			assert.NoError(err)
			assert.NotEmpty(output)

			var tasks []mockPersistedTask
			err = json.Unmarshal([]byte(output), &tasks)
			assert.NoError(err)
			assert.Greater(len(tasks), 0)
		})

		Context("Flag validation", func() {
			lo.ForEach([]string{
				FILTER_PRIORITY,
				SORT_DATE,
				SORT_PRIORITY,
			}, func(item string, index int) {
				It(fmt.Sprintf("returns an error when invalid value is provided for %s", item), func() {
					output, err := executeCommand(
						rootCmd,
						"list",
						createFlag(item),
						"invalid-value",
					)

					assert.Empty(output)
					assert.Error(err)
				})
			})

			Context("Task Organization", func() {
				It("sorts tasks by highest priority", func() {
					tasks, err := extractPersistedTasksFromOutput(
						executeCommand(rootCmd, "list", createFlag(SORT_PRIORITY), HIGHEST),
					)

					assert.NoError(err)
					assert.Greater(len(tasks), 1)

					priorityMap := map[string]int{
						"high":   3,
						"medium": 2,
						"low":    1,
					}

					allTasksAreSortedByHighestOrder := lo.EveryBy(
						lo.Chunk(tasks, 2),
						func(item []mockPersistedTask) bool {
							if len(item) < 2 {
								return true
							}

							first, second := item[0], item[1]
							return priorityMap[first.Priority] >= priorityMap[second.Priority]
						})

					assert.True(allTasksAreSortedByHighestOrder)
				})

				// Continue with sort and filter tests...
				It("sorts tasks by lowest priority", func() {
					tasks, err := extractPersistedTasksFromOutput(
						executeCommand(rootCmd, "list", createFlag(SORT_PRIORITY), LOWEST),
					)

					assert.NoError(err)
					assert.Greater(len(tasks), 1)

					priorityMap := map[string]int{
						"high":   3,
						"medium": 2,
						"low":    1,
					}

					allTasksAreSortedByLowestOrder := lo.EveryBy(
						lo.Chunk(tasks, 2),
						func(item []mockPersistedTask) bool {
							if len(item) < 2 {
								return true
							}
							return priorityMap[item[0].Priority] <= priorityMap[item[1].Priority]
						})

					assert.True(allTasksAreSortedByLowestOrder)
				})

				It("sorts tasks by most recent date", func() {
					tasks, err := extractPersistedTasksFromOutput(
						executeCommand(rootCmd, "list", createFlag(SORT_DATE), LATEST),
					)

					assert.NoError(err)
					assert.Greater(len(tasks), 1)

					tasksAreSortedByLatest := lo.EveryBy(
						lo.Chunk(tasks, 2),
						func(item []mockPersistedTask) bool {
							if len(item) < 2 {
								return true
							}
							return item[0].CreatedAt > item[1].CreatedAt
						})

					assert.True(tasksAreSortedByLatest)
				})

				It("sorts tasks by oldest date", func() {
					tasks, err := extractPersistedTasksFromOutput(
						executeCommand(rootCmd, "list", createFlag(SORT_DATE), EARLIEST),
					)

					assert.NoError(err)
					assert.Greater(len(tasks), 1)

					tasksAreSortedByEarliest := lo.EveryBy(
						lo.Chunk(tasks, 2),
						func(item []mockPersistedTask) bool {
							if len(item) < 2 {
								return true
							}
							return item[0].CreatedAt < item[1].CreatedAt
						})

					assert.True(tasksAreSortedByEarliest)
				})
			})
		})

	})

	Context("Editing tasks", Ordered, func() {
		var mockTasks []mockPersistedTask
		var mockTask mockPersistedTask

		BeforeAll(func() {
			tasks, err := getMockPersistedTasks()
			assert.NoError(err)
			assert.NotEmpty(tasks)
			mockTasks = tasks
		})

		BeforeEach(func() {
			storageTask, err := getRandomPersistedTask(mockTasks)
			assert.NoError(err)
			assert.NotEmpty(storageTask)
			mockTask = storageTask
		})

		type EditCase struct {
			FlagName string
			Argument string
		}

		lo.ForEach([]EditCase{
			{FlagName: TITLE,
				Argument: gofakeit.Sentence(gofakeit.Number(1, 16)),
			},
			{
				FlagName: DESCRIPTION,
				Argument: gofakeit.Paragraph(1, gofakeit.Number(3, 11), gofakeit.Number(1, 16), " "),
			},
			{
				FlagName: PRIORITY,
				Argument: gofakeit.RandomString(task.AllowedProrities),
			},
			{
				FlagName: COMPLETE,
				Argument: gofakeit.RandomString([]string{"true", "false"}),
			},
		}, func(editCase EditCase, index int) {
			It(
				fmt.Sprintf("successfully updates task's %s field", editCase.FlagName),
				func() {
					taskFromOutput, err := getMockPersistedTaskBasedOnOutput(
						executeCommand(
							rootCmd,
							"edit",
							mockTask.Id,
							createFlag(editCase.FlagName),
							editCase.Argument,
						),
					)

					assert.NoError(err)
					assert.NotEmpty(taskFromOutput)

					capitalizedFlagName := lo.Capitalize(editCase.FlagName)
					oldValue := reflect.ValueOf(mockTask).FieldByName(capitalizedFlagName).Interface()
					newValue := reflect.ValueOf(taskFromOutput).FieldByName(capitalizedFlagName).Interface()

					message := fmt.Sprintf(
						"Field %s update validation failed:\nOld value: %v\nNew value: %v\nUpdated at (old): %d\nUpdated at (new): %d",
						capitalizedFlagName,
						oldValue,
						newValue,
						mockTask.UpdatedAt,
						taskFromOutput.UpdatedAt,
					)

					assert.True(
						lo.Ternary(
							oldValue != newValue,
							mockTask.UpdatedAt != taskFromOutput.UpdatedAt,
							mockTask.UpdatedAt == taskFromOutput.UpdatedAt,
						),
						message,
					)
				})
		})

		It("returns an error when invalid priority value is provided", func() {
			taskFromOutput, err := getMockPersistedTaskBasedOnOutput(
				executeCommand(
					rootCmd,
					"edit",
					mockTask.Id,
					createFlag(PRIORITY),
					"invalid-priority",
				),
			)

			assert.Error(err)
			assert.Empty(taskFromOutput)
		})

		It("returns an error when invalid complete value is provided", func() {
			taskFromOutput, err := getMockPersistedTaskBasedOnOutput(
				executeCommand(
					rootCmd,
					"edit",
					mockTask.Id,
					createFlag(COMPLETE),
					"invalid-boolean",
				),
			)

			assert.Error(err)
			assert.Empty(taskFromOutput)
		})
	})

	Context("Adding tasks", Ordered, func() {
		generateFakeTitle := func() string {
			return gofakeit.Sentence(gofakeit.IntRange(1, 5))
		}

		generateFakeDescription := func() string {
			return gofakeit.Paragraph(
				gofakeit.IntRange(1, 5),
				gofakeit.IntRange(5, 10),
				gofakeit.IntRange(7, 12),
				"\n",
			)
		}

		It("successfully adds a new task", func() {
			previousTasks, err := getMockPersistedTasks()
			assert.NoError(err)
			assert.NotEmpty(previousTasks)

			newTask, err := getMockPersistedTaskBasedOnOutput(
				executeCommand(
					rootCmd,
					"add",
					generateFakeTitle(),
				),
			)
			assert.NoError(err)
			assert.NotEmpty(newTask)

			currentTasks, err := getMockPersistedTasks()
			assert.NoError(err)
			assert.NotEmpty(currentTasks)

			assert.Greater(
				len(currentTasks),
				len(previousTasks),
				fmt.Sprintf(
					"Expected task count to increase.\nPrevious count: %d\nCurrent count: %d",
					len(previousTasks),
					len(currentTasks),
				),
			)
		})

		It("adds a task with description when provided", func() {
			task, err := getMockPersistedTaskBasedOnOutput(
				executeCommand(
					rootCmd,
					"add",
					generateFakeTitle(),
					generateFakeDescription(),
				),
			)

			assert.NoError(err)
			assert.NotEmpty(task)
			assert.NotEmpty(task.Description, "Task description should not be empty")
		})

		lo.ForEach([]struct {
			priority    string
			description string
		}{
			{task.HIGH.Value(), "high priority"},
			{task.MEDIUM.Value(), "medium priority"},
			{task.LOW.Value(), "low priority"},
		}, func(testCase struct {
			priority    string
			description string
		}, index int) {
			It(fmt.Sprintf("sets %s when specified", testCase.description), func() {
				task, err := getMockPersistedTaskBasedOnOutput(
					executeCommand(
						rootCmd,
						"add",
						generateFakeTitle(),
						createFlag(PRIORITY),
						testCase.priority,
					),
				)

				assert.NoError(err)
				assert.NotEmpty(task)
				assert.Equal(
					task.Priority,
					testCase.priority,
					fmt.Sprintf("Task priority should be %s", testCase.priority),
				)
			})
		})
	})

	Context("Deleting tasks", func() {
		oldPersistedTasks := []mockPersistedTask{}

		assertTasksAreDeleted := func(oldTasks, newTasks []mockPersistedTask, contextMessage string) {
			oldTaskCount := len(oldTasks)
			newTaskCount := len(newTasks)

			assert.Greater(
				oldTaskCount,
				newTaskCount,
				fmt.Sprintf(
					"%s\nOld task count: %d\nNew task count: %d",
					contextMessage,
					oldTaskCount,
					newTaskCount,
				),
			)
		}

		BeforeEach(func() {
			tasks, err := getMockPersistedTasks()
			assert.NoError(err)
			assert.NotEmpty(tasks)
			oldPersistedTasks = tasks
		})

		AfterEach(func() {
			seedTasks()
		})

		It("successfully deletes a task by ID", func() {
			randomTask, err := getRandomPersistedTask(oldPersistedTasks)
			assert.NoError(err)
			assert.NotEmpty(randomTask)

			output, err := executeCommand(rootCmd, "delete", randomTask.Id)
			assert.NoError(err)
			assert.NotEmpty(output)

			currentTasks, err := getMockPersistedTasks()
			assert.NoError(err)

			assertTasksAreDeleted(
				oldPersistedTasks,
				currentTasks,
				fmt.Sprintf("Failed to delete task with ID: %s", randomTask.Id),
			)
		})

		lo.ForEach([]string{
			PRIORITY,
			COMPLETION,
		}, func(flagName string, index int) {
			It(fmt.Sprintf("returns an error when invalid value is provided for %s flag", flagName), func() {
				output, err := executeCommand(
					rootCmd,
					"delete",
					"invalid-value",
					createFlag(flagName),
				)

				assert.Error(err)
				assert.Empty(output)
			})
		})

		lo.ForEach(task.AllowedProrities, func(priority string, index int) {
			It(fmt.Sprintf("deletes all tasks with %s priority", priority), func() {
				output, err := executeCommand(
					rootCmd,
					"delete",
					priority,
					createFlag(PRIORITY),
				)
				assert.NoError(err)
				assert.NotEmpty(output)

				currentTasks, err := getMockPersistedTasks()
				assert.NoError(err)

				assertTasksAreDeleted(
					oldPersistedTasks,
					currentTasks,
					fmt.Sprintf("Failed to delete tasks with priority: %s", priority),
				)

				noTasksWithPriority := lo.EveryBy(
					currentTasks,
					func(item mockPersistedTask) bool {
						return item.Priority != priority
					},
				)
				assert.True(
					noTasksWithPriority,
					fmt.Sprintf("Found tasks with priority %s after deletion", priority),
				)
			})
		})

		It("deletes all tasks with matching title", func() {
			randomTask, err := getRandomPersistedTask(oldPersistedTasks)
			assert.NoError(err)

			output, err := executeCommand(
				rootCmd,
				"delete",
				randomTask.Title,
				createFlag(TITLE),
			)
			assert.NoError(err)
			assert.NotEmpty(output)

			currentTasks, err := getMockPersistedTasks()
			assert.NoError(err)

			assertTasksAreDeleted(
				oldPersistedTasks,
				currentTasks,
				fmt.Sprintf("Failed to delete tasks with title: %s", randomTask.Title),
			)
		})

		It("deletes all completed tasks", func() {
			output, err := executeCommand(
				rootCmd,
				"delete",
				"complete",
				createFlag(COMPLETION),
			)
			assert.NoError(err)
			assert.NotEmpty(output)

			currentTasks, err := getMockPersistedTasks()
			assert.NoError(err)

			assertTasksAreDeleted(
				oldPersistedTasks,
				currentTasks,
				"No completed tasks were deleted",
			)
		})

		It("deletes all incomplete tasks", func() {
			output, err := executeCommand(
				rootCmd,
				"delete",
				"incomplete",
				createFlag(COMPLETION),
			)
			assert.NoError(err)
			assert.NotEmpty(output)

			currentTasks, err := getMockPersistedTasks()
			assert.NoError(err)

			assertTasksAreDeleted(
				oldPersistedTasks,
				currentTasks,
				"No incomplete tasks were deleted",
			)
		})
	})
})
