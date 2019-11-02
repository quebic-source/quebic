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
	"log"
	"path/filepath"
	"quebic-faas/auth"
	"quebic-faas/common"
	"quebic-faas/config"

	"k8s.io/client-go/util/homedir"
)

//AppConfig appConfig
type AppConfig struct {
	AppID              string                `json:"appID"`
	Auth               AuthConfig            `json:"auth"`
	ServerConfig       config.ServerConfig   `json:"serverConfig"`
	DockerConfig       DockerConfig          `json:"dockerConfig"`
	KubernetesConfig   KubeConfig            `json:"kubernetesConfig"`
	EventBusConfig     config.EventBusConfig `json:"eventBusConfig"`
	APIGatewayConfig   APIGatewayConfig      `json:"apiGatewayConfig"`
	IngressConfig      IngressConfig         `json:"ingressConfig" yaml:"ingressConfig"`
	MgrDashboardConfig MgrDashboardConfig    `json:"mgrDashboardConfig"`
	InCluster          bool                  `json:"inCluster"`
	Deployment         string                `json:"deployment"`
}

//SavingConfig configuration going for save
type SavingConfig struct {
	Auth               AuthConfig            `json:"auth" yaml:"auth"`
	ServerConfig       config.ServerConfig   `json:"serverConfig" yaml:"serverConfig"`
	DockerConfig       DockerConfig          `json:"dockerConfig" yaml:"dockerConfig"`
	KubernetesConfig   KubeConfig            `json:"kubernetesConfig" yaml:"kubernetesConfig"`
	EventBusConfig     config.EventBusConfig `json:"eventBusConfig" yaml:"eventBusConfig"`
	APIGatewayConfig   APIGatewayConfig      `json:"apiGatewayConfig" yaml:"apiGatewayConfig"`
	IngressConfig      IngressConfig         `json:"ingressConfig" yaml:"ingressConfig"`
	MgrDashboardConfig MgrDashboardConfig    `json:"mgrDashboardConfig" yaml:"mgrDashboardConfig"`
	InCluster          bool                  `json:"inCluster"`
	Deployment         string                `json:"deployment" yaml:"deployment"`
}

//SetDefault set default values
func (appConfig *AppConfig) SetDefault() {

	appConfig.AppID = "quebic-faas-mgr"

	appConfig.Auth = prepareDefaultAuthConfig()

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
		Replicas: common.ComponentAPIGatewayDefaultReplicas,
	}

	appConfig.MgrDashboardConfig = MgrDashboardConfig{
		ServerConfig: config.ServerConfig{
			Host: common.HostMachineIP,
			Port: common.MgrDashboardPort,
		},
	}

	appConfig.IngressConfig = IngressConfig{}

	appConfig.DockerConfig = DockerConfig{RegistryAddress: ""}

	appConfig.KubernetesConfig = KubeConfig{ConfigPath: filepath.Join(homedir.HomeDir(), ".kube", "config")}

	appConfig.InCluster = true

	appConfig.Deployment = Deployment_Kubernetes

}

//AuthConfig authConfig
//Auth for connect manager
type AuthConfig struct {
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	JWTSecret string `json:"jwtSecret" yaml:"jwtSecret"`
}

//APIGatewayConfig apigateway config
type APIGatewayConfig struct {
	ServerConfig config.ServerConfig `json:"serverConfig" yaml:"serverConfig"`
	Replicas     int                 `json:"replicas" yaml:"replicas"`
}

//IngressConfig config for ingress controller
type IngressConfig struct {
	Provider string `json:"provider" yaml:"provider"`
	StaticIP string `json:"staticIP" yaml:"staticIP"`
}

//DockerConfig docker confog
type DockerConfig struct {
	RegistryAddress string `json:"registryAddress" yaml:"registryAddress"`
}

//KubeConfig kube confog
type KubeConfig struct {
	ConfigPath string `json:"configPath" yaml:"configPath"`
}

//MgrDashboardConfig mgrDashboardConfig config
type MgrDashboardConfig struct {
	ServerConfig config.ServerConfig `json:"serverConfig" yaml:"serverConfig"`
}

func prepareDefaultAuthConfig() AuthConfig {
	jwtSecret, err := auth.GenerateRandomString(32)
	if err != nil {
		log.Fatalf("unable to create auth secret")
	}
	return AuthConfig{
		Username:  auth.DefaultUsername,
		Password:  auth.DefaultPassword,
		JWTSecret: jwtSecret,
	}
}
