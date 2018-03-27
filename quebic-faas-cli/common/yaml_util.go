package common

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//ParseYAMLToObject yaml to object
func ParseYAMLToObject(file string, object interface{}) error {

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlContent, object)
	if err != nil {
		return err
	}

	return nil
}
