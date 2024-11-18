package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	HTTPServer HTTPServer
	Postgres   Postgres
	JWT        JWT
}

type HTTPServer struct {
	Address string `env:"HTTP_HOST"`
	Port    string `env:"HTTP_PORT"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type JWT struct {
	Secret string        `env:"JWT_SECRET"`
	Expire time.Duration `env:"JWT_EXPIRE"`
}

func MustLoad() *Config {
	var config Config

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalln(err)
	}

	return &config
}
