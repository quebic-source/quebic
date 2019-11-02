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

//Used by manager dashboard
const EnvKey_mgrAPI = "mgr_api"
const EnvKey_ingressConfig_staticIP = "ingressConfig_staticIP"

const EnvKey_appID = "appID"
const EnvKey_deploymentID = "deploymentID"
const EnvKey_version = "version"

const EnvKey_rabbitmq_exchange = "rabbitmq_exchange"
const EnvKey_rabbitmq_host = "rabbitmq_host"
const EnvKey_rabbitmq_port = "rabbitmq_port"
const EnvKey_rabbitmq_management_username = "rabbitmq_management_username"
const EnvKey_rabbitmq_management_password = "rabbitmq_management_password"

const EnvKey_eventConst_eventPrefixUserDefined = "eventConst_eventPrefixUserDefined"

const EnvKey_eventConst_eventLog = "eventConst_eventLog"
const EnvKey_eventConst_eventFunctionAwake = "eventConst_eventFunctionAwake"
const EnvKey_eventConst_eventDataFetch = "eventConst_eventDataFetch"
const EnvKey_eventConst_eventNewVersion = "eventConst_eventNewVersion"
const EnvKey_eventConst_eventShutDownRequest = "eventConst_eventShutDownRequest"

const EnvKey_events = "events"
const EnvKey_artifactLocation = "artifactLocation"
const EnvKey_functionPath = "functionPath"
const EnvKey_functionAge = "functionAge"

//EventBox
const EnvKey_mongo_host = "mongo.host"
const EnvKey_mongo_port = "mongo.port"
const EnvKey_mongo_db = "mongo.db"
const EnvKey_mongo_username = "mongo.username"
const EnvKey_mongo_password = "mongo.password"
const EnvKey_eventbox_uri = "EVENTBOX_URI"
