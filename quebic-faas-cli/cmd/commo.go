package cmd

import (
	"os"
	"quebic-faas/types"
	"reflect"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

//PrepareError prepare error
func prepareError(cmd *cobra.Command, err error) {

	color.Red("%s failed", cmd.CommandPath())
	color.Red("cause: %v", err)

	os.Exit(1)
}

//PrepareErrorResponse prepare errorResponse
func prepareErrorResponse(cmd *cobra.Command, errorResponse *types.ErrorResponse) {

	color.Red("%s failed", cmd.CommandPath())

	if errorResponse != nil {
		//color.Red("status : %v", errorResponse.Status)
		color.Red("cause: %v", errorResponse.Cause)
		if errorResponse.Message != nil {
			checkObj := make(map[string]interface{})
			if reflect.TypeOf(errorResponse.Message).Kind() == reflect.TypeOf(checkObj).Kind() {
				mapObj := errorResponse.Message.(map[string]interface{})
				for key, val := range mapObj {
					color.Red("%s : %s", key, val)
				}

			} else {
				color.Red("message : %v", errorResponse.Message)
			}

		}
	}

	os.Exit(1)
}
