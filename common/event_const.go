package common

const eventPrefix = "quebic-faas-event"

//EventJOIN event id join
const EventJOIN = "."

//EventPrefixUserDefined event id prefix
const EventPrefixUserDefined = eventPrefix + EventJOIN + "user-defined"

//EventPrefixFunction function prefix
const EventPrefixFunction = eventPrefix + EventJOIN + "function"

//EventPrefixExecutionACK execution ack prefix
const EventPrefixExecutionACK = eventPrefix + EventJOIN + "execution-ack"

//EventPrefixInternal execution ack prefix
const EventPrefixInternal = eventPrefix + EventJOIN + "internal"

//ConsumerMgr consumer
const ConsumerMgr = "quebic-faas-mgr"

//ConsumerApigateway consumer
const ConsumerApigateway = "quebic-faas-apigateway"

//ConsumerRequestTracker consumer
const ConsumerRequestTracker = "quebic-faas-request-tracker"

//ConsumerRequestTrackerDataFetch consumer data-fetch
const ConsumerRequestTrackerDataFetch = "quebic-faas-request-tracker-data-fetch"

//EventApigatewayDataFetch used to share data between manage ans apigateway
const EventApigatewayDataFetch = EventPrefixInternal + EventJOIN + "apigateway-data-fetch"

//EventRequestTracker used to listing apigatewat logs
const EventRequestTracker = EventPrefixInternal + EventJOIN + "request-tracker"

//EventRequestTrackerDataFetch used to listing apigateways data fetch request
const EventRequestTrackerDataFetch = EventPrefixInternal + EventJOIN + "request-tracker-data-fetch"
