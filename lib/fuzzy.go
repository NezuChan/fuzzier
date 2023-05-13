package lib

import (
	"fmt"
	"github.com/disgoorg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/redis"
	"github.com/nezuchan/fuzzier/routes"
)

type Fuzzy struct {
	Redis *redis.Redis
	App   *fiber.App
}

func InitFuzzy() {
	redis.InitRedis(config.Conf.Redis)

	App := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		})

	App.Use(
		func(c *fiber.Ctx) error {
			log.Info(fmt.Sprintf("%s:%s %d - %s %s",
				c.IP(),
				c.Port(),
				c.Response().StatusCode(),
				c.Method(),
				c.Path(),
			))
			return c.Next()
		})

	App.Post("/match", routes.FuzzyMatch)
	App.Post("/", routes.FuzzySave)

	log.Fatal(App.Listen(fmt.Sprintf(":%s", config.Conf.Port)))

}
