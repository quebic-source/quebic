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
package cmd

var authUserName string
var authPassword string

var serverHost string
var serverPort int

var dockerregistryAddress string
var kubeConfigPath string

var eventBusAMQPHost string
var eventBusAMQPPort int
var eventBusManagementHost string
var eventBusManagementPort int
var eventBusManagementUserName string
var eventBusManagementPassword string

var apigatewayServerHost string
var apigatewayServerPort int

var deployment string

var configSave bool

func setupFlags() {

	rootCmd.PersistentFlags().StringVarP(&authUserName, "auth-username", "", "", "auth-username")
	rootCmd.PersistentFlags().StringVarP(&authPassword, "auth-password", "", "", "auth-password")

	rootCmd.PersistentFlags().StringVarP(&serverHost, "server-host", "", "", "server-host")
	rootCmd.PersistentFlags().IntVarP(&serverPort, "server-port", "", 0, "server-port")

	rootCmd.PersistentFlags().StringVarP(&dockerregistryAddress, "docker-registry", "", "", "docker-registry")

	rootCmd.PersistentFlags().StringVarP(&kubeConfigPath, "kube-config-path", "", "", "kube-config-path")

	rootCmd.PersistentFlags().StringVarP(&eventBusAMQPHost, "eventbus-amqp-host", "", "", "eventbus-amqp-host")
	rootCmd.PersistentFlags().IntVarP(&eventBusAMQPPort, "eventbus-amqp-port", "", 0, "eventbus-amqp-port")
	rootCmd.PersistentFlags().StringVarP(&eventBusManagementHost, "eventbus-management-host", "", "", "eventbus-management-host")
	rootCmd.PersistentFlags().IntVarP(&eventBusManagementPort, "eventbus-management-port", "", 0, "eventbus-management-port")
	rootCmd.PersistentFlags().StringVarP(&eventBusManagementUserName, "eventbus-management-username", "", "", "eventbus-management-username")
	rootCmd.PersistentFlags().StringVarP(&eventBusManagementPassword, "eventbus-management-password", "", "", "eventbus-management-password")

	rootCmd.PersistentFlags().StringVarP(&apigatewayServerHost, "apigateway-server-host", "", "", "apigateway-server-host")
	rootCmd.PersistentFlags().IntVarP(&apigatewayServerPort, "apigateway-server-port", "", 0, "apigateway-server-port")

	rootCmd.PersistentFlags().StringVarP(&deployment, "deployment", "", "", "deployment")

	rootCmd.PersistentFlags().BoolVarP(&configSave, "config-save", "", false, "allow configurations save in .config")

}

func setFlagsToConfig() {

	appConfig := appContainer.GetAppConfig()

	if authUserName != "" {
		appConfig.Auth.Username = authUserName
	}
	if authPassword != "" {
		appConfig.Auth.Password = authPassword
	}

	if serverHost != "" {
		appConfig.ServerConfig.Host = serverHost
	}
	if serverPort != 0 {
		appConfig.ServerConfig.Port = serverPort
	}

	if dockerregistryAddress != "" {
		appConfig.DockerConfig.RegistryAddress = dockerregistryAddress
	}

	if kubeConfigPath != "" {
		appConfig.KubernetesConfig.ConfigPath = kubeConfigPath
	}

	if eventBusAMQPHost != "" {
		appConfig.EventBusConfig.AMQPHost = eventBusAMQPHost
	}

	if eventBusAMQPPort != 0 {
		appConfig.EventBusConfig.AMQPPort = eventBusAMQPPort
	}

	if eventBusManagementHost != "" {
		appConfig.EventBusConfig.ManagementHost = eventBusManagementHost
	}

	if eventBusManagementPort != 0 {
		appConfig.EventBusConfig.ManagementPort = eventBusManagementPort
	}

	if eventBusManagementUserName != "" {
		appConfig.EventBusConfig.ManagementUserName = eventBusManagementUserName
	}

	if eventBusManagementPassword != "" {
		appConfig.EventBusConfig.ManagementPassword = eventBusManagementPassword
	}

	if apigatewayServerHost != "" {
		appConfig.APIGatewayConfig.ServerConfig.Host = apigatewayServerHost
	}

	if apigatewayServerPort != 0 {
		appConfig.APIGatewayConfig.ServerConfig.Port = apigatewayServerPort
	}

	if deployment != "" {
		appConfig.Deployment = deployment
	}

}
