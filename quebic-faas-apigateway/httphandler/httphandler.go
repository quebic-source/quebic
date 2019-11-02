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
	"quebic-faas/common"
	_messenger "quebic-faas/messenger"
	"quebic-faas/quebic-faas-apigateway/config"
	"quebic-faas/types"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

//Httphandler handlers
type Httphandler struct {
	Config        config.AppConfig
	Messenger     _messenger.Messenger
	AppStatusList map[string]string
	Usage         int32
}

//SuccessResponse successResponse
type SuccessResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

//ErrorResponse errorResponse
type ErrorResponse struct {
	Status int    `json:"status"`
	Cause  string `json:"cause"`
}

//RequestTrackerResponse requestTrackerResponse
type RequestTrackerResponse struct {
	RequestID string `json:"requestID"`
}

//SetUpHTTPHandlers setUpHTTPHandlers
func SetUpHTTPHandlers(
	config config.AppConfig,
	resources []types.Resource,
	router *mux.Router,
	messenger _messenger.Messenger,
	appStatusList map[string]string) {

	httphandler := &Httphandler{
		Config:        config,
		Messenger:     messenger,
		AppStatusList: appStatusList,
	}

	//check wether this deployment is latest version
	httphandler.checkForShutdown(config.CurrentDeploymentVersion)

	httphandler.listeningForNewVersion()

	httphandler.healthCheckEndpointHandler(router)
	httphandler.requestTrackerHandler(router)

	//go throught each resource
	for _, resource := range resources {
		httphandler.createHandler(router, resource)
	}

}

func (httphandler *Httphandler) healthCheckEndpointHandler(router *mux.Router) {

	type StatusResponse struct {
		Status  interface{} `json:"status"`
		Version string      `json:"version"`
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httphandler.usageUp()
		writeResponse(w, StatusResponse{
			Status:  fmt.Sprintf("%v is running", common.ComponentAPIGateway),
			Version: httphandler.Config.Version,
		}, 200, nil)
		httphandler.usageDown()
	}).Methods("GET")

	router.HandleFunc("/manage/health", func(w http.ResponseWriter, r *http.Request) {
		httphandler.usageUp()
		if len(httphandler.AppStatusList) == 0 {
			writeResponse(w, StatusResponse{Status: "UP"}, 200, nil)
		} else {
			writeResponse(w, StatusResponse{Status: httphandler.AppStatusList}, 500, nil)
		}
		httphandler.usageDown()
	}).Methods("GET")

}

func (httphandler *Httphandler) requestTrackerHandler(router *mux.Router) {

	router.HandleFunc("/request-tracker/{requestID}", func(w http.ResponseWriter, r *http.Request) {

		httphandler.usageUp()

		params := mux.Vars(r)
		requestID := params["requestID"]

		managerAccessKey := httphandler.Config.Auth.Accesstoken
		requestHeaders := make(map[string]interface{})
		requestHeaders[common.HeaderAccessKey] = managerAccessKey

		messenger := httphandler.Messenger

		messenger.PublishBlocking(
			common.EventRequestTrackerDataFetch,
			requestID,
			requestHeaders,
			func(message _messenger.BaseEvent, status int, context _messenger.Context) {

				rtMap := message.GetPayloadAsObject()
				rtJSON, _ := json.Marshal(rtMap)

				rt := types.RequestTracker{}
				err := json.Unmarshal(rtJSON, &rt)
				if err != nil {
					makeErrorResponse(w, http.StatusNotFound, fmt.Errorf("response value not found for this request-id"))
					return
				}

				makeSuccessResponse(w, rt.Response.Status, rt.Response.Message)

			},
			func(err string, statuscode int, context _messenger.Context) {
				makeErrorResponse(w, statuscode, fmt.Errorf(err))
			},
			time.Second*5,
		)

		httphandler.usageDown()

	}).Methods("GET")

}

func processRequest(r *http.Request, requestMap *map[string]interface{}) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return makeError("unable to read request %v", err)
	}

	err = json.Unmarshal(body, requestMap)
	if err != nil {
		return makeError("unable to parse json request to requestMap %v", err)
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

func makeSuccessResponse(w http.ResponseWriter, status int, message interface{}) {
	successResponse := SuccessResponse{Status: status, Message: message}
	writeResponse(w, &successResponse, status, nil)
}

func makeErrorResponse(w http.ResponseWriter, status int, cause error) {
	errorResponse := ErrorResponse{Status: status, Cause: cause.Error()}
	writeResponse(w, &errorResponse, status, nil)
}

func writeResponse(w http.ResponseWriter, v interface{}, status int, headers map[string]string) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if headers != nil {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
	}

}

func (httphandler *Httphandler) listeningForNewVersion() {

	appID := httphandler.Config.DeploymentID + "-" + common.UUIDGen()
	messenger := _messenger.Messenger{
		AppID:          appID,
		EventBusConfig: httphandler.Config.EventBusConfig,
	}
	err := messenger.Init()
	if err != nil {
		log.Fatalf("listening new version messenger setup failed : %v", err)
		return
	}

	consumerID := common.ConsumerNewVersionAPIGateway
	err = messenger.Subscribe(common.EventNewVersionAPIGateway, func(event _messenger.BaseEvent) {

		newVersion := types.NewVersionMessage{}
		event.ParsePayloadAsObject(&newVersion)

		log.Printf("new version received my : %v, new : %v", httphandler.Config.Version, newVersion.Version)

		httphandler.checkForShutdown(newVersion.Version)

	}, consumerID)
	if err != nil {
		log.Fatalf("new-version-listener setup failed. cause : %v", err.Error())
	}

	log.Printf("new-version-listener setup")

}

func (httphandler *Httphandler) checkForShutdown(newVersion string) {

	if httphandler.Config.Version != newVersion {

		for true {
			if !httphandler.isUse() {
				break
			}

			time.Sleep(2 * time.Second)

			//TODO remove this log later
			log.Printf("waiting for apigateway is become free. usage %v", httphandler.getUsage())
		}

		//TODO remove this log later
		log.Printf("apigateway is become free")

		_, err := httphandler.Messenger.Publish(
			common.EventShutDownRequest,
			types.ShutDownRequest{
				DeploymentID: httphandler.Config.DeploymentID,
			},
			nil,
			func(event _messenger.BaseEvent, statuscode int, context _messenger.Context) {
				log.Printf("shutdown request published")
			},
			nil,
			0,
		)
		if err != nil {
			log.Printf("shutdown request publish failed : %v", err)
			return
		}
	}
}

func (httphandler *Httphandler) usageUp() {
	atomic.AddInt32(&httphandler.Usage, 1)
}

func (httphandler *Httphandler) usageDown() {
	atomic.AddInt32(&httphandler.Usage, -1)
}

func (httphandler *Httphandler) isUse() bool {
	return (httphandler.getUsage() > 0)
}

func (httphandler *Httphandler) getUsage() int32 {
	return atomic.LoadInt32(&httphandler.Usage)
}
