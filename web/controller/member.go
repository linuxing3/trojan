package controller

import (
	"encoding/base64"
	"fmt"
	"time"
	"trojan/core"
	"trojan/xray"

	"github.com/google/uuid"
)

// MemberList 获取用户列表
func MemberList(findMember string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sqlite := core.GetSqlite()
	memberList := sqlite.GetDataORM()
	if findMember != "" {
		for _, member := range memberList {
			if member.Membername == findMember {
				memberList = []core.Member{member}
				break
			}
		}
	}
	domain, port := xray.GetDomainAndPort()
	responseBody.Data = map[string]interface{}{
		"domain":   domain,
		"port":     port,
		"memberList": memberList,
	}
	return &responseBody
}

// PageMemberList 分页查询获取用户列表
func PageMemberList(curPage int, pageSize int) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sqlite := core.GetSqlite()
	pageData, err := sqlite.PageList(curPage, pageSize)
	if err != nil {
		responseBody.Msg = err.Error()
		return &responseBody
	}
	domain, port := xray.GetDomainAndPort()
	responseBody.Data = map[string]interface{}{
		"domain":   domain,
		"port":     port,
		"pageData": pageData,
	}
	return &responseBody
}

// CreateMember 创建用户
func CreateMember(membername string, password string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	if membername == "admin" {
		responseBody.Msg = "不能创建用户名为admin的用户!"
		return &responseBody
	}
	sqlite := core.GetSqlite()
	if member := sqlite.GetMemberByName(membername); member != nil {
		responseBody.Msg = "已存在用户。[用户名]:" + membername + "。[uuid]:" + uuid
		return &responseBody
	}
	pass, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		responseBody.Msg = "Base64解码失败: " + err.Error()
		return &responseBody
	}
	if member := sqlite.GetMemberByPass(password); member != nil {
		responseBody.Msg = "已存在密码为: " + string(pass) + " 的用户!"
		return &responseBody
	}
	if err := sqlite.CreateMember("", membername, password, string(pass)); err != nil {
		responseBody.Msg = err.Error()
	}
	return &responseBody
}

// UpdateMember 更新用户
func UpdateMember(id string, membername string, password string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	if membername == "admin" {
		responseBody.Msg = "不能更改用户名为admin的用户!"
		return &responseBody
	}
	sqlite := core.GetSqlite()
	memberList := sqlite.GetDataORM(id)
	if memberList[0].Membername != membername {
		if member := sqlite.GetMemberByName(membername); member != nil {
			responseBody.Msg = "已存在用户名为: " + membername + " 的用户!"
			return &responseBody
		}
	}
	pass, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		responseBody.Msg = "Base64解码失败: " + err.Error()
		return &responseBody
	}
	if memberList[0].Password != password {
		if member := sqlite.GetMemberByPass(password); member != nil {
			responseBody.Msg = "已存在密码为: " + string(pass) + " 的用户!"
			return &responseBody
		}
	}
	if err := sqlite.UpdateMember(id, membername, password, string(pass)); err != nil {
		responseBody.Msg = err.Error()
	}
	return &responseBody
}

// DelMember 删除用户
func DelMember(id string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sqlite := core.GetSqlite()
	sqlite.DeleteMemberORM(id)
	responseBody.Msg = "Deleted"
	return &responseBody
}

// SetExpire 设置用户过期
func SetMemberExpire(id string, useDays uint) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sqlite := core.GetSqlite()
	if err := sqlite.SetExpire(id, useDays); err != nil {
		responseBody.Msg = err.Error()
	}
	return &responseBody
}

// CancelExpire 取消设置用户过期
func CancelMemberExpire(id string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sqlite := core.GetSqlite()
	if err := sqlite.CancelExpire(id); err != nil {
		responseBody.Msg = err.Error()
	}
	return &responseBody
}
