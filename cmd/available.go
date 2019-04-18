package cmd

import (
	"github.com/gobuffalo/buffalo-plugins/plugins/plugcmds"
	"mmgitl.mattclark.guru/MM/buffalo-token-auth/tokenauth"
)

var Available = plugcmds.NewAvailable()

func init() {
	Available.Add("root", tokenauthCmd)
	Available.Listen(tokenauth.Listen)
	Available.Add("generate", generateCmd)
	Available.Mount(rootCmd)
}
