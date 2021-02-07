package core

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"trojan/util"

	_ "github.com/mattn/go-sqlite3"
)

// Sqlite 结构体
type Sqlite struct {
	Enabled  bool   `json:"enabled"`
	Path     string `json:"path"`
	Password string `json:"password"`
	Table    string `json:"table"`
}

var defaultPath string = "./xray.db"

// GetDB 获取sqlite数据库连接
func (sqlite *Sqlite) GetDB() *sql.DB {
	// 屏蔽sqlite驱动包的日志输出
	log.Println("Creating xray.db...")
	if _, err := os.Lstat(sqlite.Path); err != nil {
		// os.Remove(sqlite.Path)
		file, err := os.Create(sqlite.Path) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
	}
	db, err := sql.Open("sqlite3", sqlite.Path)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// CreateDefaultTable create table in db
func (sqlite *Sqlite) CreateDefaultTable() bool {
	db := sqlite.GetDB()
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS members (
			id CHAR(56) PRIMARY KEY NOT NULL,
			membername CHAR(64) NOT NULL,
			password CHAR(56) NOT NULL,
			passwordShow CHAR(255) NOT NULL,
			email CHAR(56),
			level CHAR(56),
			quota REAL,
			download REAL,
			upload REAL,
			useDays INT(10),
			expiryDate CHAR(10)
	);
			`); err != nil {
		fmt.Println(err)
	}
	return true
}

// 查询全部的用户列表，用于GetData，可以同时用于mysql和sqlite
func queryMemberList(db *sql.DB, sql string) ([]*Member, error) {
	var (
		id         uint
		membername   string
		originPass string
		level      string
		email      string
		passShow   string
		download   uint64
		upload     uint64
		quota      int64
		useDays    uint
		expiryDate string
	)
	var memberList []*Member
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &membername, &originPass, &level, &email, &passShow, &quota, &download, &upload, &useDays, &expiryDate); err != nil {
			return nil, err
		}
		fmt.Printf("用户名:" + membername)
		memberList = append(memberList, &Member{
			ID:         id,
			Membername:   membername,
			Password:   passShow,
			Level:      level,
			Email:      email,
			Download:   download,
			Upload:     upload,
			Quota:      quota,
			UseDays:    useDays,
			ExpiryDate: expiryDate,
		})
	}
	return memberList, nil
}
// 查询用户，用于GetData
func queryMember(db *sql.DB, sql string) (*Member, error) {
	var (
		id         uint
		membername   string
		originPass string
		passShow   string
		level      string
		email      string
		download   uint64
		upload     uint64
		quota      int64
		useDays    uint
		expiryDate string
	)
	row := db.QueryRow(sql)
	if err := row.Scan(&id, &membername, &originPass, &level, &email, &passShow, &quota, &download, &upload, &useDays, &expiryDate); err != nil {
		return nil, err
	}
	return &Member{ID: id, Membername: membername, Password: originPass, Download: download, Upload: upload, Quota: quota, UseDays: useDays, ExpiryDate: expiryDate}, nil
}

// CreateTable create table in db with fields array
func (sqlite *Sqlite) CreateTable(dbName string, fields []string) bool {
	db := sqlite.GetDB()
	var defaultFields = []string{"membername", "password", "passwordShow", "email", "level", "quota", "download", "upload", "useDays", "expiryDate"}
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

// CreateMember 创建Xray用户
func (sqlite *Sqlite) CreateMember(id string, membername string, base64Pass string, originPass string) error {

	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	encryPass := sha256.Sum224([]byte(originPass))
	if _, err := db.Exec(fmt.Sprintf("INSERT INTO members(id, membername, password, passwordShow, quota) VALUES ('%s', '%s', '%x', '%s', -1);", id, membername, encryPass, base64Pass)); err != nil {
		fmt.Println(err)
		return err
	}
	// TODO
	// return sqlite.CreateMemberORM(id, membername, base64Pass, originPass)
	return nil
}

// UpdateMember 更新Xray用户名和密码
func (sqlite *Sqlite) UpdateMember(id string, membername string, base64Pass string, originPass string) error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	encryPass := sha256.Sum224([]byte(originPass))
	if _, err := db.Exec(fmt.Sprintf("UPDATE members SET membername='%s', password='%x', passwordShow='%s' WHERE id='%s';", membername, encryPass, base64Pass, id)); err != nil {
		fmt.Println(err)
		return err
	}
	// TODO
	// return sqlite.UpdateMemberORM(id, membername, base64Pass, originPass)
	return nil
}

// DeleteMember 删除用户
func (sqlite *Sqlite) DeleteMember(id string) error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	if memberList, err := sqlite.GetData(id); err != nil {
		return err
	} else if memberList != nil && len(memberList) == 0 {
		return fmt.Errorf("不存在id为%s的用户", id)
	}
	if _, err := db.Exec(fmt.Sprintf("DELETE FROM members WHERE id='%s';", id)); err != nil {
		fmt.Println(err)
		return err
	}
	// TODO
	// return sqlite.DeleteMemberORM(id)
	return nil
}

// MonthlyResetData 设置了过期时间的用户，每月定时清空使用流量
func (sqlite *Sqlite) MonthlyResetData() error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	memberList, err := queryMemberList(db, "SELECT * FROM members WHERE useDays != 0 AND quota != 0")
	if err != nil {
		return err
	}
	for _, member := range memberList {
		if _, err := db.Exec(fmt.Sprintf("UPDATE members SET download=0, upload=0 WHERE id='%s';", member.ID)); err != nil {
			return err
		}
	}
	return nil
}

// DailyCheckExpire 检查是否有过期，过期了设置流量上限为0
func (sqlite *Sqlite) DailyCheckExpire() (bool, error) {
	needRestart := false
	now := time.Now()
	utc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return false, err
	}
	addDay, _ := time.ParseDuration("-24h")
	todayDay := now.Add(addDay).In(utc).Format("2006-01-02")
	db := sqlite.GetDB()
	if db == nil {
		return false, errors.New("can't connect sqlite")
	}
	defer db.Close()
	memberList, err := queryMemberList(db, "SELECT * FROM users WHERE useDays != 0 AND quota != 0")
	if err != nil {
		return false, err
	}
	for _, member := range memberList {
		if member.ExpiryDate == todayDay {
			if _, err := db.Exec(fmt.Sprintf("UPDATE members SET quota=0 WHERE id='%s';", member.ID)); err != nil {
				return false, err
			}
			if !needRestart {
				needRestart = true
			}
		}
	}
	return needRestart, nil
}

// CancelExpire 取消过期时间
func (sqlite *Sqlite) CancelExpire(id string) error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE members SET useDays=0, expiryDate='' WHERE id='%s';", id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SetExpire 设置过期时间
func (sqlite *Sqlite) SetExpire(id string, useDays uint) error {
	now := time.Now()
	utc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return err
	}
	addDay, _ := time.ParseDuration(strconv.Itoa(int(24*useDays)) + "h")
	expiryDate := now.Add(addDay).In(utc).Format("2006-01-02")

	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE members SET useDays=%d, expiryDate='%s' WHERE id='%s';", useDays, expiryDate, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SetQuota 限制流量
func (sqlite *Sqlite) SetQuota(id string, quota int) error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE members SET quota=%d WHERE id='%s';", quota, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// UpgradeDB 升级数据库表结构以及迁移数据
func (sqlite *Sqlite) UpgradeDB() error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	var field string
	error := db.QueryRow("SHOW COLUMNS FROM members LIKE 'passwordShow';").Scan(&field)
	if error == sql.ErrNoRows {
		fmt.Println(util.Yellow("正在进行数据库升级, 请稍等.."))
		if _, err := db.Exec("ALTER TABLE members ADD COLUMN passwordShow VARCHAR(255) NOT NULL AFTER password;"); err != nil {
			fmt.Println(err)
			return err
		}
		memberList, err := sqlite.GetData()
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, member := range memberList {
			pass, _ := GetValue(fmt.Sprintf("%s_pass", member.Membername))
			if pass != "" {
				base64Pass := base64.StdEncoding.EncodeToString([]byte(pass))
				if _, err := db.Exec(fmt.Sprintf("UPDATE members SET passwordShow='%s' WHERE id='%s';", base64Pass, member.ID)); err != nil {
					fmt.Println(err)
					return err
				}
				DelValue(fmt.Sprintf("%s_pass", member.Membername))
			}
		}
	}
	error = db.QueryRow("SHOW COLUMNS FROM members LIKE 'useDays';").Scan(&field)
	if error == sql.ErrNoRows {
		fmt.Println(util.Yellow("正在进行数据库升级, 请稍等.."))
		if _, err := db.Exec(`
ALTER TABLE members
ADD COLUMN useDays int(10) DEFAULT 0,
ADD COLUMN expiryDate char(10) DEFAULT '';
`); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

// CleanData 清空流量统计
func (sqlite *Sqlite) CleanData(id string) error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE members SET download=0, upload=0 WHERE id='%s';", id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// CleanDataByName 清空指定用户名流量统计数据
func (sqlite *Sqlite) CleanDataByName(membernames []string) error {
	db := sqlite.GetDB()
	if db == nil {
		return errors.New("can't connect sqlite")
	}
	defer db.Close()
	runSql := "UPDATE members SET download=0, upload=0 WHERE membername in ("
	for i, name := range membernames {
		runSql = runSql + "'" + name + "'"
		if i == len(membernames)-1 {
			runSql = runSql + ")"
		} else {
			runSql = runSql + ","
		}
	}
	if _, err := db.Exec(runSql); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetMemberByName 通过用户名来获取用户
func (sqlite *Sqlite) GetMemberByName(name string) *Member {
	db := sqlite.GetDB()
	if db == nil {
		return nil
	}
	defer db.Close()
	member, err := queryMember(db, fmt.Sprintf("SELECT * FROM members WHERE membername='%s'", name))
	if err != nil {
		return nil
	}
	return member
}

// GetMemberByPass 通过密码来获取用户
func (sqlite *Sqlite) GetMemberByPass(pass string) *Member {
	db := sqlite.GetDB()
	if db == nil {
		return nil
	}
	defer db.Close()
	member, err := queryMember(db, fmt.Sprintf("SELECT * FROM members WHERE passwordShow='%s'", pass))
	if err != nil {
		return nil
	}
	return member
}

// PageList 通过分页获取用户记录
func (sqlite *Sqlite) PageList(curPage int, pageSize int) (*PageQuery, error) {
	var (
		total int
	)

	db := sqlite.GetDB()
	if db == nil {
		return nil, errors.New("连接sqlite失败")
	}
	defer db.Close()
	offset := (curPage - 1) * pageSize
	querySQL := fmt.Sprintf("SELECT * FROM members LIMIT %d, %d", offset, pageSize)
	memberList, err := queryMemberList(db, querySQL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db.QueryRow("SELECT COUNT(id) FROM members").Scan(&total)
	return &PageQuery{
		CurPage:  curPage,
		PageSize: pageSize,
		Total:    total,
		DataList: memberList,
		PageNum:  (total + pageSize - 1) / pageSize,
	}, nil
}


// GetData 获取用户记录
func (sqlite *Sqlite) GetData(ids ...string) ([]*Member, error) {
	querySQL := "SELECT * FROM members"
	db := sqlite.GetDB()
	if db == nil {
		return nil, errors.New("连接sqlite失败")
	}
	defer db.Close()
	if len(ids) > 0 {
		querySQL = querySQL + " WHERE id in ('" + strings.Join(ids, "','") + "')"
	}
	fmt.Printf("[querySQL]: Get Data")
	fmt.Printf(querySQL)
	memberList, err := queryMemberList(db, querySQL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return memberList, nil
}

// =================================================
// 以下是通用的基本方法
// =================================================

// Insert 使用字段名和数据插入
func (sqlite *Sqlite) Insert(dbName string, fields []string, values []interface{}) sql.Result {

	db := sqlite.GetDB()
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
func (sqlite *Sqlite) InsertMany(dbName string, fields []string, values []interface{}) {

	db := sqlite.GetDB()
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
func (sqlite *Sqlite) QueryAll(dbName string, fields []string, values []interface{}) {
	db := sqlite.GetDB()
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
func (sqlite *Sqlite) QueryOneByID(dbName string, id uint, fields []string, values []interface{}) {
	db := sqlite.GetDB()
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
func (sqlite *Sqlite) QueryManyWithFields(dbName string, fields []string, values []interface{}) {
	db := sqlite.GetDB()
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
