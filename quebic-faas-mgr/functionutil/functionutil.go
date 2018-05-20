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

package functionutil

import (
	"fmt"
	"log"
	"os"
	"quebic-faas/common"
	"quebic-faas/messenger"
	mgrconfig "quebic-faas/quebic-faas-mgr/config"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/quebic-faas-mgr/functionutil/function_create"
	"quebic-faas/quebic-faas-mgr/functionutil/function_image"
	quebicFaasTypes "quebic-faas/types"

	"github.com/docker/docker/api/types"
)

const functionServicePrefix string = "quebic-faas-function-"

//FunctionCreate create function
func FunctionCreate(
	authConfig types.AuthConfig,
	functionDTO quebicFaasTypes.FunctionDTO) (string, error) {

	function := functionDTO.Function
	options := functionDTO.Options

	buildContextLocation, err := function_create.CreateFunction(
		function.Name,
		functionDTO.SourceFile,
		common.Runtime(function.Runtime))

	if err != nil {
		return "", err
	}

	imageID, err := function_image.FunctionImageBuild(
		authConfig,
		buildContextLocation,
		function,
		options.Publish)

	//remove function dir
	os.RemoveAll(function_create.GetFunctionDir(function.GetID()))

	return imageID, err

}

//FunctionDeploy create-or-update function
func FunctionDeploy(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
	function *quebicFaasTypes.Function) (string, error) {

	//validate runtime
	if !common.RuntimeValidate(common.Runtime(function.Runtime)) {
		return "", fmt.Errorf("runtime not match")
	}

	functionService := GetServiceID(function.GetID())
	functionImage := function.DockerImageID
	functionReplicas := function.Replicas

	//set accesskey
	envkeys := prepareEnvKeys(appConfig, function)

	portConfigs := []dep.PortConfig{}

	if deployment.DeploymentType() == mgrconfig.Deployment_Kubernetes {

		portConfigs = []dep.PortConfig{
			dep.PortConfig{
				Name:       "http",
				Port:       80,
				TargetPort: 80,
			},
		}
	}

	deploymentSpec := dep.Spec{
		Name:        functionService,
		Dockerimage: functionImage,
		PortConfigs: portConfigs,
		Envkeys:     envkeys,
		Replicas:    dep.Replicas(functionReplicas),
	}

	_, err := deployment.CreateOrUpdate(deploymentSpec)
	if err != nil {
		return "", err
	}

	log.Printf("%s : deployed", functionService)

	return functionService, nil

}

//StopFunction stop function
func StopFunction(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
	function *quebicFaasTypes.Function) error {

	functionService := GetServiceID(function.GetID())

	err := deployment.Delete(functionService)
	if err != nil {
		return err
	}

	log.Printf("%s : stopped", functionService)

	return nil

}

func prepareEnvKeys(
	appConfig mgrconfig.AppConfig,
	function *quebicFaasTypes.Function) map[string]string {

	eventBusConfig := appConfig.EventBusConfig

	envkeys := make(map[string]string)

	envkeys[common.EnvKey_appID] = function.GetID()

	//eventbus configuration
	if appConfig.Deployment == mgrconfig.Deployment_Docker {
		envkeys[common.EnvKey_rabbitmq_host] = common.DockerServiceEventBus
		envkeys[common.EnvKey_rabbitmq_port] = "0"
	} else {
		envkeys[common.EnvKey_rabbitmq_host] = eventBusConfig.AMQPHost
		envkeys[common.EnvKey_rabbitmq_port] = common.IntToStr(eventBusConfig.AMQPPort)
	}

	envkeys[common.EnvKey_rabbitmq_exchange] = messenger.Exchange
	envkeys[common.EnvKey_rabbitmq_management_username] = eventBusConfig.ManagementUserName
	envkeys[common.EnvKey_rabbitmq_management_password] = eventBusConfig.ManagementPassword
	envkeys[common.EnvKey_eventConst_eventPrefixUserDefined] = common.EventPrefixUserDefined
	envkeys[common.EnvKey_eventConst_eventLogListener] = common.EventRequestTracker

	//events
	eventsStr := ""
	for _, event := range function.Events {
		if eventsStr == "" {
			eventsStr = event
		} else {
			eventsStr = eventsStr + "," + event
		}
	}
	envkeys[common.EnvKey_events] = eventsStr
	envkeys[common.EnvKey_artifactLocation] = function.HandlerFile
	envkeys[common.EnvKey_functionPath] = function.HandlerPath

	return envkeys

}

//GetServiceID get function service name
func GetServiceID(functionID string) string {
	return functionServicePrefix + functionID
}
