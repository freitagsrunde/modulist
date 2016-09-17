package main

import (
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
	"github.com/leebenson/conform"
)

// Structs

type SearchPayload struct {
	Query string `conform:"trim,lower"`
}

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

	// Load all modules beginning with 'A' from database.
	var Modules []db.Module
	app.DB.Where("lower(\"title\") LIKE ?", "a%").Find(&Modules)

	c.HTML(http.StatusOK, "modules-list.html", gin.H{
		"PageTitle":   "Übersicht der Modulbeschreibungen",
		"User":        User,
		"FirstLetter": "A",
		"Modules":     Modules,
	})
}

func (app *App) SearchModules(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Retrieve query term from URL.
	Payload := SearchPayload{
		Query: c.Request.URL.Query().Get("query"),
	}

	// Let it be conformant.
	conform.Strings(&Payload)

	// Find all modules in main database that contain
	// the query in their module titles.
	var Modules []db.Module
	app.DB.Where("lower(\"title\") LIKE ? OR lower(\"title_english\") LIKE ?", ("%" + Payload.Query + "%"), ("%" + Payload.Query + "%")).Find(&Modules)

	c.HTML(http.StatusOK, "modules-list.html", gin.H{
		"PageTitle":   "Übersicht der Modulbeschreibungen",
		"User":        User,
		"FirstLetter": "all",
		"Modules":     Modules,
	})
}

func (app *App) FilterModulesByLetter(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	var Modules []db.Module
	firstLetter := c.Param("firstLetter")

	// Retrieve filter letter from URL.
	Payload := SearchPayload{
		Query: firstLetter,
	}

	// Let it be conformant.
	conform.Strings(&Payload)

	if Payload.Query == "all" {
		app.DB.Find(&Modules)
	} else {
		app.DB.Where("lower(\"title\") LIKE ?", (Payload.Query + "%")).Find(&Modules)
	}

	c.HTML(http.StatusOK, "modules-list.html", gin.H{
		"PageTitle":   "Übersicht der Modulbeschreibungen",
		"User":        User,
		"FirstLetter": firstLetter,
		"Modules":     Modules,
	})
}

func (app *App) MarkModuleDone(c *gin.Context) {}
