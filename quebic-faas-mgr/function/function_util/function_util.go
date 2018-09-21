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

package function_util

import (
	"fmt"
	"log"
	"os"
	"quebic-faas/common"
	"quebic-faas/messenger"
	quebic_messenger "quebic-faas/messenger"
	mgrconfig "quebic-faas/quebic-faas-mgr/config"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/quebic-faas-mgr/function/function_common"
	"quebic-faas/quebic-faas-mgr/function/function_create"
	"quebic-faas/quebic-faas-mgr/function/function_image"
	"quebic-faas/quebic-faas-mgr/function/function_runtime"
	quebicFaasTypes "quebic-faas/types"
	"time"

	"github.com/docker/docker/api/types"
	uuid "github.com/satori/go.uuid"
)

const functionServicePrefix string = "quebic-faas-function-"

//FunctionCreate create function
func FunctionCreate(
	authConfig types.AuthConfig,
	functionDTO quebicFaasTypes.FunctionDTO,
	functionRunTime function_runtime.FunctionRunTime) (string, error) {

	function := functionDTO.Function
	options := functionDTO.Options

	buildContextLocation, err := function_create.CreateFunction(
		function.Name,
		functionDTO.SourceFile,
		functionRunTime)

	if err != nil {
		return "", err
	}

	imageID, err := function_image.FunctionImageBuild(
		authConfig,
		buildContextLocation,
		function,
		options.Publish)

	//remove function dir
	os.RemoveAll(function_common.GetFunctionDir(function.GetID()))

	return imageID, err

}

//FunctionDeploy create-or-update function
func FunctionDeploy(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
	messenger quebic_messenger.Messenger,
	function *quebicFaasTypes.Function) (string, error) {

	if function.Life.Awake != common.FunctionLifeAwakeTypeRequest {
		return functionDeploy(appConfig, deployment, function)
	}

	//when awake type request
	functionService := GetServiceID(function.GetID())
	for _, eventID := range function.Events {

		consumerUUID, err := uuid.NewV4()
		if err != nil {
			return "", fmt.Errorf("unable to create consumerID %v", err)
		}
		consumerID := common.ConsumerFunctionRequestPrefix + consumerUUID.String()

		//create event-listener for each functions events
		err = messenger.Subscribe(eventID, func(baseEvent quebic_messenger.BaseEvent) {

			//check function status
			details, err := deployment.ListByName(functionService)
			if err == nil {
				//if function is allready running, nothing to deploy
				if common.KubeStatusTrue == details.Status {
					return
				}
			}

			//function is not running, deploy it
			_, err = functionDeploy(appConfig, deployment, function)
			if err != nil {
				log.Printf("function deploy failed : %s", err.Error())
				return
			}

			consumerID := common.ConsumerFunctionAwakePrefix + consumerUUID.String()
			functionAwakeEvent := common.EventPrefixFunctionAwake + common.EventJOIN + function.GetID()
			err = messenger.Subscribe(functionAwakeEvent, func(_ quebic_messenger.BaseEvent) {

				//publish message
				_, err = messenger.Publish(
					baseEvent.GetEventID(),
					baseEvent.GetPayloadAsObject(),
					baseEvent.GetHeaders(),
					func(response quebic_messenger.BaseEvent, statuscode int, context quebic_messenger.Context) {

						//reply success message caller
						err := messenger.ReplySuccess(baseEvent, response.GetPayloadAsObject(), statuscode)
						if err != nil {
							log.Printf("proxy success reply failed for %s, cause : %s", baseEvent.GetRequestID(), err.Error())
						}

						//release functionAwakeEvent
						messenger.ReleseQueue(functionAwakeEvent)

					},
					func(errorResponse string, statuscode int, context quebic_messenger.Context) {

						//reply error message caller
						err := messenger.ReplyError(baseEvent, errorResponse, statuscode)
						if err != nil {
							log.Printf("proxy error reply failed for %s, cause : %s", baseEvent.GetRequestID(), err.Error())
						}

						//release functionAwakeEvent
						messenger.ReleseQueue(functionAwakeEvent)

					},
					2*time.Hour,
				)
				if err != nil {
					log.Printf("unable to publish for %s, cause %v", baseEvent.GetEventID(), err)
				}

			}, consumerID)
			if err != nil {
				log.Printf("unable to subscribe for %s, cause %v", functionAwakeEvent, err)
				return
			}

		}, consumerID)
		if err != nil {
			return "", fmt.Errorf("unable to subscribe for %s, cause %v", eventID, err)
		}

	}

	return "", nil

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

//GetFunctionStatus get function current status
func GetFunctionStatus(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
	function quebicFaasTypes.Function) (string, error) {

	details, err := GetFunctionDetails(appConfig, deployment, function)
	if err != nil {
		return common.KubeStatusFalse, nil
	}

	return details.Status, nil

}

//GetFunctionDetails get function details
func GetFunctionDetails(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
	function quebicFaasTypes.Function) (dep.Details, error) {

	functionService := GetServiceID(function.GetID())

	details, err := deployment.ListByName(functionService)
	if err != nil {
		return dep.Details{}, fmt.Errorf(common.KubeStatusFalse)
	}

	return details, nil

}

//GetServiceID get function service name
func GetServiceID(functionID string) string {
	return functionServicePrefix + functionID
}

//functionDeploy callback
func functionDeploy(
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
	envkeys := prepareEnvKeys(appConfig, deployment, function)

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

func prepareEnvKeys(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
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
	envkeys[common.EnvKey_eventConst_eventPrefixFunctionAwake] = common.EventPrefixFunctionAwake
	envkeys[common.EnvKey_eventConst_eventLogListener] = common.EventRequestTracker

	//events eg: e1,e2,e3,
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

	//functionAge eg: minutes:4
	envkeys[common.EnvKey_functionAge] = function.Life.IdleState.Timeunit + ":" + common.IntToStr(function.Life.IdleState.Timeout)

	//function env
	for _, environmentVariable := range function.EnvironmentVariables {
		if environmentVariable.Name != "" {
			envkeys[environmentVariable.Name] = environmentVariable.Value
		}
	}

	//eventbox
	envkeys[common.EnvKey_eventbox_uri], _ = getEventBoxConnStr(deployment)

	return envkeys

}

func getEventBoxConnStr(deployment dep.Deployment) (string, error) {

	eventBoxDetails, err := deployment.ListByName(common.ComponentEventBox)
	if err != nil {
		return "", err
	}

	host := eventBoxDetails.Host
	port := int(eventBoxDetails.PortConfigs[0].Port)

	return fmt.Sprintf("http://%s:%d", host, port), nil

}
