package config

import (
	"quebic-faas/common"
	"quebic-faas/config"
)

//AppConfig appConfig
type AppConfig struct {
	AppID           string                `json:"appID"`
	Auth            AuthConfig            `json:"auth"`
	MgrServerConfig config.ServerConfig   `json:"mgrServerConfig"`
	EventBusConfig  config.EventBusConfig `json:"eventBusConfig"`
}

//SavingConfig configuration going for save
type SavingConfig struct {
	Auth            AuthConfig          `json:"auth"`
	MgrServerConfig config.ServerConfig `json:"mgrServerConfig"`
}

//SetDefault set default values
func (appConfig *AppConfig) SetDefault() {

	appConfig.AppID = "quebic-lawa-cli"

	appConfig.Auth = AuthConfig{AuthToken: ""}

	appConfig.MgrServerConfig = config.ServerConfig{
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

}

//AuthConfig auth-token for connect manager
type AuthConfig struct {
	AuthToken string `json:"authToken"`
}
