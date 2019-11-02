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

//EventPrefixFunctionAwakePrefix quebic-faas-event.internal.function-awake.<function-name>
const EventPrefixFunctionAwakePrefix = EventPrefixInternal + EventJOIN + "function-awake" + EventJOIN

//EventApigatewayDataFetch used to share data between manage ans apigateway
const EventApigatewayDataFetch = EventPrefixInternal + EventJOIN + "apigateway-data-fetch"

//EventFunctionDataFetch used to share data between manage ans function
const EventFunctionDataFetch = EventPrefixInternal + EventJOIN + "function-data-fetch"

//EventRequestTracker used to listing apigatewat logs
const EventRequestTracker = EventPrefixInternal + EventJOIN + "request-tracker"

//EventRequestTrackerDataFetch used to listing apigateways data fetch request
const EventRequestTrackerDataFetch = EventPrefixInternal + EventJOIN + "request-tracker-data-fetch"

//EventNewVersionPrefix brodcast event for new versions
const EventNewVersionPrefix = EventPrefixInternal + EventJOIN + "new-version" + EventJOIN

//EventNewVersionAPIGateway api-gateway new version
const EventNewVersionAPIGateway = EventNewVersionPrefix + ComponentAPIGateway

//EventNewVersionFunctionPrefix function new version
const EventNewVersionFunctionPrefix = EventNewVersionPrefix + "function" + EventJOIN

//EventShutDownRequest event send from component that request to delete deployment
const EventShutDownRequest = EventPrefixInternal + EventJOIN + "shutdown-request"
