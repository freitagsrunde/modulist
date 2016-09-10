package main

import (
	"fmt"

	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
)

// Structs

type ReviewModulePayload struct {
	ID string `conform:"trim" validate:"required,uuid4"`
}

// Functions

func (app *App) ReviewModule(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Extract ID of module to review from URL.
	Payload := ReviewModulePayload{
		ID: c.Param("moduleID"),
	}

	// Check supplied ID for conformity and validity.
	if errs := app.ConformAndValidate(&Payload); errs != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Retrieve module information for supplied ID from database.
	var Module db.Module
	app.DB.First(&Module, "\"id\" = ?", Payload.ID)
	app.DB.Model(&Module).Related(&Module.ReferencePerson, "ReferencePersonID")
	app.DB.Model(&Module).Related(&Module.ResponsiblePerson, "ResponsiblePersonID")

	c.HTML(http.StatusOK, "module-feedback.html", gin.H{
		"PageTitle": fmt.Sprintf("Feedback zu Modul #%d", Module.ModuleID),
		"User":      User,
		"Module":    Module,
	})
}

func (app *App) AddFeedback(c *gin.Context) {}

func (app *App) DeleteFeedback(c *gin.Context) {}
