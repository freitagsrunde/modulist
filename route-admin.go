package main

import (
	"fmt"
	"log"

	"crypto/rand"
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Structs

type CreateUserPayload struct {
	FirstName   string `form:"user-first-name" conform:"trim" validate:"required,excludesall=!@#$%^&*()_+-=:;?/0x2C0x7C"`
	LastName    string `form:"user-last-name" conform:"trim" validate:"required,excludesall=!@#$%^&*()_+-=:;?/0x2C0x7C"`
	Mail        string `form:"user-mail" conform:"trim,email" validate:"required,email"`
	StatusGroup int    `form:"user-status-group" conform:"trim" validate:"min=0"`
	Privileges  int    `form:"user-privileges" conform:"trim" validate:"min=0"`
}

// Functions

func (app *App) ListUsers(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_ADMIN)
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Fetch all users registered in database.
	var Users []db.User
	app.DB.Find(&Users)

	c.HTML(http.StatusOK, "admin-users.html", gin.H{
		"PageTitle": "Admin - Nutzerverwaltung",
		"User":      User,
		"Users":     Users,
	})
}

func (app *App) CreateUser(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_ADMIN)
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	var Payload CreateUserPayload
	var Users []db.User

	err = c.BindWith(&Payload, binding.FormPost)
	if err != nil {

		app.DB.Find(&Users)

		c.HTML(http.StatusBadRequest, "admin-users.html", gin.H{
			"PageTitle":  "Admin - Nutzerverwaltung",
			"User":       User,
			"Users":      Users,
			"FatalError": "Gesendete Daten für neuen Nutzer konnten nicht verarbeitet werden. Bitte erneut versuchen.",
		})

		return
	}

	// Check sent content for validity.
	ErrorDesc := app.ConformAndValidate(&Payload)
	if ErrorDesc != nil {

		app.DB.Find(&Users)

		// If payload did not pass, report errors to user.
		c.HTML(http.StatusBadRequest, "admin-users.html", gin.H{
			"PageTitle": "Admin - Nutzerverwaltung",
			"User":      User,
			"Users":     Users,
			"Errors":    ErrorDesc,
		})

		return
	}

	// Based on payload from request, create new user struct.
	var NewUser db.User

	NewUser.ID = fmt.Sprintf("%s", uuid.NewV4())
	NewUser.FirstName = Payload.FirstName
	NewUser.LastName = Payload.LastName
	NewUser.Mail = Payload.Mail
	NewUser.MailVerified = false
	NewUser.StatusGroup = Payload.StatusGroup
	NewUser.Privileges = Payload.Privileges
	NewUser.Enabled = false

	// Generate a random, secure password hash for new user.
	// Will never be used but we have to satisfy the database
	// constraints and also have some worst case help.
	randomBytes := make([]byte, 24)

	_, err = rand.Read(randomBytes)
	if err != nil {

		log.Printf("[CreateUser] Generating random bytes went wrong: %s.\n", err.Error())

		app.DB.Find(&Users)

		// Report fatal error to user.
		c.HTML(http.StatusInternalServerError, "admin-users.html", gin.H{
			"PageTitle":  "Admin - Nutzerverwaltung",
			"User":       User,
			"Users":      Users,
			"FatalError": "Auf dem Server ist ein Fehler aufgetreten. Erneut versuchen oder Admin kontaktieren.",
		})

		return
	}

	// Generate a secure bcrypt hash from generated random bytes.
	hash, err := bcrypt.GenerateFromPassword(randomBytes, app.HashCost)
	if err != nil {

		log.Printf("[CreateUser] Creating bcrypt password hash went wrong: %s.\n", err.Error())

		app.DB.Find(&Users)

		// Report fatal error to user.
		c.HTML(http.StatusInternalServerError, "admin-users.html", gin.H{
			"PageTitle":  "Admin - Nutzerverwaltung",
			"User":       User,
			"Users":      Users,
			"FatalError": "Auf dem Server ist ein Fehler aufgetreten. Erneut versuchen oder Admin kontaktieren.",
		})

		return
	}

	// Finally save random password in new user's struct.
	NewUser.PasswordHash = string(hash)

	// TODO: Generate secret link to initial password site.

	// TODO: Send out link to new user.

	// Save new user to database.
	app.DB.Create(&NewUser)

	// Retrieve an updated list of all users to display.
	app.DB.Find(&Users)

	c.HTML(http.StatusOK, "admin-users.html", gin.H{
		"PageTitle": "Admin - Nutzerverwaltung",
		"User":      User,
		"Users":     Users,
	})
}

func (app *App) DeleteUser(c *gin.Context) {}

func (app *App) SendFeedback(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_ADMIN)
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// TODO: Correct admin send feedback page behaviour.

	c.HTML(http.StatusOK, "admin-send-feedback.html", gin.H{
		"PageTitle":     "Admin - Feedback versenden",
		"User":          User,
		"FeedbackMails": struct{}{},
		"Error":         "",
		"MailHeader":    "A",
		"MailFooter":    "B",
	})
}

func (app *App) UpdateMailTemplate(c *gin.Context) {}
