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
		return functionDeploy(appConfig, deployment, messenger, function)
	}

	//when awake type request
	functionDeploymentID := GetDeploymentID(*function)
	for _, eventID := range function.Events {

		consumerID := common.ConsumerFunctionRequestPrefix

		//create event-listener for each functions events
		err := messenger.Subscribe(eventID, func(baseEvent quebic_messenger.BaseEvent) {

			//check function status
			functionStatus, err := deployment.GetStatus(functionDeploymentID)
			if err == nil {
				//if function is allready running, nothing to deploy
				if common.KubeStatusTrue == functionStatus {
					return
				}
			}

			//function is not running, deploy it
			_, err = functionDeploy(appConfig, deployment, messenger, function)
			if err != nil {
				log.Printf("function deploy failed : %s", err.Error())
				return
			}

			consumerID := common.ConsumerFunctionAwakePrefix
			functionAwakeEvent := GetFunctionAwakeEvent(*function)
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

	functionDeploymentID := GetDeploymentID(*function)

	err := deployment.DeleteDeployment(functionDeploymentID)
	if err != nil {
		return err
	}

	log.Printf("%s : deleted", functionDeploymentID)

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

	functionDeploymentID := GetDeploymentID(function)

	details, err := deployment.GetDeployment(functionDeploymentID)
	if err != nil {
		return dep.Details{}, fmt.Errorf(common.KubeStatusFalse)
	}

	return details, nil

}

//GetID get function id
func GetID(function quebicFaasTypes.Function) string {
	return functionServicePrefix + function.GetID()
}

//GetDeploymentID get function deployment id
func GetDeploymentID(function quebicFaasTypes.Function) string {
	return GetID(function) + "-" + function.Version
}

//GetFunctionAwakeEvent get function awake event
func GetFunctionAwakeEvent(function quebicFaasTypes.Function) string {
	return common.EventPrefixFunctionAwakePrefix + function.GetID()
}

//GetNewVersionEvent get new version notify event
func GetNewVersionEvent(function quebicFaasTypes.Function) string {
	return common.EventNewVersionFunctionPrefix + function.GetID()
}

//functionDeploy callback
func functionDeploy(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
	msg messenger.Messenger,
	function *quebicFaasTypes.Function) (string, error) {

	//validate runtime
	if !common.RuntimeValidate(common.Runtime(function.Runtime)) {
		return "", fmt.Errorf("runtime not match")
	}

	functionID := GetID(*function)
	functionDeploymentID := GetDeploymentID(*function)
	functionVersion := function.Version

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
		Name:            functionID,
		DeploymentName:  functionDeploymentID,
		Version:         functionVersion,
		Dockerimage:     functionImage,
		PortConfigs:     portConfigs,
		Envkeys:         envkeys,
		Replicas:        dep.Replicas(functionReplicas),
		ImagePullPolicy: "IfNotPresent",
	}

	err := deployment.CreateOrUpdateDeployment(deploymentSpec)
	if err != nil {
		return "", err
	}

	_, err = deployment.CreateService(deploymentSpec)
	if err != nil {
		return "", err
	}

	publishFunctionNewVersion(deployment, msg, *function)

	log.Printf("%s : function is deployed", functionDeploymentID)

	return functionDeploymentID, nil

}

func prepareEnvKeys(
	appConfig mgrconfig.AppConfig,
	deployment dep.Deployment,
	function *quebicFaasTypes.Function) map[string]string {

	eventBusConfig := appConfig.EventBusConfig

	envkeys := make(map[string]string)

	envkeys[common.EnvKey_appID] = function.GetID()
	envkeys[common.EnvKey_deploymentID] = GetDeploymentID(*function)
	envkeys[common.EnvKey_version] = function.Version

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
	envkeys[common.EnvKey_eventConst_eventLog] = common.EventRequestTracker
	envkeys[common.EnvKey_eventConst_eventFunctionAwake] = GetFunctionAwakeEvent(*function)
	envkeys[common.EnvKey_eventConst_eventDataFetch] = common.EventFunctionDataFetch
	envkeys[common.EnvKey_eventConst_eventNewVersion] = GetNewVersionEvent(*function)
	envkeys[common.EnvKey_eventConst_eventShutDownRequest] = common.EventShutDownRequest

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

func publishFunctionNewVersion(
	deployment dep.Deployment,
	msg messenger.Messenger,
	function quebicFaasTypes.Function) {

	functionDeploymentID := GetDeploymentID(function)
	waitForResponse := make(chan bool)

	go func() {

		for {

			status, err := deployment.GetStatus(functionDeploymentID)
			if err != nil {
				log.Printf("%s new version published failed : %v", functionDeploymentID, err)
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
				GetNewVersionEvent(function),
				quebicFaasTypes.NewVersionMessage{
					Version: function.Version,
				},
				nil,
				nil,
				nil,
				0,
			)
			if err != nil {
				log.Printf("%s new version published failed : %v", functionDeploymentID, err)
				return
			}

			log.Printf("%s new version published verion : %v", functionDeploymentID, function.Version)

		case <-time.After(time.Minute * 30):
			log.Printf("%s new version published failed : still not available", functionDeploymentID)
		}

	}()

}
