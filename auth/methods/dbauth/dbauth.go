package dbauth

import (
	"gobbs/auth"

	"gorm.io/gorm"
)

type DBAuth struct {
	DB *gorm.DB
}

func NewDBAuth(db *gorm.DB) auth.Authenticator {
	return &DBAuth{DB: db}
}

func (d *DBAuth) Authenticate(username, password string) (bool, error) {
	// Implement authentication logic using the database
	return true, nil
}
