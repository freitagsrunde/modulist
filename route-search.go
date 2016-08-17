package main

import (
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
)

// Functions

func (app *App) Search(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// TODO: Correct search page behaviour.

	c.HTML(http.StatusOK, "search.html", gin.H{
		"PageTitle": "Suche",
		"User":      User,
		"Query":     "BERND",
		"Modules":   struct{}{},
	})
}
