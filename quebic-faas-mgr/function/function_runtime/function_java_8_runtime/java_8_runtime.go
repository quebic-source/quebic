package function_java_8_runtime

import (
	"fmt"
	"path/filepath"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/function/function_common"
	"quebic-faas/types"
)

//BuildContextJar jar file
const BuildContextJar string = "function.jar"

//FunctionRunTime function runtime
type FunctionRunTime struct {
}

func (functionRunTime FunctionRunTime) RuntimeType() string {
	return common.RuntimeJava
}

func (functionRunTime FunctionRunTime) SetFunctionHandler(
	function *types.Function,
	functionSourceFile types.FunctionSourceFile,
) error {

	functionArtifactFilename := functionSourceFile.FileHeader.Filename
	fileExt := filepath.Ext(functionArtifactFilename)

	if fileExt == ".jar" {
		function.HandlerFile = getDockerFunctionJAR()
		function.HandlerPath = function.Handler
		return nil
	}

	return fmt.Errorf("invalide artifact file type %s", fileExt)

}

func (functionRunTime FunctionRunTime) GetFunctionDockerFileContent() string {
	return common.DockerFileContent_Java
}

func (functionRunTime FunctionRunTime) GetTargetFunctionArtifactPath(functionID string) string {
	return getBuildContextJar(functionID)
}

func (functionRunTime FunctionRunTime) CopyFunctionIntoBuildContextLocation(
	functionID string,
	functionSource types.FunctionSourceFile,
) error {

	functionArtifactFile := functionSource.File

	targetFunctionArtifactPath := functionRunTime.GetTargetFunctionArtifactPath(functionID)
	return function_common.CopyArtifactSourceToTarget(functionArtifactFile, targetFunctionArtifactPath)

}

//GetDockerFunctionJAR jar file location inside docker
func getDockerFunctionJAR() string {
	return function_common.GetDockerFunctionDir() + common.FilepathSeparator + BuildContextJar
}

func getBuildContextJar(functionID string) string {
	return function_common.GetFunctionDir(functionID) + common.FilepathSeparator + BuildContextJar
}
