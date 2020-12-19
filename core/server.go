package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"trojan/util"
)

var configPath = "/usr/local/etc/xray/config.json"

// ServerConfig 结构体
type ServerConfig struct {
	Config
	Mysql Mysql `json: "mysql"`
}

// TrojanServerConfig 结构体
type TrojanServerConfig struct {
	TrojanConfig
	SSl   ServerSSL `json:"ssl"`
	Tcp   ServerTCP `json:"tcp"`
	Mysql Mysql     `json: "mysql"`
}

// ServerSSL 结构体
type ServerSSL struct {
	SSL
	Key                string `json:"key"`
	KeyPassword        string `json:"key_password"`
	PreferServerCipher bool   `json:"prefer_server_cipher"`
	SessionTimeout     int    `json:"session_timeout"`
	PlainHttpResponse  string `json:"plain_http_response"`
	Dhparam            string `json:"dhparam"`
}

// ServerTCP 结构体
type ServerTCP struct {
	TCP
	PreferIPv4 bool `json:"prefer_ipv4"`
}

// Load 加载服务端配置文件
func Load(path string) *ServerConfig {
	if path == "" {
		path = configPath
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	config := ServerConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(err)
		return nil
	}
	return &config
}

// Save 保存服务端配置文件
func Save(config *ServerConfig, path string) bool {
	if path == "" {
		path = configPath
	}
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Println(err)
		return false
	}
	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GetMysql 获取mysql连接
func GetMysql() *Mysql {
	config := Load("")
	return &config.Mysql
}

// WriteMysql 写mysql配置
func WriteMysql(mysql *Mysql) bool {
	mysql.Enabled = true
	config := Load("")
	config.Mysql = *mysql
	return Save(config, "")
}

// WriteTls 写tls配置
func WriteTls(cert, key, domain string) bool {
	config := Load("")
	// 入站层的设置
	certificates := config.Inbounds[0].StreamSettings.XtlsSettings.Certificates[0]
	certificates.CertificateFile = cert
	certificates.KeyFile = key
	config.Inbounds[0].StreamSettings.SNI = domain
	return Save(config, "")
}

// WriteDomain 写域名
func WriteDomain(domain string) bool {
	config := Load("")
	config.Config.Inbounds[0].StreamSettings.SNI = domain
	return Save(config, "")
}

// WriteInbloudClient 写入站客户端信息
func WriteInbloudClient(ids []string, flag string) bool {
	config := Load("")
	// 获取xray入站的client列表
	clients := config.Inbounds[0].Settings.Clients
	// 生成包括现有clinet的id的列表
	var clientsKeys []string
	for _, client := range clients {
		clientsKeys = append(clientsKeys, client.Id)
	}
	// CRUD更新
	if flag == "create" {
		fmt.Printf("如果没有就插入新的client")
		for _, id := range ids {
			if !util.Contains(clientsKeys, id) {
				var newClient InBoundSettingClientConfig
				newClient.Id = id
				newClient.Email = "love@example.com"
				newClient.Level = 0
				newClient.Flow = "xtls-rprx-direct"
				clients = append(clients, newClient)
			}
		}
	} else if flag == "delete" {
		fmt.Printf("如果有就删除client")
		for _, id := range ids {
			if util.Contains(clientsKeys, id) {
				for i, k := range clients {
					if k.Id == id {
						clients = append(clients[:i], clients[i+1:]...)
					}
				}
			}
		}
	} else {
		fmt.Printf("无操作符，忽略。。。")
	}
	fmt.Printf("重新写回到配置文件中")
	config.Inbounds[0].Settings.Clients = clients
	return Save(config, "")
}

// WriteLogLevel 写日志等级
func WriteLogLevel(level string) bool {
	config := Load("")
	config.Log.LogLevel = level
	return Save(config, "")
}
