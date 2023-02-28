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
	Loms           Loms           `yaml:"loms"`
	ProductService ProductService `yaml:"productService"`
	Checkout       Checkout       `yaml:"checkout"`
}

type Loms struct {
	Url string `yaml:"url"`
}

type ProductService struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type Checkout struct {
	Port string `yaml:"port"`
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
