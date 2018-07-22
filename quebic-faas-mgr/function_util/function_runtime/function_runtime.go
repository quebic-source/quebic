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

package function_runtime

import (
	"quebic-faas/types"
)

//FunctionRunTime Function RunTime
type FunctionRunTime interface {

	//Runtime type eg: Java 8, Node JS, Python 2, Python 3
	RuntimeType() string

	//Prepare function.handlerFile and function.HandlerPath
	SetFunctionHandler(
		function *types.Function,
		functionSource types.FunctionSourceFile,
	) error

	//Get function Docker file content according to the runtime
	GetFunctionDockerFileContent() string

	//Get target file path of the uploaded artifact

	//Ex:
	//Java ==> local-artifact.jar --upload--> quebic-function-dir/{function-id}-function/function.jar
	//Node Package ==> local-artifact.tar --upload--> quebic-function-dir/{function-id}-function/function_handler.tar
	//Node Single File ==> local-artifact.js
	// --upload--> quebic-function-dir/{function-id}-function/handler.js
	// --create handler-tar--> function_handler.tar
	GetTargetFunctionArtifactPath(functionID string) string

	//Copy function artifact into docker build context location
	CopyFunctionIntoBuildContextLocation(
		functionID string,
		functionSource types.FunctionSourceFile,
	) error
}
