package main

import (
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

// Structs

type UpdatePasswordPayload struct {
	OldPassword         string `form:"old-password" validate:"required"`
	NewPassword         string `form:"new-password" validate:"required,min=16,containsany=0123456789,containsany=!@#$%^&*()_+-=:;?/0x2C0x7C"`
	RepeatedNewPassword string `form:"repeat-new-password" validate:"required,min=16,containsany=0123456789,containsany=!@#$%^&*()_+-=:;?/0x2C0x7C"`
}

// Functions

func (app *App) ListSettings(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"PageTitle": "Einstellungen",
		"User":      User,
	})
}

func (app *App) UpdateSettings(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	var Payload UpdatePasswordPayload

	err = c.BindWith(&Payload, binding.FormPost)
	if err != nil {

		c.HTML(http.StatusBadRequest, "settings.html", gin.H{
			"PageTitle":  "Einstellungen",
			"User":       User,
			"FatalError": "Gesendete Daten zum Aktualisieren des Passworts konnten nicht verarbeitet werden. Bitte erneut versuchen.",
		})

		return
	}

	// Check sent content for validity.
	ErrorDesc := app.ConformAndValidate(&Payload)
	if ErrorDesc != nil {

		// If payload did not pass, report errors to user.
		c.HTML(http.StatusBadRequest, "settings.html", gin.H{
			"PageTitle": "Einstellungen",
			"User":      User,
			"Errors":    ErrorDesc,
		})

		return
	}

	// Compare password hash from database with possible plaintext
	// password from submitted update form. Compares in constant time.
	err = bcrypt.CompareHashAndPassword([]byte(User.PasswordHash), []byte(Payload.OldPassword))
	if (User.ID == "") || (err != nil) {

		// Signal client that an error occured.
		c.HTML(http.StatusBadRequest, "settings.html", gin.H{
			"PageTitle":  "Einstellungen",
			"User":       User,
			"FatalError": "Das bisherige Passwort ist falsch.",
		})

		return
	}

	// Compare both supplied new password candidates for equality.
	if Payload.NewPassword != Payload.RepeatedNewPassword {

		// Signal client that an error occured.
		c.HTML(http.StatusBadRequest, "settings.html", gin.H{
			"PageTitle":  "Einstellungen",
			"User":       User,
			"FatalError": "Die beiden Zeichenketten des neuen Passworts stimmen nicht überein. Bitte dasselbe Passwort zweimal eingeben.",
		})

		return
	}

	// Everything until now checks out. Generate new password hash.
	hash, err := bcrypt.GenerateFromPassword([]byte(Payload.NewPassword), app.HashCost)
	if err != nil {

		// Signal client that an error occured.
		c.HTML(http.StatusBadRequest, "settings.html", gin.H{
			"PageTitle":  "Einstellungen",
			"User":       User,
			"FatalError": "Es ist etwas schiefgegangen. Bitte erneut versuchen.",
		})

		return
	}

	// Update user element in database to new password hash.
	app.DB.Model(&User).Select("password_hash").Update("PasswordHash", string(hash))

	// Create a JWT and store it as a cookie.
	app.CreateSession(c, *User)

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"PageTitle": "Einstellungen",
		"User":      User,
		"Success":   "Neues Passwort gespeichert!",
	})
}
