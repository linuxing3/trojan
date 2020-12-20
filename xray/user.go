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

// UserMenu 用户管理菜单
func UserMenu() {
	fmt.Println()
	menu := []string{"新增用户", "删除用户", "限制流量", "清空流量", "设置限期", "取消限期"}
	switch util.LoopInput("请选择: ", menu, false) {
	case 1:
		AddUser()
	case 2:
		DelUser()
	case 3:
		SetUserQuota()
	case 4:
		CleanData()
	case 5:
		SetupExpire()
	case 6:
		CancelExpire()
	}
}

// AddUser 添加用户
func AddUser() {

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

	// 3. Optional 创建Sqlite新用户
	// FIXED 这里的配置是用硬盘配置文件中读取的，所以记得先写入配置文件才能正常使用
	sqlite := core.GetSqlite()
	if err := sqlite.CreateUser(uuid, inputUser, base64Pass, inputPass); err == nil {
		fmt.Println("新增Sqlite用户成功!")
		fmt.Println("")
	} else {
		fmt.Println(err)
	}

	// 4. 获取数据库配置
	mysql := core.GetMysql()
	if user := mysql.GetUserByName(inputUser); user != nil {
		fmt.Println(util.Yellow("已存在用户。[用户名]:" + inputUser + "。[uuid]:" + uuid))
		return
	}
	if user := mysql.GetUserByPass(base64Pass); user != nil {
		fmt.Println(util.Yellow("已存在密码为: " + inputPass + " 的用户!"))
		return
	}

	// 创建Mysql新用户
	if err := mysql.CreateUser(uuid, inputUser, base64Pass, inputPass); err == nil {
		fmt.Println("新增Mysql用户成功!")
		fmt.Println("")
	} else {
		fmt.Println(err)
	}

}

// DelUser 删除用户
func DelUser() {

	userList := UserList()
	mysql := core.GetMysql()

	choice := util.LoopInput("请选择要删除的用户序号: ", userList, true)
	if choice == -1 {
		return
	}
	// FIXME
	if mysql.DeleteUser(userList[choice-1].ID) == nil {
		fmt.Println("删除mysql用户成功!")
	}

	// Sqlite用户
	sqlite := core.GetSqlite()
	if sqlite.DeleteUser(userList[choice-1].ID) == nil {
		fmt.Println("删除sqlite用户成功!")
	}
}

// SetUserQuota 限制用户流量
func SetUserQuota() {
	var (
		limit int
		err   error
	)
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("请选择要限制流量的用户序号: ", userList, true)
	if choice == -1 {
		return
	}
	for {
		quota := util.Input("请输入用户"+userList[choice-1].Username+"限制的流量大小(单位byte)", "")
		limit, err = strconv.Atoi(quota)
		if err != nil {
			fmt.Printf("%s 不是数字, 请重新输入!\n", quota)
		} else {
			break
		}
	}
	if mysql.SetQuota(userList[choice-1].ID, limit) == nil {
		fmt.Println("成功设置mysql用户" + userList[choice-1].Username + "限制流量" + util.Bytefmt(uint64(limit)))
	}

	sqlite := core.GetSqlite()
	if sqlite.SetQuota(userList[choice-1].ID, limit) == nil {
		fmt.Println("成功设置sqlite用户" + userList[choice-1].Username + "限制流量" + util.Bytefmt(uint64(limit)))
	}
}

// CleanData 清空用户流量
func CleanData() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("请选择要清空流量的用户序号: ", userList, true)
	if choice == -1 {
		return
	}
	if mysql.CleanData(userList[choice-1].ID) == nil {
		fmt.Println("清空mysql流量成功!")
	}
	sqlite := core.GetSqlite()
	if sqlite.CleanData(userList[choice-1].ID) == nil {
		fmt.Println("清空sqlite流量成功!")
	}
}

// CancelExpire 取消限期
func CancelExpire() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("请选择要取消限期的用户序号: ", userList, true)
	if choice == -1 {
		return
	}
	if userList[choice-1].UseDays == 0 {
		fmt.Println(util.Yellow("选择的用户未设置限期!"))
		return
	}
	if mysql.CancelExpire(userList[choice-1].ID) == nil {
		fmt.Println("取消mysql限期成功!")
	}
	sqlite := core.GetSqlite()
	if sqlite.CancelExpire(userList[choice-1].ID) == nil {
		fmt.Println("取消sqlite限期成功!")
	}
}

// SetupExpire 设置限期
func SetupExpire() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("请选择要设置限期的用户序号: ", userList, true)
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
	if mysql.SetExpire(userList[choice-1].ID, uint(useDays)) == nil {
		fmt.Println("设置mysql限期成功!")
	}
	sqlite := core.GetSqlite()
	if sqlite.SetExpire(userList[choice-1].ID, uint(useDays)) == nil {
		fmt.Println("设置sqlite限期成功!")
	}

}

// CleanDataByName 清空指定用户流量
func CleanDataByName(usernames []string) {
	mysql := core.GetMysql()
	if err := mysql.CleanDataByName(usernames); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("清空mysql流量成功!")
	}

	sqlite := core.GetSqlite()
	if err := sqlite.CleanDataByName(usernames); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("清空sqlite流量成功!")
	}

}

// UserList 获取用户列表并打印显示
func UserList(ids ...string) []*core.User {
	var (
		userList []*core.User
		err      error
	)
	// mysql 作为主数据库
	mysql := core.GetMysql()
	if mysql == nil {
		// sqlite 作为备用
		sqlite := core.GetSqlite()
		userList, err = sqlite.GetData(ids...)
	} else {
		userList, err = mysql.GetData(ids...)
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	domain, port := GetDomainAndPort()
	for i, k := range userList {
		pass, err := base64.StdEncoding.DecodeString(k.Password)
		if err != nil {
			pass = []byte("")
		}
		fmt.Printf("%d.\n", i+1)
		fmt.Println("用户名: " + k.Username)
		fmt.Println("客户端id: " + k.ID)
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
		fmt.Println("分享链接: " + util.Green(fmt.Sprintf("xray://%s@%s:%d", string(pass), domain, port)))
		fmt.Println()
	}
	return userList
}
