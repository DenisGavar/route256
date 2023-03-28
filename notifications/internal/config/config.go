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
	Kafka Kafka `yaml:"kafka"`
}

type Kafka struct {
	BalanceStrategy string   `yaml:"balanceStrategy"`
	GroupName       string   `yaml:"groupName"`
	TopicForOrders  string   `yaml:"topicForOrders"`
	Brokers         []string `yaml:"brokers"`
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
