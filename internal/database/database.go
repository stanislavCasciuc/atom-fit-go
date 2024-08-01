package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stanislavCasciuc/atom-fit-go/internal/config"
	"log"
)

func New(dbConfig config.DbConfig) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host,
		dbConfig.Name,
	)
	log.Print(dbURL)
	db, err := sqlx.Connect("postgres", dbURL)

	if err != nil {
		return nil, err
	}

	return db, nil
}
