// Golang API.
//
//     Schemes: https
//     Host: localhost
//     BasePath: /api/v1
//     Version: 0.0.1-alpha
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"github.com/Improwised/golang-api/cli"
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/routinewrapper"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"time"
)

func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	var logger *zap.Logger
	var err error

	if config.GetConfigByName("IS_DEVELOPMENT") == "true" {
		logger, err = zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	} else {
		logger, err = zap.NewProduction(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	}

	// this function will logged error log in sentry
	sentryLoggedFunc := func() {
		err := recover()

		if err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 2)
		}
	}

	// routine wrapper will handle go routine error also an log into sentry
	routinewrapper.Init(sentryLoggedFunc)
	defer sentryLoggedFunc()

	err = cli.Init(cfg, logger)
	if err != nil {
		panic(err)
	}

}
