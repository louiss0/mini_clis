package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mini-clis/pass-gen/cmd"
	. "github.com/onsi/ginkgo/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(command *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	errorBuf := new(bytes.Buffer)
	command.SetOut(buf)
	command.SetErr(errorBuf)
	command.SetArgs(args)
	err = command.Execute()

	if errorBuf.Len() > 0 {
		err = fmt.Errorf("Command failed: %s", errorBuf.String())
	}

	return buf.String(), err
}

var _ = Describe("Cmd", func() {

	assert := assert.New(GinkgoT())

	var rootCmd = cmd.RootCmd()

	BeforeEach(func() {

		rootCmd.AddCommand(
			cmd.CreateWordsCmd(),
		)
	})

	AfterEach(func() {

		rootCmd.ResetCommands()
	})

	// This test only works when a string is passed as an argument
	It("should execute successfully", func() {
		output, err := executeCommand(rootCmd, "")
		assert.NoError(err)
		assert.NotEmpty(output)
	})

	Context("Words", func() {
		It("should execute successfully", func() {
			output, err := executeCommand(rootCmd, "words")
			assert.NoError(err)
			assert.NotEmpty(output)
		})

		It("generates three words by default separated by dashes", func() {
			output, err := executeCommand(rootCmd, "words")
			assert.NoError(err)
			assert.NotEmpty(output)
		})

		It("generates words with a length of 5", Focus, func() {
			output, err := executeCommand(rootCmd, "words")
			assert.NoError(err)

			allWordsAreTheLengthOfFive := lo.EveryBy(
				strings.Split(output, "-"),
				func(word string) bool {

					return utf8.RuneCountInString(word) == 5
				})
			assert.True(allWordsAreTheLengthOfFive)
		})

	})

})
