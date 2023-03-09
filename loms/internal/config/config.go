package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Services Services `yaml:"services"`
}

type Services struct {
	Loms   Loms   `yaml:"loms"`
	LomsDB LomsDB `yaml:"lomsDB"`
}

type Loms struct {
	Port string `yaml:"port"`
}

type LomsDB struct {
	Host     string `yaml:"hostDB"`
	Port     string `yaml:"portDB"`
	User     string `yaml:"userDB"`
	Password string `yaml:"passwordDB"`
	DBName   string `yaml:"nameDB"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return errors.WithMessage(err, "rparsing yaml")
	}

	return nil
}
