package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/disgoorg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/redis"
	_struct "github.com/nezuchan/fuzzier/struct"
)

type FindMatch struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Duration int    `json:"duration"`
}

func FuzzyMatch(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader != config.Conf.Authorization {
        return c.Status(403).JSON(_struct.BaseJSONResponse{
            Message: "You are not authorized to to that !",
            Status:  403,
        })
    }

    var body FindMatch
    err := c.BodyParser(&body)
    if err != nil && c.Query("query") == "" {
        log.Infof("Failed to parse body from client")
        return c.Status(400).JSON(_struct.BaseJSONResponse{
            Message: "Failed to parse body you sent",
            Status:  400,
        })
    }

    values, _ := redis.Client.LRange(context.Background(), config.Conf.RedisKey, 0, -1).Result()
    resultsMap := make(map[string]bool) // Use a map to keep track of unique strings

    for _, v := range values {
        var p FindMatch
        err := json.Unmarshal([]byte(v), &p)
        if err != nil {
            panic(err)
        }

        // Convert all strings to lowercase
        title := strings.ToLower(p.Title)
        artist := strings.ToLower(p.Artist)
        duration := p.Duration

        s := fmt.Sprintf("%s - %s - %d", title, artist, duration)
        resultsMap[s] = true // Add the string to the map
    }

    // Convert the map keys back to a slice
    results := make([]string, 0, len(resultsMap))
    for s := range resultsMap {
        results = append(results, s)
    }

    // Convert the search strings to lowercase
    title := strings.ToLower(body.Title)
    artist := strings.ToLower(body.Artist)
    duration := body.Duration

    query := c.Query("query", fmt.Sprintf("%s - %s - %d", title, artist, duration))

    // Perform the search
    fuzzyResults := fuzzy.Find(query, results)

    var lastResults []FindMatch

    for _, value := range fuzzyResults {
        parts := strings.Split(value, " - ")
        if len(parts) != 3 {
            continue
        }

        // Convert duration to int
        duration := 0
        _, err := fmt.Sscanf(parts[2], "%d", &duration)
        if err != nil {
            continue
        }

        // Find the original case of the title and artist
        originalTitle, originalArtist := "", ""
        for _, v := range values {
            var p FindMatch
            err := json.Unmarshal([]byte(v), &p)
            if err != nil {
                continue
            }
            if strings.ToLower(p.Title) == parts[0] {
                originalTitle = p.Title
            }
            if strings.ToLower(p.Artist) == parts[1] {
                originalArtist = p.Artist
            }
            if originalTitle != "" && originalArtist != "" {
                break
            }
        }

        // Create a new FindMatch object with the original case of the title and artist
        s := FindMatch{
            Title:    originalTitle,
            Artist:   originalArtist,
            Duration: duration,
        }

        lastResults = append(lastResults, s)
    }

    if len(lastResults) > 0 {
        return c.JSON(lastResults)
    }

    // Return an empty array if no results were found
    return c.JSON([]FindMatch{})
}
