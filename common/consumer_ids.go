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

package common

const consumerPrefix = "quebic-faas-event-consumer"

//ConsumerJOIN consumer id join
const ConsumerJOIN = "-"

//ConsumerApigatewayDataFetch apigateway-data-fetch
var ConsumerApigatewayDataFetch = prepareConsumerID("apigateway-data-fetch")

//ConsumerFunctionDataFetch function-data-fetch
var ConsumerFunctionDataFetch = prepareConsumerID("function-data-fetch")

//ConsumerApigateway apigateway
var ConsumerApigateway = prepareConsumerID("apigateway")

//ConsumerRequestTracker request-tracker
var ConsumerRequestTracker = prepareConsumerID("request-tracker")

//ConsumerRequestTrackerDataFetch request-tracker-data-fetch
var ConsumerRequestTrackerDataFetch = prepareConsumerID("request-tracker-data-fetch")

//ConsumerFunctionRequestPrefix function-request
var ConsumerFunctionRequestPrefix = prepareConsumerID("function-request")

//ConsumerFunctionAwakePrefix function-awake
var ConsumerFunctionAwakePrefix = prepareConsumerID("function-awake")

//ConsumerNewVersionAPIGateway new-version-apigateway
var ConsumerNewVersionAPIGateway = prepareConsumerID("new-version-apigateway")

//ConsumerShutDownRequest shutdown-request
var ConsumerShutDownRequest = prepareConsumerID("shutdown-request")

func prepareConsumerID(id string) string {
	return consumerPrefix + ConsumerJOIN + id + ConsumerJOIN + UUIDGen()
}
