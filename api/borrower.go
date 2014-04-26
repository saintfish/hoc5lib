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
