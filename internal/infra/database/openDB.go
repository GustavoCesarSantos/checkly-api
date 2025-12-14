package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"GustavoCesarSantos/checkly-api/internal/shared/configs"
)

func OpenDB() (*sql.DB, error) {
	var databaseConfig = configs.LoadDatabaseConfig()
	db, err := sql.Open("postgres", databaseConfig.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(databaseConfig.MaxOpenConns)
	db.SetMaxIdleConns(databaseConfig.MaxIdleConns)
	db.SetConnMaxIdleTime(databaseConfig.MaxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		db.Close()
		return nil, pingErr
	}
	return db, nil
}
