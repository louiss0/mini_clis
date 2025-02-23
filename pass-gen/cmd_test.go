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
			'a': {"4", "@", "^", "/\\", "4"},
			'b': {"8", "|3", "ß", "13"},
			'c': {"(", "{", "[", "©"},
			'd': {"|)", "[)", "Ð", "6"},
			'e': {"3", "&", "€", "13"},
			'f': {"|=", "ƒ", "ph"},
			'g': {"6", "9", "&", "5"},
			'h': {"#", "[-]", "|-|", "4"},
			'i': {"1", "!", "|", "L"},
			'j': {"_|", "_/", "J"},
			'k': {"|<", "|{", "X", "K"},
			'l': {"1", "|", "I"},
			'm': {"|V|", "/\\/\\", "[V]", "M"},
			'n': {"^/", "/\\/", "Ñ", "n"},
			'o': {"0", "()", "[]", "<>", "O"},
			'p': {"|>", "9", "Þ", "P"},
			'q': {"9", "O_", "(,)", "Q"},
			'r': {"|2", "®", "/2", "R"},
			's': {"5", "$", "z", "§", "S"},
			't': {"7", "+", "†", "T"},
			'u': {"|_|", "[_]", "\\_/", "U"},
			'v': {"\\/", "√", "V", "V"},
			'w': {"\\/\\/", "VV", "µ", "W"},
			'x': {"%", "><", "*", "×"},
			'y': {"`/", "¥", "Y"},
			'z': {"2", "%", "7_", "Z"},
		}

		// leetSpeakNumberMap := map[rune]{
		// 	'0': {"O", "o"},
		// 	'1': {"I", "l", "!", "L"},
		// 	'2': {"Z", "z"},
		// 	'3': {"E", "e"},
		// 	'4': {"A", "a"},
		// 	'5': {"S", "s"},
		// 	'6': {"G", "g"},
		// 	'7': {"T", "t"},
		// 	'8': {"B", "b"},
		// 	'9': {"P", "p"},
		// }

		assertEveryCharacterIsALeetSpeakVersion := func(
			leetSpeakMap map[rune][]string,
			input,
			output string,
		) {
			everyCharacterIsALeetSpeakVersion := lo.EveryBy(
				strings.Split(input, ""),
				func(character string) bool {

					symbolsForCharacter := leetSpeakMap[rune(character[0])]

					return lo.SomeBy(
						symbolsForCharacter,
						func(symbol string) bool {
							return strings.Contains(output, symbol)
						},
					)

				})

			assert.True(everyCharacterIsALeetSpeakVersion)
		}

		It("generates a leet speak password", func() {
			const word = "hello"
			output, err := executeCommand(rootCmd, "leetspeak", word)

			assert.NoError(err)
			assert.NotEmpty(output)

			assertEveryCharacterIsALeetSpeakVersion(leetSpeakLetterMap, word, output)

		})

	})

})
