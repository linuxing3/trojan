package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var configPath = "/usr/local/etc/xray/config.json"

// ServerConfig 结构体
type ServerConfig struct {
	Config
	SSl   ServerSSL `json:"ssl"`
	Tcp   ServerTCP `json:"tcp"`
	Mysql Mysql     `json: "mysql"`
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

// LoadTrojanConfig 加载服务端配置文件
func LoadTrojanConfig(path string) *TrojanServerConfig {
	if path == "" {
		path = configPath
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	config := TrojanServerConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(err)
		return nil
	}
	return &config
}

// SaveTrojanConfig 保存服务端配置文件
func SaveTrojanConfig(config *TrojanServerConfig, path string) bool {
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
	// 第一层
	config.SSl.Cert = cert
	config.SSl.Key = key
	config.SSl.Sni = domain
	// 入站层的设置
	config.Inbounds[0].StreamSettings.XtlsSettings.Certificates[0].CertificateFile = cert
	config.Inbounds[0].StreamSettings.XtlsSettings.Certificates[0].KeyFile = key
	config.Inbounds[0].StreamSettings.SNI = domain
	return Save(config, "")
}

// WriteDomain 写域名
func WriteDomain(domain string) bool {
	config := Load("")
	config.SSl.Sni = domain
	config.Config.Inbounds[0].StreamSettings.SNI = domain
	return Save(config, "")
}

// WritePassword 写密码
func WritePassword(pass []string) bool {
	config := Load("")
	config.Inbounds[0].Settings[0].Clients[0].Id = pass[0]
	return Save(config, "")
}

// WriteLogLevel 写日志等级
func WriteLogLevel(level string) bool {
	config := Load("")
	config.Log.LogLevel = level
	return Save(config, "")
}
