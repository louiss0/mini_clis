package printer

import (
	"fmt"

	"github.com/spf13/cobra"
)

func PrintUsingCommmand(cmd *cobra.Command, string string) error {

	_, err := fmt.Fprint(cmd.OutOrStdout(), string)

	return err
}
