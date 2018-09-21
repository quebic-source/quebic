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

const eventPrefix = "quebic-faas-event"

//EventJOIN event id join
const EventJOIN = "."

//EventPrefixUserDefined event id prefix
const EventPrefixUserDefined = eventPrefix + EventJOIN + "defined"

//EventPrefixFunction function prefix
const EventPrefixFunction = EventPrefixUserDefined + EventJOIN + "function"

//EventPrefixExecutionACK execution ack prefix
const EventPrefixExecutionACK = eventPrefix + EventJOIN + "execution-ack"

//EventPrefixInternal execution ack prefix
const EventPrefixInternal = eventPrefix + EventJOIN + "internal"

//EventPrefixFunctionAwake quebic-faas-event.internal.function-awake.<function-name>
const EventPrefixFunctionAwake = eventPrefix + EventJOIN + "function-awake"

//ConsumerMgr consumer
const ConsumerMgr = "quebic-faas-mgr"

//ConsumerApigateway consumer
const ConsumerApigateway = "quebic-faas-apigateway"

//ConsumerRequestTracker consumer
const ConsumerRequestTracker = "quebic-faas-request-tracker"

//ConsumerRequestTrackerDataFetch consumer data-fetch
const ConsumerRequestTrackerDataFetch = "quebic-faas-request-tracker-data-fetch"

//ConsumerFunctionRequestPrefix quebic-faas-function-awake-
const ConsumerFunctionRequestPrefix = "quebic-faas-function-request-"

//ConsumerFunctionAwakePrefix quebic-faas-function-awake-
const ConsumerFunctionAwakePrefix = "quebic-faas-function-awake-"

//EventApigatewayDataFetch used to share data between manage ans apigateway
const EventApigatewayDataFetch = EventPrefixInternal + EventJOIN + "apigateway-data-fetch"

//EventRequestTracker used to listing apigatewat logs
const EventRequestTracker = EventPrefixInternal + EventJOIN + "request-tracker"

//EventRequestTrackerDataFetch used to listing apigateways data fetch request
const EventRequestTrackerDataFetch = EventPrefixInternal + EventJOIN + "request-tracker-data-fetch"
