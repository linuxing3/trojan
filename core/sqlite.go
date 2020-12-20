package core

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var defaultPath string = "./xray.db"

// Init init sqlite db
func Init(path string) (*DB, error) {
	if path == "" {
		path = defaultPath
	}
	os.Remove(path)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return db, err
}

// CreateTable create table in db
func CreateTable(dbName string, fields []string) bool {
	db, _ := init()
	sqlStmt := fmt.Sprintf("create table %s (id integer not null primary key, %s text);", dbName, strings.Join(fields, " text, "))
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return false
	}
	return true
}

func Insert(dbName string, fields []string, values []string) bool {

	db, _ := init()
	tx, err := db.Begin()

	if len(fields) != len(values) {
		return false
	}

	if err != nil {
		log.Fatal(err)
		return false
	}
	questions := []string{}
	for i := range values {
		questions[i] = "?"
	}

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", dbName, strings.Join(fields, ","), strings.Join(questions, ","))
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer stmt.Close()
	tx.Commit()
	return true
}

func InsertMany(dbName string, fields []string, values []string) {

	db, _ := init()
	tx, err := db.Begin()
	if len(fields) != len(values) {
		return false
	}

	if err != nil {
		log.Fatal(err)
		return false
	}
	questions := []string{}
	for i := range values {
		questions[i] = "?"
	}

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", dbName, strings.Join(fields, ","), strings.Join(questions, ","))
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(&values...)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
	return true
}

func QueryAll(dbName string, fields []string) {
	db, _ := init()
	rows, err := db.Query(fmt.Sprintf("select %s from %s", strings.Join(fields, ","), dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		// FIXME 如何扫描多个域
		err = rows.Scan(&id, &fields...)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}

func QueryOneByID(dbName string, id uint, fields []string) {
	db, _ := init()
	stmt, err = db.Prepare(fmt.Sprintf("select * from %s where id = ?", dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&fields...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}

func QueryManyWithFields(dbName string, fields []string) {
	db, _ := init()
	rows, err = db.Query("select %s from %s", strings.Join(fields, ","), dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		// FIXME 如何扫描多个域
		err = rows.Scan(&fields...)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
