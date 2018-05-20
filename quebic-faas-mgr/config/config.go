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
	"path/filepath"
	"quebic-faas/common"
	"quebic-faas/config"

	"k8s.io/client-go/util/homedir"
)

//AppConfig appConfig
type AppConfig struct {
	AppID            string                `json:"appID"`
	Auth             AuthConfig            `json:"auth"`
	ServerConfig     config.ServerConfig   `json:"serverConfig"`
	DockerConfig     DockerConfig          `json:"dockerConfig"`
	KubernetesConfig KubeConfig            `json:"kubernetesConfig"`
	EventBusConfig   config.EventBusConfig `json:"eventBusConfig"`
	APIGatewayConfig APIGatewayConfig      `json:"apiGatewayConfig"`
	Deployment       string                `json:"deployment"`
}

//SavingConfig configuration going for save
type SavingConfig struct {
	Auth             AuthConfig            `json:"auth"`
	ServerConfig     config.ServerConfig   `json:"serverConfig"`
	DockerConfig     DockerConfig          `json:"dockerConfig"`
	KubernetesConfig KubeConfig            `json:"kubernetesConfig"`
	EventBusConfig   config.EventBusConfig `json:"eventBusConfig"`
	APIGatewayConfig APIGatewayConfig      `json:"apiGatewayConfig"`
	Deployment       string                `json:"deployment"`
}

//SetDefault set default values
func (appConfig *AppConfig) SetDefault() {

	appConfig.AppID = "quebic-faas-mgr"

	appConfig.Auth = AuthConfig{Username: "admin", Password: "admin"}

	appConfig.ServerConfig = config.ServerConfig{
		Host: common.HostMachineIP,
		Port: common.MgrServerPort,
	}

	appConfig.EventBusConfig = config.EventBusConfig{
		AMQPHost:           common.HostMachineIP,
		AMQPPort:           common.RabbitmqAMQPPort,
		ManagementHost:     common.HostMachineIP,
		ManagementPort:     common.RabbitmqManagementPort,
		ManagementUserName: common.RabbitmqManagementUserName,
		ManagementPassword: common.RabbitmqManagementPassword,
	}

	appConfig.APIGatewayConfig = APIGatewayConfig{
		ServerConfig: config.ServerConfig{
			Host: common.HostMachineIP,
			Port: common.ApigatewayServerPort,
		},
	}

	appConfig.DockerConfig = DockerConfig{RegistryAddress: ""}

	appConfig.KubernetesConfig = KubeConfig{ConfigPath: filepath.Join(homedir.HomeDir(), ".kube", "config")}

	appConfig.Deployment = Deployment_Kubernetes

}

//AuthConfig authConfig
//Auth for connect manager
type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//APIGatewayConfig apigateway config
type APIGatewayConfig struct {
	ServerConfig config.ServerConfig `json:"serverConfig"`
}

//DockerConfig docker confog
type DockerConfig struct {
	RegistryAddress string `json:"registryAddress"`
}

//KubeConfig kube confog
type KubeConfig struct {
	ConfigPath string `json:"configPath"`
}
