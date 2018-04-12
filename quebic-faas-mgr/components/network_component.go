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
	"log"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"

	bolt "github.com/coreos/bbolt"
)

//NetworkSetup docker network setup
func NetworkSetup(db *bolt.DB, appConfig config.AppConfig) error {

	_, err := common.DockerNetworkCreate(common.DockerNetworkID)
	if err != nil {
		log.Printf("docker network %s is allready exists", common.DockerNetworkID)
		return err
	}

	log.Printf("docker network %s is created", common.DockerNetworkID)

	return nil

}