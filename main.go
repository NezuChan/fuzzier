package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/lib"
)

func main() {
	config.Init()
	lib.InitFuzzy()

	select {}
}
