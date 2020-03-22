package tokenauth

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gobuffalo/genny/v2/gogen"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/pkg/errors"
)

var fields attrs.Attrs

func extraAttrs(args []string) attrs.Attrs {
	var extraAttrs attrs.Attrs

	var names = map[string]string{
		"id":               "id",
		"created_at":       "created_at",
		"updated_at":       "updated_at",
		"email":            "email",
		"password_hash":    "password_hash",
		"refresh_token":    "refresh_token",
		"password":         "password",
		"password_confirm": "password_confirm",
	}

	for _, field := range args {
		attr, _ := attrs.Parse(field)

		field = attr.Name.Underscore().String()

		if names[field] != "" {
			continue
		}

		names[field] = field
		extraAttrs = append(extraAttrs, attr)
	}

	return extraAttrs
}

// New generates and modifies files via plush and genny
func New(args *Options) (*genny.Generator, error) {
	g := genny.New()

	var err error
	fields = extraAttrs(args.UserFields)
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

	g.Transformer(plushTransformer(ctx))
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
			`app.Middleware.Skip(middleware.TokenMiddleware, HomeHandler, UsersCreate, AuthCreate)`,
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

// Transformer will plushify any file that has a ".plush" extension
func plushTransformer(ctx *plush.Context) genny.Transformer {
	t := genny.NewTransformer(".plush", func(f genny.File) (genny.File, error) {
		s, err := plush.RenderR(f, ctx)
		if err != nil {
			return f, errors.Wrap(err, f.Name())
		}
		return genny.NewFileS(f.Name(), s), nil
	})
	t.StripExt = true
	return t
}
