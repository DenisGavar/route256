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
	Loms              Loms              `yaml:"loms"`
	ProductService    ProductService    `yaml:"productService"`
	Checkout          Checkout          `yaml:"checkout"`
	CheckoutPgBouncer CheckoutPgBouncer `yaml:"checkoutPgBouncer"`
}

type Loms struct {
	Address string `yaml:"address"`
}

type ProductService struct {
	Address              string `yaml:"address"`
	Token                string `yaml:"token"`
	ListCartWorkersCount int    `yaml:"listCartWorkersCount"`
	RateLimit            int    `yaml:"rateLimit"`
}

type Checkout struct {
	HTTPPort string `yaml:"httpPort"`
	GRPCPort string `yaml:"grpcPort"`
}

type CheckoutPgBouncer struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	UserDB     string `yaml:"userDB"`
	PasswordDB string `yaml:"passwordDB"`
	NameDB     string `yaml:"nameDB"`
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
