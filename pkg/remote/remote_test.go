package remote

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"log"
	"testing"
	"time"
)

func TestNewSSHClient(t *testing.T) {
	conf := config.SSHConfig{
		Address:  "192.168.8.49",
		Port:     22,
		Timeout:  "1s",
		Username: "katz.zhang",
		Password: "nebula",
	}
	go testSSHClient(conf)
	//go testSSHClient1(conf)
	//testSSHClient2(conf)
	time.Sleep(3 * time.Second)
}

func testSSHClient(conf config.SSHConfig) {
	c, err := GetSSHClient(conf.Username, conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", c)

	ch := make(chan ExecuteResult)
	go c.ExecuteAsync("vmstat 1 1", ch)
	// sudo du -sh /home/*
	// df -H | grep -vE '^Filesystem|tmpfs|udev' | awk '{ print $1 " " $2 " " $3 " " $4 " " $5 }'
	for {
		select {
		case res := <-ch:
			log.Println("\n", string(res.StdOut))
			break
		default:
		}
	}
}

func testSSHClient1(conf config.SSHConfig) {
	c1, err := GetSSHClient(conf.Username, conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", c1)

	ch1 := make(chan ExecuteResult)
	go c1.ExecuteAsync("ls", ch1)

	for {
		select {
		case res := <-ch1:
			log.Println(string(res.StdOut))
			break
		default:

		}
	}
}

func testSSHClient2(conf config.SSHConfig) {
	c2, err := GetSSHClient(conf.Username, conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", c2)

	res2, _ := c2.Execute("ls")

	log.Println(string(res2.StdOut))
}
