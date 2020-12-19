package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// upgradeCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新xray",
	Run: func(cmd *cobra.Command, args []string) {
		xray.InstallXray()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
