package config

import (
	"quebic-faas/common"
	"quebic-faas/config"
)

//AppConfig appConfig
type AppConfig struct {
	AppID          string                `json:"appID"`
	Auth           AuthConfig            `json:"auth"`
	ServerConfig   config.ServerConfig   `json:"serverConfig"`
	EventBusConfig config.EventBusConfig `json:"eventBusConfig"`
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
