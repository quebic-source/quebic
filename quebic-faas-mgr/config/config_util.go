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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"quebic-faas/common"
	"quebic-faas/config"
	"strings"

	"github.com/docker/docker/api/types"
)

//GetDockerAuthConfig get docker auth data
func (dockerConfig *DockerConfig) GetDockerAuthConfig() (types.AuthConfig, error) {

	//read docker config
	dockerConfigFilePath := GetDockerConfigFilePath()
	readDockerConfig, err := ioutil.ReadFile(dockerConfigFilePath)
	if err != nil {
		return types.AuthConfig{}, fmt.Errorf("unable to open docker config")
	}

	//get encoded token
	type DockerConfig struct {
		Auths map[string]map[string]string `json:"auths"`
	}
	readDockerConfigMap := DockerConfig{}
	err = json.Unmarshal(readDockerConfig, &readDockerConfigMap)
	if err != nil {
		return types.AuthConfig{}, fmt.Errorf("unable to parse docker config %v", err)
	}

	authMap := readDockerConfigMap.Auths
	var token string
	for k, v := range authMap {
		if strings.Contains(k, "https://index.docker.io") {
			token = v["auth"]
			break
		}
	}
	if token == "" {
		return types.AuthConfig{}, fmt.Errorf("unable to get docker auth token")
	}

	//decode tocken
	decodedToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return types.AuthConfig{}, fmt.Errorf("failed to decode docker auth token %v", err)
	}

	splitetdToken := strings.Split(string(decodedToken), ":")
	if len(splitetdToken) != 2 {
		return types.AuthConfig{}, fmt.Errorf("failed to access decodedToken")
	}
	username := splitetdToken[0]
	password := splitetdToken[1]

	return types.AuthConfig{Username: username, Password: password}, nil

}

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
