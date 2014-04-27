package api

import (
	"github.com/hoisie/web"
	"github.com/saintfish/hoc5lib/model"
	"github.com/saintfish/webutil"
	"strconv"
)

const (
	numBorrowersPerPage = 5
)

func BorrowerGet(ctx *web.Context, phone string) {
	borrower, err := model.GetBorrower(phone)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, struct {
		Borrower *model.Borrower
	}{borrower})
}

func BorrowerAdd(ctx *web.Context) {
	borrower := &model.Borrower{}
	err := webutil.ReadJson(ctx, borrower)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.IsBorrowerValid(borrower)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.AddBorrower(borrower)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, borrower)
	return
}

func BorrowerUpdate(ctx *web.Context, phone string) {
	newValue := &model.Borrower{}
	err := webutil.ReadJson(ctx, newValue)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.IsBorrowerValid(newValue)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	origValue, err := model.GetBorrower(phone)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.UpdateBorrower(origValue, newValue)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, origValue)
	return
}

type borrowerSearchResult struct {
	NumPage, Page int
	Borrowers     []model.Borrower
}

func BorrowerSearch(ctx *web.Context) {
	query := ctx.Request.FormValue("q")
	page, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	numResults, err := model.SearchBorrowersCount(query)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	numPage := numResults / numBorrowersPerPage
	if numPage*numBorrowersPerPage < numResults {
		numPage++
	}
	if page < 1 {
		page = 1
	} else if page > numPage {
		page = numPage
	}
	borrowers, err := model.SearchBorrowers(query, numBorrowersPerPage, (page-1)*numBorrowersPerPage)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, borrowerSearchResult{
		NumPage:   numPage,
		Page:      page,
		Borrowers: borrowers,
	})
}

func BorrowerList(ctx *web.Context) {
	page, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	numResults, err := model.ListBorrowersCount()
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	numPage := numResults / numBorrowersPerPage
	if numPage*numBorrowersPerPage < numResults {
		numPage++
	}
	if page < 1 {
		page = 1
	} else if page > numPage {
		page = numPage
	}
	borrowers, err := model.ListBorrowers(numBorrowersPerPage, (page-1)*numBorrowersPerPage)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, borrowerSearchResult{
		NumPage:   numPage,
		Page:      page,
		Borrowers: borrowers,
	})
	return
}
