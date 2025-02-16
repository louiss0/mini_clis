package task

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
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
	Title, Description, id string
	createdAt              int64
	Priority               priority
	Complete               bool
	UpdatedAt              time.Time
}

func NewTask(title, description string) Task {

	return Task{
		Title:       title,
		Description: description,
		id:          uuid.NewString(),
		Priority:    LOW,
		createdAt:   time.Now().UnixMicro(),
		UpdatedAt:   time.Now(),
	}
}

func (self Task) CreatedAt() int64 {

	return self.createdAt
}

func (self Task) Id() string {

	return self.id
}

func (self Task) UpdatedAtTimeStamp() int64 {
	return self.UpdatedAt.UnixMicro()
}

func (self Task) ToJSON() (string, error) {

	byte, error := json.Marshal(persistedTask{
		Id:          self.id,
		Title:       self.Title,
		Description: self.Description,
		Priority:    self.Priority.Value(),
		Complete:    self.Complete,
		CreatedAt:   self.CreatedAt(),
		UpdatedAt:   self.UpdatedAtTimeStamp(),
	})

	return string(byte), error
}

type persistedTask struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Complete    bool   `json:"complete"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
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
				UpdatedAt:   item.UpdatedAtTimeStamp(),
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

		parsedPriority, _ := ParsePriority(item.Priority)

		return Task{
			Title:       item.Title,
			Description: item.Description,
			Priority:    parsedPriority,
			Complete:    item.Complete,
			UpdatedAt:   time.UnixMicro(item.UpdatedAt),
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
				UpdatedAt:   item.UpdatedAtTimeStamp(),
			}
		},
	)

	byte, error := json.Marshal(&persistedTasks)

	if error != nil {

		return "", error

	}

	return string(byte), nil

}
