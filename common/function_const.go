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

package common

//Request field
const FunctionSaveField_SPEC = "spec"
const FunctionSaveField_SOURCE = "source"

//Runtime function runtime
type Runtime string

//RuntimeJava java
const RuntimeJava = "java"

//RuntimeNodeJS nodejs
const RuntimeNodeJS = "nodejs"

//RuntimePython_2_7 python 2.7
const RuntimePython_2_7 = "python_2.7"

//RuntimePython_3_6 python 3_6
const RuntimePython_3_6 = "python_3.6"

//KubeStatusTrue True
const KubeStatusTrue = "True"

//KubeStatusFalse False
const KubeStatusFalse = "False"

//FunctionStatusRunning Running
const FunctionStatusRunning = "Running"

//FunctionStatusPending Pending
const FunctionStatusPending = "Pending"

//FunctionStatusNotFound NotFound
const FunctionStatusNotFound = "NotFound"

//FunctionStatusFailed Failed
const FunctionStatusFailed = "Failed"

//RuntimeValidate runtime validate
func RuntimeValidate(runtime Runtime) bool {

	runtimesAviable := [4]string{RuntimeJava, RuntimeNodeJS, RuntimePython_2_7, RuntimePython_3_6}

	for _, runtimeAviable := range runtimesAviable {

		if Runtime(runtimeAviable) == runtime {
			return true
		}

	}

	return false

}

//GetFunctionStatus map kube status-value to quebic status-value
func GetFunctionStatus(kubeStatus string) string {

	if kubeStatus == KubeStatusTrue {
		return FunctionStatusRunning
	}

	return FunctionStatusPending

}
