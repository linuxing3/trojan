package cmd

import (
	"trojan/util"

	"github.com/spf13/cobra"
)

// updateWebCmd represents the update command
var updateWebCmd = &cobra.Command{
	Use:   "updateWeb",
	Short: "更新xray管理程序",
	Run: func(cmd *cobra.Command, args []string) {
		util.RunWebShell("https://git.io/xray-install")
	},
}

func init() {
	rootCmd.AddCommand(updateWebCmd)
}
