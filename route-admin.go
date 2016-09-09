package main

import (
	"fmt"
	"log"
	"time"

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

type ActDeactUserPayload struct {
	ID string `conform:"trim" validate:"required,uuid4"`
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

		log.Printf("[CreateUser] Generating random bytes for temporary user password went wrong: %s.\n", err.Error())

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

	pwd := fmt.Sprintf("%x", randomBytes)

	// Generate a secure bcrypt hash from generated random bytes.
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), app.HashCost)
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

	// Generate secret link to initial password site.
	var PasswordLink db.PasswordLink
	PasswordLink.ID = fmt.Sprintf("%s", uuid.NewV4())
	PasswordLink.UserID = NewUser.ID
	PasswordLink.User = NewUser

	// Link to password site is valid for 5 days.
	PasswordLink.Expires = time.Now().Add(5 * 24 * time.Hour)

	// Again, generate some amount of random bytes.
	randomBytes = make([]byte, 36)
	_, err = rand.Read(randomBytes)
	if err != nil {

		log.Printf("[CreateUser] Generating random bytes for password link went wrong: %s.\n", err.Error())

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

	// Store string of random bytes as secret token.
	PasswordLink.SecretToken = fmt.Sprintf("%x", randomBytes)

	// Save password link element to database.
	app.DB.Create(&PasswordLink)

	// TODO: Send out link to new user with password
	//       link and expiration date of that link.

	// Save new user to database.
	app.DB.Create(&NewUser)

	// Retrieve an updated list of all users to display.
	app.DB.Find(&Users)

	c.HTML(http.StatusOK, "admin-users.html", gin.H{
		"PageTitle": "Admin - Nutzerverwaltung",
		"User":      User,
		"Users":     Users,
		"Success":   "Nutzer angelegt! Eine Mail mit einem Link zum Setzen des Passworts wurde versandt.",
	})
}

// DeactivateUser takes in an user's ID and after
// some conformity and validity checks disables that
// user's account. This will prevent the user from
// logging in to the service.
func (app *App) DeactivateUser(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_ADMIN)
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Retrieve ID of user to deactivate from URL.
	Payload := ActDeactUserPayload{
		ID: c.Param("id"),
	}

	// Check if sent ID is conform and valid.
	if errs := app.ConformAndValidate(&Payload); errs != nil {
		c.Redirect(http.StatusFound, "/admin/users")

		return
	}

	// Set user account specified by supplied ID
	// to disabled in database.
	app.DB.Model(&db.User{ID: Payload.ID}).Update("enabled", false)

	// Redirect if everything was successful.
	c.Redirect(http.StatusFound, "/admin/users")
}

// ActivateUser is the counterpart to DeactivateUser.
// It takes in an user's ID, too. After checks, it
// generates a password link again, much like the behaviour
// when user accounts are first created. At the end
// the user will be notified with a mail containing
// the password link.
func (app *App) ActivateUser(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_ADMIN)
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Retrieve ID of user to activate from URL.
	Payload := ActDeactUserPayload{
		ID: c.Param("id"),
	}

	// Check if sent ID is conform and valid.
	if errs := app.ConformAndValidate(&Payload); errs != nil {
		c.Redirect(http.StatusFound, "/admin/users")

		return
	}

	// Load concerned user from database.
	var ActivatedUser db.User
	app.DB.First(&ActivatedUser, "\"id\" = ?", Payload.ID)

	// Generate a random, secure password hash for reactivated user.
	// Will never be used but we have to satisfy the database
	// constraints and also have some worst case help.
	randomBytes := make([]byte, 24)
	_, err = rand.Read(randomBytes)
	if err != nil {

		log.Printf("[ActivateUser] Generating random bytes for temporary user password went wrong: %s.\n", err.Error())
		c.Redirect(http.StatusFound, "/admin/users")

		return
	}

	pwd := fmt.Sprintf("%x", randomBytes)

	// Generate a secure bcrypt hash from generated random bytes.
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), app.HashCost)
	if err != nil {

		log.Printf("[ActivateUser] Creating bcrypt password hash went wrong: %s.\n", err.Error())
		c.Redirect(http.StatusFound, "/admin/users")

		return
	}

	// Save random password in to-be-reactivated user's struct.
	ActivatedUser.PasswordHash = string(hash)

	// Generate secret link to site for setting the password again.
	var PasswordLink db.PasswordLink
	PasswordLink.ID = fmt.Sprintf("%s", uuid.NewV4())
	PasswordLink.UserID = ActivatedUser.ID
	PasswordLink.User = ActivatedUser

	// Link to password site is valid for 5 days.
	PasswordLink.Expires = time.Now().Add(5 * 24 * time.Hour)

	// Again, generate some amount of random bytes.
	randomBytes = make([]byte, 36)
	_, err = rand.Read(randomBytes)
	if err != nil {

		log.Printf("[ActivateUser] Generating random bytes for password link went wrong: %s.\n", err.Error())
		c.Redirect(http.StatusFound, "/admin/users")

		return
	}

	// Store string of random bytes as secret token.
	PasswordLink.SecretToken = fmt.Sprintf("%x", randomBytes)

	// Save password link element to database.
	app.DB.Create(&PasswordLink)

	// TODO: Send out mail to user containing link to password
	//       reset site and expiration notification.

	// Redirect if everything was successful.
	c.Redirect(http.StatusFound, "/admin/users")
}

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
		"MailHeader":    "A",
		"MailFooter":    "B",
	})
}

func (app *App) UpdateMailTemplate(c *gin.Context) {}
