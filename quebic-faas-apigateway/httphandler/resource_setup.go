package httphandler

import (
	"net/http"
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
		httphandler.eventInvoke(w, r, resource)
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
	requestHeaders := make(map[string]string)
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
			func(message messenger.BaseEvent, statuscode int) {

				if statuscode == 0 {
					statuscode = resource.SuccessResponseStatus
				}

				makeSuccessResponse(w, statuscode, message.GetPayloadAsObject())

			},
			func(message string, statuscode int) {

				makeErrorStrResponse(w, statuscode, message)

			},
			requestTimeout,
		)
		if err != nil {
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
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		makeSuccessResponse(w, 200, RequestTrackerResponse{RequestID: requestID})

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
