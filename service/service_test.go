package service

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"gitlab.com/renodesper/spenmo-test/repository/postgre"
	"gitlab.com/renodesper/spenmo-test/util/logger"
	"gitlab.com/renodesper/spenmo-test/util/logger/noop"
	"gitlab.com/renodesper/spenmo-test/util/logger/zap"
)

var Log *logger.Loggers
var DB *pg.DB
var M *migrate.Migrate

func TestMain(m *testing.M) {
	loadConfig()
	initService()

	code := m.Run()
	os.Exit(code)
}

func loadConfig() {
	// Command line flag
	configFile := flag.String("config", "../config/env/test.toml", "configuration path")
	flag.Parse()

	viper.SetConfigFile(*configFile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initService() {
	if testing.Verbose() {
		z, _ := zap.CreateLogger("development", "debug")
		Log = logger.New(z)
	} else {
		Log = logger.New(noop.CreateLogger())
	}

	if !testing.Short() {
		dbUsername := viper.GetString("db.username")
		dbPassword := viper.GetString("db.password")
		dbHost := viper.GetString("db.host")
		dbPort := viper.GetInt("db.port")
		dbName := viper.GetString("db.name")

		if testing.Verbose() {
			DB = postgre.NewPostgreClient(Log, dbUsername, dbPassword, dbHost, dbPort, dbName)
		} else {
			DB = postgre.NewPostgreClient(Log, dbUsername, dbPassword, dbHost, dbPort, dbName)
		}

		var err error
		M, err = migrate.New(
			"file://../config/db/migrations",
			fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUsername, dbPassword, dbHost, dbPort, dbName),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func reinitializeDB() {
	if err := M.Down(); err != nil {
		log.Fatal(err)
	}
	if err := M.Up(); err != nil {
		log.Fatal(err)
	}
}
