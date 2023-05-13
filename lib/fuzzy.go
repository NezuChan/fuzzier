package lib

import (
	"fmt"
	"github.com/disgoorg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/redis"
	"time"
)

type Fuzzy struct {
	Redis *redis.Redis
	App   *fiber.App
}

func InitFuzzy(conf *config.Config) *Fuzzy {
	fuzzy := Fuzzy{
		Redis: redis.InitRedis(conf.Redis),
		App: fiber.New(
			fiber.Config{
				DisableStartupMessage: true,
			}),
	}

	fuzzy.App.Use(
		func(c *fiber.Ctx) error {
			start := time.Now()

			err := c.Next()

			latency := time.Since(start)

			log.Info(fmt.Sprintf("%s:%s [%d ms] %d - %s %s",
				c.IP(),
				c.Port(),
				latency/100,
				c.Response().StatusCode(),
				c.Method(),
				c.Path(),
			))
			return err
		})

	log.Fatal(fuzzy.App.Listen(fmt.Sprintf(":%s", conf.Port)))

	return &fuzzy
}
