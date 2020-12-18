package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "重启xray",
	Run: func(cmd *cobra.Command, args []string) {
		xray.Restart()
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
