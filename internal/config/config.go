package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

const (
	dev  string = "local"
	prod string = "prod"
)

type Config struct {
	Environment string `env:"ENVIRONMENT"`
	Admin       Admin
	HTTPServer  HTTPServer
	Postgres    Postgres
	JWT         JWT
}

type Admin struct {
	Email    string `env:"ADMIN_EMAIL"`
	Password string `env:"ADMIN_PASSWORD"`
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

func (c *Config) IsProd() bool {
	return c.Environment == prod
}

func (c *Config) IsLocal() bool {
	return c.Environment == dev
}
