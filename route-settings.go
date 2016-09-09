package main

import (
	"time"

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

type PasswordLinkViewPayload struct {
	SecretToken string `conform:"trim" validate:"required,len=72,alphanum"`
}

type UsePasswordLinkPayload struct {
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

// PasswordLinkView prompts the user requesting this
// page to set her or his initial password after a new
// account was created and the setup link was sent out.
func (app *App) PasswordLinkView(c *gin.Context) {

	// Extract supposed secret token from URL.
	Payload := PasswordLinkViewPayload{
		SecretToken: c.Param("secretToken"),
	}

	// Check secret token for conformity and validity.
	if errs := app.ConformAndValidate(&Payload); errs != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Attempt to find secret token in database for password links.
	var PasswordLink db.PasswordLink
	app.DB.First(&PasswordLink, "\"secret_token\" = ?", Payload.SecretToken)

	// If no element with matching token could be found, redirect.
	if PasswordLink.ID == "" {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Check if token is not yet expired.
	if time.Now().After(PasswordLink.Expires) {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Token checks out, display password page.
	c.HTML(http.StatusOK, "password-link.html", gin.H{
		"PageTitle":   "Passwort setzen",
		"MainTitle":   "Passwort setzen",
		"SecretToken": Payload.SecretToken,
	})
}

// UsePasswordLink takes the user's choice for her or his
// password, validates it and sets it as the new password
// for the new user in the database.
func (app *App) UsePasswordLink(c *gin.Context) {

	// Extract supposed secret token from URL.
	TokenPayload := PasswordLinkViewPayload{
		SecretToken: c.Param("secretToken"),
	}

	// Check secret token for conformity and validity.
	if errs := app.ConformAndValidate(&TokenPayload); errs != nil {

		// If payload did not pass, redirect user to start page.
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Attempt to find secret token in database for password links.
	var PasswordLink db.PasswordLink
	app.DB.First(&PasswordLink, "\"secret_token\" = ?", TokenPayload.SecretToken)
	app.DB.Model(&PasswordLink).Related(&PasswordLink.User)

	// Check if token is not yet expired.
	if time.Now().After(PasswordLink.Expires) {

		// Token expired, redirect to start page.
		c.Redirect(http.StatusFound, "/")

		return
	}

	var PasswordPayload UsePasswordLinkPayload

	err := c.BindWith(&PasswordPayload, binding.FormPost)
	if err != nil {

		c.HTML(http.StatusBadRequest, "password-link.html", gin.H{
			"PageTitle":   "Passwort setzen",
			"MainTitle":   "Passwort setzen",
			"SecretToken": TokenPayload.SecretToken,
			"FatalError":  "Gesendete Daten zum Aktualisieren des Passworts konnten nicht verarbeitet werden. Bitte erneut versuchen.",
		})

		return
	}

	// Check sent content for validity.
	ErrorDesc := app.ConformAndValidate(&PasswordPayload)
	if ErrorDesc != nil {

		// If payload did not pass, report errors to user.
		c.HTML(http.StatusBadRequest, "password-link.html", gin.H{
			"PageTitle":   "Passwort setzen",
			"MainTitle":   "Passwort setzen",
			"SecretToken": TokenPayload.SecretToken,
			"Errors":      ErrorDesc,
		})

		return
	}

	// Compare both supplied new password candidates for equality.
	if PasswordPayload.NewPassword != PasswordPayload.RepeatedNewPassword {

		// Signal client that an error occured.
		c.HTML(http.StatusBadRequest, "password-link.html", gin.H{
			"PageTitle":   "Passwort setzen",
			"MainTitle":   "Passwort setzen",
			"SecretToken": TokenPayload.SecretToken,
			"FatalError":  "Die beiden Zeichenketten des neuen Passworts stimmen nicht überein. Bitte dasselbe Passwort zweimal eingeben.",
		})

		return
	}

	// Everything until now checks out. Generate new password hash.
	hash, err := bcrypt.GenerateFromPassword([]byte(PasswordPayload.NewPassword), app.HashCost)
	if err != nil {

		// Signal client that an error occured.
		c.HTML(http.StatusBadRequest, "password-link.html", gin.H{
			"PageTitle":   "Passwort setzen",
			"MainTitle":   "Passwort setzen",
			"SecretToken": TokenPayload.SecretToken,
			"FatalError":  "Es ist etwas schiefgegangen. Bitte erneut versuchen.",
		})

		return
	}

	// Delete entry in table containing tokens for password reset.
	app.DB.Delete(&PasswordLink)

	// Update user element in database to new password hash,
	// verified mail address and enable the account.
	app.DB.Model(&db.User{ID: PasswordLink.UserID}).Updates(&db.User{
		MailVerified: true,
		PasswordHash: string(hash),
		Enabled:      true,
	})

	// Everything went fine. Signal success to user.
	c.HTML(http.StatusOK, "password-link.html", gin.H{
		"PageTitle":   "Passwort setzen",
		"MainTitle":   "Passwort setzen",
		"SecretToken": TokenPayload.SecretToken,
		"Success":     "Passwort wurde gesetzt! Es kann losgehen!",
	})
}
