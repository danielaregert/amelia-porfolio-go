package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/hook"

	_ "porfolio-amelia/migrations"

	"porfolio-amelia/internal/handlers"
)

//go:embed all:static
var staticFS embed.FS

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangGo,
		Automigrate:  true,
		Dir:          "migrations",
	})

	app.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
		Func: func(e *core.ServeEvent) error {
			handlers.RegisterRoutes(e.Router, app)

			if !e.Router.HasRoute(http.MethodGet, "/static/{path...}") {
				staticSub, err := fs.Sub(staticFS, "static")
				if err != nil {
					return err
				}
				e.Router.GET("/static/{path...}", apis.Static(staticSub, false))
			}
			return e.Next()
		},
		Priority: 100,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
