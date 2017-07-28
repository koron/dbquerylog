package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "vagrant:db1234@tcp(127.0.0.1:3306)/vagrant")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = test(db)
	if err != nil {
		log.Printf("WARN: %s", err)
	}
}

func test(db *sql.DB) error {
	rows, err := db.Query("show databases")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return err
		}
		fmt.Printf("table: %s\n", name)
	}
	return nil
}
