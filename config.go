// Package modulist is a web-based review tool for module
// descriptions developed by Freitagsrunde at TU Berlin.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/numbleroot/modulist/db"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	TLS         bool
	TLSCertFile string
	TLSKeyFile  string
	IP          string
	Port        string
	Stage       string
	Router      *gin.Engine
	DB          *gorm.DB
}

// DefineRoutes registers all frontend endpoints.
// These will make up the web interface of MODULIST.
func (app *App) DefineRoutes() {

	// Route 'index'.
	app.Router.GET("/", app.Index)
	app.Router.POST("/login", app.Login)
	app.Router.GET("/logout", app.Logout)

	// Route 'list'.
	app.Router.GET("/modules", app.ListModules)
	app.Router.POST("/modules", app.FilterModules)
	app.Router.GET("/modules/:firstLetter", app.FilterModulesByLetter)
	app.Router.GET("/done/:id", app.MarkModuleDone)

	// Route 'search'.
	app.Router.GET("/search", app.Search)

	// Route 'feedback'.
	app.Router.GET("/review/:moduleID", app.ReviewModule)
	app.Router.POST("/review/:moduleID/addFeedback/:category", app.AddFeedback)
	app.Router.GET("/deleteFeedback/:id", app.DeleteFeedback)

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

	// Before starting gin, check if we are running in
	// production and do not want to log everything.
	if app.Stage == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Fill router variable with default gin router.
	app.Router = gin.Default()

	// Append database connection.
	app.DB = db.InitDB()

	// Register frontend routes.
	app.DefineRoutes()

	// Check if a command line flag was provided.
	initFlag := flag.Bool("init", false, "Append this flag in order to initialize a new setup of MODULIST. This includes the interactive creation of the default admin user.")
	flag.Parse()

	if *initFlag {

		// Create all needed database tables.
		db.SetUpTables(app.DB)

		// Transfer modules from SQLite database specified
		// in .env file to main database.
		// TODO: do that.

		// Ask user for data of default admin user to create.

		var NewAdmin db.User
		NewAdmin.ID = fmt.Sprintf("%s", uuid.NewV4())

		fmt.Printf("\nInitializing MODULIST.\n\nCreate default admin user.\nFirst name: ")
		fmt.Scanln(&NewAdmin.FirstName)

		fmt.Printf("Last name: ")
		fmt.Scanln(&NewAdmin.LastName)

		fmt.Printf("Mail address: ")
		fmt.Scanln(&NewAdmin.Mail)

		NewAdmin.MailVerified = false

		fmt.Printf("Password: ")
		pwd, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal("\n[InitApp] Error during read of admin's password from STDIN during initilization.")
		}

		fmt.Print("\nGenerating secure hash of provided password. Please wait... ")

		// TODO: use .env for hash cost.
		hash, err := bcrypt.GenerateFromPassword(pwd, 16)
		if err != nil {
			// If there was an error during hash creation - terminate immediately.
			log.Fatal("[InitApp] Error while generating hash in user creation. Terminating.")
		}
		NewAdmin.PasswordHash = string(hash)

		fmt.Print("Done!\n\n")

		// TODO: make useful.
		NewAdmin.Privileges = 1
		NewAdmin.Enabled = true

		fmt.Printf("New admin:\n%v\n", NewAdmin)

		// Save new admin to database.
		app.DB.Create(&NewAdmin)
		fmt.Printf("\nCreated new admin user with mail address '%s'.\n\n", NewAdmin.Mail)
	}

	return app
}
