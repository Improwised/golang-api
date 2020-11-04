package main

import (
	"github.com/Improwised/golang-api/cli"
	"github.com/Improwised/golang-api/config"
)

func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	err := cli.Init(cfg)

	if err != nil {
		panic(err)
	}

}
