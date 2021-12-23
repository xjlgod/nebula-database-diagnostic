package main

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"log"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	nodeConf := &config.NodeConfig{
		SSH: config.SSHConfig{
			Address:  "192.168.8.49",
			Port:     22,
			Timeout:  "1s",
			Username: "katz.zhang",
			Password: "nebula",
		},
	}
	serviceConf := &config.ServiceConfig{
		Type: config.GraphService,
		HTTPPort: 19559,
	}
	info, _ := GetMetricsInfo(nodeConf, serviceConf)
	log.Printf("%+v", info)
}