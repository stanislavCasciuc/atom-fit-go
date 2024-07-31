package config

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

var Envs = initConfig()

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret string
	Exp    time.Duration
}

type Config struct {
	DbCfg  DbConfig
	JwtCfg JWTConfig
	Env    string
	HttpServer
}

type HttpServer struct {
	Addr       string
	IdleTimout time.Duration
	Timeout    time.Duration
}

func initConfig() Config {
	var envPath string

	flag.StringVar(&envPath, "env-path", "", "path to .env file")
	flag.Parse()

	if envPath == "" {
		panic("env-path is required")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		panic("cannot to load the .env file")
	}

	dbCfg := DbConfig{
		Host: os.Getenv("DB_HOST"),
		Port: func() int {
			port, err := strconv.Atoi(os.Getenv("DB_PORT"))
			if err != nil {
				panic("cannot to convert DB_PORT of database")
			}
			return port
		}(),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	jwtCfg := JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
		Exp: func() time.Duration {
			expStr := os.Getenv("JWT_EXP")
			exp, err := time.ParseDuration(expStr)
			if err != nil {
				panic("cannot to convert JWT_EXP to time duration")
			}
			return exp
		}(),
	}

	httpServer := HttpServer{
		Addr: os.Getenv("HTTP_ADDR"),
		Timeout: func() time.Duration {
			timeout, err := time.ParseDuration(os.Getenv("HTTP_TIMEOUT"))
			if err != nil {
				panic("cannot convert HTTP_TIMEOUT to time duration")
			}
			return timeout
		}(),

		IdleTimout: func() time.Duration {
			idleTimeout, err := time.ParseDuration(os.Getenv("HTTP_IDLE_TIMEOUT"))
			if err != nil {
				panic("cannot convert HTTP_IDLE_TIMEOUT to time duration")
			}
			return idleTimeout
		}(),
	}

	env := os.Getenv("ENV")
	return Config{
		DbCfg:      dbCfg,
		JwtCfg:     jwtCfg,
		Env:        env,
		HttpServer: httpServer,
	}
}
