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
		log.Printf("%+v\n", node)
	}
}
