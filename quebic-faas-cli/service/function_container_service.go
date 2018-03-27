package service

import (
	"quebic-faas/types"
)

const api_function_container = "/function_containers"

//FunctionContainerDeploy function container create-or-update
func (mgrService *MgrService) FunctionContainerDeploy(function *types.Function) *types.ErrorResponse {

	response, err := mgrService.POST(api_function_container+"/deploy", function, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, function)

	return nil

}

//FunctionContainerStop function container stop
func (mgrService *MgrService) FunctionContainerStop(function *types.Function) *types.ErrorResponse {

	response, err := mgrService.POST(api_function_container+"/stop", function, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, function)

	return nil

}
