package config

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"os"
)

type Counter struct {
	Id    int         `json:"id" db:"id"`
	Url   pgtype.Text `json:"url" db:"url"`
	Name  string      `json:"name"`
	Code  uuid.UUID   `json:"code"`
	Count int         `json:"count" `
}

func ConnectDB(cfg *Config) (*pgxpool.Pool, error) {

	ConnectStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	dbpool, err := pgxpool.New(context.Background(), ConnectStr)
	if err != nil {
		logrus.Errorf("Unable to create connection pool: %v", err)
		os.Exit(1)
	}

	return dbpool, nil
}
