package controller

import (
	"fmt"
	"log"
	"time"
	"trojan/core"
	websocket "trojan/util"
	"trojan/xray"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

// Start 启动xray
func Start() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	xray.Start()
	return &responseBody
}

// Stop 停止xray
func Stop() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	xray.Stop()
	return &responseBody
}

// Restart 重启xray
func Restart() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	xray.Restart()
	return &responseBody
}

// Update xray更新
func Update() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	xray.InstallXray()
	return &responseBody
}

// SetLogLevel 修改xray日志等级
func SetLogLevel(level int) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	core.WriteLogLevel(level)
	xray.Restart()
	return &responseBody
}

// GetLogLevel 获取xray日志等级
func GetLogLevel() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	config := core.Load("")
	responseBody.Data = map[string]interface{}{
		"loglevel": config.LogLevel,
	}
	return &responseBody
}

// Log 通过ws查看trojan实时日志
func Log(c *gin.Context) {
	var (
		wsConn *websocket.WsConnection
		err    error
	)
	if wsConn, err = websocket.InitWebsocket(c.Writer, c.Request); err != nil {
		fmt.Println(err)
		return
	}
	defer wsConn.WsClose()
	param := c.DefaultQuery("line", "300")
	if param == "-1" {
		param = "--no-tail"
	} else {
		param = "-n " + param
	}
	result, err := xray.LogChan(param, wsConn.CloseChan)
	if err != nil {
		fmt.Println(err)
		wsConn.WsClose()
		return
	}
	for line := range result {
		if err := wsConn.WsWrite(ws.TextMessage, []byte(line+"\n")); err != nil {
			log.Println("can't send: ", line)
			break
		}
	}
}
