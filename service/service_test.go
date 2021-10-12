package service

import (
	"flag"
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
	"gitlab.com/renodesper/spenmo-test/repository/postgre"
	"gitlab.com/renodesper/spenmo-test/util/logger"
	"gitlab.com/renodesper/spenmo-test/util/logger/noop"
	"gitlab.com/renodesper/spenmo-test/util/logger/zap"
)

var Log *logger.Loggers
var DB *pg.DB

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
	dbUsername := viper.GetString("db.username")
	dbPassword := viper.GetString("db.password")
	dbHost := viper.GetString("db.host")
	dbPort := viper.GetInt("db.port")
	dbName := viper.GetString("db.name")

	if testing.Verbose() {
		z, _ := zap.CreateLogger("development", "debug")
		Log = logger.New(z)
		DB = postgre.NewPostgreClient(Log, dbUsername, dbPassword, dbHost, dbPort, dbName)
	} else {
		Log = logger.New(noop.CreateLogger())
		DB = postgre.NewPostgreClient(Log, dbUsername, dbPassword, dbHost, dbPort, dbName)
	}
}
