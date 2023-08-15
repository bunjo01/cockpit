package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete")
	},
}
