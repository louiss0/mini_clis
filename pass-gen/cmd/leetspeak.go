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
	"math/rand"
	"strings"

	"github.com/mini-clis/pass-gen/printer"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// CreateLeetspeakCmd returns the leetspeak command
func CreateLeetspeakCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "leetspeak",
		Short: "Generate leetspeak password",
		Args:  cobra.ExactArgs(1),
		Long: `Generate a leetspeak password using a combination of letters and numbers.
	A leetspeak password is a password that uses your input to create passwords.With similar symbols`,
		Run: func(cmd *cobra.Command, args []string) {

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

			// leetSpeakNumberMap := map[rune][]string{
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

			leetSpeakSlice := lo.Map(
				strings.Split(args[0], ""),
				func(item string, index int) string {

					key := []rune(item)[0]

					symbols, ok := leetSpeakLetterMap[key]

					if ok {
						return symbols[rand.Intn(len(symbols))]
					}
					return item

				})

			printer.PrintUsingCommmand(cmd, strings.Join(leetSpeakSlice, ""))
		},
	}
}

func init() {
	rootCmd.AddCommand(CreateLeetspeakCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// leetspeakCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// leetspeakCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
