package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// bloquearCmd represents the bloquear command
var bloquearCmd = &cobra.Command{
	Use:   "bloquear",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bloquear called")
	},
}

func init() {
	rootCmd.AddCommand(bloquearCmd)
}
