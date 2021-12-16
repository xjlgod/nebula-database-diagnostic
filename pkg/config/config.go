package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Nodes map[string]NodeConfig `mapstructure:"nodes"`
}

type NodeConfig struct {
	Host     HostConfig   `mapstructure:"host"`
	SSH      SSHConfig    `mapstructure:"ssh"`
	Duration string       `mapstructure:"duration"` // TODO parse into time.Duration
	Interval string       `mapstructure:"interval"` // TODO parse into time.Duration
	Output   string       `mapstructure:"output"`
	Infos    []InfoOption `mapstructure:"info,omitempty"`
	Diags    []DiagOption `mapstructure:"diag,omitempty"`
}

type HostConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type SSHConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
}

type InfoOption string

const (
	Metrics        InfoOption = "metrics"
	Physical       InfoOption = "physical"
	GraphService   InfoOption = "graphService"
	MetaService    InfoOption = "metaService"
	StorageService InfoOption = "storageService"
	Stats          InfoOption = "stats"
	AllInfo        InfoOption = "all"
)

type DiagOption string

const (
	Partition DiagOption = "partition"
	AllDiag   DiagOption = "all"
)

func NewConfig(confPath string, configType string) (*Config, error) {
	var viperConfig = viper.New()
	viperConfig.SetConfigName(confPath)
	viperConfig.SetConfigFile(confPath)
	viperConfig.SetConfigType(configType)
	if err := viperConfig.ReadInConfig(); err != nil {
		return nil, err
	}

	conf := new(Config)
	if err := viperConfig.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
