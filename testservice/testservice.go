package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"quebic-faas/quebic-faas-mgr/config"
	dep "quebic-faas/quebic-faas-mgr/deployment"
	"quebic-faas/quebic-faas-mgr/deployment/kube_deployment"
	"quebic-faas/types"
	"strings"
)

func main() {

	url := "http://192.168.1.105/manager/api/functions"
	method := "POST"

	type User struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	type Response struct {
		Message string `json:"message"`
	}

	payload := User{
		UserName: "u1",
		Password: "zxc",
	}

	jsonPayload, _ := json.Marshal(payload)

	requestBody := strings.NewReader(string(jsonPayload))

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		fmt.Print("error #1")
		panic(err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "qaz")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		fmt.Print("error #2")
		panic(err.Error())
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Print("error #3")
		panic(err.Error())
	}

	println("got response")

	resp := types.ErrorResponse{}
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		fmt.Print("error #4")
		panic(err.Error())
	}
	fmt.Printf("%v", resp)

}

func getDeployment(appConfig config.AppConfig) dep.Deployment {
	return kube_deployment.Deployment{
		Config: kube_deployment.Config{
			ConfigPath: appConfig.KubernetesConfig.ConfigPath,
		},
	}
}
