package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gobuffalo/packr/v2"
)

// ClientConfig 结构体
type ClientConfig struct {
	Config
}

type TrojanClientConfig struct {
	TrojanConfig
	SSl ClientSSL `json:"ssl"`
	Tcp ClientTCP `json:"tcp"`
}

// ClientSSL 结构体
type ClientSSL struct {
	SSL
	Verify         bool `json:"verify"`
	VerifyHostname bool `json:"verify_hostname"`
}

// ClientTCP 结构体
type ClientTCP struct {
	TCP
}

// WriteClient 生成客户端json
func WriteClient(port int, id, domain, writePath string) bool {
	box := packr.New("client.json", "../asset")
	data, err := box.Find("client.json")
	if err != nil {
		fmt.Println(err)
		return false
	}
	config := ClientConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(err)
		return false
	}
	// 写入客户端配置文件
	config.Inbounds[0].Port = port
	config.Inbounds[0].Settings.Clients[0].Id = id
	outData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Println(err)
		return false
	}
	if err = ioutil.WriteFile(writePath, outData, 0644); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// WriteClient 生成客户端json
func WriteTrojanClient(port int, id, domain, writePath string) bool {
	box := packr.New("client.json", "../asset")
	data, err := box.Find("client.json")
	if err != nil {
		fmt.Println(err)
		return false
	}
	config := ClientConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(err)
		return false
	}
	config.Inbounds[0].Port = port
	config.Inbounds[0].Settings.Clients[0].Id = id
	outData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Println(err)
		return false
	}
	if err = ioutil.WriteFile(writePath, outData, 0644); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
