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
	"quebic-faas/quebic-faas-cli/common"
	"quebic-faas/types"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var routeSpecFile string
var routeName string

func init() {
	setupRouteCmds()
	setupRouteFlags()
}

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Route commonds",
	Long:  `Route commonds`,
}

func setupRouteCmds() {

	routeCmd.AddCommand(routeCreateCmd)
	routeCmd.AddCommand(routeUpdateCmd)
	routeCmd.AddCommand(routeGetALLCmd)
	routeCmd.AddCommand(routeInspectCmd)

}

func setupRouteFlags() {

	//route-create
	routeCreateCmd.PersistentFlags().StringVarP(&routeSpecFile, "spec", "f", "route.yml", "route input file")

	//route-update
	routeUpdateCmd.PersistentFlags().StringVarP(&routeSpecFile, "spec", "f", "route.yml", "route input file")

	//route-inspect
	routeInspectCmd.PersistentFlags().StringVarP(&routeName, "name", "n", "", "route name")

}

var routeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "route : create",
	Long:  `route : create`,
	Run: func(cmd *cobra.Command, args []string) {
		routeSave(cmd, args, true)
	},
}

var routeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "route : update",
	Long:  `route : update`,
	Run: func(cmd *cobra.Command, args []string) {
		routeSave(cmd, args, false)
	},
}

var routeGetALLCmd = &cobra.Command{
	Use:   "ls",
	Short: "route : get-all",
	Long:  `route : get-all`,
	Run: func(cmd *cobra.Command, args []string) {
		routeGetALL(cmd, args)
	},
}

var routeInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "route : inspect route details",
	Long:  `route : inspect route details`,
	Run: func(cmd *cobra.Command, args []string) {
		routeGetByName(cmd, args)
	},
}

func routeSave(cmd *cobra.Command, args []string, isAdd bool) {

	route := &types.Resource{}
	err := common.ParseYAMLFileToObject(routeSpecFile, route)
	if err != nil {
		prepareError(cmd, err)
	}

	mgrService := appContainer.GetMgrService()

	var errResponse *types.ErrorResponse
	if isAdd {
		errResponse = mgrService.RouteCreate(route)
	} else {
		errResponse = mgrService.RouteUpdate(route)
	}

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("%s route saved", route.GetID())

}

func routeGetALL(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	routes, err := mgrService.RouteGetALL()
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "URL", "Method", "Event"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(prepareRouteTable(routes))
	table.Render()

}

func routeGetByName(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	route, err := mgrService.RouteGetByName(routeName)
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	ymlStr, _ := yaml.Marshal(route)
	fmt.Printf("%s", ymlStr)

}

func prepareRouteTable(data []types.Resource) [][]string {

	var rows [][]string

	for _, val := range data {

		name := val.Name
		url := val.URL
		requestMethod := val.RequestMethod
		event := val.Event

		rows = append(rows, []string{name, url, requestMethod, event})

	}

	return rows

}
