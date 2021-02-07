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

// CreateMemberORM 使用给定信息创建成员
func (s *Sqlite) CreateMemberORM(id string, membername string, base64Pass string, originPass string) error {

	db := ORMOpen(s.Path)

	encryPass := sha256.Sum224([]byte(originPass))
	if err := db.Create(&Member{Membername: membername, Password: fmt.Sprintf("%x", encryPass), PasswordShow: base64Pass}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateMemberORM 使用给定信息更新用户名和密码
func (s *Sqlite) UpdateMemberORM(id string, membername string, base64Pass string, originPass string) error {
	var member Member
	db := ORMOpen(s.Path)

	encryPass := sha256.Sum224([]byte(originPass))
	if err := db.Where(&Member{Membername: membername}).First(&member).Error; err != nil {
		return err
	}
	if err := db.Model(&member).Updates(&Member{Password: fmt.Sprintf("%x", encryPass), PasswordShow: base64Pass}).Error; err != nil {
		return err
	}
	return nil

}

// DeleteMemberORM 使用给定信息删除用户
func (s *Sqlite) DeleteMemberORM(id string) error {
	var member Member
	db := ORMOpen(s.Path)
	fmt.Println("Deleteing record:")
	fmt.Println(id)
	idInt, _ := strconv.Atoi(id)
	if err := db.Delete(&member, idInt).Error; err != nil {
		return err
	}
	return nil
}

// QueryMemberORM 用id查询数据
func (s *Sqlite) QueryMemberORM(id string) (*Member, error) {
	var member Member
	db := ORMOpen(s.Path)
	idInt, _ := strconv.Atoi(id)
	if err := db.Find(&member, idInt).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetDataORM 根据指定多个id获取用户记录
func (s *Sqlite) GetDataORM(ids ...string) ([]*Member, error) {
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
		if err := db.Find(&member, idsInt).Error; err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Find all records:")
		if err := db.Where("id > ?", 0).Find(&member).Error; err != nil {
			return nil, err
		}
	}
	// 更改为指针数组
	for _, e := range member {
		memberList = append(memberList, &e)
	}
	fmt.Println(member)
	return memberList, nil
}
