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

package deployment

import "quebic-faas/types"

//Deployment deployment
type Deployment interface {
	Init() error
	Create(deploySpec Spec) (Details, error)
	Update(deploySpec Spec) (Details, error)
	CreateOrUpdate(deploySpec Spec) (Details, error)
	Delete(name string) error
	ListAll(filters ListFilters) ([]Details, error)
	ListByName(name string) (Details, error)
	LogsByName(name string, options types.FunctionContainerLogOptions) (string, error)
	DeploymentType() string
}

//Details about deployment details
type Details struct {
	ID          string
	Name        string
	Dockerimage string
	Replicas    Replicas
	Envkeys     map[string]string
	Host        string
	PortConfigs []PortConfig
	Pods        []Pod
	Status      string
}

//Spec deployment spec
type Spec struct {
	Name        string
	Dockerimage string
	Replicas    Replicas
	Envkeys     map[string]string
	PortConfigs []PortConfig
}

//PortConfig portConfig
type PortConfig struct {
	Name       string
	Protocol   PortProtocol
	Port       Port //expose port
	TargetPort Port //container port
}

//Pod pod spec
type Pod struct {
	Name string
}

//Port outside exposed port
type Port int

//PortProtocol port protocol
type PortProtocol string

//Replicas replicas type
type Replicas int

//ListFilters list filters
type ListFilters map[string]string
