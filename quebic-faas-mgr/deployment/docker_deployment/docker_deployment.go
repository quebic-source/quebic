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

package docker_deployment

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/deployment"
	"time"

	quebictypes "quebic-faas/types"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

//Config docker config
type Config struct {
	NetworkName string
}

//Deployment docker deployment
type Deployment struct {
	Config Config
}

const logPrefix = "docker-deployment"
const waitTimeForPreviousServiceStop = time.Minute * 5

//DeploymentType implementation for deployment.DeploymentType()
func (kubeDeployment Deployment) DeploymentType() string {
	return config.Deployment_Docker
}

//Init implementation for deployment.Init()
func (dockerDeployment Deployment) Init() error {

	//network create
	err := dockerDeployment.networkCreate(dockerDeployment.Config.NetworkName)
	if err != nil {
		return err
	}

	log.Printf("%s network-created", logPrefix)

	return nil
}

//Create implementation for deployment.Create()
func (dockerDeployment Deployment) Create(deploySpec deployment.Spec) (deployment.Details, error) {

	_, err := dockerDeployment.serviceCreate(deploySpec)
	if err != nil {
		return deployment.Details{}, err
	}

	return deployment.Details{
		Host: common.HostMachineIP,
	}, nil

}

//Update implementation for deployment.Update()
func (dockerDeployment Deployment) Update(deploySpec deployment.Spec) (deployment.Details, error) {

	err := dockerDeployment.serviceUpdate(deploySpec)
	if err != nil {
		return deployment.Details{}, err
	}

	return deployment.Details{
		Host: common.HostMachineIP,
	}, nil

}

//CreateOrUpdate implementation for deployment.CreateOrUpdate()
func (dockerDeployment Deployment) CreateOrUpdate(deploySpec deployment.Spec) (deployment.Details, error) {

	serviceName := deploySpec.Name

	service, _ := dockerDeployment.serviceGetByName(serviceName)

	if service == nil {
		_, err := dockerDeployment.Create(deploySpec)
		if err != nil {
			return deployment.Details{}, err
		}
	} else {
		err := dockerDeployment.serviceUpdate(deploySpec)
		if err != nil {
			return deployment.Details{}, err
		}
	}

	return deployment.Details{
		Host: common.HostMachineIP,
	}, nil

}

//Delete implementation for deployment.Delete()
func (dockerDeployment Deployment) Delete(name string) error {

	if name == "" {
		return fmt.Errorf("%s service name should not be empty", logPrefix)
	}

	err := dockerDeployment.serviceRemove(name)
	if err != nil {
		return err
	}

	return nil

}

//ListAll implementation for deployment.ListAll()
func (dockerDeployment Deployment) ListAll(filters deployment.ListFilters) ([]deployment.Details, error) {
	services, err := dockerDeployment.serviceList(filters)
	if err != nil {
		return nil, err
	}

	var serviceDetails []deployment.Details

	for _, service := range services {

		serviceID := service.ID
		name := service.Spec.Name
		dockerImage := service.Spec.TaskTemplate.ContainerSpec.Image
		portConfig := getPortConfigForDetails(service.Spec.EndpointSpec.Ports)
		replicas := deployment.Replicas(*service.Spec.Mode.Replicated.Replicas)

		serviceDetails = append(serviceDetails, deployment.Details{
			ID:          serviceID,
			Name:        name,
			Dockerimage: dockerImage,
			PortConfigs: portConfig,
			Replicas:    replicas,
		})

	}

	return serviceDetails, nil
}

//ListByName implementation for deployment.ListByName()
func (dockerDeployment Deployment) ListByName(name string) (deployment.Details, error) {
	return deployment.Details{}, nil
}

//LogsByName implementation for deployment.LogsByName()
func (dockerDeployment Deployment) LogsByName(name string, options quebictypes.FunctionContainerLogOptions) (string, error) {

	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return "", err
	}

	dockerLogsOptions := types.ContainerLogsOptions{
		Details:    options.Details,
		Follow:     options.Follow,
		ShowStderr: options.ShowStderr,
		ShowStdout: options.ShowStdout,
		Since:      options.Since,
		Tail:       options.Tail,
		Timestamps: options.Timestamps,
		Until:      options.Until,
	}

	if !options.ShowStdout && !options.ShowStderr {
		options.ShowStdout = true
	}

	reader, err := cli.ServiceLogs(ctx, name, dockerLogsOptions)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	return buf.String(), nil
}

func (dockerDeployment Deployment) getContextAndClient() (context.Context, *client.Client, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return ctx, cli, fmt.Errorf("%s client-get-failed %v", logPrefix, err)
	}

	return ctx, cli, nil

}

func (dockerDeployment Deployment) networkCreate(networkID string) error {

	if networkID == "" {
		return fmt.Errorf("%s network-id should not be empty", logPrefix)
	}

	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return err
	}

	options := types.NetworkCreate{
		Driver: "overlay"}

	cli.NetworkCreate(ctx, networkID, options)

	return nil

}

func (dockerDeployment Deployment) serviceCreate(deploySpec deployment.Spec) (string, error) {

	serviceSpec, err := dockerDeployment.getServiceSpec(deploySpec)
	if err != nil {
		return "", err
	}

	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return "", err
	}

	options := types.ServiceCreateOptions{}

	serviceCreateResponse, err := cli.ServiceCreate(ctx, serviceSpec, options)
	if err != nil {
		return "", fmt.Errorf("%s service-create-failed %v", logPrefix, err)
	}

	if len(serviceCreateResponse.Warnings) != 0 {
		log.Printf("%s service-create-warn %v", logPrefix, serviceCreateResponse.Warnings[0])
	}

	return serviceCreateResponse.ID, nil

}

func (dockerDeployment Deployment) serviceUpdate(deploySpec deployment.Spec) error {

	serviceSpec, err := dockerDeployment.getServiceSpec(deploySpec)
	if err != nil {
		return err
	}

	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return err
	}

	//delete previous running service
	err = cli.ServiceRemove(ctx, serviceSpec.Name)
	if err != nil {
		return err
	}

	waitForPreviousServiceStop(deploySpec.PortConfigs)

	options := types.ServiceCreateOptions{}

	//create as new service
	serviceUpdateResponse, err := cli.ServiceCreate(ctx, serviceSpec, options)
	if err != nil {
		return fmt.Errorf("%s service-update-failed %v", logPrefix, err)
	}

	if len(serviceUpdateResponse.Warnings) != 0 {
		log.Printf("%s service-update-warn %v", logPrefix, serviceUpdateResponse.Warnings[0])
	}

	return nil

}

//WaitForPreviousServiceStop waitForPreviousServiceStop
func waitForPreviousServiceStop(portConfigs []deployment.PortConfig) error {

	waitForResponse := make(chan bool)

	go func() {

		log.Printf("waiting for previous services stop ...")

		var connectedList []bool

		for {

			for _, portConfig := range portConfigs {

				connectionStr := fmt.Sprintf("%s:%d", common.HostMachineIP, portConfig.Port)

				_, err := net.Dial("tcp", connectionStr)

				//connected
				if err == nil {
					connectedList = append(connectedList, true)
				}

			}

			//check there are no any connected ports
			if len(connectedList) == 0 {
				waitForResponse <- true
				break
			}

			connectedList = nil

			time.Sleep(time.Second * 1)

		}

	}()

	select {
	case <-waitForResponse:
		break
	case <-time.After(waitTimeForPreviousServiceStop):
		return fmt.Errorf("previous service is still running. please try again")
	}

	return nil
}

func (dockerDeployment Deployment) serviceRemove(serviceName string) error {
	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return err
	}

	return cli.ServiceRemove(ctx, serviceName)

}

func (dockerDeployment Deployment) serviceGetByName(serviceName string) (*swarm.Service, error) {

	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return nil, err
	}

	options := types.ServiceInspectOptions{}
	service, _, err := cli.ServiceInspectWithRaw(ctx, serviceName, options)
	if err != nil {
		return nil, fmt.Errorf("%s service-inspect-failed %v", logPrefix, err)
	}
	return &service, nil

}

func (dockerDeployment Deployment) serviceList(filters deployment.ListFilters) ([]swarm.Service, error) {

	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return nil, err
	}

	//TODO filters
	options := types.ServiceListOptions{}

	return cli.ServiceList(ctx, options)

}

func (dockerDeployment Deployment) serviceLogs(
	serviceName string,
	options types.ContainerLogsOptions,
) (io.ReadCloser, error) {

	ctx, cli, err := dockerDeployment.getContextAndClient()
	if err != nil {
		return nil, err
	}

	if !options.ShowStdout && !options.ShowStderr {
		options.ShowStdout = true
	}

	return cli.ServiceLogs(ctx, serviceName, options)

}

func (dockerDeployment Deployment) getServiceSpec(deploySpec deployment.Spec) (swarm.ServiceSpec, error) {

	serviceName := deploySpec.Name
	dockerimage := deploySpec.Dockerimage
	portConfig := deploySpec.PortConfigs
	replicas := deploySpec.Replicas
	envkeys := deploySpec.Envkeys

	//validation
	if serviceName == "" {
		return swarm.ServiceSpec{}, fmt.Errorf("%s service-name should not be empty", logPrefix)
	}

	if dockerimage == "" {
		return swarm.ServiceSpec{}, fmt.Errorf("%s docker-image should not be empty", logPrefix)
	}

	if replicas < 1 {
		return swarm.ServiceSpec{}, fmt.Errorf("%s replicas should not be empty", logPrefix)
	}
	//validation

	networkName := dockerDeployment.Config.NetworkName
	replicasValue := uint64(replicas)

	networkAttachmentConfigs := []swarm.NetworkAttachmentConfig{
		swarm.NetworkAttachmentConfig{Target: networkName},
	}

	containerSpec := &swarm.ContainerSpec{
		Image:    dockerimage,
		Hostname: serviceName,
	}

	//env set
	var env []string
	if envkeys != nil {
		for key, value := range envkeys {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}
	if env != nil {
		containerSpec.Env = env
	}

	taskSpec := swarm.TaskSpec{
		ContainerSpec: containerSpec,
		Networks:      networkAttachmentConfigs,
	}

	return swarm.ServiceSpec{
		Annotations: swarm.Annotations{Name: serviceName},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicasValue,
			},
		},
		TaskTemplate: taskSpec,
		EndpointSpec: &swarm.EndpointSpec{Ports: getPortConfig(portConfig)},
	}, nil

}

func getPortConfig(portConfigs []deployment.PortConfig) []swarm.PortConfig {

	var swarmPortConfigs []swarm.PortConfig
	for _, portConfig := range portConfigs {

		swarmPortConfig := swarm.PortConfig{
			Name:          portConfig.Name,
			Protocol:      swarm.PortConfigProtocol(portConfig.Protocol),
			PublishedPort: uint32(portConfig.Port),
			TargetPort:    uint32(portConfig.TargetPort),
		}

		swarmPortConfigs = append(swarmPortConfigs, swarmPortConfig)

	}

	if swarmPortConfigs == nil {
		swarmPortConfigs = append(swarmPortConfigs, swarm.PortConfig{})
	}

	return swarmPortConfigs

}

func getPortConfigForDetails(portConfigs []swarm.PortConfig) []deployment.PortConfig {

	var deploymentPortConfigs []deployment.PortConfig
	for _, portConfig := range portConfigs {

		deploymentPortConfig := deployment.PortConfig{
			Name:       portConfig.Name,
			Protocol:   deployment.PortProtocol(portConfig.Protocol),
			Port:       deployment.Port(portConfig.PublishedPort),
			TargetPort: deployment.Port(portConfig.TargetPort),
		}

		deploymentPortConfigs = append(deploymentPortConfigs, deploymentPortConfig)

	}

	return deploymentPortConfigs

}
