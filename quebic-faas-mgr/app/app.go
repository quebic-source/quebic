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

package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"quebic-faas/auth"
	"quebic-faas/common"
	"quebic-faas/config"
	_messenger "quebic-faas/messenger"
	"quebic-faas/quebic-faas-mgr/components"
	mgrconfig "quebic-faas/quebic-faas-mgr/config"
	dao "quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/quebic-faas-mgr/db"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/quebic-faas-mgr/deployment/kube_deployment"
	"quebic-faas/quebic-faas-mgr/function/function_util"
	_gc "quebic-faas/quebic-faas-mgr/gc"
	"quebic-faas/quebic-faas-mgr/httphandler"
	"quebic-faas/quebic-faas-mgr/logger"
	"quebic-faas/types"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"

	bolt "github.com/coreos/bbolt"
)

//App app container
type App struct {
	config     mgrconfig.AppConfig
	db         *bolt.DB
	messenger  _messenger.Messenger
	router     *mux.Router
	loggerUtil logger.Logger
	deployment dep.Deployment
	gc         _gc.GC
}

//Start init app
func (app *App) Start() {

	log.Printf("quebic-faas-manager starting")

	//get db
	db, err := db.GetDb()
	defer db.Close()
	if err != nil {
		log.Fatalf("init failed. unable to connect db, error : %v", err)
	}
	app.db = db

	//setup root admin user to loging manager
	app.setupAdminUser()

	//setup deployment
	app.setupDeployment()

	//setup gc
	app.setupGC()

	//setup manager components
	app.setupIngress()

	//setup manager components
	app.setupManagerComponents()

	//setup functions
	app.setupFunctions()

	//setup router
	app.router = mux.NewRouter()
	app.setUpHTTPHandlers()

}

//DoConfiguration app configuration
func (app *App) DoConfiguration() {
	//setup configuration
	app.config.SetDefault()
	app.setupConfiguration()
	app.setupConfigurationFromEnvVars()
}

func (app *App) setupConfiguration() {

	configFilePath := mgrconfig.GetConfigFilePath()

	readConfigJSON, err := ioutil.ReadFile(configFilePath)
	if err != nil {

		app.SaveConfiguration()

	} else {

		//if found config file set those configuration into App.config object
		savingConfig := mgrconfig.SavingConfig{}
		yaml.Unmarshal(readConfigJSON, &savingConfig)

		app.config.Auth = savingConfig.Auth
		app.config.ServerConfig = savingConfig.ServerConfig
		app.config.DockerConfig = savingConfig.DockerConfig
		app.config.KubernetesConfig = savingConfig.KubernetesConfig
		app.config.EventBusConfig = savingConfig.EventBusConfig
		app.config.APIGatewayConfig = savingConfig.APIGatewayConfig
		app.config.IngressConfig = savingConfig.IngressConfig
		app.config.MgrDashboardConfig = savingConfig.MgrDashboardConfig
		app.config.InCluster = savingConfig.InCluster
		app.config.Deployment = savingConfig.Deployment

	}

}

func (app *App) setupConfigurationFromEnvVars() {
	ingressConfigStaticIP := os.Getenv(common.EnvKey_ingressConfig_staticIP)

	if ingressConfigStaticIP != "" {
		app.config.IngressConfig.StaticIP = ingressConfigStaticIP
	}
}

//SaveConfiguration saveConfiguration in .config file
func (app *App) SaveConfiguration() {

	configFilePath := mgrconfig.GetConfigFilePath()

	//creating config dir
	os.Mkdir(config.GetConfigDirPath(), os.FileMode.Perm(0777))

	//write default configurations
	wrireConfigJSON, _ := yaml.Marshal(mgrconfig.SavingConfig{
		Auth:               app.config.Auth,
		ServerConfig:       app.config.ServerConfig,
		DockerConfig:       app.config.DockerConfig,
		EventBusConfig:     app.config.EventBusConfig,
		APIGatewayConfig:   app.config.APIGatewayConfig,
		IngressConfig:      app.config.IngressConfig,
		MgrDashboardConfig: app.config.MgrDashboardConfig,
		KubernetesConfig:   app.config.KubernetesConfig,
		InCluster:          app.config.InCluster,
		Deployment:         app.config.Deployment,
	})

	//write default configurations into config file
	err := ioutil.WriteFile(configFilePath, wrireConfigJSON, 0777)
	if err != nil {
		log.Fatalf("unable to create config file %v\n", err)
	}

}

func (app *App) setupAdminUser() {

	adminUser := types.User{
		Username:  app.config.Auth.Username,
		Password:  app.config.Auth.Password,
		Firstname: app.config.Auth.Username,
		Role:      auth.DefaultRole,
	}

	dao.AddUser(app.db, &adminUser)

}

func (app *App) setupDeployment() {

	kubeDeployment := kube_deployment.Deployment{
		Config: kube_deployment.Config{
			ConfigPath: app.config.KubernetesConfig.ConfigPath,
		},
		InCluser: app.config.InCluster,
	}

	app.deployment = kubeDeployment

	if !app.config.InCluster {
		err := app.deployment.Init()
		if err != nil {
			log.Fatalf("deployment init failed : %v", err)
		}
	}

}

func (app *App) setupGC() {
	gc := _gc.GC{}
	gc.Init(app.config, app.db, app.deployment)
	app.gc = gc
}

func (app *App) setupIngress() {

	log.Printf("ingress starting")

	deployment := app.deployment
	err := deployment.IngressCreateOrUpdate(dep.IngressSpec{
		StaticIPName: app.config.IngressConfig.StaticIP,
	})
	if err != nil {
		log.Fatalf("ingress setup failed : %v", err)
	}

	log.Printf("ingress started")

}

func (app *App) setupManagerComponents() {

	err := components.EventbusSetup(&app.config, app.deployment)
	if err != nil {
		log.Fatalf("eventbus setup failed : %v", err)
	}

	//setup messenger
	app.setupMessenger()

	//apigateway data serve
	app.setupApigatewayDataFetchListener()

	//function data serve
	app.setupFunctionDataFetchListener()

	//shutDownRequest
	app.setupShutDownRequestListener()

	//setup logger
	app.setupLogger()

	//apigateway
	err = components.ApigatewaySetup(&app.config, app.db, app.deployment, app.messenger)
	if err != nil {
		log.Fatalf("apigateway setup failed : %v", err)
	}

}

func (app *App) setupMessenger() {

	config := app.config
	messenger := _messenger.Messenger{AppID: config.AppID, EventBusConfig: config.EventBusConfig}
	err := messenger.WaitInit(time.Minute * 10)
	if err != nil {
		log.Fatalf("unable to connect eventbus %v\n", err)
	}

	app.messenger = messenger

}

func (app *App) setupApigatewayDataFetchListener() {

	messenger := app.messenger

	err := messenger.Subscribe(common.EventApigatewayDataFetch, func(event _messenger.BaseEvent) {

		apigatewayData := types.ApigatewayData{}

		//resources
		var resources []types.Resource
		err := dao.GetAll(app.db, &types.Resource{}, func(k, v []byte) error {

			resource := types.Resource{}
			json.Unmarshal(v, &resource)
			resources = append(resources, resource)
			return nil
		})
		if err != nil {
			messenger.ReplyError(event, err.Error(), 500)
			return
		}

		//manager-components
		allowComponents := [1]string{
			common.ComponentEventBus,
		}
		var components []types.ManagerComponent
		err = dao.GetAll(app.db, &types.ManagerComponent{}, func(k, v []byte) error {

			component := types.ManagerComponent{}
			json.Unmarshal(v, &component)

			//check this component is allow to serve
			for _, a := range allowComponents {
				if a == component.ID {
					components = append(components, component)
				}
			}

			return nil

		})
		if err != nil {
			messenger.ReplyError(event, err.Error(), 500)
			return
		}

		//api-gateway version
		apiGateway, err := dao.ManagerComponentGetAPIGateway(app.db)
		if err != nil {
			messenger.ReplyError(event, err.Error(), 500)
			return
		}

		//remove nil from response
		if components == nil {
			components = make([]types.ManagerComponent, 0)
		}

		if resources == nil {
			resources = make([]types.Resource, 0)
		}

		//assign data
		apigatewayData.Resources = resources
		apigatewayData.ManagerComponents = components
		apigatewayData.CurrentDeploymentVersion = apiGateway.Version

		messenger.ReplySuccess(event, apigatewayData, 200)

	}, common.ConsumerApigatewayDataFetch)
	if err != nil {
		log.Fatalf("unable to subscribe internal message listen %v\n", err)
	}

}

func (app *App) setupFunctionDataFetchListener() {

	messenger := app.messenger

	err := messenger.Subscribe(common.EventFunctionDataFetch, func(event _messenger.BaseEvent) {

		function := &types.Function{}
		err := event.ParsePayloadAsObject(function)
		if err != nil {
			messenger.ReplyError(event, "unable parse request object", 500)
			return
		}

		err = dao.GetByID(app.db, function, func(v []byte) error {

			if v == nil {
				return fmt.Errorf("function not found")
			}

			json.Unmarshal(v, function)
			return nil
		})
		if err != nil {
			log.Printf("function-data-sent faild %v", err.Error())
			messenger.ReplyError(event, err.Error(), 500)
			return
		}

		functionData := types.FunctionData{Version: function.Version}

		err = messenger.ReplySuccess(event, functionData, 200)
		if err != nil {
			log.Printf("function-data-sent reply faild %v", err.Error())
		}

	}, common.ConsumerFunctionDataFetch)

	if err != nil {
		log.Fatalf("unable to subscribe internal message listen %v\n", err)
	}

}

func (app *App) setupShutDownRequestListener() {

	messenger := app.messenger
	gc := app.gc

	err := messenger.Subscribe(common.EventShutDownRequest, func(event _messenger.BaseEvent) {

		shutDownRequest := types.ShutDownRequest{}
		event.ParsePayloadAsObject(&shutDownRequest)

		gc.SubmitJob(shutDownRequest.DeploymentID)

		messenger.ReplySuccess(event, shutDownRequest, 200)

	}, common.ConsumerShutDownRequest)
	if err != nil {
		log.Fatalf("unable to subscribe internal message listen %v", err)
	}

}

func (app *App) setupLogger() {

	loggerUtil := logger.Logger{}
	loggerUtil.Init(app.db, app.messenger)
	loggerUtil.Listen()

	app.loggerUtil = loggerUtil

}

func (app *App) setupFunctions() {

	appConfig := app.config
	db := app.db
	messenger := app.messenger
	deployment := app.deployment

	err := dao.GetAll(db, &types.Function{}, func(k, v []byte) error {

		function := &types.Function{}
		json.Unmarshal(v, function)

		_, err := function_util.FunctionDeploy(
			appConfig,
			deployment,
			messenger,
			function)
		if err != nil {
			log.Printf("unable to deploy function %s. cause : %v\n", function.Name, err)
		}

		return nil

	})
	if err != nil {
		log.Fatalf("unable to setup functions %v\n", err)
	}
}

//GetDB get app db connection
func (app *App) GetDB() *bolt.DB {
	return app.db
}

func (app *App) setUpHTTPHandlers() {

	router := app.router
	db := app.db
	messenger := app.messenger
	loggerUtil := app.loggerUtil
	deployment := app.deployment

	httphandler.SetUpHTTPHandlers(
		app.config,
		router,
		db,
		messenger,
		loggerUtil,
		deployment)

	address := app.config.ServerConfig.Host + ":" + common.IntToStr(app.config.ServerConfig.Port)

	log.Printf("quebic-faas-manager running on %s\n", address)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

	err := http.ListenAndServe(address, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router))
	if err != nil {
		log.Fatalf("quebic-faas-manager failed. error : %v", err)
	}

}

//GetAppConfig get app config
func (app *App) GetAppConfig() *mgrconfig.AppConfig {
	return &app.config
}
