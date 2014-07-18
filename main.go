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

func RequireAuth(d *webutil.Digest, f func(*web.Context)) func(*web.Context) {
	return func(ctx *web.Context) {
		DisableCache(ctx)
		if !webutil.HandleAuth(d, ctx) {
			return
		}
		f(ctx)
	}
}

func RequireAuth1(d *webutil.Digest, f func(*web.Context, string)) func(*web.Context, string) {
	return func(ctx *web.Context, arg string) {
		DisableCache(ctx)
		if !webutil.HandleAuth(d, ctx) {
			return
		}
		f(ctx, arg)
	}
}

func RequireAuth2(d *webutil.Digest, f func(*web.Context, string, string)) func(*web.Context, string, string) {
	return func(ctx *web.Context, arg1, arg2 string) {
		DisableCache(ctx)
		if !webutil.HandleAuth(d, ctx) {
			return
		}
		f(ctx, arg1, arg2)
	}
}

func main() {
	libAuth, err := webutil.NewDigest("hoclib", "htdigest", 10*time.Minute)
	if err != nil {
		panic(err)
	}

	web.Get("/logout", func(ctx *web.Context) {
		DisableCache(ctx)
		webutil.Logout(libAuth, ctx)
		webutil.ExecuteTemplateWithContext(ctx, "logout.html", nil)
	})
	web.Get("/", RequireAuth(libAuth, webutil.HandleTemplate("index.html", nil)))
	web.Get("/(partials/.*)", RequireAuth1(libAuth, func(ctx *web.Context, path string) {
		webutil.ExecuteTemplateWithContext(ctx, path, nil)
	}))
	web.Get("/api/book/(\\d+)", RequireAuth1(libAuth, api.BookGet))
	web.Get("/api/book/(\\d+)/borrower", RequireAuth1(libAuth, api.BookGetBorrower))
	web.Get("/api/book/list", RequireAuth(libAuth, api.BookList))
	web.Get("/api/book/search", RequireAuth(libAuth, api.BookSearch))
	web.Get("/api/borrower/(\\d+)", RequireAuth1(libAuth, api.BorrowerGet))
	web.Get("/api/borrower/list", RequireAuth(libAuth, api.BorrowerList))
	web.Get("/api/borrower/search", RequireAuth(libAuth, api.BorrowerSearch))
	web.Post("/api/book", RequireAuth(libAuth, api.BookAdd))
	web.Post("/api/book/(\\d+)", RequireAuth1(libAuth, api.BookUpdate))
	web.Post("/api/book/borrow", RequireAuth(libAuth, api.BookBorrow))
	web.Post("/api/book/return", RequireAuth(libAuth, api.BookReturn))
	web.Post("/api/borrower", RequireAuth(libAuth, api.BorrowerAdd))
	web.Post("/api/borrower/(\\d+)", RequireAuth1(libAuth, api.BorrowerUpdate))
	web.Get("/api/book/overdue", RequireAuth(libAuth, api.BookOverdueList))
	web.Get("/api/book/(\\d+)/(\\d+)", RequireAuth2(libAuth, api.BookRange))
	web.Get("/barcode/(\\d+)", api.HandleEan13)
	web.Get("/barcode/book/(\\d+)/(\\d+)", RequireAuth2(libAuth, api.HandleBookBarcodeRange))
	web.Post("/barcode/book/print.pdf", RequireAuth(libAuth, api.HandleBookBarcodePDF))

	web.Run("127.0.0.1:9000")
}
