package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use: "restore",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restore called")
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
