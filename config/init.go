package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Router struct {
		RequestID bool  `yaml:"request_id"`
		RealIP    bool  `yaml:"real_ip"`
		Logger    bool  `yaml:"real_ip"`
		Recoverer bool  `yaml:"recoverer"`
		Timeout   int64 `yaml:"timeout"`
	} `yaml:"router"`
}

func NewConfig() (conf Config, err error) {
	file, err := ioutil.ReadFile("cofig/config.yml")
	if err != nil {
		return Config{}, err
	}

	return conf, yaml.Unmarshal(file, &conf)
}
