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
	setupAPIGatewayCmds()
	setupAPIGatewayFlags()
}

var apigatewayCmd = &cobra.Command{
	Use:   "api-gateway",
	Short: "ApiGateway commonds",
	Long:  `ApiGateway commonds`,
}

func setupAPIGatewayCmds() {
	apigatewayCmd.AddCommand(apigatewayInfoCmd)
}

func setupAPIGatewayFlags() {
}

var apigatewayInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "api-gateway : info",
	Long:  `api-gateway : info`,
	Run: func(cmd *cobra.Command, args []string) {
		getAPIGatewayInfo(cmd, args)
	},
}

func getAPIGatewayInfo(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	comps, err := mgrService.ManagerComponentGetALL()
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	apigateway := comps[0]

	color.Green("api-gateway running at => %s:%d", apigateway.Deployment.Host, apigateway.Deployment.Port)

}
