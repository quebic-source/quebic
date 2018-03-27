/*
Copyright 2018 Tharanga Nilupul Thennakoon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
/*
	common utils
*/

package common

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

//FilepathSeparator filepathSeparator
const FilepathSeparator = string(filepath.Separator)

//GetUserHomeDir get user home dir
func GetUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

//StrToInt string to int
func StrToInt(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("ParseInt failed for %s, err %v", value, err)
		return 0
	}
	return v
}

//IntToStr int to string
func IntToStr(value int) string {
	return strconv.Itoa(value)
}

//StrParseToPrimitive convert string to matching primitve type
func StrParseToPrimitive(str string) interface{} {

	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return f
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return i
	}

	u, err := strconv.ParseUint(str, 10, 64)
	if err == nil {
		return u
	}

	b, err := strconv.ParseBool(str)
	if err == nil {
		return b
	}

	return str

}
