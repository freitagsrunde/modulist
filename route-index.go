package main

import (
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

// Structs

type LoginPayload struct {
	Mail     string `form:"login-mail" conform:"trim,email" validate:"required,email"`
	Password string `form:"login-password" validate:"required"`
}

// Functions

// Index renders the page first visible when
// navigating to MODULIST start page.
func (app *App) Index(c *gin.Context) {

	// Check if user is already logged in.
	_, err := app.Authorize(c.Request)
	if err == nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"PageTitle": "Willkommen bei MODULIST",
		"MainTitle": "Willkommen bei MODULIST",
	})

	return
}

// Login provides all necessary functionality
// in order to log in a user.
func (app *App) Login(c *gin.Context) {

	// Check if user is already logged in.
	_, err := app.Authorize(c.Request)
	if err == nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	var Payload LoginPayload

	err = c.BindWith(&Payload, binding.FormPost)
	if err != nil {

		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"PageTitle":  "Willkommen bei MODULIST",
			"MainTitle":  "Willkommen bei MODULIST",
			"FatalError": "Sent data could not be interpreted. Please try again.",
		})

		return
	}

	// Check sent content for validity.
	ErrorDesc := app.ConformAndValidate(Payload)
	if ErrorDesc != nil {

		// If payload did not pass, report errors to user.
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"PageTitle": "Willkommen bei MODULIST",
			"MainTitle": "Willkommen bei MODULIST",
			"Errors":    ErrorDesc,
		})

		return
	}

	// Data is valid, try to locate user in database.
	var User db.User
	app.DB.First(&User, "\"mail\" = ?", Payload.Mail)

	// Compare password hash from database with possible plaintext
	// password from submitted login form. Compares in constant time.
	err = bcrypt.CompareHashAndPassword([]byte(User.PasswordHash), []byte(Payload.Password))
	if (User.ID == "") || (err != nil) {

		// Signal client that an error occured.
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"PageTitle":  "Willkommen bei MODULIST",
			"MainTitle":  "Willkommen bei MODULIST",
			"FatalError": "Mail and/or password wrong.",
		})

		return
	}

	// Create a JWT and store it as a cookie.
	app.CreateSession(c, User)

	// Redirect to first authorized page.
	c.Redirect(http.StatusFound, "/modules")
}

// Logout destroys the user's session by storing
// garbage in the current session cookie and instructing
// the browser to delete that cookie.
func (app *App) Logout(c *gin.Context) {

	// Check if user is authorized.
	_, err := app.Authorize(c.Request)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Set token cookie content to garbage and
	// expiration date to a date in the past.
	c.SetCookie("Token", "", -1, "", "", false, true)

	// Redirect back to index page.
	c.Redirect(http.StatusFound, "/")
}
