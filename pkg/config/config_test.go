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

	for name, node := range conf.Nodes {
		for _, service := range node.Services {
			log.Printf("%s, %+v\n", name, service)
		}
		log.Printf("%s, %+v\n", name, node)
	}

}
