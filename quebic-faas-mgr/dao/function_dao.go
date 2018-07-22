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

//AddFunctionDockerImageID set DockerImageID
func AddFunctionDockerImageID(db *bolt.DB, function *types.Function, dockerImageID string) error {

	err := getByID(db, function, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found function")
		}

		json.Unmarshal(savedObj, function)

		return nil
	})

	if err != nil {
		return err
	}

	function.DockerImageID = dockerImageID
	return Save(db, function)
}

//AddFunctionLog add function log
func AddFunctionLog(db *bolt.DB, function *types.Function, log types.EntityLog, status string) error {

	log.Time = time.Now().String()

	err := getByID(db, function, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found function")
		}

		json.Unmarshal(savedObj, function)

		return nil
	})

	if err != nil {
		return err
	}

	function.Log = log
	function.Status = status
	return Save(db, function)
}
