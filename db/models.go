package db

// Constants

const (
	// Privileges have to be kept monotonic.
	// That means, a user with a lower integer
	// privilege value will also have all privileges
	// numerically greater than that value.
	// CAUTION: Changes here will need to be reflected
	// to other places, e.g. template and handler
	// functions of admin's users site.
	PRIVILEGE_ADMIN = iota
	PRIVILEGE_REVIEWER

	// Status groups as increasing integer.
	// CAUTION: Changes here will need to be reflected
	// to other places, e.g. template and handler
	// functions of admin's users site and initial user.
	STATUS_GROUP_PROF = iota
	STATUS_GROUP_WIMI
	STATUS_GROUP_STUDI
	STATUS_GROUP_OTHER
)

// Structs

type User struct {
	ID           string `gorm:"primary_key"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Mail         string `gorm:"index;not null;unique"`
	MailVerified bool   `gorm:"not null"`
	PasswordHash string `gorm:"not null;unique"`
	StatusGroup  int    `gorm:"not null"`
	Privileges   int    `gorm:"not null"`
	Enabled      bool   `gorm:"not null"`
}
