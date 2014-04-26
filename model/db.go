package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saintfish/orm.go"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	o := orm.New(db)
	o.CreateTable(&Book{}, true)
}
