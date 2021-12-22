package config

import (
	"log"
	"testing"
)

func TestNewConfig(t *testing.T) {
	confPath := "./example_config/example.yaml"

	conf, err := NewConfig(confPath, "yaml")

	if err != nil {
		t.Error(err.Error())
	}

	for _, node := range conf.Nodes {
		for _, service := range node.Services {
			log.Printf("%+v\n", service)
		}
		log.Printf("%+v\n", node)
	}

}
