package cmd

import (
	"fmt"
	"os"
	"quebic-faas/types"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var requestID string

func init() {
	setupRequestTrackerCmds()
	setupRequestTrackerFlags()

}

var requestTrackerCmd = &cobra.Command{
	Use:   "request-tracker",
	Short: "Request-Tracker commonds",
	Long:  `Request-Tracker commonds`,
}

func setupRequestTrackerCmds() {

	requestTrackerCmd.AddCommand(requestTrackerGetALLCmd)
	requestTrackerCmd.AddCommand(requestTrackerInspectCmd)
	requestTrackerCmd.AddCommand(requestTrackerLogsCmd)

}

func setupRequestTrackerFlags() {

	requestTrackerInspectCmd.PersistentFlags().StringVarP(&requestID, "request-id", "i", "", "request id")
	requestTrackerLogsCmd.PersistentFlags().StringVarP(&requestID, "request-id", "i", "", "request id")

}

var requestTrackerGetALLCmd = &cobra.Command{
	Use:   "ls",
	Short: "request-tracker : get-all",
	Long:  `request-tracker : get-all`,
	Run: func(cmd *cobra.Command, args []string) {
		requestTrackerGetALL(cmd, args)
	},
}

var requestTrackerInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "request-tracker : inspect",
	Long:  `request-tracker : inspect`,
	Run: func(cmd *cobra.Command, args []string) {
		requestTrackerGetByRequestID(cmd, args)
	},
}

var requestTrackerLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "request-tracker : logs",
	Long:  `request-tracker : logs`,
	Run: func(cmd *cobra.Command, args []string) {
		requestTrackerLogs(cmd, args)
	},
}

func requestTrackerGetALL(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	requestTrackers, err := mgrService.RequestTrackerGetALL()
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Request_ID", "Source", "Created"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(prepareRequestTrackerTable(requestTrackers))
	table.Render()

}

func requestTrackerGetByRequestID(cmd *cobra.Command, args []string) {

	if requestID == "" {
		prepareErrorResponse(cmd, &types.ErrorResponse{Cause: "request-id is empty"})
	}

	mgrService := appContainer.GetMgrService()
	rt, err := mgrService.RequestTrackerGetByID(requestID)
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	ymlStr, _ := yaml.Marshal(rt)
	fmt.Printf("%s", ymlStr)

}

func requestTrackerLogs(cmd *cobra.Command, args []string) {

	if requestID == "" {
		prepareErrorResponse(cmd, &types.ErrorResponse{Cause: "request-id is empty"})
	}

	mgrService := appContainer.GetMgrService()
	rt, err := mgrService.RequestTrackerGetByID(requestID)
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Request_ID", "Type", "Time", "Message", "Source"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(prepareRequestTrackerLogsTable(*rt))
	table.Render()

}

func prepareRequestTrackerTable(data []types.RequestTracker) [][]string {

	var rows [][]string

	for _, val := range data {

		requestID := val.RequestID
		source := val.Source
		createdAt := val.CreatedAt

		rows = append(rows, []string{requestID, source, createdAt})

	}

	return rows

}

func prepareRequestTrackerLogsTable(requestTracker types.RequestTracker) [][]string {

	var rows [][]string

	for _, val := range requestTracker.Logs {

		requestID := requestTracker.RequestID
		logType := val.Type
		time := val.Time
		message := val.Message
		source := val.Source

		rows = append(rows, []string{requestID, logType, time, message, source})

	}

	return rows

}
