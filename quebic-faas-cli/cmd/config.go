//    Copyright 2018 Tharanga Nilupul Thennakoon
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

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
