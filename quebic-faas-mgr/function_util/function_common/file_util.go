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

package function_common

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"quebic-faas/common"
)

//GetFunctionDir get function artifact creating location
func GetFunctionDir(functionID string) string {
	return functionsStoredLocation() + common.FilepathSeparator + functionDirPrefix + functionID
}

func GetDockerFilePath(functionID string) string {
	return GetFunctionDir(functionID) + common.FilepathSeparator + "Dockerfile"
}

// Dir inside docker container
func GetDockerFunctionDir() string {
	return common.FilepathSeparator + DockerFunctionDIR
}

func CopyArtifactSourceToTarget(
	functionArtifact multipart.File,
	targetFunctionArtifactPath string) error {

	from := functionArtifact
	defer from.Close()

	to, err := os.OpenFile(targetFunctionArtifactPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("unable to access targetFunctionArtifactPath %[1]s", err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return fmt.Errorf("unable to copy functionArtifact %[1]s", err)
	}

	log.Printf("copied function into buildContext location %[1]s", targetFunctionArtifactPath)

	return nil

}

func PrepareBuildContextLocation(functionID string) (string, error) {

	buildContextTar := getBuildContextTar(functionID)
	functionDirPath := GetFunctionDir(functionID)

	//removing previously created buildContextTar
	os.Remove(buildContextTar)

	//open function dir
	functionDir, err := os.Open(functionDirPath)
	if err != nil {
		return "", fmt.Errorf("unable to open function dir %[1]s", err)
	}
	defer functionDir.Close()

	//get all files from function dir
	files, err := functionDir.Readdir(0)
	if err != nil {
		return "", fmt.Errorf("unable to read function dir's files %[1]s", err)
	}

	// set up the buildContextTarFile file
	buildContextTarFile, err := os.Create(buildContextTar)
	if err != nil {
		return "", fmt.Errorf("unable to create buildContextTarFile %[1]s", err)
	}
	defer buildContextTarFile.Close()

	// set up the gzip writer for buildContextTarFile
	gw := gzip.NewWriter(buildContextTarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	//adding each file which stored inside functionDir
	for _, file := range files {

		if !file.IsDir() {

			filePath := functionDirPath + common.FilepathSeparator + file.Name()

			openedFile, err := os.Open(filePath)
			if err != nil {
				return "", fmt.Errorf("unable to open file %[1]s", err)
			}
			defer openedFile.Close()

			tarHeader := &tar.Header{
				Name: file.Name(),
				Mode: 0600,
				Size: file.Size(),
			}

			errorWriteHdr := tw.WriteHeader(tarHeader)
			if err != nil {
				return "", fmt.Errorf("unable to write tar header %[1]s", errorWriteHdr)
			}

			_, errorWriteTar := io.Copy(tw, openedFile)
			if err != nil {
				return "", fmt.Errorf("unable to write file into tar %[1]s", errorWriteTar)
			}

			log.Printf("added %[1]s into %[2]s", filePath, buildContextTar)

		}

	}

	return buildContextTar, nil

}

//copy original artifact into funtionDir and compress it as .tar
func CreateHandlerTar(artifactPath string, targetArtifactPath string) error {

	// set up the buildContextTarFile file
	tarFile, err := os.Create(targetArtifactPath)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	// set up the gzip writer for buildContextTarFile
	gw := gzip.NewWriter(tarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	//open artifactFile file
	artifactFile, err := os.Open(artifactPath)
	if err != nil {
		return fmt.Errorf("unable to found artifactPath %v", err)
	}
	defer artifactFile.Close()

	fileInfo, _ := artifactFile.Stat()

	tarHeader := &tar.Header{
		Name: fileInfo.Name(),
		Mode: 0600,
		Size: fileInfo.Size(),
	}

	errorWriteHdr := tw.WriteHeader(tarHeader)
	if err != nil {
		return fmt.Errorf("unable to write tar header %vs", errorWriteHdr)
	}

	_, errorWriteTar := io.Copy(tw, artifactFile)
	if err != nil {
		return fmt.Errorf("unable to write file into tar %vs", errorWriteTar)
	}

	return nil
}

func getBuildContextTar(functionID string) string {
	return GetFunctionDir(functionID) + common.FilepathSeparator + buildContextTar
}

func functionsStoredLocation() string {
	return common.GetUserHomeDir() + common.FilepathSeparator + functionsStoredDir
}
