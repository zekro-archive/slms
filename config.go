package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type MySqlCreds struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Cert struct {
	CertFile string `yaml:"certfile"`
	KeyFile  string `yaml:"keyfile"`
}

type Config struct {
	Port          string      `yaml:"port"`
	CreationToken string      `yaml:"creationtoken"`
	RandShortLen  int         `yaml:"randshortlen"`
	Cert          *Cert       `yaml:"cert"`
	MySql         *MySqlCreds `yaml:"mysql"`
}

func ConfigOpen(path string) (*Config, error) {
	fhandler, err := os.Open(path)
	if os.IsNotExist(err) {
		fhandler, err = os.Create(path)
		if err != nil {
			return nil, err
		}
		config := &Config{
			Port:         "8080",
			RandShortLen: 8,
			Cert:         new(Cert),
			MySql:        new(MySqlCreds),
		}
		encoder := yaml.NewEncoder(fhandler)
		defer encoder.Close()
		encoder.Encode(config)
		return config, err
	} else if err != nil {
		return nil, err
	}
	config := new(Config)
	decoder := yaml.NewDecoder(fhandler)
	err = decoder.Decode(config)
	if config.RandShortLen == 0 {
		config.RandShortLen = 8
	}
	return config, err
}
