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
	"log"
	"net/http"
	"quebic-faas/common"
	"quebic-faas/messenger"
	"quebic-faas/types"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (httphandler *Httphandler) createHandler(router *mux.Router, resource types.Resource) {

	url := resource.URL
	requestMethod := resource.RequestMethod

	router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		httphandler.usageUp()
		httphandler.eventInvoke(w, r, resource)
		httphandler.usageDown()
	}).Methods(requestMethod)

}

func (httphandler *Httphandler) eventInvoke(w http.ResponseWriter, r *http.Request, resource types.Resource) {

	payload := make(map[string]interface{})
	processRequest(r, &payload)

	for _, mappingTemplate := range resource.RequestMapping {
		httphandler.processForRequestMapping(r, mappingTemplate, payload)
	}

	for _, mappingTemplate := range resource.HeaderMapping {
		httphandler.processForHeaderMapping(r, mappingTemplate, payload)
	}

	//user's headers
	requestHeaders := make(map[string]interface{})
	for _, header := range resource.HeadersToPass {
		requestHeaders[header] = r.Header.Get(header)
	}

	//set request http method
	requestHeaders["requestPath"] = r.RequestURI
	requestHeaders["requestHTTPMethod"] = r.Method
	requestHeaders["async"] = strconv.FormatBool(resource.Async)

	m := httphandler.Messenger

	if !resource.Async {

		requestTimeout := time.Second * time.Duration(resource.RequestTimeout)

		_, err := m.PublishBlocking(
			resource.Event,
			payload,
			requestHeaders,
			func(message messenger.BaseEvent, statuscode int, context messenger.Context) {

				if statuscode == 0 {
					statuscode = resource.SuccessResponseStatus
				}

				makeAPIGatewaySuccessResponse(w, statuscode, message.GetPayloadAsObject(), context.RequestID)

			},
			func(message string, statuscode int, context messenger.Context) {

				makeAPIGatewayErrorResponse(w, statuscode, message, context.RequestID)

			},
			requestTimeout,
		)
		if err != nil {

			log.Printf("internal server error, cause : %s\n", err.Error())

			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

	} else {

		requestID, err := m.Publish(
			resource.Event,
			payload,
			requestHeaders,
			nil,
			nil,
			0,
		)

		if err != nil {

			log.Printf("internal server error, cause : %s\n", err.Error())

			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		makeAPIGatewaySuccessResponse(w, 200, RequestTrackerResponse{RequestID: requestID}, requestID)

	}

}

//check wether requestExpressions value is came as request parms
func (httphandler *Httphandler) processForRequestMapping(r *http.Request, mappingTemplate types.RequestMappingTemplate, payload map[string]interface{}) {

	requestAttribute := mappingTemplate.RequestAttribute
	matchedRequestValue := r.FormValue(requestAttribute)
	if matchedRequestValue != "" {
		payload[mappingTemplate.EventAttribute] = matchedRequestValue
		return
	}

	matchedRequestValue = mux.Vars(r)[requestAttribute]
	if matchedRequestValue != "" {
		payload[mappingTemplate.EventAttribute] = matchedRequestValue
		return
	}

}

//check wether requestExpressions value is came in request header
func (httphandler *Httphandler) processForHeaderMapping(r *http.Request, mappingTemplate types.HeaderMappingTemplate, payload map[string]interface{}) {

	headerAttribute := mappingTemplate.HeaderAttribute

	matchedHeaderValue := r.Header.Get(headerAttribute)

	if matchedHeaderValue != "" {
		payload[mappingTemplate.EventAttribute] = matchedHeaderValue
	}

}

func makeAPIGatewaySuccessResponse(w http.ResponseWriter, status int, message interface{}, requestID string) {
	successResponse := SuccessResponse{Status: status, Message: message}
	writeResponse(w, &successResponse, status, prepareAPIGatewayHeaders(requestID))
}

func makeAPIGatewayErrorResponse(w http.ResponseWriter, status int, err string, requestID string) {
	errorResponse := ErrorResponse{Status: status, Cause: err}
	writeResponse(w, &errorResponse, status, prepareAPIGatewayHeaders(requestID))
}

func prepareAPIGatewayHeaders(requestID string) map[string]string {
	header := make(map[string]string)
	header[common.HeaderRequestID] = requestID
	return header
}
