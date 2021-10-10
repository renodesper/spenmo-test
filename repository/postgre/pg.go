package postgre

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type DbLogger struct {
	Log logger.Logger
}

func (d DbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d DbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	d.Log.Info(fmt.Sprintf("sql: %s", string(fq)))
	return nil
}

func NewDbLogger(log logger.Logger) *DbLogger {
	return &DbLogger{
		Log: log,
	}
}

// NewPostgreClient ...
func NewPostgreClient(log logger.Logger, username, password, host string, port int, dbName string) *pg.DB {
	addr := fmt.Sprintf("%s:%d", host, port)
	db := pg.Connect(&pg.Options{
		Addr:     addr,
		User:     username,
		Password: password,
		Database: dbName,
	})

	if viper.GetString("app.env") != "production" {
		dbLogger := NewDbLogger(log)
		db.AddQueryHook(dbLogger)
	}

	_, err := db.Exec("SELECT 1")
	if err != nil {
		panic(err)
	}

	return db
}
