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
	"quebic-faas/quebic-faas-mgr/config"

	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"
)

const apigatewayReplicas = 1

//ApigatewaySetup apigateway setup
func ApigatewaySetup(db *bolt.DB, appConfig *config.AppConfig) error {

	componentID := common.ComponentAPIGateway
	log.Printf("%s-component : starting", componentID)

	kubeConfig := appConfig.KubernetesConfig

	//port details always come from config settings
	apiGatewayPort := appConfig.APIGatewayConfig.ServerConfig.Port

	portConfigs := getApigatewayPortConfigs(apiGatewayPort)

	envkeys, err := getEnvVars(appConfig.EventBusConfig)
	if err != nil {
		return fmt.Errorf("%s-component-setup-failed when preparing env vars : %v", componentID, err)
	}

	spec := common.KubeServiceCreateSpec{
		AppName:     componentID,
		Dockerimage: common.ApigatewayImage,
		Envkeys:     envkeys,
		PortConfigs: portConfigs,
		Replicas:    apigatewayReplicas,
	}
	_, err = common.KubeDeploy(kubeConfig, spec)
	if err != nil {
		return fmt.Errorf("%s : %v", componentID, err)
	}

	//update config details
	service, err := common.KubeServiceGetByAppName(kubeConfig, componentID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	clusterIP := service.Spec.ClusterIP
	appConfig.APIGatewayConfig.ServerConfig.Host = clusterIP

	log.Printf("%s-component : started", componentID)

	return nil

}

func getApigatewayPortConfigs(apigatewayPort int) []common.PortConfig {

	publishApigatewayPort := common.Port(apigatewayPort)

	targetApigatewayPort := common.TargetPort(common.ApigatewayServerPort)

	portConfigs := []common.PortConfig{common.PortConfig{
		Name:       "apigateway",
		Port:       publishApigatewayPort,
		TargetPort: targetApigatewayPort,
	}}

	return portConfigs

}

func getEnvVars(eventBusConfig commonconfig.EventBusConfig) (map[string]string, error) {

	envkeys := make(map[string]string)

	//access key
	accessKeyUUID, err := uuid.NewV4()
	if err != nil {
		return envkeys, fmt.Errorf("unable to assign access_key Apigateway %v", err)
	}
	envkeys[common.EnvKeyAPIGateWayAccessKey] = accessKeyUUID.String()

	//rabbitmq config
	envkeys[common.EnvKey_rabbitmq_host] = eventBusConfig.AMQPHost
	envkeys[common.EnvKey_rabbitmq_port] = common.IntToStr(eventBusConfig.AMQPPort)
	envkeys[common.EnvKey_rabbitmq_management_username] = eventBusConfig.ManagementUserName
	envkeys[common.EnvKey_rabbitmq_management_password] = eventBusConfig.ManagementPassword

	return envkeys, nil
}
