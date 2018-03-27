package service

import "quebic-faas/types"

const api_request_tracker = "/request-trackers"

//RequestTrackerGetALL get all
func (mgrService *MgrService) RequestTrackerGetALL() ([]types.RequestTracker, *types.ErrorResponse) {

	response, err := mgrService.GET(api_request_tracker, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	type requestTrackerArray []types.RequestTracker
	var requestTrackers requestTrackerArray
	parseResponseData(response.Data, &requestTrackers)

	return requestTrackers, nil

}

//RequestTrackerGetByID get logs by requestID
func (mgrService *MgrService) RequestTrackerGetByID(requestID string) (*types.RequestTracker, *types.ErrorResponse) {

	response, err := mgrService.GET(api_request_tracker+"/"+requestID, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	rt := &types.RequestTracker{}
	parseResponseData(response.Data, rt)

	return rt, nil

}
