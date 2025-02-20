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
	"strings"

	"math/rand"

	"github.com/mini-clis/pass-gen/printer"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// wordsCmd represents the words command
func CreateWordsCmd() *cobra.Command {
	var wordsCmd = &cobra.Command{
		Use:   "words",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			amountOfWords := 3

			wordLength := 5

			allLetters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

			allNumbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

			keyboardSymbols := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "+", "=", "{", "}", "[", "]", "|", ":", ";", ".", "/", "?"}

			numberToCharsMap := map[int][]string{
				0: allLetters,
				1: allNumbers,
				2: keyboardSymbols,
				3: lo.Map(allLetters, func(item string, index int) string {
					return strings.ToUpper(item)
				}),
			}

			values :=
				lo.Map(
					lo.Range(amountOfWords),
					func(item int, index int) string {

						return strings.Join(
							lo.Map(lo.Range(wordLength),
								func(item int, index int) string {

									randomNumberFromZeroToTwo := rand.Intn(len(numberToCharsMap))

									randomIntFromCharSetlength := rand.Intn(len(numberToCharsMap[randomNumberFromZeroToTwo]))

									return numberToCharsMap[randomNumberFromZeroToTwo][randomIntFromCharSetlength]
								}),
							"")

					},
				)

			return printer.PrintUsingCommmand(cmd, strings.Join(values, "-"))
		},
	}

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
