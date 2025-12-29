package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormDB *gorm.DB

func ConnectGorm() {

	dsn := "postgres://postgres:waleedhaider@localhost:5432/jobboard?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal(" Failed to connect to database with GORM:", err)
	}

	GormDB = db
	log.Println(" GORM database connected successfully")
}
