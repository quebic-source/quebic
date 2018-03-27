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
	"encoding/json"
	"log"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"

	bolt "github.com/coreos/bbolt"
	"github.com/docker/go-connections/nat"
)

const dockerRegistryImage string = "registry:2"

//DockerRegistrySetup docker registry setup
func DockerRegistrySetup(db *bolt.DB, appConfig config.AppConfig) error {

	//check for saved dockerRegistry
	var found bool

	id := common.ComponentDockerRegistry
	log.Printf("%s component setup starting", id)

	dockerRegistry := &types.ManagerComponent{ID: id}

	//port details always come from config settings
	dockerRegistryPort := "5000"

	_ = dao.GetByID(db, &types.ManagerComponent{ID: id}, func(savedObj []byte) error {

		if savedObj == nil {
			found = false
		} else {
			found = true
			json.Unmarshal(savedObj, dockerRegistry)
		}

		return nil

	})

	if found {

		err := ComponentContainerReStart(
			db,
			dockerRegistry,
			appConfig.DockerConfig,
			getDockerRegistryExposedPort(),
			getDockerRegistryPortMap(dockerRegistryPort),
			false)

		if err != nil {
			return err
		}

	} else {

		dockerRegistry := &types.ManagerComponent{
			ID:            id,
			DockerImageID: dockerRegistryImage,
		}

		//save
		err := dao.Add(db, dockerRegistry)
		if err != nil {
			return err
		}

		//log saved
		entityLog := types.EntityLog{State: common.LogStateSaved}
		err = dao.AddManagerComponentLog(db, dockerRegistry, entityLog)
		if err != nil {
			return err
		}

		err = ComponentContainerStart(
			db,
			dockerRegistry,
			appConfig.DockerConfig,
			getDockerRegistryExposedPort(),
			getDockerRegistryPortMap(dockerRegistryPort),
			false,
		)
		if err != nil {
			return err
		}

	}

	return nil

}

func getDockerRegistryPortMap(dockerRegistryPort string) nat.PortMap {

	portMap := nat.PortMap{
		common.DockerRegistryPort: []nat.PortBinding{
			{
				HostPort: dockerRegistryPort,
			},
		},
	}

	return portMap

}

func getDockerRegistryExposedPort() nat.PortSet {

	exposedPort := nat.Port(common.DockerRegistryPort)

	exposedPortSet := nat.PortSet{
		exposedPort: struct{}{},
	}

	return exposedPortSet

}
