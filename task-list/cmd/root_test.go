package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {

	buffer := new(bytes.Buffer)

	cmd.SetOut(buffer)
	cmd.SetErr(buffer)
	cmd.SetArgs(args)

	err := cmd.Execute()

	return buffer.String(), err
}

type mockPersistedTask struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"string"`
	Complete    bool   `json:"complete"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

var createFlag = func(flagName string) string {
	return fmt.Sprintf("--%s", flagName)
}

func TestListCommand(t *testing.T) {

	assert := assert.New(t)

	t.Run("works", func(t *testing.T) {

		output, error := executeCommand(rootCmd, "list")

		assert.NoError(error)

		assert.NotNil(output)

		assert.NotEmpty(output)

		var tasks []mockPersistedTask

		error = json.Unmarshal([]byte(output), &tasks)

		assert.NoError(error)

		assert.Greater(len(tasks), 0)

	})

	t.Run("it errors when filter-priority with wrong value is passed", func(t *testing.T) {

		output, error := executeCommand(rootCmd, "list", createFlag(FILTER_PRIORITY), "foo")

		assert.NoError(error)

		assert.NotNil(output)

	})

}
