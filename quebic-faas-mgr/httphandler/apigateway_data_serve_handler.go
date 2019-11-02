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
	"encoding/json"
	"net/http"
	"quebic-faas/auth"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"

	"github.com/gorilla/mux"
)

//ApigatewayDataServe handler
func (httphandler *Httphandler) ApigatewayDataServe(router *mux.Router) {

	authConfig := httphandler.config.Auth
	db := httphandler.db

	router.HandleFunc("/apigateway-data-serve", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		apigatewayData := types.ApigatewayData{}

		//resources
		var resources []types.Resource
		err := dao.GetAll(db, &types.Resource{}, func(k, v []byte) error {

			resource := types.Resource{}
			json.Unmarshal(v, &resource)
			resources = append(resources, resource)
			return nil
		})
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		//manager-components
		allowComponents := [1]string{
			common.ComponentEventBus,
		}
		var components []types.ManagerComponent
		err = dao.GetAll(db, &types.ManagerComponent{}, func(k, v []byte) error {

			component := types.ManagerComponent{}
			json.Unmarshal(v, &component)

			//check this component is allow to serve
			for _, a := range allowComponents {
				if a == component.ID {
					components = append(components, component)
				}
			}

			return nil

		})
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		//remove nil from response
		if components == nil {
			components = make([]types.ManagerComponent, 0)
		}

		if resources == nil {
			resources = make([]types.Resource, 0)
		}

		//assign data
		apigatewayData.Resources = resources
		apigatewayData.ManagerComponents = components

		writeResponse(w, apigatewayData, http.StatusOK)

	}, auth.RoleAny, authConfig)).Methods("GET")

}
