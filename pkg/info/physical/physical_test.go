package physical

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"log"
	"testing"
)

func TestGetInfo(t *testing.T) {
	conf := config.SSHConfig{
		Address:  "192.168.8.49",
		Port:     22,
		Timeout:  "1s",
		Username: "katz.zhang",
		Password: "nebula",
	}
	info, _ := GetPhyInfo(conf)
	log.Printf("%+v", info)
}
