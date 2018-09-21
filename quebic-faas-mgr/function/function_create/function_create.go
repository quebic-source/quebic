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

package function_create

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"quebic-faas/quebic-faas-mgr/function/function_common"
	"quebic-faas/quebic-faas-mgr/function/function_runtime"
	quebicFaasTypes "quebic-faas/types"
)

const functionsStoredDir string = ".quebic-faas-functions"
const functionDirPrefix string = "function-"
const buildContextTar string = "function.tar"

const buildContextHandlerTar string = "function_handler.tar"

//DockerFunctionDIR dir which function is locted
const DockerFunctionDIR string = "app"

//CreateFunction create function
func CreateFunction(
	functionID string,
	functionSource quebicFaasTypes.FunctionSourceFile,
	functionRunTime function_runtime.FunctionRunTime) (string, error) {

	err := createFunctionDir(functionID)
	if err != nil {
		return "", err
	}

	err = createFunctionDockerFile(functionID, functionRunTime)
	if err != nil {
		return "", err
	}

	err = functionRunTime.CopyFunctionIntoBuildContextLocation(functionID, functionSource)
	if err != nil {
		return "", err
	}

	return function_common.PrepareBuildContextLocation(functionID)

}

func createFunctionDir(functionID string) error {

	functionDir := function_common.GetFunctionDir(functionID)

	err := os.MkdirAll(functionDir, os.FileMode.Perm(0777))
	if err != nil {
		return fmt.Errorf("function dir creation failed %[1]s", err)
	}

	log.Printf("created function dir %[1]s", functionDir)

	return nil

}

func createFunctionDockerFile(
	functionID string,
	functionRunTime function_runtime.FunctionRunTime) error {

	functionDockerfilePath := function_common.GetDockerFilePath(functionID)

	dockerFileContent := functionRunTime.GetFunctionDockerFileContent()
	if dockerFileContent == "" {
		return fmt.Errorf("function Dockerfile read failed. Empty docker file")
	}

	err := ioutil.WriteFile(functionDockerfilePath, []byte(dockerFileContent), os.FileMode.Perm(0777))
	if err != nil {
		return fmt.Errorf("function Dockerfile creation failed %[1]s", err)
	}

	log.Printf("created Dockerfile %[1]s", functionDockerfilePath)

	return nil

}
