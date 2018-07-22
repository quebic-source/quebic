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

/*
	quebic-faas apigateway.
*/

package main

import (
	"log"
	"net/http"
	"os"
	"quebic-faas/common"
	"quebic-faas/messenger"
	_messenger "quebic-faas/messenger"
	"quebic-faas/quebic-faas-apigateway/config"
	"quebic-faas/quebic-faas-apigateway/httphandler"
	"quebic-faas/types"
	"time"

	"github.com/gorilla/mux"
)

const appStatusKeyEventbusConnect = "eventbus-connect-failed"
const appStatusKeyApigatewatDataFetch = "apigateway-data-fetch-failed"

//App app container
type App struct {
	config        config.AppConfig
	resources     []types.Resource
	router        *mux.Router
	messenger     _messenger.Messenger
	appStatusList map[string]string
}

//Start init app
func (app *App) Start() {

	//setup configuration
	app.setupConfiguration()

	//setup messenger
	app.setupMessenger()

	//Load all other configurations from manager rest endpoint.
	app.loadAPIGatewayData()

	//setup router
	app.router = mux.NewRouter()
	app.setUpHTTPHandlers()

}

func (app *App) setupConfiguration() {

	app.config.SetDefault()

	app.setupEnvVariables()

	app.appStatusList = make(map[string]string)

}

func (app *App) setupEnvVariables() {

	app.config.Auth.Accesstoken = os.Getenv(common.EnvKeyAPIGateWayAccessKey)

	rabbitmq_host := os.Getenv(common.EnvKey_rabbitmq_host)
	rabbitmq_port := os.Getenv(common.EnvKey_rabbitmq_port)
	rabbitmq_management_username := os.Getenv(common.EnvKey_rabbitmq_management_username)
	rabbitmq_management_password := os.Getenv(common.EnvKey_rabbitmq_management_password)

	if rabbitmq_host != "" {
		app.config.EventBusConfig.AMQPHost = rabbitmq_host
		log.Printf("## load from env AMQPHost : %s ##", app.config.EventBusConfig.AMQPHost)
	}

	if rabbitmq_port != "" {
		app.config.EventBusConfig.AMQPPort = common.StrToInt(rabbitmq_port)
		log.Printf("## load from env AMQPPort : %d ##", app.config.EventBusConfig.AMQPPort)
	}

	if rabbitmq_management_username != "" {
		app.config.EventBusConfig.ManagementUserName = rabbitmq_management_username
		log.Printf("## load from env ManagementUserName : %s ##", app.config.EventBusConfig.ManagementUserName)
	}

	if rabbitmq_management_password != "" {
		app.config.EventBusConfig.ManagementPassword = rabbitmq_management_password
		log.Printf("## load from env ManagementPassword : %s ##", app.config.EventBusConfig.ManagementPassword)
	}

}

func (app *App) setUpHTTPHandlers() {

	router := app.router
	messenger := app.messenger
	appStatusList := app.appStatusList

	httphandler.SetUpHTTPHandlers(
		app.config,
		app.resources,
		router,
		messenger,
		appStatusList)

	address := app.config.ServerConfig.Host + ":" + common.IntToStr(app.config.ServerConfig.Port)

	log.Printf("quebic-faas apigateway start %s\n", address)

	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Panicf("quebic-faas apigateway start failed. error : %v", err)
	}

}

func (app *App) setupMessenger() {

	config := app.config
	messenger := _messenger.Messenger{AppID: config.AppID, EventBusConfig: config.EventBusConfig}
	err := messenger.Init()
	if err != nil {
		app.addStatus(appStatusKeyEventbusConnect, err.Error())
	}

	app.messenger = messenger

}

func (app *App) loadAPIGatewayData() {

	managerAccessKey := app.config.Auth.Accesstoken
	requestHeaders := make(map[string]string)
	requestHeaders[common.HeaderAccessKey] = managerAccessKey

	_, err := app.messenger.PublishBlocking(
		common.EventApigatewayDataFetch,
		"",
		requestHeaders,
		func(message messenger.BaseEvent, status int, context messenger.Context) {

			apigatewayData := &types.ApigatewayData{}
			message.ParsePayloadAsObject(apigatewayData)
			app.resources = apigatewayData.Resources

			log.Printf("apigateway-data fetched")

		},
		func(err string, statuscode int, context messenger.Context) {

			log.Printf("apigateway-data fetched failed : %v", err)
			app.addStatus(appStatusKeyApigatewatDataFetch, err)

		},
		time.Second*5,
	)
	if err != nil {
		app.addStatus(appStatusKeyApigatewatDataFetch, err.Error())
		return
	}

}

func (app *App) addStatus(key string, message string) {
	app.appStatusList[key] = message
}

func main() {

	app := App{}
	app.Start()

}
