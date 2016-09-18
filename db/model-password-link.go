package db

import (
	"time"
)

// Structs

type PasswordLink struct {
	ID          string    `gorm:"primary_key"`
	UserID      string    `gorm:"index;not null"`
	User        User      `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer;"`
	SecretToken string    `gorm:"not null;unique"`
	Expires     time.Time `gorm:"not null"`
}
