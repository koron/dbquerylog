package main

// 1. Start MySQL with docker:
//
//		$ docker run --rm -e MYSQL_ROOT_PASSWORD=abcd1234 -p 3306:3306 --name test_mysql -d mysql:5.7.44
//
// 2. Run this program: `go run .` or `go run ./cmd/test-conn`
// 3. Stop MySQL: `docker stop test_mysql`

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	stInsert *sql.Stmt
	stSelect *sql.Stmt
	stUpdate *sql.Stmt
	stDelete *sql.Stmt
)

func main() {
	db, err := sql.Open("mysql", "root:abcd1234@tcp(127.0.0.1:3306)/mysql")
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
	if err := test4(db); err != nil {
		return err
	}
	if err := test5(db); err != nil {
		return err
	}
	// FIXME: implement other tests
	if err := test99(db); err != nil {
		return err
	}
	return nil
}

func test0(db *sql.DB) error {
	fmt.Println("test0: tests to list databases")
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
		fmt.Printf("    table:%s found\n", name)
	}
	return nil
}

func test1(db *sql.DB) error {
	fmt.Println("test1: create a \"users\" table use in tests")
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
	fmt.Println("test2: prepare statements")
	var err error
	stInsert, err = db.Prepare(
		`INSERT INTO users (name, password) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	stSelect, err = db.Prepare(`SELECT * FROM users WHERE name LIKE ?`)
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
	fmt.Println("test3: tests an error in the statements")
	var err error
	_, err = db.Prepare(
		`INSERT INTO users (name, password) VALUES (?, ?`)
	if err == nil {
		panic("prepare in test3 should be failed")
	}
	fmt.Printf("    IGNORED ERROR: %s\n", err)
	return nil
}

func test4(db *sql.DB) error {
	fmt.Println("test4: insert 6 records")
	insert := func(u, p string) error {
		r, err := stInsert.Exec(u, p)
		if err != nil {
			return err
		}
		id, _ := r.LastInsertId()
		if err != nil {
			return err
		}
		fmt.Printf("    inserted %q as %d\n", u, id)
		return nil
	}
	if err := insert("foo", "pass1234"); err != nil {
		return err
	}
	if err := insert("baz", "pass1234"); err != nil {
		return err
	}
	if err := insert("bar", "pass1234"); err != nil {
		return err
	}
	if err := insert("user001", "pass1234"); err != nil {
		return err
	}
	if err := insert("user002", "pass1234"); err != nil {
		return err
	}
	if err := insert("user003", "pass1234"); err != nil {
		return err
	}
	return nil
}

func test5(db *sql.DB) error {
	fmt.Println("test5: query with wildcard \"user%\"")
	rows, err := stSelect.Query("user%")
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func test99(db *sql.DB) error {
	fmt.Println("test99: cleanup test resources")
	if stDelete != nil {
		stDelete.Close()
	}
	if stUpdate != nil {
		stUpdate.Close()
	}
	if stSelect != nil {
		stSelect.Close()
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
