package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBUser     string `mapstructure:"DB_USER"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	DBDriver   string `mapstructure:"DB_DRIVER"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	AppEnv     string `mapstructure:"APP_ENV"`
}

/*
type Config struct {
	Db     DB
	Server Server
	Jwt    Jwt
	AppEnv string
}

type DB struct {
	Username string
	Password string
	DBName   string
	Host     string
	Port     string
}

type Jwt struct {
	ATPrivateKey        string
	ATExpirationMinutes int
}

type Server struct {
	Port string
}
*/

func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBDriver:   os.Getenv("DB_DRIVER"),
		ServerPort: os.Getenv("SERVER_PORT"),
		AppEnv:     os.Getenv("APP_ENV"),
	}
}

func (c *Config) Print() {
	fmt.Println(*c)
}
