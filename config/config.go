package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/disgoorg/log"
	"github.com/nezuchan/fuzzier/redis"
)

type Config struct {
	Redis redis.Config
	Port  string `env:"PORT" envDefault:"3000"`
}

func Init() (conf Config, err error) {
	conf = Config{}
	if err := env.Parse(&conf); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return conf, nil
}
