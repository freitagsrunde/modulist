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

	// Load all modules from database.
	var Modules []db.Module
	app.DB.Where("\"title\" LIKE ?", "A%").Find(&Modules)

	c.HTML(http.StatusOK, "modules-list.html", gin.H{
		"PageTitle":   "Übersicht der Modulbeschreibungen",
		"User":        User,
		"FirstLetter": "A",
		"Modules":     Modules,
	})
}

func (app *App) FilterModules(c *gin.Context) {}

func (app *App) FilterModulesByLetter(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Load all modules with first letter taken from
	// GET parameters or all from database.
	var Modules []db.Module
	firstLetter := c.Param("firstLetter")

	if firstLetter == "all" {
		app.DB.Order("lower(\"title\") asc").Find(&Modules)
	} else {
		app.DB.Where("\"title\" LIKE ?", (firstLetter + "%")).Find(&Modules)
	}

	c.HTML(http.StatusOK, "modules-list.html", gin.H{
		"PageTitle":   "Übersicht der Modulbeschreibungen",
		"User":        User,
		"FirstLetter": firstLetter,
		"Modules":     Modules,
	})
}

func (app *App) MarkModuleDone(c *gin.Context) {}
