package service

import "quebic-faas/types"

const api_mgr_components = "/mgr-components"

//ManagerComponentGetALL get all
func (mgrService *MgrService) ManagerComponentGetALL() ([]types.ManagerComponent, *types.ErrorResponse) {

	response, err := mgrService.GET(api_mgr_components, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	type managerComponentArray []types.ManagerComponent
	var managerComponents managerComponentArray
	parseResponseData(response.Data, &managerComponents)

	return managerComponents, nil

}
