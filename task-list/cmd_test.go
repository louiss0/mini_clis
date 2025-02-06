package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"

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
	Priority    string `json:"string"`
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

		It("works", func() {

			output, error := executeCommand(rootCmd, "list")

			assert.NotEmpty(output)

			assert.Nil(error)

			var tasks []mockPersistedTask

			json.Unmarshal([]byte(output), &tasks)

			assert.Greater(len(tasks), 0)

		})
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

})
