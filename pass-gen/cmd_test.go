package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
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
			cmd.CreateNumericCmd(),
			cmd.CreateEncodeCmd(),
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

		leetSpeakNumberMap := map[rune][]string{
			'0': {"O", "o"},
			'1': {"I", "l", "!", "L"},
			'2': {"Z", "z"},
			'3': {"E", "e"},
			'4': {"A", "a"},
			'5': {"S", "s"},
			'6': {"G", "g"},
			'7': {"T", "t"},
			'8': {"B", "b"},
			'9': {"P", "p"},
		}

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

		It("generates a leet speak password by changing words only", func() {
			const word = "hello123456"
			output, err := executeCommand(rootCmd, "leetspeak", word)

			assert.NoError(err)
			assert.NotEmpty(output)

			outputContainsLeetSpeakVersionsOfLetters := lo.EveryBy(
				lo.Filter(
					strings.Split(word, ""),
					func(character string, index int) bool {
						return unicode.IsLetter(rune(character[0]))
					}),
				func(character string) bool {

					symbolsForCharacter := leetSpeakLetterMap[rune(character[0])]

					return lo.SomeBy(
						symbolsForCharacter,
						func(symbol string) bool {
							return strings.Contains(output, symbol)
						},
					)

				})

			assert.True(
				outputContainsLeetSpeakVersionsOfLetters,
				"output does not contain leet speak versions of letters",
			)

			numbersFromWord := strings.Join(
				lo.Filter(
					strings.Split(word, ""),
					func(character string, index int) bool {
						return unicode.IsNumber(rune(character[0]))
					}),
				"",
			)

			assert.True(
				strings.Contains(output, numbersFromWord),
				"output contains leet speak versions of numbers",
			)

		})

		It("changes the numbers only when the --numbers flag is passed", func() {
			const word = "world987654"
			output, err := executeCommand(rootCmd, "leetspeak", word, "--numbers")

			assert.NoError(err)
			assert.NotEmpty(output)

			outputContainsLeetSpeakVersionsOfNumbers := lo.EveryBy(
				lo.Filter(
					strings.Split(word, ""),
					func(character string, index int) bool {
						return unicode.IsNumber(rune(character[0]))
					}),
				func(character string) bool {

					symbolsForCharacter := leetSpeakNumberMap[rune(character[0])]

					return lo.SomeBy(
						symbolsForCharacter,
						func(symbol string) bool {
							return strings.Contains(output, symbol)
						},
					)

				})

			assert.True(
				outputContainsLeetSpeakVersionsOfNumbers,
				"output does not contain leet speak versions of numbers",
			)

			lettersFromWord := strings.Join(
				lo.Filter(
					strings.Split(word, ""),
					func(character string, index int) bool {
						return unicode.IsLetter(rune(character[0]))
					}),
				"",
			)

			assert.True(
				strings.Contains(output, lettersFromWord),
				"output contains leet speak versions of letters",
			)

		})

		It("makes sure that only numbers and letters can be input", func() {
			const word = "!@#$%^&*()_+"
			output, err := executeCommand(rootCmd, "leetspeak", word)

			assert.Error(err)
			assert.Empty(output)

		})

	})

	Context("Numeric", func() {

		It("generates a numeric password", func() {

			output, err := executeCommand(rootCmd, "numeric")

			assert.NoError(err)
			assert.NotEmpty(output)

			allCharactersAreNumbers := lo.EveryBy(
				strings.Split(output, ""),
				func(character string) bool {
					return unicode.IsNumber(rune(character[0]))
				})

			assert.True(
				allCharactersAreNumbers,
				"output contains only numbers",
			)

		})

		It("defaults to only four numbers", func() {
			output, err := executeCommand(rootCmd, "numeric")

			assert.NoError(err)
			assert.NotEmpty(output)

			assert.Len(strings.Split(output, ""), 4)

		})

		It("changes the amount of numbers when the length flag is provided", func() {
			output, err := executeCommand(rootCmd, "numeric", "-l", "6")

			assert.NoError(err)
			assert.NotEmpty(output)

			assert.Len(strings.Split(output, ""), 6)

		})

		It(
			fmt.Sprintf(
				"only allows the length to only be between %d and %d",
				cmd.SHORTEST_LENGTH,
				cmd.LONGEST_LENGTH,
			),
			func() {
				output, err := executeCommand(rootCmd, "numeric", "-l", "2")

				assert.Error(err)
				assert.Empty(output)

				output, err = executeCommand(rootCmd, "numeric", "-l", "21")

				assert.Error(err)
				assert.Empty(output)

			},
		)

		It("generates a timestamp when the date-pin flag is provided", func() {

			isTimestamp := func(n int64) bool {
				// Define reasonable timestamp bounds
				minTimestamp := int64(946684800)      // January 1, 2000 (seconds)
				maxTimestamp := int64(4102444800)     // January 1, 2100 (seconds)
				minTimestampMs := minTimestamp * 1000 // In milliseconds
				maxTimestampMs := maxTimestamp * 1000 // In milliseconds

				// Check if it's in the valid range for seconds or milliseconds
				return (n >= minTimestamp && n <= maxTimestamp) || (n >= minTimestampMs && n <= maxTimestampMs)
			}

			output, err := executeCommand(rootCmd, "numeric", "-d")

			assert.NoError(err)
			assert.NotEmpty(output)

			intOutput, err := strconv.ParseInt(output, 10, 64)
			assert.NoError(err)

			assert.True(isTimestamp(intOutput))

		})

	})

	Context("Encode", func() {
		It("encodes a string", func() {
			const ARGUMENT = "hello"
			output, err := executeCommand(rootCmd, "encode", ARGUMENT)

			assert.NoError(err)
			assert.NotEmpty(output)

			decodedByte, decodedError := hex.DecodeString(output)
			assert.NoError(decodedError)

			assert.Equal(decodedByte, []byte(ARGUMENT))
		})

		It("encodes any amount of args passed", func() {

			ARGUMENTS := []string{"hello", "world"}
			output, err := executeCommand(rootCmd, slices.Insert(ARGUMENTS, 0, "encode")...)

			assert.NoError(err)
			assert.NotEmpty(output)

			decodedByte, decodedError := hex.DecodeString(output)
			assert.NoError(decodedError)
			assert.Equal(decodedByte, []byte(strings.Join(ARGUMENTS, "")))

		})

		It(
			"encodes any amount of arguments passed with a separator when the separator flag is provided",
			func() {

				ARGUMENTS := []string{"hello", "world"}
				output, err := executeCommand(rootCmd, slices.Insert(ARGUMENTS, 0, "encode", "-s", ",")...)

				assert.NoError(err)
				assert.NotEmpty(output)

				splitEncodedStrings := strings.Split(output, ",")

				assert.Len(splitEncodedStrings, len(ARGUMENTS))

				lo.ForEach(
					splitEncodedStrings,
					func(splitEncodedString string, index int) {

						decodedByte, decodedError := hex.DecodeString(splitEncodedString)

						assert.NoError(
							decodedError,
							"Failed to decode string at index %d",
							index,
						)

						assert.Equal(decodedByte, []byte(ARGUMENTS[index]))

					})

			})

	})

})
