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
	CheckoutDB     CheckoutDB     `yaml:"checkoutDB"`
}

type Loms struct {
	Address string `yaml:"address"`
}

type ProductService struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

type Checkout struct {
	HTTPPort string `yaml:"httpPort"`
	GRPCPort string `yaml:"grpcPort"`
}

type CheckoutDB struct {
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
