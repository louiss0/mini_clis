package task

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/samber/lo"
)

type Task struct {
	Title,
	Description,
	id,
	createdAt string
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

type persistentTask struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
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
		func(item Task, index int) persistentTask {

			return persistentTask{
				Id:          item.id,
				Title:       item.Title,
				Description: item.Description,
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

func ReadTasks() ([]persistentTask, error) {

	var persistentTasks []persistentTask

	taskListJSONFilePath, error := getTaskListJSONFilePath()

	if error != nil {

		return persistentTasks, error

	}

	byte, error := os.ReadFile(taskListJSONFilePath)

	if error != nil {

		return persistentTasks, error

	}

	unmarshalError := json.Unmarshal(byte, &persistentTasks)

	if unmarshalError != nil {
		return persistentTasks, unmarshalError
	}
	return persistentTasks, unmarshalError

}
