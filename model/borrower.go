package model

import (
	"errors"
	"fmt"
	"github.com/saintfish/orm.go"
	"log"
)

type Borrower struct {
	Id          int
	Phone       string
	EnglishName string
	ChineseName string
	MorePhones  string
	NumBorrowed int
}

var borrowerSpec = orm.NewStructSpecBuilder(&Borrower{}).
	GenericOtherFields().
	SetPrimaryKey("Id").
	SetConstraints(map[string]string{"Phone": "UNIQUE"}).
	Build()

func (*Borrower) TableSpec() orm.TableSpec {
	return borrowerSpec
}

func SearchBorrowersCount(query string) (int, error) {
	if len(query) <= 0 {
		return 0, errors.New("No query")
	}
	o := orm.New(db)
	prefixPattern := query + "%"
	fullPattern := "%" + query + "%"
	count, err := o.Select().
		Where("Phone like ? OR EnglishName like ? OR ChineseName like ? OR MorePhones like ?", prefixPattern, fullPattern, fullPattern, prefixPattern).
		Count(&Borrower{})
	if err != nil {
		log.Println(err)
		return 0, errors.New("No result")
	}
	return count, nil
}

func SearchBorrowers(query string, limit, offset int) ([]Borrower, error) {
	if len(query) <= 0 {
		return nil, errors.New("No query")
	}
	o := orm.New(db)
	borrowers := []Borrower{}
	prefixPattern := query + "%"
	fullPattern := "%" + query + "%"
	err := o.Select().
		Where("Phone like ? OR EnglishName like ? OR ChineseName like ? OR MorePhones like ?", prefixPattern, fullPattern, fullPattern, prefixPattern).
		Limit(fmt.Sprintf("%d OFFSET %d", limit, offset)).
		Order("Phone").
		FindAll(&borrowers)
	if err != nil {
		log.Println(err)
		return nil, errors.New("No result")
	}
	return borrowers, nil
}

func ListBorrowersCount() (int, error) {
	o := orm.New(db)
	count, err := o.Select().
		Count(&Borrower{})
	if err != nil {
		log.Println(err)
		return 0, errors.New("No result")
	}
	return count, nil
}

func ListBorrowers(limit, offset int) ([]Borrower, error) {
	o := orm.New(db)
	borrowers := []Borrower{}
	err := o.Select().
		Limit(fmt.Sprintf("%d OFFSET %d", limit, offset)).
		Order("Phone").
		FindAll(&borrowers)
	if err != nil {
		log.Println(err)
		return nil, errors.New("No result")
	}
	return borrowers, nil
}
