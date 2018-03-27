package config

import (
	"quebic-faas/common"
)

//GetConfigDirPath quebic-faas config dir path
func GetConfigDirPath() string {
	return common.GetUserHomeDir() + common.FilepathSeparator + ConfigFileDir
}
