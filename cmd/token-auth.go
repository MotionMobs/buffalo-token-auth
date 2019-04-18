package cmd

import (
	"github.com/spf13/cobra"
)

// token-authCmd represents the buffalo-token-auth command
var tokenauthCmd = &cobra.Command{
	Use:   "buffalo-token-auth",
	Short: "tools for working with token-auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenauthCmd)
}
