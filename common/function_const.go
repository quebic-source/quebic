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

//RuntimeValidate runtime validate
func RuntimeValidate(runtime Runtime) bool {

	runtimesAviable := [2]string{RuntimeJava, RuntimeNodeJS}

	for _, runtimeAviable := range runtimesAviable {

		if Runtime(runtimeAviable) == runtime {
			return true
		}

	}

	return false

}
