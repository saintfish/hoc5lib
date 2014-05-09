package api

import (
	"github.com/hoisie/web"
	"github.com/saintfish/hoc5lib/model"
	"github.com/saintfish/webutil"
	"strconv"
	"time"
)

const (
	numBooksPerPage = 10
)

func BookGet(ctx *web.Context, barcode string) {
	book, err := model.GetBook(barcode)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, struct {
		Book *model.Book
	}{book})
}

func BookRange(ctx *web.Context, start, end string) {
	books, err := model.RangeBook(start, end)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, books)
}

func BookAdd(ctx *web.Context) {
	book := &model.Book{}
	err := webutil.ReadJson(ctx, book)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.IsBookValid(book)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.AddBook(book)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, book)
	return
}

func BookUpdate(ctx *web.Context, barcode string) {
	newValue := &model.Book{}
	err := webutil.ReadJson(ctx, newValue)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.IsBookValid(newValue)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	origValue, err := model.GetBook(barcode)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	err = model.UpdateBook(origValue, newValue)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, origValue)
	return

}

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

type overdueResultEntry struct {
	Borrower   *model.Borrower
	Book       *model.Book
	BorrowDate time.Time
	BorrowDays int
}
type overdueResult struct {
	Entry []overdueResultEntry
}

func BookOverdueList(ctx *web.Context) {
	records, err := model.ListOverdue()
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	result := []overdueResultEntry{}
	now := time.Now()
	for _, r := range records {
		book, err := model.GetBookById(r.BookId)
		if err != nil {
			webutil.Error(ctx, err)
			return
		}
		borrower, err := model.GetBorrowerById(r.BorrowerId)
		if err != nil {
			webutil.Error(ctx, err)
			return
		}
		borrowDays := int(now.Sub(r.BorrowDate) / (24 * time.Hour))
		entry := overdueResultEntry{
			Book:       book,
			Borrower:   borrower,
			BorrowDate: r.BorrowDate,
			BorrowDays: borrowDays,
		}
		result = append(result, entry)
	}
	webutil.Json(ctx, overdueResult{result})
	return
}
