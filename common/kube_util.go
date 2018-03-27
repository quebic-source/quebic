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
package common

import (
	"fmt"
	"log"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	kubev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeSelecterKey = "quebic-faas-app"
const kubeNamespace = "quebic-faas"

//KubeServiceCreateSpec kubeServiceCreateSpec
type KubeServiceCreateSpec struct {
	AppName     string
	Dockerimage string
	Replicas    int
	Envkeys     map[string]string
	PortConfigs []PortConfig
}

//KubeConfig kubernates config
type KubeConfig struct {
	ConfigPath string `json:"configPath"`
}

//Port outside exposed port
type Port int

//TargetPort container exposed port
type TargetPort int

//PortConfig portConfig
type PortConfig struct {
	Name       string
	Port       Port
	TargetPort TargetPort
}

//KubeNameSpaceCreate create namespace
func KubeNameSpaceCreate(
	kubeConfig KubeConfig) error {

	clientset, err := kubeClient(kubeConfig)
	if err != nil {
		return err
	}

	_, err = clientset.Core().Namespaces().Get(kubeNamespace, metav1.GetOptions{})
	if err != nil {

		if errors.IsNotFound(err) {

			_, err = clientset.Core().Namespaces().Create(&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: kubeNamespace,
				},
			})
			if err != nil {
				return fmt.Errorf("kube-namespace-create failed : %v", err)
			}

			log.Printf("kube-namespace-created : %s", kubeNamespace)

		} else {
			return err
		}
	}

	return nil

}

//KubeDeploy create or update
func KubeDeploy(
	kubeConfig KubeConfig,
	spec KubeServiceCreateSpec) (*v1.Service, error) {

	clientset, err := kubeClient(kubeConfig)
	if err != nil {
		return nil, err
	}

	err = kubeDeploymentUpdate(clientset, spec)
	if err != nil {

		if errors.IsNotFound(err) {

			err = kubeDeploymentCreate(clientset, spec)
			if err != nil {
				log.Printf("kube-deployment-create-failed : %v", err)
				return nil, err
			}

		} else {
			log.Printf("kube-deployment-update-failed : %v", err)
			return nil, err
		}

	}

	service, err := kubeServiceUpdate(clientset, spec)
	if err != nil {
		if errors.IsNotFound(err) {

			svc, err := kubeServiceCrete(clientset, spec)
			if err != nil {
				log.Printf("kube-service-create-failed : %v", err)
				return nil, err
			}

			return svc, err

		}

		log.Printf("kube-service-update-failed : %v", err)
		return nil, err

	}

	return service, nil

}

//KubeDeployCreate kube-update
func KubeDeployCreate(
	kubeConfig KubeConfig,
	spec KubeServiceCreateSpec) (*v1.Service, error) {

	clientset, err := kubeClient(kubeConfig)
	if err != nil {
		return nil, err
	}

	err = kubeDeploymentCreate(clientset, spec)
	if err != nil {
		return nil, err
	}

	service, err := kubeServiceCrete(clientset, spec)
	if err != nil {
		return nil, err
	}

	return service, nil

}

//KubeDeployUpdate kube-update
func KubeDeployUpdate(
	kubeConfig KubeConfig,
	spec KubeServiceCreateSpec) (*v1.Service, error) {

	clientset, err := kubeClient(kubeConfig)
	if err != nil {
		return nil, err
	}

	err = kubeDeploymentUpdate(clientset, spec)
	if err != nil {
		return nil, err
	}

	service, err := kubeServiceUpdate(clientset, spec)
	if err != nil {
		return nil, err
	}

	return service, nil

}

//KubeServiceGetByAppName get by appName
func KubeServiceGetByAppName(
	kubeConfig KubeConfig,
	appName string) (*v1.Service, error) {

	clientset, err := kubeClient(kubeConfig)
	if err != nil {
		return nil, err
	}

	service := clientset.Core().Services(kubeNamespace)

	return kubeServiceGetByAppName(service, appName)

}

//KubeServiceDeleteByAppName delete by appName
func KubeServiceDeleteByAppName(
	kubeConfig KubeConfig,
	appName string) error {

	clientset, err := kubeClient(kubeConfig)
	if err != nil {
		return err
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)
	err = deploymentsClient.Delete(appName, &metav1.DeleteOptions{})
	if err != nil {
		log.Printf("kube-deployment delete failed : %v", err)
	}

	service := clientset.Core().Services(kubeNamespace)
	err = service.Delete(appName, &metav1.DeleteOptions{})
	if err != nil {
		log.Printf("kube-service delete failed : %v", err)
	}

	return nil

}

func kubeClient(kubeConfig KubeConfig) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("kube-config config getting failed %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("kube-config config failed %v", err)
	}

	return clientset, nil

}

func getKubeDeploymentSpec(
	clientset *kubernetes.Clientset,
	spec KubeServiceCreateSpec) *appsv1beta1.Deployment {

	appName := spec.AppName
	dockerimage := spec.Dockerimage
	replicas := int32Ptr(int32(spec.Replicas))

	//ports
	var ports []apiv1.ContainerPort
	for _, portConfig := range spec.PortConfigs {
		ports = append(ports, apiv1.ContainerPort{
			Name:          portConfig.Name,
			Protocol:      apiv1.ProtocolTCP,
			ContainerPort: int32(portConfig.TargetPort),
		})
	}

	//env variabls
	var envVar []v1.EnvVar
	for k, v := range spec.Envkeys {
		envVar = append(envVar, v1.EnvVar{Name: k, Value: v})
	}

	return &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   appName,
			Labels: map[string]string{kubeSelecterKey: appName},
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: replicas,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						kubeSelecterKey: appName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  appName,
							Image: dockerimage,
							Ports: ports,
							Env:   envVar,
						},
					},
				},
			},
		},
	}

}

func kubeDeploymentCreate(
	clientset *kubernetes.Clientset,
	spec KubeServiceCreateSpec) error {

	deploymentSpec := getKubeDeploymentSpec(clientset, spec)

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	_, err := deploymentsClient.Create(deploymentSpec)
	if err != nil {
		return err
	}

	log.Println("kube-deployment : created")

	return nil

}

func kubeDeploymentUpdate(
	clientset *kubernetes.Clientset,
	spec KubeServiceCreateSpec) error {

	deploymentSpec := getKubeDeploymentSpec(clientset, spec)

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	_, err := deploymentsClient.Update(deploymentSpec)
	if err != nil {
		return err
	}

	log.Println("kube-deployment : updated")

	return nil

}

func getKubeServiceSpec(
	clientset *kubernetes.Clientset,
	spec KubeServiceCreateSpec) *v1.Service {

	appName := spec.AppName

	//ports
	var ports []v1.ServicePort
	for _, portConfig := range spec.PortConfigs {
		ports = append(ports, v1.ServicePort{
			Name:     portConfig.Name,
			Protocol: v1.ProtocolTCP,
			Port:     int32(portConfig.Port),
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(portConfig.TargetPort),
			},
		})
	}

	// Define service spec.
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Selector: map[string]string{kubeSelecterKey: appName},
			Ports:    ports,
		},
	}

}

func kubeServiceCrete(
	clientset *kubernetes.Clientset,
	spec KubeServiceCreateSpec) (*v1.Service, error) {

	serviceSpec := getKubeServiceSpec(clientset, spec)

	service := clientset.Core().Services(kubeNamespace)

	serviceCreated, err := service.Create(serviceSpec)
	if err != nil {
		return nil, err
	}

	log.Println("kube-service : created")

	return serviceCreated, nil

}

func kubeServiceUpdate(
	clientset *kubernetes.Clientset,
	spec KubeServiceCreateSpec) (*v1.Service, error) {

	appName := spec.AppName

	serviceSpec := getKubeServiceSpec(clientset, spec)

	service := clientset.Core().Services(kubeNamespace)
	svc, err := kubeServiceGetByAppName(service, appName)
	if err != nil {
		return nil, err
	}

	serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
	serviceSpec.Spec.ClusterIP = svc.Spec.ClusterIP

	serviceUpdated, err := service.Update(serviceSpec)
	if err != nil {
		return nil, err
	}

	log.Println("kube-service : updated")

	return serviceUpdated, nil

}

func kubeServiceGetByAppName(service kubev1.ServiceInterface, appName string) (*v1.Service, error) {

	svc, err := service.Get(appName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return svc, nil

}

func int32Ptr(i int32) *int32 { return &i }
