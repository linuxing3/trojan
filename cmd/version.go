package cmd

import (
	"fmt"
	"trojan/util"
	"trojan/xray"

	"github.com/spf13/cobra"
)

// versionCmd represents the Version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本号",
	Run: func(cmd *cobra.Command, args []string) {
		runTime := xray.RunTime()
		xrayVersion := xray.Version()
		fmt.Println()
		fmt.Printf("Version: %s\n\n", util.Cyan(xray.MVersion))
		fmt.Printf("BuildDate: %s\n\n", util.Cyan(xray.BuildDate))
		fmt.Printf("GoVersion: %s\n\n", util.Cyan(xray.GoVersion))
		fmt.Printf("GitVersion: %s\n\n", util.Cyan(xray.GitVersion))
		fmt.Printf("XrayVersion: %s\n\n", util.Cyan(xrayVersion))
		fmt.Printf("XrayRunTime: %s\n\n", util.Cyan(runTime))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
