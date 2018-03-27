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
package kube_components

import (
	"fmt"
	"log"
	"quebic-faas/common"
	commonconfig "quebic-faas/config"
	"quebic-faas/quebic-faas-mgr/components"
	"quebic-faas/quebic-faas-mgr/config"
	"time"

	bolt "github.com/coreos/bbolt"
)

const eventbusReplicas = 1

//EventbusSetup eventbus setup
func EventbusSetup(db *bolt.DB, appConfig *config.AppConfig) error {

	componentID := common.ComponentEventBus
	log.Printf("%s-component : starting", componentID)

	eventBusConfig := appConfig.EventBusConfig

	kubeConfig := appConfig.KubernetesConfig

	//port details always come from config settings
	eventbusAMQPPort := eventBusConfig.AMQPPort
	eventbusManagementPort := eventBusConfig.ManagementPort

	portConfigs := getEventbusPortConfigs(eventbusAMQPPort, eventbusManagementPort)

	envkeys := make(map[string]string)

	connectedEventBus := tryingToConnectPreviousEventBus(kubeConfig, componentID, &eventBusConfig)
	if connectedEventBus {

		spec := common.KubeServiceCreateSpec{
			AppName:     componentID,
			Dockerimage: common.EventbusImage,
			Envkeys:     envkeys,
			PortConfigs: portConfigs,
			Replicas:    eventbusReplicas,
		}
		_, err := common.KubeDeployUpdate(kubeConfig, spec)
		if err != nil {
			return fmt.Errorf("%s : %v", componentID, err)
		}

	} else {

		spec := common.KubeServiceCreateSpec{
			AppName:     componentID,
			Dockerimage: common.EventbusImage,
			Envkeys:     envkeys,
			PortConfigs: portConfigs,
			Replicas:    eventbusReplicas,
		}
		_, err := common.KubeDeploy(kubeConfig, spec)
		if err != nil {
			return fmt.Errorf("%s : %v", componentID, err)
		}

	}

	//update config details
	service, err := common.KubeServiceGetByAppName(kubeConfig, componentID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	clusterIP := service.Spec.ClusterIP
	appConfig.EventBusConfig.AMQPHost = clusterIP
	appConfig.EventBusConfig.ManagementHost = clusterIP

	log.Printf("%s-component : started", componentID)

	return nil

}

//(if connect return true, otherwise false)
func tryingToConnectPreviousEventBus(
	kubeConfig common.KubeConfig,
	componentID string,
	eventBusConfig *commonconfig.EventBusConfig) bool {

	service, err := common.KubeServiceGetByAppName(kubeConfig, componentID)
	if err != nil {
		return false
	}

	clusterIP := service.Spec.ClusterIP
	eventBusConfig.AMQPHost = clusterIP
	eventBusConfig.ManagementHost = clusterIP

	log.Printf("trying to connect running eventbus...")

	err = components.WaitForEventbusConnect(*eventBusConfig, time.Second*3)
	if err != nil {
		return false
	}

	return true

}

func loadPreviousEventBusConfig(
	kubeConfig common.KubeConfig,
	componentID string,
	eventBusConfig *commonconfig.EventBusConfig) {

	service, err := common.KubeServiceGetByAppName(kubeConfig, componentID)
	if err != nil {
		return
	}

	clusterIP := service.Spec.ClusterIP
	eventBusConfig.AMQPHost = clusterIP
	eventBusConfig.ManagementHost = clusterIP

}

func getEventbusPortConfigs(amqpPortPort int, managementPort int) []common.PortConfig {

	publishAMQPPort := common.Port(amqpPortPort)
	targetAMQPPort := common.TargetPort(common.RabbitmqAMQPPort)

	publishManagementPort := common.Port(managementPort)
	targetManagementPort := common.TargetPort(common.RabbitmqManagementPort)

	portConfigs := []common.PortConfig{
		common.PortConfig{
			Name:       "amqp",
			Port:       publishAMQPPort,
			TargetPort: targetAMQPPort,
		},
		common.PortConfig{
			Name:       "management",
			Port:       publishManagementPort,
			TargetPort: targetManagementPort,
		},
	}

	return portConfigs

}
