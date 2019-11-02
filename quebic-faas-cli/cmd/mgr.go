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

package cmd

import (
	"fmt"
	"path/filepath"
	"quebic-faas/common"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/quebic-faas-mgr/deployment/kube_deployment"
	"quebic-faas/types"

	"k8s.io/client-go/util/homedir"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"k8s.io/api/core/v1"
)

const quebicManagerComponentID = common.ComponentMgrAPI
const quebicManagerDockerImage = "quebicdocker/quebic-faas-mgr:0.1.0"
const quebicManagerPort = 1028

const dockerSockVolume = "docker-sock-volume"

const defaultDockerSockVolumePath = "/var/run/docker.sock"
const defaultWaitForAvailable = true

var dockerSockVolumePath string
var ingressStaticIP string

var waitForAvailable bool

func init() {
	setupManagerCompCmds()
	setupManagerCompFlags()
}

var mgrCmd = &cobra.Command{
	Use:   "manager",
	Short: "Manager commonds",
	Long:  `Manager commonds`,
}

func setupManagerCompCmds() {

	mgrCmd.AddCommand(managerStartCmd)
	mgrCmd.AddCommand(managerConnectCmd)
	mgrCmd.AddCommand(managerStatusCmd)
	mgrCmd.AddCommand(managerLogsCmd)

}

func setupManagerCompFlags() {
	mgrCmd.PersistentFlags().StringVarP(&dockerSockVolumePath, "docker_sock_path", "d", defaultDockerSockVolumePath, "docker sock path. eg: /var/run/docker.sock")
	mgrCmd.PersistentFlags().StringVarP(&ingressStaticIP, "static_ip_name", "n", "", "gce static_ip_name.")
	mgrCmd.PersistentFlags().BoolVarP(&waitForAvailable, "wait_for_available", "w", defaultWaitForAvailable, "manager status wait for available")
}

var managerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "manager : start",
	Long:  `manager : start`,
	Run: func(cmd *cobra.Command, args []string) {
		managerStart(cmd, args)
	},
}

var managerConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "manager : connect",
	Long:  `manager : connect`,
	Run: func(cmd *cobra.Command, args []string) {
		managerConnect(cmd, args)
	},
}

var managerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "manager : status",
	Long:  `manager : status`,
	Run: func(cmd *cobra.Command, args []string) {
		managerStatus(cmd, args)
	},
}

var managerLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "manager : logs",
	Long:  `manager : logs`,
	Run: func(cmd *cobra.Command, args []string) {
		managerLogs(cmd, args)
	},
}

func managerStart(cmd *cobra.Command, args []string) {
	deployment, err := getDeployment()
	if err != nil {
		prepareError(cmd, err)
	}

	err = managerDeploy(cmd, deployment)
	if err != nil {
		prepareError(cmd, err)
	}
}

func managerConnect(cmd *cobra.Command, args []string) {

	ingressDetails := managerStatus(cmd, args)

	host := ingressDetails.IP
	appContainer.GetAppConfig().MgrServerConfig.Host = host
	appContainer.GetAppConfig().MgrServerConfig.Port = 0
	appContainer.SaveConfiguration()

	color.Green("successfully connected with quebic-manager. lets start !!!")

}

func managerStatus(cmd *cobra.Command, args []string) dep.IngressDetails {

	deployment, err := getDeployment()
	if err != nil {
		prepareError(cmd, err)
	}

	details, err := deployment.ListByName(quebicManagerComponentID)
	if err != nil {
		prepareError(cmd, err)
	}

	if details.Status == common.KubeStatusFalse {
		prepareError(cmd, fmt.Errorf("manager not ready yet"))
	}

	ingressDetails, err := deployment.IngressDescribe(waitForAvailable)
	if err != nil {
		prepareError(cmd, fmt.Errorf("manager not ready yet"))
	}

	color.Green("%s is running on http://%s Host: %s", quebicManagerComponentID, ingressDetails.IP, common.IngressHostManager)

	return ingressDetails

}

func managerLogs(cmd *cobra.Command, args []string) {

	deployment, err := getDeployment()
	if err != nil {
		prepareError(cmd, err)
	}

	err = deployment.LogsByName(quebicManagerComponentID, types.FunctionContainerLogOptions{
		ReplicaIndex: 0,
	})
	if err != nil {
		prepareError(cmd, err)
	}
}

func getDeployment() (dep.Deployment, error) {

	deployment := kube_deployment.Deployment{
		Config: kube_deployment.Config{
			ConfigPath: filepath.Join(homedir.HomeDir(), ".kube", "config"),
		},
		InCluser: false,
	}

	err := deployment.Init()
	if err != nil {
		return deployment, err
	}

	return deployment, nil

}

//EventbusSetup eventbus setup
func managerDeploy(cmd *cobra.Command, deployment dep.Deployment) error {

	color.Green("%s : starting", quebicManagerComponentID)

	portConfig := getManagerPortConfigs()

	envkeys := make(map[string]string)
	envkeys[common.EnvKey_ingressConfig_staticIP] = ingressStaticIP

	volumes := []dep.Volume{
		{
			HostPath:      dockerSockVolumePath,
			ContainerPath: dockerSockVolumePath,
			HostPathType:  v1.HostPathFile,
		},
	}

	deploymentSpec := dep.Spec{
		Name:        quebicManagerComponentID,
		Dockerimage: quebicManagerDockerImage,
		PortConfigs: portConfig,
		Envkeys:     envkeys,
		Replicas:    1,
		Volumes:     volumes,
	}

	_, err := deployment.CreateOrUpdate(deploymentSpec)
	if err != nil {
		prepareError(cmd, err)
	}

	color.Green("%s deployed", quebicManagerComponentID)

	return nil

}

func getManagerPortConfigs() []dep.PortConfig {

	publishPort := dep.Port(quebicManagerPort)
	targetPort := dep.Port(quebicManagerPort)

	portConfigs := []dep.PortConfig{
		dep.PortConfig{
			Name:       "manager",
			Port:       publishPort,
			TargetPort: targetPort,
		},
	}

	return portConfigs

}
