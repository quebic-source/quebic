package service

import (
	"quebic-faas/types"
)

const api_function = "/functions"

//FunctionCreate create function
func (mgrService *MgrService) FunctionCreate(functionDTO *types.FunctionDTO) *types.ErrorResponse {

	return mgrService.functionSave(functionDTO, request_post)

}

//FunctionUpdate update function
func (mgrService *MgrService) FunctionUpdate(functionDTO *types.FunctionDTO) *types.ErrorResponse {

	return mgrService.functionSave(functionDTO, request_put)

}

func (mgrService *MgrService) functionSave(functionDTO *types.FunctionDTO, requestMethod string) *types.ErrorResponse {

	response, err := mgrService.makeRequest(api_function, requestMethod, functionDTO, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, functionDTO)

	return nil
}

//FunctionsGetALL get all functions
func (mgrService *MgrService) FunctionsGetALL() ([]types.Function, *types.ErrorResponse) {

	response, err := mgrService.GET(api_function, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	type functionArray []types.Function
	var functions functionArray
	parseResponseData(response.Data, &functions)

	return functions, nil

}

//FunctionsGetByName get by name
func (mgrService *MgrService) FunctionsGetByName(name string) (*types.Function, *types.ErrorResponse) {

	response, err := mgrService.GET(api_function+"/"+name, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	function := &types.Function{}
	parseResponseData(response.Data, function)

	return function, nil

}
