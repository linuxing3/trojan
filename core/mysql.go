package core

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"trojan/util"

	mysqlDriver "github.com/go-sql-driver/mysql"

	"strconv"
	"strings"

	// mysql sql驱动
	_ "github.com/go-sql-driver/mysql"
)

// Mysql 结构体
type Mysql struct {
	Enabled    bool   `json:"enabled"`
	ServerAddr string `json:"server_addr"`
	ServerPort int    `json:"server_port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Cafile     string `json:"cafile"`
}

// User 用户表记录结构体
type User struct {
	ID         string
	Username   string
	Password   string
	Level      string
	Email      string
	Quota      int64
	Download   uint64
	Upload     uint64
	UseDays    uint
	ExpiryDate string
}

// PageQuery 分页查询的结构体
type PageQuery struct {
	PageNum  int
	CurPage  int
	Total    int
	PageSize int
	DataList []*User
}

// GetDB 获取mysql数据库连接
func (mysql *Mysql) GetDB() *sql.DB {
	// 屏蔽mysql驱动包的日志输出
	mysqlDriver.SetLogger(log.New(ioutil.Discard, "", 0))
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysql.Username, mysql.Password, mysql.ServerAddr, mysql.ServerPort, mysql.Database)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

// CreateTable 不存在Xray user表则自动创建
func (mysql *Mysql) CreateTable() {
	db := mysql.GetDB()
	defer db.Close()
	if _, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id CHAR(56) NOT NULL,
    username VARCHAR(64) NOT NULL,
    password CHAR(56) NOT NULL,
    passwordShow VARCHAR(255) NOT NULL,
    email CHAR(56) NOT NULL DEFAULT 'love@example',
    level CHAR(56) NOT NULL DEFAULT 0,
    quota BIGINT NOT NULL DEFAULT 0,
    download BIGINT UNSIGNED NOT NULL DEFAULT 0,
    upload BIGINT UNSIGNED NOT NULL DEFAULT 0,
    useDays int(10) DEFAULT 0,
    expiryDate char(10) DEFAULT '',
    PRIMARY KEY (id),
    INDEX (password)
);
    `); err != nil {
		fmt.Println(err)
	}
}

// 数据库查询全部的用户列表，用于GetData
func queryUserList(db *sql.DB, sql string) ([]*User, error) {
	var (
		id         string
		username   string
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
	var userList []*User
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &username, &originPass, &level, &email, &passShow, &quota, &download, &upload, &useDays, &expiryDate); err != nil {
			return nil, err
		}
		userList = append(userList, &User{
			ID:         id,
			Username:   username,
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
	return userList, nil
}

func queryUser(db *sql.DB, sql string) (*User, error) {
	var (
		id         string
		username   string
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
	row := db.QueryRow(sql)
	if err := row.Scan(&id, &username, &originPass, &level, &email, &passShow, &quota, &download, &upload, &useDays, &expiryDate); err != nil {
		return nil, err
	}
	return &User{ID: id, Username: username, Password: originPass, Download: download, Upload: upload, Quota: quota, UseDays: useDays, ExpiryDate: expiryDate}, nil
}

// CreateUser 创建Xray用户
func (mysql *Mysql) CreateUser(id string, username string, base64Pass string, originPass string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	encryPass := sha256.Sum224([]byte(originPass))
	if _, err := db.Exec(fmt.Sprintf("INSERT INTO users(id, username, password, passwordShow, quota) VALUES ('%s', '%s', '%x', '%s', -1);", id, username, encryPass, base64Pass)); err != nil {
		fmt.Println(err)
		return err
	}
	// if ok write to configuration file
	if success := WriteInbloudClient([]string{id}, "create"); success == true {
		fmt.Println("成功在配置文件中假如客户端信息，请重启xray服务器")
	}
	return nil
}

// UpdateUser 更新Xray用户名和密码
func (mysql *Mysql) UpdateUser(id string, username string, base64Pass string, originPass string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	encryPass := sha256.Sum224([]byte(originPass))
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET username='%s', password='%x', passwordShow='%s' WHERE id=%d;", username, encryPass, base64Pass, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// DeleteUser 删除用户
func (mysql *Mysql) DeleteUser(id string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	if userList, err := mysql.GetData(id); err != nil {
		return err
	} else if userList != nil && len(userList) == 0 {
		return fmt.Errorf("不存在id为%d的用户", id)
	}
	if _, err := db.Exec(fmt.Sprintf("DELETE FROM users WHERE id=%d;", id)); err != nil {
		fmt.Println(err)
		return err
	}
	// if ok write to configuration file
	if success := WriteInbloudClient([]string{id}, "delete"); success == true {
		fmt.Println("成功删除配置文件中客户端信息，请重启xray服务器")
	}
	return nil
}

// MonthlyResetData 设置了过期时间的用户，每月定时清空使用流量
func (mysql *Mysql) MonthlyResetData() error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	userList, err := queryUserList(db, "SELECT * FROM users WHERE useDays != 0 AND quota != 0")
	if err != nil {
		return err
	}
	for _, user := range userList {
		if _, err := db.Exec(fmt.Sprintf("UPDATE users SET download=0, upload=0 WHERE id=%d;", user.ID)); err != nil {
			return err
		}
	}
	return nil
}

// DailyCheckExpire 检查是否有过期，过期了设置流量上限为0
func (mysql *Mysql) DailyCheckExpire() (bool, error) {
	needRestart := false
	now := time.Now()
	utc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return false, err
	}
	addDay, _ := time.ParseDuration("-24h")
	todayDay := now.Add(addDay).In(utc).Format("2006-01-02")
	db := mysql.GetDB()
	if db == nil {
		return false, errors.New("can't connect mysql")
	}
	defer db.Close()
	userList, err := queryUserList(db, "SELECT * FROM users WHERE useDays != 0 AND quota != 0")
	if err != nil {
		return false, err
	}
	for _, user := range userList {
		if user.ExpiryDate == todayDay {
			if _, err := db.Exec(fmt.Sprintf("UPDATE users SET quota=0 WHERE id=%d;", user.ID)); err != nil {
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
func (mysql *Mysql) CancelExpire(id string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET useDays=0, expiryDate='' WHERE id=%d;", id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SetExpire 设置过期时间
func (mysql *Mysql) SetExpire(id string, useDays uint) error {
	now := time.Now()
	utc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return err
	}
	addDay, _ := time.ParseDuration(strconv.Itoa(int(24*useDays)) + "h")
	expiryDate := now.Add(addDay).In(utc).Format("2006-01-02")

	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET useDays=%d, expiryDate='%s' WHERE id=%d;", useDays, expiryDate, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SetQuota 限制流量
func (mysql *Mysql) SetQuota(id string, quota int) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET quota=%d WHERE id=%d;", quota, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// UpgradeDB 升级数据库表结构以及迁移数据
func (mysql *Mysql) UpgradeDB() error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	var field string
	error := db.QueryRow("SHOW COLUMNS FROM users LIKE 'passwordShow';").Scan(&field)
	if error == sql.ErrNoRows {
		fmt.Println(util.Yellow("正在进行数据库升级, 请稍等.."))
		if _, err := db.Exec("ALTER TABLE users ADD COLUMN passwordShow VARCHAR(255) NOT NULL AFTER password;"); err != nil {
			fmt.Println(err)
			return err
		}
		userList, err := mysql.GetData()
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, user := range userList {
			pass, _ := GetValue(fmt.Sprintf("%s_pass", user.Username))
			if pass != "" {
				base64Pass := base64.StdEncoding.EncodeToString([]byte(pass))
				if _, err := db.Exec(fmt.Sprintf("UPDATE users SET passwordShow='%s' WHERE id=%d;", base64Pass, user.ID)); err != nil {
					fmt.Println(err)
					return err
				}
				DelValue(fmt.Sprintf("%s_pass", user.Username))
			}
		}
	}
	error = db.QueryRow("SHOW COLUMNS FROM users LIKE 'useDays';").Scan(&field)
	if error == sql.ErrNoRows {
		fmt.Println(util.Yellow("正在进行数据库升级, 请稍等.."))
		if _, err := db.Exec(`
ALTER TABLE users
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
func (mysql *Mysql) CleanData(id string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET download=0, upload=0 WHERE id=%d;", id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// CleanDataByName 清空指定用户名流量统计数据
func (mysql *Mysql) CleanDataByName(usernames []string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("can't connect mysql")
	}
	defer db.Close()
	runSql := "UPDATE users SET download=0, upload=0 WHERE username in ("
	for i, name := range usernames {
		runSql = runSql + "'" + name + "'"
		if i == len(usernames)-1 {
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

// GetUserByName 通过用户名来获取用户
func (mysql *Mysql) GetUserByName(name string) *User {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	defer db.Close()
	user, err := queryUser(db, fmt.Sprintf("SELECT * FROM users WHERE username='%s'", name))
	if err != nil {
		return nil
	}
	return user
}

// GetUserByPass 通过密码来获取用户
func (mysql *Mysql) GetUserByPass(pass string) *User {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	defer db.Close()
	user, err := queryUser(db, fmt.Sprintf("SELECT * FROM users WHERE passwordShow='%s'", pass))
	if err != nil {
		return nil
	}
	return user
}

// PageList 通过分页获取用户记录
func (mysql *Mysql) PageList(curPage int, pageSize int) (*PageQuery, error) {
	var (
		total int
	)

	db := mysql.GetDB()
	if db == nil {
		return nil, errors.New("连接mysql失败")
	}
	defer db.Close()
	offset := (curPage - 1) * pageSize
	querySQL := fmt.Sprintf("SELECT * FROM users LIMIT %d, %d", offset, pageSize)
	userList, err := queryUserList(db, querySQL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db.QueryRow("SELECT COUNT(id) FROM users").Scan(&total)
	return &PageQuery{
		CurPage:  curPage,
		PageSize: pageSize,
		Total:    total,
		DataList: userList,
		PageNum:  (total + pageSize - 1) / pageSize,
	}, nil
}

// GetData 获取用户记录
func (mysql *Mysql) GetData(ids ...string) ([]*User, error) {
	querySQL := "SELECT * FROM users"
	db := mysql.GetDB()
	if db == nil {
		return nil, errors.New("连接mysql失败")
	}
	defer db.Close()
	if len(ids) > 0 {
		querySQL = querySQL + " WHERE id in (" + strings.Join(ids, ",") + ")"
	}
	userList, err := queryUserList(db, querySQL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return userList, nil
}
