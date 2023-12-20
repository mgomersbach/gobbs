package dbauth

import (
	"gobbs/auth"
	"gobbs/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DBAuth struct {
	DB *gorm.DB
}

func NewDBAuth(db *gorm.DB) auth.Authenticator {
	return &DBAuth{DB: db}
}

func (d *DBAuth) Authenticate(username, password string) (bool, error) {
	var user model.User // Using the User struct from the model package

	// Fetch the user from the database
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // User not found
		}
		return false, err // Database error
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, nil // Password does not match
	}

	return true, nil // Authentication successful
}
