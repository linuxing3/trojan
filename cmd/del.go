package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "删除用户",
	Run: func(cmd *cobra.Command, args []string) {
		xray.DelUser()
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
