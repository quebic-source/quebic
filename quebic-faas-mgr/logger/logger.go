/*
Copyright 2018 Tharanga Nilupul Thennakoon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"quebic-faas/common"
	_messenger "quebic-faas/messenger"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"

	bolt "github.com/coreos/bbolt"
)

//Logger apigateway-logger
type Logger struct {
	db        *bolt.DB
	messenger _messenger.Messenger
}

//Init init
func (logger *Logger) Init(db *bolt.DB, messenger _messenger.Messenger) {

	logger.db = db
	logger.messenger = messenger

}

//Listen start listen for logs for give requestID
func (logger *Logger) Listen() error {

	db := logger.db
	messenger := logger.messenger

	if db == nil {
		return fmt.Errorf("log-listener setup failed error : logger is not initialized")
	}

	//setup log listener
	err := messenger.Subscribe(common.EventRequestTracker, func(event _messenger.BaseEvent) {
		err := logger.saveLog(event)
		if err != nil {
			log.Printf("log-listener log event save failed error : %v", err)
		}
	}, common.ConsumerRequestTracker)
	if err != nil {
		return err
	}

	//setup data fetch listener
	err = messenger.Subscribe(common.EventRequestTrackerDataFetch, func(event _messenger.BaseEvent) {

		requestID := event.GetPayloadAsString()
		if requestID != "" {

			requestTracker, err := logger.GetRequestTrackerByRequestID(requestID)

			if err != nil {
				messenger.ReplyError(event, err.Error(), 404)
				return
			}

			messenger.ReplySuccess(event, requestTracker, 200)

		} else {

			messenger.ReplyError(event, "request-id not found", 400)

		}

	}, common.ConsumerRequestTrackerDataFetch)
	if err != nil {
		return err
	}

	log.Printf("log-listener  setup successfully")

	return nil

}

//GetRequestTrackers get all logs for requestID
func (logger *Logger) GetRequestTrackers() ([]types.RequestTracker, error) {

	var requestTrackers []types.RequestTracker
	err := dao.GetAll(logger.db, &types.RequestTracker{}, func(k, v []byte) error {

		requestTracker := types.RequestTracker{}
		json.Unmarshal(v, &requestTracker)
		requestTrackers = append(requestTrackers, requestTracker)
		return nil

	})
	if err != nil {
		return nil, err
	}

	if requestTrackers == nil {
		requestTrackers = make([]types.RequestTracker, 0)
	}

	return requestTrackers, nil

}

//GetRequestTrackerByRequestID get by requestID
func (logger *Logger) GetRequestTrackerByRequestID(requestID string) (*types.RequestTracker, error) {

	requestTracker := &types.RequestTracker{RequestID: requestID}
	err := logger.getRequestTrackerByID(requestTracker)
	if err != nil {
		return nil, err
	}

	return requestTracker, nil

}

func (logger *Logger) saveLog(event _messenger.BaseEvent) error {

	requestTrackerMessage := types.RequestTrackerMessage{}
	err := event.ParsePayloadAsObject(&requestTrackerMessage)
	if err != nil {
		return err
	}

	db := logger.db

	requestTracker := &types.RequestTracker{RequestID: requestTrackerMessage.RequestID}
	log := requestTrackerMessage.Log

	err = logger.getRequestTrackerByID(requestTracker)
	if err != nil {
		//this is new request tracker
		requestTracker.Source = log.Source
		requestTracker.CreatedAt = log.Time
	}

	requestTracker.Logs = append(requestTracker.Logs, log)

	//check for completed
	if requestTrackerMessage.Completed {
		requestTracker.Response = requestTrackerMessage.Response
		requestTracker.CompletedAt = log.Time
	}

	dao.Save(db, requestTracker)

	return nil
}

func (logger *Logger) getRequestTrackerByID(requestTracker *types.RequestTracker) error {

	err := dao.GetByID(logger.db, requestTracker, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("request-tracker not found")
		}

		json.Unmarshal(savedObj, requestTracker)

		return nil
	})
	if err != nil {
		return err
	}

	return nil

}

func makeError(format string, err error) error {

	if err != nil {
		return fmt.Errorf(format, err)
	}

	return fmt.Errorf(format)

}
