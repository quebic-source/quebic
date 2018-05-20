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
	"github.com/spf13/cobra"
)

var functionLogDetails bool
var functionLogFollow bool
var functionLogShowStderr bool
var functionLogShowStdout bool
var functionLogTimestamps bool
var functionLogSince string
var functionLogUntil string
var functionLogTail string

func init() {
	setupFunctionLogsFlags()
}

func setupFunctionLogsFlags() {

	functionLogsCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")

	functionLogsCmd.PersistentFlags().BoolVarP(&functionLogDetails, "details", "d", false, "detailes")
	functionLogsCmd.PersistentFlags().BoolVarP(&functionLogFollow, "follow", "f", false, "follow")
	functionLogsCmd.PersistentFlags().BoolVarP(&functionLogShowStderr, "stderr", "e", false, "stderr")
	functionLogsCmd.PersistentFlags().BoolVarP(&functionLogShowStdout, "stdout", "o", true, "stdout")
	functionLogsCmd.PersistentFlags().BoolVarP(&functionLogTimestamps, "timestamps", "t", true, "timestamps")

	functionLogsCmd.PersistentFlags().StringVarP(&functionLogSince, "since", "s", "", "since")
	functionLogsCmd.PersistentFlags().StringVarP(&functionLogUntil, "until", "u", "", "until")
	functionLogsCmd.PersistentFlags().StringVarP(&functionLogTail, "tail", "a", "", "tail")

}

var functionLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "function : logs",
	Long:  `function : logs`,
	Run: func(cmd *cobra.Command, args []string) {
		functionLogs(cmd, args)
	},
}

func functionLogs(cmd *cobra.Command, args []string) {

	/*mgrService := appContainer.GetMgrService()
	function, err := mgrService.FunctionsGetByName(functionName)
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	options := dockerTypes.ContainerLogsOptions{
		Details:    functionLogDetails,
		Follow:     functionLogFollow,
		ShowStderr: functionLogShowStderr,
		ShowStdout: functionLogShowStdout,
		Since:      functionLogSince,
		Tail:       functionLogTail,
		Timestamps: functionLogTimestamps,
		Until:      functionLogUntil,
	}

	functionService := functionutil.GetServiceID(function.Name)
	out, errO := common.DockerServiceLogs(functionService, options)
	if errO != nil {
		prepareError(cmd, errO)
	}

	io.Copy(os.Stdout, out)*/

}
