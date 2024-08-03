package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/cetec")
	if err != nil {
		fmt.Println("no connection: %v", err)
	}
	if error := db.Ping(); err != nil {
		fmt.Println("no connection %v", error)
	}
}

func getDB() *sql.DB {
	return db
}
