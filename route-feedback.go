package main

import (
	"fmt"
	"strings"

	"html/template"
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

	// Only use this field if it contains a valid value.
	if Module.ReferencePersonID.Valid {
		app.DB.Model(&Module).Related(&Module.ReferencePerson, "ReferencePersonID")
	}

	// Only use this field if it contains a valid value.
	if Module.ResponsiblePersonID.Valid {
		app.DB.Model(&Module).Related(&Module.ResponsiblePerson, "ResponsiblePersonID")
	}

	Module.LearningOutcomesHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.LearningOutcomes.String), "\n", "<br />", -1))
	Module.LearningOutcomesEnglishHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.LearningOutcomesEnglish.String), "\n", "<br />", -1))
	Module.TeachingContentsHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.TeachingContents.String), "\n", "<br />", -1))
	Module.TeachingContentsEnglishHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.TeachingContentsEnglish.String), "\n", "<br />", -1))
	Module.InstructiveFormHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.InstructiveForm), "\n", "<br />", -1))
	Module.OptionalRequirementsHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.OptionalRequirements), "\n", "<br />", -1))
	Module.MandatoryRequirementsHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.MandatoryRequirements.String), "\n", "<br />", -1))
	Module.ExaminationDescriptionHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.ExaminationDescription.String), "\n", "<br />", -1))
	Module.MiscellaneousHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.Miscellaneous.String), "\n", "<br />", -1))
	Module.LiteratureHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.Literature), "\n", "<br />", -1))
	Module.RegistrationFormalitiesHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.RegistrationFormalities.String), "\n", "<br />", -1))

	c.HTML(http.StatusOK, "module-feedback.html", gin.H{
		"PageTitle": fmt.Sprintf("Feedback zu Modul #%d", Module.ModuleID),
		"User":      User,
		"Module":    Module,
	})
}

func (app *App) AddFeedback(c *gin.Context) {}

func (app *App) DeleteFeedback(c *gin.Context) {}
