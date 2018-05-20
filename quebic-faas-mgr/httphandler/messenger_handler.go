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
	"io/ioutil"
	"log"
	"net/http"
	"quebic-faas/messenger"
	"time"

	"github.com/gorilla/mux"
)

//ConsumerDTO consumer DTO
type ConsumerDTO struct {
	EventID    string
	ConsumerID string
}

//Log log
type Log struct {
	RequestID string `json:"requestID"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Source    string `json:"source"` //executed by function-container-id or app-id
	Time      string `json:"time"`   //executed time
}

//MessengerHandler handler
func (httphandler *Httphandler) MessengerHandler(router *mux.Router) {

	messengerContext := httphandler.messenger

	router.HandleFunc("/messenger/consume", func(w http.ResponseWriter, r *http.Request) {

		consumerDTO := ConsumerDTO{}
		err := processRequestBaseEvent(r, &consumerDTO)

		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		//c1
		err = messengerContext.Subscribe(consumerDTO.EventID, func(event messenger.BaseEvent) {

			log.Printf("payload %v", event.GetPayloadAsString())

			messengerContext.ReplyError(event, "Fuck off", 403)

		}, consumerDTO.ConsumerID)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		writeResponse(w, consumerDTO, 200)

	}).Methods("POST")

	router.HandleFunc("/messenger/publish", func(w http.ResponseWriter, r *http.Request) {

		consumerDTO := ConsumerDTO{}
		err := processRequestBaseEvent(r, &consumerDTO)

		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		type Request struct {
			Message string
		}

		logMessage := Log{
			RequestID: consumerDTO.EventID,
			Message:   "test log 1",
			Type:      "INFO",
			Source:    "message_handler",
			Time:      time.Now().String(),
		}

		_, err = messengerContext.Publish(
			consumerDTO.EventID,
			logMessage,
			nil,
			nil,
			nil,
			0)

		if err != nil {
			makeErrorResponse(w, http.StatusRequestTimeout, err)
			return
		}

		writeResponse(w, consumerDTO, 200)

	}).Methods("POST")

}

func processRequestBaseEvent(r *http.Request, consumerDTO *ConsumerDTO) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return makeError("unable to read request %v", err)
	}

	err = json.Unmarshal(body, consumerDTO)
	if err != nil {
		return makeError("unable to parse json request to consumerDTO %v", err)
	}

	return nil

}
