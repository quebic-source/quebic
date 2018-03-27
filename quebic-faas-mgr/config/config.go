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
	KubernetesConfig common.KubeConfig     `json:"kubernetesConfig"`
	EventBusConfig   config.EventBusConfig `json:"eventBusConfig"`
	APIGatewayConfig APIGatewayConfig      `json:"apiGatewayConfig"`
	Deployment       string                `json:"deployment"`
}

//SavingConfig configuration going for save
type SavingConfig struct {
	Auth             AuthConfig            `json:"auth"`
	ServerConfig     config.ServerConfig   `json:"serverConfig"`
	DockerConfig     DockerConfig          `json:"dockerConfig"`
	KubernetesConfig common.KubeConfig     `json:"kubernetesConfig"`
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

	appConfig.KubernetesConfig = common.KubeConfig{ConfigPath: filepath.Join(homedir.HomeDir(), ".kube", "config")}

	appConfig.Deployment = Deployment_Docker

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
