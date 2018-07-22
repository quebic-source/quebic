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

package messenger

import (
	"fmt"
	"log"
	"quebic-faas/common"
	"quebic-faas/types"
	"time"

	"github.com/streadway/amqp"
)

const defaultRequestTimeout = 40 * time.Second

//Context publish handler current context
type Context struct {
	RequestID string //self reuest-id
}

//Publish publish
func (messenger *Messenger) Publish(
	eventID string,
	payload interface{},
	requestHeaders map[string]string,
	successHandler func(message BaseEvent, statuscode int, context Context),
	errHandler func(err string, statuscode int, context Context),
	requestTimeout time.Duration) (string, error) {

	return messenger.publish(
		eventID,
		payload,
		requestHeaders,
		defaultStatuscode,
		emptyError,
		successHandler,
		errHandler,
		requestTimeout,
		false)

}

//PublishBlocking publish and wait caller thread until responce came
func (messenger *Messenger) PublishBlocking(
	eventID string,
	payload interface{},
	requestHeaders map[string]string,
	successHandler func(message BaseEvent, statuscode int, context Context),
	errHandler func(err string, statuscode int, context Context),
	requestTimeout time.Duration) (string, error) {

	return messenger.publish(
		eventID,
		payload,
		requestHeaders,
		defaultStatuscode,
		emptyError,
		successHandler,
		errHandler,
		requestTimeout,
		true)

}

//publish internal publish
func (messenger *Messenger) publish(
	eventID string,
	payload interface{},
	requestHeaders map[string]string,
	statuscode int,
	errorMessage string,
	successHandler func(message BaseEvent, statuscode int, context Context),
	errHandler func(err string, statuscode int, context Context),
	requestTimeout time.Duration,
	blocking bool) (string, error) {

	if eventID == "" {
		return "", fmt.Errorf("eventID should not be empty")
	}

	//create BaseEvent
	baseEvent := BaseEvent{}
	baseEvent.init()
	baseEvent.setEventID(eventID)
	baseEvent.setPayloadObject(payload)
	baseEvent.setError(errorMessage)

	baseEvent.setStatuscode(statuscode)

	for k, v := range requestHeaders {
		baseEvent.setHeaderData(k, v)
	}

	//eventID become the routingKey
	routingKey := baseEvent.GetEventID()

	//requestID
	err := createRequestID(&baseEvent)
	if err != nil {
		return "", err
	}
	requestID := baseEvent.GetRequestID()

	//request-tracker setup
	if eventID != common.EventRequestTracker {
		messenger.setUpRequestTracker(requestID)
	}

	//created date
	baseEvent.setCreated(time.Now().String())

	//responseHandler setup
	waitForResponse := make(chan bool)
	if successHandler != nil || errHandler != nil {

		err := messenger.Subscribe(requestID, func(be BaseEvent) {

			//reply
			err := be.GetError()
			statuscode := be.GetStatuscode()
			if err == "" {

				successHandler(be, statuscode, Context{RequestID: requestID})
				waitForResponse <- true

			} else {

				if errHandler != nil {
					errHandler(err, statuscode, Context{RequestID: requestID})
				}

				waitForResponse <- false

			}

			//relese
			messenger.ReleseQueue(requestID)

		}, common.ConsumerRequestTracker)
		if err != nil {
			log.Printf("publish responseHandler setup failed %v", err)
			waitForResponse <- false
		}

		//default timeout
		if requestTimeout <= 0 {
			requestTimeout = defaultRequestTimeout
		}

	}

	//publish
	err = messenger.channel.Publish(
		Exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        baseEvent.eventPayload,
			Headers:     baseEvent.headers,
		},
	)
	if err != nil {
		messenger.ReleseQueue(requestID)
		return "", fmt.Errorf("failed to publish message, error : %v", err)
	}

	if successHandler != nil || errHandler != nil {

		if blocking {
			wait(messenger, requestID, waitForResponse, requestTimeout, errHandler)
		} else {
			go func() {
				wait(messenger, requestID, waitForResponse, requestTimeout, errHandler)
			}()
		}

	}

	return requestID, nil

}

func wait(
	messenger *Messenger,
	requestID string,
	waitForResponse chan bool,
	requestTimeout time.Duration,
	errHandler func(err string, statuscode int, context Context),
) {
	//wait for response
	select {
	case <-waitForResponse:
		break
	case <-time.After(requestTimeout):
		messenger.ReleseQueue(requestID)
		if errHandler != nil {
			errHandler(fmt.Sprintf("request timeout for %s", requestID), 500, Context{RequestID: requestID})
		}
	}
}

func (messenger *Messenger) setUpRequestTracker(requestID string) {

	requestTrackerMessage := types.RequestTrackerMessage{
		RequestID: requestID,
		Log: types.Log{
			Message: "request-tracker created",
			Source:  messenger.AppID,
			Time:    time.Now().Format(common.DefaultTimeLayout),
			Type:    "INFO",
		},
	}

	messenger.publish(
		common.EventRequestTracker,
		requestTrackerMessage,
		nil,
		0,
		"",
		nil,
		nil,
		0,
		false,
	)

}
