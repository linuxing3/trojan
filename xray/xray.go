package xray

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"trojan/core"
	"trojan/util"
)

// ControllMenu Trojan控制菜单
func ControllMenu() {
	fmt.Println()

	menu := []string{"启动xray", "停止xray", "重启xray", "查看xray状态", "查看xray日志"}
	switch util.LoopInput("请选择: ", menu, true) {
	case 1:
		Start()
	case 2:
		Stop()
	case 3:
		Restart()
	case 4:
		Status(true)
	case 5:
		go Log(300)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		//阻塞
		<-c
	}
}

// Restart 重启xray
func Restart() {
	if err := util.ExecCommand("systemctl restart xray"); err != nil {
		fmt.Println(util.Red("重启xray失败!"))
	} else {
		fmt.Println(util.Green("重启xray成功!"))
	}
}

// Start 启动xray
func Start() {
	if err := util.ExecCommand("systemctl start xray"); err != nil {
		fmt.Println(util.Red("启动xray失败!"))
	} else {
		fmt.Println(util.Green("启动xray成功!"))
	}
}

// Stop 停止xray
func Stop() {
	if err := util.ExecCommand("systemctl stop xray"); err != nil {
		fmt.Println(util.Red("停止xray失败!"))
	} else {
		fmt.Println(util.Green("停止xray成功!"))
	}
}

// Status 获取xray状态
func Status(isPrint bool) string {
	result := util.ExecCommandWithResult("systemctl status xray")
	if isPrint {
		fmt.Println(result)
	}
	return result
}

// RunTime xray运行时间
func RunTime() string {
	result := strings.TrimSpace(util.ExecCommandWithResult("ps -Ao etime,args|grep -v grep|grep /usr/local/etc/xray/config.json"))
	resultSlice := strings.Split(result, " ")
	if len(resultSlice) > 0 {
		return resultSlice[0]
	}
	return ""
}

// Version xray原始文件的版本
func Version() string {
	flag := "-v"
	result := strings.TrimSpace(util.ExecCommandWithResult("/usr/bin/xray/xray " + flag))
	if len(result) == 0 {
		return ""
	}
	firstLine := strings.Split(result, "\n")[0]
	tempSlice := strings.Split(firstLine, " ")
	return tempSlice[len(tempSlice)-1]
}

// Type Xray类型
func Type() string {
	tType, _ := core.GetValue("xrayType")
	if tType == "" {
		_ = core.SetValue("xrayType", tType)
	}
	return tType
}

// Log 实时打印日志
func Log(line int) {
	result, _ := LogChan("-n "+strconv.Itoa(line), make(chan byte))
	for line := range result {
		fmt.Println(line)
	}
}

// LogChan 实时日志, 返回chan
func LogChan(param string, closeChan chan byte) (chan string, error) {
	cmd := exec.Command("bash", "-c", "journalctl -f -u xray "+param)

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err: ", err.Error())
		return nil, err
	}
	ch := make(chan string, 100)
	stdoutScan := bufio.NewScanner(stdout)
	go func() {
		for stdoutScan.Scan() {
			select {
			case <-closeChan:
				stdout.Close()
				return
			default:
				ch <- stdoutScan.Text()
			}
		}
	}()
	return ch, nil
}

// SetDomain 设置显示的域名
func SetDomain(domain string) {
	if domain == "" {
		domain = util.Input("请输入要显示的域名地址: ", "")
	}
	if domain == "" {
		fmt.Println("撤销更改!")
	} else {
		core.WriteDomain(domain)
		Restart()
		fmt.Println("修改domain成功!")
	}
}

// GetDomainAndPort 获取域名和端口
func GetDomainAndPort() (string, int) {
	config := core.Load("")
	return config.Inbounds[0].StreamSettings.SNI, config.Inbounds[0].Port
}
