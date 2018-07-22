package service

import (
	"quebic-faas/types"
)

const api_eventbox = "/eventbox"

//EventBoxStart eventbox start
func (mgrService *MgrService) EventBoxStart() (*types.ManagerComponent, *types.ErrorResponse) {
	return mgrService.eventBoxStart()
}

//EventBoxInfo eventbox info
func (mgrService *MgrService) EventBoxInfo() (*types.ManagerComponent, *types.ErrorResponse) {
	return mgrService.eventBoxInfo()
}

func (mgrService *MgrService) eventBoxStart() (*types.ManagerComponent, *types.ErrorResponse) {

	response, err := mgrService.makeRequest(
		api_eventbox+"/start",
		request_post,
		make(map[string]string), nil)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	eventbox := &types.ManagerComponent{}
	parseResponseData(response.Data, eventbox)

	return eventbox, nil
}

func (mgrService *MgrService) eventBoxInfo() (*types.ManagerComponent, *types.ErrorResponse) {

	response, err := mgrService.makeRequest(
		api_eventbox+"/info",
		request_get,
		make(map[string]string), nil)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	componemt := &types.ManagerComponent{}
	parseResponseData(response.Data, componemt)

	return componemt, nil

}
