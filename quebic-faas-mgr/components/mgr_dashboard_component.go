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

package components

import (
	"fmt"
	"log"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"

	"github.com/docker/go-connections/nat"
)

//MgrDashboardSetup MgrDashboardSetup
func MgrDashboardSetup(appConfig *config.AppConfig) error {

	componentID := common.ComponentMgrDashboard
	log.Printf("%s : starting", componentID)

	mgrDashboardPort := appConfig.MgrDashboardConfig.ServerConfig.Port

	containerDetails, err := common.DockerContainerInspect(componentID)
	if err == nil {

		//if found previous running mgr-dashboard remove it
		err = common.DockerContainerRemove(containerDetails.ID)
		if err != nil {
			return fmt.Errorf("%s-component : previous container removed failed", componentID)
		}

	}

	portMap := getMgrDashboardPortMap(mgrDashboardPort)
	exposedPort := getMgrDashboardExposedPort()

	authConfig, err := appConfig.DockerConfig.GetDockerAuthConfig()
	if err != nil {
		return err
	}

	//mgr-dashboard never run before
	_, err = common.DockerContainerStart(
		authConfig,
		componentID,
		common.MgrDashboardImage,
		exposedPort,
		portMap,
		"",
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s container start failed : %s", componentID, err.Error())
	}

	log.Printf("%s : started", componentID)

	return nil
}

func getMgrDashboardPortMap(mgrDashboardPort int) nat.PortMap {

	commonMgrDashboardPort := nat.Port(common.IntToStr(common.MgrDashboardPort))

	portMap := nat.PortMap{
		commonMgrDashboardPort: []nat.PortBinding{
			{
				HostPort: common.IntToStr(mgrDashboardPort),
			},
		},
	}

	return portMap

}

func getMgrDashboardExposedPort() nat.PortSet {

	exposedPort := nat.Port(common.IntToStr(common.MgrDashboardPort))

	exposedPortSet := nat.PortSet{
		exposedPort: struct{}{},
	}

	return exposedPortSet

}
