package remote

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"log"
	"testing"
	"time"
)

func TestGetFilesInRemoteDir(t *testing.T) {
	path := "/home/katz.zhang/logs"
	localDir := "logs"
	conf := config.SSHConfig{
		Address:  "192.168.8.49",
		Port:     22,
		Timeout:  "1s",
		Username: "katz.zhang",
		Password: "nebula",
	}
	err := GetFilesInRemoteDir("nihaoa", conf, path, localDir)
	if err != nil {
		log.Fatal(err.Error())
	}
	//go testSSHClient1(conf)
	//testSSHClient2(conf)
	time.Sleep(3 * time.Second)
}
