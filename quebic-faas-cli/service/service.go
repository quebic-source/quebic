package service

import (
	"encoding/json"
	"fmt"
	"quebic-faas/common"
	quebicFaas "quebic-faas/config"
	"quebic-faas/quebic-faas-cli/config"
	"quebic-faas/types"
)

//MgrService mgr-service
type MgrService struct {
	MgrServerConfig quebicFaas.ServerConfig
	Auth            config.AuthConfig
}

//ResponseMessage response
type ResponseMessage struct {
	StatusCode int
	Data       []byte
}

//getApiURL path=> /users/....
func (mgrService *MgrService) prepareURL(path string) string {

	if mgrService.MgrServerConfig.Port == 0 {
		return "http://" + mgrService.MgrServerConfig.Host + path
	}

	return "http://" + mgrService.MgrServerConfig.Host + ":" + common.IntToStr(mgrService.MgrServerConfig.Port) + path
}

func parseResponseData(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

func makeError(format string, err error) error {

	if err != nil {
		return fmt.Errorf(format, err)
	}

	return fmt.Errorf(format)

}

func makeErrorToErrorResponse(err error) *types.ErrorResponse {
	return &types.ErrorResponse{Cause: err.Error()}
}

func processErrorResponse(responseMessage *ResponseMessage) *types.ErrorResponse {
	errResponse := &types.ErrorResponse{}
	parseResponseData(responseMessage.Data, errResponse)
	return errResponse
}
