package config

import (
	"github.com/labstack/echo/v4"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormDB *gorm.DB

func ConnectGorm(e *echo.Echo) {
    dsn := "postgres://postgres:waleedhaider@localhost:5432/jobboard?sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        e.Logger.Fatal("Failed to connect to database with GORM:", err)
    }

    GormDB = db
    e.Logger.Info("GORM database connected successfully")
}
