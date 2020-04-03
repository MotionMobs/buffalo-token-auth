package tokenauth

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	g := genny.New()
	appPath := filepath.Join("actions", "app.go")
	g.File(genny.NewFileS(appPath, (appBefore)))
	run.With(g)

	gg, err := New(&Options{})
	r.NoError(err)
	run.With(gg)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 8)

	f := res.Files[0]
	r.Equal("actions/app.go", f.Name())
	r.Equal(appAfter, f.String())

	f = res.Files[1]
	r.Equal("actions/auth.go", f.Name())

	f = res.Files[2]
	r.Equal("actions/users.go", f.Name())

	f = res.Files[3]
	r.Equal("middleware/token_middleware.go", f.Name())

	f = res.Files[4]
	r.Contains(f.Name(), "_create_users.down.fizz")

	f = res.Files[5]
	r.Contains(f.Name(), "_create_users.up.fizz")

	f = res.Files[6]
	r.Equal("models/user.go", f.Name())

	f = res.Files[7]
	r.Equal("models/user_test.go", f.Name())
}

const appBefore = `package actions
import (
	"github.com/gobuffalo/buffalo"
)
var app *buffalo.App
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{})
	}
	return app
}`

const appAfter = `package actions
import (
	"github.com/gobuffalo/buffalo"
)
var app *buffalo.App
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{})
		app.Use(middleware.TokenMiddleware)
		app.POST("/users", UsersCreate)
		app.PUT("/users/{user_id}/", UsersUpdate)
		app.DELETE("/users/{user_id}/", UsersDestroy)
		app.POST("/signin", AuthCreate)
		app.DELETE("/signout", AuthDestroy)
		app.Middleware.Skip(middleware.TokenMiddleware, HomeHandler, UsersCreate, AuthCreate)
	}
	return app
}`
