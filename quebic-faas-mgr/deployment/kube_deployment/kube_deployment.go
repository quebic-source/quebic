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
	"time"

	"quebic-faas/quebic-faas-mgr/config"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	kubev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	typedextv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"

	"quebic-faas/common"
	quebictypes "quebic-faas/types"
)

const kubeSelecterKey = "quebic-faas-app"
const kubeNamespace = "quebic-faas"
const kubeIngressName = "quebic-faas-ingress"
const kubeIngressAvailableWaitTime = time.Minute * 10

//Config kubernates config
type Config struct {
	ConfigPath string `json:"configPath"`
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

	//deployment details
	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return deployment.Details{}, err
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	deploymentDetails, err := deploymentsClient.Get(name, metav1.GetOptions{})
	if err != nil {
		return deployment.Details{}, err
	}

	conditionStatus := common.KubeStatusFalse
	conditions := deploymentDetails.Status.Conditions
	if len(conditions) > 0 {
		for _, spec := range conditions {
			if spec.Type == "Available" {
				conditionStatus = string(spec.Status)
				break
			}
		}
	}

	//pod details
	/*deploymentLabels := labels.Set(deploymentDetails.Labels).AsSelector().String()
	podSpecs, err := clientset.Core().Pods(kubeNamespace).List(metav1.ListOptions{
		LabelSelector: deploymentLabels,
	})
	if err != nil {
		return deployment.Details{}, err
	}
	var pods []deployment.Pod
	for _, pod := range podSpecs.Items {
		pods = append(pods, deployment.Pod{Name: pod.Name})
	}*/

	//service details
	service := clientset.Core().Services(kubeNamespace)
	serviceDetails, err := service.Get(name, metav1.GetOptions{})
	if err != nil {
		return deployment.Details{}, err
	}

	host := serviceDetails.Spec.ClusterIP

	//ports
	var ports []deployment.PortConfig
	for _, specPorts := range serviceDetails.Spec.Ports {
		ports = append(ports, deployment.PortConfig{
			Name:     specPorts.Name,
			Protocol: deployment.PortProtocol(specPorts.Protocol),
			Port:     deployment.Port(specPorts.Port),
		})
	}

	//append details
	return deployment.Details{
		Status:      string(conditionStatus),
		Host:        host,
		PortConfigs: ports,
	}, nil

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

	service := clientset.Core().Services(kubeNamespace)
	err = service.Delete(appName, &metav1.DeleteOptions{})
	if err != nil {
		log.Printf("kube-service delete failed : %v", err)
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)
	err = deploymentsClient.Delete(appName, &metav1.DeleteOptions{})
	if err != nil {
		log.Printf("kube-deployment delete failed : %v", err)
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

//IngressInspect ingress get details
func IngressInspect(kubeDeployment Deployment) (*extv1beta1.Ingress, error) {

	ingress, err := ingressInspect(kubeDeployment)
	if err != nil {
		log.Printf("kube-ingress-inspect-failed : %v", err)
		return nil, err
	}

	return ingress, nil

}

//IngressLoadbalancerIP ingress get lb ip
func IngressLoadbalancerIP(kubeDeployment Deployment) (string, error) {

	ingress, err := ingressInspect(kubeDeployment)
	if err != nil {
		log.Printf("kube-ingress-inspect-failed : %v", err)
		return "", err
	}

	return getIngressLoadbalancerIP(ingress), nil

}

//IngressCreateOrUpdate ingress update or create
func IngressCreateOrUpdate(
	kubeDeployment Deployment,
	apiGatewayConfig config.APIGatewayConfig) (*extv1beta1.Ingress, error) {

	ingress, err := ingressUpdate(kubeDeployment, apiGatewayConfig)
	if err != nil {

		if errors.IsNotFound(err) {

			ingress, err = ingressCreate(kubeDeployment, apiGatewayConfig)
			if err != nil {
				log.Printf("kube-ingress-create-failed : %v", err)
				return nil, err
			}

		} else {
			log.Printf("kube-ingress-update-failed : %v", err)
			return nil, err
		}

	}

	return ingress, nil

}

func ingressInspect(kubeDeployment Deployment) (*extv1beta1.Ingress, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	ingressesClient := clientset.ExtensionsV1beta1().Ingresses(kubeNamespace)

	return ingressesClient.Get(kubeIngressName, metav1.GetOptions{})

}

func ingressCreate(
	kubeDeployment Deployment,
	apiGatewayConfig config.APIGatewayConfig) (*extv1beta1.Ingress, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	ingressesClient := clientset.ExtensionsV1beta1().Ingresses(kubeNamespace)

	ingressSpec := getIngressSpec(apiGatewayConfig)

	ingress, err := ingressesClient.Create(ingressSpec)
	if err != nil {
		return nil, err
	}

	err = waitingForIngressAvailable(ingressesClient)
	if err != nil {
		return nil, err
	}

	log.Println("kube-ingress : created")

	return ingress, nil

}

func ingressUpdate(
	kubeDeployment Deployment,
	apiGatewayConfig config.APIGatewayConfig) (*extv1beta1.Ingress, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	ingressesClient := clientset.ExtensionsV1beta1().Ingresses(kubeNamespace)

	ingressSpec := getIngressSpec(apiGatewayConfig)

	ingress, err := ingressesClient.Update(ingressSpec)
	if err != nil {
		return nil, err
	}

	log.Println("kube-ingress : updated")

	return ingress, nil

}

func getIngressSpec(apiGatewayConfig config.APIGatewayConfig) *extv1beta1.Ingress {

	apiGatewayPort := int32(apiGatewayConfig.ServerConfig.Port)
	ingressConfig := apiGatewayConfig.IngressConfig

	annotations := make(map[string]string)

	if ingressConfig.Provider == "" {
		annotations[common.IngressClass] = common.IngressNginx
	} else {
		annotations[ingressConfig.Provider] = ingressConfig.StaticIP
	}

	ingressBackend := extv1beta1.IngressBackend{
		ServiceName: common.ComponentAPIGateway,
		ServicePort: intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: apiGatewayPort,
		},
	}

	httpIngressPath := extv1beta1.HTTPIngressPath{
		Backend: ingressBackend,
		Path:    ingressConfig.RoutePrefix,
	}

	ingressRule := extv1beta1.IngressRule{
		Host: "quebic-faas-api.io",
		IngressRuleValue: extv1beta1.IngressRuleValue{
			HTTP: &extv1beta1.HTTPIngressRuleValue{
				Paths: []extv1beta1.HTTPIngressPath{httpIngressPath},
			},
		},
	}

	ingressSpec := &extv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        kubeIngressName,
			Annotations: annotations,
		},
		Spec: extv1beta1.IngressSpec{
			Rules: []extv1beta1.IngressRule{ingressRule},
		},
	}

	return ingressSpec

}

func waitingForIngressAvailable(ingressClient typedextv1beta1.IngressInterface) error {

	waitForResponse := make(chan bool)

	go func() {

		log.Printf("waiting for ingress available")

		for {

			ingress, err := ingressClient.Get(kubeIngressName, metav1.GetOptions{})
			if err != nil {
				waitForResponse <- false
				break
			}

			lbIP := getIngressLoadbalancerIP(ingress)
			if lbIP != "" {
				waitForResponse <- true
				break
			}

			time.Sleep(time.Second * 2)

			fmt.Print(".")

		}
		fmt.Println("")
	}()

	select {
	case status := <-waitForResponse:

		if !status {
			return fmt.Errorf("kube-ingress-inspect failed")
		}

		return nil

	case <-time.After(kubeIngressAvailableWaitTime):
		return fmt.Errorf("ingress take longtime to available. please try again. reason : timeout")
	}

}

func getIngressLoadbalancerIP(ingress *extv1beta1.Ingress) string {

	lbs := ingress.Status.LoadBalancer.Ingress
	if len(lbs) > 0 {
		lb := lbs[0]
		return lb.IP
	}

	return ""
}

func int32Ptr(i int32) *int32 { return &i }
