package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

// Structs

// App struct contains all relevant information read
// from .env file and pointers to connectors of middleware.
type App struct {
	TLS         bool
	TLSCertFile string
	TLSKeyFile  string
	IP          string
	Port        string
	Stage       string
	HashCost    int
	JWTValidFor time.Duration
	Router      *gin.Engine
	DB          *gorm.DB
	Validator   *validator.Validate
}

// Functions

// DefineRoutes registers all frontend endpoints.
// These will make up the web interface of MODULIST.
func (app *App) DefineRoutes() {

	// Route 'index'.
	app.Router.GET("/", app.Index)
	app.Router.POST("/", app.Login)
	app.Router.GET("/logout", app.Logout)

	// Route 'list'.
	app.Router.GET("/modules", app.ListModules)
	app.Router.GET("/modules/search", app.SearchModules)
	app.Router.GET("/modules/filter/:firstLetter", app.FilterModulesByLetter)
	app.Router.GET("/modules/done/:id", app.MarkModuleDone)

	// Route 'feedback'.
	app.Router.GET("/review/module/:moduleID", app.ReviewModule)
	app.Router.POST("/review/module/:moduleID/add/:category", app.AddFeedback)
	app.Router.GET("/review/delete/:id", app.DeleteFeedback)

	// Route 'settings'.
	app.Router.GET("/settings", app.ListSettings)
	app.Router.POST("/settings", app.UpdateSettings)

	// Route 'admin'.
	app.Router.GET("/admin/users", app.ListUsers)
	app.Router.POST("/admin/users", app.CreateUser)
	app.Router.GET("/admin/users/delete/:id", app.DeleteUser)
	app.Router.GET("/admin/send-feedback", app.SendFeedback)
	app.Router.POST("/admin/send-feedback/:where", app.UpdateMailTemplate)

	// Serve static files and HTML templates.
	app.Router.Static("/static", "./static")
	app.Router.LoadHTMLGlob("templates/*")
}

// InitApp walks through the necessary steps to set up a correct
// environment for MODULIST. This includes parsing command line
// flags, connecting to the database, registering HTTP endpoints
// and saving everything to a portable struct.
func InitApp() *App {

	// Load .env configuration file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[InitApp] Failed to load .env file. Terminating.")
	}

	// Create an empty App struct to fill from .env.
	app := new(App)

	// Read TLS parameters.
	app.TLS, err = strconv.ParseBool(os.Getenv("HTTP_TLS"))
	app.TLSCertFile = os.Getenv("HTTP_TLS_CERT_FILE")
	app.TLSKeyFile = os.Getenv("HTTP_TLS_KEY_FILE")
	if err != nil {
		log.Fatal("[InitApp] Unrecognized TLS indicator in .env, expecting bool. Terminating.")
	}

	// Save on which IP and port the frontend is to be run.
	app.IP = os.Getenv("HTTP_IP")
	app.Port = os.Getenv("HTTP_PORT")

	// Store stage mode the application is running in.
	app.Stage = os.Getenv("DEPLOY_STAGE")

	// Save bcrypt hash cost from .env.
	// In production, this value should be at least '16'.
	app.HashCost, err = strconv.Atoi(os.Getenv("APP_PASSWORD_HASH_COST"))
	if err != nil {
		log.Fatal("[InitApp] Could not load APP_PASSWORD_HASH_COST from .env file. Missing or not an integer?")
	}

	// Set JWT session token validity to the duration in minutes loaded from environment.
	validFor, err := strconv.Atoi(os.Getenv("APP_JWT_VALID_FOR"))
	if err != nil {
		log.Fatal("[InitApp] Could not load APP_JWT_VALID_FOR from .env file. Missing or not an integer?")
	}
	app.JWTValidFor = time.Duration(validFor) * time.Minute

	// Before starting gin, check if we are running in
	// production and do not want to log everything.
	if app.Stage == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Fill router variable with default gin router.
	app.Router = gin.Default()

	// Append database connection.
	app.DB = db.InitDB()

	// Initialize the validator instance to validate fields with tag 'validate'.
	app.Validator = validator.New()

	// Register frontend routes.
	app.DefineRoutes()

	// Check if an initialization command line flag was provided.
	initFlag := flag.Bool("init", false, "Append this flag in order to initialize a new setup of MODULIST. This includes the interactive creation of the default admin user.")
	flag.Parse()

	if *initFlag {

		// Create all needed database tables.
		db.SetUpTables(app.DB)

		// Transfer modules from SQLite database specified
		// in .env file to main database.
		db.TransferModules(app.DB, os.Getenv("MODULES_SQLITE_PATH"))

		// Default admin user creation.
		fmt.Printf("\n========== Begin initializing MODULIST ==========\n\nCreate default admin user.\n")

		var NewAdmin db.User

		// Assign new UUID v4.
		NewAdmin.ID = fmt.Sprintf("%s", uuid.NewV4())

		// Ask for first name.
		fmt.Printf("First name: ")
		fmt.Scanln(&NewAdmin.FirstName)

		// Ask for last name.
		fmt.Printf("Last name: ")
		fmt.Scanln(&NewAdmin.LastName)

		// Ask for mail address.
		fmt.Printf("Mail address: ")
		fmt.Scanln(&NewAdmin.Mail)

		// Set mail verification flag of this user to false, initially.
		NewAdmin.MailVerified = false

		// Ask for status group of new admin.
		fmt.Printf("Status group (0 for Prof, 1 for WiMi, 2 for Studi, 3 for other): ")
		fmt.Scanln(&NewAdmin.StatusGroup)

		// Ask for password.
		fmt.Printf("Password: ")
		pwd, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal("\n[InitApp] Error during read of admin's password from STDIN during initilization.")
		}

		fmt.Print("\nGenerating secure hash of provided password. Please wait... ")

		// Generate a secure bcrypt hash from supplied password.
		hash, err := bcrypt.GenerateFromPassword(pwd, app.HashCost)
		if err != nil {
			// If there was an error during hash creation - terminate immediately.
			log.Fatal("[InitApp] Error while generating hash in user creation. Terminating.")
		}
		NewAdmin.PasswordHash = string(hash)

		fmt.Print("Done!\n\n")

		// Give this user admin permissions.
		NewAdmin.Privileges = db.PRIVILEGE_ADMIN

		// Set account to enabled, initially.
		NewAdmin.Enabled = true

		fmt.Printf("\n\nNEW USER: %v\n\n", NewAdmin)

		// Save new admin to database.
		app.DB.Create(&NewAdmin)
		fmt.Printf("Created new admin user with mail address '%s'.\n\n", NewAdmin.Mail)
		fmt.Printf("==========  End initializing MODULIST  ==========\n\n\n")
	}

	return app
}
