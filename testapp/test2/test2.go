package main

import (
	"fmt"
	"log"
	mgrconfig "quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/deployment/kube_deployment"
)

func main() {
	deployment := setupDeployment()
	err := deployment.CreateAdminRole()
	if err != nil {
		fmt.Printf("err %v", err.Error())
		return
	}

	fmt.Printf("created")

}

func setupDeployment() kube_deployment.Deployment {

	appConfig := mgrconfig.AppConfig{}
	appConfig.SetDefault()

	kubeDeployment := kube_deployment.Deployment{
		Config: kube_deployment.Config{
			ConfigPath: appConfig.KubernetesConfig.ConfigPath,
		},
		InCluser: false,
	}

	err := kubeDeployment.Init()
	if err != nil {
		log.Fatalf("deployment init failed : %v", err)
	}

	return kubeDeployment

}
