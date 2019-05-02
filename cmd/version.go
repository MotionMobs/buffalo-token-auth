package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/MotionMobs/buffalo-token-auth/tokenauth"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of token-auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("token-auth", tokenauth.Version)
		return nil
	},
}

func init() {
	tokenauthCmd.AddCommand(versionCmd)
}
