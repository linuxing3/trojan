package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加用户",
	Run: func(cmd *cobra.Command, args []string) {
		xray.AddUser()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
