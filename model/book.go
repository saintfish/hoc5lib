package model

import (
	"errors"
	"fmt"
	"github.com/saintfish/orm.go"
	"log"
	"regexp"
	"time"
)

type Book struct {
	Id             int
	Barcode        string
	Title          string
	Authors        string
	PublishDate    *time.Time
	TypeCode       string
	Availability   bool
	AvailableAfter *time.Time
}

var bookSpec = orm.NewStructSpecBuilder(&Book{}).
	GenericOtherFields().
	SetPrimaryKey("Id").
	SetConstraints(map[string]string{"Barcode": "UNIQUE"}).
	Build()

func (*Book) TableSpec() orm.TableSpec {
	return bookSpec
}

var bookBarcodePattern = regexp.MustCompile(`^\d{13}$`)

func isBarcodeValid(barcode string) bool {
	return bookBarcodePattern.MatchString(barcode)
}

func IsBookValid(book *Book) error {
	if !isBarcodeValid(book.Barcode) {
		return errors.New("Invalid barcode.")
	}
	if len(book.Title) == 0 {
		return errors.New("Empty title.")
	}
	return nil
}

func GetBook(barcode string) (*Book, error) {
	if !isBarcodeValid(barcode) {
		return nil, errors.New("Invalid barcode.")
	}
	o := orm.New(db)
	b := Book{}
	err := o.Select().
		Where("Barcode = ?", barcode).
		Find(&b)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Book not found.")
	}
	return &b, nil
}

func AddBook(book *Book) error {
	err := IsBookValid(book)
	if err != nil {
		return err
	}
	book.Id = 0
	book.Availability = true
	book.AvailableAfter = nil
	b, _ := GetBook(book.Barcode)
	if b != nil {
		return errors.New("Barcode is used.")
	}
	o := orm.New(db)
	err = o.Insert(book)
	if err != nil {
		return errors.New("Error in adding book")
	}
	return nil
}

func UpdateBook(origValue, newValue *Book) error {
	if newValue.Barcode != origValue.Barcode {
		b, _ := GetBook(newValue.Barcode)
		if b != nil {
			return errors.New("Barcode is used.")
		}
	}
	origValue.Title = newValue.Title
	origValue.Authors = newValue.Authors
	origValue.PublishDate = newValue.PublishDate
	origValue.TypeCode = newValue.TypeCode
	err := IsBookValid(origValue)
	if err != nil {
		return err
	}
	o := orm.New(db)
	err = o.UpdateByPrimaryKey(origValue)
	if err != nil {
		return errors.New("Error in updating book")
	}
	return nil
}

func SearchBooksCount(query string) (int, error) {
	if len(query) <= 0 {
		return 0, errors.New("No query")
	}
	o := orm.New(db)
	prefixPattern := query + "%"
	fullPattern := "%" + query + "%"
	count, err := o.Select().
		Where("Barcode like ? OR Title like ? OR Authors like ? OR TypeCode like ?", prefixPattern, fullPattern, fullPattern, fullPattern).
		Count(&Book{})
	if err != nil {
		log.Println(err)
		return 0, errors.New("No result")
	}
	return count, nil
}

func SearchBooks(query string, limit, offset int) ([]Book, error) {
	if len(query) <= 0 {
		return nil, errors.New("No query")
	}
	o := orm.New(db)
	books := []Book{}
	prefixPattern := query + "%"
	fullPattern := "%" + query + "%"
	err := o.Select().
		Where("Barcode like ? OR Title like ? OR Authors like ? OR TypeCode like ?", prefixPattern, fullPattern, fullPattern, fullPattern).
		Limit(fmt.Sprintf("%d OFFSET %d", limit, offset)).
		Order("Barcode").
		FindAll(&books)
	if err != nil {
		log.Println(err)
		return nil, errors.New("No result")
	}
	return books, nil
}

func ListBooksCount() (int, error) {
	o := orm.New(db)
	count, err := o.Select().
		Count(&Book{})
	if err != nil {
		log.Println(err)
		return 0, errors.New("No result")
	}
	return count, nil
}

func ListBooks(limit, offset int) ([]Book, error) {
	o := orm.New(db)
	books := []Book{}
	err := o.Select().
		Limit(fmt.Sprintf("%d OFFSET %d", limit, offset)).
		Order("Barcode").
		FindAll(&books)
	if err != nil {
		log.Println(err)
		return nil, errors.New("No result")
	}
	return books, nil
}

func GetBookById(id int) (*Book, error) {
	book := &Book{Id: id}
	o := orm.New(db)
	err := o.FindByPrimaryKey(book)
	if err != nil {
		return nil, errors.New("Error in finding book")
	}
	return book, nil
}
