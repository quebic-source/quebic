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
	"quebic-faas/messenger"
	"quebic-faas/quebic-faas-mgr/config"
	dep "quebic-faas/quebic-faas-mgr/deployment"

	uuid "github.com/satori/go.uuid"
)

const eventboxDBReplicas = 1
const eventboxReplicas = 1

//EventBoxSetup eventbox setup
//return eventbox-db, eventbox-server, error
func EventBoxSetup(appConfig config.AppConfig, deployment dep.Deployment) (dep.Details, dep.Details, error) {

	eventBoxDB, err := eventBoxDBSetup(deployment)
	if err != nil {
		return dep.Details{}, dep.Details{}, err
	}

	eventBoxServer, err := eventBoxServerSetup(appConfig, deployment, eventBoxDB)
	if err != nil {
		return dep.Details{}, dep.Details{}, err
	}

	return eventBoxDB, eventBoxServer, nil

}

func eventBoxDBSetup(deployment dep.Deployment) (dep.Details, error) {

	componentID := common.ComponentEventBoxDB
	log.Printf("%s : starting", componentID)

	portConfig := getEventBoxDBPortConfig()

	//always when starting manager, component are setup
	accessKeyUUID, err := uuid.NewV4()
	if err != nil {
		return dep.Details{}, fmt.Errorf("unable to assign access_key EventBoxDB %v", err)
	}

	envkeys := make(map[string]string)
	envkeys[common.EnvKeyAPIGateWayAccessKey] = accessKeyUUID.String()

	deploymentSpec := dep.Spec{
		Name:        componentID,
		Dockerimage: common.EventBoxDBImage,
		PortConfigs: portConfig,
		Envkeys:     envkeys,
		Replicas:    eventboxDBReplicas,
	}

	details, err := deployment.CreateOrUpdate(deploymentSpec)
	if err != nil {
		return dep.Details{}, fmt.Errorf("%s setup-failed : %v", componentID, err)
	}

	log.Printf("%s : started", componentID)

	return details, nil

}

func getEventBoxDBPortConfig() []dep.PortConfig {

	targetPort := dep.Port(common.EventBoxDBServerPort)

	publishPort := dep.Port(common.EventBoxDBServerPort)

	portConfigs := []dep.PortConfig{
		dep.PortConfig{
			Name:       "eventbox-db",
			Port:       publishPort,
			TargetPort: targetPort,
		},
	}

	return portConfigs

}

func eventBoxServerSetup(appConfig config.AppConfig, deployment dep.Deployment, eventboxDB dep.Details) (dep.Details, error) {

	eventBusConfig := appConfig.EventBusConfig

	componentID := common.ComponentEventBox
	log.Printf("%s : starting", componentID)

	portConfig := getEventBoxPortConfig()

	//always when starting manager, component are setup
	accessKeyUUID, err := uuid.NewV4()
	if err != nil {
		return dep.Details{}, fmt.Errorf("unable to assign access_key EventBox %v", err)
	}

	envkeys := make(map[string]string)
	envkeys[common.EnvKeyAPIGateWayAccessKey] = accessKeyUUID.String()

	//eventbus
	envkeys[common.EnvKey_rabbitmq_exchange] = messenger.Exchange
	envkeys[common.EnvKey_rabbitmq_host] = eventBusConfig.AMQPHost
	envkeys[common.EnvKey_rabbitmq_port] = common.IntToStr(eventBusConfig.AMQPPort)
	envkeys[common.EnvKey_rabbitmq_management_username] = eventBusConfig.ManagementUserName
	envkeys[common.EnvKey_rabbitmq_management_password] = eventBusConfig.ManagementPassword
	envkeys[common.EnvKey_eventConst_eventPrefixUserDefined] = common.EventPrefixUserDefined
	envkeys[common.EnvKey_eventConst_eventLogListener] = common.EventRequestTracker

	//db
	envkeys[common.EnvKey_mongo_host] = eventboxDB.Host
	envkeys[common.EnvKey_mongo_port] = common.IntToStr(common.EventBoxDBServerPort)
	//envkeys[common.EnvKey_mongo_db] = ""
	//envkeys[common.EnvKey_mongo_username] = ""
	//envkeys[common.EnvKey_mongo_password] = ""

	deploymentSpec := dep.Spec{
		Name:        componentID,
		Dockerimage: common.EventBoxImage,
		PortConfigs: portConfig,
		Envkeys:     envkeys,
		Replicas:    eventboxReplicas,
	}

	details, err := deployment.CreateOrUpdate(deploymentSpec)
	if err != nil {
		return dep.Details{}, fmt.Errorf("%s setup-failed : %v", componentID, err)
	}

	log.Printf("%s : started", componentID)

	return details, nil

}

func getEventBoxPortConfig() []dep.PortConfig {

	targetPort := dep.Port(common.EventBoxServerPort)

	publishPort := dep.Port(common.EventBoxServerPort)

	portConfigs := []dep.PortConfig{
		dep.PortConfig{
			Name:       "eventbox",
			Port:       publishPort,
			TargetPort: targetPort,
		},
	}

	return portConfigs

}
