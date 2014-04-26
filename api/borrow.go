package api

import (
	"github.com/hoisie/web"
	"github.com/saintfish/hoc5lib/model"
	"github.com/saintfish/webutil"
	"log"
)

func BookBorrow(ctx *web.Context) {
	barcode := ctx.Request.FormValue("barcode")
	phone := ctx.Request.FormValue("phone")

	book, err := model.GetBook(barcode)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	borrower, err := model.GetBorrower(phone)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	entry, err := model.BorrowBook(borrower, book)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	log.Printf("Borrow: book %s, borrower %s", book.Barcode, borrower.Phone)
	webutil.Json(ctx, entry)
}

func BookGetBorrower(ctx *web.Context, barcode string) {
	book, err := model.GetBook(barcode)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	borrower, err := model.GetBorrowerFromBook(book)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, struct {
		Book     *model.Book
		Borrower *model.Borrower
	}{book, borrower})
}

func BookReturn(ctx *web.Context) {
	barcode := ctx.Request.FormValue("barcode")

	book, err := model.GetBook(barcode)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	borrower, err := model.GetBorrowerFromBook(book)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	entry, err := model.ReturnBook(borrower, book)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, struct {
		Book     *model.Book
		Borrower *model.Borrower
		Entry    *model.BorrowEntry
	}{book, borrower, entry})
}
