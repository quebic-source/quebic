package service

import (
	"quebic-faas/types"
)

const api_route = "/routes"

//RouteGetALL get all route
func (mgrService *MgrService) RouteGetALL() ([]types.Resource, *types.ErrorResponse) {

	response, err := mgrService.GET(api_route, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	type routeArray []types.Resource
	var routes routeArray
	parseResponseData(response.Data, &routes)

	return routes, nil

}

//RouteGetByName get route by name
func (mgrService *MgrService) RouteGetByName(name string) (*types.Resource, *types.ErrorResponse) {

	response, err := mgrService.GET(api_route+"/"+name, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	route := new(types.Resource)
	parseResponseData(response.Data, route)

	return route, nil

}

//RouteCreate create route
func (mgrService *MgrService) RouteCreate(route *types.Resource) *types.ErrorResponse {

	return mgrService.routeSave(route, request_post)

}

//RouteUpdate update route
func (mgrService *MgrService) RouteUpdate(route *types.Resource) *types.ErrorResponse {

	return mgrService.routeSave(route, request_put)

}

func (mgrService *MgrService) routeSave(route *types.Resource, requestMethod string) *types.ErrorResponse {

	response, err := mgrService.makeRequest(api_route, requestMethod, route, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, route)

	return nil
}
