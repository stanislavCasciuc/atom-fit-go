package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stanislavCasciuc/atom-fit-go/internal/config"
)

func New(dbConfig config.DbConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host,
		dbConfig.Name,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
