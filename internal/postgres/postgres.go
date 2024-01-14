package postgres

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(config config.Config) (*gorm.DB, error) {
	fmt.Println("DSN", createDSN(config))
	db, err := gorm.Open(postgres.Open(createDSN(config)), &gorm.Config{})
	if err != nil {
		fmt.Println("ERR", err.Error())
		log.Fatalf("can not connect db %v", err)
	}

	return db, nil
}

func createDSN(config config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode)
}
