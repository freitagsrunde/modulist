package main

import (
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
)

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

	// TODO: Correct settings page behaviour.

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"PageTitle": "Einstellungen",
		"User":      User,
		"Success":   false,
		"Error":     "",
	})
}

func (app *App) UpdateSettings(c *gin.Context) {}
