/*
Copyright © 2025 Shelton Louis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/mini-clis/pass-gen/printer"
	"github.com/mini-clis/shared/custom_errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type SeparatorFlag struct {
	value    string
	flagName string
}

func (self SeparatorFlag) String() string {
	return self.value
}

func (self *SeparatorFlag) Set(value string) error {
	allowedSeparators := []string{"-", "_", "=", ":"}

	if !lo.Contains(allowedSeparators, value) {
		return custom_errors.CreateInvalidFlagErrorWithMessage(
			custom_errors.FlagName(self.flagName),
			fmt.Sprintf("must be one of %v", allowedSeparators),
		)
	}

	self.value = value
	return nil
}

func (self SeparatorFlag) Type() string {
	return "string"
}

type ZeroCheckFlag struct {
	value    int
	flagName string
}

func (self *ZeroCheckFlag) Set(value string) error {

	int, err := strconv.Atoi(value)

	if err != nil {
		return err
	}

	if int <= 0 {
		return custom_errors.CreateInvalidFlagErrorWithMessage(custom_errors.FlagName("count"), "word length must be greater than 0")
	}

	self.value = int
	return nil

}

func (self ZeroCheckFlag) Value() int {
	return self.value
}

func (self ZeroCheckFlag) String() string {
	return fmt.Sprintf("%v", self.value)
}

func (self ZeroCheckFlag) Type() string {
	return "number"
}

// wordsCmd represents the words command
func CreateWordsCmd() *cobra.Command {

	amountOfWordsFlag := ZeroCheckFlag{
		value:    3,
		flagName: "count",
	}

	wordLengthFlag := ZeroCheckFlag{
		value:    5,
		flagName: "length",
	}

	separatorFlag := SeparatorFlag{
		value:    "-",
		flagName: "separator",
	}

	var wordsCmd = &cobra.Command{
		Use:   "words",
		Short: "Generate random words",
		Long: `Generate random words. This command generates a specified number of random words, each with a specified length, separated by a specified separator.
		The words can be generated with different character sets, such as letters, numbers, and symbols.
		The default amount of words is 3. The default word length is 5.
		The default separator is "-".
		You can specify a custom separator using the --separator flag.
		You can specify a custom length using the --length flag.
		You can specify how many words using the --count flag.
		`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			allLetters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

			allNumbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

			keyboardSymbols := []string{"!", "#", "$", "%", "^", "&", ",", ":", ".", "~", "@", "|", "<", ">", "?", "*", ";", "[", "]"}

			allLettersCapitalized := lo.Map(
				allLetters,
				func(item string, index int) string {
					return strings.ToUpper(item)
				})

			numberToCharsMap := map[int][]string{
				0: allLetters,
				1: allNumbers,
				2: keyboardSymbols,
				3: allLettersCapitalized,
			}

			values :=
				lo.Map(
					lo.Range(amountOfWordsFlag.Value()),
					func(outerItem int, index int) string {

						return strings.Join(
							lo.Map(lo.Range(wordLengthFlag.Value()),
								func(innerItem int, index int) string {

									randomNumberFromZeroToTwo := rand.Intn(len(numberToCharsMap))

									charSet := numberToCharsMap[randomNumberFromZeroToTwo]

									randomIntFromCharSetLength := rand.Intn(len(charSet))

									return charSet[randomIntFromCharSetLength]
								}),
							"")

					},
				)

			printer.PrintUsingCommmand(cmd, strings.Join(values, separatorFlag.String()))
		},
	}

	wordsCmd.Flags().Var(&amountOfWordsFlag, "count", "How many words are created ")
	wordsCmd.Flags().Var(&wordLengthFlag, "length", "How many characters are in each word")
	wordsCmd.Flags().Var(&separatorFlag, "separator", "The separator to use between words")

	return wordsCmd
}
func init() {
	rootCmd.AddCommand(CreateWordsCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wordsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wordsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
