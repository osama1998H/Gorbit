package gorbit

import (
	"gorbit/cmd/gorbit/version"
)

func main() {
	rootCmd.AddCommand(version.Cmd)
	Execute()
}
