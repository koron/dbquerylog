package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	stInsert *sql.Stmt
	stUpdate *sql.Stmt
	stDelete *sql.Stmt
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
	if err := test0(db); err != nil {
		return err
	}
	if err := test1(db); err != nil {
		return err
	}
	if err := test2(db); err != nil {
		return err
	}
	if err := test3(db); err != nil {
		return err
	}
	// TODO:
	if err := test99(db); err != nil {
		return err
	}
	return nil
}

func test0(db *sql.DB) error {
	rows, err := db.Query(`SHOW DATABASES`)
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

func test1(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) UNIQUE,
		password VARCHAR(255)
	)`)
	if err != nil {
		return err
	}
	return nil
}

func test2(db *sql.DB) error {
	var err error
	stInsert, err = db.Prepare(
		`INSERT INTO users (name, password) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	stUpdate, err = db.Prepare(
		`UPDATE users SET name = ?, password = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	stDelete, err = db.Prepare(`DELETE FROM users WHERE id = ?`)
	if err != nil {
		return err
	}
	return nil
}

func test3(db *sql.DB) error {
	var err error
	_, err = db.Prepare(
		`INSERT INTO users (name, password) VALUES (?, ?`)
	if err == nil {
		panic("prepare in test3 should be failed")
	}
	log.Printf("test3: %s", err)
	return nil
}

func test99(db *sql.DB) error {
	if stDelete != nil {
		stDelete.Close()
	}
	if stUpdate != nil {
		stUpdate.Close()
	}
	if stInsert != nil {
		stInsert.Close()
	}
	_, err := db.Exec(`DROP TABLE users`)
	if err != nil {
		return err
	}
	return nil
}
