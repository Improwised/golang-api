package main

import (
	"github.com/Improwised/golang-api/cli"
	"github.com/Improwised/golang-api/config"
	"go.uber.org/zap"
)

func main() {

	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	var logger *zap.Logger

	cli.Init(cfg, logger)
}
