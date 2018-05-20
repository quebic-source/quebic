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

package dao

import (
	"encoding/json"
	"fmt"
	"quebic-faas/types"
	"time"

	bolt "github.com/coreos/bbolt"
)

//AddManagerComponentDockerImageID set DockerImageID
func AddManagerComponentDockerImageID(db *bolt.DB, component *types.ManagerComponent, dockerImageID string) error {

	err := getByID(db, component, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found component")
		}

		json.Unmarshal(savedObj, component)

		return nil
	})

	if err != nil {
		return err
	}

	component.DockerImageID = dockerImageID
	return Save(db, component)

}

//AddManagerComponentDockerContainerID set ContainerID
func AddManagerComponentDockerContainerID(db *bolt.DB, component *types.ManagerComponent, containerID string) error {

	err := getByID(db, component, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found component")
		}

		json.Unmarshal(savedObj, component)

		return nil
	})

	if err != nil {
		return err
	}

	component.DockerContainerID = containerID
	return Save(db, component)

}

//AddManagerComponentLog add log
func AddManagerComponentLog(db *bolt.DB, component *types.ManagerComponent, log types.EntityLog) error {

	log.Time = time.Now().String()

	err := getByID(db, component, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found component")
		}

		json.Unmarshal(savedObj, component)

		return nil
	})

	if err != nil {
		return err
	}

	component.Log = log
	return Save(db, component)
}
