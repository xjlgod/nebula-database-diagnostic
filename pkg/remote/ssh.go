package remote

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	*ssh.Client
}

type ExecuteResult struct {
	CMD     string
	Err     error
	StdOut  []byte
	StdErr  []byte
	Latency time.Duration
}

var clients = make(map[string]*SSHClient)
var mux sync.RWMutex

func GetSSHClient(scid string, conf config.SSHConfig) (*SSHClient, error) {
	mux.Lock()
	if _, ok := clients[scid]; !ok {
		c, err := newSSHClient(conf)
		if err != nil {
			return nil, err
		}

		clients[scid] = c
	}
	mux.Unlock()

	mux.RLock()
	c := clients[scid]
	mux.RUnlock()

	return c, nil
}

func newSSHClient(conf config.SSHConfig) (*SSHClient, error) {
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return nil, err
	}

	sshConf := &ssh.ClientConfig{
		Timeout:         timeout,
		User:            conf.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConf.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}

	sshHost := fmt.Sprintf("%s:%d", conf.Address, conf.Port)
	sshClient, err := ssh.Dial("tcp", sshHost, sshConf)
	if err != nil {
		return nil, err
	}

	return &SSHClient{
		Client: sshClient,
	}, nil
}

func (c *SSHClient) ExecuteAsync(cmd string, ch chan<- ExecuteResult) {
	res, _ := c.Execute(cmd)
	ch <- res
}

func (c *SSHClient) Execute(cmd string) (ExecuteResult, bool) {
	now := time.Now()
	session, err := c.NewSession()
	if err != nil {
		return ExecuteResult{cmd, err, []byte{}, []byte{}, time.Since(now)}, false
	}
	defer session.Close()

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	session.Stdout = &stdOut
	session.Stderr = &stdErr

	err = session.Run(cmd)
	if err != nil || len(stdErr.Bytes()) > 0 {
		return ExecuteResult{cmd, err, stdOut.Bytes(), stdErr.Bytes(), time.Since(now)}, false
	}

	return ExecuteResult{cmd, err, stdOut.Bytes(), stdErr.Bytes(), time.Since(now)}, true
}
