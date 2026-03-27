package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the TaskFix version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("taskfix v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
