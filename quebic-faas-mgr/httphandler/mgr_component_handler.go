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
	"quebic-faas/auth"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/types"

	"github.com/gorilla/mux"
)

//MgrComponentHandler handler
func (httphandler *Httphandler) MgrComponentHandler(router *mux.Router) {

	appConfig := httphandler.config
	authConfig := appConfig.Auth

	router.HandleFunc("/mgr-components", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		//components
		var components []types.ManagerComponent

		apiGateway, err := getAPIGateway(appConfig)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		components = append(components, apiGateway)

		if components == nil {
			var emptyStr [0]string
			writeResponse(w, emptyStr, http.StatusOK)
		} else {
			writeResponse(w, components, http.StatusOK)
		}

	}, auth.RoleAny, authConfig)).Methods("GET")

}

func getAPIGateway(appConfig config.AppConfig) (types.ManagerComponent, error) {

	component := types.ManagerComponent{ID: common.ComponentAPIGateway}

	component.Deployment.Host = appConfig.APIGatewayConfig.ServerConfig.Host
	component.Deployment.Port = appConfig.APIGatewayConfig.ServerConfig.Port

	return component, nil
}
