package main

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/lib"
	"github.com/nezuchan/fuzzier/redis"
)

func main() {
	config.Init()

	log.Info(fmt.Sprintf("Fuzzier now listening on port %s", config.Conf.Port))
	log.Info(fmt.Sprintf("Fuzzier will delete cached results every %d ms", config.Conf.StoreTimeout))

	go func () {
		intervalDuration := time.Duration(config.Conf.StoreTimeout) * time.Millisecond
    	ticker := time.NewTicker(intervalDuration)
    	defer ticker.Stop()

		for range ticker.C {
			redis.Client.Unlink(context.Background(), config.Conf.RedisKey)
			log.Infof("Deleted cached fuzzier results")
		}	
	}()

	lib.InitFuzzy()

	select {}
}
