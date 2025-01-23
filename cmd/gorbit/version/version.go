package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.0-beta"

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gorbit CLI",
	Long:  `Display the current version of the Gorbit command-line tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Gorbit CLI Version: %s\n", version)
	},
}
