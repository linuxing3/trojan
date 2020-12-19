package cmd

import (
	"trojan/xray"

	"github.com/spf13/cobra"
)

// tlsCmd represents the tls command
var tlsCmd = &cobra.Command{
	Use:   "tls",
	Short: "证书安装",
	Run: func(cmd *cobra.Command, args []string) {
		xray.InstallTls()
	},
}

func init() {
	rootCmd.AddCommand(tlsCmd)
}
