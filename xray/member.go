package xray

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"trojan/core"
	"trojan/util"

	"github.com/google/uuid"
)

// MemberMenu 用户管理菜单
func MemberMenu() {
	fmt.Println()
	menu := []string{"新增用户", "删除用户", "限制流量", "清空流量", "设置限期", "取消限期", "换数据库"}
	switch util.LoopInput("请选择: ", menu, false) {
	case 1:
		AddMember()
	case 2:
		DelMember()
	case 3:
		SetMemberQuota()
	case 4:
		CleanMemberData()
	case 5:
		SetupMemberExpire()
	case 6:
		CancelMemberExpire()
	case 7:
		UserMenu()
	}
}

// AddMember 直接后台添加成员
func AddMember() {
	// randomUser name and pass
	randomUser := util.RandString(4)
	randomPass := util.RandString(8)
	inputUser := util.Input(fmt.Sprintf("生成随机用户名: %s, 使用直接回车, 否则输入自定义用户名: ", randomUser), randomUser)
	if inputUser == "admin" {
		fmt.Println(util.Yellow("不能新建用户名为'admin'的用户!"))
		return
	}
	// 1. uuid，用于xray
	uuid := fmt.Sprintf("%s", uuid.New())
	fmt.Println(util.Yellow("[uuid]:" + uuid))

	// 2. 生成随机密码，通过密码获取用户，存在报错
	inputPass := util.Input(fmt.Sprintf("生成随机密码: %s, 使用直接回车, 否则输入自定义密码: ", randomPass), randomPass)
	base64Pass := base64.StdEncoding.EncodeToString([]byte(inputPass))

	// 创建Sqlite新用户
	// FIXED 这里的配置是用硬盘配置文件中读取的，所以记得先写入配置文件才能正常使用
	sqlite := core.GetSqlite()
	if err := sqlite.CreateMemberORM(uuid, inputUser, base64Pass, inputPass); err != nil {
		fmt.Println("新增Sqlite用户成功!")
		fmt.Println("")
	} else {
		fmt.Println(err)
	}
}

// DelMember 后台删除成员
func DelMember() {
	memberList := MemberList()
	fmt.Println("Record list:")
	fmt.Println(memberList)
	choice := util.LoopInput("请选择要删除的用户序号: ", memberList, true)
	if choice == -1 {
		return
	}
	sqlite := core.GetSqlite()
	if err := sqlite.DeleteMemberORM(fmt.Sprint(memberList[choice-1].ID)); err != nil {
		fmt.Println("删除Sqlite用户成功!")
		fmt.Println("")
	} else {
		fmt.Println(err)
	}
}

// SetMemberQuota 限制用户流量
func SetMemberQuota() {
	var (
		limit int
		err   error
	)
	memberList := MemberList()
	choice := util.LoopInput("请选择要限制流量的用户序号: ", memberList, true)
	if choice == -1 {
		return
	}
	for {
		quota := util.Input("请输入用户"+memberList[choice-1].Membername+"限制的流量大小(单位byte)", "")
		limit, err = strconv.Atoi(quota)
		if err != nil {
			fmt.Printf("%s 不是数字, 请重新输入!\n", quota)
		} else {
			break
		}
	}
	sqlite := core.GetSqlite()
	id := fmt.Sprint(memberList[choice-1].ID)
	if sqlite.SetQuota(id, limit) == nil {
		fmt.Println("成功设置sqlite用户" + memberList[choice-1].Membername + "限制流量" + util.Bytefmt(uint64(limit)))
	}
}

// CleanMemberData 清空用户流量
func CleanMemberData() {
	memberList := MemberList()
	choice := util.LoopInput("请选择要清空流量的用户序号: ", memberList, true)
	if choice == -1 {
		return
	}
	sqlite := core.GetSqlite()

	id := fmt.Sprint(memberList[choice-1].ID)
	if sqlite.CleanData(id) == nil {
		fmt.Println("清空sqlite流量成功!")
	}
}

// CancelMemberExpire 取消限期
func CancelMemberExpire() {
	memberList := MemberList()
	choice := util.LoopInput("请选择要取消限期的用户序号: ", memberList, true)
	if choice == -1 {
		return
	}
	if memberList[choice-1].UseDays == 0 {
		fmt.Println(util.Yellow("选择的用户未设置限期!"))
		return
	}
	sqlite := core.GetSqlite()
	id := fmt.Sprint(memberList[choice-1].ID)
	if sqlite.CancelExpire(id) == nil {
		fmt.Println("取消mysql限期成功!")
	}
}

// SetupMemberExpire 设置限期
func SetupMemberExpire() {
	memberList := MemberList()
	choice := util.LoopInput("请选择要设置限期的用户序号: ", memberList, true)
	if choice == -1 {
		return
	}
	useDayStr := util.Input("请输入要限制使用的天数: ", "")
	if useDayStr == "" {
		return
	} else if strings.Contains(useDayStr, "-") {
		fmt.Println(util.Yellow("天数不能为负数"))
		return
	} else if !util.IsInteger(useDayStr) {
		fmt.Println(util.Yellow("输入为非整数!"))
		return
	}
	useDays, _ := strconv.Atoi(useDayStr)
	sqlite := core.GetSqlite()
	id := fmt.Sprint(memberList[choice-1].ID)
	if sqlite.SetExpire(id, uint(useDays)) == nil {
		fmt.Println("设置sqlite限期成功!")
	}

}

// CleanDataByMemberName 清空指定用户流量
func CleanDataByMemberName(usernames []string) {

	sqlite := core.GetSqlite()
	if err := sqlite.CleanDataByName(usernames); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("清空sqlite流量成功!")
	}

}

// MemberList 返回并打印选定的成员指针数组
func MemberList(ids ...string) []*core.Member {
	sqlite := core.GetSqlite()
	memberList, err := sqlite.GetDataORM(ids...)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	fmt.Println(memberList)
	for i, k := range memberList {
		pass, err := base64.StdEncoding.DecodeString(k.Password)
		if err != nil {
			pass = []byte("")
		}
		fmt.Printf("%d.\n", i+1)
		fmt.Println("用户名: " + k.Membername)
		fmt.Println("密码: " + util.Green(fmt.Sprintf(string(pass))))
		fmt.Println("上传流量: " + util.Cyan(util.Bytefmt(k.Upload)))
		fmt.Println("下载流量: " + util.Cyan(util.Bytefmt(k.Download)))
		if k.Quota < 0 {
			fmt.Println("流量限额: " + util.Cyan("无限制"))
		} else {
			fmt.Println("流量限额: " + util.Cyan(util.Bytefmt(uint64(k.Quota))))
		}
		if k.UseDays == 0 {
			fmt.Println("到期日期: " + util.Cyan("无限制"))
		} else {
			fmt.Println("到期日期: " + util.Cyan(k.ExpiryDate))
		}
		fmt.Println()
	}
	return memberList
}