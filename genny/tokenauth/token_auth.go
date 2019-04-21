package tokenauth

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gobuffalo/gogen"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/pkg/errors"
)

var fields attrs.Attrs

func extraAttrs(args []string) []string {
	var names = map[string]string{
		"email":    "email",
		"password": "password",
		"id":       "id",
	}

	var result = []string{}
	for _, field := range args {
		attr, _ := attrs.Parse(field)
		field = attr.Name.Underscore().String()

		if names[field] != "" {
			continue
		}

		names[field] = field
		result = append(result, field)
	}

	return result
}

// New
func New(args *Options) (*genny.Generator, error) {
	g := genny.New()

	var err error
	fields, err = attrs.ParseArgs(extraAttrs(args.UserFields)...)
	if err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.New("templates", "../tokenauth/templates")); err != nil {
		return g, errors.WithStack(err)
	}

	app := meta.New(".")
	ctx := plush.NewContext()
	ctx.Set("app", app)
	ctx.Set("fields", fields)
	ctx.Set("token_prefix", args.Prefix)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.NewTransformer(".fizz", migrationsTransformer(time.Now())))

	g.RunFn(func(r *genny.Runner) error {
		path := filepath.Join("actions", "app.go")
		gf, err := r.FindFile(path)
		if err != nil {
			return err
		}

		gf, err = gogen.AddInsideBlock(
			gf,
			`if app == nil {`,
			`app.Use(middleware.TokenMiddleware)`,
			`app.POST("/users", UsersCreate)`,
			`app.PUT("/users/{user_id}/", UsersUpdate)`,
			`app.DELETE("/users/{user_id}/", UsersDestroy)`,
			`app.POST("/signin", AuthCreate)`,
			`app.DELETE("/signout", AuthDestroy)`,
			`app.Middleware.Skip(middleware.TokenMiddleware, HomeHandler, UsersCreate, UsersLogin)`,
		)
		return r.File(gf)
	})

	return g, nil
}

func migrationsTransformer(t time.Time) genny.TransformerFn {
	v := t.UTC().Format("20060102150405")
	return func(f genny.File) (genny.File, error) {
		p := filepath.Base(f.Name())
		return genny.NewFile(filepath.Join("migrations", fmt.Sprintf("%s_%s", v, p)), f), nil
	}
}
