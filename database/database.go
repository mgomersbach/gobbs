package database

import (
	"gobbs/config"
	"gobbs/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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

// InitializeDatabase sets up the required tables in the database
func InitializeDatabase(db *gorm.DB, cfg *config.Config) error {
	// AutoMigrate as before
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	// Add predefined users from the config
	for _, u := range cfg.Users {
		var count int64
		db.Model(&model.User{}).Where("username = ?", u.Username).Count(&count)

		if count == 0 {
			// User does not exist, create a new one
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			newUser := model.User{
				Username: u.Username,
				Password: string(hashedPassword),
			}

			db.Create(&newUser)
		}
	}

	return nil
}
