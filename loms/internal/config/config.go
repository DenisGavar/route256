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
	LomsPgBouncer     LomsPgBouncer     `yaml:"lomsPgBouncer"`
	CancelOrderDaemon CancelOrderDaemon `yaml:"cancelOrderDaemon"`
	Kafka             Kafka             `yaml:"kafka"`
}

type Loms struct {
	Name           string `yaml:"name"`
	Port           string `yaml:"port"`
	PrometheusPort string `yaml:"prometheusPort"`
}

type CancelOrderDaemon struct {
	WorkersCount             int `yaml:"workersCount"`
	CancelOrderTimeInMinutes int `yaml:"cancelOrderTimeInMinutes"`
}

type LomsPgBouncer struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	UserDB     string `yaml:"userDB"`
	PasswordDB string `yaml:"passwordDB"`
	NameDB     string `yaml:"nameDB"`
}

type Kafka struct {
	WorkersCount   int      `yaml:"workersCount"`
	TopicForOrders string   `yaml:"topicForOrders"`
	Brokers        []string `yaml:"brokers"`
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
