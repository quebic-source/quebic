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
package components

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"

	bolt "github.com/coreos/bbolt"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const imageRepo string = "quebic-faas-components"
const buildContextTar string = "buildContextTar.tar"

func checkComponentISAllreadySetup(db *bolt.DB, component types.Entity) (bool, error) {

	found := true

	err := dao.GetByID(db, component, func(savedObj []byte) error {

		if savedObj == nil {
			found = false
		}

		return nil

	})
	if err != nil {
		return false, err
	}

	return found, nil
}

func createComponent(componentDirType string, dockerFileContent string) (string, error) {

	err := creatDir(componentDirType)
	if err != nil {
		return "", err
	}

	err = createDockerFile(componentDirType, dockerFileContent)
	if err != nil {
		return "", err
	}

	return prepareBuildContextLocation(componentDirType)

}

func creatDir(componentDirType string) error {

	componentDir := getComponentStoredLocation(componentDirType)

	err := os.MkdirAll(componentDir, os.FileMode.Perm(0777))
	if err != nil {
		return fmt.Errorf("%s dir creation failed %v", componentDirType, err)
	}

	log.Printf("created dir %s", componentDir)

	return nil

}

func deleteDir(componentDirType string) error {

	componentDir := getComponentStoredLocation(componentDirType)

	err := os.RemoveAll(componentDir)
	if err != nil {
		return fmt.Errorf("%s dir remove failed %v", componentDirType, err)
	}

	log.Printf("removed dir %s", componentDir)

	return nil

}

func createDockerFile(componentDirType string, dockerFileContent string) error {

	componentDockerfilePath := getDockerFilePath(componentDirType)

	err := ioutil.WriteFile(componentDockerfilePath, []byte(dockerFileContent), os.FileMode.Perm(0777))
	if err != nil {
		return fmt.Errorf("%s Dockerfile creation failed %v", componentDirType, err)
	}

	log.Printf("created Dockerfile %s", componentDockerfilePath)

	return nil

}

func prepareBuildContextLocation(componentDirType string) (string, error) {

	buildContextTar := getBuildContextTar(componentDirType)
	componentDirPath := getComponentStoredLocation(componentDirType)

	//removing previously created buildContextTar
	os.Remove(buildContextTar)

	//open function dir
	functionDir, err := os.Open(componentDirPath)
	defer functionDir.Close()
	if err != nil {
		return "", fmt.Errorf("unable to open %s dir %v", componentDirPath, err)
	}

	//get all files from function dir
	files, err := functionDir.Readdir(0)
	if err != nil {
		return "", fmt.Errorf("unable to read %s dir's files %v", componentDirPath, err)
	}

	// set up the buildContextTarFile file
	buildContextTarFile, err := os.Create(buildContextTar)
	defer buildContextTarFile.Close()
	if err != nil {
		return "", fmt.Errorf("unable to create %s buildContextTarFile %v", componentDirPath, err)
	}

	// set up the gzip writer for buildContextTarFile
	gw := gzip.NewWriter(buildContextTarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	//adding each file which stored inside functionDir
	for _, file := range files {

		if !file.IsDir() {

			filePath := componentDirPath + common.FilepathSeparator + file.Name()

			openedFile, err := os.Open(filePath)
			defer openedFile.Close()
			if err != nil {
				return "", fmt.Errorf("unable to open file %[1]s", err)
			}

			tarHeader := &tar.Header{
				Name: file.Name(),
				Mode: 0600,
				Size: file.Size(),
			}

			errorWriteHdr := tw.WriteHeader(tarHeader)
			if err != nil {
				return "", fmt.Errorf("unable to write tar header %[1]s", errorWriteHdr)
			}

			_, errorWriteTar := io.Copy(tw, openedFile)
			if err != nil {
				return "", fmt.Errorf("unable to write file into tar %[1]s", errorWriteTar)
			}

			log.Printf("added %[1]s into %[2]s", filePath, buildContextTar)

		}

	}

	return buildContextTar, nil

}

func getBuildContextTar(componentDirType string) string {
	return getComponentStoredLocation(componentDirType) + common.FilepathSeparator + buildContextTar
}

func getDockerFilePath(componentDirType string) string {
	return getComponentStoredLocation(componentDirType) + common.FilepathSeparator + "Dockerfile"
}

func getComponentStoredLocation(componentDirType string) string {
	return config.GetConfigDir() + common.FilepathSeparator + componentDirType
}

//componentImageBuild function image build
func componentImageBuild(authConfig dockerTypes.AuthConfig, buildContextLocation string, componentImageTag string, buildArgs map[string]*string) (string, error) {

	image := getImage(authConfig, componentImageTag)

	cli, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("docker client get failed %v", err)
	}

	functionImageBuildContext, err := os.Open(buildContextLocation)
	defer functionImageBuildContext.Close()
	if err != nil {
		return "", fmt.Errorf("unable to open buildContextLocation %v", err)
	}

	options := dockerTypes.ImageBuildOptions{
		Tags: []string{image},
	}

	if buildArgs != nil {
		options.BuildArgs = buildArgs
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

	return image, nil

}

//FunctionImagePublish function image publish
func componentImagePublish(authConfig dockerTypes.AuthConfig, dockerImage string) error {

	if authConfig.Username == "" {
		log.Printf("docker auth configuration not found. not going to publish")
		return nil
	}

	authStr := common.GetAuthStr(authConfig)

	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("docker client get failed %v", err)
	}

	options := dockerTypes.ImagePushOptions{
		RegistryAuth: authStr,
	}

	imagePushResponse, err := cli.ImagePush(context.Background(), dockerImage, options)
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

func getImage(authConfig dockerTypes.AuthConfig, componentImageTag string) string {

	if authConfig.Username == "" {
		return imageRepo + ":" + componentImageTag
	}

	return authConfig.Username + "/" + imageRepo + ":" + componentImageTag

}

//ComponentContainerStart component container start
func ComponentContainerStart(
	db *bolt.DB,
	component *types.ManagerComponent,
	dockerConfig config.DockerConfig,
	exposedPortSet nat.PortSet,
	portMap nat.PortMap,
	check bool,
) error {

	if check {
		//get saved component details
		err := dao.GetByID(db, component, func(savedObj []byte) error {

			if savedObj == nil {
				return fmt.Errorf("component not found : %s", component.ID)
			}

			json.Unmarshal(savedObj, component)

			return nil

		})
		if err != nil {
			return err
		}
	}

	//docker container start
	authConfig, err := dockerConfig.GetDockerAuthConfig()
	if err != nil {
		return err
	}

	dockerContainerID, err := common.DockerContainerStart(
		authConfig,
		component.DockerImageID,
		exposedPortSet,
		portMap,
		"",
		nil)

	if err != nil {
		entityLog := types.EntityLog{
			State:   common.LogStateDockerContainerStartFailed,
			Message: err.Error(),
		}
		dao.AddManagerComponentLog(db, component, entityLog)
		return err
	}

	dao.AddManagerComponentDockerContainerID(db, component, dockerContainerID)

	entityLog := types.EntityLog{
		State:   common.LogStateDockerContainerStarted,
		Message: dockerContainerID,
	}
	dao.AddManagerComponentLog(db, component, entityLog)

	log.Printf("successfully started docker-container : %s", component.ID)

	return nil

}

//ComponentContainerStop component container stop
func ComponentContainerStop(
	db *bolt.DB,
	component *types.ManagerComponent,
	check bool) error {

	if check {
		//get saved component details
		err := dao.GetByID(db, component, func(savedObj []byte) error {

			if savedObj == nil {
				return fmt.Errorf("component not found : %s", component.ID)
			}

			json.Unmarshal(savedObj, component)

			return nil

		})
		if err != nil {
			return err
		}
	}

	//stop
	if component.DockerContainerID != "" {
		//container stop
		err := common.DockerContainerStop(component.DockerContainerID)
		if err != nil {
			return err
		}

		entityLog := types.EntityLog{State: common.LogStateDockerContainerStoped}
		dao.AddManagerComponentLog(db, component, entityLog)

		dao.AddManagerComponentDockerContainerID(db, component, "")

		log.Printf("successfully stopped docker-container : %s", component.ID)

	}

	return nil

}

//ComponentContainerReStart component container re-start
func ComponentContainerReStart(
	db *bolt.DB,
	component *types.ManagerComponent,
	dockerConfig config.DockerConfig,
	exposedPortSet nat.PortSet,
	portMap nat.PortMap,
	check bool,
) error {

	if check {
		//get saved component details
		err := dao.GetByID(db, component, func(savedObj []byte) error {

			if savedObj == nil {
				return fmt.Errorf("component not found : %s", component.ID)
			}

			json.Unmarshal(savedObj, component)

			return nil

		})
		if err != nil {
			return err
		}
	}

	ComponentContainerStop(db, component, false)

	err := ComponentContainerStart(db, component, dockerConfig, exposedPortSet, portMap, check)
	if err != nil {
		return err
	}

	log.Printf("successfully re-started docker-container : %s", component.ID)

	return nil

}

//ComponentServiceStart component service start
func ComponentServiceStart(
	db *bolt.DB,
	component *types.ManagerComponent,
	dockerConfig config.DockerConfig,
	networkID string,
	portConfig []swarm.PortConfig,
	replicas int,
	check bool,
) error {

	if check {
		//get saved component details
		err := dao.GetByID(db, component, func(savedObj []byte) error {

			if savedObj == nil {
				return fmt.Errorf("component not found : %s", component.ID)
			}

			json.Unmarshal(savedObj, component)

			return nil

		})
		if err != nil {
			return err
		}
	}

	//set accesskey
	var envkeys map[string]string
	if component.AccessKey != "" {
		envkeys = make(map[string]string)
		envkeys[common.EnvKeyAPIGateWayAccessKey] = component.AccessKey
	}

	_, err := common.DockerServiceCreate(
		networkID,
		component.ID,
		component.DockerImageID,
		portConfig,
		replicas,
		envkeys,
	)
	if err != nil {
		entityLog := types.EntityLog{
			State:   common.LogStateDockerServiceStartFailed,
			Message: err.Error(),
		}
		dao.AddManagerComponentLog(db, component, entityLog)
		return err
	}

	entityLog := types.EntityLog{
		State: common.LogStateDockerServiceStarted,
	}
	dao.AddManagerComponentLog(db, component, entityLog)

	log.Printf("successfully started docker-service : %s", component.ID)

	return nil

}

//ComponentServiceReStart component service re-start
func ComponentServiceReStart(
	db *bolt.DB,
	component *types.ManagerComponent,
	dockerConfig config.DockerConfig,
	networkID string,
	portConfig []swarm.PortConfig,
	replicas int,
	check bool,
) error {

	if check {
		//get saved component details
		err := dao.GetByID(db, component, func(savedObj []byte) error {

			if savedObj == nil {
				return fmt.Errorf("component not found : %s", component.ID)
			}

			json.Unmarshal(savedObj, component)

			return nil

		})
		if err != nil {
			return err
		}
	}

	ComponentServiceStop(db, component, dockerConfig, networkID, portConfig, replicas, false)
	err := ComponentServiceStart(db, component, dockerConfig, networkID, portConfig, replicas, check)
	if err != nil {
		return err
	}

	log.Printf("successfully re-started docker-service : %s", component.ID)

	return nil

}

//ComponentServiceStop component service stop
func ComponentServiceStop(
	db *bolt.DB,
	component *types.ManagerComponent,
	dockerConfig config.DockerConfig,
	networkID string,
	portConfig []swarm.PortConfig,
	replicas int,
	check bool,
) error {

	if check {
		//get saved component details
		err := dao.GetByID(db, component, func(savedObj []byte) error {

			if savedObj == nil {
				return fmt.Errorf("component not found : %s", component.ID)
			}

			json.Unmarshal(savedObj, component)

			return nil

		})
		if err != nil {
			return err
		}
	}

	err := common.DockerServiceStop(component.ID)
	if err != nil {
		entityLog := types.EntityLog{
			State:   common.LogStateDockerServiceStartFailed,
			Message: err.Error(),
		}
		dao.AddManagerComponentLog(db, component, entityLog)
		return err
	}

	entityLog := types.EntityLog{
		State: common.LogStateDockerServiceStopped,
	}
	dao.AddManagerComponentLog(db, component, entityLog)

	log.Printf("successfully stoped docker-service : %s", component.ID)

	return nil

}
