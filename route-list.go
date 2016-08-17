package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Functions

func (app *App) ListModules(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)
}

func (app *App) FilterModules(c *gin.Context) {}

func (app *App) FilterModulesByLetter(c *gin.Context) {}

func (app *App) MarkModuleDone(c *gin.Context) {}
