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

import "quebic-faas/common"

//BuildContextJar jar file
const BuildContextJar string = "function.jar"

//BuildContextJS system create target js file name
const BuildContextJS string = "handler.js"

//BuildContextHandlerDir script dir inside docker
const BuildContextHandlerDir string = "function_handler"

func getDockerFunctionDir() string {
	return common.FilepathSeparator + DockerFunctionDIR
}

// ############################ JAVA ###################################

//GetDockerFunctionJAR jar file location inside docker
func GetDockerFunctionJAR() string {
	return getDockerFunctionDir() + common.FilepathSeparator + BuildContextJar
}

func getBuildContextJar(functionID string) string {
	return GetFunctionDir(functionID) + common.FilepathSeparator + BuildContextJar
}

// ######################################################################

// ############################ NodeJS ###################################

//GetDockerFunctionJS node js file location inside docker
func GetDockerFunctionJS() string {
	return getDockerFunctionDir() + common.FilepathSeparator + BuildContextHandlerDir + common.FilepathSeparator + BuildContextJS
}

// ############################ NodeJS Package ###################################

//GetDockerFunctionJSPackage node js file location inside docker
func GetDockerFunctionJSPackage(handlerFile string) string {
	return getDockerFunctionDir() + common.FilepathSeparator + BuildContextHandlerDir + common.FilepathSeparator + handlerFile
}

func getBuildContextJS(functionID string) string {
	return GetFunctionDir(functionID) + common.FilepathSeparator + BuildContextJS
}

// ######################################################################
