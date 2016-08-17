package main

import (
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
)

// Functions

func (app *App) ListModules(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	c.HTML(http.StatusOK, "modules-list.html", gin.H{
		"PageTitle":   "Übersicht der Modulbeschreibungen",
		"User":        User,
		"FirstLetter": "A",
		"Alphabet":    map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true, "F": true, "G": true, "H": true, "I": true, "J": true, "K": true, "L": true, "M": true, "N": true, "O": true, "P": true, "Q": true, "R": true, "S": true, "T": true, "U": true, "V": true, "W": true, "X": true, "Y": true, "Z": true},
		"ModuleList":  struct{}{},
	})
}

func (app *App) FilterModules(c *gin.Context) {}

func (app *App) FilterModulesByLetter(c *gin.Context) {}

func (app *App) MarkModuleDone(c *gin.Context) {}
