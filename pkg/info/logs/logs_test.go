package logs

import (
	"fmt"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	configinfo "github.com/xjlgod/nebula-database-diagnostic/pkg/info/config"
	"testing"
)

func TestGetAllLog(t *testing.T) {

	nodeConf := &config.NodeConfig{
		SSH: config.SSHConfig{
			Address:  "192.168.8.49",
			Port:     22,
			Timeout:  "1s",
			Username: "katz.zhang",
			Password: "nebula",
		},
		Output: config.OutputConfig{
			DirPath: "",
		},

	}
	serviceConf := &config.ServiceConfig{
		Type: config.GraphService,
		HTTPPort: 19669,
		RuntimeDir: "/home/katz.zhang/.nebula/clusters/graphd",
	}
	configinfo.GetConfigInfo(nodeConf, serviceConf)
	err := GetAllLog(nodeConf, serviceConf)
	if err != nil {
		fmt.Println(err.Error())
	}

}
