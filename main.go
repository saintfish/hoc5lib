package main

import (
	"github.com/hoisie/web"
	"github.com/saintfish/webutil"
)

func main() {
	web.Get("/", webutil.HandleTemplate("index.html", nil))
	web.Get("/(partials/.*)", func(ctx *web.Context, path string) {
		webutil.ExecuteTemplateWithContext(ctx, path, nil)
	})
	web.Run("0.0.0.0:9000")
}
