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
	"quebic-faas/common"
	"quebic-faas/types"

	bolt "github.com/coreos/bbolt"
)

//ManagerComponentSetupAPIGateway setupAPIGateway
func ManagerComponentSetupAPIGateway(db *bolt.DB) (*types.ManagerComponent, error) {

	apiGateway := &types.ManagerComponent{ID: common.ComponentAPIGateway}
	var version int

	getByID(db, apiGateway, func(savedObj []byte) error {

		if savedObj == nil {
			version = common.ComponentAPIGatewayVersionDefaultStart
			return nil
		}

		json.Unmarshal(savedObj, apiGateway)
		version = common.StrToInt(apiGateway.Version) + 1
		return nil
	})

	apiGateway.Version = common.IntToStr(version)
	err := Save(db, apiGateway)
	if err != nil {
		return nil, err
	}

	return apiGateway, nil

}

//ManagerComponentGetAPIGateway getAPIGateway
func ManagerComponentGetAPIGateway(db *bolt.DB) (*types.ManagerComponent, error) {

	apiGateway := &types.ManagerComponent{ID: common.ComponentAPIGateway}

	getByID(db, apiGateway, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("%v not found", apiGateway.GetID())
		}

		json.Unmarshal(savedObj, apiGateway)
		return nil
	})

	return apiGateway, nil

}
