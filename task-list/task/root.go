package task

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/tidwall/pretty"
)

type priority string

const HIGH = priority("high")
const MEDIUM = priority("medium")
const LOW = priority("low")

func ParsePriority(input string) (priority, error) {

	allowedProrities := []string{
		string(LOW),
		string(HIGH),
		string(MEDIUM),
	}
	if !lo.Contains(allowedProrities, input) {

		return "", fmt.Errorf(
			"Wrong option %s a priority is supposed to be %s",
			input,
			strings.Join(allowedProrities, ","),
		)
	}

	return priority(input), nil

}

type Task struct {
	Title,
	Description,
	id,
	createdAt string
	Priority  priority
	Complete  bool
	UpdatedAt time.Time
}

func NewTask(title, description string) Task {

	// generateRandomString creates a random string of given length
	generateRandomString := func(length int) string {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		rand.New(rand.NewSource(int64(length))) // Seed the random number generator
		b := make([]byte, length)

		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}

	return Task{
		Title:       title,
		Description: description,
		id:          generateRandomString(12),
		Priority:    LOW,
		createdAt:   time.Now().Format(time.UnixDate),
		UpdatedAt:   time.Now(),
	}
}

func (self Task) CreatedAt() string {

	return self.createdAt
}

func (self Task) Id() string {

	return self.id
}

func (self Task) UpdatedAtAsUnixDateFormat() string {
	return self.UpdatedAt.Format(time.UnixDate)
}

type persistedTask struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    priority `json:"string"`
	Complete    bool     `json:"complete"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

func getTaskListJSONFilePath() (string, error) {

	dir, error := os.Getwd()

	if error != nil {

		return "", error
	}

	return fmt.Sprintf("%s/%s", dir, "task-list.json"), nil

}

func SaveTasks(tasks []Task) error {

	taskListJSONFilePath, error := getTaskListJSONFilePath()

	if error != nil {

		return error
	}

	byte, error := json.Marshal(lo.Map(
		tasks,
		func(item Task, index int) persistedTask {

			return persistedTask{
				Id:          item.id,
				Title:       item.Title,
				Description: item.Description,
				Priority:    item.Priority,
				Complete:    item.Complete,
				CreatedAt:   item.createdAt,
				UpdatedAt:   item.UpdatedAtAsUnixDateFormat(),
			}
		},
	),
	)

	if error != nil {

		return error

	}

	return os.WriteFile(taskListJSONFilePath, byte, os.ModeDevice)

}

func ReadTasks() ([]Task, error) {

	var tasks []Task

	taskListJSONFilePath, error := getTaskListJSONFilePath()

	if error != nil {

		return tasks, error

	}

	byte, error := os.ReadFile(taskListJSONFilePath)

	if error != nil {

		return tasks, error

	}

	var persistedTasks []persistedTask

	unmarshalError := json.Unmarshal(byte, &persistedTasks)

	if unmarshalError != nil {
		return tasks, unmarshalError
	}

	tasks = lo.Map(persistedTasks, func(item persistedTask, index int) Task {

		updatedAtTime, _ := time.Parse(time.UnixDate, item.UpdatedAt)

		return Task{
			Title:       item.Title,
			Description: item.Description,
			Priority:    item.Priority,
			Complete:    item.Complete,
			UpdatedAt:   updatedAtTime,
			createdAt:   item.CreatedAt,
			id:          item.Id,
		}

	})

	return tasks, nil

}

func MarshallTasks(tasks []Task) (string, error) {

	persistedTasks := lo.Map(
		tasks,
		func(item Task, index int) persistedTask {

			return persistedTask{
				Id:          item.id,
				Title:       item.Title,
				Description: item.Description,
				Priority:    item.Priority,
				Complete:    item.Complete,
				CreatedAt:   item.createdAt,
				UpdatedAt:   item.UpdatedAtAsUnixDateFormat(),
			}
		},
	)

	byte, error := json.Marshal(&persistedTasks)

	if error != nil {

		return "", error

	}

	return string(
		pretty.Color(
			pretty.PrettyOptions(byte, nil),
			nil,
		),
	), nil

}
