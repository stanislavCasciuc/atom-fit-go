package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	var migrationsPath, envPath string
	var forceVersion int

	flag.StringVar(&migrationsPath, "migrations-path", "", "migrations path")
	flag.StringVar(&envPath, "env-path", "", "path of .env file")
	flag.IntVar(&forceVersion, "force", 0, "force set version")
	flag.Parse()
	if migrationsPath == "" {
		panic("migrations path is required")
	}
	if envPath == "" {
		panic("env path is required")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		panic("Error loading .env file" + envPath + err.Error())
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName + "?sslmode=disable"

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		panic(err)
	}

	if forceVersion != 0 {
		if err := m.Force(forceVersion); err != nil {
			panic(err)
		}
		fmt.Println("Forced version", forceVersion)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "down" {
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to roll back")
				return
			}
			panic(err)
		}
		fmt.Println("migrations rolled back")
	} else {
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")
				return
			}
			panic(err)
		}
		fmt.Println("migrations applied")
	}
}
