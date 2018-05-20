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

package components

import (
	"fmt"
	"log"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	dep "quebic-faas/quebic-faas-mgr/deployment"
)

const eventbusReplicas = 1

//EventbusSetup eventbus setup
func EventbusSetup(appConfig *config.AppConfig, deployment dep.Deployment) error {

	componentID := common.ComponentEventBus
	log.Printf("%s-component : starting", componentID)

	eventBusConfig := appConfig.EventBusConfig

	//port details always come from config settings
	eventbusAMQPPort := eventBusConfig.AMQPPort
	eventbusManagementPort := eventBusConfig.ManagementPort

	portConfig := getEventbusPortConfigs(eventbusAMQPPort, eventbusManagementPort)

	envkeys := make(map[string]string)

	deploymentSpec := dep.Spec{
		Name:        componentID,
		Dockerimage: common.EventbusImage,
		PortConfigs: portConfig,
		Envkeys:     envkeys,
		Replicas:    eventbusReplicas,
	}

	details, err := deployment.CreateOrUpdate(deploymentSpec)
	if err != nil {
		return fmt.Errorf("%s-component-setup-failed : %v", componentID, err)
	}

	host := details.Host
	appConfig.EventBusConfig.AMQPHost = host
	appConfig.EventBusConfig.ManagementHost = host

	log.Printf("%s-component : started", componentID)

	return nil

}

func getEventbusPortConfigs(amqpPortPort int, managementPort int) []dep.PortConfig {

	publishAMQPPort := dep.Port(amqpPortPort)
	targetAMQPPort := dep.Port(common.RabbitmqAMQPPort)

	publishManagementPort := dep.Port(managementPort)
	targetManagementPort := dep.Port(common.RabbitmqManagementPort)

	portConfigs := []dep.PortConfig{
		dep.PortConfig{
			Name:       "amqp",
			Port:       publishAMQPPort,
			TargetPort: targetAMQPPort,
		},
		dep.PortConfig{
			Name:       "management",
			Port:       publishManagementPort,
			TargetPort: targetManagementPort,
		},
	}

	return portConfigs

}
