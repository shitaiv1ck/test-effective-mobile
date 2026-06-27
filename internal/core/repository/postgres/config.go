package postgres

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	User     string        `envconfig:"USER"`
	Password string        `envconfig:"PASSWORD"`
	Host     string        `envconfig:"HOST"`
	Port     string        `envconfig:"PORT"`
	DB       string        `envconfig:"DB"`
	Timeout  time.Duration `envconfig:"TIMEOUT"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return config
}
