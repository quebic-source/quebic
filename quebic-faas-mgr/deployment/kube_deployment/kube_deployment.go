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

package kube_deployment

import (
	"fmt"
	"log"
	"quebic-faas/quebic-faas-mgr/deployment"

	"quebic-faas/quebic-faas-mgr/config"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	kubev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	quebictypes "quebic-faas/types"
)

const kubeSelecterKey = "quebic-faas-app"
const kubeNamespace = "quebic-faas"

//Config kubernates config
type Config struct {
	ConfigPath    string        `json:"configPath"`
	IngressConfig IngressConfig `json:"ingressConfig"`
}

//IngressConfig config for ingress controller
type IngressConfig struct {
	StaticIP string `json:"staticIP"`
}

//Deployment kube deployment
type Deployment struct {
	Config Config
}

//DeploymentType implementation for deployment.DeploymentType()
func (kubeDeployment Deployment) DeploymentType() string {
	return config.Deployment_Kubernetes
}

//Init implementation for deployment.Init()
func (kubeDeployment Deployment) Init() error {

	err := kubeDeployment.nameSpaceCreate()
	if err != nil {
		return err
	}

	return nil

}

//Create implementation for deployment.Create()
func (kubeDeployment Deployment) Create(deploySpec deployment.Spec) (deployment.Details, error) {

	service, err := kubeDeployment.deployCreate(deploySpec)
	if err != nil {
		return deployment.Details{}, err
	}

	return deployment.Details{
		Host: service.Spec.ClusterIP,
	}, nil

}

//Update implementation for deployment.Update()
func (kubeDeployment Deployment) Update(deploySpec deployment.Spec) (deployment.Details, error) {

	service, err := kubeDeployment.deployUpdate(deploySpec)
	if err != nil {
		return deployment.Details{}, err
	}

	return deployment.Details{
		Host: service.Spec.ClusterIP,
	}, nil

}

//CreateOrUpdate implementation for deployment.CreateOrUpdate()
func (kubeDeployment Deployment) CreateOrUpdate(deploySpec deployment.Spec) (deployment.Details, error) {

	service, err := kubeDeployment.deployCreateOrUpdate(deploySpec)
	if err != nil {
		return deployment.Details{}, err
	}

	return deployment.Details{
		Host: service.Spec.ClusterIP,
	}, nil

}

//Delete implementation for deployment.Delete()
func (kubeDeployment Deployment) Delete(name string) error {

	err := kubeDeployment.serviceDeleteByAppName(name)
	if err != nil {
		return err
	}

	return nil
}

//ListAll implementation for deployment.ListAll()
func (kubeDeployment Deployment) ListAll(filters deployment.ListFilters) ([]deployment.Details, error) {
	return nil, nil
}

//ListByName implementation for deployment.ListByName()
func (kubeDeployment Deployment) ListByName(name string) (deployment.Details, error) {
	return deployment.Details{}, nil
}

//LogsByName implementation for deployment.LogsByName()
func (kubeDeployment Deployment) LogsByName(name string, options quebictypes.FunctionContainerLogOptions) (string, error) {
	return "", nil
}

//KubeNameSpaceCreate create namespace
func (kubeDeployment Deployment) nameSpaceCreate() error {

	clientset, err := kubeDeployment.getClient()
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

//deployCreateOrUpdate create or update
//kube-deployment create or update
//kube-service create or update
func (kubeDeployment Deployment) deployCreateOrUpdate(spec deployment.Spec) (*v1.Service, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	//deployment
	err = kubeDeployment.deploymentUpdate(clientset, spec)
	if err != nil {

		if errors.IsNotFound(err) {

			err = kubeDeployment.deploymentCreate(clientset, spec)
			if err != nil {
				log.Printf("kube-deployment-create-failed : %v", err)
				return nil, err
			}

		} else {
			log.Printf("kube-deployment-update-failed : %v", err)
			return nil, err
		}

	}

	//service
	service, err := kubeDeployment.serviceUpdate(clientset, spec)
	if err != nil {
		if errors.IsNotFound(err) {

			svc, err := kubeDeployment.serviceCreate(clientset, spec)
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

//deployCreate create
//kube-deployment create
//kube-service create
func (kubeDeployment Deployment) deployCreate(spec deployment.Spec) (*v1.Service, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	err = kubeDeployment.deploymentCreate(clientset, spec)
	if err != nil {
		return nil, err
	}

	service, err := kubeDeployment.serviceCreate(clientset, spec)
	if err != nil {
		return nil, err
	}

	return service, nil

}

//deployUpdate update
//kube-deployment update
//kube-service update
func (kubeDeployment Deployment) deployUpdate(spec deployment.Spec) (*v1.Service, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	err = kubeDeployment.deploymentUpdate(clientset, spec)
	if err != nil {
		return nil, err
	}

	service, err := kubeDeployment.serviceUpdate(clientset, spec)
	if err != nil {
		return nil, err
	}

	return service, nil

}

//KubeServiceGetByAppName get by appName
/*func (kubeDeployment Deployment) KubeServiceGetByAppName(
	appName string) (*v1.Service, error) {

	clientset, err := kubeDeployment.kubeClient()
	if err != nil {
		return nil, err
	}

	service := clientset.Core().Services(kubeNamespace)

	return kubeDeployment.kubeServiceGetByAppName(service, appName)

}*/

//serviceDeleteByAppName delete by appName
func (kubeDeployment Deployment) serviceDeleteByAppName(appName string) error {

	clientset, err := kubeDeployment.getClient()
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

func (kubeDeployment Deployment) getClient() (*kubernetes.Clientset, error) {

	kubeConfig := kubeDeployment.Config

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

func getDeploymentSpec(
	clientset *kubernetes.Clientset,
	spec deployment.Spec) *appsv1beta1.Deployment {

	appName := spec.Name
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

//kube-deployment create
func (kubeDeployment Deployment) deploymentCreate(
	clientset *kubernetes.Clientset,
	spec deployment.Spec) error {

	deploymentSpec := getDeploymentSpec(clientset, spec)

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	_, err := deploymentsClient.Create(deploymentSpec)
	if err != nil {
		return err
	}

	log.Println("kube-deployment : created")

	return nil

}

//kube-deployment update
func (kubeDeployment Deployment) deploymentUpdate(
	clientset *kubernetes.Clientset,
	spec deployment.Spec) error {

	deploymentSpec := getDeploymentSpec(clientset, spec)

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	_, err := deploymentsClient.Update(deploymentSpec)
	if err != nil {
		return err
	}

	log.Println("kube-deployment : updated")

	return nil

}

func getServiceSpec(
	clientset *kubernetes.Clientset,
	spec deployment.Spec) *v1.Service {

	appName := spec.Name

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

//kube-service create
func (kubeDeployment Deployment) serviceCreate(
	clientset *kubernetes.Clientset,
	spec deployment.Spec) (*v1.Service, error) {

	serviceSpec := getServiceSpec(clientset, spec)

	service := clientset.Core().Services(kubeNamespace)

	serviceCreated, err := service.Create(serviceSpec)
	if err != nil {
		return nil, err
	}

	log.Println("kube-service : created")

	return serviceCreated, nil

}

//kube-service update
func (kubeDeployment Deployment) serviceUpdate(
	clientset *kubernetes.Clientset,
	spec deployment.Spec) (*v1.Service, error) {

	appName := spec.Name

	serviceSpec := getServiceSpec(clientset, spec)

	service := clientset.Core().Services(kubeNamespace)
	svc, err := kubeDeployment.getServiceByAppName(service, appName)
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

func (kubeDeployment Deployment) getServiceByAppName(service kubev1.ServiceInterface, appName string) (*v1.Service, error) {

	svc, err := service.Get(appName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return svc, nil

}

func int32Ptr(i int32) *int32 { return &i }
