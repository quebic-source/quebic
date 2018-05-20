package common

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//ParseYAMLFileToObject yaml file to object
func ParseYAMLFileToObject(file string, object interface{}) error {

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

//ParseObjectToYAML object to yaml
func ParseObjectToYAML(object interface{}) (string, error) {

	bytes, err := yaml.Marshal(object)
	if err != nil {
		return "", err
	}

	return string(bytes), nil

}
