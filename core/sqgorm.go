package core

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Member Model
type Member struct {
	gorm.Model
	ID           uint `gorm:"primarykey"`
	Membername   string
	Password     string
	PasswordShow string
	Level        string
	Email        string
	Quota        int64
	Download     uint64
	Upload       uint64
	UseDays      uint
	ExpiryDate   string
}

// ORMOpen Sqlite for MemberModel
func ORMOpen(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Member{})
	return db
}

func (s *Sqlite) CreateMemberORM(id string, membername string, base64Pass string, originPass string) error {

	db := ORMOpen(s.Path)

	encryPass := sha256.Sum224([]byte(originPass))
	db.Create(&Member{Membername: membername, Password: fmt.Sprintf("%x", encryPass), PasswordShow: base64Pass})

	return nil
}

// UpdateMemberORM 更新Xray用户名和密码
func (s *Sqlite) UpdateMemberORM(id string, membername string, base64Pass string, originPass string) error {
	var member Member
	db := ORMOpen(s.Path)

	encryPass := sha256.Sum224([]byte(originPass))
	db.Where(&Member{Membername: membername}).First(&member)
	db.Model(&member).Updates(&Member{Password: fmt.Sprintf("%x", encryPass), PasswordShow: base64Pass})
	return nil

}

// DeleteMemberORM 删除用户
func (s *Sqlite) DeleteMemberORM(id string) error {
	var member Member
	db := ORMOpen(s.Path)
	fmt.Println("Deleteing record:")
	fmt.Println(id)
	idInt, _ := strconv.Atoi(id)
	db.Delete(&member, idInt)
	return nil
}

// ReadOneMemberORM 读取部分数据
func (s *Sqlite) QueryMemberORM(id string) *Member {
	var member Member
	db := ORMOpen(s.Path)
	idInt, _ := strconv.Atoi(id)
	db.Find(&member, idInt)
	return &member
}

// GetData 获取用户记录
func (s *Sqlite) GetDataORM(ids ...string) []*Member {
	var member []Member
	var memberList []*Member
	db := ORMOpen(s.Path)

	fmt.Println("Got records:")
	fmt.Println(len(ids))
	if len(ids) > 0 {
		fmt.Println("Find some records:")
		var idsInt []int
		for i, e := range ids {
			idInt, _ := strconv.Atoi(e)
			idsInt[i] = idInt
		}
		db.Find(&member, idsInt)
		fmt.Println(member)
		return member
	} else {
		fmt.Println("Find all records:")
		db.Where("id > ?", 0).Find(&member)
		for _, e := range member {
			memberList = append(memberList, &e)
		}
		fmt.Println(member)
		return memberList
	}
}
