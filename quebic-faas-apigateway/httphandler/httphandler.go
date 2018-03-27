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
	"time"

	"github.com/gorilla/mux"
)

//Httphandler handlers
type Httphandler struct {
	Config        config.AppConfig
	Messenger     _messenger.Messenger
	AppStatusList map[string]string
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

	httphandler.healthCheckEndpointHandler(router)
	httphandler.requestTrackerHandler(router)

	//go throught each resource
	for _, resource := range resources {
		httphandler.createHandler(router, resource)
	}

}

func (httphandler *Httphandler) healthCheckEndpointHandler(router *mux.Router) {

	type StatusResponse struct {
		Status interface{} `json:"status"`
	}

	router.HandleFunc("/manage/health", func(w http.ResponseWriter, r *http.Request) {
		if len(httphandler.AppStatusList) == 0 {
			writeResponse(w, StatusResponse{Status: "UP"}, 200)
		} else {
			writeResponse(w, StatusResponse{Status: httphandler.AppStatusList}, 500)
		}
	}).Methods("GET")

}

func (httphandler *Httphandler) requestTrackerHandler(router *mux.Router) {

	router.HandleFunc("/request-tracker/{requestID}", func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		requestID := params["requestID"]

		managerAccessKey := httphandler.Config.Auth.Accesstoken
		requestHeaders := make(map[string]string)
		requestHeaders[common.HeaderAccessKey] = managerAccessKey

		messenger := httphandler.Messenger

		messenger.PublishBlocking(
			common.EventRequestTrackerDataFetch,
			requestID,
			requestHeaders,
			func(message _messenger.BaseEvent, status int) {

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
			func(err string, statuscode int) {
				makeErrorResponse(w, statuscode, fmt.Errorf(err))
			},
			time.Second*5,
		)

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
	writeResponse(w, &successResponse, status)
}

func makeErrorResponse(w http.ResponseWriter, status int, cause error) {
	errorResponse := ErrorResponse{Status: status, Cause: cause.Error()}
	writeResponse(w, &errorResponse, status)
}

func makeErrorStrResponse(w http.ResponseWriter, status int, err string) {
	log.Printf("error response status %v", status)
	errorResponse := ErrorResponse{Status: status, Cause: err}
	writeResponse(w, &errorResponse, status)
}

func writeResponse(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
	}
}
