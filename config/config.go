package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/disgoorg/log"
	"github.com/nezuchan/fuzzier/redis"
)

type Config struct {
	Redis         redis.Config
	Port          string `env:"PORT" envDefault:"3000"`
	Authorization string `env:"AUTHORIZATION,required"`
	RedisKey      string `env:"REDIS_KEY,required"`
	StoreTimeout  int    `env:"STORE_TIMEOUT" envDefault:"86400000"`
}

var Conf Config

func Init() {
	Conf = Config{}
	if err := env.Parse(&Conf); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
