package service

import "quebic-faas/types"

const api_event = "/events"

//EventGetALL get all events
func (mgrService *MgrService) EventGetALL() (*[]types.Event, *types.ErrorResponse) {

	response, err := mgrService.GET(api_event, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	events := new([]types.Event)
	parseResponseData(response.Data, events)

	return events, nil

}
