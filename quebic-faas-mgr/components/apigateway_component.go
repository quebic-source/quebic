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
	"quebic-faas/quebic-faas-mgr/dao"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/types"
	"time"

	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"
)

//ApigatewaySetup apigateway setup
func ApigatewaySetup(appConfig *config.AppConfig, db *bolt.DB, deployment dep.Deployment, msg messenger.Messenger) error {

	componentID := common.ComponentAPIGateway
	log.Printf("%s : starting", componentID)

	apiGateway, err := dao.ManagerComponentSetupAPIGateway(db)
	if err != nil {
		return fmt.Errorf("%s setup-failed : %v", componentID, err)
	}

	apiGatewayVersion := apiGateway.Version
	apiGatewayDeploymentName := prepareAPIGatewayDeplymentName(apiGatewayVersion)
	apiGatewayReplicas := dep.Replicas(appConfig.APIGatewayConfig.Replicas)

	//port details always come from config settings
	apiGatewayPort := appConfig.APIGatewayConfig.ServerConfig.Port

	portConfig := getAPIGatewayPortConfig(apiGatewayPort)

	//always when starting manager, component are setup
	accessKeyUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("unable to assign access_key Apigateway %v", err)
	}

	envkeys := make(map[string]string)

	envkeys[common.EnvKey_appID] = componentID
	envkeys[common.EnvKey_deploymentID] = apiGatewayDeploymentName
	envkeys[common.EnvKey_version] = apiGatewayVersion
	envkeys[common.EnvKeyAPIGateWayAccessKey] = accessKeyUUID.String()

	//eventbus
	eventBusConfig := appConfig.EventBusConfig
	envkeys[common.EnvKey_rabbitmq_exchange] = messenger.Exchange
	envkeys[common.EnvKey_rabbitmq_host] = eventBusConfig.AMQPHost
	envkeys[common.EnvKey_rabbitmq_port] = common.IntToStr(eventBusConfig.AMQPPort)
	envkeys[common.EnvKey_rabbitmq_management_username] = eventBusConfig.ManagementUserName
	envkeys[common.EnvKey_rabbitmq_management_password] = eventBusConfig.ManagementPassword

	deploymentSpec := dep.Spec{
		Name:           componentID,
		Version:        apiGatewayVersion,
		DeploymentName: apiGatewayDeploymentName,
		Dockerimage:    common.ApigatewayImage,
		PortConfigs:    portConfig,
		Envkeys:        envkeys,
		Replicas:       apiGatewayReplicas,
	}

	err = deployment.CreateOrUpdateDeployment(deploymentSpec)
	if err != nil {
		return fmt.Errorf("%s setup-failed : %v", componentID, err)
	}

	details, err := deployment.CreateService(deploymentSpec)
	if err != nil {
		return fmt.Errorf("%s setup-failed : %v", componentID, err)
	}

	publishAPIGatewayNewVersion(apiGatewayDeploymentName, apiGatewayVersion, deployment, msg)

	host := details.Host
	appConfig.APIGatewayConfig.ServerConfig.Host = host

	log.Printf("%s : started", componentID)

	return nil

}

func prepareAPIGatewayDeplymentName(version string) string {
	return common.ComponentAPIGateway + "-v" + version
}

func getAPIGatewayPortConfig(apigatewayPort int) []dep.PortConfig {

	targetApigatewayPort := dep.Port(common.ApigatewayServerPort)

	publishApigatewayPort := dep.Port(apigatewayPort)

	portConfigs := []dep.PortConfig{
		dep.PortConfig{
			Name:       "apigateway",
			Port:       publishApigatewayPort,
			TargetPort: targetApigatewayPort,
		},
	}

	return portConfigs

}

func publishAPIGatewayNewVersion(deploymentName string, version string, deployment dep.Deployment, msg messenger.Messenger) {

	waitForResponse := make(chan bool)

	go func() {

		for {

			status, err := deployment.GetStatus(deploymentName)
			if err != nil {
				log.Printf("%s new version published failed : %v", deploymentName, err)
				continue
			}

			if common.KubeStatusTrue == status {
				waitForResponse <- true
				break
			}

			time.Sleep(time.Second * 2)

		}

	}()

	go func() {
		//wait until new version aviable
		select {
		case <-waitForResponse:

			//publish about new version
			_, err := msg.Publish(
				common.EventNewVersionAPIGateway,
				types.NewVersionMessage{
					Version: version,
				},
				nil,
				nil,
				nil,
				0,
			)
			if err != nil {
				log.Printf("%s new version published failed : %v", deploymentName, err)
				return
			}

			log.Printf("%s new version published verion : %v", common.ComponentAPIGateway, version)

		case <-time.After(time.Minute * 30):
			log.Printf("%s new version published failed : still not available", deploymentName)
		}

	}()

}
