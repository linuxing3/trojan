package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var defaultPath string = "./xray.db"

// Init init sqlite db
func Init(path string) (*sql.DB, error) {
	log.Println("Creating sqlite-database.db...")
	if path == "" {
		path = defaultPath
	}
	os.Remove(path)
	file, err := os.Create(path) // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return db, err
}

// CreateDefaultTable create table in db
func CreateDefaultTable(dbName sql.NullString) bool {
	db, _ := Init("")
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
			id INT NOT NULL PRIMARY KEY AUTO INCREMENT,
			username VARCHAR(64) NOT NULL,
			password CHAR(56) NOT NULL,
			passwordShow VARCHAR(255) NOT NULL,
			email CHAR(56) NOT NULL DEFAULT 'love@example',
			level CHAR(56) NOT NULL DEFAULT 0,
			quota BIGINT NOT NULL DEFAULT 0,
			download BIGINT UNSIGNED NOT NULL DEFAULT 0,
			upload BIGINT UNSIGNED NOT NULL DEFAULT 0,
			useDays INT(10) DEFAULT 0,
			expiryDate CHAR(10) DEFAULT '',
			PRIMARY KEY (id),
			INDEX (password)
	);
			`); err != nil {
		fmt.Println(err)
	}
	return true
}

// CreateTable create table in db
func CreateTable(dbName string, fields []string) bool {
	db, _ := Init("")
	var defaultFields = []string{"username", "password", "passwordShow", "email", "level", "quota", "download", "upload", "useDays", "expiryDate"}
	if len(fields) == 0 {
		fields = defaultFields
	}
	sqlStmt := fmt.Sprintf("create table %s (id integer not null primary key, %s text);", dbName, strings.Join(fields, " text, "))
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return false
	}
	return true
}

// Insert 使用字段名和数据插入
func Insert(dbName string, fields []string, values []interface{}) sql.Result {

	db, _ := Init("")
	tx, err := db.Begin()

	if len(fields) != len(values) {
	}

	if err != nil {
		log.Fatal(err)
	}
	questions := []string{}
	for i := range values {
		questions[i] = "?"
	}

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", dbName, strings.Join(fields, ","), strings.Join(questions, ","))
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec(values...)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	tx.Commit()
	return result
}

// InsertMany 批量插入数据
func InsertMany(dbName string, fields []string, values []interface{}) {

	db, _ := Init("")
	tx, err := db.Begin()
	if len(fields) != len(values) {
		return
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	questions := []string{}
	for i := range fields {
		questions[i] = "?"
	}
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", dbName, strings.Join(fields, ","), strings.Join(questions, ","))
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(values...)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
	return
}

// QueryAll 使用字段名查询
func QueryAll(dbName string, fields []string, values []interface{}) {
	db, _ := Init("")
	rows, err := db.Query(fmt.Sprintf("select %s from %s", strings.Join(fields, ","), dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		// FIXME 如何扫描多个域
		err = rows.Scan(values...)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fields[0])
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}

// QueryOneByID 使用id查询
func QueryOneByID(dbName string, id uint, fields []string, values []interface{}) {
	db, _ := Init("")
	stmt, err := db.Prepare(fmt.Sprintf("select * from %s where id = ?", dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(values...)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// QueryManyWithFields 查询多个字段
func QueryManyWithFields(dbName string, fields []string, values []interface{}) {
	db, _ := Init("")
	rows, err := db.Query("select %s from %s", strings.Join(fields, ","), dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		// FIXME 如何扫描多个域
		err = rows.Scan(values)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fields[0])
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
