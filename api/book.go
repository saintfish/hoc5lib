package api

import (
	"github.com/hoisie/web"
	"github.com/saintfish/hoc5lib/model"
	"github.com/saintfish/webutil"
	"strconv"
)

const (
	numBooksPerPage = 10
)

type bookSearchResult struct {
	NumResults, NumPage, Page int
	Books                     []model.Book
}

func BookSearch(ctx *web.Context) {
	query := ctx.Request.FormValue("q")
	page, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	numResults, err := model.SearchBooksCount(query)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	numPage := numResults / numBooksPerPage
	if numPage*numBooksPerPage < numResults {
		numPage++
	}
	if page < 1 {
		page = 1
	} else if page > numPage {
		page = numPage
	}
	books, err := model.SearchBooks(query, numBooksPerPage, (page-1)*numBooksPerPage)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, bookSearchResult{
		NumResults: numResults,
		NumPage:    numPage,
		Page:       page,
		Books:      books,
	})
}

func BookList(ctx *web.Context) {
	page, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	numResults, err := model.ListBooksCount()
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	numPage := numResults / numBooksPerPage
	if numPage*numBooksPerPage < numResults {
		numPage++
	}
	if page < 1 {
		page = 1
	} else if page > numPage {
		page = numPage
	}
	books, err := model.ListBooks(numBooksPerPage, (page-1)*numBooksPerPage)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, bookSearchResult{
		NumResults: numResults,
		NumPage:    numPage,
		Page:       page,
		Books:      books,
	})
	return
}
