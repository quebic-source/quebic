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
	"os"
	"quebic-faas/quebic-faas-cli/common"
	"quebic-faas/types"

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
		color.Red("\ncause: %v", errorResponse.Cause)
		if errorResponse.Message != nil {
			/*checkObj := make(map[string]interface{})
			if reflect.TypeOf(errorResponse.Message).Kind() == reflect.TypeOf(checkObj).Kind() {
				mapObj := errorResponse.Message.(map[string]interface{})
				for key, val := range mapObj {
					color.Red("%s : %s", key, val)
				}

			} else {
				color.Red("message : %v", errorResponse.Message)
			}*/

			ymlStr, err := common.ParseObjectToYAML(errorResponse.Message)
			if err != nil {
				color.Red("message : %v", errorResponse.Message)
			} else {
				color.Red("%s", ymlStr)
			}

		}
	}

	os.Exit(1)
}
