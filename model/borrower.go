package model

import (
	"errors"
	"fmt"
	"github.com/saintfish/orm.go"
	"log"
	"regexp"
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

var borrowerPhonePattern = regexp.MustCompile(`^\d{10}$`)

func isPhoneValid(phone string) bool {
	return borrowerPhonePattern.MatchString(phone)
}

func IsBorrowerValid(borrower *Borrower) error {
	if !isPhoneValid(borrower.Phone) {
		return errors.New("Invalid phone number.")
	}
	if len(borrower.ChineseName) == 0 && len(borrower.EnglishName) == 0 {
		return errors.New("Both names are empty.")
	}
	return nil
}

func GetBorrower(phone string) (*Borrower, error) {
	if !isPhoneValid(phone) {
		return nil, errors.New("Invalid phone number.")
	}
	o := orm.New(db)
	b := Borrower{}
	err := o.Select().
		Where("Phone = ?", phone).
		Find(&b)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Borrower not found.")
	}
	return &b, nil
}

func UpdateBorrower(origValue, newValue *Borrower) error {
	if newValue.Phone != origValue.Phone {
		b, _ := GetBorrower(newValue.Phone)
		if b != nil {
			return errors.New("Barcode is used.")
		}
	}
	origValue.Phone = newValue.Phone
	origValue.EnglishName = newValue.EnglishName
	origValue.ChineseName = newValue.ChineseName
	origValue.MorePhones = newValue.MorePhones
	o := orm.New(db)
	err := o.UpdateByPrimaryKey(origValue)
	if err != nil {
		return errors.New("Error in updating borrower")
	}
	return nil
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
