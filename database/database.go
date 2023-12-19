package database

import (
	"gobbs/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect establishes a database connection based on the provided configuration
func Connect(cfg *config.Config, log *logrus.Logger) (*gorm.DB, error) {
	// Database Setup with switch case for different database types
	var db *gorm.DB
	var err error
	switch viper.GetString("database.type") {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(viper.GetString("database.dsn")), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(viper.GetString("database.dsn")), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(viper.GetString("database.dsn")), &gorm.Config{})
	// Add other cases as needed
	default:
		log.Fatal("Unsupported database type")
	}

	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	return db, nil
}
