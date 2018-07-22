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

package config

//ServerConfig manager server confog
type ServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

//EventBusConfig rabbitmq server confog
type EventBusConfig struct {
	AMQPHost           string `json:"amqpHost" yaml:"amqpHost"`
	AMQPPort           int    `json:"amqpPort" yaml:"amqpPort"`
	ManagementHost     string `json:"managementHost" yaml:"managementHost"`
	ManagementPort     int    `json:"managementPort" yaml:"managementPort"`
	ManagementUserName string `json:"managementUserName" yaml:"managementUserName"`
	ManagementPassword string `json:"managementPassword" yaml:"managementPassword"`
}
