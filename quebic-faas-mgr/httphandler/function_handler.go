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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/quebic-faas-mgr/functionutil"
	"quebic-faas/quebic-faas-mgr/functionutil/functioncreate"
	"quebic-faas/types"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

//FunctionHandler handler
func (httphandler *Httphandler) FunctionHandler(router *mux.Router) {

	db := httphandler.db
	appConfig := httphandler.config

	router.HandleFunc("/functions", func(w http.ResponseWriter, r *http.Request) {
		getAllFunctions(w, r, db, &types.Function{})
	}).Methods("GET")

	router.HandleFunc("/functions/{name}", func(w http.ResponseWriter, r *http.Request) {
		getByID(w, r, db, processRequestParmForID(r, "name", &types.Function{}))
	}).Methods("GET")

	router.HandleFunc("/functions", func(w http.ResponseWriter, r *http.Request) {

		functionDTO := &types.FunctionDTO{}
		err := processRequest(r, functionDTO)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		saveFunctionDTO(w, r, db, functionDTO, appConfig, true)

	}).Methods("POST")

	router.HandleFunc("/functions", func(w http.ResponseWriter, r *http.Request) {

		functionDTO := &types.FunctionDTO{}
		err := processRequest(r, functionDTO)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		saveFunctionDTO(w, r, db, functionDTO, appConfig, false)

	}).Methods("PUT")

}

func getAllFunctions(w http.ResponseWriter, r *http.Request, db *bolt.DB, entity types.Entity) {

	var functions []types.Function
	err := dao.GetAll(db, entity, func(k, v []byte) error {

		function := types.Function{}
		json.Unmarshal(v, &function)
		functions = append(functions, function)
		return nil
	})

	if err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if functions == nil {
		var emptyStr [0]string
		writeResponse(w, emptyStr, http.StatusOK)
	} else {
		writeResponse(w, functions, http.StatusOK)
	}

}

func saveFunctionDTO(
	w http.ResponseWriter,
	r *http.Request,
	db *bolt.DB,
	functionDTO *types.FunctionDTO,
	appConfig config.AppConfig,
	isCreate bool) {

	function := &functionDTO.Function
	route := &functionDTO.Route
	fillRouteDataFromFunction(function, route)

	//trim
	trimStringFieldsFunction(function)
	trimStringFieldsResource(route)

	errors := validateFunctionDTO(db, functionDTO, isCreate)
	if len(errors) > 0 {
		status := http.StatusBadRequest
		writeResponse(w, types.ErrorResponse{Cause: "validation-failed", Message: errors, Status: status}, status)
		return
	}

	err := saveFunction(db, functionDTO, appConfig, isCreate)
	if err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = saveRoute(db, route, appConfig, isCreate)
	if err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	writeResponse(w, functionDTO, 200)

}

func fillRouteDataFromFunction(function *types.Function, route *types.Resource) {
	//fill function's data
	route.Name = function.Name
	route.Event = prepareFunctionEvent(function.GetID())
}

func saveFunction(
	db *bolt.DB,
	functionDTO *types.FunctionDTO,
	appConfig config.AppConfig,
	isCreate bool) error {

	function := &functionDTO.Function

	err := preProcessFunction(function)
	if err != nil {
		return err
	}

	if function.GetID() == "" {
		return fmt.Errorf("unable find id field")
	}

	if isCreate {
		err = dao.Add(db, function)
		if err != nil {
			return err
		}
	} else {
		err = dao.Update(db, function)
		if err != nil {
			return err
		}
	}

	return postProcessFunction(db, appConfig.DockerConfig, function, functionDTO.Options)

}

func saveRoute(db *bolt.DB, route *types.Resource, appConfig config.AppConfig, isCreate bool) error {

	err := preProcessResource(db, route)
	if err != nil {
		return err
	}

	if route.GetID() == "" {
		return nil
	}

	if isCreate {
		err = dao.Add(db, route)
		if err != nil {
			return err
		}
	} else {
		err = dao.Update(db, route)
		if err != nil {
			return err
		}
	}

	//re-start apigateway
	restartAPIGateway(db, appConfig)

	return nil

}

func validateFunctionDTO(db *bolt.DB, functionDTO *types.FunctionDTO, isCreate bool) map[string][]string {

	var errors = make(map[string][]string)

	functionErrors := validationFunction(db, &functionDTO.Function, isCreate)
	if functionErrors != nil {
		errors["function-validation-errors"] = functionErrors
	}

	if functionDTO.Route.URL != "" {
		routeErrors := validationRoute(db, &functionDTO.Route, true, isCreate)
		if routeErrors != nil {
			errors["route-validation-errors"] = routeErrors
		}
	}

	functionDTO.Function.Route = functionDTO.Route.GetID()

	return errors

}

func trimStringFieldsFunction(function *types.Function) {
	function.Name = Trim(function.Name)
}

func validationFunction(db *bolt.DB, function *types.Function, isCreate bool) []string {

	var errors []string

	if function.Name == "" {
		errors = append(errors, "name field should not be empty")
	}

	if strings.Contains(function.Name, " ") {
		errors = append(errors, "name field not allow to contain spaces")
	}

	if function.ArtifactStoredLocation == "" {
		errors = append(errors, "artifactStoredLocation field should not be empty")
	}

	if function.ArtifactStoredLocation != "" {
		file, err := os.Open(function.ArtifactStoredLocation)
		defer file.Close()
		if err != nil {
			errors = append(errors, fmt.Sprintf("unable to found function-artifact : %s", function.ArtifactStoredLocation))
		}
	}

	if function.HandlerPath == "" {
		errors = append(errors, "handlerPath field should not be empty")
	}

	if function.Runtime == "" {
		errors = append(errors, "runtime field should not be empty")
	}

	if function.Runtime != "" {

		if !common.RuntimeValidate(common.Runtime(function.Runtime)) {
			errors = append(errors, "runtime not match")
		} else {
			err := prepareHandlerFile(function)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}

	}

	if function.Replicas < 0 {
		errors = append(errors, "replicas value is invalide")
	}

	if function.Events != nil {

		for i, eventID := range function.Events {

			event, eventErrors := createEventFromEventID(db, eventID)
			if eventErrors != nil {
				errors = append(errors, eventErrors...)
			} else {
				function.Events[i] = event.GetID()
			}

		}
	}

	if isCreate {

		if checkFunctionISAlreadyExists(db, function) {
			errors = append(errors, "function is already exists")
		}

	} else {

		if !checkFunctionISAlreadyExists(db, function) {
			errors = append(errors, "function not found")
		}

	}

	return errors

}

func checkFunctionISAlreadyExists(db *bolt.DB, function *types.Function) bool {

	found := false

	_ = dao.GetByID(db, function, func(savedObj []byte) error {

		if savedObj != nil {
			found = true
		}

		return nil
	})

	return found
}

func preProcessFunction(
	function *types.Function) error {

	secretKeyUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("unable to assign function secretKey %v", err)
	}

	function.SecretKey = secretKeyUUID.String()

	if function.Replicas <= 0 {
		function.Replicas = 1
	}

	function.Events = append(function.Events, prepareFunctionEvent(function.GetID()))

	return nil
}

func prepareFunctionEvent(functionID string) string {
	return common.EventPrefixFunction + common.EventJOIN + functionID
}

func postProcessFunction(
	db *bolt.DB,
	dockerConfig config.DockerConfig,
	function *types.Function,
	options types.FunctionCreateOptions) error {

	entityLog := types.EntityLog{State: common.LogStateSaved}
	dao.AddFunctionLog(db, function, entityLog)

	authConfig, err := dockerConfig.GetDockerAuthConfig()
	if err != nil {
		log.Print(err)
		return err
	}

	dockerImageID, err := functionutil.FunctionCreate(authConfig, function, options)

	if err != nil {
		entityLog = types.EntityLog{State: common.LogStateDockerImageCreatingFailed, Message: err.Error()}
		dao.AddFunctionLog(db, function, entityLog)
		return err
	}

	dao.AddFunctionDockerImageID(db, function, dockerImageID)
	entityLog = types.EntityLog{State: common.LogStateDockerImageCreated, Message: dockerImageID}
	dao.AddFunctionLog(db, function, entityLog)

	function.DockerImageID = dockerImageID

	return nil

}

func prepareHandlerFile(function *types.Function) error {

	runtime := function.Runtime
	functionArtifactPath := function.ArtifactStoredLocation
	functionFile := function.HandlerFile

	if runtime == common.RuntimeJava {
		// ############## JAVA ##############
		function.HandlerFile = functioncreate.GetDockerFunctionJAR()
		return nil

	} else if runtime == common.RuntimeNodeJS {

		// ############## NodeJS ##############

		//if provided artifact path is not a tar/tar.gz. put handler.js into .tar
		if filepath.Ext(functionArtifactPath) != ".tar" && filepath.Ext(functionArtifactPath) != ".gz" {
			function.HandlerFile = functioncreate.GetDockerFunctionJS()
			return nil
		}

		if functionFile == "" {
			return fmt.Errorf("handler cannot be empty for node package")
		}

		function.HandlerFile = functioncreate.GetDockerFunctionJSPackage(functionFile)
		return nil

	} else {
		return fmt.Errorf("unable to get handler file. runtime not match")
	}

}
