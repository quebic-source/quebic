package function_nodejs_runtime

import (
	"fmt"
	"path/filepath"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/function_util/function_common"
	"quebic-faas/types"
	"strings"
)

//BuildContextJS system create target js file name
const BuildContextJS string = "handler.js"

//BuildContextHandlerDir script dir inside docker
const BuildContextHandlerDir string = "function_handler"

//BuildContext target location tar
const buildContextHandlerTar string = "function_handler.tar"

//FunctionRunTime function runtime
type FunctionRunTime struct {
}

func (functionRunTime FunctionRunTime) RuntimeType() string {
	return common.RuntimeNodeJS
}

func (functionRunTime FunctionRunTime) SetFunctionHandler(
	function *types.Function,
	functionSourceFile types.FunctionSourceFile,
) error {

	functionArtifactFilename := functionSourceFile.FileHeader.Filename
	fileExt := filepath.Ext(functionArtifactFilename)

	// Ex : handler = index.myHandler
	//handler = {handlerFile}.{handlerPath}
	h := strings.Split(function.Handler, ".")

	if len(h) < 2 {
		return fmt.Errorf("handler is invalide. unable to found module")
	}

	handlerFile := h[0] + ".js"
	handlerPath := h[1]

	function.HandlerPath = handlerPath

	//node package
	if fileExt == ".tar" || fileExt == ".gz" {

		if handlerFile == "" {
			return fmt.Errorf("handler cannot be empty for node package")
		}

		function.HandlerFile = getDockerFunctionJSPackage(handlerFile)

		return nil

	} else if fileExt == ".js" {

		function.HandlerFile = getDockerFunctionJS()
		return nil

	}

	return fmt.Errorf("invalide artifact file type %s", fileExt)

}

func (functionRunTime FunctionRunTime) GetFunctionDockerFileContent() string {
	return common.DockerFileContent_NodeJS
}

func (functionRunTime FunctionRunTime) GetTargetFunctionArtifactPath(functionID string) string {
	return getBuildContextHandlerTar(functionID)
}

func (functionRunTime FunctionRunTime) CopyFunctionIntoBuildContextLocation(
	functionID string,
	functionSource types.FunctionSourceFile,
) error {

	functionArtifactFile := functionSource.File
	functionArtifactFilename := functionSource.FileHeader.Filename
	fileExt := filepath.Ext(functionArtifactFilename)

	targetFunctionArtifactPath := getBuildContextHandlerTar(functionID)

	//nodejs package
	if fileExt == ".tar" || fileExt == ".gz" {

		err := function_common.CopyArtifactSourceToTarget(functionArtifactFile, targetFunctionArtifactPath)
		if err != nil {
			return fmt.Errorf("unable to copy package into build-context location %v", err)
		}

	} else {

		//copy js file to function storing location
		//then create nodejs package .tar
		targetArtifactPathJS := getBuildContextJS(functionID)

		//copy original js into function dir as handler.js
		err := function_common.CopyArtifactSourceToTarget(functionArtifactFile, targetArtifactPathJS)
		if err != nil {
			return fmt.Errorf("unable to copy handler js into build-context location %v", err)
		}

		err = function_common.CreateHandlerTar(targetArtifactPathJS, targetFunctionArtifactPath)
		if err != nil {
			return fmt.Errorf("unable to create handler tar in build-context location %v", err)
		}

	}

	return nil

}

//nodejs file location inside docker
func getDockerFunctionJS() string {
	return function_common.GetDockerFunctionDir() + common.FilepathSeparator + BuildContextHandlerDir + common.FilepathSeparator + BuildContextJS
}

//nodejs package file location inside docker
func getDockerFunctionJSPackage(handlerFile string) string {
	return function_common.GetDockerFunctionDir() + common.FilepathSeparator + BuildContextHandlerDir + common.FilepathSeparator + handlerFile
}

func getBuildContextJS(functionID string) string {
	return function_common.GetFunctionDir(functionID) + common.FilepathSeparator + BuildContextJS
}

func getBuildContextHandlerTar(functionID string) string {
	return function_common.GetFunctionDir(functionID) + common.FilepathSeparator + buildContextHandlerTar
}
