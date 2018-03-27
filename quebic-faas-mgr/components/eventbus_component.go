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
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/docker/docker/api/types/swarm"
)

const eventbusReplicas = 1

//EventbusSetup eventbus setup
func EventbusSetup(db *bolt.DB, appConfig config.AppConfig) error {

	componentID := common.ComponentEventBus
	log.Printf("%s-component : starting", componentID)

	eventBusConfig := appConfig.EventBusConfig

	//port details always come from config settings
	eventbusAMQPPort := eventBusConfig.AMQPPort
	eventbusManagementPort := eventBusConfig.ManagementPort

	portConfig := getEventbusPortConfig(eventbusAMQPPort, eventbusManagementPort)

	envkeys := make(map[string]string)

	common.DockerServiceStop(componentID)

	//wait until current running eventbus stop
	err := WaitForEventbusStop(eventBusConfig, time.Minute*5)
	if err != nil {
		return err
	}

	_, err = common.DockerServiceCreate(
		common.DockerNetworkID,
		common.ComponentEventBus,
		common.EventbusImage,
		portConfig,
		eventbusReplicas,
		envkeys)
	if err != nil {
		return fmt.Errorf("%s-component-setup-failed : %v", componentID, err)
	}

	log.Printf("%s-component : started", componentID)

	return nil

}

func getEventbusPortConfig(amqpPortPort int, managementPort int) []swarm.PortConfig {

	targetAMQPPort := uint32(common.RabbitmqAMQPPort)
	targetManagementPort := uint32(common.RabbitmqManagementPort)

	publishAMQPPort := uint32(amqpPortPort)
	publishManagementPort := uint32(managementPort)

	portConfig := []swarm.PortConfig{
		swarm.PortConfig{
			PublishedPort: publishAMQPPort,
			TargetPort:    targetAMQPPort,
		},
		swarm.PortConfig{
			PublishedPort: publishManagementPort,
			TargetPort:    targetManagementPort,
		},
	}

	return portConfig

}
