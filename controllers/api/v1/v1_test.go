package v1_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Improwised/golang-api/cli"
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"
	"github.com/Improwised/golang-api/logger"
	"github.com/doug-martin/goqu/v9"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var client *resty.Client = nil
var db *goqu.Database = nil

func TestMain(m *testing.M) {
	os.Chdir("../../../")
	cfg := config.LoadTestEnv()
	logger, err := logger.NewRootLogger(true, true)
	if err != nil {
		log.Fatal(err)
	}

	db, err = database.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://%s", cfg.Port)

	client = resty.New().SetBaseURL(url)

	cmd := cli.GetAPICommandDef(cfg, logger)

	// execute migration in sqlite
	migrationCmd := cli.GetMigrationCommandDef(cfg)
	migrationCmd.SetArgs([]string{"up"})
	err = migrationCmd.Execute()
	if err != nil {
		logger.Fatal("error while execute migration", zap.Error(err))
	}

	go func() {
		err = cmd.Execute()
		if err != nil {
			log.Fatal(err)
		}
	}()

	serverRunning := false
	for count := 0; count < 100; count += 1 {
		client = client.SetTimeout(time.Second * 2)
		res, _ := client.R().EnableTrace().Get("/healthz")
		if err == nil {
			log.Println("received status code", res.StatusCode())
		}
		if err == nil && res.StatusCode() == http.StatusOK {
			serverRunning = true
			break
		}
	}

	if !serverRunning {
		log.Fatal("program exit due to server is not running...")
	}

	client = client.SetTimeout(time.Second * 10)
	log.Println("server is running...")
	os.Exit(m.Run())
}
