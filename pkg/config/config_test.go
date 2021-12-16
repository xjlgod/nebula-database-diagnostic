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

	log.Printf("%+v\n", conf)
}
