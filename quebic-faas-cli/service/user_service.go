package service

import (
	"quebic-faas/types"
)

const api_auth = "/auth"
const api_user = "/users"

//UserLogin user_login
func (mgrService *MgrService) UserLogin(authDTO *types.AuthDTO) (*types.JWTToken, *types.ErrorResponse) {

	response, err := mgrService.makeRequest(api_auth, request_post, authDTO, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	jwtToken := new(types.JWTToken)
	parseResponseData(response.Data, jwtToken)

	return jwtToken, nil

}

//UserCurrentAuth user_auth_current
func (mgrService *MgrService) UserCurrentAuth() (*types.User, *types.ErrorResponse) {

	response, err := mgrService.makeRequest(api_auth+"/current", request_get, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	user := new(types.User)
	parseResponseData(response.Data, user)

	return user, nil

}

//UserMe user_me
func (mgrService *MgrService) UserMe() (*types.User, *types.ErrorResponse) {

	response, err := mgrService.makeRequest(api_auth+"/me", request_get, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	user := new(types.User)
	parseResponseData(response.Data, user)

	return user, nil

}

//UserCreate user_create
func (mgrService *MgrService) UserCreate(user *types.User) *types.ErrorResponse {

	response, err := mgrService.makeRequest(api_auth+api_user, request_post, user, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, user)

	return nil

}

//UserUpdate user_update
func (mgrService *MgrService) UserUpdate(user *types.User) *types.ErrorResponse {

	response, err := mgrService.makeRequest(api_auth+api_user, request_put, user, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, user)

	return nil

}

//UserChangePassword user_change_password
func (mgrService *MgrService) UserChangePassword(user *types.User) *types.ErrorResponse {

	response, err := mgrService.makeRequest(api_auth+api_user+"/change-password", request_post, user, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, user)

	return nil

}

//UserGetAll user_get_all
func (mgrService *MgrService) UserGetAll() ([]types.User, *types.ErrorResponse) {

	response, err := mgrService.makeRequest(api_auth+api_user, request_get, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	type userArray []types.User
	var users userArray
	parseResponseData(response.Data, &users)

	return users, nil

}
