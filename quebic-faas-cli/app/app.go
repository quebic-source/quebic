package app

import (
	"io/ioutil"
	"log"
	"os"
	quebicfaas "quebic-faas/config"
	"quebic-faas/quebic-faas-cli/config"
	"quebic-faas/quebic-faas-cli/service"

	yaml "gopkg.in/yaml.v2"
)

//App app container
type App struct {
	config     config.AppConfig
	mgrService service.MgrService
}

//Start init app
func (app *App) Start() {

	//setup configuration
	app.config.SetDefault()
	app.setupConfiguration()

	app.mgrService = service.MgrService{
		MgrServerConfig: app.config.MgrServerConfig,
		Auth:            app.config.Auth,
	}

}

func (app *App) setupConfiguration() {

	configFilePath := config.GetConfigFilePath()

	readConfigJSON, err := ioutil.ReadFile(configFilePath)
	if err != nil {

		app.SaveConfiguration()

	} else {

		//if found config file set those configuration into App.config object
		savingConfig := config.SavingConfig{}
		yaml.Unmarshal(readConfigJSON, &savingConfig)

		app.config.Auth = savingConfig.Auth
		app.config.MgrServerConfig = savingConfig.MgrServerConfig

	}

}

//GetAppConfig get app config
func (app *App) GetAppConfig() *config.AppConfig {
	return &app.config
}

//GetMgrService get mgrService
func (app *App) GetMgrService() service.MgrService {
	return app.mgrService
}

//SaveConfiguration saveConfiguration in .config file
func (app *App) SaveConfiguration() {

	configFilePath := config.GetConfigFilePath()

	//creating config dir
	os.Mkdir(quebicfaas.GetConfigDirPath(), os.FileMode.Perm(0777))

	//write default configurations
	wrireConfigJSON, _ := yaml.Marshal(config.SavingConfig{
		Auth:            app.config.Auth,
		MgrServerConfig: app.config.MgrServerConfig,
	})

	//write default configurations into config file
	err := ioutil.WriteFile(configFilePath, wrireConfigJSON, 0777)
	if err != nil {
		log.Fatalf("unable to create config file %v\n", err)
	}

}
