/*
Copyright 2018 Tharanga Nilupul Thennakoon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package components

import (
	"fmt"
	"log"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"

	bolt "github.com/coreos/bbolt"
	"github.com/docker/docker/api/types/swarm"

	uuid "github.com/satori/go.uuid"
)

const apigatewayReplicas = 1

//ApigatewaySetup apigateway setup
func ApigatewaySetup(db *bolt.DB, appConfig config.AppConfig) error {

	componentID := common.ComponentAPIGateway
	log.Printf("%s-component : starting", componentID)

	//port details always come from config settings
	apiGatewayPort := appConfig.APIGatewayConfig.ServerConfig.Port

	portConfig := getApigatewayPortConfig(apiGatewayPort)

	//always when starting manager, component are setup
	accessKeyUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("unable to assign access_key Apigateway %v", err)
	}

	envkeys := make(map[string]string)
	envkeys[common.EnvKeyAPIGateWayAccessKey] = accessKeyUUID.String()

	_, err = common.DockerServiceReStart(
		common.DockerNetworkID,
		common.ComponentAPIGateway,
		common.ApigatewayImage,
		portConfig,
		apigatewayReplicas,
		envkeys)
	if err != nil {
		return fmt.Errorf("%s-component-setup-failed : %v", componentID, err)
	}

	log.Printf("%s-component : started", componentID)

	return nil

}

func getApigatewayPortConfig(apigatewayPort int) []swarm.PortConfig {

	targetApigatewayPort := uint32(common.ApigatewayServerPort)

	publishApigatewayPort := uint32(apigatewayPort)

	portConfig := []swarm.PortConfig{
		swarm.PortConfig{
			PublishedPort: publishApigatewayPort,
			TargetPort:    targetApigatewayPort,
		},
	}

	return portConfig

}
