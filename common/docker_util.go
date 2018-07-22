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
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

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

//DockerContainerInspect get container de
func DockerContainerInspect(containerID string) (*types.ContainerJSON, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("docker-client-get-failed : %v", err)
	}

	details, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("docker-container-inspect-failed : %v", err)
	}

	return &details, nil
}

//DockerImageAvilableCheck check image is avilable on locally
func DockerImageAvilableCheck(image string) error {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker-client-get-failed : %v", err)
	}

	_, _, err = cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return fmt.Errorf("docker-image-notfound : %v", err)
	}

	return nil
}

//DockerContainerStart start docker container by imageName
func DockerContainerStart(
	authConfig types.AuthConfig,
	containerName string,
	image string,
	exposedPortSet nat.PortSet,
	portMap nat.PortMap,
	networkMode container.NetworkMode,
	cmd strslice.StrSlice) (string, error) {

	authStr := GetAuthStr(authConfig)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("docker-client-get-failed : %v", err)
	}

	err = DockerImageAvilableCheck(image)
	if err != nil {

		//image unable to found locally
		log.Printf("unable to find image locally for %s", containerName)
		log.Printf("image-pull starting")

		options := types.ImagePullOptions{
			RegistryAuth: authStr,
		}

		out, err := cli.ImagePull(ctx, image, options)
		if err != nil {
			return "", fmt.Errorf("docker-image-pull-failed : %v", err)
		}
		io.Copy(os.Stdout, out)

		log.Printf("image-pull completed for %s", containerName)

	}

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

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, containerName)
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

//DockerContainerRemove remove container
func DockerContainerRemove(containerID string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker-client-get-failed : %v", err)
	}

	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("docker-container-remove-failed : %v for %s", err, containerID)
	}

	return nil
}

//GetAuthStr authStr
func GetAuthStr(authConfig types.AuthConfig) string {

	encodedJSON, _ := json.Marshal(authConfig)

	return base64.URLEncoding.EncodeToString(encodedJSON)

}
