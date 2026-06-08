package database

import (
    "fmt"
    "log"
    "api_hr/internal/config"
    "api_hr/internal/models"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Автомиграция
    if err := autoMigrate(db); err != nil {
        return nil, err
    }

    log.Println("Successfully connected to database")
    return db, nil
}

func autoMigrate(db *gorm.DB) error {
    // автоматические миграции GORM
    return db.AutoMigrate(
        &models.UserT{},
        // другие модели...
    )
}