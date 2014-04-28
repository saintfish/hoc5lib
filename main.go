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
	web.Get("/api/book/(\\d+)", api.BookGet)
	web.Post("/api/book/(\\d+)", api.BookUpdate)
	web.Get("/api/book/(\\d+)/borrower", api.BookGetBorrower)
	web.Get("/api/book/list", api.BookList)
	web.Get("/api/book/search", api.BookSearch)
	web.Get("/api/borrower/(\\d+)", api.BorrowerGet)
	web.Get("/api/borrower/list", api.BorrowerList)
	web.Get("/api/borrower/search", api.BorrowerSearch)
	web.Post("/api/book/borrow", api.BookBorrow)
	web.Post("/api/book/return", api.BookReturn)
	web.Post("/api/borrower", api.BorrowerAdd)
	web.Post("/api/borrower/(\\d+)", api.BorrowerUpdate)

	web.Run("0.0.0.0:9000")
}
