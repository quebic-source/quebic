package config

import (
	"quebic-faas/common"
	"quebic-faas/config"
)

//GetConfigDir quebic-faas config dir
func GetConfigDir() string {
	return common.GetUserHomeDir() + common.FilepathSeparator + config.ConfigFileDir
}

//GetConfigFilePath quebic-faas config file path
func GetConfigFilePath() string {
	return GetConfigDir() + common.FilepathSeparator + ConfigFile
}

//GetDockerConfigFilePath get path .docker/config.json
func GetDockerConfigFilePath() string {
	return common.GetUserHomeDir() + common.FilepathSeparator + ".docker" + common.FilepathSeparator + "config.json"
}
