package cmd

import (
	"fmt"
	"quebic-faas/types"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const function_container_op_deploy = "deploy"
const function_container_op_stop = "stop"

var functionName string

func init() {
	setupFunctionContainerFlags()
}

func setupFunctionContainerFlags() {

	functionDeployCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")
	functionStopCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")

}

var functionDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "function : deploy",
	Long:  `function : deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		functionContainerOP(cmd, args, function_container_op_deploy, functionName)
	},
}

var functionStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "function : stop",
	Long:  `function : stop`,
	Run: func(cmd *cobra.Command, args []string) {
		functionContainerOP(cmd, args, function_container_op_stop, functionName)
	},
}

func functionContainerOP(cmd *cobra.Command, args []string, op string, fID string) {

	if fID == "" {
		prepareError(cmd, fmt.Errorf("function-id should not be empty"))
	}

	function := &types.Function{Name: fID}

	mgrService := appContainer.GetMgrService()

	var errResponse *types.ErrorResponse
	if op == function_container_op_deploy {
		errResponse = mgrService.FunctionContainerDeploy(function)
		color.Green("%s function-container : deployed", function.Name)
	} else {
		errResponse = mgrService.FunctionContainerStop(function)
		color.Green("%s function-container : stopped", function.Name)
	}

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

}
