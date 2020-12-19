package cmd

import (
	"fmt"
	"os"
	"trojan/core"
	"trojan/util"
	"trojan/xray"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "xray",
	Run: func(cmd *cobra.Command, args []string) {
		mainMenu()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check() {
	if !util.IsExists("/usr/local/etc/xray/config.json") {
		fmt.Println("本机未安装xray, 正在自动安装...")
		xray.InstallXray()
		core.WriteInbloudClient(nil, "create")
		xray.InstallTls()
		xray.InstallMysql(xray.XrayDbDockerRun, "xray")
		util.ExecCommand("systemctl restart xray-web")
	}
}

func mainMenu() {
	check()
exit:
	for {
		fmt.Println()
		fmt.Println(util.Cyan("欢迎使用xray管理程序"))
		fmt.Println()
		menuList := []string{"xray管理", "用户管理", "安装管理", "web管理", "查看配置", "生成json"}
		switch util.LoopInput("请选择: ", menuList, false) {
		case 1:
			xray.ControllMenu()
		case 2:
			xray.UserMenu()
		case 3:
			xray.InstallMenu()
		case 4:
			xray.WebMenu()
		case 5:
			xray.UserList()
		case 6:
			xray.GenClientJson()
		default:
			break exit
		}
	}
}
