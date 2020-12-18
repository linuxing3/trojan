package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止xray",
	Run: func(cmd *cobra.Command, args []string) {
		xray.Stop()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
