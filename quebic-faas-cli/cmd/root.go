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
	"fmt"
	"os"
	"quebic-faas/quebic-faas-cli/app"

	"github.com/spf13/cobra"
)

var appContainer app.App

func init() {
	setupCmds()
	setupFlags()
}

func setupCmds() {
	rootCmd.AddCommand(functionCmd)
	rootCmd.AddCommand(routeCmd)
	rootCmd.AddCommand(requestTrackerCmd)
	rootCmd.AddCommand(mgrCompCmd)
	rootCmd.AddCommand(apigatewayCmd)
	rootCmd.AddCommand(eventboxCmd)
	rootCmd.AddCommand(configCmd)
}

func setupFlags() {

}

var rootCmd = &cobra.Command{
	Use:   "quebic",
	Short: "quebic is a development kit for write serverless function",
	Long:  `quebic is a development kit for write serverless function. Complete documentation is available at http://quebic.io`,
}

//Execute execute cmd
func Execute() {

	appContainer = app.App{}
	appContainer.Start()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
