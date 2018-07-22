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

//DockerNetworkID network id in service
const DockerNetworkID = "quebic-faas-network"

//DockerServiceEventBus network id in service
const DockerServiceEventBus = "quebic-faas-eventbus"

//DockerServiceApigateway network id in service
const DockerServiceApigateway = "quebic-faas-apigateway"

const ApigatewayImage = "quebicdocker/quebic-faas-apigateway:1.0.0"

const EventbusImage = "rabbitmq:3.7-management-alpine"

const MgrDashboardImage = "quebicdocker/quebic-faas-mgr-dashboard:0.1.0"

const EventBoxImage = "quebicdocker/quebic-eventbox:0.1.0"

const EventBoxDBImage = "mongo"
