package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看xray状态",
	Run: func(cmd *cobra.Command, args []string) {
		xray.Status(true)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
