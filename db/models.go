// Package db provides us with everything database related.
// Connection creation, model representation and other things.
package db

const (
	PRIVILEGE_ADMIN = iota
	PRIVILEGE_REVIEWER
)

type User struct {
	ID           string `gorm:"primary_key"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Mail         string `gorm:"index;not null;unique"`
	MailVerified bool   `gorm:"not null"`
	PasswordHash string `gorm:"not null;unique"`
	Privileges   int    `gorm:"not null"`
	Enabled      bool   `gorm:"not null"`
}
