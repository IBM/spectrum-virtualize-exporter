// Copyright 2021-2024 IBM Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Targets         []Target        `yaml:"targets"`
	ExtraLabels     []Label         `yaml:"extra_labels"`
	TlsServerConfig TlsServerConfig `yaml:"tls_server_config"`
	filename        string
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

type TlsServerConfig struct {
	CaCert     string `yaml:"ca_cert"`
	ServerCert string `yaml:"server_cert"`
	ServerKey  string `yaml:"server_key"`
}

func (cfg *Config) SetFilename(filename string) {
	cfg.filename = filename
}

// Load loads a config from filename
func (cfg *Config) _Init() (*Config, error) {

	content, err := os.ReadFile(cfg.filename)
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
	cfg.SetFilename(filename)
	return cfg._Init()
}
