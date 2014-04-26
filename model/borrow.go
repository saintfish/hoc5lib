package model

import (
	"errors"
	"github.com/saintfish/orm.go"
	"log"
	"time"
)

const MaxBooksOneCanBorrow = 3
const MaxDaysBorrowPeriod = 15

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
	now := time.Now()
	returnDate := now.AddDate(0, 0, MaxDaysBorrowPeriod)
	book.AvailableAfter = &returnDate
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

func GetBorrowerFromBook(book *Book) (*Borrower, error) {
	if book.Availability {
		return nil, errors.New("The book is already returned.")
	}
	o := orm.New(db)
	entry := BorrowEntry{}
	err := o.Select().
		Where("BookId = ? AND ReturnDate IS NULL", book.Id).
		Find(&entry)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error in querying borrow record.")
	}
	borrower := Borrower{Id: entry.BorrowerId}
	err = o.FindByPrimaryKey(&borrower)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Unable to find the borrower for the book.")
	}
	return &borrower, nil
}

func ReturnBook(borrower *Borrower, book *Book) (*BorrowEntry, error) {
	if book.Availability {
		return nil, errors.New("The book is already returned.")
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("Error in returning book.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	o := orm.New(tx)

	entry := BorrowEntry{}
	err = o.Select().
		Where("BorrowerId = ? AND BookId = ? AND ReturnDate IS NULL", borrower.Id, book.Id).
		Find(&entry)
	if err != nil {
		return nil, errors.New("Borrow record not found.")
	}

	now := time.Now()
	entry.ReturnDate = &now
	err = o.UpdateByPrimaryKey(&entry)
	if err != nil {
		return nil, errors.New("Error in returning book.")
	}
	book.Availability = true
	book.AvailableAfter = nil
	err = o.UpdateByPrimaryKey(book)
	if err != nil {
		return nil, errors.New("Error in returning book.")
	}
	borrower.NumBorrowed--
	err = o.UpdateByPrimaryKey(borrower)
	if err != nil {
		return nil, errors.New("Error in returning book.")
	}
	return &entry, nil
}
