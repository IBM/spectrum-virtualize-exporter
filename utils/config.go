package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Targets     []Target `yaml:"targets"`
	ExtraLabels []Label  `yaml:"extra_labels"`
}

type Target struct {
	IpAddress  string `yaml:"ipAddress"`
	Userid     string `yaml:"userid"`
	Password   string `yaml:"password"`
	VerifyCert bool   `yaml:"verifyCert"`
}

type Label struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

//Load loads a config from filename
func (cfg *Config) _Init(filename string) (*Config, error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
func GetConfig(filename string) (*Config, error) {
	var cfg Config
	return cfg._Init(filename)
}
