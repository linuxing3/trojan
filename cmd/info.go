package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "用户信息列表",
	Run: func(cmd *cobra.Command, args []string) {
		xray.MemberList()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
