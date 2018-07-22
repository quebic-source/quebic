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
	"encoding/json"
	"fmt"
	"os"
	quebic_common "quebic-faas/common"
	"quebic-faas/quebic-faas-cli/common"
	"quebic-faas/types"
	"reflect"

	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var functionInputFile string
var functionName string
var functionVersion string
var functionReplicas int
var functionTestPayload string
var functionStart bool

func init() {
	setupFunctionCmds()
	setupFunctionFlags()
}

var functionCmd = &cobra.Command{
	Use:   "function",
	Short: "Function commonds",
	Long:  `Function commonds`,
}

func setupFunctionCmds() {

	//function
	functionCmd.AddCommand(functionCreateCmd)
	functionCmd.AddCommand(functionUpdateCmd)
	functionCmd.AddCommand(functionDeployCmd)
	functionCmd.AddCommand(functionScaleCmd)
	functionCmd.AddCommand(functionTestCmd)
	functionCmd.AddCommand(functionDeleteCmd)
	functionCmd.AddCommand(functionGetALLCmd)
	functionCmd.AddCommand(functionInspectCmd)

	//function-logs
	functionCmd.AddCommand(functionLogsCmd)

}

func setupFunctionFlags() {

	//function-create
	functionCreateCmd.PersistentFlags().StringVarP(&functionInputFile, "file", "f", "function.yml", "function spec file")
	functionCreateCmd.PersistentFlags().BoolVarP(&functionStart, "start", "s", true, "if true function-container will start. otherwise not")

	//function-update
	functionUpdateCmd.PersistentFlags().StringVarP(&functionInputFile, "file", "f", "function.yml", "function spec file")
	functionUpdateCmd.PersistentFlags().BoolVarP(&functionStart, "start", "s", true, "if true function-container will start. otherwise not")

	//function-deploy
	functionDeployCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")
	functionDeployCmd.PersistentFlags().StringVarP(&functionVersion, "version", "v", "", "function version")

	//function-scale
	functionScaleCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")
	functionScaleCmd.PersistentFlags().IntVarP(&functionReplicas, "replicas", "r", 1, "function replicas value")

	//function-deploy
	functionTestCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")
	functionTestCmd.PersistentFlags().StringVarP(&functionTestPayload, "payload", "p", "", "test payload")

	//function-delete
	functionDeleteCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")

	//function-inspect
	functionInspectCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")

}

var functionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "function : create",
	Long:  `function : create`,
	Run: func(cmd *cobra.Command, args []string) {
		functionSave(cmd, args, true)
	},
}

var functionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "function : update",
	Long:  `function : update`,
	Run: func(cmd *cobra.Command, args []string) {
		functionSave(cmd, args, false)
	},
}

var functionDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "function : deploy",
	Long:  `function : deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		functionDeploy(cmd, args)
	},
}

var functionScaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "function : scale",
	Long:  `function : scale`,
	Run: func(cmd *cobra.Command, args []string) {
		functionScale(cmd, args)
	},
}

var functionTestCmd = &cobra.Command{
	Use:   "test",
	Short: "function : test",
	Long:  `function : test`,
	Run: func(cmd *cobra.Command, args []string) {
		functionTest(cmd, args)
	},
}

var functionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "function : delete",
	Long:  `function : delete`,
	Run: func(cmd *cobra.Command, args []string) {
		functionDelete(cmd, args, functionName)
	},
}

var functionGetALLCmd = &cobra.Command{
	Use:   "ls",
	Short: "function : get-all",
	Long:  `function : get-all`,
	Run: func(cmd *cobra.Command, args []string) {
		functionGetALL(cmd, args)
	},
}

var functionInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "function : inspect function details",
	Long:  `function : inspect function details`,
	Run: func(cmd *cobra.Command, args []string) {
		functionGetByName(cmd, args)
	},
}

func functionSave(cmd *cobra.Command, args []string, isAdd bool) {

	functionDTO := &types.FunctionDTO{}
	err := common.ParseYAMLFileToObject(functionInputFile, functionDTO)
	if err != nil {
		prepareError(cmd, err)
	}

	mgrService := appContainer.GetMgrService()

	var errResponse *types.ErrorResponse
	if isAdd {
		errResponse = mgrService.FunctionCreate(functionDTO)
	} else {
		errResponse = mgrService.FunctionUpdate(functionDTO)
	}

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("%s:%s function is saved", functionDTO.Function.GetID(), functionDTO.Function.Version)

}

func functionDeploy(cmd *cobra.Command, args []string) {

	function := &types.Function{Name: functionName, Version: functionVersion}

	mgrService := appContainer.GetMgrService()

	errResponse := mgrService.FunctionDeploy(function)

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("%s:%s function is deployed", function.GetID(), function.Version)

}

func functionScale(cmd *cobra.Command, args []string) {

	function := &types.Function{Name: functionName, Replicas: functionReplicas}

	mgrService := appContainer.GetMgrService()

	errResponse := mgrService.FunctionScale(function)

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("%s:%s function is scaled", function.GetID(), function.Version)

}

func functionTest(cmd *cobra.Command, args []string) {

	payloadMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(functionTestPayload), &payloadMap)
	if err != nil {
		prepareErrorResponse(cmd,
			&types.ErrorResponse{
				Status:  400,
				Cause:   "invalid-request",
				Message: "please enter valid json formatted payload",
			},
		)
		return
	}

	functionTest := &types.FunctionTest{Name: functionName, Payload: payloadMap}

	mgrService := appContainer.GetMgrService()

	functionTestResponse, errResponse := mgrService.FunctionTest(functionTest)

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	jsonResponse, _ := json.Marshal(functionTestResponse)

	color.Green("%s", jsonResponse)

}

func functionDelete(cmd *cobra.Command, args []string, fID string) {

	if fID == "" {
		prepareError(cmd, fmt.Errorf("function-id should not be empty"))
	}

	function := &types.Function{Name: fID}

	mgrService := appContainer.GetMgrService()

	errResponse := mgrService.FunctionDelete(function)
	color.Green("%s function is deleted", function.Name)

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

}

func functionGetALL(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	functions, err := mgrService.FunctionsGetALL()
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"Runtime",
		"Replicas",
		"Route",
		"Status",
	})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(prepareFunctionsTable(functions))
	table.Render()

}

func functionGetByName(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	function, err := mgrService.FunctionsGetByName(functionName)
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	ymlStr, _ := yaml.Marshal(function)
	fmt.Printf("%s", ymlStr)

}

func prepareFunctionsTable(data []types.Function) [][]string {

	var rows [][]string

	for _, val := range data {

		name := val.Name
		runtime := val.Runtime
		replicas := quebic_common.IntToStr(val.Replicas)
		route := val.Route
		status := val.Status

		rows = append(rows, []string{name, runtime, replicas, route, status})

	}

	return rows

}

func getField(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}
