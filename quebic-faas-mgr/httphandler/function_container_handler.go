/*
Copyright 2018 Tharanga Nilupul Thennakoon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/quebic-faas-mgr/functionutil"
	"quebic-faas/types"

	dockerTypes "github.com/docker/docker/api/types"

	bolt "github.com/coreos/bbolt"
	"github.com/gorilla/mux"
)

//FunctionContainerHandler handler
func (httphandler *Httphandler) FunctionContainerHandler(router *mux.Router) {

	appConfig := httphandler.config
	db := httphandler.db

	router.HandleFunc("/function_containers/deploy", func(w http.ResponseWriter, r *http.Request) {

		function := &types.Function{}
		err := processRequest(r, function)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		trimStringFieldsFunction(function)

		errors := validationFunctionContainer(db, function)
		if errors != nil {
			status := http.StatusBadRequest
			writeResponse(w, types.ErrorResponse{Cause: "validation-failed", Message: errors, Status: status}, status)
			return
		}

		_, err = functionutil.FunctionDeploy(
			appConfig,
			function)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		writeResponse(w, function, 200)

	}).Methods("POST")

	router.HandleFunc("/function_containers/stop", func(w http.ResponseWriter, r *http.Request) {

		function := &types.Function{}
		err := processRequest(r, function)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		trimStringFieldsFunction(function)

		errors := validationFunctionContainer(db, function)
		if errors != nil {
			status := http.StatusBadRequest
			writeResponse(w, types.ErrorResponse{Cause: "validation-failed", Message: errors, Status: status}, status)
			return
		}

		//stop container
		err = functionutil.StopFunction(appConfig, function)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		writeResponse(w, function, 200)

	}).Methods("POST")

	router.HandleFunc("/function_containers/logs", func(w http.ResponseWriter, r *http.Request) {

		if appConfig.Deployment != config.Deployment_Docker {
			makeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("funtion-logs access only allow for docker deployment mode"))
			return
		}

		logDTO := &types.FunctionContainerLogDTO{}
		err := processRequest(r, logDTO)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		errors := validationFunctionContainer(db, &logDTO.Function)
		if errors != nil {
			status := http.StatusBadRequest
			writeResponse(w, types.ErrorResponse{Cause: "validation-failed", Message: errors, Status: status}, status)
			return
		}

		options := dockerTypes.ContainerLogsOptions{
			Details:    logDTO.Options.Details,
			Follow:     logDTO.Options.Follow,
			ShowStderr: logDTO.Options.ShowStderr,
			ShowStdout: logDTO.Options.ShowStdout,
			Since:      logDTO.Options.Since,
			Tail:       logDTO.Options.Tail,
			Timestamps: logDTO.Options.Timestamps,
			Until:      logDTO.Options.Until,
		}

		functionService := functionutil.GetServiceID(logDTO.Function.Name)
		reader, err := common.DockerServiceLogs(functionService, options)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		logsStr := buf.String()

		writeResponse(w, logsStr, 200)

	}).Methods("POST")

}

func validationFunctionContainer(db *bolt.DB, function *types.Function) []string {

	var errors []string

	if function.Name == "" {
		errors = append(errors, "function name should not be empty")
	}

	if function.Name != "" {
		err := getFunctionByID(db, function)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	return errors

}

func getFunctionByID(db *bolt.DB, function *types.Function) error {

	err := dao.GetByID(db, function, func(savedObj []byte) error {

		if savedObj == nil {
			return makeError("function not found", nil)
		}

		json.Unmarshal(savedObj, function)

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}
