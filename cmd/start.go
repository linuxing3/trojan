package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动xray",
	Run: func(cmd *cobra.Command, args []string) {
		xray.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
