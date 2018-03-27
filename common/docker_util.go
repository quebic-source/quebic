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
package common

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

//DockerContainerStart start docker container by imageName
func DockerContainerStart(
	authConfig types.AuthConfig,
	image string, exposedPortSet nat.PortSet,
	portMap nat.PortMap,
	networkMode container.NetworkMode,
	cmd strslice.StrSlice) (string, error) {

	authStr := GetAuthStr(authConfig)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("docker-client-get-failed : %v", err)
	}

	options := types.ImagePullOptions{
		RegistryAuth: authStr,
	}

	out, err := cli.ImagePull(ctx, image, options)
	if err != nil {
		return "", fmt.Errorf("docker-image-pull-failed : %v", err)
	}
	io.Copy(os.Stdout, out)

	config := &container.Config{
		Image:        image,
		ExposedPorts: exposedPortSet,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portMap,
	}

	if cmd != nil {
		config.Cmd = cmd
	}

	if networkMode != "" {
		hostConfig.NetworkMode = networkMode
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "")
	if err != nil {
		return "", fmt.Errorf("docker-container-create-failed : %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("docker-container-start-failed : %v", err)
	}

	return resp.ID, nil

}

//DockerContainerReStart re-start docker container by imageName
func DockerContainerReStart(containerID string) error {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker-client-get-failed : %v", err)
	}

	requestTimeout := time.Second * time.Duration(5000)

	if err := cli.ContainerRestart(ctx, containerID, &requestTimeout); err != nil {
		return fmt.Errorf("docker-container-restart-failed : %v", err)
	}

	return nil

}

//DockerContainerStop stop docker container by containerID
func DockerContainerStop(containerID string) error {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker-client-get-failed : %v", err)
	}

	err = cli.ContainerStop(ctx, containerID, nil)
	if err != nil {
		return fmt.Errorf("docker-container-stop-failed : %v", err)
	}

	return nil

}

//DockerNetworkCreate create natework
func DockerNetworkCreate(networkID string) (string, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("docker-client-get-failed : %v", err)
	}

	options := types.NetworkCreate{
		Driver: "overlay"}

	networkResponse, err := cli.NetworkCreate(ctx, networkID, options)
	if err != nil {
		return "", fmt.Errorf("docker-network-create-failed : %v", err)
	}

	if networkResponse.Warning != "" {
		log.Printf("docker-network-create-warn : %v", networkResponse.Warning)
	}

	return networkResponse.ID, nil

}

//DockerServiceCreate service create
func DockerServiceCreate(
	networkID string,
	serviceName string,
	dockerimage string,
	portConfig []swarm.PortConfig,
	replicas int,
	envkeys map[string]string) (string, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("docker client get failed %v", err)
	}

	replicasValue := uint64(replicas)

	networkAttachmentConfigs := []swarm.NetworkAttachmentConfig{
		swarm.NetworkAttachmentConfig{Target: networkID},
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

	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{Name: serviceName},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicasValue,
			},
		},
		TaskTemplate: taskSpec,
		EndpointSpec: &swarm.EndpointSpec{Ports: portConfig},
	}

	options := types.ServiceCreateOptions{}

	serviceCreateResponse, err := cli.ServiceCreate(ctx, serviceSpec, options)
	if err != nil {
		return "", fmt.Errorf("docker-service-create-failed : %v", err)
	}

	if len(serviceCreateResponse.Warnings) != 0 {
		log.Printf("docker-service-create-warn : %v", serviceCreateResponse.Warnings[0])
	}

	return serviceCreateResponse.ID, nil

}

//DockerServiceStop service create
func DockerServiceStop(serviceName string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker-client-get failed : %v", err)
	}

	return cli.ServiceRemove(ctx, serviceName)

}

//DockerServiceReStart service re-start
func DockerServiceReStart(
	networkID string,
	serviceName string,
	dockerimage string,
	portConfig []swarm.PortConfig,
	replicas int,
	envkeys map[string]string) (string, error) {

	err := DockerServiceStop(serviceName)
	if err != nil {
		log.Printf("docker-service-stop-failed : %v", err)
	}

	return DockerServiceCreate(networkID, serviceName, dockerimage, portConfig, replicas, envkeys)

}

//DockerServiceGetByName get service by name
func DockerServiceGetByName(serviceName string) (swarm.Service, []byte, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return swarm.Service{}, nil, fmt.Errorf("docker-client-get-failed : %v", err)
	}

	options := types.ServiceInspectOptions{}
	service, serviceData, err := cli.ServiceInspectWithRaw(ctx, serviceName, options)
	if err != nil {
		return swarm.Service{}, nil, fmt.Errorf("docker-service-inspect-failed : %v", err)
	}
	return service, serviceData, nil

}

//DockerServiceList service list
func DockerServiceList(serviceName string) ([]swarm.Service, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("docker-client-get-failed : %v", err)
	}

	options := types.ServiceListOptions{}
	return cli.ServiceList(ctx, options)

}

//DockerServiceLogs service list
func DockerServiceLogs(
	serviceName string,
	options types.ContainerLogsOptions,
) (io.ReadCloser, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("docker-client-get-failed : %v", err)
	}

	if !options.ShowStdout && !options.ShowStderr {
		options.ShowStdout = true
	}

	return cli.ServiceLogs(ctx, serviceName, options)

}

//GetFunctionDockerFileContent getFunctionDockerFileContent
/*func GetFunctionDockerFileContent(runtime Runtime) string {

	var dockerfile string

	fileTemplate := "docker-files/quebic-faas-function-%s-docker-file"

	if runtime == RuntimeJava {
		dockerfile = fmt.Sprintf(fileTemplate, RuntimeJava)
	} else if runtime == RuntimeNodeJS {
		dockerfile = fmt.Sprintf(fileTemplate, RuntimeNodeJS)
	} else {
		return ""
	}

	//TODO remove later this log
	log.Printf("going to read dockerfile : %s", dockerfile)

	readDockerFile, err := ioutil.ReadFile(dockerfile)
	if err != nil {
	}

	return string(readDockerFile)

}*/

//GetFunctionDockerFileContent getFunctionDockerFileContent
func GetFunctionDockerFileContent(runtime Runtime) string {

	if runtime == RuntimeJava {
		return DockerFileContent_Java
	} else if runtime == RuntimeNodeJS {
		return DockerFileContent_NodeJS
	} else {
		return ""
	}

}

//GetAuthStr authStr
func GetAuthStr(authConfig types.AuthConfig) string {

	encodedJSON, _ := json.Marshal(authConfig)

	return base64.URLEncoding.EncodeToString(encodedJSON)

}
