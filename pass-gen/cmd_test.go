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
			cmd.CreateLeetspeakCmd(),
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

		It("generates words with a length of 5", func() {
			output, err := executeCommand(rootCmd, "words")
			assert.NoError(err)

			allWordsAreTheLengthOfFive := lo.EveryBy(
				strings.Split(output, "-"),
				func(word string) bool {

					return utf8.RuneCountInString(word) == 5
				})
			assert.True(allWordsAreTheLengthOfFive)
		})

		It("allows the user to specify the number of words using the count flag", func() {
			output, err := executeCommand(rootCmd, "words", "--count", "5")
			assert.NoError(err)
			assert.NotEmpty(output)

			thereAreFiveWords := len(strings.Split(output, "-")) == 5
			assert.True(thereAreFiveWords)

		})

		It("allows the user to specify the length of words using the length flag", func() {
			output, err := executeCommand(rootCmd, "words", "--length", "3")
			assert.NoError(err)
			assert.NotEmpty(output)

			allWordsAreTheLengthOfThree := lo.EveryBy(
				strings.Split(output, "-"),
				func(word string) bool {

					return utf8.RuneCountInString(word) == 3
				})
			assert.True(allWordsAreTheLengthOfThree)
		})

		Context("Proper Separators", func() {

			allowedSeparators := []string{"-", "_", "=", ":"}

			lo.ForEach(allowedSeparators, func(separator string, index int) {
				It("allows the user to specify the separator using the sep flag", func() {
					output, err := executeCommand(
						rootCmd,
						"words",
						"--separator",
						separator,
					)
					assert.NoError(err)
					assert.NotEmpty(output)

					allWordsAreSeparatedBySeperator := strings.Contains(output, separator)
					assert.True(allWordsAreSeparatedBySeperator)
				})

			})

		})

	})

	Context("Leet Speak", func() {

		leetSpeakLetterMap := map[rune][]string{
			'a': []string{"4", "@", "^", "/\\", "4"},
			'b': []string{"8", "|3", "ß", "13"},
			'c': []string{"(", "{", "[", "©"},
			'd': []string{"|)", "[)", "Ð", "6"},
			'e': []string{"3", "&", "€", "13"},
			'f': []string{"|=", "ƒ", "ph"},
			'g': []string{"6", "9", "&", "5"},
			'h': []string{"#", "[-]", "|-|", "4"},
			'i': []string{"1", "!", "|", "L"},
			'j': []string{"_|", "_/", "J"},
			'k': []string{"|<", "|{", "X", "K"},
			'l': []string{"1", "|", "I"},
			'm': []string{"|V|", "/\\/\\", "[V]", "M"},
			'n': []string{"^/", "/\\/", "Ñ", "n"},
			'o': []string{"0", "()", "[]", "<>", "O"},
			'p': []string{"|>", "9", "Þ", "P"},
			'q': []string{"9", "O_", "(,)", "Q"},
			'r': []string{"|2", "®", "/2", "R"},
			's': []string{"5", "$", "z", "§", "S"},
			't': []string{"7", "+", "†", "T"},
			'u': []string{"|_|", "[_]", "\\_/", "U"},
			'v': []string{"\\/", "√", "V", "V"},
			'w': []string{"\\/\\/", "VV", "µ", "W"},
			'x': []string{"%", "><", "*", "×"},
			'y': []string{"`/", "¥", "Y"},
			'z': []string{"2", "%", "7_", "Z"},
		}

		// leetSpeakNumberMap := map[rune][]string{
		// 	'0': []string{"O", "o"},
		// 	'1': []string{"I", "l", "!", "L"},
		// 	'2': []string{"Z", "z"},
		// 	'3': []string{"E", "e"},
		// 	'4': []string{"A", "a"},
		// 	'5': []string{"S", "s"},
		// 	'6': []string{"G", "g"},
		// 	'7': []string{"T", "t"},
		// 	'8': []string{"B", "b"},
		// 	'9': []string{"P", "p"},
		// }

		It("generates a leet speak password", func() {
			const word = "hello"
			output, err := executeCommand(rootCmd, "leetspeak", word)

			splitWord := strings.Split(word, "")

			possibleValues := lo.Reduce(
				splitWord,
				func(agg []string, item string, index int) []string {
					return append(agg, leetSpeakLetterMap[[]rune(item)[0]]...)
				},
				[]string{},
			)

			assert.NoError(err)
			assert.NotEmpty(output)

			allPossibleValuesAreInSplitWord := lo.Every(
				possibleValues,
				strings.Split(output, ""),
			)

			assert.True(allPossibleValuesAreInSplitWord)

		})

	})

})
