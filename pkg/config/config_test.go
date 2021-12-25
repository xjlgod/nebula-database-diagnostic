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

	for name, info := range conf.Infos {
		for _, service := range info.Services {
			log.Printf("%s, %+v\n", name, service)
		}
		log.Printf("%s, %+v\n", name, info)
	}
	log.Printf("%+v\n", conf.Diag)
}
