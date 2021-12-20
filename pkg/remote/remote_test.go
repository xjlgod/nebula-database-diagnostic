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
	go testSSHClient1(conf)
	testSSHClient2(conf)
	time.Sleep(3 * time.Second)
}

func testSSHClient(conf config.SSHConfig) {
	c, err := GetSSHClient(conf.Username, conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", c)

	ch := make(chan ExecuteResult)
	go c.ExecuteAsync("whoami", ch)

	res := <-ch
	log.Println(string(res.StdOut))
}

func testSSHClient1(conf config.SSHConfig) {
	c1, err := GetSSHClient(conf.Username, conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", c1)

	ch1 := make(chan ExecuteResult)
	go c1.ExecuteAsync("ls", ch1)

	res1 := <-ch1
	log.Println(string(res1.StdOut))
}

func testSSHClient2(conf config.SSHConfig) {
	c1, err := GetSSHClient(conf.Username, conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", c1)

	res, _ := c1.Execute("ls")

	log.Println(string(res.StdOut))
}
