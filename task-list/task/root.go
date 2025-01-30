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

func (self Task) ToJSON() string {

	byte, _ := json.Marshal(
		struct {
			Id          string `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		}{
			Id:          self.id,
			Title:       self.Title,
			Description: self.Description,
			CreatedAt:   self.createdAt,
			UpdatedAt:   self.UpdatedAtAsUnixDateFormat(),
		})

	return string(byte)

}

func SaveTasks(tasks []Task) error {

	dir, error := os.Getwd()

	if error != nil {

		return error
	}

	byte, error := json.Marshal(lo.Map(
		tasks,
		func(item Task, index int) string {
			return item.ToJSON()
		},
	),
	)

	if error != nil {

		return error

	}

	return os.WriteFile(fmt.Sprintf("%s/%s", dir, "task-list.json"), byte, os.ModeDevice)

}

func ReadTasks() {}
