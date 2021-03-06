package cmd

import (
	"context"

	"github.com/MotionMobs/buffalo-token-auth/genny/tokenauth"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateOptions = struct {
	*tokenauth.Options
	dryRun bool
}{
	Options: &tokenauth.Options{},
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "token-auth",
	Short: "generates a new token-auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if generateOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		g, err := tokenauth.New(generateOptions.Options)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		g, err = gogen.Fmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&generateOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	generateCmd.Flags().StringSliceVarP(&generateOptions.UserFields, "user-fields", "u", []string{}, "comma separated list of fields to add to the user model, other than email, password_hash, refresh_token, and id")
	generateCmd.Flags().StringVarP(&generateOptions.Prefix, "token-prefix", "p", "", "prefix for use in the JWT generated by the middleware")
	tokenauthCmd.AddCommand(generateCmd)
}
