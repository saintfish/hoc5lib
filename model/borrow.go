package model

import (
	"errors"
	"github.com/saintfish/orm.go"
	"log"
	"time"
)

const MaxBooksOneCanBorrow = 3

type BorrowEntry struct {
	Id         int
	BookId     int
	BorrowerId int
	BorrowDate time.Time
	ReturnDate *time.Time
}

var borrowSpec = orm.NewStructSpecBuilder(&BorrowEntry{}).
	GenericOtherFields().
	SetTable("Borrow").
	SetPrimaryKey("Id").
	Build()

func (*BorrowEntry) TableSpec() orm.TableSpec {
	return borrowSpec
}

func BorrowBook(borrower *Borrower, book *Book) (*BorrowEntry, error) {
	if !book.Availability {
		return nil, errors.New("The book is not available for borrow.")
	}
	if borrower.NumBorrowed >= MaxBooksOneCanBorrow {
		return nil, errors.New("Too many books already borrowed.")
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error in borrowing the book.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	o := orm.New(tx)

	entry := BorrowEntry{
		BookId:     book.Id,
		BorrowerId: borrower.Id,
		BorrowDate: time.Now(),
		ReturnDate: nil,
	}
	err = o.Insert(&entry)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error in borrowing the book.")
	}
	book.Availability = false
	err = o.UpdateByPrimaryKey(book)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error in borrowing the book.")
	}
	borrower.NumBorrowed++
	err = o.UpdateByPrimaryKey(borrower)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error in borrowing the book.")
	}
	return &entry, nil
}
