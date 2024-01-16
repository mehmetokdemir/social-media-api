package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBPassword             string `mapstructure:"DB_PASSWORD"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBName                 string `mapstructure:"DB_NAME"`
	DBSSLMode              string `mapstructure:"DB_SSLMODE"`
	DBDriver               string `mapstructure:"DB_DRIVER"`
	ServerPort             string `mapstructure:"SERVER_PORT"`
	AppEnv                 string `mapstructure:"APP_ENV"`
	JwtATPrivateKey        string `mapstructure:"JWT_AT_PRIVATE_KEY"`
	JwtATExpirationMinutes int    `mapstructure:"JWT_AT_EXPIRATION_MIN"`
	CloudinaryCloudName    string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryApiKey       string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryApiSecret    string `mapstructure:"CLOUDINARY_API_SECRET"`
}

func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	fmt.Println("APP ENV", os.Getenv("APP_ENV"))

	jwtExpirationMin, _ := strconv.Atoi(os.Getenv("JWT_AT_EXPIRATION_MIN"))

	return Config{
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBUser:                 os.Getenv("DB_USER"),
		DBName:                 os.Getenv("DB_NAME"),
		DBSSLMode:              os.Getenv("DB_SSLMODE"),
		DBDriver:               os.Getenv("DB_DRIVER"),
		ServerPort:             os.Getenv("SERVER_PORT"),
		AppEnv:                 os.Getenv("APP_ENV"),
		JwtATPrivateKey:        os.Getenv("JWT_AT_PRIVATE_KEY"),
		JwtATExpirationMinutes: jwtExpirationMin,
		CloudinaryCloudName:    os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryApiKey:       os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryApiSecret:    os.Getenv("CLOUDINARY_API_SECRET"),
	}
}

func (c *Config) Print() {
	fmt.Println(*c)
}
