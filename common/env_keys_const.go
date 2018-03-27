package common

//EnvKeyHostIP env key
const EnvKeyHostIP = "host_ip"

//EnvKeyMgrPORT env key
/*
this is env key which stored in docker image to access manager.
*/
const EnvKeyMgrPORT = "mgr_port"

//EnvKeyAPIGateWayAccessKey env key
/*
this is env key which stored in docker image to access manager restricted endpoints.
apigateway will make http request with the access_key.
*/
const EnvKeyAPIGateWayAccessKey = "access_key"

//EnvKeyFunctionContainerSecret env key
/*
this is env key which stored in docker image when creating new function. when start function container
manager pass jwt encoded meta data which required to start function. this encode data open using this env key
*/
const EnvKeyFunctionContainerSecret = "secret"

//Used by function docker service
const EnvKey_appID = "appID"
const EnvKey_rabbitmq_exchange = "rabbitmq.exchange"
const EnvKey_rabbitmq_host = "rabbitmq.host"
const EnvKey_rabbitmq_port = "rabbitmq.port"
const EnvKey_rabbitmq_management_username = "rabbitmq.management.username"
const EnvKey_rabbitmq_management_password = "rabbitmq.management.password"
const EnvKey_eventConst_eventPrefixUserDefined = "eventConst.eventPrefixUserDefined"
const EnvKey_eventConst_eventLogListener = "eventConst.eventLogListener"

const EnvKey_events = "events"
const EnvKey_artifactLocation = "artifactLocation"
const EnvKey_functionPath = "functionPath"
