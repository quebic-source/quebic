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
package messenger

import (
	"encoding/json"
	"fmt"
	"quebic-faas/common"
	"reflect"
)

const headerEventID = "eventID"
const headerRequestID = "requestID"
const headerCreated = "created"
const headerStatuscode = "statuscode"
const headerError = "error"

//BaseEvent type
type BaseEvent struct {
	eventPayload []byte
	headers      map[string]interface{}
}

func (baseEvent *BaseEvent) init() {
	baseEvent.headers = make(map[string]interface{})
}

//SetEventID set eventID
func (baseEvent *BaseEvent) setEventID(eventID string) {
	baseEvent.headers[headerEventID] = eventID
}

//GetEventID get eventID
func (baseEvent *BaseEvent) GetEventID() string {

	headerval := baseEvent.headers[headerEventID]
	if headerval == nil {
		return ""
	}

	return headerval.(string)
}

//SetRequestID set requestID
func (baseEvent *BaseEvent) setRequestID(requestID string) {
	baseEvent.headers[headerRequestID] = requestID
}

//GetRequestID get requestID
func (baseEvent *BaseEvent) GetRequestID() string {

	headerval := baseEvent.headers[headerRequestID]
	if headerval == nil {
		return ""
	}

	return headerval.(string)

}

//SetCreated set created
func (baseEvent *BaseEvent) setCreated(created string) {
	baseEvent.headers[headerCreated] = created
}

//GetCreated get created
func (baseEvent *BaseEvent) GetCreated() string {

	headerval := baseEvent.headers[headerCreated]
	if headerval == nil {
		return ""
	}

	return headerval.(string)

}

//SetStatuscode set requestHeaders
func (baseEvent *BaseEvent) setStatuscode(statuscode int) {
	baseEvent.headers[headerStatuscode] = common.IntToStr(statuscode)
}

//GetStatuscode get requestHeaders
func (baseEvent *BaseEvent) GetStatuscode() int {

	headerval := baseEvent.headers[headerStatuscode]
	if headerval == nil {
		return common.StrToInt("0")
	}

	return common.StrToInt(headerval.(string))

}

//SetStatuscode set requestHeaders
func (baseEvent *BaseEvent) setError(err string) {
	baseEvent.headers[headerError] = err
}

//GetError get requestHeaders
func (baseEvent *BaseEvent) GetError() string {

	headerval := baseEvent.headers[headerError]
	if headerval == nil {
		return ""
	}

	return headerval.(string)

}

//SetHeaderData set requestHeaders
func (baseEvent *BaseEvent) setHeaderData(key string, value string) {
	baseEvent.headers[key] = value
}

//GetHeaderData get requestHeaders
func (baseEvent *BaseEvent) GetHeaderData(key string) string {

	headerval := baseEvent.headers[key]
	if headerval == nil {
		return ""
	}

	return headerval.(string)

}

//GetPayload get payload []byte
func (baseEvent *BaseEvent) GetPayload() []byte {
	return baseEvent.eventPayload
}

//GetPayloadAsString get payload string
func (baseEvent *BaseEvent) GetPayloadAsString() string {
	return string(baseEvent.eventPayload)
}

//GetPayloadAsObject get payloadObject as a perticular object
func (baseEvent *BaseEvent) GetPayloadAsObject() interface{} {

	payload := make(map[string]interface{})
	baseEvent.ParsePayloadAsObject(&payload)

	//check responce is json formated
	if len(payload) > 0 {
		return payload
	}

	//otherwise response is a primitive data type.
	return common.StrParseToPrimitive(baseEvent.GetPayloadAsString())

}

//ParsePayloadAsObject parse to object you provided
func (baseEvent *BaseEvent) ParsePayloadAsObject(payloadObject interface{}) error {
	err := json.Unmarshal(baseEvent.eventPayload, payloadObject)
	if err != nil {
		return fmt.Errorf("failed when convet json eventPayload to object %v", err)
	}
	return nil
}

//SetPayloadObject set payloadObject
func (baseEvent *BaseEvent) setPayloadObject(payloadObject interface{}) {

	if payloadObject == nil {
		return
	}

	//check only for string
	if reflect.TypeOf(payloadObject).String() == "string" {
		baseEvent.eventPayload = []byte(payloadObject.(string))
		return
	}

	jsonPayload, _ := json.Marshal(payloadObject)
	baseEvent.eventPayload = jsonPayload

}
