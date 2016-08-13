package db

// Constants

const (
	PRIVILEGE_ADMIN = iota
	PRIVILEGE_REVIEWER
)

// Structs

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
