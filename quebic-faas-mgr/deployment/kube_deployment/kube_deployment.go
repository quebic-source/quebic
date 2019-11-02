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
	"io"
	"log"
	"os"
	"quebic-faas/quebic-faas-mgr/deployment"
	"time"

	"quebic-faas/quebic-faas-mgr/config"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	kubev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	typedextv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"quebic-faas/common"
	quebictypes "quebic-faas/types"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const kubeSelecterKeyApp = "quebic-faas-app"
const kubeSelecterKeyVersion = "quebic-faas-version"
const kubeDefaultVersion = "v1"
const kubeNamespace = "quebic-faas"
const kubeIngressName = "quebic-faas-ingress"
const kubeRoleAdmin = "quebic-faas-role-admin"
const kubeVolumeFunctionDir = "quebic-faas-function-dir"
const kubeIngressAvailableWaitTime = time.Minute * 10

//Config kubernates config
type Config struct {
	ConfigPath string `json:"configPath"`
}

//Deployment kube deployment
type Deployment struct {
	Config   Config
	InCluser bool //is kube client is inside the cluster
}

//DeploymentType implementation for deployment.DeploymentType()
func (kubeDeployment Deployment) DeploymentType() string {
	return config.Deployment_Kubernetes
}

//Init implementation for deployment.Init()
func (kubeDeployment Deployment) Init() error {
	err := kubeDeployment.createAdminRole()
	if err != nil {
		return err
	}

	err = kubeDeployment.createNameSpace()
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

//CreateOrUpdateDeployment implementation for deployment.CreateOrUpdateDeployment()
func (kubeDeployment Deployment) CreateOrUpdateDeployment(spec deployment.Spec) error {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return err
	}

	err = kubeDeployment.deploymentUpdate(clientset, spec)
	if err != nil {

		if errors.IsNotFound(err) {

			err = kubeDeployment.deploymentCreate(clientset, spec)
			if err != nil {
				return err
			}

		} else {
			return err
		}

	}

	return nil
}

//CreateOrUpdateService implementation for deployment.CreateOrUpdateService()
func (kubeDeployment Deployment) CreateOrUpdateService(deploySpec deployment.Spec) (deployment.Details, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return deployment.Details{}, err
	}

	var service *v1.Service

	//service
	service, err = kubeDeployment.serviceUpdate(clientset, deploySpec)
	if err != nil {
		if errors.IsNotFound(err) {

			service, err = kubeDeployment.serviceCreate(clientset, deploySpec)
			if err != nil {
				return deployment.Details{}, err
			}

		} else {
			return deployment.Details{}, err
		}

	}

	return deployment.Details{
		Host: service.Spec.ClusterIP,
	}, nil
}

//CreateService implementation for deployment.CreateService()
func (kubeDeployment Deployment) CreateService(deploySpec deployment.Spec) (deployment.Details, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return deployment.Details{}, err
	}

	service, err := kubeDeployment.serviceCreate(clientset, deploySpec)
	if err != nil {

		if errors.IsAlreadyExists(err) {
			return kubeDeployment.GetService(deploySpec.Name)
		}
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

//DeleteDeployment implementation for deployment.DeleteDeployment()
func (kubeDeployment Deployment) DeleteDeployment(name string) error {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return err
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	deletePolicy := metav1.DeletePropagationForeground
	return deploymentsClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

//ListAll implementation for deployment.ListAll()
func (kubeDeployment Deployment) ListAll(filters deployment.ListFilters) ([]deployment.Details, error) {
	return nil, nil
}

//ListByName implementation for deployment.ListByName()
func (kubeDeployment Deployment) ListByName(name string) (deployment.Details, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return deployment.Details{}, err
	}

	details := deployment.Details{}

	err = kubeDeployment.getService(name, &details, clientset)
	if err != nil {
		return deployment.Details{}, err
	}

	err = kubeDeployment.getDeployment(name, &details, clientset)
	if err != nil {
		return deployment.Details{}, err
	}

	//append details
	return details, nil

}

//GetStatus implementation for deployment.GetStatus()
func (kubeDeployment Deployment) GetStatus(name string) (string, error) {

	details, err := kubeDeployment.GetDeployment(name)
	if err != nil {
		return "", err
	}

	return details.Status, nil

}

//GetDeployment implementation for deployment.GetDeployment()
func (kubeDeployment Deployment) GetDeployment(name string) (deployment.Details, error) {

	//deployment details
	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return deployment.Details{}, err
	}

	details := deployment.Details{}

	err = kubeDeployment.getDeployment(name, &details, clientset)
	if err != nil {
		return deployment.Details{}, err
	}

	return details, nil

}

//GetService implementation for deployment.GetService()
func (kubeDeployment Deployment) GetService(name string) (deployment.Details, error) {

	//deployment details
	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return deployment.Details{}, err
	}

	details := deployment.Details{}
	err = kubeDeployment.getService(name, &details, clientset)
	if err != nil {
		return deployment.Details{}, err
	}

	return details, nil

}

func (kubeDeployment Deployment) getDeployment(name string, details *deployment.Details, clientset *kubernetes.Clientset) error {

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	deploymentDetails, err := deploymentsClient.Get(name, metav1.GetOptions{})
	if err != nil {
		return err
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

	details.Status = string(conditionStatus)
	details.Replicas = deployment.Replicas(int(*deploymentDetails.Spec.Replicas))

	return nil
}

func (kubeDeployment Deployment) getService(name string, details *deployment.Details, clientset *kubernetes.Clientset) error {
	//service details
	service := clientset.Core().Services(kubeNamespace)
	serviceDetails, err := service.Get(name, metav1.GetOptions{})
	if err != nil {
		return err
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

	details.Host = host
	details.PortConfigs = ports

	return nil
}

//LogsByName implementation for deployment.LogsByName()
func (kubeDeployment Deployment) LogsByName(name string, options quebictypes.FunctionContainerLogOptions) error {

	replicaIndex := options.ReplicaIndex

	containers, err := kubeDeployment.ListContainersByName(name)
	if err != nil {
		return err
	}

	containerID := containers[replicaIndex].ID

	err = kubeDeployment.LogsByContainerID(containerID, options)
	if err != nil {
		return err
	}

	return nil

}

//ListContainersByName implementation for deployment.ListContainersByName()
func (kubeDeployment Deployment) ListContainersByName(name string) ([]deployment.Container, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(kubeNamespace)

	deploymentDetails, err := deploymentsClient.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	deploymentLabels := labels.Set(deploymentDetails.Labels).AsSelector().String()
	podSpecs, err := clientset.Core().Pods(kubeNamespace).List(metav1.ListOptions{
		LabelSelector: deploymentLabels,
	})
	if err != nil {
		return nil, err
	}

	var containers []deployment.Container
	for _, p := range podSpecs.Items {
		containers = append(containers, deployment.Container{ID: p.Name})
	}

	return containers, nil

}

//LogsByContainerID implementation for deployment.LogsByContainerID()
func (kubeDeployment Deployment) LogsByContainerID(id string, options quebictypes.FunctionContainerLogOptions) error {

	logOptions := v1.PodLogOptions{}

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return err
	}

	req := clientset.Core().Pods(kubeNamespace).GetLogs(id, &logOptions)

	readCloser, err := req.Stream()
	if err != nil {
		return err
	}

	defer readCloser.Close()
	_, err = io.Copy(os.Stdout, readCloser)
	if err != nil {
		return err
	}

	return nil

}

func (kubeDeployment Deployment) createNameSpace() error {

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

func (kubeDeployment Deployment) createAdminRole() error {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return err
	}

	roleSpec := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: kubeRoleAdmin,
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: kubeNamespace,
			},
		},
	}
	_, err = clientset.RbacV1().ClusterRoleBindings().Create(roleSpec)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return fmt.Errorf("kube-role-create failed : %v", err)
		}

		return nil
	}

	log.Printf("kube-role-created : %s", kubeRoleAdmin)

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

	var config *rest.Config

	if kubeDeployment.InCluser {
		c, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("kube-in-cluster-config getting failed %v", err)
		}
		config = c
	} else {
		kubeConfig := kubeDeployment.Config
		c, err := clientcmd.BuildConfigFromFlags("", kubeConfig.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("kube-config config getting failed %v", err)
		}
		config = c
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
	version := spec.Version

	if "" == spec.DeploymentName {
		spec.DeploymentName = appName
	}

	if "" == version {
		version = kubeDefaultVersion
	}

	dockerimage := spec.Dockerimage
	replicas := int32Ptr(int32(spec.Replicas))

	var imagePullPolicy = spec.ImagePullPolicy
	if imagePullPolicy == "" {
		imagePullPolicy = string(apiv1.PullAlways)
	}

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

	//volumes
	var volumeMounts []v1.VolumeMount
	var volumes []v1.Volume
	for _, volume := range spec.Volumes {

		//volumeMounts
		containerPath := volume.ContainerPath
		volumeMounts = append(
			volumeMounts,
			v1.VolumeMount{
				Name:      kubeVolumeFunctionDir,
				MountPath: containerPath,
			})

		//volumes
		hostPath := volume.HostPath
		hostPathDirectoryType := volume.HostPathType
		volumeSource := v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: hostPath,
				Type: &hostPathDirectoryType,
			},
		}

		volumes = append(
			volumes, v1.Volume{
				Name:         kubeVolumeFunctionDir,
				VolumeSource: volumeSource,
			})

	}

	terminationGracePeriodSeconds := int64(0)

	return &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.DeploymentName,
			Labels: map[string]string{
				kubeSelecterKeyApp:     appName,
				kubeSelecterKeyVersion: version,
			},
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: replicas,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						kubeSelecterKeyApp:     appName,
						kubeSelecterKeyVersion: version,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            appName,
							Image:           dockerimage,
							Ports:           ports,
							Env:             envVar,
							VolumeMounts:    volumeMounts,
							Command:         spec.Command,
							ImagePullPolicy: apiv1.PullPolicy(imagePullPolicy),
						},
					},
					Volumes:                       volumes,
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
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
	version := spec.Version

	if "" == version {
		version = kubeDefaultVersion
	}

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
			Type: v1.ServiceTypeNodePort,
			Selector: map[string]string{
				kubeSelecterKeyApp: appName,
				//kubeSelecterKeyVersion: version,
			},
			Ports: ports,
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

//IngressDescribe ingress get details
func (kubeDeployment Deployment) IngressDescribe(waitForAvailable bool) (deployment.IngressDetails, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return deployment.IngressDetails{}, err
	}

	ingressesClient := clientset.ExtensionsV1beta1().Ingresses(kubeNamespace)

	ingress, err := ingressesClient.Get(kubeIngressName, metav1.GetOptions{})
	if err != nil {
		log.Printf("kube-ingress inspect-failed : %v", err)
		return deployment.IngressDetails{}, err
	}

	if waitForAvailable {
		err = waitingForIngressAvailable(ingressesClient)
		if err != nil {
			return deployment.IngressDetails{}, err
		}
	} else {
		chk := checkIngressLoadbalancerIP(ingress)
		if !chk {
			return deployment.IngressDetails{}, fmt.Errorf("kube-ingress not ready yet")
		}
	}

	lb := ingress.Status.LoadBalancer.Ingress[0]
	hostname := lb.Hostname
	ip := lb.IP

	return deployment.IngressDetails{
		Hostname: hostname,
		IP:       ip,
	}, nil

}

//IngressCreateOrUpdate ingress update or create
func (kubeDeployment Deployment) IngressCreateOrUpdate(spec deployment.IngressSpec) error {

	_, err := kubeDeployment.ingressUpdate(spec)
	if err != nil {

		if errors.IsNotFound(err) {

			_, err = kubeDeployment.ingressCreate(spec)
			if err != nil {
				log.Printf("kube-ingress create-failed : %v", err)
				return err
			}

		} else {
			log.Printf("kube-ingress update-failed : %v", err)
			return err
		}

	}

	return nil

}

func (kubeDeployment Deployment) ingressInspect() (*extv1beta1.Ingress, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	ingressesClient := clientset.ExtensionsV1beta1().Ingresses(kubeNamespace)

	return ingressesClient.Get(kubeIngressName, metav1.GetOptions{})

}

func (kubeDeployment Deployment) ingressCreate(spec deployment.IngressSpec) (*extv1beta1.Ingress, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	ingressesClient := clientset.ExtensionsV1beta1().Ingresses(kubeNamespace)

	ingressSpec := kubeDeployment.getIngressSpec(spec)

	ingress, err := ingressesClient.Create(ingressSpec)
	if err != nil {
		return nil, err
	}

	log.Println("kube-ingress : created")

	return ingress, nil

}

func (kubeDeployment Deployment) ingressUpdate(spec deployment.IngressSpec) (*extv1beta1.Ingress, error) {

	clientset, err := kubeDeployment.getClient()
	if err != nil {
		return nil, err
	}

	ingressesClient := clientset.ExtensionsV1beta1().Ingresses(kubeNamespace)

	ingressSpec := kubeDeployment.getIngressSpec(spec)

	ingress, err := ingressesClient.Update(ingressSpec)
	if err != nil {
		return nil, err
	}

	log.Println("kube-ingress : updated")

	return ingress, nil

}

func (kubeDeployment Deployment) getIngressSpec(spec deployment.IngressSpec) *extv1beta1.Ingress {

	annotations := make(map[string]string)

	annotations[common.IngressAnnotationRewriteTarget] = common.IngressAnnotationRewriteTargetVal

	if spec.StaticIPName != "" {
		annotations[common.IngressAnnotationStaticIP] = spec.StaticIPName
	}

	managerAPIIngressRule := getIngressRule(common.IngressHostManager, common.ComponentMgrAPI, common.MgrServerPort, common.IngressRoutePrefixManagerAPI)
	apiGatewayIngressRule := getIngressRule(common.IngressHostAPIGateway, common.ComponentAPIGateway, common.ApigatewayServerPort, common.IngressRoutePrefixAPIGateway)

	ingressSpec := &extv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        kubeIngressName,
			Annotations: annotations,
		},
		Spec: extv1beta1.IngressSpec{
			Rules: []extv1beta1.IngressRule{
				managerAPIIngressRule,
				apiGatewayIngressRule,
			},
		},
	}

	return ingressSpec

}

func getIngressRule(host string, service string, port int, routePrefix string) extv1beta1.IngressRule {

	ingressBackend := extv1beta1.IngressBackend{
		ServiceName: service,
		ServicePort: intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: int32(port),
		},
	}

	httpIngressPath := extv1beta1.HTTPIngressPath{
		Backend: ingressBackend,
	}

	return extv1beta1.IngressRule{
		Host: host,
		IngressRuleValue: extv1beta1.IngressRuleValue{
			HTTP: &extv1beta1.HTTPIngressRuleValue{
				Paths: []extv1beta1.HTTPIngressPath{httpIngressPath},
			},
		},
	}

}

func waitingForIngressAvailable(ingressClient typedextv1beta1.IngressInterface) error {

	waitForResponse := make(chan bool)

	go func() {

		log.Printf("waiting for ingress available...")

		for {

			ingress, err := ingressClient.Get(kubeIngressName, metav1.GetOptions{})
			if err != nil {
				waitForResponse <- false
				break
			}

			chkLbIP := checkIngressLoadbalancerIP(ingress)
			if chkLbIP {
				waitForResponse <- true
				break
			}

			time.Sleep(time.Second * 2)

		}
	}()

	select {
	case status := <-waitForResponse:

		if !status {
			return fmt.Errorf("kube-ingress inspect failed")
		}

		return nil

	case <-time.After(kubeIngressAvailableWaitTime):
		return fmt.Errorf("ingress take longtime to available. please try again. reason : timeout")
	}

}

func checkIngressLoadbalancerIP(ingress *extv1beta1.Ingress) bool {

	lbs := ingress.Status.LoadBalancer.Ingress
	if len(lbs) > 0 {
		return true
	}

	return false
}

func int32Ptr(i int32) *int32 { return &i }
