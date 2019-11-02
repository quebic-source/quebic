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

package config

import (
	"quebic-faas/common"
	"quebic-faas/config"
)

//AppConfig appConfig
type AppConfig struct {
	AppID                    string                `json:"appID"`
	DeploymentID             string                `json:"deploymentID"`
	Version                  string                `json:"version"`
	CurrentDeploymentVersion string                `json:"currentDeploymentVersion"`
	Auth                     AuthConfig            `json:"auth"`
	ServerConfig             config.ServerConfig   `json:"serverConfig"`
	EventBusConfig           config.EventBusConfig `json:"eventBusConfig"`
}

//SetDefault set default
func (appConfig *AppConfig) SetDefault() {

	appConfig.AppID = common.ComponentAPIGateway

	appConfig.ServerConfig = config.ServerConfig{
		Host: common.HostMachineIP,
		Port: common.ApigatewayServerPort,
	}

	//No need port configurations. All services are run under same network
	appConfig.EventBusConfig = config.EventBusConfig{
		AMQPHost:           common.DockerServiceEventBus,
		ManagementHost:     common.DockerServiceEventBus,
		ManagementUserName: common.RabbitmqManagementUserName,
		ManagementPassword: common.RabbitmqManagementPassword,
	}

	//used in development mode
	/*appConfig.EventBusConfig = config.EventBusConfig{
		AMQPHost:           common.HostMachineIP,
		AMQPPort:           common.RabbitmqAMQPPort,
		ManagementHost:     common.HostMachineIP,
		ManagementPort:     common.RabbitmqManagementPort,
		ManagementUserName: common.RabbitmqManagementUserName,
		ManagementPassword: common.RabbitmqManagementPassword,
	}*/

}

//AuthConfig authConfig
type AuthConfig struct {
	Accesstoken string `json:"accesstoken"`
}
