package gorbit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gorbit",
	Short: "Gorbit CLI - The official command-line tool for Gorbit framework",
	Long: `Gorbit CLI provides tools to create, manage, and deploy 
Gorbit-based applications. Complete documentation is available at 
https://gorbit.io/docs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Gorbit CLI!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
