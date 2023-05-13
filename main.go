package main

import (
	"fmt"
	"github.com/nezuchan/fuzzier/config"
	"github.com/nezuchan/fuzzier/lib"
)

func main() {
	conf, err := config.Init()
	if err != nil {
		panic(fmt.Sprintf("couldn't initialize config: %v", err))
	}

	lib.InitFuzzy(&conf)

	select {}
}
