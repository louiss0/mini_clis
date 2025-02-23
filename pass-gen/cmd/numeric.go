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
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/mini-clis/pass-gen/printer"
	"github.com/spf13/cobra"
)

const (
	SHORTEST_LENGTH = 3
	DEFAULT_LENGTH  = 4
	LONGEST_LENGTH  = 20
)

const LENGTH = "length"

func CreateNumericCmd() *cobra.Command {
	numericCmd := &cobra.Command{
		Use:   "numeric",
		Short: "Generate a random numeric string",
		Long: `Generate a random numeric string of a specified length.
	You can specify the length of the numeric string using the --length flag.
	You get 4 digits by default.
	`,
		Run: func(cmd *cobra.Command, args []string) {

			generateSecureNDigitNumber := func(n int) int {

				// intPow computes integer exponentiation (10^n)
				intPow := func(base, exp int) int {
					result := 1
					for i := 0; i < exp; i++ {
						result *= base
					}
					return result
				}

				if n <= 0 {
					return 0
				}

				// Define the range for n-digit numbers
				min := intPow(10, n-1)   // Smallest n-digit number
				max := intPow(10, n) - 1 // Largest n-digit number

				// Generate a cryptographically secure random number
				num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
				return int(num.Int64()) + min
			}

			length, _ := cmd.Flags().GetInt(LENGTH)

			printer.PrintUsingCommmand(cmd, fmt.Sprintf("%d", generateSecureNDigitNumber(length)))

		},
	}

	numericCmd.Flags().IntP(LENGTH, "l", DEFAULT_LENGTH, "Length of the numeric string")
	return numericCmd
}

func init() {
	rootCmd.AddCommand(CreateNumericCmd())
}
