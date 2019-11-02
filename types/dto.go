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

package types

import "mime/multipart"

//FunctionDTO dto
type FunctionDTO struct {
	Options    FunctionCreateOptions `json:"options" yaml:"options"`
	Function   Function              `json:"function" yaml:"function"`
	Route      Resource              `json:"route" yaml:"route"`
	SourceFile FunctionSourceFile    `json:"sourceFile" yaml:"sourceFile"`
}

//FunctionSourceFile extra options
type FunctionSourceFile struct {
	File       multipart.File        `json:"file" yaml:"file"`
	FileHeader *multipart.FileHeader `json:"fileHeader" yaml:"fileHeader"`
}

//FunctionCreateOptions extra options
type FunctionCreateOptions struct {
	Publish bool `json:"publish"`
}

//FunctionContainerLogDTO dto
type FunctionContainerLogDTO struct {
	Options  FunctionContainerLogOptions `json:"options" yaml:"options"`
	Function Function                    `json:"function" yaml:"function"`
}

//FunctionContainerLogOptions extra options
type FunctionContainerLogOptions struct {
	ReplicaIndex int    `json:"replicaIndex"`
	ShowStdout   bool   `json:"showStdout"`
	ShowStderr   bool   `json:"showStderr"`
	Since        string `json:"since"`
	Until        string `json:"until"`
	Timestamps   bool   `json:"timestamps"`
	Follow       bool   `json:"follow"`
	Tail         string `json:"tail"`
	Details      bool   `json:"details"`
}

//FunctionTest function test
type FunctionTest struct {
	Name    string                 `json:"name" yaml:"name"`
	Payload map[string]interface{} `json:"payload" yaml:"payload"`
}

//FunctionTestResponse functionTestResponse
type FunctionTestResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

//ErrorResponse errorResponse
type ErrorResponse struct {
	Status  int         `json:"status"`
	Cause   string      `json:"cause"`
	Message interface{} `json:"message"`
}

//AuthDTO ######################################
type AuthDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//JWTToken ######################################
type JWTToken struct {
	Token string `json:"token"`
}

//NewVersionMessage new version message
type NewVersionMessage struct {
	Version string `json:"version"`
}
