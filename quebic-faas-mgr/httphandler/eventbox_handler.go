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

/*
	this is used to server data which are requested by apigateway
*/

import (
	"net/http"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/components"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/types"

	"github.com/gorilla/mux"
)

//EventBoxHandler handler
func (httphandler *Httphandler) EventBoxHandler(router *mux.Router) {

	appConfig := httphandler.config
	deployment := httphandler.deployment

	router.HandleFunc("/eventbox/start", func(w http.ResponseWriter, r *http.Request) {

		_, _, err := components.EventBoxSetup(appConfig, deployment)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		component, err := getEventBox(deployment)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		writeResponse(w, component, http.StatusOK)

	}).Methods("POST")

	router.HandleFunc("/eventbox/info", func(w http.ResponseWriter, r *http.Request) {

		component, err := getEventBox(deployment)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		writeResponse(w, component, http.StatusOK)

	}).Methods("GET")

}

func getEventBox(deployment dep.Deployment) (types.ManagerComponent, error) {

	component := types.ManagerComponent{ID: common.ComponentEventBox}

	eventBoxDetails, err := deployment.ListByName(component.ID)
	if err != nil {
		return component, err
	}

	component.Deployment.Host = eventBoxDetails.Host
	component.Deployment.Port = int(eventBoxDetails.PortConfigs[0].Port)

	return component, nil
}
