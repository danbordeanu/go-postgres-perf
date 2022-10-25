package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type CSwagger struct {
	Version     string `yaml:"version"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	BasePath    string `yaml:"basepath"`
}

// LoadSwaggerConf load swagger related configs
func (c *Configuration) LoadSwaggerConf() {
	yamlFile, err := ioutil.ReadFile("swagger.yaml")
	if err != nil {
		log.Fatalf("Error opening swagger configuration file swagger.yaml: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c.Swagger)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}
}
