package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"quebic-faas/common"
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

//FunctionDeploy deploy function
func (mgrService *MgrService) FunctionDeploy(function *types.Function) *types.ErrorResponse {

	return mgrService.functionDeploy(function)

}

//FunctionScale scale function
func (mgrService *MgrService) FunctionScale(function *types.Function) *types.ErrorResponse {

	return mgrService.functionScale(function)

}

//FunctionTest test function
func (mgrService *MgrService) FunctionTest(functionTest *types.FunctionTest) (*types.FunctionTestResponse, *types.ErrorResponse) {

	return mgrService.functionTest(functionTest)

}

//FunctionDelete function delete
func (mgrService *MgrService) FunctionDelete(function *types.Function) *types.ErrorResponse {

	response, err := mgrService.DELETE(api_function, function, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, function)

	return nil

}

func (mgrService *MgrService) functionDeploy(function *types.Function) *types.ErrorResponse {

	response, err := mgrService.makeRequest(api_function+"/deploy", request_post, function, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, function)

	return nil
}

func (mgrService *MgrService) functionScale(function *types.Function) *types.ErrorResponse {

	response, err := mgrService.makeRequest(api_function+"/scale", request_post, function, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, function)

	return nil
}

func (mgrService *MgrService) functionTest(functionTest *types.FunctionTest) (*types.FunctionTestResponse, *types.ErrorResponse) {

	response, err := mgrService.makeRequest(api_function+"/test", request_post, functionTest, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, processErrorResponse(response)
	}

	functionTestResponse := &types.FunctionTestResponse{}
	parseResponseData(response.Data, functionTestResponse)

	return functionTestResponse, nil
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

func (mgrService *MgrService) functionSave(functionDTO *types.FunctionDTO, requestMethod string) *types.ErrorResponse {

	response, err := mgrService.makeMultipartFormRequest(functionDTO, requestMethod, nil)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return processErrorResponse(response)
	}

	parseResponseData(response.Data, functionDTO)

	return nil
}

func (mgrService *MgrService) makeMultipartFormRequest(functionDTO *types.FunctionDTO, method string, header map[string]string) (*ResponseMessage, *types.ErrorResponse) {

	url := mgrService.prepareURL(api_function)

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	//Spec data
	specDataProcess(writer, functionDTO)

	//Artifact file
	err := artifactFileProcess(writer, functionDTO.Function.Source)
	if err != nil {
		return nil, makeErrorToErrorResponse(err)
	}

	err = writer.Close()
	if err != nil {
		return nil, makeErrorToErrorResponse(err)
	}

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, makeErrorToErrorResponse(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", mgrService.Auth.AuthToken)

	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, makeErrorToErrorResponse(err)
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, makeErrorToErrorResponse(err)
	}

	return &ResponseMessage{StatusCode: res.StatusCode, Data: responseBody}, nil

}

func specDataProcess(writer *multipart.Writer, functionDTO *types.FunctionDTO) error {

	field := common.FunctionSaveField_SPEC

	functionDTOJson, err := json.Marshal(functionDTO)
	if err != nil {
		return err
	}

	return writer.WriteField(field, string(functionDTOJson))
}

func artifactFileProcess(writer *multipart.Writer, filePath string) error {

	field := common.FunctionSaveField_SOURCE

	if filePath == "" {
		return fmt.Errorf("%s file path cannot be empty", field)
	}

	artifactFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer artifactFile.Close()

	part, err := writer.CreateFormFile(field, filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, artifactFile)

	return err
}
