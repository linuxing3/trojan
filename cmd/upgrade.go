package cmd

import (
	"fmt"
	"trojan/core"
	"trojan/xray"

	"github.com/spf13/cobra"
)

// upgradeCmd represents the update command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级数据库和xray配置文件",
}

func upgradeConfig() {
	domain, _ := core.GetValue("domain")
	if domain == "" {
		return
	}
	config := core.Load("")
	// config.SSl.Sni = domain
	core.Save(config, "")
	xray.Restart()
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.AddCommand(&cobra.Command{Use: "db", Short: "升级数据库", Run: func(cmd *cobra.Command, args []string) {
		if err := core.GetMysql().UpgradeDB(); err != nil {
			fmt.Println(err)
		}
	}})
	upgradeCmd.AddCommand(&cobra.Command{Use: "config", Short: "升级配置文件", Run: func(cmd *cobra.Command, args []string) {
		upgradeConfig()
	}})
}
