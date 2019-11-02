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

const componentPrefix = "quebic-faas-"

//ComponentDockerRegistry docker_registry
const ComponentDockerRegistry = componentPrefix + "docker-registry"

//ComponentEventBus eventbus
const ComponentEventBus = componentPrefix + "eventbus"

//ComponentAPIGateway apigateway
const ComponentAPIGateway = componentPrefix + "apigateway"

//ComponentMgrAPI mgr-api
const ComponentMgrAPI = componentPrefix + "mgr"

//ComponentMgrDashboard mgr-dashboard
const ComponentMgrDashboard = componentPrefix + "mgr-dashboard"

//ComponentEventBox eventbox
const ComponentEventBox = componentPrefix + "eventbox"

//ComponentEventBoxDB eventbox-db
const ComponentEventBoxDB = componentPrefix + "eventbox-db"

//ComponentAPIGatewayDefaultReplicas apigateway default replicas
const ComponentAPIGatewayDefaultReplicas = 1

//ComponentAPIGatewayVersionDefaultStart default start value of a apigateway
const ComponentAPIGatewayVersionDefaultStart = 1
