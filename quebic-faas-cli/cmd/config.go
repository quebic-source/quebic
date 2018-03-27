package cmd

import (
	cliconfig "quebic-faas/quebic-faas-cli/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var mgrServerHost string
var mgrServerPort int

func init() {
	setupConfigFlags()
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config cli",
	Long:  `Config cli`,
	Run: func(cmd *cobra.Command, args []string) {
		config(cmd, args)
	},
}

func setupConfigFlags() {

	defaultConfig := cliconfig.AppConfig{}
	defaultConfig.SetDefault()

	configCmd.PersistentFlags().StringVarP(&mgrServerHost, "mgr-host", "a", defaultConfig.MgrServerConfig.Host, "manager host")
	configCmd.PersistentFlags().IntVarP(&mgrServerPort, "mgr-port", "b", defaultConfig.MgrServerConfig.Port, "manager port")

}

func config(cmd *cobra.Command, args []string) {

	appContainer.GetAppConfig().MgrServerConfig.Host = mgrServerHost
	appContainer.GetAppConfig().MgrServerConfig.Port = mgrServerPort

	appContainer.SaveConfiguration()
	color.Green("configurations saved")

}
