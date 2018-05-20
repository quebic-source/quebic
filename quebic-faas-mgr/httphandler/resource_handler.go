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
	"net/http"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"
	"sort"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/gorilla/mux"
)

//ResourceHandler handler
func (httphandler *Httphandler) ResourceHandler(router *mux.Router) {

	db := httphandler.db
	appConfig := httphandler.config
	deployment := httphandler.deployment

	router.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {

		qID := r.FormValue("id")
		if qID == "" {
			getAllResources(w, r, db, &types.Resource{})
			return
		}

		resource := &types.Resource{}
		resource.ID = qID
		getByID(w, r, db, resource)

	}).Methods("GET")

	router.HandleFunc("/routes/{name}", func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		routeName := params["name"]

		found := false
		var route *types.Resource
		_ = dao.GetAll(db, &types.Resource{}, func(k, v []byte) error {

			route = &types.Resource{}
			json.Unmarshal(v, route)

			if routeName == route.Name {
				found = true
				return fmt.Errorf("break loop")
			}

			return nil
		})

		if !found {
			makeErrorResponse(w, http.StatusNotFound, fmt.Errorf("resource not found"))
			return
		}

		writeResponse(w, route, 200)

	}).Methods("GET")

	router.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {

		resource := &types.Resource{}
		err := processRequest(r, resource)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		trimStringFieldsResource(resource)

		errors := validationRoute(db, resource, false, true)
		if errors != nil {
			status := http.StatusBadRequest
			writeResponse(w, types.ErrorResponse{Cause: "validation-Failed", Message: errors, Status: status}, status)
			return
		}

		err = preProcessResource(db, resource)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		add(w, r, db, resource)

		restartAPIGateway(appConfig, deployment)

	}).Methods("POST")

	router.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {

		resource := &types.Resource{}
		err := processRequest(r, resource)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		trimStringFieldsResource(resource)

		errors := validationRoute(db, resource, false, false)
		if errors != nil {
			status := http.StatusBadRequest
			writeResponse(w, types.ErrorResponse{Cause: "validation-Failed", Message: errors, Status: status}, status)
			return
		}

		err = preProcessResource(db, resource)
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		update(w, r, db, resource)

		restartAPIGateway(appConfig, deployment)

	}).Methods("PUT")

}

func getAllResources(w http.ResponseWriter, r *http.Request, db *bolt.DB, entity types.Entity) {

	var resources []types.Resource
	err := dao.GetAll(db, entity, func(k, v []byte) error {

		resource := types.Resource{}
		json.Unmarshal(v, &resource)
		resources = append(resources, resource)
		return nil
	})

	if err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if resources == nil {
		var emptyStr [0]string
		writeResponse(w, emptyStr, http.StatusOK)
	} else {

		doSort(r, resources)

		writeResponse(w, resources, http.StatusOK)
	}

}

func trimStringFieldsResource(resource *types.Resource) {
	resource.Name = Trim(resource.Name)
	resource.URL = Trim(resource.URL)
	resource.RequestMethod = Trim(resource.RequestMethod)
}

func validationRoute(db *bolt.DB, resource *types.Resource, isFounctionInvokeRoute bool, isCreate bool) []string {

	var errors []string

	if resource.Name == "" {
		errors = append(errors, "name field should not be empty")
	}

	if strings.Contains(resource.Name, " ") {
		errors = append(errors, "name field not allow to contain spaces")
	}

	if resource.URL == "" {
		errors = append(errors, "url field should not be empty")
	}

	if strings.Contains(resource.URL, " ") {
		errors = append(errors, "url field not allow to contain spaces")
	}

	if resource.RequestMethod == "" {
		errors = append(errors, "requestMethod field should not be empty")
	}

	if resource.RequestMethod != "" {

		if !(resource.RequestMethod == common.ResourceRequestMethodGET ||
			resource.RequestMethod == common.ResourceRequestMethodPOST ||
			resource.RequestMethod == common.ResourceRequestMethodPUT ||
			resource.RequestMethod == common.ResourceRequestMethodDELETE) {
			errors = append(errors, "requestMethod is not match to any of thses methods ( GET , POST , PUT , DELETE )")
		}

	}

	if resource.Event == "" {
		errors = append(errors, "event field should not be empty")
	}

	if resource.Event != "" {

		//event is saved only for user defined route. function invokers no need to save
		if !isFounctionInvokeRoute {

			event, eventErrors := createEventFromEventID(db, resource.Event)

			if eventErrors != nil {
				errors = append(errors, eventErrors...)
			} else {
				resource.Event = event.GetID()
			}

		}

	}

	prepareResourceID(resource)

	if isCreate {

		if checkRouteISAlreadyExists(db, resource) {
			errors = append(errors, "route is already exists")
		}

	} else {

		if !checkRouteISAlreadyExists(db, resource) {
			errors = append(errors, "route is not found")
		}

	}

	return errors

}

func checkRouteISAlreadyExists(db *bolt.DB, resource *types.Resource) bool {

	found := false
	_ = dao.GetByID(db, resource, func(savedObj []byte) error {

		if savedObj != nil {
			found = true
		}

		return nil
	})

	return found
}

func preProcessResource(db *bolt.DB, resource *types.Resource) error {

	if !strings.HasPrefix(resource.URL, "/") {
		resource.URL = "/" + resource.URL
	}

	if strings.HasSuffix(resource.URL, "/") {
		slashLastIndex := strings.LastIndex(resource.URL, "/")
		resource.URL = resource.URL[0:slashLastIndex]
	}

	return nil
}

func prepareResourceID(resource *types.Resource) {
	resource.ID = resource.URL + common.ResourceJOIN + resource.RequestMethod
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(r1, r2 *types.Resource) bool

type routesSorter struct {
	routes []types.Resource
	by     func(r1, r2 *types.Resource) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *routesSorter) Len() int {
	return len(s.routes)
}

// Swap is part of sort.Interface.
func (s *routesSorter) Swap(i, j int) {
	s.routes[i], s.routes[j] = s.routes[j], s.routes[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *routesSorter) Less(i, j int) bool {
	return s.by(&s.routes[i], &s.routes[j])
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(routes []types.Resource) {
	rs := &routesSorter{
		routes: routes,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(rs)
}

func doSort(r *http.Request, resources []types.Resource) {

	const fieldID = "id"
	const fieldModifiedDate = "modifiedDate"
	const fieldURL = "url"

	sortField := r.FormValue("sortBy")

	var sort func(r1, r2 *types.Resource) bool

	if sortField == fieldModifiedDate {
		sort = func(r1, r2 *types.Resource) bool {
			return r1.ModifiedAt < r2.ModifiedAt
		}
	} else if sortField == fieldURL {
		sort = func(r1, r2 *types.Resource) bool {
			return r1.URL < r2.URL
		}
	} else {
		sort = func(r1, r2 *types.Resource) bool {
			return r1.ID < r2.ID
		}
	}

	By(sort).Sort(resources)

}
