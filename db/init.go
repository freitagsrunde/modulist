// Package db provides us with everything database related.
// Connection creation, model representation and other things.
package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Functions

// InitDB connects to the database specified by the .env file.
// It returns the correctly configured connector.
func InitDB() *gorm.DB {

	var db *gorm.DB

	// Read from .env file what database we are connecting to.
	dbType := os.Getenv("DB_TYPE")

	if dbType == "postgres" {

		// Read database port and convert to integer.
		port, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			log.Fatal("[InitDB] Unrecognized port type in .env file, integer expected. Terminating.")
		}

		// Fill db variable with real connection.
		db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PW"), os.Getenv("DB_HOST"),
			port, os.Getenv("DB_DBNAME"), os.Getenv("DB_SSLMODE")))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("[InitDB] Unsupported database type in environment file, please use PostgreSQL. Did you forget to specify a database in your .env file? Terminating.")
	}

	// Check connection to database in order to be sure.
	err := db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	// If app runs in development mode, log SQL queries.
	if os.Getenv("DEPLOY_STAGE") == "dev" {
		db.LogMode(true)
	}

	return db
}

// CreateTables sets up the connected database correctly by first
// deleting all considered tables and afterwards setting new ones
// up correctly.
func SetUpTables(db *gorm.DB) {

	// Delete all tables corresponding to models if they exist.
	db.DropTableIfExists(&User{})
	db.DropTableIfExists(&PasswordLink{})
	db.DropTableIfExists(&Module{})
	db.DropTableIfExists(&Person{})

	// Create new ones for all models.
	db.CreateTable(&User{})
	db.CreateTable(&PasswordLink{})
	db.CreateTable(&Module{})
	db.CreateTable(&Person{})
}

// TransferPersons connects to the provided SQLite database
// containing the persons involved in the faculty's modules
// and exports them into the services's main database.
func TransferPersons(db *gorm.DB, modulesDBPath string) {

	// Connect to SQLite database.
	modulesDB, err := gorm.Open("sqlite3", modulesDBPath)
	if err != nil {
		log.Fatal(err)
	}

	// Select all persons from SQLite database.
	var sqlitePersons []SQLitePerson
	modulesDB.Find(&sqlitePersons)

	for _, sqlitePerson := range sqlitePersons {

		// Convert each person from SQLite database to
		// person made to be stored in main database.
		var person Person
		person = sqlitePerson.ToPerson()

		// Save person to main database.
		db.Create(&person)
	}

	// Close connection to database.
	modulesDB.Close()
}

// TransferModules connects to the provided SQLite database containing
// the modules and exports them into the services's main database.
func TransferModules(db *gorm.DB, modulesDBPath string) {

	// Connect to SQLite database.
	modulesDB, err := gorm.Open("sqlite3", modulesDBPath)
	if err != nil {
		log.Fatal(err)
	}

	// Select all MTS modules from SQLite database.
	var sqliteModules []SQLiteModule
	modulesDB.Find(&sqliteModules)

	for _, sqliteModule := range sqliteModules {

		// Convert each module from SQLite database to
		// module made to be stored in main database.
		var module Module
		module = sqliteModule.ToModule(db)

		// Save it to main database.
		db.Create(&module)
	}

	// Close connection to database.
	modulesDB.Close()
}
