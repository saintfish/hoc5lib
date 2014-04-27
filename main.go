package main

import (
	"github.com/hoisie/web"
	"github.com/saintfish/hoc5lib/api"
	"github.com/saintfish/webutil"
)

func main() {
	web.Get("/", webutil.HandleTemplate("index.html", nil))
	web.Get("/(partials/.*)", func(ctx *web.Context, path string) {
		webutil.ExecuteTemplateWithContext(ctx, path, nil)
	})
	web.Get("/api/book/search", api.BookSearch)
	web.Get("/api/book/list", api.BookList)
	web.Get("/api/borrower/search", api.BorrowerSearch)
	web.Get("/api/borrower/list", api.BorrowerList)
	web.Post("/api/book/borrow", api.BookBorrow)
	web.Get("/api/book/(\\d+)/borrower", api.BookGetBorrower)
	web.Post("/api/book/return", api.BookReturn)
	web.Get("/api/borrower/(\\d+)", api.BorrowerGet)
	web.Post("/api/borrower/(\\d+)/update", api.BorrowerUpdate)
	web.Post("/api/borrower", api.BorrowerAdd)

	web.Run("0.0.0.0:9000")
}
