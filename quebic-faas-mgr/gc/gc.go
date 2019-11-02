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

//Package gc gargage-collector cleanup deployments
package gc

import (
	"fmt"
	"log"
	"quebic-faas/quebic-faas-mgr/config"

	dep "quebic-faas/quebic-faas-mgr/deployment"

	bolt "github.com/coreos/bbolt"
)

//GC gargage-collector
type GC struct {
	config               config.AppConfig
	db                   *bolt.DB
	deployment           dep.Deployment
	completedWorkersPool map[string]int // deployment-id : completed workes count
}

//Init init
func (gc *GC) Init(appConfig config.AppConfig, db *bolt.DB, deployment dep.Deployment) {
	gc.config = appConfig
	gc.db = db
	gc.deployment = deployment
	gc.completedWorkersPool = make(map[string]int)
}

//SubmitJob submit deployment shutdown job
func (gc *GC) SubmitJob(deploymentID string) {

	gc.assignWorker(deploymentID)
	completedWorkersCount := gc.getCompletedWorkersCount(deploymentID)
	replicasCount, err := gc.getReplicasCount(deploymentID)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}

	log.Printf(
		"%v - gc job submitted completedWorkersCount(%v), replicasCount(%v)",
		deploymentID, completedWorkersCount,
		replicasCount,
	)

	if completedWorkersCount >= replicasCount {
		log.Printf("%v - deployment deleting", deploymentID)
		gc.deleteDeployment(deploymentID)
	}

}

func (gc *GC) getReplicasCount(deploymentID string) (int, error) {
	details, err := gc.deployment.GetDeployment(deploymentID)
	if err != nil {
		return -1, fmt.Errorf("%v - unable to get replica count", deploymentID)
	}
	return int(details.Replicas), nil
}

func (gc *GC) assignWorker(deploymentID string) {
	var workersCount int
	if gc.isDeploymentAllreadyAssigned(deploymentID) {
		workersCount = gc.getCompletedWorkersCount(deploymentID) + 1
	} else {
		workersCount = 1
	}
	gc.setCompletedWorkersCount(deploymentID, workersCount)
	log.Printf("%v - %v workers completed task", deploymentID, workersCount)
}

func (gc *GC) isDeploymentAllreadyAssigned(deploymentID string) bool {
	if gc.completedWorkersPool[deploymentID] == 0 {
		return false
	}
	return true
}

func (gc *GC) getCompletedWorkersCount(deploymentID string) int {
	return gc.completedWorkersPool[deploymentID]
}

func (gc *GC) setCompletedWorkersCount(deploymentID string, cnt int) {
	gc.completedWorkersPool[deploymentID] = cnt
}

func (gc *GC) deleteDeployment(deploymentID string) {
	go func() {
		err := gc.deployment.DeleteDeployment(deploymentID)
		if err != nil {
			log.Printf("%v - unable to deleting deployment %v", deploymentID, err)
			return
		}
		log.Printf("%v - deployment successfully deleted", deploymentID)
	}()
}

func (gc *GC) releaseDeployment(deploymentID string) {
	delete(gc.completedWorkersPool, deploymentID)
}
