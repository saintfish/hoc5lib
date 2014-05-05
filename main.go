package main

import (
	"github.com/hoisie/web"
	"github.com/saintfish/hoc5lib/api"
	"github.com/saintfish/webutil"
	"time"
)

func DisableCache(ctx *web.Context) {
	ctx.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header().Set("Pragma", "no-cache")
	ctx.Header().Set("Expires", "0")
}

func RequireAuth(f func(*web.Context)) func(*web.Context) {
	return func(ctx *web.Context) {
		DisableCache(ctx)
		if !webutil.HandleAuth(ctx) {
			return
		}
		f(ctx)
	}
}

func RequireAuth1(f func(*web.Context, string)) func(*web.Context, string) {
	return func(ctx *web.Context, arg string) {
		DisableCache(ctx)
		if !webutil.HandleAuth(ctx) {
			return
		}
		f(ctx, arg)
	}
}

func main() {
	webutil.InitHtdigest("hoclib", "/home/yusheng/htdigest", 1*time.Minute)

	web.Get("/logout", func(ctx *web.Context) {
		webutil.Logout(ctx)
		webutil.ExecuteTemplateWithContext(ctx, "logout.html", nil)
	})
	web.Get("/", RequireAuth(webutil.HandleTemplate("index.html", nil)))
	web.Get("/(partials/.*)", RequireAuth1(func(ctx *web.Context, path string) {
		webutil.ExecuteTemplateWithContext(ctx, path, nil)
	}))
	web.Get("/api/book/(\\d+)", RequireAuth1(api.BookGet))
	web.Get("/api/book/(\\d+)/borrower", RequireAuth1(api.BookGetBorrower))
	web.Get("/api/book/list", RequireAuth(api.BookList))
	web.Get("/api/book/search", RequireAuth(api.BookSearch))
	web.Get("/api/borrower/(\\d+)", RequireAuth1(api.BorrowerGet))
	web.Get("/api/borrower/list", RequireAuth(api.BorrowerList))
	web.Get("/api/borrower/search", RequireAuth(api.BorrowerSearch))
	web.Post("/api/book", RequireAuth(api.BookAdd))
	web.Post("/api/book/(\\d+)", RequireAuth1(api.BookUpdate))
	web.Post("/api/book/borrow", RequireAuth(api.BookBorrow))
	web.Post("/api/book/return", RequireAuth(api.BookReturn))
	web.Post("/api/borrower", RequireAuth(api.BorrowerAdd))
	web.Post("/api/borrower/(\\d+)", RequireAuth1(api.BorrowerUpdate))
	web.Get("/api/book/overdue", RequireAuth(api.BookOverdueList))

	web.Run("0.0.0.0:9000")
}
