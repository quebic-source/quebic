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
	"net/http"
	"quebic-faas/auth"

	"github.com/gorilla/mux"
)

//RequestTrackerHandler request-tracker handler
func (httphandler *Httphandler) RequestTrackerHandler(router *mux.Router) {

	authConfig := httphandler.config.Auth

	router.HandleFunc("/request-trackers", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		qID := r.FormValue("id")
		if qID == "" {

			requestTrackers, err := httphandler.loggerUtil.GetRequestTrackers()
			if err != nil {
				makeErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			writeResponse(w, requestTrackers, http.StatusOK)

			return

		}

		rt, err := httphandler.loggerUtil.GetRequestTrackerByRequestID(qID)
		if err != nil {
			makeErrorResponse(w, http.StatusNotFound, err)
			return
		}

		writeResponse(w, rt, http.StatusOK)

	}, auth.RoleAny, authConfig)).Methods("GET")

	router.HandleFunc("/request-trackers/{requestID}", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		requestID := params["requestID"]
		rt, err := httphandler.loggerUtil.GetRequestTrackerByRequestID(requestID)
		if err != nil {
			makeErrorResponse(w, http.StatusNotFound, err)
			return
		}

		writeResponse(w, rt, http.StatusOK)

	}, auth.RoleAny, authConfig)).Methods("GET")

}
