//    Copyright 2018 Tharanga Nilupul Thennakoon
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_messenger "quebic-faas/messenger"
	"quebic-faas/quebic-faas-mgr/components"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/dao"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/quebic-faas-mgr/logger"
	"quebic-faas/types"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/gorilla/mux"
)

//Httphandler handlers
type Httphandler struct {
	config     config.AppConfig
	db         *bolt.DB
	messenger  _messenger.Messenger
	loggerUtil logger.Logger
	deployment dep.Deployment
}

//SetUpHTTPHandlers setUpHTTPHandlers
func SetUpHTTPHandlers(
	config config.AppConfig,
	router *mux.Router,
	db *bolt.DB,
	messenger _messenger.Messenger,
	loggerUtil logger.Logger,
	deployment dep.Deployment) {

	http := &Httphandler{
		config:     config,
		db:         db,
		messenger:  messenger,
		loggerUtil: loggerUtil,
		deployment: deployment,
	}
	http.EventHandler(router)
	http.ResourceHandler(router)
	http.FunctionHandler(router)
	http.MessengerHandler(router)
	http.AuthHandler(router)
	http.ApigatewayDataServe(router)
	http.MgrComponentHandler(router)
	http.RequestTrackerHandler(router)
	http.EventBoxHandler(router)

}

func getByID(w http.ResponseWriter, r *http.Request, db *bolt.DB, entity types.Entity) {

	if entity.GetID() == "" {
		makeErrorResponse(w, http.StatusBadRequest, makeError("id is empty", nil))
		return
	}

	err := dao.GetByID(db, entity, func(savedObj []byte) error {

		if savedObj == nil {
			return makeError("resource not found", nil)
		}

		json.Unmarshal(savedObj, entity)

		writeResponse(w, entity, http.StatusOK)
		return nil
	})

	if err != nil {
		makeErrorResponse(w, http.StatusNotFound, err)
	}

}

func add(w http.ResponseWriter, r *http.Request, db *bolt.DB, entity types.Entity) {

	if entity.GetID() == "" {
		makeErrorResponse(w, http.StatusBadRequest, makeError("unable find id field", nil))
		return
	}

	err := dao.Add(db, entity)
	if err != nil {
		makeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	writeResponse(w, entity, http.StatusCreated)

}

func update(w http.ResponseWriter, r *http.Request, db *bolt.DB, entity types.Entity) {

	if entity.GetID() == "" {
		makeErrorResponse(w, http.StatusBadRequest, makeError("unable find id field", nil))
		return
	}

	err := dao.Update(db, entity)
	if err != nil {
		makeErrorResponse(w, http.StatusNotFound, err)
		return
	}

	writeResponse(w, entity, http.StatusAccepted)

}

func processRequest(r *http.Request, entity interface{}) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return makeError("unable to read request %v", err)
	}

	err = json.Unmarshal(body, entity)
	if err != nil {
		return makeError("unable to parse json request to entity %v", err)
	}

	return nil

}

func processRequestParmForID(r *http.Request, parmkey string, entity types.Entity) types.Entity {
	params := mux.Vars(r)
	entity.SetID(params[parmkey])
	return entity
}

func makeError(format string, err error) error {

	if err != nil {
		return fmt.Errorf(format, err)
	}

	return fmt.Errorf(format)

}

func makeErrorResponse(w http.ResponseWriter, status int, cause error) {
	errorResponse := types.ErrorResponse{Status: status, Cause: cause.Error()}
	writeResponse(w, &errorResponse, status)
}

func writeResponse(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
	}
}

//Trim trim string for spaces
func Trim(s string) string {
	return strings.Trim(s, " ")
}

func restartAPIGateway(appConfig config.AppConfig, deployment dep.Deployment) {

	go func() {

		err := components.ApigatewaySetup(&appConfig, deployment)
		if err != nil {
			log.Printf("apigateway re-start failed : %v", err)
		}

	}()

}
