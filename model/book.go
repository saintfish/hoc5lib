package model

import (
	"errors"
	"fmt"
	"github.com/saintfish/orm.go"
	"log"
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
