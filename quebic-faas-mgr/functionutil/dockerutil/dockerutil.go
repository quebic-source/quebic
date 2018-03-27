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
package dockerutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"quebic-faas/common"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const defaultTag string = "1.0.0"
const imageTagPrefix string = "quebic-faas-function-"

//FunctionImageBuild function image build
func FunctionImageBuild(
	authConfig types.AuthConfig,
	buildContextLocation string,
	functionID string,
	secretKey string,
	publish bool) (string, error) {

	image := getImage(authConfig, functionID)

	cli, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("docker client get failed %v", err)
	}

	functionImageBuildContext, err := os.Open(buildContextLocation)
	defer functionImageBuildContext.Close()
	if err != nil {
		return "", fmt.Errorf("unable to open buildContextLocation %v", err)
	}

	//set accessKey into function container
	buildArgs := make(map[string]*string)
	buildArgs[common.EnvKeyAPIGateWayAccessKey] = &secretKey

	options := types.ImageBuildOptions{
		Tags:      []string{image},
		BuildArgs: buildArgs,
	}

	imageBuildResponse, err := cli.ImageBuild(context.Background(), functionImageBuildContext, options)
	defer imageBuildResponse.Body.Close()
	if err != nil {
		return "", fmt.Errorf("image-build failed %v", err)
	}

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return "", fmt.Errorf("image-build-response-read failed %v", err)
	}

	if publish {
		err = functionImagePublish(authConfig, functionID)
		if err != nil {
			return "", err
		}
	}

	return image, nil

}

//FunctionImagePublish function image publish
func functionImagePublish(authConfig types.AuthConfig, functionID string) error {

	if authConfig.Username == "" {
		log.Printf("docker auth configuration not found. not going to publish")
		return nil
	}

	image := getImage(authConfig, functionID)
	authStr := getAuthStr(authConfig)

	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker client get failed %v", err)
	}

	options := types.ImagePushOptions{
		RegistryAuth: authStr,
	}

	imagePushResponse, err := cli.ImagePush(context.Background(), image, options)
	defer imagePushResponse.Close()
	if err != nil {
		return fmt.Errorf("image-push failed %v", err)
	}

	_, err = io.Copy(os.Stdout, imagePushResponse)
	if err != nil {
		return fmt.Errorf("image-push-response-read failed %v", err)
	}

	return nil

}

//FunctionContainerStart function container start
func FunctionContainerStart(authConfig types.AuthConfig, image string) error {

	authStr := getAuthStr(authConfig)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker client get failed %v", err)
	}

	options := types.ImagePullOptions{
		RegistryAuth: authStr,
	}

	_, err = cli.ImagePull(ctx, image, options)
	if err != nil {
		return fmt.Errorf("docker image pull failed %v", err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("docker container create failed %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("docker container start failed %v", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("docker container wait failed %v", err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return fmt.Errorf("docker container logs failed %v", err)
	}

	io.Copy(os.Stdout, out)

	return nil

}

func getImage(authConfig types.AuthConfig, functionID string) string {

	if authConfig.Username == "" {
		return imageTagPrefix + functionID + ":" + defaultTag
	}

	return authConfig.Username + "/" + imageTagPrefix + functionID + ":" + defaultTag

}

func getAuthStr(authConfig types.AuthConfig) string {

	encodedJSON, _ := json.Marshal(authConfig)

	return base64.URLEncoding.EncodeToString(encodedJSON)

}
