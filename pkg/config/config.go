package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Infos map[string]*InfoConfig `mapstructure:"infos"`
	Diag  *DiagConfig            `mapstructure:"diag"`
}

type InfoConfig struct {
	Node     NodeConfig               `mapstructure:"node"`
	Output   OutputConfig             `mapstructure:"output"`           // output location
	Duration string                   `mapstructure:"duration"`         // TODO parse into time.Duration, default is -1
	Period   string                   `mapstructure:"period"`           // TODO parse into time.Duration, default is 0
	Options  []InfoOption             `mapstructure:"option,omitempty"` // info to fetch, default is all
	Services map[string]ServiceConfig `mapstructure:"services"`
}

type DiagConfig struct {
	Output  OutputConfig `mapstructure:"output"` // output location
	Input   InputConfig  `mapstructure:"input,omitempty"`
	Options []DiagOption `mapstructure:"option,omitempty"` // diag result to analyze, default is no
}

type NodeConfig struct {
	Host HostConfig `mapstructure:"host"` // node host
	SSH  SSHConfig  `mapstructure:"ssh"`  // node ssh
}

type HostConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func (c HostConfig) IsValid() bool {
	return c.Address != "" && c.Port > 0 // TODO add more exactly verify: address, port
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
	DirPath   string `mapstructure:"dirPath"` // output dir included logs, info, diag, etc., default is ./out, and will auto create if not existed
	LogToFile bool   `mapstructure:"logToFile"`
}

func (c OutputConfig) IsValid() bool {
	return c.DirPath != "" // TODO add more exactly verify: dirPath
}

type InputConfig struct {
	DirPath string `mapstructure:"dirPath"` // input dir included logs, info, etc.
	Remote  bool   `mapstructure:"remote"`  // remote = true means that input dir is located at the remote node
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
	defaultDuration      = "-1"
	defaultPeriod        = "5s"
	defaultOutputDirPath = "./output"
	defaultInfoOptions   = []InfoOption{AllInfo}
	defaultDiagOptions   = []DiagOption{NoDiag}
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

	ConfigComplete(conf)
	if ids, ok := ConfigValidate(conf); !ok {
		err := fmt.Errorf("such nodes config is invalid: %+v", ids)
		return nil, err
	}
	return conf, nil
}

func (c *Config) String() string {
	sb := strings.Builder{}
	for k, info := range c.Infos {
		sb.WriteString(k + ": ")
		sb.WriteString(fmt.Sprintf("%+v\n", info))
	}
	sb.WriteString(fmt.Sprintf("%+v", c.Diag))
	return sb.String()
}

func ConfigComplete(conf *Config) {
	for _, info := range conf.Infos {
		node := info.Node
		if node.SSH.Timeout == "" {
			node.SSH.Timeout = "3s"
		}
		if info.Duration == "" {
			info.Duration = defaultDuration
		}
		if info.Period == "" {
			info.Period = defaultPeriod
		}
		if info.Output.DirPath == "" {
			info.Output.DirPath = defaultOutputDirPath
		}
		if len(info.Options) == 0 {
			info.Options = defaultInfoOptions
		}
	}
	if conf.Diag.Output.DirPath == "" {
		conf.Diag.Output.DirPath = defaultOutputDirPath
	}
	if len(conf.Diag.Options) == 0 {
		conf.Diag.Options = defaultDiagOptions
	}
}

func ConfigValidate(conf *Config) ([]string, bool) {
	ids := make([]string, 0)
	ok := true
	for k, info := range conf.Infos {
		node := info.Node
		if !node.Host.IsValid() || !node.SSH.IsValid() || !info.Output.IsValid() || !isValidDuration(info.Duration) || !isValidPeriod(info.Period) {
			ids = append(ids, k)
			ok = false
		}
	}

	if !conf.Diag.Output.IsValid() {
		ok = false
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
