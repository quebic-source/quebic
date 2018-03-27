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

var routeInputFile string
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
	routeCreateCmd.PersistentFlags().StringVarP(&routeInputFile, "file", "f", "route.yml", "route input file")

	//route-update
	routeUpdateCmd.PersistentFlags().StringVarP(&routeInputFile, "file", "f", "route.yml", "route input file")

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
	err := common.ParseYAMLToObject(routeInputFile, route)
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

	color.Green("route saved : %s", route.GetID())

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
