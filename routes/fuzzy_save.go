package routes

import (
	"context"
	"encoding/json"
	"github.com/disgoorg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/redis"
	_struct "github.com/nezuchan/fuzzier/struct"
)

func FuzzySave(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != config.Conf.Authorization {
		return c.Status(403).JSON(_struct.BaseJSONResponse{
			Message: "You are not authorized to to that !",
			Status:  403,
		})
	}

	var body FindMatch
	err := c.BodyParser(&body)
	if err != nil {
		log.Infof("Failed to parse body from client")
		return c.Status(400).JSON(_struct.BaseJSONResponse{
			Message: "Failed to parse body you sent",
			Status:  400,
		})
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Infof("Failed to marshal body to JSON: %v", err)
		return c.Status(500).JSON(_struct.BaseJSONResponse{
			Message: "Failed to marshal body to JSON",
			Status:  500,
		})
	}

	_, err = redis.Client.RPush(context.Background(), config.Conf.RedisKey, jsonBody).Result()

	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(_struct.BaseJSONResponse{
		Message: "Success added track to redis",
		Status:  200,
	})
}
