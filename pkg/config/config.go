package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Nodes map[string]*NodeConfig `mapstructure:"nodes"`
}

type NodeConfig struct {
	Host     HostConfig               `mapstructure:"host"`           // node host
	SSH      SSHConfig                `mapstructure:"ssh"`            // node ssh
	Output   OutputConfig             `mapstructure:"output"`         // output location
	Duration string                   `mapstructure:"duration"`       // TODO parse into time.Duration, default is -1
	Period   string                   `mapstructure:"period"`         // TODO parse into time.Duration, default is 0
	Infos    []InfoOption             `mapstructure:"info,omitempty"` // info to fetch, default is all
	Diags    []DiagOption             `mapstructure:"diag,omitempty"` // diag result to analyze, default is no
	Services map[string]ServiceConfig `mapstructure:"services"`
}

type HostConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type ServiceConfig struct {
	Type       ComponentType `mapstructure:"type"`
	DeployDir  string        `mapstructure:"deploy_dir"`
	RuntimeDir string        `mapstructure:"runtime_dir"`
	Port       int           `mapstructure:"port"`
	HTTPPort   int           `mapstructure:"http_port"`
	HTTP2Port  int           `mapstructure:"http_2_port"`
}

func (s ServiceConfig) IsValid() bool {
	return s.HTTPPort > 0 // TODO add more exactly verify: DeployDir, RuntimeDir
}

func (c HostConfig) IsValid() bool {
	return c.Address != "" && c.Port > 0 // TODO add more exactly verify: address, port
}

type SSHConfig struct {
	// ssh address equals to service address
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Timeout  string `mapstructure:"timeout"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`

	// TODO support the private key login
}

func (c SSHConfig) IsValid() bool {
	return c.Address != "" && c.Port > 0 && isValidTimeout(c.Timeout) && c.Username != "" && c.Password != "" // TODO add more exactly verify: port
}

type OutputConfig struct {
	DirPath string `mapstructure:"dirPath"` // output dir included logs, info, diag, etc., default is ./out, and will auto create if not existed
	Remote  bool   `mapstructure:"remote"`  // remote = true means that output dir is located at the remote node
}

func (c OutputConfig) IsValid() bool {
	return c.DirPath != "" // TODO add more exactly verify: dirPath
}

type InfoOption string

type ComponentType string

const (
	Metrics        InfoOption    = "metrics"
	Physical       InfoOption    = "physical"
	GraphService   ComponentType = "graphService"
	MetaService    ComponentType = "metaService"
	StorageService ComponentType = "storageService"
	Stats          InfoOption    = "stats"
	AllInfo        InfoOption    = "all"
	NoInfo         InfoOption    = "no"
)

type DiagOption string

const (
	Partition DiagOption = "partition"
	AllDiag   DiagOption = "all"
	NoDiag    DiagOption = "no"
)

var (
	defaultDuration = "-1"
	defaultPeriod   = "5s"
	defaultDirPath  = ""
	defaultInfos    = []InfoOption{AllInfo}
	defaultDiags    = []DiagOption{NoDiag}
)

func GetConfigType(confPath string) string {
	if strings.HasSuffix(confPath, "yaml") {
		return "yaml"
	}

	return ""
}

func NewConfig(confPath string, configType string) (*Config, error) {
	var viperConfig = viper.New()
	viperConfig.SetConfigName(confPath)
	viperConfig.SetConfigFile(confPath)
	viperConfig.SetConfigType(configType)
	if err := viperConfig.ReadInConfig(); err != nil {
		return nil, err
	}

	conf := new(Config)
	err := viperConfig.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	configComplete(conf)
	if ids, ok := configValidate(conf); !ok {
		err := fmt.Errorf("such nodes config is invalid: %+v", ids)
		return nil, err
	}
	return conf, nil
}

func (c *Config) String() string {
	sb := strings.Builder{}
	for k, node := range c.Nodes {
		sb.WriteString(k + ": ")
		sb.WriteString(fmt.Sprintf("%+v\n", node))
	}
	return sb.String()
}

func configComplete(conf *Config) {
	for _, node := range conf.Nodes {
		if node.SSH.Timeout == "" {
			node.SSH.Timeout = "3s"
		}
		if node.Duration == "" {
			node.Duration = defaultDuration
		}
		if node.Period == "" {
			node.Period = defaultPeriod
		}
		if node.Output.DirPath == "" {
			node.Output.DirPath = defaultDirPath
		}
		if len(node.Infos) == 0 {
			node.Infos = defaultInfos
		}
		if len(node.Diags) == 0 {
			node.Diags = defaultDiags
		}
	}
}

func configValidate(conf *Config) ([]string, bool) {
	ids := make([]string, 0)
	ok := true
	for k, node := range conf.Nodes {
		if !node.Host.IsValid() || !node.SSH.IsValid() || !node.Output.IsValid() || !isValidDuration(node.Duration) || !isValidPeriod(node.Period) {
			ids = append(ids, k)
			ok = false
		}
	}

	return ids, ok
}

func isValidDuration(duration string) bool {
	_, err := time.ParseDuration(duration)
	return duration == "-1" || err == nil
}

func isValidPeriod(period string) bool {
	d, err := time.ParseDuration(period)
	return d > 0 && err == nil
}

func isValidTimeout(timeout string) bool {
	d, err := time.ParseDuration(timeout)
	return d >= 0 && err == nil
}
