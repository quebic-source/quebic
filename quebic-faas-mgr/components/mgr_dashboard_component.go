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
	dep "quebic-faas/quebic-faas-mgr/deployment"
)

const mgrDashboardReplicas = 1

//MgrDashboardSetup mgr ui setup
func MgrDashboardSetup(appConfig *config.AppConfig, deployment dep.Deployment) error {

	componentID := common.ComponentMgrDashboard
	log.Printf("%s : starting", componentID)

	portConfig := getMgrDashboardPortConfig()

	mgrAPIDetails, err := deployment.ListByName(common.ComponentMgrAPI)
	if err != nil {
		return fmt.Errorf("unable to get manager-api host : %v", err)
	}

	mgrAPIHost := mgrAPIDetails.Host
	mgrAPIPort := mgrAPIDetails.PortConfigs[0].Port
	mgrAPI := mgrAPIHost + ":" + common.IntToStr(int(mgrAPIPort))

	envkeys := make(map[string]string)
	envkeys[common.EnvKey_mgrAPI] = mgrAPI

	deploymentSpec := dep.Spec{
		Name:        componentID,
		Dockerimage: common.MgrDashboardImage,
		PortConfigs: portConfig,
		Envkeys:     envkeys,
		Replicas:    mgrDashboardReplicas,
	}

	_, err = deployment.CreateOrUpdate(deploymentSpec)
	if err != nil {
		return fmt.Errorf("%s setup-failed : %v", componentID, err)
	}

	log.Printf("%s : started", componentID)

	return nil

}

func getMgrDashboardPortConfig() []dep.PortConfig {

	mgrDashboardPort := dep.Port(common.MgrDashboardPort)

	portConfigs := []dep.PortConfig{
		dep.PortConfig{
			Name:       "mgr-dashboard",
			Port:       mgrDashboardPort,
			TargetPort: mgrDashboardPort,
		},
	}

	return portConfigs

}
