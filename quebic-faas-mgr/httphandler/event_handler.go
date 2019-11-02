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
	"net/http"
	"quebic-faas/auth"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/gorilla/mux"
)

//EventHandler handler
func (httphandler *Httphandler) EventHandler(router *mux.Router) {

	db := httphandler.db
	authConfig := httphandler.config.Auth

	router.HandleFunc("/events", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {
		getAllEvents(w, r, db, &types.Event{})
	}, auth.RoleAny, authConfig)).Methods("GET")

	router.HandleFunc("/events", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		event := &types.Event{}
		err := processRequest(r, event)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		trimStringFieldsEvent(event)

		errors := validationEvent(event)
		if errors != nil {
			status := http.StatusBadRequest
			writeResponse(w, types.ErrorResponse{Cause: "validation Failed", Message: errors, Status: status}, status)
			return
		}

		processEvent(event)

		add(w, r, db, event)

	}, auth.RoleAny, authConfig)).Methods("POST")

}

func createEventFromEventID(db *bolt.DB, eventID string) (*types.Event, []string) {

	event := extractEventFromEventID(eventID)

	trimStringFieldsEvent(event)

	errors := validationEvent(event)

	if errors != nil || len(errors) > 0 {
		return nil, errors
	}

	processEvent(event)

	dao.Add(db, event)

	return event, nil

}

func getAllEvents(w http.ResponseWriter, r *http.Request, db *bolt.DB, entity types.Entity) {

	var events []types.Event
	err := dao.GetAll(db, entity, func(k, v []byte) error {

		event := types.Event{}
		json.Unmarshal(v, &event)
		events = append(events, event)
		return nil
	})

	if err != nil {
		makeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if events == nil {
		var emptyStr [0]string
		writeResponse(w, emptyStr, http.StatusOK)
	} else {
		writeResponse(w, events, http.StatusOK)
	}

}

func trimStringFieldsEvent(event *types.Event) {
	event.Name = Trim(event.Name)
}

func validationEvent(event *types.Event) []string {

	var errors []string

	//TODO needs to re think about this validation
	/*
		if event.Group == "" {
			errors = append(errors, "event group field should not be empty")
		}
	*/

	if strings.Contains(event.Group, " ") {
		errors = append(errors, "event group field not allow to contain spaces")
	}

	if strings.Contains(event.Group, ".") || strings.Contains(event.Group, "*") || strings.Contains(event.Group, "#") {
		errors = append(errors, "event group field not allow to contain */#/.")
	}

	if event.Name == "" {
		errors = append(errors, "event name field should not be empty")
	}

	if strings.Contains(event.Name, " ") {
		errors = append(errors, "event name field not allow to contain spaces")
	}

	return errors

}

func processEvent(event *types.Event) {
	if event.Group == "" {
		event.ID = common.EventPrefixUserDefined + common.EventJOIN + event.Name
	} else {
		event.ID = common.EventPrefixUserDefined + common.EventJOIN + event.Group + common.EventJOIN + event.Name
	}
}

/*
	name = event.Group + event.Name
*/
func checkEventGroupAndNameIsValide(db *bolt.DB, name string) bool {

	id := common.EventPrefixUserDefined + common.EventJOIN + name

	err := dao.GetByID(db, &types.Event{ID: id}, func(savedObj []byte) error {

		if savedObj == nil {
			return makeError("resource not found", nil)
		}

		return nil
	})

	if err != nil {
		return false
	}

	return true

}

/*
	name = event.Group + event.Name
	concat EventPrefix to eventGroup + name
*/
func prepareEventGroupAndNameToID(name string) string {
	return common.EventPrefixUserDefined + common.EventJOIN + name
}

func extractEventFromEventID(eventID string) *types.Event {

	// <group>.<EventName>
	eArray := strings.Split(eventID, ".")

	event := &types.Event{}

	if len(eArray) >= 2 {
		event.Group = eArray[0]
		event.Name = strings.Join(eArray[1:len(eArray)], ".")
	} else {
		event.Name = eventID
	}

	return event

}
