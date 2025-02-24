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
	"encoding/hex"
	"strings"

	"github.com/mini-clis/pass-gen/printer"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const SEPERATOR = "separator"

// CreateEncodeCmd creates and returns the encode command
func CreateEncodeCmd() *cobra.Command {
	encodeCmd := &cobra.Command{
		Use:   "encode",
		Short: "Encode a string",
		Long: `Encode a series of characters into a base64 encoded string.
		You can only use letters, numbers, and symbols.
		Spaces are not allowed.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			separator, _ := cmd.Flags().GetString(SEPERATOR)

			hexcodedArgs := lo.Map(args, func(arg string, _ int) string {
				return hex.EncodeToString([]byte(arg))
			})

			printer.PrintUsingCommmand(cmd, strings.Join(hexcodedArgs, separator))

		},
	}

	encodeCmd.Flags().StringP(SEPERATOR, "s", "", "Separator to use between encoded arguments")

	return encodeCmd
}

func init() {
	rootCmd.AddCommand(CreateEncodeCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
