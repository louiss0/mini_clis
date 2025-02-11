package task

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mini-clis/task-list/custom_errors"
	"github.com/samber/lo"
)

type priority string

const HIGH = priority("high")
const MEDIUM = priority("medium")
const LOW = priority("low")

func (self priority) Order() int {

	return map[priority]int{
		HIGH:   3,
		MEDIUM: 2,
		LOW:    1,
	}[self]

}

func (self priority) Value() string {

	return string(self)
}

var AllowedProrities = []string{
	string(LOW),
	string(HIGH),
	string(MEDIUM),
}

func ParsePriority(input string) (priority, error) {

	if !lo.Contains(AllowedProrities, input) {

		return "", custom_errors.
			CreateInvalidFlagErrorWithMessage(
				fmt.Sprintf(
					"Wrong option %s a priority is supposed to be %s",
					input,
					strings.Join(AllowedProrities, ","),
				),
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
		createdAt:   time.Now().Format(time.DateTime),
		UpdatedAt:   time.Now(),
	}
}

func (self Task) CreatedAt() string {

	return self.createdAt
}

func (self Task) Id() string {

	return self.id
}

func (self Task) UpdatedAtDateString() string {
	return self.UpdatedAt.Format(time.DateTime)
}

func (self Task) ToJSON() (string, error) {

	byte, error := json.Marshal(persistedTask{
		Id:          self.id,
		Title:       self.Title,
		Description: self.Description,
		Priority:    self.Priority.Value(),
		Complete:    self.Complete,
		CreatedAt:   self.CreatedAt(),
		UpdatedAt:   self.UpdatedAtDateString(),
	})

	return string(byte), error
}

type persistedTask struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Complete    bool   `json:"complete"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const TASK_LIST_STORAGE_PATH = "/home/shelton-louis/Desktop/cli-projects/mini-clis/task-list/task-list.json"

func SaveTasks(tasks []Task) error {

	byte, error := json.Marshal(lo.Map(
		tasks,
		func(item Task, index int) persistedTask {

			return persistedTask{
				Id:          item.id,
				Title:       item.Title,
				Description: item.Description,
				Priority:    item.Priority.Value(),
				Complete:    item.Complete,
				CreatedAt:   item.createdAt,
				UpdatedAt:   item.UpdatedAtDateString(),
			}
		},
	),
	)

	if error != nil {

		return error

	}

	return os.WriteFile(TASK_LIST_STORAGE_PATH, byte, os.ModeDevice)

}

func ReadTasks() ([]Task, error) {

	var tasks []Task

	byte, error := os.ReadFile(TASK_LIST_STORAGE_PATH)

	if error != nil {

		return tasks, error

	}

	var persistedTasks []persistedTask

	unmarshalError := json.Unmarshal(byte, &persistedTasks)

	if unmarshalError != nil {
		return tasks, unmarshalError
	}

	tasks = lo.Map(persistedTasks, func(item persistedTask, index int) Task {

		updatedAtTime, _ := time.Parse(time.DateTime, item.UpdatedAt)

		parsedPriority, _ := ParsePriority(item.Priority)

		return Task{
			Title:       item.Title,
			Description: item.Description,
			Priority:    parsedPriority,
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
				Priority:    item.Priority.Value(),
				Complete:    item.Complete,
				CreatedAt:   item.createdAt,
				UpdatedAt:   item.UpdatedAtDateString(),
			}
		},
	)

	byte, error := json.Marshal(&persistedTasks)

	if error != nil {

		return "", error

	}

	return string(byte), nil

}
