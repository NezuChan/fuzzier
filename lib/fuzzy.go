package lib

import (
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/redis"
)

type Fuzzy struct {
	Redis *redis.Redis
}

func InitFuzzy(conf *config.Config) *Fuzzy {
	fuzzy := Fuzzy{
		Redis: redis.InitRedis(conf.Redis),
	}

	return &fuzzy
}
