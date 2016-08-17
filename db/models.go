package db

// Constants

const (
	// Privileges have to be kept monotonic.
	// That means, a user with a lower integer
	// privilege value will also have all privileges
	// numerically greater than that value.
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
