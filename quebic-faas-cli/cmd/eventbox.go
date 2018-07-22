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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	setupEventBoxCmds()
	setupEventBoxFlags()
}

var eventboxCmd = &cobra.Command{
	Use:   "eventbox",
	Short: "Eventbox commonds",
	Long:  `Eventbox commonds`,
}

func setupEventBoxCmds() {
	eventboxCmd.AddCommand(eventBoxStartCmd)
	eventboxCmd.AddCommand(eventBoxInfoCmd)
}

func setupEventBoxFlags() {
}

var eventBoxStartCmd = &cobra.Command{
	Use:   "start",
	Short: "eventbox : start",
	Long:  `eventbox : start`,
	Run: func(cmd *cobra.Command, args []string) {
		eventBoxStart(cmd, args)
	},
}

var eventBoxInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "eventbox : info",
	Long:  `eventbox : info`,
	Run: func(cmd *cobra.Command, args []string) {
		eventBoxInfo(cmd, args)
	},
}

func eventBoxStart(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()

	eventbox, errResponse := mgrService.EventBoxStart()
	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("eventbox running at => %s:%d", eventbox.Deployment.Host, eventbox.Deployment.Port)

}

func eventBoxInfo(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()

	eventbox, errResponse := mgrService.EventBoxInfo()
	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("eventbox running at => %s:%d", eventbox.Deployment.Host, eventbox.Deployment.Port)

}
